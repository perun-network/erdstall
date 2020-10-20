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
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	peruneth "perun.network/go-perun/backend/ethereum/channel"
	perunhd "perun.network/go-perun/backend/ethereum/wallet/hd"

	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/tee"
)

const (
	DefaultGasLimit = 2000000
	defaultPowDepth = 0
)

type (
	// Client is an Ethereum client to interact with the Erdstall contract.
	Client struct {
		peruneth.ContractBackend
		account accounts.Account
	}

	// BlockSubscription represents a subscription to Ethereum blocks.
	BlockSubscription struct {
		sub    ethereum.Subscription
		blocks chan *tee.Block
		quit   chan struct{}
	}
)

// NewClient creates a new Erdstall Ethereum client.
func NewClient(
	cb peruneth.ContractBackend,
	a accounts.Account,
) *Client {
	return &Client{cb, a}
}

// NewClientForWalletAndAccount creates a new Erdstall Ethereum client for the
// given wallet and account.
func NewClientForWalletAndAccount(
	ci peruneth.ContractInterface,
	w accounts.Wallet,
	a accounts.Account,
) *Client {
	tr := NewDefaultTransactor(w)
	cb := peruneth.NewContractBackend(ci, tr)
	return &Client{cb, a}
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
	cb := peruneth.NewContractBackend(ci, tr)

	return NewClient(cb, acc.Account), nil
}

// CreateEthereumClient creates and connects a new ethereum client.
func CreateEthereumClient(url string, wallet accounts.Wallet, a accounts.Account) (*Client, error) {
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("dialing ethereum: %w", err)
	}
	return NewClientForWalletAndAccount(ethClient, wallet, a), nil
}

func (cl *Client) Account() accounts.Account {
	return cl.account
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

	quit := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Errorf("EthClient: Header subscription error: %v", err)
				sub.Unsubscribe()
				return
			case header := <-headers:
				log.WithFields(log.Fields{
					"blockNum": header.Number.Uint64(),
					"hash":     header.Hash().Hex()}).
					Debugf("EthClient: New header.")

				ctx, cancel := NewDefaultContext()

				block, err := cl.BlockByHash(ctx, header.Hash())
				cancel()
				if err != nil {
					log.Errorf("EthClient: Error retrieving block: %v", err)
				}

				select {
				case blocks <- block:
				case <-quit:
					close(blocks)
					log.Debug("EthClient: subscription closed")
					return
				}

			case <-quit:
				close(blocks)
				log.Debug("EthClient: subscription closed")
				return
			}
		}
	}()

	return &BlockSubscription{sub, blocks, quit}, nil
}

// SubscribeToDeposited writes past Deposited events and newly received ones
// into the sink. The `epochs` and `accs` arguments can be used to filter for
// specific events. Passing `nil` will skip the filtering.
// Should be started in a go-routine, since it blocks.
// Can be cancelled via context.
func (cl *Client) SubscribeToDeposited(ctx context.Context, contract *bindings.Erdstall, epochs []uint64, accs []common.Address, sink chan *bindings.ErdstallDeposited) error {
	// sub to new events
	wOpts, err := cl.NewWatchOpts(ctx)
	if err != nil {
		return err
	}
	sub, err := contract.WatchDeposited(wOpts, sink, epochs, accs)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()
	// filter past events
	fOpts, err := cl.NewFilterOpts(ctx)
	if err != nil {
		return err
	}
	it, err := contract.FilterDeposited(fOpts, epochs, accs)
	defer it.Close()
	if err != nil {
		return err
	}
	for it.Next() {
		sink <- it.Event
	}
	// Wait for error or ctx done.
	select {
	case err := <-sub.Err():
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
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

	tr, err := cl.NewTransactor(ctx, big.NewInt(0), DefaultGasLimit, cl.account)
	if err != nil {
		return fmt.Errorf("creating transactor: %w", err)
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

func (cl *Client) BindContract(ctx context.Context, addr common.Address) (*tee.Parameters, *bindings.Erdstall, error) {
	contract, err := bindings.NewErdstall(addr, cl.ContractInterface)
	if err != nil {
		return nil, nil, err
	}
	opts := &bind.CallOpts{Context: ctx}
	phaseDuration, err := contract.PhaseDuration(opts)
	if err != nil {
		return nil, nil, err
	}
	responseDuration, err := contract.ResponseDuration(opts)
	if err != nil {
		return nil, nil, err
	}
	bigBang, err := contract.BigBang(opts)
	if err != nil {
		return nil, nil, err
	}
	teeAddr, err := contract.Tee(opts)
	if err != nil {
		return nil, nil, err
	}
	return &tee.Parameters{
		PowDepth:         defaultPowDepth,
		PhaseDuration:    phaseDuration,
		ResponseDuration: responseDuration,
		InitBlock:        bigBang,
		TEE:              teeAddr,
		Contract:         addr,
	}, contract, nil
}

// Blocks returns the channel on which to receive subscribed blocks.
func (sub *BlockSubscription) Blocks() <-chan *tee.Block { return sub.blocks }

// Unsubscribe ends the subscription.
func (sub *BlockSubscription) Unsubscribe() {
	sub.sub.Unsubscribe()
	select {
	case <-sub.quit: // already closed
	default:
		close(sub.quit)
	}
}
