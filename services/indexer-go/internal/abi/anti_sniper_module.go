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

// AntiSniperModuleMetaData contains all meta data concerning the AntiSniperModule contract.
var AntiSniperModuleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"computeTax\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"tax\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"netAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configs\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"duration\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"startBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"endBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"configured\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configure\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"duration\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"startBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"endBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentBps\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Configured\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"duration\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"startBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"},{\"name\":\"endBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidBps\",\"inputs\":[{\"name\":\"bps\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"type\":\"error\",\"name\":\"InvalidDuration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// AntiSniperModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use AntiSniperModuleMetaData.ABI instead.
var AntiSniperModuleABI = AntiSniperModuleMetaData.ABI

// AntiSniperModule is an auto generated Go binding around an Ethereum contract.
type AntiSniperModule struct {
	AntiSniperModuleCaller     // Read-only binding to the contract
	AntiSniperModuleTransactor // Write-only binding to the contract
	AntiSniperModuleFilterer   // Log filterer for contract events
}

// AntiSniperModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type AntiSniperModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AntiSniperModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AntiSniperModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AntiSniperModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AntiSniperModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AntiSniperModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AntiSniperModuleSession struct {
	Contract     *AntiSniperModule // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AntiSniperModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AntiSniperModuleCallerSession struct {
	Contract *AntiSniperModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// AntiSniperModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AntiSniperModuleTransactorSession struct {
	Contract     *AntiSniperModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// AntiSniperModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type AntiSniperModuleRaw struct {
	Contract *AntiSniperModule // Generic contract binding to access the raw methods on
}

// AntiSniperModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AntiSniperModuleCallerRaw struct {
	Contract *AntiSniperModuleCaller // Generic read-only contract binding to access the raw methods on
}

// AntiSniperModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AntiSniperModuleTransactorRaw struct {
	Contract *AntiSniperModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAntiSniperModule creates a new instance of AntiSniperModule, bound to a specific deployed contract.
func NewAntiSniperModule(address common.Address, backend bind.ContractBackend) (*AntiSniperModule, error) {
	contract, err := bindAntiSniperModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AntiSniperModule{AntiSniperModuleCaller: AntiSniperModuleCaller{contract: contract}, AntiSniperModuleTransactor: AntiSniperModuleTransactor{contract: contract}, AntiSniperModuleFilterer: AntiSniperModuleFilterer{contract: contract}}, nil
}

// NewAntiSniperModuleCaller creates a new read-only instance of AntiSniperModule, bound to a specific deployed contract.
func NewAntiSniperModuleCaller(address common.Address, caller bind.ContractCaller) (*AntiSniperModuleCaller, error) {
	contract, err := bindAntiSniperModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AntiSniperModuleCaller{contract: contract}, nil
}

// NewAntiSniperModuleTransactor creates a new write-only instance of AntiSniperModule, bound to a specific deployed contract.
func NewAntiSniperModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*AntiSniperModuleTransactor, error) {
	contract, err := bindAntiSniperModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AntiSniperModuleTransactor{contract: contract}, nil
}

// NewAntiSniperModuleFilterer creates a new log filterer instance of AntiSniperModule, bound to a specific deployed contract.
func NewAntiSniperModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*AntiSniperModuleFilterer, error) {
	contract, err := bindAntiSniperModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AntiSniperModuleFilterer{contract: contract}, nil
}

// bindAntiSniperModule binds a generic wrapper to an already deployed contract.
func bindAntiSniperModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AntiSniperModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AntiSniperModule *AntiSniperModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AntiSniperModule.Contract.AntiSniperModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AntiSniperModule *AntiSniperModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AntiSniperModule.Contract.AntiSniperModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AntiSniperModule *AntiSniperModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AntiSniperModule.Contract.AntiSniperModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AntiSniperModule *AntiSniperModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AntiSniperModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AntiSniperModule *AntiSniperModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AntiSniperModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AntiSniperModule *AntiSniperModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AntiSniperModule.Contract.contract.Transact(opts, method, params...)
}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint16)
func (_AntiSniperModule *AntiSniperModuleCaller) MAXBPS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _AntiSniperModule.contract.Call(opts, &out, "MAX_BPS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint16)
func (_AntiSniperModule *AntiSniperModuleSession) MAXBPS() (uint16, error) {
	return _AntiSniperModule.Contract.MAXBPS(&_AntiSniperModule.CallOpts)
}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint16)
func (_AntiSniperModule *AntiSniperModuleCallerSession) MAXBPS() (uint16, error) {
	return _AntiSniperModule.Contract.MAXBPS(&_AntiSniperModule.CallOpts)
}

// ComputeTax is a free data retrieval call binding the contract method 0xa723afc9.
//
// Solidity: function computeTax(address agent, uint256 amount) view returns(uint256 tax, uint256 netAmount)
func (_AntiSniperModule *AntiSniperModuleCaller) ComputeTax(opts *bind.CallOpts, agent common.Address, amount *big.Int) (struct {
	Tax       *big.Int
	NetAmount *big.Int
}, error) {
	var out []interface{}
	err := _AntiSniperModule.contract.Call(opts, &out, "computeTax", agent, amount)

	outstruct := new(struct {
		Tax       *big.Int
		NetAmount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Tax = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.NetAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ComputeTax is a free data retrieval call binding the contract method 0xa723afc9.
//
// Solidity: function computeTax(address agent, uint256 amount) view returns(uint256 tax, uint256 netAmount)
func (_AntiSniperModule *AntiSniperModuleSession) ComputeTax(agent common.Address, amount *big.Int) (struct {
	Tax       *big.Int
	NetAmount *big.Int
}, error) {
	return _AntiSniperModule.Contract.ComputeTax(&_AntiSniperModule.CallOpts, agent, amount)
}

// ComputeTax is a free data retrieval call binding the contract method 0xa723afc9.
//
// Solidity: function computeTax(address agent, uint256 amount) view returns(uint256 tax, uint256 netAmount)
func (_AntiSniperModule *AntiSniperModuleCallerSession) ComputeTax(agent common.Address, amount *big.Int) (struct {
	Tax       *big.Int
	NetAmount *big.Int
}, error) {
	return _AntiSniperModule.Contract.ComputeTax(&_AntiSniperModule.CallOpts, agent, amount)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps, bool configured)
func (_AntiSniperModule *AntiSniperModuleCaller) Configs(opts *bind.CallOpts, agent common.Address) (struct {
	StartTime  uint64
	Duration   uint32
	StartBps   uint16
	EndBps     uint16
	Configured bool
}, error) {
	var out []interface{}
	err := _AntiSniperModule.contract.Call(opts, &out, "configs", agent)

	outstruct := new(struct {
		StartTime  uint64
		Duration   uint32
		StartBps   uint16
		EndBps     uint16
		Configured bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StartTime = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.Duration = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.StartBps = *abi.ConvertType(out[2], new(uint16)).(*uint16)
	outstruct.EndBps = *abi.ConvertType(out[3], new(uint16)).(*uint16)
	outstruct.Configured = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps, bool configured)
func (_AntiSniperModule *AntiSniperModuleSession) Configs(agent common.Address) (struct {
	StartTime  uint64
	Duration   uint32
	StartBps   uint16
	EndBps     uint16
	Configured bool
}, error) {
	return _AntiSniperModule.Contract.Configs(&_AntiSniperModule.CallOpts, agent)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps, bool configured)
func (_AntiSniperModule *AntiSniperModuleCallerSession) Configs(agent common.Address) (struct {
	StartTime  uint64
	Duration   uint32
	StartBps   uint16
	EndBps     uint16
	Configured bool
}, error) {
	return _AntiSniperModule.Contract.Configs(&_AntiSniperModule.CallOpts, agent)
}

// CurrentBps is a free data retrieval call binding the contract method 0x430ade8c.
//
// Solidity: function currentBps(address agent) view returns(uint16)
func (_AntiSniperModule *AntiSniperModuleCaller) CurrentBps(opts *bind.CallOpts, agent common.Address) (uint16, error) {
	var out []interface{}
	err := _AntiSniperModule.contract.Call(opts, &out, "currentBps", agent)

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// CurrentBps is a free data retrieval call binding the contract method 0x430ade8c.
//
// Solidity: function currentBps(address agent) view returns(uint16)
func (_AntiSniperModule *AntiSniperModuleSession) CurrentBps(agent common.Address) (uint16, error) {
	return _AntiSniperModule.Contract.CurrentBps(&_AntiSniperModule.CallOpts, agent)
}

// CurrentBps is a free data retrieval call binding the contract method 0x430ade8c.
//
// Solidity: function currentBps(address agent) view returns(uint16)
func (_AntiSniperModule *AntiSniperModuleCallerSession) CurrentBps(agent common.Address) (uint16, error) {
	return _AntiSniperModule.Contract.CurrentBps(&_AntiSniperModule.CallOpts, agent)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AntiSniperModule *AntiSniperModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AntiSniperModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AntiSniperModule *AntiSniperModuleSession) Owner() (common.Address, error) {
	return _AntiSniperModule.Contract.Owner(&_AntiSniperModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AntiSniperModule *AntiSniperModuleCallerSession) Owner() (common.Address, error) {
	return _AntiSniperModule.Contract.Owner(&_AntiSniperModule.CallOpts)
}

// Configure is a paid mutator transaction binding the contract method 0x82e16172.
//
// Solidity: function configure(address agent, uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps) returns()
func (_AntiSniperModule *AntiSniperModuleTransactor) Configure(opts *bind.TransactOpts, agent common.Address, startTime uint64, duration uint32, startBps uint16, endBps uint16) (*types.Transaction, error) {
	return _AntiSniperModule.contract.Transact(opts, "configure", agent, startTime, duration, startBps, endBps)
}

// Configure is a paid mutator transaction binding the contract method 0x82e16172.
//
// Solidity: function configure(address agent, uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps) returns()
func (_AntiSniperModule *AntiSniperModuleSession) Configure(agent common.Address, startTime uint64, duration uint32, startBps uint16, endBps uint16) (*types.Transaction, error) {
	return _AntiSniperModule.Contract.Configure(&_AntiSniperModule.TransactOpts, agent, startTime, duration, startBps, endBps)
}

// Configure is a paid mutator transaction binding the contract method 0x82e16172.
//
// Solidity: function configure(address agent, uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps) returns()
func (_AntiSniperModule *AntiSniperModuleTransactorSession) Configure(agent common.Address, startTime uint64, duration uint32, startBps uint16, endBps uint16) (*types.Transaction, error) {
	return _AntiSniperModule.Contract.Configure(&_AntiSniperModule.TransactOpts, agent, startTime, duration, startBps, endBps)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AntiSniperModule *AntiSniperModuleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AntiSniperModule.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AntiSniperModule *AntiSniperModuleSession) RenounceOwnership() (*types.Transaction, error) {
	return _AntiSniperModule.Contract.RenounceOwnership(&_AntiSniperModule.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AntiSniperModule *AntiSniperModuleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AntiSniperModule.Contract.RenounceOwnership(&_AntiSniperModule.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AntiSniperModule *AntiSniperModuleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AntiSniperModule.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AntiSniperModule *AntiSniperModuleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AntiSniperModule.Contract.TransferOwnership(&_AntiSniperModule.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AntiSniperModule *AntiSniperModuleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AntiSniperModule.Contract.TransferOwnership(&_AntiSniperModule.TransactOpts, newOwner)
}

// AntiSniperModuleConfiguredIterator is returned from FilterConfigured and is used to iterate over the raw logs and unpacked data for Configured events raised by the AntiSniperModule contract.
type AntiSniperModuleConfiguredIterator struct {
	Event *AntiSniperModuleConfigured // Event containing the contract specifics and raw log

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
func (it *AntiSniperModuleConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AntiSniperModuleConfigured)
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
		it.Event = new(AntiSniperModuleConfigured)
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
func (it *AntiSniperModuleConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AntiSniperModuleConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AntiSniperModuleConfigured represents a Configured event raised by the AntiSniperModule contract.
type AntiSniperModuleConfigured struct {
	Agent     common.Address
	StartTime uint64
	Duration  uint32
	StartBps  uint16
	EndBps    uint16
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterConfigured is a free log retrieval operation binding the contract event 0x48f3b813c2853e7ceefe8271c4667cd6503447be4a04b5ceeb4f4c5814446f1c.
//
// Solidity: event Configured(address indexed agent, uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps)
func (_AntiSniperModule *AntiSniperModuleFilterer) FilterConfigured(opts *bind.FilterOpts, agent []common.Address) (*AntiSniperModuleConfiguredIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _AntiSniperModule.contract.FilterLogs(opts, "Configured", agentRule)
	if err != nil {
		return nil, err
	}
	return &AntiSniperModuleConfiguredIterator{contract: _AntiSniperModule.contract, event: "Configured", logs: logs, sub: sub}, nil
}

// WatchConfigured is a free log subscription operation binding the contract event 0x48f3b813c2853e7ceefe8271c4667cd6503447be4a04b5ceeb4f4c5814446f1c.
//
// Solidity: event Configured(address indexed agent, uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps)
func (_AntiSniperModule *AntiSniperModuleFilterer) WatchConfigured(opts *bind.WatchOpts, sink chan<- *AntiSniperModuleConfigured, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _AntiSniperModule.contract.WatchLogs(opts, "Configured", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AntiSniperModuleConfigured)
				if err := _AntiSniperModule.contract.UnpackLog(event, "Configured", log); err != nil {
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

// ParseConfigured is a log parse operation binding the contract event 0x48f3b813c2853e7ceefe8271c4667cd6503447be4a04b5ceeb4f4c5814446f1c.
//
// Solidity: event Configured(address indexed agent, uint64 startTime, uint32 duration, uint16 startBps, uint16 endBps)
func (_AntiSniperModule *AntiSniperModuleFilterer) ParseConfigured(log types.Log) (*AntiSniperModuleConfigured, error) {
	event := new(AntiSniperModuleConfigured)
	if err := _AntiSniperModule.contract.UnpackLog(event, "Configured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AntiSniperModuleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AntiSniperModule contract.
type AntiSniperModuleOwnershipTransferredIterator struct {
	Event *AntiSniperModuleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AntiSniperModuleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AntiSniperModuleOwnershipTransferred)
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
		it.Event = new(AntiSniperModuleOwnershipTransferred)
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
func (it *AntiSniperModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AntiSniperModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AntiSniperModuleOwnershipTransferred represents a OwnershipTransferred event raised by the AntiSniperModule contract.
type AntiSniperModuleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AntiSniperModule *AntiSniperModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AntiSniperModuleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AntiSniperModule.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AntiSniperModuleOwnershipTransferredIterator{contract: _AntiSniperModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AntiSniperModule *AntiSniperModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AntiSniperModuleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AntiSniperModule.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AntiSniperModuleOwnershipTransferred)
				if err := _AntiSniperModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_AntiSniperModule *AntiSniperModuleFilterer) ParseOwnershipTransferred(log types.Log) (*AntiSniperModuleOwnershipTransferred, error) {
	event := new(AntiSniperModuleOwnershipTransferred)
	if err := _AntiSniperModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
