// SPDX-License-Identifier: Apache-2.0

package prototype

import (
	"fmt"
	"sync"

	"github.com/perun-network/erdstall/tee"
)

type epochchain struct {
	sync.RWMutex
	offset tee.Epoch
	epochs []*Epoch
}

func (ec *epochchain) Head() *Epoch {
	if len(ec.epochs) == 0 {
		return nil
	}

	return ec.epochs[len(ec.epochs)-1]
}

func (ec *epochchain) Push(e *Epoch) {
	ec.Lock()
	defer ec.Unlock()

	if len(ec.epochs) == 0 {
		ec.offset = e.Number
		ec.epochs = []*Epoch{e}
		return
	}

	headNum := ec.Head().Number
	if e.Number != headNum+1 {
		panic(fmt.Sprintf("Push: non-consecutive Epoch (got: %d, head %d)", e.Number, headNum))
	}

	ec.epochs = append(ec.epochs, e)
}

func (ec *epochchain) PruneUntil(n tee.Epoch) {
	if ec.offset > n {
		return
	}

	diff := n - ec.offset
	// GC
	for i := uint64(0); i < diff; i++ {
		ec.epochs[i] = nil
	}
	ec.epochs = ec.epochs[diff : len(ec.epochs)-1]
	ec.offset = n
}
