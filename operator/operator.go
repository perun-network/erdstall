// SPDX-License-Identifier: Apache-2.0

package operator

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	log "github.com/sirupsen/logrus"
	perrors "perun.network/go-perun/pkg/errors"
	pkgsync "perun.network/go-perun/pkg/sync"

	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
	"github.com/perun-network/erdstall/tee/prototype"
	"github.com/perun-network/erdstall/tee/rpc"
)

// Operator resprents a TEE Plasma operator.
type Operator struct {
	pkgsync.Closer
	enclave   tee.Enclave
	params    tee.Parameters
	EthClient *eth.Client
	*depositProofs
	*balanceProofs
	rpcOperator *RPCOperator
	contract    *bindings.Erdstall
	cfg         Config
}

// EnclaveParams returns the enclave parameters.
func (operator *Operator) EnclaveParams() tee.Parameters {
	return operator.params
}

// New instantiates an operator with the given parameters.
func New(
	enclave tee.Enclave,
	params tee.Parameters,
	client *eth.Client,
	cfg Config,
) (*Operator, error) {
	_contract, err := bindings.NewErdstall(params.Contract, client)
	if err != nil {
		return nil, fmt.Errorf("loading contract: %w", err)
	}

	return &Operator{
		enclave:       enclave,
		params:        params,
		EthClient:     client,
		depositProofs: newDepositProofs(),
		balanceProofs: newBalanceProofs(),
		contract:      _contract,
		cfg:           cfg,
	}, nil
}

// SetupWithPrototypeEnclave creates an operator from the given configuration
// using a prototype.Enclave.
func SetupWithPrototypeEnclave(cfg *Config) *Operator {
	wallet, err := hdwallet.NewFromMnemonic(cfg.Mnemonic)
	AssertNoError(err)
	enclaveAccountDerivationPath := hdwallet.MustParseDerivationPath(cfg.EnclaveDerivationPath)
	enclaveAccount, err := wallet.Derive(enclaveAccountDerivationPath, true)
	AssertNoError(err)
	log.WithField("enclave", enclaveAccount.Address.Hex()).
		Debug("Operator.Setup: account loaded")
	enclave := prototype.NewEnclaveWithAccount(wallet, enclaveAccount)
	return Setup(cfg, enclave)
}

// Setup creates an operator with the requested enclave. If the
// enclave is nil, creates a new prototype enclave.
func Setup(cfg *Config, enclave tee.Enclave) *Operator {
	wallet, err := hdwallet.NewFromMnemonic(cfg.Mnemonic)
	AssertNoError(err)

	enclavePublicKey, _, err := enclave.Init()
	AssertNoError(err)
	log.Info("Operator.Setup: Enclave created")

	operatorAccountDerivationPath := hdwallet.MustParseDerivationPath(cfg.OperatorDerivationPath)
	operatorAccount, err := wallet.Derive(operatorAccountDerivationPath, true)
	AssertNoError(err)

	log.WithFields(log.Fields{
		"op": operatorAccount.Address.Hex()}).Info("Operator.Setup: Accounts loaded")

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	client, err := eth.CreateEthereumClient(ctx, cfg.EthereumNodeURL, wallet, operatorAccount)
	AssertNoError(err)
	// Skip retrieving other receipts as they're currently not checked in the
	// prototype enclave...
	client.OnlyErdstallReceipts()
	log.Info("Operator.Setup: Ethereum client initialized")

	var params tee.Parameters
	if cfg.ContractAddr != "" {
		if !common.IsHexAddress(cfg.ContractAddr) {
			log.Fatalf("Config: No hex address: %s", cfg.ContractAddr)
		}
		contract := common.HexToAddress(cfg.ContractAddr)
		log.Infof("Operator.Setup: Binding contract at %s", contract.String())
		ctx, cancel := eth.ContextNodeReq()
		defer cancel()
		p, _, err := client.BindContract(ctx, contract)
		AssertNoError(err)
		params = *p
	} else {
		log.Infof("Operator.Setup: Deploying contract...")
		params = tee.Parameters{
			TEE:              enclavePublicKey,
			PhaseDuration:    cfg.PhaseDuration,
			ResponseDuration: cfg.ResponseDuration,
			PowDepth:         cfg.PowDepth,
		}
		err = client.DeployContracts(&params)
		AssertNoError(err)
		log.Infof("Operator.Setup: Contract deployed at %s", params.Contract.String())
	}

	if !cfg.RespondChallenges || !cfg.SendDepositProofs || !cfg.SendBalanceProofs {
		log.Warnf(
			"Operator will respond to challenges: %t, send out deposit proofs: %t, send out balance proofs: %t",
			cfg.RespondChallenges, cfg.SendDepositProofs, cfg.SendBalanceProofs)
	}

	operator, err := New(enclave, params, client, *cfg)
	AssertNoError(err)

	return operator
}

// SetupWithRemoteEnclave creates an operator setup that dials the specified remote
// enclave.
func SetupWithRemoteEnclave(
	cfg *Config,
	enclaveAddr string,
) (op *Operator, err error) {
	e, err := rpc.DialEnclave(enclaveAddr)
	if err != nil {
		return nil, err
	}

	return Setup(cfg, e), nil
}

// Serve starts the operator's main routine.
func (operator *Operator) Serve(port uint16) error {
	// Handle errors, print them as they occurr
	errg := perrors.NewGatherer()
	errGo := func(name string, fn func() error) {
		errg.Go(func() error {
			err := fn()
			if err != nil {
				log.Errorf("Error in %s: %v", name, err)
			} else {
				log.Debugf("%s returned.", name)
			}
			return err
		})
	}

	// Start enclave
	errGo("Enclave.Run", func() error {
		defer operator.Close() // Make sure that operator closes on a failing enclave.
		return operator.enclave.Run(operator.params)
	})
	log.Info("Operator.Serve: Enclave running")

	operator.rpcOperator = NewRPCOperator(operator.enclave)

	clientConfig := ClientConfig{
		Contract:  operator.params.Contract,
		NetworkID: operator.EthClient.NetworkID.String(),
		POWDepth:  operator.params.PowDepth,
	}
	osc := OpServerConfig{
		Host:         "0.0.0.0",
		Port:         port,
		CertFilePath: operator.cfg.CertFile,
		KeyFilePath:  operator.cfg.KeyFile,
		ClientConfig: clientConfig,
	}

	// Handle RPC
	errGo("Op.RPCServe", func() error {
		rpc := NewRPC(operator.rpcOperator, osc)
		operator.OnClose(func() {
			rpc.Close()
		})
		return rpc.Serve()
	})
	log.Info("Operator.Serve: RPC handling started")

	// Handle Ethereum blocks
	errGo("Op.BlockSub", operator.handleBlocks)
	log.Info("Operator.Serve: Block subcription started")

	// Handle on-chain challenges
	errGo("Op.Challenges", operator.handleChallenges)
	log.Info("Operator.Serve: Challenge handling started")

	// Handle deposit proofs
	errGo("Op.DepositProofs", operator.handleDepositProofs)
	log.Info("Operator.Serve: Deposit proof handling started")

	// Handle balance proofs
	errGo("Op.BalanceProofs", operator.handleBalanceProofs)
	log.Info("Operator.Serve: Balance proof handling started")

	return errg.Wait()
}

func (operator *Operator) handleBlocks() error {
	blockSub, err := operator.EthClient.SubscribeVerifiedBlocksFrom(operator.params.InitBlock)
	if err != nil {
		return fmt.Errorf("creating block subscription: %w", err)
	}
	defer blockSub.Unsubscribe()
	for {
		select {
		case b := <-blockSub.Blocks():
			log.Infof("Operator.Serve: incoming block %d", b.NumberU64())
			if err := operator.enclave.ProcessBlocks(b); err != nil {
				//TODO check for ErrEnclaveStopped error, see enclave internal tests
				return err
			}
			log.Debugf("Operator.Serve: processed block %d", b.NumberU64())
		case <-operator.Closed():
			return nil
		}
	}
}

func (operator *Operator) handleChallenges() error {
	challenges := make(chan *bindings.ErdstallChallenged)
	sub, err := operator.contract.WatchChallenged(nil, challenges, nil, nil)
	if err != nil {
		return fmt.Errorf("creating challenge subcription: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err)

		case _c := <-challenges:
			c := challengedEvent(*_c)
			log.Warnf("Operator.handleChallenges: Incoming challenge %v", c)

			if err := operator.handleChallengedEvent(c); err != nil {
				log.Errorf("Operator.handleChallenges: Failed to handle challenged event %v: %v", c, err)
			}
		case <-operator.Closed():
			return nil
		}
	}
}

func (operator *Operator) handleChallengedEvent(c challengedEvent) error {
	if !operator.cfg.RespondChallenges {
		log.Warn("Operator.handleChallengedEvent: ignoring challenges, returning.")
		return nil
	}

	balanceProof, ok := operator.balanceProofs.Get(c.Account)
	if !ok {
		return errors.New("getting balance proof")
	}

	ctx, cancel := eth.ContextNodeReq()
	defer cancel()

	tr, err := operator.EthClient.NewTransactor(ctx)
	if err != nil {
		return fmt.Errorf("creating transactor: %w", err)
	}

	tx, err := operator.contract.Exit(tr, balanceProof.Balance.ToEthBal(), balanceProof.Sig)
	if err != nil {
		return fmt.Errorf("sending challenge response: %w", err)
	}

	// Track challenge response transaction status
	go func() {
		ctx, cancel := eth.ContextWaitMined()
		defer cancel()

		if _, err := operator.EthClient.ConfirmTransaction(ctx, tx, operator.EthClient.Account()); err != nil {
			log.Errorf("Operator.handleChallengedEvent: Failed to wait for mining of response to challenge %v: %v", c, err)
			return
		}

		log.Infof("Operator.handleChallengedEvent: Resolved dispute for challenge %v", c)
	}()

	return nil
}

func (operator *Operator) handleDepositProofs() error {
	if !operator.cfg.SendDepositProofs {
		log.Warn("Ignoring deposit proofs")
		return nil
	}

	for {
		dps, err := operator.enclave.DepositProofs()
		if err != nil {
			return fmt.Errorf("retrieving deposit proofs: %w", err)
		}
		if len(dps) > 0 {
			log.Debugf("Operator.Serve: Retrieved %d deposit proofs", len(dps))
		}
		operator.depositProofs.AddAll(dps)
		for _, dp := range dps {
			operator.rpcOperator.PushDepositProof(*dp)
		}
	}
}

func (operator *Operator) handleBalanceProofs() error {
	if !operator.cfg.SendDepositProofs {
		log.Warn("Ignoring balance proofs")
		return nil
	}

	for {
		bps, err := operator.enclave.BalanceProofs()
		if err != nil {
			return fmt.Errorf("retrieving balance proofs: %w", err)
		}
		if len(bps) > 0 {
			log.Debugf("Operator.Serve: Retrieved %d balance proofs", len(bps))
		}
		operator.balanceProofs.AddAll(bps)
		for _, bp := range bps {
			operator.rpcOperator.PushBalanceProof(*bp)
		}
	}
}

// AssertNoError logs the error and exits if the error is not nil.
func AssertNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// challengedEvent provides formatting for bindings.ErdstallChallenged.
type challengedEvent bindings.ErdstallChallenged

func (c challengedEvent) String() string {
	return fmt.Sprintf("ChallengedEvent{Account: %v, Epoch: %v}", c.Account.String(), c.Epoch)
}
