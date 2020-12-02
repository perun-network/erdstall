// SPDX-License-Identifier: Apache-2.0

package tee

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	abiUint256, _ = abi.NewType("uint256", "", nil)
	abiUint64, _  = abi.NewType("uint64", "", nil)
	abiAddress, _ = abi.NewType("address", "", nil)
	abiString, _  = abi.NewType("string", "", nil)
)

func EncodeDepositProof(contract common.Address, balance Balance) ([]byte, error) {
	return encodeBalance("ErdstallDeposit", contract, balance)
}

func EncodeBalanceProof(contract common.Address, balance Balance) ([]byte, error) {
	return encodeBalance("ErdstallBalance", contract, balance)
}

func encodeBalance(tag string, contract common.Address, balance Balance) ([]byte, error) {
	return abi.Arguments{
		{Type: abiString},  // tag
		{Type: abiAddress}, // contract
		{Type: abiUint64},  // epoch
		{Type: abiAddress}, // account
		{Type: abiUint256}, // balance
	}.Pack(
		tag,
		contract,
		balance.Epoch,
		balance.Account,
		(*big.Int)(balance.Value),
	)
}

// EncodeTransaction abi-encodes an off-chain transaction. We use abi encoding
// for consistency even though this message is never used on-chain.
// Should only be used for signing purposes.
func EncodeTransaction(contract common.Address, tx Transaction) ([]byte, error) {
	return abi.Arguments{
		{Type: abiString},  // tag
		{Type: abiAddress}, // contract
		{Type: abiUint64},  // nonce
		{Type: abiUint64},  // epoch
		{Type: abiAddress}, // sender
		{Type: abiAddress}, // recipient
		{Type: abiUint256}, // amount
	}.Pack(
		"ErdstallTransaction",
		contract,
		tx.Nonce,
		tx.Epoch,
		tx.Sender,
		tx.Recipient,
		(*big.Int)(tx.Amount),
	)
}
