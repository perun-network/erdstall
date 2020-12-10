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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

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
	Config       config.ClientConfig
	conn         *RPC
	proofSub     *Subscription
	ethClient    *eth.Client
	contractAddr common.Address
	signer       tee.TextSigner
	txNonce      uint64
	balances     map[uint64]EpochBalance // epoch => balance
	// Initialized in Run()
	lastBlock uint64 // Atomic
	contract  *bindings.Erdstall
	params    *tee.Parameters
	events    chan *Event
	// balMtx protects balances.
	balMtx            sync.Mutex
	stopFrozenWatcher context.CancelFunc
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
	c.setOpTrust(TRUSTED)

	c.params = params
	c.contract = contract
	c.proofSub, err = c.conn.Subscribe(newCtx(5*time.Second), c.Address())
	if err != nil {
		return fmt.Errorf("subscribing to proofs: %w", err)
	}
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
	tx, err := c.createTransfer(receiver, amount)
	if err != nil {
		status <- &CmdStatus{Err: err}
		return
	}
	status <- &CmdStatus{Msg: "Forwarding to Operator"}
	if err := c.conn.SendTx(shortCtx(), tx); err != nil {
		status <- &CmdStatus{Err: err}
	}
}

func (c *Client) createTransfer(receiver common.Address, amount *big.Int) (tee.Transaction, error) {
	block := atomic.LoadUint64(&c.lastBlock) + 1

	tx := tee.Transaction{
		Nonce:     c.txNonce,
		Epoch:     c.params.TxEpoch(block),
		Sender:    c.Address(),
		Recipient: receiver,
		Amount:    (*tee.Amount)(amount),
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
	if len(args) > 3 {
		status <- &CmdStatus{Err: errors.New("Command 'bench' can have arguments: <runs> <address> <amount>")}
		return
	}
	var err error
	n := 1000
	if len(args) >= 1 {
		n, err = strconv.Atoi(args[0])
		if err != nil {
			status <- &CmdStatus{Err: fmt.Errorf("Could not parse <runs> as int: %s", args[0])}
			return
		}
	}
	a := c.Address()
	if len(args) >= 2 {
		a = common.HexToAddress(args[1])
	}
	amount := big.NewInt(1)
	if len(args) >= 3 {
		totalAmount, err := strconv.Atoi(args[2])
		if err != nil {
			status <- &CmdStatus{Err: fmt.Errorf("Could not parse <amount> as int: %s: %w", args[2], err)}
			return
		}
		amount = new(big.Int).Div(eth.EthToWeiInt(int64(totalAmount)), big.NewInt(int64(n)))
	}
	status <- &CmdStatus{Msg: fmt.Sprintf("Sending %d payments", n)}
	result, err := Benchmark(n, func() error {
		tx, err := c.createTransfer(a, amount)
		if err != nil {
			return err
		}
		return c.conn.SendTx(shortCtx(), tx)
	})
	c.events <- &Event{Type: BENCH, Result: result}
	if err != nil {
		status <- &CmdStatus{Err: err}
	}
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
	rec, err := c.sendTx("Deposit", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		opts.Value = amount
		return c.contract.Deposit(opts)
	}, status)
	if err != nil {
		status <- &CmdStatus{Err: fmt.Errorf("Deposit TX: %w", err)}
		return
	}

	c.logOnChain("Deposit TX mined in Block #%d", rec.BlockNumber.Uint64())
	status <- &CmdStatus{Msg: "Deposit proof: Waiting for Operator"}
	blockNum := rec.BlockNumber.Uint64()
	// The epoch that we want to do the deposit in.
	epoch := c.params.DepositEpoch(blockNum)

	proof := make(chan tee.DepositProof)
	proofErr := make(chan error)
	waitErr := make(chan error)
	// no specific wait ctx, since we do not know what the block time is.
	ctx, cancel := context.WithCancel(c.Ctx())
	defer cancel()
	go func() {
		if p, err := c.proofSub.DepositProof(ctx); err != nil {
			if p.Balance.Epoch != epoch {
				proofErr <- fmt.Errorf("Got proof for wrong epoch #%d", p.Balance.Epoch)
			} else {
				proofErr <- err
			}
		} else {
			proof <- p
		}
	}()
	go func() {
		depEndBlock := c.params.TxDoneBlock(epoch)
		c.logOffChain("Waiting for deposit proof until block #%d", depEndBlock)
		err := c.ethClient.WaitForBlock(ctx, depEndBlock) // Add PowDepth here if needed
		if err == nil {
			c.logOffChain("Deposit Phase for Epoch %d ended in block %d", epoch, depEndBlock)
		}
		waitErr <- err
	}()
	select {
	case e := <-waitErr:
		if e != nil {
			c.logOffChain("WaitForEpoch: %s", e.Error())
			status <- &CmdStatus{War: "Deposit proof: Chain error - resuming protocol"}
			c.setOpTrust(UNKNOWN)
		} else {
			status <- &CmdStatus{War: "Deposit proof: Operator timed out - resuming protocol"}
			c.setOpTrust(UNTRUSTED)
		}
	case p := <-proof:
		status <- &CmdStatus{Msg: "Deposit proof: Verifying"}
		ok, err := tee.VerifyDepositProof(*c.params, p)
		if !ok || err != nil {
			status <- &CmdStatus{War: "Deposit proof: Invalid Signature - resuming protocol"}
			c.setOpTrust(UNTRUSTED)
		} else if (*big.Int)(p.Balance.Value).Cmp(amount) != 0 || p.Balance.Epoch != epoch || p.Balance.Account != c.Address() {
			status <- &CmdStatus{War: "Deposit proof: Wrong Proof - resuming protocol"}
			c.setOpTrust(UNTRUSTED)
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
		c.setOpTrust(UNKNOWN)
	case <-ctx.Done():
		status <- &CmdStatus{War: "Deposit proof: Operator timed out - resuming protocol"}
		c.setOpTrust(UNKNOWN)
	}

	c.challengeDeposit(status)
}

// BalanceProofWatcher waits for the balance proof of an epoch and disputes
// if non was received after the TxPhase + balanceProofGrace.
// Can currently only deal with one deposit per deposit-phase.
// Should be started in a go-routine.
func (c *Client) BalanceProofWatcher() {
	oldEpoch := uint64(0)
	for !c.IsClosed() {
		proof, err := c.proofSub.BalanceProof(c.Ctx())
		if err != nil {
			c.logProof("Balance Proof error: %v", err)
			c.setOpTrust(UNKNOWN)
		} else if proof.Balance.Epoch > oldEpoch {
			c.logProof("Got Balance Proof for %v ETH in epoch %d", eth.WeiToEthFloat((*big.Int)(proof.Balance.Value)), proof.Balance.Epoch)
			oldEpoch = proof.Balance.Epoch
			c.balMtx.Lock()
			c.balances[proof.Balance.Epoch] = EpochBalance{Balance: proof.Balance, Bal: &proof}
			c.balMtx.Unlock()

			ok, err := tee.VerifyBalanceProof(*c.params, proof)
			if !ok || err != nil {
				c.setOpTrust(UNKNOWN)
				c.logProof("Invalid balance proof: err=%v ok=%t", err, ok)
				return
			}

			c.events <- &Event{Type: SET_BALANCE, Report: BalanceReport{Balance: new(big.Int).Set((*big.Int)(proof.Balance.Value))}}
			c.setOpTrust(TRUSTED)
		}
		time.Sleep(time.Second)
	}
}

// FrozenWatcher listens for Frozen events and calls WithdrawFrozen
// if a balance proof is available.
func (c *Client) FrozenWatcher(_ctx context.Context) {
	ctx, cancel := context.WithCancel(_ctx)
	defer cancel()
	sub, err := c.ethClient.SubscribeFrozen(ctx, c.contract, nil)
	if err != nil {
		c.logError("Frozen event subscription: %v", err)
		return
	}
	defer sub.Unsubscribe()

	select {
	case <-ctx.Done():
		c.logOnChain("Frozen watcher stopped")
		return
	case err := <-sub.Err():
		c.logError("Frozen event subscription: %v", err)
		return
	case event := <-sub.Events():
		c.handleFrozen(event)
	}
}

func (c *Client) handleFrozen(event *bindings.ErdstallFrozen) {
	epoch := event.Epoch
	c.setOpTrust(UNTRUSTED)
	c.log("❄️Contract Frozen in epoch #%d", epoch)

	status := make(chan *CmdStatus)
	defer close(status)
	go func() {
		for msg := range status {
			if msg.Err != nil {
				c.logOnChain(msg.Err.Error())
			} else {
				c.logOnChain(msg.Msg)
			}
		}
	}()
	c.withdrawFrozen(status, epoch)
}

func (c *Client) withdrawFrozen(status chan *CmdStatus, epoch uint64) {
	c.balMtx.Lock()
	defer c.balMtx.Unlock()
	bal, ok := c.balances[epoch]
	if !ok {
		c.logOffChain("No balance-proof available for freeze")
		return
	}

	_, err := c.sendTx("WithdrawFrozen", func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.contract.WithdrawFrozen(auth, bal.ToEthBals(), bal.Bal.Sig)
	}, status)
	if err != nil {
		status <- &CmdStatus{Err: fmt.Errorf("WithdrawFrozen TX: %w", err)}
		return
	}
	c.logOnChain("❄️WithdrawFrozen: Complete")
}

func (c *Client) lastBal() *EpochBalance {
	c.balMtx.Lock()
	defer c.balMtx.Unlock()
	epoch := c.params.ExitEpoch(atomic.LoadUint64(&c.lastBlock))
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

	rec, err := c.sendTx("Exit", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return c.contract.Exit(opts, bal.ToEthBals(), bal.Bal.Sig)
	}, status)
	if err != nil {
		status <- &CmdStatus{Err: fmt.Errorf("Exit TX: %w", err)}
		return
	}
	c.logOnChain("Exit mined in block #%d", rec.BlockNumber.Uint64())
	c.withdraw(bal.Epoch, status)
}

func (c *Client) CmdChallenge(status chan *CmdStatus, args ...string) {
	defer close(status)
	if len(args) != 0 {
		status <- &CmdStatus{Err: errors.New("Command 'challenge' does not accept arguments.")}
		return
	}

	// Wait for the next epoch, when we are currently in the response phase.
	lastBlock := atomic.LoadUint64(&c.lastBlock)
	if c.params.IsChallengeResponsePhase(lastBlock + 1) {
		status <- &CmdStatus{Msg: "Waiting for next epoch"}
		next := c.params.DepositDoneBlock(c.params.DepositEpoch(lastBlock + 1))
		if err := c.ethClient.WaitForBlock(c.Ctx(), next); err != nil {
			status <- &CmdStatus{Err: err}
			return
		}
	}

	// Get the BP of the sealed epoch.
	lastBlock = atomic.LoadUint64(&c.lastBlock) + 1
	sealed := c.params.SealedEpoch(lastBlock)
	c.balMtx.Lock()
	bal, ok := c.balances[sealed]
	c.balMtx.Unlock()
	if !ok || bal.Bal == nil {
		status <- &CmdStatus{Err: fmt.Errorf("No Balance proof available")}
		return
	}

	c.challenge(status, bal.Bal)
}

func (c *Client) challengeDeposit(status chan *CmdStatus) {
	tx, err := c.sendTx("ChallengeDeposit", func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.contract.ChallengeDeposit(auth)
	}, status)
	if err != nil {
		status <- &CmdStatus{Err: fmt.Errorf("ChallengeDeposit TX: %w", err)}
		return
	}
	c.waitForReponse(status, tx.BlockNumber.Uint64())
}

func (c *Client) challenge(status chan *CmdStatus, bp *tee.BalanceProof) {
	tx, err := c.sendTx("Challenge", func(auth *bind.TransactOpts) (*types.Transaction, error) {
		return c.contract.Challenge(auth, bp.Balance.ToEthBals(), bp.Sig)
	}, status)
	if err != nil {
		status <- &CmdStatus{Err: fmt.Errorf("Challenge TX: %w", err)}
		return
	}
	c.waitForReponse(status, tx.BlockNumber.Uint64())
}

func (c *Client) waitForReponse(status chan *CmdStatus, block uint64) {
	// Wait for the operator til the end of the epoch and Freeze otherwise.
	status <- &CmdStatus{Msg: "Exiting event: Waiting"}
	subCtx, cancel := context.WithCancel(c.Ctx())
	defer cancel()
	exitEpoch := c.params.ExitEpoch(block)
	sub, err := c.ethClient.SubscribeExiting(subCtx, c.contract, []uint64{exitEpoch}, []common.Address{c.Address()})
	if err != nil {
		status <- &CmdStatus{Err: fmt.Errorf("Exiting subscription: %w", err)}
		return
	}
	defer sub.Unsubscribe()

	done := make(chan struct{})
	go func() {
		next := c.params.DepositDoneBlock(c.params.DepositEpoch(block))
		c.ethClient.WaitForBlock(subCtx, next) // nolint: errcheck
		close(done)
	}()

	select {
	case <-done: // ChallengeReponse phase is over, freeze.
		c.logOnChain("Freezin in epoch %d", exitEpoch)
		c.stopFrozenWatcher()
		_, err := c.sendTx("WithdrawChallenge", func(auth *bind.TransactOpts) (*types.Transaction, error) {
			return c.contract.WithdrawChallenge(auth)
		}, status)
		if err != nil {
			status <- &CmdStatus{Err: fmt.Errorf("WithdrawChallenge TX: %w", err)}
			return
		}
	case err := <-sub.Err():
		if err != nil {
			c.logError("Exiting subscription: %v", err)
		}
		return
	case <-sub.Events(): // Challenge was posted on-chain by OP.
		c.logOnChain("Received on-chain challenge response")
		c.withdraw(exitEpoch, status)
	}
}

// sendTx already checks the receipt status and returns an error
// when it failed.
func (c *Client) sendTx(name string, f func(*bind.TransactOpts) (*types.Transaction, error), status chan *CmdStatus) (*types.Receipt, error) {
	status <- &CmdStatus{Msg: name + " TX: Preparing"}
	opts, err := c.ethClient.NewTransactor(txCtx())
	if err != nil {
		return nil, err
	}
	status <- &CmdStatus{Msg: name + " TX: Sending"}
	tx, err := f(opts)
	if err != nil {
		return nil, err
	}
	status <- &CmdStatus{Msg: name + " TX: Mining"}

	return c.ethClient.ConfirmTransaction(txCtx(), tx, c.ethClient.Account())
}

func (c *Client) withdraw(exitEpoch uint64, status chan *CmdStatus) {
	// Wait for the end of the Exit epoch before sending the Withdraw TX.
	endExitBlock := c.params.DepositStartBlock(exitEpoch + 1)
	if err := c.ethClient.WaitForBlock(c.Ctx(), endExitBlock); err != nil {
		status <- &CmdStatus{Err: err}
		return
	}

	rec, err := c.sendTx("Withdraw", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return c.contract.Withdraw(opts, exitEpoch)
	}, status)
	if err != nil {
		status <- &CmdStatus{Err: fmt.Errorf("Withdraw TX: %w", err)}
		return
	}

	c.logOnChain("Withdraw mined in block #%d", rec.BlockNumber.Uint64())
	c.events <- &Event{Type: SET_BALANCE, Report: BalanceReport{Balance: big.NewInt(0)}}
	c.setOpTrust(TRUSTED)
	c.stopFrozenWatcher()
}

// writes to chainVvents
func (c *Client) listenOnChain() error {
	epochs := make(chan uint64)
	blocks := make(chan uint64)
	subError := make(chan error)

	go func() {
		subError <- c.ethClient.SubscribeEpochs(c.Ctx(), *c.params, epochs, blocks)
	}()
	go c.BalanceProofWatcher()
	frozenCtx, cancel := context.WithCancel(c.Ctx())
	c.stopFrozenWatcher = cancel
	go c.FrozenWatcher(frozenCtx)
	c.OnCloseAlways(cancel)

	for !c.IsClosed() {
		select {
		case epoch := <-epochs:
			c.events <- &Event{Type: NEW_EPOCH, EpochNum: epoch}
		case block := <-blocks:
			atomic.StoreUint64(&c.lastBlock, block)
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

func (c *Client) logError(format string, args ...interface{}) {
	c.events <- &Event{Type: CHAIN_MSG, Message: "⚠ " + fmt.Sprintf(format, args...)}
}

func (c *Client) log(format string, args ...interface{}) {
	c.events <- &Event{Type: CHAIN_MSG, Message: fmt.Sprintf(format, args...)}
}

func (c *Client) setOpTrust(trust Trust) {
	c.events <- &Event{Type: SET_OP_TRUST, OpTrust: trust}
}

func (c *Client) Address() common.Address {
	return c.ethClient.Account().Address
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
