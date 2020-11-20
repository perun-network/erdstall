// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"fmt"

	"github.com/perun-network/erdstall/tee"
)

type blockchain struct {
	offset uint64 // block.number of first block in blocks slice
	blocks []*tee.Block
}

func (b *blockchain) Head() *tee.Block {
	if len(b.blocks) == 0 {
		return nil
	}
	return b.blocks[len(b.blocks)-1]
}

func (b *blockchain) Len() uint64 {
	return uint64(len(b.blocks))
}

// PushVerify pushes the block onto the chain, verifying that it is indeed a
// correct next block.
func (b *blockchain) PushVerify(block *tee.Block) error {
	blockNum := block.NumberU64()
	if len(b.blocks) == 0 {
		// first block
		b.offset = blockNum
		b.blocks = []*tee.Block{block}
		return nil
	}

	headNum := b.Head().NumberU64()
	if headNum+1 != blockNum {
		return fmt.Errorf("not next block, head: %d, block: %d", headNum, blockNum)
	}

	if err := verifyBlock(block, b.Head()); err != nil {
		return fmt.Errorf("verifying block: %v", err)
	}

	b.blocks = append(b.blocks, block)
	return nil
}

// PruneUntil can be used to discard all blocks until the given block number
func (b *blockchain) PruneUntil(blockNum uint64) {
	if b.offset >= blockNum {
		return
	}

	diff := blockNum - b.offset
	b.blocks = b.blocks[diff : len(b.blocks)-1]
	b.offset = blockNum
}

// verifyBlock verifies that block is a valid next block after parent.
func verifyBlock(block, parent *tee.Block) error {
	bHash := block.Block.Header().ParentHash
	pHash := parent.Block.Hash()

	// TODO: Extend validation to test consensus (PoW).
	//   This function then probably becomes a method on blockchain after we add a
	//   consensus engine to blockchain.
	if bHash != pHash {
		return fmt.Errorf("parent header mismatch, expected %x, got: %x", pHash, bHash)
	}
	return nil
}
