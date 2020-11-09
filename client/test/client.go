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
	}, nil
}

func (c *Client) NextNonce() uint64 {
	c.Nonce++
	return c.Nonce
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
	return c.call(ctx, func(tr *bind.TransactOpts) (*types.Transaction, error) {
		tr.Value = amount
		return c.contract.Deposit(tr)
	})
}

// Exit exits the Erdstall contract and waits until the exit TX is mined. The
// minedBlockNum is updated automatically.
func (c *Client) Exit(ctx context.Context, bal *tee.BalanceProof) error {
	return c.call(ctx, func(tr *bind.TransactOpts) (*types.Transaction, error) {
		return c.contract.Exit(tr, bindings.ErdstallBalance{
			Epoch:   bal.Balance.Epoch,
			Account: bal.Balance.Account,
			Value:   bal.Balance.Value,
		}, bal.Sig)
	})
}

// Withdraw withdraws the balance with the given `tee.BalanceProof` and waits
// until the withdraw TX is mined. The minedBlockNum is updated automatically.
func (c *Client) Withdraw(ctx context.Context, bal *tee.BalanceProof) error {
	return c.call(ctx, func(tr *bind.TransactOpts) (*types.Transaction, error) {
		return c.contract.Withdraw(tr, bal.Balance.Epoch)
	})
}

func (c *Client) Send(recipient common.Address, amount *big.Int) error {
	tx := &tee.Transaction{
		Nonce:     c.NextNonce(),
		Epoch:     c.params.TxEpoch(c.minedBlockNum + 1),
		Sender:    c.ethClient.Account().Address,
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
	return c.tr.Send(tx)
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
