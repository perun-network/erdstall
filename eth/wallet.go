// SPDX-License-Identifier: Apache-2.0

package eth

import (
	"io"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	log "github.com/sirupsen/logrus"
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
