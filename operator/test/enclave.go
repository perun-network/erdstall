// SPDX-License-Identifier: Apache-2.0

package test

import (
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/tee"
)

type Enclave struct {
	processTXsError error
	txs             chan *tee.Transaction

	deps chan *tee.DepositProof
	bals chan *tee.BalanceProof
}

var _ tee.Enclave = (*Enclave)(nil)

func NewMockedEnclave() *Enclave {
	return &Enclave{
		// Buffer some TX to make testing easier.
		txs:  make(chan *tee.Transaction, 10),
		deps: make(chan *tee.DepositProof, 10),
		bals: make(chan *tee.BalanceProof, 10),
	}
}

func (e *Enclave) Log() *log.Entry {
	return log.WithField("role", "enclave")
}

func (e *Enclave) Init() (common.Address, []byte, error) {
	return common.Address{}, nil, nil
}

func (e *Enclave) Run(tee.Parameters) error {
	return nil
}

func (e *Enclave) ProcessBlocks(...*tee.Block) error {
	return nil
}

func (e *Enclave) ProcessTXs(txs ...*tee.Transaction) error {
	if e.processTXsError != nil {
		return e.processTXsError
	}
	for _, tx := range txs {
		e.txs <- tx
	}
	return nil
}

func (e *Enclave) DepositProofs() (ret []*tee.DepositProof, err error) {
	for {
		select {
		case proof := <-e.deps:
			ret = append(ret, proof)
		default:
			return
		}
	}
}

func (e *Enclave) BalanceProofs() (ret []*tee.BalanceProof, err error) {
	for {
		select {
		case proof := <-e.bals:
			ret = append(ret, proof)
		default:
			return
		}
	}
}

func (e *Enclave) Shutdown() {}

func (e *Enclave) SetProcessTXsError(err error) {
	e.processTXsError = err
}

func (e *Enclave) PushDepositProof(proof *tee.DepositProof) {
	e.Log().WithField("acc", proof.Balance.Account.Hex()).Debug("Produced Deposit proof")
	e.deps <- proof
}

func (e *Enclave) PushBalanceProof(proof *tee.BalanceProof) {
	e.Log().WithField("acc", proof.Balance.Account.Hex()).Debug("Produced Balance proof")
	e.bals <- proof
}

func (e *Enclave) Transactions() <-chan *tee.Transaction {
	return e.txs
}
