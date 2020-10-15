package operator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config represents the operator's configuration.
type Config struct {
	EthereumNodeURL        string
	Mnemonic               string
	EnclaveDerivationPath  string
	OperatorDerivationPath string
	PhaseDuration          uint64
	ResponseDuration       uint64
	PowDepth               uint64
	Port                   int
}

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

	return &config, nil
}
