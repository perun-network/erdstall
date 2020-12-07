// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gorilla "github.com/gorilla/websocket"
	"github.com/perun-network/erdstall/tee"
	"github.com/perun-network/erdstall/wire"
	perunlog "perun.network/go-perun/log"
	pkgsync "perun.network/go-perun/pkg/sync"
)

type (
	// RPC connects the client with the operator over websockets.
	RPC struct {
		pkgsync.Closer
		perunlog.Embedding

		connMtx sync.Mutex // protects conn.send.
		conn    *gorilla.Conn

		cbMtx     sync.RWMutex // protects callbacks.
		id        uint64
		callbacks map[wire.ID]callback

		subscription *Subscription
	}

	// Subscription is returned by Subscribe() and can be used to iterate
	// over the deposit and balance proofs.
	Subscription struct {
		perunlog.Embedding
		// deposit proofs from the OP will be written into this channel.
		depProofs chan tee.DepositProof
		// balance proofs from the OP will be written into this channel.
		balProofs chan tee.BalanceProof
	}

	// callback is used for handling call results.
	callback func(wire.Result, []byte)
)

// NewRPC returns a new RPC object.
// RPC immediately tries to connect to the operator and starts to handle
// incomming data.
// You may want to call Subscribe afterwards if you need balance and/or
// deposit proofs.
func NewRPC(host string, port uint16) (*RPC, error) {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", host, port), Path: "/ws"}
	conn, _, err := gorilla.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	rpc := &RPC{
		conn:      conn,
		callbacks: make(map[wire.ID]callback),
		Embedding: perunlog.MakeEmbedding(perunlog.WithField("role", "client")),
	}
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		return nil, err
	}
	conn.SetCloseHandler(func(code int, err string) error {
		rpc.Log().WithField("code", code).Errorf("WS connection was closed: %s", err)
		return rpc.Close()
	})
	go rpc.handleConnections()

	return rpc, nil
}

// SendTx sends one transaction to the operator.
func (r *RPC) SendTx(ctx context.Context, tx tee.Transaction) error {
	call := wire.NewSendTx(r.nextID(), tx)
	errChan := make(chan error)
	// Setup async response cb.
	r.registerCallback(call.Call.ID, func(result wire.Result, msg []byte) {
		if result.Error != "" {
			errChan <- fmt.Errorf("SendTx RPC result: %s", result.Error)
		} else {
			errChan <- nil
		}
	})
	// Make the call.
	if err := r.sendJSON(call); err != nil {
		return fmt.Errorf("sending json object: %w", err)
	}
	// Return error from async response cb.
	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Subscribe subscribes to the Balance and Deposit proof topic.
// The user must always read the proofs via `DepositProof` and
// `BalanceProof`. Calling this function more than once if it did not
// error will cause undefined behaviour.
func (r *RPC) Subscribe(ctx context.Context, user common.Address) (*Subscription, error) {
	r.subscription = &Subscription{
		Embedding: perunlog.MakeEmbedding(perunlog.WithField("role", "proofSub")),
		// Buffer the proofs here, otherwise the client has to read them
		// immediately to prevent that they get reordered by a race condition
		// from the go-routines writing them to the channel since a mutex is
		// not FIFO.
		balProofs: make(chan tee.BalanceProof, 10),
		depProofs: make(chan tee.DepositProof, 10),
	}

	call := wire.NewSubscribe(r.nextID(), user)
	errChan := make(chan error)
	// Setup async response cb.
	r.registerCallback(call.Call.ID, func(result wire.Result, msg []byte) {
		if result.Error != "" {
			errChan <- fmt.Errorf("Subscribe RPC result: %s", result.Error)
		} else {
			errChan <- nil
		}
	})
	// Make the call.
	if err := r.sendJSON(call); err != nil {
		return nil, fmt.Errorf("sending json object: %w", err)
	}
	// Return error from async response cb.
	select {
	case err := <-errChan:
		return r.subscription, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r *RPC) handleConnections() {
	for !r.IsClosed() {
		// gorilla has no async read method?!
		_, data, err := r.conn.ReadMessage()
		if err != nil {
			r.Log().Error("reading ws message: ", err)
			time.Sleep(1 * time.Second)
			continue
		}
		r.Log().Trace("client received: ", string(data))
		var msg wire.Result
		if err := json.Unmarshal(data, &msg); err != nil {
			r.Log().Error("decoding message: ", err)
			continue
		}

		switch {
		case msg.ID != "":
			r.callCallback(msg, data)
		case msg.Topic != "":
			if r.subscription == nil {
				r.Log().Error("Received proof without subscription")
				return
			}
			r.subscription.handleTopic(msg.Topic, data)
		default:
			r.Log().Error("Received result without ID or Topic")
		}
	}
}

func (r *RPC) registerCallback(id wire.ID, cb callback) {
	r.cbMtx.Lock()
	defer r.cbMtx.Unlock()
	if _, ok := r.callbacks[id]; ok {
		r.Log().WithField("id", id).Error("Callback already registered, skipping.")
		return
	}
	r.callbacks[id] = cb
}

func (r *RPC) callCallback(result wire.Result, data []byte) {
	r.cbMtx.RLock()
	cb, ok := r.callbacks[result.ID]
	r.cbMtx.RUnlock()

	if !ok {
		r.Log().WithField("id", result.ID).Error("unknown result id")
		return
	}
	cb(result, data)

	r.cbMtx.Lock()
	delete(r.callbacks, result.ID)
	r.cbMtx.Unlock()
}

func (r *RPC) nextID() wire.ID {
	// This does not ensure that ID always increments in messages that are
	// sent out, but it ensures that they are always different.
	id := atomic.AddUint64(&r.id, 1)
	return wire.ID(strconv.FormatUint(id, 10))
}

func (r *RPC) sendJSON(obj interface{}) error {
	r.connMtx.Lock()
	defer r.connMtx.Unlock()
	return r.conn.WriteJSON(obj)
}

func (s *Subscription) handleTopic(topic wire.Topic, data []byte) {
	switch topic {
	case wire.DepositProofs:
		var msg wire.DepositProof
		if err := json.Unmarshal(data, &msg); err != nil {
			s.Log().WithError(err).Error("decoding json")
			return
		}
		s.depProofs <- msg.Proof
		s.Log().WithField("epoch", msg.Proof.Balance.Epoch).Trace("Received deposit proof")
	case wire.BalanceProofs:
		var msg wire.BalanceProof
		if err := json.Unmarshal(data, &msg); err != nil {
			s.Log().WithError(err).Error("decoding json")
			return
		}
		s.balProofs <- msg.Proof
		s.Log().WithField("epoch", msg.Proof.Balance.Epoch).Trace("Received balance proof")
	default:
		s.Log().WithField("topic", topic).Error("unknown result topic")
	}
}

// DepositProof blocks until it can return the next deposit proof from the
// operator or an error if the context ran out.
func (s *Subscription) DepositProof(ctx context.Context) (tee.DepositProof, error) {
	select {
	case <-ctx.Done():
		return tee.DepositProof{}, ctx.Err()
	case proof := <-s.depProofs:
		return proof, nil
	}
}

// BalanceProof blocks until it can return the next balance proof from the
// operator or an error if the context ran out.
func (s *Subscription) BalanceProof(ctx context.Context) (tee.BalanceProof, error) {
	select {
	case <-ctx.Done():
		return tee.BalanceProof{}, ctx.Err()
	case proof := <-s.balProofs:
		return proof, nil
	}
}
