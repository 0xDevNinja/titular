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

// FundTransferHookMetaData contains all meta data concerning the FundTransferHook contract.
var FundTransferHookMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hookName\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onAccept\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onApprove\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onCancel\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onReject\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onSubmit\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"FundTransferred\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agentAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"pnlAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"settlementAddr\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InvalidBps\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// FundTransferHookABI is the input ABI used to generate the binding from.
// Deprecated: Use FundTransferHookMetaData.ABI instead.
var FundTransferHookABI = FundTransferHookMetaData.ABI

// FundTransferHook is an auto generated Go binding around an Ethereum contract.
type FundTransferHook struct {
	FundTransferHookCaller     // Read-only binding to the contract
	FundTransferHookTransactor // Write-only binding to the contract
	FundTransferHookFilterer   // Log filterer for contract events
}

// FundTransferHookCaller is an auto generated read-only Go binding around an Ethereum contract.
type FundTransferHookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FundTransferHookTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FundTransferHookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FundTransferHookFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FundTransferHookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FundTransferHookSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FundTransferHookSession struct {
	Contract     *FundTransferHook // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FundTransferHookCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FundTransferHookCallerSession struct {
	Contract *FundTransferHookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// FundTransferHookTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FundTransferHookTransactorSession struct {
	Contract     *FundTransferHookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// FundTransferHookRaw is an auto generated low-level Go binding around an Ethereum contract.
type FundTransferHookRaw struct {
	Contract *FundTransferHook // Generic contract binding to access the raw methods on
}

// FundTransferHookCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FundTransferHookCallerRaw struct {
	Contract *FundTransferHookCaller // Generic read-only contract binding to access the raw methods on
}

// FundTransferHookTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FundTransferHookTransactorRaw struct {
	Contract *FundTransferHookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFundTransferHook creates a new instance of FundTransferHook, bound to a specific deployed contract.
func NewFundTransferHook(address common.Address, backend bind.ContractBackend) (*FundTransferHook, error) {
	contract, err := bindFundTransferHook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FundTransferHook{FundTransferHookCaller: FundTransferHookCaller{contract: contract}, FundTransferHookTransactor: FundTransferHookTransactor{contract: contract}, FundTransferHookFilterer: FundTransferHookFilterer{contract: contract}}, nil
}

// NewFundTransferHookCaller creates a new read-only instance of FundTransferHook, bound to a specific deployed contract.
func NewFundTransferHookCaller(address common.Address, caller bind.ContractCaller) (*FundTransferHookCaller, error) {
	contract, err := bindFundTransferHook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FundTransferHookCaller{contract: contract}, nil
}

// NewFundTransferHookTransactor creates a new write-only instance of FundTransferHook, bound to a specific deployed contract.
func NewFundTransferHookTransactor(address common.Address, transactor bind.ContractTransactor) (*FundTransferHookTransactor, error) {
	contract, err := bindFundTransferHook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FundTransferHookTransactor{contract: contract}, nil
}

// NewFundTransferHookFilterer creates a new log filterer instance of FundTransferHook, bound to a specific deployed contract.
func NewFundTransferHookFilterer(address common.Address, filterer bind.ContractFilterer) (*FundTransferHookFilterer, error) {
	contract, err := bindFundTransferHook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FundTransferHookFilterer{contract: contract}, nil
}

// bindFundTransferHook binds a generic wrapper to an already deployed contract.
func bindFundTransferHook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FundTransferHookMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FundTransferHook *FundTransferHookRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FundTransferHook.Contract.FundTransferHookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FundTransferHook *FundTransferHookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FundTransferHook.Contract.FundTransferHookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FundTransferHook *FundTransferHookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FundTransferHook.Contract.FundTransferHookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FundTransferHook *FundTransferHookCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FundTransferHook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FundTransferHook *FundTransferHookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FundTransferHook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FundTransferHook *FundTransferHookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FundTransferHook.Contract.contract.Transact(opts, method, params...)
}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_FundTransferHook *FundTransferHookCaller) BPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FundTransferHook.contract.Call(opts, &out, "BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_FundTransferHook *FundTransferHookSession) BPS() (*big.Int, error) {
	return _FundTransferHook.Contract.BPS(&_FundTransferHook.CallOpts)
}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_FundTransferHook *FundTransferHookCallerSession) BPS() (*big.Int, error) {
	return _FundTransferHook.Contract.BPS(&_FundTransferHook.CallOpts)
}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_FundTransferHook *FundTransferHookCaller) HookName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FundTransferHook.contract.Call(opts, &out, "hookName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_FundTransferHook *FundTransferHookSession) HookName() (string, error) {
	return _FundTransferHook.Contract.HookName(&_FundTransferHook.CallOpts)
}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_FundTransferHook *FundTransferHookCallerSession) HookName() (string, error) {
	return _FundTransferHook.Contract.HookName(&_FundTransferHook.CallOpts)
}

// OnAccept is a free data retrieval call binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookCaller) OnAccept(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _FundTransferHook.contract.Call(opts, &out, "onAccept", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnAccept is a free data retrieval call binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookSession) OnAccept(arg0 *big.Int, arg1 []byte) error {
	return _FundTransferHook.Contract.OnAccept(&_FundTransferHook.CallOpts, arg0, arg1)
}

// OnAccept is a free data retrieval call binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookCallerSession) OnAccept(arg0 *big.Int, arg1 []byte) error {
	return _FundTransferHook.Contract.OnAccept(&_FundTransferHook.CallOpts, arg0, arg1)
}

// OnCancel is a free data retrieval call binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookCaller) OnCancel(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _FundTransferHook.contract.Call(opts, &out, "onCancel", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnCancel is a free data retrieval call binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookSession) OnCancel(arg0 *big.Int, arg1 []byte) error {
	return _FundTransferHook.Contract.OnCancel(&_FundTransferHook.CallOpts, arg0, arg1)
}

// OnCancel is a free data retrieval call binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookCallerSession) OnCancel(arg0 *big.Int, arg1 []byte) error {
	return _FundTransferHook.Contract.OnCancel(&_FundTransferHook.CallOpts, arg0, arg1)
}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookCaller) OnReject(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _FundTransferHook.contract.Call(opts, &out, "onReject", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookSession) OnReject(arg0 *big.Int, arg1 []byte) error {
	return _FundTransferHook.Contract.OnReject(&_FundTransferHook.CallOpts, arg0, arg1)
}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookCallerSession) OnReject(arg0 *big.Int, arg1 []byte) error {
	return _FundTransferHook.Contract.OnReject(&_FundTransferHook.CallOpts, arg0, arg1)
}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookCaller) OnSubmit(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _FundTransferHook.contract.Call(opts, &out, "onSubmit", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookSession) OnSubmit(arg0 *big.Int, arg1 []byte) error {
	return _FundTransferHook.Contract.OnSubmit(&_FundTransferHook.CallOpts, arg0, arg1)
}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_FundTransferHook *FundTransferHookCallerSession) OnSubmit(arg0 *big.Int, arg1 []byte) error {
	return _FundTransferHook.Contract.OnSubmit(&_FundTransferHook.CallOpts, arg0, arg1)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes context) returns()
func (_FundTransferHook *FundTransferHookTransactor) OnApprove(opts *bind.TransactOpts, jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _FundTransferHook.contract.Transact(opts, "onApprove", jobId, context)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes context) returns()
func (_FundTransferHook *FundTransferHookSession) OnApprove(jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _FundTransferHook.Contract.OnApprove(&_FundTransferHook.TransactOpts, jobId, context)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes context) returns()
func (_FundTransferHook *FundTransferHookTransactorSession) OnApprove(jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _FundTransferHook.Contract.OnApprove(&_FundTransferHook.TransactOpts, jobId, context)
}

// FundTransferHookFundTransferredIterator is returned from FilterFundTransferred and is used to iterate over the raw logs and unpacked data for FundTransferred events raised by the FundTransferHook contract.
type FundTransferHookFundTransferredIterator struct {
	Event *FundTransferHookFundTransferred // Event containing the contract specifics and raw log

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
func (it *FundTransferHookFundTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FundTransferHookFundTransferred)
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
		it.Event = new(FundTransferHookFundTransferred)
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
func (it *FundTransferHookFundTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FundTransferHookFundTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FundTransferHookFundTransferred represents a FundTransferred event raised by the FundTransferHook contract.
type FundTransferHookFundTransferred struct {
	JobId          *big.Int
	Agent          common.Address
	Token          common.Address
	AgentAmount    *big.Int
	PnlAmount      *big.Int
	SettlementAddr common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterFundTransferred is a free log retrieval operation binding the contract event 0xa33cd98019ad28fa37fbbc444267ba860abc1c449e0f9bdad074d09cdae5ebb1.
//
// Solidity: event FundTransferred(uint256 indexed jobId, address indexed agent, address indexed token, uint256 agentAmount, uint256 pnlAmount, address settlementAddr)
func (_FundTransferHook *FundTransferHookFilterer) FilterFundTransferred(opts *bind.FilterOpts, jobId []*big.Int, agent []common.Address, token []common.Address) (*FundTransferHookFundTransferredIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FundTransferHook.contract.FilterLogs(opts, "FundTransferred", jobIdRule, agentRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &FundTransferHookFundTransferredIterator{contract: _FundTransferHook.contract, event: "FundTransferred", logs: logs, sub: sub}, nil
}

// WatchFundTransferred is a free log subscription operation binding the contract event 0xa33cd98019ad28fa37fbbc444267ba860abc1c449e0f9bdad074d09cdae5ebb1.
//
// Solidity: event FundTransferred(uint256 indexed jobId, address indexed agent, address indexed token, uint256 agentAmount, uint256 pnlAmount, address settlementAddr)
func (_FundTransferHook *FundTransferHookFilterer) WatchFundTransferred(opts *bind.WatchOpts, sink chan<- *FundTransferHookFundTransferred, jobId []*big.Int, agent []common.Address, token []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FundTransferHook.contract.WatchLogs(opts, "FundTransferred", jobIdRule, agentRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FundTransferHookFundTransferred)
				if err := _FundTransferHook.contract.UnpackLog(event, "FundTransferred", log); err != nil {
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

// ParseFundTransferred is a log parse operation binding the contract event 0xa33cd98019ad28fa37fbbc444267ba860abc1c449e0f9bdad074d09cdae5ebb1.
//
// Solidity: event FundTransferred(uint256 indexed jobId, address indexed agent, address indexed token, uint256 agentAmount, uint256 pnlAmount, address settlementAddr)
func (_FundTransferHook *FundTransferHookFilterer) ParseFundTransferred(log types.Log) (*FundTransferHookFundTransferred, error) {
	event := new(FundTransferHookFundTransferred)
	if err := _FundTransferHook.contract.UnpackLog(event, "FundTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
