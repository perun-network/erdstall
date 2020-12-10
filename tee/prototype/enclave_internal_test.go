// SPDX-License-Identifier: Apache-2.0

package prototype_test

import (
	"testing"

	ptest "perun.network/go-perun/pkg/test"

	"github.com/perun-network/erdstall/eth"
	. "github.com/perun-network/erdstall/tee/prototype"
	ttest "github.com/perun-network/erdstall/tee/test"
)

func TestEnclave(t *testing.T) {
	rng := ptest.Prng(t)
	encWallet := eth.NewHdWallet(rng)
	ttest.GenericEnclaveTest(t, NewEnclave(encWallet))
}
