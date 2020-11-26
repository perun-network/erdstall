// SPDX-License-Identifier: Apache-2.0

package wire

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/tee"
)

type (
	// Call is a generic call.
	Call struct {
		// ID is an unique id.
		ID ID `json:"id"`
		// Method describes which action should be executed.
		Method Method `json:"method"`
	}

	// Result is a generic result.
	// It has either ID or Topic set, depending on whether it is a result
	// for a call or a subscription update for a topic.
	Result struct {
		// ID can be used to map the Result fo its corresponding `Call`.
		ID ID `json:"id,omitempty"`
		// Topic is set iff this Result to the specified Topic subscription.
		Topic Topic `json:"topic,omitempty"`
		// Error is set iff an error occurred with a Call or subscription.
		Error string `json:"error,omitempty"`
	}

	// SendTx sends one Transaction to the remote operator.
	SendTx struct {
		Call
		Tx tee.Transaction `json:"tx"`
	}

	// Subscribe sets up a client subscription.
	Subscribe struct {
		Call
		// The Address this subscription is for.
		Who common.Address `json:"who"`
	}

	// DepositProof contains one DepositProof from the according subscription.
	DepositProof struct {
		Result
		Proof tee.DepositProof `json:"proof"`
	}

	// BalanceProof contains one BalanceProof from the according subscription.
	BalanceProof struct {
		Result
		Proof tee.BalanceProof `json:"proof"`
	}

	// Method describes a method to RPC call.
	Method string
	// Topic describes a topic to RPC subscribe on.
	Topic string
	// ID unanimously identifies a RPC call and result.
	ID string
)

const (
	MethodSendTx    Method = "sendTx"
	MethodSubscribe Method = "subscribe"

	BalanceProofs Topic = "balanceProofs"
	DepositProofs Topic = "depositProofs"
)

// NewSendTx returns a `SendTx` object.
func NewSendTx(id ID, tx tee.Transaction) *SendTx {
	return &SendTx{
		Call: Call{
			ID:     id,
			Method: MethodSendTx,
		},
		Tx: tx,
	}
}

// NewSubscribe returns a `Subscribe` object.
func NewSubscribe(id ID, who common.Address) *Subscribe {
	return &Subscribe{
		Call: Call{
			ID:     id,
			Method: MethodSubscribe,
		},
		Who: who,
	}
}
