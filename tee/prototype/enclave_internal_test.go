// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"context"
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
	operator := eth.NewClient(*setup.CB, setup.Accounts[0])

	sub, err := operator.SubscribeToBlocks()
	requiree.NoError(err)
	defer sub.Unsubscribe()
	t.Log("Subscribed to new blocks")

	// Start mini-operator
	ct := test.NewConcurrent(t)
	go ct.Stage("operator", func(t test.ConcT) {
		for b := range sub.Blocks() {
			require.NoError(t, enc.ProcessBlocks(b))
		}
	})

	requiree.NoError(operator.DeployContracts(&params))
	requiree.NoError(enc.SetParams(params))

	// Create clients
	encTr := &cltest.EnclaveTransactor{enc} // transact directly on the enclave, bypassing the operator
	aliceEthCl := eth.NewClient(*setup.CB, setup.Accounts[1])
	alice, err := cltest.NewClient(params, setup.HdWallet, aliceEthCl, encTr)
	bobEthCl := eth.NewClient(*setup.CB, setup.Accounts[2])
	bob, err := cltest.NewClient(params, setup.HdWallet, bobEthCl, encTr)

	// Do deposits
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	requiree.NoError(alice.Deposit(ctx, eth.EthToWeiInt(1)))
	requiree.NoError(bob.Deposit(ctx, eth.EthToWeiInt(1)))
	t.Log("Deposits made!")

	sub.Unsubscribe()
	ct.Wait("operator")
}
