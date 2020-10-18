// SPDX-License-Identifier: Apache-2.0

package tee

import (
	"encoding/binary"

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
