// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/tee"
)

// generateBalanceProofs creates balance proofs for an epoch outcome.
func (e *Enclave) generateBalanceProofs(o Outcome) []*tee.BalanceProof {
	proofs := make([]*tee.BalanceProof, 0, len(o.Accounts))
	for addr, bal := range o.Accounts {
		proofs = append(proofs, e.signBalanceProof(tee.Balance{
			Epoch:   o.TxEpoch,
			Account: addr,
			Value:   (*tee.Amount)(bal.Value),
		}))
	}
	return proofs
}

// generateDepositProofs generates deposit proofs for the given deposit events.
func (e *Enclave) generateDepositProofs(
	deposit ...*erdstallDepEvent,
) []*tee.DepositProof {
	proofs := make([]*tee.DepositProof, len(deposit))
	for i, d := range deposit {
		proofs[i] = e.signDepositProof(tee.Balance{
			Epoch:   d.Epoch,
			Account: d.Account,
			Value:   (*tee.Amount)(new(big.Int).Set(d.Value)),
		})
	}
	return proofs
}

func (e *Enclave) signBalanceProof(b tee.Balance) *tee.BalanceProof {
	msg, err := tee.EncodeBalanceProof(e.params.Contract, b)
	if err != nil {
		log.WithError(err).Panic("encoding balance proof")
	}

	sig, err := e.wallet.SignText(*e.account, crypto.Keccak256(msg))
	if err != nil {
		log.WithError(err).Panic("Signing balance proof")
	}
	sig[64] += 27
	return &tee.BalanceProof{
		Balance: b,
		Sig:     sig,
	}
}

func (e *Enclave) signDepositProof(b tee.Balance) *tee.DepositProof {
	msg, err := tee.EncodeDepositProof(e.Params.Contract, b)
	if err != nil {
		log.WithError(err).Panic("encoding deposit proof")
	}

	sig, err := e.wallet.SignText(*e.account, crypto.Keccak256(msg))
	if err != nil {
		log.WithError(err).Panic("Signing deposit proof")
	}
	sig[64] += 27

	return &tee.DepositProof{
		Balance: b,
		Sig:     sig,
	}
}
