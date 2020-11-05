// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/tee/test"
	"github.com/stretchr/testify/require"
	ethtest "perun.network/go-perun/backend/ethereum/wallet/test"
	pkgtest "perun.network/go-perun/pkg/test"
)

func TestTxCache(t *testing.T) {
	rng := pkgtest.Prng(t)
	requiree := require.New(t)
	N := 10
	userA := common.Address(ethtest.NewRandomAddress(rng))
	userB := common.Address(ethtest.NewRandomAddress(rng))

	txc := makeTxCache()
	for j := 0; j < N; j++ {
		txc.cacheTx(test.NewTxFrom(rng, userA))
	}
	txs, ok := txc.senders[userA]
	requiree.True(ok)
	requiree.Equal(N, len(txs))
	for j := 0; j < N; j++ {
		txc.cacheTx(test.NewTxFromTo(rng, userA, userB))
	}
	txs, ok = txc.recipients[userB]
	requiree.True(ok)
	requiree.Equal(N, len(txs))
}

func TestInconsistentExits(t *testing.T) {
	rng := pkgtest.Prng(t)
	requiree := require.New(t)
	userA := common.Address(ethtest.NewRandomAddress(rng))
	userB := common.Address(ethtest.NewRandomAddress(rng))
	txc := makeTxCache()
	txc.cacheTx(test.NewTxFromTo(rng, userA, userB))
	requiree.NoError(noInconsistentExits(&txc, exitersSet{}))
	requiree.Error(noInconsistentExits(&txc, exitersSet{userA}))
	requiree.Error(noInconsistentExits(&txc, exitersSet{userB}))
	requiree.Error(noInconsistentExits(&txc, exitersSet{userA, userB}))
}
