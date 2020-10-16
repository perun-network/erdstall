// SPDX-License-Identifier: Apache-2.0

package tee

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type (
	Enclave interface {
		// Init initializes the enclave, generating a new secp256k1 ECDSA key and
		// storing it as the enclave's signing key.
		//
		// It returns the public key derived Ethereum address and attestation of
		// correct initialization of the enclave with the generated address. The
		// attestation can be used to verify the enclave with the TEE vendor.
		//
		// The Operator must deploy the contract with the Enclave's address after
		// calling Init.
		Init() (common.Address, []byte, error)

		// SetParams should be called by the operator after they deployed the
		// contract to set the system parameters, including the contract address.
		// The Enclave verifies the parameters upon receival of the first block.
		SetParams(Parameters) error

		// ProcessBlocks should be called by the Operator to cause the enclave to
		// process the given block(s), logging deposits and exits.
		//
		// Note that BalanceProofs requires an additional k blocks to be known to
		// the Enclave before it reveals an epoch's balance proofs to the operator.
		// k is a security parameter to guarantee enough PoW depth.
		ProcessBlocks(...*Block) error

		// ProcessTXs should be called by the Operator whenever they receive new
		// transactions from users. After a transaction epoch has finished and an
		// additional k blocks made known to the Enclave, the epoch's balance proofs
		// can be received by calling BalanceProofs.
		ProcessTXs(...*Transaction) error

		// DepositProofs returns the deposit proofs of all deposits made in an epoch
		// at the end of the deposit phase.
		//
		// Note that all blocks of the epoch's deposit phase (+k) need to be known
		// to the Enclave. This call blocks until all necessary blocks are received
		// and processed.
		//
		// It should be called in a loop by the operator.
		DepositProofs() ([]*DepositProof, error)

		// BalanceProofs returns all balance proofs at the end of each transaction
		// phase.
		//
		// Note that all blocks of the epoch's transaction phase (+k) need to be
		// known to the Enclave. This call blocks until all necessary blocks are
		// received and processed.
		//
		// It should be called in a loop by the operator.
		BalanceProofs() ([]*BalanceProof, error)
	}

	// Epoch is the epoch counter type.
	Epoch = uint64

	// A Block is a go-ethereum block together with its receipts. go-ethereum's
	// types.Block type doesn't store the receipts...
	Block struct {
		types.Block
		Receipts types.Receipts
	}

	// Transaction is a payment transaction from Sender to Recipient, signed by
	// the Sender. The epoch must match the current transaction epoch.
	//
	// Nonce tracking allows to send multiple transactions per epoch, each only
	// stating the amount of the individual transaction.
	Transaction struct {
		Nonce     uint64 // starts at 0, each tx must increase by one, across epochs
		Epoch     Epoch
		Sender    common.Address
		Recipient common.Address
		Amount    *big.Int
		Sig       []byte
	}

	// A DepositProof is generated by the Enclave at the end of each deposit phase
	// for each account that made a deposit. The Operator has to forward those to
	// the users or risks facing an on-chain challenge.
	DepositProof struct {
		Balance Balance
		Sig     []byte
	}

	// A BalanceProof is generated by the Enclave at the end of each transaction
	// phase for each account in the system. The Operator has to forward those to
	// the users or risks facing an on-chain challenge.
	BalanceProof struct {
		Balance Balance
		Sig     []byte
	}

	// A Balance states the balance of the user with address Account in the system
	// at epoch Epoch.
	//
	// Its Solidity struct type is (uint64, address, uint256).
	Balance struct {
		Epoch   Epoch          // sol: uint64
		Account common.Address // sol: address
		Value   *big.Int       // sol: uint256
	}
)
