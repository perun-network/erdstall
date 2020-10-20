package operator

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/perun-network/erdstall/contracts/bindings"
	"github.com/perun-network/erdstall/eth"
	"github.com/perun-network/erdstall/tee"
	"github.com/perun-network/erdstall/tee/prototype"
	log "github.com/sirupsen/logrus"
	"perun.network/go-perun/pkg/errors"
)

// Operator resprents a TEE Plasma operator.
type Operator struct {
	enclave   *prototype.Enclave
	ethClient *eth.Client
	*depositProofs
	*balanceProofs
	contract *bindings.Erdstall
}

// New instantiates an operator with the given parameters.
func New(
	enclave *prototype.Enclave,
	client *eth.Client,
	contract common.Address,
) (*Operator, error) {
	_contract, err := bindings.NewErdstall(contract, client)
	if err != nil {
		return nil, fmt.Errorf("loading contract: %w", err)
	}

	return &Operator{
		enclave,
		client,
		newDepositProofs(),
		newBalanceProofs(),
		_contract,
	}, nil
}

// Setup creates an operator from the given configuration.
func Setup(cfg *Config) (*Operator, tee.Parameters) {
	wallet, err := hdwallet.NewFromMnemonic(cfg.Mnemonic)
	AssertNoError(err)

	enclaveAccountDerivationPath := hdwallet.MustParseDerivationPath(cfg.EnclaveDerivationPath)
	enclaveAccount, err := wallet.Derive(enclaveAccountDerivationPath, true)
	AssertNoError(err)
	log.Debug("Operator.Setup: Enclave account loaded")

	enclave := prototype.NewEnclaveWithAccount(wallet, enclaveAccount)
	enclavePublicKey, _, err := enclave.Init()
	AssertNoError(err)
	log.Debug("Operator.Setup: Enclave created")

	operatorAccountDerivationPath := hdwallet.MustParseDerivationPath(cfg.OperatorDerivationPath)
	operatorAccount, err := wallet.Derive(operatorAccountDerivationPath, true)
	AssertNoError(err)

	client, err := eth.CreateEthereumClient(cfg.EthereumNodeURL, wallet, operatorAccount)
	AssertNoError(err)
	log.Debug("Operator.Setup: Ethereum client initialized")

	enclaveParameters := tee.Parameters{
		TEE:              enclavePublicKey,
		PhaseDuration:    cfg.PhaseDuration,
		ResponseDuration: cfg.ResponseDuration,
		PowDepth:         cfg.PowDepth,
	}

	err = client.DeployContracts(&enclaveParameters)
	AssertNoError(err)
	log.Debugf("Operator.Setup: Contract deployed at %x", enclaveParameters.Contract)

	err = enclave.SetParams(enclaveParameters)
	AssertNoError(err)
	log.Debug("Operator.Setup: Enclave initialized")

	operator, err := New(enclave, client, enclaveParameters.Contract)
	AssertNoError(err)

	return operator, enclaveParameters
}

// Serve starts the operator's main routine.
func (operator *Operator) Serve(port int) error {
	errg := errors.NewGatherer()

	// Start enclave
	errg.Go(operator.enclave.Start)
	log.Debug("Operator.Serve: Enclave started")

	bigBang, err := operator.contract.BigBang(nil)
	if err != nil {
		return fmt.Errorf("reading BigBang: %w", err)
	}

	// Ethereum block handling
	blockSub, err := operator.ethClient.SubscribeToBlocksStartingFrom(new(big.Int).SetUint64(bigBang))
	if err != nil {
		return fmt.Errorf("creating block subscription: %w", err)
	}
	errg.Go(func() error {
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
	})

	// Handle deposit proofs
	errg.Go(func() error {
		for {
			dps, err := operator.enclave.DepositProofs()
			if err != nil {
				return fmt.Errorf("retrieving deposit proofs: %w", err)
			}
			if len(dps) > 0 {
				log.Infof("Operator.Serve: Retrieved %d deposit proofs", len(dps))
			}
			operator.depositProofs.AddAll(dps)
		}
	})

	// Handle balance proofs
	errg.Go(func() error {
		for {
			bps, err := operator.enclave.BalanceProofs()
			if err != nil {
				return fmt.Errorf("retrieving balance proofs: %w", err)
			}
			if len(bps) > 0 {
				log.Infof("Operator.Serve: Retrieved %d balance proofs", len(bps))
			}
			operator.balanceProofs.AddAll(bps)
		}
	})

	//TODO: operator handles on-chain challenge events

	// RPC handling
	remoteEnclave := newRemoteEnclave(operator)
	err = rpc.Register(remoteEnclave)
	if err != nil {
		return fmt.Errorf("registering remote enclave interface: %w", err)
	}
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("binding to socket: %w", err)
	}

	errg.Go(func() error { return http.Serve(l, nil) })

	return errg.Wait()
}

type depositProofs struct {
	mu      sync.RWMutex
	entries map[common.Address]*tee.DepositProof
}

func newDepositProofs() *depositProofs {
	return &depositProofs{entries: make(map[common.Address]*tee.DepositProof)}
}

// Get gets the deposit proof for the given user, threadsafe.
func (dps *depositProofs) Get(user common.Address) (*tee.DepositProof, bool) {
	dps.mu.RLock()
	defer dps.mu.RUnlock()

	dp, ok := dps.entries[user]

	return dp, ok
}

// AddAll adds the given deposit proofs, threadsafe.
func (dps *depositProofs) AddAll(in []*tee.DepositProof) {
	dps.mu.Lock()
	defer dps.mu.Unlock()

	for _, dp := range in {
		dps.entries[dp.Balance.Account] = dp
	}
}

type balanceProofs struct {
	mu      sync.RWMutex
	entries map[common.Address]*tee.BalanceProof
}

func newBalanceProofs() *balanceProofs {
	return &balanceProofs{entries: make(map[common.Address]*tee.BalanceProof)}
}

// Get gets the balance proof for the given user, threadsafe.
func (bps *balanceProofs) Get(user common.Address) (*tee.BalanceProof, bool) {
	bps.mu.RLock()
	defer bps.mu.RUnlock()

	bp, ok := bps.entries[user]

	return bp, ok
}

// AddAll adds the given balance proofs, threadsafe.
func (bps *balanceProofs) AddAll(in []*tee.BalanceProof) {
	bps.mu.Lock()
	defer bps.mu.Unlock()

	for _, bp := range in {
		bps.entries[bp.Balance.Account] = bp
	}
}

const gasLimit = 2000000
const defaultContextTimeout = 10 * time.Second

// AssertNoError logs the error and exits if the error is not nil.
func AssertNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// NewDefaultContext creates a default context for the operator.
func NewDefaultContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), defaultContextTimeout)
}
