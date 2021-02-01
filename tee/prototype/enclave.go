// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	perrors "perun.network/go-perun/pkg/errors"
	"perun.network/go-perun/pkg/sync/atomic"

	"github.com/perun-network/erdstall/tee"
)

type (
	// State contains all essential enclave state.
	State struct {
		Params *tee.Parameters // The fixed enclave parameters.

		Epoch         tee.Epoch   // The current epoch.
		LastBlock     uint64      // The last known block height.
		LastBlockHash common.Hash // The last known block's hash.

		Accounts map[common.Address]*Acc // The users' accounts.
	}

	Enclave struct {
		params *tee.Parameters // The enclave's parameters.
		*State                 // The enclave's essential state.

		account *accounts.Account
		wallet  accounts.Wallet

		chain blockchain
		epoch *Epoch // Epoch manager, nil until first block is known.

		// Incoming commands are queued here to be executed in order.
		commands chan command

		// Outgoing data: these proofs have to be consumed within one epoch.
		depositProofs chan []*tee.DepositProof
		balanceProofs chan []*tee.BalanceProof

		// cache
		depositProofCache []*tee.DepositProof // Accumulated until phase shift.

		// Running/stopping
		shutdownRequested bool // User requested shutdown.
		shutdownApproved  bool // Enclave wants to shut down.
		running           atomic.Bool
		stopped           chan struct{} //
	}
)

var _ (tee.Enclave) = (*Enclave)(nil) // compile-time check

// enclaveMaxCommandQueue is the number of commands that can be enqueued to the
// enclave simultaneously. Once this number is surpassed, it is no longer
// possible to ensure a FIFO execution of commands, and excess commands will
// result in an error.
const enclaveMaxCommandQueue = 256 // for good measure.

func NewEnclave(wallet accounts.Wallet) *Enclave {
	return &Enclave{
		wallet:        wallet,
		chain:         blockchain{},
		commands:      make(chan command, enclaveMaxCommandQueue),
		depositProofs: make(chan []*tee.DepositProof, 1),
		balanceProofs: make(chan []*tee.BalanceProof, 1),
		stopped:       make(chan struct{}),
	}
}

func NewEnclaveWithAccount(wallet accounts.Wallet, account accounts.Account) *Enclave {
	e := NewEnclave(wallet)
	e.account = &account
	return e
}

func (e *Enclave) BlockNum() uint64 {
	return e.chain.Head().NumberU64()
}

func (e *Enclave) IsAtPhaseEnd() bool {
	if e.chain.empty() {
		return false
	}
	return e.params.IsLastPhaseBlock(e.BlockNum())
}

func (e *Enclave) mainLoop() (err error) {
	panicked := true
	defer func() {
		if panicked {
			err = fmt.Errorf("panic: %v", recover())
		}
	}()

	for cmd := range e.commands {
		switch cmd := cmd.(type) {
		case *processBlocksCmd:
			errs := make([]error, len(cmd.blocks))
			for i, b := range cmd.blocks {
				errs[i] = e.processBlock(b)
			}
			cmd.result <- errs
		case *processTxsCmd:
			errs := make([]error, len(cmd.txs))
			for i, tx := range cmd.txs {
				errs[i] = e.epoch.ProcessTx(e.Params.Contract, tx)
			}
			cmd.result <- errs
		case *shutdownCmd:
			e.shutdownRequested = true
		default:
			log.Panicf("unhandled command type %T", cmd)
		}

		if e.shutdownApproved {
			break
		}
	}
	panicked = false
	return
}

// ProcessBlocks instantaneously processes a list of blocks. If any of the
// blocks are erroneous, they are simply ignored, without affecting the rest of
// the blocks. Returns the accumulated error messages.
func (e *Enclave) ProcessBlocks(blocks ...*tee.Block) error {
	if e.shutdownApproved {
		return tee.ErrEnclaveStopped
	}

	errCh := make(chan []error, 1)
	select {
	case e.commands <- &processBlocksCmd{blocks: blocks, result: errCh}:
		select {
		case errs := <-errCh:
			errg := perrors.NewGatherer()
			for _, e := range errs {
				errg.Add(e)
			}
			return errg.Err()
		case <-e.stopped:
			return tee.ErrEnclaveStopped
		}
	case <-e.stopped:
		return tee.ErrEnclaveStopped
	}
}

// ProcessTXs should be called by the Operator whenever they receive new
// transactions from users. After a transaction epoch has finished and an
// additional k blocks made known to the Enclave, the epoch's balance proofs
// can be received by calling BalanceProofs.
func (e *Enclave) ProcessTXs(txs ...*tee.Transaction) error {
	if e.shutdownApproved {
		return tee.ErrEnclaveStopped
	}

	errCh := make(chan []error, 1)
	select {
	case e.commands <- &processTxsCmd{txs: txs, result: errCh}:
		select {
		case errs := <-errCh:
			errg := perrors.NewGatherer()
			for _, e := range errs {
				errg.Add(e)
			}
			return errg.Err()
		case <-e.stopped:
			return tee.ErrEnclaveStopped
		}
	case <-e.stopped:
		return tee.ErrEnclaveStopped
	}
}

// Shutdown lets the Enclave gracefully shutdown after the next phase is sealed. It
// will continue receiving transactions and blocks until the last block of the
// current phase is received via ProcessBlocks.
//
// The Enclave Interface methods will return an tee.ErrEnclaveStopped error after
// the Enclave shut down.
func (e *Enclave) Shutdown() {
	log.Info("Enclave: shutting down when phase is finished")
	select {
	case e.commands <- &shutdownCmd{}:
	case <-e.stopped:
		log.Panic("Shutdown has been called multiple times")
	}
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
	case <-e.stopped:
		return nil, tee.ErrEnclaveStopped
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
	case <-e.stopped:
		return nil, tee.ErrEnclaveStopped
	}
}
