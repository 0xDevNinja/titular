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

// SubscriptionHookMetaData contains all meta data concerning the SubscriptionHook contract.
var SubscriptionHookMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"hookName\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onAccept\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onApprove\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onCancel\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onReject\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onSubmit\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"subscriptions\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"nextRenewalAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"SubscriptionCancelled\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SubscriptionRenewed\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"subscriber\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"nextRenewalAt\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidDuration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SubscriptionNotActive\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// SubscriptionHookABI is the input ABI used to generate the binding from.
// Deprecated: Use SubscriptionHookMetaData.ABI instead.
var SubscriptionHookABI = SubscriptionHookMetaData.ABI

// SubscriptionHook is an auto generated Go binding around an Ethereum contract.
type SubscriptionHook struct {
	SubscriptionHookCaller     // Read-only binding to the contract
	SubscriptionHookTransactor // Write-only binding to the contract
	SubscriptionHookFilterer   // Log filterer for contract events
}

// SubscriptionHookCaller is an auto generated read-only Go binding around an Ethereum contract.
type SubscriptionHookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubscriptionHookTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SubscriptionHookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubscriptionHookFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SubscriptionHookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SubscriptionHookSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SubscriptionHookSession struct {
	Contract     *SubscriptionHook // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SubscriptionHookCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SubscriptionHookCallerSession struct {
	Contract *SubscriptionHookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// SubscriptionHookTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SubscriptionHookTransactorSession struct {
	Contract     *SubscriptionHookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// SubscriptionHookRaw is an auto generated low-level Go binding around an Ethereum contract.
type SubscriptionHookRaw struct {
	Contract *SubscriptionHook // Generic contract binding to access the raw methods on
}

// SubscriptionHookCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SubscriptionHookCallerRaw struct {
	Contract *SubscriptionHookCaller // Generic read-only contract binding to access the raw methods on
}

// SubscriptionHookTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SubscriptionHookTransactorRaw struct {
	Contract *SubscriptionHookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSubscriptionHook creates a new instance of SubscriptionHook, bound to a specific deployed contract.
func NewSubscriptionHook(address common.Address, backend bind.ContractBackend) (*SubscriptionHook, error) {
	contract, err := bindSubscriptionHook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SubscriptionHook{SubscriptionHookCaller: SubscriptionHookCaller{contract: contract}, SubscriptionHookTransactor: SubscriptionHookTransactor{contract: contract}, SubscriptionHookFilterer: SubscriptionHookFilterer{contract: contract}}, nil
}

// NewSubscriptionHookCaller creates a new read-only instance of SubscriptionHook, bound to a specific deployed contract.
func NewSubscriptionHookCaller(address common.Address, caller bind.ContractCaller) (*SubscriptionHookCaller, error) {
	contract, err := bindSubscriptionHook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SubscriptionHookCaller{contract: contract}, nil
}

// NewSubscriptionHookTransactor creates a new write-only instance of SubscriptionHook, bound to a specific deployed contract.
func NewSubscriptionHookTransactor(address common.Address, transactor bind.ContractTransactor) (*SubscriptionHookTransactor, error) {
	contract, err := bindSubscriptionHook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SubscriptionHookTransactor{contract: contract}, nil
}

// NewSubscriptionHookFilterer creates a new log filterer instance of SubscriptionHook, bound to a specific deployed contract.
func NewSubscriptionHookFilterer(address common.Address, filterer bind.ContractFilterer) (*SubscriptionHookFilterer, error) {
	contract, err := bindSubscriptionHook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SubscriptionHookFilterer{contract: contract}, nil
}

// bindSubscriptionHook binds a generic wrapper to an already deployed contract.
func bindSubscriptionHook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SubscriptionHookMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubscriptionHook *SubscriptionHookRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubscriptionHook.Contract.SubscriptionHookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubscriptionHook *SubscriptionHookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubscriptionHook.Contract.SubscriptionHookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubscriptionHook *SubscriptionHookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubscriptionHook.Contract.SubscriptionHookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SubscriptionHook *SubscriptionHookCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SubscriptionHook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SubscriptionHook *SubscriptionHookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SubscriptionHook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SubscriptionHook *SubscriptionHookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SubscriptionHook.Contract.contract.Transact(opts, method, params...)
}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_SubscriptionHook *SubscriptionHookCaller) HookName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SubscriptionHook.contract.Call(opts, &out, "hookName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_SubscriptionHook *SubscriptionHookSession) HookName() (string, error) {
	return _SubscriptionHook.Contract.HookName(&_SubscriptionHook.CallOpts)
}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_SubscriptionHook *SubscriptionHookCallerSession) HookName() (string, error) {
	return _SubscriptionHook.Contract.HookName(&_SubscriptionHook.CallOpts)
}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_SubscriptionHook *SubscriptionHookCaller) OnReject(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _SubscriptionHook.contract.Call(opts, &out, "onReject", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_SubscriptionHook *SubscriptionHookSession) OnReject(arg0 *big.Int, arg1 []byte) error {
	return _SubscriptionHook.Contract.OnReject(&_SubscriptionHook.CallOpts, arg0, arg1)
}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_SubscriptionHook *SubscriptionHookCallerSession) OnReject(arg0 *big.Int, arg1 []byte) error {
	return _SubscriptionHook.Contract.OnReject(&_SubscriptionHook.CallOpts, arg0, arg1)
}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_SubscriptionHook *SubscriptionHookCaller) OnSubmit(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _SubscriptionHook.contract.Call(opts, &out, "onSubmit", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_SubscriptionHook *SubscriptionHookSession) OnSubmit(arg0 *big.Int, arg1 []byte) error {
	return _SubscriptionHook.Contract.OnSubmit(&_SubscriptionHook.CallOpts, arg0, arg1)
}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_SubscriptionHook *SubscriptionHookCallerSession) OnSubmit(arg0 *big.Int, arg1 []byte) error {
	return _SubscriptionHook.Contract.OnSubmit(&_SubscriptionHook.CallOpts, arg0, arg1)
}

// Subscriptions is a free data retrieval call binding the contract method 0x2d5bbf60.
//
// Solidity: function subscriptions(uint256 ) view returns(uint64 nextRenewalAt, bool active)
func (_SubscriptionHook *SubscriptionHookCaller) Subscriptions(opts *bind.CallOpts, arg0 *big.Int) (struct {
	NextRenewalAt uint64
	Active        bool
}, error) {
	var out []interface{}
	err := _SubscriptionHook.contract.Call(opts, &out, "subscriptions", arg0)

	outstruct := new(struct {
		NextRenewalAt uint64
		Active        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NextRenewalAt = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.Active = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// Subscriptions is a free data retrieval call binding the contract method 0x2d5bbf60.
//
// Solidity: function subscriptions(uint256 ) view returns(uint64 nextRenewalAt, bool active)
func (_SubscriptionHook *SubscriptionHookSession) Subscriptions(arg0 *big.Int) (struct {
	NextRenewalAt uint64
	Active        bool
}, error) {
	return _SubscriptionHook.Contract.Subscriptions(&_SubscriptionHook.CallOpts, arg0)
}

// Subscriptions is a free data retrieval call binding the contract method 0x2d5bbf60.
//
// Solidity: function subscriptions(uint256 ) view returns(uint64 nextRenewalAt, bool active)
func (_SubscriptionHook *SubscriptionHookCallerSession) Subscriptions(arg0 *big.Int) (struct {
	NextRenewalAt uint64
	Active        bool
}, error) {
	return _SubscriptionHook.Contract.Subscriptions(&_SubscriptionHook.CallOpts, arg0)
}

// OnAccept is a paid mutator transaction binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 jobId, bytes context) returns()
func (_SubscriptionHook *SubscriptionHookTransactor) OnAccept(opts *bind.TransactOpts, jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _SubscriptionHook.contract.Transact(opts, "onAccept", jobId, context)
}

// OnAccept is a paid mutator transaction binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 jobId, bytes context) returns()
func (_SubscriptionHook *SubscriptionHookSession) OnAccept(jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _SubscriptionHook.Contract.OnAccept(&_SubscriptionHook.TransactOpts, jobId, context)
}

// OnAccept is a paid mutator transaction binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 jobId, bytes context) returns()
func (_SubscriptionHook *SubscriptionHookTransactorSession) OnAccept(jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _SubscriptionHook.Contract.OnAccept(&_SubscriptionHook.TransactOpts, jobId, context)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes context) returns()
func (_SubscriptionHook *SubscriptionHookTransactor) OnApprove(opts *bind.TransactOpts, jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _SubscriptionHook.contract.Transact(opts, "onApprove", jobId, context)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes context) returns()
func (_SubscriptionHook *SubscriptionHookSession) OnApprove(jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _SubscriptionHook.Contract.OnApprove(&_SubscriptionHook.TransactOpts, jobId, context)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes context) returns()
func (_SubscriptionHook *SubscriptionHookTransactorSession) OnApprove(jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _SubscriptionHook.Contract.OnApprove(&_SubscriptionHook.TransactOpts, jobId, context)
}

// OnCancel is a paid mutator transaction binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 jobId, bytes ) returns()
func (_SubscriptionHook *SubscriptionHookTransactor) OnCancel(opts *bind.TransactOpts, jobId *big.Int, arg1 []byte) (*types.Transaction, error) {
	return _SubscriptionHook.contract.Transact(opts, "onCancel", jobId, arg1)
}

// OnCancel is a paid mutator transaction binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 jobId, bytes ) returns()
func (_SubscriptionHook *SubscriptionHookSession) OnCancel(jobId *big.Int, arg1 []byte) (*types.Transaction, error) {
	return _SubscriptionHook.Contract.OnCancel(&_SubscriptionHook.TransactOpts, jobId, arg1)
}

// OnCancel is a paid mutator transaction binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 jobId, bytes ) returns()
func (_SubscriptionHook *SubscriptionHookTransactorSession) OnCancel(jobId *big.Int, arg1 []byte) (*types.Transaction, error) {
	return _SubscriptionHook.Contract.OnCancel(&_SubscriptionHook.TransactOpts, jobId, arg1)
}

// SubscriptionHookSubscriptionCancelledIterator is returned from FilterSubscriptionCancelled and is used to iterate over the raw logs and unpacked data for SubscriptionCancelled events raised by the SubscriptionHook contract.
type SubscriptionHookSubscriptionCancelledIterator struct {
	Event *SubscriptionHookSubscriptionCancelled // Event containing the contract specifics and raw log

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
func (it *SubscriptionHookSubscriptionCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionHookSubscriptionCancelled)
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
		it.Event = new(SubscriptionHookSubscriptionCancelled)
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
func (it *SubscriptionHookSubscriptionCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubscriptionHookSubscriptionCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubscriptionHookSubscriptionCancelled represents a SubscriptionCancelled event raised by the SubscriptionHook contract.
type SubscriptionHookSubscriptionCancelled struct {
	JobId *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionCancelled is a free log retrieval operation binding the contract event 0xbd2bcea75d16a85f005cd83447e0de57341bf926fe7419e6d553663e91ab4da7.
//
// Solidity: event SubscriptionCancelled(uint256 indexed jobId)
func (_SubscriptionHook *SubscriptionHookFilterer) FilterSubscriptionCancelled(opts *bind.FilterOpts, jobId []*big.Int) (*SubscriptionHookSubscriptionCancelledIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}

	logs, sub, err := _SubscriptionHook.contract.FilterLogs(opts, "SubscriptionCancelled", jobIdRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionHookSubscriptionCancelledIterator{contract: _SubscriptionHook.contract, event: "SubscriptionCancelled", logs: logs, sub: sub}, nil
}

// WatchSubscriptionCancelled is a free log subscription operation binding the contract event 0xbd2bcea75d16a85f005cd83447e0de57341bf926fe7419e6d553663e91ab4da7.
//
// Solidity: event SubscriptionCancelled(uint256 indexed jobId)
func (_SubscriptionHook *SubscriptionHookFilterer) WatchSubscriptionCancelled(opts *bind.WatchOpts, sink chan<- *SubscriptionHookSubscriptionCancelled, jobId []*big.Int) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}

	logs, sub, err := _SubscriptionHook.contract.WatchLogs(opts, "SubscriptionCancelled", jobIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubscriptionHookSubscriptionCancelled)
				if err := _SubscriptionHook.contract.UnpackLog(event, "SubscriptionCancelled", log); err != nil {
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

// ParseSubscriptionCancelled is a log parse operation binding the contract event 0xbd2bcea75d16a85f005cd83447e0de57341bf926fe7419e6d553663e91ab4da7.
//
// Solidity: event SubscriptionCancelled(uint256 indexed jobId)
func (_SubscriptionHook *SubscriptionHookFilterer) ParseSubscriptionCancelled(log types.Log) (*SubscriptionHookSubscriptionCancelled, error) {
	event := new(SubscriptionHookSubscriptionCancelled)
	if err := _SubscriptionHook.contract.UnpackLog(event, "SubscriptionCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SubscriptionHookSubscriptionRenewedIterator is returned from FilterSubscriptionRenewed and is used to iterate over the raw logs and unpacked data for SubscriptionRenewed events raised by the SubscriptionHook contract.
type SubscriptionHookSubscriptionRenewedIterator struct {
	Event *SubscriptionHookSubscriptionRenewed // Event containing the contract specifics and raw log

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
func (it *SubscriptionHookSubscriptionRenewedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SubscriptionHookSubscriptionRenewed)
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
		it.Event = new(SubscriptionHookSubscriptionRenewed)
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
func (it *SubscriptionHookSubscriptionRenewedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SubscriptionHookSubscriptionRenewedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SubscriptionHookSubscriptionRenewed represents a SubscriptionRenewed event raised by the SubscriptionHook contract.
type SubscriptionHookSubscriptionRenewed struct {
	JobId         *big.Int
	Subscriber    common.Address
	Provider      common.Address
	Amount        *big.Int
	NextRenewalAt uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSubscriptionRenewed is a free log retrieval operation binding the contract event 0x4a40031fc0683564de624c4cee39202baa4d6ebe5d29b6c82b11e5b6cf185ef3.
//
// Solidity: event SubscriptionRenewed(uint256 indexed jobId, address indexed subscriber, address indexed provider, uint256 amount, uint64 nextRenewalAt)
func (_SubscriptionHook *SubscriptionHookFilterer) FilterSubscriptionRenewed(opts *bind.FilterOpts, jobId []*big.Int, subscriber []common.Address, provider []common.Address) (*SubscriptionHookSubscriptionRenewedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _SubscriptionHook.contract.FilterLogs(opts, "SubscriptionRenewed", jobIdRule, subscriberRule, providerRule)
	if err != nil {
		return nil, err
	}
	return &SubscriptionHookSubscriptionRenewedIterator{contract: _SubscriptionHook.contract, event: "SubscriptionRenewed", logs: logs, sub: sub}, nil
}

// WatchSubscriptionRenewed is a free log subscription operation binding the contract event 0x4a40031fc0683564de624c4cee39202baa4d6ebe5d29b6c82b11e5b6cf185ef3.
//
// Solidity: event SubscriptionRenewed(uint256 indexed jobId, address indexed subscriber, address indexed provider, uint256 amount, uint64 nextRenewalAt)
func (_SubscriptionHook *SubscriptionHookFilterer) WatchSubscriptionRenewed(opts *bind.WatchOpts, sink chan<- *SubscriptionHookSubscriptionRenewed, jobId []*big.Int, subscriber []common.Address, provider []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var subscriberRule []interface{}
	for _, subscriberItem := range subscriber {
		subscriberRule = append(subscriberRule, subscriberItem)
	}
	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _SubscriptionHook.contract.WatchLogs(opts, "SubscriptionRenewed", jobIdRule, subscriberRule, providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SubscriptionHookSubscriptionRenewed)
				if err := _SubscriptionHook.contract.UnpackLog(event, "SubscriptionRenewed", log); err != nil {
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

// ParseSubscriptionRenewed is a log parse operation binding the contract event 0x4a40031fc0683564de624c4cee39202baa4d6ebe5d29b6c82b11e5b6cf185ef3.
//
// Solidity: event SubscriptionRenewed(uint256 indexed jobId, address indexed subscriber, address indexed provider, uint256 amount, uint64 nextRenewalAt)
func (_SubscriptionHook *SubscriptionHookFilterer) ParseSubscriptionRenewed(log types.Log) (*SubscriptionHookSubscriptionRenewed, error) {
	event := new(SubscriptionHookSubscriptionRenewed)
	if err := _SubscriptionHook.contract.UnpackLog(event, "SubscriptionRenewed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
