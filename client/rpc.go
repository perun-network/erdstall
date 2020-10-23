// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"net/rpc"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/tee"
	"github.com/pkg/errors"
)

// RPC warps go rpc calls to connect the client with the operator.
type RPC struct {
	conn *rpc.Client
}

func NewRPC(host string, port uint16) *RPC {
	conn, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}
	return &RPC{conn}
}

func (r *RPC) AddTX(ctx context.Context, tx tee.Transaction) error {
	var _reply struct{}
	return r.call(ctx, "RemoteEnclave.AddTransaction", tx, &_reply)
}

func (r *RPC) GetDepositProof(ctx context.Context, epoch uint64, user common.Address) (reply tee.DepositProof, err error) {
	for i := 0; ; i++ {
		if e := r.call(ctx, "RemoteEnclave.GetDepositProof", user, &reply); e != nil {
			return tee.DepositProof{}, e
		}
		if reply.Balance.Epoch < epoch {
			time.Sleep(time.Second * 1)
			continue
		} else if reply.Balance.Epoch > epoch {
			return tee.DepositProof{}, errors.New("Newer epoch")
		}
		return
	}
}

func (r *RPC) GetBalanceProof(ctx context.Context, user common.Address) (reply tee.BalanceProof, err error) {
	return reply, r.call(ctx, "RemoteEnclave.GetBalanceProof", user, &reply)
}

// call polls a function until it succeedes or ctx is cancelled.
func (r *RPC) call(ctx context.Context, f string, arg interface{}, reply interface{}) (err error) {
	for i := 0; ; i++ {
		e := make(chan error)
		go func() { e <- r.conn.Call(f, arg, reply) }()

		select {
		case e := <-e:
			if e == nil {
				return nil
			}
			err = fmt.Errorf("Retry %d: %s", i, errors.Cause(e))
			time.Sleep(time.Second * 1)
		case <-ctx.Done():
			if err != nil {
				return err
			}
			return ctx.Err()
		}
	}
}
