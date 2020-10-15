// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"

	"github.com/perun-network/erdstall/tee"
)

type (
	Enclave struct {
		params  tee.Parameters
		wallet  accounts.Wallet
		account accounts.Account

		bc     blockchain // TODO: do we need? we process on the fly...
		epochs epochchain
	}
)

// Init initializes the enclave, generating a new secp256k1 ECDSA key and
// storing it as the enclave's signing key.
//
// It returns the public key derived Ethereum address and attestation of
// correct initialization of the enclave with the generated address. The
// attestation can be used to verify the enclave with the TEE vendor.
//
// The Operator must deploy the contract with the Enclave's address after
// calling Init.
func (e *Enclave) Init() (common.Address, []byte, error) {
	panic("not implemented") // TODO: Implement
}

// ProcessBlocks should be called by the Operator to cause the enclave to
// process the given block(s), logging deposits and exits.
//
// Note that BalanceProofs requires an additional k blocks to be known to
// the Enclave before it reveals an epoch's balance proofs to the operator.
// k is a security parameter to guarantee enough PoW depth.
func (e *Enclave) ProcessBlocks(_ ...*tee.Block) error {
	panic("not implemented") // TODO: Implement
}

// ProcessTXs should be called by the Operator whenever they receive new
// transactions from users. After a transaction epoch has finished and an
// additional k blocks made known to the Enclave, the epoch's balance proofs
// can be received by calling BalanceProofs.
func (e *Enclave) ProcessTXs(_ ...*tee.Transaction) error {
	panic("not implemented") // TODO: Implement
}

// BalanceProofs returns the balance proofs of the given epoch. Note that
// all blocks of the epoch's transaction phase (+k) need to be known to the
// Enclave.
func (e *Enclave) BalanceProofs(_ tee.Epoch) ([]tee.BalanceProof, error) {
	panic("not implemented") // TODO: Implement
}
