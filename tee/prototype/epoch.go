// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/tee"
)

// An Epoch stores all balances of an epoch and contains signalling
// infrastructure.
type (
	Epoch struct {
		Number tee.Epoch

		balances map[common.Address]*Bal
	}

	Bal struct {
		Nonce uint64
		Value *big.Int
	}
)

func NewEpoch(n tee.Epoch) *Epoch {
	return &Epoch{
		Number:   n,
		balances: make(map[common.Address]*Bal),
	}
}

// Next returns a clone of the current Epoch, with the Epoch counter increased
// by one.
func (e *Epoch) NewNext() *Epoch {
	next := NewEpoch(e.Number + 1)
	return next
}

func cloneBalances(a map[common.Address]*Bal) map[common.Address]*Bal {
	b := make(map[common.Address]*Bal, len(a))
	for k, v := range a {
		b[k] = &Bal{
			Nonce: v.Nonce,
			Value: new(big.Int).Set(v.Value),
		}
	}
	return b
}

func (e *Epoch) applyExits(exiters exitersSet) {
	if e == nil {
		return
	}
	log := log.WithField("epoch", e.Number)
	log.Debugf("removing %d exiters", len(exiters))
	for _, exiter := range exiters {
		log.Tracef("removed %s", exiter.String())
		delete(e.balances, exiter)
	}
}

// merge merges `e` with `p` and returns a clone.
func (e *Epoch) merge(p *Epoch) *Epoch {
	mEpoch := NewEpoch(e.Number)
	if p == nil {
		mEpoch.balances = cloneBalances(e.balances)
		return mEpoch
	}

	mEpoch.balances = cloneBalances(p.balances)
	for acc, bal := range e.balances {
		if pBal, ok := mEpoch.balances[acc]; ok {
			bal.Value.Add(bal.Value, pBal.Value)
		} else {
			mEpoch.balances[acc] = bal
		}
	}
	return mEpoch
}
