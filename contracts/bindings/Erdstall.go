// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ErdstallBalance is an auto generated low-level Go binding around an user-defined struct.
type ErdstallBalance struct {
	Epoch   uint64
	Account common.Address
	Value   *big.Int
}

// ECDSAABI is the input ABI used to generate the binding from.
const ECDSAABI = "[]"

// ECDSABin is the compiled bytecode used for deploying new contracts.
var ECDSABin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212200b98ef93566f8dc78bee4ef49b2fe5aa148020ee6fd1682827427527903c9bda64736f6c63430007040033"

// DeployECDSA deploys a new Ethereum contract, binding an instance of ECDSA to it.
func DeployECDSA(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ECDSA, error) {
	parsed, err := abi.JSON(strings.NewReader(ECDSAABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ECDSABin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ECDSA{ECDSACaller: ECDSACaller{contract: contract}, ECDSATransactor: ECDSATransactor{contract: contract}, ECDSAFilterer: ECDSAFilterer{contract: contract}}, nil
}

// ECDSA is an auto generated Go binding around an Ethereum contract.
type ECDSA struct {
	ECDSACaller     // Read-only binding to the contract
	ECDSATransactor // Write-only binding to the contract
	ECDSAFilterer   // Log filterer for contract events
}

// ECDSACaller is an auto generated read-only Go binding around an Ethereum contract.
type ECDSACaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSATransactor is an auto generated write-only Go binding around an Ethereum contract.
type ECDSATransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSAFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ECDSAFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ECDSASession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ECDSASession struct {
	Contract     *ECDSA            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ECDSACallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ECDSACallerSession struct {
	Contract *ECDSACaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ECDSATransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ECDSATransactorSession struct {
	Contract     *ECDSATransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ECDSARaw is an auto generated low-level Go binding around an Ethereum contract.
type ECDSARaw struct {
	Contract *ECDSA // Generic contract binding to access the raw methods on
}

// ECDSACallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ECDSACallerRaw struct {
	Contract *ECDSACaller // Generic read-only contract binding to access the raw methods on
}

// ECDSATransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ECDSATransactorRaw struct {
	Contract *ECDSATransactor // Generic write-only contract binding to access the raw methods on
}

// NewECDSA creates a new instance of ECDSA, bound to a specific deployed contract.
func NewECDSA(address common.Address, backend bind.ContractBackend) (*ECDSA, error) {
	contract, err := bindECDSA(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ECDSA{ECDSACaller: ECDSACaller{contract: contract}, ECDSATransactor: ECDSATransactor{contract: contract}, ECDSAFilterer: ECDSAFilterer{contract: contract}}, nil
}

// NewECDSACaller creates a new read-only instance of ECDSA, bound to a specific deployed contract.
func NewECDSACaller(address common.Address, caller bind.ContractCaller) (*ECDSACaller, error) {
	contract, err := bindECDSA(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ECDSACaller{contract: contract}, nil
}

// NewECDSATransactor creates a new write-only instance of ECDSA, bound to a specific deployed contract.
func NewECDSATransactor(address common.Address, transactor bind.ContractTransactor) (*ECDSATransactor, error) {
	contract, err := bindECDSA(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ECDSATransactor{contract: contract}, nil
}

// NewECDSAFilterer creates a new log filterer instance of ECDSA, bound to a specific deployed contract.
func NewECDSAFilterer(address common.Address, filterer bind.ContractFilterer) (*ECDSAFilterer, error) {
	contract, err := bindECDSA(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ECDSAFilterer{contract: contract}, nil
}

// bindECDSA binds a generic wrapper to an already deployed contract.
func bindECDSA(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ECDSAABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECDSA *ECDSARaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ECDSA.Contract.ECDSACaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECDSA *ECDSARaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECDSA.Contract.ECDSATransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECDSA *ECDSARaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECDSA.Contract.ECDSATransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ECDSA *ECDSACallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ECDSA.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ECDSA *ECDSATransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ECDSA.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ECDSA *ECDSATransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ECDSA.Contract.contract.Transact(opts, method, params...)
}

// ErdstallABI is the input ABI used to generate the binding from.
const ErdstallABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tee\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_phaseDuration\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_responseDuration\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Challenged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Exiting\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"}],\"name\":\"Frozen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bigBang\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"challenges\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"}],\"name\":\"encodeBalanceProof\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"exit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"exits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"freeze\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"frozenEpoch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"frozenWithdraws\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"numChallenges\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"phaseDuration\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"responseDuration\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"verifyBalance\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"withdrawFrozen\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ErdstallFuncSigs maps the 4-byte function signature to its string representation.
var ErdstallFuncSigs = map[string]string{
	"03cf0678": "bigBang()",
	"d2ef7398": "challenge()",
	"234c49a0": "challenges(uint64,address)",
	"d0e30db0": "deposit()",
	"9b7c7725": "deposits(uint64,address)",
	"0b7042d2": "encodeBalanceProof((uint64,address,uint256))",
	"63a3a27f": "exit((uint64,address,uint256),bytes)",
	"70e4a2c4": "exits(uint64,address)",
	"62a5af3b": "freeze()",
	"585db72a": "frozenEpoch()",
	"c49abb21": "frozenWithdraws(uint64,address)",
	"f2910773": "numChallenges(uint64)",
	"ac5553ce": "phaseDuration()",
	"423d703a": "recoverDeposit()",
	"854b86d9": "responseDuration()",
	"67eeb62b": "tee()",
	"a608911d": "verifyBalance((uint64,address,uint256),bytes)",
	"750f0acc": "withdraw(uint64)",
	"f4a85043": "withdrawFrozen((uint64,address,uint256),bytes)",
}

// ErdstallBin is the compiled bytecode used for deploying new contracts.
var ErdstallBin = "0x610100604052600480546001600160401b0319166002600160401b031790553480156200002b57600080fd5b50604051620018f6380380620018f68339810160408190526200004e91620000e9565b816001600160401b0316816002026001600160401b031611156200008f5760405162461bcd60e51b815260040162000086906200013e565b60405180910390fd5b60609290921b6001600160601b0319166080524360c090811b6001600160c01b031990811660a05291811b821681529190911b1660e05262000175565b80516001600160401b0381168114620000e457600080fd5b919050565b600080600060608486031215620000fe578283fd5b83516001600160a01b038116811462000115578384fd5b92506200012560208501620000cc565b91506200013560408501620000cc565b90509250925092565b60208082526019908201527f726573706f6e73654475726174696f6e20746f6f206c6f6e6700000000000000604082015260600190565b60805160601c60a05160c01c60c05160c01c60e05160c01c61171c620001da6000398061091f5280610de35250806109b05280610e0d5280610e6b5280610ecc5250806103575280610e375280610ef65250806107b8528061096a525061171c6000f3fe6080604052600436106101145760003560e01c8063750f0acc116100a0578063c49abb2111610064578063c49abb21146102d8578063d0e30db0146102f8578063d2ef739814610300578063f291077314610315578063f4a850431461033557610114565b8063750f0acc1461024e578063854b86d91461026e5780639b7c772514610283578063a608911d146102a3578063ac5553ce146102c357610114565b8063585db72a116100e7578063585db72a146101b557806362a5af3b146101ca57806363a3a27f146101df57806367eeb62b146101ff57806370e4a2c41461022157610114565b806303cf0678146101195780630b7042d214610144578063234c49a014610171578063423d703a1461019e575b600080fd5b34801561012557600080fd5b5061012e610355565b60405161013b919061166b565b60405180910390f35b34801561015057600080fd5b5061016461015f36600461129c565b610379565b60405161013b91906113be565b34801561017d57600080fd5b5061019161018c36600461136d565b6103b5565b60405161013b91906113b3565b3480156101aa57600080fd5b506101b36103d5565b005b3480156101c157600080fd5b5061012e61040e565b3480156101d657600080fd5b506101b361041d565b3480156101eb57600080fd5b506101b36101fa366004611218565b6104c4565b34801561020b57600080fd5b506102146107b6565b60405161013b919061139f565b34801561022d57600080fd5b5061024161023c36600461136d565b6107da565b60405161013b9190611662565b34801561025a57600080fd5b506101b3610269366004611353565b6107f7565b34801561027a57600080fd5b5061012e61091d565b34801561028f57600080fd5b5061024161029e36600461136d565b610941565b3480156102af57600080fd5b506101b36102be3660046112b7565b61095b565b3480156102cf57600080fd5b5061012e6109ae565b3480156102e457600080fd5b506101916102f336600461136d565b6109d2565b6101b36109f2565b34801561030c57600080fd5b506101b3610aae565b34801561032157600080fd5b50610241610330366004611353565b610bde565b34801561034157600080fd5b506101b3610350366004611218565b610bf0565b7f000000000000000000000000000000000000000000000000000000000000000081565b60603082600001518360200151846040015160405160200161039e949392919061154b565b60405160208183030381529060405290505b919050565b600260209081526000928352604080842090915290825290205460ff1681565b6103dd610c38565b6104025760405162461bcd60e51b81526004016103f99061143c565b60405180910390fd5b61040c6000610c52565b565b6004546001600160401b031681565b610425610c38565b156104425760405162461bcd60e51b81526004016103f990611523565b61044a610d41565b6104665760405162461bcd60e51b81526004016103f9906115cb565b60006001610472610d77565b6004805467ffffffffffffffff1916929091036001600160401b0381169283179091556040519092507f5e20151a99b0432a9ac06d33b91b77d3134ce0638cc70d7df042947ca48a2caf90600090a250565b6104cc610c38565b156104e95760405162461bcd60e51b81526004016103f9906115a4565b6104f1610d41565b1561050e5760405162461bcd60e51b81526004016103f990611639565b610516610d89565b6001600160401b031661052c6020850185611353565b6001600160401b03161461053f57600080fd5b61058d6105513685900385018561129c565b83838080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061095b92505050565b6002600061059e6020860186611353565b6001600160401b03166001600160401b0316815260200190815260200160002060008460200160208101906105d391906111f7565b6001600160a01b0316815260208101919091526040016000205460ff16610630573361060560408501602086016111f7565b6001600160a01b03161461062b5760405162461bcd60e51b81526004016103f9906114ca565b6106d5565b60006002816106426020870187611353565b6001600160401b03166001600160401b03168152602001908152602001600020600085602001602081019061067791906111f7565b6001600160a01b031681526020808201929092526040016000908120805460ff191693151593909317909255600391906106b390860186611353565b6001600160401b03168152602081019190915260400160002080546000190190555b6040830135600160006106eb6020870187611353565b6001600160401b03166001600160401b03168152602001908152602001600020600085602001602081019061072091906111f7565b6001600160a01b03166001600160a01b031681526020019081526020016000208190555082602001602081019061075791906111f7565b6001600160a01b031661076d6020850185611353565b6001600160401b03167f874e6a4ac09c210cf4cd123caaf949f43c3c6f07f2f46f26ccc5b0fd881c3d0485604001356040516107a99190611662565b60405180910390a3505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600160209081526000928352604080842090915290825290205481565b6107ff610c38565b1561081c5760405162461bcd60e51b81526004016103f9906115a4565b610824610d41565b156108415760405162461bcd60e51b81526004016103f990611639565b610849610d89565b6001600160401b0316816001600160401b0316106108795760405162461bcd60e51b81526004016103f9906114f6565b6001600160401b0381166000908152600160209081526040808320338085529252808320805490849055905190926108fc841502918491818181858888f193505050501580156108cd573d6000803e3d6000fd5b50336001600160a01b0316826001600160401b03167f0ff23c4cdc2733f56d8f04d7a351c4332a1cd3334287ed5b2e9c6a28da9d3533836040516109119190611662565b60405180910390a35050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600060208181529281526040808220909352908152205481565b61098e61096783610379565b827f0000000000000000000000000000000000000000000000000000000000000000610d95565b6109aa5760405162461bcd60e51b81526004016103f990611411565b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600560209081526000928352604080842090915290825290205460ff1681565b6109fa610c38565b15610a175760405162461bcd60e51b81526004016103f9906115a4565b610a1f610d41565b15610a3c5760405162461bcd60e51b81526004016103f990611639565b6000610a46610dd0565b6001600160401b038116600081815260208181526040808320338085529252918290208054349081019091559151939450927fe007c38a05fbf2010d1c1ed20f91e675c91d41699926124738a8c3fe9fc791b491610aa391611662565b60405180910390a350565b610ab6610c38565b15610ad35760405162461bcd60e51b81526004016103f9906115a4565b610adb610d41565b15610af85760405162461bcd60e51b81526004016103f990611639565b610b00610ddf565b15610b1d5760405162461bcd60e51b81526004016103f990611467565b6000610b27610d89565b6001600160401b038116600090815260026020908152604080832033845290915290205490915060ff1615610b6e5760405162461bcd60e51b81526004016103f99061149e565b6001600160401b038116600081815260026020908152604080832033808552908352818420805460ff1916600190811790915585855260039093528184208054909301909255519092917f9f71686e9e2eed0a0a99340b1c3b230369f255b1d452130cead54f8308654dfd91a350565b60036020526000908152604090205481565b610bf8610c38565b610c145760405162461bcd60e51b81526004016103f99061143c565b610c266105513685900385018561129c565b610c338360400135610c52565b505050565b6004546001600160401b031667fffffffffffffffe141590565b6004546001600160401b0316600090815260056020908152604080832033845290915290205460ff1615610c985760405162461bcd60e51b81526004016103f990611602565b6000610ca2610e9b565b6004546001600160401b03166000908152600560209081526040808320338085529252808320805460ff191660011790555192850193509183156108fc0291849190818181858888f19350505050158015610d01573d6000803e3d6000fd5b5060045460405133916001600160401b0316907f0ff23c4cdc2733f56d8f04d7a351c4332a1cd3334287ed5b2e9c6a28da9d353390610911908590611662565b60008060036000610d50610d77565b6001600160401b03166001600160401b031681526020019081526020016000205411905090565b60006003610d83610ec8565b03905090565b60006002610d83610ec8565b600080610da88580519060200120610f2d565b90506000610db68286610f7e565b6001600160a01b0390811690851614925050509392505050565b6000610dda610ec8565b905090565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160401b03167f00000000000000000000000000000000000000000000000000000000000000006001600160401b03167f000000000000000000000000000000000000000000000000000000000000000043036001600160401b031681610e6857fe5b067f0000000000000000000000000000000000000000000000000000000000000000036001600160401b03161115905090565b6004546001600160401b039081166001011660009081526020818152604080832033845290915290205490565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160401b03167f000000000000000000000000000000000000000000000000000000000000000043036001600160401b031681610f2757fe5b04905090565b604080517f19457468657265756d205369676e6564204d6573736167653a0a333200000000602080830191909152603c8083019490945282518083039094018452605c909101909152815191012090565b60008151604114610fd6576040805162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e67746800604482015290519081900360640190fd5b60208201516040830151606084015160001a7f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08211156110475760405162461bcd60e51b81526004018080602001828103825260228152602001806116a36022913960400191505060405180910390fd5b8060ff16601b1415801561105f57508060ff16601c14155b1561109b5760405162461bcd60e51b81526004018080602001828103825260228152602001806116c56022913960400191505060405180910390fd5b600060018783868660405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156110f7573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b03811661115f576040805162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e61747572650000000000000000604482015290519081900360640190fd5b9695505050505050565b80356001600160a01b03811681146103b057600080fd5b600060608284031215611191578081fd5b604051606081018181106001600160401b03821117156111ad57fe5b6040529050806111bc836111e0565b81526111ca60208401611169565b6020820152604083013560408201525092915050565b80356001600160401b03811681146103b057600080fd5b600060208284031215611208578081fd5b61121182611169565b9392505050565b6000806000838503608081121561122d578283fd5b606081121561123a578283fd5b5083925060608401356001600160401b0380821115611257578384fd5b818601915086601f83011261126a578384fd5b813581811115611278578485fd5b876020828501011115611289578485fd5b6020830194508093505050509250925092565b6000606082840312156112ad578081fd5b6112118383611180565b600080608083850312156112c9578182fd5b6112d38484611180565b915060608301356001600160401b03808211156112ee578283fd5b818501915085601f830112611301578283fd5b81358181111561130d57fe5b611320601f8201601f191660200161167f565b9150808252866020828501011115611336578384fd5b806020840160208401378101602001929092525090939092509050565b600060208284031215611364578081fd5b611211826111e0565b6000806040838503121561137f578182fd5b611388836111e0565b915061139660208401611169565b90509250929050565b6001600160a01b0391909116815260200190565b901515815260200190565b6000602080835283518082850152825b818110156113ea578581018301518582016040015282016113ce565b818111156113fb5783604083870101525b50601f01601f1916929092016040019392505050565b602080825260119082015270696e76616c6964207369676e617475726560781b604082015260600190565b602080825260119082015270383630b9b6b0903737ba10333937bd32b760791b604082015260600190565b6020808252601b908201527f696e206368616c6c656e676520726573706f6e73652070686173650000000000604082015260600190565b602080825260129082015271185b1c9958591e4818da185b1b195b99d95960721b604082015260600190565b60208082526012908201527132bc34ba1d103bb937b7339039b2b73232b960711b604082015260600190565b60208082526013908201527277697468647261773a20746f6f206561726c7960681b604082015260600190565b6020808252600e908201526d30b63932b0b23c90333937bd32b760911b604082015260600190565b60a0808252600f908201526e4572647374616c6c42616c616e636560881b60c08201526001600160a01b0394851660208201526001600160401b0393909316604084015292166060820152608081019190915260e00190565b6020808252600d908201526c383630b9b6b090333937bd32b760991b604082015260600190565b6020808252601a908201527f6e6f206368616c6c656e676520696e206c6173742065706f6368000000000000604082015260600190565b6020808252601a908201527f616c72656164792077697468647261776e202866726f7a656e29000000000000604082015260600190565b6020808252600f908201526e706c61736d6120667265657a696e6760881b604082015260600190565b90815260200190565b6001600160401b0391909116815260200190565b6040518181016001600160401b038111828210171561169a57fe5b60405291905056fe45434453413a20696e76616c6964207369676e6174757265202773272076616c756545434453413a20696e76616c6964207369676e6174757265202776272076616c7565a2646970667358221220f7697909125f8fc27e62135afaf233c9223ec32d099f50b12afc994fbbee5b6e64736f6c63430007040033"

// DeployErdstall deploys a new Ethereum contract, binding an instance of Erdstall to it.
func DeployErdstall(auth *bind.TransactOpts, backend bind.ContractBackend, _tee common.Address, _phaseDuration uint64, _responseDuration uint64) (common.Address, *types.Transaction, *Erdstall, error) {
	parsed, err := abi.JSON(strings.NewReader(ErdstallABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ErdstallBin), backend, _tee, _phaseDuration, _responseDuration)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Erdstall{ErdstallCaller: ErdstallCaller{contract: contract}, ErdstallTransactor: ErdstallTransactor{contract: contract}, ErdstallFilterer: ErdstallFilterer{contract: contract}}, nil
}

// Erdstall is an auto generated Go binding around an Ethereum contract.
type Erdstall struct {
	ErdstallCaller     // Read-only binding to the contract
	ErdstallTransactor // Write-only binding to the contract
	ErdstallFilterer   // Log filterer for contract events
}

// ErdstallCaller is an auto generated read-only Go binding around an Ethereum contract.
type ErdstallCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ErdstallTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ErdstallTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ErdstallFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ErdstallFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ErdstallSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ErdstallSession struct {
	Contract     *Erdstall         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ErdstallCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ErdstallCallerSession struct {
	Contract *ErdstallCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ErdstallTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ErdstallTransactorSession struct {
	Contract     *ErdstallTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ErdstallRaw is an auto generated low-level Go binding around an Ethereum contract.
type ErdstallRaw struct {
	Contract *Erdstall // Generic contract binding to access the raw methods on
}

// ErdstallCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ErdstallCallerRaw struct {
	Contract *ErdstallCaller // Generic read-only contract binding to access the raw methods on
}

// ErdstallTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ErdstallTransactorRaw struct {
	Contract *ErdstallTransactor // Generic write-only contract binding to access the raw methods on
}

// NewErdstall creates a new instance of Erdstall, bound to a specific deployed contract.
func NewErdstall(address common.Address, backend bind.ContractBackend) (*Erdstall, error) {
	contract, err := bindErdstall(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erdstall{ErdstallCaller: ErdstallCaller{contract: contract}, ErdstallTransactor: ErdstallTransactor{contract: contract}, ErdstallFilterer: ErdstallFilterer{contract: contract}}, nil
}

// NewErdstallCaller creates a new read-only instance of Erdstall, bound to a specific deployed contract.
func NewErdstallCaller(address common.Address, caller bind.ContractCaller) (*ErdstallCaller, error) {
	contract, err := bindErdstall(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ErdstallCaller{contract: contract}, nil
}

// NewErdstallTransactor creates a new write-only instance of Erdstall, bound to a specific deployed contract.
func NewErdstallTransactor(address common.Address, transactor bind.ContractTransactor) (*ErdstallTransactor, error) {
	contract, err := bindErdstall(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ErdstallTransactor{contract: contract}, nil
}

// NewErdstallFilterer creates a new log filterer instance of Erdstall, bound to a specific deployed contract.
func NewErdstallFilterer(address common.Address, filterer bind.ContractFilterer) (*ErdstallFilterer, error) {
	contract, err := bindErdstall(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ErdstallFilterer{contract: contract}, nil
}

// bindErdstall binds a generic wrapper to an already deployed contract.
func bindErdstall(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ErdstallABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erdstall *ErdstallRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erdstall.Contract.ErdstallCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erdstall *ErdstallRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erdstall.Contract.ErdstallTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erdstall *ErdstallRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erdstall.Contract.ErdstallTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erdstall *ErdstallCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erdstall.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erdstall *ErdstallTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erdstall.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erdstall *ErdstallTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erdstall.Contract.contract.Transact(opts, method, params...)
}

// BigBang is a free data retrieval call binding the contract method 0x03cf0678.
//
// Solidity: function bigBang() view returns(uint64)
func (_Erdstall *ErdstallCaller) BigBang(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "bigBang")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// BigBang is a free data retrieval call binding the contract method 0x03cf0678.
//
// Solidity: function bigBang() view returns(uint64)
func (_Erdstall *ErdstallSession) BigBang() (uint64, error) {
	return _Erdstall.Contract.BigBang(&_Erdstall.CallOpts)
}

// BigBang is a free data retrieval call binding the contract method 0x03cf0678.
//
// Solidity: function bigBang() view returns(uint64)
func (_Erdstall *ErdstallCallerSession) BigBang() (uint64, error) {
	return _Erdstall.Contract.BigBang(&_Erdstall.CallOpts)
}

// Challenges is a free data retrieval call binding the contract method 0x234c49a0.
//
// Solidity: function challenges(uint64 , address ) view returns(bool)
func (_Erdstall *ErdstallCaller) Challenges(opts *bind.CallOpts, arg0 uint64, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "challenges", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Challenges is a free data retrieval call binding the contract method 0x234c49a0.
//
// Solidity: function challenges(uint64 , address ) view returns(bool)
func (_Erdstall *ErdstallSession) Challenges(arg0 uint64, arg1 common.Address) (bool, error) {
	return _Erdstall.Contract.Challenges(&_Erdstall.CallOpts, arg0, arg1)
}

// Challenges is a free data retrieval call binding the contract method 0x234c49a0.
//
// Solidity: function challenges(uint64 , address ) view returns(bool)
func (_Erdstall *ErdstallCallerSession) Challenges(arg0 uint64, arg1 common.Address) (bool, error) {
	return _Erdstall.Contract.Challenges(&_Erdstall.CallOpts, arg0, arg1)
}

// Deposits is a free data retrieval call binding the contract method 0x9b7c7725.
//
// Solidity: function deposits(uint64 , address ) view returns(uint256)
func (_Erdstall *ErdstallCaller) Deposits(opts *bind.CallOpts, arg0 uint64, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "deposits", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Deposits is a free data retrieval call binding the contract method 0x9b7c7725.
//
// Solidity: function deposits(uint64 , address ) view returns(uint256)
func (_Erdstall *ErdstallSession) Deposits(arg0 uint64, arg1 common.Address) (*big.Int, error) {
	return _Erdstall.Contract.Deposits(&_Erdstall.CallOpts, arg0, arg1)
}

// Deposits is a free data retrieval call binding the contract method 0x9b7c7725.
//
// Solidity: function deposits(uint64 , address ) view returns(uint256)
func (_Erdstall *ErdstallCallerSession) Deposits(arg0 uint64, arg1 common.Address) (*big.Int, error) {
	return _Erdstall.Contract.Deposits(&_Erdstall.CallOpts, arg0, arg1)
}

// EncodeBalanceProof is a free data retrieval call binding the contract method 0x0b7042d2.
//
// Solidity: function encodeBalanceProof((uint64,address,uint256) balance) view returns(bytes)
func (_Erdstall *ErdstallCaller) EncodeBalanceProof(opts *bind.CallOpts, balance ErdstallBalance) ([]byte, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "encodeBalanceProof", balance)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// EncodeBalanceProof is a free data retrieval call binding the contract method 0x0b7042d2.
//
// Solidity: function encodeBalanceProof((uint64,address,uint256) balance) view returns(bytes)
func (_Erdstall *ErdstallSession) EncodeBalanceProof(balance ErdstallBalance) ([]byte, error) {
	return _Erdstall.Contract.EncodeBalanceProof(&_Erdstall.CallOpts, balance)
}

// EncodeBalanceProof is a free data retrieval call binding the contract method 0x0b7042d2.
//
// Solidity: function encodeBalanceProof((uint64,address,uint256) balance) view returns(bytes)
func (_Erdstall *ErdstallCallerSession) EncodeBalanceProof(balance ErdstallBalance) ([]byte, error) {
	return _Erdstall.Contract.EncodeBalanceProof(&_Erdstall.CallOpts, balance)
}

// Exits is a free data retrieval call binding the contract method 0x70e4a2c4.
//
// Solidity: function exits(uint64 , address ) view returns(uint256)
func (_Erdstall *ErdstallCaller) Exits(opts *bind.CallOpts, arg0 uint64, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "exits", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Exits is a free data retrieval call binding the contract method 0x70e4a2c4.
//
// Solidity: function exits(uint64 , address ) view returns(uint256)
func (_Erdstall *ErdstallSession) Exits(arg0 uint64, arg1 common.Address) (*big.Int, error) {
	return _Erdstall.Contract.Exits(&_Erdstall.CallOpts, arg0, arg1)
}

// Exits is a free data retrieval call binding the contract method 0x70e4a2c4.
//
// Solidity: function exits(uint64 , address ) view returns(uint256)
func (_Erdstall *ErdstallCallerSession) Exits(arg0 uint64, arg1 common.Address) (*big.Int, error) {
	return _Erdstall.Contract.Exits(&_Erdstall.CallOpts, arg0, arg1)
}

// FrozenEpoch is a free data retrieval call binding the contract method 0x585db72a.
//
// Solidity: function frozenEpoch() view returns(uint64)
func (_Erdstall *ErdstallCaller) FrozenEpoch(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "frozenEpoch")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// FrozenEpoch is a free data retrieval call binding the contract method 0x585db72a.
//
// Solidity: function frozenEpoch() view returns(uint64)
func (_Erdstall *ErdstallSession) FrozenEpoch() (uint64, error) {
	return _Erdstall.Contract.FrozenEpoch(&_Erdstall.CallOpts)
}

// FrozenEpoch is a free data retrieval call binding the contract method 0x585db72a.
//
// Solidity: function frozenEpoch() view returns(uint64)
func (_Erdstall *ErdstallCallerSession) FrozenEpoch() (uint64, error) {
	return _Erdstall.Contract.FrozenEpoch(&_Erdstall.CallOpts)
}

// FrozenWithdraws is a free data retrieval call binding the contract method 0xc49abb21.
//
// Solidity: function frozenWithdraws(uint64 , address ) view returns(bool)
func (_Erdstall *ErdstallCaller) FrozenWithdraws(opts *bind.CallOpts, arg0 uint64, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "frozenWithdraws", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// FrozenWithdraws is a free data retrieval call binding the contract method 0xc49abb21.
//
// Solidity: function frozenWithdraws(uint64 , address ) view returns(bool)
func (_Erdstall *ErdstallSession) FrozenWithdraws(arg0 uint64, arg1 common.Address) (bool, error) {
	return _Erdstall.Contract.FrozenWithdraws(&_Erdstall.CallOpts, arg0, arg1)
}

// FrozenWithdraws is a free data retrieval call binding the contract method 0xc49abb21.
//
// Solidity: function frozenWithdraws(uint64 , address ) view returns(bool)
func (_Erdstall *ErdstallCallerSession) FrozenWithdraws(arg0 uint64, arg1 common.Address) (bool, error) {
	return _Erdstall.Contract.FrozenWithdraws(&_Erdstall.CallOpts, arg0, arg1)
}

// NumChallenges is a free data retrieval call binding the contract method 0xf2910773.
//
// Solidity: function numChallenges(uint64 ) view returns(uint256)
func (_Erdstall *ErdstallCaller) NumChallenges(opts *bind.CallOpts, arg0 uint64) (*big.Int, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "numChallenges", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumChallenges is a free data retrieval call binding the contract method 0xf2910773.
//
// Solidity: function numChallenges(uint64 ) view returns(uint256)
func (_Erdstall *ErdstallSession) NumChallenges(arg0 uint64) (*big.Int, error) {
	return _Erdstall.Contract.NumChallenges(&_Erdstall.CallOpts, arg0)
}

// NumChallenges is a free data retrieval call binding the contract method 0xf2910773.
//
// Solidity: function numChallenges(uint64 ) view returns(uint256)
func (_Erdstall *ErdstallCallerSession) NumChallenges(arg0 uint64) (*big.Int, error) {
	return _Erdstall.Contract.NumChallenges(&_Erdstall.CallOpts, arg0)
}

// PhaseDuration is a free data retrieval call binding the contract method 0xac5553ce.
//
// Solidity: function phaseDuration() view returns(uint64)
func (_Erdstall *ErdstallCaller) PhaseDuration(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "phaseDuration")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PhaseDuration is a free data retrieval call binding the contract method 0xac5553ce.
//
// Solidity: function phaseDuration() view returns(uint64)
func (_Erdstall *ErdstallSession) PhaseDuration() (uint64, error) {
	return _Erdstall.Contract.PhaseDuration(&_Erdstall.CallOpts)
}

// PhaseDuration is a free data retrieval call binding the contract method 0xac5553ce.
//
// Solidity: function phaseDuration() view returns(uint64)
func (_Erdstall *ErdstallCallerSession) PhaseDuration() (uint64, error) {
	return _Erdstall.Contract.PhaseDuration(&_Erdstall.CallOpts)
}

// ResponseDuration is a free data retrieval call binding the contract method 0x854b86d9.
//
// Solidity: function responseDuration() view returns(uint64)
func (_Erdstall *ErdstallCaller) ResponseDuration(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "responseDuration")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ResponseDuration is a free data retrieval call binding the contract method 0x854b86d9.
//
// Solidity: function responseDuration() view returns(uint64)
func (_Erdstall *ErdstallSession) ResponseDuration() (uint64, error) {
	return _Erdstall.Contract.ResponseDuration(&_Erdstall.CallOpts)
}

// ResponseDuration is a free data retrieval call binding the contract method 0x854b86d9.
//
// Solidity: function responseDuration() view returns(uint64)
func (_Erdstall *ErdstallCallerSession) ResponseDuration() (uint64, error) {
	return _Erdstall.Contract.ResponseDuration(&_Erdstall.CallOpts)
}

// Tee is a free data retrieval call binding the contract method 0x67eeb62b.
//
// Solidity: function tee() view returns(address)
func (_Erdstall *ErdstallCaller) Tee(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "tee")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Tee is a free data retrieval call binding the contract method 0x67eeb62b.
//
// Solidity: function tee() view returns(address)
func (_Erdstall *ErdstallSession) Tee() (common.Address, error) {
	return _Erdstall.Contract.Tee(&_Erdstall.CallOpts)
}

// Tee is a free data retrieval call binding the contract method 0x67eeb62b.
//
// Solidity: function tee() view returns(address)
func (_Erdstall *ErdstallCallerSession) Tee() (common.Address, error) {
	return _Erdstall.Contract.Tee(&_Erdstall.CallOpts)
}

// VerifyBalance is a free data retrieval call binding the contract method 0xa608911d.
//
// Solidity: function verifyBalance((uint64,address,uint256) balance, bytes sig) view returns()
func (_Erdstall *ErdstallCaller) VerifyBalance(opts *bind.CallOpts, balance ErdstallBalance, sig []byte) error {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "verifyBalance", balance, sig)

	if err != nil {
		return err
	}

	return err

}

// VerifyBalance is a free data retrieval call binding the contract method 0xa608911d.
//
// Solidity: function verifyBalance((uint64,address,uint256) balance, bytes sig) view returns()
func (_Erdstall *ErdstallSession) VerifyBalance(balance ErdstallBalance, sig []byte) error {
	return _Erdstall.Contract.VerifyBalance(&_Erdstall.CallOpts, balance, sig)
}

// VerifyBalance is a free data retrieval call binding the contract method 0xa608911d.
//
// Solidity: function verifyBalance((uint64,address,uint256) balance, bytes sig) view returns()
func (_Erdstall *ErdstallCallerSession) VerifyBalance(balance ErdstallBalance, sig []byte) error {
	return _Erdstall.Contract.VerifyBalance(&_Erdstall.CallOpts, balance, sig)
}

// Challenge is a paid mutator transaction binding the contract method 0xd2ef7398.
//
// Solidity: function challenge() returns()
func (_Erdstall *ErdstallTransactor) Challenge(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "challenge")
}

// Challenge is a paid mutator transaction binding the contract method 0xd2ef7398.
//
// Solidity: function challenge() returns()
func (_Erdstall *ErdstallSession) Challenge() (*types.Transaction, error) {
	return _Erdstall.Contract.Challenge(&_Erdstall.TransactOpts)
}

// Challenge is a paid mutator transaction binding the contract method 0xd2ef7398.
//
// Solidity: function challenge() returns()
func (_Erdstall *ErdstallTransactorSession) Challenge() (*types.Transaction, error) {
	return _Erdstall.Contract.Challenge(&_Erdstall.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Erdstall *ErdstallTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Erdstall *ErdstallSession) Deposit() (*types.Transaction, error) {
	return _Erdstall.Contract.Deposit(&_Erdstall.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_Erdstall *ErdstallTransactorSession) Deposit() (*types.Transaction, error) {
	return _Erdstall.Contract.Deposit(&_Erdstall.TransactOpts)
}

// Exit is a paid mutator transaction binding the contract method 0x63a3a27f.
//
// Solidity: function exit((uint64,address,uint256) balance, bytes sig) returns()
func (_Erdstall *ErdstallTransactor) Exit(opts *bind.TransactOpts, balance ErdstallBalance, sig []byte) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "exit", balance, sig)
}

// Exit is a paid mutator transaction binding the contract method 0x63a3a27f.
//
// Solidity: function exit((uint64,address,uint256) balance, bytes sig) returns()
func (_Erdstall *ErdstallSession) Exit(balance ErdstallBalance, sig []byte) (*types.Transaction, error) {
	return _Erdstall.Contract.Exit(&_Erdstall.TransactOpts, balance, sig)
}

// Exit is a paid mutator transaction binding the contract method 0x63a3a27f.
//
// Solidity: function exit((uint64,address,uint256) balance, bytes sig) returns()
func (_Erdstall *ErdstallTransactorSession) Exit(balance ErdstallBalance, sig []byte) (*types.Transaction, error) {
	return _Erdstall.Contract.Exit(&_Erdstall.TransactOpts, balance, sig)
}

// Freeze is a paid mutator transaction binding the contract method 0x62a5af3b.
//
// Solidity: function freeze() returns()
func (_Erdstall *ErdstallTransactor) Freeze(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "freeze")
}

// Freeze is a paid mutator transaction binding the contract method 0x62a5af3b.
//
// Solidity: function freeze() returns()
func (_Erdstall *ErdstallSession) Freeze() (*types.Transaction, error) {
	return _Erdstall.Contract.Freeze(&_Erdstall.TransactOpts)
}

// Freeze is a paid mutator transaction binding the contract method 0x62a5af3b.
//
// Solidity: function freeze() returns()
func (_Erdstall *ErdstallTransactorSession) Freeze() (*types.Transaction, error) {
	return _Erdstall.Contract.Freeze(&_Erdstall.TransactOpts)
}

// RecoverDeposit is a paid mutator transaction binding the contract method 0x423d703a.
//
// Solidity: function recoverDeposit() returns()
func (_Erdstall *ErdstallTransactor) RecoverDeposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "recoverDeposit")
}

// RecoverDeposit is a paid mutator transaction binding the contract method 0x423d703a.
//
// Solidity: function recoverDeposit() returns()
func (_Erdstall *ErdstallSession) RecoverDeposit() (*types.Transaction, error) {
	return _Erdstall.Contract.RecoverDeposit(&_Erdstall.TransactOpts)
}

// RecoverDeposit is a paid mutator transaction binding the contract method 0x423d703a.
//
// Solidity: function recoverDeposit() returns()
func (_Erdstall *ErdstallTransactorSession) RecoverDeposit() (*types.Transaction, error) {
	return _Erdstall.Contract.RecoverDeposit(&_Erdstall.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x750f0acc.
//
// Solidity: function withdraw(uint64 epoch) returns()
func (_Erdstall *ErdstallTransactor) Withdraw(opts *bind.TransactOpts, epoch uint64) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "withdraw", epoch)
}

// Withdraw is a paid mutator transaction binding the contract method 0x750f0acc.
//
// Solidity: function withdraw(uint64 epoch) returns()
func (_Erdstall *ErdstallSession) Withdraw(epoch uint64) (*types.Transaction, error) {
	return _Erdstall.Contract.Withdraw(&_Erdstall.TransactOpts, epoch)
}

// Withdraw is a paid mutator transaction binding the contract method 0x750f0acc.
//
// Solidity: function withdraw(uint64 epoch) returns()
func (_Erdstall *ErdstallTransactorSession) Withdraw(epoch uint64) (*types.Transaction, error) {
	return _Erdstall.Contract.Withdraw(&_Erdstall.TransactOpts, epoch)
}

// WithdrawFrozen is a paid mutator transaction binding the contract method 0xf4a85043.
//
// Solidity: function withdrawFrozen((uint64,address,uint256) balance, bytes sig) returns()
func (_Erdstall *ErdstallTransactor) WithdrawFrozen(opts *bind.TransactOpts, balance ErdstallBalance, sig []byte) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "withdrawFrozen", balance, sig)
}

// WithdrawFrozen is a paid mutator transaction binding the contract method 0xf4a85043.
//
// Solidity: function withdrawFrozen((uint64,address,uint256) balance, bytes sig) returns()
func (_Erdstall *ErdstallSession) WithdrawFrozen(balance ErdstallBalance, sig []byte) (*types.Transaction, error) {
	return _Erdstall.Contract.WithdrawFrozen(&_Erdstall.TransactOpts, balance, sig)
}

// WithdrawFrozen is a paid mutator transaction binding the contract method 0xf4a85043.
//
// Solidity: function withdrawFrozen((uint64,address,uint256) balance, bytes sig) returns()
func (_Erdstall *ErdstallTransactorSession) WithdrawFrozen(balance ErdstallBalance, sig []byte) (*types.Transaction, error) {
	return _Erdstall.Contract.WithdrawFrozen(&_Erdstall.TransactOpts, balance, sig)
}

// ErdstallChallengedIterator is returned from FilterChallenged and is used to iterate over the raw logs and unpacked data for Challenged events raised by the Erdstall contract.
type ErdstallChallengedIterator struct {
	Event *ErdstallChallenged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ErdstallChallengedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ErdstallChallenged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ErdstallChallenged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ErdstallChallengedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ErdstallChallengedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ErdstallChallenged represents a Challenged event raised by the Erdstall contract.
type ErdstallChallenged struct {
	Epoch   uint64
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterChallenged is a free log retrieval operation binding the contract event 0x9f71686e9e2eed0a0a99340b1c3b230369f255b1d452130cead54f8308654dfd.
//
// Solidity: event Challenged(uint64 indexed epoch, address indexed account)
func (_Erdstall *ErdstallFilterer) FilterChallenged(opts *bind.FilterOpts, epoch []uint64, account []common.Address) (*ErdstallChallengedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Erdstall.contract.FilterLogs(opts, "Challenged", epochRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &ErdstallChallengedIterator{contract: _Erdstall.contract, event: "Challenged", logs: logs, sub: sub}, nil
}

// WatchChallenged is a free log subscription operation binding the contract event 0x9f71686e9e2eed0a0a99340b1c3b230369f255b1d452130cead54f8308654dfd.
//
// Solidity: event Challenged(uint64 indexed epoch, address indexed account)
func (_Erdstall *ErdstallFilterer) WatchChallenged(opts *bind.WatchOpts, sink chan<- *ErdstallChallenged, epoch []uint64, account []common.Address) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Erdstall.contract.WatchLogs(opts, "Challenged", epochRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ErdstallChallenged)
				if err := _Erdstall.contract.UnpackLog(event, "Challenged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChallenged is a log parse operation binding the contract event 0x9f71686e9e2eed0a0a99340b1c3b230369f255b1d452130cead54f8308654dfd.
//
// Solidity: event Challenged(uint64 indexed epoch, address indexed account)
func (_Erdstall *ErdstallFilterer) ParseChallenged(log types.Log) (*ErdstallChallenged, error) {
	event := new(ErdstallChallenged)
	if err := _Erdstall.contract.UnpackLog(event, "Challenged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ErdstallDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the Erdstall contract.
type ErdstallDepositedIterator struct {
	Event *ErdstallDeposited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ErdstallDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ErdstallDeposited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ErdstallDeposited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ErdstallDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ErdstallDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ErdstallDeposited represents a Deposited event raised by the Erdstall contract.
type ErdstallDeposited struct {
	Epoch   uint64
	Account common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0xe007c38a05fbf2010d1c1ed20f91e675c91d41699926124738a8c3fe9fc791b4.
//
// Solidity: event Deposited(uint64 indexed epoch, address indexed account, uint256 value)
func (_Erdstall *ErdstallFilterer) FilterDeposited(opts *bind.FilterOpts, epoch []uint64, account []common.Address) (*ErdstallDepositedIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Erdstall.contract.FilterLogs(opts, "Deposited", epochRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &ErdstallDepositedIterator{contract: _Erdstall.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0xe007c38a05fbf2010d1c1ed20f91e675c91d41699926124738a8c3fe9fc791b4.
//
// Solidity: event Deposited(uint64 indexed epoch, address indexed account, uint256 value)
func (_Erdstall *ErdstallFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *ErdstallDeposited, epoch []uint64, account []common.Address) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Erdstall.contract.WatchLogs(opts, "Deposited", epochRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ErdstallDeposited)
				if err := _Erdstall.contract.UnpackLog(event, "Deposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposited is a log parse operation binding the contract event 0xe007c38a05fbf2010d1c1ed20f91e675c91d41699926124738a8c3fe9fc791b4.
//
// Solidity: event Deposited(uint64 indexed epoch, address indexed account, uint256 value)
func (_Erdstall *ErdstallFilterer) ParseDeposited(log types.Log) (*ErdstallDeposited, error) {
	event := new(ErdstallDeposited)
	if err := _Erdstall.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ErdstallExitingIterator is returned from FilterExiting and is used to iterate over the raw logs and unpacked data for Exiting events raised by the Erdstall contract.
type ErdstallExitingIterator struct {
	Event *ErdstallExiting // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ErdstallExitingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ErdstallExiting)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ErdstallExiting)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ErdstallExitingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ErdstallExitingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ErdstallExiting represents a Exiting event raised by the Erdstall contract.
type ErdstallExiting struct {
	Epoch   uint64
	Account common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterExiting is a free log retrieval operation binding the contract event 0x874e6a4ac09c210cf4cd123caaf949f43c3c6f07f2f46f26ccc5b0fd881c3d04.
//
// Solidity: event Exiting(uint64 indexed epoch, address indexed account, uint256 value)
func (_Erdstall *ErdstallFilterer) FilterExiting(opts *bind.FilterOpts, epoch []uint64, account []common.Address) (*ErdstallExitingIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Erdstall.contract.FilterLogs(opts, "Exiting", epochRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &ErdstallExitingIterator{contract: _Erdstall.contract, event: "Exiting", logs: logs, sub: sub}, nil
}

// WatchExiting is a free log subscription operation binding the contract event 0x874e6a4ac09c210cf4cd123caaf949f43c3c6f07f2f46f26ccc5b0fd881c3d04.
//
// Solidity: event Exiting(uint64 indexed epoch, address indexed account, uint256 value)
func (_Erdstall *ErdstallFilterer) WatchExiting(opts *bind.WatchOpts, sink chan<- *ErdstallExiting, epoch []uint64, account []common.Address) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Erdstall.contract.WatchLogs(opts, "Exiting", epochRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ErdstallExiting)
				if err := _Erdstall.contract.UnpackLog(event, "Exiting", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseExiting is a log parse operation binding the contract event 0x874e6a4ac09c210cf4cd123caaf949f43c3c6f07f2f46f26ccc5b0fd881c3d04.
//
// Solidity: event Exiting(uint64 indexed epoch, address indexed account, uint256 value)
func (_Erdstall *ErdstallFilterer) ParseExiting(log types.Log) (*ErdstallExiting, error) {
	event := new(ErdstallExiting)
	if err := _Erdstall.contract.UnpackLog(event, "Exiting", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ErdstallFrozenIterator is returned from FilterFrozen and is used to iterate over the raw logs and unpacked data for Frozen events raised by the Erdstall contract.
type ErdstallFrozenIterator struct {
	Event *ErdstallFrozen // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ErdstallFrozenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ErdstallFrozen)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ErdstallFrozen)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ErdstallFrozenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ErdstallFrozenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ErdstallFrozen represents a Frozen event raised by the Erdstall contract.
type ErdstallFrozen struct {
	Epoch uint64
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterFrozen is a free log retrieval operation binding the contract event 0x5e20151a99b0432a9ac06d33b91b77d3134ce0638cc70d7df042947ca48a2caf.
//
// Solidity: event Frozen(uint64 indexed epoch)
func (_Erdstall *ErdstallFilterer) FilterFrozen(opts *bind.FilterOpts, epoch []uint64) (*ErdstallFrozenIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Erdstall.contract.FilterLogs(opts, "Frozen", epochRule)
	if err != nil {
		return nil, err
	}
	return &ErdstallFrozenIterator{contract: _Erdstall.contract, event: "Frozen", logs: logs, sub: sub}, nil
}

// WatchFrozen is a free log subscription operation binding the contract event 0x5e20151a99b0432a9ac06d33b91b77d3134ce0638cc70d7df042947ca48a2caf.
//
// Solidity: event Frozen(uint64 indexed epoch)
func (_Erdstall *ErdstallFilterer) WatchFrozen(opts *bind.WatchOpts, sink chan<- *ErdstallFrozen, epoch []uint64) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}

	logs, sub, err := _Erdstall.contract.WatchLogs(opts, "Frozen", epochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ErdstallFrozen)
				if err := _Erdstall.contract.UnpackLog(event, "Frozen", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFrozen is a log parse operation binding the contract event 0x5e20151a99b0432a9ac06d33b91b77d3134ce0638cc70d7df042947ca48a2caf.
//
// Solidity: event Frozen(uint64 indexed epoch)
func (_Erdstall *ErdstallFilterer) ParseFrozen(log types.Log) (*ErdstallFrozen, error) {
	event := new(ErdstallFrozen)
	if err := _Erdstall.contract.UnpackLog(event, "Frozen", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ErdstallWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the Erdstall contract.
type ErdstallWithdrawnIterator struct {
	Event *ErdstallWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ErdstallWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ErdstallWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ErdstallWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ErdstallWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ErdstallWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ErdstallWithdrawn represents a Withdrawn event raised by the Erdstall contract.
type ErdstallWithdrawn struct {
	Epoch   uint64
	Account common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x0ff23c4cdc2733f56d8f04d7a351c4332a1cd3334287ed5b2e9c6a28da9d3533.
//
// Solidity: event Withdrawn(uint64 indexed epoch, address indexed account, uint256 value)
func (_Erdstall *ErdstallFilterer) FilterWithdrawn(opts *bind.FilterOpts, epoch []uint64, account []common.Address) (*ErdstallWithdrawnIterator, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Erdstall.contract.FilterLogs(opts, "Withdrawn", epochRule, accountRule)
	if err != nil {
		return nil, err
	}
	return &ErdstallWithdrawnIterator{contract: _Erdstall.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x0ff23c4cdc2733f56d8f04d7a351c4332a1cd3334287ed5b2e9c6a28da9d3533.
//
// Solidity: event Withdrawn(uint64 indexed epoch, address indexed account, uint256 value)
func (_Erdstall *ErdstallFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *ErdstallWithdrawn, epoch []uint64, account []common.Address) (event.Subscription, error) {

	var epochRule []interface{}
	for _, epochItem := range epoch {
		epochRule = append(epochRule, epochItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _Erdstall.contract.WatchLogs(opts, "Withdrawn", epochRule, accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ErdstallWithdrawn)
				if err := _Erdstall.contract.UnpackLog(event, "Withdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawn is a log parse operation binding the contract event 0x0ff23c4cdc2733f56d8f04d7a351c4332a1cd3334287ed5b2e9c6a28da9d3533.
//
// Solidity: event Withdrawn(uint64 indexed epoch, address indexed account, uint256 value)
func (_Erdstall *ErdstallFilterer) ParseWithdrawn(log types.Log) (*ErdstallWithdrawn, error) {
	event := new(ErdstallWithdrawn)
	if err := _Erdstall.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SigABI is the input ABI used to generate the binding from.
const SigABI = "[]"

// SigBin is the compiled bytecode used for deploying new contracts.
var SigBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212205c0a198de4538b86f09c2b85f0f8cf5cba2c48312423b80a3c38fd793889aba864736f6c63430007040033"

// DeploySig deploys a new Ethereum contract, binding an instance of Sig to it.
func DeploySig(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Sig, error) {
	parsed, err := abi.JSON(strings.NewReader(SigABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SigBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Sig{SigCaller: SigCaller{contract: contract}, SigTransactor: SigTransactor{contract: contract}, SigFilterer: SigFilterer{contract: contract}}, nil
}

// Sig is an auto generated Go binding around an Ethereum contract.
type Sig struct {
	SigCaller     // Read-only binding to the contract
	SigTransactor // Write-only binding to the contract
	SigFilterer   // Log filterer for contract events
}

// SigCaller is an auto generated read-only Go binding around an Ethereum contract.
type SigCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SigTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SigTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SigFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SigFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SigSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SigSession struct {
	Contract     *Sig              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SigCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SigCallerSession struct {
	Contract *SigCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SigTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SigTransactorSession struct {
	Contract     *SigTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SigRaw is an auto generated low-level Go binding around an Ethereum contract.
type SigRaw struct {
	Contract *Sig // Generic contract binding to access the raw methods on
}

// SigCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SigCallerRaw struct {
	Contract *SigCaller // Generic read-only contract binding to access the raw methods on
}

// SigTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SigTransactorRaw struct {
	Contract *SigTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSig creates a new instance of Sig, bound to a specific deployed contract.
func NewSig(address common.Address, backend bind.ContractBackend) (*Sig, error) {
	contract, err := bindSig(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Sig{SigCaller: SigCaller{contract: contract}, SigTransactor: SigTransactor{contract: contract}, SigFilterer: SigFilterer{contract: contract}}, nil
}

// NewSigCaller creates a new read-only instance of Sig, bound to a specific deployed contract.
func NewSigCaller(address common.Address, caller bind.ContractCaller) (*SigCaller, error) {
	contract, err := bindSig(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SigCaller{contract: contract}, nil
}

// NewSigTransactor creates a new write-only instance of Sig, bound to a specific deployed contract.
func NewSigTransactor(address common.Address, transactor bind.ContractTransactor) (*SigTransactor, error) {
	contract, err := bindSig(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SigTransactor{contract: contract}, nil
}

// NewSigFilterer creates a new log filterer instance of Sig, bound to a specific deployed contract.
func NewSigFilterer(address common.Address, filterer bind.ContractFilterer) (*SigFilterer, error) {
	contract, err := bindSig(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SigFilterer{contract: contract}, nil
}

// bindSig binds a generic wrapper to an already deployed contract.
func bindSig(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SigABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sig *SigRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Sig.Contract.SigCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sig *SigRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sig.Contract.SigTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sig *SigRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sig.Contract.SigTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sig *SigCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Sig.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sig *SigTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sig.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sig *SigTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sig.Contract.contract.Transact(opts, method, params...)
}
