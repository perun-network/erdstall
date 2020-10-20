package test

import (
	"context"
	"fmt"
	"math/big"
	"net/rpc"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

const gasLimit = 2000000
const defaultContextTimeout = 10 * time.Second
const proofResponseTimeout = 10 * time.Second
const proofRequestInterval = 500 * time.Millisecond

func newDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultContextTimeout)
}

// User represents a TEE Plasma user.
type User struct {
	*testing.T
	wallet          accounts.Wallet
	account         accounts.Account
	ethClient       *eth.Client
	rpcClient       *rpc.Client
	contract        *bindings.Erdstall
	contractAddress common.Address
	nonceCounter    uint64
}

// Address returns the user's account address.
func (u *User) Address() common.Address {
	return u.ethClient.Account().Address
}

// Account returns the user's account.
func (u *User) Account() accounts.Account {
	return u.ethClient.Account()
}

// CreateUser creates a new user and connects it to the operator.
func CreateUser(t *testing.T, ethURL string, wallet accounts.Wallet, account accounts.Account, rpcURL string, contractAddress common.Address) *User {
	ethClient, err := eth.CreateEthereumClient(ethURL, wallet, account)
	if err != nil {
		t.Fatal("creating ethereum wallet and client:", err)
	}

	rpcClient, err := rpc.DialHTTP("tcp", rpcURL)
	if err != nil {
		t.Fatal("dialing rpc:", err)
	}

	contract, err := bindings.NewErdstall(contractAddress, ethClient)
	if err != nil {
		t.Fatal("loading contract:", err)
	}

	return &User{t, wallet, account, ethClient, rpcClient, contract, contractAddress, 0}
}

// Deposit deposits the specified amount at the TEE Plasma.
func (u *User) Deposit(amount int64) {
	ctx, cancel := newDefaultContext()
	defer cancel()

	t, err := u.ethClient.NewTransactor(ctx, big.NewInt(0), gasLimit, u.ethClient.Account())
	if err != nil {
		u.Fatal("creating transactor:", err)
	}

	t.Value = new(big.Int).SetInt64(amount)

	tx, err := u.contract.Deposit(t)
	if err != nil {
		u.Fatal("depositing:", err)
	}

	_, err = bind.WaitMined(ctx, u.ethClient, tx)
	if err != nil {
		u.Fatal("waiting for transaction confirmation:", err)
	}
}

// DepositProof returns the deposit proof for the last epoch.
func (u *User) DepositProof() *tee.DepositProof {
	var dp *tee.DepositProof
	err := u.rpcClient.Call("RemoteEnclave.GetDepositProof", u.Address(), dp)
	if err != nil {
		u.Fatal("calling RemoteEnclave.GetDepositProof:", err)
	}
	return dp
}

// NextDepositProof returns the deposit proof for the last epoch.
func (u *User) NextDepositProof() *tee.DepositProof {
	dpChan := make(chan *tee.DepositProof)

	go func() {
		var dp *tee.DepositProof
		for {
			if err := u.rpcClient.Call("RemoteEnclave.GetDepositProof", u.Address(), dp); err == nil {
				dpChan <- dp
			}
			time.Sleep(proofRequestInterval)
		}
	}()

	var dp *tee.DepositProof
	select {
	case dp = <-dpChan:
	case <-time.After(proofResponseTimeout):
		u.Fatal("deposit proof timeout")
	}

	return dp
}

// Transfer transfers the specified amount to the specified receiver.
func (u *User) Transfer(receiver common.Address, amount int64) {
	epoch, err := u.TransactionEpoch()
	if err != nil {
		u.Fatal("reading transaction epoch:", err)
	}

	tx := tee.Transaction{
		Nonce:     u.Nonce(),
		Epoch:     epoch,
		Sender:    u.Address(),
		Recipient: receiver,
		Amount:    big.NewInt(amount),
	}

	tx.Sign(u.contractAddress, u.Account(), u.wallet)

	err = u.rpcClient.Call("RemoteEnclave.AddTransaction", tx, nil)
	if err != nil {
		u.Fatal("RemoteEnclave.AddTransaction error:", err)
	}
}

// BalanceProof returns the balance proof for the last epoch.
func (u *User) BalanceProof() *tee.BalanceProof {
	var bp *tee.BalanceProof
	err := u.rpcClient.Call("RemoteEnclave.GetBalanceProof", u.Address(), bp)
	if err != nil {
		u.Fatal("calling RemoteEnclave.GetBalanceProof:", err)
	}
	return bp
}

// Nonce returns the next nonce.
func (u *User) Nonce() uint64 {
	u.nonceCounter++
	return u.nonceCounter
}

// TransactionEpoch returns the current transaction epoch.
func (u *User) TransactionEpoch() (uint64, error) {
	ctx, cancel := newDefaultContext()
	defer cancel()

	blockHeader, err := u.ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("reading block header: %w", err)
	}

	p := tee.Parameters{}

	return p.TxEpoch(blockHeader.NumberU64()), nil

}
