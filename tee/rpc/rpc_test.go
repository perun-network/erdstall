// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
	"github.com/perun-network/erdstall/tee/prototype"
	ttest "github.com/perun-network/erdstall/tee/test"

	ctxtest "perun.network/go-perun/pkg/context/test"
	pkgtest "perun.network/go-perun/pkg/test"
)

var _ tee.Enclave = (*mockEnclave)(nil)

type mockEnclave struct{}

func (*mockEnclave) Init() (_ common.Address, _ []byte, _ error)     { return }
func (*mockEnclave) Run(tee.Parameters) (_ error)                    { return }
func (*mockEnclave) Shutdown()                                       {}
func (*mockEnclave) ProcessBlocks(...*tee.Block) (_ error)           { return }
func (*mockEnclave) ProcessTXs(...*tee.Transaction) (_ error)        { return }
func (*mockEnclave) DepositProofs() (_ []*tee.DepositProof, _ error) { return }
func (*mockEnclave) BalanceProofs() (_ []*tee.BalanceProof, _ error) { return }

var _ net.Listener = (*mockListener)(nil)

type mockListener struct{ conn chan net.Conn }

func (l *mockListener) Accept() (net.Conn, error) {
	c := <-l.conn
	if c == nil {
		return nil, errors.New("closed")
	}
	return c, nil
}

func (l *mockListener) Close() (_ error) { close(l.conn); return }
func (*mockListener) Addr() net.Addr     { panic(nil) }

func (l *mockListener) dial() (_ net.Conn, err error) {
	a, b := net.Pipe()
	err = errors.New("closed")
	func() {
		defer func() { _ = recover() }()
		l.conn <- b
		err = nil
	}()
	return a, err
}

func newMockListener() *mockListener {
	return &mockListener{conn: make(chan net.Conn)}
}

func TestConn(t *testing.T) {
	l := newMockListener()

	node := NewServer(&mockEnclave{})
	node.Start(l)

	conn, err := l.dial()
	require.NoError(t, err)
	enc := NewRPCEnclave(conn)

	ctxtest.AssertTerminates(t, 10*time.Second, func() {
		_, _, err = enc.Init()
		assert.NoError(t, err)
		assert.NoError(t, enc.Run(tee.Parameters{}))
		assert.NoError(t, enc.ProcessBlocks())
		assert.NoError(t, enc.ProcessTXs())
		_, err = enc.DepositProofs()
		assert.NoError(t, err)
		_, err = enc.BalanceProofs()
		assert.NoError(t, err)
		enc.Shutdown()
		assert.NoError(t, enc.Stop())
	})
}

func TestRPCEnclave(t *testing.T) {
	rng := pkgtest.Prng(t)
	encWallet := eth.NewHdWallet(rng)
	l := newMockListener()

	node := NewServer(prototype.NewEnclave(encWallet))
	node.Start(l)

	conn, err := l.dial()
	require.NoError(t, err)
	cfg := ttest.EnclaveTestCfg{IsCachingEnclave: false}
	ttest.GenericEnclaveTest(t, NewRPCEnclave(conn), cfg)
}
