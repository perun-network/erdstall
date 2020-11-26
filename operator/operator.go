package operator

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	log "github.com/sirupsen/logrus"
	perrors "perun.network/go-perun/pkg/errors"

	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
	"github.com/perun-network/erdstall/tee/prototype"
)

// Operator resprents a TEE Plasma operator.
type Operator struct {
	enclave   tee.Enclave
	params    tee.Parameters
	ethClient *eth.Client
	*depositProofs
	*balanceProofs
	contract          *bindings.Erdstall
	respondChallenges bool
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
	respondChallenges bool,
) (*Operator, error) {
	_contract, err := bindings.NewErdstall(params.Contract, client)
	if err != nil {
		return nil, fmt.Errorf("loading contract: %w", err)
	}

	if !respondChallenges {
		log.Warn("Operator will not respond to on-chain challenges.")
	}

	return &Operator{
		enclave:           enclave,
		params:            params,
		ethClient:         client,
		depositProofs:     newDepositProofs(),
		balanceProofs:     newBalanceProofs(),
		contract:          _contract,
		respondChallenges: respondChallenges,
	}, nil
}

// Setup creates an operator from the given configuration.
func Setup(cfg *Config) *Operator {
	wallet, err := hdwallet.NewFromMnemonic(cfg.Mnemonic)
	AssertNoError(err)

	enclaveAccountDerivationPath := hdwallet.MustParseDerivationPath(cfg.EnclaveDerivationPath)
	enclaveAccount, err := wallet.Derive(enclaveAccountDerivationPath, true)
	AssertNoError(err)
	log.Debug("Operator.Setup: Enclave account loaded")

	enclave := prototype.NewEnclaveWithAccount(wallet, enclaveAccount)
	enclavePublicKey, _, err := enclave.Init()
	AssertNoError(err)
	log.Info("Operator.Setup: Enclave created")

	operatorAccountDerivationPath := hdwallet.MustParseDerivationPath(cfg.OperatorDerivationPath)
	operatorAccount, err := wallet.Derive(operatorAccountDerivationPath, true)
	AssertNoError(err)

	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()
	client, err := eth.CreateEthereumClient(ctx, cfg.EthereumNodeURL, wallet, operatorAccount)
	AssertNoError(err)
	log.Info("Operator.Setup: Ethereum client initialized")

	params := tee.Parameters{
		TEE:              enclavePublicKey,
		PhaseDuration:    cfg.PhaseDuration,
		ResponseDuration: cfg.ResponseDuration,
		PowDepth:         cfg.PowDepth,
	}

	err = client.DeployContracts(&params)
	AssertNoError(err)
	log.Infof("Operator.Setup: Contract deployed at %s", params.Contract.String())

	operator, err := New(enclave, params, client, cfg.RespondChallenges)
	AssertNoError(err)

	return operator
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
			}
			return err
		})
	}

	// Start enclave
	errGo("Enclave.Run", func() error { return operator.enclave.Run(operator.params) })
	log.Info("Operator.Serve: Enclave running")

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

	// Handle RPC
	errGo("Op.RPCServe", func() error { return operator.handleRPC(port) })
	log.Info("Operator.Serve: RPC handling started")

	return errg.Wait()
}

func (operator *Operator) handleBlocks() error {
	bigBang, err := operator.contract.BigBang(nil)
	if err != nil {
		return fmt.Errorf("reading BigBang: %w", err)
	}
	blockSub, err := operator.ethClient.SubscribeBlocksStartingFrom(new(big.Int).SetUint64(bigBang))
	if err != nil {
		return fmt.Errorf("creating block subscription: %w", err)
	}
	defer blockSub.Unsubscribe()
	for b := range blockSub.Blocks() {
		log.Debugf("Operator.Serve: incoming block %d", b.NumberU64())
		if err := operator.enclave.ProcessBlocks(b); err != nil {
			//TODO check for ErrEnclaveStopped error, see enclave internal tests
			return err
		}
		log.Debugf("Operator.Serve: processed block %d", b.NumberU64())
	}
	return nil
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
		}
	}
}

func (operator *Operator) handleChallengedEvent(c challengedEvent) error {
	if !operator.respondChallenges {
		log.Warn("Operator.handleChallengedEvent: ignoring challenges, returning.")
		return nil
	}

	ctx, cancel := createDefaultContext()
	defer cancel()

	tr, err := operator.ethClient.NewTransactor(ctx)
	if err != nil {
		return fmt.Errorf("creating transactor: %w", err)
	}

	balanceProof, ok := operator.balanceProofs.Get(c.Account)
	if !ok {
		return errors.New("getting balance proof")
	}

	balance := bindings.ErdstallBalance{
		Epoch:   balanceProof.Balance.Epoch,
		Account: balanceProof.Balance.Account,
		Value:   balanceProof.Balance.Value,
	}

	tx, err := operator.contract.Exit(tr, balance, balanceProof.Sig)
	if err != nil {
		return fmt.Errorf("sending challenge response: %w", err)
	}

	// Track challenge response transaction status
	go func() {
		ctx, cancel := createOnChainContext()
		defer cancel()

		r, err := bind.WaitMined(ctx, operator.ethClient.ContractBackend, tx)
		if err != nil {
			log.Errorf("Operator.handleChallengedEvent: Failed to wait for mining of response to challenge %v: %v", c, err)
			return
		}

		if r.Status != types.ReceiptStatusSuccessful {
			log.Errorf("Operator.handleChallengedEvent: Failed to complete response transaction for challenge %v", c)
			return
		}

		log.Infof("Operator.handleChallengedEvent: Resolved dispute for challenge %v", c)
	}()

	return nil
}

func (operator *Operator) handleDepositProofs() error {
	for {
		dps, err := operator.enclave.DepositProofs()
		if err != nil {
			return fmt.Errorf("retrieving deposit proofs: %w", err)
		}
		if len(dps) > 0 {
			log.Debugf("Operator.Serve: Retrieved %d deposit proofs", len(dps))
		}
		operator.depositProofs.AddAll(dps)
	}
}

func (operator *Operator) handleBalanceProofs() error {
	for {
		bps, err := operator.enclave.BalanceProofs()
		if err != nil {
			return fmt.Errorf("retrieving balance proofs: %w", err)
		}
		if len(bps) > 0 {
			log.Debugf("Operator.Serve: Retrieved %d balance proofs", len(bps))
		}
		operator.balanceProofs.AddAll(bps)
	}
}

func (operator *Operator) handleRPC(port int) error {
	remoteEnclave := newRemoteEnclave(operator)
	if err := rpc.Register(remoteEnclave); err != nil {
		return fmt.Errorf("registering remote enclave interface: %w", err)
	}
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("binding to socket: %w", err)
	}
	return http.Serve(l, nil)
}

// AssertNoError logs the error and exits if the error is not nil.
func AssertNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createDefaultContext() (context.Context, context.CancelFunc) {
	//todo: make on-chain context configurable
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func createOnChainContext() (context.Context, context.CancelFunc) {
	//todo: make on-chain context configurable
	return context.WithTimeout(context.Background(), 60*time.Second)
}

// challengedEvent provides formatting for bindings.ErdstallChallenged.
type challengedEvent bindings.ErdstallChallenged

func (c challengedEvent) String() string {
	return fmt.Sprintf("ChallengedEvent{Account: %v, Epoch: %v}", c.Account.String(), c.Epoch)
}
