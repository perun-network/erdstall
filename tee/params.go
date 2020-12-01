// SPDX-License-Identifier: Apache-2.0

package tee

import "github.com/ethereum/go-ethereum/common"

type Parameters struct {
	PowDepth         uint64         // required confirmed block depth
	PhaseDuration    uint64         // number of blocks of one phase (not epoch)
	ResponseDuration uint64         // challenge response grace period for operator at end of exit phase
	InitBlock        uint64         // block at which Erdstall contract was deployed
	TEE              common.Address // Enclave's public key address
	Contract         common.Address // Erdstall contract address
}

// DepositEpoch returns the deposit epoch at the given block number.
func (p Parameters) DepositEpoch(blockNum uint64) Epoch {
	return p.epoch(blockNum)
}

// TxEpoch returns the transaction epoch at the given block number.
func (p Parameters) TxEpoch(blockNum uint64) Epoch {
	return p.epoch(blockNum) - 1
}

// ExitEpoch returns the exit epoch at the given block number.
func (p Parameters) ExitEpoch(blockNum uint64) Epoch {
	return p.epoch(blockNum) - 2
}

// FreezingEpoch returns the freezing epoch at the given block number.
func (p Parameters) FreezingEpoch(blockNum uint64) Epoch {
	return p.epoch(blockNum) - 3
}

// Don't use this, use the specific FooEpoch methods.
func (p Parameters) epoch(blockNum uint64) Epoch {
	return (blockNum - p.InitBlock) / p.PhaseDuration
}

func (p Parameters) IsChallengeResponsePhase(blockNum uint64) bool {
	return p.PhaseDuration-((blockNum-p.InitBlock)%p.PhaseDuration) <= p.ResponseDuration
}

// IsLastPhaseBlock tells whether this block is the last block of a phase.
func (p Parameters) IsLastPhaseBlock(blockNum uint64) bool {
	return ((blockNum - p.InitBlock) % p.PhaseDuration) == p.PhaseDuration-1
}

func (p Parameters) DepositStartBlock(epoch uint64) uint64 {
	return p.InitBlock + epoch*p.PhaseDuration
}

func (p Parameters) DepositDoneBlock(epoch uint64) uint64 {
	return p.DepositStartBlock(epoch) + p.PhaseDuration
}

func (p Parameters) TxDoneBlock(epoch uint64) uint64 {
	return p.DepositStartBlock(epoch) + 2*p.PhaseDuration
}

func (p Parameters) ExitDoneBlock(epoch uint64) uint64 {
	return p.DepositStartBlock(epoch) + 3*p.PhaseDuration
}
