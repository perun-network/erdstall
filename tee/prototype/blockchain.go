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

// Block returns the block with the given block number. Panics if not present.
func (b *blockchain) Block(blockNum uint64) *tee.Block {
	return b.blocks[blockNum-b.offset]
}

func (b *blockchain) Len() uint64 {
	return uint64(len(b.blocks))
}

// PushVerify pushes the block onto the chain, verifying that it is indeed a
// correct next block.
//
// If the chain's current head has the same or a larger block number, a chain
// reorg is assumed and the current block becomes the new head, discarding all
// previous blocks.
func (b *blockchain) PushVerify(block *tee.Block) error {
	blockNum := block.NumberU64()
	if len(b.blocks) == 0 {
		// first block
		b.offset = blockNum
		b.blocks = []*tee.Block{block}
		return nil
	}

	prev := b.Head()
	prevNum := prev.NumberU64()
	if blockNum > prevNum+1 {
		return fmt.Errorf("intermediate blocks missing, head: %d, block: %d", prevNum, blockNum)
	} else if blockNum <= prevNum {
		// reorg
		prev = b.Block(blockNum - 1)
		prevNum = prev.NumberU64()
	}

	if err := verifyBlock(block, prev); err != nil {
		return fmt.Errorf("verifying block: %v", err)
	}

	b.blocks = append(b.blocks[:prevNum+1-b.offset], block)
	return nil
}

// PruneUntil can be used to discard all blocks until the given block number
func (b *blockchain) PruneUntil(blockNum uint64) {
	if b.offset >= blockNum {
		return
	}

	diff := blockNum - b.offset
	b.blocks = b.blocks[diff:len(b.blocks)]
	b.offset = blockNum
}

// verifyBlock verifies that block is a valid next block after parent.
func verifyBlock(block, parent *tee.Block) error {
	bHash := block.Header().ParentHash
	pHash := parent.Hash()

	// TODO: Extend validation to test consensus (PoW).
	//   This function then probably becomes a method on blockchain after we add a
	//   consensus engine to blockchain.
	if bHash != pHash {
		return fmt.Errorf("parent header mismatch, expected %x, got: %x", pHash, bHash)
	}
	return nil
}
