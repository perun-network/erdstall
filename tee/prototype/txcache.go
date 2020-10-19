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

// cacheTx caches the given transaction.
func (txc *txCache) cacheTx(tx *tee.Transaction) *txCache {
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
	sTxs = append(sTxs, tx)
	rTxs = append(rTxs, tx)

	return txc
}
