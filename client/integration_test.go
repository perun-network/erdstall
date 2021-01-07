// SPDX-License-Identifier: Apache-2.0
// +build integration

package client_test

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/perun-network/erdstall/client"
	"github.com/perun-network/erdstall/config"
	"github.com/perun-network/erdstall/eth"
	op "github.com/perun-network/erdstall/operator"
	"github.com/perun-network/erdstall/wallet"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	perunchannel "perun.network/go-perun/backend/ethereum/channel"
	perunhd "perun.network/go-perun/backend/ethereum/wallet/hd"
	pkgtest "perun.network/go-perun/pkg/test"
)

const (
	mnemonic  = "myth like bonus scare over problem client lizard pioneer submit female collect"
	ethUrl    = "ws://127.0.0.1:8545"
	rpcPort   = 8401
	blockTime = 2 * time.Second
)

func TestWalkthroughs(t *testing.T) {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	require.NoError(t, startChain(t))
	time.Sleep(5 * time.Second) // wait for chain start

	t.Run("honest", func(t *testing.T) {
		testWalkthrough(t, &op.Config{
			RespondChallenges: true,
			SendDepositProofs: true,
			SendBalanceProofs: true,
		})
	})

	t.Run("dishonest-no-CR", func(t *testing.T) {
		testWalkthrough(t, &op.Config{
			RespondChallenges: false,
			SendDepositProofs: true,
			SendBalanceProofs: true,
		})
	})
}

func testWalkthrough(t *testing.T, honesty *op.Config) {
	t.Run("deposit-send-leave", func(t *testing.T) {
		operator, clients := setup(t, honesty, 2)
		testDeposit(t, clients...)
		waitBlocks(4.5)
		testSend(t, 0, clients...)
		testSend(t, 1, clients...)
		waitBlocks(1)
		testLeave(t, clients...)
		cleanup(operator, clients)
	})

	t.Run("deposit-send-challenge", func(t *testing.T) {
		operator, clients := setup(t, honesty, 2)
		testDeposit(t, clients...)
		waitBlocks(4.5)
		testSend(t, 0, clients...)
		testSend(t, 1, clients...)
		waitBlocks(8.5)
		// Bug in OP does not allow for parallel challenges.
		testChallenge(t, clients[0])
		cleanup(operator, clients)
	})
}

func testDeposit(t *testing.T, clients ...*client.Client) {
	testPhase(t, "deposit", func(client *client.Client, status chan *client.CmdStatus) {
		client.CmdDeposit(status, "100")
	}, clients...)
}

func testSend(t *testing.T, recipient int, clients ...*client.Client) {
	testPhase(t, "send", func(client *client.Client, status chan *client.CmdStatus) {
		client.CmdSend(status, clients[recipient].Address().Hex(), "10")
	}, clients...)
}

func testLeave(t *testing.T, clients ...*client.Client) {
	testPhase(t, "leave", func(client *client.Client, status chan *client.CmdStatus) {
		client.CmdLeave(status)
	}, clients...)
}

func testChallenge(t *testing.T, clients ...*client.Client) {
	testPhase(t, "challenge", func(client *client.Client, status chan *client.CmdStatus) {
		client.CmdChallenge(status)
	}, clients...)
}

func testPhase(t *testing.T, phase string, fn func(client *client.Client, status chan *client.CmdStatus), clients ...*client.Client) {
	log.Info("Phase: ", phase)
	ct := pkgtest.NewConcurrent(t)
	for _, client := range clients {
		client := client
		go ct.StageN(phase, len(clients), func(t pkgtest.ConcT) {
			status := filterStatus(t)
			fn(client, status)
		})
	}
	ct.Wait(phase)
}

func setup(t *testing.T, honesty *op.Config, numClients int) (operator *op.Operator, clients []*client.Client) {
	// Start Operator.
	operator = startOp(honesty)
	// Start clients.
	for i := 0; i < numClients; i++ {
		client, err := startClient(i, operator.EnclaveParams().Contract)
		time.Sleep(3 * time.Second) // wait for chain connection
		require.NoError(t, err)
		clients = append(clients, client)
	}
	return
}

func cleanup(operator *op.Operator, clients []*client.Client) {
	for _, client := range clients {
		client.Close()
	}
	operator.Close()
}

// filterStatus returns a channel that the client can use to report the
// status of a command. The channel is filtered for errors and all errors
// are required to be nil.
func filterStatus(t pkgtest.ConcT) chan *client.CmdStatus {
	status := make(chan *client.CmdStatus)
	go func() {
		for update := range status {
			require.NoError(t, update.Err)
		}
	}()
	return status
}

func startClient(index int, contract common.Address) (*client.Client, error) {
	cfg := config.ClientConfig{
		ChainURL:     ethUrl,
		OpHost:       "127.0.0.1",
		OpPort:       rpcPort,
		Mnemonic:     mnemonic,
		AccountIndex: index + 2,
		Contract:     contract.String(),
		UserName:     fmt.Sprintf("client-%d", index),
	}

	wallet := wallet.NewWallet(cfg.Mnemonic, uint(cfg.AccountIndex)) // HD Wallet
	eb, err := ethclient.Dial(cfg.ChainURL)
	if err != nil {
		return nil, err
	}
	events := make(chan *client.Event, 10) // GUI event pipe
	go func() {
		for e := range events {
			if e.Message != "" {
				log.Info(cfg.UserName, e.Message)
			}
		}
	}()

	cb := perunchannel.NewContractBackend(eb, perunhd.NewTransactor(wallet.Wallet.Wallet()))
	rpc, err := client.NewRPC(cfg.OpHost, uint16(cfg.OpPort))
	if err != nil {
		return nil, err
	}
	chain := eth.NewClient(cb, wallet.Acc.Account) // ETHChain conn
	log.Info(cfg.UserName, " address: ", wallet.Acc.Address())
	client := client.NewClient(cfg, rpc, events, chain, wallet)
	go func() {
		if err := client.Run(); err != nil {
			log.Info("client stopped: ", err)
		}
	}()
	return client, nil
}

func startChain(t *testing.T) error {
	prog, args := op.GanacheCommand()
	args = append(args,
		"-a 4",
		"-e 1000",
		"-d", // deterministic
		fmt.Sprintf("-b %f", blockTime.Seconds()),
	)
	cmd := exec.Command(prog, args...)
	t.Cleanup(func() {
		if err := cmd.Process.Kill(); err != nil {
			log.Error("stopping ganache: ", err)
		}
	})
	return cmd.Start()
}

func startOp(honesty *op.Config) *op.Operator {
	cfg := &op.Config{
		EthereumNodeURL:        ethUrl,
		Mnemonic:               mnemonic,
		EnclaveDerivationPath:  "m/44'/60'/0'/0/1",
		OperatorDerivationPath: "m/44'/60'/0'/0/0",
		PhaseDuration:          3,
		ResponseDuration:       1,
		PowDepth:               0,
		RPCPort:                8401,
		RPCHost:                "0.0.0.0",
		RespondChallenges:      honesty.RespondChallenges,
		SendDepositProofs:      honesty.SendDepositProofs,
		SendBalanceProofs:      honesty.SendBalanceProofs,
	}

	operator := op.SetupWithPrototypeEnclave(cfg)
	go func() {
		if err := operator.Serve(rpcPort); err != nil {
			panic(fmt.Sprintf("Operator.Serve: %v", err))
		}
	}()
	// Wait until we can be sure that the server is up and running.
	time.Sleep(5 * time.Second)
	return operator
}

func waitBlocks(amount float64) {
	time.Sleep(time.Duration(float64(blockTime) * amount))
}
