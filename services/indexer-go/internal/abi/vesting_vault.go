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

// VestingVaultMetaData contains all meta data concerning the VestingVault contract.
var VestingVaultMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TOKEN\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIERC20\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addGrant\",\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"start\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"cliff\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"duration\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"grants\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"total\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"released\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"start\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"cliff\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"duration\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"releasable\",\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"release\",\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"release\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revoke\",\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"vested\",\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"GrantAdded\",\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"start\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"cliff\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"duration\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Released\",\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Revoked\",\"inputs\":[{\"name\":\"beneficiary\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"vestedKept\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"unvestedReturned\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"CliffExceedsDuration\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GrantExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoGrant\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NothingToRelease\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroDuration\",\"inputs\":[]}]",
}

// VestingVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use VestingVaultMetaData.ABI instead.
var VestingVaultABI = VestingVaultMetaData.ABI

// VestingVault is an auto generated Go binding around an Ethereum contract.
type VestingVault struct {
	VestingVaultCaller     // Read-only binding to the contract
	VestingVaultTransactor // Write-only binding to the contract
	VestingVaultFilterer   // Log filterer for contract events
}

// VestingVaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type VestingVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VestingVaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VestingVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VestingVaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VestingVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VestingVaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VestingVaultSession struct {
	Contract     *VestingVault     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VestingVaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VestingVaultCallerSession struct {
	Contract *VestingVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// VestingVaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VestingVaultTransactorSession struct {
	Contract     *VestingVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// VestingVaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type VestingVaultRaw struct {
	Contract *VestingVault // Generic contract binding to access the raw methods on
}

// VestingVaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VestingVaultCallerRaw struct {
	Contract *VestingVaultCaller // Generic read-only contract binding to access the raw methods on
}

// VestingVaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VestingVaultTransactorRaw struct {
	Contract *VestingVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVestingVault creates a new instance of VestingVault, bound to a specific deployed contract.
func NewVestingVault(address common.Address, backend bind.ContractBackend) (*VestingVault, error) {
	contract, err := bindVestingVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VestingVault{VestingVaultCaller: VestingVaultCaller{contract: contract}, VestingVaultTransactor: VestingVaultTransactor{contract: contract}, VestingVaultFilterer: VestingVaultFilterer{contract: contract}}, nil
}

// NewVestingVaultCaller creates a new read-only instance of VestingVault, bound to a specific deployed contract.
func NewVestingVaultCaller(address common.Address, caller bind.ContractCaller) (*VestingVaultCaller, error) {
	contract, err := bindVestingVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VestingVaultCaller{contract: contract}, nil
}

// NewVestingVaultTransactor creates a new write-only instance of VestingVault, bound to a specific deployed contract.
func NewVestingVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*VestingVaultTransactor, error) {
	contract, err := bindVestingVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VestingVaultTransactor{contract: contract}, nil
}

// NewVestingVaultFilterer creates a new log filterer instance of VestingVault, bound to a specific deployed contract.
func NewVestingVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*VestingVaultFilterer, error) {
	contract, err := bindVestingVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VestingVaultFilterer{contract: contract}, nil
}

// bindVestingVault binds a generic wrapper to an already deployed contract.
func bindVestingVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VestingVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VestingVault *VestingVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VestingVault.Contract.VestingVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VestingVault *VestingVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VestingVault.Contract.VestingVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VestingVault *VestingVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VestingVault.Contract.VestingVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VestingVault *VestingVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VestingVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VestingVault *VestingVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VestingVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VestingVault *VestingVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VestingVault.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_VestingVault *VestingVaultCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _VestingVault.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_VestingVault *VestingVaultSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _VestingVault.Contract.DEFAULTADMINROLE(&_VestingVault.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_VestingVault *VestingVaultCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _VestingVault.Contract.DEFAULTADMINROLE(&_VestingVault.CallOpts)
}

// TOKEN is a free data retrieval call binding the contract method 0x82bfefc8.
//
// Solidity: function TOKEN() view returns(address)
func (_VestingVault *VestingVaultCaller) TOKEN(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _VestingVault.contract.Call(opts, &out, "TOKEN")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TOKEN is a free data retrieval call binding the contract method 0x82bfefc8.
//
// Solidity: function TOKEN() view returns(address)
func (_VestingVault *VestingVaultSession) TOKEN() (common.Address, error) {
	return _VestingVault.Contract.TOKEN(&_VestingVault.CallOpts)
}

// TOKEN is a free data retrieval call binding the contract method 0x82bfefc8.
//
// Solidity: function TOKEN() view returns(address)
func (_VestingVault *VestingVaultCallerSession) TOKEN() (common.Address, error) {
	return _VestingVault.Contract.TOKEN(&_VestingVault.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_VestingVault *VestingVaultCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _VestingVault.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_VestingVault *VestingVaultSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _VestingVault.Contract.GetRoleAdmin(&_VestingVault.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_VestingVault *VestingVaultCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _VestingVault.Contract.GetRoleAdmin(&_VestingVault.CallOpts, role)
}

// Grants is a free data retrieval call binding the contract method 0xb869cea3.
//
// Solidity: function grants(address ) view returns(uint256 total, uint256 released, uint64 start, uint64 cliff, uint64 duration)
func (_VestingVault *VestingVaultCaller) Grants(opts *bind.CallOpts, arg0 common.Address) (struct {
	Total    *big.Int
	Released *big.Int
	Start    uint64
	Cliff    uint64
	Duration uint64
}, error) {
	var out []interface{}
	err := _VestingVault.contract.Call(opts, &out, "grants", arg0)

	outstruct := new(struct {
		Total    *big.Int
		Released *big.Int
		Start    uint64
		Cliff    uint64
		Duration uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Total = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Released = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Start = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Cliff = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.Duration = *abi.ConvertType(out[4], new(uint64)).(*uint64)

	return *outstruct, err

}

// Grants is a free data retrieval call binding the contract method 0xb869cea3.
//
// Solidity: function grants(address ) view returns(uint256 total, uint256 released, uint64 start, uint64 cliff, uint64 duration)
func (_VestingVault *VestingVaultSession) Grants(arg0 common.Address) (struct {
	Total    *big.Int
	Released *big.Int
	Start    uint64
	Cliff    uint64
	Duration uint64
}, error) {
	return _VestingVault.Contract.Grants(&_VestingVault.CallOpts, arg0)
}

// Grants is a free data retrieval call binding the contract method 0xb869cea3.
//
// Solidity: function grants(address ) view returns(uint256 total, uint256 released, uint64 start, uint64 cliff, uint64 duration)
func (_VestingVault *VestingVaultCallerSession) Grants(arg0 common.Address) (struct {
	Total    *big.Int
	Released *big.Int
	Start    uint64
	Cliff    uint64
	Duration uint64
}, error) {
	return _VestingVault.Contract.Grants(&_VestingVault.CallOpts, arg0)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_VestingVault *VestingVaultCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _VestingVault.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_VestingVault *VestingVaultSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _VestingVault.Contract.HasRole(&_VestingVault.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_VestingVault *VestingVaultCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _VestingVault.Contract.HasRole(&_VestingVault.CallOpts, role, account)
}

// Releasable is a free data retrieval call binding the contract method 0xa3f8eace.
//
// Solidity: function releasable(address beneficiary) view returns(uint256)
func (_VestingVault *VestingVaultCaller) Releasable(opts *bind.CallOpts, beneficiary common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VestingVault.contract.Call(opts, &out, "releasable", beneficiary)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Releasable is a free data retrieval call binding the contract method 0xa3f8eace.
//
// Solidity: function releasable(address beneficiary) view returns(uint256)
func (_VestingVault *VestingVaultSession) Releasable(beneficiary common.Address) (*big.Int, error) {
	return _VestingVault.Contract.Releasable(&_VestingVault.CallOpts, beneficiary)
}

// Releasable is a free data retrieval call binding the contract method 0xa3f8eace.
//
// Solidity: function releasable(address beneficiary) view returns(uint256)
func (_VestingVault *VestingVaultCallerSession) Releasable(beneficiary common.Address) (*big.Int, error) {
	return _VestingVault.Contract.Releasable(&_VestingVault.CallOpts, beneficiary)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_VestingVault *VestingVaultCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _VestingVault.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_VestingVault *VestingVaultSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _VestingVault.Contract.SupportsInterface(&_VestingVault.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_VestingVault *VestingVaultCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _VestingVault.Contract.SupportsInterface(&_VestingVault.CallOpts, interfaceId)
}

// Vested is a free data retrieval call binding the contract method 0x7102b728.
//
// Solidity: function vested(address beneficiary) view returns(uint256)
func (_VestingVault *VestingVaultCaller) Vested(opts *bind.CallOpts, beneficiary common.Address) (*big.Int, error) {
	var out []interface{}
	err := _VestingVault.contract.Call(opts, &out, "vested", beneficiary)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Vested is a free data retrieval call binding the contract method 0x7102b728.
//
// Solidity: function vested(address beneficiary) view returns(uint256)
func (_VestingVault *VestingVaultSession) Vested(beneficiary common.Address) (*big.Int, error) {
	return _VestingVault.Contract.Vested(&_VestingVault.CallOpts, beneficiary)
}

// Vested is a free data retrieval call binding the contract method 0x7102b728.
//
// Solidity: function vested(address beneficiary) view returns(uint256)
func (_VestingVault *VestingVaultCallerSession) Vested(beneficiary common.Address) (*big.Int, error) {
	return _VestingVault.Contract.Vested(&_VestingVault.CallOpts, beneficiary)
}

// AddGrant is a paid mutator transaction binding the contract method 0x3b89403b.
//
// Solidity: function addGrant(address beneficiary, uint256 amount, uint64 start, uint64 cliff, uint64 duration) returns()
func (_VestingVault *VestingVaultTransactor) AddGrant(opts *bind.TransactOpts, beneficiary common.Address, amount *big.Int, start uint64, cliff uint64, duration uint64) (*types.Transaction, error) {
	return _VestingVault.contract.Transact(opts, "addGrant", beneficiary, amount, start, cliff, duration)
}

// AddGrant is a paid mutator transaction binding the contract method 0x3b89403b.
//
// Solidity: function addGrant(address beneficiary, uint256 amount, uint64 start, uint64 cliff, uint64 duration) returns()
func (_VestingVault *VestingVaultSession) AddGrant(beneficiary common.Address, amount *big.Int, start uint64, cliff uint64, duration uint64) (*types.Transaction, error) {
	return _VestingVault.Contract.AddGrant(&_VestingVault.TransactOpts, beneficiary, amount, start, cliff, duration)
}

// AddGrant is a paid mutator transaction binding the contract method 0x3b89403b.
//
// Solidity: function addGrant(address beneficiary, uint256 amount, uint64 start, uint64 cliff, uint64 duration) returns()
func (_VestingVault *VestingVaultTransactorSession) AddGrant(beneficiary common.Address, amount *big.Int, start uint64, cliff uint64, duration uint64) (*types.Transaction, error) {
	return _VestingVault.Contract.AddGrant(&_VestingVault.TransactOpts, beneficiary, amount, start, cliff, duration)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_VestingVault *VestingVaultTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VestingVault.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_VestingVault *VestingVaultSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VestingVault.Contract.GrantRole(&_VestingVault.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_VestingVault *VestingVaultTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VestingVault.Contract.GrantRole(&_VestingVault.TransactOpts, role, account)
}

// Release is a paid mutator transaction binding the contract method 0x19165587.
//
// Solidity: function release(address beneficiary) returns()
func (_VestingVault *VestingVaultTransactor) Release(opts *bind.TransactOpts, beneficiary common.Address) (*types.Transaction, error) {
	return _VestingVault.contract.Transact(opts, "release", beneficiary)
}

// Release is a paid mutator transaction binding the contract method 0x19165587.
//
// Solidity: function release(address beneficiary) returns()
func (_VestingVault *VestingVaultSession) Release(beneficiary common.Address) (*types.Transaction, error) {
	return _VestingVault.Contract.Release(&_VestingVault.TransactOpts, beneficiary)
}

// Release is a paid mutator transaction binding the contract method 0x19165587.
//
// Solidity: function release(address beneficiary) returns()
func (_VestingVault *VestingVaultTransactorSession) Release(beneficiary common.Address) (*types.Transaction, error) {
	return _VestingVault.Contract.Release(&_VestingVault.TransactOpts, beneficiary)
}

// Release0 is a paid mutator transaction binding the contract method 0x86d1a69f.
//
// Solidity: function release() returns()
func (_VestingVault *VestingVaultTransactor) Release0(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VestingVault.contract.Transact(opts, "release0")
}

// Release0 is a paid mutator transaction binding the contract method 0x86d1a69f.
//
// Solidity: function release() returns()
func (_VestingVault *VestingVaultSession) Release0() (*types.Transaction, error) {
	return _VestingVault.Contract.Release0(&_VestingVault.TransactOpts)
}

// Release0 is a paid mutator transaction binding the contract method 0x86d1a69f.
//
// Solidity: function release() returns()
func (_VestingVault *VestingVaultTransactorSession) Release0() (*types.Transaction, error) {
	return _VestingVault.Contract.Release0(&_VestingVault.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_VestingVault *VestingVaultTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _VestingVault.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_VestingVault *VestingVaultSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _VestingVault.Contract.RenounceRole(&_VestingVault.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_VestingVault *VestingVaultTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _VestingVault.Contract.RenounceRole(&_VestingVault.TransactOpts, role, callerConfirmation)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address beneficiary) returns()
func (_VestingVault *VestingVaultTransactor) Revoke(opts *bind.TransactOpts, beneficiary common.Address) (*types.Transaction, error) {
	return _VestingVault.contract.Transact(opts, "revoke", beneficiary)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address beneficiary) returns()
func (_VestingVault *VestingVaultSession) Revoke(beneficiary common.Address) (*types.Transaction, error) {
	return _VestingVault.Contract.Revoke(&_VestingVault.TransactOpts, beneficiary)
}

// Revoke is a paid mutator transaction binding the contract method 0x74a8f103.
//
// Solidity: function revoke(address beneficiary) returns()
func (_VestingVault *VestingVaultTransactorSession) Revoke(beneficiary common.Address) (*types.Transaction, error) {
	return _VestingVault.Contract.Revoke(&_VestingVault.TransactOpts, beneficiary)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_VestingVault *VestingVaultTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VestingVault.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_VestingVault *VestingVaultSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VestingVault.Contract.RevokeRole(&_VestingVault.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_VestingVault *VestingVaultTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _VestingVault.Contract.RevokeRole(&_VestingVault.TransactOpts, role, account)
}

// VestingVaultGrantAddedIterator is returned from FilterGrantAdded and is used to iterate over the raw logs and unpacked data for GrantAdded events raised by the VestingVault contract.
type VestingVaultGrantAddedIterator struct {
	Event *VestingVaultGrantAdded // Event containing the contract specifics and raw log

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
func (it *VestingVaultGrantAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VestingVaultGrantAdded)
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
		it.Event = new(VestingVaultGrantAdded)
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
func (it *VestingVaultGrantAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VestingVaultGrantAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VestingVaultGrantAdded represents a GrantAdded event raised by the VestingVault contract.
type VestingVaultGrantAdded struct {
	Beneficiary common.Address
	Amount      *big.Int
	Start       uint64
	Cliff       uint64
	Duration    uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGrantAdded is a free log retrieval operation binding the contract event 0x504ecf986b8e5a88b539ab65c8f375584c39fe9def76092420ba4e90e283eae7.
//
// Solidity: event GrantAdded(address indexed beneficiary, uint256 amount, uint64 start, uint64 cliff, uint64 duration)
func (_VestingVault *VestingVaultFilterer) FilterGrantAdded(opts *bind.FilterOpts, beneficiary []common.Address) (*VestingVaultGrantAddedIterator, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _VestingVault.contract.FilterLogs(opts, "GrantAdded", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return &VestingVaultGrantAddedIterator{contract: _VestingVault.contract, event: "GrantAdded", logs: logs, sub: sub}, nil
}

// WatchGrantAdded is a free log subscription operation binding the contract event 0x504ecf986b8e5a88b539ab65c8f375584c39fe9def76092420ba4e90e283eae7.
//
// Solidity: event GrantAdded(address indexed beneficiary, uint256 amount, uint64 start, uint64 cliff, uint64 duration)
func (_VestingVault *VestingVaultFilterer) WatchGrantAdded(opts *bind.WatchOpts, sink chan<- *VestingVaultGrantAdded, beneficiary []common.Address) (event.Subscription, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _VestingVault.contract.WatchLogs(opts, "GrantAdded", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VestingVaultGrantAdded)
				if err := _VestingVault.contract.UnpackLog(event, "GrantAdded", log); err != nil {
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

// ParseGrantAdded is a log parse operation binding the contract event 0x504ecf986b8e5a88b539ab65c8f375584c39fe9def76092420ba4e90e283eae7.
//
// Solidity: event GrantAdded(address indexed beneficiary, uint256 amount, uint64 start, uint64 cliff, uint64 duration)
func (_VestingVault *VestingVaultFilterer) ParseGrantAdded(log types.Log) (*VestingVaultGrantAdded, error) {
	event := new(VestingVaultGrantAdded)
	if err := _VestingVault.contract.UnpackLog(event, "GrantAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VestingVaultReleasedIterator is returned from FilterReleased and is used to iterate over the raw logs and unpacked data for Released events raised by the VestingVault contract.
type VestingVaultReleasedIterator struct {
	Event *VestingVaultReleased // Event containing the contract specifics and raw log

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
func (it *VestingVaultReleasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VestingVaultReleased)
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
		it.Event = new(VestingVaultReleased)
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
func (it *VestingVaultReleasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VestingVaultReleasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VestingVaultReleased represents a Released event raised by the VestingVault contract.
type VestingVaultReleased struct {
	Beneficiary common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReleased is a free log retrieval operation binding the contract event 0xb21fb52d5749b80f3182f8c6992236b5e5576681880914484d7f4c9b062e619e.
//
// Solidity: event Released(address indexed beneficiary, uint256 amount)
func (_VestingVault *VestingVaultFilterer) FilterReleased(opts *bind.FilterOpts, beneficiary []common.Address) (*VestingVaultReleasedIterator, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _VestingVault.contract.FilterLogs(opts, "Released", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return &VestingVaultReleasedIterator{contract: _VestingVault.contract, event: "Released", logs: logs, sub: sub}, nil
}

// WatchReleased is a free log subscription operation binding the contract event 0xb21fb52d5749b80f3182f8c6992236b5e5576681880914484d7f4c9b062e619e.
//
// Solidity: event Released(address indexed beneficiary, uint256 amount)
func (_VestingVault *VestingVaultFilterer) WatchReleased(opts *bind.WatchOpts, sink chan<- *VestingVaultReleased, beneficiary []common.Address) (event.Subscription, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _VestingVault.contract.WatchLogs(opts, "Released", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VestingVaultReleased)
				if err := _VestingVault.contract.UnpackLog(event, "Released", log); err != nil {
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

// ParseReleased is a log parse operation binding the contract event 0xb21fb52d5749b80f3182f8c6992236b5e5576681880914484d7f4c9b062e619e.
//
// Solidity: event Released(address indexed beneficiary, uint256 amount)
func (_VestingVault *VestingVaultFilterer) ParseReleased(log types.Log) (*VestingVaultReleased, error) {
	event := new(VestingVaultReleased)
	if err := _VestingVault.contract.UnpackLog(event, "Released", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VestingVaultRevokedIterator is returned from FilterRevoked and is used to iterate over the raw logs and unpacked data for Revoked events raised by the VestingVault contract.
type VestingVaultRevokedIterator struct {
	Event *VestingVaultRevoked // Event containing the contract specifics and raw log

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
func (it *VestingVaultRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VestingVaultRevoked)
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
		it.Event = new(VestingVaultRevoked)
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
func (it *VestingVaultRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VestingVaultRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VestingVaultRevoked represents a Revoked event raised by the VestingVault contract.
type VestingVaultRevoked struct {
	Beneficiary      common.Address
	VestedKept       *big.Int
	UnvestedReturned *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRevoked is a free log retrieval operation binding the contract event 0x486e02409c32040cba4d1675ead85cba092960333a3e8b18cb06e235cef29bfc.
//
// Solidity: event Revoked(address indexed beneficiary, uint256 vestedKept, uint256 unvestedReturned)
func (_VestingVault *VestingVaultFilterer) FilterRevoked(opts *bind.FilterOpts, beneficiary []common.Address) (*VestingVaultRevokedIterator, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _VestingVault.contract.FilterLogs(opts, "Revoked", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return &VestingVaultRevokedIterator{contract: _VestingVault.contract, event: "Revoked", logs: logs, sub: sub}, nil
}

// WatchRevoked is a free log subscription operation binding the contract event 0x486e02409c32040cba4d1675ead85cba092960333a3e8b18cb06e235cef29bfc.
//
// Solidity: event Revoked(address indexed beneficiary, uint256 vestedKept, uint256 unvestedReturned)
func (_VestingVault *VestingVaultFilterer) WatchRevoked(opts *bind.WatchOpts, sink chan<- *VestingVaultRevoked, beneficiary []common.Address) (event.Subscription, error) {

	var beneficiaryRule []interface{}
	for _, beneficiaryItem := range beneficiary {
		beneficiaryRule = append(beneficiaryRule, beneficiaryItem)
	}

	logs, sub, err := _VestingVault.contract.WatchLogs(opts, "Revoked", beneficiaryRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VestingVaultRevoked)
				if err := _VestingVault.contract.UnpackLog(event, "Revoked", log); err != nil {
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

// ParseRevoked is a log parse operation binding the contract event 0x486e02409c32040cba4d1675ead85cba092960333a3e8b18cb06e235cef29bfc.
//
// Solidity: event Revoked(address indexed beneficiary, uint256 vestedKept, uint256 unvestedReturned)
func (_VestingVault *VestingVaultFilterer) ParseRevoked(log types.Log) (*VestingVaultRevoked, error) {
	event := new(VestingVaultRevoked)
	if err := _VestingVault.contract.UnpackLog(event, "Revoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VestingVaultRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the VestingVault contract.
type VestingVaultRoleAdminChangedIterator struct {
	Event *VestingVaultRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *VestingVaultRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VestingVaultRoleAdminChanged)
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
		it.Event = new(VestingVaultRoleAdminChanged)
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
func (it *VestingVaultRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VestingVaultRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VestingVaultRoleAdminChanged represents a RoleAdminChanged event raised by the VestingVault contract.
type VestingVaultRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_VestingVault *VestingVaultFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*VestingVaultRoleAdminChangedIterator, error) {

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

	logs, sub, err := _VestingVault.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &VestingVaultRoleAdminChangedIterator{contract: _VestingVault.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_VestingVault *VestingVaultFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *VestingVaultRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _VestingVault.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VestingVaultRoleAdminChanged)
				if err := _VestingVault.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_VestingVault *VestingVaultFilterer) ParseRoleAdminChanged(log types.Log) (*VestingVaultRoleAdminChanged, error) {
	event := new(VestingVaultRoleAdminChanged)
	if err := _VestingVault.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VestingVaultRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the VestingVault contract.
type VestingVaultRoleGrantedIterator struct {
	Event *VestingVaultRoleGranted // Event containing the contract specifics and raw log

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
func (it *VestingVaultRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VestingVaultRoleGranted)
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
		it.Event = new(VestingVaultRoleGranted)
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
func (it *VestingVaultRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VestingVaultRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VestingVaultRoleGranted represents a RoleGranted event raised by the VestingVault contract.
type VestingVaultRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_VestingVault *VestingVaultFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*VestingVaultRoleGrantedIterator, error) {

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

	logs, sub, err := _VestingVault.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VestingVaultRoleGrantedIterator{contract: _VestingVault.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_VestingVault *VestingVaultFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *VestingVaultRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VestingVault.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VestingVaultRoleGranted)
				if err := _VestingVault.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_VestingVault *VestingVaultFilterer) ParseRoleGranted(log types.Log) (*VestingVaultRoleGranted, error) {
	event := new(VestingVaultRoleGranted)
	if err := _VestingVault.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VestingVaultRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the VestingVault contract.
type VestingVaultRoleRevokedIterator struct {
	Event *VestingVaultRoleRevoked // Event containing the contract specifics and raw log

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
func (it *VestingVaultRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VestingVaultRoleRevoked)
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
		it.Event = new(VestingVaultRoleRevoked)
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
func (it *VestingVaultRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VestingVaultRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VestingVaultRoleRevoked represents a RoleRevoked event raised by the VestingVault contract.
type VestingVaultRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_VestingVault *VestingVaultFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*VestingVaultRoleRevokedIterator, error) {

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

	logs, sub, err := _VestingVault.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &VestingVaultRoleRevokedIterator{contract: _VestingVault.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_VestingVault *VestingVaultFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *VestingVaultRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _VestingVault.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VestingVaultRoleRevoked)
				if err := _VestingVault.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_VestingVault *VestingVaultFilterer) ParseRoleRevoked(log types.Log) (*VestingVaultRoleRevoked, error) {
	event := new(VestingVaultRoleRevoked)
	if err := _VestingVault.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
