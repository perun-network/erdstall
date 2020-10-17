// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"fmt"

	"github.com/perun-network/erdstall/tee"
)

func (e *Enclave) blockProcessor(
	newBlocks <-chan *tee.Block,
	verifiedBlocks chan<- *tee.Block,
) error {
	// read blocks from newBlocks
	k := e.params.PowDepth
	for b := range newBlocks {
		// verify block chain
		if err := e.bc.PushVerify(b); err != nil {
			return fmt.Errorf("pushing block to local blockchain: %w", err)
		}

		// TODO: * Handle Reorgs...
		//   * Receipts verification
		//   * First block verification for `params`

		// write verified block (up to PoW-depth) to verifiedBlocks
		l := e.bc.Len()
		if l > k {
			verifiedBlocks <- e.bc.blocks[l-k]
		}
	}
	return nil
}
