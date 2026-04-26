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

// RoyaltyHookMetaData contains all meta data concerning the RoyaltyHook contract.
var RoyaltyHookMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hookName\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onAccept\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onApprove\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"context\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"onCancel\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onReject\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"onSubmit\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"RoyaltyPaid\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"NoRecipients\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RecipientZeroAddress\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ShareSumMismatch\",\"inputs\":[{\"name\":\"sum\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// RoyaltyHookABI is the input ABI used to generate the binding from.
// Deprecated: Use RoyaltyHookMetaData.ABI instead.
var RoyaltyHookABI = RoyaltyHookMetaData.ABI

// RoyaltyHook is an auto generated Go binding around an Ethereum contract.
type RoyaltyHook struct {
	RoyaltyHookCaller     // Read-only binding to the contract
	RoyaltyHookTransactor // Write-only binding to the contract
	RoyaltyHookFilterer   // Log filterer for contract events
}

// RoyaltyHookCaller is an auto generated read-only Go binding around an Ethereum contract.
type RoyaltyHookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoyaltyHookTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RoyaltyHookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoyaltyHookFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RoyaltyHookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RoyaltyHookSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RoyaltyHookSession struct {
	Contract     *RoyaltyHook      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RoyaltyHookCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RoyaltyHookCallerSession struct {
	Contract *RoyaltyHookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// RoyaltyHookTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RoyaltyHookTransactorSession struct {
	Contract     *RoyaltyHookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// RoyaltyHookRaw is an auto generated low-level Go binding around an Ethereum contract.
type RoyaltyHookRaw struct {
	Contract *RoyaltyHook // Generic contract binding to access the raw methods on
}

// RoyaltyHookCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RoyaltyHookCallerRaw struct {
	Contract *RoyaltyHookCaller // Generic read-only contract binding to access the raw methods on
}

// RoyaltyHookTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RoyaltyHookTransactorRaw struct {
	Contract *RoyaltyHookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRoyaltyHook creates a new instance of RoyaltyHook, bound to a specific deployed contract.
func NewRoyaltyHook(address common.Address, backend bind.ContractBackend) (*RoyaltyHook, error) {
	contract, err := bindRoyaltyHook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RoyaltyHook{RoyaltyHookCaller: RoyaltyHookCaller{contract: contract}, RoyaltyHookTransactor: RoyaltyHookTransactor{contract: contract}, RoyaltyHookFilterer: RoyaltyHookFilterer{contract: contract}}, nil
}

// NewRoyaltyHookCaller creates a new read-only instance of RoyaltyHook, bound to a specific deployed contract.
func NewRoyaltyHookCaller(address common.Address, caller bind.ContractCaller) (*RoyaltyHookCaller, error) {
	contract, err := bindRoyaltyHook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RoyaltyHookCaller{contract: contract}, nil
}

// NewRoyaltyHookTransactor creates a new write-only instance of RoyaltyHook, bound to a specific deployed contract.
func NewRoyaltyHookTransactor(address common.Address, transactor bind.ContractTransactor) (*RoyaltyHookTransactor, error) {
	contract, err := bindRoyaltyHook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RoyaltyHookTransactor{contract: contract}, nil
}

// NewRoyaltyHookFilterer creates a new log filterer instance of RoyaltyHook, bound to a specific deployed contract.
func NewRoyaltyHookFilterer(address common.Address, filterer bind.ContractFilterer) (*RoyaltyHookFilterer, error) {
	contract, err := bindRoyaltyHook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RoyaltyHookFilterer{contract: contract}, nil
}

// bindRoyaltyHook binds a generic wrapper to an already deployed contract.
func bindRoyaltyHook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RoyaltyHookMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RoyaltyHook *RoyaltyHookRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RoyaltyHook.Contract.RoyaltyHookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RoyaltyHook *RoyaltyHookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RoyaltyHook.Contract.RoyaltyHookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RoyaltyHook *RoyaltyHookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RoyaltyHook.Contract.RoyaltyHookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RoyaltyHook *RoyaltyHookCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _RoyaltyHook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RoyaltyHook *RoyaltyHookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RoyaltyHook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RoyaltyHook *RoyaltyHookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RoyaltyHook.Contract.contract.Transact(opts, method, params...)
}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_RoyaltyHook *RoyaltyHookCaller) BPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _RoyaltyHook.contract.Call(opts, &out, "BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_RoyaltyHook *RoyaltyHookSession) BPS() (*big.Int, error) {
	return _RoyaltyHook.Contract.BPS(&_RoyaltyHook.CallOpts)
}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_RoyaltyHook *RoyaltyHookCallerSession) BPS() (*big.Int, error) {
	return _RoyaltyHook.Contract.BPS(&_RoyaltyHook.CallOpts)
}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_RoyaltyHook *RoyaltyHookCaller) HookName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _RoyaltyHook.contract.Call(opts, &out, "hookName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_RoyaltyHook *RoyaltyHookSession) HookName() (string, error) {
	return _RoyaltyHook.Contract.HookName(&_RoyaltyHook.CallOpts)
}

// HookName is a free data retrieval call binding the contract method 0x144a07c0.
//
// Solidity: function hookName() pure returns(string)
func (_RoyaltyHook *RoyaltyHookCallerSession) HookName() (string, error) {
	return _RoyaltyHook.Contract.HookName(&_RoyaltyHook.CallOpts)
}

// OnAccept is a free data retrieval call binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookCaller) OnAccept(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _RoyaltyHook.contract.Call(opts, &out, "onAccept", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnAccept is a free data retrieval call binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookSession) OnAccept(arg0 *big.Int, arg1 []byte) error {
	return _RoyaltyHook.Contract.OnAccept(&_RoyaltyHook.CallOpts, arg0, arg1)
}

// OnAccept is a free data retrieval call binding the contract method 0xabe7b240.
//
// Solidity: function onAccept(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookCallerSession) OnAccept(arg0 *big.Int, arg1 []byte) error {
	return _RoyaltyHook.Contract.OnAccept(&_RoyaltyHook.CallOpts, arg0, arg1)
}

// OnCancel is a free data retrieval call binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookCaller) OnCancel(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _RoyaltyHook.contract.Call(opts, &out, "onCancel", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnCancel is a free data retrieval call binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookSession) OnCancel(arg0 *big.Int, arg1 []byte) error {
	return _RoyaltyHook.Contract.OnCancel(&_RoyaltyHook.CallOpts, arg0, arg1)
}

// OnCancel is a free data retrieval call binding the contract method 0x5b6d048f.
//
// Solidity: function onCancel(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookCallerSession) OnCancel(arg0 *big.Int, arg1 []byte) error {
	return _RoyaltyHook.Contract.OnCancel(&_RoyaltyHook.CallOpts, arg0, arg1)
}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookCaller) OnReject(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _RoyaltyHook.contract.Call(opts, &out, "onReject", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookSession) OnReject(arg0 *big.Int, arg1 []byte) error {
	return _RoyaltyHook.Contract.OnReject(&_RoyaltyHook.CallOpts, arg0, arg1)
}

// OnReject is a free data retrieval call binding the contract method 0x87e477dd.
//
// Solidity: function onReject(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookCallerSession) OnReject(arg0 *big.Int, arg1 []byte) error {
	return _RoyaltyHook.Contract.OnReject(&_RoyaltyHook.CallOpts, arg0, arg1)
}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookCaller) OnSubmit(opts *bind.CallOpts, arg0 *big.Int, arg1 []byte) error {
	var out []interface{}
	err := _RoyaltyHook.contract.Call(opts, &out, "onSubmit", arg0, arg1)

	if err != nil {
		return err
	}

	return err

}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookSession) OnSubmit(arg0 *big.Int, arg1 []byte) error {
	return _RoyaltyHook.Contract.OnSubmit(&_RoyaltyHook.CallOpts, arg0, arg1)
}

// OnSubmit is a free data retrieval call binding the contract method 0x25f0968e.
//
// Solidity: function onSubmit(uint256 , bytes ) pure returns()
func (_RoyaltyHook *RoyaltyHookCallerSession) OnSubmit(arg0 *big.Int, arg1 []byte) error {
	return _RoyaltyHook.Contract.OnSubmit(&_RoyaltyHook.CallOpts, arg0, arg1)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes context) returns()
func (_RoyaltyHook *RoyaltyHookTransactor) OnApprove(opts *bind.TransactOpts, jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _RoyaltyHook.contract.Transact(opts, "onApprove", jobId, context)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes context) returns()
func (_RoyaltyHook *RoyaltyHookSession) OnApprove(jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _RoyaltyHook.Contract.OnApprove(&_RoyaltyHook.TransactOpts, jobId, context)
}

// OnApprove is a paid mutator transaction binding the contract method 0x0657118e.
//
// Solidity: function onApprove(uint256 jobId, bytes context) returns()
func (_RoyaltyHook *RoyaltyHookTransactorSession) OnApprove(jobId *big.Int, context []byte) (*types.Transaction, error) {
	return _RoyaltyHook.Contract.OnApprove(&_RoyaltyHook.TransactOpts, jobId, context)
}

// RoyaltyHookRoyaltyPaidIterator is returned from FilterRoyaltyPaid and is used to iterate over the raw logs and unpacked data for RoyaltyPaid events raised by the RoyaltyHook contract.
type RoyaltyHookRoyaltyPaidIterator struct {
	Event *RoyaltyHookRoyaltyPaid // Event containing the contract specifics and raw log

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
func (it *RoyaltyHookRoyaltyPaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RoyaltyHookRoyaltyPaid)
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
		it.Event = new(RoyaltyHookRoyaltyPaid)
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
func (it *RoyaltyHookRoyaltyPaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RoyaltyHookRoyaltyPaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RoyaltyHookRoyaltyPaid represents a RoyaltyPaid event raised by the RoyaltyHook contract.
type RoyaltyHookRoyaltyPaid struct {
	JobId     *big.Int
	Token     common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRoyaltyPaid is a free log retrieval operation binding the contract event 0xf3f83c9d88dc83497cc705129934b4aa146d32ecd677e6ed32c57009a077067e.
//
// Solidity: event RoyaltyPaid(uint256 indexed jobId, address indexed token, address indexed recipient, uint256 amount)
func (_RoyaltyHook *RoyaltyHookFilterer) FilterRoyaltyPaid(opts *bind.FilterOpts, jobId []*big.Int, token []common.Address, recipient []common.Address) (*RoyaltyHookRoyaltyPaidIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _RoyaltyHook.contract.FilterLogs(opts, "RoyaltyPaid", jobIdRule, tokenRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &RoyaltyHookRoyaltyPaidIterator{contract: _RoyaltyHook.contract, event: "RoyaltyPaid", logs: logs, sub: sub}, nil
}

// WatchRoyaltyPaid is a free log subscription operation binding the contract event 0xf3f83c9d88dc83497cc705129934b4aa146d32ecd677e6ed32c57009a077067e.
//
// Solidity: event RoyaltyPaid(uint256 indexed jobId, address indexed token, address indexed recipient, uint256 amount)
func (_RoyaltyHook *RoyaltyHookFilterer) WatchRoyaltyPaid(opts *bind.WatchOpts, sink chan<- *RoyaltyHookRoyaltyPaid, jobId []*big.Int, token []common.Address, recipient []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _RoyaltyHook.contract.WatchLogs(opts, "RoyaltyPaid", jobIdRule, tokenRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RoyaltyHookRoyaltyPaid)
				if err := _RoyaltyHook.contract.UnpackLog(event, "RoyaltyPaid", log); err != nil {
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

// ParseRoyaltyPaid is a log parse operation binding the contract event 0xf3f83c9d88dc83497cc705129934b4aa146d32ecd677e6ed32c57009a077067e.
//
// Solidity: event RoyaltyPaid(uint256 indexed jobId, address indexed token, address indexed recipient, uint256 amount)
func (_RoyaltyHook *RoyaltyHookFilterer) ParseRoyaltyPaid(log types.Log) (*RoyaltyHookRoyaltyPaid, error) {
	event := new(RoyaltyHookRoyaltyPaid)
	if err := _RoyaltyHook.contract.UnpackLog(event, "RoyaltyPaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
