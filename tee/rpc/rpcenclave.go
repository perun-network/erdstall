// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"io"
	"net"
	"net/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/tee"

	"perun.network/go-perun/log"
)

var _ tee.Enclave = (*RPCEnclave)(nil)

// RPCEnclave communicates with a rpc.Server's enclave.
type RPCEnclave struct {
	client *rpc.Client
}

// NewRPCEnclave communicates with a rpc.Server via conn.
func NewRPCEnclave(conn io.ReadWriteCloser) *RPCEnclave {
	return &RPCEnclave{
		client: rpc.NewClient(conn),
	}
}

// DialEnclave dials an enclave at the given TCP/IP address.
func DialEnclave(address string) (*RPCEnclave, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return NewRPCEnclave(conn), nil
}

// Close closes the connection to the remote enclave.
func (re *RPCEnclave) Close() error {
	return re.client.Close()
}

func (re *RPCEnclave) Run(p tee.Parameters) error {
	return re.client.Call("Server.Run", p, &Void{})
}

func (re *RPCEnclave) Shutdown() {
	if err := re.client.Call("Server.Shutdown", Void{}, &Void{}); err != nil {
		log.Error(err)
	}
}

func (re *RPCEnclave) Init() (common.Address, []byte, error) {
	var res TeeInitRes
	err := re.client.Call("Server.Init", &Void{}, &res)
	return res.Addr, res.Sig, err
}

func (re *RPCEnclave) ProcessBlocks(blocks ...*tee.Block) error {
	return re.client.Call("Server.ProcessBlocks", blocks, &Void{})
}

func (re *RPCEnclave) ProcessTXs(txs ...*tee.Transaction) error {
	return re.client.Call("Server.ProcessTXs", &txs, &Void{})
}

func (re *RPCEnclave) DepositProofs() (res []*tee.DepositProof, err error) {
	err = re.client.Call("Server.DepositProofs", Void{}, &res)
	return
}

func (re *RPCEnclave) BalanceProofs() (res []*tee.BalanceProof, err error) {
	err = re.client.Call("Server.BalanceProofs", Void{}, &res)
	return
}

func (re *RPCEnclave) Stop() error {
	return re.client.Call("Server.Stop", Void{}, &Void{})
}
