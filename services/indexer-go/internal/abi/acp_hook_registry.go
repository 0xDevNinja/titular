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

// HookRegistryHookInfo is an auto generated low-level Go binding around an user-defined struct.
type HookRegistryHookInfo struct {
	Approved bool
	Name     string
}

// HookRegistryMetaData contains all meta data concerning the HookRegistry contract.
var HookRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allHooks\",\"inputs\":[],\"outputs\":[{\"name\":\"addresses\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deregisterHook\",\"inputs\":[{\"name\":\"hook\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getHookInfo\",\"inputs\":[{\"name\":\"hook\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"info\",\"type\":\"tuple\",\"internalType\":\"structHookRegistry.HookInfo\",\"components\":[{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isApproved\",\"inputs\":[{\"name\":\"hook\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerHook\",\"inputs\":[{\"name\":\"hook\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"HookDeregistered\",\"inputs\":[{\"name\":\"hook\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HookRegistered\",\"inputs\":[{\"name\":\"hook\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AlreadyRegistered\",\"inputs\":[{\"name\":\"hook\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotRegistered\",\"inputs\":[{\"name\":\"hook\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// HookRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use HookRegistryMetaData.ABI instead.
var HookRegistryABI = HookRegistryMetaData.ABI

// HookRegistry is an auto generated Go binding around an Ethereum contract.
type HookRegistry struct {
	HookRegistryCaller     // Read-only binding to the contract
	HookRegistryTransactor // Write-only binding to the contract
	HookRegistryFilterer   // Log filterer for contract events
}

// HookRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type HookRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HookRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HookRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HookRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HookRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HookRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HookRegistrySession struct {
	Contract     *HookRegistry     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HookRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HookRegistryCallerSession struct {
	Contract *HookRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// HookRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HookRegistryTransactorSession struct {
	Contract     *HookRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// HookRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type HookRegistryRaw struct {
	Contract *HookRegistry // Generic contract binding to access the raw methods on
}

// HookRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HookRegistryCallerRaw struct {
	Contract *HookRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// HookRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HookRegistryTransactorRaw struct {
	Contract *HookRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHookRegistry creates a new instance of HookRegistry, bound to a specific deployed contract.
func NewHookRegistry(address common.Address, backend bind.ContractBackend) (*HookRegistry, error) {
	contract, err := bindHookRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HookRegistry{HookRegistryCaller: HookRegistryCaller{contract: contract}, HookRegistryTransactor: HookRegistryTransactor{contract: contract}, HookRegistryFilterer: HookRegistryFilterer{contract: contract}}, nil
}

// NewHookRegistryCaller creates a new read-only instance of HookRegistry, bound to a specific deployed contract.
func NewHookRegistryCaller(address common.Address, caller bind.ContractCaller) (*HookRegistryCaller, error) {
	contract, err := bindHookRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HookRegistryCaller{contract: contract}, nil
}

// NewHookRegistryTransactor creates a new write-only instance of HookRegistry, bound to a specific deployed contract.
func NewHookRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*HookRegistryTransactor, error) {
	contract, err := bindHookRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HookRegistryTransactor{contract: contract}, nil
}

// NewHookRegistryFilterer creates a new log filterer instance of HookRegistry, bound to a specific deployed contract.
func NewHookRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*HookRegistryFilterer, error) {
	contract, err := bindHookRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HookRegistryFilterer{contract: contract}, nil
}

// bindHookRegistry binds a generic wrapper to an already deployed contract.
func bindHookRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HookRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HookRegistry *HookRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HookRegistry.Contract.HookRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HookRegistry *HookRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HookRegistry.Contract.HookRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HookRegistry *HookRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HookRegistry.Contract.HookRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HookRegistry *HookRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HookRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HookRegistry *HookRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HookRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HookRegistry *HookRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HookRegistry.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_HookRegistry *HookRegistryCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _HookRegistry.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_HookRegistry *HookRegistrySession) DEFAULTADMINROLE() ([32]byte, error) {
	return _HookRegistry.Contract.DEFAULTADMINROLE(&_HookRegistry.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_HookRegistry *HookRegistryCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _HookRegistry.Contract.DEFAULTADMINROLE(&_HookRegistry.CallOpts)
}

// AllHooks is a free data retrieval call binding the contract method 0x3d5522f2.
//
// Solidity: function allHooks() view returns(address[] addresses)
func (_HookRegistry *HookRegistryCaller) AllHooks(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _HookRegistry.contract.Call(opts, &out, "allHooks")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// AllHooks is a free data retrieval call binding the contract method 0x3d5522f2.
//
// Solidity: function allHooks() view returns(address[] addresses)
func (_HookRegistry *HookRegistrySession) AllHooks() ([]common.Address, error) {
	return _HookRegistry.Contract.AllHooks(&_HookRegistry.CallOpts)
}

// AllHooks is a free data retrieval call binding the contract method 0x3d5522f2.
//
// Solidity: function allHooks() view returns(address[] addresses)
func (_HookRegistry *HookRegistryCallerSession) AllHooks() ([]common.Address, error) {
	return _HookRegistry.Contract.AllHooks(&_HookRegistry.CallOpts)
}

// GetHookInfo is a free data retrieval call binding the contract method 0x25c7f789.
//
// Solidity: function getHookInfo(address hook) view returns((bool,string) info)
func (_HookRegistry *HookRegistryCaller) GetHookInfo(opts *bind.CallOpts, hook common.Address) (HookRegistryHookInfo, error) {
	var out []interface{}
	err := _HookRegistry.contract.Call(opts, &out, "getHookInfo", hook)

	if err != nil {
		return *new(HookRegistryHookInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(HookRegistryHookInfo)).(*HookRegistryHookInfo)

	return out0, err

}

// GetHookInfo is a free data retrieval call binding the contract method 0x25c7f789.
//
// Solidity: function getHookInfo(address hook) view returns((bool,string) info)
func (_HookRegistry *HookRegistrySession) GetHookInfo(hook common.Address) (HookRegistryHookInfo, error) {
	return _HookRegistry.Contract.GetHookInfo(&_HookRegistry.CallOpts, hook)
}

// GetHookInfo is a free data retrieval call binding the contract method 0x25c7f789.
//
// Solidity: function getHookInfo(address hook) view returns((bool,string) info)
func (_HookRegistry *HookRegistryCallerSession) GetHookInfo(hook common.Address) (HookRegistryHookInfo, error) {
	return _HookRegistry.Contract.GetHookInfo(&_HookRegistry.CallOpts, hook)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_HookRegistry *HookRegistryCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _HookRegistry.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_HookRegistry *HookRegistrySession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _HookRegistry.Contract.GetRoleAdmin(&_HookRegistry.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_HookRegistry *HookRegistryCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _HookRegistry.Contract.GetRoleAdmin(&_HookRegistry.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_HookRegistry *HookRegistryCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _HookRegistry.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_HookRegistry *HookRegistrySession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _HookRegistry.Contract.HasRole(&_HookRegistry.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_HookRegistry *HookRegistryCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _HookRegistry.Contract.HasRole(&_HookRegistry.CallOpts, role, account)
}

// IsApproved is a free data retrieval call binding the contract method 0x673448dd.
//
// Solidity: function isApproved(address hook) view returns(bool approved)
func (_HookRegistry *HookRegistryCaller) IsApproved(opts *bind.CallOpts, hook common.Address) (bool, error) {
	var out []interface{}
	err := _HookRegistry.contract.Call(opts, &out, "isApproved", hook)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApproved is a free data retrieval call binding the contract method 0x673448dd.
//
// Solidity: function isApproved(address hook) view returns(bool approved)
func (_HookRegistry *HookRegistrySession) IsApproved(hook common.Address) (bool, error) {
	return _HookRegistry.Contract.IsApproved(&_HookRegistry.CallOpts, hook)
}

// IsApproved is a free data retrieval call binding the contract method 0x673448dd.
//
// Solidity: function isApproved(address hook) view returns(bool approved)
func (_HookRegistry *HookRegistryCallerSession) IsApproved(hook common.Address) (bool, error) {
	return _HookRegistry.Contract.IsApproved(&_HookRegistry.CallOpts, hook)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_HookRegistry *HookRegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _HookRegistry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_HookRegistry *HookRegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _HookRegistry.Contract.SupportsInterface(&_HookRegistry.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_HookRegistry *HookRegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _HookRegistry.Contract.SupportsInterface(&_HookRegistry.CallOpts, interfaceId)
}

// DeregisterHook is a paid mutator transaction binding the contract method 0xc58e9eb0.
//
// Solidity: function deregisterHook(address hook) returns()
func (_HookRegistry *HookRegistryTransactor) DeregisterHook(opts *bind.TransactOpts, hook common.Address) (*types.Transaction, error) {
	return _HookRegistry.contract.Transact(opts, "deregisterHook", hook)
}

// DeregisterHook is a paid mutator transaction binding the contract method 0xc58e9eb0.
//
// Solidity: function deregisterHook(address hook) returns()
func (_HookRegistry *HookRegistrySession) DeregisterHook(hook common.Address) (*types.Transaction, error) {
	return _HookRegistry.Contract.DeregisterHook(&_HookRegistry.TransactOpts, hook)
}

// DeregisterHook is a paid mutator transaction binding the contract method 0xc58e9eb0.
//
// Solidity: function deregisterHook(address hook) returns()
func (_HookRegistry *HookRegistryTransactorSession) DeregisterHook(hook common.Address) (*types.Transaction, error) {
	return _HookRegistry.Contract.DeregisterHook(&_HookRegistry.TransactOpts, hook)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_HookRegistry *HookRegistryTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _HookRegistry.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_HookRegistry *HookRegistrySession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _HookRegistry.Contract.GrantRole(&_HookRegistry.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_HookRegistry *HookRegistryTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _HookRegistry.Contract.GrantRole(&_HookRegistry.TransactOpts, role, account)
}

// RegisterHook is a paid mutator transaction binding the contract method 0x6354b661.
//
// Solidity: function registerHook(address hook) returns()
func (_HookRegistry *HookRegistryTransactor) RegisterHook(opts *bind.TransactOpts, hook common.Address) (*types.Transaction, error) {
	return _HookRegistry.contract.Transact(opts, "registerHook", hook)
}

// RegisterHook is a paid mutator transaction binding the contract method 0x6354b661.
//
// Solidity: function registerHook(address hook) returns()
func (_HookRegistry *HookRegistrySession) RegisterHook(hook common.Address) (*types.Transaction, error) {
	return _HookRegistry.Contract.RegisterHook(&_HookRegistry.TransactOpts, hook)
}

// RegisterHook is a paid mutator transaction binding the contract method 0x6354b661.
//
// Solidity: function registerHook(address hook) returns()
func (_HookRegistry *HookRegistryTransactorSession) RegisterHook(hook common.Address) (*types.Transaction, error) {
	return _HookRegistry.Contract.RegisterHook(&_HookRegistry.TransactOpts, hook)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_HookRegistry *HookRegistryTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _HookRegistry.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_HookRegistry *HookRegistrySession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _HookRegistry.Contract.RenounceRole(&_HookRegistry.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_HookRegistry *HookRegistryTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _HookRegistry.Contract.RenounceRole(&_HookRegistry.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_HookRegistry *HookRegistryTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _HookRegistry.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_HookRegistry *HookRegistrySession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _HookRegistry.Contract.RevokeRole(&_HookRegistry.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_HookRegistry *HookRegistryTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _HookRegistry.Contract.RevokeRole(&_HookRegistry.TransactOpts, role, account)
}

// HookRegistryHookDeregisteredIterator is returned from FilterHookDeregistered and is used to iterate over the raw logs and unpacked data for HookDeregistered events raised by the HookRegistry contract.
type HookRegistryHookDeregisteredIterator struct {
	Event *HookRegistryHookDeregistered // Event containing the contract specifics and raw log

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
func (it *HookRegistryHookDeregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HookRegistryHookDeregistered)
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
		it.Event = new(HookRegistryHookDeregistered)
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
func (it *HookRegistryHookDeregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HookRegistryHookDeregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HookRegistryHookDeregistered represents a HookDeregistered event raised by the HookRegistry contract.
type HookRegistryHookDeregistered struct {
	Hook common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterHookDeregistered is a free log retrieval operation binding the contract event 0x5f04ed00bba892fd984432cd96beabe664a2de332ff480f08ec0c135c0372d55.
//
// Solidity: event HookDeregistered(address indexed hook)
func (_HookRegistry *HookRegistryFilterer) FilterHookDeregistered(opts *bind.FilterOpts, hook []common.Address) (*HookRegistryHookDeregisteredIterator, error) {

	var hookRule []interface{}
	for _, hookItem := range hook {
		hookRule = append(hookRule, hookItem)
	}

	logs, sub, err := _HookRegistry.contract.FilterLogs(opts, "HookDeregistered", hookRule)
	if err != nil {
		return nil, err
	}
	return &HookRegistryHookDeregisteredIterator{contract: _HookRegistry.contract, event: "HookDeregistered", logs: logs, sub: sub}, nil
}

// WatchHookDeregistered is a free log subscription operation binding the contract event 0x5f04ed00bba892fd984432cd96beabe664a2de332ff480f08ec0c135c0372d55.
//
// Solidity: event HookDeregistered(address indexed hook)
func (_HookRegistry *HookRegistryFilterer) WatchHookDeregistered(opts *bind.WatchOpts, sink chan<- *HookRegistryHookDeregistered, hook []common.Address) (event.Subscription, error) {

	var hookRule []interface{}
	for _, hookItem := range hook {
		hookRule = append(hookRule, hookItem)
	}

	logs, sub, err := _HookRegistry.contract.WatchLogs(opts, "HookDeregistered", hookRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HookRegistryHookDeregistered)
				if err := _HookRegistry.contract.UnpackLog(event, "HookDeregistered", log); err != nil {
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

// ParseHookDeregistered is a log parse operation binding the contract event 0x5f04ed00bba892fd984432cd96beabe664a2de332ff480f08ec0c135c0372d55.
//
// Solidity: event HookDeregistered(address indexed hook)
func (_HookRegistry *HookRegistryFilterer) ParseHookDeregistered(log types.Log) (*HookRegistryHookDeregistered, error) {
	event := new(HookRegistryHookDeregistered)
	if err := _HookRegistry.contract.UnpackLog(event, "HookDeregistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HookRegistryHookRegisteredIterator is returned from FilterHookRegistered and is used to iterate over the raw logs and unpacked data for HookRegistered events raised by the HookRegistry contract.
type HookRegistryHookRegisteredIterator struct {
	Event *HookRegistryHookRegistered // Event containing the contract specifics and raw log

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
func (it *HookRegistryHookRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HookRegistryHookRegistered)
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
		it.Event = new(HookRegistryHookRegistered)
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
func (it *HookRegistryHookRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HookRegistryHookRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HookRegistryHookRegistered represents a HookRegistered event raised by the HookRegistry contract.
type HookRegistryHookRegistered struct {
	Hook common.Address
	Name string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterHookRegistered is a free log retrieval operation binding the contract event 0x8096ae10f94d9c9266fe844741548a8afe2a55900eb4607c7a5c3dd5560651a9.
//
// Solidity: event HookRegistered(address indexed hook, string name)
func (_HookRegistry *HookRegistryFilterer) FilterHookRegistered(opts *bind.FilterOpts, hook []common.Address) (*HookRegistryHookRegisteredIterator, error) {

	var hookRule []interface{}
	for _, hookItem := range hook {
		hookRule = append(hookRule, hookItem)
	}

	logs, sub, err := _HookRegistry.contract.FilterLogs(opts, "HookRegistered", hookRule)
	if err != nil {
		return nil, err
	}
	return &HookRegistryHookRegisteredIterator{contract: _HookRegistry.contract, event: "HookRegistered", logs: logs, sub: sub}, nil
}

// WatchHookRegistered is a free log subscription operation binding the contract event 0x8096ae10f94d9c9266fe844741548a8afe2a55900eb4607c7a5c3dd5560651a9.
//
// Solidity: event HookRegistered(address indexed hook, string name)
func (_HookRegistry *HookRegistryFilterer) WatchHookRegistered(opts *bind.WatchOpts, sink chan<- *HookRegistryHookRegistered, hook []common.Address) (event.Subscription, error) {

	var hookRule []interface{}
	for _, hookItem := range hook {
		hookRule = append(hookRule, hookItem)
	}

	logs, sub, err := _HookRegistry.contract.WatchLogs(opts, "HookRegistered", hookRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HookRegistryHookRegistered)
				if err := _HookRegistry.contract.UnpackLog(event, "HookRegistered", log); err != nil {
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

// ParseHookRegistered is a log parse operation binding the contract event 0x8096ae10f94d9c9266fe844741548a8afe2a55900eb4607c7a5c3dd5560651a9.
//
// Solidity: event HookRegistered(address indexed hook, string name)
func (_HookRegistry *HookRegistryFilterer) ParseHookRegistered(log types.Log) (*HookRegistryHookRegistered, error) {
	event := new(HookRegistryHookRegistered)
	if err := _HookRegistry.contract.UnpackLog(event, "HookRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HookRegistryRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the HookRegistry contract.
type HookRegistryRoleAdminChangedIterator struct {
	Event *HookRegistryRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *HookRegistryRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HookRegistryRoleAdminChanged)
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
		it.Event = new(HookRegistryRoleAdminChanged)
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
func (it *HookRegistryRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HookRegistryRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HookRegistryRoleAdminChanged represents a RoleAdminChanged event raised by the HookRegistry contract.
type HookRegistryRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_HookRegistry *HookRegistryFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*HookRegistryRoleAdminChangedIterator, error) {

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

	logs, sub, err := _HookRegistry.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &HookRegistryRoleAdminChangedIterator{contract: _HookRegistry.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_HookRegistry *HookRegistryFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *HookRegistryRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _HookRegistry.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HookRegistryRoleAdminChanged)
				if err := _HookRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_HookRegistry *HookRegistryFilterer) ParseRoleAdminChanged(log types.Log) (*HookRegistryRoleAdminChanged, error) {
	event := new(HookRegistryRoleAdminChanged)
	if err := _HookRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HookRegistryRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the HookRegistry contract.
type HookRegistryRoleGrantedIterator struct {
	Event *HookRegistryRoleGranted // Event containing the contract specifics and raw log

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
func (it *HookRegistryRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HookRegistryRoleGranted)
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
		it.Event = new(HookRegistryRoleGranted)
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
func (it *HookRegistryRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HookRegistryRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HookRegistryRoleGranted represents a RoleGranted event raised by the HookRegistry contract.
type HookRegistryRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_HookRegistry *HookRegistryFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*HookRegistryRoleGrantedIterator, error) {

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

	logs, sub, err := _HookRegistry.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &HookRegistryRoleGrantedIterator{contract: _HookRegistry.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_HookRegistry *HookRegistryFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *HookRegistryRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _HookRegistry.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HookRegistryRoleGranted)
				if err := _HookRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_HookRegistry *HookRegistryFilterer) ParseRoleGranted(log types.Log) (*HookRegistryRoleGranted, error) {
	event := new(HookRegistryRoleGranted)
	if err := _HookRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HookRegistryRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the HookRegistry contract.
type HookRegistryRoleRevokedIterator struct {
	Event *HookRegistryRoleRevoked // Event containing the contract specifics and raw log

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
func (it *HookRegistryRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HookRegistryRoleRevoked)
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
		it.Event = new(HookRegistryRoleRevoked)
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
func (it *HookRegistryRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HookRegistryRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HookRegistryRoleRevoked represents a RoleRevoked event raised by the HookRegistry contract.
type HookRegistryRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_HookRegistry *HookRegistryFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*HookRegistryRoleRevokedIterator, error) {

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

	logs, sub, err := _HookRegistry.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &HookRegistryRoleRevokedIterator{contract: _HookRegistry.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_HookRegistry *HookRegistryFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *HookRegistryRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _HookRegistry.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HookRegistryRoleRevoked)
				if err := _HookRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_HookRegistry *HookRegistryFilterer) ParseRoleRevoked(log types.Log) (*HookRegistryRoleRevoked, error) {
	event := new(HookRegistryRoleRevoked)
	if err := _HookRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
