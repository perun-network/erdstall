// SPDX-License-Identifier: Apache-2.0

package test

import (
	"io"
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"perun.network/go-perun/backend/ethereum/wallet/hd"
	wtest "perun.network/go-perun/backend/ethereum/wallet/test"

	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

func RandomTx(t *testing.T, rng *rand.Rand) *tee.Transaction {
	hdw := eth.NewHdWallet(rng)
	w, err := hd.NewWallet(hdw, hd.DefaultRootDerivationPath.String(), 0)
	require.NoError(t, err)
	sender, err := w.NewAccount()
	require.NoError(t, err)

	contract := NewRandomAddress(rng)
	tx := NewTxFrom(rng, sender.Account.Address)

	require.NoError(t, tx.Sign(contract, sender.Account, hdw))
	return tx
}

func RandomBP(rng *rand.Rand) *tee.BalanceProof {
	return &tee.BalanceProof{
		Balance: RandomBalance(rng),
		Sig:     RandomSig(rng),
	}
}

func RandomDP(rng *rand.Rand) *tee.DepositProof {
	return &tee.DepositProof{
		Balance: RandomBalance(rng),
		Sig:     RandomSig(rng),
	}
}

// RandomSig produces invalid ECDSA signatures, but does not need a wallet
// to work.
func RandomSig(rng io.Reader) tee.Sig {
	sig := make([]byte, 65)
	if _, err := rng.Read(sig); err != nil {
		panic("Could not read from rng.")
	}
	return sig
}

func RandomBalance(rng *rand.Rand) tee.Balance {
	return tee.Balance{
		Epoch:   rng.Uint64(),
		Account: common.Address(wtest.NewRandomAddress(rng)),
		Value:   big.NewInt(rng.Int63()),
	}
}
