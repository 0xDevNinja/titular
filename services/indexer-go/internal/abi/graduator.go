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

// GraduatorMetaData contains all meta data concerning the Graduator contract.
var GraduatorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"router_\",\"type\":\"address\",\"internalType\":\"contractIUniswapV2Router02\"},{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BPS_DENOMINATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_LIQUIDITY_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"graduate\",\"inputs\":[{\"name\":\"curve\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"graduated\",\"inputs\":[{\"name\":\"curve\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lpLocks\",\"inputs\":[{\"name\":\"curve\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"lpLock\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerLPLock\",\"inputs\":[{\"name\":\"curve\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpLock\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"router\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUniswapV2Router02\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"uniFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUniswapV2Factory\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Graduated\",\"inputs\":[{\"name\":\"curve\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"pair\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"liquidity\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"quoteAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"agentAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"LPLockRegistered\",\"inputs\":[{\"name\":\"curve\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"lock\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyGraduated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LPLockNotRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotYetGraduatedOnCurve\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// GraduatorABI is the input ABI used to generate the binding from.
// Deprecated: Use GraduatorMetaData.ABI instead.
var GraduatorABI = GraduatorMetaData.ABI

// Graduator is an auto generated Go binding around an Ethereum contract.
type Graduator struct {
	GraduatorCaller     // Read-only binding to the contract
	GraduatorTransactor // Write-only binding to the contract
	GraduatorFilterer   // Log filterer for contract events
}

// GraduatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type GraduatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GraduatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GraduatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GraduatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GraduatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GraduatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GraduatorSession struct {
	Contract     *Graduator        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GraduatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GraduatorCallerSession struct {
	Contract *GraduatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// GraduatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GraduatorTransactorSession struct {
	Contract     *GraduatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// GraduatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type GraduatorRaw struct {
	Contract *Graduator // Generic contract binding to access the raw methods on
}

// GraduatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GraduatorCallerRaw struct {
	Contract *GraduatorCaller // Generic read-only contract binding to access the raw methods on
}

// GraduatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GraduatorTransactorRaw struct {
	Contract *GraduatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGraduator creates a new instance of Graduator, bound to a specific deployed contract.
func NewGraduator(address common.Address, backend bind.ContractBackend) (*Graduator, error) {
	contract, err := bindGraduator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Graduator{GraduatorCaller: GraduatorCaller{contract: contract}, GraduatorTransactor: GraduatorTransactor{contract: contract}, GraduatorFilterer: GraduatorFilterer{contract: contract}}, nil
}

// NewGraduatorCaller creates a new read-only instance of Graduator, bound to a specific deployed contract.
func NewGraduatorCaller(address common.Address, caller bind.ContractCaller) (*GraduatorCaller, error) {
	contract, err := bindGraduator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GraduatorCaller{contract: contract}, nil
}

// NewGraduatorTransactor creates a new write-only instance of Graduator, bound to a specific deployed contract.
func NewGraduatorTransactor(address common.Address, transactor bind.ContractTransactor) (*GraduatorTransactor, error) {
	contract, err := bindGraduator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GraduatorTransactor{contract: contract}, nil
}

// NewGraduatorFilterer creates a new log filterer instance of Graduator, bound to a specific deployed contract.
func NewGraduatorFilterer(address common.Address, filterer bind.ContractFilterer) (*GraduatorFilterer, error) {
	contract, err := bindGraduator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GraduatorFilterer{contract: contract}, nil
}

// bindGraduator binds a generic wrapper to an already deployed contract.
func bindGraduator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GraduatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Graduator *GraduatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Graduator.Contract.GraduatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Graduator *GraduatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Graduator.Contract.GraduatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Graduator *GraduatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Graduator.Contract.GraduatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Graduator *GraduatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Graduator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Graduator *GraduatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Graduator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Graduator *GraduatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Graduator.Contract.contract.Transact(opts, method, params...)
}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint256)
func (_Graduator *GraduatorCaller) BPSDENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Graduator.contract.Call(opts, &out, "BPS_DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint256)
func (_Graduator *GraduatorSession) BPSDENOMINATOR() (*big.Int, error) {
	return _Graduator.Contract.BPSDENOMINATOR(&_Graduator.CallOpts)
}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint256)
func (_Graduator *GraduatorCallerSession) BPSDENOMINATOR() (*big.Int, error) {
	return _Graduator.Contract.BPSDENOMINATOR(&_Graduator.CallOpts)
}

// MINLIQUIDITYBPS is a free data retrieval call binding the contract method 0x75f0c2f0.
//
// Solidity: function MIN_LIQUIDITY_BPS() view returns(uint256)
func (_Graduator *GraduatorCaller) MINLIQUIDITYBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Graduator.contract.Call(opts, &out, "MIN_LIQUIDITY_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MINLIQUIDITYBPS is a free data retrieval call binding the contract method 0x75f0c2f0.
//
// Solidity: function MIN_LIQUIDITY_BPS() view returns(uint256)
func (_Graduator *GraduatorSession) MINLIQUIDITYBPS() (*big.Int, error) {
	return _Graduator.Contract.MINLIQUIDITYBPS(&_Graduator.CallOpts)
}

// MINLIQUIDITYBPS is a free data retrieval call binding the contract method 0x75f0c2f0.
//
// Solidity: function MIN_LIQUIDITY_BPS() view returns(uint256)
func (_Graduator *GraduatorCallerSession) MINLIQUIDITYBPS() (*big.Int, error) {
	return _Graduator.Contract.MINLIQUIDITYBPS(&_Graduator.CallOpts)
}

// Graduated is a free data retrieval call binding the contract method 0x9b083ec9.
//
// Solidity: function graduated(address curve) view returns(bool)
func (_Graduator *GraduatorCaller) Graduated(opts *bind.CallOpts, curve common.Address) (bool, error) {
	var out []interface{}
	err := _Graduator.contract.Call(opts, &out, "graduated", curve)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Graduated is a free data retrieval call binding the contract method 0x9b083ec9.
//
// Solidity: function graduated(address curve) view returns(bool)
func (_Graduator *GraduatorSession) Graduated(curve common.Address) (bool, error) {
	return _Graduator.Contract.Graduated(&_Graduator.CallOpts, curve)
}

// Graduated is a free data retrieval call binding the contract method 0x9b083ec9.
//
// Solidity: function graduated(address curve) view returns(bool)
func (_Graduator *GraduatorCallerSession) Graduated(curve common.Address) (bool, error) {
	return _Graduator.Contract.Graduated(&_Graduator.CallOpts, curve)
}

// LpLocks is a free data retrieval call binding the contract method 0xb05b1484.
//
// Solidity: function lpLocks(address curve) view returns(address lpLock)
func (_Graduator *GraduatorCaller) LpLocks(opts *bind.CallOpts, curve common.Address) (common.Address, error) {
	var out []interface{}
	err := _Graduator.contract.Call(opts, &out, "lpLocks", curve)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LpLocks is a free data retrieval call binding the contract method 0xb05b1484.
//
// Solidity: function lpLocks(address curve) view returns(address lpLock)
func (_Graduator *GraduatorSession) LpLocks(curve common.Address) (common.Address, error) {
	return _Graduator.Contract.LpLocks(&_Graduator.CallOpts, curve)
}

// LpLocks is a free data retrieval call binding the contract method 0xb05b1484.
//
// Solidity: function lpLocks(address curve) view returns(address lpLock)
func (_Graduator *GraduatorCallerSession) LpLocks(curve common.Address) (common.Address, error) {
	return _Graduator.Contract.LpLocks(&_Graduator.CallOpts, curve)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Graduator *GraduatorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Graduator.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Graduator *GraduatorSession) Owner() (common.Address, error) {
	return _Graduator.Contract.Owner(&_Graduator.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Graduator *GraduatorCallerSession) Owner() (common.Address, error) {
	return _Graduator.Contract.Owner(&_Graduator.CallOpts)
}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_Graduator *GraduatorCaller) Router(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Graduator.contract.Call(opts, &out, "router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_Graduator *GraduatorSession) Router() (common.Address, error) {
	return _Graduator.Contract.Router(&_Graduator.CallOpts)
}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_Graduator *GraduatorCallerSession) Router() (common.Address, error) {
	return _Graduator.Contract.Router(&_Graduator.CallOpts)
}

// UniFactory is a free data retrieval call binding the contract method 0x76771d4b.
//
// Solidity: function uniFactory() view returns(address)
func (_Graduator *GraduatorCaller) UniFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Graduator.contract.Call(opts, &out, "uniFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UniFactory is a free data retrieval call binding the contract method 0x76771d4b.
//
// Solidity: function uniFactory() view returns(address)
func (_Graduator *GraduatorSession) UniFactory() (common.Address, error) {
	return _Graduator.Contract.UniFactory(&_Graduator.CallOpts)
}

// UniFactory is a free data retrieval call binding the contract method 0x76771d4b.
//
// Solidity: function uniFactory() view returns(address)
func (_Graduator *GraduatorCallerSession) UniFactory() (common.Address, error) {
	return _Graduator.Contract.UniFactory(&_Graduator.CallOpts)
}

// Graduate is a paid mutator transaction binding the contract method 0xff6d8d05.
//
// Solidity: function graduate(address curve) returns()
func (_Graduator *GraduatorTransactor) Graduate(opts *bind.TransactOpts, curve common.Address) (*types.Transaction, error) {
	return _Graduator.contract.Transact(opts, "graduate", curve)
}

// Graduate is a paid mutator transaction binding the contract method 0xff6d8d05.
//
// Solidity: function graduate(address curve) returns()
func (_Graduator *GraduatorSession) Graduate(curve common.Address) (*types.Transaction, error) {
	return _Graduator.Contract.Graduate(&_Graduator.TransactOpts, curve)
}

// Graduate is a paid mutator transaction binding the contract method 0xff6d8d05.
//
// Solidity: function graduate(address curve) returns()
func (_Graduator *GraduatorTransactorSession) Graduate(curve common.Address) (*types.Transaction, error) {
	return _Graduator.Contract.Graduate(&_Graduator.TransactOpts, curve)
}

// RegisterLPLock is a paid mutator transaction binding the contract method 0xe3cd6627.
//
// Solidity: function registerLPLock(address curve, address lpLock) returns()
func (_Graduator *GraduatorTransactor) RegisterLPLock(opts *bind.TransactOpts, curve common.Address, lpLock common.Address) (*types.Transaction, error) {
	return _Graduator.contract.Transact(opts, "registerLPLock", curve, lpLock)
}

// RegisterLPLock is a paid mutator transaction binding the contract method 0xe3cd6627.
//
// Solidity: function registerLPLock(address curve, address lpLock) returns()
func (_Graduator *GraduatorSession) RegisterLPLock(curve common.Address, lpLock common.Address) (*types.Transaction, error) {
	return _Graduator.Contract.RegisterLPLock(&_Graduator.TransactOpts, curve, lpLock)
}

// RegisterLPLock is a paid mutator transaction binding the contract method 0xe3cd6627.
//
// Solidity: function registerLPLock(address curve, address lpLock) returns()
func (_Graduator *GraduatorTransactorSession) RegisterLPLock(curve common.Address, lpLock common.Address) (*types.Transaction, error) {
	return _Graduator.Contract.RegisterLPLock(&_Graduator.TransactOpts, curve, lpLock)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Graduator *GraduatorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Graduator.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Graduator *GraduatorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Graduator.Contract.RenounceOwnership(&_Graduator.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Graduator *GraduatorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Graduator.Contract.RenounceOwnership(&_Graduator.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Graduator *GraduatorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Graduator.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Graduator *GraduatorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Graduator.Contract.TransferOwnership(&_Graduator.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Graduator *GraduatorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Graduator.Contract.TransferOwnership(&_Graduator.TransactOpts, newOwner)
}

// GraduatorGraduatedIterator is returned from FilterGraduated and is used to iterate over the raw logs and unpacked data for Graduated events raised by the Graduator contract.
type GraduatorGraduatedIterator struct {
	Event *GraduatorGraduated // Event containing the contract specifics and raw log

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
func (it *GraduatorGraduatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GraduatorGraduated)
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
		it.Event = new(GraduatorGraduated)
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
func (it *GraduatorGraduatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GraduatorGraduatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GraduatorGraduated represents a Graduated event raised by the Graduator contract.
type GraduatorGraduated struct {
	Curve       common.Address
	Pair        common.Address
	Liquidity   *big.Int
	QuoteAmount *big.Int
	AgentAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGraduated is a free log retrieval operation binding the contract event 0x487dc7f66c623fb0ff13f9024a3ff9675453d069e075eceb12d9f8d7870e2374.
//
// Solidity: event Graduated(address indexed curve, address indexed pair, uint256 liquidity, uint256 quoteAmount, uint256 agentAmount)
func (_Graduator *GraduatorFilterer) FilterGraduated(opts *bind.FilterOpts, curve []common.Address, pair []common.Address) (*GraduatorGraduatedIterator, error) {

	var curveRule []interface{}
	for _, curveItem := range curve {
		curveRule = append(curveRule, curveItem)
	}
	var pairRule []interface{}
	for _, pairItem := range pair {
		pairRule = append(pairRule, pairItem)
	}

	logs, sub, err := _Graduator.contract.FilterLogs(opts, "Graduated", curveRule, pairRule)
	if err != nil {
		return nil, err
	}
	return &GraduatorGraduatedIterator{contract: _Graduator.contract, event: "Graduated", logs: logs, sub: sub}, nil
}

// WatchGraduated is a free log subscription operation binding the contract event 0x487dc7f66c623fb0ff13f9024a3ff9675453d069e075eceb12d9f8d7870e2374.
//
// Solidity: event Graduated(address indexed curve, address indexed pair, uint256 liquidity, uint256 quoteAmount, uint256 agentAmount)
func (_Graduator *GraduatorFilterer) WatchGraduated(opts *bind.WatchOpts, sink chan<- *GraduatorGraduated, curve []common.Address, pair []common.Address) (event.Subscription, error) {

	var curveRule []interface{}
	for _, curveItem := range curve {
		curveRule = append(curveRule, curveItem)
	}
	var pairRule []interface{}
	for _, pairItem := range pair {
		pairRule = append(pairRule, pairItem)
	}

	logs, sub, err := _Graduator.contract.WatchLogs(opts, "Graduated", curveRule, pairRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GraduatorGraduated)
				if err := _Graduator.contract.UnpackLog(event, "Graduated", log); err != nil {
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

// ParseGraduated is a log parse operation binding the contract event 0x487dc7f66c623fb0ff13f9024a3ff9675453d069e075eceb12d9f8d7870e2374.
//
// Solidity: event Graduated(address indexed curve, address indexed pair, uint256 liquidity, uint256 quoteAmount, uint256 agentAmount)
func (_Graduator *GraduatorFilterer) ParseGraduated(log types.Log) (*GraduatorGraduated, error) {
	event := new(GraduatorGraduated)
	if err := _Graduator.contract.UnpackLog(event, "Graduated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GraduatorLPLockRegisteredIterator is returned from FilterLPLockRegistered and is used to iterate over the raw logs and unpacked data for LPLockRegistered events raised by the Graduator contract.
type GraduatorLPLockRegisteredIterator struct {
	Event *GraduatorLPLockRegistered // Event containing the contract specifics and raw log

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
func (it *GraduatorLPLockRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GraduatorLPLockRegistered)
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
		it.Event = new(GraduatorLPLockRegistered)
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
func (it *GraduatorLPLockRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GraduatorLPLockRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GraduatorLPLockRegistered represents a LPLockRegistered event raised by the Graduator contract.
type GraduatorLPLockRegistered struct {
	Curve common.Address
	Lock  common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLPLockRegistered is a free log retrieval operation binding the contract event 0x84bb7f323dbdd64b7f1f9418890d70f99fe52649e0164c64a9335fa0cc941ce6.
//
// Solidity: event LPLockRegistered(address indexed curve, address indexed lock)
func (_Graduator *GraduatorFilterer) FilterLPLockRegistered(opts *bind.FilterOpts, curve []common.Address, lock []common.Address) (*GraduatorLPLockRegisteredIterator, error) {

	var curveRule []interface{}
	for _, curveItem := range curve {
		curveRule = append(curveRule, curveItem)
	}
	var lockRule []interface{}
	for _, lockItem := range lock {
		lockRule = append(lockRule, lockItem)
	}

	logs, sub, err := _Graduator.contract.FilterLogs(opts, "LPLockRegistered", curveRule, lockRule)
	if err != nil {
		return nil, err
	}
	return &GraduatorLPLockRegisteredIterator{contract: _Graduator.contract, event: "LPLockRegistered", logs: logs, sub: sub}, nil
}

// WatchLPLockRegistered is a free log subscription operation binding the contract event 0x84bb7f323dbdd64b7f1f9418890d70f99fe52649e0164c64a9335fa0cc941ce6.
//
// Solidity: event LPLockRegistered(address indexed curve, address indexed lock)
func (_Graduator *GraduatorFilterer) WatchLPLockRegistered(opts *bind.WatchOpts, sink chan<- *GraduatorLPLockRegistered, curve []common.Address, lock []common.Address) (event.Subscription, error) {

	var curveRule []interface{}
	for _, curveItem := range curve {
		curveRule = append(curveRule, curveItem)
	}
	var lockRule []interface{}
	for _, lockItem := range lock {
		lockRule = append(lockRule, lockItem)
	}

	logs, sub, err := _Graduator.contract.WatchLogs(opts, "LPLockRegistered", curveRule, lockRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GraduatorLPLockRegistered)
				if err := _Graduator.contract.UnpackLog(event, "LPLockRegistered", log); err != nil {
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

// ParseLPLockRegistered is a log parse operation binding the contract event 0x84bb7f323dbdd64b7f1f9418890d70f99fe52649e0164c64a9335fa0cc941ce6.
//
// Solidity: event LPLockRegistered(address indexed curve, address indexed lock)
func (_Graduator *GraduatorFilterer) ParseLPLockRegistered(log types.Log) (*GraduatorLPLockRegistered, error) {
	event := new(GraduatorLPLockRegistered)
	if err := _Graduator.contract.UnpackLog(event, "LPLockRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GraduatorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Graduator contract.
type GraduatorOwnershipTransferredIterator struct {
	Event *GraduatorOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GraduatorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GraduatorOwnershipTransferred)
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
		it.Event = new(GraduatorOwnershipTransferred)
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
func (it *GraduatorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GraduatorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GraduatorOwnershipTransferred represents a OwnershipTransferred event raised by the Graduator contract.
type GraduatorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Graduator *GraduatorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GraduatorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Graduator.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GraduatorOwnershipTransferredIterator{contract: _Graduator.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Graduator *GraduatorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GraduatorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Graduator.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GraduatorOwnershipTransferred)
				if err := _Graduator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Graduator *GraduatorFilterer) ParseOwnershipTransferred(log types.Log) (*GraduatorOwnershipTransferred, error) {
	event := new(GraduatorOwnershipTransferred)
	if err := _Graduator.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
