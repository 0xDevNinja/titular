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

// FeeSplitterMetaData contains all meta data concerning the FeeSplitter contract.
var FeeSplitterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_treasury\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_buybackBurner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"CALLER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SCHEDULE_A_BUYBACK_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SCHEDULE_A_TREASURY_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SCHEDULE_B_BUYBACK_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SCHEDULE_B_TREASURY_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buybackBurner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewSplit\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"schedule\",\"type\":\"uint8\",\"internalType\":\"enumFeeSplitter.Schedule\"}],\"outputs\":[{\"name\":\"primaryAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"treasuryAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"buybackAmt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setBuybackBurner\",\"inputs\":[{\"name\":\"newBurner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTreasury\",\"inputs\":[{\"name\":\"newTreasury\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"split\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"primary\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint8\",\"internalType\":\"enumFeeSplitter.Schedule\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BuybackBurnerUpdated\",\"inputs\":[{\"name\":\"oldBurner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newBurner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeSplit\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"primary\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"schedule\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumFeeSplitter.Schedule\"},{\"name\":\"primaryAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"treasuryAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"buybackAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TreasuryUpdated\",\"inputs\":[{\"name\":\"oldTreasury\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newTreasury\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// FeeSplitterABI is the input ABI used to generate the binding from.
// Deprecated: Use FeeSplitterMetaData.ABI instead.
var FeeSplitterABI = FeeSplitterMetaData.ABI

// FeeSplitter is an auto generated Go binding around an Ethereum contract.
type FeeSplitter struct {
	FeeSplitterCaller     // Read-only binding to the contract
	FeeSplitterTransactor // Write-only binding to the contract
	FeeSplitterFilterer   // Log filterer for contract events
}

// FeeSplitterCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeeSplitterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeSplitterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeeSplitterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeSplitterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeeSplitterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeSplitterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeeSplitterSession struct {
	Contract     *FeeSplitter      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeeSplitterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeeSplitterCallerSession struct {
	Contract *FeeSplitterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// FeeSplitterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeeSplitterTransactorSession struct {
	Contract     *FeeSplitterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// FeeSplitterRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeeSplitterRaw struct {
	Contract *FeeSplitter // Generic contract binding to access the raw methods on
}

// FeeSplitterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeeSplitterCallerRaw struct {
	Contract *FeeSplitterCaller // Generic read-only contract binding to access the raw methods on
}

// FeeSplitterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeeSplitterTransactorRaw struct {
	Contract *FeeSplitterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeeSplitter creates a new instance of FeeSplitter, bound to a specific deployed contract.
func NewFeeSplitter(address common.Address, backend bind.ContractBackend) (*FeeSplitter, error) {
	contract, err := bindFeeSplitter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeSplitter{FeeSplitterCaller: FeeSplitterCaller{contract: contract}, FeeSplitterTransactor: FeeSplitterTransactor{contract: contract}, FeeSplitterFilterer: FeeSplitterFilterer{contract: contract}}, nil
}

// NewFeeSplitterCaller creates a new read-only instance of FeeSplitter, bound to a specific deployed contract.
func NewFeeSplitterCaller(address common.Address, caller bind.ContractCaller) (*FeeSplitterCaller, error) {
	contract, err := bindFeeSplitter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeSplitterCaller{contract: contract}, nil
}

// NewFeeSplitterTransactor creates a new write-only instance of FeeSplitter, bound to a specific deployed contract.
func NewFeeSplitterTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeSplitterTransactor, error) {
	contract, err := bindFeeSplitter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeSplitterTransactor{contract: contract}, nil
}

// NewFeeSplitterFilterer creates a new log filterer instance of FeeSplitter, bound to a specific deployed contract.
func NewFeeSplitterFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeSplitterFilterer, error) {
	contract, err := bindFeeSplitter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeSplitterFilterer{contract: contract}, nil
}

// bindFeeSplitter binds a generic wrapper to an already deployed contract.
func bindFeeSplitter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeSplitterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeSplitter *FeeSplitterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeSplitter.Contract.FeeSplitterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeSplitter *FeeSplitterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeSplitter.Contract.FeeSplitterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeSplitter *FeeSplitterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeSplitter.Contract.FeeSplitterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeSplitter *FeeSplitterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeSplitter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeSplitter *FeeSplitterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeSplitter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeSplitter *FeeSplitterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeSplitter.Contract.contract.Transact(opts, method, params...)
}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterCaller) BPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterSession) BPS() (*big.Int, error) {
	return _FeeSplitter.Contract.BPS(&_FeeSplitter.CallOpts)
}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterCallerSession) BPS() (*big.Int, error) {
	return _FeeSplitter.Contract.BPS(&_FeeSplitter.CallOpts)
}

// CALLERROLE is a free data retrieval call binding the contract method 0x774237fc.
//
// Solidity: function CALLER_ROLE() view returns(bytes32)
func (_FeeSplitter *FeeSplitterCaller) CALLERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "CALLER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CALLERROLE is a free data retrieval call binding the contract method 0x774237fc.
//
// Solidity: function CALLER_ROLE() view returns(bytes32)
func (_FeeSplitter *FeeSplitterSession) CALLERROLE() ([32]byte, error) {
	return _FeeSplitter.Contract.CALLERROLE(&_FeeSplitter.CallOpts)
}

// CALLERROLE is a free data retrieval call binding the contract method 0x774237fc.
//
// Solidity: function CALLER_ROLE() view returns(bytes32)
func (_FeeSplitter *FeeSplitterCallerSession) CALLERROLE() ([32]byte, error) {
	return _FeeSplitter.Contract.CALLERROLE(&_FeeSplitter.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_FeeSplitter *FeeSplitterCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_FeeSplitter *FeeSplitterSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _FeeSplitter.Contract.DEFAULTADMINROLE(&_FeeSplitter.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_FeeSplitter *FeeSplitterCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _FeeSplitter.Contract.DEFAULTADMINROLE(&_FeeSplitter.CallOpts)
}

// SCHEDULEABUYBACKBPS is a free data retrieval call binding the contract method 0x80ca6fd5.
//
// Solidity: function SCHEDULE_A_BUYBACK_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterCaller) SCHEDULEABUYBACKBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "SCHEDULE_A_BUYBACK_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SCHEDULEABUYBACKBPS is a free data retrieval call binding the contract method 0x80ca6fd5.
//
// Solidity: function SCHEDULE_A_BUYBACK_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterSession) SCHEDULEABUYBACKBPS() (*big.Int, error) {
	return _FeeSplitter.Contract.SCHEDULEABUYBACKBPS(&_FeeSplitter.CallOpts)
}

// SCHEDULEABUYBACKBPS is a free data retrieval call binding the contract method 0x80ca6fd5.
//
// Solidity: function SCHEDULE_A_BUYBACK_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterCallerSession) SCHEDULEABUYBACKBPS() (*big.Int, error) {
	return _FeeSplitter.Contract.SCHEDULEABUYBACKBPS(&_FeeSplitter.CallOpts)
}

// SCHEDULEATREASURYBPS is a free data retrieval call binding the contract method 0x8d623d05.
//
// Solidity: function SCHEDULE_A_TREASURY_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterCaller) SCHEDULEATREASURYBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "SCHEDULE_A_TREASURY_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SCHEDULEATREASURYBPS is a free data retrieval call binding the contract method 0x8d623d05.
//
// Solidity: function SCHEDULE_A_TREASURY_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterSession) SCHEDULEATREASURYBPS() (*big.Int, error) {
	return _FeeSplitter.Contract.SCHEDULEATREASURYBPS(&_FeeSplitter.CallOpts)
}

// SCHEDULEATREASURYBPS is a free data retrieval call binding the contract method 0x8d623d05.
//
// Solidity: function SCHEDULE_A_TREASURY_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterCallerSession) SCHEDULEATREASURYBPS() (*big.Int, error) {
	return _FeeSplitter.Contract.SCHEDULEATREASURYBPS(&_FeeSplitter.CallOpts)
}

// SCHEDULEBBUYBACKBPS is a free data retrieval call binding the contract method 0x50a43c09.
//
// Solidity: function SCHEDULE_B_BUYBACK_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterCaller) SCHEDULEBBUYBACKBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "SCHEDULE_B_BUYBACK_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SCHEDULEBBUYBACKBPS is a free data retrieval call binding the contract method 0x50a43c09.
//
// Solidity: function SCHEDULE_B_BUYBACK_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterSession) SCHEDULEBBUYBACKBPS() (*big.Int, error) {
	return _FeeSplitter.Contract.SCHEDULEBBUYBACKBPS(&_FeeSplitter.CallOpts)
}

// SCHEDULEBBUYBACKBPS is a free data retrieval call binding the contract method 0x50a43c09.
//
// Solidity: function SCHEDULE_B_BUYBACK_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterCallerSession) SCHEDULEBBUYBACKBPS() (*big.Int, error) {
	return _FeeSplitter.Contract.SCHEDULEBBUYBACKBPS(&_FeeSplitter.CallOpts)
}

// SCHEDULEBTREASURYBPS is a free data retrieval call binding the contract method 0x011275fd.
//
// Solidity: function SCHEDULE_B_TREASURY_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterCaller) SCHEDULEBTREASURYBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "SCHEDULE_B_TREASURY_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SCHEDULEBTREASURYBPS is a free data retrieval call binding the contract method 0x011275fd.
//
// Solidity: function SCHEDULE_B_TREASURY_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterSession) SCHEDULEBTREASURYBPS() (*big.Int, error) {
	return _FeeSplitter.Contract.SCHEDULEBTREASURYBPS(&_FeeSplitter.CallOpts)
}

// SCHEDULEBTREASURYBPS is a free data retrieval call binding the contract method 0x011275fd.
//
// Solidity: function SCHEDULE_B_TREASURY_BPS() view returns(uint256)
func (_FeeSplitter *FeeSplitterCallerSession) SCHEDULEBTREASURYBPS() (*big.Int, error) {
	return _FeeSplitter.Contract.SCHEDULEBTREASURYBPS(&_FeeSplitter.CallOpts)
}

// BuybackBurner is a free data retrieval call binding the contract method 0x82d14fd9.
//
// Solidity: function buybackBurner() view returns(address)
func (_FeeSplitter *FeeSplitterCaller) BuybackBurner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "buybackBurner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BuybackBurner is a free data retrieval call binding the contract method 0x82d14fd9.
//
// Solidity: function buybackBurner() view returns(address)
func (_FeeSplitter *FeeSplitterSession) BuybackBurner() (common.Address, error) {
	return _FeeSplitter.Contract.BuybackBurner(&_FeeSplitter.CallOpts)
}

// BuybackBurner is a free data retrieval call binding the contract method 0x82d14fd9.
//
// Solidity: function buybackBurner() view returns(address)
func (_FeeSplitter *FeeSplitterCallerSession) BuybackBurner() (common.Address, error) {
	return _FeeSplitter.Contract.BuybackBurner(&_FeeSplitter.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_FeeSplitter *FeeSplitterCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_FeeSplitter *FeeSplitterSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _FeeSplitter.Contract.GetRoleAdmin(&_FeeSplitter.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_FeeSplitter *FeeSplitterCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _FeeSplitter.Contract.GetRoleAdmin(&_FeeSplitter.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_FeeSplitter *FeeSplitterCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_FeeSplitter *FeeSplitterSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _FeeSplitter.Contract.HasRole(&_FeeSplitter.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_FeeSplitter *FeeSplitterCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _FeeSplitter.Contract.HasRole(&_FeeSplitter.CallOpts, role, account)
}

// PreviewSplit is a free data retrieval call binding the contract method 0x345290f1.
//
// Solidity: function previewSplit(uint256 amount, uint8 schedule) pure returns(uint256 primaryAmt, uint256 treasuryAmt, uint256 buybackAmt)
func (_FeeSplitter *FeeSplitterCaller) PreviewSplit(opts *bind.CallOpts, amount *big.Int, schedule uint8) (struct {
	PrimaryAmt  *big.Int
	TreasuryAmt *big.Int
	BuybackAmt  *big.Int
}, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "previewSplit", amount, schedule)

	outstruct := new(struct {
		PrimaryAmt  *big.Int
		TreasuryAmt *big.Int
		BuybackAmt  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PrimaryAmt = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TreasuryAmt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BuybackAmt = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// PreviewSplit is a free data retrieval call binding the contract method 0x345290f1.
//
// Solidity: function previewSplit(uint256 amount, uint8 schedule) pure returns(uint256 primaryAmt, uint256 treasuryAmt, uint256 buybackAmt)
func (_FeeSplitter *FeeSplitterSession) PreviewSplit(amount *big.Int, schedule uint8) (struct {
	PrimaryAmt  *big.Int
	TreasuryAmt *big.Int
	BuybackAmt  *big.Int
}, error) {
	return _FeeSplitter.Contract.PreviewSplit(&_FeeSplitter.CallOpts, amount, schedule)
}

// PreviewSplit is a free data retrieval call binding the contract method 0x345290f1.
//
// Solidity: function previewSplit(uint256 amount, uint8 schedule) pure returns(uint256 primaryAmt, uint256 treasuryAmt, uint256 buybackAmt)
func (_FeeSplitter *FeeSplitterCallerSession) PreviewSplit(amount *big.Int, schedule uint8) (struct {
	PrimaryAmt  *big.Int
	TreasuryAmt *big.Int
	BuybackAmt  *big.Int
}, error) {
	return _FeeSplitter.Contract.PreviewSplit(&_FeeSplitter.CallOpts, amount, schedule)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_FeeSplitter *FeeSplitterCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_FeeSplitter *FeeSplitterSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FeeSplitter.Contract.SupportsInterface(&_FeeSplitter.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_FeeSplitter *FeeSplitterCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _FeeSplitter.Contract.SupportsInterface(&_FeeSplitter.CallOpts, interfaceId)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_FeeSplitter *FeeSplitterCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeSplitter.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_FeeSplitter *FeeSplitterSession) Treasury() (common.Address, error) {
	return _FeeSplitter.Contract.Treasury(&_FeeSplitter.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_FeeSplitter *FeeSplitterCallerSession) Treasury() (common.Address, error) {
	return _FeeSplitter.Contract.Treasury(&_FeeSplitter.CallOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_FeeSplitter *FeeSplitterTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FeeSplitter.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_FeeSplitter *FeeSplitterSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FeeSplitter.Contract.GrantRole(&_FeeSplitter.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_FeeSplitter *FeeSplitterTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FeeSplitter.Contract.GrantRole(&_FeeSplitter.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_FeeSplitter *FeeSplitterTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _FeeSplitter.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_FeeSplitter *FeeSplitterSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _FeeSplitter.Contract.RenounceRole(&_FeeSplitter.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_FeeSplitter *FeeSplitterTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _FeeSplitter.Contract.RenounceRole(&_FeeSplitter.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_FeeSplitter *FeeSplitterTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FeeSplitter.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_FeeSplitter *FeeSplitterSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FeeSplitter.Contract.RevokeRole(&_FeeSplitter.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_FeeSplitter *FeeSplitterTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _FeeSplitter.Contract.RevokeRole(&_FeeSplitter.TransactOpts, role, account)
}

// SetBuybackBurner is a paid mutator transaction binding the contract method 0x344e74b2.
//
// Solidity: function setBuybackBurner(address newBurner) returns()
func (_FeeSplitter *FeeSplitterTransactor) SetBuybackBurner(opts *bind.TransactOpts, newBurner common.Address) (*types.Transaction, error) {
	return _FeeSplitter.contract.Transact(opts, "setBuybackBurner", newBurner)
}

// SetBuybackBurner is a paid mutator transaction binding the contract method 0x344e74b2.
//
// Solidity: function setBuybackBurner(address newBurner) returns()
func (_FeeSplitter *FeeSplitterSession) SetBuybackBurner(newBurner common.Address) (*types.Transaction, error) {
	return _FeeSplitter.Contract.SetBuybackBurner(&_FeeSplitter.TransactOpts, newBurner)
}

// SetBuybackBurner is a paid mutator transaction binding the contract method 0x344e74b2.
//
// Solidity: function setBuybackBurner(address newBurner) returns()
func (_FeeSplitter *FeeSplitterTransactorSession) SetBuybackBurner(newBurner common.Address) (*types.Transaction, error) {
	return _FeeSplitter.Contract.SetBuybackBurner(&_FeeSplitter.TransactOpts, newBurner)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address newTreasury) returns()
func (_FeeSplitter *FeeSplitterTransactor) SetTreasury(opts *bind.TransactOpts, newTreasury common.Address) (*types.Transaction, error) {
	return _FeeSplitter.contract.Transact(opts, "setTreasury", newTreasury)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address newTreasury) returns()
func (_FeeSplitter *FeeSplitterSession) SetTreasury(newTreasury common.Address) (*types.Transaction, error) {
	return _FeeSplitter.Contract.SetTreasury(&_FeeSplitter.TransactOpts, newTreasury)
}

// SetTreasury is a paid mutator transaction binding the contract method 0xf0f44260.
//
// Solidity: function setTreasury(address newTreasury) returns()
func (_FeeSplitter *FeeSplitterTransactorSession) SetTreasury(newTreasury common.Address) (*types.Transaction, error) {
	return _FeeSplitter.Contract.SetTreasury(&_FeeSplitter.TransactOpts, newTreasury)
}

// Split is a paid mutator transaction binding the contract method 0x0f41ff3d.
//
// Solidity: function split(address token, uint256 amount, address primary, uint8 schedule) returns()
func (_FeeSplitter *FeeSplitterTransactor) Split(opts *bind.TransactOpts, token common.Address, amount *big.Int, primary common.Address, schedule uint8) (*types.Transaction, error) {
	return _FeeSplitter.contract.Transact(opts, "split", token, amount, primary, schedule)
}

// Split is a paid mutator transaction binding the contract method 0x0f41ff3d.
//
// Solidity: function split(address token, uint256 amount, address primary, uint8 schedule) returns()
func (_FeeSplitter *FeeSplitterSession) Split(token common.Address, amount *big.Int, primary common.Address, schedule uint8) (*types.Transaction, error) {
	return _FeeSplitter.Contract.Split(&_FeeSplitter.TransactOpts, token, amount, primary, schedule)
}

// Split is a paid mutator transaction binding the contract method 0x0f41ff3d.
//
// Solidity: function split(address token, uint256 amount, address primary, uint8 schedule) returns()
func (_FeeSplitter *FeeSplitterTransactorSession) Split(token common.Address, amount *big.Int, primary common.Address, schedule uint8) (*types.Transaction, error) {
	return _FeeSplitter.Contract.Split(&_FeeSplitter.TransactOpts, token, amount, primary, schedule)
}

// FeeSplitterBuybackBurnerUpdatedIterator is returned from FilterBuybackBurnerUpdated and is used to iterate over the raw logs and unpacked data for BuybackBurnerUpdated events raised by the FeeSplitter contract.
type FeeSplitterBuybackBurnerUpdatedIterator struct {
	Event *FeeSplitterBuybackBurnerUpdated // Event containing the contract specifics and raw log

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
func (it *FeeSplitterBuybackBurnerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeSplitterBuybackBurnerUpdated)
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
		it.Event = new(FeeSplitterBuybackBurnerUpdated)
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
func (it *FeeSplitterBuybackBurnerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeSplitterBuybackBurnerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeSplitterBuybackBurnerUpdated represents a BuybackBurnerUpdated event raised by the FeeSplitter contract.
type FeeSplitterBuybackBurnerUpdated struct {
	OldBurner common.Address
	NewBurner common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBuybackBurnerUpdated is a free log retrieval operation binding the contract event 0x2413e4696976dabb94237cb746c767cd0118c63bef644fdd08141569bee6331c.
//
// Solidity: event BuybackBurnerUpdated(address indexed oldBurner, address indexed newBurner)
func (_FeeSplitter *FeeSplitterFilterer) FilterBuybackBurnerUpdated(opts *bind.FilterOpts, oldBurner []common.Address, newBurner []common.Address) (*FeeSplitterBuybackBurnerUpdatedIterator, error) {

	var oldBurnerRule []interface{}
	for _, oldBurnerItem := range oldBurner {
		oldBurnerRule = append(oldBurnerRule, oldBurnerItem)
	}
	var newBurnerRule []interface{}
	for _, newBurnerItem := range newBurner {
		newBurnerRule = append(newBurnerRule, newBurnerItem)
	}

	logs, sub, err := _FeeSplitter.contract.FilterLogs(opts, "BuybackBurnerUpdated", oldBurnerRule, newBurnerRule)
	if err != nil {
		return nil, err
	}
	return &FeeSplitterBuybackBurnerUpdatedIterator{contract: _FeeSplitter.contract, event: "BuybackBurnerUpdated", logs: logs, sub: sub}, nil
}

// WatchBuybackBurnerUpdated is a free log subscription operation binding the contract event 0x2413e4696976dabb94237cb746c767cd0118c63bef644fdd08141569bee6331c.
//
// Solidity: event BuybackBurnerUpdated(address indexed oldBurner, address indexed newBurner)
func (_FeeSplitter *FeeSplitterFilterer) WatchBuybackBurnerUpdated(opts *bind.WatchOpts, sink chan<- *FeeSplitterBuybackBurnerUpdated, oldBurner []common.Address, newBurner []common.Address) (event.Subscription, error) {

	var oldBurnerRule []interface{}
	for _, oldBurnerItem := range oldBurner {
		oldBurnerRule = append(oldBurnerRule, oldBurnerItem)
	}
	var newBurnerRule []interface{}
	for _, newBurnerItem := range newBurner {
		newBurnerRule = append(newBurnerRule, newBurnerItem)
	}

	logs, sub, err := _FeeSplitter.contract.WatchLogs(opts, "BuybackBurnerUpdated", oldBurnerRule, newBurnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeSplitterBuybackBurnerUpdated)
				if err := _FeeSplitter.contract.UnpackLog(event, "BuybackBurnerUpdated", log); err != nil {
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

// ParseBuybackBurnerUpdated is a log parse operation binding the contract event 0x2413e4696976dabb94237cb746c767cd0118c63bef644fdd08141569bee6331c.
//
// Solidity: event BuybackBurnerUpdated(address indexed oldBurner, address indexed newBurner)
func (_FeeSplitter *FeeSplitterFilterer) ParseBuybackBurnerUpdated(log types.Log) (*FeeSplitterBuybackBurnerUpdated, error) {
	event := new(FeeSplitterBuybackBurnerUpdated)
	if err := _FeeSplitter.contract.UnpackLog(event, "BuybackBurnerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeSplitterFeeSplitIterator is returned from FilterFeeSplit and is used to iterate over the raw logs and unpacked data for FeeSplit events raised by the FeeSplitter contract.
type FeeSplitterFeeSplitIterator struct {
	Event *FeeSplitterFeeSplit // Event containing the contract specifics and raw log

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
func (it *FeeSplitterFeeSplitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeSplitterFeeSplit)
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
		it.Event = new(FeeSplitterFeeSplit)
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
func (it *FeeSplitterFeeSplitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeSplitterFeeSplitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeSplitterFeeSplit represents a FeeSplit event raised by the FeeSplitter contract.
type FeeSplitterFeeSplit struct {
	Token          common.Address
	Primary        common.Address
	Schedule       uint8
	PrimaryAmount  *big.Int
	TreasuryAmount *big.Int
	BuybackAmount  *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterFeeSplit is a free log retrieval operation binding the contract event 0x396f3c0c854fc40380f8870c3e6c468a076d3510a63ecbf8c82c38fa647c098c.
//
// Solidity: event FeeSplit(address indexed token, address indexed primary, uint8 indexed schedule, uint256 primaryAmount, uint256 treasuryAmount, uint256 buybackAmount)
func (_FeeSplitter *FeeSplitterFilterer) FilterFeeSplit(opts *bind.FilterOpts, token []common.Address, primary []common.Address, schedule []uint8) (*FeeSplitterFeeSplitIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var primaryRule []interface{}
	for _, primaryItem := range primary {
		primaryRule = append(primaryRule, primaryItem)
	}
	var scheduleRule []interface{}
	for _, scheduleItem := range schedule {
		scheduleRule = append(scheduleRule, scheduleItem)
	}

	logs, sub, err := _FeeSplitter.contract.FilterLogs(opts, "FeeSplit", tokenRule, primaryRule, scheduleRule)
	if err != nil {
		return nil, err
	}
	return &FeeSplitterFeeSplitIterator{contract: _FeeSplitter.contract, event: "FeeSplit", logs: logs, sub: sub}, nil
}

// WatchFeeSplit is a free log subscription operation binding the contract event 0x396f3c0c854fc40380f8870c3e6c468a076d3510a63ecbf8c82c38fa647c098c.
//
// Solidity: event FeeSplit(address indexed token, address indexed primary, uint8 indexed schedule, uint256 primaryAmount, uint256 treasuryAmount, uint256 buybackAmount)
func (_FeeSplitter *FeeSplitterFilterer) WatchFeeSplit(opts *bind.WatchOpts, sink chan<- *FeeSplitterFeeSplit, token []common.Address, primary []common.Address, schedule []uint8) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var primaryRule []interface{}
	for _, primaryItem := range primary {
		primaryRule = append(primaryRule, primaryItem)
	}
	var scheduleRule []interface{}
	for _, scheduleItem := range schedule {
		scheduleRule = append(scheduleRule, scheduleItem)
	}

	logs, sub, err := _FeeSplitter.contract.WatchLogs(opts, "FeeSplit", tokenRule, primaryRule, scheduleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeSplitterFeeSplit)
				if err := _FeeSplitter.contract.UnpackLog(event, "FeeSplit", log); err != nil {
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

// ParseFeeSplit is a log parse operation binding the contract event 0x396f3c0c854fc40380f8870c3e6c468a076d3510a63ecbf8c82c38fa647c098c.
//
// Solidity: event FeeSplit(address indexed token, address indexed primary, uint8 indexed schedule, uint256 primaryAmount, uint256 treasuryAmount, uint256 buybackAmount)
func (_FeeSplitter *FeeSplitterFilterer) ParseFeeSplit(log types.Log) (*FeeSplitterFeeSplit, error) {
	event := new(FeeSplitterFeeSplit)
	if err := _FeeSplitter.contract.UnpackLog(event, "FeeSplit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeSplitterRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the FeeSplitter contract.
type FeeSplitterRoleAdminChangedIterator struct {
	Event *FeeSplitterRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *FeeSplitterRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeSplitterRoleAdminChanged)
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
		it.Event = new(FeeSplitterRoleAdminChanged)
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
func (it *FeeSplitterRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeSplitterRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeSplitterRoleAdminChanged represents a RoleAdminChanged event raised by the FeeSplitter contract.
type FeeSplitterRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_FeeSplitter *FeeSplitterFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*FeeSplitterRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _FeeSplitter.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &FeeSplitterRoleAdminChangedIterator{contract: _FeeSplitter.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_FeeSplitter *FeeSplitterFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *FeeSplitterRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _FeeSplitter.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeSplitterRoleAdminChanged)
				if err := _FeeSplitter.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_FeeSplitter *FeeSplitterFilterer) ParseRoleAdminChanged(log types.Log) (*FeeSplitterRoleAdminChanged, error) {
	event := new(FeeSplitterRoleAdminChanged)
	if err := _FeeSplitter.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeSplitterRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the FeeSplitter contract.
type FeeSplitterRoleGrantedIterator struct {
	Event *FeeSplitterRoleGranted // Event containing the contract specifics and raw log

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
func (it *FeeSplitterRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeSplitterRoleGranted)
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
		it.Event = new(FeeSplitterRoleGranted)
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
func (it *FeeSplitterRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeSplitterRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeSplitterRoleGranted represents a RoleGranted event raised by the FeeSplitter contract.
type FeeSplitterRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_FeeSplitter *FeeSplitterFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*FeeSplitterRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _FeeSplitter.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &FeeSplitterRoleGrantedIterator{contract: _FeeSplitter.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_FeeSplitter *FeeSplitterFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *FeeSplitterRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _FeeSplitter.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeSplitterRoleGranted)
				if err := _FeeSplitter.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_FeeSplitter *FeeSplitterFilterer) ParseRoleGranted(log types.Log) (*FeeSplitterRoleGranted, error) {
	event := new(FeeSplitterRoleGranted)
	if err := _FeeSplitter.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeSplitterRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the FeeSplitter contract.
type FeeSplitterRoleRevokedIterator struct {
	Event *FeeSplitterRoleRevoked // Event containing the contract specifics and raw log

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
func (it *FeeSplitterRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeSplitterRoleRevoked)
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
		it.Event = new(FeeSplitterRoleRevoked)
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
func (it *FeeSplitterRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeSplitterRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeSplitterRoleRevoked represents a RoleRevoked event raised by the FeeSplitter contract.
type FeeSplitterRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_FeeSplitter *FeeSplitterFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*FeeSplitterRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _FeeSplitter.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &FeeSplitterRoleRevokedIterator{contract: _FeeSplitter.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_FeeSplitter *FeeSplitterFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *FeeSplitterRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _FeeSplitter.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeSplitterRoleRevoked)
				if err := _FeeSplitter.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_FeeSplitter *FeeSplitterFilterer) ParseRoleRevoked(log types.Log) (*FeeSplitterRoleRevoked, error) {
	event := new(FeeSplitterRoleRevoked)
	if err := _FeeSplitter.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeSplitterTreasuryUpdatedIterator is returned from FilterTreasuryUpdated and is used to iterate over the raw logs and unpacked data for TreasuryUpdated events raised by the FeeSplitter contract.
type FeeSplitterTreasuryUpdatedIterator struct {
	Event *FeeSplitterTreasuryUpdated // Event containing the contract specifics and raw log

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
func (it *FeeSplitterTreasuryUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeSplitterTreasuryUpdated)
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
		it.Event = new(FeeSplitterTreasuryUpdated)
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
func (it *FeeSplitterTreasuryUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeSplitterTreasuryUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeSplitterTreasuryUpdated represents a TreasuryUpdated event raised by the FeeSplitter contract.
type FeeSplitterTreasuryUpdated struct {
	OldTreasury common.Address
	NewTreasury common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTreasuryUpdated is a free log retrieval operation binding the contract event 0x4ab5be82436d353e61ca18726e984e561f5c1cc7c6d38b29d2553c790434705a.
//
// Solidity: event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury)
func (_FeeSplitter *FeeSplitterFilterer) FilterTreasuryUpdated(opts *bind.FilterOpts, oldTreasury []common.Address, newTreasury []common.Address) (*FeeSplitterTreasuryUpdatedIterator, error) {

	var oldTreasuryRule []interface{}
	for _, oldTreasuryItem := range oldTreasury {
		oldTreasuryRule = append(oldTreasuryRule, oldTreasuryItem)
	}
	var newTreasuryRule []interface{}
	for _, newTreasuryItem := range newTreasury {
		newTreasuryRule = append(newTreasuryRule, newTreasuryItem)
	}

	logs, sub, err := _FeeSplitter.contract.FilterLogs(opts, "TreasuryUpdated", oldTreasuryRule, newTreasuryRule)
	if err != nil {
		return nil, err
	}
	return &FeeSplitterTreasuryUpdatedIterator{contract: _FeeSplitter.contract, event: "TreasuryUpdated", logs: logs, sub: sub}, nil
}

// WatchTreasuryUpdated is a free log subscription operation binding the contract event 0x4ab5be82436d353e61ca18726e984e561f5c1cc7c6d38b29d2553c790434705a.
//
// Solidity: event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury)
func (_FeeSplitter *FeeSplitterFilterer) WatchTreasuryUpdated(opts *bind.WatchOpts, sink chan<- *FeeSplitterTreasuryUpdated, oldTreasury []common.Address, newTreasury []common.Address) (event.Subscription, error) {

	var oldTreasuryRule []interface{}
	for _, oldTreasuryItem := range oldTreasury {
		oldTreasuryRule = append(oldTreasuryRule, oldTreasuryItem)
	}
	var newTreasuryRule []interface{}
	for _, newTreasuryItem := range newTreasury {
		newTreasuryRule = append(newTreasuryRule, newTreasuryItem)
	}

	logs, sub, err := _FeeSplitter.contract.WatchLogs(opts, "TreasuryUpdated", oldTreasuryRule, newTreasuryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeSplitterTreasuryUpdated)
				if err := _FeeSplitter.contract.UnpackLog(event, "TreasuryUpdated", log); err != nil {
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

// ParseTreasuryUpdated is a log parse operation binding the contract event 0x4ab5be82436d353e61ca18726e984e561f5c1cc7c6d38b29d2553c790434705a.
//
// Solidity: event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury)
func (_FeeSplitter *FeeSplitterFilterer) ParseTreasuryUpdated(log types.Log) (*FeeSplitterTreasuryUpdated, error) {
	event := new(FeeSplitterTreasuryUpdated)
	if err := _FeeSplitter.contract.UnpackLog(event, "TreasuryUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
