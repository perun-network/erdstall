// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/perun-network/erdstall/tee"
)

type txCache struct {
	senders    map[common.Address][]*tee.Transaction
	recipients map[common.Address][]*tee.Transaction
}

func makeTxCache() txCache {
	return txCache{
		senders:    make(map[common.Address][]*tee.Transaction),
		recipients: make(map[common.Address][]*tee.Transaction),
	}
}

// cacheTx caches the given transaction.
func (txc *txCache) cacheTx(tx *tee.Transaction) {
	sender := tx.Sender
	recipient := tx.Recipient

	if _, sOk := txc.senders[sender]; !sOk {
		txc.senders[sender] = make([]*tee.Transaction, 0)
	}
	if _, rOk := txc.recipients[recipient]; !rOk {
		txc.recipients[recipient] = make([]*tee.Transaction, 0)
	}
	sTxs := txc.senders[sender]
	rTxs := txc.recipients[recipient]
	txc.senders[sender] = append(sTxs, tx)
	txc.recipients[recipient] = append(rTxs, tx)
}
