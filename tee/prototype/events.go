// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/perun-network/erdstall/contracts/bindings"
)

var (
	contractAbi, _ = abi.JSON(strings.NewReader(bindings.ErdstallABI))
	depositedEvent = contractAbi.Events["Deposited"].ID
	exitingEvent   = contractAbi.Events["Exiting"].ID
)

type (
	// erdstallDepEvent is a generic wrapper type for `Deposited` events.
	erdstallDepEvent struct {
		Epoch   uint64
		Account common.Address
		Value   *big.Int
	}

	// erdstallDepEvent is wrapper type for `Exiting` events.
	erdstallExitEvent struct {
		Epoch   uint64
		Account common.Address
		Value   *big.Int
	}

	logPredicate = func(l *types.Log) bool
)

// parseDepEvent parses a given `log` and returns an `erdstallEvent`.
func parseDepEvent(l *types.Log) (*erdstallDepEvent, error) {
	name := "Deposited"
	event := new(erdstallDepEvent)
	err := unpackLog(event, name, *l)
	if err != nil {
		return nil, fmt.Errorf("unpacking %s: %w", name, err)
	}
	return event, nil
}

// parseExitEvent parses a given `log` and returns an `erdstallEvent`.
func parseExitEvent(l *types.Log) (*erdstallExitEvent, error) {
	name := "Exiting"
	event := new(erdstallExitEvent)
	err := unpackLog(event, name, *l)
	if err != nil {
		return nil, fmt.Errorf("unpacking %v: %w", name, err)
	}

	return event, nil
}

// UnpackLog unpacks a retrieved log into the provided output structure.
func unpackLog(out interface{}, event string, log types.Log) error {
	if len(log.Data) > 0 {
		if err := contractAbi.Unpack(out, event, log.Data); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contractAbi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return abi.ParseTopics(out, indexed, log.Topics[1:])
}

func logIsDepositEvt(l *types.Log) bool {
	return l.Topics[0].String() == depositedEvent.String()
}

func logIsExitEvt(l *types.Log) bool {
	return l.Topics[0].String() == exitingEvent.String()
}
