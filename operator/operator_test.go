// SPDX-License-Identifier: Apache-2.0

package operator_test

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"perun.network/go-perun/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	op "github.com/perun-network/erdstall/operator"
	"github.com/perun-network/erdstall/operator/test"
	"github.com/perun-network/erdstall/tee"
)

func TestOperator(t *testing.T) {
	environment := initEnvironment(t)
	t.Cleanup(environment.Shutdown)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	errg := environment.errg
	user1 := environment.user1
	user2 := environment.user2

	op.AssertNoError(errg.Err())

	user1.TargetBalance = int64(10)
	user2.TargetBalance = int64(5)

	environment.WaitPhase()

	// deposit
	user1.Deposit()
	user2.Deposit()
	log.Info("operator_test.TestOperator: Deposited funds at contract")

	environment.WaitBlocks(1)

	// get deposit proofs
	user1.DepositProof(ctx)
	user2.DepositProof(ctx)
	log.Info("operator_test.TestOperator: Retrieved deposit proofs")

	// transfer from user1 to user2
	user1.Transfer(ctx, user2, 3)
	log.Info("operator_test.TestOperator: Transfer from user1 to user2")

	environment.WaitPhase()

	// get balance proof
	user1.BalanceProof(ctx)
	user2.BalanceProof(ctx)
	log.Info("operator_test.TestOperator: Retrieved balance proofs")

	// transfer from user2 to user1
	user2.Transfer(ctx, user1, 2)
	log.Info("operator_test.TestOperator: Transfer from user2 to user1")

	environment.WaitPhase()

	// get balance proof
	user1.BalanceProof(ctx)
	user2.BalanceProof(ctx)
	log.Info("operator_test.TestOperator: Retrieved balance proofs")

	// transfer from user1 to user2 and transfer from user2 to user1
	user1.Transfer(ctx, user2, 1)
	user2.Transfer(ctx, user1, 1)
	log.Info("operator_test.TestOperator: Transfer from user1 to user2 and transfer from user2 to user1")

	environment.WaitPhase()

	// get balance proofs
	user1.BalanceProof(ctx)
	user2.BalanceProof(ctx)
	log.Info("operator_test.TestOperator: Retrieved balance proofs")

	// challenge response
	sub, exitEvents := user1.SubscribeExitEvents()
	defer sub.Unsubscribe()
	user1.Challenge()
	onChainTransactionTimeout := time.Duration(blockTime*environment.operator.EnclaveParams().PhaseDuration) * time.Second

	select {
	case err := <-sub.Err():
		user1.Fatalf("exit event subscription error: %v", err)
	case exitEvent := <-exitEvents:
		if exitEvent.Account != user1.Address() {
			user1.Errorf("invalid account, expected %v, got %v", user1.Address().String(), exitEvent.Account.String())
		}
	case <-time.After(onChainTransactionTimeout):
		user1.Fatalf("exit event timeout")
	}

	// // exit
	// user1.Exit()
	// user2.Exit()

	// // withdraw
	// user1.Withdraw()
	// user2.Withdraw()
}

type environment struct {
	*testing.T
	cfg               *op.Config
	cmd               *exec.Cmd
	operator          *op.Operator
	user1             *test.User
	user2             *test.User
	errg              *errors.Gatherer
	enclaveParameters tee.Parameters
}

const blockTime = 1 // block mining interval in seconds

func initEnvironment(t *testing.T) *environment {
	errg := errors.NewGatherer()

	cfg := newDefaultConfig()

	prog, args := op.GanacheCommand()
	args = append(args,
		"--accounts=10",
		"--defaultBalanceEther=100",
		fmt.Sprintf("--mnemonic=\"%s\"", cfg.Mnemonic),
		fmt.Sprintf("--blockTime=%d", blockTime),
	)
	cmd := exec.Command(prog, args...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	operator := op.SetupWithPrototypeEnclave(cfg)
	params := operator.EnclaveParams()
	log.Info("operator_test.initEnvironment: Created operator")
	errg.Go(func() error {
		return operator.Serve(cfg.RPCPort)
	})
	time.Sleep(1 * time.Second)

	w, err := hdwallet.NewFromMnemonic(cfg.Mnemonic)
	op.AssertNoError(err)

	createUserAccount := func(userIndex int) accounts.Account {
		derivationPathUser := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", 2+userIndex))
		userAccount, err := w.Derive(derivationPathUser, true)
		op.AssertNoError(err)
		return userAccount
	}

	userAccount1 := createUserAccount(1)
	userAccount2 := createUserAccount(2)

	user1 := test.CreateUser(t, cfg.EthereumNodeURL, w, userAccount1, cfg.RPCHost, cfg.RPCPort, params.Contract, params)
	user2 := test.CreateUser(t, cfg.EthereumNodeURL, w, userAccount2, cfg.RPCHost, cfg.RPCPort, params.Contract, params)

	log.Info("operator_test.initEnvironment: Created users")

	return &environment{t, cfg, cmd, operator, user1, user2, errg, params}
}

func (e *environment) WaitPhase() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(blockTime*e.cfg.PhaseDuration+1)*time.Second)
	defer cancel()

	e.WaitBlockPredicate(ctx, func(block uint64) bool {
		return e.enclaveParameters.IsLastPhaseBlock(block)
	})
}

func (e *environment) WaitBlocks(n int) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(blockTime*n+1)*time.Second)
	defer cancel()

	e.WaitBlockPredicate(ctx, func(uint64) bool {
		n--
		return n <= 0
	})
}

func (e *environment) WaitBlockPredicate(ctx context.Context, p func(uint64) bool) {
	heads := make(chan *types.Header)
	sub, err := e.operator.EthClient.SubscribeNewHead(ctx, heads)
	if err != nil {
		e.T.Fatal("subscribing to header: ", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case head := <-heads:
			if p(head.Number.Uint64()) {
				time.Sleep(200 * time.Millisecond)
				return
			}
		case <-ctx.Done():
			e.T.Fatal("context: ", ctx.Err())
		case err := <-sub.Err():
			e.T.Fatal("header subscription: ", err)
		}
	}
}

func (e *environment) Shutdown() {
	if e.cmd != nil {
		if err := e.cmd.Process.Kill(); err != nil {
			log.Warn("Could not kill process:", err)
		}
	}
}

func newDefaultConfig() *op.Config {
	return &op.Config{
		"ws://127.0.0.1:8545",
		"pistol kiwi shrug future ozone ostrich match remove crucial oblige cream critic",
		"m/44'/60'/0'/0/0",
		"m/44'/60'/0'/0/1",
		3,
		1,
		0,
		8401,
		"0.0.0.0",
		true,
	}
}
