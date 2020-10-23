// SPDX-License-Identifier: Apache-2.0

package config

import "flag"

type ClientConfig struct {
	ChainURL     string
	Mnemonic     string
	AccountIndex int
	Contract     string
}

func ParseClientConfig() (cfg ClientConfig) {
	flag.StringVar(&cfg.ChainURL, "chain-url", "ws://127.0.0.1:8545", "'protocol://ip:port'")
	flag.StringVar(&cfg.Mnemonic, "mnemonic", "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic", "Wallet mnemonic.")
	flag.IntVar(&cfg.AccountIndex, "account-index", 0, "Account derivation index.")
	flag.StringVar(&cfg.Contract, "contract", "", "Address of the Erdstall contract. Empty means deploy.")
	flag.Parse()
	return
}
