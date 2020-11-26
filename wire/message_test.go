// SPDX-License-Identifier: Apache-2.0

package wire_test

import (
	"testing"

	pkgtest "perun.network/go-perun/pkg/test"

	"github.com/perun-network/erdstall/tee"
	ttest "github.com/perun-network/erdstall/tee/test"
	"github.com/perun-network/erdstall/wire"
	"github.com/perun-network/erdstall/wire/test"
)

func TestJson_Marshalling(t *testing.T) {
	var id = wire.ID("an-id")
	rng := pkgtest.Prng(t)

	t.Run("SendTx", func(t *testing.T) {
		tx := tee.Transaction{
			Sig: ttest.RandomSig(rng),
		}
		obj := wire.NewSendTx(id, tx)
		test.GenericJSONMarshallingTest(t, *obj, &wire.SendTx{})
	})

	t.Run("Subscribe", func(t *testing.T) {
		addr := ttest.NewRandomAddress(rng)
		obj := wire.NewSubscribe(id, addr)
		test.GenericJSONMarshallingTest(t, *obj, &wire.Subscribe{})
	})

	t.Run("CallResult", func(t *testing.T) {
		obj := wire.Result{
			ID:    id,
			Error: "could not get proof",
		}
		test.GenericJSONMarshallingTest(t, obj, &wire.Result{})
	})

	t.Run("TopicResult", func(t *testing.T) {
		obj := wire.Result{
			Topic: wire.BalanceProofs,
			Error: "could not setup sub",
		}
		test.GenericJSONMarshallingTest(t, obj, &wire.Result{})
	})
}
