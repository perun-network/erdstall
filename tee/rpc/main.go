// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"net"
	"os"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/perun-network/erdstall/tee/prototype"

	"perun.network/go-perun/log"
)

// main(listenAddr:port, mnemonic, derivation path)
// nolint:unused,deadcode
func main() {
	if len(os.Args) != 4 {
		log.Fatalf("main: expected <listen addr> <mnemonic> <derivation path>")
	}

	listenAddr := os.Args[1]
	mnemonic := os.Args[2]
	derivationPath := os.Args[3]

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Panicf("net.Listen: %v", err)
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Panicf("Mnemonic: %v", err)
	}

	enclaveAccountDerivationPath := hdwallet.MustParseDerivationPath(derivationPath)
	enclaveAccount, err := wallet.Derive(enclaveAccountDerivationPath, true)
	if err != nil {
		log.Panicf("Derive: %v", err)
	}

	node := NewNode(prototype.NewEnclaveWithAccount(wallet, enclaveAccount))
	node.Start(l)
	log.Infof("Started node on %s", listenAddr)

	<-node.Stopped()
}
