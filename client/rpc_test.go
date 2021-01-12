// SPDX-License-Identifier: Apache-2.0

package client_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pkgtest "perun.network/go-perun/pkg/test"

	"github.com/perun-network/erdstall/client"
	"github.com/perun-network/erdstall/operator"
	optest "github.com/perun-network/erdstall/operator/test"
	"github.com/perun-network/erdstall/tee"
	ttest "github.com/perun-network/erdstall/tee/test"
)

var shortWait = 300 * time.Millisecond
var longWait = 10 * time.Second

const opRPCPort = 8401

// TestRPC_ClientOp is a test which tests all functions of the Operator
// over a websocket connection.
//
// It uses a mocked remote enclave to simulate failing functions.
// The construction looks a bit like this:
// Client <-> RPC <-> MockedRPCOperator <-> OP
func TestRPC_ClientOp(t *testing.T) {
	rng := pkgtest.Prng(t)
	enclave := optest.NewMockedEnclave()
	op := optest.NewRPROperator(enclave)
	op.Run()
	osc := operator.OpServerConfig{
		Host:         "",
		Port:         opRPCPort,
		ClientConfig: operator.ClientConfig{},
	}
	rpcServer := operator.NewRPC(op, osc)
	go func() {
		if err := rpcServer.Serve(); err != nil {
			panic(err)
		}
	}()

	time.Sleep(shortWait)
	rpcClient, err := client.NewRPC("0.0.0.0", opRPCPort)
	myErr := errors.New("Should fail")
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), longWait)
	defer cancel()

	t.Run("SendTx-error", func(t *testing.T) {
		tx := ttest.RandomTx(t, rng)
		enclave.SetProcessTXsError(myErr)
		err = rpcClient.SendTx(ctx, *tx)
		assert.Error(t, err)
		enclave.SetProcessTXsError(nil)
	})

	t.Run("SendTx-ok", func(t *testing.T) {
		tx := ttest.RandomTx(t, rng)
		err = rpcClient.SendTx(ctx, *tx)
		require.NoError(t, err)
		tx2 := <-enclave.Transactions()
		assert.Equal(t, tx, tx2)
	})

	dp := ttest.RandomDP(rng)
	bp := ttest.RandomBP(rng)
	user1 := dp.Balance.Account
	bp.Balance.Account = user1
	var sub *client.Subscription

	for i := int64(0); i < 5; i++ {
		// Push in proofs for our client.
		dp.Balance.Value = (*tee.Amount)(big.NewInt(i))
		bp.Balance.Value = (*tee.Amount)(big.NewInt(i))
		enclave.PushDepositProof(dp)
		enclave.PushBalanceProof(bp)
		// Push in proofs for other clients.
		enclave.PushBalanceProof(ttest.RandomBP(rng))
		enclave.PushDepositProof(ttest.RandomDP(rng))
	}

	t.Run("Subscribe-error", func(t *testing.T) {
		op.SetSubscribeProofsError(myErr)
		_, err := rpcClient.Subscribe(ctx, user1)
		assert.Error(t, err)
		op.SetSubscribeProofsError(nil)
	})

	// Wait for the OP to process the enclave's proofs.
	time.Sleep(shortWait)
	t.Run("Subscribe-ok", func(t *testing.T) {
		sub, err = rpcClient.Subscribe(ctx, user1)
		assert.NoError(t, err)
		require.NotNil(t, sub)
	})

	// Test that the client receives the latest buffered proofs.
	t.Run("Subscription-DP-buffer-latest", func(t *testing.T) {
		proof, err := sub.DepositProof(ctx)
		require.NoError(t, err)
		assert.Equal(t, *dp, proof)
	})

	t.Run("Subscription-BP-buffer-latest", func(t *testing.T) {
		proof, err := sub.BalanceProof(ctx)
		require.NoError(t, err)
		assert.Equal(t, *bp, proof)
	})

	// Test that the client receives the latest unbuffered proofs.
	t.Run("Subscription-BP-latest", func(t *testing.T) {
		bp.Balance.Epoch = 1234
		enclave.PushBalanceProof(bp)
		time.Sleep(shortWait) // Wait for OP.
		proof, err := sub.BalanceProof(ctx)
		require.NoError(t, err)
		assert.Equal(t, *bp, proof)
	})

	t.Run("Subscription-DP-latest", func(t *testing.T) {
		dp.Balance.Epoch = 1234
		enclave.PushDepositProof(dp)
		time.Sleep(shortWait) // Wait for OP.
		proof, err := sub.DepositProof(ctx)
		require.NoError(t, err)
		assert.Equal(t, *dp, proof)
	})

	assert.NoError(t, rpcClient.Close())
	assert.NoError(t, rpcServer.Close())
}
