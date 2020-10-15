// SPDX-License-Identifier: Apache-2.0

package tee

import "github.com/ethereum/go-ethereum/common"

type Parameters struct {
	PowDepth      uint64         // required confirmed block depth
	PhaseDuration uint64         // number of blocks of one phase (not epoch)
	InitBlock     uint64         // block at which Erdstall contract was deployed
	Contract      common.Address // Erdstall contract address
}

// DepositEpoch returns the deposit epoch at the given block number.
func (p Parameters) DepositEpoch(blockNum uint64) Epoch {
	return p.epoch(blockNum)
}

// TxEpoch returns the transaction epoch at the given block number.
func (p Parameters) TxEpoch(params Parameters, blockNum uint64) Epoch {
	return p.epoch(blockNum) + 1
}

// ExitEpoch returns the exit epoch at the given block number.
func (p Parameters) ExitEpoch(params Parameters, blockNum uint64) Epoch {
	return p.epoch(blockNum) + 2
}

// FreezingEpoch returns the freezing epoch at the given block number.
func (p Parameters) FreezingEpoch(params Parameters, blockNum uint64) Epoch {
	return p.epoch(blockNum) + 3
}

// Don't use this, use the specific FooEpoch methods.
func (p Parameters) epoch(blockNum uint64) Epoch {
	return (blockNum - p.InitBlock) / p.PhaseDuration
}
