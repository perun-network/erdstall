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
	err := contractAbi.Unpack(event, name, l.Data)
	if err != nil {
		return nil, fmt.Errorf("unpacking %v : %w", name, err)
	}

	return event, nil
}

// parseExitEvent parses a given `log` and returns an `erdstallEvent`.
func parseExitEvent(l *types.Log) (*erdstallExitEvent, error) {
	name := "Exiting"
	event := new(erdstallExitEvent)
	err := contractAbi.Unpack(event, name, l.Data)
	if err != nil {
		return nil, fmt.Errorf("unpacking %v : %w", name, err)
	}

	return event, nil
}

func logIsDepositEvt(l *types.Log) bool {
	return l.Topics[0].String() == depositedEvent.String()
}

func logIsExitEvt(l *types.Log) bool {
	return l.Topics[0].String() == exitingEvent.String()
}
