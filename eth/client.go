// SPDX-License-Identifier: Apache-2.0

package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	log "github.com/sirupsen/logrus"
	peruneth "perun.network/go-perun/backend/ethereum/channel"
	perunhd "perun.network/go-perun/backend/ethereum/wallet/hd"

	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/tee"
)

const defaultGasLimit = 2000000

type (
	// Client is an Ethereum client to interact with the Erdstall contract.
	Client struct {
		peruneth.ContractInterface
		tr      peruneth.Transactor
		account accounts.Account
	}

	// BlockSubscription represents a subscription to Ethereum blocks.
	BlockSubscription struct {
		sub    ethereum.Subscription
		blocks chan *tee.Block
	}
)

// NewClient creates a new Erdstall Ethereum client.
func NewClient(
	ci peruneth.ContractInterface,
	tr peruneth.Transactor,
	a accounts.Account,
) *Client {
	return &Client{ci, tr, a}
}

// NewClientForWallet returns a new Client using the given wallet as transactor.
// The first account is derived from the wallet and used as the account in the
// client.
func NewClientForWallet(
	ci peruneth.ContractInterface,
	w accounts.Wallet,
) (*Client, error) {
	hdw, err := perunhd.NewWallet(w, perunhd.DefaultRootDerivationPath.String(), 0)
	if err != nil {
		return nil, fmt.Errorf("creating perun hd wallet wrapper: %w", err)
	}
	acc, err := hdw.NewAccount()
	if err != nil {
		return nil, fmt.Errorf("deriving account: %w", err)
	}
	tr := perunhd.NewTransactor(hdw)

	return NewClient(ci, tr, acc.Account), nil
}

// SubscribeToBlocks subscribes the client to the mined Ethereum blocks.
func (cl *Client) SubscribeToBlocks() (*BlockSubscription, error) {
	headers := make(chan *types.Header)
	blocks := make(chan *tee.Block)

	ctx, cancel := NewDefaultContext()
	defer cancel()

	sub, err := cl.SubscribeNewHead(ctx, headers)
	if err != nil {
		return nil, fmt.Errorf("subscribing to blockchain head: %w", err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Errorf("Header subscription error: %v", err)
				sub.Unsubscribe()
				return
			case header := <-headers:
				log.WithFields(log.Fields{
					"num":  header.Number.Uint64(),
					"hash": header.Hash().Hex()}).
					Debugf("New header. Num %x", header.Hash().Hex())

				ctx, cancel := NewDefaultContext()

				block, err := cl.BlockByHash(ctx, header.Hash())
				if err != nil {
					log.Errorf("Error retrieving block: %v", err)
				}

				cancel()
				blocks <- block
			}
		}
	}()

	return &BlockSubscription{sub, blocks}, nil
}

// BlockByHash returns the block for the given block hash together with the block's transaction receipts.
func (cl *Client) BlockByHash(ctx context.Context, hash common.Hash) (*tee.Block, error) {
	block, err := cl.ContractInterface.BlockByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("retrieving block: %w", err)
	}

	var receipts []*types.Receipt
	for _, t := range block.Transactions() {
		r, err := cl.TransactionReceipt(ctx, t.Hash())
		if err != nil {
			return nil, fmt.Errorf("retrieving receipt: %w", err)
		}
		receipts = append(receipts, r)
	}

	return &tee.Block{Block: *block, Receipts: receipts}, nil
}

// NewTransactor creates a new transactor.
func (cl *Client) NewTransactor(ctx context.Context, gasLimit uint64) (*bind.TransactOpts, error) {
	// Determine nonce
	nonce, err := cl.NextNonce(ctx)
	if err != nil {
		return nil, fmt.Errorf("determining nonce: %w", err)
	}

	// Determine gas price
	gasPrice, err := cl.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("determining gas price: %w", err)
	}

	auth, err := cl.tr.NewTransactor(cl.account)
	if err != nil {
		return nil, fmt.Errorf("creating transactor: %w", err)
	}
	auth.Nonce = new(big.Int).SetUint64(nonce)
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	return auth, nil
}

// NextNonce returns the next Ethereum nonce.
func (cl *Client) NextNonce(ctx context.Context) (uint64, error) {
	nonce, err := cl.PendingNonceAt(ctx, cl.account.Address)
	if err != nil {
		return 0, fmt.Errorf("retrieving nonce: %w", err)
	}

	return nonce, nil
}

// DeployContracts deploys the Erdstall contract to the blockchain. It updates
// the passed parameter's InitBlock and Contract fields.
func (cl *Client) DeployContracts(params *tee.Parameters) error {
	ctx, cancel := NewDefaultContext()
	defer cancel()

	tr, err := cl.NewTransactor(ctx, defaultGasLimit)
	if err != nil {
		return fmt.Errorf("creating keyed transactor: %w", err)
	}

	address, tx, _, err := bindings.DeployErdstall(tr,
		cl.ContractInterface,
		params.TEE,
		params.PhaseDuration,
		params.ResponseDuration)
	if err != nil {
		return fmt.Errorf("deploying contract: %w", err)
	}

	_, err = bind.WaitDeployed(ctx, cl, tx)
	if err != nil {
		return fmt.Errorf("waiting for contract deployment: %w", err)
	}
	params.Contract = address

	receipt, err := cl.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return fmt.Errorf("getting tx receipt: %w", err)
	}
	params.InitBlock = receipt.BlockNumber.Uint64()

	return nil
}

// Blocks returns the channel on which to receive subscribed blocks.
func (sub *BlockSubscription) Blocks() <-chan *tee.Block { return sub.blocks }

// Unsubscribe ends the subscription.
func (sub *BlockSubscription) Unsubscribe() {
	sub.sub.Unsubscribe()
}
