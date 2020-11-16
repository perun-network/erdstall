package operator

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"perun.network/go-perun/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"

	"github.com/perun-network/erdstall/operator/test"
	"github.com/perun-network/erdstall/tee"
)

func TestOperator(t *testing.T) {
	environment := initEnvironment(t)
	defer environment.Shutdown()

	errg := environment.errg
	user1 := environment.user1
	user2 := environment.user2

	AssertNoError(errg.Err())

	user1.TargetBalance = int64(10)
	user2.TargetBalance = int64(5)

	// deposit
	user1.Deposit()
	user2.Deposit()
	log.Info("operator_test.TestOperator: Deposited funds at contract")

	environment.WaitPhase()

	// get deposit proofs
	user1.DepositProof()
	user2.DepositProof()
	log.Info("operator_test.TestOperator: Retrieved deposit proofs")

	// transfer from user1 to user2
	user1.Transfer(user2, 3)
	log.Info("operator_test.TestOperator: Transfer from user1 to user2")

	environment.WaitPhase()

	// get balance proof
	user1.BalanceProof()
	user2.BalanceProof()
	log.Info("operator_test.TestOperator: Retrieved balance proofs")

	// transfer from user2 to user1
	user2.Transfer(user1, 2)
	log.Info("operator_test.TestOperator: Transfer from user2 to user1")

	environment.WaitPhase()

	// get balance proof
	user1.BalanceProof()
	user2.BalanceProof()
	log.Info("operator_test.TestOperator: Retrieved balance proofs")

	// transfer from user1 to user2 and transfer from user2 to user1
	user1.Transfer(user2, 1)
	user2.Transfer(user1, 1)
	log.Info("operator_test.TestOperator: Transfer from user1 to user2 and transfer from user2 to user1")

	environment.WaitPhase()

	// get balance proofs
	user1.BalanceProof()
	user2.BalanceProof()
	log.Info("operator_test.TestOperator: Retrieved balance proofs")

	// challenge response
	sub, exitEvents := user1.SubscribeToExitEvents()
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
	cfg               *Config
	cmd               *exec.Cmd
	operator          *Operator
	user1             *test.User
	user2             *test.User
	errg              *errors.Gatherer
	enclaveParameters tee.Parameters
}

const blockTime = 1 // block mining interval in seconds

func initEnvironment(t *testing.T) *environment {
	errg := errors.NewGatherer()

	cfg := newDefaultConfig()

	var cmd *exec.Cmd
	runGanache := true
	if runGanache {
		cmd = exec.Command(
			"ganache-cli",
			"--accounts=10",
			"--defaultBalanceEther=100",
			fmt.Sprintf("--mnemonic=\"%s\"", cfg.Mnemonic),
			fmt.Sprintf("--blockTime=%d", blockTime),
		)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(3 * time.Second)
	}

	operator := Setup(cfg)
	params := operator.EnclaveParams()
	log.Info("operator_test.initEnvironment: Created operator")
	errg.Go(func() error {
		return operator.Serve(cfg.Port)
	})
	time.Sleep(1 * time.Second)

	w, err := hdwallet.NewFromMnemonic(cfg.Mnemonic)
	AssertNoError(err)

	createUserAccount := func(userIndex int) accounts.Account {
		derivationPathUser := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", 2+userIndex))
		userAccount, err := w.Derive(derivationPathUser, true)
		AssertNoError(err)
		return userAccount
	}

	userAccount1 := createUserAccount(1)
	userAccount2 := createUserAccount(2)

	rpcURL := fmt.Sprintf("127.0.0.1:%d", cfg.Port)

	user1 := test.CreateUser(t, cfg.EthereumNodeURL, w, userAccount1, rpcURL, params.Contract, params)
	user2 := test.CreateUser(t, cfg.EthereumNodeURL, w, userAccount2, rpcURL, params.Contract, params)

	log.Info("operator_test.initEnvironment: Created users")

	return &environment{t, cfg, cmd, operator, user1, user2, errg, params}
}

func (e *environment) WaitPhase() {
	waitBlock := func() {
		time.Sleep(blockTime * time.Second)
	}
	for i := uint64(0); i < e.enclaveParameters.PhaseDuration; i++ {
		waitBlock()
	}
}

func (e *environment) Shutdown() {
	if e.cmd != nil {
		e.cmd.Process.Kill()
	}
}

func newDefaultConfig() *Config {
	return &Config{
		"ws://127.0.0.1:8545",
		"tag volcano eight thank tide danger coast health above argue embrace heavy",
		"m/44'/60'/0'/0/0",
		"m/44'/60'/0'/0/1",
		3,
		1,
		0,
		8080,
		true,
	}
}
