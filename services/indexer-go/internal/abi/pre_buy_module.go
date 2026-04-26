// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PreBuyModuleMetaData contains all meta data concerning the PreBuyModule contract.
var PreBuyModuleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"configs\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"vestingClone\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"vestAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cliffSeconds\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"durationSeconds\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"start\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"configured\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configure\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"vestAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cliffSeconds\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"durationSeconds\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"titanIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"curve\",\"type\":\"address\",\"internalType\":\"contractIBondingCurve\"}],\"outputs\":[{\"name\":\"clone\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"predictCloneAddress\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"release\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sweep\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"vestingImplementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Configured\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"vestingClone\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"vestAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"cliffSeconds\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"durationSeconds\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"QuoteSwept\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Released\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"clone\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AgentMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CliffExceedsDuration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedDeployment\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientBuyOutput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// PreBuyModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use PreBuyModuleMetaData.ABI instead.
var PreBuyModuleABI = PreBuyModuleMetaData.ABI

// PreBuyModule is an auto generated Go binding around an Ethereum contract.
type PreBuyModule struct {
	PreBuyModuleCaller     // Read-only binding to the contract
	PreBuyModuleTransactor // Write-only binding to the contract
	PreBuyModuleFilterer   // Log filterer for contract events
}

// PreBuyModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type PreBuyModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreBuyModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PreBuyModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreBuyModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PreBuyModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PreBuyModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PreBuyModuleSession struct {
	Contract     *PreBuyModule     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PreBuyModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PreBuyModuleCallerSession struct {
	Contract *PreBuyModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// PreBuyModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PreBuyModuleTransactorSession struct {
	Contract     *PreBuyModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// PreBuyModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type PreBuyModuleRaw struct {
	Contract *PreBuyModule // Generic contract binding to access the raw methods on
}

// PreBuyModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PreBuyModuleCallerRaw struct {
	Contract *PreBuyModuleCaller // Generic read-only contract binding to access the raw methods on
}

// PreBuyModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PreBuyModuleTransactorRaw struct {
	Contract *PreBuyModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPreBuyModule creates a new instance of PreBuyModule, bound to a specific deployed contract.
func NewPreBuyModule(address common.Address, backend bind.ContractBackend) (*PreBuyModule, error) {
	contract, err := bindPreBuyModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PreBuyModule{PreBuyModuleCaller: PreBuyModuleCaller{contract: contract}, PreBuyModuleTransactor: PreBuyModuleTransactor{contract: contract}, PreBuyModuleFilterer: PreBuyModuleFilterer{contract: contract}}, nil
}

// NewPreBuyModuleCaller creates a new read-only instance of PreBuyModule, bound to a specific deployed contract.
func NewPreBuyModuleCaller(address common.Address, caller bind.ContractCaller) (*PreBuyModuleCaller, error) {
	contract, err := bindPreBuyModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PreBuyModuleCaller{contract: contract}, nil
}

// NewPreBuyModuleTransactor creates a new write-only instance of PreBuyModule, bound to a specific deployed contract.
func NewPreBuyModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*PreBuyModuleTransactor, error) {
	contract, err := bindPreBuyModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PreBuyModuleTransactor{contract: contract}, nil
}

// NewPreBuyModuleFilterer creates a new log filterer instance of PreBuyModule, bound to a specific deployed contract.
func NewPreBuyModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*PreBuyModuleFilterer, error) {
	contract, err := bindPreBuyModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PreBuyModuleFilterer{contract: contract}, nil
}

// bindPreBuyModule binds a generic wrapper to an already deployed contract.
func bindPreBuyModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PreBuyModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PreBuyModule *PreBuyModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PreBuyModule.Contract.PreBuyModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PreBuyModule *PreBuyModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PreBuyModule.Contract.PreBuyModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PreBuyModule *PreBuyModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PreBuyModule.Contract.PreBuyModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PreBuyModule *PreBuyModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PreBuyModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PreBuyModule *PreBuyModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PreBuyModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PreBuyModule *PreBuyModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PreBuyModule.Contract.contract.Transact(opts, method, params...)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(address creator, address vestingClone, uint256 vestAmount, uint64 cliffSeconds, uint64 durationSeconds, uint64 start, bool configured)
func (_PreBuyModule *PreBuyModuleCaller) Configs(opts *bind.CallOpts, agent common.Address) (struct {
	Creator         common.Address
	VestingClone    common.Address
	VestAmount      *big.Int
	CliffSeconds    uint64
	DurationSeconds uint64
	Start           uint64
	Configured      bool
}, error) {
	var out []interface{}
	err := _PreBuyModule.contract.Call(opts, &out, "configs", agent)

	outstruct := new(struct {
		Creator         common.Address
		VestingClone    common.Address
		VestAmount      *big.Int
		CliffSeconds    uint64
		DurationSeconds uint64
		Start           uint64
		Configured      bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Creator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.VestingClone = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.VestAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.CliffSeconds = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.DurationSeconds = *abi.ConvertType(out[4], new(uint64)).(*uint64)
	outstruct.Start = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.Configured = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err

}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(address creator, address vestingClone, uint256 vestAmount, uint64 cliffSeconds, uint64 durationSeconds, uint64 start, bool configured)
func (_PreBuyModule *PreBuyModuleSession) Configs(agent common.Address) (struct {
	Creator         common.Address
	VestingClone    common.Address
	VestAmount      *big.Int
	CliffSeconds    uint64
	DurationSeconds uint64
	Start           uint64
	Configured      bool
}, error) {
	return _PreBuyModule.Contract.Configs(&_PreBuyModule.CallOpts, agent)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(address creator, address vestingClone, uint256 vestAmount, uint64 cliffSeconds, uint64 durationSeconds, uint64 start, bool configured)
func (_PreBuyModule *PreBuyModuleCallerSession) Configs(agent common.Address) (struct {
	Creator         common.Address
	VestingClone    common.Address
	VestAmount      *big.Int
	CliffSeconds    uint64
	DurationSeconds uint64
	Start           uint64
	Configured      bool
}, error) {
	return _PreBuyModule.Contract.Configs(&_PreBuyModule.CallOpts, agent)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PreBuyModule *PreBuyModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PreBuyModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PreBuyModule *PreBuyModuleSession) Owner() (common.Address, error) {
	return _PreBuyModule.Contract.Owner(&_PreBuyModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PreBuyModule *PreBuyModuleCallerSession) Owner() (common.Address, error) {
	return _PreBuyModule.Contract.Owner(&_PreBuyModule.CallOpts)
}

// PredictCloneAddress is a free data retrieval call binding the contract method 0xa37c6aab.
//
// Solidity: function predictCloneAddress(address agent) view returns(address)
func (_PreBuyModule *PreBuyModuleCaller) PredictCloneAddress(opts *bind.CallOpts, agent common.Address) (common.Address, error) {
	var out []interface{}
	err := _PreBuyModule.contract.Call(opts, &out, "predictCloneAddress", agent)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PredictCloneAddress is a free data retrieval call binding the contract method 0xa37c6aab.
//
// Solidity: function predictCloneAddress(address agent) view returns(address)
func (_PreBuyModule *PreBuyModuleSession) PredictCloneAddress(agent common.Address) (common.Address, error) {
	return _PreBuyModule.Contract.PredictCloneAddress(&_PreBuyModule.CallOpts, agent)
}

// PredictCloneAddress is a free data retrieval call binding the contract method 0xa37c6aab.
//
// Solidity: function predictCloneAddress(address agent) view returns(address)
func (_PreBuyModule *PreBuyModuleCallerSession) PredictCloneAddress(agent common.Address) (common.Address, error) {
	return _PreBuyModule.Contract.PredictCloneAddress(&_PreBuyModule.CallOpts, agent)
}

// VestingImplementation is a free data retrieval call binding the contract method 0x8882a2f6.
//
// Solidity: function vestingImplementation() view returns(address)
func (_PreBuyModule *PreBuyModuleCaller) VestingImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PreBuyModule.contract.Call(opts, &out, "vestingImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VestingImplementation is a free data retrieval call binding the contract method 0x8882a2f6.
//
// Solidity: function vestingImplementation() view returns(address)
func (_PreBuyModule *PreBuyModuleSession) VestingImplementation() (common.Address, error) {
	return _PreBuyModule.Contract.VestingImplementation(&_PreBuyModule.CallOpts)
}

// VestingImplementation is a free data retrieval call binding the contract method 0x8882a2f6.
//
// Solidity: function vestingImplementation() view returns(address)
func (_PreBuyModule *PreBuyModuleCallerSession) VestingImplementation() (common.Address, error) {
	return _PreBuyModule.Contract.VestingImplementation(&_PreBuyModule.CallOpts)
}

// Configure is a paid mutator transaction binding the contract method 0xc5bb705f.
//
// Solidity: function configure(address agent, address creator, uint256 vestAmount, uint64 cliffSeconds, uint64 durationSeconds, uint256 titanIn, address curve) returns(address clone)
func (_PreBuyModule *PreBuyModuleTransactor) Configure(opts *bind.TransactOpts, agent common.Address, creator common.Address, vestAmount *big.Int, cliffSeconds uint64, durationSeconds uint64, titanIn *big.Int, curve common.Address) (*types.Transaction, error) {
	return _PreBuyModule.contract.Transact(opts, "configure", agent, creator, vestAmount, cliffSeconds, durationSeconds, titanIn, curve)
}

// Configure is a paid mutator transaction binding the contract method 0xc5bb705f.
//
// Solidity: function configure(address agent, address creator, uint256 vestAmount, uint64 cliffSeconds, uint64 durationSeconds, uint256 titanIn, address curve) returns(address clone)
func (_PreBuyModule *PreBuyModuleSession) Configure(agent common.Address, creator common.Address, vestAmount *big.Int, cliffSeconds uint64, durationSeconds uint64, titanIn *big.Int, curve common.Address) (*types.Transaction, error) {
	return _PreBuyModule.Contract.Configure(&_PreBuyModule.TransactOpts, agent, creator, vestAmount, cliffSeconds, durationSeconds, titanIn, curve)
}

// Configure is a paid mutator transaction binding the contract method 0xc5bb705f.
//
// Solidity: function configure(address agent, address creator, uint256 vestAmount, uint64 cliffSeconds, uint64 durationSeconds, uint256 titanIn, address curve) returns(address clone)
func (_PreBuyModule *PreBuyModuleTransactorSession) Configure(agent common.Address, creator common.Address, vestAmount *big.Int, cliffSeconds uint64, durationSeconds uint64, titanIn *big.Int, curve common.Address) (*types.Transaction, error) {
	return _PreBuyModule.Contract.Configure(&_PreBuyModule.TransactOpts, agent, creator, vestAmount, cliffSeconds, durationSeconds, titanIn, curve)
}

// Release is a paid mutator transaction binding the contract method 0x19165587.
//
// Solidity: function release(address agent) returns(uint256 amount)
func (_PreBuyModule *PreBuyModuleTransactor) Release(opts *bind.TransactOpts, agent common.Address) (*types.Transaction, error) {
	return _PreBuyModule.contract.Transact(opts, "release", agent)
}

// Release is a paid mutator transaction binding the contract method 0x19165587.
//
// Solidity: function release(address agent) returns(uint256 amount)
func (_PreBuyModule *PreBuyModuleSession) Release(agent common.Address) (*types.Transaction, error) {
	return _PreBuyModule.Contract.Release(&_PreBuyModule.TransactOpts, agent)
}

// Release is a paid mutator transaction binding the contract method 0x19165587.
//
// Solidity: function release(address agent) returns(uint256 amount)
func (_PreBuyModule *PreBuyModuleTransactorSession) Release(agent common.Address) (*types.Transaction, error) {
	return _PreBuyModule.Contract.Release(&_PreBuyModule.TransactOpts, agent)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PreBuyModule *PreBuyModuleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PreBuyModule.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PreBuyModule *PreBuyModuleSession) RenounceOwnership() (*types.Transaction, error) {
	return _PreBuyModule.Contract.RenounceOwnership(&_PreBuyModule.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PreBuyModule *PreBuyModuleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PreBuyModule.Contract.RenounceOwnership(&_PreBuyModule.TransactOpts)
}

// Sweep is a paid mutator transaction binding the contract method 0x62c06767.
//
// Solidity: function sweep(address token, address to, uint256 amount) returns()
func (_PreBuyModule *PreBuyModuleTransactor) Sweep(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PreBuyModule.contract.Transact(opts, "sweep", token, to, amount)
}

// Sweep is a paid mutator transaction binding the contract method 0x62c06767.
//
// Solidity: function sweep(address token, address to, uint256 amount) returns()
func (_PreBuyModule *PreBuyModuleSession) Sweep(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PreBuyModule.Contract.Sweep(&_PreBuyModule.TransactOpts, token, to, amount)
}

// Sweep is a paid mutator transaction binding the contract method 0x62c06767.
//
// Solidity: function sweep(address token, address to, uint256 amount) returns()
func (_PreBuyModule *PreBuyModuleTransactorSession) Sweep(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PreBuyModule.Contract.Sweep(&_PreBuyModule.TransactOpts, token, to, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PreBuyModule *PreBuyModuleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PreBuyModule.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PreBuyModule *PreBuyModuleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PreBuyModule.Contract.TransferOwnership(&_PreBuyModule.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PreBuyModule *PreBuyModuleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PreBuyModule.Contract.TransferOwnership(&_PreBuyModule.TransactOpts, newOwner)
}

// PreBuyModuleConfiguredIterator is returned from FilterConfigured and is used to iterate over the raw logs and unpacked data for Configured events raised by the PreBuyModule contract.
type PreBuyModuleConfiguredIterator struct {
	Event *PreBuyModuleConfigured // Event containing the contract specifics and raw log

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
func (it *PreBuyModuleConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreBuyModuleConfigured)
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
		it.Event = new(PreBuyModuleConfigured)
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
func (it *PreBuyModuleConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreBuyModuleConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreBuyModuleConfigured represents a Configured event raised by the PreBuyModule contract.
type PreBuyModuleConfigured struct {
	Agent           common.Address
	Creator         common.Address
	VestingClone    common.Address
	VestAmount      *big.Int
	CliffSeconds    uint64
	DurationSeconds uint64
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConfigured is a free log retrieval operation binding the contract event 0x75cb32fcc4a1565883669cd138d6cb9a58f1ed87799912a136ba5af436c54fde.
//
// Solidity: event Configured(address indexed agent, address indexed creator, address indexed vestingClone, uint256 vestAmount, uint64 cliffSeconds, uint64 durationSeconds)
func (_PreBuyModule *PreBuyModuleFilterer) FilterConfigured(opts *bind.FilterOpts, agent []common.Address, creator []common.Address, vestingClone []common.Address) (*PreBuyModuleConfiguredIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var vestingCloneRule []interface{}
	for _, vestingCloneItem := range vestingClone {
		vestingCloneRule = append(vestingCloneRule, vestingCloneItem)
	}

	logs, sub, err := _PreBuyModule.contract.FilterLogs(opts, "Configured", agentRule, creatorRule, vestingCloneRule)
	if err != nil {
		return nil, err
	}
	return &PreBuyModuleConfiguredIterator{contract: _PreBuyModule.contract, event: "Configured", logs: logs, sub: sub}, nil
}

// WatchConfigured is a free log subscription operation binding the contract event 0x75cb32fcc4a1565883669cd138d6cb9a58f1ed87799912a136ba5af436c54fde.
//
// Solidity: event Configured(address indexed agent, address indexed creator, address indexed vestingClone, uint256 vestAmount, uint64 cliffSeconds, uint64 durationSeconds)
func (_PreBuyModule *PreBuyModuleFilterer) WatchConfigured(opts *bind.WatchOpts, sink chan<- *PreBuyModuleConfigured, agent []common.Address, creator []common.Address, vestingClone []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var vestingCloneRule []interface{}
	for _, vestingCloneItem := range vestingClone {
		vestingCloneRule = append(vestingCloneRule, vestingCloneItem)
	}

	logs, sub, err := _PreBuyModule.contract.WatchLogs(opts, "Configured", agentRule, creatorRule, vestingCloneRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreBuyModuleConfigured)
				if err := _PreBuyModule.contract.UnpackLog(event, "Configured", log); err != nil {
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

// ParseConfigured is a log parse operation binding the contract event 0x75cb32fcc4a1565883669cd138d6cb9a58f1ed87799912a136ba5af436c54fde.
//
// Solidity: event Configured(address indexed agent, address indexed creator, address indexed vestingClone, uint256 vestAmount, uint64 cliffSeconds, uint64 durationSeconds)
func (_PreBuyModule *PreBuyModuleFilterer) ParseConfigured(log types.Log) (*PreBuyModuleConfigured, error) {
	event := new(PreBuyModuleConfigured)
	if err := _PreBuyModule.contract.UnpackLog(event, "Configured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreBuyModuleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PreBuyModule contract.
type PreBuyModuleOwnershipTransferredIterator struct {
	Event *PreBuyModuleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PreBuyModuleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreBuyModuleOwnershipTransferred)
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
		it.Event = new(PreBuyModuleOwnershipTransferred)
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
func (it *PreBuyModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreBuyModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreBuyModuleOwnershipTransferred represents a OwnershipTransferred event raised by the PreBuyModule contract.
type PreBuyModuleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PreBuyModule *PreBuyModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PreBuyModuleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PreBuyModule.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PreBuyModuleOwnershipTransferredIterator{contract: _PreBuyModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PreBuyModule *PreBuyModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PreBuyModuleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PreBuyModule.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreBuyModuleOwnershipTransferred)
				if err := _PreBuyModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PreBuyModule *PreBuyModuleFilterer) ParseOwnershipTransferred(log types.Log) (*PreBuyModuleOwnershipTransferred, error) {
	event := new(PreBuyModuleOwnershipTransferred)
	if err := _PreBuyModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreBuyModuleQuoteSweptIterator is returned from FilterQuoteSwept and is used to iterate over the raw logs and unpacked data for QuoteSwept events raised by the PreBuyModule contract.
type PreBuyModuleQuoteSweptIterator struct {
	Event *PreBuyModuleQuoteSwept // Event containing the contract specifics and raw log

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
func (it *PreBuyModuleQuoteSweptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreBuyModuleQuoteSwept)
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
		it.Event = new(PreBuyModuleQuoteSwept)
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
func (it *PreBuyModuleQuoteSweptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreBuyModuleQuoteSweptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreBuyModuleQuoteSwept represents a QuoteSwept event raised by the PreBuyModule contract.
type PreBuyModuleQuoteSwept struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterQuoteSwept is a free log retrieval operation binding the contract event 0xec0cdf079961a285e45ef82b1256f7705eee61a2f12255dffd1277d333a968e6.
//
// Solidity: event QuoteSwept(address indexed token, address indexed to, uint256 amount)
func (_PreBuyModule *PreBuyModuleFilterer) FilterQuoteSwept(opts *bind.FilterOpts, token []common.Address, to []common.Address) (*PreBuyModuleQuoteSweptIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PreBuyModule.contract.FilterLogs(opts, "QuoteSwept", tokenRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PreBuyModuleQuoteSweptIterator{contract: _PreBuyModule.contract, event: "QuoteSwept", logs: logs, sub: sub}, nil
}

// WatchQuoteSwept is a free log subscription operation binding the contract event 0xec0cdf079961a285e45ef82b1256f7705eee61a2f12255dffd1277d333a968e6.
//
// Solidity: event QuoteSwept(address indexed token, address indexed to, uint256 amount)
func (_PreBuyModule *PreBuyModuleFilterer) WatchQuoteSwept(opts *bind.WatchOpts, sink chan<- *PreBuyModuleQuoteSwept, token []common.Address, to []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PreBuyModule.contract.WatchLogs(opts, "QuoteSwept", tokenRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreBuyModuleQuoteSwept)
				if err := _PreBuyModule.contract.UnpackLog(event, "QuoteSwept", log); err != nil {
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

// ParseQuoteSwept is a log parse operation binding the contract event 0xec0cdf079961a285e45ef82b1256f7705eee61a2f12255dffd1277d333a968e6.
//
// Solidity: event QuoteSwept(address indexed token, address indexed to, uint256 amount)
func (_PreBuyModule *PreBuyModuleFilterer) ParseQuoteSwept(log types.Log) (*PreBuyModuleQuoteSwept, error) {
	event := new(PreBuyModuleQuoteSwept)
	if err := _PreBuyModule.contract.UnpackLog(event, "QuoteSwept", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PreBuyModuleReleasedIterator is returned from FilterReleased and is used to iterate over the raw logs and unpacked data for Released events raised by the PreBuyModule contract.
type PreBuyModuleReleasedIterator struct {
	Event *PreBuyModuleReleased // Event containing the contract specifics and raw log

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
func (it *PreBuyModuleReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PreBuyModuleReleased)
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
		it.Event = new(PreBuyModuleReleased)
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
func (it *PreBuyModuleReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PreBuyModuleReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PreBuyModuleReleased represents a Released event raised by the PreBuyModule contract.
type PreBuyModuleReleased struct {
	Agent  common.Address
	Clone  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReleased is a free log retrieval operation binding the contract event 0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52.
//
// Solidity: event Released(address indexed agent, address indexed clone, uint256 amount)
func (_PreBuyModule *PreBuyModuleFilterer) FilterReleased(opts *bind.FilterOpts, agent []common.Address, clone []common.Address) (*PreBuyModuleReleasedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var cloneRule []interface{}
	for _, cloneItem := range clone {
		cloneRule = append(cloneRule, cloneItem)
	}

	logs, sub, err := _PreBuyModule.contract.FilterLogs(opts, "Released", agentRule, cloneRule)
	if err != nil {
		return nil, err
	}
	return &PreBuyModuleReleasedIterator{contract: _PreBuyModule.contract, event: "Released", logs: logs, sub: sub}, nil
}

// WatchReleased is a free log subscription operation binding the contract event 0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52.
//
// Solidity: event Released(address indexed agent, address indexed clone, uint256 amount)
func (_PreBuyModule *PreBuyModuleFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *PreBuyModuleReleased, agent []common.Address, clone []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var cloneRule []interface{}
	for _, cloneItem := range clone {
		cloneRule = append(cloneRule, cloneItem)
	}

	logs, sub, err := _PreBuyModule.contract.WatchLogs(opts, "Released", agentRule, cloneRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PreBuyModuleReleased)
				if err := _PreBuyModule.contract.UnpackLog(event, "Released", log); err != nil {
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

// ParseReleased is a log parse operation binding the contract event 0x2d87480f50083e2b2759522a8fdda59802650a8055e609a7772cf70c07748f52.
//
// Solidity: event Released(address indexed agent, address indexed clone, uint256 amount)
func (_PreBuyModule *PreBuyModuleFilterer) ParseReleased(log types.Log) (*PreBuyModuleReleased, error) {
	event := new(PreBuyModuleReleased)
	if err := _PreBuyModule.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
