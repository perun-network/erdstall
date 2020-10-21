// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

// TODO: order function after call-stack please.

func (e *Enclave) epochProcessor(
	verifiedBlocks <-chan *tee.Block,
	txs <-chan *tee.Transaction,
) error {
	var (
		depositEpoch = NewEpoch(0)
		txEpoch      *Epoch
		exitEpoch    *Epoch
		depExErr     = make(chan error)
		txErr        = make(chan error)
	)
	// push first epoch
	e.epochs.Push(depositEpoch)

	for {
		numProcs := 2
		log := log.WithField("depositEpoch", depositEpoch.Number)

		log.Info("epochProcessor: starting new phase")

		done := make(chan exitersSet)
		go func() {
			depExErr <- e.depositExitRoutine(verifiedBlocks, done, depositEpoch, exitEpoch)
		}()

		go func() {
			txErr <- e.txRoutine(done, txs, txEpoch)
		}()

		log.Debug("epochProcessor: waiting for depositExitRoutine and txRoutine...")
		for numProcs != 0 {
			select {
			case err := <-depExErr:
				if err != nil {
					return fmt.Errorf("depositExitRoutine: %w", err)
				}
			case err := <-txErr:
				if err != nil {
					return fmt.Errorf("txRoutine: %w", err)
				}
			}
			numProcs--
		}

		if !e.running.IsSet() {
			log.Info("epochProcessor: routines returned, shutting down")
			close(e.done) // signal to external callers that Enclave processes are done.
			return nil
		}

		log.Debug("epochProcessor: routines returned, shifting epochs")
		e.processEpochShift(&depositEpoch, &txEpoch, &exitEpoch)
	}
}

// TODO: racing transactions.
func (e *Enclave) txRoutine(
	exits <-chan exitersSet,
	txs <-chan *tee.Transaction,
	txEpoch *Epoch,
) error {
	log := log.WithField("epoch", "<not started>")
	if txEpoch != nil {
		log = log.WithField("epoch", txEpoch.Number)
	}
	for {
		stagedTxs := makeTxCache()
		select {
		case exiters := <-exits:
			log.Trace("txRoutine: exiters received")
			// TODO: check for inconsistent deposits.
			if err := noInconsistentExits(&stagedTxs, exiters); err != nil {
				log.Errorf("handling exiters: %v", err)
			}
			bps, err := e.generateBalanceProofs(txEpoch)
			if err != nil {
				return fmt.Errorf("generating balance proofs: %w", err)
			}

			log.Debug("txRoutine: tx phase done, pushing balance proofs")
			e.balanceProofs <- bps

			log.Trace("txRoutine: return")
			return nil
		case tx := <-txs:
			stagedTxs.cacheTx(tx)
			err := e.applyEpochTx(txEpoch, tx)
			if err != nil {
				return fmt.Errorf("adjusting Epoch %v Balances: %w", txEpoch.Number, err)
			}
		}
	}
}

// noInconsistentExits checks that none of the exited parties tried to submit a
// transaction.
// TODO: check for inconsistent deposits.
func noInconsistentExits(txc *txCache, exiters exitersSet) error {
	for _, e := range exiters {
		sTxs, sOk := txc.senders[e]
		rTxs, rOk := txc.recipients[e]
		if sOk {
			return fmt.Errorf("exited user %v is sender for %v transactions", e.String(), len(sTxs))
		}
		if rOk {
			return fmt.Errorf("exited user %v is recipient for %v transactions", e.String(), len(rTxs))
		}
	}
	return nil
}

func (e *Enclave) depositExitRoutine(
	verifiedBlocks <-chan *tee.Block,
	exits chan<- exitersSet,
	depositEpoch, exitEpoch *Epoch,
) error {
	var exiters exitersSet
	// read blocks from verifiedBlocks (deposit phase).
	for vb := range verifiedBlocks {
		log := log.WithField("blockNum", vb.NumberU64())
		log.Trace("depositExitRoutine: received block")
		exs, err := e.handleVerifiedBlock(depositEpoch, exitEpoch, vb)
		if err != nil {
			return fmt.Errorf("handling verified blocknr %v: %w", vb.NumberU64(), err)
		}
		exiters = append(exiters, exs...)
		if e.params.IsLastPhaseBlock(vb.NumberU64()) {
			log.Debug("depositExitRoutine: last block of phase, pushing deposit proofs and return")
			e.pushDepositProofs()
			exits <- exiters // phase done, sending exiters to TX processor.
			return nil
		}
	}
	return errors.New("depositExitRoutine: verifiedBlocks channel closed")
}

func (e *Enclave) pushDepositProofs() {
	e.depositProofs <- asDepProofs(e.depositProofCache)
	e.depositProofCache = make(map[common.Address]*tee.DepositProof)
}

// asDepProofs reduces the deposit proof cache to a slice of `tee.DepositProof`s.
func asDepProofs(cache map[common.Address]*tee.DepositProof) []*tee.DepositProof {
	dps := make([]*tee.DepositProof, len(cache))
	i := 0
	for _, dp := range cache {
		dps[i] = dp
		i++
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
	if txEpoch == nil {
		// during the very first phase of the first epoch, there is no tx phase
		return nil, nil
	}

	balProofs := make([]*tee.BalanceProof, len(txEpoch.balances))
	i := 0
	for acc, bal := range txEpoch.balances {
		b := tee.Balance{
			Epoch:   txEpoch.Number,
			Account: acc,
			Value:   bal.Value,
		}

		bp, err := e.signBalanceProof(b)
		if err != nil {
			return nil, fmt.Errorf("generating balance proofs: %w", err)
		}
		balProofs[i] = bp
		i++
	}
	return balProofs, nil
}

// signBalanceProof signs the given `tee.Balance` and returns a `tee.BalanceProof`
// containing the corresponding signature w.r.t. the enclave.
func (e *Enclave) signBalanceProof(b tee.Balance) (*tee.BalanceProof, error) {
	msg, err := tee.EncodeBalanceProof(e.params.Contract, b)
	if err != nil {
		return nil, fmt.Errorf("encoding balance proof: %w", err)
	}

	sig, err := e.wallet.SignText(e.account, crypto.Keccak256(msg))
	if err != nil {
		return nil, fmt.Errorf("signing balance proof: %w", err)
	}
	return &tee.BalanceProof{
		Balance: b,
		Sig:     sig,
	}, nil
}

// processEpochShift shifts the given three epochs by one phase.
func (e *Enclave) processEpochShift(depositEpoch, txEpoch, exitEpoch **Epoch) {
	*exitEpoch = *txEpoch
	*txEpoch = (*depositEpoch).merge(*txEpoch)
	*depositEpoch = (*txEpoch).NewNext()
	log.Tracef("epochProc: pushing new deposit epoch: %d", (*depositEpoch).Number)
	e.epochs.Push(*depositEpoch)
	log.Tracef("epochProc: epochs shifted, new deposit epoch: %d", (*depositEpoch).Number)
}

// handleVerifiedBlock receives a verified block and adjusts the transaction
// Epoch as well as the exit Epoch.
func (e *Enclave) handleVerifiedBlock(depEpoch, exitEpoch *Epoch, vb *tee.Block) (exitersSet, error) {
	var exiters exitersSet
	for _, r := range vb.Receipts {
		logss := e.filterLogs(r.Logs, []logPredicate{logIsDepositEvt, logIsExitEvt})
		depLogs, exLogs := logss[0], logss[1]
		log.WithField("blockNum", vb.Block.NumberU64()).
			Tracef("New block has %d Deposits, %d Exits", len(depLogs), len(exLogs))
		// extract deposits, adjust epoch's balances.
		if err := e.applyEpochDeposit(depEpoch, depLogs); err != nil {
			return nil, fmt.Errorf("applying epoch %d deposits: %w", depEpoch.Number, err)
		}
		// extract exits, adjust epoch's balances.
		exits, err := e.applyEpochExit(exitEpoch, exLogs)
		if err != nil {
			return nil, fmt.Errorf("applying epoch %d exits: %w", exitEpoch.Number, err)
		}
		exiters = append(exiters, exits...)
	}
	return exiters, nil
}

// filterLogs partitions logs into different buckets of matching predicates.
// Only logs from the Erdstall contract are filtered and other logs, as well as
// those without a matching predicate, are discarded.
func (e *Enclave) filterLogs(logs []*types.Log, preds []logPredicate) [][]*types.Log {
	logss := make([][]*types.Log, len(preds))
	for _, l := range logs {
		if l.Address != e.params.Contract {
			// only parse Erdstall logs
			continue
		}
		for i, p := range preds {
			if p(l) {
				logss[i] = append(logss[i], l)
			}
		}
	}
	return logss
}

// applyEpochDeposit adjusts `e.balances` according to the deposits done in
// the given block.
func (e *Enclave) applyEpochDeposit(ep *Epoch, depLogs []*types.Log) error {
	for _, depLog := range depLogs {
		deposit, err := parseDepEvent(depLog)
		if err != nil {
			return fmt.Errorf("parsing withdraw event: %w", err)
		}

		if deposit.Epoch != ep.Number {
			return fmt.Errorf("deposit-event Epoch %v != current deposit Epoch %v",
				deposit.Epoch, ep.Number)
		}

		if accBal, ok := ep.balances[deposit.Account]; ok {
			accBal.Value.Add(accBal.Value, deposit.Value)
		} else {
			ep.balances[deposit.Account] = &Bal{
				Nonce: 0,
				Value: new(big.Int).Set(deposit.Value),
			}
		}

		depProof, err := e.generateDepositProof(ep, deposit.Account)
		if err != nil {
			return fmt.Errorf("generating deposit proof: %w", err)
		}

		log.WithFields(log.Fields{
			"account": deposit.Account.String(),
			"value":   eth.WeiToEthInt(deposit.Value).Int64(),
			"epoch":   deposit.Epoch,
		}).Trace("applyEpochDeposit: Caching deposit proof.")
		e.depositProofCache[deposit.Account] = depProof
	}
	return nil
}

// exitersSet is the set of exiting participants.
type exitersSet []common.Address

// applyEpochExit handles the exit phase of given Epoch.
func (e *Enclave) applyEpochExit(ep *Epoch, exLogs []*types.Log) (exitersSet, error) {
	var exiters exitersSet
	for _, exLog := range exLogs {
		exit, err := parseExitEvent(exLog)
		if err != nil {
			return nil, fmt.Errorf("parsing exiting event: %w", err)
		}
		if exit.Epoch != ep.Number {
			return nil, fmt.Errorf("exit-event Epoch %v != current exit Epoch %v",
				exit.Epoch, ep.Number)
		}
		// only full exits supported currently
		accBal := ep.balances[exit.Account].Value
		if accBal.Cmp(exit.Value) != 0 {
			return nil, fmt.Errorf("unexpected partial exit for %v, expected %v", exit.Value, accBal)
		}
		accBal.SetUint64(0)
		delete(ep.balances, exit.Account)
		exiters = append(exiters, exit.Account)
	}

	return exiters, nil
}

// applyEpochTx adjusts `e.balances` according to given transactions.
func (e *Enclave) applyEpochTx(ep *Epoch, tx *tee.Transaction) error {
	const LT = -1

	if valid, err := tee.VerifyTransaction(e.params.Contract, *tx); err != nil {
		return fmt.Errorf("verifying tx signature: %w", err)
	} else if !valid {
		return fmt.Errorf("invalid tx signature")
	}

	if tx.Epoch != ep.Number {
		log.Errorf("wrong TX Epoch: expected %d, got %d", ep.Number, tx.Epoch)
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
		return fmt.Errorf("sender tx nonce: %v, expected %v",
			tx.Nonce, sender.Nonce+1)
	}

	sender.Value.Sub(sender.Value, tx.Amount)
	recipient.Value.Add(recipient.Value, tx.Amount)

	sender.Nonce = tx.Nonce

	return nil
}
