// SPDX-License-Identifier: Apache-2.0

package rpc

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/tee"

	"perun.network/go-perun/log"
)

var _ tee.Enclave = (*RemoteEnclave)(nil)

// RemoteEnclave communicates with a rpc.Node's enclave.
type RemoteEnclave struct {
	client *rpc.Client
}

// NewRemoteEnclave communicates with a rpc.Node via conn.
func NewRemoteEnclave(conn io.ReadWriteCloser) *RemoteEnclave {
	return &RemoteEnclave{
		client: rpc.NewClient(conn),
	}
}

// Dials an enclave three times to allow three concurrent accesses to the enclave.
func DialEnclave(address string) (*RemoteEnclave, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return NewRemoteEnclave(conn), nil
}

// Close closes the connection to the remote enclave.
func (re *RemoteEnclave) Close() error {
	return re.client.Close()
}

func (re *RemoteEnclave) Run(p tee.Parameters) error {
	return re.client.Call("Node.Run", p, &Void{})
}

func (re *RemoteEnclave) Shutdown() {
	if err := re.client.Call("Node.Shutdown", &Void{}, &Void{}); err != nil {
		log.Error(err)
	}
}

func (re *RemoteEnclave) Init() (common.Address, []byte, error) {
	var res TeeInitRes
	return res.Addr, res.Sig, re.client.Call("Node.Init", &Void{}, &res)
}

func (re *RemoteEnclave) ProcessBlocks(blocks ...*tee.Block) error {
	encodedBlocks := make([][]byte, len(blocks))
	for i, b := range blocks {
		var buf bytes.Buffer
		if err := b.EncodeRLP(&buf); err != nil {
			return fmt.Errorf("encoding block: %w", err)
		}
		encodedBlocks[i] = buf.Bytes()
	}
	return re.client.Call("Node.ProcessBlocks", encodedBlocks, &Void{})
}

func (re *RemoteEnclave) ProcessTXs(txs ...*tee.Transaction) error {
	return re.client.Call("Node.ProcessTXs", &txs, &Void{})
}

func (re *RemoteEnclave) DepositProofs() (res []*tee.DepositProof, err error) {
	err = re.client.Call("Node.DepositProofs", &Void{}, &res)
	return
}

func (re *RemoteEnclave) BalanceProofs() (res []*tee.BalanceProof, err error) {
	err = re.client.Call("Node.BalanceProofs", &Void{}, &res)
	return
}

func (re *RemoteEnclave) Stop() error {
	return re.client.Call("Node.Stop", &Void{}, &Void{})
}
