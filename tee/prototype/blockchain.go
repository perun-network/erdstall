// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/perun-network/erdstall/tee"
)

// Blockchain keeps track of the blockchain and ensures the integrity of new
// blocks entered in to the enclave.
type blockchain struct {
	head *tee.Block
}

// Head returns the latest verified and processed block or nil.
func (b *blockchain) Head() *tee.Block { return b.head }
func (b *blockchain) empty() bool      { return b.head == nil }

// PushVerify pushes the block onto the chain, verifying that it is indeed a
// valid successor block of the previous head. If the block is valid, returns
// all of the block's deposit and exit events. If the block is invalid, the
// chain remains unchanged.
func (b *blockchain) PushVerify(
	block *tee.Block,
	params *tee.Parameters,
	epoch *Epoch,
) ([]*erdstallDepEvent, []*erdstallExitEvent, error) {
	if !b.empty() {
		if err := verifySuccessorBlock(block, b.head); err != nil {
			return nil, nil, fmt.Errorf("verifying successor block: %v", err)
		}
	}

	deps, exits, err := extractEvents(block, params)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid block events: %w", err)
	}

	blockN := block.NumberU64()
	depEpoch, exitEpoch := params.DepositEpoch(blockN), params.ExitEpoch(blockN)

	if err = verifyDeposits(deps, depEpoch); err != nil {
		return nil, nil, fmt.Errorf("invalid deposits: %w", err)
	} else if err = verifyExits(exits, exitEpoch, epoch); err != nil {
		return nil, nil, fmt.Errorf("invalid exits: %w", err)
	}

	b.head = block
	return deps, exits, err
}

// verifySuccessorBlock verifies that next is a valid successor of head.
func verifySuccessorBlock(next, head *tee.Block) error {
	nextHash, headHash := next.Header().ParentHash, head.Hash()
	nextN, headN := next.NumberU64(), head.NumberU64()

	// TODO: Extend validation to test consensus (PoW).
	//   This function then probably becomes a method on nextchain after we add a
	//   consensus engine to nextchain.
	switch {
	case nextN != headN+1:
		return fmt.Errorf("next is not successor, head: %d, next: %d", headN, nextN)
	case nextHash != headHash:
		return fmt.Errorf("head header mismatch, expected %x, got: %x", nextHash, headHash)
	}
	return nil
}

// extractEvents extracts all deposit and exit events from a block.
func extractEvents(
	block *tee.Block,
	params *tee.Parameters,
) (deps []*erdstallDepEvent, exits []*erdstallExitEvent, _ error) {
	predicates := []logPredicate{logIsDepositEvt, logIsExitEvt}

	for _, r := range block.Receipts {
		logs := filterLogs(r.Logs, predicates, params.Contract)
		depLogs, exitLogs := logs[0], logs[1]
		for _, depLog := range depLogs {
			if dep, err := parseDepEvent(depLog); err != nil {
				return nil, nil, fmt.Errorf("failed to parse deposit: %w", err)
			} else {
				deps = append(deps, dep)
			}
		}
		for _, exitLog := range exitLogs {
			if exit, err := parseExitEvent(exitLog); err != nil {
				return nil, nil, fmt.Errorf("failed to parse exit: %w", err)
			} else {
				exits = append(exits, exit)
			}
		}
	}
	return
}

// filterLogs partitions logs into different buckets of matching predicates.
// Only logs from the Erdstall contract are filtered and other logs, as well as
// those without a matching predicate, are discarded.
func filterLogs(logs []*types.Log, preds []logPredicate, contract common.Address) [][]*types.Log {
	buckets := make([][]*types.Log, len(preds))
	for _, l := range logs {
		if l.Address != contract {
			// only parse Erdstall logs
			continue
		}
		for i, p := range preds {
			if p(l) {
				buckets[i] = append(buckets[i], l)
			}
		}
	}
	return buckets
}

func verifyDeposits(
	deposits []*erdstallDepEvent,
	depEpoch tee.Epoch,
) error {
	for i, dep := range deposits {
		if dep.Epoch != depEpoch {
			return fmt.Errorf("invalid epoch %d != %d in deposit[%d]", dep.Epoch, depEpoch, i)
		}
	}
	return nil
}

func verifyExits(
	exits []*erdstallExitEvent,
	exitEpoch tee.Epoch,
	epoch *Epoch,
) error {
	for i, exit := range exits {
		if exit.Epoch != exitEpoch {
			return fmt.Errorf("invalid epoch %d != %d in exit[%d]", exit.Epoch, exitEpoch, i)
		} else if exit.Value.Cmp(epoch.Balance(exit.Account)) != 0 {
			return fmt.Errorf("balance mismatch in exit[%d]", i)
		} // TODO: assert frozen.
	}
	return nil
}
