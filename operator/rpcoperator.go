// SPDX-License-Identifier: Apache-2.0

package operator

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/tee"
)

type (
	// WireAPI contains all functions that a client is interested in.
	// It is called `WireAPI` since the client will interact with an
	// implementation of it over the wire.
	WireAPI interface {
		Send(tee.Transaction) error
		SubscribeProofs(common.Address) (ProofSub, error)
	}

	// RPCOperator will be exposed to the client as a websocket RPC Server.
	// It implements `WireAPI`.
	RPCOperator struct {
		mtx     sync.Mutex // protects all
		enclave tee.Enclave

		subs map[common.Address]*ProofSub
	}

	// ProofSub is a subscription on TEE proofs.
	ProofSub struct {
		deposits chan tee.DepositProof
		balances chan tee.BalanceProof
		quit     chan struct{}
	}
)

var _ WireAPI = (*RPCOperator)(nil)

func NewRPCOperator(enclave tee.Enclave) *RPCOperator {
	return &RPCOperator{
		enclave: enclave,
		subs:    make(map[common.Address]*ProofSub),
	}
}

func (o *RPCOperator) Send(tx tee.Transaction) error {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	return o.enclave.ProcessTXs(&tx)
}

// SubscribeProofs returns a subscription on TEE proofs for the given address.
// The subscription buffers the most recent proof until the client retrieves it.
func (o *RPCOperator) SubscribeProofs(addr common.Address) (ProofSub, error) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	log.WithField("who", addr.Hex()).Debug("Subscribed to proofs")
	sub, ok := o.subs[addr]
	if !ok {
		o.subscribe(addr)
		sub = o.subs[addr]
	}
	return *sub, nil
}

// subscribe is an internal implementation detail and should not be called.
func (o *RPCOperator) subscribe(addr common.Address) ProofSub {
	o.subs[addr] = newProofSub()
	return *o.subs[addr]
}

func (o *RPCOperator) PushDepositProof(proof tee.DepositProof) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	who := proof.Balance.Account
	sub, ok := o.subs[who]
	if !ok {
		log.WithField("who", who.Hex()).Debug("Received DP without subscription - buffering")
		o.subscribe(who)
		sub = o.subs[who]
	}
	// Clear out old buffered proof.
	select {
	case <-sub.deposits:
	default:
	}
	// Write new proof in.
	select {
	case <-sub.quit:
		delete(o.subs, who)
	case sub.deposits <- proof:
	}
}

func (o *RPCOperator) PushBalanceProof(proof tee.BalanceProof) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	who := proof.Balance.Account
	sub, ok := o.subs[who]
	if !ok {
		log.WithField("who", who.Hex()).Debug("Received BP without subscription - buffering")
		o.subscribe(who)
		sub = o.subs[who]
	}
	// Clear out old buffered proof.
	select {
	case <-sub.balances:
	default:
	}
	// Write new proof in.
	select {
	case <-sub.quit:
		delete(o.subs, who)
	case sub.balances <- proof:
	}
}

// newProofSub returns a new proofSub. The proof channels have buffer size 1.
func newProofSub() *ProofSub {
	return &ProofSub{
		deposits: make(chan tee.DepositProof, 1),
		balances: make(chan tee.BalanceProof, 1),
		quit:     make(chan struct{}),
	}
}

func (sub *ProofSub) Deposits() <-chan tee.DepositProof {
	return sub.deposits
}

func (sub *ProofSub) Balances() <-chan tee.BalanceProof {
	return sub.balances
}

func (sub *ProofSub) Closed() <-chan struct{} {
	return sub.quit
}

func (sub *ProofSub) Unsubscribe() {
	select {
	case <-sub.quit:
	default:
		close(sub.quit)
	}
}
