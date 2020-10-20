package operator

import (
	"context"
	"fmt"
	"log"
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
	"github.com/pkg/errors"
)

// Operator resprents a TEE Plasma operator.
type Operator struct {
	enclave   tee.Enclave
	ethClient *eth.Client
	*depositProofs
	*balanceProofs
	contract *bindings.Erdstall
}

// New instantiates an operator with the given parameters.
func New(
	enclave tee.Enclave,
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
	log.Println("Enclave account loaded")

	enclave := prototype.NewEnclaveWithAccount(wallet, enclaveAccount)
	enclavePublicKey, _, err := enclave.Init()
	AssertNoError(err)
	log.Println("Enclave created")

	operatorAccountDerivationPath := hdwallet.MustParseDerivationPath(cfg.OperatorDerivationPath)
	operatorAccount, err := wallet.Derive(operatorAccountDerivationPath, true)
	AssertNoError(err)

	client, err := eth.CreateEthereumClient(cfg.EthereumNodeURL, wallet, operatorAccount)
	AssertNoError(err)
	log.Println("Ethereum client initialized")

	enclaveParameters := tee.Parameters{
		TEE:              enclavePublicKey,
		PhaseDuration:    cfg.PhaseDuration,
		ResponseDuration: cfg.ResponseDuration,
		PowDepth:         cfg.PowDepth,
	}

	err = client.DeployContracts(&enclaveParameters)
	AssertNoError(err)
	log.Printf("Contract deployed at %x\n", enclaveParameters.Contract)

	err = enclave.SetParams(enclaveParameters)
	AssertNoError(err)
	go enclave.Start()
	log.Println("Enclave initialized")

	operator, err := New(enclave, client, enclaveParameters.Contract)
	AssertNoError(err)

	return operator, enclaveParameters
}

// Serve starts the operator's main routine.
func (operator *Operator) Serve(port int) error {
	// Ethereum block handling
	blockSub, err := operator.ethClient.SubscribeToBlocks()
	if err != nil {
		return errors.WithMessage(err, "creating block subscription")
	}
	go func() {
		defer blockSub.Unsubscribe()
		for {
			b := <-blockSub.Blocks()
			operator.enclave.ProcessBlocks(b)
			log.Println("processed new block")
		}
	}()

	// Handle deposit proofs
	go func() {
		for {
			dps, err := operator.enclave.DepositProofs()
			if err != nil {
				log.Fatal("failed to retrieve deposit proofs", err)
			}
			log.Printf("retrieved %d deposit proofs\n", len(dps))
			operator.depositProofs.AddAll(dps)
		}
	}()

	// Handle balance proofs
	go func() {
		for {
			bps, err := operator.enclave.BalanceProofs()
			if err != nil {
				log.Fatal("failed to retrieve balance proofs", err)
			}
			log.Printf("retrieved %d balance proofs\n", len(bps))
			operator.balanceProofs.AddAll(bps)
		}
	}()

	//TODO: operator handles on-chain challenge events

	// RPC handling
	remoteEnclave := newRemoteEnclave(operator)
	err = rpc.Register(remoteEnclave)
	if err != nil {
		return errors.WithMessage(err, "registering remote enclave interface")
	}
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return errors.WithMessage(err, "binding to socket")
	}

	err = http.Serve(l, nil)
	return err
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

	dps.entries = make(map[common.Address]*tee.DepositProof)

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

	bps.entries = make(map[common.Address]*tee.BalanceProof)

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
