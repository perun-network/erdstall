package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jroimartin/gocui"
	"github.com/perun-network/erdstall/eth"
)

type namedAccount struct {
	name    string
	account common.Address
}

const timeoutDuration = 10 * time.Second
const updateInterval = 1 * time.Second

func main() {
	ethNodeURL := flag.String("ethurl", "ws://127.0.0.1:8545", "URL of Ethereum node")
	accountsFilePath := flag.String("accounts", "accounts.json", "list of accounts to be tracked")
	flag.Parse()

	ethClient, err := ethclient.Dial(*ethNodeURL)
	if err != nil {
		log.Fatalf("dialing ethereum node: %v", err)
	}

	accounts, err := loadAccounts(*accountsFilePath)
	if err != nil {
		log.Fatalf("loading accounts: %v", err)
	}

	// configure gui

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatalf("creating gui: %v", err)
	}
	defer g.Close()

	g.SetManagerFunc(func(g *gocui.Gui) error {
		maxX, maxY := g.Size()
		if v, err := g.SetView("main", 0, 0, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				log.Fatalf("creating main view: %v", err)
			}

			if err := printAccounts(v, ethClient, accounts); err != nil {
				log.Fatalf("printing accounts: %v", err)
			}
		}
		return nil
	})

	go func() {
		for {
			g.Update(func(g *gocui.Gui) error {
				v, err := g.View("main")
				if err != nil {
					log.Fatalf("getting main view: %v", err)
				}
				v.Clear()

				if err := printAccounts(v, ethClient, accounts); err != nil {
					log.Fatalf("printing accounts: %v", err)
				}

				return nil
			})
			time.Sleep(updateInterval)
		}
	}()

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }); err != nil {
		log.Fatalf("setting key binding: %v", err)
	}

	// start gui

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("running gui: %v", err)
	}
}

func loadAccounts(accountsFilePath string) ([]namedAccount, error) {
	fileContent, err := ioutil.ReadFile(accountsFilePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var accountsAsStrings [][]string
	err = json.Unmarshal(fileContent, &accountsAsStrings)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	accounts := make([]namedAccount, len(accountsAsStrings))
	for i, a := range accountsAsStrings {
		if len(a) != 2 {
			return nil, fmt.Errorf("account entry %d has invalid length", i)
		}
		accounts[i] = namedAccount{name: a[0], account: common.HexToAddress(a[1])}
	}

	return accounts, nil
}

func printAccounts(v *gocui.View, ethClient *ethclient.Client, accounts []namedAccount) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	fmt.Fprintf(v, "%-16s%-48s%s\n\n", "Name", "Address", "ETH Balance")

	for _, a := range accounts {

		balance, err := ethClient.BalanceAt(ctx, a.account, nil)
		if err != nil {
			return fmt.Errorf("retrieving balance for account %x: %w", a.account, err)
		}

		fmt.Fprintf(v, "%-16s%-48s%v\n", a.name, fmt.Sprintf("0x%x", a.account), eth.WeiToEthFloat(balance))
	}

	fmt.Fprintf(v, "\nTime: %v\n", time.Now())

	return nil
}
