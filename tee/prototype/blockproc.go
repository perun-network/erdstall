// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/tee"
)

func (e *Enclave) blockProcessor(
	newBlocks <-chan *tee.Block,
	verifiedBlocks chan<- *tee.Block,
) error {
	k := e.params.PowDepth
	log.Debug("blockProc: starting...")
	for {
		select {
		case b := <-newBlocks:
			if e.bc.Len() == 0 && e.params.InitBlock != b.NumberU64() {
				return fmt.Errorf("first block (%d) not initial Erdstall block (%d)", b.NumberU64(), e.params.InitBlock)
			}

			// verify block chain
			if err := e.bc.PushVerify(b); err != nil {
				return fmt.Errorf("pushing block to local blockchain: %w", err)
			}

			// TODO: * Handle Reorgs...
			//   * Receipts verification

			// write verified block (up to PoW-depth) to verifiedBlocks
			l := e.bc.Len()
			if l > k {
				vblock := e.bc.blocks[l-k-1]
				log.WithField("blockNum", vblock.NumberU64()).Trace("blockProc: forwarding block to epochProc")
				verifiedBlocks <- vblock
			}
		case <-e.quit:
			return nil
		}
	}
}
