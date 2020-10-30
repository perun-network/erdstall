// SPDX-License-Identifier: Apache-2.0

package gui

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/perun-network/erdstall/client"
	"github.com/perun-network/erdstall/eth"
)

type BalanceMeter struct {
	g    *gocui.Gui
	view string

	exitPossible *client.EpochBalance
	balance      *big.Int // wei

	mtx *sync.Mutex // for the whole object
}

func NewBalanceMeter(g *gocui.Gui, view string) *BalanceMeter {
	return &BalanceMeter{g: g, view: view, balance: big.NewInt(0), mtx: new(sync.Mutex)}
}

func (m *BalanceMeter) SetExitPossible(value *client.EpochBalance) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.exitPossible = value
}

func (m *BalanceMeter) SetBalance(wei *big.Int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.balance = wei
}

func (m BalanceMeter) Draw() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	v, err := m.g.View(m.view)
	if err != nil {
		panic(err)
	}
	v.Clear()
	fmt.Fprintln(v, m.render())
}

func (m BalanceMeter) render() string {
	str := fmt.Sprintf("You hold %s ETH and can ", bold(eth.WeiToEthFloat(m.balance).String()))

	if m.exitPossible == nil {
		return str + color("not withdraw", RED) + "."
	}
	with := m.exitPossible.Bal.Balance.Value
	if with.Cmp(m.balance) == 0 {
		return str + "withdraw " + bold("everything") + "."
	}
	return str + fmt.Sprintf("withdraw %s ETH.", bold(eth.WeiToEthFloat(with).String()))
}
