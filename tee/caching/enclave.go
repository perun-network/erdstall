// SPDX-License-Identifier: Apache-2.0

package caching

import (
	"math"
	"sync"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/tee"
)

type (
	Enclave struct {
		tee.Enclave

		params     tee.Parameters
		txEpoch    tee.Epoch
		epochMu    sync.Mutex
		pendingTxs map[tee.Epoch][]txReq
	}

	txReq struct {
		txs []*tee.Transaction
		res chan error
	}
)

func NewEnclave(e tee.Enclave) *Enclave {
	return &Enclave{
		Enclave:    e,
		txEpoch:    math.MaxUint64, // -1 during deposit epoch 0
		pendingTxs: make(map[tee.Epoch][]txReq),
	}
}

func (e *Enclave) Run(p tee.Parameters) error {
	e.params = p
	return e.Enclave.Run(p)
}

func (e *Enclave) ProcessBlocks(bs ...*tee.Block) error {
	log.Tracef("CachingEnclave: ProcessBlocks called with %d blocks", len(bs))
	e.epochMu.Lock()
	defer e.epochMu.Unlock()
	log.Trace("CachingEnclave: epoch lock acquired")

	for _, b := range bs {
		if err := e.Enclave.ProcessBlocks(b); err != nil {
			return err
		}

		txEpoch := e.params.TxEpoch(b.NumberU64() + 1) // active block is +1
		log := log.WithField("txEpoch", txEpoch)
		if txEpoch == e.txEpoch {
			log.Debug("CachingEnclave: same epoch")
			continue
		} else if txEpoch != e.txEpoch+1 {
			log.Panicf("tx epoch not incremented by 1, %d->%d", e.txEpoch, txEpoch)
		}
		log.Debug("CachingEnclave: tx epoch shifted")
		e.txEpoch = txEpoch

		// process cached transactions on epoch shift
		txReqs, ok := e.pendingTxs[txEpoch] // loop over nil ok
		if ok {
			delete(e.pendingTxs, txEpoch)
		}
		for _, txReq := range txReqs {
			log.Debugf("CachingEnclave: processing %d cached txs", len(txReq.txs))
			if err := e.Enclave.ProcessTXs(txReq.txs...); err != nil {
				txReq.res <- errors.WithMessagef(err, "processing cached txs of epoch %d", txEpoch)
			}
		}
	}

	return nil
}

func (e *Enclave) ProcessTXs(txs ...*tee.Transaction) error {
	var ep tee.Epoch
	for i, tx := range txs {
		if i == 0 {
			ep = tx.Epoch
			continue
		}
		if tx.Epoch != ep {
			return errors.Errorf(
				"different epochs in transaction batch, [0]: %d, [%d]: %d", ep, i, tx.Epoch)
		}
	}

	e.epochMu.Lock()

	switch {
	case ep+1 < e.txEpoch+1: // because of initial tx epoch uint64(-1)
		defer e.epochMu.Unlock()
		return errors.Errorf("TXs epoch %d already sealed, current is %d", ep, e.txEpoch)
	case ep == e.txEpoch:
		defer e.epochMu.Unlock()
		return e.Enclave.ProcessTXs(txs...)
	}

	// Case ep > e.txEpoch
	log.Debugf("Caching txs for future epoch %d", ep)
	res := make(chan error)
	e.pendingTxs[ep] = append(e.pendingTxs[ep], txReq{txs, res})
	e.epochMu.Unlock()

	// Now waiting for ProcessBlocks to receive block that will shift epoch
	// and process all pending TXs.
	return <-res
}
