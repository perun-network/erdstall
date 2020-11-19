// SPDX-License-Identifier: Apache-2.0

package main

import (
	log "github.com/sirupsen/logrus"
	perunchannel "perun.network/go-perun/backend/ethereum/channel"
	perunhd "perun.network/go-perun/backend/ethereum/wallet/hd"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/perun-network/erdstall/client"
	"github.com/perun-network/erdstall/config"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/gui"
	"github.com/perun-network/erdstall/wallet"
)

func main() {
	cfg := config.ParseClientConfig()

	wallet := wallet.NewWallet(cfg.Mnemonic, uint(cfg.AccountIndex)) // HD Wallet
	eb, err := ethclient.Dial(cfg.ChainURL)
	if err != nil {
		panic(err)
	}
	events := make(chan *client.Event, 10) // GUI event pipe

	cb := perunchannel.NewContractBackend(eb, perunhd.NewTransactor(wallet.Wallet.Wallet()))
	rpc := client.NewRPC("127.0.0.1", 8080)                     // Operator conn
	chain := eth.NewClient(cb, wallet.Acc.Account)              // ETHChain conn
	client := client.NewClient(cfg, rpc, events, chain, wallet) // Erdstall protocol client
	go gui.RunGui(client, events)                               // Run the GUI

	if err := client.Run(); err != nil { // Run the protocol
		log.WithError(err).Fatal("Client crashed")
	}
}
