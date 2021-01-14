// SPDX-License-Identifier: Apache-2.0

package config

import "flag"

type ClientConfig struct {
	ChainURL     string
	OpHost       string
	OpPort       int
	Mnemonic     string
	AccountIndex int
	Contract     string
	UserName     string
}

func ParseClientConfig() (cfg ClientConfig) {
	flag.StringVar(&cfg.ChainURL, "chain-url", "ws://127.0.0.1:8545", "'protocol://ip:port'")
	flag.StringVar(&cfg.OpHost, "op-host", "127.0.0.1", "IP/host name of operator")
	flag.IntVar(&cfg.OpPort, "op-port", 8401, "Port of operator.")
	flag.StringVar(&cfg.Mnemonic, "mnemonic", "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic", "Wallet mnemonic.")
	flag.IntVar(&cfg.AccountIndex, "account-index", 0, "Account derivation index.")
	flag.StringVar(&cfg.Contract, "contract", "", "Address of the Erdstall contract. Empty means deploy.")
	flag.StringVar(&cfg.UserName, "username", "<anonymous>", "Set an optional username.")
	flag.Parse()
	return
}
