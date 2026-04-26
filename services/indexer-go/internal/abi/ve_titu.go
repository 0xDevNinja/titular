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

// VeTITUPoint is an auto generated low-level Go binding around an user-defined struct.
type VeTITUPoint struct {
	Bias  *big.Int
	Slope *big.Int
	Ts    *big.Int
}

// VeTITUMetaData contains all meta data concerning the VeTITU contract.
var VeTITUMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"CLOCK_MODE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"MAXTIME\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TOKEN\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"WEEK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"clock\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createLock\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"unlockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getPastTotalSupply\",\"inputs\":[{\"name\":\"timepoint\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPastVotes\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timepoint\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getVotes\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"globalPointCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"increaseAmount\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"increaseUnlockTime\",\"inputs\":[{\"name\":\"newUnlockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"locked\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"int128\",\"internalType\":\"int128\"},{\"name\":\"end\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalLocked\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"userPointAt\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"i\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structVeTITU.Point\",\"components\":[{\"name\":\"bias\",\"type\":\"int128\",\"internalType\":\"int128\"},{\"name\":\"slope\",\"type\":\"int128\",\"internalType\":\"int128\"},{\"name\":\"ts\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"userPointCount\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Deposit\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"locktime\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"depositType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumVeTITU.DepositType\"},{\"name\":\"ts\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Supply\",\"inputs\":[{\"name\":\"prevSupply\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"supply\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"ts\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"LockExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LockNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoLock\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoTransfersAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"UnlockTimeNotInFuture\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnlockTimeNotLater\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnlockTimeTooLong\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnlockTimeTooShort\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// VeTITUABI is the input ABI used to generate the binding from.
// Deprecated: Use VeTITUMetaData.ABI instead.
var VeTITUABI = VeTITUMetaData.ABI

// VeTITU is an auto generated Go binding around an Ethereum contract.
type VeTITU struct {
	VeTITUCaller     // Read-only binding to the contract
	VeTITUTransactor // Write-only binding to the contract
	VeTITUFilterer   // Log filterer for contract events
}

// VeTITUCaller is an auto generated read-only Go binding around an Ethereum contract.
type VeTITUCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeTITUTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VeTITUTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeTITUFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VeTITUFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VeTITUSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VeTITUSession struct {
	Contract     *VeTITU           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VeTITUCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VeTITUCallerSession struct {
	Contract *VeTITUCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VeTITUTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VeTITUTransactorSession struct {
	Contract     *VeTITUTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VeTITURaw is an auto generated low-level Go binding around an Ethereum contract.
type VeTITURaw struct {
	Contract *VeTITU // Generic contract binding to access the raw methods on
}

// VeTITUCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VeTITUCallerRaw struct {
	Contract *VeTITUCaller // Generic read-only contract binding to access the raw methods on
}

// VeTITUTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VeTITUTransactorRaw struct {
	Contract *VeTITUTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVeTITU creates a new instance of VeTITU, bound to a specific deployed contract.
func NewVeTITU(address common.Address, backend bind.ContractBackend) (*VeTITU, error) {
	contract, err := bindVeTITU(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VeTITU{VeTITUCaller: VeTITUCaller{contract: contract}, VeTITUTransactor: VeTITUTransactor{contract: contract}, VeTITUFilterer: VeTITUFilterer{contract: contract}}, nil
}

// NewVeTITUCaller creates a new read-only instance of VeTITU, bound to a specific deployed contract.
func NewVeTITUCaller(address common.Address, caller bind.ContractCaller) (*VeTITUCaller, error) {
	contract, err := bindVeTITU(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VeTITUCaller{contract: contract}, nil
}

// NewVeTITUTransactor creates a new write-only instance of VeTITU, bound to a specific deployed contract.
func NewVeTITUTransactor(address common.Address, transactor bind.ContractTransactor) (*VeTITUTransactor, error) {
	contract, err := bindVeTITU(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VeTITUTransactor{contract: contract}, nil
}

// NewVeTITUFilterer creates a new log filterer instance of VeTITU, bound to a specific deployed contract.
func NewVeTITUFilterer(address common.Address, filterer bind.ContractFilterer) (*VeTITUFilterer, error) {
	contract, err := bindVeTITU(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VeTITUFilterer{contract: contract}, nil
}

// bindVeTITU binds a generic wrapper to an already deployed contract.
func bindVeTITU(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VeTITUMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VeTITU *VeTITURaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VeTITU.Contract.VeTITUCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VeTITU *VeTITURaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VeTITU.Contract.VeTITUTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VeTITU *VeTITURaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VeTITU.Contract.VeTITUTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VeTITU *VeTITUCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VeTITU.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VeTITU *VeTITUTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VeTITU.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VeTITU *VeTITUTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VeTITU.Contract.contract.Transact(opts, method, params...)
}

// CLOCKMODE is a free data retrieval call binding the contract method 0x4bf5d7e9.
//
// Solidity: function CLOCK_MODE() pure returns(string)
func (_VeTITU *VeTITUCaller) CLOCKMODE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "CLOCK_MODE")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// CLOCKMODE is a free data retrieval call binding the contract method 0x4bf5d7e9.
//
// Solidity: function CLOCK_MODE() pure returns(string)
func (_VeTITU *VeTITUSession) CLOCKMODE() (string, error) {
	return _VeTITU.Contract.CLOCKMODE(&_VeTITU.CallOpts)
}

// CLOCKMODE is a free data retrieval call binding the contract method 0x4bf5d7e9.
//
// Solidity: function CLOCK_MODE() pure returns(string)
func (_VeTITU *VeTITUCallerSession) CLOCKMODE() (string, error) {
	return _VeTITU.Contract.CLOCKMODE(&_VeTITU.CallOpts)
}

// MAXTIME is a free data retrieval call binding the contract method 0xee00ef3a.
//
// Solidity: function MAXTIME() view returns(uint256)
func (_VeTITU *VeTITUCaller) MAXTIME(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "MAXTIME")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXTIME is a free data retrieval call binding the contract method 0xee00ef3a.
//
// Solidity: function MAXTIME() view returns(uint256)
func (_VeTITU *VeTITUSession) MAXTIME() (*big.Int, error) {
	return _VeTITU.Contract.MAXTIME(&_VeTITU.CallOpts)
}

// MAXTIME is a free data retrieval call binding the contract method 0xee00ef3a.
//
// Solidity: function MAXTIME() view returns(uint256)
func (_VeTITU *VeTITUCallerSession) MAXTIME() (*big.Int, error) {
	return _VeTITU.Contract.MAXTIME(&_VeTITU.CallOpts)
}

// TOKEN is a free data retrieval call binding the contract method 0x82bfefc8.
//
// Solidity: function TOKEN() view returns(address)
func (_VeTITU *VeTITUCaller) TOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "TOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TOKEN is a free data retrieval call binding the contract method 0x82bfefc8.
//
// Solidity: function TOKEN() view returns(address)
func (_VeTITU *VeTITUSession) TOKEN() (common.Address, error) {
	return _VeTITU.Contract.TOKEN(&_VeTITU.CallOpts)
}

// TOKEN is a free data retrieval call binding the contract method 0x82bfefc8.
//
// Solidity: function TOKEN() view returns(address)
func (_VeTITU *VeTITUCallerSession) TOKEN() (common.Address, error) {
	return _VeTITU.Contract.TOKEN(&_VeTITU.CallOpts)
}

// WEEK is a free data retrieval call binding the contract method 0xf4359ce5.
//
// Solidity: function WEEK() view returns(uint256)
func (_VeTITU *VeTITUCaller) WEEK(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "WEEK")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WEEK is a free data retrieval call binding the contract method 0xf4359ce5.
//
// Solidity: function WEEK() view returns(uint256)
func (_VeTITU *VeTITUSession) WEEK() (*big.Int, error) {
	return _VeTITU.Contract.WEEK(&_VeTITU.CallOpts)
}

// WEEK is a free data retrieval call binding the contract method 0xf4359ce5.
//
// Solidity: function WEEK() view returns(uint256)
func (_VeTITU *VeTITUCallerSession) WEEK() (*big.Int, error) {
	return _VeTITU.Contract.WEEK(&_VeTITU.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_VeTITU *VeTITUCaller) BalanceOf(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "balanceOf", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_VeTITU *VeTITUSession) BalanceOf(user common.Address) (*big.Int, error) {
	return _VeTITU.Contract.BalanceOf(&_VeTITU.CallOpts, user)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address user) view returns(uint256)
func (_VeTITU *VeTITUCallerSession) BalanceOf(user common.Address) (*big.Int, error) {
	return _VeTITU.Contract.BalanceOf(&_VeTITU.CallOpts, user)
}

// Clock is a free data retrieval call binding the contract method 0x91ddadf4.
//
// Solidity: function clock() view returns(uint48)
func (_VeTITU *VeTITUCaller) Clock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "clock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Clock is a free data retrieval call binding the contract method 0x91ddadf4.
//
// Solidity: function clock() view returns(uint48)
func (_VeTITU *VeTITUSession) Clock() (*big.Int, error) {
	return _VeTITU.Contract.Clock(&_VeTITU.CallOpts)
}

// Clock is a free data retrieval call binding the contract method 0x91ddadf4.
//
// Solidity: function clock() view returns(uint48)
func (_VeTITU *VeTITUCallerSession) Clock() (*big.Int, error) {
	return _VeTITU.Contract.Clock(&_VeTITU.CallOpts)
}

// GetPastTotalSupply is a free data retrieval call binding the contract method 0x8e539e8c.
//
// Solidity: function getPastTotalSupply(uint256 timepoint) view returns(uint256)
func (_VeTITU *VeTITUCaller) GetPastTotalSupply(opts *bind.CallOpts, timepoint *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "getPastTotalSupply", timepoint)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPastTotalSupply is a free data retrieval call binding the contract method 0x8e539e8c.
//
// Solidity: function getPastTotalSupply(uint256 timepoint) view returns(uint256)
func (_VeTITU *VeTITUSession) GetPastTotalSupply(timepoint *big.Int) (*big.Int, error) {
	return _VeTITU.Contract.GetPastTotalSupply(&_VeTITU.CallOpts, timepoint)
}

// GetPastTotalSupply is a free data retrieval call binding the contract method 0x8e539e8c.
//
// Solidity: function getPastTotalSupply(uint256 timepoint) view returns(uint256)
func (_VeTITU *VeTITUCallerSession) GetPastTotalSupply(timepoint *big.Int) (*big.Int, error) {
	return _VeTITU.Contract.GetPastTotalSupply(&_VeTITU.CallOpts, timepoint)
}

// GetPastVotes is a free data retrieval call binding the contract method 0x3a46b1a8.
//
// Solidity: function getPastVotes(address user, uint256 timepoint) view returns(uint256)
func (_VeTITU *VeTITUCaller) GetPastVotes(opts *bind.CallOpts, user common.Address, timepoint *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "getPastVotes", user, timepoint)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPastVotes is a free data retrieval call binding the contract method 0x3a46b1a8.
//
// Solidity: function getPastVotes(address user, uint256 timepoint) view returns(uint256)
func (_VeTITU *VeTITUSession) GetPastVotes(user common.Address, timepoint *big.Int) (*big.Int, error) {
	return _VeTITU.Contract.GetPastVotes(&_VeTITU.CallOpts, user, timepoint)
}

// GetPastVotes is a free data retrieval call binding the contract method 0x3a46b1a8.
//
// Solidity: function getPastVotes(address user, uint256 timepoint) view returns(uint256)
func (_VeTITU *VeTITUCallerSession) GetPastVotes(user common.Address, timepoint *big.Int) (*big.Int, error) {
	return _VeTITU.Contract.GetPastVotes(&_VeTITU.CallOpts, user, timepoint)
}

// GetVotes is a free data retrieval call binding the contract method 0x9ab24eb0.
//
// Solidity: function getVotes(address user) view returns(uint256)
func (_VeTITU *VeTITUCaller) GetVotes(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "getVotes", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVotes is a free data retrieval call binding the contract method 0x9ab24eb0.
//
// Solidity: function getVotes(address user) view returns(uint256)
func (_VeTITU *VeTITUSession) GetVotes(user common.Address) (*big.Int, error) {
	return _VeTITU.Contract.GetVotes(&_VeTITU.CallOpts, user)
}

// GetVotes is a free data retrieval call binding the contract method 0x9ab24eb0.
//
// Solidity: function getVotes(address user) view returns(uint256)
func (_VeTITU *VeTITUCallerSession) GetVotes(user common.Address) (*big.Int, error) {
	return _VeTITU.Contract.GetVotes(&_VeTITU.CallOpts, user)
}

// GlobalPointCount is a free data retrieval call binding the contract method 0xc1a93a79.
//
// Solidity: function globalPointCount() view returns(uint256)
func (_VeTITU *VeTITUCaller) GlobalPointCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "globalPointCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalPointCount is a free data retrieval call binding the contract method 0xc1a93a79.
//
// Solidity: function globalPointCount() view returns(uint256)
func (_VeTITU *VeTITUSession) GlobalPointCount() (*big.Int, error) {
	return _VeTITU.Contract.GlobalPointCount(&_VeTITU.CallOpts)
}

// GlobalPointCount is a free data retrieval call binding the contract method 0xc1a93a79.
//
// Solidity: function globalPointCount() view returns(uint256)
func (_VeTITU *VeTITUCallerSession) GlobalPointCount() (*big.Int, error) {
	return _VeTITU.Contract.GlobalPointCount(&_VeTITU.CallOpts)
}

// Locked is a free data retrieval call binding the contract method 0xcbf9fe5f.
//
// Solidity: function locked(address ) view returns(int128 amount, uint256 end)
func (_VeTITU *VeTITUCaller) Locked(opts *bind.CallOpts, arg0 common.Address) (struct {
	Amount *big.Int
	End    *big.Int
}, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "locked", arg0)

	outstruct := new(struct {
		Amount *big.Int
		End    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.End = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Locked is a free data retrieval call binding the contract method 0xcbf9fe5f.
//
// Solidity: function locked(address ) view returns(int128 amount, uint256 end)
func (_VeTITU *VeTITUSession) Locked(arg0 common.Address) (struct {
	Amount *big.Int
	End    *big.Int
}, error) {
	return _VeTITU.Contract.Locked(&_VeTITU.CallOpts, arg0)
}

// Locked is a free data retrieval call binding the contract method 0xcbf9fe5f.
//
// Solidity: function locked(address ) view returns(int128 amount, uint256 end)
func (_VeTITU *VeTITUCallerSession) Locked(arg0 common.Address) (struct {
	Amount *big.Int
	End    *big.Int
}, error) {
	return _VeTITU.Contract.Locked(&_VeTITU.CallOpts, arg0)
}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_VeTITU *VeTITUCaller) TotalLocked(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "totalLocked")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_VeTITU *VeTITUSession) TotalLocked() (*big.Int, error) {
	return _VeTITU.Contract.TotalLocked(&_VeTITU.CallOpts)
}

// TotalLocked is a free data retrieval call binding the contract method 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (_VeTITU *VeTITUCallerSession) TotalLocked() (*big.Int, error) {
	return _VeTITU.Contract.TotalLocked(&_VeTITU.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VeTITU *VeTITUCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VeTITU *VeTITUSession) TotalSupply() (*big.Int, error) {
	return _VeTITU.Contract.TotalSupply(&_VeTITU.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_VeTITU *VeTITUCallerSession) TotalSupply() (*big.Int, error) {
	return _VeTITU.Contract.TotalSupply(&_VeTITU.CallOpts)
}

// Transfer is a free data retrieval call binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) pure returns(bool)
func (_VeTITU *VeTITUCaller) Transfer(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "transfer", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Transfer is a free data retrieval call binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) pure returns(bool)
func (_VeTITU *VeTITUSession) Transfer(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _VeTITU.Contract.Transfer(&_VeTITU.CallOpts, arg0, arg1)
}

// Transfer is a free data retrieval call binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) pure returns(bool)
func (_VeTITU *VeTITUCallerSession) Transfer(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _VeTITU.Contract.Transfer(&_VeTITU.CallOpts, arg0, arg1)
}

// TransferFrom is a free data retrieval call binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) pure returns(bool)
func (_VeTITU *VeTITUCaller) TransferFrom(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (bool, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "transferFrom", arg0, arg1, arg2)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// TransferFrom is a free data retrieval call binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) pure returns(bool)
func (_VeTITU *VeTITUSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (bool, error) {
	return _VeTITU.Contract.TransferFrom(&_VeTITU.CallOpts, arg0, arg1, arg2)
}

// TransferFrom is a free data retrieval call binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) pure returns(bool)
func (_VeTITU *VeTITUCallerSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (bool, error) {
	return _VeTITU.Contract.TransferFrom(&_VeTITU.CallOpts, arg0, arg1, arg2)
}

// UserPointAt is a free data retrieval call binding the contract method 0x4b3b9395.
//
// Solidity: function userPointAt(address user, uint256 i) view returns((int128,int128,uint256))
func (_VeTITU *VeTITUCaller) UserPointAt(opts *bind.CallOpts, user common.Address, i *big.Int) (VeTITUPoint, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "userPointAt", user, i)

	if err != nil {
		return *new(VeTITUPoint), err
	}

	out0 := *abi.ConvertType(out[0], new(VeTITUPoint)).(*VeTITUPoint)

	return out0, err

}

// UserPointAt is a free data retrieval call binding the contract method 0x4b3b9395.
//
// Solidity: function userPointAt(address user, uint256 i) view returns((int128,int128,uint256))
func (_VeTITU *VeTITUSession) UserPointAt(user common.Address, i *big.Int) (VeTITUPoint, error) {
	return _VeTITU.Contract.UserPointAt(&_VeTITU.CallOpts, user, i)
}

// UserPointAt is a free data retrieval call binding the contract method 0x4b3b9395.
//
// Solidity: function userPointAt(address user, uint256 i) view returns((int128,int128,uint256))
func (_VeTITU *VeTITUCallerSession) UserPointAt(user common.Address, i *big.Int) (VeTITUPoint, error) {
	return _VeTITU.Contract.UserPointAt(&_VeTITU.CallOpts, user, i)
}

// UserPointCount is a free data retrieval call binding the contract method 0xf6a609fb.
//
// Solidity: function userPointCount(address user) view returns(uint256)
func (_VeTITU *VeTITUCaller) UserPointCount(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VeTITU.contract.Call(opts, &out, "userPointCount", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserPointCount is a free data retrieval call binding the contract method 0xf6a609fb.
//
// Solidity: function userPointCount(address user) view returns(uint256)
func (_VeTITU *VeTITUSession) UserPointCount(user common.Address) (*big.Int, error) {
	return _VeTITU.Contract.UserPointCount(&_VeTITU.CallOpts, user)
}

// UserPointCount is a free data retrieval call binding the contract method 0xf6a609fb.
//
// Solidity: function userPointCount(address user) view returns(uint256)
func (_VeTITU *VeTITUCallerSession) UserPointCount(user common.Address) (*big.Int, error) {
	return _VeTITU.Contract.UserPointCount(&_VeTITU.CallOpts, user)
}

// CreateLock is a paid mutator transaction binding the contract method 0xb52c05fe.
//
// Solidity: function createLock(uint256 amount, uint256 unlockTime) returns()
func (_VeTITU *VeTITUTransactor) CreateLock(opts *bind.TransactOpts, amount *big.Int, unlockTime *big.Int) (*types.Transaction, error) {
	return _VeTITU.contract.Transact(opts, "createLock", amount, unlockTime)
}

// CreateLock is a paid mutator transaction binding the contract method 0xb52c05fe.
//
// Solidity: function createLock(uint256 amount, uint256 unlockTime) returns()
func (_VeTITU *VeTITUSession) CreateLock(amount *big.Int, unlockTime *big.Int) (*types.Transaction, error) {
	return _VeTITU.Contract.CreateLock(&_VeTITU.TransactOpts, amount, unlockTime)
}

// CreateLock is a paid mutator transaction binding the contract method 0xb52c05fe.
//
// Solidity: function createLock(uint256 amount, uint256 unlockTime) returns()
func (_VeTITU *VeTITUTransactorSession) CreateLock(amount *big.Int, unlockTime *big.Int) (*types.Transaction, error) {
	return _VeTITU.Contract.CreateLock(&_VeTITU.TransactOpts, amount, unlockTime)
}

// IncreaseAmount is a paid mutator transaction binding the contract method 0x15456eba.
//
// Solidity: function increaseAmount(uint256 amount) returns()
func (_VeTITU *VeTITUTransactor) IncreaseAmount(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _VeTITU.contract.Transact(opts, "increaseAmount", amount)
}

// IncreaseAmount is a paid mutator transaction binding the contract method 0x15456eba.
//
// Solidity: function increaseAmount(uint256 amount) returns()
func (_VeTITU *VeTITUSession) IncreaseAmount(amount *big.Int) (*types.Transaction, error) {
	return _VeTITU.Contract.IncreaseAmount(&_VeTITU.TransactOpts, amount)
}

// IncreaseAmount is a paid mutator transaction binding the contract method 0x15456eba.
//
// Solidity: function increaseAmount(uint256 amount) returns()
func (_VeTITU *VeTITUTransactorSession) IncreaseAmount(amount *big.Int) (*types.Transaction, error) {
	return _VeTITU.Contract.IncreaseAmount(&_VeTITU.TransactOpts, amount)
}

// IncreaseUnlockTime is a paid mutator transaction binding the contract method 0x7c616fe6.
//
// Solidity: function increaseUnlockTime(uint256 newUnlockTime) returns()
func (_VeTITU *VeTITUTransactor) IncreaseUnlockTime(opts *bind.TransactOpts, newUnlockTime *big.Int) (*types.Transaction, error) {
	return _VeTITU.contract.Transact(opts, "increaseUnlockTime", newUnlockTime)
}

// IncreaseUnlockTime is a paid mutator transaction binding the contract method 0x7c616fe6.
//
// Solidity: function increaseUnlockTime(uint256 newUnlockTime) returns()
func (_VeTITU *VeTITUSession) IncreaseUnlockTime(newUnlockTime *big.Int) (*types.Transaction, error) {
	return _VeTITU.Contract.IncreaseUnlockTime(&_VeTITU.TransactOpts, newUnlockTime)
}

// IncreaseUnlockTime is a paid mutator transaction binding the contract method 0x7c616fe6.
//
// Solidity: function increaseUnlockTime(uint256 newUnlockTime) returns()
func (_VeTITU *VeTITUTransactorSession) IncreaseUnlockTime(newUnlockTime *big.Int) (*types.Transaction, error) {
	return _VeTITU.Contract.IncreaseUnlockTime(&_VeTITU.TransactOpts, newUnlockTime)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_VeTITU *VeTITUTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VeTITU.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_VeTITU *VeTITUSession) Withdraw() (*types.Transaction, error) {
	return _VeTITU.Contract.Withdraw(&_VeTITU.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_VeTITU *VeTITUTransactorSession) Withdraw() (*types.Transaction, error) {
	return _VeTITU.Contract.Withdraw(&_VeTITU.TransactOpts)
}

// VeTITUDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the VeTITU contract.
type VeTITUDepositIterator struct {
	Event *VeTITUDeposit // Event containing the contract specifics and raw log

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
func (it *VeTITUDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeTITUDeposit)
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
		it.Event = new(VeTITUDeposit)
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
func (it *VeTITUDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeTITUDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeTITUDeposit represents a Deposit event raised by the VeTITU contract.
type VeTITUDeposit struct {
	Provider    common.Address
	Value       *big.Int
	Locktime    *big.Int
	DepositType uint8
	Ts          *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xbe9cf0e939c614fad640a623a53ba0a807c8cb503c4c4c8dacabe27b86ff2dd5.
//
// Solidity: event Deposit(address indexed provider, uint256 value, uint256 indexed locktime, uint8 depositType, uint256 ts)
func (_VeTITU *VeTITUFilterer) FilterDeposit(opts *bind.FilterOpts, provider []common.Address, locktime []*big.Int) (*VeTITUDepositIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	var locktimeRule []interface{}
	for _, locktimeItem := range locktime {
		locktimeRule = append(locktimeRule, locktimeItem)
	}

	logs, sub, err := _VeTITU.contract.FilterLogs(opts, "Deposit", providerRule, locktimeRule)
	if err != nil {
		return nil, err
	}
	return &VeTITUDepositIterator{contract: _VeTITU.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xbe9cf0e939c614fad640a623a53ba0a807c8cb503c4c4c8dacabe27b86ff2dd5.
//
// Solidity: event Deposit(address indexed provider, uint256 value, uint256 indexed locktime, uint8 depositType, uint256 ts)
func (_VeTITU *VeTITUFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *VeTITUDeposit, provider []common.Address, locktime []*big.Int) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	var locktimeRule []interface{}
	for _, locktimeItem := range locktime {
		locktimeRule = append(locktimeRule, locktimeItem)
	}

	logs, sub, err := _VeTITU.contract.WatchLogs(opts, "Deposit", providerRule, locktimeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeTITUDeposit)
				if err := _VeTITU.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xbe9cf0e939c614fad640a623a53ba0a807c8cb503c4c4c8dacabe27b86ff2dd5.
//
// Solidity: event Deposit(address indexed provider, uint256 value, uint256 indexed locktime, uint8 depositType, uint256 ts)
func (_VeTITU *VeTITUFilterer) ParseDeposit(log types.Log) (*VeTITUDeposit, error) {
	event := new(VeTITUDeposit)
	if err := _VeTITU.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VeTITUSupplyIterator is returned from FilterSupply and is used to iterate over the raw logs and unpacked data for Supply events raised by the VeTITU contract.
type VeTITUSupplyIterator struct {
	Event *VeTITUSupply // Event containing the contract specifics and raw log

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
func (it *VeTITUSupplyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeTITUSupply)
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
		it.Event = new(VeTITUSupply)
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
func (it *VeTITUSupplyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeTITUSupplyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeTITUSupply represents a Supply event raised by the VeTITU contract.
type VeTITUSupply struct {
	PrevSupply *big.Int
	Supply     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSupply is a free log retrieval operation binding the contract event 0x5e2aa66efd74cce82b21852e317e5490d9ecc9e6bb953ae24d90851258cc2f5c.
//
// Solidity: event Supply(uint256 prevSupply, uint256 supply)
func (_VeTITU *VeTITUFilterer) FilterSupply(opts *bind.FilterOpts) (*VeTITUSupplyIterator, error) {

	logs, sub, err := _VeTITU.contract.FilterLogs(opts, "Supply")
	if err != nil {
		return nil, err
	}
	return &VeTITUSupplyIterator{contract: _VeTITU.contract, event: "Supply", logs: logs, sub: sub}, nil
}

// WatchSupply is a free log subscription operation binding the contract event 0x5e2aa66efd74cce82b21852e317e5490d9ecc9e6bb953ae24d90851258cc2f5c.
//
// Solidity: event Supply(uint256 prevSupply, uint256 supply)
func (_VeTITU *VeTITUFilterer) WatchSupply(opts *bind.WatchOpts, sink chan<- *VeTITUSupply) (event.Subscription, error) {

	logs, sub, err := _VeTITU.contract.WatchLogs(opts, "Supply")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeTITUSupply)
				if err := _VeTITU.contract.UnpackLog(event, "Supply", log); err != nil {
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

// ParseSupply is a log parse operation binding the contract event 0x5e2aa66efd74cce82b21852e317e5490d9ecc9e6bb953ae24d90851258cc2f5c.
//
// Solidity: event Supply(uint256 prevSupply, uint256 supply)
func (_VeTITU *VeTITUFilterer) ParseSupply(log types.Log) (*VeTITUSupply, error) {
	event := new(VeTITUSupply)
	if err := _VeTITU.contract.UnpackLog(event, "Supply", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VeTITUWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the VeTITU contract.
type VeTITUWithdrawIterator struct {
	Event *VeTITUWithdraw // Event containing the contract specifics and raw log

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
func (it *VeTITUWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VeTITUWithdraw)
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
		it.Event = new(VeTITUWithdraw)
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
func (it *VeTITUWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VeTITUWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VeTITUWithdraw represents a Withdraw event raised by the VeTITU contract.
type VeTITUWithdraw struct {
	Provider common.Address
	Value    *big.Int
	Ts       *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xf279e6a1f5e320cca91135676d9cb6e44ca8a08c0b88342bcdb1144f6511b568.
//
// Solidity: event Withdraw(address indexed provider, uint256 value, uint256 ts)
func (_VeTITU *VeTITUFilterer) FilterWithdraw(opts *bind.FilterOpts, provider []common.Address) (*VeTITUWithdrawIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _VeTITU.contract.FilterLogs(opts, "Withdraw", providerRule)
	if err != nil {
		return nil, err
	}
	return &VeTITUWithdrawIterator{contract: _VeTITU.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xf279e6a1f5e320cca91135676d9cb6e44ca8a08c0b88342bcdb1144f6511b568.
//
// Solidity: event Withdraw(address indexed provider, uint256 value, uint256 ts)
func (_VeTITU *VeTITUFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *VeTITUWithdraw, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _VeTITU.contract.WatchLogs(opts, "Withdraw", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VeTITUWithdraw)
				if err := _VeTITU.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0xf279e6a1f5e320cca91135676d9cb6e44ca8a08c0b88342bcdb1144f6511b568.
//
// Solidity: event Withdraw(address indexed provider, uint256 value, uint256 ts)
func (_VeTITU *VeTITUFilterer) ParseWithdraw(log types.Log) (*VeTITUWithdraw, error) {
	event := new(VeTITUWithdraw)
	if err := _VeTITU.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
