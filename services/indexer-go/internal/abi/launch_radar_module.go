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

// LaunchRadarModuleMetaData contains all meta data concerning the LaunchRadarModule contract.
var LaunchRadarModuleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"agentAdmin\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"canTrade\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"ok\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configs\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"whitelistOn\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"configured\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configure\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"whitelistOn\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"agentAdmin_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setStartTime\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newStart\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWhitelist\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"users\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"allowed\",\"type\":\"bool[]\",\"internalType\":\"bool[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWhitelistOn\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"on\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"whitelist\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Configured\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"whitelistOn\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"},{\"name\":\"agentAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StartTimeSet\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WhitelistToggled\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"on\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WhitelistUpdated\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"count\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyLive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAgentAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// LaunchRadarModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use LaunchRadarModuleMetaData.ABI instead.
var LaunchRadarModuleABI = LaunchRadarModuleMetaData.ABI

// LaunchRadarModule is an auto generated Go binding around an Ethereum contract.
type LaunchRadarModule struct {
	LaunchRadarModuleCaller     // Read-only binding to the contract
	LaunchRadarModuleTransactor // Write-only binding to the contract
	LaunchRadarModuleFilterer   // Log filterer for contract events
}

// LaunchRadarModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type LaunchRadarModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LaunchRadarModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LaunchRadarModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LaunchRadarModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LaunchRadarModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LaunchRadarModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LaunchRadarModuleSession struct {
	Contract     *LaunchRadarModule // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// LaunchRadarModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LaunchRadarModuleCallerSession struct {
	Contract *LaunchRadarModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// LaunchRadarModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LaunchRadarModuleTransactorSession struct {
	Contract     *LaunchRadarModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// LaunchRadarModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type LaunchRadarModuleRaw struct {
	Contract *LaunchRadarModule // Generic contract binding to access the raw methods on
}

// LaunchRadarModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LaunchRadarModuleCallerRaw struct {
	Contract *LaunchRadarModuleCaller // Generic read-only contract binding to access the raw methods on
}

// LaunchRadarModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LaunchRadarModuleTransactorRaw struct {
	Contract *LaunchRadarModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLaunchRadarModule creates a new instance of LaunchRadarModule, bound to a specific deployed contract.
func NewLaunchRadarModule(address common.Address, backend bind.ContractBackend) (*LaunchRadarModule, error) {
	contract, err := bindLaunchRadarModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LaunchRadarModule{LaunchRadarModuleCaller: LaunchRadarModuleCaller{contract: contract}, LaunchRadarModuleTransactor: LaunchRadarModuleTransactor{contract: contract}, LaunchRadarModuleFilterer: LaunchRadarModuleFilterer{contract: contract}}, nil
}

// NewLaunchRadarModuleCaller creates a new read-only instance of LaunchRadarModule, bound to a specific deployed contract.
func NewLaunchRadarModuleCaller(address common.Address, caller bind.ContractCaller) (*LaunchRadarModuleCaller, error) {
	contract, err := bindLaunchRadarModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LaunchRadarModuleCaller{contract: contract}, nil
}

// NewLaunchRadarModuleTransactor creates a new write-only instance of LaunchRadarModule, bound to a specific deployed contract.
func NewLaunchRadarModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*LaunchRadarModuleTransactor, error) {
	contract, err := bindLaunchRadarModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LaunchRadarModuleTransactor{contract: contract}, nil
}

// NewLaunchRadarModuleFilterer creates a new log filterer instance of LaunchRadarModule, bound to a specific deployed contract.
func NewLaunchRadarModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*LaunchRadarModuleFilterer, error) {
	contract, err := bindLaunchRadarModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LaunchRadarModuleFilterer{contract: contract}, nil
}

// bindLaunchRadarModule binds a generic wrapper to an already deployed contract.
func bindLaunchRadarModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LaunchRadarModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LaunchRadarModule *LaunchRadarModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LaunchRadarModule.Contract.LaunchRadarModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LaunchRadarModule *LaunchRadarModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.LaunchRadarModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LaunchRadarModule *LaunchRadarModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.LaunchRadarModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LaunchRadarModule *LaunchRadarModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LaunchRadarModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LaunchRadarModule *LaunchRadarModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LaunchRadarModule *LaunchRadarModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.contract.Transact(opts, method, params...)
}

// AgentAdmin is a free data retrieval call binding the contract method 0x4045accf.
//
// Solidity: function agentAdmin(address agent) view returns(address)
func (_LaunchRadarModule *LaunchRadarModuleCaller) AgentAdmin(opts *bind.CallOpts, agent common.Address) (common.Address, error) {
	var out []interface{}
	err := _LaunchRadarModule.contract.Call(opts, &out, "agentAdmin", agent)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentAdmin is a free data retrieval call binding the contract method 0x4045accf.
//
// Solidity: function agentAdmin(address agent) view returns(address)
func (_LaunchRadarModule *LaunchRadarModuleSession) AgentAdmin(agent common.Address) (common.Address, error) {
	return _LaunchRadarModule.Contract.AgentAdmin(&_LaunchRadarModule.CallOpts, agent)
}

// AgentAdmin is a free data retrieval call binding the contract method 0x4045accf.
//
// Solidity: function agentAdmin(address agent) view returns(address)
func (_LaunchRadarModule *LaunchRadarModuleCallerSession) AgentAdmin(agent common.Address) (common.Address, error) {
	return _LaunchRadarModule.Contract.AgentAdmin(&_LaunchRadarModule.CallOpts, agent)
}

// CanTrade is a free data retrieval call binding the contract method 0x657157e5.
//
// Solidity: function canTrade(address agent, address user) view returns(bool ok, string reason)
func (_LaunchRadarModule *LaunchRadarModuleCaller) CanTrade(opts *bind.CallOpts, agent common.Address, user common.Address) (struct {
	Ok     bool
	Reason string
}, error) {
	var out []interface{}
	err := _LaunchRadarModule.contract.Call(opts, &out, "canTrade", agent, user)

	outstruct := new(struct {
		Ok     bool
		Reason string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Ok = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Reason = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

// CanTrade is a free data retrieval call binding the contract method 0x657157e5.
//
// Solidity: function canTrade(address agent, address user) view returns(bool ok, string reason)
func (_LaunchRadarModule *LaunchRadarModuleSession) CanTrade(agent common.Address, user common.Address) (struct {
	Ok     bool
	Reason string
}, error) {
	return _LaunchRadarModule.Contract.CanTrade(&_LaunchRadarModule.CallOpts, agent, user)
}

// CanTrade is a free data retrieval call binding the contract method 0x657157e5.
//
// Solidity: function canTrade(address agent, address user) view returns(bool ok, string reason)
func (_LaunchRadarModule *LaunchRadarModuleCallerSession) CanTrade(agent common.Address, user common.Address) (struct {
	Ok     bool
	Reason string
}, error) {
	return _LaunchRadarModule.Contract.CanTrade(&_LaunchRadarModule.CallOpts, agent, user)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(uint64 startTime, bool whitelistOn, bool configured)
func (_LaunchRadarModule *LaunchRadarModuleCaller) Configs(opts *bind.CallOpts, agent common.Address) (struct {
	StartTime   uint64
	WhitelistOn bool
	Configured  bool
}, error) {
	var out []interface{}
	err := _LaunchRadarModule.contract.Call(opts, &out, "configs", agent)

	outstruct := new(struct {
		StartTime   uint64
		WhitelistOn bool
		Configured  bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StartTime = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.WhitelistOn = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.Configured = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(uint64 startTime, bool whitelistOn, bool configured)
func (_LaunchRadarModule *LaunchRadarModuleSession) Configs(agent common.Address) (struct {
	StartTime   uint64
	WhitelistOn bool
	Configured  bool
}, error) {
	return _LaunchRadarModule.Contract.Configs(&_LaunchRadarModule.CallOpts, agent)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(uint64 startTime, bool whitelistOn, bool configured)
func (_LaunchRadarModule *LaunchRadarModuleCallerSession) Configs(agent common.Address) (struct {
	StartTime   uint64
	WhitelistOn bool
	Configured  bool
}, error) {
	return _LaunchRadarModule.Contract.Configs(&_LaunchRadarModule.CallOpts, agent)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LaunchRadarModule *LaunchRadarModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LaunchRadarModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LaunchRadarModule *LaunchRadarModuleSession) Owner() (common.Address, error) {
	return _LaunchRadarModule.Contract.Owner(&_LaunchRadarModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LaunchRadarModule *LaunchRadarModuleCallerSession) Owner() (common.Address, error) {
	return _LaunchRadarModule.Contract.Owner(&_LaunchRadarModule.CallOpts)
}

// Whitelist is a free data retrieval call binding the contract method 0xb092145e.
//
// Solidity: function whitelist(address agent, address user) view returns(bool)
func (_LaunchRadarModule *LaunchRadarModuleCaller) Whitelist(opts *bind.CallOpts, agent common.Address, user common.Address) (bool, error) {
	var out []interface{}
	err := _LaunchRadarModule.contract.Call(opts, &out, "whitelist", agent, user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Whitelist is a free data retrieval call binding the contract method 0xb092145e.
//
// Solidity: function whitelist(address agent, address user) view returns(bool)
func (_LaunchRadarModule *LaunchRadarModuleSession) Whitelist(agent common.Address, user common.Address) (bool, error) {
	return _LaunchRadarModule.Contract.Whitelist(&_LaunchRadarModule.CallOpts, agent, user)
}

// Whitelist is a free data retrieval call binding the contract method 0xb092145e.
//
// Solidity: function whitelist(address agent, address user) view returns(bool)
func (_LaunchRadarModule *LaunchRadarModuleCallerSession) Whitelist(agent common.Address, user common.Address) (bool, error) {
	return _LaunchRadarModule.Contract.Whitelist(&_LaunchRadarModule.CallOpts, agent, user)
}

// Configure is a paid mutator transaction binding the contract method 0x5420faab.
//
// Solidity: function configure(address agent, uint64 startTime, bool whitelistOn, address agentAdmin_) returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactor) Configure(opts *bind.TransactOpts, agent common.Address, startTime uint64, whitelistOn bool, agentAdmin_ common.Address) (*types.Transaction, error) {
	return _LaunchRadarModule.contract.Transact(opts, "configure", agent, startTime, whitelistOn, agentAdmin_)
}

// Configure is a paid mutator transaction binding the contract method 0x5420faab.
//
// Solidity: function configure(address agent, uint64 startTime, bool whitelistOn, address agentAdmin_) returns()
func (_LaunchRadarModule *LaunchRadarModuleSession) Configure(agent common.Address, startTime uint64, whitelistOn bool, agentAdmin_ common.Address) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.Configure(&_LaunchRadarModule.TransactOpts, agent, startTime, whitelistOn, agentAdmin_)
}

// Configure is a paid mutator transaction binding the contract method 0x5420faab.
//
// Solidity: function configure(address agent, uint64 startTime, bool whitelistOn, address agentAdmin_) returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactorSession) Configure(agent common.Address, startTime uint64, whitelistOn bool, agentAdmin_ common.Address) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.Configure(&_LaunchRadarModule.TransactOpts, agent, startTime, whitelistOn, agentAdmin_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LaunchRadarModule.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LaunchRadarModule *LaunchRadarModuleSession) RenounceOwnership() (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.RenounceOwnership(&_LaunchRadarModule.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.RenounceOwnership(&_LaunchRadarModule.TransactOpts)
}

// SetStartTime is a paid mutator transaction binding the contract method 0x4d0a29be.
//
// Solidity: function setStartTime(address agent, uint64 newStart) returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactor) SetStartTime(opts *bind.TransactOpts, agent common.Address, newStart uint64) (*types.Transaction, error) {
	return _LaunchRadarModule.contract.Transact(opts, "setStartTime", agent, newStart)
}

// SetStartTime is a paid mutator transaction binding the contract method 0x4d0a29be.
//
// Solidity: function setStartTime(address agent, uint64 newStart) returns()
func (_LaunchRadarModule *LaunchRadarModuleSession) SetStartTime(agent common.Address, newStart uint64) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.SetStartTime(&_LaunchRadarModule.TransactOpts, agent, newStart)
}

// SetStartTime is a paid mutator transaction binding the contract method 0x4d0a29be.
//
// Solidity: function setStartTime(address agent, uint64 newStart) returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactorSession) SetStartTime(agent common.Address, newStart uint64) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.SetStartTime(&_LaunchRadarModule.TransactOpts, agent, newStart)
}

// SetWhitelist is a paid mutator transaction binding the contract method 0xb32986b7.
//
// Solidity: function setWhitelist(address agent, address[] users, bool[] allowed) returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactor) SetWhitelist(opts *bind.TransactOpts, agent common.Address, users []common.Address, allowed []bool) (*types.Transaction, error) {
	return _LaunchRadarModule.contract.Transact(opts, "setWhitelist", agent, users, allowed)
}

// SetWhitelist is a paid mutator transaction binding the contract method 0xb32986b7.
//
// Solidity: function setWhitelist(address agent, address[] users, bool[] allowed) returns()
func (_LaunchRadarModule *LaunchRadarModuleSession) SetWhitelist(agent common.Address, users []common.Address, allowed []bool) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.SetWhitelist(&_LaunchRadarModule.TransactOpts, agent, users, allowed)
}

// SetWhitelist is a paid mutator transaction binding the contract method 0xb32986b7.
//
// Solidity: function setWhitelist(address agent, address[] users, bool[] allowed) returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactorSession) SetWhitelist(agent common.Address, users []common.Address, allowed []bool) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.SetWhitelist(&_LaunchRadarModule.TransactOpts, agent, users, allowed)
}

// SetWhitelistOn is a paid mutator transaction binding the contract method 0xd24fea56.
//
// Solidity: function setWhitelistOn(address agent, bool on) returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactor) SetWhitelistOn(opts *bind.TransactOpts, agent common.Address, on bool) (*types.Transaction, error) {
	return _LaunchRadarModule.contract.Transact(opts, "setWhitelistOn", agent, on)
}

// SetWhitelistOn is a paid mutator transaction binding the contract method 0xd24fea56.
//
// Solidity: function setWhitelistOn(address agent, bool on) returns()
func (_LaunchRadarModule *LaunchRadarModuleSession) SetWhitelistOn(agent common.Address, on bool) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.SetWhitelistOn(&_LaunchRadarModule.TransactOpts, agent, on)
}

// SetWhitelistOn is a paid mutator transaction binding the contract method 0xd24fea56.
//
// Solidity: function setWhitelistOn(address agent, bool on) returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactorSession) SetWhitelistOn(agent common.Address, on bool) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.SetWhitelistOn(&_LaunchRadarModule.TransactOpts, agent, on)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _LaunchRadarModule.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LaunchRadarModule *LaunchRadarModuleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.TransferOwnership(&_LaunchRadarModule.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LaunchRadarModule *LaunchRadarModuleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _LaunchRadarModule.Contract.TransferOwnership(&_LaunchRadarModule.TransactOpts, newOwner)
}

// LaunchRadarModuleConfiguredIterator is returned from FilterConfigured and is used to iterate over the raw logs and unpacked data for Configured events raised by the LaunchRadarModule contract.
type LaunchRadarModuleConfiguredIterator struct {
	Event *LaunchRadarModuleConfigured // Event containing the contract specifics and raw log

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
func (it *LaunchRadarModuleConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LaunchRadarModuleConfigured)
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
		it.Event = new(LaunchRadarModuleConfigured)
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
func (it *LaunchRadarModuleConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LaunchRadarModuleConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LaunchRadarModuleConfigured represents a Configured event raised by the LaunchRadarModule contract.
type LaunchRadarModuleConfigured struct {
	Agent       common.Address
	StartTime   uint64
	WhitelistOn bool
	AgentAdmin  common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterConfigured is a free log retrieval operation binding the contract event 0x3e8ab6864a473c81f9d5268a2ef00c16babd5a401270598c6437fe3aa0edc797.
//
// Solidity: event Configured(address indexed agent, uint64 startTime, bool whitelistOn, address agentAdmin)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) FilterConfigured(opts *bind.FilterOpts, agent []common.Address) (*LaunchRadarModuleConfiguredIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _LaunchRadarModule.contract.FilterLogs(opts, "Configured", agentRule)
	if err != nil {
		return nil, err
	}
	return &LaunchRadarModuleConfiguredIterator{contract: _LaunchRadarModule.contract, event: "Configured", logs: logs, sub: sub}, nil
}

// WatchConfigured is a free log subscription operation binding the contract event 0x3e8ab6864a473c81f9d5268a2ef00c16babd5a401270598c6437fe3aa0edc797.
//
// Solidity: event Configured(address indexed agent, uint64 startTime, bool whitelistOn, address agentAdmin)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) WatchConfigured(opts *bind.WatchOpts, sink chan<- *LaunchRadarModuleConfigured, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _LaunchRadarModule.contract.WatchLogs(opts, "Configured", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LaunchRadarModuleConfigured)
				if err := _LaunchRadarModule.contract.UnpackLog(event, "Configured", log); err != nil {
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

// ParseConfigured is a log parse operation binding the contract event 0x3e8ab6864a473c81f9d5268a2ef00c16babd5a401270598c6437fe3aa0edc797.
//
// Solidity: event Configured(address indexed agent, uint64 startTime, bool whitelistOn, address agentAdmin)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) ParseConfigured(log types.Log) (*LaunchRadarModuleConfigured, error) {
	event := new(LaunchRadarModuleConfigured)
	if err := _LaunchRadarModule.contract.UnpackLog(event, "Configured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LaunchRadarModuleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the LaunchRadarModule contract.
type LaunchRadarModuleOwnershipTransferredIterator struct {
	Event *LaunchRadarModuleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *LaunchRadarModuleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LaunchRadarModuleOwnershipTransferred)
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
		it.Event = new(LaunchRadarModuleOwnershipTransferred)
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
func (it *LaunchRadarModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LaunchRadarModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LaunchRadarModuleOwnershipTransferred represents a OwnershipTransferred event raised by the LaunchRadarModule contract.
type LaunchRadarModuleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*LaunchRadarModuleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LaunchRadarModule.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &LaunchRadarModuleOwnershipTransferredIterator{contract: _LaunchRadarModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LaunchRadarModuleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LaunchRadarModule.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LaunchRadarModuleOwnershipTransferred)
				if err := _LaunchRadarModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_LaunchRadarModule *LaunchRadarModuleFilterer) ParseOwnershipTransferred(log types.Log) (*LaunchRadarModuleOwnershipTransferred, error) {
	event := new(LaunchRadarModuleOwnershipTransferred)
	if err := _LaunchRadarModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LaunchRadarModuleStartTimeSetIterator is returned from FilterStartTimeSet and is used to iterate over the raw logs and unpacked data for StartTimeSet events raised by the LaunchRadarModule contract.
type LaunchRadarModuleStartTimeSetIterator struct {
	Event *LaunchRadarModuleStartTimeSet // Event containing the contract specifics and raw log

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
func (it *LaunchRadarModuleStartTimeSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LaunchRadarModuleStartTimeSet)
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
		it.Event = new(LaunchRadarModuleStartTimeSet)
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
func (it *LaunchRadarModuleStartTimeSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LaunchRadarModuleStartTimeSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LaunchRadarModuleStartTimeSet represents a StartTimeSet event raised by the LaunchRadarModule contract.
type LaunchRadarModuleStartTimeSet struct {
	Agent     common.Address
	StartTime uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStartTimeSet is a free log retrieval operation binding the contract event 0xaa1bcd70b51a29ec01d817f60e40a1ce21afc80546629beda8dd88e01ea0856d.
//
// Solidity: event StartTimeSet(address indexed agent, uint64 startTime)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) FilterStartTimeSet(opts *bind.FilterOpts, agent []common.Address) (*LaunchRadarModuleStartTimeSetIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _LaunchRadarModule.contract.FilterLogs(opts, "StartTimeSet", agentRule)
	if err != nil {
		return nil, err
	}
	return &LaunchRadarModuleStartTimeSetIterator{contract: _LaunchRadarModule.contract, event: "StartTimeSet", logs: logs, sub: sub}, nil
}

// WatchStartTimeSet is a free log subscription operation binding the contract event 0xaa1bcd70b51a29ec01d817f60e40a1ce21afc80546629beda8dd88e01ea0856d.
//
// Solidity: event StartTimeSet(address indexed agent, uint64 startTime)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) WatchStartTimeSet(opts *bind.WatchOpts, sink chan<- *LaunchRadarModuleStartTimeSet, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _LaunchRadarModule.contract.WatchLogs(opts, "StartTimeSet", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LaunchRadarModuleStartTimeSet)
				if err := _LaunchRadarModule.contract.UnpackLog(event, "StartTimeSet", log); err != nil {
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

// ParseStartTimeSet is a log parse operation binding the contract event 0xaa1bcd70b51a29ec01d817f60e40a1ce21afc80546629beda8dd88e01ea0856d.
//
// Solidity: event StartTimeSet(address indexed agent, uint64 startTime)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) ParseStartTimeSet(log types.Log) (*LaunchRadarModuleStartTimeSet, error) {
	event := new(LaunchRadarModuleStartTimeSet)
	if err := _LaunchRadarModule.contract.UnpackLog(event, "StartTimeSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LaunchRadarModuleWhitelistToggledIterator is returned from FilterWhitelistToggled and is used to iterate over the raw logs and unpacked data for WhitelistToggled events raised by the LaunchRadarModule contract.
type LaunchRadarModuleWhitelistToggledIterator struct {
	Event *LaunchRadarModuleWhitelistToggled // Event containing the contract specifics and raw log

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
func (it *LaunchRadarModuleWhitelistToggledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LaunchRadarModuleWhitelistToggled)
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
		it.Event = new(LaunchRadarModuleWhitelistToggled)
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
func (it *LaunchRadarModuleWhitelistToggledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LaunchRadarModuleWhitelistToggledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LaunchRadarModuleWhitelistToggled represents a WhitelistToggled event raised by the LaunchRadarModule contract.
type LaunchRadarModuleWhitelistToggled struct {
	Agent common.Address
	On    bool
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterWhitelistToggled is a free log retrieval operation binding the contract event 0xdb17f7dcc6b9d3c099107ae3eea3dcc0ce959cda4b7b164cac7f8f9c84139719.
//
// Solidity: event WhitelistToggled(address indexed agent, bool on)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) FilterWhitelistToggled(opts *bind.FilterOpts, agent []common.Address) (*LaunchRadarModuleWhitelistToggledIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _LaunchRadarModule.contract.FilterLogs(opts, "WhitelistToggled", agentRule)
	if err != nil {
		return nil, err
	}
	return &LaunchRadarModuleWhitelistToggledIterator{contract: _LaunchRadarModule.contract, event: "WhitelistToggled", logs: logs, sub: sub}, nil
}

// WatchWhitelistToggled is a free log subscription operation binding the contract event 0xdb17f7dcc6b9d3c099107ae3eea3dcc0ce959cda4b7b164cac7f8f9c84139719.
//
// Solidity: event WhitelistToggled(address indexed agent, bool on)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) WatchWhitelistToggled(opts *bind.WatchOpts, sink chan<- *LaunchRadarModuleWhitelistToggled, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _LaunchRadarModule.contract.WatchLogs(opts, "WhitelistToggled", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LaunchRadarModuleWhitelistToggled)
				if err := _LaunchRadarModule.contract.UnpackLog(event, "WhitelistToggled", log); err != nil {
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

// ParseWhitelistToggled is a log parse operation binding the contract event 0xdb17f7dcc6b9d3c099107ae3eea3dcc0ce959cda4b7b164cac7f8f9c84139719.
//
// Solidity: event WhitelistToggled(address indexed agent, bool on)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) ParseWhitelistToggled(log types.Log) (*LaunchRadarModuleWhitelistToggled, error) {
	event := new(LaunchRadarModuleWhitelistToggled)
	if err := _LaunchRadarModule.contract.UnpackLog(event, "WhitelistToggled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LaunchRadarModuleWhitelistUpdatedIterator is returned from FilterWhitelistUpdated and is used to iterate over the raw logs and unpacked data for WhitelistUpdated events raised by the LaunchRadarModule contract.
type LaunchRadarModuleWhitelistUpdatedIterator struct {
	Event *LaunchRadarModuleWhitelistUpdated // Event containing the contract specifics and raw log

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
func (it *LaunchRadarModuleWhitelistUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LaunchRadarModuleWhitelistUpdated)
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
		it.Event = new(LaunchRadarModuleWhitelistUpdated)
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
func (it *LaunchRadarModuleWhitelistUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LaunchRadarModuleWhitelistUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LaunchRadarModuleWhitelistUpdated represents a WhitelistUpdated event raised by the LaunchRadarModule contract.
type LaunchRadarModuleWhitelistUpdated struct {
	Agent common.Address
	Count *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterWhitelistUpdated is a free log retrieval operation binding the contract event 0x226d670a329c4a93cf8c1a5baeceda320e89031fe0a65343c51678bd8b5a652e.
//
// Solidity: event WhitelistUpdated(address indexed agent, uint256 count)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) FilterWhitelistUpdated(opts *bind.FilterOpts, agent []common.Address) (*LaunchRadarModuleWhitelistUpdatedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _LaunchRadarModule.contract.FilterLogs(opts, "WhitelistUpdated", agentRule)
	if err != nil {
		return nil, err
	}
	return &LaunchRadarModuleWhitelistUpdatedIterator{contract: _LaunchRadarModule.contract, event: "WhitelistUpdated", logs: logs, sub: sub}, nil
}

// WatchWhitelistUpdated is a free log subscription operation binding the contract event 0x226d670a329c4a93cf8c1a5baeceda320e89031fe0a65343c51678bd8b5a652e.
//
// Solidity: event WhitelistUpdated(address indexed agent, uint256 count)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) WatchWhitelistUpdated(opts *bind.WatchOpts, sink chan<- *LaunchRadarModuleWhitelistUpdated, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _LaunchRadarModule.contract.WatchLogs(opts, "WhitelistUpdated", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LaunchRadarModuleWhitelistUpdated)
				if err := _LaunchRadarModule.contract.UnpackLog(event, "WhitelistUpdated", log); err != nil {
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

// ParseWhitelistUpdated is a log parse operation binding the contract event 0x226d670a329c4a93cf8c1a5baeceda320e89031fe0a65343c51678bd8b5a652e.
//
// Solidity: event WhitelistUpdated(address indexed agent, uint256 count)
func (_LaunchRadarModule *LaunchRadarModuleFilterer) ParseWhitelistUpdated(log types.Log) (*LaunchRadarModuleWhitelistUpdated, error) {
	event := new(LaunchRadarModuleWhitelistUpdated)
	if err := _LaunchRadarModule.contract.UnpackLog(event, "WhitelistUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
