package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

// DefaultTransactor can be used to make TransactOpts for accounts stored in a HD wallet.
type DefaultTransactor struct {
	wallet accounts.Wallet
}

// NewTransactor returns a TransactOpts for the given account. It errors if the account is
// not contained in the wallet used for initializing transactor backend.
func (t *DefaultTransactor) NewTransactor(account accounts.Account) (*bind.TransactOpts, error) {
	if !t.wallet.Contains(account) {
		return nil, errors.New("account not found in wallet")
	}
	return &bind.TransactOpts{
		From: account.Address,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != account.Address {
				return nil, errors.New("not authorized to sign this account")
			}
			// Last parameter (chainID) is only relevant when making EIP155 compliant signatures.
			// Since we use only non EIP155 signatures, set this to zero value.
			// For more details, see here (https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md).
			return t.wallet.SignTx(account, tx, big.NewInt(0))
		},
	}, nil
}

// NewDefaultTransactor returns a backend that can make TransactOpts for accounts contained in the given HD wallet.
func NewDefaultTransactor(w accounts.Wallet) *DefaultTransactor {
	return &DefaultTransactor{wallet: w}
}
