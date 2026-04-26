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

// LPLockMetaData contains all meta data concerning the LPLock contract.
var LPLockMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"lpToken_\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"beneficiary_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"LOCK_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"beneficiary\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"lockedAmount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lpToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"timeRemaining\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unlockTime\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawn\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Deposited\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdrawn\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyWithdrawn\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotBeneficiary\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"StillLocked\",\"inputs\":[{\"name\":\"secondsRemaining\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// LPLockABI is the input ABI used to generate the binding from.
// Deprecated: Use LPLockMetaData.ABI instead.
var LPLockABI = LPLockMetaData.ABI

// LPLock is an auto generated Go binding around an Ethereum contract.
type LPLock struct {
	LPLockCaller     // Read-only binding to the contract
	LPLockTransactor // Write-only binding to the contract
	LPLockFilterer   // Log filterer for contract events
}

// LPLockCaller is an auto generated read-only Go binding around an Ethereum contract.
type LPLockCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LPLockTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LPLockTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LPLockFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LPLockFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LPLockSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LPLockSession struct {
	Contract     *LPLock           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LPLockCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LPLockCallerSession struct {
	Contract *LPLockCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// LPLockTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LPLockTransactorSession struct {
	Contract     *LPLockTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LPLockRaw is an auto generated low-level Go binding around an Ethereum contract.
type LPLockRaw struct {
	Contract *LPLock // Generic contract binding to access the raw methods on
}

// LPLockCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LPLockCallerRaw struct {
	Contract *LPLockCaller // Generic read-only contract binding to access the raw methods on
}

// LPLockTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LPLockTransactorRaw struct {
	Contract *LPLockTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLPLock creates a new instance of LPLock, bound to a specific deployed contract.
func NewLPLock(address common.Address, backend bind.ContractBackend) (*LPLock, error) {
	contract, err := bindLPLock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LPLock{LPLockCaller: LPLockCaller{contract: contract}, LPLockTransactor: LPLockTransactor{contract: contract}, LPLockFilterer: LPLockFilterer{contract: contract}}, nil
}

// NewLPLockCaller creates a new read-only instance of LPLock, bound to a specific deployed contract.
func NewLPLockCaller(address common.Address, caller bind.ContractCaller) (*LPLockCaller, error) {
	contract, err := bindLPLock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LPLockCaller{contract: contract}, nil
}

// NewLPLockTransactor creates a new write-only instance of LPLock, bound to a specific deployed contract.
func NewLPLockTransactor(address common.Address, transactor bind.ContractTransactor) (*LPLockTransactor, error) {
	contract, err := bindLPLock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LPLockTransactor{contract: contract}, nil
}

// NewLPLockFilterer creates a new log filterer instance of LPLock, bound to a specific deployed contract.
func NewLPLockFilterer(address common.Address, filterer bind.ContractFilterer) (*LPLockFilterer, error) {
	contract, err := bindLPLock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LPLockFilterer{contract: contract}, nil
}

// bindLPLock binds a generic wrapper to an already deployed contract.
func bindLPLock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LPLockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LPLock *LPLockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LPLock.Contract.LPLockCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LPLock *LPLockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LPLock.Contract.LPLockTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LPLock *LPLockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LPLock.Contract.LPLockTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LPLock *LPLockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LPLock.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LPLock *LPLockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LPLock.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LPLock *LPLockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LPLock.Contract.contract.Transact(opts, method, params...)
}

// LOCKDURATION is a free data retrieval call binding the contract method 0x485d3834.
//
// Solidity: function LOCK_DURATION() view returns(uint256)
func (_LPLock *LPLockCaller) LOCKDURATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPLock.contract.Call(opts, &out, "LOCK_DURATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LOCKDURATION is a free data retrieval call binding the contract method 0x485d3834.
//
// Solidity: function LOCK_DURATION() view returns(uint256)
func (_LPLock *LPLockSession) LOCKDURATION() (*big.Int, error) {
	return _LPLock.Contract.LOCKDURATION(&_LPLock.CallOpts)
}

// LOCKDURATION is a free data retrieval call binding the contract method 0x485d3834.
//
// Solidity: function LOCK_DURATION() view returns(uint256)
func (_LPLock *LPLockCallerSession) LOCKDURATION() (*big.Int, error) {
	return _LPLock.Contract.LOCKDURATION(&_LPLock.CallOpts)
}

// Beneficiary is a free data retrieval call binding the contract method 0x38af3eed.
//
// Solidity: function beneficiary() view returns(address)
func (_LPLock *LPLockCaller) Beneficiary(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LPLock.contract.Call(opts, &out, "beneficiary")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Beneficiary is a free data retrieval call binding the contract method 0x38af3eed.
//
// Solidity: function beneficiary() view returns(address)
func (_LPLock *LPLockSession) Beneficiary() (common.Address, error) {
	return _LPLock.Contract.Beneficiary(&_LPLock.CallOpts)
}

// Beneficiary is a free data retrieval call binding the contract method 0x38af3eed.
//
// Solidity: function beneficiary() view returns(address)
func (_LPLock *LPLockCallerSession) Beneficiary() (common.Address, error) {
	return _LPLock.Contract.Beneficiary(&_LPLock.CallOpts)
}

// LockedAmount is a free data retrieval call binding the contract method 0x6ab28bc8.
//
// Solidity: function lockedAmount() view returns(uint256)
func (_LPLock *LPLockCaller) LockedAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPLock.contract.Call(opts, &out, "lockedAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockedAmount is a free data retrieval call binding the contract method 0x6ab28bc8.
//
// Solidity: function lockedAmount() view returns(uint256)
func (_LPLock *LPLockSession) LockedAmount() (*big.Int, error) {
	return _LPLock.Contract.LockedAmount(&_LPLock.CallOpts)
}

// LockedAmount is a free data retrieval call binding the contract method 0x6ab28bc8.
//
// Solidity: function lockedAmount() view returns(uint256)
func (_LPLock *LPLockCallerSession) LockedAmount() (*big.Int, error) {
	return _LPLock.Contract.LockedAmount(&_LPLock.CallOpts)
}

// LpToken is a free data retrieval call binding the contract method 0x5fcbd285.
//
// Solidity: function lpToken() view returns(address)
func (_LPLock *LPLockCaller) LpToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LPLock.contract.Call(opts, &out, "lpToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LpToken is a free data retrieval call binding the contract method 0x5fcbd285.
//
// Solidity: function lpToken() view returns(address)
func (_LPLock *LPLockSession) LpToken() (common.Address, error) {
	return _LPLock.Contract.LpToken(&_LPLock.CallOpts)
}

// LpToken is a free data retrieval call binding the contract method 0x5fcbd285.
//
// Solidity: function lpToken() view returns(address)
func (_LPLock *LPLockCallerSession) LpToken() (common.Address, error) {
	return _LPLock.Contract.LpToken(&_LPLock.CallOpts)
}

// TimeRemaining is a free data retrieval call binding the contract method 0xe3cfef60.
//
// Solidity: function timeRemaining() view returns(uint256)
func (_LPLock *LPLockCaller) TimeRemaining(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPLock.contract.Call(opts, &out, "timeRemaining")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TimeRemaining is a free data retrieval call binding the contract method 0xe3cfef60.
//
// Solidity: function timeRemaining() view returns(uint256)
func (_LPLock *LPLockSession) TimeRemaining() (*big.Int, error) {
	return _LPLock.Contract.TimeRemaining(&_LPLock.CallOpts)
}

// TimeRemaining is a free data retrieval call binding the contract method 0xe3cfef60.
//
// Solidity: function timeRemaining() view returns(uint256)
func (_LPLock *LPLockCallerSession) TimeRemaining() (*big.Int, error) {
	return _LPLock.Contract.TimeRemaining(&_LPLock.CallOpts)
}

// UnlockTime is a free data retrieval call binding the contract method 0x251c1aa3.
//
// Solidity: function unlockTime() view returns(uint256)
func (_LPLock *LPLockCaller) UnlockTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LPLock.contract.Call(opts, &out, "unlockTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UnlockTime is a free data retrieval call binding the contract method 0x251c1aa3.
//
// Solidity: function unlockTime() view returns(uint256)
func (_LPLock *LPLockSession) UnlockTime() (*big.Int, error) {
	return _LPLock.Contract.UnlockTime(&_LPLock.CallOpts)
}

// UnlockTime is a free data retrieval call binding the contract method 0x251c1aa3.
//
// Solidity: function unlockTime() view returns(uint256)
func (_LPLock *LPLockCallerSession) UnlockTime() (*big.Int, error) {
	return _LPLock.Contract.UnlockTime(&_LPLock.CallOpts)
}

// Withdrawn is a free data retrieval call binding the contract method 0xc80ec522.
//
// Solidity: function withdrawn() view returns(bool)
func (_LPLock *LPLockCaller) Withdrawn(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _LPLock.contract.Call(opts, &out, "withdrawn")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Withdrawn is a free data retrieval call binding the contract method 0xc80ec522.
//
// Solidity: function withdrawn() view returns(bool)
func (_LPLock *LPLockSession) Withdrawn() (bool, error) {
	return _LPLock.Contract.Withdrawn(&_LPLock.CallOpts)
}

// Withdrawn is a free data retrieval call binding the contract method 0xc80ec522.
//
// Solidity: function withdrawn() view returns(bool)
func (_LPLock *LPLockCallerSession) Withdrawn() (bool, error) {
	return _LPLock.Contract.Withdrawn(&_LPLock.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_LPLock *LPLockTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _LPLock.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_LPLock *LPLockSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _LPLock.Contract.Deposit(&_LPLock.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_LPLock *LPLockTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _LPLock.Contract.Deposit(&_LPLock.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_LPLock *LPLockTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LPLock.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_LPLock *LPLockSession) Withdraw() (*types.Transaction, error) {
	return _LPLock.Contract.Withdraw(&_LPLock.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_LPLock *LPLockTransactorSession) Withdraw() (*types.Transaction, error) {
	return _LPLock.Contract.Withdraw(&_LPLock.TransactOpts)
}

// LPLockDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the LPLock contract.
type LPLockDepositedIterator struct {
	Event *LPLockDeposited // Event containing the contract specifics and raw log

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
func (it *LPLockDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPLockDeposited)
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
		it.Event = new(LPLockDeposited)
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
func (it *LPLockDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPLockDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPLockDeposited represents a Deposited event raised by the LPLock contract.
type LPLockDeposited struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed from, uint256 amount)
func (_LPLock *LPLockFilterer) FilterDeposited(opts *bind.FilterOpts, from []common.Address) (*LPLockDepositedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LPLock.contract.FilterLogs(opts, "Deposited", fromRule)
	if err != nil {
		return nil, err
	}
	return &LPLockDepositedIterator{contract: _LPLock.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed from, uint256 amount)
func (_LPLock *LPLockFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *LPLockDeposited, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _LPLock.contract.WatchLogs(opts, "Deposited", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPLockDeposited)
				if err := _LPLock.contract.UnpackLog(event, "Deposited", log); err != nil {
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

// ParseDeposited is a log parse operation binding the contract event 0x2da466a7b24304f47e87fa2e1e5a81b9831ce54fec19055ce277ca2f39ba42c4.
//
// Solidity: event Deposited(address indexed from, uint256 amount)
func (_LPLock *LPLockFilterer) ParseDeposited(log types.Log) (*LPLockDeposited, error) {
	event := new(LPLockDeposited)
	if err := _LPLock.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LPLockWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the LPLock contract.
type LPLockWithdrawnIterator struct {
	Event *LPLockWithdrawn // Event containing the contract specifics and raw log

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
func (it *LPLockWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LPLockWithdrawn)
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
		it.Event = new(LPLockWithdrawn)
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
func (it *LPLockWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LPLockWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LPLockWithdrawn represents a Withdrawn event raised by the LPLock contract.
type LPLockWithdrawn struct {
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed to, uint256 amount)
func (_LPLock *LPLockFilterer) FilterWithdrawn(opts *bind.FilterOpts, to []common.Address) (*LPLockWithdrawnIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPLock.contract.FilterLogs(opts, "Withdrawn", toRule)
	if err != nil {
		return nil, err
	}
	return &LPLockWithdrawnIterator{contract: _LPLock.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed to, uint256 amount)
func (_LPLock *LPLockFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *LPLockWithdrawn, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _LPLock.contract.WatchLogs(opts, "Withdrawn", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LPLockWithdrawn)
				if err := _LPLock.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0x7084f5476618d8e60b11ef0d7d3f06914655adb8793e28ff7f018d4c76d505d5.
//
// Solidity: event Withdrawn(address indexed to, uint256 amount)
func (_LPLock *LPLockFilterer) ParseWithdrawn(log types.Log) (*LPLockWithdrawn, error) {
	event := new(LPLockWithdrawn)
	if err := _LPLock.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
