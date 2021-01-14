// SPDX-License-Identifier: Apache-2.0

package gui

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/perun-network/erdstall/client"
)

type GUI struct {
	g       *gocui.Gui
	client  *client.Client
	bars    *ProgressMng
	meter   *PhaseMeter
	balance *BalanceMeter
}

func RunGui(client *client.Client, events chan *client.Event) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	g.SetManagerFunc(func(g *gocui.Gui) error { return layout(g, client) })
	gui := &GUI{g, client, &ProgressMng{g, "cmds", 10, nil, &sync.Mutex{}}, nil,
		NewBalanceMeter(g, "balance")}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		view := gui.getView("input")
		view.Clear()
		return view.SetCursor(0, 0)
	}); err != nil {
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
		for _e := range events {
			e := _e
			gui.g.Update(func(g *gocui.Gui) error {
				return gui.handleClientEvent(e)
			})
		}
	}()

	if err := gui.g.MainLoop(); err != nil {
		panic(err)
	}
}

// must be called in a g.Update
func (gui *GUI) handleClientEvent(e *client.Event) error {
	switch e.Type {
	case client.SET_PARAMS:
		out := gui.getView("status")
		out.Clear()
		fmt.Fprintf(out, "Contract %s\nTEE %s\nPowDepth %d PhaseDuration %d ResponseDuration %d\nInitBlock %d", e.Params.Contract.Hex(), e.Params.TEE.Hex(), e.Params.PowDepth, e.Params.PhaseDuration, e.Params.ResponseDuration, e.Params.InitBlock)
		gui.meter = NewPhaseMeter(e.Params, gui.g, "phase")
	case client.SET_BALANCE:
		gui.balance.SetBalance(e.Report.Balance)
		gui.balance.Draw()
	case client.SET_OP_TRUST:
		out := gui.getView("status")
		out.Title = fmt.Sprintf("Operator is %s", e.OpTrust)
	case client.BENCH:
		gui.logOut(e.Result.String(), "\n")
	case client.NEW_BLOCK:
		if gui.meter != nil {
			gui.meter.SetBlock(e.BlockNum)
			gui.meter.Draw()
		}
	case client.CHAIN_MSG:
		gui.logChain(e.Message, "\n")
	case client.NEW_EPOCH:
		// ignored
	case client.SET_EXIT_AVAIL:
		gui.balance.SetExitPossible(e.ExitAvailable.Clone())
		gui.balance.Draw()
	default:
		panic(fmt.Sprintf("Unhandled enum case: %d", e.Type))
	}
	return nil
}

func (gui *GUI) enter(v *gocui.View) error {
	input := strings.TrimRight(v.Buffer(), "\n")
	gui.logOut("$ ", input, "\n")
	if strings.HasPrefix(input, "help") {
		gui.logOut(helpText)
		v.Clear()
		return v.SetCursor(0, 0)
	} else if strings.HasPrefix(input, "credits") {
		gui.logOut(creditText)
		v.Clear()
		return v.SetCursor(0, 0)
	}

	// Responsive gui -> do not block till command is finished.
	go func() {
		if status, err := gui.eval(input); err != nil {
			gui.logOut("✗ ", color(err.Error(), RED), "\n")
		} else {
			bar := gui.bars.Add(input)
			for s := range status {
				if s == nil {
					break
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
			bar.Finish(color("Done", GREEN))
			gui.logOut("✓ ", color(input, GREEN), "\n")
			gui.bars.Render()
		}
	}()

	v.Clear()
	return v.SetCursor(0, 0)
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
	case "challenge":
		go gui.client.CmdChallenge(status, fs[1:]...)
	case "leave":
		go gui.client.CmdLeave(status, fs[1:]...)
	case "exit":
		fallthrough
	case "quit":
		gui.client.Close()
		gui.g.Close()
		select {}
	default:
		return nil, fmt.Errorf("Unknow command: '%s'", cmd)
	}
	return status, nil
}

func layout(g *gocui.Gui, client *client.Client) error {
	maxW, maxH := g.Size()

	outW := (maxW / 3) * 2 // TODO golden ratio
	if v, err := g.SetView("status", 0, 0, outW, 4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Status"
	}
	if v, err := g.SetView("out", 0, 5, outW, maxH-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		addr := client.Address().Hex()[:10]
		v.Title = client.Config.UserName + fmt.Sprintf(" %s on %s", addr, client.ChainURL())
		v.Autoscroll = true
	}
	chainH := (maxH / 3) * 2
	if v, err := g.SetView("balance", outW+1, 0, maxW-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Balance"
	}
	if v, err := g.SetView("phase", outW+1, 4, maxW-1, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Phase meter"
	}
	if v, err := g.SetView("chain", outW+1, 10, maxW-1, chainH); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Log"
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

// ANSI Escape colors.
const (
	BLACK = iota
	RED
	GREEN
	ORANGE
	BLUE
	PURPLE
	CYAN
	GRAY
)

func bold(msg string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", msg)
}

func color(msg string, clr int) string {
	return fmt.Sprintf("\033[3%d;%dm%s\033[0m", clr, 1, msg)
}

func colorBg(msg string, clr int) string {
	return fmt.Sprintf("\033[4%d;%dm%s\033[0m", clr, 1, msg)
}

const helpText = `Available commands:
 help
   Prints this page.
 credits
   Hommage to the creators.
 bench <runs>
   Runs a benchmark of sending <runs> off-chain TX.
 deposit <amount>
   Deposits <amount> into the Erdstall contract.
 send <receiver> <amount>
   Sends <amount> to <receiver>.
 leave
   Withdraws all funds and exits the network.
 exit, quit
   Close the client.
`

const creditText = `Developed by Perun Network for the 2020 ETH Hackathon.

  Norbert Dzikowski  - Enclave, Intel SGX
  Matthias Geihs	 - Operator Node, Presentation
  Steffen Rattay	 - GrapheneOS, Presentation
  Sebastian Stammler - Enclave, Contract, Speaker
  Oliver Tale-Yazdi  - Client, GUI

Find the code at https://github.com/perun-network/erdstall
and checkout our projects https://perun.network
`
