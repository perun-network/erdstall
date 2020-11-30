// SPDX-License-Identifier: Apache-2.0

package gui

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/perun-network/erdstall/tee"
)

type PhaseMeter struct {
	g        *gocui.Gui
	view     string
	p        tee.Parameters
	blockNum uint64
	width    int

	mtx *sync.Mutex // for the whole object
}

func NewPhaseMeter(p tee.Parameters, g *gocui.Gui, view string) *PhaseMeter {
	return &PhaseMeter{g: g, view: view, p: p, width: 3, mtx: new(sync.Mutex)}
}

func (m *PhaseMeter) SetBlock(blockNum uint64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.blockNum = blockNum
}

func (m PhaseMeter) Draw() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	v, err := m.g.View(m.view)
	if err != nil {
		panic(err)
	}
	v.Clear()
	fmt.Fprintln(v, m.render())
}

func (m PhaseMeter) render() string {
	w := m.width * int(m.p.PhaseDuration)
	//clrs := []int{BLUE, PURPLE, CYAN}
	deposit := m.p.DepositEpoch(m.blockNum)
	depositStart := m.p.DepositStartBlock(deposit)
	//leftInEpoch := m.p.DepositEndBlock(deposit) - m.blockNum
	//tx := m.p.TxEpoch(m.blockNum)
	//info := fmt.Sprintf("Deposit epoch %d has %d block(s) left.", deposit, leftInEpoch)
	epochs := fmtEpoch(deposit, w, 0) + "\n" + fmtEpoch(deposit, w, 1) + "\n" + fmtEpoch(deposit, w, 2)
	blocks := fmtBlockNum(depositStart, m.width, m.blockNum == depositStart) +
		fmtBlockNum(depositStart+1, m.width, m.blockNum == (depositStart+1)) +
		fmtBlockNum(depositStart+2, m.width, m.blockNum == (depositStart+2))
	return epochs + "\n" + strings.Repeat(" ", 6+w*2) + blocks

	/*phase := strings.Repeat(colorBg(" ", clrs[deposit%3]), int(leftInEpoch)*m.width) +
		strings.Repeat(colorBg(" ", clrs[(deposit+1)%3]), int(m.p.PhaseDuration)*m.width) +
		strings.Repeat(colorBg(" ", clrs[(deposit+2)%3]), int(m.p.PhaseDuration)*m.width) +
		strings.Repeat(" ", int(m.p.PhaseDuration-leftInEpoch))
	return fmt.Sprintf("%s %s\nD  T  E", phase, info)*/

	/*phase := strings.Repeat(colorBg(" ", ORANGE), int(m.p.PhaseDuration-leftInEpoch)*m.width) +
		strings.Repeat(color(colorBg("#", RED), BLACK), m.width) +
		strings.Repeat(colorBg(" ", GREEN), int(leftInEpoch-1)*m.width)
	return fmt.Sprintf("%s%s%s\nDeposit           Tx                Exit", phase, phase, phase)*/
}

func fmtEpoch(num uint64, w int, offset int) string {
	clrs := []int{ORANGE, PURPLE, CYAN}
	names := []string{"D", "T", "E"}
	ret := fmt.Sprintf("%s %03d %s", colorBg(names[offset], clrs[offset]), int(num)-offset, strings.Repeat(" ", w*(2-offset)))
	for _, clr := range clrs {
		ret = ret + strings.Repeat(colorBg(" ", clr), w)
	}
	return ret
}

func fmtBlockNum(num uint64, w int, active bool) string {
	n := fmt.Sprintf("%d", num)
	if len(n) > int(w) {
		return strings.Repeat("?", int(w))
	}
	fill := strings.Repeat(" ", int(w)-len(n))
	if active {
		return color(n+fill, RED)
	}
	return n + fill
}
