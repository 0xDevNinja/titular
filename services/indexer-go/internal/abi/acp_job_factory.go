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

// JobFactoryMetaData contains all meta data concerning the JobFactory contract.
var JobFactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_jobImpl\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_registry\",\"type\":\"address\",\"internalType\":\"contractAgentRegistry\"},{\"name\":\"_defaultArbiter\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAUSER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createJob\",\"inputs\":[{\"name\":\"targetAgentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"budget\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"jobType\",\"type\":\"uint8\",\"internalType\":\"enumIJob.JobType\"},{\"name\":\"evaluator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"arbiter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"clone\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"defaultArbiter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getJob\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"clone\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"jobImplementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"jobs\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"jobsByPrincipal\",\"inputs\":[{\"name\":\"principal\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"jobIds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractAgentRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDefaultArbiter\",\"inputs\":[{\"name\":\"newArbiter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setJobImplementation\",\"inputs\":[{\"name\":\"newImpl\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalJobs\",\"inputs\":[],\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"DefaultArbiterUpdated\",\"inputs\":[{\"name\":\"oldArbiter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newArbiter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ImplementationUpdated\",\"inputs\":[{\"name\":\"oldImpl\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newImpl\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"JobCreated\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"clone\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"principal\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"jobType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIJob.JobType\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"budget\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedDeployment\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidDeadline\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// JobFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use JobFactoryMetaData.ABI instead.
var JobFactoryABI = JobFactoryMetaData.ABI

// JobFactory is an auto generated Go binding around an Ethereum contract.
type JobFactory struct {
	JobFactoryCaller     // Read-only binding to the contract
	JobFactoryTransactor // Write-only binding to the contract
	JobFactoryFilterer   // Log filterer for contract events
}

// JobFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type JobFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// JobFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type JobFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// JobFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type JobFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// JobFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type JobFactorySession struct {
	Contract     *JobFactory       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// JobFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type JobFactoryCallerSession struct {
	Contract *JobFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// JobFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type JobFactoryTransactorSession struct {
	Contract     *JobFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// JobFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type JobFactoryRaw struct {
	Contract *JobFactory // Generic contract binding to access the raw methods on
}

// JobFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type JobFactoryCallerRaw struct {
	Contract *JobFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// JobFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type JobFactoryTransactorRaw struct {
	Contract *JobFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewJobFactory creates a new instance of JobFactory, bound to a specific deployed contract.
func NewJobFactory(address common.Address, backend bind.ContractBackend) (*JobFactory, error) {
	contract, err := bindJobFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &JobFactory{JobFactoryCaller: JobFactoryCaller{contract: contract}, JobFactoryTransactor: JobFactoryTransactor{contract: contract}, JobFactoryFilterer: JobFactoryFilterer{contract: contract}}, nil
}

// NewJobFactoryCaller creates a new read-only instance of JobFactory, bound to a specific deployed contract.
func NewJobFactoryCaller(address common.Address, caller bind.ContractCaller) (*JobFactoryCaller, error) {
	contract, err := bindJobFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &JobFactoryCaller{contract: contract}, nil
}

// NewJobFactoryTransactor creates a new write-only instance of JobFactory, bound to a specific deployed contract.
func NewJobFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*JobFactoryTransactor, error) {
	contract, err := bindJobFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &JobFactoryTransactor{contract: contract}, nil
}

// NewJobFactoryFilterer creates a new log filterer instance of JobFactory, bound to a specific deployed contract.
func NewJobFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*JobFactoryFilterer, error) {
	contract, err := bindJobFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &JobFactoryFilterer{contract: contract}, nil
}

// bindJobFactory binds a generic wrapper to an already deployed contract.
func bindJobFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := JobFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_JobFactory *JobFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _JobFactory.Contract.JobFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_JobFactory *JobFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _JobFactory.Contract.JobFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_JobFactory *JobFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _JobFactory.Contract.JobFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_JobFactory *JobFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _JobFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_JobFactory *JobFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _JobFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_JobFactory *JobFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _JobFactory.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_JobFactory *JobFactoryCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_JobFactory *JobFactorySession) DEFAULTADMINROLE() ([32]byte, error) {
	return _JobFactory.Contract.DEFAULTADMINROLE(&_JobFactory.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_JobFactory *JobFactoryCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _JobFactory.Contract.DEFAULTADMINROLE(&_JobFactory.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_JobFactory *JobFactoryCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_JobFactory *JobFactorySession) PAUSERROLE() ([32]byte, error) {
	return _JobFactory.Contract.PAUSERROLE(&_JobFactory.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_JobFactory *JobFactoryCallerSession) PAUSERROLE() ([32]byte, error) {
	return _JobFactory.Contract.PAUSERROLE(&_JobFactory.CallOpts)
}

// DefaultArbiter is a free data retrieval call binding the contract method 0xb67a2637.
//
// Solidity: function defaultArbiter() view returns(address)
func (_JobFactory *JobFactoryCaller) DefaultArbiter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "defaultArbiter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DefaultArbiter is a free data retrieval call binding the contract method 0xb67a2637.
//
// Solidity: function defaultArbiter() view returns(address)
func (_JobFactory *JobFactorySession) DefaultArbiter() (common.Address, error) {
	return _JobFactory.Contract.DefaultArbiter(&_JobFactory.CallOpts)
}

// DefaultArbiter is a free data retrieval call binding the contract method 0xb67a2637.
//
// Solidity: function defaultArbiter() view returns(address)
func (_JobFactory *JobFactoryCallerSession) DefaultArbiter() (common.Address, error) {
	return _JobFactory.Contract.DefaultArbiter(&_JobFactory.CallOpts)
}

// GetJob is a free data retrieval call binding the contract method 0xbf22c457.
//
// Solidity: function getJob(uint256 jobId) view returns(address clone)
func (_JobFactory *JobFactoryCaller) GetJob(opts *bind.CallOpts, jobId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "getJob", jobId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetJob is a free data retrieval call binding the contract method 0xbf22c457.
//
// Solidity: function getJob(uint256 jobId) view returns(address clone)
func (_JobFactory *JobFactorySession) GetJob(jobId *big.Int) (common.Address, error) {
	return _JobFactory.Contract.GetJob(&_JobFactory.CallOpts, jobId)
}

// GetJob is a free data retrieval call binding the contract method 0xbf22c457.
//
// Solidity: function getJob(uint256 jobId) view returns(address clone)
func (_JobFactory *JobFactoryCallerSession) GetJob(jobId *big.Int) (common.Address, error) {
	return _JobFactory.Contract.GetJob(&_JobFactory.CallOpts, jobId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_JobFactory *JobFactoryCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_JobFactory *JobFactorySession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _JobFactory.Contract.GetRoleAdmin(&_JobFactory.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_JobFactory *JobFactoryCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _JobFactory.Contract.GetRoleAdmin(&_JobFactory.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_JobFactory *JobFactoryCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_JobFactory *JobFactorySession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _JobFactory.Contract.HasRole(&_JobFactory.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_JobFactory *JobFactoryCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _JobFactory.Contract.HasRole(&_JobFactory.CallOpts, role, account)
}

// JobImplementation is a free data retrieval call binding the contract method 0x25594a40.
//
// Solidity: function jobImplementation() view returns(address)
func (_JobFactory *JobFactoryCaller) JobImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "jobImplementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// JobImplementation is a free data retrieval call binding the contract method 0x25594a40.
//
// Solidity: function jobImplementation() view returns(address)
func (_JobFactory *JobFactorySession) JobImplementation() (common.Address, error) {
	return _JobFactory.Contract.JobImplementation(&_JobFactory.CallOpts)
}

// JobImplementation is a free data retrieval call binding the contract method 0x25594a40.
//
// Solidity: function jobImplementation() view returns(address)
func (_JobFactory *JobFactoryCallerSession) JobImplementation() (common.Address, error) {
	return _JobFactory.Contract.JobImplementation(&_JobFactory.CallOpts)
}

// Jobs is a free data retrieval call binding the contract method 0x180aedf3.
//
// Solidity: function jobs(uint256 ) view returns(address)
func (_JobFactory *JobFactoryCaller) Jobs(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "jobs", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Jobs is a free data retrieval call binding the contract method 0x180aedf3.
//
// Solidity: function jobs(uint256 ) view returns(address)
func (_JobFactory *JobFactorySession) Jobs(arg0 *big.Int) (common.Address, error) {
	return _JobFactory.Contract.Jobs(&_JobFactory.CallOpts, arg0)
}

// Jobs is a free data retrieval call binding the contract method 0x180aedf3.
//
// Solidity: function jobs(uint256 ) view returns(address)
func (_JobFactory *JobFactoryCallerSession) Jobs(arg0 *big.Int) (common.Address, error) {
	return _JobFactory.Contract.Jobs(&_JobFactory.CallOpts, arg0)
}

// JobsByPrincipal is a free data retrieval call binding the contract method 0xa666d299.
//
// Solidity: function jobsByPrincipal(address principal) view returns(uint256[] jobIds)
func (_JobFactory *JobFactoryCaller) JobsByPrincipal(opts *bind.CallOpts, principal common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "jobsByPrincipal", principal)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// JobsByPrincipal is a free data retrieval call binding the contract method 0xa666d299.
//
// Solidity: function jobsByPrincipal(address principal) view returns(uint256[] jobIds)
func (_JobFactory *JobFactorySession) JobsByPrincipal(principal common.Address) ([]*big.Int, error) {
	return _JobFactory.Contract.JobsByPrincipal(&_JobFactory.CallOpts, principal)
}

// JobsByPrincipal is a free data retrieval call binding the contract method 0xa666d299.
//
// Solidity: function jobsByPrincipal(address principal) view returns(uint256[] jobIds)
func (_JobFactory *JobFactoryCallerSession) JobsByPrincipal(principal common.Address) ([]*big.Int, error) {
	return _JobFactory.Contract.JobsByPrincipal(&_JobFactory.CallOpts, principal)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_JobFactory *JobFactoryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_JobFactory *JobFactorySession) Paused() (bool, error) {
	return _JobFactory.Contract.Paused(&_JobFactory.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_JobFactory *JobFactoryCallerSession) Paused() (bool, error) {
	return _JobFactory.Contract.Paused(&_JobFactory.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_JobFactory *JobFactoryCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_JobFactory *JobFactorySession) Registry() (common.Address, error) {
	return _JobFactory.Contract.Registry(&_JobFactory.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_JobFactory *JobFactoryCallerSession) Registry() (common.Address, error) {
	return _JobFactory.Contract.Registry(&_JobFactory.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_JobFactory *JobFactoryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_JobFactory *JobFactorySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _JobFactory.Contract.SupportsInterface(&_JobFactory.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_JobFactory *JobFactoryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _JobFactory.Contract.SupportsInterface(&_JobFactory.CallOpts, interfaceId)
}

// TotalJobs is a free data retrieval call binding the contract method 0x1ace87b3.
//
// Solidity: function totalJobs() view returns(uint256 count)
func (_JobFactory *JobFactoryCaller) TotalJobs(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _JobFactory.contract.Call(opts, &out, "totalJobs")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalJobs is a free data retrieval call binding the contract method 0x1ace87b3.
//
// Solidity: function totalJobs() view returns(uint256 count)
func (_JobFactory *JobFactorySession) TotalJobs() (*big.Int, error) {
	return _JobFactory.Contract.TotalJobs(&_JobFactory.CallOpts)
}

// TotalJobs is a free data retrieval call binding the contract method 0x1ace87b3.
//
// Solidity: function totalJobs() view returns(uint256 count)
func (_JobFactory *JobFactoryCallerSession) TotalJobs() (*big.Int, error) {
	return _JobFactory.Contract.TotalJobs(&_JobFactory.CallOpts)
}

// CreateJob is a paid mutator transaction binding the contract method 0xaa6110ab.
//
// Solidity: function createJob(uint256 targetAgentId, address token, uint256 budget, uint64 deadline, uint8 jobType, address evaluator, address arbiter) returns(uint256 jobId, address clone)
func (_JobFactory *JobFactoryTransactor) CreateJob(opts *bind.TransactOpts, targetAgentId *big.Int, token common.Address, budget *big.Int, deadline uint64, jobType uint8, evaluator common.Address, arbiter common.Address) (*types.Transaction, error) {
	return _JobFactory.contract.Transact(opts, "createJob", targetAgentId, token, budget, deadline, jobType, evaluator, arbiter)
}

// CreateJob is a paid mutator transaction binding the contract method 0xaa6110ab.
//
// Solidity: function createJob(uint256 targetAgentId, address token, uint256 budget, uint64 deadline, uint8 jobType, address evaluator, address arbiter) returns(uint256 jobId, address clone)
func (_JobFactory *JobFactorySession) CreateJob(targetAgentId *big.Int, token common.Address, budget *big.Int, deadline uint64, jobType uint8, evaluator common.Address, arbiter common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.CreateJob(&_JobFactory.TransactOpts, targetAgentId, token, budget, deadline, jobType, evaluator, arbiter)
}

// CreateJob is a paid mutator transaction binding the contract method 0xaa6110ab.
//
// Solidity: function createJob(uint256 targetAgentId, address token, uint256 budget, uint64 deadline, uint8 jobType, address evaluator, address arbiter) returns(uint256 jobId, address clone)
func (_JobFactory *JobFactoryTransactorSession) CreateJob(targetAgentId *big.Int, token common.Address, budget *big.Int, deadline uint64, jobType uint8, evaluator common.Address, arbiter common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.CreateJob(&_JobFactory.TransactOpts, targetAgentId, token, budget, deadline, jobType, evaluator, arbiter)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_JobFactory *JobFactoryTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _JobFactory.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_JobFactory *JobFactorySession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.GrantRole(&_JobFactory.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_JobFactory *JobFactoryTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.GrantRole(&_JobFactory.TransactOpts, role, account)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_JobFactory *JobFactoryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _JobFactory.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_JobFactory *JobFactorySession) Pause() (*types.Transaction, error) {
	return _JobFactory.Contract.Pause(&_JobFactory.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_JobFactory *JobFactoryTransactorSession) Pause() (*types.Transaction, error) {
	return _JobFactory.Contract.Pause(&_JobFactory.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_JobFactory *JobFactoryTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _JobFactory.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_JobFactory *JobFactorySession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.RenounceRole(&_JobFactory.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_JobFactory *JobFactoryTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.RenounceRole(&_JobFactory.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_JobFactory *JobFactoryTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _JobFactory.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_JobFactory *JobFactorySession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.RevokeRole(&_JobFactory.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_JobFactory *JobFactoryTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.RevokeRole(&_JobFactory.TransactOpts, role, account)
}

// SetDefaultArbiter is a paid mutator transaction binding the contract method 0xbbdf7243.
//
// Solidity: function setDefaultArbiter(address newArbiter) returns()
func (_JobFactory *JobFactoryTransactor) SetDefaultArbiter(opts *bind.TransactOpts, newArbiter common.Address) (*types.Transaction, error) {
	return _JobFactory.contract.Transact(opts, "setDefaultArbiter", newArbiter)
}

// SetDefaultArbiter is a paid mutator transaction binding the contract method 0xbbdf7243.
//
// Solidity: function setDefaultArbiter(address newArbiter) returns()
func (_JobFactory *JobFactorySession) SetDefaultArbiter(newArbiter common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.SetDefaultArbiter(&_JobFactory.TransactOpts, newArbiter)
}

// SetDefaultArbiter is a paid mutator transaction binding the contract method 0xbbdf7243.
//
// Solidity: function setDefaultArbiter(address newArbiter) returns()
func (_JobFactory *JobFactoryTransactorSession) SetDefaultArbiter(newArbiter common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.SetDefaultArbiter(&_JobFactory.TransactOpts, newArbiter)
}

// SetJobImplementation is a paid mutator transaction binding the contract method 0x83c7cfc9.
//
// Solidity: function setJobImplementation(address newImpl) returns()
func (_JobFactory *JobFactoryTransactor) SetJobImplementation(opts *bind.TransactOpts, newImpl common.Address) (*types.Transaction, error) {
	return _JobFactory.contract.Transact(opts, "setJobImplementation", newImpl)
}

// SetJobImplementation is a paid mutator transaction binding the contract method 0x83c7cfc9.
//
// Solidity: function setJobImplementation(address newImpl) returns()
func (_JobFactory *JobFactorySession) SetJobImplementation(newImpl common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.SetJobImplementation(&_JobFactory.TransactOpts, newImpl)
}

// SetJobImplementation is a paid mutator transaction binding the contract method 0x83c7cfc9.
//
// Solidity: function setJobImplementation(address newImpl) returns()
func (_JobFactory *JobFactoryTransactorSession) SetJobImplementation(newImpl common.Address) (*types.Transaction, error) {
	return _JobFactory.Contract.SetJobImplementation(&_JobFactory.TransactOpts, newImpl)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_JobFactory *JobFactoryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _JobFactory.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_JobFactory *JobFactorySession) Unpause() (*types.Transaction, error) {
	return _JobFactory.Contract.Unpause(&_JobFactory.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_JobFactory *JobFactoryTransactorSession) Unpause() (*types.Transaction, error) {
	return _JobFactory.Contract.Unpause(&_JobFactory.TransactOpts)
}

// JobFactoryDefaultArbiterUpdatedIterator is returned from FilterDefaultArbiterUpdated and is used to iterate over the raw logs and unpacked data for DefaultArbiterUpdated events raised by the JobFactory contract.
type JobFactoryDefaultArbiterUpdatedIterator struct {
	Event *JobFactoryDefaultArbiterUpdated // Event containing the contract specifics and raw log

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
func (it *JobFactoryDefaultArbiterUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobFactoryDefaultArbiterUpdated)
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
		it.Event = new(JobFactoryDefaultArbiterUpdated)
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
func (it *JobFactoryDefaultArbiterUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobFactoryDefaultArbiterUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobFactoryDefaultArbiterUpdated represents a DefaultArbiterUpdated event raised by the JobFactory contract.
type JobFactoryDefaultArbiterUpdated struct {
	OldArbiter common.Address
	NewArbiter common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDefaultArbiterUpdated is a free log retrieval operation binding the contract event 0x8f4ec914918528531be3ff56f55a057f73eb540484e6fe6f63c7bd02a483a957.
//
// Solidity: event DefaultArbiterUpdated(address indexed oldArbiter, address indexed newArbiter)
func (_JobFactory *JobFactoryFilterer) FilterDefaultArbiterUpdated(opts *bind.FilterOpts, oldArbiter []common.Address, newArbiter []common.Address) (*JobFactoryDefaultArbiterUpdatedIterator, error) {

	var oldArbiterRule []interface{}
	for _, oldArbiterItem := range oldArbiter {
		oldArbiterRule = append(oldArbiterRule, oldArbiterItem)
	}
	var newArbiterRule []interface{}
	for _, newArbiterItem := range newArbiter {
		newArbiterRule = append(newArbiterRule, newArbiterItem)
	}

	logs, sub, err := _JobFactory.contract.FilterLogs(opts, "DefaultArbiterUpdated", oldArbiterRule, newArbiterRule)
	if err != nil {
		return nil, err
	}
	return &JobFactoryDefaultArbiterUpdatedIterator{contract: _JobFactory.contract, event: "DefaultArbiterUpdated", logs: logs, sub: sub}, nil
}

// WatchDefaultArbiterUpdated is a free log subscription operation binding the contract event 0x8f4ec914918528531be3ff56f55a057f73eb540484e6fe6f63c7bd02a483a957.
//
// Solidity: event DefaultArbiterUpdated(address indexed oldArbiter, address indexed newArbiter)
func (_JobFactory *JobFactoryFilterer) WatchDefaultArbiterUpdated(opts *bind.WatchOpts, sink chan<- *JobFactoryDefaultArbiterUpdated, oldArbiter []common.Address, newArbiter []common.Address) (event.Subscription, error) {

	var oldArbiterRule []interface{}
	for _, oldArbiterItem := range oldArbiter {
		oldArbiterRule = append(oldArbiterRule, oldArbiterItem)
	}
	var newArbiterRule []interface{}
	for _, newArbiterItem := range newArbiter {
		newArbiterRule = append(newArbiterRule, newArbiterItem)
	}

	logs, sub, err := _JobFactory.contract.WatchLogs(opts, "DefaultArbiterUpdated", oldArbiterRule, newArbiterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobFactoryDefaultArbiterUpdated)
				if err := _JobFactory.contract.UnpackLog(event, "DefaultArbiterUpdated", log); err != nil {
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

// ParseDefaultArbiterUpdated is a log parse operation binding the contract event 0x8f4ec914918528531be3ff56f55a057f73eb540484e6fe6f63c7bd02a483a957.
//
// Solidity: event DefaultArbiterUpdated(address indexed oldArbiter, address indexed newArbiter)
func (_JobFactory *JobFactoryFilterer) ParseDefaultArbiterUpdated(log types.Log) (*JobFactoryDefaultArbiterUpdated, error) {
	event := new(JobFactoryDefaultArbiterUpdated)
	if err := _JobFactory.contract.UnpackLog(event, "DefaultArbiterUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobFactoryImplementationUpdatedIterator is returned from FilterImplementationUpdated and is used to iterate over the raw logs and unpacked data for ImplementationUpdated events raised by the JobFactory contract.
type JobFactoryImplementationUpdatedIterator struct {
	Event *JobFactoryImplementationUpdated // Event containing the contract specifics and raw log

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
func (it *JobFactoryImplementationUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobFactoryImplementationUpdated)
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
		it.Event = new(JobFactoryImplementationUpdated)
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
func (it *JobFactoryImplementationUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobFactoryImplementationUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobFactoryImplementationUpdated represents a ImplementationUpdated event raised by the JobFactory contract.
type JobFactoryImplementationUpdated struct {
	OldImpl common.Address
	NewImpl common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterImplementationUpdated is a free log retrieval operation binding the contract event 0xaa3f731066a578e5f39b4215468d826cdd15373cbc0dfc9cb9bdc649718ef7da.
//
// Solidity: event ImplementationUpdated(address indexed oldImpl, address indexed newImpl)
func (_JobFactory *JobFactoryFilterer) FilterImplementationUpdated(opts *bind.FilterOpts, oldImpl []common.Address, newImpl []common.Address) (*JobFactoryImplementationUpdatedIterator, error) {

	var oldImplRule []interface{}
	for _, oldImplItem := range oldImpl {
		oldImplRule = append(oldImplRule, oldImplItem)
	}
	var newImplRule []interface{}
	for _, newImplItem := range newImpl {
		newImplRule = append(newImplRule, newImplItem)
	}

	logs, sub, err := _JobFactory.contract.FilterLogs(opts, "ImplementationUpdated", oldImplRule, newImplRule)
	if err != nil {
		return nil, err
	}
	return &JobFactoryImplementationUpdatedIterator{contract: _JobFactory.contract, event: "ImplementationUpdated", logs: logs, sub: sub}, nil
}

// WatchImplementationUpdated is a free log subscription operation binding the contract event 0xaa3f731066a578e5f39b4215468d826cdd15373cbc0dfc9cb9bdc649718ef7da.
//
// Solidity: event ImplementationUpdated(address indexed oldImpl, address indexed newImpl)
func (_JobFactory *JobFactoryFilterer) WatchImplementationUpdated(opts *bind.WatchOpts, sink chan<- *JobFactoryImplementationUpdated, oldImpl []common.Address, newImpl []common.Address) (event.Subscription, error) {

	var oldImplRule []interface{}
	for _, oldImplItem := range oldImpl {
		oldImplRule = append(oldImplRule, oldImplItem)
	}
	var newImplRule []interface{}
	for _, newImplItem := range newImpl {
		newImplRule = append(newImplRule, newImplItem)
	}

	logs, sub, err := _JobFactory.contract.WatchLogs(opts, "ImplementationUpdated", oldImplRule, newImplRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobFactoryImplementationUpdated)
				if err := _JobFactory.contract.UnpackLog(event, "ImplementationUpdated", log); err != nil {
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

// ParseImplementationUpdated is a log parse operation binding the contract event 0xaa3f731066a578e5f39b4215468d826cdd15373cbc0dfc9cb9bdc649718ef7da.
//
// Solidity: event ImplementationUpdated(address indexed oldImpl, address indexed newImpl)
func (_JobFactory *JobFactoryFilterer) ParseImplementationUpdated(log types.Log) (*JobFactoryImplementationUpdated, error) {
	event := new(JobFactoryImplementationUpdated)
	if err := _JobFactory.contract.UnpackLog(event, "ImplementationUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobFactoryJobCreatedIterator is returned from FilterJobCreated and is used to iterate over the raw logs and unpacked data for JobCreated events raised by the JobFactory contract.
type JobFactoryJobCreatedIterator struct {
	Event *JobFactoryJobCreated // Event containing the contract specifics and raw log

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
func (it *JobFactoryJobCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobFactoryJobCreated)
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
		it.Event = new(JobFactoryJobCreated)
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
func (it *JobFactoryJobCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobFactoryJobCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobFactoryJobCreated represents a JobCreated event raised by the JobFactory contract.
type JobFactoryJobCreated struct {
	JobId     *big.Int
	Clone     common.Address
	Principal common.Address
	JobType   uint8
	Token     common.Address
	Budget    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterJobCreated is a free log retrieval operation binding the contract event 0xf956496ab0e37cf09cc3df7fd2b643c4a458b72e0d06e6cb9544a68553fabdd5.
//
// Solidity: event JobCreated(uint256 indexed jobId, address indexed clone, address indexed principal, uint8 jobType, address token, uint256 budget)
func (_JobFactory *JobFactoryFilterer) FilterJobCreated(opts *bind.FilterOpts, jobId []*big.Int, clone []common.Address, principal []common.Address) (*JobFactoryJobCreatedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var cloneRule []interface{}
	for _, cloneItem := range clone {
		cloneRule = append(cloneRule, cloneItem)
	}
	var principalRule []interface{}
	for _, principalItem := range principal {
		principalRule = append(principalRule, principalItem)
	}

	logs, sub, err := _JobFactory.contract.FilterLogs(opts, "JobCreated", jobIdRule, cloneRule, principalRule)
	if err != nil {
		return nil, err
	}
	return &JobFactoryJobCreatedIterator{contract: _JobFactory.contract, event: "JobCreated", logs: logs, sub: sub}, nil
}

// WatchJobCreated is a free log subscription operation binding the contract event 0xf956496ab0e37cf09cc3df7fd2b643c4a458b72e0d06e6cb9544a68553fabdd5.
//
// Solidity: event JobCreated(uint256 indexed jobId, address indexed clone, address indexed principal, uint8 jobType, address token, uint256 budget)
func (_JobFactory *JobFactoryFilterer) WatchJobCreated(opts *bind.WatchOpts, sink chan<- *JobFactoryJobCreated, jobId []*big.Int, clone []common.Address, principal []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var cloneRule []interface{}
	for _, cloneItem := range clone {
		cloneRule = append(cloneRule, cloneItem)
	}
	var principalRule []interface{}
	for _, principalItem := range principal {
		principalRule = append(principalRule, principalItem)
	}

	logs, sub, err := _JobFactory.contract.WatchLogs(opts, "JobCreated", jobIdRule, cloneRule, principalRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobFactoryJobCreated)
				if err := _JobFactory.contract.UnpackLog(event, "JobCreated", log); err != nil {
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

// ParseJobCreated is a log parse operation binding the contract event 0xf956496ab0e37cf09cc3df7fd2b643c4a458b72e0d06e6cb9544a68553fabdd5.
//
// Solidity: event JobCreated(uint256 indexed jobId, address indexed clone, address indexed principal, uint8 jobType, address token, uint256 budget)
func (_JobFactory *JobFactoryFilterer) ParseJobCreated(log types.Log) (*JobFactoryJobCreated, error) {
	event := new(JobFactoryJobCreated)
	if err := _JobFactory.contract.UnpackLog(event, "JobCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobFactoryPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the JobFactory contract.
type JobFactoryPausedIterator struct {
	Event *JobFactoryPaused // Event containing the contract specifics and raw log

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
func (it *JobFactoryPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobFactoryPaused)
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
		it.Event = new(JobFactoryPaused)
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
func (it *JobFactoryPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobFactoryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobFactoryPaused represents a Paused event raised by the JobFactory contract.
type JobFactoryPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_JobFactory *JobFactoryFilterer) FilterPaused(opts *bind.FilterOpts) (*JobFactoryPausedIterator, error) {

	logs, sub, err := _JobFactory.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &JobFactoryPausedIterator{contract: _JobFactory.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_JobFactory *JobFactoryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *JobFactoryPaused) (event.Subscription, error) {

	logs, sub, err := _JobFactory.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobFactoryPaused)
				if err := _JobFactory.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_JobFactory *JobFactoryFilterer) ParsePaused(log types.Log) (*JobFactoryPaused, error) {
	event := new(JobFactoryPaused)
	if err := _JobFactory.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobFactoryRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the JobFactory contract.
type JobFactoryRoleAdminChangedIterator struct {
	Event *JobFactoryRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *JobFactoryRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobFactoryRoleAdminChanged)
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
		it.Event = new(JobFactoryRoleAdminChanged)
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
func (it *JobFactoryRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobFactoryRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobFactoryRoleAdminChanged represents a RoleAdminChanged event raised by the JobFactory contract.
type JobFactoryRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_JobFactory *JobFactoryFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*JobFactoryRoleAdminChangedIterator, error) {

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

	logs, sub, err := _JobFactory.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &JobFactoryRoleAdminChangedIterator{contract: _JobFactory.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_JobFactory *JobFactoryFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *JobFactoryRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _JobFactory.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobFactoryRoleAdminChanged)
				if err := _JobFactory.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_JobFactory *JobFactoryFilterer) ParseRoleAdminChanged(log types.Log) (*JobFactoryRoleAdminChanged, error) {
	event := new(JobFactoryRoleAdminChanged)
	if err := _JobFactory.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobFactoryRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the JobFactory contract.
type JobFactoryRoleGrantedIterator struct {
	Event *JobFactoryRoleGranted // Event containing the contract specifics and raw log

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
func (it *JobFactoryRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobFactoryRoleGranted)
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
		it.Event = new(JobFactoryRoleGranted)
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
func (it *JobFactoryRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobFactoryRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobFactoryRoleGranted represents a RoleGranted event raised by the JobFactory contract.
type JobFactoryRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_JobFactory *JobFactoryFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*JobFactoryRoleGrantedIterator, error) {

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

	logs, sub, err := _JobFactory.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &JobFactoryRoleGrantedIterator{contract: _JobFactory.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_JobFactory *JobFactoryFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *JobFactoryRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _JobFactory.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobFactoryRoleGranted)
				if err := _JobFactory.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_JobFactory *JobFactoryFilterer) ParseRoleGranted(log types.Log) (*JobFactoryRoleGranted, error) {
	event := new(JobFactoryRoleGranted)
	if err := _JobFactory.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobFactoryRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the JobFactory contract.
type JobFactoryRoleRevokedIterator struct {
	Event *JobFactoryRoleRevoked // Event containing the contract specifics and raw log

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
func (it *JobFactoryRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobFactoryRoleRevoked)
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
		it.Event = new(JobFactoryRoleRevoked)
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
func (it *JobFactoryRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobFactoryRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobFactoryRoleRevoked represents a RoleRevoked event raised by the JobFactory contract.
type JobFactoryRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_JobFactory *JobFactoryFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*JobFactoryRoleRevokedIterator, error) {

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

	logs, sub, err := _JobFactory.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &JobFactoryRoleRevokedIterator{contract: _JobFactory.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_JobFactory *JobFactoryFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *JobFactoryRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _JobFactory.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobFactoryRoleRevoked)
				if err := _JobFactory.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_JobFactory *JobFactoryFilterer) ParseRoleRevoked(log types.Log) (*JobFactoryRoleRevoked, error) {
	event := new(JobFactoryRoleRevoked)
	if err := _JobFactory.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobFactoryUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the JobFactory contract.
type JobFactoryUnpausedIterator struct {
	Event *JobFactoryUnpaused // Event containing the contract specifics and raw log

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
func (it *JobFactoryUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobFactoryUnpaused)
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
		it.Event = new(JobFactoryUnpaused)
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
func (it *JobFactoryUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobFactoryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobFactoryUnpaused represents a Unpaused event raised by the JobFactory contract.
type JobFactoryUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_JobFactory *JobFactoryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*JobFactoryUnpausedIterator, error) {

	logs, sub, err := _JobFactory.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &JobFactoryUnpausedIterator{contract: _JobFactory.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_JobFactory *JobFactoryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *JobFactoryUnpaused) (event.Subscription, error) {

	logs, sub, err := _JobFactory.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobFactoryUnpaused)
				if err := _JobFactory.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_JobFactory *JobFactoryFilterer) ParseUnpaused(log types.Log) (*JobFactoryUnpaused, error) {
	event := new(JobFactoryUnpaused)
	if err := _JobFactory.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
