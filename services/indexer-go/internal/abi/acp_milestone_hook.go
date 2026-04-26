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

// MilestoneHookMetaData contains all meta data concerning the MilestoneHook contract.
var MilestoneHookMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"hookName\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"milestones\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stageAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastStageRemainder\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalStages\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"completedStages\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"initialised\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onAccept\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onApprove\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onCancel\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onReject\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onSubmit\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"MilestoneCancelled\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"remainingStages\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MilestoneReleased\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"stage\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"totalStages\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AllStagesComplete\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidStages\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitialised\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// MilestoneHookABI is the input ABI used to generate the binding from.
// Deprecated: Use MilestoneHookMetaData.ABI instead.
var MilestoneHookABI = MilestoneHookMetaData.ABI

// MilestoneHook is an auto generated Go binding around an Ethereum contract.
type MilestoneHook struct {
	MilestoneHookCaller     // Read-only binding to the contract
	MilestoneHookTransactor // Write-only binding to the contract
	MilestoneHookFilterer   // Log filterer for contract events
}

// MilestoneHookCaller is an auto generated read-only Go binding around an Ethereum contract.
type MilestoneHookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MilestoneHookTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MilestoneHookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MilestoneHookFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MilestoneHookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MilestoneHookSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MilestoneHookSession struct {
	Contract     *MilestoneHook    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MilestoneHookCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MilestoneHookCallerSession struct {
	Contract *MilestoneHookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// MilestoneHookTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MilestoneHookTransactorSession struct {
	Contract     *MilestoneHookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// MilestoneHookRaw is an auto generated low-level Go binding around an Ethereum contract.
type MilestoneHookRaw struct {
	Contract *MilestoneHook // Generic contract binding to access the raw methods on
}

// MilestoneHookCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MilestoneHookCallerRaw struct {
	Contract *MilestoneHookCaller // Generic read-only contract binding to access the raw methods on
}

// MilestoneHookTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MilestoneHookTransactorRaw struct {
	Contract *MilestoneHookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMilestoneHook creates a new instance of MilestoneHook, bound to a specific deployed contract.
func NewMilestoneHook(address common.Address, backend bind.ContractBackend) (*MilestoneHook, error) {
	contract, err := bindMilestoneHook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MilestoneHook{MilestoneHookCaller: MilestoneHookCaller{contract: contract}, MilestoneHookTransactor: MilestoneHookTransactor{contract: contract}, MilestoneHookFilterer: MilestoneHookFilterer{contract: contract}}, nil
}

// NewMilestoneHookCaller creates a new read-only instance of MilestoneHook, bound to a specific deployed contract.
func NewMilestoneHookCaller(address common.Address, caller bind.ContractCaller) (*MilestoneHookCaller, error) {
	contract, err := bindMilestoneHook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MilestoneHookCaller{contract: contract}, nil
}

// NewMilestoneHookTransactor creates a new write-only instance of MilestoneHook, bound to a specific deployed contract.
func NewMilestoneHookTransactor(address common.Address, transactor bind.ContractTransactor) (*MilestoneHookTransactor, error) {
	contract, err := bindMilestoneHook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MilestoneHookTransactor{contract: contract}, nil
}

// NewMilestoneHookFilterer creates a new log filterer instance of MilestoneHook, bound to a specific deployed contract.
func NewMilestoneHookFilterer(address common.Address, filterer bind.ContractFilterer) (*MilestoneHookFilterer, error) {
	contract, err := bindMilestoneHook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MilestoneHookFilterer{contract: contract}, nil
}

// bindMilestoneHook binds a generic wrapper to an already deployed contract.
func bindMilestoneHook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MilestoneHookMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MilestoneHook *MilestoneHookRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MilestoneHook.Contract.MilestoneHookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MilestoneHook *MilestoneHookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MilestoneHook.Contract.MilestoneHookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MilestoneHook *MilestoneHookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MilestoneHook.Contract.MilestoneHookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MilestoneHook *MilestoneHookCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MilestoneHook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MilestoneHook *MilestoneHookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MilestoneHook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MilestoneHook *MilestoneHookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MilestoneHook.Contract.contract.Transact(opts, method, params...)
}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_MilestoneHook *MilestoneHookCaller) HookName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _MilestoneHook.contract.Call(opts, &out, "hookName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_MilestoneHook *MilestoneHookSession) HookName() (string, error) {
	return _MilestoneHook.Contract.HookName(&_MilestoneHook.CallOpts)
}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_MilestoneHook *MilestoneHookCallerSession) HookName() (string, error) {
	return _MilestoneHook.Contract.HookName(&_MilestoneHook.CallOpts)
}

// Milestones is a free data retrieval call binding the contract method 0xe89e4ed6.
//
// Solidity: function milestones(uint256 ) view returns(address agent, address token, uint256 stageAmount, uint256 lastStageRemainder, uint8 totalStages, uint8 completedStages, bool initialised)
func (_MilestoneHook *MilestoneHookCaller) Milestones(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Agent              common.Address
	Token              common.Address
	StageAmount        *big.Int
	LastStageRemainder *big.Int
	TotalStages        uint8
	CompletedStages    uint8
	Initialised        bool
}, error) {
	var out []interface{}
	err := _MilestoneHook.contract.Call(opts, &out, "milestones", arg0)

	outstruct := new(struct {
		Agent              common.Address
		Token              common.Address
		StageAmount        *big.Int
		LastStageRemainder *big.Int
		TotalStages        uint8
		CompletedStages    uint8
		Initialised        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Agent = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Token = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.StageAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LastStageRemainder = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.TotalStages = *abi.ConvertType(out[4], new(uint8)).(*uint8)
	outstruct.CompletedStages = *abi.ConvertType(out[5], new(uint8)).(*uint8)
	outstruct.Initialised = *abi.ConvertType(out[6], new(bool)).(*bool)

	return *outstruct, err

}

// Milestones is a free data retrieval call binding the contract method 0xe89e4ed6.
//
// Solidity: function milestones(uint256 ) view returns(address agent, address token, uint256 stageAmount, uint256 lastStageRemainder, uint8 totalStages, uint8 completedStages, bool initialised)
func (_MilestoneHook *MilestoneHookSession) Milestones(arg0 *big.Int) (struct {
	Agent              common.Address
	Token              common.Address
	StageAmount        *big.Int
	LastStageRemainder *big.Int
	TotalStages        uint8
	CompletedStages    uint8
	Initialised        bool
}, error) {
	return _MilestoneHook.Contract.Milestones(&_MilestoneHook.CallOpts, arg0)
}

// Milestones is a free data retrieval call binding the contract method 0xe89e4ed6.
//
// Solidity: function milestones(uint256 ) view returns(address agent, address token, uint256 stageAmount, uint256 lastStageRemainder, uint8 totalStages, uint8 completedStages, bool initialised)
func (_MilestoneHook *MilestoneHookCallerSession) Milestones(arg0 *big.Int) (struct {
	Agent              common.Address
	Token              common.Address
	StageAmount        *big.Int
	LastStageRemainder *big.Int
	TotalStages        uint8
	CompletedStages    uint8
	Initialised        bool
}, error) {
	return _MilestoneHook.Contract.Milestones(&_MilestoneHook.CallOpts, arg0)
}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_MilestoneHook *MilestoneHookCaller) OnReject(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _MilestoneHook.contract.Call(opts, &out, "onReject", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_MilestoneHook *MilestoneHookSession) OnReject(arg0 *big.Int, arg1 []byte) error {
	return _MilestoneHook.Contract.OnReject(&_MilestoneHook.CallOpts, arg0, arg1)
}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_MilestoneHook *MilestoneHookCallerSession) OnReject(arg0 *big.Int, arg1 []byte) error {
	return _MilestoneHook.Contract.OnReject(&_MilestoneHook.CallOpts, arg0, arg1)
}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_MilestoneHook *MilestoneHookCaller) OnSubmit(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _MilestoneHook.contract.Call(opts, &out, "onSubmit", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_MilestoneHook *MilestoneHookSession) OnSubmit(arg0 *big.Int, arg1 []byte) error {
	return _MilestoneHook.Contract.OnSubmit(&_MilestoneHook.CallOpts, arg0, arg1)
}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_MilestoneHook *MilestoneHookCallerSession) OnSubmit(arg0 *big.Int, arg1 []byte) error {
	return _MilestoneHook.Contract.OnSubmit(&_MilestoneHook.CallOpts, arg0, arg1)
}

// OnAccept is a paid mutator transaction binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 jobId, bytes context) returns()
func (_MilestoneHook *MilestoneHookTransactor) OnAccept(opts *bind.TransactOpts, jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _MilestoneHook.contract.Transact(opts, "onAccept", jobId, context)
}

// OnAccept is a paid mutator transaction binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 jobId, bytes context) returns()
func (_MilestoneHook *MilestoneHookSession) OnAccept(jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _MilestoneHook.Contract.OnAccept(&_MilestoneHook.TransactOpts, jobId, context)
}

// OnAccept is a paid mutator transaction binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 jobId, bytes context) returns()
func (_MilestoneHook *MilestoneHookTransactorSession) OnAccept(jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _MilestoneHook.Contract.OnAccept(&_MilestoneHook.TransactOpts, jobId, context)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes ) returns()
func (_MilestoneHook *MilestoneHookTransactor) OnApprove(opts *bind.TransactOpts, jobId *big.Int, arg1 []byte) (*types.Transaction, error) {
	return _MilestoneHook.contract.Transact(opts, "onApprove", jobId, arg1)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes ) returns()
func (_MilestoneHook *MilestoneHookSession) OnApprove(jobId *big.Int, arg1 []byte) (*types.Transaction, error) {
	return _MilestoneHook.Contract.OnApprove(&_MilestoneHook.TransactOpts, jobId, arg1)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes ) returns()
func (_MilestoneHook *MilestoneHookTransactorSession) OnApprove(jobId *big.Int, arg1 []byte) (*types.Transaction, error) {
	return _MilestoneHook.Contract.OnApprove(&_MilestoneHook.TransactOpts, jobId, arg1)
}

// OnCancel is a paid mutator transaction binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 jobId, bytes ) returns()
func (_MilestoneHook *MilestoneHookTransactor) OnCancel(opts *bind.TransactOpts, jobId *big.Int, arg1 []byte) (*types.Transaction, error) {
	return _MilestoneHook.contract.Transact(opts, "onCancel", jobId, arg1)
}

// OnCancel is a paid mutator transaction binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 jobId, bytes ) returns()
func (_MilestoneHook *MilestoneHookSession) OnCancel(jobId *big.Int, arg1 []byte) (*types.Transaction, error) {
	return _MilestoneHook.Contract.OnCancel(&_MilestoneHook.TransactOpts, jobId, arg1)
}

// OnCancel is a paid mutator transaction binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 jobId, bytes ) returns()
func (_MilestoneHook *MilestoneHookTransactorSession) OnCancel(jobId *big.Int, arg1 []byte) (*types.Transaction, error) {
	return _MilestoneHook.Contract.OnCancel(&_MilestoneHook.TransactOpts, jobId, arg1)
}

// MilestoneHookMilestoneCancelledIterator is returned from FilterMilestoneCancelled and is used to iterate over the raw logs and unpacked data for MilestoneCancelled events raised by the MilestoneHook contract.
type MilestoneHookMilestoneCancelledIterator struct {
	Event *MilestoneHookMilestoneCancelled // Event containing the contract specifics and raw log

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
func (it *MilestoneHookMilestoneCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MilestoneHookMilestoneCancelled)
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
		it.Event = new(MilestoneHookMilestoneCancelled)
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
func (it *MilestoneHookMilestoneCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MilestoneHookMilestoneCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MilestoneHookMilestoneCancelled represents a MilestoneCancelled event raised by the MilestoneHook contract.
type MilestoneHookMilestoneCancelled struct {
	JobId           *big.Int
	RemainingStages uint8
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterMilestoneCancelled is a free log retrieval operation binding the contract event 0x1a0b394b2791bb8cb3c7490e14c4a3617f85d3c9b0bcf966fc8154c0618cfaef.
//
// Solidity: event MilestoneCancelled(uint256 indexed jobId, uint8 remainingStages)
func (_MilestoneHook *MilestoneHookFilterer) FilterMilestoneCancelled(opts *bind.FilterOpts, jobId []*big.Int) (*MilestoneHookMilestoneCancelledIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}

	logs, sub, err := _MilestoneHook.contract.FilterLogs(opts, "MilestoneCancelled", jobIdRule)
	if err != nil {
		return nil, err
	}
	return &MilestoneHookMilestoneCancelledIterator{contract: _MilestoneHook.contract, event: "MilestoneCancelled", logs: logs, sub: sub}, nil
}

// WatchMilestoneCancelled is a free log subscription operation binding the contract event 0x1a0b394b2791bb8cb3c7490e14c4a3617f85d3c9b0bcf966fc8154c0618cfaef.
//
// Solidity: event MilestoneCancelled(uint256 indexed jobId, uint8 remainingStages)
func (_MilestoneHook *MilestoneHookFilterer) WatchMilestoneCancelled(opts *bind.WatchOpts, sink chan<- *MilestoneHookMilestoneCancelled, jobId []*big.Int) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}

	logs, sub, err := _MilestoneHook.contract.WatchLogs(opts, "MilestoneCancelled", jobIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MilestoneHookMilestoneCancelled)
				if err := _MilestoneHook.contract.UnpackLog(event, "MilestoneCancelled", log); err != nil {
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

// ParseMilestoneCancelled is a log parse operation binding the contract event 0x1a0b394b2791bb8cb3c7490e14c4a3617f85d3c9b0bcf966fc8154c0618cfaef.
//
// Solidity: event MilestoneCancelled(uint256 indexed jobId, uint8 remainingStages)
func (_MilestoneHook *MilestoneHookFilterer) ParseMilestoneCancelled(log types.Log) (*MilestoneHookMilestoneCancelled, error) {
	event := new(MilestoneHookMilestoneCancelled)
	if err := _MilestoneHook.contract.UnpackLog(event, "MilestoneCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MilestoneHookMilestoneReleasedIterator is returned from FilterMilestoneReleased and is used to iterate over the raw logs and unpacked data for MilestoneReleased events raised by the MilestoneHook contract.
type MilestoneHookMilestoneReleasedIterator struct {
	Event *MilestoneHookMilestoneReleased // Event containing the contract specifics and raw log

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
func (it *MilestoneHookMilestoneReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MilestoneHookMilestoneReleased)
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
		it.Event = new(MilestoneHookMilestoneReleased)
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
func (it *MilestoneHookMilestoneReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MilestoneHookMilestoneReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MilestoneHookMilestoneReleased represents a MilestoneReleased event raised by the MilestoneHook contract.
type MilestoneHookMilestoneReleased struct {
	JobId       *big.Int
	Agent       common.Address
	Stage       uint8
	TotalStages uint8
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMilestoneReleased is a free log retrieval operation binding the contract event 0xd25357de9f9b59472a26c94ce8935500c7af97338a72f82d8c0d5ea5dde3f27d.
//
// Solidity: event MilestoneReleased(uint256 indexed jobId, address indexed agent, uint8 stage, uint8 totalStages, uint256 amount)
func (_MilestoneHook *MilestoneHookFilterer) FilterMilestoneReleased(opts *bind.FilterOpts, jobId []*big.Int, agent []common.Address) (*MilestoneHookMilestoneReleasedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _MilestoneHook.contract.FilterLogs(opts, "MilestoneReleased", jobIdRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &MilestoneHookMilestoneReleasedIterator{contract: _MilestoneHook.contract, event: "MilestoneReleased", logs: logs, sub: sub}, nil
}

// WatchMilestoneReleased is a free log subscription operation binding the contract event 0xd25357de9f9b59472a26c94ce8935500c7af97338a72f82d8c0d5ea5dde3f27d.
//
// Solidity: event MilestoneReleased(uint256 indexed jobId, address indexed agent, uint8 stage, uint8 totalStages, uint256 amount)
func (_MilestoneHook *MilestoneHookFilterer) WatchMilestoneReleased(opts *bind.WatchOpts, sink chan<- *MilestoneHookMilestoneReleased, jobId []*big.Int, agent []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _MilestoneHook.contract.WatchLogs(opts, "MilestoneReleased", jobIdRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MilestoneHookMilestoneReleased)
				if err := _MilestoneHook.contract.UnpackLog(event, "MilestoneReleased", log); err != nil {
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

// ParseMilestoneReleased is a log parse operation binding the contract event 0xd25357de9f9b59472a26c94ce8935500c7af97338a72f82d8c0d5ea5dde3f27d.
//
// Solidity: event MilestoneReleased(uint256 indexed jobId, address indexed agent, uint8 stage, uint8 totalStages, uint256 amount)
func (_MilestoneHook *MilestoneHookFilterer) ParseMilestoneReleased(log types.Log) (*MilestoneHookMilestoneReleased, error) {
	event := new(MilestoneHookMilestoneReleased)
	if err := _MilestoneHook.contract.UnpackLog(event, "MilestoneReleased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
