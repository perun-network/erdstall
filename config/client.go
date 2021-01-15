// SPDX-License-Identifier: Apache-2.0

package config

import (
	"encoding/json"
	"flag"

	log "github.com/sirupsen/logrus"
)

type ClientConfig struct {
	ChainURLs    map[string]string // map[ChainID]url
	OpHost       string
	OpPort       int
	Mnemonic     string
	AccountIndex int
	UserName     string
}

func ParseClientConfig() (cfg ClientConfig) {
	var urlsJson string
	flag.StringVar(&urlsJson, "chain-urls", `{"1337": "ws://127.0.0.1:8545"}`, `JSON dictionary {"chainID": Ethereum node URL}`)
	flag.StringVar(&cfg.OpHost, "op-host", "127.0.0.1", "IP/host name of operator")
	flag.IntVar(&cfg.OpPort, "op-port", 8401, "Port of operator.")
	flag.StringVar(&cfg.Mnemonic, "mnemonic", "pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic", "Wallet mnemonic.")
	flag.IntVar(&cfg.AccountIndex, "account-index", 0, "Account derivation index.")
	flag.StringVar(&cfg.UserName, "username", "<anonymous>", "Set an optional username.")
	flag.Parse()

	cfg.ChainURLs = parseChainURLs(urlsJson)

	return
}

func parseChainURLs(jsonText string) map[string]string {
	m := make(map[string]string)
	var urls map[string]interface{}
	if err := json.Unmarshal([]byte(jsonText), &urls); err != nil {
		log.WithError(err).Error("Unmarshal chain URL")
	}

	for k, v := range urls {
		vs, ok := v.(string)
		if !ok {
			log.Panicf("Client config: chain URLs: expected string, got `%v`", v)
		}
		m[k] = vs
	}
	return m
}

// ChainURL looks up the configured ethereum network URL for a network ID.
func (c *ClientConfig) ChainURL(networkID string) string {
	chURL, ok := c.ChainURLs[networkID]
	if !ok {
		log.WithField("network ID", networkID).
			Fatal("Unknown ethereum network ID")
	}

	return chURL
}
