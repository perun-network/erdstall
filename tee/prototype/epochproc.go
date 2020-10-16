// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/tee"
	err "perun.network/go-perun/pkg/errors"
)

func (e *Enclave) epochProcessor(
	verifiedBlocks <-chan *tee.Block,
	txs <-chan *tee.Transaction,
) error {
	var (
		depositEpoch *Epoch = NewEpoch(tee.Epoch(0))
		txEpoch      *Epoch
		exitEpoch    *Epoch
		shift        chan struct{}
		errg         *err.Gatherer = err.NewGatherer()
	)

	for {
		select {
		case <-e.quit:
			return errg.Err()
		default:
			done := make(chan struct{})

			// Deposit- && Exit-Goroutine.
			errg.Go(func() error {
				// read blocks from verifiedBlocks (deposit phase).
				for vb := range verifiedBlocks {
					if err := e.handleVerifiedBlock(txEpoch, exitEpoch, vb); err != nil {
						return fmt.Errorf("handling verified blocknr %v : %w", vb.NumberU64(), err)
					}
					if e.phaseDone(vb.NumberU64()) {
						// TODO: Enclave signature for deposit to Operator.
						//		  * Use map[common.Address]DepositProof and proof absolute bals
						//			for incoming Deposit logs.
						//		  * Feed to `DepositProof()` via `e.depositProofs` chan.
						close(done) // phase done, stop TX processor.
						return nil
					}
				}
				return nil
			})

			// TX-Goroutine.
			errg.Go(func() error {
				// read tx from txs (tx phase).
				for tx := range txs {
					select {
					case <-done:
						balProofs, err := e.generateBalanceProofs(txEpoch)
						if err != nil {
							return fmt.Errorf("generating balance proofs: %w", err)
						}
						e.balanceProofs <- balProofs
						shift <- struct{}{} // phase done, stop processing TXs and signal epoch shift.
						return nil
					default:
						// extract tx-changes, adjust epoch's balances.
						if err := e.adjustEpochTx(txEpoch, tx); err != nil {
							return fmt.Errorf("adjusting Epoch %v Balances: %w", txEpoch.Number, err)
						}
					}
				}
				return nil
			})

			// responsible for closing end-of-phase signalling channels.
			select {
			case <-shift:
				e.shiftEpochs(depositEpoch, txEpoch, exitEpoch)
			case <-errg.Failed():
				return errg.Err()
			}
		}
	}
}

// generateBalanceProofs generates the balance proofs for all users in the given
// transaction Epoch.
func (e *Enclave) generateBalanceProofs(txEpoch *Epoch) ([]*tee.BalanceProof, error) {
	balProofs := make([]*tee.BalanceProof, 0, len(txEpoch.balances))
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

		sig, err := e.wallet.SignText(e.account, msg)
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
func (e *Enclave) shiftEpochs(depositEpoch, txEpoch, exitEpoch *Epoch) {
	close(depositEpoch.depositDone)
	if txEpoch != nil {
		close(txEpoch.txDone)
	}
	if exitEpoch != nil {
		close(exitEpoch.exitDone)
	}
	exitEpoch = txEpoch
	txEpoch = depositEpoch
	depositEpoch = txEpoch.NewNext()
	e.epochs.Push(depositEpoch)
}

// TODO: correct or off by one error?
func (e *Enclave) phaseDone(blocknr uint64) bool {
	return (blocknr % e.params.PhaseDuration) == 0
}

// handleVerifiedBlock receives a verified block and adjusts the transaction
// Epoch as well as the exit Epoch.
func (e *Enclave) handleVerifiedBlock(txEpoch, exitEpoch *Epoch, vb *tee.Block) error {
	for _, r := range vb.Receipts {
		exLogs, depLogs := partition(r.Logs, logIsDepositEvt)
		// extract deposits, adjust epoch's balances.
		if err := e.adjustEpochDeposit(e.params.Contract, txEpoch, depLogs); err != nil {
			return fmt.Errorf("adjusting Epoch %v Balances: %w", txEpoch.Number, err)
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
		accBal := ep.balances[deposit.Account].Value
		accBal.Add(accBal, deposit.Value)
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
		accBal := ep.balances[exit.Account].Value
		accBal.Sub(accBal, exit.Value)
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

// parseEvent parses a given `log` and returns an `erdstallEvent`.
func parseEvent(l *types.Log, name string) (*erdstallEvent, error) {
	contractAbi, err := abi.JSON(strings.NewReader(bindings.ErdstallABI))
	if err != nil {
		return nil, fmt.Errorf("creating contractAbi: %w", err)
	}
	var event *erdstallEvent
	err = contractAbi.Unpack(event, name, l.Data)
	if err != nil {
		return nil, fmt.Errorf("unpacking %v : %w", name, err)
	}

	return event, nil
}

// adjustEpochTx adjusts `e.balances` according to given transactions.
func (e *Enclave) adjustEpochTx(ep *Epoch, tx *tee.Transaction) error {
	if err := validateTx(e.params, ep, tx); err != nil {
		return fmt.Errorf("validating tx: %w", err)
	}

	from := ep.balances[tx.Sender]
	to := ep.balances[tx.Recipient]

	from.Value.Sub(from.Value, tx.Amount)
	to.Value.Add(to.Value, tx.Amount)

	from.Nonce = tx.Nonce
	to.Nonce = tx.Nonce
	return nil
}

// validateTx validates a `tee.Transaction` and performs sanity checks.
func validateTx(p tee.Parameters, e *Epoch, tx *tee.Transaction) error {
	const LT = -1

	valid, err := tee.VerifyTransaction(p, *tx)
	if err != nil {
		return fmt.Errorf("verifying tx signature: %w", err)
	}
	if !valid {
		return fmt.Errorf("invalid tx signature")
	}

	if tx.Epoch != e.Number {
		return fmt.Errorf("unexpected epoch nr, got: %v want: %v", tx.Epoch, e.Number)
	}

	sender := e.balances[tx.Sender]
	recipient := e.balances[tx.Recipient]
	if sender.Value.Cmp(tx.Amount) == LT {
		return fmt.Errorf("tx amount exceeds senders balance: has: %v, needs: %v", sender.Value, tx.Amount)
	}

	if tx.Nonce <= sender.Nonce || tx.Nonce <= recipient.Nonce {
		return fmt.Errorf("comparing tx nonce: %v, sender nonce: %v, recipient nonce: %v",
			tx.Nonce, sender.Nonce, recipient.Nonce)
	}

	return nil
}
