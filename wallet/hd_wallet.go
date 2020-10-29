// SPDX-License-Identifier: Apache-2.0

package wallet

import (
	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/perun-network/erdstall/tee"

	phd "perun.network/go-perun/backend/ethereum/wallet/hd"
)

// Wallet stores a hdwallet and an unlocked account which will be used for signing.
type Wallet struct {
	Wallet *phd.Wallet
	Acc    *phd.Account
}

var _ tee.TextSigner = &Wallet{}

// Sig is always 65 bytes long.
type Sig = []byte

// NewWallet creates a wallet from a `mnemonic` string.
// The secrete key of the account is deterministically derived from
// the `mnemonic`.
func NewWallet(mnemonic string, accountIndex uint) *Wallet {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		panic(err)
	}

	perunWallet, err := phd.NewWallet(wallet, "m/44'/60'/0'/0/0", accountIndex)
	if err != nil {
		panic(err)
	}
	acc, err := perunWallet.NewAccount()
	if err != nil {
		panic(err)
	}

	return &Wallet{
		Wallet: perunWallet,
		Acc:    acc,
	}
}

// SignText implements the tee.TextSigner interface.
func (w *Wallet) SignText(acc accounts.Account, text []byte) (Sig, error) {
	sig, err := w.Wallet.Wallet().SignText(acc, text)
	if sig[64] != 0 && sig[64] != 1 {
		panic("Invalid v")
	}
	return sig, err
}
