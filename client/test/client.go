// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"

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
	tr        tee.Transactor

	Nonce         uint64
	minedBlockNum uint64
	balance       *big.Int // local balance tracking
}

// NewClient creates a testing Erdstall client.
// The wallet w must contain the Account of ethClient for Erdstall TX sending.
func NewClient(params tee.Parameters, wallet tee.TextSigner, ethClient *eth.Client, tr tee.Transactor) (*Client, error) {
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

func (c *Client) AddBalance(amount *big.Int) {
	c.balance.Add(c.balance, amount)
}

func (c *Client) Address() common.Address {
	return c.ethClient.Account().Address
}

func (c *Client) SetMinedBlockNum(n uint64) {
	c.minedBlockNum = n
}

// TxEpoch return the current transaction epoch as calculated from the currently
// mined block number known to the client and the enclave parameters.
func (c *Client) TxEpoch() tee.Epoch {
	return c.params.TxEpoch(c.minedBlockNum + 1)
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
		return c.contract.Exit(tr, bal.Balance.ToEthBals(), bal.Sig)
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
		Epoch:     c.TxEpoch(),
		Sender:    c.Address(),
		Recipient: recipient,
		Amount:    (*tee.Amount)(amount),
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
		recipient.AddBalance(amount)
	}
	return err
}

// SendInvalidTxs sends several invalid transactions. It uses InvalidTxs, see
// its documentation for which invalid transactions are used.
func (c *Client) SendInvalidTxs(rng *rand.Rand, validRecipient common.Address) (errs []error) {
	for _, tx := range c.InvalidTxs(rng, validRecipient) {
		errs = append(errs, c.tr.Send(tx))
	}
	return errs
}

// InvalidTxs generates a list of several invalid transactions, each having one
// of the following fields set to an invalid value:
// - Sig (random)
// - Nonce (+-1)
// - Epoch (+-1)
// - Recipient (random)
// - Sender (random)
// - Amount (1 above max, -1)
// The remaining fields are set to a valid value.
func (c *Client) InvalidTxs(rng *rand.Rand, validRecipient common.Address) (txs []*tee.Transaction) {
	invalidators := []func(*tee.Transaction){
		func(tx *tee.Transaction) {
			tx.Sig = make([]byte, 65)
			rng.Read(tx.Sig)
		},
		func(tx *tee.Transaction) {
			tx.Nonce += 1
			c.SignTx(tx)
		},
		func(tx *tee.Transaction) {
			tx.Nonce -= 1
			c.SignTx(tx)
		},
		func(tx *tee.Transaction) {
			tx.Epoch += 1
			c.SignTx(tx)
		},
		func(tx *tee.Transaction) {
			tx.Epoch -= 1
			c.SignTx(tx)
		},
		func(tx *tee.Transaction) {
			tx.Recipient = eth.NewRandomAddress(rng)
			c.SignTx(tx)
		},
		func(tx *tee.Transaction) {
			tx.Sender = eth.NewRandomAddress(rng)
			if err := tx.SignAlien(c.params.Contract, c.ethClient.Account(), c.wallet); err != nil {
				log.Panicf("Error signing alien tx: %v", err)
			}
		},
		func(tx *tee.Transaction) {
			tx.Amount = (*tee.Amount)(new(big.Int).Add((*big.Int)(tx.Amount), big.NewInt(1)))
			c.SignTx(tx)
		},
		func(tx *tee.Transaction) {
			tx.Amount = (*tee.Amount)(big.NewInt(-1))
			c.SignTx(tx)
		},
	}

	for _, invalidate := range invalidators {
		// set to valid and then invalidate
		tx := &tee.Transaction{
			Nonce:     c.Nonce + 1, // only increment locally
			Epoch:     c.TxEpoch(),
			Sender:    c.ethClient.Account().Address,
			Recipient: validRecipient,
			Amount:    (*tee.Amount)(c.Balance()), // total balance is valid
		}
		invalidate(tx)
		txs = append(txs, tx)
	}

	return txs
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
