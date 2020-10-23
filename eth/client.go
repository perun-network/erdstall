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
	// DefaultGasLimit represents the default gas limit for this application.
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

	// BlockSubscription2 represents a subscription to Ethereum blocks, version 2.
	BlockSubscription2 struct {
		headerSub       ethereum.Subscription
		blocks          chan *tee.Block
		err             chan error
		nextBlockNumber *big.Int
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
	tr := perunhd.NewTransactor(hdw.Wallet())
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

// Account returns the account of the client.
func (cl *Client) Account() accounts.Account {
	return cl.account
}

func (cl *Client) WaitForBlock(ctx context.Context, target uint64) error {
	sub, err := cl.SubscribeToBlocks()
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case b := <-sub.Blocks():
			if b.NumberU64() >= target {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// SubscribeToEpochs writes new epochs into `sink`. Can be cancelled via ctx.
// Should be called in a go-routine, since it blocks.
// Future InitBlock numbers are not supported, so it will not return the 0. th epoch.
func (cl *Client) SubscribeToEpochs(ctx context.Context, params tee.Parameters, sink chan uint64, blockSink chan uint64) error {
	sub, err := cl.SubscribeToBlocks()
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	oldEpoch := uint64(0)
	for {
		select {
		case block := <-sub.Blocks():
			blockSink <- block.NumberU64()
			if newEpoch := params.DepositEpoch(block.NumberU64()); newEpoch > oldEpoch {
				sink <- newEpoch
				oldEpoch = newEpoch
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
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
				if err != nil {
					log.Errorf("EthClient: Header subscription error: %v", err)
				} else {
					log.Debug("EthClient: Header subscription closed")
				}
				close(blocks)
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

// SubscribeToBlocksStartingFrom subscribes the client to the mined Ethereum
// blocks starting from the given block number.
func (cl *Client) SubscribeToBlocksStartingFrom(startBlockNumber *big.Int) (*BlockSubscription2, error) {
	headers := make(chan *types.Header)
	blocks := make(chan *tee.Block)
	errChan := make(chan error)

	ctx, cancel := NewDefaultContext()
	defer cancel()

	headerSub, err := cl.SubscribeNewHead(ctx, headers)
	if err != nil {
		return nil, fmt.Errorf("subscribing to blockchain head: %w", err)
	}

	blockSub := &BlockSubscription2{headerSub, blocks, errChan, startBlockNumber}

	run := func() error {
		for {
			select {
			case err := <-headerSub.Err():
				if err != nil {
					return fmt.Errorf("subscription error: %w", err)
				}
				log.Debug("EthClient: Header subscription closed")
				close(blocks)
				return nil

			case header := <-headers:
				log.WithFields(log.Fields{
					"blockNum": header.Number.Uint64(),
					"hash":     header.Hash().Hex()}).
					Debugf("EthClient: New header.")

				for blockSub.nextBlockNumber.Cmp(header.Number) < 0 {
					if err := blockSub.pushNextBlock(cl); err != nil {
						return fmt.Errorf("pushing block: %w", err)
					}
				}

				if err := blockSub.pushNextBlock(cl); err != nil {
					return fmt.Errorf("pushing block: %w", err)
				}
			}
		}
	}

	go func() {
		err := run()
		errChan <- err
	}()

	return blockSub, nil
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

	receipts, err := cl.TransactionReceipts(ctx, block)
	if err != nil {
		return nil, fmt.Errorf("retrieving receipts: %w", err)
	}

	return &tee.Block{Block: *block, Receipts: receipts}, nil
}

// TransactionReceipts returns the transaction receipts for the given block.
func (cl *Client) TransactionReceipts(ctx context.Context, block *types.Block) (types.Receipts, error) {
	var receipts []*types.Receipt
	for _, t := range block.Transactions() {
		r, err := cl.TransactionReceipt(ctx, t.Hash())
		if err != nil {
			return nil, fmt.Errorf("retrieving receipt: %w", err)
		}
		receipts = append(receipts, r)
	}
	return receipts, nil
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

// Blocks returns the channel on which to receive subscribed blocks.
func (blockSub *BlockSubscription2) Blocks() <-chan *tee.Block { return blockSub.blocks }

// Unsubscribe ends the subscription.
func (blockSub *BlockSubscription2) Unsubscribe() {
	blockSub.headerSub.Unsubscribe()
}

func (blockSub *BlockSubscription2) pushNextBlock(cl *Client) error {
	ctx, cancel := NewDefaultContext()
	defer cancel()

	block, err := cl.ContractBackend.BlockByNumber(ctx, blockSub.nextBlockNumber)
	if err != nil {
		return fmt.Errorf("retrieving block by number: %w", err)
	}

	receipts, err := cl.TransactionReceipts(ctx, block)
	if err != nil {
		return fmt.Errorf("retrieving block receipts: %w", err)
	}

	log.Tracef("eth.Client: Pushing block number %d", block.NumberU64())

	blockSub.blocks <- &tee.Block{Block: *block, Receipts: receipts}
	blockSub.nextBlockNumber = new(big.Int).Add(blockSub.nextBlockNumber, big.NewInt(1))
	return nil
}
