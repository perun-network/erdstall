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
			return fmt.Errorf("invalid block: %v", err)
		}

		// TODO: Handle Reorg...
		// TODO: * Receipts verification? Maybe already done in `blockProcessor`.
		//		 * First block verification for `params`.

		// write verified block (up to PoW-depth) to verifiedBlocks
		l := uint64(len(e.bc.blocks))
		if l > k {
			verifiedBlocks <- e.bc.blocks[l-k]
		}
	}
	return nil
}

// verifyBlock verifies a given block by checking if it's hashes are consistent
// with the blockchain.
func verifyBlock(block, parent *tee.Block) error {
	bHash := block.Header().ParentHash.String()
	pHash := parent.Hash().String()

	// TODO: Extend validation.
	if bHash != pHash {
		return fmt.Errorf("comparing header hashes, expected %v, got: %v", pHash, bHash)
	}
	return nil
}
