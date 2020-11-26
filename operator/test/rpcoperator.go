// SPDX-License-Identifier: Apache-2.0

package test

import (
	"time"

	"github.com/ethereum/go-ethereum/common"

	op "github.com/perun-network/erdstall/operator"
	"github.com/perun-network/erdstall/tee"
)

type (
	// RPCOperator is a mocked operator.WireAPI.
	// It can be used as an WireAPI for testing the operators
	// RPC module, since it expects an operator.WireAPI.
	//
	// There are two functions for setting errors that should be returned
	// when the mocked functions are called.
	RPCOperator struct {
		enclave *Enclave
		op      *op.RPCOperator

		subscribeError error
	}
)

var _ op.WireAPI = (*RPCOperator)(nil)

// NewRPROperator returns a new mocked WireAPI.
func NewRPROperator(enclave *Enclave) *RPCOperator {
	return &RPCOperator{
		enclave: enclave,
		op:      op.NewRPCOperator(enclave),
	}
}

// Run pulls the proofs from the enclave and writes them to the subscriptions.
// This is normally done by the Operator. Returns immediately.
func (r *RPCOperator) Run() {
	go func() {
		for {
			// Mocked subscriptions never error.
			proofs, _ := r.enclave.DepositProofs()
			for _, proof := range proofs {
				r.op.PushDepositProof(*proof)
			}
			time.Sleep(100 * time.Millisecond) // poll interval
		}
	}()
	go func() {
		for {
			// Mocked subscriptions never error.
			proofs, _ := r.enclave.BalanceProofs()
			for _, proof := range proofs {
				r.op.PushBalanceProof(*proof)
			}
			time.Sleep(100 * time.Millisecond) // poll interval
		}
	}()
}

// Send is part of the operator.WireAPI interface and adds
// a Transaction to the enclave.
// Can be read back from Transactions() which bufferes one TX.
// Returns the error that was set by SetSendError.
func (r *RPCOperator) Send(tx tee.Transaction) error {
	return r.op.Send(tx)
}

// SubscribeProofs subscribed to the proofs that can be added via
// PushDepositProof and PushBalanceProof which buffers one proof.
// Returns the error that was set by SetSubscribeProofsError.
func (r *RPCOperator) SubscribeProofs(addr common.Address) (op.ProofSub, error) {
	if r.subscribeError != nil {
		return op.ProofSub{}, r.subscribeError
	}
	return r.op.SubscribeProofs(addr)
}

// SetSubscribeProofsError sets the error that should be returned
// by SubscribeProofs.
func (r *RPCOperator) SetSubscribeProofsError(err error) {
	r.subscribeError = err
}
