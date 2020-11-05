package tee_test

import (
	"testing"

	ttest "github.com/perun-network/erdstall/tee/test"
	"perun.network/go-perun/backend/ethereum/wallet/hd"
	"perun.network/go-perun/pkg/test"

	"github.com/stretchr/testify/require"

	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

func TestTransaction_SignVerify(t *testing.T) {
	require := require.New(t)
	rng := test.Prng(t)
	hdw := eth.NewHdWallet(rng)
	w, err := hd.NewWallet(hdw, hd.DefaultRootDerivationPath.String(), 0)
	require.NoError(err)
	sender, err := w.NewAccount()
	require.NoError(err)

	contract := ttest.NewRandomAddress(rng)

	tx := ttest.NewTxFrom(rng, sender.Account.Address)

	require.NoError(tx.Sign(contract, sender.Account, hdw))

	ok, err := tee.VerifyTransaction(contract, *tx)
	require.True(ok)
	require.NoError(err)
}
