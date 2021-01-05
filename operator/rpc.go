// SPDX-License-Identifier: Apache-2.0

package operator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gorilla "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	pkgsync "perun.network/go-perun/pkg/sync"

	"github.com/perun-network/erdstall/wire"
)

type (
	// RPCServer handels RPC requests and forwards them to the enclave.
	RPCServer struct {
		pkgsync.Closer
		op     WireAPI
		server *opServer
	}

	// opServer is an implementation detail and wraps a `http.Server` and holds
	// necessary configuration values which hide the explicit instantiation of
	// an encrypted or unencrypted connection. Also holds client configuration
	// which is transmitted over the wire.
	opServer struct {
		server       *http.Server
		serveMux     *http.ServeMux
		cert         string // Maybe value for ssl certificate path.
		key          string // Maybe value for ssl key path.
		clientConfig ClientConfig
	}

	// OpServerConfig holds all necessary configuration values for the
	// `OpServer` which the RPCServer uses to communicate over the wire.
	OpServerConfig struct {
		Host         string
		Port         uint16
		CertFilePath string
		KeyFilePath  string
		ClientConfig ClientConfig
	}

	// ClientConfig describes the configuration, which is send to connecting
	// Clients.
	ClientConfig struct {
		NetworkID string // Network-ID
		Contract  common.Address
		POWDepth  uint64
	}

	// Peer is a connected client.
	Peer struct {
		pkgsync.Closer
		op WireAPI

		connMtx sync.Mutex // protects conn.
		conn    *gorilla.Conn

		sub *ProofSub
	}
)

// newOpServer returns an `opServer` instance.
func newOpServer(osc OpServerConfig) *opServer {
	serveMux := http.NewServeMux()
	return &opServer{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", osc.Host, osc.Port),
			Handler: serveMux,
		},
		serveMux:     serveMux,
		key:          osc.KeyFilePath,
		cert:         osc.CertFilePath,
		clientConfig: osc.ClientConfig,
	}
}

func (s *opServer) ListenAndServe() error {
	if s.cert != "" && s.key != "" {
		return s.server.ListenAndServeTLS(s.cert, s.key)
	} else {
		return s.server.ListenAndServe()
	}
}

func (s *opServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// NewRPC returns a new RPC object with the given `opServer`. Call Serve to start it.
func NewRPC(op WireAPI, osc OpServerConfig) *RPCServer {
	server := newOpServer(osc)
	rpc := &RPCServer{
		op:     op,
		server: server,
	}
	server.serveMux.HandleFunc("/ws", rpc.connectionHandler)
	return rpc
}

func (r *RPCServer) Log() *log.Entry {
	return log.WithField("role", "op")
}

// Serve serves RPC requests on the specified host and port.
// Should be called in a go-routine since it blocks.
func (r *RPCServer) Serve() error {
	if !r.OnClose(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		r.server.Shutdown(ctx) // nolint: errcheck
	}) {
		panic("Could not add OnClose function")
	}

	if err := r.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (r *RPCServer) connectionHandler(out http.ResponseWriter, in *http.Request) {
	upgrader := gorilla.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(out, in, nil)
	if err != nil {
		r.Log().WithError(err).Error("WS upgrade:", err)
		return
	}

	peer := &Peer{conn: conn, op: r.op}
	if err := peer.sendJSON(r.server.clientConfig); err != nil {
		r.Log().WithError(err).Error("Pushing ClientConfig:", err)
	}

	conn.SetCloseHandler(func(int, string) error {
		return peer.Close()
	})
	go peer.readMessages()
}

func (p *Peer) Log() *log.Entry {
	return log.WithField("role", "peer")
}

func (p *Peer) readMessages() {
	for !p.IsClosed() {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			p.Log().WithError(err).Error("Could not read msg")
			time.Sleep(1 * time.Second)
			continue
		}

		var call wire.Call
		var sendErr error
		if err := json.Unmarshal(msg, &call); err != nil {
			sendErr = p.sendResult("", fmt.Errorf("Invalid json: %w", err))
		} else {
			sendErr = p.sendResult(call.ID, p.handleCall(call.ID, call.Method, msg))
		}
		if sendErr != nil {
			p.Log().WithField("id", call.ID).WithError(sendErr).Error("Could not send result")
		}
	}
}

func (p *Peer) Close() error {
	if err := p.Closer.Close(); pkgsync.IsAlreadyClosedError(err) {
		return err
	}
	p.conn.Close()
	if p.sub != nil {
		p.sub.Unsubscribe()
	}
	return nil
}

func (p *Peer) handleCall(id wire.ID, method wire.Method, msg []byte) error {
	p.Log().Trace("Server received ", string(msg))
	switch method {
	case wire.MethodSendTx:
		var call wire.SendTx
		if err := json.Unmarshal(msg, &call); err != nil {
			return fmt.Errorf("unmarshalling SendTx: %w", err)
		}
		return p.op.Send(call.Tx)
	case wire.MethodSubscribe:
		var call wire.Subscribe
		if err := json.Unmarshal(msg, &call); err != nil {
			return fmt.Errorf("unmarshalling Subscribe: %w", err)
		}
		return p.subscribe(call.Who)
	default:
		return fmt.Errorf("unknown method '%s'", method)
	}
}

func (p *Peer) subscribe(who common.Address) error {
	if p.sub != nil {
		return errors.New("subscribed twice to proofs")
	}

	sub, err := p.op.SubscribeProofs(who)
	if err != nil {
		return fmt.Errorf("subscribing to proofs: %w", err)
	}
	p.sub = &sub

	go func() {
		for {
			var update interface{}

			select {
			case proof := <-p.sub.Deposits():
				update = &wire.DepositProof{
					Result: wire.Result{
						Topic: wire.DepositProofs,
					},
					Proof: proof,
				}
			case proof := <-p.sub.Balances():
				update = &wire.BalanceProof{
					Result: wire.Result{
						Topic: wire.BalanceProofs,
					},
					Proof: proof,
				}
			case <-p.sub.Closed():
				p.Log().Debug("Subscription routine returns due to closed sub.")
				return
			}

			if err := p.sendJSON(update); err != nil {
				p.Log().WithError(err).Error("Could not send topic update.")
			}
		}
	}()
	return nil
}

// sendResult sends an `error` that occurred while handling the message with
// `id` back to the user. The error can be nil-
func (p *Peer) sendResult(id wire.ID, _err error) error {
	errorMsg := ""
	if _err != nil {
		errorMsg = _err.Error()
	}
	return p.sendJSON(wire.Result{ID: id, Error: errorMsg})
}

func (p *Peer) sendJSON(obj interface{}) error {
	p.connMtx.Lock()
	defer p.connMtx.Unlock()
	return p.conn.WriteJSON(obj)
}
