// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
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
		bc blockchain
	)

	ethbc.Commit()
	require.NoError(bc.PushVerify(head()))

	ethbc.Commit()
	b1 := head()
	require.NoError(bc.PushVerify(b1))
	require.NoError(bc.PushVerify(b1))
	require.Len(bc.blocks, 2)

	ethbc.Commit()
	b2 := head()
	require.NoError(bc.PushVerify(b2))
	require.Len(bc.blocks, 3)
	require.NoError(bc.PushVerify(b1))
	require.Len(bc.blocks, 2, "pushing previous block should shorten bc")
	require.Equal(bc.Head(), b1)

	ethbc.Commit()
	b3 := head()
	require.Error(bc.PushVerify(b3))
	require.NoError(bc.PushVerify(b2))
	require.NoError(bc.PushVerify(b3))
	require.Len(bc.blocks, 4)
	require.Equal(bc.Head(), b3)
}
