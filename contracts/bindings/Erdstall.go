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
const ErdstallABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_tee\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_phaseDuration\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_responseDuration\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Challenged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Exiting\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"}],\"name\":\"Frozen\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bigBang\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"challengeDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"challenges\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"}],\"name\":\"encodeBalanceProof\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ensureFrozen\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"exit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"exits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"frozenEpoch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"frozenWithdrawals\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"numChallenges\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"phaseDuration\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"responseDuration\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tee\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"verifyBalance\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawChallenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"epoch\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"internalType\":\"structErdstall.Balance\",\"name\":\"balance\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"withdrawFrozen\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ErdstallFuncSigs maps the 4-byte function signature to its string representation.
var ErdstallFuncSigs = map[string]string{
	"03cf0678": "bigBang()",
	"778a2707": "challenge((uint64,address,uint256),bytes)",
	"0d13fd7b": "challengeDeposit()",
	"234c49a0": "challenges(uint64,address)",
	"d0e30db0": "deposit()",
	"9b7c7725": "deposits(uint64,address)",
	"0b7042d2": "encodeBalanceProof((uint64,address,uint256))",
	"64c38ddd": "ensureFrozen()",
	"63a3a27f": "exit((uint64,address,uint256),bytes)",
	"70e4a2c4": "exits(uint64,address)",
	"585db72a": "frozenEpoch()",
	"3f48a2a8": "frozenWithdrawals(address)",
	"f2910773": "numChallenges(uint64)",
	"ac5553ce": "phaseDuration()",
	"854b86d9": "responseDuration()",
	"67eeb62b": "tee()",
	"a608911d": "verifyBalance((uint64,address,uint256),bytes)",
	"750f0acc": "withdraw(uint64)",
	"3de970e3": "withdrawChallenge()",
	"f4a85043": "withdrawFrozen((uint64,address,uint256),bytes)",
}

// ErdstallBin is the compiled bytecode used for deploying new contracts.
var ErdstallBin = "0x610100604052600580546001600160401b0319166002600160401b031790553480156200002b57600080fd5b5060405162001c0738038062001c078339810160408190526200004e91620000e9565b816001600160401b0316816002026001600160401b031611156200008f5760405162461bcd60e51b815260040162000086906200013e565b60405180910390fd5b60609290921b6001600160601b0319166080524360c090811b6001600160c01b031990811660a05291811b821681529190911b1660e05262000175565b80516001600160401b0381168114620000e457600080fd5b919050565b600080600060608486031215620000fe578283fd5b83516001600160a01b038116811462000115578384fd5b92506200012560208501620000cc565b91506200013560408501620000cc565b90509250925092565b60208082526019908201527f726573706f6e73654475726174696f6e20746f6f206c6f6e6700000000000000604082015260600190565b60805160601c60a05160c01c60c05160c01c60e05160c01c611a2d620001da60003980610af65280610fd0525080610b875280610ffa5280611058528061108c525080610382528061102452806110b65250806108695280610b415250611a2d6000f3fe60806040526004361061011f5760003560e01c806370e4a2c4116100a0578063a608911d11610064578063a608911d14610303578063ac5553ce14610323578063d0e30db014610338578063f291077314610340578063f4a85043146103605761011f565b806370e4a2c41461026e578063750f0acc1461028e578063778a2707146102ae578063854b86d9146102ce5780639b7c7725146102e35761011f565b80633f48a2a8116100e75780633f48a2a8146101d5578063585db72a1461020257806363a3a27f1461021757806364c38ddd1461023757806367eeb62b1461024c5761011f565b806303cf0678146101245780630b7042d21461014f5780630d13fd7b1461017c578063234c49a0146101935780633de970e3146101c0575b600080fd5b34801561013057600080fd5b50610139610380565b604051610146919061197c565b60405180910390f35b34801561015b57600080fd5b5061016f61016a36600461145c565b6103a4565b604051610146919061157e565b34801561018857600080fd5b506101916103e0565b005b34801561019f57600080fd5b506101b36101ae36600461152d565b61043f565b6040516101469190611973565b3480156101cc57600080fd5b5061019161045c565b3480156101e157600080fd5b506101f56101f03660046113b7565b6104b9565b6040516101469190611573565b34801561020e57600080fd5b506101396104ce565b34801561022357600080fd5b506101916102323660046113d8565b6104dd565b34801561024357600080fd5b506101916107d3565b34801561025857600080fd5b50610261610867565b604051610146919061155f565b34801561027a57600080fd5b506101b361028936600461152d565b61088b565b34801561029a57600080fd5b506101916102a9366004611513565b6108a8565b3480156102ba57600080fd5b506101916102c93660046113d8565b610a0b565b3480156102da57600080fd5b50610139610af4565b3480156102ef57600080fd5b506101b36102fe36600461152d565b610b18565b34801561030f57600080fd5b5061019161031e366004611477565b610b32565b34801561032f57600080fd5b50610139610b85565b610191610ba9565b34801561034c57600080fd5b506101b361035b366004611513565b610c65565b34801561036c57600080fd5b5061019161037b3660046113d8565b610c77565b7f000000000000000000000000000000000000000000000000000000000000000081565b6060308260000151836020015184604001516040516020016103c99493929190611782565b60405160208183030381529060405290505b919050565b6103e8610d47565b1561040e5760405162461bcd60e51b815260040161040590611812565b60405180910390fd5b610416610d61565b156104335760405162461bcd60e51b81526004016104059061194a565b61043d6000610d97565b565b600260209081526000928352604080842090915290825290205481565b6104646107d3565b6005546001600160401b03908116600101166000908152600260209081526040808320338452909152902054806104ad5760405162461bcd60e51b815260040161040590611839565b6104b681610eb1565b50565b60046020526000908152604090205460ff1681565b6005546001600160401b031681565b6104e5610d47565b156105025760405162461bcd60e51b815260040161040590611812565b61050a610d61565b156105275760405162461bcd60e51b81526004016104059061194a565b61052f610f64565b6001600160401b03166105456020850185611513565b6001600160401b03161461056b5760405162461bcd60e51b81526004016104059061191f565b6105b961057d3685900385018561145c565b83838080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610b3292505050565b600260006105ca6020860186611513565b6001600160401b03166001600160401b0316815260200190815260200160002060008460200160208101906105ff91906113b7565b6001600160a01b03168152602081019190915260400160002054610659573361062e60408501602086016113b7565b6001600160a01b0316146106545760405162461bcd60e51b815260040161040590611729565b6106f2565b600060028161066b6020870187611513565b6001600160401b03166001600160401b0316815260200190815260200160002060008560200160208101906106a091906113b7565b6001600160a01b03168152602080820192909252604001600090812092909255600391906106d090860186611513565b6001600160401b03168152602081019190915260400160002080546000190190555b6040830135600160006107086020870187611513565b6001600160401b03166001600160401b03168152602001908152602001600020600085602001602081019061073d91906113b7565b6001600160a01b03166001600160a01b031681526020019081526020016000208190555082602001602081019061077491906113b7565b6001600160a01b031661078a6020850185611513565b6001600160401b03167f874e6a4ac09c210cf4cd123caaf949f43c3c6f07f2f46f26ccc5b0fd881c3d0485604001356040516107c69190611973565b60405180910390a3505050565b6107db610d47565b156107e55761043d565b6107ed610d61565b6108095760405162461bcd60e51b8152600401610405906118b1565b60006001610815610f76565b6005805467ffffffffffffffff1916929091036001600160401b0381169283179091556040519092507f5e20151a99b0432a9ac06d33b91b77d3134ce0638cc70d7df042947ca48a2caf90600090a250565b7f000000000000000000000000000000000000000000000000000000000000000081565b600160209081526000928352604080842090915290825290205481565b6108b0610d47565b156108cd5760405162461bcd60e51b815260040161040590611812565b6108d5610d61565b156108f25760405162461bcd60e51b81526004016104059061194a565b6108fa610f64565b6001600160401b0316816001600160401b03161061092a5760405162461bcd60e51b815260040161040590611755565b6001600160401b03811660009081526001602090815260408083203384529091529020548061096b5760405162461bcd60e51b81526004016104059061165f565b6001600160401b038216600090815260016020908152604080832033808552925280832083905551909183156108fc02918491818181858888f193505050501580156109bb573d6000803e3d6000fd5b50336001600160a01b0316826001600160401b03167f0ff23c4cdc2733f56d8f04d7a351c4332a1cd3334287ed5b2e9c6a28da9d3533836040516109ff9190611973565b60405180910390a35050565b610a13610d47565b15610a305760405162461bcd60e51b815260040161040590611812565b610a38610d61565b15610a555760405162461bcd60e51b81526004016104059061194a565b33610a6660408501602086016113b7565b6001600160a01b031614610a8c5760405162461bcd60e51b815260040161040590611696565b610a94610f76565b6001600160401b0316610aaa6020850185611513565b6001600160401b031614610ad05760405162461bcd60e51b8152600401610405906116cd565b610ae261057d3685900385018561145c565b610aef8360400135610d97565b505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b600060208181529281526040808220909352908152205481565b610b65610b3e836103a4565b827f0000000000000000000000000000000000000000000000000000000000000000610f82565b610b815760405162461bcd60e51b8152600401610405906115d1565b5050565b7f000000000000000000000000000000000000000000000000000000000000000081565b610bb1610d47565b15610bce5760405162461bcd60e51b815260040161040590611812565b610bd6610d61565b15610bf35760405162461bcd60e51b81526004016104059061194a565b6000610bfd610fbd565b6001600160401b038116600081815260208181526040808320338085529252918290208054349081019091559151939450927fe007c38a05fbf2010d1c1ed20f91e675c91d41699926124738a8c3fe9fc791b491610c5a91611973565b60405180910390a350565b60036020526000908152604090205481565b610c7f6107d3565b33610c9060408501602086016113b7565b6001600160a01b031614610cb65760405162461bcd60e51b81526004016104059061187a565b6005546001600160401b0316610ccf6020850185611513565b6001600160401b031614610cf55760405162461bcd60e51b8152600401610405906117db565b610d0761057d3685900385018561145c565b6005546001600160401b039081166001011660009081526020818152604080832033845290915290819020549084013501610d4181610eb1565b50505050565b6005546001600160401b031667fffffffffffffffe141590565b60008060036000610d70610f76565b6001600160401b03166001600160401b031681526020019081526020016000205411905090565b610d9f610fcc565b15610dbc5760405162461bcd60e51b8152600401610405906115fc565b6000610dc6610f64565b6001600160401b038116600090815260026020908152604080832033845290915290205490915015610e0a5760405162461bcd60e51b8152600401610405906116fd565b6001600160401b038116600090815260208181526040808320338452909152902054820180610e4b5760405162461bcd60e51b815260040161040590611633565b6001600160401b038216600081815260026020908152604080832033808552908352818420869055848452600390925280832080546001019055519092917f9f71686e9e2eed0a0a99340b1c3b230369f255b1d452130cead54f8308654dfd91a3505050565b3360009081526004602052604090205460ff1615610ee15760405162461bcd60e51b8152600401610405906118e8565b33600081815260046020526040808220805460ff191660011790555183156108fc0291849190818181858888f19350505050158015610f24573d6000803e3d6000fd5b5060055460405133916001600160401b0316907f0ff23c4cdc2733f56d8f04d7a351c4332a1cd3334287ed5b2e9c6a28da9d353390610c5a908590611973565b60006002610f70611088565b03905090565b60006003610f70611088565b600080610f9585805190602001206110ed565b90506000610fa3828661113e565b6001600160a01b0390811690851614925050509392505050565b6000610fc7611088565b905090565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160401b03167f00000000000000000000000000000000000000000000000000000000000000006001600160401b03167f000000000000000000000000000000000000000000000000000000000000000043036001600160401b03168161105557fe5b067f0000000000000000000000000000000000000000000000000000000000000000036001600160401b03161115905090565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160401b03167f000000000000000000000000000000000000000000000000000000000000000043036001600160401b0316816110e757fe5b04905090565b604080517f19457468657265756d205369676e6564204d6573736167653a0a333200000000602080830191909152603c8083019490945282518083039094018452605c909101909152815191012090565b60008151604114611196576040805162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e67746800604482015290519081900360640190fd5b60208201516040830151606084015160001a7f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a08211156112075760405162461bcd60e51b81526004018080602001828103825260228152602001806119b46022913960400191505060405180910390fd5b8060ff16601b1415801561121f57508060ff16601c14155b1561125b5760405162461bcd60e51b81526004018080602001828103825260228152602001806119d66022913960400191505060405180910390fd5b600060018783868660405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156112b7573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b03811661131f576040805162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e61747572650000000000000000604482015290519081900360640190fd5b9695505050505050565b80356001600160a01b03811681146103db57600080fd5b600060608284031215611351578081fd5b604051606081018181106001600160401b038211171561136d57fe5b60405290508061137c836113a0565b815261138a60208401611329565b6020820152604083013560408201525092915050565b80356001600160401b03811681146103db57600080fd5b6000602082840312156113c8578081fd5b6113d182611329565b9392505050565b600080600083850360808112156113ed578283fd5b60608112156113fa578283fd5b5083925060608401356001600160401b0380821115611417578384fd5b818601915086601f83011261142a578384fd5b813581811115611438578485fd5b876020828501011115611449578485fd5b6020830194508093505050509250925092565b60006060828403121561146d578081fd5b6113d18383611340565b60008060808385031215611489578182fd5b6114938484611340565b915060608301356001600160401b03808211156114ae578283fd5b818501915085601f8301126114c1578283fd5b8135818111156114cd57fe5b6114e0601f8201601f1916602001611990565b91508082528660208285010111156114f6578384fd5b806020840160208401378101602001929092525090939092509050565b600060208284031215611524578081fd5b6113d1826113a0565b6000806040838503121561153f578182fd5b611548836113a0565b915061155660208401611329565b90509250929050565b6001600160a01b0391909116815260200190565b901515815260200190565b6000602080835283518082850152825b818110156115aa5785810183015185820160400152820161158e565b818111156115bb5783604083870101525b50601f01601f1916929092016040019392505050565b602080825260119082015270696e76616c6964207369676e617475726560781b604082015260600190565b6020808252601b908201527f696e206368616c6c656e676520726573706f6e73652070686173650000000000604082015260600190565b6020808252601290820152716e6f2076616c756520696e2073797374656d60701b604082015260600190565b60208082526018908201527f6e6f7468696e67206c65667420746f2077697468647261770000000000000000604082015260600190565b60208082526017908201527f6368616c6c656e67653a2077726f6e672073656e646572000000000000000000604082015260600190565b6020808252601690820152750c6d0c2d8d8cadcceca7440eee4dedcce40cae0dec6d60531b604082015260600190565b602080825260129082015271185b1c9958591e4818da185b1b195b99d95960721b604082015260600190565b60208082526012908201527132bc34ba1d103bb937b7339039b2b73232b960711b604082015260600190565b60208082526013908201527277697468647261773a20746f6f206561726c7960681b604082015260600190565b60a0808252600f908201526e4572647374616c6c42616c616e636560881b60c08201526001600160a01b0394851660208201526001600160401b0393909316604084015292166060820152608081019190915260e00190565b6020808252601b908201527f776974686472617746726f7a656e3a2077726f6e672065706f63680000000000604082015260600190565b6020808252600d908201526c383630b9b6b090333937bd32b760991b604082015260600190565b60208082526021908201527f6e6f7468696e67206c65667420746f207769746864726177202866726f7a656e6040820152602960f81b606082015260800190565b6020808252601c908201527f776974686472617746726f7a656e3a2077726f6e672073656e64657200000000604082015260600190565b6020808252601a908201527f6e6f206368616c6c656e676520696e206c6173742065706f6368000000000000604082015260600190565b6020808252601a908201527f616c72656164792077697468647261776e202866726f7a656e29000000000000604082015260600190565b6020808252601190820152700caf0d2e87440eee4dedcce40cae0dec6d607b1b604082015260600190565b6020808252600f908201526e706c61736d6120667265657a696e6760881b604082015260600190565b90815260200190565b6001600160401b0391909116815260200190565b6040518181016001600160401b03811182821017156119ab57fe5b60405291905056fe45434453413a20696e76616c6964207369676e6174757265202773272076616c756545434453413a20696e76616c6964207369676e6174757265202776272076616c7565a264697066735822122021141d4d04ceab273065b3c35cd6c62bb9fb93b529f2aaf72b215841f4cbb22664736f6c63430007040033"

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
// Solidity: function challenges(uint64 , address ) view returns(uint256)
func (_Erdstall *ErdstallCaller) Challenges(opts *bind.CallOpts, arg0 uint64, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "challenges", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Challenges is a free data retrieval call binding the contract method 0x234c49a0.
//
// Solidity: function challenges(uint64 , address ) view returns(uint256)
func (_Erdstall *ErdstallSession) Challenges(arg0 uint64, arg1 common.Address) (*big.Int, error) {
	return _Erdstall.Contract.Challenges(&_Erdstall.CallOpts, arg0, arg1)
}

// Challenges is a free data retrieval call binding the contract method 0x234c49a0.
//
// Solidity: function challenges(uint64 , address ) view returns(uint256)
func (_Erdstall *ErdstallCallerSession) Challenges(arg0 uint64, arg1 common.Address) (*big.Int, error) {
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

// FrozenWithdrawals is a free data retrieval call binding the contract method 0x3f48a2a8.
//
// Solidity: function frozenWithdrawals(address ) view returns(bool)
func (_Erdstall *ErdstallCaller) FrozenWithdrawals(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Erdstall.contract.Call(opts, &out, "frozenWithdrawals", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// FrozenWithdrawals is a free data retrieval call binding the contract method 0x3f48a2a8.
//
// Solidity: function frozenWithdrawals(address ) view returns(bool)
func (_Erdstall *ErdstallSession) FrozenWithdrawals(arg0 common.Address) (bool, error) {
	return _Erdstall.Contract.FrozenWithdrawals(&_Erdstall.CallOpts, arg0)
}

// FrozenWithdrawals is a free data retrieval call binding the contract method 0x3f48a2a8.
//
// Solidity: function frozenWithdrawals(address ) view returns(bool)
func (_Erdstall *ErdstallCallerSession) FrozenWithdrawals(arg0 common.Address) (bool, error) {
	return _Erdstall.Contract.FrozenWithdrawals(&_Erdstall.CallOpts, arg0)
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

// Challenge is a paid mutator transaction binding the contract method 0x778a2707.
//
// Solidity: function challenge((uint64,address,uint256) balance, bytes sig) returns()
func (_Erdstall *ErdstallTransactor) Challenge(opts *bind.TransactOpts, balance ErdstallBalance, sig []byte) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "challenge", balance, sig)
}

// Challenge is a paid mutator transaction binding the contract method 0x778a2707.
//
// Solidity: function challenge((uint64,address,uint256) balance, bytes sig) returns()
func (_Erdstall *ErdstallSession) Challenge(balance ErdstallBalance, sig []byte) (*types.Transaction, error) {
	return _Erdstall.Contract.Challenge(&_Erdstall.TransactOpts, balance, sig)
}

// Challenge is a paid mutator transaction binding the contract method 0x778a2707.
//
// Solidity: function challenge((uint64,address,uint256) balance, bytes sig) returns()
func (_Erdstall *ErdstallTransactorSession) Challenge(balance ErdstallBalance, sig []byte) (*types.Transaction, error) {
	return _Erdstall.Contract.Challenge(&_Erdstall.TransactOpts, balance, sig)
}

// ChallengeDeposit is a paid mutator transaction binding the contract method 0x0d13fd7b.
//
// Solidity: function challengeDeposit() returns()
func (_Erdstall *ErdstallTransactor) ChallengeDeposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "challengeDeposit")
}

// ChallengeDeposit is a paid mutator transaction binding the contract method 0x0d13fd7b.
//
// Solidity: function challengeDeposit() returns()
func (_Erdstall *ErdstallSession) ChallengeDeposit() (*types.Transaction, error) {
	return _Erdstall.Contract.ChallengeDeposit(&_Erdstall.TransactOpts)
}

// ChallengeDeposit is a paid mutator transaction binding the contract method 0x0d13fd7b.
//
// Solidity: function challengeDeposit() returns()
func (_Erdstall *ErdstallTransactorSession) ChallengeDeposit() (*types.Transaction, error) {
	return _Erdstall.Contract.ChallengeDeposit(&_Erdstall.TransactOpts)
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

// EnsureFrozen is a paid mutator transaction binding the contract method 0x64c38ddd.
//
// Solidity: function ensureFrozen() returns()
func (_Erdstall *ErdstallTransactor) EnsureFrozen(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "ensureFrozen")
}

// EnsureFrozen is a paid mutator transaction binding the contract method 0x64c38ddd.
//
// Solidity: function ensureFrozen() returns()
func (_Erdstall *ErdstallSession) EnsureFrozen() (*types.Transaction, error) {
	return _Erdstall.Contract.EnsureFrozen(&_Erdstall.TransactOpts)
}

// EnsureFrozen is a paid mutator transaction binding the contract method 0x64c38ddd.
//
// Solidity: function ensureFrozen() returns()
func (_Erdstall *ErdstallTransactorSession) EnsureFrozen() (*types.Transaction, error) {
	return _Erdstall.Contract.EnsureFrozen(&_Erdstall.TransactOpts)
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

// WithdrawChallenge is a paid mutator transaction binding the contract method 0x3de970e3.
//
// Solidity: function withdrawChallenge() returns()
func (_Erdstall *ErdstallTransactor) WithdrawChallenge(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erdstall.contract.Transact(opts, "withdrawChallenge")
}

// WithdrawChallenge is a paid mutator transaction binding the contract method 0x3de970e3.
//
// Solidity: function withdrawChallenge() returns()
func (_Erdstall *ErdstallSession) WithdrawChallenge() (*types.Transaction, error) {
	return _Erdstall.Contract.WithdrawChallenge(&_Erdstall.TransactOpts)
}

// WithdrawChallenge is a paid mutator transaction binding the contract method 0x3de970e3.
//
// Solidity: function withdrawChallenge() returns()
func (_Erdstall *ErdstallTransactorSession) WithdrawChallenge() (*types.Transaction, error) {
	return _Erdstall.Contract.WithdrawChallenge(&_Erdstall.TransactOpts)
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
