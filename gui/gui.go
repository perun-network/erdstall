// SPDX-License-Identifier: Apache-2.0

package gui

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/perun-network/erdstall/client"
	"github.com/perun-network/erdstall/eth"
)

type GUI struct {
	g      *gocui.Gui
	client *client.Client
	bars   *ProgressMng
}

func RunGui(client *client.Client, clEvents chan *client.ClientEvent, chainEvents chan string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.SetManagerFunc(func(g *gocui.Gui) error { return layout(g, client) })
	gui := &GUI{g, client, &ProgressMng{g, "cmds", 3, nil, &sync.Mutex{}}}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if g != gui.g {
			panic("Wrong gui object")
		}
		return gui.enter(v)
	}); err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case e := <-chainEvents:
				gui.g.Update(func(g *gocui.Gui) error {
					if g != gui.g {
						panic("Wrong gui object")
					}
					if strings.HasPrefix(e, "new Block") {

					} else if strings.HasPrefix(e, "New Epoch") {
						gui.logChain("🍵  ", e, "\n")
					} else {
						gui.logChain("⛓  ", e, "\n")
					}

					return nil
				})
			case e := <-clEvents:
				gui.g.Update(func(g *gocui.Gui) error {
					return gui.handleClientEvent(e)
				})
			}
		}
	}()

	if err := gui.g.MainLoop(); err != nil {
		panic(err)
	}
}

// must be called in a g.Update
func (gui *GUI) handleClientEvent(e *client.ClientEvent) error {
	switch e.Type {
	case client.SET_PARAMS:
		out := gui.getView("status")
		out.Clear()
		fmt.Fprint(out, fmt.Sprintf("Contract %s TEE %s\nPowDepth %d PhaseDuration %d ResponseDuration %d\nInitBlock %d", e.Params.Contract.Hex(), e.Params.TEE.Hex(), e.Params.PowDepth, e.Params.PhaseDuration, e.Params.ResponseDuration, e.Params.InitBlock))
	case client.SET_BALANCE:
		out := gui.getView("balance")
		out.Clear()
		fmt.Fprintf(out, "You hold %s ETH.", bold(eth.WeiToEthFloat(e.Report.Balance).String()))
	case client.SET_OP_TRUST:
		out := gui.getView("status")
		out.Title = fmt.Sprintf("Operator is %s", e.OpTrust)
	case client.BENCH:
		gui.logOut(e.Result.String(), "\n")
	}
	return nil
}

func (gui *GUI) enter(v *gocui.View) error {
	input := strings.TrimRight(v.Buffer(), "\n")
	gui.logOut("$ ", input, "\n")
	if strings.HasPrefix(input, "help") {
		gui.logOut(helpText)
		v.Clear()
		v.SetCursor(0, 0)
		return nil
	} else if strings.HasPrefix(input, "credits") {
		gui.logOut(creditText)
		v.Clear()
		v.SetCursor(0, 0)
		return nil
	}

	// Responsive gui -> do not block till command is finished.
	go func() {
		if status, err := gui.eval(input); err != nil {
			gui.logOut("✗ ", color(err.Error(), RED), "\n")
		} else {
			bar := gui.bars.Add(input)
			for {
				select {
				case s := <-status:
					if s == nil {
						bar.Finish(color("Done", GREEN))
						gui.logOut("✓ ", color(input, GREEN), "\n")
						gui.bars.Render()
						return
					}
					if s.Err != nil {
						bar.Finish(color("Error", RED))
						gui.logOut("✗ ", color(s.Err.Error(), RED), "\n")
						gui.bars.Render()
						return
					} else if len(s.War) != 0 {
						bar.Add(color(s.War, ORANGE))
						gui.bars.Render()
					} else {
						bar.Add(s.Msg)
						gui.bars.Render()
					}
				}
			}
		}
	}()

	v.Clear()
	v.SetCursor(0, 0)
	return nil
}

func (gui *GUI) eval(input string) (chan *client.CmdStatus, error) {
	fs := strings.Fields(strings.TrimSpace(input))
	if len(fs) < 1 {
		return nil, errors.New("Empty input")
	}
	cmd := fs[0]
	status := make(chan *client.CmdStatus, 2)

	switch cmd {
	case "deposit":
		go gui.client.CmdDeposit(status, fs[1:]...)
	case "send":
		go gui.client.CmdSend(status, fs[1:]...)
	case "bench":
		go gui.client.CmdBench(status, fs[1:]...)
	case "leave":
		go gui.client.CmdLeave(status, fs[1:]...)
	default:
		return nil, fmt.Errorf("Unknow command: '%s'", cmd)
	}
	return status, nil
}

func layout(g *gocui.Gui, client *client.Client) error {
	maxW, maxH := g.Size()

	outW := (maxW / 3) * 2 // TODO golden ratio
	if v, err := g.SetView("status", 0, 0, outW, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Status"
	}
	if v, err := g.SetView("out", 0, 4, outW, maxH-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Log"
		v.Autoscroll = true
	}
	chainH := (maxH / 3) * 2
	if v, err := g.SetView("balance", outW+1, 0, maxW-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Balance"
	}
	if v, err := g.SetView("chain", outW+1, 4, maxW-1, chainH); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = fmt.Sprintf("PK: %s @ %s", client.Address().Hex(), client.Config.ChainURL)
		v.Autoscroll = true
	}

	if v, err := g.SetView("cmds", outW+1, chainH+1, maxW-1, maxH-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Commands"
		v.Autoscroll = true
	}

	if v, err := g.SetView("input", 0, maxH-3, maxW-1, maxH-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}

		v.Editor = gocui.DefaultEditor
		v.Editable = true
		g.Cursor = true
	}
	return nil
}

func (gui *GUI) logOut(v ...interface{}) {
	gui.log("out", v...)
}

func (gui *GUI) logChain(v ...interface{}) {
	gui.log("chain", v...)
}

func (gui *GUI) log(where string, v ...interface{}) {
	out, err := gui.g.View(where)
	if err != nil {
		panic(err)
	}
	if _, err := fmt.Fprint(out, v...); err != nil {
		panic(err)
	}
}

func (gui *GUI) getView(name string) *gocui.View {
	out, err := gui.g.View(name)
	if err != nil {
		panic(err)
	}
	return out
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// ANSI Escape colors.
const (
	RED    = 31
	GREEN  = 32
	ORANGE = 33
)

func bold(msg string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", msg)
}

func color(msg string, clr int) string {
	return fmt.Sprintf("\033[%d;%dm%s\033[0m", clr, 1, msg)
}

const helpText = `Erdstall client CLI
Available commands:
 help
   Prints this page.
 bench <runs>
   Runs a benchmark of sending <runs> off-chain TX.
 deposit <amount>
   Deposits <amount> into the Erdstall contract.
 send <receiver> <amount>
   Sends <amount> to <receiver>.
 leave
   Withdraws all funds and exits the network.
 version
`

const creditText = `Developed by Perun Network for the 2020 ETH Hackathon.

  Norbert Dzikowski  - Enclave, Intel SGX
  Matthias Geihs	 - Operator Node, Presantation
  Steffen Rattay	 - GrapheneOS, Presantation
  Sebastian Stammler - Enclave, Contract, Speaker
  Oliver Tale-Yazdi  - Client, GUI

Checkout our other projects at https://perun.network
`
