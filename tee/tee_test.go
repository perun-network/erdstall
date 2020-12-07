// SPDX-License-Identifier: Apache-2.0

package tee_test

import (
	"testing"

	pkgtest "perun.network/go-perun/pkg/test"

	"github.com/perun-network/erdstall/tee"
	"github.com/perun-network/erdstall/tee/test"
	wiretest "github.com/perun-network/erdstall/wire/test"
)

func TestBalance_Json(t *testing.T) {
	rng := pkgtest.Prng(t)
	bal := test.RandomBalance(rng)

	wiretest.GenericJSONMarshallingTest(t, bal, &tee.Balance{})
}

func TestBalanceProof_Json(t *testing.T) {
	rng := pkgtest.Prng(t)
	proof := test.RandomBP(rng)

	wiretest.GenericJSONMarshallingTest(t, *proof, &tee.BalanceProof{})
}

func TestDepositProof_Json(t *testing.T) {
	rng := pkgtest.Prng(t)
	proof := test.RandomDP(rng)

	wiretest.GenericJSONMarshallingTest(t, *proof, &tee.DepositProof{})
}
