// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/client"
	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

const defaultContextTimeout = 10 * time.Second

func newDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultContextTimeout)
}

// User represents a TEE Plasma user.
type User struct {
	*testing.T
	wallet            accounts.Wallet
	account           accounts.Account
	ethClient         *eth.Client
	rpcClient         *client.RPC
	proofSub          *client.Subscription
	contract          *bindings.Erdstall
	contractAddress   common.Address
	nonceCounter      uint64
	TargetBalance     int64
	enclaveParameters tee.Parameters
	dp                tee.DepositProof
	bp                tee.BalanceProof
	epoch             tee.Epoch
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
func CreateUser(
	t *testing.T,
	ethURL string,
	wallet accounts.Wallet,
	account accounts.Account,
	rpcHost string, rpcPort uint16,
	contractAddress common.Address,
	enclaveParameters tee.Parameters,
) *User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ethClient, err := eth.CreateEthereumClient(ctx, ethURL, wallet, account)
	if err != nil {
		t.Fatal("creating ethereum wallet and client:", err)
	}

	rpcClient, err := client.NewRPC(rpcHost, rpcPort)
	if err != nil {
		t.Fatal("dialing rpc:", err)
	}
	sub, err := rpcClient.Subscribe(ctx, account.Address)
	if err != nil {
		t.Fatal("subscribing:", err)
	}

	contract, err := bindings.NewErdstall(contractAddress, ethClient)
	if err != nil {
		t.Fatal("loading contract:", err)
	}

	return &User{t, wallet, account, ethClient, rpcClient, sub, contract, contractAddress, 0, 0, enclaveParameters, tee.DepositProof{}, tee.BalanceProof{}, 0}
}

// Deposit deposits the current target balance at the TEE Plasma.
func (u *User) Deposit() {
	ctx, cancel := newDefaultContext()
	defer cancel()

	t, err := u.ethClient.NewTransactor(ctx)
	if err != nil {
		u.Fatal("creating transactor:", err)
	}

	t.Value = new(big.Int).SetInt64(u.TargetBalance)

	tx, err := u.contract.Deposit(t)
	if err != nil {
		u.Fatal("depositing:", err)
	}

	r, err := bind.WaitMined(ctx, u.ethClient, tx)
	if err != nil {
		u.Fatal("waiting for transaction confirmation:", err)
	}

	if r.Status != types.ReceiptStatusSuccessful {
		u.Fatal("deposit transaction failed:", err)
	}
	log.Debugf("Deposited %d in block %d", u.TargetBalance, r.BlockNumber.Uint64())
}

// DepositProof returns the deposit proof for the last epoch.
func (u *User) DepositProof(ctx context.Context) {
	proof, err := u.proofSub.DepositProof(ctx)
	if err != nil {
		u.Fatal("calling Subscription.DepositProof:", err)
	}
	u.dp = proof

	if (*big.Int)(u.dp.Balance.Value).Int64() != u.TargetBalance {
		u.FailNow()
	}

	log.Debug("Got deposit proof for epoch #", u.dp.Balance.Epoch)
	u.epoch = u.dp.Balance.Epoch
}

// Transfer transfers the specified amount to the specified receiver.
func (u *User) Transfer(ctx context.Context, receiver *User, amount int64) {
	log.Debug("Sending transfer in epoch #", u.dp.Balance.Epoch)
	tx := tee.Transaction{
		Nonce:     u.Nonce(),
		Epoch:     u.epoch,
		Sender:    u.Address(),
		Recipient: receiver.Address(),
		Amount:    (*tee.Amount)(big.NewInt(amount)),
	}

	if err := tx.Sign(u.contractAddress, u.Account(), u.wallet); err != nil {
		u.Fatal("Signing transaction:", err)
	}

	err := u.rpcClient.SendTx(ctx, tx)
	if err != nil {
		u.Fatal("RPC.Send error:", err)
	}

	u.TargetBalance -= amount
	receiver.TargetBalance += amount
}

// BalanceProof returns the balance proof for the last epoch.
func (u *User) BalanceProof(ctx context.Context) {
	proof, err := u.proofSub.BalanceProof(ctx)
	if err != nil {
		u.Fatal("calling Subscription.BalanceProof:", err)
	}
	u.bp = proof

	if balance := (*big.Int)(u.bp.Balance.Value).Int64(); balance != u.TargetBalance {
		u.Errorf("incorrect balance, got %d, expected %d", balance, u.TargetBalance)
	}

	log.Debug("Got balance proof for epoch #", u.bp.Balance.Epoch)
	u.epoch = u.bp.Balance.Epoch + 1
}

// Nonce returns the next nonce.
func (u *User) Nonce() uint64 {
	u.nonceCounter++
	return u.nonceCounter
}

// SubscribeExitEvents subscribes the user the exit events.
func (u *User) SubscribeExitEvents() (event.Subscription, chan *bindings.ErdstallExiting) {
	exitEvents := make(chan *bindings.ErdstallExiting)
	sub, err := u.contract.WatchExiting(nil, exitEvents, nil, nil)
	if err != nil {
		u.Fatal("subscribing to exit events:", err)
	}

	return sub, exitEvents
}

// Challenge challenges the operator for the balance proof of the current epoch.
func (u *User) Challenge() {
	ctx, cancel := newDefaultContext()
	defer cancel()

	tr, err := u.ethClient.NewTransactor(ctx)
	if err != nil {
		u.Fatal("creating transactor:", err)
	}

	tx, err := u.contract.Challenge(tr)
	if err != nil {
		u.Fatal("sending challenge transaction:", err)
	}

	r, err := bind.WaitMined(ctx, u.ethClient, tx)
	if err != nil {
		u.Fatal("waiting for transaction confirmation:", err)
	}

	if r.Status != types.ReceiptStatusSuccessful {
		u.Fatal("challenge transaction failed:", err)
	}
}
