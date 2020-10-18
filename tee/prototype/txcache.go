// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/tee"
)

type txCache struct {
	senders    map[common.Address][]*tee.Transaction
	recipients map[common.Address][]*tee.Transaction
	txs        map[common.Hash]*tee.Transaction
}

// TODO: Pointer receiver?
// cacheTx caches the given transaction.
func cacheTx(cache txCache, tx *tee.Transaction) txCache {
	sender := tx.Sender
	recipient := tx.Recipient

	if _, sOk := cache.senders[sender]; !sOk {
		cache.senders[sender] = make([]*tee.Transaction, 0)
	}
	if _, rOk := cache.recipients[recipient]; !rOk {
		cache.recipients[recipient] = make([]*tee.Transaction, 0)
	}
	sTxs := cache.senders[sender]
	rTxs := cache.recipients[recipient]
	sTxs = append(sTxs, tx)
	rTxs = append(rTxs, tx)

	cache.txs[tx.Hash()] = tx
	return cache
}
