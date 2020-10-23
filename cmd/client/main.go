// SPDX-License-Identifier: Apache-2.0

package main

import (
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
	chainEvents := make(chan string, 10)           // GUI event pipe
	clEvents := make(chan *client.ClientEvent, 10) // GUI event pipe

	cb := perunchannel.NewContractBackend(eb, perunhd.NewTransactor(wallet.Wallet.Wallet()))
	conn := client.NewRPC("127.0.0.1", 8080)                 // Operator conn
	ethClient := eth.NewClient(cb, wallet.Acc.Account)       // ETHChain conn
	client := client.NewClient(cfg, conn, ethClient, wallet) // Erdstall protocol client
	go gui.RunGui(client, clEvents, chainEvents)             // Run the GUI

	client.Run(clEvents, chainEvents) // Run the protocol
}
