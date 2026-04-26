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

// ExistingTokenModuleConfig is an auto generated low-level Go binding around an user-defined struct.
type ExistingTokenModuleConfig struct {
	Curve      common.Address
	AgentAdmin common.Address
	Supply     *big.Int
	Wrapped    bool
}

// ExistingTokenModuleMetaData contains all meta data concerning the ExistingTokenModule contract.
var ExistingTokenModuleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"configs\",\"inputs\":[{\"name\":\"externalToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"curve\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"supply\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"wrapped\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configure\",\"inputs\":[{\"name\":\"externalToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"curve\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"supply\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getConfig\",\"inputs\":[{\"name\":\"externalToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"cfg\",\"type\":\"tuple\",\"internalType\":\"structExistingTokenModule.Config\",\"components\":[{\"name\":\"curve\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"supply\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"wrapped\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"wrap\",\"inputs\":[{\"name\":\"externalToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Configured\",\"inputs\":[{\"name\":\"externalToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"curve\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"supply\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"agentAdmin\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Wrapped\",\"inputs\":[{\"name\":\"externalToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"depositor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientPreDeposit\",\"inputs\":[{\"name\":\"held\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"required\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidExternalToken\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroSupply\",\"inputs\":[]}]",
}

// ExistingTokenModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use ExistingTokenModuleMetaData.ABI instead.
var ExistingTokenModuleABI = ExistingTokenModuleMetaData.ABI

// ExistingTokenModule is an auto generated Go binding around an Ethereum contract.
type ExistingTokenModule struct {
	ExistingTokenModuleCaller     // Read-only binding to the contract
	ExistingTokenModuleTransactor // Write-only binding to the contract
	ExistingTokenModuleFilterer   // Log filterer for contract events
}

// ExistingTokenModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type ExistingTokenModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExistingTokenModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ExistingTokenModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExistingTokenModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ExistingTokenModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ExistingTokenModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ExistingTokenModuleSession struct {
	Contract     *ExistingTokenModule // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ExistingTokenModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ExistingTokenModuleCallerSession struct {
	Contract *ExistingTokenModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ExistingTokenModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ExistingTokenModuleTransactorSession struct {
	Contract     *ExistingTokenModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ExistingTokenModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type ExistingTokenModuleRaw struct {
	Contract *ExistingTokenModule // Generic contract binding to access the raw methods on
}

// ExistingTokenModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ExistingTokenModuleCallerRaw struct {
	Contract *ExistingTokenModuleCaller // Generic read-only contract binding to access the raw methods on
}

// ExistingTokenModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ExistingTokenModuleTransactorRaw struct {
	Contract *ExistingTokenModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewExistingTokenModule creates a new instance of ExistingTokenModule, bound to a specific deployed contract.
func NewExistingTokenModule(address common.Address, backend bind.ContractBackend) (*ExistingTokenModule, error) {
	contract, err := bindExistingTokenModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ExistingTokenModule{ExistingTokenModuleCaller: ExistingTokenModuleCaller{contract: contract}, ExistingTokenModuleTransactor: ExistingTokenModuleTransactor{contract: contract}, ExistingTokenModuleFilterer: ExistingTokenModuleFilterer{contract: contract}}, nil
}

// NewExistingTokenModuleCaller creates a new read-only instance of ExistingTokenModule, bound to a specific deployed contract.
func NewExistingTokenModuleCaller(address common.Address, caller bind.ContractCaller) (*ExistingTokenModuleCaller, error) {
	contract, err := bindExistingTokenModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ExistingTokenModuleCaller{contract: contract}, nil
}

// NewExistingTokenModuleTransactor creates a new write-only instance of ExistingTokenModule, bound to a specific deployed contract.
func NewExistingTokenModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*ExistingTokenModuleTransactor, error) {
	contract, err := bindExistingTokenModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ExistingTokenModuleTransactor{contract: contract}, nil
}

// NewExistingTokenModuleFilterer creates a new log filterer instance of ExistingTokenModule, bound to a specific deployed contract.
func NewExistingTokenModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*ExistingTokenModuleFilterer, error) {
	contract, err := bindExistingTokenModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ExistingTokenModuleFilterer{contract: contract}, nil
}

// bindExistingTokenModule binds a generic wrapper to an already deployed contract.
func bindExistingTokenModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ExistingTokenModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExistingTokenModule *ExistingTokenModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExistingTokenModule.Contract.ExistingTokenModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExistingTokenModule *ExistingTokenModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.ExistingTokenModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExistingTokenModule *ExistingTokenModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.ExistingTokenModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ExistingTokenModule *ExistingTokenModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ExistingTokenModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ExistingTokenModule *ExistingTokenModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ExistingTokenModule *ExistingTokenModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.contract.Transact(opts, method, params...)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address externalToken) view returns(address curve, address agentAdmin, uint256 supply, bool wrapped)
func (_ExistingTokenModule *ExistingTokenModuleCaller) Configs(opts *bind.CallOpts, externalToken common.Address) (struct {
	Curve      common.Address
	AgentAdmin common.Address
	Supply     *big.Int
	Wrapped    bool
}, error) {
	var out []interface{}
	err := _ExistingTokenModule.contract.Call(opts, &out, "configs", externalToken)

	outstruct := new(struct {
		Curve      common.Address
		AgentAdmin common.Address
		Supply     *big.Int
		Wrapped    bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Curve = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.AgentAdmin = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Supply = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Wrapped = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address externalToken) view returns(address curve, address agentAdmin, uint256 supply, bool wrapped)
func (_ExistingTokenModule *ExistingTokenModuleSession) Configs(externalToken common.Address) (struct {
	Curve      common.Address
	AgentAdmin common.Address
	Supply     *big.Int
	Wrapped    bool
}, error) {
	return _ExistingTokenModule.Contract.Configs(&_ExistingTokenModule.CallOpts, externalToken)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address externalToken) view returns(address curve, address agentAdmin, uint256 supply, bool wrapped)
func (_ExistingTokenModule *ExistingTokenModuleCallerSession) Configs(externalToken common.Address) (struct {
	Curve      common.Address
	AgentAdmin common.Address
	Supply     *big.Int
	Wrapped    bool
}, error) {
	return _ExistingTokenModule.Contract.Configs(&_ExistingTokenModule.CallOpts, externalToken)
}

// GetConfig is a free data retrieval call binding the contract method 0xe48a5f7b.
//
// Solidity: function getConfig(address externalToken) view returns((address,address,uint256,bool) cfg)
func (_ExistingTokenModule *ExistingTokenModuleCaller) GetConfig(opts *bind.CallOpts, externalToken common.Address) (ExistingTokenModuleConfig, error) {
	var out []interface{}
	err := _ExistingTokenModule.contract.Call(opts, &out, "getConfig", externalToken)

	if err != nil {
		return *new(ExistingTokenModuleConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(ExistingTokenModuleConfig)).(*ExistingTokenModuleConfig)

	return out0, err

}

// GetConfig is a free data retrieval call binding the contract method 0xe48a5f7b.
//
// Solidity: function getConfig(address externalToken) view returns((address,address,uint256,bool) cfg)
func (_ExistingTokenModule *ExistingTokenModuleSession) GetConfig(externalToken common.Address) (ExistingTokenModuleConfig, error) {
	return _ExistingTokenModule.Contract.GetConfig(&_ExistingTokenModule.CallOpts, externalToken)
}

// GetConfig is a free data retrieval call binding the contract method 0xe48a5f7b.
//
// Solidity: function getConfig(address externalToken) view returns((address,address,uint256,bool) cfg)
func (_ExistingTokenModule *ExistingTokenModuleCallerSession) GetConfig(externalToken common.Address) (ExistingTokenModuleConfig, error) {
	return _ExistingTokenModule.Contract.GetConfig(&_ExistingTokenModule.CallOpts, externalToken)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExistingTokenModule *ExistingTokenModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ExistingTokenModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExistingTokenModule *ExistingTokenModuleSession) Owner() (common.Address, error) {
	return _ExistingTokenModule.Contract.Owner(&_ExistingTokenModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ExistingTokenModule *ExistingTokenModuleCallerSession) Owner() (common.Address, error) {
	return _ExistingTokenModule.Contract.Owner(&_ExistingTokenModule.CallOpts)
}

// Configure is a paid mutator transaction binding the contract method 0xf1f58d8a.
//
// Solidity: function configure(address externalToken, address curve, uint256 supply, address admin) returns()
func (_ExistingTokenModule *ExistingTokenModuleTransactor) Configure(opts *bind.TransactOpts, externalToken common.Address, curve common.Address, supply *big.Int, admin common.Address) (*types.Transaction, error) {
	return _ExistingTokenModule.contract.Transact(opts, "configure", externalToken, curve, supply, admin)
}

// Configure is a paid mutator transaction binding the contract method 0xf1f58d8a.
//
// Solidity: function configure(address externalToken, address curve, uint256 supply, address admin) returns()
func (_ExistingTokenModule *ExistingTokenModuleSession) Configure(externalToken common.Address, curve common.Address, supply *big.Int, admin common.Address) (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.Configure(&_ExistingTokenModule.TransactOpts, externalToken, curve, supply, admin)
}

// Configure is a paid mutator transaction binding the contract method 0xf1f58d8a.
//
// Solidity: function configure(address externalToken, address curve, uint256 supply, address admin) returns()
func (_ExistingTokenModule *ExistingTokenModuleTransactorSession) Configure(externalToken common.Address, curve common.Address, supply *big.Int, admin common.Address) (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.Configure(&_ExistingTokenModule.TransactOpts, externalToken, curve, supply, admin)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExistingTokenModule *ExistingTokenModuleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ExistingTokenModule.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExistingTokenModule *ExistingTokenModuleSession) RenounceOwnership() (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.RenounceOwnership(&_ExistingTokenModule.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ExistingTokenModule *ExistingTokenModuleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.RenounceOwnership(&_ExistingTokenModule.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExistingTokenModule *ExistingTokenModuleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ExistingTokenModule.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExistingTokenModule *ExistingTokenModuleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.TransferOwnership(&_ExistingTokenModule.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ExistingTokenModule *ExistingTokenModuleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.TransferOwnership(&_ExistingTokenModule.TransactOpts, newOwner)
}

// Wrap is a paid mutator transaction binding the contract method 0xbf376c7a.
//
// Solidity: function wrap(address externalToken, uint256 amount) returns()
func (_ExistingTokenModule *ExistingTokenModuleTransactor) Wrap(opts *bind.TransactOpts, externalToken common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ExistingTokenModule.contract.Transact(opts, "wrap", externalToken, amount)
}

// Wrap is a paid mutator transaction binding the contract method 0xbf376c7a.
//
// Solidity: function wrap(address externalToken, uint256 amount) returns()
func (_ExistingTokenModule *ExistingTokenModuleSession) Wrap(externalToken common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.Wrap(&_ExistingTokenModule.TransactOpts, externalToken, amount)
}

// Wrap is a paid mutator transaction binding the contract method 0xbf376c7a.
//
// Solidity: function wrap(address externalToken, uint256 amount) returns()
func (_ExistingTokenModule *ExistingTokenModuleTransactorSession) Wrap(externalToken common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ExistingTokenModule.Contract.Wrap(&_ExistingTokenModule.TransactOpts, externalToken, amount)
}

// ExistingTokenModuleConfiguredIterator is returned from FilterConfigured and is used to iterate over the raw logs and unpacked data for Configured events raised by the ExistingTokenModule contract.
type ExistingTokenModuleConfiguredIterator struct {
	Event *ExistingTokenModuleConfigured // Event containing the contract specifics and raw log

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
func (it *ExistingTokenModuleConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExistingTokenModuleConfigured)
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
		it.Event = new(ExistingTokenModuleConfigured)
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
func (it *ExistingTokenModuleConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExistingTokenModuleConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExistingTokenModuleConfigured represents a Configured event raised by the ExistingTokenModule contract.
type ExistingTokenModuleConfigured struct {
	ExternalToken common.Address
	Curve         common.Address
	Supply        *big.Int
	AgentAdmin    common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterConfigured is a free log retrieval operation binding the contract event 0x0222b19aaf167ba785a1c0f00f2c772ca36639f19257568a756da9c4cb32a3e3.
//
// Solidity: event Configured(address indexed externalToken, address indexed curve, uint256 supply, address indexed agentAdmin)
func (_ExistingTokenModule *ExistingTokenModuleFilterer) FilterConfigured(opts *bind.FilterOpts, externalToken []common.Address, curve []common.Address, agentAdmin []common.Address) (*ExistingTokenModuleConfiguredIterator, error) {

	var externalTokenRule []interface{}
	for _, externalTokenItem := range externalToken {
		externalTokenRule = append(externalTokenRule, externalTokenItem)
	}
	var curveRule []interface{}
	for _, curveItem := range curve {
		curveRule = append(curveRule, curveItem)
	}

	var agentAdminRule []interface{}
	for _, agentAdminItem := range agentAdmin {
		agentAdminRule = append(agentAdminRule, agentAdminItem)
	}

	logs, sub, err := _ExistingTokenModule.contract.FilterLogs(opts, "Configured", externalTokenRule, curveRule, agentAdminRule)
	if err != nil {
		return nil, err
	}
	return &ExistingTokenModuleConfiguredIterator{contract: _ExistingTokenModule.contract, event: "Configured", logs: logs, sub: sub}, nil
}

// WatchConfigured is a free log subscription operation binding the contract event 0x0222b19aaf167ba785a1c0f00f2c772ca36639f19257568a756da9c4cb32a3e3.
//
// Solidity: event Configured(address indexed externalToken, address indexed curve, uint256 supply, address indexed agentAdmin)
func (_ExistingTokenModule *ExistingTokenModuleFilterer) WatchConfigured(opts *bind.WatchOpts, sink chan<- *ExistingTokenModuleConfigured, externalToken []common.Address, curve []common.Address, agentAdmin []common.Address) (event.Subscription, error) {

	var externalTokenRule []interface{}
	for _, externalTokenItem := range externalToken {
		externalTokenRule = append(externalTokenRule, externalTokenItem)
	}
	var curveRule []interface{}
	for _, curveItem := range curve {
		curveRule = append(curveRule, curveItem)
	}

	var agentAdminRule []interface{}
	for _, agentAdminItem := range agentAdmin {
		agentAdminRule = append(agentAdminRule, agentAdminItem)
	}

	logs, sub, err := _ExistingTokenModule.contract.WatchLogs(opts, "Configured", externalTokenRule, curveRule, agentAdminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExistingTokenModuleConfigured)
				if err := _ExistingTokenModule.contract.UnpackLog(event, "Configured", log); err != nil {
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

// ParseConfigured is a log parse operation binding the contract event 0x0222b19aaf167ba785a1c0f00f2c772ca36639f19257568a756da9c4cb32a3e3.
//
// Solidity: event Configured(address indexed externalToken, address indexed curve, uint256 supply, address indexed agentAdmin)
func (_ExistingTokenModule *ExistingTokenModuleFilterer) ParseConfigured(log types.Log) (*ExistingTokenModuleConfigured, error) {
	event := new(ExistingTokenModuleConfigured)
	if err := _ExistingTokenModule.contract.UnpackLog(event, "Configured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExistingTokenModuleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ExistingTokenModule contract.
type ExistingTokenModuleOwnershipTransferredIterator struct {
	Event *ExistingTokenModuleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ExistingTokenModuleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExistingTokenModuleOwnershipTransferred)
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
		it.Event = new(ExistingTokenModuleOwnershipTransferred)
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
func (it *ExistingTokenModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExistingTokenModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExistingTokenModuleOwnershipTransferred represents a OwnershipTransferred event raised by the ExistingTokenModule contract.
type ExistingTokenModuleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExistingTokenModule *ExistingTokenModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ExistingTokenModuleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExistingTokenModule.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ExistingTokenModuleOwnershipTransferredIterator{contract: _ExistingTokenModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ExistingTokenModule *ExistingTokenModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ExistingTokenModuleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ExistingTokenModule.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExistingTokenModuleOwnershipTransferred)
				if err := _ExistingTokenModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ExistingTokenModule *ExistingTokenModuleFilterer) ParseOwnershipTransferred(log types.Log) (*ExistingTokenModuleOwnershipTransferred, error) {
	event := new(ExistingTokenModuleOwnershipTransferred)
	if err := _ExistingTokenModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ExistingTokenModuleWrappedIterator is returned from FilterWrapped and is used to iterate over the raw logs and unpacked data for Wrapped events raised by the ExistingTokenModule contract.
type ExistingTokenModuleWrappedIterator struct {
	Event *ExistingTokenModuleWrapped // Event containing the contract specifics and raw log

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
func (it *ExistingTokenModuleWrappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ExistingTokenModuleWrapped)
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
		it.Event = new(ExistingTokenModuleWrapped)
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
func (it *ExistingTokenModuleWrappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ExistingTokenModuleWrappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ExistingTokenModuleWrapped represents a Wrapped event raised by the ExistingTokenModule contract.
type ExistingTokenModuleWrapped struct {
	ExternalToken common.Address
	Amount        *big.Int
	Depositor     common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterWrapped is a free log retrieval operation binding the contract event 0x1dcb5a508170be4d49f81a2e679aea5a1a064f5e7c34d73cd8bc670305d74b16.
//
// Solidity: event Wrapped(address indexed externalToken, uint256 amount, address indexed depositor)
func (_ExistingTokenModule *ExistingTokenModuleFilterer) FilterWrapped(opts *bind.FilterOpts, externalToken []common.Address, depositor []common.Address) (*ExistingTokenModuleWrappedIterator, error) {

	var externalTokenRule []interface{}
	for _, externalTokenItem := range externalToken {
		externalTokenRule = append(externalTokenRule, externalTokenItem)
	}

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _ExistingTokenModule.contract.FilterLogs(opts, "Wrapped", externalTokenRule, depositorRule)
	if err != nil {
		return nil, err
	}
	return &ExistingTokenModuleWrappedIterator{contract: _ExistingTokenModule.contract, event: "Wrapped", logs: logs, sub: sub}, nil
}

// WatchWrapped is a free log subscription operation binding the contract event 0x1dcb5a508170be4d49f81a2e679aea5a1a064f5e7c34d73cd8bc670305d74b16.
//
// Solidity: event Wrapped(address indexed externalToken, uint256 amount, address indexed depositor)
func (_ExistingTokenModule *ExistingTokenModuleFilterer) WatchWrapped(opts *bind.WatchOpts, sink chan<- *ExistingTokenModuleWrapped, externalToken []common.Address, depositor []common.Address) (event.Subscription, error) {

	var externalTokenRule []interface{}
	for _, externalTokenItem := range externalToken {
		externalTokenRule = append(externalTokenRule, externalTokenItem)
	}

	var depositorRule []interface{}
	for _, depositorItem := range depositor {
		depositorRule = append(depositorRule, depositorItem)
	}

	logs, sub, err := _ExistingTokenModule.contract.WatchLogs(opts, "Wrapped", externalTokenRule, depositorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ExistingTokenModuleWrapped)
				if err := _ExistingTokenModule.contract.UnpackLog(event, "Wrapped", log); err != nil {
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

// ParseWrapped is a log parse operation binding the contract event 0x1dcb5a508170be4d49f81a2e679aea5a1a064f5e7c34d73cd8bc670305d74b16.
//
// Solidity: event Wrapped(address indexed externalToken, uint256 amount, address indexed depositor)
func (_ExistingTokenModule *ExistingTokenModuleFilterer) ParseWrapped(log types.Log) (*ExistingTokenModuleWrapped, error) {
	event := new(ExistingTokenModuleWrapped)
	if err := _ExistingTokenModule.contract.UnpackLog(event, "Wrapped", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
