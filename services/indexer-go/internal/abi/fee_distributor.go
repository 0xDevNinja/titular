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

// FeeDistributorMetaData contains all meta data concerning the FeeDistributor contract.
var FeeDistributorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"ve\",\"type\":\"address\",\"internalType\":\"contractIVeTITU\"},{\"name\":\"treasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"START_WEEK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TREASURY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIVeTITU\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WEEK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"checkpointToken\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokens\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimable\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastTokenBalance\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastTokenCheckpoint\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokensPerWeek\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userNextClaimWeek\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Claimed\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"weeksAdvanced\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TokenCheckpointed\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"week\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"NothingToClaim\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// FeeDistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use FeeDistributorMetaData.ABI instead.
var FeeDistributorABI = FeeDistributorMetaData.ABI

// FeeDistributor is an auto generated Go binding around an Ethereum contract.
type FeeDistributor struct {
	FeeDistributorCaller     // Read-only binding to the contract
	FeeDistributorTransactor // Write-only binding to the contract
	FeeDistributorFilterer   // Log filterer for contract events
}

// FeeDistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeeDistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeDistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeeDistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeDistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeeDistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeDistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeeDistributorSession struct {
	Contract     *FeeDistributor   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeeDistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeeDistributorCallerSession struct {
	Contract *FeeDistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// FeeDistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeeDistributorTransactorSession struct {
	Contract     *FeeDistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// FeeDistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeeDistributorRaw struct {
	Contract *FeeDistributor // Generic contract binding to access the raw methods on
}

// FeeDistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeeDistributorCallerRaw struct {
	Contract *FeeDistributorCaller // Generic read-only contract binding to access the raw methods on
}

// FeeDistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeeDistributorTransactorRaw struct {
	Contract *FeeDistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeeDistributor creates a new instance of FeeDistributor, bound to a specific deployed contract.
func NewFeeDistributor(address common.Address, backend bind.ContractBackend) (*FeeDistributor, error) {
	contract, err := bindFeeDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeDistributor{FeeDistributorCaller: FeeDistributorCaller{contract: contract}, FeeDistributorTransactor: FeeDistributorTransactor{contract: contract}, FeeDistributorFilterer: FeeDistributorFilterer{contract: contract}}, nil
}

// NewFeeDistributorCaller creates a new read-only instance of FeeDistributor, bound to a specific deployed contract.
func NewFeeDistributorCaller(address common.Address, caller bind.ContractCaller) (*FeeDistributorCaller, error) {
	contract, err := bindFeeDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeDistributorCaller{contract: contract}, nil
}

// NewFeeDistributorTransactor creates a new write-only instance of FeeDistributor, bound to a specific deployed contract.
func NewFeeDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeDistributorTransactor, error) {
	contract, err := bindFeeDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeDistributorTransactor{contract: contract}, nil
}

// NewFeeDistributorFilterer creates a new log filterer instance of FeeDistributor, bound to a specific deployed contract.
func NewFeeDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeDistributorFilterer, error) {
	contract, err := bindFeeDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeDistributorFilterer{contract: contract}, nil
}

// bindFeeDistributor binds a generic wrapper to an already deployed contract.
func bindFeeDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeDistributorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeDistributor *FeeDistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeDistributor.Contract.FeeDistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeDistributor *FeeDistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeDistributor.Contract.FeeDistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeDistributor *FeeDistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeDistributor.Contract.FeeDistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeDistributor *FeeDistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeDistributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeDistributor *FeeDistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeDistributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeDistributor *FeeDistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeDistributor.Contract.contract.Transact(opts, method, params...)
}

// STARTWEEK is a free data retrieval call binding the contract method 0x0e057047.
//
// Solidity: function START_WEEK() view returns(uint256)
func (_FeeDistributor *FeeDistributorCaller) STARTWEEK(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeDistributor.contract.Call(opts, &out, "START_WEEK")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// STARTWEEK is a free data retrieval call binding the contract method 0x0e057047.
//
// Solidity: function START_WEEK() view returns(uint256)
func (_FeeDistributor *FeeDistributorSession) STARTWEEK() (*big.Int, error) {
	return _FeeDistributor.Contract.STARTWEEK(&_FeeDistributor.CallOpts)
}

// STARTWEEK is a free data retrieval call binding the contract method 0x0e057047.
//
// Solidity: function START_WEEK() view returns(uint256)
func (_FeeDistributor *FeeDistributorCallerSession) STARTWEEK() (*big.Int, error) {
	return _FeeDistributor.Contract.STARTWEEK(&_FeeDistributor.CallOpts)
}

// TREASURY is a free data retrieval call binding the contract method 0x2d2c5565.
//
// Solidity: function TREASURY() view returns(address)
func (_FeeDistributor *FeeDistributorCaller) TREASURY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeDistributor.contract.Call(opts, &out, "TREASURY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TREASURY is a free data retrieval call binding the contract method 0x2d2c5565.
//
// Solidity: function TREASURY() view returns(address)
func (_FeeDistributor *FeeDistributorSession) TREASURY() (common.Address, error) {
	return _FeeDistributor.Contract.TREASURY(&_FeeDistributor.CallOpts)
}

// TREASURY is a free data retrieval call binding the contract method 0x2d2c5565.
//
// Solidity: function TREASURY() view returns(address)
func (_FeeDistributor *FeeDistributorCallerSession) TREASURY() (common.Address, error) {
	return _FeeDistributor.Contract.TREASURY(&_FeeDistributor.CallOpts)
}

// VE is a free data retrieval call binding the contract method 0xc863657d.
//
// Solidity: function VE() view returns(address)
func (_FeeDistributor *FeeDistributorCaller) VE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeDistributor.contract.Call(opts, &out, "VE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VE is a free data retrieval call binding the contract method 0xc863657d.
//
// Solidity: function VE() view returns(address)
func (_FeeDistributor *FeeDistributorSession) VE() (common.Address, error) {
	return _FeeDistributor.Contract.VE(&_FeeDistributor.CallOpts)
}

// VE is a free data retrieval call binding the contract method 0xc863657d.
//
// Solidity: function VE() view returns(address)
func (_FeeDistributor *FeeDistributorCallerSession) VE() (common.Address, error) {
	return _FeeDistributor.Contract.VE(&_FeeDistributor.CallOpts)
}

// WEEK is a free data retrieval call binding the contract method 0xf4359ce5.
//
// Solidity: function WEEK() view returns(uint256)
func (_FeeDistributor *FeeDistributorCaller) WEEK(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeDistributor.contract.Call(opts, &out, "WEEK")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WEEK is a free data retrieval call binding the contract method 0xf4359ce5.
//
// Solidity: function WEEK() view returns(uint256)
func (_FeeDistributor *FeeDistributorSession) WEEK() (*big.Int, error) {
	return _FeeDistributor.Contract.WEEK(&_FeeDistributor.CallOpts)
}

// WEEK is a free data retrieval call binding the contract method 0xf4359ce5.
//
// Solidity: function WEEK() view returns(uint256)
func (_FeeDistributor *FeeDistributorCallerSession) WEEK() (*big.Int, error) {
	return _FeeDistributor.Contract.WEEK(&_FeeDistributor.CallOpts)
}

// Claimable is a free data retrieval call binding the contract method 0xd4570c1c.
//
// Solidity: function claimable(address user, address token) view returns(uint256)
func (_FeeDistributor *FeeDistributorCaller) Claimable(opts *bind.CallOpts, user common.Address, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeDistributor.contract.Call(opts, &out, "claimable", user, token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Claimable is a free data retrieval call binding the contract method 0xd4570c1c.
//
// Solidity: function claimable(address user, address token) view returns(uint256)
func (_FeeDistributor *FeeDistributorSession) Claimable(user common.Address, token common.Address) (*big.Int, error) {
	return _FeeDistributor.Contract.Claimable(&_FeeDistributor.CallOpts, user, token)
}

// Claimable is a free data retrieval call binding the contract method 0xd4570c1c.
//
// Solidity: function claimable(address user, address token) view returns(uint256)
func (_FeeDistributor *FeeDistributorCallerSession) Claimable(user common.Address, token common.Address) (*big.Int, error) {
	return _FeeDistributor.Contract.Claimable(&_FeeDistributor.CallOpts, user, token)
}

// LastTokenBalance is a free data retrieval call binding the contract method 0xe49e133e.
//
// Solidity: function lastTokenBalance(address ) view returns(uint256)
func (_FeeDistributor *FeeDistributorCaller) LastTokenBalance(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeDistributor.contract.Call(opts, &out, "lastTokenBalance", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastTokenBalance is a free data retrieval call binding the contract method 0xe49e133e.
//
// Solidity: function lastTokenBalance(address ) view returns(uint256)
func (_FeeDistributor *FeeDistributorSession) LastTokenBalance(arg0 common.Address) (*big.Int, error) {
	return _FeeDistributor.Contract.LastTokenBalance(&_FeeDistributor.CallOpts, arg0)
}

// LastTokenBalance is a free data retrieval call binding the contract method 0xe49e133e.
//
// Solidity: function lastTokenBalance(address ) view returns(uint256)
func (_FeeDistributor *FeeDistributorCallerSession) LastTokenBalance(arg0 common.Address) (*big.Int, error) {
	return _FeeDistributor.Contract.LastTokenBalance(&_FeeDistributor.CallOpts, arg0)
}

// LastTokenCheckpoint is a free data retrieval call binding the contract method 0xb08a7451.
//
// Solidity: function lastTokenCheckpoint(address ) view returns(uint256)
func (_FeeDistributor *FeeDistributorCaller) LastTokenCheckpoint(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeDistributor.contract.Call(opts, &out, "lastTokenCheckpoint", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastTokenCheckpoint is a free data retrieval call binding the contract method 0xb08a7451.
//
// Solidity: function lastTokenCheckpoint(address ) view returns(uint256)
func (_FeeDistributor *FeeDistributorSession) LastTokenCheckpoint(arg0 common.Address) (*big.Int, error) {
	return _FeeDistributor.Contract.LastTokenCheckpoint(&_FeeDistributor.CallOpts, arg0)
}

// LastTokenCheckpoint is a free data retrieval call binding the contract method 0xb08a7451.
//
// Solidity: function lastTokenCheckpoint(address ) view returns(uint256)
func (_FeeDistributor *FeeDistributorCallerSession) LastTokenCheckpoint(arg0 common.Address) (*big.Int, error) {
	return _FeeDistributor.Contract.LastTokenCheckpoint(&_FeeDistributor.CallOpts, arg0)
}

// TokensPerWeek is a free data retrieval call binding the contract method 0x00d440c1.
//
// Solidity: function tokensPerWeek(address , uint256 ) view returns(uint256)
func (_FeeDistributor *FeeDistributorCaller) TokensPerWeek(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FeeDistributor.contract.Call(opts, &out, "tokensPerWeek", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokensPerWeek is a free data retrieval call binding the contract method 0x00d440c1.
//
// Solidity: function tokensPerWeek(address , uint256 ) view returns(uint256)
func (_FeeDistributor *FeeDistributorSession) TokensPerWeek(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _FeeDistributor.Contract.TokensPerWeek(&_FeeDistributor.CallOpts, arg0, arg1)
}

// TokensPerWeek is a free data retrieval call binding the contract method 0x00d440c1.
//
// Solidity: function tokensPerWeek(address , uint256 ) view returns(uint256)
func (_FeeDistributor *FeeDistributorCallerSession) TokensPerWeek(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _FeeDistributor.Contract.TokensPerWeek(&_FeeDistributor.CallOpts, arg0, arg1)
}

// UserNextClaimWeek is a free data retrieval call binding the contract method 0xd9fde0bb.
//
// Solidity: function userNextClaimWeek(address , address ) view returns(uint256)
func (_FeeDistributor *FeeDistributorCaller) UserNextClaimWeek(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeDistributor.contract.Call(opts, &out, "userNextClaimWeek", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserNextClaimWeek is a free data retrieval call binding the contract method 0xd9fde0bb.
//
// Solidity: function userNextClaimWeek(address , address ) view returns(uint256)
func (_FeeDistributor *FeeDistributorSession) UserNextClaimWeek(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _FeeDistributor.Contract.UserNextClaimWeek(&_FeeDistributor.CallOpts, arg0, arg1)
}

// UserNextClaimWeek is a free data retrieval call binding the contract method 0xd9fde0bb.
//
// Solidity: function userNextClaimWeek(address , address ) view returns(uint256)
func (_FeeDistributor *FeeDistributorCallerSession) UserNextClaimWeek(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _FeeDistributor.Contract.UserNextClaimWeek(&_FeeDistributor.CallOpts, arg0, arg1)
}

// CheckpointToken is a paid mutator transaction binding the contract method 0x3902b9bc.
//
// Solidity: function checkpointToken(address token) returns()
func (_FeeDistributor *FeeDistributorTransactor) CheckpointToken(opts *bind.TransactOpts, token common.Address) (*types.Transaction, error) {
	return _FeeDistributor.contract.Transact(opts, "checkpointToken", token)
}

// CheckpointToken is a paid mutator transaction binding the contract method 0x3902b9bc.
//
// Solidity: function checkpointToken(address token) returns()
func (_FeeDistributor *FeeDistributorSession) CheckpointToken(token common.Address) (*types.Transaction, error) {
	return _FeeDistributor.Contract.CheckpointToken(&_FeeDistributor.TransactOpts, token)
}

// CheckpointToken is a paid mutator transaction binding the contract method 0x3902b9bc.
//
// Solidity: function checkpointToken(address token) returns()
func (_FeeDistributor *FeeDistributorTransactorSession) CheckpointToken(token common.Address) (*types.Transaction, error) {
	return _FeeDistributor.Contract.CheckpointToken(&_FeeDistributor.TransactOpts, token)
}

// Claim is a paid mutator transaction binding the contract method 0x8e2eba09.
//
// Solidity: function claim(address user, address[] tokens) returns()
func (_FeeDistributor *FeeDistributorTransactor) Claim(opts *bind.TransactOpts, user common.Address, tokens []common.Address) (*types.Transaction, error) {
	return _FeeDistributor.contract.Transact(opts, "claim", user, tokens)
}

// Claim is a paid mutator transaction binding the contract method 0x8e2eba09.
//
// Solidity: function claim(address user, address[] tokens) returns()
func (_FeeDistributor *FeeDistributorSession) Claim(user common.Address, tokens []common.Address) (*types.Transaction, error) {
	return _FeeDistributor.Contract.Claim(&_FeeDistributor.TransactOpts, user, tokens)
}

// Claim is a paid mutator transaction binding the contract method 0x8e2eba09.
//
// Solidity: function claim(address user, address[] tokens) returns()
func (_FeeDistributor *FeeDistributorTransactorSession) Claim(user common.Address, tokens []common.Address) (*types.Transaction, error) {
	return _FeeDistributor.Contract.Claim(&_FeeDistributor.TransactOpts, user, tokens)
}

// FeeDistributorClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the FeeDistributor contract.
type FeeDistributorClaimedIterator struct {
	Event *FeeDistributorClaimed // Event containing the contract specifics and raw log

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
func (it *FeeDistributorClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeDistributorClaimed)
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
		it.Event = new(FeeDistributorClaimed)
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
func (it *FeeDistributorClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeDistributorClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeDistributorClaimed represents a Claimed event raised by the FeeDistributor contract.
type FeeDistributorClaimed struct {
	User          common.Address
	Token         common.Address
	Amount        *big.Int
	WeeksAdvanced *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0x2f6639d24651730c7bf57c95ddbf96d66d11477e4ec626876f92c22e5f365e68.
//
// Solidity: event Claimed(address indexed user, address indexed token, uint256 amount, uint256 weeksAdvanced)
func (_FeeDistributor *FeeDistributorFilterer) FilterClaimed(opts *bind.FilterOpts, user []common.Address, token []common.Address) (*FeeDistributorClaimedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeDistributor.contract.FilterLogs(opts, "Claimed", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeDistributorClaimedIterator{contract: _FeeDistributor.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0x2f6639d24651730c7bf57c95ddbf96d66d11477e4ec626876f92c22e5f365e68.
//
// Solidity: event Claimed(address indexed user, address indexed token, uint256 amount, uint256 weeksAdvanced)
func (_FeeDistributor *FeeDistributorFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *FeeDistributorClaimed, user []common.Address, token []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeDistributor.contract.WatchLogs(opts, "Claimed", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeDistributorClaimed)
				if err := _FeeDistributor.contract.UnpackLog(event, "Claimed", log); err != nil {
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

// ParseClaimed is a log parse operation binding the contract event 0x2f6639d24651730c7bf57c95ddbf96d66d11477e4ec626876f92c22e5f365e68.
//
// Solidity: event Claimed(address indexed user, address indexed token, uint256 amount, uint256 weeksAdvanced)
func (_FeeDistributor *FeeDistributorFilterer) ParseClaimed(log types.Log) (*FeeDistributorClaimed, error) {
	event := new(FeeDistributorClaimed)
	if err := _FeeDistributor.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeDistributorTokenCheckpointedIterator is returned from FilterTokenCheckpointed and is used to iterate over the raw logs and unpacked data for TokenCheckpointed events raised by the FeeDistributor contract.
type FeeDistributorTokenCheckpointedIterator struct {
	Event *FeeDistributorTokenCheckpointed // Event containing the contract specifics and raw log

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
func (it *FeeDistributorTokenCheckpointedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeDistributorTokenCheckpointed)
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
		it.Event = new(FeeDistributorTokenCheckpointed)
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
func (it *FeeDistributorTokenCheckpointedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeDistributorTokenCheckpointedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeDistributorTokenCheckpointed represents a TokenCheckpointed event raised by the FeeDistributor contract.
type FeeDistributorTokenCheckpointed struct {
	Token  common.Address
	Week   *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTokenCheckpointed is a free log retrieval operation binding the contract event 0x9b7f1a85a4c9b4e59e1b6527d9969c50cdfb3a1a467d0c4a51fb0ed8bf07f130.
//
// Solidity: event TokenCheckpointed(address indexed token, uint256 indexed week, uint256 amount)
func (_FeeDistributor *FeeDistributorFilterer) FilterTokenCheckpointed(opts *bind.FilterOpts, token []common.Address, week []*big.Int) (*FeeDistributorTokenCheckpointedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var weekRule []interface{}
	for _, weekItem := range week {
		weekRule = append(weekRule, weekItem)
	}

	logs, sub, err := _FeeDistributor.contract.FilterLogs(opts, "TokenCheckpointed", tokenRule, weekRule)
	if err != nil {
		return nil, err
	}
	return &FeeDistributorTokenCheckpointedIterator{contract: _FeeDistributor.contract, event: "TokenCheckpointed", logs: logs, sub: sub}, nil
}

// WatchTokenCheckpointed is a free log subscription operation binding the contract event 0x9b7f1a85a4c9b4e59e1b6527d9969c50cdfb3a1a467d0c4a51fb0ed8bf07f130.
//
// Solidity: event TokenCheckpointed(address indexed token, uint256 indexed week, uint256 amount)
func (_FeeDistributor *FeeDistributorFilterer) WatchTokenCheckpointed(opts *bind.WatchOpts, sink chan<- *FeeDistributorTokenCheckpointed, token []common.Address, week []*big.Int) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var weekRule []interface{}
	for _, weekItem := range week {
		weekRule = append(weekRule, weekItem)
	}

	logs, sub, err := _FeeDistributor.contract.WatchLogs(opts, "TokenCheckpointed", tokenRule, weekRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeDistributorTokenCheckpointed)
				if err := _FeeDistributor.contract.UnpackLog(event, "TokenCheckpointed", log); err != nil {
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

// ParseTokenCheckpointed is a log parse operation binding the contract event 0x9b7f1a85a4c9b4e59e1b6527d9969c50cdfb3a1a467d0c4a51fb0ed8bf07f130.
//
// Solidity: event TokenCheckpointed(address indexed token, uint256 indexed week, uint256 amount)
func (_FeeDistributor *FeeDistributorFilterer) ParseTokenCheckpointed(log types.Log) (*FeeDistributorTokenCheckpointed, error) {
	event := new(FeeDistributorTokenCheckpointed)
	if err := _FeeDistributor.contract.UnpackLog(event, "TokenCheckpointed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
