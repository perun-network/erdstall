package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/perun-network/erdstall/eth"
)

// Config represents the operator's configuration.
type Config struct {
	EthereumNodeURL        string
	ContractAddr           string
	Mnemonic               string
	EnclaveDerivationPath  string
	OperatorDerivationPath string
	PhaseDuration          uint64
	ResponseDuration       uint64
	PowDepth               uint64
	RPCPort                uint16
	RPCHost                string
	KeyFile                string
	CertFile               string
	RespondChallenges      bool
	SendDepositProofs      bool
	SendBalanceProofs      bool
	NodeReqTimeout         uint64 // Node request timeout in seconds.
	WaitMinedTimeout       uint64 // Transaction mining timeout in seconds.
}

// dialTimeout will be used a timeout when dialing to the ethereum node.
const dialTimeout = 20 * time.Second

// LoadConfig loads an operator configuration from the given file path.
func LoadConfig(path string) (*Config, error) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var config Config
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	if config.NodeReqTimeout != 0 {
		eth.SetNodeReqTimeout(time.Duration(config.NodeReqTimeout) * time.Second)
	}
	if config.WaitMinedTimeout != 0 {
		eth.SetWaitMinedTimeout(time.Duration(config.WaitMinedTimeout) * time.Second)
	}

	return &config, nil
}

// envGanacheCmd is the environment variable to set the command
// for spawning a `ganache-cli`. It is set "ERDSTALL_GANACHE_CMD".
// See `ganacheCommand()`.
const envGanacheCmd = "ERDSTALL_GANACHE_CMD"

// GanacheCommand returns the command to spawn a ganache-cli.
// The default value is `ganache-cli` and can be set by an ENV
// variable. Example:
// ERDSTALL_GANACHE_CMD="./my_ganache.sh --seed 123" go test ./...
// The arguments for configuring the ganache will be appended
// to the value of the ENV variable.
func GanacheCommand() (cmd string, args []string) {
	cmd = "ganache-cli"
	if env := os.Getenv(envGanacheCmd); len(env) != 0 {
		splits := strings.Split(env, " ")
		cmd = splits[0]
		args = splits[1:]
	}
	return
}
