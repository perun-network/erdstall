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

func (t *Transaction) Sign(contract common.Address, account accounts.Account, w TextSigner) error {
	if account.Address != t.Sender {
		return errors.New("not Sender's account")
	}
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
