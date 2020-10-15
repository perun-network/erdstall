package operator

import (
	"fmt"
	"log"
	"os/exec"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/perun-network/erdstall/operator/test"
)

func TestMain(t *testing.T) {
	environment := initEnvironment(t)
	defer environment.Shutdown()

	user1 := environment.user1
	user2 := environment.user2

	// deposit
	user1.Deposit(10)
	user2.Deposit(10)

	// wait deposit phase
	environment.WaitPhase()

	// get deposit proofs
	dp1 := user1.NextDepositProof()
	log.Println("user1 deposit proof: ", dp1)
	dp2 := user2.NextDepositProof()
	log.Println("user2 deposit proof: ", dp2)

	// add transaction
	user1.Transfer(user2.Address(), 3)

	// wait phase
	environment.WaitPhase()

	// get balance proof
	bp1 := user1.BalanceProof()
	log.Println("user1 balance proof: ", bp1)

	// add transactions
	user1.Transfer(user2.Address(), 2)
	user2.Transfer(user1.Address(), 1)

	// wait phase
	environment.WaitPhase()

	// get balance proofs
	bp1 = user1.BalanceProof()
	log.Println("user1 balance proof: ", bp1)
	bp2 := user2.BalanceProof()
	log.Println("user2 balance proof: ", bp2)

	// // exit
	// user1.Exit()
	// user2.Exit()

	// // withdraw
	// user1.Withdraw()
	// user2.Withdraw()
}

type environment struct {
	*testing.T
	cfg      *Config
	cmd      *exec.Cmd
	operator *Operator
	user1    *test.User
	user2    *test.User
}

const blockTime = 1

func initEnvironment(t *testing.T) *environment {
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

	operator, enclaveParameters := Setup(cfg)
	log.Println("created operator")
	go func() {
		err := operator.Serve(cfg.Port)
		AssertNoError(err)
	}()
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

	user1 := test.CreateUser(t, cfg.EthereumNodeURL, w, userAccount1, rpcURL, enclaveParameters.Contract)
	user2 := test.CreateUser(t, cfg.EthereumNodeURL, w, userAccount2, rpcURL, enclaveParameters.Contract)

	log.Println("created users")

	return &environment{t, cfg, cmd, operator, user1, user2}
}

func (e *environment) WaitPhase() {
	/*
		waitBlock := func() {
			dummyTransaction := func() {
				ctx, cancel := NewDefaultContext()
				defer cancel()

				signedTx, err := e.operator.NewSignedTransaction(ctx, common.Address{}, big.NewInt(1))
				if err != nil {
					e.Error("creating signed transaction", err)
				}

				err = e.operator.SendTransaction(context.Background(), signedTx)
				if err != nil {
					e.Error("sending transaction", err)
				}
			}

			dummyTransaction()
		}
	*/

	waitBlock := func() {
		time.Sleep(blockTime * time.Second)
	}
	waitBlock()
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
	}
}
