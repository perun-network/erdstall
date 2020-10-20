// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

// Deposit deposits amount to the Erdstall contract and waits until the deposit
// tx is mined. The minedBlockNum is updated automatically.
func (c *Client) Deposit(ctx context.Context, amount *big.Int) error {
	tr, err := c.ethClient.NewTransactor(ctx, amount, eth.DefaultGasLimit, c.ethClient.Account())
	if err != nil {
		return fmt.Errorf("creating transactor: %w", err)
	}

	tx, err := c.contract.Deposit(tr)
	if err != nil {
		return fmt.Errorf("calling deposit: %w", err)
	}

	rec, err := bind.WaitMined(ctx, c.ethClient, tx)
	if err != nil {
		return fmt.Errorf("waiting for deposit block: %w", err)
	}
	c.SetMinedBlockNum(uint64(rec.BlockNumber.Int64()))

	return nil
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

	return c.tr.Send(tx)
}

func (c *Client) SignTx(tx *tee.Transaction) {
	if err := tx.Sign(c.params.Contract, c.ethClient.Account(), c.wallet); err != nil {
		log.Panicf("Error signing tx: %v", err)
	}
}
