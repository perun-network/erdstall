package tee_test

import (
	"math/big"
	"math/rand"
	"testing"

	"perun.network/go-perun/backend/ethereum/wallet/hd"
	wtest "perun.network/go-perun/backend/ethereum/wallet/test"
	"perun.network/go-perun/pkg/test"

	"github.com/ethereum/go-ethereum/common"
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

	contract := newRandomAddress(rng)

	tx := tee.Transaction{
		Nonce:     uint64(rng.Int63()),
		Epoch:     uint64(rng.Int63()),
		Sender:    sender.Account.Address,
		Recipient: newRandomAddress(rng),
		Amount:    big.NewInt(rng.Int63()),
	}

	require.NoError(tx.Sign(contract, sender.Account, hdw))

	ok, err := tee.VerifyTransaction(contract, tx)
	require.True(ok)
	require.NoError(err)
}

func newRandomAddress(rng *rand.Rand) common.Address {
	return common.Address(wtest.NewRandomAddress(rng))
}
