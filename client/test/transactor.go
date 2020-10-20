// SPDX-License-Identifier: Apache-2.0

package test

import "github.com/perun-network/erdstall/tee"

type Transactor interface {
	Send(*tee.Transaction) error
}

var _ Transactor = (*EnclaveTransactor)(nil)

type EnclaveTransactor struct {
	Enclave tee.Enclave
}

func (et *EnclaveTransactor) Send(tx *tee.Transaction) error {
	return et.Enclave.ProcessTXs(tx)
}
