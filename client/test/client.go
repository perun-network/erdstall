// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"

	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

type Client struct {
	params    tee.Parameters
	ethClient *eth.Client
	contract  *bindings.Erdstall
	wallet    tee.TextSigner
	tr        Transactor

	Nonce         uint64
	minedBlockNum uint64
	balance       *big.Int // local balance tracking
}

// NewClient creates a testing Erdstall client.
// The wallet w must contain the Account of ethClient for Erdstall TX sending.
func NewClient(params tee.Parameters, wallet tee.TextSigner, ethClient *eth.Client, tr Transactor) (*Client, error) {
	contract, err := bindings.NewErdstall(params.Contract, ethClient)
	if err != nil {
		return nil, fmt.Errorf("binding Erdstall contract: %w", err)
	}

	return &Client{
		params:    params,
		ethClient: ethClient,
		contract:  contract,
		wallet:    wallet,
		tr:        tr,
		balance:   new(big.Int),
	}, nil
}

func (c *Client) nextNonce() uint64 {
	c.Nonce++
	return c.Nonce
}

func (c *Client) Balance() *big.Int {
	return new(big.Int).Set(c.balance)
}

func (c *Client) Address() common.Address {
	return c.ethClient.Account().Address
}

func (c *Client) SetMinedBlockNum(n uint64) {
	c.minedBlockNum = n
}

func (c *Client) UpdateLastBlockNum() {
	ctx, cancel := eth.NewDefaultContext()
	defer cancel()
	h, err := c.ethClient.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Panicf("Error getting latest head: %v", err)
	}
	c.SetMinedBlockNum(uint64(h.Number.Int64()))
}

// Deposit deposits amount to the Erdstall contract and waits until the deposit
// TX is mined. The minedBlockNum is updated automatically.
func (c *Client) Deposit(ctx context.Context, amount *big.Int) error {
	err := c.call(ctx, func(tr *bind.TransactOpts) (*types.Transaction, error) {
		tr.Value = amount
		return c.contract.Deposit(tr)
	})
	if err == nil {
		c.balance.Add(c.balance, amount)
	}
	return err
}

// Exit exits the Erdstall contract and waits until the exit TX is mined. The
// minedBlockNum is updated automatically.
func (c *Client) Exit(ctx context.Context, bal *tee.BalanceProof) error {
	err := c.call(ctx, func(tr *bind.TransactOpts) (*types.Transaction, error) {
		return c.contract.Exit(tr, bindings.ErdstallBalance{
			Epoch:   bal.Balance.Epoch,
			Account: bal.Balance.Account,
			Value:   bal.Balance.Value,
		}, bal.Sig)
	})
	if err == nil {
		c.balance.SetUint64(0)
	}
	return err
}

// Withdraw withdraws the balance with the given `tee.BalanceProof` and waits
// until the withdraw TX is mined. The minedBlockNum is updated automatically.
func (c *Client) Withdraw(ctx context.Context, bal *tee.BalanceProof) error {
	return c.call(ctx, func(tr *bind.TransactOpts) (*types.Transaction, error) {
		return c.contract.Withdraw(tr, bal.Balance.Epoch)
	})
}

// Send sends amount to the given recipient.
// If you need proper balance tracking for testing, use SendToClient instead.
func (c *Client) Send(recipient common.Address, amount *big.Int) error {
	tx := &tee.Transaction{
		Nonce:     c.nextNonce(),
		Epoch:     c.params.TxEpoch(c.minedBlockNum + 1),
		Sender:    c.Address(),
		Recipient: recipient,
		Amount:    amount,
	}

	c.SignTx(tx)

	log := log.WithFields(log.Fields{
		"nonce":  tx.Nonce,
		"epoch":  tx.Epoch,
		"amount": tx.Amount,
		"sender": tx.Sender.String(),
	})
	log.Debug("Sending TX")
	defer log.Trace("TX sent")
	err := c.tr.Send(tx)
	if err == nil {
		c.balance.Sub(c.balance, amount)
	}
	return err
}

// SendToClient sends amount to the given recipient.
// Use this if you need proper balance tracking in your tests.
func (c *Client) SendToClient(recipient *Client, amount *big.Int) error {
	err := c.Send(recipient.Address(), amount)
	if err == nil {
		recipient.balance.Add(recipient.balance, amount)
	}
	return err
}

func (c *Client) SignTx(tx *tee.Transaction) {
	if err := tx.Sign(c.params.Contract, c.ethClient.Account(), c.wallet); err != nil {
		log.Panicf("Error signing tx: %v", err)
	}
}

// call calls the given function, waits for the TX to be mined and updates the
// clients last mined blocknumber.
func (c *Client) call(ctx context.Context, call func(*bind.TransactOpts) (*types.Transaction, error)) error {
	tr, err := c.ethClient.NewTransactor(ctx)
	if err != nil {
		return err
	}

	tx, err := call(tr)
	if err != nil {
		return fmt.Errorf("calling contract: %w", err)
	}

	rec, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		return fmt.Errorf("waiting for block containing TX: %w", err)
	}
	if rec.Status == types.ReceiptStatusFailed {
		return fmt.Errorf("execution of contract call failed")
	}
	c.SetMinedBlockNum(uint64(rec.BlockNumber.Int64()))

	return nil
}
