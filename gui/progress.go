// SPDX-License-Identifier: Apache-2.0

package gui

import (
	"strings"
	"sync"
)

type ProgressBar struct {
	Title string
	// < 0 means unknown
	Status   []string
	finished bool

	mtx *sync.Mutex
}

func (p *ProgressBar) Add(status string) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if p.finished {
		return
	}
	p.Status = append(p.Status, "├ "+status)
}

func (p *ProgressBar) Finish(status string) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if p.finished {
		return
	}
	p.Status = append(p.Status, "└─"+status)
	p.finished = true
}

func (p *ProgressBar) Render() string {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	return "┌─" + p.Title + "\n" + strings.Join(p.Status, "\n")
}
