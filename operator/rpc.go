// SPDX-License-Identifier: Apache-2.0

package operator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gorilla "github.com/gorilla/websocket"
	"github.com/perun-network/erdstall/wire"
	perunlog "perun.network/go-perun/log"
	pkgsync "perun.network/go-perun/pkg/sync"
)

type (
	// RPCServer handels RPC requests and forwards them to the enclave.
	RPCServer struct {
		pkgsync.Closer
		perunlog.Embedding
		op WireAPI

		mtx   sync.Mutex // protects peers
		peers []*Peer
	}

	// Peer is a connected client.
	Peer struct {
		pkgsync.Closer
		perunlog.Embedding
		op WireAPI

		connMtx sync.Mutex // protects conn.
		conn    *gorilla.Conn

		sub *ProofSub
	}
)

// NewRPC returns a new RPC object. Call Serve to start it.
func NewRPC(op WireAPI) *RPCServer {
	rpc := &RPCServer{op: op, Embedding: perunlog.MakeEmbedding(perunlog.WithField("role", "op"))}
	http.HandleFunc("/ws", rpc.connectionHandler)
	return rpc
}

// Serve serves RPC requests on the specified host and port.
// Should be called in a go-routine since it blocks.
func (r *RPCServer) Serve(host string, port uint16) error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), nil)
}

func (r *RPCServer) connectionHandler(out http.ResponseWriter, in *http.Request) {
	upgrader := gorilla.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(out, in, nil)
	if err != nil {
		r.Log().WithError(err).Error("WS upgrade:", err)
		return
	}

	peer := &Peer{conn: conn, op: r.op,
		Embedding: perunlog.MakeEmbedding(perunlog.WithField("role", "peer"))}
	conn.SetCloseHandler(func(int, string) error {
		return peer.Close()
	})
	r.addPeer(peer)
	peer.OnCloseAlways(func() {
		r.removePeer(peer)
	})
	go peer.readMessages()
}

func (r *RPCServer) addPeer(p *Peer) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.peers = append(r.peers, p)
}

func (r *RPCServer) removePeer(p *Peer) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.peers = append(r.peers, p)
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
		return fmt.Errorf("subscribed twice to proofs")
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
