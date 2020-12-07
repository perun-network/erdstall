// SPDX-License-Identifier: Apache-2.0

package test

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/tee"
	wtest "perun.network/go-perun/backend/ethereum/wallet/test"
)

// NewTx creates a new unsigned random transaction.
func NewTx(rng *rand.Rand) *tee.Transaction {
	return NewTxFromTo(rng, NewRandomAddress(rng), NewRandomAddress(rng))
}

// NewTxFrom creates a new unsigned random transaction with the given sender.
func NewTxFrom(rng *rand.Rand, sender common.Address) *tee.Transaction {
	return NewTxFromTo(rng, sender, NewRandomAddress(rng))
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

// NewRandomAddress creates a new random address.
func NewRandomAddress(rng *rand.Rand) common.Address {
	return common.Address(wtest.NewRandomAddress(rng))
}
