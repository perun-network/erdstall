// SPDX-License-Identifier: Apache-2.0

package eth

import (
	"io"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	log "github.com/sirupsen/logrus"
	wtest "perun.network/go-perun/backend/ethereum/wallet/test"
)

func NewHdWallet(rng io.Reader) *hdwallet.Wallet {
	seed := make([]byte, 20)
	_, err := rng.Read(seed)
	if err != nil {
		log.Panicf("Error reading randomness: %v", err)
	}
	hdw, err := hdwallet.NewFromSeed(seed)
	if err != nil {
		log.Panicf("Error creating random hdwallet: %v", err)
	}
	return hdw
}

// NewRandomAddress creates a new random address.
func NewRandomAddress(rng *rand.Rand) common.Address {
	return common.Address(wtest.NewRandomAddress(rng))
}
