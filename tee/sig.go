// SPDX-License-Identifier: Apache-2.0

package tee

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"perun.network/go-perun/backend/ethereum/wallet"
)

func VerifyDepositProof(params Parameters, proof DepositProof) (bool, error) {
	msg, err := EncodeDepositProof(params.Contract, proof.Balance)
	if err != nil {
		return false, fmt.Errorf("encoding balance: %w", err)
	}
	return wallet.VerifySignature(msg, proof.Sig, (*wallet.Address)(&params.TEE))
}

func VerifyBalanceProof(params Parameters, proof BalanceProof) (bool, error) {
	msg, err := EncodeBalanceProof(params.Contract, proof.Balance)
	if err != nil {
		return false, fmt.Errorf("encoding balance: %w", err)
	}
	return wallet.VerifySignature(msg, proof.Sig, (*wallet.Address)(&params.TEE))
}

func VerifyTransaction(contract common.Address, tx Transaction) (bool, error) {
	msg, err := EncodeTransaction(contract, tx)
	if err != nil {
		return false, fmt.Errorf("encoding tx: %w", err)
	}
	return wallet.VerifySignature(msg, tx.Sig, (*wallet.Address)(&tx.Sender))
}
