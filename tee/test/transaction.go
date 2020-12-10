// SPDX-License-Identifier: Apache-2.0

package test

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

// NewTx creates a new unsigned random transaction.
func NewTx(rng *rand.Rand) *tee.Transaction {
	return NewTxFromTo(rng, eth.NewRandomAddress(rng), eth.NewRandomAddress(rng))
}

// NewTxFrom creates a new unsigned random transaction with the given sender.
func NewTxFrom(rng *rand.Rand, sender common.Address) *tee.Transaction {
	return NewTxFromTo(rng, sender, eth.NewRandomAddress(rng))
}

// NewTxFromTo creates a new unsigned transaction from a random sender to
// another random recipient.
func NewTxFromTo(rng *rand.Rand, sender, recipient common.Address) *tee.Transaction {
	return &tee.Transaction{
		Nonce:     uint64(rng.Int63()),
		Epoch:     uint64(rng.Int63()),
		Sender:    sender,
		Recipient: recipient,
		Amount:    (*tee.Amount)(big.NewInt(rng.Int63())),
	}
}
