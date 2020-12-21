// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (e *Enclave) blockProcessor(
	newBlocks <-chan blockReq,
	verifiedBlocks chan<- blockReq,
) error {
	log.Debug("blockProc: starting...")

	var vn uint64 // last verified block number

	process := func(b blockReq) (error, bool) {
		k := uint64(0) //e.params.PowDepth

		if n := b.block.NumberU64(); e.bc.Len() == 0 && e.params.InitBlock != n {
			return fmt.Errorf("first block (%d) not initial Erdstall block (%d)", n, e.params.InitBlock), true
		}

		// verify blockchain
		if err := e.bc.PushVerify(b.block); err != nil {
			return fmt.Errorf("pushing block to local blockchain: %w", err), true
		}

		// TODO: Receipts verification

		// Now write verified block (up to PoW-depth) to verifiedBlocks

		l := e.bc.Len()
		if l <= k {
			log.Trace("blockProc: less than k blocks")
			return nil, false
		}

		vblock := e.bc.blocks[l-k-1]
		log := log.WithField("blockNum", vblock.NumberU64())
		if vn >= vblock.NumberU64() {
			log.Trace("blockProc: skipping forwarding reorg block to epochProc")
			return nil, false
		} else if vn != 0 && vn+1 != vblock.NumberU64() {
			log.Panic("blockProc: next verified block should increment blockNum by 1")
		}
		vn = vblock.NumberU64()

		log.Trace("blockProc: forwarding verified block to epochProc")
		verifiedBlocks <- blockReq{block: vblock, result: b.result}
		if e.params.IsLastPhaseBlock(vn) && !e.running.IsSet() {
			// graceful shutdown
			log.Info("blockProc: last verified phase block forwarded, shutting down")
			return nil, true
		}

		return nil, false
	}

	for b := range newBlocks {
		err, shutdown := process(b)
		b.result <- err
		if err != nil || shutdown {
			return err
		}
	}

	return errors.New("blockProc: newBlocks channel closed")
}
