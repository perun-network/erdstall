// SPDX-License-Identifier: Apache-2.0

package operator

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/tee"
)

type (
	// WireAPI contains all functions that a client is interested in.
	// It is called `WireAPI` since the client will interact with an
	// implementation of it over the wire.
	WireAPI interface {
		Send(tee.Transaction) error
		SubscribeProofs(common.Address) (ClientSub, error)
	}

	// RPCOperator will be exposed to the client as a websocket RPC Server.
	// It implements `WireAPI`.
	RPCOperator struct {
		mtx     sync.Mutex // protects all
		enclave tee.Enclave

		txReceipts *txReceipts
		subs       map[common.Address]*BufferedClientSubs
	}

	// BufferedClientSub describes a slice of subs, which also holds the latest
	// available deposit- and/or balance-proof.
	BufferedClientSubs struct {
		subs     []ClientSub
		latestDP *tee.DepositProof
		latestBP *tee.BalanceProof
	}

	// ClientSub is a subscription on TEE proofs.
	ClientSub struct {
		deposits chan tee.DepositProof
		balances chan tee.BalanceProof
		receipts chan tee.Transaction
		quit     chan struct{}
	}
)

var _ WireAPI = (*RPCOperator)(nil)

var txReceiptDeliveryTimeout = 20 * time.Second

func NewRPCOperator(enclave tee.Enclave, txReceipts *txReceipts) *RPCOperator {
	if txReceipts == nil {
		txReceipts = newTXReceipts()
	}
	return &RPCOperator{
		enclave:    enclave,
		subs:       make(map[common.Address]*BufferedClientSubs),
		txReceipts: txReceipts,
	}
}

func (o *RPCOperator) Send(tx tee.Transaction) error {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	log.Infof("Sending %d WEI 0x%s…->0x%s…", (*big.Int)(tx.Amount).Uint64(), tx.Sender.Hex()[:5], tx.Recipient.Hex()[:5])
	if err := o.enclave.ProcessTXs(&tx); err != nil {
		return err
	}
	if bsub, ok := o.subs[tx.Recipient]; ok {
		timeout := time.After(txReceiptDeliveryTimeout)
		for _, sub := range bsub.subs {
			sub := sub
			go func() {
				select {
				case sub.receipts <- tx:
				case <-timeout:
				}
			}()
		}
	}
	return nil
}

// SubscribeProofs returns a subscription on TEE proofs for the given address.
// The subscription buffers the most recent proof until the client retrieves it.
func (o *RPCOperator) SubscribeProofs(addr common.Address) (ClientSub, error) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	log.WithField("who", addr.Hex()).Debug("Subscribed to proofs")
	sub := o.subscribe(addr)
	return sub, nil
}

// subscribe is an internal implementation detail and should not be called.
func (o *RPCOperator) subscribe(addr common.Address) ClientSub {
	sub := *newProofSub(o.txReceipts.AddPeer(addr))
	if _, ok := o.subs[addr]; !ok {
		o.subs[addr] = new(BufferedClientSubs)
	} else {
		if dp := o.subs[addr].latestDP; dp != nil {
			sub.deposits <- *dp
		}
		if bp := o.subs[addr].latestBP; bp != nil {
			sub.balances <- *bp
		}
	}
	o.subs[addr].subs = append(o.subs[addr].subs, sub)
	return sub
}

func (o *RPCOperator) PushDepositProof(proof tee.DepositProof) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	who := proof.Balance.Account
	bufsub, ok := o.subs[who]
	if !ok {
		log.WithField("who", who.Hex()).Debug("Received DP without subscription - buffering")
		bufsub = new(BufferedClientSubs)
		o.subs[who] = bufsub
	}
	bufsub.latestDP = &proof
	subs := bufsub.subs

	for i := 0; i < len(subs); i++ {
		// Clear out old buffered proof.
		select {
		case <-subs[i].deposits:
		default:
		}
		// Write new proof in.
		select {
		case <-subs[i].quit:
			l := len(subs) - 1
			subs[i] = subs[l]
			subs = subs[:l]
			i--
		case subs[i].deposits <- proof:
		}
	}
	bufsub.subs = subs
}

func (o *RPCOperator) PushBalanceProof(proof tee.BalanceProof) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	who := proof.Balance.Account
	bufsub, ok := o.subs[who]
	if !ok {
		log.WithField("who", who.Hex()).Debug("Received BP without subscription - buffering")
		bufsub = new(BufferedClientSubs)
		o.subs[who] = bufsub
	}
	bufsub.latestBP = &proof
	subs := bufsub.subs

	for i := 0; i < len(subs); i++ {
		// Clear out old buffered proof.
		select {
		case <-subs[i].balances:
		default:
		}
		// Write new proof in.
		select {
		case <-subs[i].quit:
			l := len(subs) - 1
			subs[i] = subs[l]
			subs = subs[:l]
			i--
		case subs[i].balances <- proof:
		}
	}
	bufsub.subs = subs
}

// newProofSub returns a new proofSub. The proof channels have buffer size 1.
func newProofSub(receipts chan tee.Transaction) *ClientSub {
	return &ClientSub{
		deposits: make(chan tee.DepositProof, 1),
		balances: make(chan tee.BalanceProof, 1),
		receipts: receipts,
		quit:     make(chan struct{}),
	}
}

func (sub ClientSub) Deposits() <-chan tee.DepositProof {
	return sub.deposits
}

func (sub ClientSub) Balances() <-chan tee.BalanceProof {
	return sub.balances
}

func (sub ClientSub) Receipts() <-chan tee.Transaction {
	return sub.receipts
}

func (sub ClientSub) Closed() <-chan struct{} {
	return sub.quit
}

func (sub ClientSub) Unsubscribe() {
	select {
	case <-sub.quit:
	default:
		close(sub.quit)
	}
}
