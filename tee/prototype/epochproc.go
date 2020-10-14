// SPDX-License-Identifier: Apache-2.0

package prototype

func (e *Enclave) epochProcessor() error {
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

	// read blocks from e.verifiedBlocks
	// read tx from e.txs
	// extract deposits, adjust epoch's balances
	// extract exits, adjust epoch's balances
	// responsible for closing phase channels
}
