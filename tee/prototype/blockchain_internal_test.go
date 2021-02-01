// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/perun-network/erdstall/tee"
)

func TestBlockchain(t *testing.T) {
	var (
		require = require.New(t)
		ethbc   = backends.NewSimulatedBackend(nil, 8000000)
		head    = func() *tee.Block {
			b, err := ethbc.BlockByNumber(context.Background(), nil)
			require.NoError(err)
			return &tee.Block{Block: *b}
		}
		bc     blockchain
		params = &tee.Parameters{PhaseDuration: 10}
		epoch  = newEpoch(0)
	)

	t.Run("empty chain", func(t *testing.T) {
		require.True(bc.empty())
		require.Nil(bc.Head())
	})

	t.Run("origin block", func(t *testing.T) {
		ethbc.Commit()
		deps, exits, err := bc.PushVerify(head(), params, epoch)
		require.Empty(deps)
		require.Empty(exits)
		require.NoError(err)
		require.Equal(bc.Head(), head())
	})

	ethbc.Commit()
	b1 := head()

	t.Run("normal successor", func(t *testing.T) {
		deps, exits, err := bc.PushVerify(b1, params, epoch)
		require.Empty(deps)
		require.Empty(exits)
		require.NoError(err)
	})

	t.Run("duplicate block", func(t *testing.T) {
		_, _, err := bc.PushVerify(b1, params, epoch)
		require.Error(err)
	})

	ethbc.Commit()

	t.Run("Invalid parent hash", func(t *testing.T) {
		h2 := head().Header()
		h2.ParentHash[0] ^= 0x3f
		b2 := &tee.Block{Block: *types.NewBlockWithHeader(h2)}
		_, _, err := bc.PushVerify(b2, params, epoch)
		require.Error(err)
		require.Equal(bc.Head(), b1)
	})

	b2 := head()
	ethbc.Commit()
	b3 := head()

	t.Run("block gap", func(t *testing.T) {
		_, _, err := bc.PushVerify(b3, params, epoch)
		require.Error(err)
		require.Equal(bc.Head(), b1)
		_, _, err = bc.PushVerify(b2, params, epoch)
		require.NoError(err)
		require.Equal(bc.Head(), b2)
	})

	t.Run("past block", func(t *testing.T) {
		_, _, err := bc.PushVerify(b1, params, epoch)
		require.Error(err)
		require.Equal(bc.Head(), b2)
	})
}
