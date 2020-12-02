// SPDX-License-Identifier: Apache-2.0

package tee

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Transaction is a payment transaction from Sender to Recipient, signed by
// the Sender. The epoch must match the current transaction epoch.
//
// Nonce tracking allows to send multiple transactions per epoch, each only
// stating the amount of the individual transaction.
type Transaction struct {
	Nonce     uint64         `json:"nonce"` // starts at 0, each tx must increase by one, across epochs
	Epoch     Epoch          `json:"epoch"`
	Sender    common.Address `json:"sender"`
	Recipient common.Address `json:"recipient"`
	Amount    *Amount        `json:"amount"`
	Sig       Sig            `json:"sig"`
	hash      common.Hash
}

// Hash hashes a transaction.
func (t *Transaction) Hash() common.Hash {
	if t.hash != [32]byte{} {
		return t.hash
	}
	bs := make([]byte, 16)
	binary.LittleEndian.PutUint64(bs, t.Nonce)
	binary.LittleEndian.PutUint64(bs[7:], t.Epoch)
	bs = append(bs, t.Sender.Bytes()...)
	bs = append(bs, t.Recipient.Bytes()...)
	bs = append(bs, t.Sig...)

	return crypto.Keccak256Hash(bs)
}

// A TextSigner can sign messages. It's usually a wallet like an accounts.Wallet
type TextSigner interface {
	// SignText returns a signature with 'v' value of 0 or 1.
	SignText(account accounts.Account, text []byte) ([]byte, error)
}

// Sign signs the transaction with the given account and signer. It checks that
// the account matches the transaction's sender.
func (t *Transaction) Sign(contract common.Address, account accounts.Account, w TextSigner) error {
	if account.Address != t.Sender {
		return errors.New("not Sender's account")
	}
	return t.SignAlien(contract, account, w)
}

// SignAlien signs the transaction with the given account and signer,
// irrespective of who's the transaction's sender. Should only be used in
// testing.
func (t *Transaction) SignAlien(contract common.Address, account accounts.Account, w TextSigner) error {
	msg, err := EncodeTransaction(contract, *t)
	if err != nil {
		return fmt.Errorf("encoding tx: %w", err)
	}
	hash := crypto.Keccak256Hash(msg)
	sig, err := w.SignText(account, hash[:])
	if err != nil {
		return fmt.Errorf("signing tx hash: %w", err)
	}
	sig[64] += 27

	t.Sig = sig
	return nil
}
