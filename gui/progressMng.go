// SPDX-License-Identifier: Apache-2.0

package gui

import (
	"fmt"
	"sync"

	"github.com/jroimartin/gocui"
)

type ProgressMng struct {
	g    *gocui.Gui
	view string
	max  int
	bars []*ProgressBar

	mtx *sync.Mutex
}

func (m *ProgressMng) Add(title string) *ProgressBar {
	// todo scroll
	m.mtx.Lock()
	defer m.mtx.Unlock()
	bar := &ProgressBar{
		Title: title,
		Width: 10,
		mtx:   &sync.Mutex{},
	}
	if len(m.bars) >= m.max {
		m.bars = m.bars[1:]
	}
	m.bars = append(m.bars, bar)
	return bar
}

func (m *ProgressMng) Render() {
	m.g.Update(func(g *gocui.Gui) error {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		v, err := m.g.View(m.view)
		if err != nil {
			panic(err)
		}
		v.Clear()
		for _, bar := range m.bars {
			fmt.Fprintf(v, "%s\n", bar.Render())
		}
		return nil
	})
}
