// SPDX-License-Identifier: Apache-2.0

package tee_test

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	_ "perun.network/go-perun/backend/ethereum" // init
	"perun.network/go-perun/backend/ethereum/channel"
	wtest "perun.network/go-perun/backend/ethereum/wallet/test"
	pkgtest "perun.network/go-perun/pkg/test"

	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

func TestErdstallBindings(t *testing.T) {
	rng := pkgtest.Prng(t)
	s := eth.NewSimSetup(rng, 2)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	params := &tee.Parameters{TEE: s.Accounts[0].Address, PhaseDuration: 3, ResponseDuration: 1}
	erdstall, err := deployErdstall(ctx, params, s.CB, s.Accounts[0])
	require.NoError(t, err)
	require.NotNil(t, erdstall)
	opts := &bind.CallOpts{Context: ctx}

	t.Run("testEncodeBalances", func(t *testing.T) {
		testEncodeBalances(t, rng, params, erdstall, opts)
	})
	t.Run("testSigVerify", func(t *testing.T) {
		testSigVerify(t, rng, params, erdstall, opts, s)
	})
}

func testEncodeBalances(t *testing.T, rng *rand.Rand, params *tee.Parameters, contr *bindings.Erdstall, opts *bind.CallOpts) {
	for i := 0; i < 100; i++ {
		bal := randomBalance(rng)
		encodedByGo, err := tee.EncodeBalanceProof(params.Contract, bal)
		require.NoError(t, err)
		encodedBySol, err := contr.EncodeBalanceProof(opts, toEthBals(bal))
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, encodedByGo, encodedBySol)
	}
}

func testSigVerify(t *testing.T, rng *rand.Rand, params *tee.Parameters, contr *bindings.Erdstall, opts *bind.CallOpts, s *eth.SimSetup) {
	for i := 0; i < 100; i++ {
		bal := randomBalance(rng)
		encodedByGo, err := tee.EncodeBalanceProof(params.Contract, bal)
		sig, err := s.HdWallet.SignText(s.Accounts[0], crypto.Keccak256(encodedByGo))
		sig[64] += 27

		require.Contains(t, []byte{27, 28}, sig[64], "v Value should be 27 or 28")
		require.NoError(t, err, "Signing failed")
		proof := tee.BalanceProof{
			Balance: bal,
			Sig:     sig,
		}
		ok, err := tee.VerifyBalanceProof(*params, proof)
		require.True(t, ok)
		require.NoError(t, err)

		err = contr.VerifyBalance(opts, toEthBals(bal), sig)
		require.NoError(t, err, "On-chain verify failed")
	}
}

func deployErdstall(ctx context.Context, params *tee.Parameters, cb *channel.ContractBackend, acc accounts.Account) (*bindings.Erdstall, error) {
	tr, err := cb.NewTransactor(context.Background(), big.NewInt(0), 6000000, acc)
	if err != nil {
		return nil, fmt.Errorf("creating transactor: %w", err)
	}

	address, tx, contract, err := bindings.DeployErdstall(tr, cb.ContractInterface,
		params.TEE,
		params.PhaseDuration,
		params.ResponseDuration)
	if err != nil {
		return nil, fmt.Errorf("deploying contract: %w", err)
	}

	_, err = bind.WaitDeployed(ctx, cb, tx)
	if err != nil {
		return nil, fmt.Errorf("waiting for contract deployment: %w", err)
	}
	params.Contract = address

	receipt, err := cb.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("getting tx receipt: %w", err)
	}
	params.InitBlock = receipt.BlockNumber.Uint64()

	return contract, nil
}

func randomBalance(rng *rand.Rand) tee.Balance {
	return tee.Balance{
		Epoch:   rng.Uint64(),
		Account: common.Address(wtest.NewRandomAddress(rng)),
		Value:   big.NewInt(rng.Int63()),
	}
}

func toEthBals(b tee.Balance) bindings.ErdstallBalance {
	return bindings.ErdstallBalance{
		Epoch:   b.Epoch,
		Account: b.Account,
		Value:   b.Value,
	}
}
