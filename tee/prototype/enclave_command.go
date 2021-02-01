// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/tee"
)

type (
	// command represents all enclave command types.
	command interface {
		command()
	}

	processBlocksCmd struct {
		blocks []*tee.Block
		result chan<- []error
	}

	processTxsCmd struct {
		txs    []*tee.Transaction
		result chan<- []error
	}

	shutdownCmd struct{}
)

var _ command = (*processBlocksCmd)(nil)
var _ command = (*processTxsCmd)(nil)
var _ command = (*shutdownCmd)(nil)

func (processBlocksCmd) command() {}
func (processTxsCmd) command()    {}
func (shutdownCmd) command()      {}

// Run starts the enclave's main loop.
//
// Run must be called after Init.
//
// Run can be stopped by calling Shutdown. However, Run will process blocks and
// transactions until the current phase has finished.
func (e *Enclave) Run(params tee.Parameters) error {
	if !e.running.TrySet() {
		log.Panic("Enclave already running")
	}

	if err := e.setParams(params); err != nil {
		return err
	}

	return e.mainLoop()
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
func (e *Enclave) Init() (_ common.Address, _ []byte, err error) {
	if e.running.IsSet() {
		err = tee.ErrEnclaveStopped
		return
	}

	if e.account != nil {
		return e.account.Address, nil, nil
	}
	e.account = new(accounts.Account)
	*e.account, err = e.wallet.Derive(accounts.DefaultRootDerivationPath, true)
	if err != nil {
		return
	}
	return e.account.Address, nil, nil
}

func (e *Enclave) setParams(p tee.Parameters) error {
	if p.TEE != e.account.Address {
		return errors.New("tee address mismatch")
	} else if e.params != nil {
		return errors.New("params already set")
	}
	e.params = &p
	return nil
}

// processBlock processes a single block.
func (e *Enclave) processBlock(block *tee.Block) error {
	if e.shutdownApproved {
		return errors.New("Enclave terminated, does not accept new blocks.")
	}

	if e.chain.empty() {
		if n := block.NumberU64(); n > e.params.InitBlock {
			return fmt.Errorf("first block (%d) not initial Erdstall block (%d)", n, e.params.InitBlock)
		} else if n < e.params.InitBlock {
			log.Warnf("Received block (%d) < first Erdstall block (%d), ignoring.", n, e.params.InitBlock)
			return nil
		}
		e.epoch = newEpoch(0)
		e.State = &State{
			Params:        e.params,
			Epoch:         0,
			LastBlock:     block.NumberU64(),
			LastBlockHash: block.Hash(),
			Accounts:      nil,
		}
	}

	// Push block and if there is a phase shift, progress into the next epoch.
	if err := e.pushBlock(block); err != nil {
		return err
	}
	if e.IsAtPhaseEnd() {
		log.Debug("end of phase, shifting epochs")

		// Progress the epoch and enclave state.
		e.epoch.progressPhase()
		outcome := e.epoch.Outcome()
		e.State.Epoch++
		e.State.Accounts = outcome.Accounts

		// Publish all balance and deposit proofs of the epoch.
		e.balanceProofs <- e.generateBalanceProofs(outcome)
		e.depositProofs <- e.depositProofCache
		e.depositProofCache = e.depositProofCache[:0] // Clear deposit proofs.

		if e.shutdownRequested {
			e.shutdownApproved = true
			close(e.stopped)
		}
	}

	return nil
}

// pushBlock pushes a new block onto the enclave's blockchain and updates the
// enclave's state to reflect the new block. On error, the enclave remains
// unchanged.
func (e *Enclave) pushBlock(block *tee.Block) error {
	deps, exits, err := e.chain.PushVerify(block, e.Params, e.epoch)
	if err != nil {
		return fmt.Errorf("pushing block to local blockchain: %w", err)
	}

	e.State.LastBlock = block.NumberU64()
	e.State.LastBlockHash = block.Hash()

	// Cache any deposit proofs until the end of the phase.
	e.depositProofCache = append(e.depositProofCache, e.generateDepositProofs(deps...)...)

	e.epoch.ApplyDeposits(deps...)
	e.epoch.RegisterExits(exits...)

	return nil
}
