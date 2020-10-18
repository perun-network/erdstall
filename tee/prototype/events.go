// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/perun-network/erdstall/contracts/bindings"
)

var depositedEvent common.Hash = crypto.Keccak256Hash([]byte("Deposited(uint64,address,uint256)"))

var exitingEvent common.Hash = crypto.Keccak256Hash([]byte("Exiting(uint64,address,uint256)"))

// erdstallDepEvent is a generic wrapper type for `Deposited` events.
type erdstallDepEvent struct {
	Epoch   uint64
	Account common.Address
	Value   *big.Int
}

// erdstallDepEvent is wrapper type for `Exiting` events.
type erdstallExitEvent struct {
	Epoch   uint64
	Account common.Address
	Value   *big.Int
}

var contractAbi, _ = abi.JSON(strings.NewReader(bindings.ErdstallABI))

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
