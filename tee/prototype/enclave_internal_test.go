// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"io"
	"testing"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/stretchr/testify/assert"
	chtest "perun.network/go-perun/backend/ethereum/channel/test"
	"perun.network/go-perun/pkg/test"

	"github.com/perun-network/erdstall/tee"
)

func TestEnclave(t *testing.T) {
	rng := test.Prng(t)

	hdw := newHdWallet(t, rng)
	enc := NewEnclave(hdw)

	teeAddr, _, err := enc.Init() // ignore attestation for now
	assert.NoError(t, err)

	params := tee.Parameters{
		PowDepth:         1,
		PhaseDuration:    3,
		ResponseDuration: 1,
		TEE:              teeAddr,
	}
	_ = params

	// Setup blockchain
	setup := chtest.NewSimSetup(rng)
	_ = setup
}

func newHdWallet(t *testing.T, rng io.Reader) *hdwallet.Wallet {
	seed := make([]byte, 20)
	rng.Read(seed)
	hdw, err := hdwallet.NewFromSeed(seed)
	assert.NoError(t, err)
	return hdw
}
