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
var ECDSABin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea264697066735822122037d288f8724b0aa95b900cf0f03a84678869ea20e35b35b7973e439e41f179ab64736f6c63430007030033"

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
const ErdstallABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tee\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_phaseDuration\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_responseDuration\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Challenged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Exiting\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"}],\"name\":\"Frozen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bigBang\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"challenges\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"}],\"name\":\"encodeBalanceProof\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"exit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"exits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"frozenEpoch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"frozenWithdraws\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"numChallenges\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"phaseDuration\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"recoverDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"responseDuration\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"verifyBalance\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"withdrawFrozen\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

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
var ErdstallBin = "0x610100604052600480546001600160401b0319166002600160401b031790553480156200002b57600080fd5b506040516200184b3803806200184b8339810160408190526200004e91620000e9565b816001600160401b0316816002026001600160401b031611156200008f5760405162461bcd60e51b815260040162000086906200013e565b60405180910390fd5b60609290921b6001600160601b0319166080524360c090811b6001600160c01b031990811660a05291811b821681529190911b1660e05262000175565b80516001600160401b0381168114620000e457600080fd5b919050565b600080600060608486031215620000fe578283fd5b83516001600160a01b038116811462000115578384fd5b92506200012560208501620000cc565b91506200013560408501620000cc565b90509250925092565b60208082526019908201527f726573706f6e73654475726174696f6e20746f6f206c6f6e6700000000000000604082015260600190565b60805160601c60a05160c01c60c05160c01c60e05160c01c611671620001da600039806108345280610cd05250806108c55280610cfa5280610d585280610e555250806103375280610d245280610e7f5250806106cd528061087f52506116716000f3fe6080604052600436106101095760003560e01c8063854b86d911610095578063c49abb2111610064578063c49abb21146102b8578063d0e30db0146102d8578063d2ef7398146102e0578063f2910773146102f5578063f4a850431461031557610109565b8063854b86d91461024e5780639b7c772514610263578063a608911d14610283578063ac5553ce146102a357610109565b8063585db72a116100dc578063585db72a146101aa57806363a3a27f146101bf57806367eeb62b146101df57806370e4a2c414610201578063750f0acc1461022e57610109565b806303cf06781461010e5780630b7042d214610139578063234c49a014610166578063423d703a14610193575b600080fd5b34801561011a57600080fd5b50610123610335565b6040516101309190611604565b60405180910390f35b34801561014557600080fd5b50610159610154366004611147565b610359565b60405161013091906112b8565b34801561017257600080fd5b50610186610181366004611218565b610395565b604051610130919061128f565b34801561019f57600080fd5b506101a86103b5565b005b3480156101b657600080fd5b506101236103c1565b3480156101cb57600080fd5b506101a86101da3660046110c3565b6103d0565b3480156101eb57600080fd5b506101f46106cb565b604051610130919061127b565b34801561020d57600080fd5b5061022161021c366004611218565b6106ef565b60405161013091906115fb565b34801561023a57600080fd5b506101a86102493660046111fe565b61070c565b34801561025a57600080fd5b50610123610832565b34801561026f57600080fd5b5061022161027e366004611218565b610856565b34801561028f57600080fd5b506101a861029e366004611162565b610870565b3480156102af57600080fd5b506101236108c3565b3480156102c457600080fd5b506101866102d3366004611218565b6108e7565b6101a8610907565b3480156102ec57600080fd5b506101a86109c3565b34801561030157600080fd5b506102216103103660046111fe565b610af3565b34801561032157600080fd5b506101a86103303660046110c3565b610b05565b7f000000000000000000000000000000000000000000000000000000000000000081565b60603082600001518360200151846040015160405160200161037e94939291906114a2565b60405160208183030381529060405290505b919050565b600260209081526000928352604080842090915290825290205460ff1681565b6103bf6000610b29565b565b6004546001600160401b031681565b6103d8610c20565b156103fe5760405162461bcd60e51b81526004016103f5906114fb565b60405180910390fd5b610406610c3a565b156104235760405162461bcd60e51b81526004016103f5906115d2565b61042b610c70565b6001600160401b031661044160208501856111fe565b6001600160401b03161461045457600080fd5b6104a261046636859003850185611147565b83838080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061087092505050565b600260006104b360208601866111fe565b6001600160401b03166001600160401b0316815260200190815260200160002060008460200160208101906104e891906110a2565b6001600160a01b0316815260208101919091526040016000205460ff16610545573361051a60408501602086016110a2565b6001600160a01b0316146105405760405162461bcd60e51b81526004016103f590611407565b6105ea565b600060028161055760208701876111fe565b6001600160401b03166001600160401b03168152602001908152602001600020600085602001602081019061058c91906110a2565b6001600160a01b031681526020808201929092526040016000908120805460ff191693151593909317909255600391906105c8908601866111fe565b6001600160401b03168152602081019190915260400160002080546000190190555b60408301356001600061060060208701876111fe565b6001600160401b03166001600160401b03168152602001908152602001600020600085602001602081019061063591906110a2565b6001600160a01b03166001600160a01b031681526020019081526020016000208190555082602001602081019061066c91906110a2565b6001600160a01b031661068260208501856111fe565b6001600160401b03167f874e6a4ac09c210cf4cd123caaf949f43c3c6f07f2f46f26ccc5b0fd881c3d0485604001356040516106be91906115fb565b60405180910390a3505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600160209081526000928352604080842090915290825290205481565b610714610c20565b156107315760405162461bcd60e51b81526004016103f5906114fb565b610739610c3a565b156107565760405162461bcd60e51b81526004016103f5906115d2565b61075e610c70565b6001600160401b0316816001600160401b03161061078e5760405162461bcd60e51b81526004016103f590611433565b6001600160401b0381166000908152600160209081526040808320338085529252808320805490849055905190926108fc841502918491818181858888f193505050501580156107e2573d6000803e3d6000fd5b50336001600160a01b0316826001600160401b03167f0ff23c4cdc2733f56d8f04d7a351c4332a1cd3334287ed5b2e9c6a28da9d35338360405161082691906115fb565b60405180910390a35050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600060208181529281526040808220909352908152205481565b6108a361087c83610359565b827f0000000000000000000000000000000000000000000000000000000000000000610c82565b6108bf5760405162461bcd60e51b81526004016103f590611342565b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600560209081526000928352604080842090915290825290205460ff1681565b61090f610c20565b1561092c5760405162461bcd60e51b81526004016103f5906114fb565b610934610c3a565b156109515760405162461bcd60e51b81526004016103f5906115d2565b600061095b610cbd565b6001600160401b038116600081815260208181526040808320338085529252918290208054349081019091559151939450927fe007c38a05fbf2010d1c1ed20f91e675c91d41699926124738a8c3fe9fc791b4916109b8916115fb565b60405180910390a350565b6109cb610c20565b156109e85760405162461bcd60e51b81526004016103f5906114fb565b6109f0610c3a565b15610a0d5760405162461bcd60e51b81526004016103f5906115d2565b610a15610ccc565b15610a325760405162461bcd60e51b81526004016103f5906113a4565b6000610a3c610c70565b6001600160401b038116600090815260026020908152604080832033845290915290205490915060ff1615610a835760405162461bcd60e51b81526004016103f5906113db565b6001600160401b038116600081815260026020908152604080832033808552908352818420805460ff1916600190811790915585855260039093528184208054909301909255519092917f9f71686e9e2eed0a0a99340b1c3b230369f255b1d452130cead54f8308654dfd91a350565b60036020526000908152604090205481565b610b1761046636859003850185611147565b610b248360400135610b29565b505050565b610b31610d88565b6004546001600160401b0316600090815260056020908152604080832033845290915290205460ff1615610b775760405162461bcd60e51b81526004016103f59061159b565b6000610b81610e1c565b6004546001600160401b03166000908152600560209081526040808320338085529252808320805460ff191660011790555192850193509183156108fc0291849190818181858888f19350505050158015610be0573d6000803e3d6000fd5b5060045460405133916001600160401b0316907f0ff23c4cdc2733f56d8f04d7a351c4332a1cd3334287ed5b2e9c6a28da9d3533906108269085906115fb565b6004546001600160401b031667fffffffffffffffe141590565b60008060036000610c49610e49565b6001600160401b03166001600160401b031681526020019081526020016000205411905090565b60006002610c7c610e51565b03905090565b600080610c958580519060200120610eb6565b90506000610ca38286610ee6565b6001600160a01b0390811690851614925050509392505050565b6000610cc7610e51565b905090565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160401b03167f00000000000000000000000000000000000000000000000000000000000000006001600160401b03167f000000000000000000000000000000000000000000000000000000000000000043036001600160401b031681610d5557fe5b067f0000000000000000000000000000000000000000000000000000000000000000036001600160401b03161115905090565b610d90610c20565b15610d9a576103bf565b610da2610c3a565b610dbe5760405162461bcd60e51b81526004016103f590611564565b60006001610dca610e49565b6004805467ffffffffffffffff1916929091036001600160401b0381169283179091556040519092507f5e20151a99b0432a9ac06d33b91b77d3134ce0638cc70d7df042947ca48a2caf90600090a250565b6004546001600160401b039081166001011660009081526020818152604080832033845290915290205490565b60006003610c7c5b60007f00000000000000000000000000000000000000000000000000000000000000006001600160401b03167f000000000000000000000000000000000000000000000000000000000000000043036001600160401b031681610eb057fe5b04905090565b600081604051602001610ec9919061124a565b604051602081830303815290604052805190602001209050919050565b60008151604114610f095760405162461bcd60e51b81526004016103f59061136d565b60208201516040830151606084015160001a7f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0821115610f5b5760405162461bcd60e51b81526004016103f590611460565b8060ff16601b14158015610f7357508060ff16601c14155b15610f905760405162461bcd60e51b81526004016103f590611522565b600060018783868660405160008152602001604052604051610fb5949392919061129a565b6020604051602081039080840390855afa158015610fd7573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b03811661100a5760405162461bcd60e51b81526004016103f59061130b565b9695505050505050565b80356001600160a01b038116811461039057600080fd5b60006060828403121561103c578081fd5b604051606081018181106001600160401b038211171561105857fe5b6040529050806110678361108b565b815261107560208401611014565b6020820152604083013560408201525092915050565b80356001600160401b038116811461039057600080fd5b6000602082840312156110b3578081fd5b6110bc82611014565b9392505050565b600080600083850360808112156110d8578283fd5b60608112156110e5578283fd5b5083925060608401356001600160401b0380821115611102578384fd5b818601915086601f830112611115578384fd5b813581811115611123578485fd5b876020828501011115611134578485fd5b6020830194508093505050509250925092565b600060608284031215611158578081fd5b6110bc838361102b565b60008060808385031215611174578182fd5b61117e848461102b565b915060608301356001600160401b0380821115611199578283fd5b818501915085601f8301126111ac578283fd5b8135818111156111b857fe5b6111cb601f8201601f1916602001611618565b91508082528660208285010111156111e1578384fd5b806020840160208401378101602001929092525090939092509050565b60006020828403121561120f578081fd5b6110bc8261108b565b6000806040838503121561122a578182fd5b6112338361108b565b915061124160208401611014565b90509250929050565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b6001600160a01b0391909116815260200190565b901515815260200190565b93845260ff9290921660208401526040830152606082015260800190565b6000602080835283518082850152825b818110156112e4578581018301518582016040015282016112c8565b818111156112f55783604083870101525b50601f01601f1916929092016040019392505050565b60208082526018908201527f45434453413a20696e76616c6964207369676e61747572650000000000000000604082015260600190565b602080825260119082015270696e76616c6964207369676e617475726560781b604082015260600190565b6020808252601f908201527f45434453413a20696e76616c6964207369676e6174757265206c656e67746800604082015260600190565b6020808252601b908201527f696e206368616c6c656e676520726573706f6e73652070686173650000000000604082015260600190565b602080825260129082015271185b1c9958591e4818da185b1b195b99d95960721b604082015260600190565b60208082526012908201527132bc34ba1d103bb937b7339039b2b73232b960711b604082015260600190565b60208082526013908201527277697468647261773a20746f6f206561726c7960681b604082015260600190565b60208082526022908201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604082015261756560f01b606082015260800190565b60a0808252600f908201526e4572647374616c6c42616c616e636560881b60c08201526001600160a01b0394851660208201526001600160401b0393909316604084015292166060820152608081019190915260e00190565b6020808252600d908201526c383630b9b6b090333937bd32b760991b604082015260600190565b60208082526022908201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604082015261756560f01b606082015260800190565b6020808252601a908201527f6e6f206368616c6c656e676520696e206c6173742065706f6368000000000000604082015260600190565b6020808252601a908201527f616c72656164792077697468647261776e202866726f7a656e29000000000000604082015260600190565b6020808252600f908201526e706c61736d6120667265657a696e6760881b604082015260600190565b90815260200190565b6001600160401b0391909116815260200190565b6040518181016001600160401b038111828210171561163357fe5b60405291905056fea26469706673582212209f56768be17dc5dd276c6ade795b3f463204cb270c76cc02f2e36cc3fdc5026d64736f6c63430007030033"

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
var SigBin = "0x60566023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220039052960ba829ddb0721a0f7cc1d8b832117abc5f32bc13d347ed647eccb36564736f6c63430007030033"

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
