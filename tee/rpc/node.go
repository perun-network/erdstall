// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"bytes"
	"fmt"
	"net"
	"net/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/perun-network/erdstall/tee"

	"perun.network/go-perun/log"
	"perun.network/go-perun/pkg/sync"
	"perun.network/go-perun/pkg/sync/atomic"
)

// Node is a server that exposes an Enclave as a slave.
type Node struct {
	enclave tee.Enclave // the slave enclave.
	running atomic.Bool // whether the server has started.
	server  *rpc.Server // accepts new connections.
	stopped sync.Closer // whether the server was commanded to stop.
}

// NewNode creates a new server which is not yet running.
func NewNode(impl tee.Enclave) *Node {
	return &Node{
		enclave: impl,
		server:  rpc.NewServer(),
	}
}

// Start starts the server and sets its listener. Call only once. Starts listening
// for new connections in the background.
func (n *Node) Start(l net.Listener) {
	if l == nil {
		log.Panic("Node.Start(): nil listener")
	}
	if !n.running.TrySet() || n.IsStopped() {
		log.Panic("Node.Start(): already running")
	}
	n.stopped.OnCloseAlways(func() { l.Close() })
	if err := n.server.Register(n); err != nil {
		log.Panicf("Node.Start() registering server: %v", err)
	}
	go n.server.Accept(l)
}

// Void indicates that a function has no arguments o no return value.
type Void struct{}

// TeeInitRes holds the result for Enclave.Init requests.
type TeeInitRes struct {
	Addr common.Address
	Sig  []byte
}

// Init wraps Enclave.Init.
func (n *Node) Init(_ Void, res *TeeInitRes) (err error) {
	res.Addr, res.Sig, err = n.enclave.Init()
	return
}

// Run wraps Enclave.Run.
func (n *Node) Run(p tee.Parameters, _ *Void) (err error) {
	return n.enclave.Run(p)
}

// ProcessBlocks wraps Enclave.ProcessBlocks.
func (n *Node) ProcessBlocks(encodedBlocks [][]byte, _ *Void) error {
	blocks := make([]*tee.Block, len(encodedBlocks))
	for i, be := range encodedBlocks {
		stream := rlp.NewStream(bytes.NewReader(be), 0)
		var b tee.Block
		if err := b.DecodeRLP(stream); err != nil {
			return fmt.Errorf("decoding block: %w", err)
		}
		blocks[i] = &b
	}
	return n.enclave.ProcessBlocks(blocks...)
}

// ProcessTXs wraps Enclave.ProcessTX.
func (n *Node) ProcessTXs(txs *[]*tee.Transaction, _ *Void) error {
	return n.enclave.ProcessTXs(*txs...)
}

// DepositProofs wraps Enclave.DepositProofs.
func (n *Node) DepositProofs(_ Void, res *[]*tee.DepositProof) (err error) {
	*res, err = n.enclave.DepositProofs()
	return
}

// BalanceProofs wraps Enclave.BalanceProofs.
func (n *Node) BalanceProofs(_ Void, res *[]*tee.BalanceProof) (err error) {
	*res, err = n.enclave.BalanceProofs()
	return
}

// Shutdown wraps Enclave.Shutdown.
func (n *Node) Shutdown(Void, *Void) error {
	n.enclave.Shutdown()
	return nil
}

// Stop stops the server from accepting new connections.
func (n *Node) Stop(Void, *Void) error {
	return n.stopped.Close()
}

// IsStopped returns whether Stop() has been called.
func (n *Node) IsStopped() bool {
	return n.stopped.IsClosed()
}

// Stopped returns a channel that will be closed once the node is stopped.
func (n *Node) Stopped() <-chan struct{} {
	return n.stopped.Closed()
}
