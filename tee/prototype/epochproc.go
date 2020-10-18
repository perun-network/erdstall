// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	err "perun.network/go-perun/pkg/errors"

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
		done := make(chan exitersSet)

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

// TODO: racing transactions.
func (e *Enclave) txRoutine(
	phaseShift chan<- struct{},
	exits <-chan exitersSet,
	txs <-chan *tee.Transaction,
	txEpoch *Epoch,
) error {
	var stagedTxs txCache
	for {
		bpCache, err := e.generateBalanceProofs(txEpoch)
		if err != nil {
			return fmt.Errorf("generating balance proofs: %w", err)
		}
		select {
		case exiters := <-exits:
			txEpoch, bpCache, ptxs := prune(txEpoch, stagedTxs, bpCache, exiters)
			bpCache, err = e.applyEpochTx(txEpoch, ptxs, bpCache)
			if err != nil {
				return fmt.Errorf("adjusting Epoch %v Balances: %w", txEpoch.Number, err)
			}
			e.balanceProofs <- getBalProofs(bpCache)
			phaseShift <- struct{}{} // phase done, stop processing TXs and signal epoch shift.
			return nil
		case tx := <-txs:
			stagedTxs = cacheTx(stagedTxs, tx)
		}
	}
}

// getBalProofs retrieves a slice of `*tee.BalanceProof`s from the balance proof cache.
func getBalProofs(bps balProofCache) []*tee.BalanceProof {
	proofs := make([]*tee.BalanceProof, len(bps))
	i := 0
	for _, bp := range bps {
		proofs[i] = bp
	}
	return proofs
}

// prune prunes a given slice of transactions for a given set of addresses, s.t.
// no element of the given set is either a recipient or sender in a transaction,
// it also makes sure to remove exited parties from the txEpoch and the
// balance proof cache.
func prune(epoch *Epoch,
	txs txCache,
	bpc balProofCache,
	forS []common.Address,
) (*Epoch, balProofCache, []*tee.Transaction) {
	for _, m := range forS {
		delete(epoch.balances, m) // remove member from txEpoch
		delete(bpc, m)            // remove member from balProofCache

		sTxs, sOk := txs.senders[m]
		rTxs, rOk := txs.recipients[m]
		if !sOk && !rOk {
			continue
		}
		for _, tx := range sTxs {
			delete(txs.txs, tx.Hash())
		}
		for _, tx := range rTxs {
			delete(txs.txs, tx.Hash())
		}
	}
	prunedTxs := make([]*tee.Transaction, 0, len(txs.txs))
	i := 0
	for _, tx := range txs.txs {
		prunedTxs[i] = tx
	}
	return epoch, bpc, prunedTxs
}

type balProofCache = map[common.Address]*tee.BalanceProof

func (e *Enclave) depositExitRoutine(
	verifiedBlocks <-chan *tee.Block,
	phaseDone chan<- exitersSet,
	depositEpoch, exitEpoch *Epoch,
) error {
	// read blocks from verifiedBlocks (deposit phase).
	for vb := range verifiedBlocks {
		exiters, err := e.handleVerifiedBlock(depositEpoch, exitEpoch, vb)
		if err != nil {
			return fmt.Errorf("handling verified blocknr %v : %w", vb.NumberU64(), err)
		}
		if e.phaseDone(vb.NumberU64()) {
			e.depositProofs <- retrieveCachedDepProofs(e.depositProofCache)
			e.depositProofCache = make(map[common.Address]*tee.DepositProof)
			phaseDone <- exiters
			close(phaseDone) // phase done, stop TX processor.
			return nil
		}
	}
	return nil
}

// retrieveCachedDepProofs reduces the deposit proof cache to a slice of `tee.DepositProof`s.
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

// updateBalanceProofs updates the balance proof cache with for given transaction.
func (e *Enclave) updateBalanceProofs(bpc balProofCache, ep *Epoch, sender, recipient common.Address) (balProofCache, error) {
	sBp := bpc[sender]
	rBp := bpc[recipient]
	sBp.Balance.Value = sBp.Balance.Value.Set(ep.balances[sender].Value)
	rBp.Balance.Value = rBp.Balance.Value.Set(ep.balances[recipient].Value)

	sBp, err := e.signBalanceProof(sBp.Balance)
	if err != nil {
		return nil, fmt.Errorf("updating balance proof for %v", sender.String())
	}
	rBp, err = e.signBalanceProof(rBp.Balance)
	if err != nil {
		return nil, fmt.Errorf("updating balance proof for %v", recipient.String())
	}

	bpc[sender] = sBp
	bpc[recipient] = rBp

	return bpc, nil
}

// generateBalanceProofs generates the balance proofs for all users in the given
// transaction Epoch.
func (e *Enclave) generateBalanceProofs(txEpoch *Epoch) (balProofCache, error) {
	balProofs := make(balProofCache)
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
		balProofs[acc] = bp
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
func (e *Enclave) handleVerifiedBlock(depEpoch, exitEpoch *Epoch, vb *tee.Block) (exitersSet, error) {
	var exiters exitersSet
	for _, r := range vb.Receipts {
		exLogs, depLogs := partition(r.Logs, logIsDepositEvt)
		// extract deposits, adjust epoch's balances.
		if err := e.applyEpochDeposit(e.params.Contract, depEpoch, depLogs); err != nil {
			return nil, fmt.Errorf("adjusting Epoch %v Balances: %w", depEpoch.Number, err)
		}
		// extract exits, adjust epoch's balances.
		exits, err := e.applyEpochExit(exitEpoch, exLogs)
		if err != nil {
			return nil, fmt.Errorf("handling exit of Epoch %v: %w", exitEpoch.Number, err)
		}
		exiters = append(exiters, exits...)
	}
	return exiters, nil
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

// applyEpochDeposit adjusts `e.balances` according to the deposits done in
// the given block.
func (e *Enclave) applyEpochDeposit(contract common.Address, ep *Epoch, depLogs []*types.Log) error {
	for _, depLog := range depLogs {
		deposit, err := parseDepEvent(depLog)
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

// exitersSet is the set of exiting participants.
type exitersSet []common.Address

// applyEpochExit handles the exit phase of given Epoch.
func (e *Enclave) applyEpochExit(ep *Epoch, exLogs []*types.Log) (exitersSet, error) {
	exiters := make(exitersSet, 0)
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
func (e *Enclave) applyEpochTx(ep *Epoch, txs []*tee.Transaction, bpc balProofCache) (balProofCache, error) {
	const LT = -1

	for _, tx := range txs {
		if valid, err := tee.VerifyTransaction(e.params, *tx); err != nil {
			return nil, fmt.Errorf("verifying tx signature: %w", err)
		} else if !valid {
			return nil, fmt.Errorf("invalid tx signature")
		}

		if tx.Epoch != ep.Number {
			// TODO: log that we drop
			return nil, nil
		}

		sender, oks := ep.balances[tx.Sender]
		recipient, okr := ep.balances[tx.Recipient]
		if !oks {
			return nil, fmt.Errorf("unknown sender: %x", tx.Sender)
		}
		if !okr {
			return nil, fmt.Errorf("unknown recipient: %x", tx.Recipient)
		}
		if sender.Value.Cmp(tx.Amount) == LT {
			return nil, fmt.Errorf("tx amount exceeds senders balance: has: %v, needs: %v", sender.Value, tx.Amount)
		}
		if tx.Nonce != sender.Nonce+1 {
			return nil, fmt.Errorf("comparing tx nonce: %v, sender nonce: %v",
				tx.Nonce, sender.Nonce)
		}

		sender.Value.Sub(sender.Value, tx.Amount)
		recipient.Value.Add(recipient.Value, tx.Amount)

		sender.Nonce = tx.Nonce

		e.updateBalanceProofs(bpc, ep, tx.Sender, tx.Recipient)
	}
	return nil, nil
}
