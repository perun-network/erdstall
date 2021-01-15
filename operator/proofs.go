package operator

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/tee"
)

type depositProofs struct {
	mu      sync.RWMutex
	entries map[common.Address]*tee.DepositProof
}

func newDepositProofs() *depositProofs {
	return &depositProofs{entries: make(map[common.Address]*tee.DepositProof)}
}

// Get gets the deposit proof for the given user, threadsafe.
func (dps *depositProofs) Get(user common.Address) (*tee.DepositProof, bool) {
	dps.mu.RLock()
	defer dps.mu.RUnlock()

	dp, ok := dps.entries[user]

	return dp, ok
}

// AddAll adds the given deposit proofs, threadsafe.
func (dps *depositProofs) AddAll(in []*tee.DepositProof) {
	dps.mu.Lock()
	defer dps.mu.Unlock()

	for _, dp := range in {
		dps.entries[dp.Balance.Account] = dp
	}
}

type balanceProofs struct {
	mu      sync.RWMutex
	entries map[common.Address]*tee.BalanceProof
}

func newBalanceProofs() *balanceProofs {
	return &balanceProofs{entries: make(map[common.Address]*tee.BalanceProof)}
}

// Get gets the balance proof for the given user, threadsafe.
func (bps *balanceProofs) Get(user common.Address) (*tee.BalanceProof, bool) {
	bps.mu.RLock()
	defer bps.mu.RUnlock()

	bp, ok := bps.entries[user]

	return bp, ok
}

// AddAll adds the given balance proofs, threadsafe.
func (bps *balanceProofs) AddAll(in []*tee.BalanceProof) {
	bps.mu.Lock()
	defer bps.mu.Unlock()

	for _, bp := range in {
		bps.entries[bp.Balance.Account] = bp
	}
}

type txReceipts struct {
	mu      sync.Mutex
	entries map[common.Address][]chan tee.Transaction
	closed  chan struct{}
}

func newTXReceipts() *txReceipts {
	return &txReceipts{
		entries: make(map[common.Address][]chan tee.Transaction),
		closed:  make(chan struct{}),
	}
}

func (trs *txReceipts) AddPeer(addr common.Address) chan tee.Transaction {
	trs.mu.Lock()
	defer trs.mu.Unlock()

	ch := make(chan tee.Transaction, 1)
	trs.entries[addr] = append(trs.entries[addr], ch)
	return ch
}
