package operator

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/tee"
)

// RemoteEnclave provides the enclave interface to the user.
type RemoteEnclave struct {
	operator *Operator
}

func newRemoteEnclave(operator *Operator) *RemoteEnclave {
	return &RemoteEnclave{operator: operator}
}

// AddTransaction adds a transaction to the enclave's transaction pool.
func (e *RemoteEnclave) AddTransaction(t tee.Transaction, _ *struct{}) error {
	return e.operator.enclave.ProcessTXs(&t)
}

// GetDepositProof returns the next deposit proof for the given user.
func (e *RemoteEnclave) GetDepositProof(user common.Address, _dp *tee.DepositProof) error {
	dp, ok := e.operator.depositProofs.Get(user)
	if !ok {
		return errors.New("deposit proof not available")
	}

	*_dp = *dp

	return nil
}

// GetBalanceProof returns the next balance proof for the given user.
func (e *RemoteEnclave) GetBalanceProof(user common.Address, _bp *tee.BalanceProof) error {
	bp, ok := e.operator.balanceProofs.Get(user)
	if !ok {
		return errors.New("balance proof not available")
	}

	*_bp = *bp

	return nil
}
