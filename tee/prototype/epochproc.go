// SPDX-License-Identifier: Apache-2.0

package prototype

import "github.com/perun-network/erdstall/tee"

func (e *Enclave) epochProcessor(
	verifiedBlocks <-chan *tee.Block,
	txs <-chan *tee.Transaction,
) error {
	var (
		depositEpoch *Epoch
		txEpoch      *Epoch
		exitEpoch    *Epoch
	)

	shiftEpochs := func() {
		// close channels (exit of n+2, tx of n+1, deposit of n)
		exitEpoch = txEpoch
		txEpoch = depositEpoch
		depositEpoch = txEpoch.NewNext()
		// write new epoch to e.epochs
	}

	// read blocks from verifiedBlocks
	// read tx from txs
	// extract deposits, adjust epoch's balances
	// extract exits, adjust epoch's balances
	// responsible for closing end-of-phase signalling channels
}
