// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"perun.network/go-perun/pkg/test"

	cltest "github.com/perun-network/erdstall/client/test"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

func TestEnclave(t *testing.T) {
	assert, requiree := assert.New(t), require.New(t)
	_ = assert
	rng := test.Prng(t)

	encWallet := eth.NewHdWallet(rng)
	enc := NewEnclave(encWallet)

	teeAddr, _, err := enc.Init() // ignore attestation for now
	requiree.NoError(err)

	params := tee.Parameters{
		PowDepth:         0,
		PhaseDuration:    3,
		ResponseDuration: 1,
		TEE:              teeAddr,
	}

	// Setup blockchain and accounts
	setup := eth.NewSimSetup(rng, 3) // 1 Operator + 2 Clients
	operatorAd, aliceAd, bobAd := setup.Accounts[0], setup.Accounts[1], setup.Accounts[2]
	operator := eth.NewClient(*setup.CB, operatorAd)

	sub, err := operator.SubscribeToBlocks()
	requiree.NoError(err)
	defer sub.Unsubscribe()
	t.Log("Subscribed to new blocks")

	requiree.NoError(operator.DeployContracts(&params))
	requiree.NoError(enc.SetParams(params))

	ct := test.NewConcurrent(t)
	// Start enclave routines
	go ct.StageN("operator", 2, func(t test.ConcT) {
		assert.NoError(enc.Start())
	})

	// Start mini-operator
	go ct.StageN("operator", 2, func(t test.ConcT) {
		for b := range sub.Blocks() {
			err := enc.ProcessBlocks(b)
			if errors.Is(err, ErrEnclaveStopped) {
				return
			}
			require.NoError(t, err)
		}
	})

	// Create clients
	encTr := &cltest.EnclaveTransactor{Enclave: enc} // transact directly on the enclave, bypassing the operator
	aliceEthCl := eth.NewClient(*setup.CB, aliceAd)
	alice, err := cltest.NewClient(params, setup.HdWallet, aliceEthCl, encTr)
	bobEthCl := eth.NewClient(*setup.CB, bobAd)
	bob, err := cltest.NewClient(params, setup.HdWallet, bobEthCl, encTr)

	// Do deposits
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	requiree.NoError(alice.Deposit(ctx, eth.EthToWeiInt(100)))
	requiree.NoError(bob.Deposit(ctx, eth.EthToWeiInt(100)))
	alice.UpdateLastBlockNum()
	t.Log("Deposits made!")

	dps, err := enc.DepositProofs()
	requiree.NoError(err)
	assert.Len(dps, 2)
	for _, dp := range dps {
		ok, err := tee.VerifyDepositProof(params, *dp)
		requiree.True(ok)
		requiree.NoError(err)
	}

	bps, err := enc.BalanceProofs()
	requiree.NoError(err)
	assert.Len(bps, 0)

	// now: deposit epoch: 1, tx epoch: 0

	t.Log("Sending two TXs.")

	requiree.NoError(alice.Send(bobAd.Address, eth.EthToWeiInt(5)))
	requiree.NoError(bob.Send(aliceAd.Address, eth.EthToWeiInt(10)))
	requiree.NoError(alice.Send(bobAd.Address, eth.EthToWeiInt(2)))

	// Signalling the enclave to stop now, so that it doesn't start new epochs on
	// the next block.
	t.Log("Set Enclave to shutdown after next phase.")
	enc.Stop()

	t.Log("Adding 3 new blocks to seal next phase.")

	for i := uint64(0); i < params.PhaseDuration; i++ {
		setup.SimBackend.Commit()
	}

	t.Log("Getting deposit proofs.")
	dps, err = enc.DepositProofs()
	requiree.NoError(err)
	assert.Len(dps, 0)

	t.Log("Getting balance proofs.")
	bps, err = enc.BalanceProofs()
	requiree.NoError(err)
	assert.Len(bps, 2)
	for _, bp := range bps {
		ok, err := tee.VerifyBalanceProof(params, *bp)
		requiree.True(ok)
		requiree.NoError(err)
	}

	sub.Unsubscribe()
	ct.Wait("operator")
}
