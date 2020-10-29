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
	events       chan *Event
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
}

// BalanceReport will be displayed by the GUI.
type BalanceReport struct {
	Balance *big.Int
}

func (e *EpochBalance) Clone() *EpochBalance {
	if e == nil {
		return nil
	}
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

type EventType = int

const (
	SET_PARAMS EventType = iota
	SET_BALANCE
	SET_OP_TRUST   // Emitted when operator becomes malicious
	SET_EXIT_AVAIL // Emitted when an exit becomes (in)available

	NEW_BLOCK // New block mined
	NEW_EPOCH
	CHAIN_MSG // Chain related message.

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

type Event struct {
	Type          EventType
	Params        tee.Parameters // SET_PARAMS
	Report        BalanceReport  // SET_BALANCE
	OpTrust       Trust
	Result        Result
	Message       string        // CHAIN_MSG
	BlockNum      uint64        // NEW_BLOCK
	EpochNum      uint64        // NEW_EPOCH
	ExitAvailable *EpochBalance // SET_EXIT_AVAIL
}

type CmdStatus struct {
	Msg string
	Err error
	War string // Warning
}

func NewClient(cfg config.ClientConfig, conn *RPC, events chan *Event, ethClient *eth.Client, signer tee.TextSigner) *Client {
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
		events:       events,
	}
}

func (c *Client) Run() error {
	c.logOnChain("Connecting to contract...")
	params, contract, err := c.ethClient.BindContract(shortCtx(), c.contractAddr)
	if err != nil {
		return err
	}
	c.logOnChain("Connected to contract")
	c.events <- &Event{Type: SET_PARAMS, Params: *params}
	c.events <- &Event{Type: SET_OP_TRUST, OpTrust: TRUSTED}

	c.params = params
	c.contract = contract
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
	c.events <- &Event{Type: BENCH, Result: result}
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
	status <- &CmdStatus{Msg: "Deposit TX: Mining"}
	c.logOnChain("Sent deposit TX with hash %s", tx.Hash().Hex())
	res, err := bind.WaitMined(txCtx(), c.ethClient, tx)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	c.logOnChain("Deposit TX mined in Block #%d", res.BlockNumber.Uint64())
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
			c.logOffChain("Deposit Phase for Epoch %d ended in block %d", epoch, depEndBlock)
			time.Sleep(depositProofGrace)
		}
		waitErr <- err
	}()
	select {
	case e := <-waitErr:
		if e != nil {
			c.logOffChain("WaitForEpoch: %s", e.Error())
			status <- &CmdStatus{War: "Deposit proof: Chain error - resuming protocol"}
			c.events <- &Event{Type: SET_OP_TRUST, OpTrust: UNKNOWN}
		} else {
			status <- &CmdStatus{War: "Deposit proof: Operator timed out - resuming protocol"}
			c.events <- &Event{Type: SET_OP_TRUST, OpTrust: UNTRUSTED}
		}
	case p := <-proof:
		status <- &CmdStatus{Msg: "Deposit proof: Verifying"}
		ok, err := tee.VerifyDepositProof(*c.params, p)
		if !ok || err != nil {
			status <- &CmdStatus{War: "Deposit proof: Invalid Signature - resuming protocol"}
			c.events <- &Event{Type: SET_OP_TRUST, OpTrust: UNTRUSTED}
		} else if p.Balance.Value.Cmp(amount) != 0 || p.Balance.Epoch != epoch || p.Balance.Account != c.Address() {
			status <- &CmdStatus{War: "Deposit proof: Wrong Proof - resuming protocol"}
			c.events <- &Event{Type: SET_OP_TRUST, OpTrust: UNTRUSTED}
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
		c.events <- &Event{Type: SET_OP_TRUST, OpTrust: UNKNOWN}
	case <-ctx.Done():
		status <- &CmdStatus{War: "Deposit proof: Operator timed out - resuming protocol"}
		c.events <- &Event{Type: SET_OP_TRUST, OpTrust: UNKNOWN}
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
			c.logProof("Balance Proof error: %v", err)
			c.events <- &Event{Type: SET_OP_TRUST, OpTrust: UNKNOWN}
		} else if proof.Balance.Epoch > oldEpoch {
			c.logProof("Got Balance Proof for %v ETH in epoch %d", eth.WeiToEthFloat(proof.Balance.Value), proof.Balance.Epoch)
			oldEpoch = proof.Balance.Epoch
			c.balMtx.Lock()
			c.balances[proof.Balance.Epoch] = EpochBalance{Balance: proof.Balance, Bal: &proof}
			c.balMtx.Unlock()

			ok, err := tee.VerifyBalanceProof(*c.params, proof)
			if !ok || err != nil {
				c.events <- &Event{Type: SET_OP_TRUST, OpTrust: UNKNOWN}
				c.logProof("Invalid balance proof: err=%v ok=%t", err, ok)
				return
			}

			c.events <- &Event{Type: SET_BALANCE, Report: BalanceReport{Balance: new(big.Int).Set(proof.Balance.Value)}}
			c.events <- &Event{Type: SET_OP_TRUST, OpTrust: TRUSTED}
		}
		time.Sleep(time.Second)
	}
}

func (c *Client) lastBal() *EpochBalance {
	c.balMtx.Lock()
	defer c.balMtx.Unlock()
	epoch := c.params.ExitEpoch(atomic.LoadUint64(&c.currentBlock))
	bal, ok := c.balances[epoch]
	if !ok || bal.Bal == nil || bal.Bal.Sig == nil {
		return nil
	}
	return &bal
}

func (c *Client) CmdLeave(status chan *CmdStatus, args ...string) {
	defer close(status)
	if len(args) != 0 {
		status <- &CmdStatus{Err: errors.New("Command 'leave' does not accept arguments.")}
		return
	}
	bal := c.lastBal()
	if bal == nil {
		status <- &CmdStatus{Err: errors.New("No balance proof available")}
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
	status <- &CmdStatus{Msg: "Exit TX: Mining"}
	rec, err := bind.WaitMined(txCtx(), c.ethClient, tx)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	if rec.Status == types.ReceiptStatusFailed {
		status <- &CmdStatus{Err: errors.New("Exit TX: Receipt failed")}
		/*reason, err := errorReason(c.Ctx(), &c.ethClient.ContractBackend, tx, rec.BlockNumber, c.ethClient.Account())
		if err != nil {
			c.logOnChain("Unknown revert reason: %v", err)
		} else {
			c.logOnChain("Exit TX revert reason: %s", reason)
		}*/
		return
	}
	c.logOnChain("Exit mined in block #%d", rec.BlockNumber.Uint64())
	// Wait for the end of the Exit epoch before sending the Withdraw TX.
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
	status <- &CmdStatus{Msg: "Withdraw TX: Mining"}
	rec, err = bind.WaitMined(txCtx(), c.ethClient, tx)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	if rec.Status == types.ReceiptStatusFailed {
		status <- &CmdStatus{Err: errors.New("Withdraw TX: Receipt failed")}
		return
	}
	c.logOnChain("Withdraw mined in block #%d", rec.BlockNumber.Uint64())
	c.events <- &Event{Type: SET_BALANCE, Report: BalanceReport{Balance: big.NewInt(0)}}
	c.events <- &Event{Type: SET_OP_TRUST, OpTrust: TRUSTED}
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
			c.events <- &Event{Type: NEW_EPOCH, EpochNum: epoch}
		case block := <-blocks:
			atomic.StoreUint64(&c.currentBlock, block)
			c.events <- &Event{Type: NEW_BLOCK, BlockNum: block}
		case err := <-subError:
			return err
		}
		c.events <- &Event{Type: SET_EXIT_AVAIL, ExitAvailable: c.lastBal()}
	}
	return nil
}

func (c *Client) logProof(format string, args ...interface{}) {
	c.events <- &Event{Type: CHAIN_MSG, Message: "🔒 " + fmt.Sprintf(format, args...)}
}

func (c *Client) logOnChain(format string, args ...interface{}) {
	c.events <- &Event{Type: CHAIN_MSG, Message: "🔗 " + fmt.Sprintf(format, args...)}
}

func (c *Client) logOffChain(format string, args ...interface{}) {
	c.events <- &Event{Type: CHAIN_MSG, Message: "🍵 " + fmt.Sprintf(format, args...)}
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
