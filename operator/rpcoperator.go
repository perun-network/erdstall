// SPDX-License-Identifier: Apache-2.0

package operator

import (
	"math/big"
	"sync"

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
		SubscribeProofs(common.Address) (ProofSub, error)
	}

	// RPCOperator will be exposed to the client as a websocket RPC Server.
	// It implements `WireAPI`.
	RPCOperator struct {
		mtx     sync.Mutex // protects all
		enclave tee.Enclave

		subs map[common.Address]*BufferedProofSubs
	}

	// BufferedProofSub describes a slice of subs, which also holds the latest
	// available deposit- and/or balance-proof.
	BufferedProofSubs struct {
		subs     []ProofSub
		latestDP *tee.DepositProof
		latestBP *tee.BalanceProof
	}

	// ProofSub is a subscription on TEE proofs.
	ProofSub struct {
		deposits chan tee.DepositProof
		balances chan tee.BalanceProof
		quit     chan struct{}
	}
)

var _ WireAPI = (*RPCOperator)(nil)

func NewRPCOperator(enclave tee.Enclave) *RPCOperator {
	return &RPCOperator{
		enclave: enclave,
		subs:    make(map[common.Address]*BufferedProofSubs),
	}
}

func (o *RPCOperator) Send(tx tee.Transaction) error {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	log.Infof("Sending %d WEI 0x%s…->0x%s…", (*big.Int)(tx.Amount).Uint64(), tx.Sender.Hex()[:5], tx.Recipient.Hex()[:5])
	return o.enclave.ProcessTXs(&tx)
}

// SubscribeProofs returns a subscription on TEE proofs for the given address.
// The subscription buffers the most recent proof until the client retrieves it.
func (o *RPCOperator) SubscribeProofs(addr common.Address) (ProofSub, error) {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	log.WithField("who", addr.Hex()).Debug("Subscribed to proofs")
	sub := o.subscribe(addr)
	return sub, nil
}

// subscribe is an internal implementation detail and should not be called.
func (o *RPCOperator) subscribe(addr common.Address) ProofSub {
	sub := *newProofSub()
	if _, ok := o.subs[addr]; !ok {
		o.subs[addr] = new(BufferedProofSubs)
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
		bufsub = new(BufferedProofSubs)
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
		bufsub = new(BufferedProofSubs)
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
func newProofSub() *ProofSub {
	return &ProofSub{
		deposits: make(chan tee.DepositProof, 1),
		balances: make(chan tee.BalanceProof, 1),
		quit:     make(chan struct{}),
	}
}

func (sub ProofSub) Deposits() <-chan tee.DepositProof {
	return sub.deposits
}

func (sub ProofSub) Balances() <-chan tee.BalanceProof {
	return sub.balances
}

func (sub ProofSub) Closed() <-chan struct{} {
	return sub.quit
}

func (sub ProofSub) Unsubscribe() {
	select {
	case <-sub.quit:
	default:
		close(sub.quit)
	}
}
