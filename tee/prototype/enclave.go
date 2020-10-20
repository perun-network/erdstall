// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	perrors "perun.network/go-perun/pkg/errors"
	"perun.network/go-perun/pkg/sync/atomic"

	"github.com/perun-network/erdstall/tee"
)

type (
	Enclave struct {
		params  tee.Parameters
		wallet  accounts.Wallet
		account accounts.Account

		bc     blockchain // TODO: do we need? we process on the fly...
		epochs epochchain

		// incoming data
		newBlocks chan *tee.Block
		newTXs    chan *tee.Transaction

		// outgoing data
		depositProofs chan []*tee.DepositProof
		balanceProofs chan []*tee.BalanceProof

		// cache
		depositProofCache map[common.Address]*tee.DepositProof

		running atomic.Bool   // if false, signals processors to return after sealing epoch
		done    chan struct{} // signal by processors that they're done
	}
)

var _ (tee.Enclave) = (*Enclave)(nil) // compile-time check

const (
	// TODO: do we even want buffering? It may also be ok if ProcessX calls block.
	bufSizeBlocks = 0 // incoming blocks buffer size
	bufSizeTXs    = 0 // incoming tx buffer size
	bufSizeProofs = 1 // proofs buffer size in #epochs
)

var ErrEnclaveStopped = errors.New("Enclave stopped")

func NewEnclave(wallet accounts.Wallet) *Enclave {
	return &Enclave{
		wallet:            wallet,
		newBlocks:         make(chan *tee.Block, bufSizeBlocks),
		newTXs:            make(chan *tee.Transaction, bufSizeTXs),
		depositProofs:     make(chan []*tee.DepositProof, bufSizeProofs),
		balanceProofs:     make(chan []*tee.BalanceProof, bufSizeProofs),
		depositProofCache: make(map[common.Address]*tee.DepositProof),
		done:              make(chan struct{}),
	}
}

func NewEnclaveWithAccount(wallet accounts.Wallet, account accounts.Account) *Enclave {
	e := NewEnclave(wallet)
	e.account = account
	return e
}

// Start starts the enclave routines and blocks until they return. It returns an
// error gatherer of all errors that the routines return, if any.
//
// Start must be called after Init and SetParams.
//
// Start can be stopped by calling Stop. However, Start will process blocks and
// transactions until the current phase has finished.
func (e *Enclave) Start() error {
	if !e.running.TrySet() {
		panic("Enclave already running")
	}

	var (
		verifiedBlocks = make(chan *tee.Block, bufSizeBlocks) // connects the block and epoch processors
		errg           = perrors.NewGatherer()
	)

	errg.Go(func() error {
		return e.blockProcessor(e.newBlocks, verifiedBlocks)
	})
	errg.Go(func() error {
		return e.epochProcessor(verifiedBlocks, e.newTXs)
	})

	return errg.Wait()
}

// Stop lets the Enclave gracefully shutdown after the next phase is sealed. It
// will continue receiving transactions and blocks until the last block of the
// current phase is received via ProcessBlocks.
//
// The Enclave Interface methods will return an ErrEnclaveStopped error after
// the Enclave shut down.
func (e *Enclave) Stop() {
	if !e.running.TryUnset() {
		panic("Enclave not running")
	}
}

// Init initializes the enclave, generating a new secp256k1 ECDSA key and
// storing it as the enclave's signing key.
//
// It returns the public key derived Ethereum address and attestation of
// correct initialization of the enclave with the generated address. The
// attestation can be used to verify the enclave with the TEE vendor.
//
// The Operator must deploy the contract with the Enclave's address after
// calling Init.
func (e *Enclave) Init() (tee common.Address, _ []byte, err error) {
	if e.account.Address != tee {
		return e.account.Address, nil, nil
	}
	e.account, err = e.wallet.Derive(accounts.DefaultRootDerivationPath, true)
	if err != nil {
		return
	}
	return e.account.Address, nil, nil
}

// SetParams should be called by the operator after they deployed the
// contract to set the system parameters, including the contract address.
// The Enclave verifies the parameters upon receival of the first block.
func (e *Enclave) SetParams(p tee.Parameters) error {
	if p.TEE != e.account.Address {
		return errors.New("tee address mismatch")
	}
	e.params = p
	return nil
}

// ProcessBlocks should be called by the Operator to cause the enclave to
// process the given block(s), logging deposits and exits.
//
// Note that BalanceProofs requires an additional k blocks to be known to
// the Enclave before it reveals an epoch's balance proofs to the operator.
// k is a security parameter to guarantee enough PoW depth.
func (e *Enclave) ProcessBlocks(blocks ...*tee.Block) error {
	for _, b := range blocks {
		select {
		case e.newBlocks <- b:
		case <-e.done:
			return ErrEnclaveStopped
		}
	}
	return nil
}

// ProcessTXs should be called by the Operator whenever they receive new
// transactions from users. After a transaction epoch has finished and an
// additional k blocks made known to the Enclave, the epoch's balance proofs
// can be received by calling BalanceProofs.
func (e *Enclave) ProcessTXs(txs ...*tee.Transaction) error {
	for _, tx := range txs {
		select {
		case e.newTXs <- tx:
		case <-e.done:
			return ErrEnclaveStopped
		}
	}
	return nil
}

// DepositProofs returns the deposit proofs of all deposits made in an epoch
// at the end of the deposit phase.
//
// Note that all blocks of the epoch's deposit phase (+k) need to be known
// to the Enclave. This call blocks until all necessary blocks are received
// and processed.
//
// It should be called in a loop by the operator.
func (e *Enclave) DepositProofs() ([]*tee.DepositProof, error) {
	select {
	case dps := <-e.depositProofs:
		return dps, nil
	case <-e.done:
		return nil, ErrEnclaveStopped
	}
}

// BalanceProofs returns all balance proofs at the end of each transaction
// phase.
//
// Note that all blocks of the epoch's transaction phase (+k) need to be
// known to the Enclave. This call blocks until all necessary blocks are
// received and processed.
//
// It should be called in a loop by the operator.
func (e *Enclave) BalanceProofs() ([]*tee.BalanceProof, error) {
	select {
	case bps := <-e.balanceProofs:
		return bps, nil
	case <-e.done:
		return nil, ErrEnclaveStopped
	}
}
