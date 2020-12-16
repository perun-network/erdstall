// SPDX-License-Identifier: Apache-2.0

package caching_test

import (
	"testing"

	ptest "perun.network/go-perun/pkg/test"

	"github.com/perun-network/erdstall/eth"
	caching "github.com/perun-network/erdstall/tee/caching"
	proto "github.com/perun-network/erdstall/tee/prototype"
	ttest "github.com/perun-network/erdstall/tee/test"
)

func TestEnclave(t *testing.T) {
	rng := ptest.Prng(t)
	encWallet := eth.NewHdWallet(rng)
	enc := caching.NewEnclave(proto.NewEnclave(encWallet))
	cfg := ttest.EnclaveTestCfg{IsCachingEnclave: true}
	ttest.GenericEnclaveTest(t, enc, cfg)
}
