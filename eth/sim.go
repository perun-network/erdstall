// SPDX-License-Identifier: Apache-2.0

package eth

import (
	"context"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	log "github.com/sirupsen/logrus"
	ethchannel "perun.network/go-perun/backend/ethereum/channel"
	chtest "perun.network/go-perun/backend/ethereum/channel/test"
	"perun.network/go-perun/backend/ethereum/wallet/hd"
)

type SimSetup struct {
	SimBackend *chtest.SimulatedBackend    // A simulated blockchain backend
	Accounts   []accounts.Account          // funded accounts
	CB         *ethchannel.ContractBackend // contract backend that can transact for all accounts
	Wallet     *hd.Wallet                  // Wallet containing accounts (Perun wrapper)
	HdWallet   *hdwallet.Wallet            // Wallet containing accounts (original)
}

func NewSimSetup(rng *rand.Rand, numAccounts int) *SimSetup {
	simBackend := chtest.NewSimulatedBackend()

	wallet := NewHdWallet(rng)
	pwallet, err := hd.NewWallet(wallet, hd.DefaultRootDerivationPath.String(), 0)
	if err != nil {
		log.Panicf("Error wrapping hd wallet: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	accs := make([]accounts.Account, numAccounts)
	for i := 0; i < numAccounts; i++ {
		acc, err := pwallet.NewAccount()
		if err != nil {
			log.Panicf("Error creating new account: %v", err)
		}
		accs[i] = acc.Account
		simBackend.FundAddress(ctx, acc.Account.Address)
	}

	cb := ethchannel.NewContractBackend(simBackend, hd.NewTransactor(pwallet.Wallet()))

	return &SimSetup{
		SimBackend: simBackend,
		Accounts:   accs,
		CB:         &cb,
		Wallet:     pwallet,
		HdWallet:   wallet,
	}
}
