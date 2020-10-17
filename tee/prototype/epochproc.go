// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	err "perun.network/go-perun/pkg/errors"

	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/tee"
)

func (e *Enclave) epochProcessor(
	verifiedBlocks <-chan *tee.Block,
	txs <-chan *tee.Transaction,
) error {
	var (
		// TODO: handle nil tx/exit epoch
		depositEpoch = NewEpoch(0)
		txEpoch      *Epoch
		exitEpoch    *Epoch
		phaseShift   = make(chan struct{})
		errg         = err.NewGatherer()
	)
	// push first epoch
	e.epochs.Push(depositEpoch)

	for {
		select {
		case <-e.quit:
			return errg.Err()
		default:
		}
		done := make(chan struct{}) // TODO: change to exiters slice

		errg.Go(func() error {
			return e.depositExitRoutine(verifiedBlocks, done, depositEpoch, exitEpoch)
		})

		errg.Go(func() error {
			return e.txRoutine(phaseShift, done, txs, txEpoch)
		})

		select {
		case <-phaseShift:
			e.shiftEpochs(&depositEpoch, &txEpoch, &exitEpoch)
		case <-errg.Failed():
			return errg.Err()
		}
	}
}

func (e *Enclave) txRoutine(
	phaseShift, done chan struct{},
	txs <-chan *tee.Transaction,
	txEpoch *Epoch,
) error {
	// read tx from txs (tx phase).
	// TODO: let enclave generate all balance proofs beforehand and just update
	//		 them for incoming transactions.
	//		  -> Right now all balance proofs are generated when a txPhase is
	//           done.
	for {
		select {
		case <-done:
			// TODO: pass set of exiting users from exit processor here
			//   then check that these users didn't perform txs. If they did, error for now now
			//   Best would be to revert all sending and receiving transactions of exiters.
			balProofs, err := e.generateBalanceProofs(txEpoch)
			if err != nil {
				return fmt.Errorf("generating balance proofs: %w", err)
			}
			e.balanceProofs <- balProofs
			phaseShift <- struct{}{} // phase done, stop processing TXs and signal epoch shift.
			return nil
		case tx := <-txs:
			// extract tx-changes, adjust epoch's balances.
			if err := e.adjustEpochTx(txEpoch, tx); err != nil {
				// TODO: racing transactions.
				//		  && not synced with validated block-input.
				//		  we cant error, else we shutdown the enclave!
				return fmt.Errorf("adjusting Epoch %v Balances: %w", txEpoch.Number, err)
			}
		}
	}
}

func (e *Enclave) depositExitRoutine(
	verifiedBlocks <-chan *tee.Block,
	phaseDone chan<- struct{},
	depositEpoch, exitEpoch *Epoch,
) error {
	// read blocks from verifiedBlocks (deposit phase).
	for vb := range verifiedBlocks {
		if err := e.handleVerifiedBlock(depositEpoch, exitEpoch, vb); err != nil {
			return fmt.Errorf("handling verified blocknr %v : %w", vb.NumberU64(), err)
		}
		if e.phaseDone(vb.NumberU64()) {
			e.depositProofs <- retrieveCachedDepProofs(e.depositProofCache)
			e.depositProofCache = make(map[common.Address]*tee.DepositProof)
			close(phaseDone) // phase done, stop TX processor.
			return nil
		}
	}
	return nil
}

func retrieveCachedDepProofs(cache map[common.Address]*tee.DepositProof) []*tee.DepositProof {
	dps := make([]*tee.DepositProof, 0, len(cache))
	for _, dp := range cache {
		dps = append(dps, dp)
	}
	return dps
}

// generateDepositProof generates the deposit proof for the given user in the given
// deposit Epoch.
func (e *Enclave) generateDepositProof(depEpoch *Epoch, acc common.Address) (*tee.DepositProof, error) {
	b := tee.Balance{
		Epoch:   depEpoch.Number,
		Account: acc,
		Value:   new(big.Int).Set(depEpoch.balances[acc].Value),
	}

	msg, err := tee.EncodeDepositProof(e.params.Contract, b)
	if err != nil {
		return nil, fmt.Errorf("encoding deposit proof: %w", err)
	}

	sig, err := e.wallet.SignText(e.account, crypto.Keccak256(msg))
	if err != nil {
		return nil, fmt.Errorf("signing deposit proof: %w", err)
	}

	return &tee.DepositProof{
		Balance: b,
		Sig:     sig,
	}, nil
}

// generateBalanceProofs generates the balance proofs for all users in the given
// transaction Epoch.
func (e *Enclave) generateBalanceProofs(txEpoch *Epoch) ([]*tee.BalanceProof, error) {
	balProofs := make([]*tee.BalanceProof, 0, len(txEpoch.balances))
	// TODO: counter plox
	for acc, bal := range txEpoch.balances {
		b := tee.Balance{
			Epoch:   txEpoch.Number,
			Account: acc,
			Value:   bal.Value,
		}
		msg, err := tee.EncodeBalanceProof(e.params.Contract, b)
		if err != nil {
			return nil, fmt.Errorf("encoding balance proof: %w", err)
		}

		sig, err := e.wallet.SignText(e.account, crypto.Keccak256(msg))
		if err != nil {
			return nil, fmt.Errorf("signing balance proof: %w", err)
		}

		balProofs = append(balProofs, &tee.BalanceProof{
			Balance: b,
			Sig:     sig,
		})
	}
	return balProofs, nil
}

// shiftEpochs shifts the given three epochs by one phase.
func (e *Enclave) shiftEpochs(depositEpoch, txEpoch, exitEpoch **Epoch) {
	*exitEpoch = *txEpoch
	*txEpoch = *depositEpoch
	*depositEpoch = (*txEpoch).NewNext()
	e.epochs.Push(*depositEpoch)
}

func (e *Enclave) phaseDone(blocknr uint64) bool {
	return (blocknr % e.params.PhaseDuration) == e.params.PhaseDuration-1
}

// handleVerifiedBlock receives a verified block and adjusts the transaction
// Epoch as well as the exit Epoch.
func (e *Enclave) handleVerifiedBlock(depEpoch, exitEpoch *Epoch, vb *tee.Block) error {
	for _, r := range vb.Receipts {
		exLogs, depLogs := partition(r.Logs, logIsDepositEvt)
		// extract deposits, adjust epoch's balances.
		if err := e.adjustEpochDeposit(e.params.Contract, depEpoch, depLogs); err != nil {
			return fmt.Errorf("adjusting Epoch %v Balances: %w", depEpoch.Number, err)
		}
		// extract exits, adjust epoch's balances.
		if err := e.adjustEpochExit(exitEpoch, exLogs); err != nil {
			return fmt.Errorf("handling exit of Epoch %v: %w", exitEpoch.Number, err)
		}
	}
	return nil
}

type noPredLogs = []*types.Log
type predLogs = []*types.Log

// partition partitions a given slice of `*types.Log` into a slice where `pred`
// holds and vice versa.
func partition(ls []*types.Log, pred func(l *types.Log) bool) (predLogs, noPredLogs) {
	var depLogs, exLogs []*types.Log
	for _, l := range ls {
		if pred(l) {
			depLogs = append(depLogs, l)
		} else {
			exLogs = append(exLogs, l)
		}
	}
	return depLogs, exLogs
}

// TODO: rephrase, Go doesn't allow for complex global consts.
var depositedEvent common.Hash = crypto.Keccak256Hash([]byte("Deposited(uint64,address,uint256)"))

var exitingEvent common.Hash = crypto.Keccak256Hash([]byte("Exiting(uint64,address,uint256)"))

func logIsDepositEvt(l *types.Log) bool {
	return l.Topics[0].String() == depositedEvent.String()
}

// adjustEpochDeposit adjusts `e.balances` according to the deposits done in
// the given block.
func (e *Enclave) adjustEpochDeposit(contract common.Address, ep *Epoch, depLogs []*types.Log) error {
	for _, depLog := range depLogs {
		deposit, err := parseEvent(depLog, "Deposited")
		if err != nil {
			return fmt.Errorf("parsing withdraw event: %w", err)
		}
		if accBal, ok := ep.balances[deposit.Account]; ok {
			accBal.Value.Add(accBal.Value, deposit.Value)
		} else {
			ep.balances[deposit.Account] = Bal{
				0,
				new(big.Int).Set(deposit.Value),
			}
		}
		if deposit.Epoch != ep.Number {
			return fmt.Errorf("deposit-event Epoch %v != current deposit Epoch %v",
				deposit.Epoch, ep.Number)
		}

		depProof, err := e.generateDepositProof(ep, deposit.Account)
		if err != nil {
			return fmt.Errorf("generating deposit proof: %w", err)
		}
		e.depositProofCache[deposit.Account] = depProof
	}
	return nil
}

// adjustEpochExit handles the exit phase of given Epoch.
func (e *Enclave) adjustEpochExit(ep *Epoch, exLogs []*types.Log) error {
	for _, exLog := range exLogs {
		exit, err := parseEvent(exLog, "Exiting")
		if err != nil {
			return fmt.Errorf("parsing exiting event: %w", err)
		}
		if exit.Epoch != ep.Number {
			return fmt.Errorf("exit-event Epoch %v != current exit Epoch %v",
				exit.Epoch, ep.Number)
		}
		// only full exits supported currently
		accBal := ep.balances[exit.Account].Value
		if accBal.Cmp(exit.Value) != 0 {
			return fmt.Errorf("unexpected partial exit for %v, expected %v", exit.Value, accBal)
		}
		accBal.SetUint64(0)
		delete(ep.balances, exit.Account)
		// add exit.Account to set of exiters ([]common.Address)
	}

	return nil
}

// erdstallEvent is a generic wrapper type for `Deposited`, `Exiting` and
// `Withdrawn` solidity events.
type erdstallEvent struct {
	Epoch   uint64
	Account common.Address
	Value   *big.Int
}

var contractAbi, _ = abi.JSON(strings.NewReader(bindings.ErdstallABI))

// parseEvent parses a given `log` and returns an `erdstallEvent`.
func parseEvent(l *types.Log, name string) (*erdstallEvent, error) {
	event := new(erdstallEvent)
	err := contractAbi.Unpack(event, name, l.Data)
	if err != nil {
		return nil, fmt.Errorf("unpacking %v : %w", name, err)
	}

	return event, nil
}

// adjustEpochTx adjusts `e.balances` according to given transactions.
func (e *Enclave) adjustEpochTx(ep *Epoch, tx *tee.Transaction) error {
	const LT = -1

	if valid, err := tee.VerifyTransaction(e.params, *tx); err != nil {
		return fmt.Errorf("verifying tx signature: %w", err)
	} else if !valid {
		return fmt.Errorf("invalid tx signature")
	}

	// TODO: cache future transaction until we process them
	if tx.Epoch != ep.Number {
		// TODO: log that we drop
		return nil
	}

	sender, oks := ep.balances[tx.Sender]
	recipient, okr := ep.balances[tx.Recipient]
	if !oks {
		return fmt.Errorf("unknown sender: %x", tx.Sender)
	}
	if !okr {
		return fmt.Errorf("unknown recipient: %x", tx.Recipient)
	}
	if sender.Value.Cmp(tx.Amount) == LT {
		return fmt.Errorf("tx amount exceeds senders balance: has: %v, needs: %v", sender.Value, tx.Amount)
	}
	if tx.Nonce != sender.Nonce+1 {
		return fmt.Errorf("comparing tx nonce: %v, sender nonce: %v",
			tx.Nonce, sender.Nonce)
	}

	sender.Value.Sub(sender.Value, tx.Amount)
	recipient.Value.Add(recipient.Value, tx.Amount)

	sender.Nonce = tx.Nonce
	return nil
}

var mismatchedTxEpochErr error = errors.New("mismatched Epochs for Tx")
