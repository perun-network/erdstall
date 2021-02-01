// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/tee"
)

type (
	// Epoch manages deposits, transactions, and exits, as well as the
	// progression of time in the system.
	Epoch struct {
		number tee.Epoch // Current deposit epoch.

		exitLocked map[common.Address]struct{} // Withdrawing at end of epoch.
		exitReqs   map[common.Address]struct{} // Requested exits.
		accs       map[common.Address]*Acc     // Latest account states.

		outcome chan Outcome // The last epoch's outcome.
	}

	// Acc contains an account's balance and transaction nonce.
	Acc struct {
		Nonce uint64
		Value *big.Int
	}

	// Outcome contains all of an epoch's final balances, as well as all
	// accounts that exited the system at the end of the epoch.
	Outcome struct {
		TxEpoch  tee.Epoch               // The finished TX epoch.
		Exits    map[common.Address]*Acc // Exited accounts.
		Accounts map[common.Address]*Acc // Final balances and nonces.
	}
)

func newEpoch(n tee.Epoch) *Epoch {
	return &Epoch{
		number:     n,
		exitLocked: make(map[common.Address]struct{}),
		exitReqs:   make(map[common.Address]struct{}),
		accs:       make(map[common.Address]*Acc),
		outcome:    make(chan Outcome, 1),
	}
}

// DepositNum is the current deposit epoch number.
func (e *Epoch) DepositNum() tee.Epoch { return e.number }

// TxNum is the current transaction epoch number.
func (e *Epoch) TxNum() tee.Epoch { return e.number - 1 }

// ExitNum is the current exit epoch number.
func (e *Epoch) ExitNum() tee.Epoch { return e.number - 2 }

// Outcome has to be called exactly once after each phase shift and contains all
// accounts that exited within this epoch, as well as all balances at the end of
// the epoch. It blocks until the epoch in which it was last called ends.
func (e *Epoch) Outcome() (o Outcome) {
	select {
	case o = <-e.outcome:
	default:
		log.Panic("Requested outcome twice or before end of epoch")
	}
	return // Go doesn't recognize that log.Panic will not return.
}

// RegisterExits registers a exit requests for the end of the current phase. The
// requested accounts are locked for one epoch before they can be withdrawn.
func (e *Epoch) RegisterExits(exitReqs ...*erdstallExitEvent) {
	for _, exit := range exitReqs {
		if exit.Epoch != e.ExitNum() {
			log.WithFields(log.Fields{
				"req. epoch": exit.Epoch, "exit epoch": e.ExitNum(),
			}).Panic("epoch mismatch")
		}
		e.exitReqs[exit.Account] = struct{}{}
	}
}

// ApplyDeposits applies a series deposit events and makes the deposited values
// available instantly.
func (e *Epoch) ApplyDeposits(deposits ...*erdstallDepEvent) {
	for _, dep := range deposits {
		e.receive(dep.Account, dep.Value)
	}
}

// receive increases an account's balance, or creates a new account if it didn't
// exist yet.
func (e *Epoch) receive(account common.Address, value *big.Int) {
	if acc, ok := e.accs[account]; ok {
		acc.Value.Add(acc.Value, value)
	} else {
		e.accs[account] = &Acc{
			Nonce: 0,
			Value: new(big.Int).Set(value),
		}
	}
}

// ProcessTx processes a transaction. If it is invalid, does nothing and returns
// an error. A transaction is invalid if either party is currently locked for
// withdrawal or the transaction is otherwise invalid (e.g., insufficient funds,
// invalid signature, ...). Valid transactions increase the sender's nonce.
func (e *Epoch) ProcessTx(contract common.Address, tx *tee.Transaction) error {
	sender, ok := e.accs[tx.Sender]
	if !ok {
		return errors.New("sender does not exist")
	}

	// Flush the tx.hash so that it is recalculated.
	*tx = tee.Transaction{
		Nonce: tx.Nonce, Epoch: tx.Epoch,
		Sender: tx.Sender, Recipient: tx.Recipient,
		Amount: tx.Amount, Sig: tx.Sig}

	// Check whether both participants are eligible for trading and that the
	// transaction is valid.
	if _, locked := e.exitLocked[tx.Sender]; locked {
		return errors.New("sender is locked for withdrawing")
	} else if _, locked := e.exitLocked[tx.Recipient]; locked {
		return errors.New("recipient is locked for withdrawing")
	} else if tx.Nonce != sender.Nonce+1 {
		return fmt.Errorf("nonce mismatch: %d != %d", tx.Nonce, sender.Nonce+1)
	} else if sender.Value.Cmp((*big.Int)(tx.Amount)) < 0 {
		return errors.New("insufficient balance")
	} else if (*big.Int)(tx.Amount).Sign() < 0 {
		return errors.New("negative amount")
	} else if tx.Epoch != e.TxNum() {
		return fmt.Errorf("epoch mismatch: %d != %d", tx.Epoch, e.TxNum())
	} else if valid, err := tee.VerifyTransaction(contract, *tx); err != nil {
		return fmt.Errorf("verifying tx signature: %w", err)
	} else if !valid {
		return fmt.Errorf("invalid tx signature")
	}

	// Execute the transaction.
	sender.Value.Sub(sender.Value, (*big.Int)(tx.Amount))
	sender.Nonce = tx.Nonce
	e.receive(tx.Recipient, (*big.Int)(tx.Amount))

	return nil
}

// progressPhase transitions from one epoch to the next and resets all
// phase-specific fields of the previous phase. It also applies any ongoing
// exits.
func (e *Epoch) progressPhase() {
	// Settle current epoch.
	exits := e.applyExits()
	bals := e.cloneBals()

	// Publish the epoch's outcome.
	e.outcome <- Outcome{TxEpoch: e.TxNum(), Exits: exits, Accounts: bals}
	e.number++
}

// applyExits deletes all accounts that exited during the previous epoch and
// returns their final balances. Exit requests for nonexisting accounts are
// ignored.
func (e *Epoch) applyExits() map[common.Address]*Acc {
	exits := make(map[common.Address]*Acc)
	// Delete and collect all accounts that were locked for withdrawing during
	// this epoch.
	for addr := range e.exitLocked {
		if bal, ok := e.accs[addr]; ok {
			delete(e.accs, addr)
			exits[addr] = bal
		}
	}
	// Lock all accounts that requested to withdraw within the next epoch.
	e.exitLocked, e.exitReqs = e.exitReqs, make(map[common.Address]struct{})
	return exits
}

func (e *Epoch) cloneBals() map[common.Address]*Acc {
	accs := make(map[common.Address]*Acc)
	for addr, acc := range e.accs {
		accs[addr] = &Acc{
			Nonce: acc.Nonce,
			Value: new(big.Int).Set(acc.Value),
		}
	}
	return accs
}

// IsExitLocked returns whether an account is exit locked.
func (e *Epoch) IsExitLocked(who common.Address) bool {
	_, ok := e.exitLocked[who]
	return ok
}

// Balance looks up an account balance in the epoch.
func (e *Epoch) Balance(who common.Address) *big.Int {
	if acc, ok := e.accs[who]; ok {
		return new(big.Int).Set(acc.Value)
	} else {
		return big.NewInt(0)
	}
}
