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

// BuybackBurnerMetaData contains all meta data concerning the BuybackBurner contract.
var BuybackBurnerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_router\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_paymentToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_titu\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_swapPath\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_minTituOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_swapDeadlineBuffer\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"EXECUTOR_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buybackAndBurn\",\"inputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountOutMin\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSwapPath\",\"inputs\":[],\"outputs\":[{\"name\":\"path\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minTituOut\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paymentToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rescueTokens\",\"inputs\":[{\"name\":\"tokenAddr\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"router\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUniswapV2Router02\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setMinTituOut\",\"inputs\":[{\"name\":\"newMin\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRouter\",\"inputs\":[{\"name\":\"newRouter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSwapPath\",\"inputs\":[{\"name\":\"path\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"swapDeadlineBuffer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"swapPath\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"titu\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20Burnable\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"BuybackAndBurn\",\"inputs\":[{\"name\":\"executor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"paymentToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amountIn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"tituBurned\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MinTituOutUpdated\",\"inputs\":[{\"name\":\"oldMin\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newMin\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouterUpdated\",\"inputs\":[{\"name\":\"oldRouter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newRouter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SwapPathUpdated\",\"inputs\":[{\"name\":\"path\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"InvalidSwapPath\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"RescuePaymentTokenForbidden\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// BuybackBurnerABI is the input ABI used to generate the binding from.
// Deprecated: Use BuybackBurnerMetaData.ABI instead.
var BuybackBurnerABI = BuybackBurnerMetaData.ABI

// BuybackBurner is an auto generated Go binding around an Ethereum contract.
type BuybackBurner struct {
	BuybackBurnerCaller     // Read-only binding to the contract
	BuybackBurnerTransactor // Write-only binding to the contract
	BuybackBurnerFilterer   // Log filterer for contract events
}

// BuybackBurnerCaller is an auto generated read-only Go binding around an Ethereum contract.
type BuybackBurnerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BuybackBurnerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BuybackBurnerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BuybackBurnerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BuybackBurnerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BuybackBurnerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BuybackBurnerSession struct {
	Contract     *BuybackBurner    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BuybackBurnerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BuybackBurnerCallerSession struct {
	Contract *BuybackBurnerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// BuybackBurnerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BuybackBurnerTransactorSession struct {
	Contract     *BuybackBurnerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// BuybackBurnerRaw is an auto generated low-level Go binding around an Ethereum contract.
type BuybackBurnerRaw struct {
	Contract *BuybackBurner // Generic contract binding to access the raw methods on
}

// BuybackBurnerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BuybackBurnerCallerRaw struct {
	Contract *BuybackBurnerCaller // Generic read-only contract binding to access the raw methods on
}

// BuybackBurnerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BuybackBurnerTransactorRaw struct {
	Contract *BuybackBurnerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBuybackBurner creates a new instance of BuybackBurner, bound to a specific deployed contract.
func NewBuybackBurner(address common.Address, backend bind.ContractBackend) (*BuybackBurner, error) {
	contract, err := bindBuybackBurner(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BuybackBurner{BuybackBurnerCaller: BuybackBurnerCaller{contract: contract}, BuybackBurnerTransactor: BuybackBurnerTransactor{contract: contract}, BuybackBurnerFilterer: BuybackBurnerFilterer{contract: contract}}, nil
}

// NewBuybackBurnerCaller creates a new read-only instance of BuybackBurner, bound to a specific deployed contract.
func NewBuybackBurnerCaller(address common.Address, caller bind.ContractCaller) (*BuybackBurnerCaller, error) {
	contract, err := bindBuybackBurner(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BuybackBurnerCaller{contract: contract}, nil
}

// NewBuybackBurnerTransactor creates a new write-only instance of BuybackBurner, bound to a specific deployed contract.
func NewBuybackBurnerTransactor(address common.Address, transactor bind.ContractTransactor) (*BuybackBurnerTransactor, error) {
	contract, err := bindBuybackBurner(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BuybackBurnerTransactor{contract: contract}, nil
}

// NewBuybackBurnerFilterer creates a new log filterer instance of BuybackBurner, bound to a specific deployed contract.
func NewBuybackBurnerFilterer(address common.Address, filterer bind.ContractFilterer) (*BuybackBurnerFilterer, error) {
	contract, err := bindBuybackBurner(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BuybackBurnerFilterer{contract: contract}, nil
}

// bindBuybackBurner binds a generic wrapper to an already deployed contract.
func bindBuybackBurner(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BuybackBurnerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BuybackBurner *BuybackBurnerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BuybackBurner.Contract.BuybackBurnerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BuybackBurner *BuybackBurnerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BuybackBurner.Contract.BuybackBurnerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BuybackBurner *BuybackBurnerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BuybackBurner.Contract.BuybackBurnerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BuybackBurner *BuybackBurnerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BuybackBurner.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BuybackBurner *BuybackBurnerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BuybackBurner.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BuybackBurner *BuybackBurnerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BuybackBurner.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BuybackBurner *BuybackBurnerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BuybackBurner *BuybackBurnerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _BuybackBurner.Contract.DEFAULTADMINROLE(&_BuybackBurner.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_BuybackBurner *BuybackBurnerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _BuybackBurner.Contract.DEFAULTADMINROLE(&_BuybackBurner.CallOpts)
}

// EXECUTORROLE is a free data retrieval call binding the contract method 0x07bd0265.
//
// Solidity: function EXECUTOR_ROLE() view returns(bytes32)
func (_BuybackBurner *BuybackBurnerCaller) EXECUTORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "EXECUTOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EXECUTORROLE is a free data retrieval call binding the contract method 0x07bd0265.
//
// Solidity: function EXECUTOR_ROLE() view returns(bytes32)
func (_BuybackBurner *BuybackBurnerSession) EXECUTORROLE() ([32]byte, error) {
	return _BuybackBurner.Contract.EXECUTORROLE(&_BuybackBurner.CallOpts)
}

// EXECUTORROLE is a free data retrieval call binding the contract method 0x07bd0265.
//
// Solidity: function EXECUTOR_ROLE() view returns(bytes32)
func (_BuybackBurner *BuybackBurnerCallerSession) EXECUTORROLE() ([32]byte, error) {
	return _BuybackBurner.Contract.EXECUTORROLE(&_BuybackBurner.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BuybackBurner *BuybackBurnerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BuybackBurner *BuybackBurnerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _BuybackBurner.Contract.GetRoleAdmin(&_BuybackBurner.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_BuybackBurner *BuybackBurnerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _BuybackBurner.Contract.GetRoleAdmin(&_BuybackBurner.CallOpts, role)
}

// GetSwapPath is a free data retrieval call binding the contract method 0x8ffb62d2.
//
// Solidity: function getSwapPath() view returns(address[] path)
func (_BuybackBurner *BuybackBurnerCaller) GetSwapPath(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "getSwapPath")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetSwapPath is a free data retrieval call binding the contract method 0x8ffb62d2.
//
// Solidity: function getSwapPath() view returns(address[] path)
func (_BuybackBurner *BuybackBurnerSession) GetSwapPath() ([]common.Address, error) {
	return _BuybackBurner.Contract.GetSwapPath(&_BuybackBurner.CallOpts)
}

// GetSwapPath is a free data retrieval call binding the contract method 0x8ffb62d2.
//
// Solidity: function getSwapPath() view returns(address[] path)
func (_BuybackBurner *BuybackBurnerCallerSession) GetSwapPath() ([]common.Address, error) {
	return _BuybackBurner.Contract.GetSwapPath(&_BuybackBurner.CallOpts)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BuybackBurner *BuybackBurnerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BuybackBurner *BuybackBurnerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BuybackBurner.Contract.HasRole(&_BuybackBurner.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_BuybackBurner *BuybackBurnerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _BuybackBurner.Contract.HasRole(&_BuybackBurner.CallOpts, role, account)
}

// MinTituOut is a free data retrieval call binding the contract method 0xff354d50.
//
// Solidity: function minTituOut() view returns(uint256)
func (_BuybackBurner *BuybackBurnerCaller) MinTituOut(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "minTituOut")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinTituOut is a free data retrieval call binding the contract method 0xff354d50.
//
// Solidity: function minTituOut() view returns(uint256)
func (_BuybackBurner *BuybackBurnerSession) MinTituOut() (*big.Int, error) {
	return _BuybackBurner.Contract.MinTituOut(&_BuybackBurner.CallOpts)
}

// MinTituOut is a free data retrieval call binding the contract method 0xff354d50.
//
// Solidity: function minTituOut() view returns(uint256)
func (_BuybackBurner *BuybackBurnerCallerSession) MinTituOut() (*big.Int, error) {
	return _BuybackBurner.Contract.MinTituOut(&_BuybackBurner.CallOpts)
}

// PaymentToken is a free data retrieval call binding the contract method 0x3013ce29.
//
// Solidity: function paymentToken() view returns(address)
func (_BuybackBurner *BuybackBurnerCaller) PaymentToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "paymentToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PaymentToken is a free data retrieval call binding the contract method 0x3013ce29.
//
// Solidity: function paymentToken() view returns(address)
func (_BuybackBurner *BuybackBurnerSession) PaymentToken() (common.Address, error) {
	return _BuybackBurner.Contract.PaymentToken(&_BuybackBurner.CallOpts)
}

// PaymentToken is a free data retrieval call binding the contract method 0x3013ce29.
//
// Solidity: function paymentToken() view returns(address)
func (_BuybackBurner *BuybackBurnerCallerSession) PaymentToken() (common.Address, error) {
	return _BuybackBurner.Contract.PaymentToken(&_BuybackBurner.CallOpts)
}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_BuybackBurner *BuybackBurnerCaller) Router(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_BuybackBurner *BuybackBurnerSession) Router() (common.Address, error) {
	return _BuybackBurner.Contract.Router(&_BuybackBurner.CallOpts)
}

// Router is a free data retrieval call binding the contract method 0xf887ea40.
//
// Solidity: function router() view returns(address)
func (_BuybackBurner *BuybackBurnerCallerSession) Router() (common.Address, error) {
	return _BuybackBurner.Contract.Router(&_BuybackBurner.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BuybackBurner *BuybackBurnerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BuybackBurner *BuybackBurnerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BuybackBurner.Contract.SupportsInterface(&_BuybackBurner.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_BuybackBurner *BuybackBurnerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _BuybackBurner.Contract.SupportsInterface(&_BuybackBurner.CallOpts, interfaceId)
}

// SwapDeadlineBuffer is a free data retrieval call binding the contract method 0xdec8f8c0.
//
// Solidity: function swapDeadlineBuffer() view returns(uint256)
func (_BuybackBurner *BuybackBurnerCaller) SwapDeadlineBuffer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "swapDeadlineBuffer")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SwapDeadlineBuffer is a free data retrieval call binding the contract method 0xdec8f8c0.
//
// Solidity: function swapDeadlineBuffer() view returns(uint256)
func (_BuybackBurner *BuybackBurnerSession) SwapDeadlineBuffer() (*big.Int, error) {
	return _BuybackBurner.Contract.SwapDeadlineBuffer(&_BuybackBurner.CallOpts)
}

// SwapDeadlineBuffer is a free data retrieval call binding the contract method 0xdec8f8c0.
//
// Solidity: function swapDeadlineBuffer() view returns(uint256)
func (_BuybackBurner *BuybackBurnerCallerSession) SwapDeadlineBuffer() (*big.Int, error) {
	return _BuybackBurner.Contract.SwapDeadlineBuffer(&_BuybackBurner.CallOpts)
}

// SwapPath is a free data retrieval call binding the contract method 0x24cc0766.
//
// Solidity: function swapPath(uint256 ) view returns(address)
func (_BuybackBurner *BuybackBurnerCaller) SwapPath(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "swapPath", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SwapPath is a free data retrieval call binding the contract method 0x24cc0766.
//
// Solidity: function swapPath(uint256 ) view returns(address)
func (_BuybackBurner *BuybackBurnerSession) SwapPath(arg0 *big.Int) (common.Address, error) {
	return _BuybackBurner.Contract.SwapPath(&_BuybackBurner.CallOpts, arg0)
}

// SwapPath is a free data retrieval call binding the contract method 0x24cc0766.
//
// Solidity: function swapPath(uint256 ) view returns(address)
func (_BuybackBurner *BuybackBurnerCallerSession) SwapPath(arg0 *big.Int) (common.Address, error) {
	return _BuybackBurner.Contract.SwapPath(&_BuybackBurner.CallOpts, arg0)
}

// Titu is a free data retrieval call binding the contract method 0x33fcb0f6.
//
// Solidity: function titu() view returns(address)
func (_BuybackBurner *BuybackBurnerCaller) Titu(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BuybackBurner.contract.Call(opts, &out, "titu")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Titu is a free data retrieval call binding the contract method 0x33fcb0f6.
//
// Solidity: function titu() view returns(address)
func (_BuybackBurner *BuybackBurnerSession) Titu() (common.Address, error) {
	return _BuybackBurner.Contract.Titu(&_BuybackBurner.CallOpts)
}

// Titu is a free data retrieval call binding the contract method 0x33fcb0f6.
//
// Solidity: function titu() view returns(address)
func (_BuybackBurner *BuybackBurnerCallerSession) Titu() (common.Address, error) {
	return _BuybackBurner.Contract.Titu(&_BuybackBurner.CallOpts)
}

// BuybackAndBurn is a paid mutator transaction binding the contract method 0x2238eee5.
//
// Solidity: function buybackAndBurn(uint256 amountIn, uint256 amountOutMin) returns()
func (_BuybackBurner *BuybackBurnerTransactor) BuybackAndBurn(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int) (*types.Transaction, error) {
	return _BuybackBurner.contract.Transact(opts, "buybackAndBurn", amountIn, amountOutMin)
}

// BuybackAndBurn is a paid mutator transaction binding the contract method 0x2238eee5.
//
// Solidity: function buybackAndBurn(uint256 amountIn, uint256 amountOutMin) returns()
func (_BuybackBurner *BuybackBurnerSession) BuybackAndBurn(amountIn *big.Int, amountOutMin *big.Int) (*types.Transaction, error) {
	return _BuybackBurner.Contract.BuybackAndBurn(&_BuybackBurner.TransactOpts, amountIn, amountOutMin)
}

// BuybackAndBurn is a paid mutator transaction binding the contract method 0x2238eee5.
//
// Solidity: function buybackAndBurn(uint256 amountIn, uint256 amountOutMin) returns()
func (_BuybackBurner *BuybackBurnerTransactorSession) BuybackAndBurn(amountIn *big.Int, amountOutMin *big.Int) (*types.Transaction, error) {
	return _BuybackBurner.Contract.BuybackAndBurn(&_BuybackBurner.TransactOpts, amountIn, amountOutMin)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BuybackBurner *BuybackBurnerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BuybackBurner.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BuybackBurner *BuybackBurnerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BuybackBurner.Contract.GrantRole(&_BuybackBurner.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_BuybackBurner *BuybackBurnerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BuybackBurner.Contract.GrantRole(&_BuybackBurner.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_BuybackBurner *BuybackBurnerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _BuybackBurner.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_BuybackBurner *BuybackBurnerSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _BuybackBurner.Contract.RenounceRole(&_BuybackBurner.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_BuybackBurner *BuybackBurnerTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _BuybackBurner.Contract.RenounceRole(&_BuybackBurner.TransactOpts, role, callerConfirmation)
}

// RescueTokens is a paid mutator transaction binding the contract method 0xcea9d26f.
//
// Solidity: function rescueTokens(address tokenAddr, address to, uint256 amount) returns()
func (_BuybackBurner *BuybackBurnerTransactor) RescueTokens(opts *bind.TransactOpts, tokenAddr common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BuybackBurner.contract.Transact(opts, "rescueTokens", tokenAddr, to, amount)
}

// RescueTokens is a paid mutator transaction binding the contract method 0xcea9d26f.
//
// Solidity: function rescueTokens(address tokenAddr, address to, uint256 amount) returns()
func (_BuybackBurner *BuybackBurnerSession) RescueTokens(tokenAddr common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BuybackBurner.Contract.RescueTokens(&_BuybackBurner.TransactOpts, tokenAddr, to, amount)
}

// RescueTokens is a paid mutator transaction binding the contract method 0xcea9d26f.
//
// Solidity: function rescueTokens(address tokenAddr, address to, uint256 amount) returns()
func (_BuybackBurner *BuybackBurnerTransactorSession) RescueTokens(tokenAddr common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BuybackBurner.Contract.RescueTokens(&_BuybackBurner.TransactOpts, tokenAddr, to, amount)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BuybackBurner *BuybackBurnerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BuybackBurner.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BuybackBurner *BuybackBurnerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BuybackBurner.Contract.RevokeRole(&_BuybackBurner.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_BuybackBurner *BuybackBurnerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _BuybackBurner.Contract.RevokeRole(&_BuybackBurner.TransactOpts, role, account)
}

// SetMinTituOut is a paid mutator transaction binding the contract method 0x904b4475.
//
// Solidity: function setMinTituOut(uint256 newMin) returns()
func (_BuybackBurner *BuybackBurnerTransactor) SetMinTituOut(opts *bind.TransactOpts, newMin *big.Int) (*types.Transaction, error) {
	return _BuybackBurner.contract.Transact(opts, "setMinTituOut", newMin)
}

// SetMinTituOut is a paid mutator transaction binding the contract method 0x904b4475.
//
// Solidity: function setMinTituOut(uint256 newMin) returns()
func (_BuybackBurner *BuybackBurnerSession) SetMinTituOut(newMin *big.Int) (*types.Transaction, error) {
	return _BuybackBurner.Contract.SetMinTituOut(&_BuybackBurner.TransactOpts, newMin)
}

// SetMinTituOut is a paid mutator transaction binding the contract method 0x904b4475.
//
// Solidity: function setMinTituOut(uint256 newMin) returns()
func (_BuybackBurner *BuybackBurnerTransactorSession) SetMinTituOut(newMin *big.Int) (*types.Transaction, error) {
	return _BuybackBurner.Contract.SetMinTituOut(&_BuybackBurner.TransactOpts, newMin)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_BuybackBurner *BuybackBurnerTransactor) SetRouter(opts *bind.TransactOpts, newRouter common.Address) (*types.Transaction, error) {
	return _BuybackBurner.contract.Transact(opts, "setRouter", newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_BuybackBurner *BuybackBurnerSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BuybackBurner.Contract.SetRouter(&_BuybackBurner.TransactOpts, newRouter)
}

// SetRouter is a paid mutator transaction binding the contract method 0xc0d78655.
//
// Solidity: function setRouter(address newRouter) returns()
func (_BuybackBurner *BuybackBurnerTransactorSession) SetRouter(newRouter common.Address) (*types.Transaction, error) {
	return _BuybackBurner.Contract.SetRouter(&_BuybackBurner.TransactOpts, newRouter)
}

// SetSwapPath is a paid mutator transaction binding the contract method 0x12ddadbc.
//
// Solidity: function setSwapPath(address[] path) returns()
func (_BuybackBurner *BuybackBurnerTransactor) SetSwapPath(opts *bind.TransactOpts, path []common.Address) (*types.Transaction, error) {
	return _BuybackBurner.contract.Transact(opts, "setSwapPath", path)
}

// SetSwapPath is a paid mutator transaction binding the contract method 0x12ddadbc.
//
// Solidity: function setSwapPath(address[] path) returns()
func (_BuybackBurner *BuybackBurnerSession) SetSwapPath(path []common.Address) (*types.Transaction, error) {
	return _BuybackBurner.Contract.SetSwapPath(&_BuybackBurner.TransactOpts, path)
}

// SetSwapPath is a paid mutator transaction binding the contract method 0x12ddadbc.
//
// Solidity: function setSwapPath(address[] path) returns()
func (_BuybackBurner *BuybackBurnerTransactorSession) SetSwapPath(path []common.Address) (*types.Transaction, error) {
	return _BuybackBurner.Contract.SetSwapPath(&_BuybackBurner.TransactOpts, path)
}

// BuybackBurnerBuybackAndBurnIterator is returned from FilterBuybackAndBurn and is used to iterate over the raw logs and unpacked data for BuybackAndBurn events raised by the BuybackBurner contract.
type BuybackBurnerBuybackAndBurnIterator struct {
	Event *BuybackBurnerBuybackAndBurn // Event containing the contract specifics and raw log

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
func (it *BuybackBurnerBuybackAndBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuybackBurnerBuybackAndBurn)
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
		it.Event = new(BuybackBurnerBuybackAndBurn)
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
func (it *BuybackBurnerBuybackAndBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuybackBurnerBuybackAndBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuybackBurnerBuybackAndBurn represents a BuybackAndBurn event raised by the BuybackBurner contract.
type BuybackBurnerBuybackAndBurn struct {
	Executor     common.Address
	PaymentToken common.Address
	AmountIn     *big.Int
	TituBurned   *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterBuybackAndBurn is a free log retrieval operation binding the contract event 0x0f995e1de8b2d83fb2db107b6a0cb5b3a5bb1cca0d54dbcb01dd63f19ffed3bc.
//
// Solidity: event BuybackAndBurn(address indexed executor, address indexed paymentToken, uint256 amountIn, uint256 tituBurned)
func (_BuybackBurner *BuybackBurnerFilterer) FilterBuybackAndBurn(opts *bind.FilterOpts, executor []common.Address, paymentToken []common.Address) (*BuybackBurnerBuybackAndBurnIterator, error) {

	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}
	var paymentTokenRule []interface{}
	for _, paymentTokenItem := range paymentToken {
		paymentTokenRule = append(paymentTokenRule, paymentTokenItem)
	}

	logs, sub, err := _BuybackBurner.contract.FilterLogs(opts, "BuybackAndBurn", executorRule, paymentTokenRule)
	if err != nil {
		return nil, err
	}
	return &BuybackBurnerBuybackAndBurnIterator{contract: _BuybackBurner.contract, event: "BuybackAndBurn", logs: logs, sub: sub}, nil
}

// WatchBuybackAndBurn is a free log subscription operation binding the contract event 0x0f995e1de8b2d83fb2db107b6a0cb5b3a5bb1cca0d54dbcb01dd63f19ffed3bc.
//
// Solidity: event BuybackAndBurn(address indexed executor, address indexed paymentToken, uint256 amountIn, uint256 tituBurned)
func (_BuybackBurner *BuybackBurnerFilterer) WatchBuybackAndBurn(opts *bind.WatchOpts, sink chan<- *BuybackBurnerBuybackAndBurn, executor []common.Address, paymentToken []common.Address) (event.Subscription, error) {

	var executorRule []interface{}
	for _, executorItem := range executor {
		executorRule = append(executorRule, executorItem)
	}
	var paymentTokenRule []interface{}
	for _, paymentTokenItem := range paymentToken {
		paymentTokenRule = append(paymentTokenRule, paymentTokenItem)
	}

	logs, sub, err := _BuybackBurner.contract.WatchLogs(opts, "BuybackAndBurn", executorRule, paymentTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuybackBurnerBuybackAndBurn)
				if err := _BuybackBurner.contract.UnpackLog(event, "BuybackAndBurn", log); err != nil {
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

// ParseBuybackAndBurn is a log parse operation binding the contract event 0x0f995e1de8b2d83fb2db107b6a0cb5b3a5bb1cca0d54dbcb01dd63f19ffed3bc.
//
// Solidity: event BuybackAndBurn(address indexed executor, address indexed paymentToken, uint256 amountIn, uint256 tituBurned)
func (_BuybackBurner *BuybackBurnerFilterer) ParseBuybackAndBurn(log types.Log) (*BuybackBurnerBuybackAndBurn, error) {
	event := new(BuybackBurnerBuybackAndBurn)
	if err := _BuybackBurner.contract.UnpackLog(event, "BuybackAndBurn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BuybackBurnerMinTituOutUpdatedIterator is returned from FilterMinTituOutUpdated and is used to iterate over the raw logs and unpacked data for MinTituOutUpdated events raised by the BuybackBurner contract.
type BuybackBurnerMinTituOutUpdatedIterator struct {
	Event *BuybackBurnerMinTituOutUpdated // Event containing the contract specifics and raw log

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
func (it *BuybackBurnerMinTituOutUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuybackBurnerMinTituOutUpdated)
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
		it.Event = new(BuybackBurnerMinTituOutUpdated)
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
func (it *BuybackBurnerMinTituOutUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuybackBurnerMinTituOutUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuybackBurnerMinTituOutUpdated represents a MinTituOutUpdated event raised by the BuybackBurner contract.
type BuybackBurnerMinTituOutUpdated struct {
	OldMin *big.Int
	NewMin *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinTituOutUpdated is a free log retrieval operation binding the contract event 0x53fd101c8f6ba0f0fb4f44cc48b8e00b504b4f1e56073a9a170509f802488d84.
//
// Solidity: event MinTituOutUpdated(uint256 oldMin, uint256 newMin)
func (_BuybackBurner *BuybackBurnerFilterer) FilterMinTituOutUpdated(opts *bind.FilterOpts) (*BuybackBurnerMinTituOutUpdatedIterator, error) {

	logs, sub, err := _BuybackBurner.contract.FilterLogs(opts, "MinTituOutUpdated")
	if err != nil {
		return nil, err
	}
	return &BuybackBurnerMinTituOutUpdatedIterator{contract: _BuybackBurner.contract, event: "MinTituOutUpdated", logs: logs, sub: sub}, nil
}

// WatchMinTituOutUpdated is a free log subscription operation binding the contract event 0x53fd101c8f6ba0f0fb4f44cc48b8e00b504b4f1e56073a9a170509f802488d84.
//
// Solidity: event MinTituOutUpdated(uint256 oldMin, uint256 newMin)
func (_BuybackBurner *BuybackBurnerFilterer) WatchMinTituOutUpdated(opts *bind.WatchOpts, sink chan<- *BuybackBurnerMinTituOutUpdated) (event.Subscription, error) {

	logs, sub, err := _BuybackBurner.contract.WatchLogs(opts, "MinTituOutUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuybackBurnerMinTituOutUpdated)
				if err := _BuybackBurner.contract.UnpackLog(event, "MinTituOutUpdated", log); err != nil {
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

// ParseMinTituOutUpdated is a log parse operation binding the contract event 0x53fd101c8f6ba0f0fb4f44cc48b8e00b504b4f1e56073a9a170509f802488d84.
//
// Solidity: event MinTituOutUpdated(uint256 oldMin, uint256 newMin)
func (_BuybackBurner *BuybackBurnerFilterer) ParseMinTituOutUpdated(log types.Log) (*BuybackBurnerMinTituOutUpdated, error) {
	event := new(BuybackBurnerMinTituOutUpdated)
	if err := _BuybackBurner.contract.UnpackLog(event, "MinTituOutUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BuybackBurnerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the BuybackBurner contract.
type BuybackBurnerRoleAdminChangedIterator struct {
	Event *BuybackBurnerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *BuybackBurnerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuybackBurnerRoleAdminChanged)
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
		it.Event = new(BuybackBurnerRoleAdminChanged)
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
func (it *BuybackBurnerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuybackBurnerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuybackBurnerRoleAdminChanged represents a RoleAdminChanged event raised by the BuybackBurner contract.
type BuybackBurnerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_BuybackBurner *BuybackBurnerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*BuybackBurnerRoleAdminChangedIterator, error) {

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

	logs, sub, err := _BuybackBurner.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &BuybackBurnerRoleAdminChangedIterator{contract: _BuybackBurner.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_BuybackBurner *BuybackBurnerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *BuybackBurnerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _BuybackBurner.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuybackBurnerRoleAdminChanged)
				if err := _BuybackBurner.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_BuybackBurner *BuybackBurnerFilterer) ParseRoleAdminChanged(log types.Log) (*BuybackBurnerRoleAdminChanged, error) {
	event := new(BuybackBurnerRoleAdminChanged)
	if err := _BuybackBurner.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BuybackBurnerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the BuybackBurner contract.
type BuybackBurnerRoleGrantedIterator struct {
	Event *BuybackBurnerRoleGranted // Event containing the contract specifics and raw log

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
func (it *BuybackBurnerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuybackBurnerRoleGranted)
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
		it.Event = new(BuybackBurnerRoleGranted)
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
func (it *BuybackBurnerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuybackBurnerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuybackBurnerRoleGranted represents a RoleGranted event raised by the BuybackBurner contract.
type BuybackBurnerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_BuybackBurner *BuybackBurnerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BuybackBurnerRoleGrantedIterator, error) {

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

	logs, sub, err := _BuybackBurner.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BuybackBurnerRoleGrantedIterator{contract: _BuybackBurner.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_BuybackBurner *BuybackBurnerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *BuybackBurnerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _BuybackBurner.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuybackBurnerRoleGranted)
				if err := _BuybackBurner.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_BuybackBurner *BuybackBurnerFilterer) ParseRoleGranted(log types.Log) (*BuybackBurnerRoleGranted, error) {
	event := new(BuybackBurnerRoleGranted)
	if err := _BuybackBurner.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BuybackBurnerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the BuybackBurner contract.
type BuybackBurnerRoleRevokedIterator struct {
	Event *BuybackBurnerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *BuybackBurnerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuybackBurnerRoleRevoked)
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
		it.Event = new(BuybackBurnerRoleRevoked)
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
func (it *BuybackBurnerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuybackBurnerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuybackBurnerRoleRevoked represents a RoleRevoked event raised by the BuybackBurner contract.
type BuybackBurnerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_BuybackBurner *BuybackBurnerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BuybackBurnerRoleRevokedIterator, error) {

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

	logs, sub, err := _BuybackBurner.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BuybackBurnerRoleRevokedIterator{contract: _BuybackBurner.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_BuybackBurner *BuybackBurnerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *BuybackBurnerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _BuybackBurner.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuybackBurnerRoleRevoked)
				if err := _BuybackBurner.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_BuybackBurner *BuybackBurnerFilterer) ParseRoleRevoked(log types.Log) (*BuybackBurnerRoleRevoked, error) {
	event := new(BuybackBurnerRoleRevoked)
	if err := _BuybackBurner.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BuybackBurnerRouterUpdatedIterator is returned from FilterRouterUpdated and is used to iterate over the raw logs and unpacked data for RouterUpdated events raised by the BuybackBurner contract.
type BuybackBurnerRouterUpdatedIterator struct {
	Event *BuybackBurnerRouterUpdated // Event containing the contract specifics and raw log

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
func (it *BuybackBurnerRouterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuybackBurnerRouterUpdated)
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
		it.Event = new(BuybackBurnerRouterUpdated)
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
func (it *BuybackBurnerRouterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuybackBurnerRouterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuybackBurnerRouterUpdated represents a RouterUpdated event raised by the BuybackBurner contract.
type BuybackBurnerRouterUpdated struct {
	OldRouter common.Address
	NewRouter common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRouterUpdated is a free log retrieval operation binding the contract event 0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684.
//
// Solidity: event RouterUpdated(address indexed oldRouter, address indexed newRouter)
func (_BuybackBurner *BuybackBurnerFilterer) FilterRouterUpdated(opts *bind.FilterOpts, oldRouter []common.Address, newRouter []common.Address) (*BuybackBurnerRouterUpdatedIterator, error) {

	var oldRouterRule []interface{}
	for _, oldRouterItem := range oldRouter {
		oldRouterRule = append(oldRouterRule, oldRouterItem)
	}
	var newRouterRule []interface{}
	for _, newRouterItem := range newRouter {
		newRouterRule = append(newRouterRule, newRouterItem)
	}

	logs, sub, err := _BuybackBurner.contract.FilterLogs(opts, "RouterUpdated", oldRouterRule, newRouterRule)
	if err != nil {
		return nil, err
	}
	return &BuybackBurnerRouterUpdatedIterator{contract: _BuybackBurner.contract, event: "RouterUpdated", logs: logs, sub: sub}, nil
}

// WatchRouterUpdated is a free log subscription operation binding the contract event 0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684.
//
// Solidity: event RouterUpdated(address indexed oldRouter, address indexed newRouter)
func (_BuybackBurner *BuybackBurnerFilterer) WatchRouterUpdated(opts *bind.WatchOpts, sink chan<- *BuybackBurnerRouterUpdated, oldRouter []common.Address, newRouter []common.Address) (event.Subscription, error) {

	var oldRouterRule []interface{}
	for _, oldRouterItem := range oldRouter {
		oldRouterRule = append(oldRouterRule, oldRouterItem)
	}
	var newRouterRule []interface{}
	for _, newRouterItem := range newRouter {
		newRouterRule = append(newRouterRule, newRouterItem)
	}

	logs, sub, err := _BuybackBurner.contract.WatchLogs(opts, "RouterUpdated", oldRouterRule, newRouterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuybackBurnerRouterUpdated)
				if err := _BuybackBurner.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
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

// ParseRouterUpdated is a log parse operation binding the contract event 0x02dc5c233404867c793b749c6d644beb2277536d18a7e7974d3f238e4c6f1684.
//
// Solidity: event RouterUpdated(address indexed oldRouter, address indexed newRouter)
func (_BuybackBurner *BuybackBurnerFilterer) ParseRouterUpdated(log types.Log) (*BuybackBurnerRouterUpdated, error) {
	event := new(BuybackBurnerRouterUpdated)
	if err := _BuybackBurner.contract.UnpackLog(event, "RouterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BuybackBurnerSwapPathUpdatedIterator is returned from FilterSwapPathUpdated and is used to iterate over the raw logs and unpacked data for SwapPathUpdated events raised by the BuybackBurner contract.
type BuybackBurnerSwapPathUpdatedIterator struct {
	Event *BuybackBurnerSwapPathUpdated // Event containing the contract specifics and raw log

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
func (it *BuybackBurnerSwapPathUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuybackBurnerSwapPathUpdated)
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
		it.Event = new(BuybackBurnerSwapPathUpdated)
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
func (it *BuybackBurnerSwapPathUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuybackBurnerSwapPathUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuybackBurnerSwapPathUpdated represents a SwapPathUpdated event raised by the BuybackBurner contract.
type BuybackBurnerSwapPathUpdated struct {
	Path []common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSwapPathUpdated is a free log retrieval operation binding the contract event 0x6ac8d8fd2cc5d6c950ab25024970e2ab3d81ff6e6e5ebeb21d992ee119b8e6ce.
//
// Solidity: event SwapPathUpdated(address[] path)
func (_BuybackBurner *BuybackBurnerFilterer) FilterSwapPathUpdated(opts *bind.FilterOpts) (*BuybackBurnerSwapPathUpdatedIterator, error) {

	logs, sub, err := _BuybackBurner.contract.FilterLogs(opts, "SwapPathUpdated")
	if err != nil {
		return nil, err
	}
	return &BuybackBurnerSwapPathUpdatedIterator{contract: _BuybackBurner.contract, event: "SwapPathUpdated", logs: logs, sub: sub}, nil
}

// WatchSwapPathUpdated is a free log subscription operation binding the contract event 0x6ac8d8fd2cc5d6c950ab25024970e2ab3d81ff6e6e5ebeb21d992ee119b8e6ce.
//
// Solidity: event SwapPathUpdated(address[] path)
func (_BuybackBurner *BuybackBurnerFilterer) WatchSwapPathUpdated(opts *bind.WatchOpts, sink chan<- *BuybackBurnerSwapPathUpdated) (event.Subscription, error) {

	logs, sub, err := _BuybackBurner.contract.WatchLogs(opts, "SwapPathUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuybackBurnerSwapPathUpdated)
				if err := _BuybackBurner.contract.UnpackLog(event, "SwapPathUpdated", log); err != nil {
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

// ParseSwapPathUpdated is a log parse operation binding the contract event 0x6ac8d8fd2cc5d6c950ab25024970e2ab3d81ff6e6e5ebeb21d992ee119b8e6ce.
//
// Solidity: event SwapPathUpdated(address[] path)
func (_BuybackBurner *BuybackBurnerFilterer) ParseSwapPathUpdated(log types.Log) (*BuybackBurnerSwapPathUpdated, error) {
	event := new(BuybackBurnerSwapPathUpdated)
	if err := _BuybackBurner.contract.UnpackLog(event, "SwapPathUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
