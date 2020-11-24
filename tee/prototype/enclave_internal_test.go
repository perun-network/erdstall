// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"perun.network/go-perun/pkg/test"

	cltest "github.com/perun-network/erdstall/client/test"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

type clientAction = func(ctx context.Context, bal *tee.BalanceProof) error

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

	seal := func(phase string, n uint64) {
		testLog("Adding %d new blocks to seal %s.", n, phase)
		for i := uint64(0); i < n; i++ {
			setup.SimBackend.Commit()
		}
	}

	sub, err := operator.SubscribeBlocks()
	requiree.NoError(err)
	defer sub.Unsubscribe()
	testLog("Subscribed to new blocks")

	requiree.NoError(operator.DeployContracts(&params))

	ct := test.NewConcurrent(t)
	// Start enclave routines
	go ct.StageN("operator", 2, func(t test.ConcT) {
		assert.NoError(enc.Run(params))
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
	require.NoError(t, err)
	bobEthCl := eth.NewClient(*setup.CB, bobAd)
	bob, err := cltest.NewClient(params, setup.HdWallet, bobEthCl, encTr)
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	aliceInitBal := setup.Balance(alice.Address())
	bobInitBal := setup.Balance(bob.Address())

	// Do deposits
	initValue := eth.EthToWeiInt(100)
	requiree.NoError(alice.Deposit(ctx, initValue))
	requiree.NoError(bob.Deposit(ctx, initValue))
	alice.UpdateLastBlockNum()
	// local tracking of balances
	balances := map[common.Address]interface{ Balance() *big.Int }{
		alice.Address(): alice,
		bob.Address():   bob,
	}
	testLog("Deposits made!")

	dps, err := enc.DepositProofs()
	requiree.NoError(err)
	assert.Len(dps, 2)
	for _, dp := range dps {
		ok, err := tee.VerifyDepositProof(params, *dp)
		requiree.True(ok)
		requiree.NoError(err)
		requiree.Zero(initValue.Cmp(dp.Balance.Value))
	}

	bps, err := enc.BalanceProofs()
	requiree.NoError(err)
	assert.Len(bps, 0)

	// now: deposit epoch: 1, tx epoch: 0

	testLog("Sending three TXs.")

	requiree.NoError(alice.SendToClient(bob, eth.EthToWeiInt(5)))
	requiree.NoError(bob.SendToClient(alice, eth.EthToWeiInt(10)))
	requiree.NoError(alice.SendToClient(bob, eth.EthToWeiInt(2)))

	// send invalid txs
	errs := alice.SendInvalidTxs(rng, bob.Address())
	for i, err := range errs {
		requiree.Errorf(err, "expected invalid tx #%d", i)
	}

	seal("txPhase", params.PhaseDuration)

	// now: deposit epoch: 2, tx epoch: 1, exit epoch: 0

	testLog("Getting deposit proofs.")
	dps, err = enc.DepositProofs()
	requiree.NoError(err)
	assert.Len(dps, 0)

	testLog("Getting balance proofs.")
	bps, err = enc.BalanceProofs()
	requiree.NoError(err)
	verifyBalanceProofs(t, params, balances, bps)

	doWith := func(bp *tee.BalanceProof, aliceDo, bobDo clientAction) {
		switch a := bp.Balance.Account; a {
		case alice.Address():
			requiree.NoError(aliceDo(ctx, bp))
		case bob.Address():
			requiree.NoError(bobDo(ctx, bp))
		default:
			t.Fatalf("Unknown account: %s", a.String())
		}
	}

	testLog("Sending two exit TXs.")
	for _, bp := range bps {
		doWith(bp, alice.Exit, bob.Exit)
	}

	// Signalling the enclave to stop now, so that it doesn't start new epochs on
	// the next block.
	testLog("Set Enclave to shutdown after next phase.")
	enc.Shutdown()

	seal("exitPhase", 1)

	// Check balance proofs after all users exited.
	// This also serves as a synchronization point so that enough blocks are
	// processed by the enclave before the block subscription is closed.
	bpsExit, err := enc.BalanceProofs()
	requiree.NoError(err)
	requiree.Len(bpsExit, 0, "all users should have exited the system")

	testLog("Closing block subscription and waiting for operator routines...")
	sub.Unsubscribe()
	ct.Wait("operator")

	testLog("Sending two withdrawal TXs")
	for _, bp := range bps {
		doWith(bp, alice.Withdraw, bob.Withdraw)
	}

	testLog("Checking final balances.")
	aliceNewBal := setup.Balance(alice.Address())
	bobNewBal := setup.Balance(bob.Address())
	requiree.NoError(checkBals(
		aliceInitBal,
		aliceNewBal,
		bobInitBal,
		bobNewBal,
		eth.EthToWeiInt(3),
	))
}

// checkBals checks whether the new balance of Alice has increased by `difference`
// and conversely whether the new balance of Bob has decreased. All parameters are
// assumed to be denominated in `Wei`.
func checkBals(aliceInit, aliceNew, bobInit, bobNew, difference *big.Int) error {
	aliceGain := new(big.Int).Sub(aliceNew, aliceInit)
	bobGain := new(big.Int).Sub(bobNew, bobInit)
	delta := big.NewInt(1000000)
	if err := checkInRange(difference, aliceGain, delta); err != nil {
		return fmt.Errorf("checking alice bals: %w", err)
	}
	if err := checkInRange(new(big.Int).Neg(difference), bobGain, delta); err != nil {
		return fmt.Errorf("checking bob bals: %w", err)
	}
	return nil
}

// checkInRange checks, that given `value` is in range of `median` +- `delta`.
func checkInRange(median, value, delta *big.Int) error {
	const LT, GT = -1, 1
	lowerBound := new(big.Int).Sub(median, delta)
	upperBound := new(big.Int).Add(median, delta)
	if lowerBound.Cmp(value) != LT || upperBound.Cmp(value) != GT {
		return fmt.Errorf("value: %v not in range of median: %v", value, median)
	}
	return nil
}

func verifyBalanceProofs(t require.TestingT,
	params tee.Parameters,
	expBalances map[common.Address]interface{ Balance() *big.Int },
	bps []*tee.BalanceProof) {
	require := require.New(t)
	require.Len(bps, len(expBalances))
	for _, bp := range bps {
		ok, err := tee.VerifyBalanceProof(params, *bp)
		require.True(ok)
		require.NoError(err)
		require.Contains(expBalances, bp.Balance.Account)
		got, exp := bp.Balance.Value, expBalances[bp.Balance.Account].Balance()
		require.Zerof(got.Cmp(exp),
			"balance mismatch for %s, got: %v, expected: %v [ETH]", bp.Balance.Account.String(), eth.WeiToEthFloat(got), eth.WeiToEthFloat(exp))
	}
}

func testLog(format string, args ...interface{}) {
	log.Infof("Test: "+format, args...)
}
