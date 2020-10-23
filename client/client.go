// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	pethchannel "perun.network/go-perun/backend/ethereum/channel"
	pethwallet "perun.network/go-perun/backend/ethereum/wallet"
	psync "perun.network/go-perun/pkg/sync"
	pwallet "perun.network/go-perun/wallet"

	"github.com/perun-network/erdstall/config"
	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
)

type Client struct {
	psync.Closer
	// Initialized in NewClient
	Config          config.ClientConfig
	conn            *RPC
	ethClient       *eth.Client
	contractAddr    common.Address
	signer          tee.TextSigner
	txNonce         uint64
	balances        map[uint64]EpochBalance // epoch => balance
	balProofWatcher sync.Once
	// Initialized in Run()
	currentBlock uint64 // Atomic
	contract     *bindings.Erdstall
	params       *tee.Parameters
	chainEvents  chan string
	clEvents     chan *ClientEvent
	// balMtx protects balances.
	balMtx sync.Mutex
}

// EpochBalance describes the balance that a specific user has/has in a epoch.
// Currently, onls only one deposit per Epoch is allowed otherwise it would need
// a slice of `*DepositProof`s.
// The Proofs are initialized to nil, the DepositProof should be set after the
// Deposit succeeded.
// The BalanceProof should be set at the end of a Transaction phase.
type EpochBalance struct {
	tee.Balance
	Dep *tee.DepositProof
	Bal *tee.BalanceProof
	// TODO isthis epoch exiting or depositing?
}

// BalanceReport will be displayed by the GUI.
type BalanceReport struct {
	Balance *big.Int
}

func (e *EpochBalance) Clone() *EpochBalance {
	var dep *tee.DepositProof
	if e.Dep != nil {
		dep = &tee.DepositProof{
			Balance: e.Dep.Balance.Clone(),
			Sig:     append([]byte(nil), e.Dep.Sig...),
		}
	}
	var bal *tee.BalanceProof
	if e.Bal != nil {
		bal = &tee.BalanceProof{
			Balance: e.Bal.Balance.Clone(),
			Sig:     append([]byte(nil), e.Bal.Sig...),
		}
	}
	return &EpochBalance{
		Balance: e.Balance.Clone(),
		Dep:     dep,
		Bal:     bal,
	}
}

type ClientEventType = int

const (
	SET_PARAMS ClientEventType = iota
	SET_BALANCE
	SET_OP_TRUST // Emitted when operator becomes malicious
	BENCH
)

const (
	// How long do we wait for the deposit proof after the deposit phase is over.
	depositProofGrace = time.Second * 30
	balanceProofGrace = time.Second * 30
	blockTime         = time.Second * 2
)

// Trust describes how we perceive the operator.
type Trust = string

const (
	UNKNOWN   = "😕 Unknown"
	TRUSTED   = "😀 Benevolent"
	UNTRUSTED = "😠 Malicious"
)

type ClientEvent struct {
	Type    ClientEventType
	Params  tee.Parameters // SET_PARAMS
	Report  BalanceReport  // SET_BALANCE
	OpTrust Trust
	Result  Result
}

type CmdStatus struct {
	Msg string
	Err error
	War string // Warning
}

func NewClient(cfg config.ClientConfig, conn *RPC, ethClient *eth.Client, signer tee.TextSigner) *Client {
	addr, err := strToCommonAddress(cfg.Contract)
	if err != nil {
		panic(err)
	}
	return &Client{
		Config:       cfg,
		conn:         conn,
		ethClient:    ethClient,
		contractAddr: addr,
		signer:       signer,
		txNonce:      1,
		balances:     make(map[uint64]EpochBalance),
	}
}

func (c *Client) Run(clEvents chan *ClientEvent, events chan string) error {
	events <- "Connecting to contract..."
	params, contract, err := c.ethClient.BindContract(shortCtx(), c.contractAddr)
	if err != nil {
		return err
	}
	events <- "Connected to contract"
	clEvents <- &ClientEvent{Type: SET_PARAMS, Params: *params}
	clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: TRUSTED}

	c.params = params
	c.contract = contract
	c.chainEvents = events
	c.clEvents = clEvents

	return c.listenOnChain()
}

func (c *Client) CmdSend(status chan *CmdStatus, args ...string) {
	defer close(status)
	if len(args) != 2 {
		status <- &CmdStatus{Err: errors.New("Command 'send' needs arguments: <receiver> <amount>")}
		return
	}
	receiver, err := strToCommonAddress(args[0])
	if args[0] == "me" {
		receiver = c.Address()
	} else if err != nil {

		status <- &CmdStatus{Err: fmt.Errorf("Invalid <receiver>: %v", err)}
		return
	}
	_amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		status <- &CmdStatus{Err: fmt.Errorf("Invalid <amount>: %v", err)}
		return
	}
	amount := eth.EthToWeiFloat(_amount)
	status <- &CmdStatus{Msg: "Creating Message"}
	/*block, err := blockNum(c.ethClient.ContractInterface)
	if err != nil {
		c.chainEvents <- fmt.Sprintf("Could not read block number: %s", err.Error())
		status <- &CmdStatus{Err: errors.New("Chain error")}
		return
	}*/
	tx, err := c.createTx(receiver, amount)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	status <- &CmdStatus{Msg: "Forwarding to Operator"}
	if err := c.conn.AddTX(shortCtx(), tx); err != nil {
		status <- &CmdStatus{Err: err}
	}
}

func (c *Client) createTx(receiver common.Address, amount *big.Int) (tee.Transaction, error) {
	block := atomic.LoadUint64(&c.currentBlock) + 1

	tx := tee.Transaction{
		Nonce:     c.txNonce,
		Epoch:     c.params.TxEpoch(block),
		Sender:    c.Address(),
		Recipient: receiver,
		Amount:    amount,
	}
	c.txNonce++
	err := tx.Sign(c.params.Contract, c.ethClient.Account(), c.signer)
	if err != nil {
		fmt.Printf("Could not sign: %v", err)
	}

	if ok, err := tee.VerifyTransaction(c.params.Contract, tx); !ok || err != nil {
		fmt.Printf("Could not sign: %v", err)
	}
	return tx, err
}

func (c *Client) CmdBench(status chan *CmdStatus, args ...string) {
	defer close(status)
	if len(args) > 1 {
		status <- &CmdStatus{Err: errors.New("Command 'bench' can have argument: <runs>")}
		return
	}
	var err error
	n := 1000
	if len(args) == 1 {
		n, err = strconv.Atoi(args[0])
		if err != nil {
			status <- &CmdStatus{Err: fmt.Errorf("Could not parse <runs> as int: %s", args[0])}
			return
		}
	}
	status <- &CmdStatus{Msg: fmt.Sprintf("Sending %d payments", n)}
	result, err := Benchmark(n, func() error {
		tx, err := c.createTx(c.Address(), big.NewInt(1)) // 1 WEI
		if err != nil {
			return err
		}
		return c.conn.AddTX(shortCtx(), tx)
	})
	c.clEvents <- &ClientEvent{Type: BENCH, Result: result}
	return
}

func (c *Client) CmdDeposit(status chan *CmdStatus, args ...string) {
	defer close(status)
	if len(args) != 1 {
		status <- &CmdStatus{Err: errors.New("Command 'deposit' needs argument: <amount>")}
		return
	}
	_amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		status <- &CmdStatus{Err: fmt.Errorf("Could not parse <amount> as float: %s", args[0])}
		return
	}
	amount := eth.EthToWeiFloat(_amount)

	status <- &CmdStatus{Msg: "Deposit TX: Preparing"}
	otps, err := c.ethClient.NewTransactor(txCtx(), amount, eth.DefaultGasLimit, c.ethClient.Account())
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	status <- &CmdStatus{Msg: "Deposit TX: Sending"}
	tx, err := c.contract.Deposit(otps)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	status <- &CmdStatus{Msg: "Deposit TX: Waiting for on-chain confirmation"}
	c.chainEvents <- fmt.Sprintf("Sent deposit TX with hash %s", tx.Hash().Hex())
	res, err := bind.WaitMined(txCtx(), c.ethClient, tx)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	c.chainEvents <- fmt.Sprintf("Deposit TX mined in Block #%d", res.BlockNumber.Uint64())
	status <- &CmdStatus{Msg: "Deposit proof: Waiting for Operator"}
	blockNum := res.BlockNumber.Uint64()
	// The epoch that we want to do the deposit in.
	epoch := c.params.DepositEpoch(blockNum)

	proof := make(chan tee.DepositProof)
	proofErr := make(chan error)
	waitErr := make(chan error)
	// no specific wait ctx, since we do not know what the block time is.
	ctx, cancel := context.WithCancel(c.Ctx())
	defer cancel()
	go func() {
		if p, err := c.conn.GetDepositProof(ctx, epoch, c.Address()); err != nil {
			proofErr <- err
		} else {
			proof <- p
		}
	}()
	go func() {
		depEndBlock := c.params.DepositEndBlock(epoch)
		err := c.ethClient.WaitForBlock(ctx, depEndBlock) // Add PowDepth here if needed
		if err == nil {
			c.chainEvents <- fmt.Sprintf("Deposit Phase for Epoch %d ended in block %d", epoch, depEndBlock)
			time.Sleep(depositProofGrace)
		}
		waitErr <- err
	}()
	select {
	case e := <-waitErr:
		if e != nil {
			c.chainEvents <- fmt.Sprintf("WaitForEpoch: %s", e.Error())
			status <- &CmdStatus{War: "Deposit proof: Chain error - resuming protocol"}
			c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNKNOWN}
		} else {
			status <- &CmdStatus{War: "Deposit proof: Operator timed out - resuming protocol"}
			c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNTRUSTED}
		}
	case p := <-proof:
		status <- &CmdStatus{Msg: "Deposit proof: Verifying"}
		ok, err := tee.VerifyDepositProof(*c.params, p)
		if !ok || err != nil {
			status <- &CmdStatus{War: "Deposit proof: Invalid Signature - resuming protocol"}
			c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNTRUSTED}
		} else if p.Balance.Value.Cmp(amount) != 0 || p.Balance.Epoch != epoch || p.Balance.Account != c.Address() {
			status <- &CmdStatus{War: "Deposit proof: Wrong Proof - resuming protocol"}
			c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNTRUSTED}
		} else {
			c.balMtx.Lock()
			defer c.balMtx.Unlock()
			if _, ok := c.balances[epoch]; ok {
				panic("Deposit can not be called twice in one epoch")
			}
			// TODO add instead of replace
			c.balances[epoch] = EpochBalance{Balance: p.Balance, Dep: &p}
			status <- &CmdStatus{Msg: "Deposit proof: Valid"}
			return
		}
	case e := <-proofErr:
		status <- &CmdStatus{War: fmt.Sprintf("Deposit proof: '%s' - resuming protocol", e.Error())}
		c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNKNOWN}
	case <-ctx.Done():
		status <- &CmdStatus{War: "Deposit proof: Operator timed out - resuming protocol"}
		c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNKNOWN}
	}
	// TODO challenge
	status <- &CmdStatus{Err: errors.New("TODO challenge")}
	return
}

func errorReason(ctx context.Context, b *pethchannel.ContractBackend, tx *types.Transaction, blockNum *big.Int, acc accounts.Account) (string, error) {
	msg := ethereum.CallMsg{
		From:     acc.Address,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	res, err := b.CallContract(ctx, msg, blockNum)
	if err != nil {
		return "", fmt.Errorf("CallContract: %v", err)
	}
	reason, err := abi.UnpackRevert(res)
	return reason, fmt.Errorf("unpacking revert reason: %v", err)
}

// BalanceProofWatcher waits for the balance proof of an epoch and disputes
// if non was received after the TxPhase + balanceProofGrace.
// Can currently only deal with one deposit per deposit-phase.
// Should be started in a go-routine.
func (c *Client) BalanceProofWatcher() {
	oldEpoch := uint64(0)
	for {
		proof, err := c.conn.GetBalanceProof(c.Ctx(), c.Address())
		if err != nil {
			c.chainEvents <- fmt.Sprintf("Balance Proof: %v", err)
			c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNKNOWN}
		} else if proof.Balance.Epoch > oldEpoch {
			c.chainEvents <- fmt.Sprintf("Balance Proof: bal=%v, epoch=%d", eth.WeiToEthFloat(proof.Balance.Value), proof.Balance.Epoch)
			oldEpoch = proof.Balance.Epoch
			c.balMtx.Lock()
			c.balances[proof.Balance.Epoch] = EpochBalance{Balance: proof.Balance, Bal: &proof}
			c.balMtx.Unlock()

			ok, err := tee.VerifyBalanceProof(*c.params, proof)
			if !ok || err != nil {
				c.chainEvents <- fmt.Sprintf("Invalid balance proof: err=%v ok=%t", err, ok)
				return
			}

			c.clEvents <- &ClientEvent{Type: SET_BALANCE, Report: BalanceReport{Balance: new(big.Int).Set(proof.Balance.Value)}}
			c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: TRUSTED}
		}
		time.Sleep(time.Second)
	}
	return
	/*
		oldEpoch := uint64(0)
		for {
			epoch := atomic.LoadUint64(&c.currentBlock)
			if oldEpoch == epoch {
				time.Sleep(time.Second)
				continue
			}
			oldEpoch = epoch

			txEndBlock := c.params.TxEndBlock(epoch)
			waitErr := make(chan error)
			proofErr := make(chan error)
			proof := make(chan tee.BalanceProof)
			ctx, cancel := context.WithCancel(c.Ctx())
			defer cancel()

			go func() {
				if p, err := c.conn.GetBalanceProof(ctx, epoch, c.Address()); err != nil {
					proofErr <- err
				} else {
					proof <- p
				}
			}()
			go func() {
				err := c.ethClient.WaitForBlock(ctx, txEndBlock)
				if err == nil {
					c.chainEvents <- fmt.Sprintf("Transaction Phase for Epoch %d ended in block %d", epoch, txEndBlock)
					time.Sleep(balanceProofGrace)
				}
				waitErr <- err
			}()

			select {
			case p := <-proof:
				c.balMtx.Lock()
				defer c.balMtx.Unlock()
				bal := c.balances[epoch]
				ok, err := tee.VerifyBalanceProof(*c.params, p)
				if !ok || err != nil {
					c.chainEvents <- "Balance proof: Invalid Signature - challenging operator"
					c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNTRUSTED}
				} else if p.Balance.Value.Cmp(bal.Value) < 0 || p.Balance.Epoch != epoch || p.Balance.Account != c.Address() {
					c.chainEvents <- "Balance proof: Wrong balance - challenging operator"
					c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNTRUSTED}
				} else {
					bal.Bal = &p
					c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: TRUSTED}
					c.clEvents <- &ClientEvent{Type: SET_BALANCE, Report: c.report()}
					continue
					//return
				}
			case e := <-waitErr:
				if e != nil {
					c.chainEvents <- fmt.Sprintf("WaitForEpoch: %s", e.Error())
					c.chainEvents <- "Chain error - challenging operator"
					c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNKNOWN}
				} else {
					c.chainEvents <- "Balance proof timed out - challenging operator"
					c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: UNTRUSTED}
				}
			case <-ctx.Done():
				return
			}

			// TODO challenge
			c.chainEvents <- "TODO: Challenge"
			//return
		}*/
}

func (c *Client) CmdLeave(status chan *CmdStatus, args ...string) {
	defer close(status)
	if len(args) != 0 {
		status <- &CmdStatus{Err: errors.New("Command 'leave' does not accept arguments.")}
		return
	}
	epoch := c.params.ExitEpoch(atomic.LoadUint64(&c.currentBlock))
	c.balMtx.Lock()
	defer c.balMtx.Unlock()
	var bal EpochBalance
	bal1, ok1 := c.balances[epoch]
	bal2, ok2 := c.balances[epoch-1]

	if ok1 {
		bal = bal1
	} else if ok2 {
		bal = bal2
	} else {
		status <- &CmdStatus{Err: errors.New("No Balance Proof available")}
		return
	}
	if bal.Bal == nil {
		status <- &CmdStatus{Err: errors.New("No Balance Proof available")}
		return
	}

	status <- &CmdStatus{Msg: "Exit TX: Preparing"}
	otps, err := c.ethClient.NewTransactor(txCtx(), big.NewInt(0), eth.DefaultGasLimit, c.ethClient.Account())
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	status <- &CmdStatus{Msg: "Exit TX: Sending"}
	tx, err := c.contract.Exit(otps, bindings.ErdstallBalance{Epoch: bal.Epoch, Account: bal.Account, Value: bal.Value}, bal.Bal.Sig)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	status <- &CmdStatus{Msg: fmt.Sprintf("Waiting for TX: %s", tx.Hash().Hex())}
	rec, err := bind.WaitMined(txCtx(), c.ethClient, tx)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	if rec.Status == types.ReceiptStatusFailed {
		status <- &CmdStatus{Err: errors.New("Exit TX: Receipt failed")}
		reason, err := errorReason(c.Ctx(), &c.ethClient.ContractBackend, tx, rec.BlockNumber, c.ethClient.Account())
		if err != nil {
			c.chainEvents <- fmt.Sprintf("Unknown revert reason: %v", err)
		} else {
			c.chainEvents <- fmt.Sprintf("Exit TX revert reason: %s", reason)
		}
		return
	}
	c.chainEvents <- fmt.Sprintf("Exit mined in block #%d", rec.BlockNumber.Uint64())
	// End of exit is begin of next
	endExitBlock := c.params.EpochStartBlock(c.params.ExitEpoch(rec.BlockNumber.Uint64()) + 1)
	if err := c.ethClient.WaitForBlock(c.Ctx(), endExitBlock); err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	status <- &CmdStatus{Msg: "Withdraw TX: Preparing"}
	otps, err = c.ethClient.NewTransactor(txCtx(), big.NewInt(0), eth.DefaultGasLimit, c.ethClient.Account())
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	status <- &CmdStatus{Msg: "Withdraw TX: Sending"}
	tx, err = c.contract.Withdraw(otps, bal.Epoch)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	status <- &CmdStatus{Msg: fmt.Sprintf("Waiting for TX: %s", tx.Hash().Hex())}
	rec, err = bind.WaitMined(txCtx(), c.ethClient, tx)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	if rec.Status == types.ReceiptStatusFailed {
		status <- &CmdStatus{Err: errors.New("Withdraw TX: Receipt failed")}
		return
	}
	c.chainEvents <- fmt.Sprintf("Withdraw mined in block #%d", rec.BlockNumber.Uint64())
	c.clEvents <- &ClientEvent{Type: SET_BALANCE, Report: BalanceReport{Balance: big.NewInt(0)}}
	c.clEvents <- &ClientEvent{Type: SET_OP_TRUST, OpTrust: TRUSTED}
}

// writes to chainVvents
func (c *Client) listenOnChain() error {
	epochs := make(chan uint64)
	blocks := make(chan uint64)
	subError := make(chan error)

	go func() {
		subError <- c.ethClient.SubscribeToEpochs(c.Ctx(), *c.params, epochs, blocks)
	}()
	go c.BalanceProofWatcher()

	for !c.IsClosed() {
		select {
		case epoch := <-epochs:
			c.chainEvents <- fmt.Sprintf("New Epoch #%d", epoch)
		case block := <-blocks:
			atomic.StoreUint64(&c.currentBlock, block)
			c.chainEvents <- fmt.Sprintf("new Block #%d", block)
		case err := <-subError:
			return err
		}
	}
	return nil
}

// must hold balMtx while calling
func (c *Client) report() BalanceReport {
	sum := new(big.Int)
	for _, b := range c.balances {
		sum.Add(sum, b.Value)
	}
	return BalanceReport{Balance: sum}
}

func (c *Client) Address() common.Address {
	return c.ethClient.Account().Address
}

func blockNum(cb pethchannel.ContractInterface) (uint64, error) {
	h, err := cb.HeaderByNumber(context.TODO(), nil)
	if err != nil {
		return 0, err
	}
	if !h.Number.IsUint64() {
		panic("Block number too big")
	}
	return h.Number.Uint64(), nil
}

func strToPerunAddress(str string) (pwallet.Address, error) {
	if len(str) != 42 {
		return nil, errors.New("Public keys must be chars 40 hex")
	}
	h, err := hex.DecodeString(str[2:])
	if err != nil {
		return nil, errors.New("Could not parse address as hexadecimal")
	}
	addr, err := pwallet.DecodeAddress(bytes.NewBuffer(h))
	return addr, err
}

func strToCommonAddress(s string) (common.Address, error) {
	var err error
	var walletAddr pwallet.Address

	if walletAddr, err = strToPerunAddress(s); err != nil {
		return common.Address{}, err
	}

	return pethwallet.AsEthAddr(walletAddr), nil
}

func txCtx() context.Context {
	return newCtx(time.Second * 30)
}

func shortCtx() context.Context {
	return newCtx(time.Second * 20)
}

func newCtx(d time.Duration) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(d)
		cancel()
	}()
	return ctx
}
