// SPDX-License-Identifier: Apache-2.0

package tee_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	_ "perun.network/go-perun/backend/ethereum" // init
	pkgtest "perun.network/go-perun/pkg/test"

	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
	"github.com/perun-network/erdstall/tee/test"
)

func TestErdstallBindings(t *testing.T) {
	rng := pkgtest.Prng(t)
	s := eth.NewSimSetup(rng, 2)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	params := &tee.Parameters{TEE: s.Accounts[0].Address, PhaseDuration: 3, ResponseDuration: 1}
	erdstall, err := deployErdstall(ctx, params, eth.NewClient(*s.CB, s.Accounts[0]))
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
		bal := test.RandomBalance(rng)
		encodedByGo, err := tee.EncodeBalanceProof(params.Contract, bal)
		require.NoError(t, err)
		encodedBySol, err := contr.EncodeBalanceProof(opts, bal.ToEthBals())
		require.NoError(t, err)

		require.NoError(t, err)
		require.Equal(t, encodedByGo, encodedBySol)
	}
}

func testSigVerify(t *testing.T, rng *rand.Rand, params *tee.Parameters, contr *bindings.Erdstall, opts *bind.CallOpts, s *eth.SimSetup) {
	for i := 0; i < 100; i++ {
		bal := test.RandomBalance(rng)
		encodedByGo, err := tee.EncodeBalanceProof(params.Contract, bal)
		require.NoError(t, err)
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

		err = contr.VerifyBalance(opts, bal.ToEthBals(), sig)
		require.NoError(t, err, "On-chain verify failed")
	}
}

func deployErdstall(ctx context.Context, params *tee.Parameters, cl *eth.Client) (*bindings.Erdstall, error) {
	tr, err := cl.NewTransactor(context.Background())
	if err != nil {
		return nil, fmt.Errorf("creating transactor: %w", err)
	}

	address, tx, contract, err := bindings.DeployErdstall(tr, cl.ContractInterface,
		params.TEE,
		params.PhaseDuration,
		params.ResponseDuration)
	if err != nil {
		return nil, fmt.Errorf("deploying contract: %w", err)
	}

	_, err = bind.WaitDeployed(ctx, cl.ContractBackend, tx)
	if err != nil {
		return nil, fmt.Errorf("waiting for contract deployment: %w", err)
	}
	params.Contract = address

	receipt, err := cl.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("getting tx receipt: %w", err)
	}
	params.InitBlock = receipt.BlockNumber.Uint64()

	return contract, nil
}
