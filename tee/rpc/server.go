// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"net"
	"net/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/tee"

	"perun.network/go-perun/log"
	"perun.network/go-perun/pkg/sync"
	"perun.network/go-perun/pkg/sync/atomic"
)

// Server is a server that exposes an Enclave as a slave.
type Server struct {
	enclave tee.Enclave // the slave enclave.
	running atomic.Bool // whether the server has started.
	server  *rpc.Server // accepts new connections.
	stopped sync.Closer // whether the server was commanded to stop.
}

// NewServer creates a new server which is not yet running.
func NewServer(impl tee.Enclave) *Server {
	return &Server{
		enclave: impl,
		server:  rpc.NewServer(),
	}
}

// Start starts the server and sets its listener. Call only once. Starts listening
// for new connections in the background.
func (n *Server) Start(l net.Listener) {
	if l == nil {
		log.Panic("Server.Start(): nil listener")
	}
	if !n.running.TrySet() || n.IsStopped() {
		log.Panic("Server.Start(): already running")
	}
	n.stopped.OnCloseAlways(func() { l.Close() })
	if err := n.server.Register(n); err != nil {
		log.Panicf("Server.Start() registering server: %v", err)
	}
	go n.server.Accept(l)
}

// Void indicates that a function has no arguments or no return value.
type Void struct{}

// TeeInitRes holds the result for Enclave.Init requests.
type TeeInitRes struct {
	Addr common.Address
	Sig  []byte
}

// Init wraps Enclave.Init.
func (n *Server) Init(_ Void, res *TeeInitRes) (err error) {
	res.Addr, res.Sig, err = n.enclave.Init()
	return
}

// Run wraps Enclave.Run.
func (n *Server) Run(p tee.Parameters, _ *Void) (err error) {
	return n.enclave.Run(p)
}

// ProcessBlocks wraps Enclave.ProcessBlocks.
func (n *Server) ProcessBlocks(blocks []*tee.Block, _ *Void) error {
	return n.enclave.ProcessBlocks(blocks...)
}

// ProcessTXs wraps Enclave.ProcessTX.
func (n *Server) ProcessTXs(txs *[]*tee.Transaction, _ *Void) error {
	return n.enclave.ProcessTXs(*txs...)
}

// DepositProofs wraps Enclave.DepositProofs.
func (n *Server) DepositProofs(_ Void, res *[]*tee.DepositProof) (err error) {
	*res, err = n.enclave.DepositProofs()
	return
}

// BalanceProofs wraps Enclave.BalanceProofs.
func (n *Server) BalanceProofs(_ Void, res *[]*tee.BalanceProof) (err error) {
	*res, err = n.enclave.BalanceProofs()
	return
}

// Shutdown wraps Enclave.Shutdown.
func (n *Server) Shutdown(Void, *Void) error {
	n.enclave.Shutdown()
	return nil
}

// Stop stops the server from accepting new connections.
func (n *Server) Stop(Void, *Void) error {
	return n.stopped.Close()
}

// IsStopped returns whether Stop() has been called.
func (n *Server) IsStopped() bool {
	return n.stopped.IsClosed()
}

// Stopped returns a channel that will be closed once the server is stopped.
func (n *Server) Stopped() <-chan struct{} {
	return n.stopped.Closed()
}
