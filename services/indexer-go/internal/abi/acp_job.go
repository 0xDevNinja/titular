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

// JobInitParams is an auto generated low-level Go binding around an user-defined struct.
type JobInitParams struct {
	JobId         *big.Int
	Principal     common.Address
	Registry      common.Address
	TargetAgentId *big.Int
	Token         common.Address
	Budget        *big.Int
	Deadline      uint64
	JobType       uint8
	Evaluator     common.Address
	Arbiter       common.Address
}

// JobMetaData contains all meta data concerning the Job contract.
var JobMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"DISPUTE_GRACE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"RELEASE_GRACE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"accept\",\"inputs\":[{\"name\":\"_agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"agent\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"agentId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approveResult\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"arbiter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"budget\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancel\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentPhase\",\"inputs\":[],\"outputs\":[{\"name\":\"current\",\"type\":\"uint8\",\"internalType\":\"enumIJob.Phase\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deadline\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"evaluator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"expireJob\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"p\",\"type\":\"tuple\",\"internalType\":\"structJob.InitParams\",\"components\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"principal\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"registry\",\"type\":\"address\",\"internalType\":\"contractAgentRegistry\"},{\"name\":\"targetAgentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"budget\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"jobType\",\"type\":\"uint8\",\"internalType\":\"enumIJob.JobType\"},{\"name\":\"evaluator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"arbiter\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"jobId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"jobType\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIJob.JobType\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"phase\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIJob.Phase\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"principal\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"raiseDispute\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractAgentRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rejectResult\",\"inputs\":[{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"release\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resolveDispute\",\"inputs\":[{\"name\":\"agentFavoured\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"resultSubmittedAt\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resultURI\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"submitResult\",\"inputs\":[{\"name\":\"_resultURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"targetAgentId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AgentAccepted\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DisputeRaised\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"raisedBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DisputeResolved\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"resolver\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agentFavoured\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EvaluatorAssigned\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"evaluator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"JobCancelled\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"cancelledBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"JobCompleted\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"releasedBy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"JobInitialised\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"principal\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"jobType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIJob.JobType\"},{\"name\":\"budget\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"deadline\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ResultApproved\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"evaluator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ResultRejected\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"evaluator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ResultSubmitted\",\"inputs\":[{\"name\":\"jobId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"resultURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AgentInactive\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"AgentNotFound\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"AlreadyInitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DirectJobNoEvaluator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EvaluatorRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GracePeriodActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidDeadline\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidJobType\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"JobNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAgent\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotArbiter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEvaluator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotPrincipal\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongPhase\",\"inputs\":[{\"name\":\"current\",\"type\":\"uint8\",\"internalType\":\"enumIJob.Phase\"},{\"name\":\"required\",\"type\":\"uint8\",\"internalType\":\"enumIJob.Phase\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// JobABI is the input ABI used to generate the binding from.
// Deprecated: Use JobMetaData.ABI instead.
var JobABI = JobMetaData.ABI

// Job is an auto generated Go binding around an Ethereum contract.
type Job struct {
	JobCaller     // Read-only binding to the contract
	JobTransactor // Write-only binding to the contract
	JobFilterer   // Log filterer for contract events
}

// JobCaller is an auto generated read-only Go binding around an Ethereum contract.
type JobCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// JobTransactor is an auto generated write-only Go binding around an Ethereum contract.
type JobTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// JobFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type JobFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// JobSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type JobSession struct {
	Contract     *Job              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// JobCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type JobCallerSession struct {
	Contract *JobCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// JobTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type JobTransactorSession struct {
	Contract     *JobTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// JobRaw is an auto generated low-level Go binding around an Ethereum contract.
type JobRaw struct {
	Contract *Job // Generic contract binding to access the raw methods on
}

// JobCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type JobCallerRaw struct {
	Contract *JobCaller // Generic read-only contract binding to access the raw methods on
}

// JobTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type JobTransactorRaw struct {
	Contract *JobTransactor // Generic write-only contract binding to access the raw methods on
}

// NewJob creates a new instance of Job, bound to a specific deployed contract.
func NewJob(address common.Address, backend bind.ContractBackend) (*Job, error) {
	contract, err := bindJob(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Job{JobCaller: JobCaller{contract: contract}, JobTransactor: JobTransactor{contract: contract}, JobFilterer: JobFilterer{contract: contract}}, nil
}

// NewJobCaller creates a new read-only instance of Job, bound to a specific deployed contract.
func NewJobCaller(address common.Address, caller bind.ContractCaller) (*JobCaller, error) {
	contract, err := bindJob(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &JobCaller{contract: contract}, nil
}

// NewJobTransactor creates a new write-only instance of Job, bound to a specific deployed contract.
func NewJobTransactor(address common.Address, transactor bind.ContractTransactor) (*JobTransactor, error) {
	contract, err := bindJob(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &JobTransactor{contract: contract}, nil
}

// NewJobFilterer creates a new log filterer instance of Job, bound to a specific deployed contract.
func NewJobFilterer(address common.Address, filterer bind.ContractFilterer) (*JobFilterer, error) {
	contract, err := bindJob(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &JobFilterer{contract: contract}, nil
}

// bindJob binds a generic wrapper to an already deployed contract.
func bindJob(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := JobMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Job *JobRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Job.Contract.JobCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Job *JobRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Job.Contract.JobTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Job *JobRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Job.Contract.JobTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Job *JobCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Job.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Job *JobTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Job.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Job *JobTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Job.Contract.contract.Transact(opts, method, params...)
}

// DISPUTEGRACE is a free data retrieval call binding the contract method 0x276a556e.
//
// Solidity: function DISPUTE_GRACE() view returns(uint64)
func (_Job *JobCaller) DISPUTEGRACE(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "DISPUTE_GRACE")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// DISPUTEGRACE is a free data retrieval call binding the contract method 0x276a556e.
//
// Solidity: function DISPUTE_GRACE() view returns(uint64)
func (_Job *JobSession) DISPUTEGRACE() (uint64, error) {
	return _Job.Contract.DISPUTEGRACE(&_Job.CallOpts)
}

// DISPUTEGRACE is a free data retrieval call binding the contract method 0x276a556e.
//
// Solidity: function DISPUTE_GRACE() view returns(uint64)
func (_Job *JobCallerSession) DISPUTEGRACE() (uint64, error) {
	return _Job.Contract.DISPUTEGRACE(&_Job.CallOpts)
}

// RELEASEGRACE is a free data retrieval call binding the contract method 0xe73ad7c6.
//
// Solidity: function RELEASE_GRACE() view returns(uint64)
func (_Job *JobCaller) RELEASEGRACE(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "RELEASE_GRACE")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// RELEASEGRACE is a free data retrieval call binding the contract method 0xe73ad7c6.
//
// Solidity: function RELEASE_GRACE() view returns(uint64)
func (_Job *JobSession) RELEASEGRACE() (uint64, error) {
	return _Job.Contract.RELEASEGRACE(&_Job.CallOpts)
}

// RELEASEGRACE is a free data retrieval call binding the contract method 0xe73ad7c6.
//
// Solidity: function RELEASE_GRACE() view returns(uint64)
func (_Job *JobCallerSession) RELEASEGRACE() (uint64, error) {
	return _Job.Contract.RELEASEGRACE(&_Job.CallOpts)
}

// Agent is a free data retrieval call binding the contract method 0xf5ff5c76.
//
// Solidity: function agent() view returns(address)
func (_Job *JobCaller) Agent(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "agent")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Agent is a free data retrieval call binding the contract method 0xf5ff5c76.
//
// Solidity: function agent() view returns(address)
func (_Job *JobSession) Agent() (common.Address, error) {
	return _Job.Contract.Agent(&_Job.CallOpts)
}

// Agent is a free data retrieval call binding the contract method 0xf5ff5c76.
//
// Solidity: function agent() view returns(address)
func (_Job *JobCallerSession) Agent() (common.Address, error) {
	return _Job.Contract.Agent(&_Job.CallOpts)
}

// AgentId is a free data retrieval call binding the contract method 0xe84f43b7.
//
// Solidity: function agentId() view returns(uint256)
func (_Job *JobCaller) AgentId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "agentId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentId is a free data retrieval call binding the contract method 0xe84f43b7.
//
// Solidity: function agentId() view returns(uint256)
func (_Job *JobSession) AgentId() (*big.Int, error) {
	return _Job.Contract.AgentId(&_Job.CallOpts)
}

// AgentId is a free data retrieval call binding the contract method 0xe84f43b7.
//
// Solidity: function agentId() view returns(uint256)
func (_Job *JobCallerSession) AgentId() (*big.Int, error) {
	return _Job.Contract.AgentId(&_Job.CallOpts)
}

// Arbiter is a free data retrieval call binding the contract method 0xfe25e00a.
//
// Solidity: function arbiter() view returns(address)
func (_Job *JobCaller) Arbiter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "arbiter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Arbiter is a free data retrieval call binding the contract method 0xfe25e00a.
//
// Solidity: function arbiter() view returns(address)
func (_Job *JobSession) Arbiter() (common.Address, error) {
	return _Job.Contract.Arbiter(&_Job.CallOpts)
}

// Arbiter is a free data retrieval call binding the contract method 0xfe25e00a.
//
// Solidity: function arbiter() view returns(address)
func (_Job *JobCallerSession) Arbiter() (common.Address, error) {
	return _Job.Contract.Arbiter(&_Job.CallOpts)
}

// Budget is a free data retrieval call binding the contract method 0xed01bf29.
//
// Solidity: function budget() view returns(uint256)
func (_Job *JobCaller) Budget(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "budget")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Budget is a free data retrieval call binding the contract method 0xed01bf29.
//
// Solidity: function budget() view returns(uint256)
func (_Job *JobSession) Budget() (*big.Int, error) {
	return _Job.Contract.Budget(&_Job.CallOpts)
}

// Budget is a free data retrieval call binding the contract method 0xed01bf29.
//
// Solidity: function budget() view returns(uint256)
func (_Job *JobCallerSession) Budget() (*big.Int, error) {
	return _Job.Contract.Budget(&_Job.CallOpts)
}

// CurrentPhase is a free data retrieval call binding the contract method 0x055ad42e.
//
// Solidity: function currentPhase() view returns(uint8 current)
func (_Job *JobCaller) CurrentPhase(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "currentPhase")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CurrentPhase is a free data retrieval call binding the contract method 0x055ad42e.
//
// Solidity: function currentPhase() view returns(uint8 current)
func (_Job *JobSession) CurrentPhase() (uint8, error) {
	return _Job.Contract.CurrentPhase(&_Job.CallOpts)
}

// CurrentPhase is a free data retrieval call binding the contract method 0x055ad42e.
//
// Solidity: function currentPhase() view returns(uint8 current)
func (_Job *JobCallerSession) CurrentPhase() (uint8, error) {
	return _Job.Contract.CurrentPhase(&_Job.CallOpts)
}

// Deadline is a free data retrieval call binding the contract method 0x29dcb0cf.
//
// Solidity: function deadline() view returns(uint64)
func (_Job *JobCaller) Deadline(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "deadline")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// Deadline is a free data retrieval call binding the contract method 0x29dcb0cf.
//
// Solidity: function deadline() view returns(uint64)
func (_Job *JobSession) Deadline() (uint64, error) {
	return _Job.Contract.Deadline(&_Job.CallOpts)
}

// Deadline is a free data retrieval call binding the contract method 0x29dcb0cf.
//
// Solidity: function deadline() view returns(uint64)
func (_Job *JobCallerSession) Deadline() (uint64, error) {
	return _Job.Contract.Deadline(&_Job.CallOpts)
}

// Evaluator is a free data retrieval call binding the contract method 0x9cb93dd1.
//
// Solidity: function evaluator() view returns(address)
func (_Job *JobCaller) Evaluator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "evaluator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Evaluator is a free data retrieval call binding the contract method 0x9cb93dd1.
//
// Solidity: function evaluator() view returns(address)
func (_Job *JobSession) Evaluator() (common.Address, error) {
	return _Job.Contract.Evaluator(&_Job.CallOpts)
}

// Evaluator is a free data retrieval call binding the contract method 0x9cb93dd1.
//
// Solidity: function evaluator() view returns(address)
func (_Job *JobCallerSession) Evaluator() (common.Address, error) {
	return _Job.Contract.Evaluator(&_Job.CallOpts)
}

// JobId is a free data retrieval call binding the contract method 0xc2939d97.
//
// Solidity: function jobId() view returns(uint256)
func (_Job *JobCaller) JobId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "jobId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// JobId is a free data retrieval call binding the contract method 0xc2939d97.
//
// Solidity: function jobId() view returns(uint256)
func (_Job *JobSession) JobId() (*big.Int, error) {
	return _Job.Contract.JobId(&_Job.CallOpts)
}

// JobId is a free data retrieval call binding the contract method 0xc2939d97.
//
// Solidity: function jobId() view returns(uint256)
func (_Job *JobCallerSession) JobId() (*big.Int, error) {
	return _Job.Contract.JobId(&_Job.CallOpts)
}

// JobType is a free data retrieval call binding the contract method 0x8080be77.
//
// Solidity: function jobType() view returns(uint8)
func (_Job *JobCaller) JobType(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "jobType")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// JobType is a free data retrieval call binding the contract method 0x8080be77.
//
// Solidity: function jobType() view returns(uint8)
func (_Job *JobSession) JobType() (uint8, error) {
	return _Job.Contract.JobType(&_Job.CallOpts)
}

// JobType is a free data retrieval call binding the contract method 0x8080be77.
//
// Solidity: function jobType() view returns(uint8)
func (_Job *JobCallerSession) JobType() (uint8, error) {
	return _Job.Contract.JobType(&_Job.CallOpts)
}

// Phase is a free data retrieval call binding the contract method 0xb1c9fe6e.
//
// Solidity: function phase() view returns(uint8)
func (_Job *JobCaller) Phase(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "phase")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Phase is a free data retrieval call binding the contract method 0xb1c9fe6e.
//
// Solidity: function phase() view returns(uint8)
func (_Job *JobSession) Phase() (uint8, error) {
	return _Job.Contract.Phase(&_Job.CallOpts)
}

// Phase is a free data retrieval call binding the contract method 0xb1c9fe6e.
//
// Solidity: function phase() view returns(uint8)
func (_Job *JobCallerSession) Phase() (uint8, error) {
	return _Job.Contract.Phase(&_Job.CallOpts)
}

// Principal is a free data retrieval call binding the contract method 0xba5d3078.
//
// Solidity: function principal() view returns(address)
func (_Job *JobCaller) Principal(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "principal")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Principal is a free data retrieval call binding the contract method 0xba5d3078.
//
// Solidity: function principal() view returns(address)
func (_Job *JobSession) Principal() (common.Address, error) {
	return _Job.Contract.Principal(&_Job.CallOpts)
}

// Principal is a free data retrieval call binding the contract method 0xba5d3078.
//
// Solidity: function principal() view returns(address)
func (_Job *JobCallerSession) Principal() (common.Address, error) {
	return _Job.Contract.Principal(&_Job.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_Job *JobCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_Job *JobSession) Registry() (common.Address, error) {
	return _Job.Contract.Registry(&_Job.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_Job *JobCallerSession) Registry() (common.Address, error) {
	return _Job.Contract.Registry(&_Job.CallOpts)
}

// ResultSubmittedAt is a free data retrieval call binding the contract method 0xd478f498.
//
// Solidity: function resultSubmittedAt() view returns(uint64)
func (_Job *JobCaller) ResultSubmittedAt(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "resultSubmittedAt")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ResultSubmittedAt is a free data retrieval call binding the contract method 0xd478f498.
//
// Solidity: function resultSubmittedAt() view returns(uint64)
func (_Job *JobSession) ResultSubmittedAt() (uint64, error) {
	return _Job.Contract.ResultSubmittedAt(&_Job.CallOpts)
}

// ResultSubmittedAt is a free data retrieval call binding the contract method 0xd478f498.
//
// Solidity: function resultSubmittedAt() view returns(uint64)
func (_Job *JobCallerSession) ResultSubmittedAt() (uint64, error) {
	return _Job.Contract.ResultSubmittedAt(&_Job.CallOpts)
}

// ResultURI is a free data retrieval call binding the contract method 0xb56fd20a.
//
// Solidity: function resultURI() view returns(string)
func (_Job *JobCaller) ResultURI(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "resultURI")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ResultURI is a free data retrieval call binding the contract method 0xb56fd20a.
//
// Solidity: function resultURI() view returns(string)
func (_Job *JobSession) ResultURI() (string, error) {
	return _Job.Contract.ResultURI(&_Job.CallOpts)
}

// ResultURI is a free data retrieval call binding the contract method 0xb56fd20a.
//
// Solidity: function resultURI() view returns(string)
func (_Job *JobCallerSession) ResultURI() (string, error) {
	return _Job.Contract.ResultURI(&_Job.CallOpts)
}

// TargetAgentId is a free data retrieval call binding the contract method 0xa63176be.
//
// Solidity: function targetAgentId() view returns(uint256)
func (_Job *JobCaller) TargetAgentId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "targetAgentId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TargetAgentId is a free data retrieval call binding the contract method 0xa63176be.
//
// Solidity: function targetAgentId() view returns(uint256)
func (_Job *JobSession) TargetAgentId() (*big.Int, error) {
	return _Job.Contract.TargetAgentId(&_Job.CallOpts)
}

// TargetAgentId is a free data retrieval call binding the contract method 0xa63176be.
//
// Solidity: function targetAgentId() view returns(uint256)
func (_Job *JobCallerSession) TargetAgentId() (*big.Int, error) {
	return _Job.Contract.TargetAgentId(&_Job.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Job *JobCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Job.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Job *JobSession) Token() (common.Address, error) {
	return _Job.Contract.Token(&_Job.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Job *JobCallerSession) Token() (common.Address, error) {
	return _Job.Contract.Token(&_Job.CallOpts)
}

// Accept is a paid mutator transaction binding the contract method 0x19b05f49.
//
// Solidity: function accept(uint256 _agentId) returns()
func (_Job *JobTransactor) Accept(opts *bind.TransactOpts, _agentId *big.Int) (*types.Transaction, error) {
	return _Job.contract.Transact(opts, "accept", _agentId)
}

// Accept is a paid mutator transaction binding the contract method 0x19b05f49.
//
// Solidity: function accept(uint256 _agentId) returns()
func (_Job *JobSession) Accept(_agentId *big.Int) (*types.Transaction, error) {
	return _Job.Contract.Accept(&_Job.TransactOpts, _agentId)
}

// Accept is a paid mutator transaction binding the contract method 0x19b05f49.
//
// Solidity: function accept(uint256 _agentId) returns()
func (_Job *JobTransactorSession) Accept(_agentId *big.Int) (*types.Transaction, error) {
	return _Job.Contract.Accept(&_Job.TransactOpts, _agentId)
}

// ApproveResult is a paid mutator transaction binding the contract method 0xc5ec9b6a.
//
// Solidity: function approveResult() returns()
func (_Job *JobTransactor) ApproveResult(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Job.contract.Transact(opts, "approveResult")
}

// ApproveResult is a paid mutator transaction binding the contract method 0xc5ec9b6a.
//
// Solidity: function approveResult() returns()
func (_Job *JobSession) ApproveResult() (*types.Transaction, error) {
	return _Job.Contract.ApproveResult(&_Job.TransactOpts)
}

// ApproveResult is a paid mutator transaction binding the contract method 0xc5ec9b6a.
//
// Solidity: function approveResult() returns()
func (_Job *JobTransactorSession) ApproveResult() (*types.Transaction, error) {
	return _Job.Contract.ApproveResult(&_Job.TransactOpts)
}

// Cancel is a paid mutator transaction binding the contract method 0x0b4f3f3d.
//
// Solidity: function cancel(string reason) returns()
func (_Job *JobTransactor) Cancel(opts *bind.TransactOpts, reason string) (*types.Transaction, error) {
	return _Job.contract.Transact(opts, "cancel", reason)
}

// Cancel is a paid mutator transaction binding the contract method 0x0b4f3f3d.
//
// Solidity: function cancel(string reason) returns()
func (_Job *JobSession) Cancel(reason string) (*types.Transaction, error) {
	return _Job.Contract.Cancel(&_Job.TransactOpts, reason)
}

// Cancel is a paid mutator transaction binding the contract method 0x0b4f3f3d.
//
// Solidity: function cancel(string reason) returns()
func (_Job *JobTransactorSession) Cancel(reason string) (*types.Transaction, error) {
	return _Job.Contract.Cancel(&_Job.TransactOpts, reason)
}

// ExpireJob is a paid mutator transaction binding the contract method 0x7ae6a045.
//
// Solidity: function expireJob() returns()
func (_Job *JobTransactor) ExpireJob(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Job.contract.Transact(opts, "expireJob")
}

// ExpireJob is a paid mutator transaction binding the contract method 0x7ae6a045.
//
// Solidity: function expireJob() returns()
func (_Job *JobSession) ExpireJob() (*types.Transaction, error) {
	return _Job.Contract.ExpireJob(&_Job.TransactOpts)
}

// ExpireJob is a paid mutator transaction binding the contract method 0x7ae6a045.
//
// Solidity: function expireJob() returns()
func (_Job *JobTransactorSession) ExpireJob() (*types.Transaction, error) {
	return _Job.Contract.ExpireJob(&_Job.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4955ffc.
//
// Solidity: function initialize((uint256,address,address,uint256,address,uint256,uint64,uint8,address,address) p) returns()
func (_Job *JobTransactor) Initialize(opts *bind.TransactOpts, p JobInitParams) (*types.Transaction, error) {
	return _Job.contract.Transact(opts, "initialize", p)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4955ffc.
//
// Solidity: function initialize((uint256,address,address,uint256,address,uint256,uint64,uint8,address,address) p) returns()
func (_Job *JobSession) Initialize(p JobInitParams) (*types.Transaction, error) {
	return _Job.Contract.Initialize(&_Job.TransactOpts, p)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4955ffc.
//
// Solidity: function initialize((uint256,address,address,uint256,address,uint256,uint64,uint8,address,address) p) returns()
func (_Job *JobTransactorSession) Initialize(p JobInitParams) (*types.Transaction, error) {
	return _Job.Contract.Initialize(&_Job.TransactOpts, p)
}

// RaiseDispute is a paid mutator transaction binding the contract method 0x6daa2d44.
//
// Solidity: function raiseDispute() returns()
func (_Job *JobTransactor) RaiseDispute(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Job.contract.Transact(opts, "raiseDispute")
}

// RaiseDispute is a paid mutator transaction binding the contract method 0x6daa2d44.
//
// Solidity: function raiseDispute() returns()
func (_Job *JobSession) RaiseDispute() (*types.Transaction, error) {
	return _Job.Contract.RaiseDispute(&_Job.TransactOpts)
}

// RaiseDispute is a paid mutator transaction binding the contract method 0x6daa2d44.
//
// Solidity: function raiseDispute() returns()
func (_Job *JobTransactorSession) RaiseDispute() (*types.Transaction, error) {
	return _Job.Contract.RaiseDispute(&_Job.TransactOpts)
}

// RejectResult is a paid mutator transaction binding the contract method 0x3868606a.
//
// Solidity: function rejectResult(string reason) returns()
func (_Job *JobTransactor) RejectResult(opts *bind.TransactOpts, reason string) (*types.Transaction, error) {
	return _Job.contract.Transact(opts, "rejectResult", reason)
}

// RejectResult is a paid mutator transaction binding the contract method 0x3868606a.
//
// Solidity: function rejectResult(string reason) returns()
func (_Job *JobSession) RejectResult(reason string) (*types.Transaction, error) {
	return _Job.Contract.RejectResult(&_Job.TransactOpts, reason)
}

// RejectResult is a paid mutator transaction binding the contract method 0x3868606a.
//
// Solidity: function rejectResult(string reason) returns()
func (_Job *JobTransactorSession) RejectResult(reason string) (*types.Transaction, error) {
	return _Job.Contract.RejectResult(&_Job.TransactOpts, reason)
}

// Release is a paid mutator transaction binding the contract method 0x86d1a69f.
//
// Solidity: function release() returns()
func (_Job *JobTransactor) Release(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Job.contract.Transact(opts, "release")
}

// Release is a paid mutator transaction binding the contract method 0x86d1a69f.
//
// Solidity: function release() returns()
func (_Job *JobSession) Release() (*types.Transaction, error) {
	return _Job.Contract.Release(&_Job.TransactOpts)
}

// Release is a paid mutator transaction binding the contract method 0x86d1a69f.
//
// Solidity: function release() returns()
func (_Job *JobTransactorSession) Release() (*types.Transaction, error) {
	return _Job.Contract.Release(&_Job.TransactOpts)
}

// ResolveDispute is a paid mutator transaction binding the contract method 0x89e1e82a.
//
// Solidity: function resolveDispute(bool agentFavoured) returns()
func (_Job *JobTransactor) ResolveDispute(opts *bind.TransactOpts, agentFavoured bool) (*types.Transaction, error) {
	return _Job.contract.Transact(opts, "resolveDispute", agentFavoured)
}

// ResolveDispute is a paid mutator transaction binding the contract method 0x89e1e82a.
//
// Solidity: function resolveDispute(bool agentFavoured) returns()
func (_Job *JobSession) ResolveDispute(agentFavoured bool) (*types.Transaction, error) {
	return _Job.Contract.ResolveDispute(&_Job.TransactOpts, agentFavoured)
}

// ResolveDispute is a paid mutator transaction binding the contract method 0x89e1e82a.
//
// Solidity: function resolveDispute(bool agentFavoured) returns()
func (_Job *JobTransactorSession) ResolveDispute(agentFavoured bool) (*types.Transaction, error) {
	return _Job.Contract.ResolveDispute(&_Job.TransactOpts, agentFavoured)
}

// SubmitResult is a paid mutator transaction binding the contract method 0xec184245.
//
// Solidity: function submitResult(string _resultURI) returns()
func (_Job *JobTransactor) SubmitResult(opts *bind.TransactOpts, _resultURI string) (*types.Transaction, error) {
	return _Job.contract.Transact(opts, "submitResult", _resultURI)
}

// SubmitResult is a paid mutator transaction binding the contract method 0xec184245.
//
// Solidity: function submitResult(string _resultURI) returns()
func (_Job *JobSession) SubmitResult(_resultURI string) (*types.Transaction, error) {
	return _Job.Contract.SubmitResult(&_Job.TransactOpts, _resultURI)
}

// SubmitResult is a paid mutator transaction binding the contract method 0xec184245.
//
// Solidity: function submitResult(string _resultURI) returns()
func (_Job *JobTransactorSession) SubmitResult(_resultURI string) (*types.Transaction, error) {
	return _Job.Contract.SubmitResult(&_Job.TransactOpts, _resultURI)
}

// JobAgentAcceptedIterator is returned from FilterAgentAccepted and is used to iterate over the raw logs and unpacked data for AgentAccepted events raised by the Job contract.
type JobAgentAcceptedIterator struct {
	Event *JobAgentAccepted // Event containing the contract specifics and raw log

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
func (it *JobAgentAcceptedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobAgentAccepted)
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
		it.Event = new(JobAgentAccepted)
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
func (it *JobAgentAcceptedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobAgentAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobAgentAccepted represents a AgentAccepted event raised by the Job contract.
type JobAgentAccepted struct {
	JobId   *big.Int
	Agent   common.Address
	AgentId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentAccepted is a free log retrieval operation binding the contract event 0xa08a9d61f61bc14cb9b5ef024d79b5dbbb271dfebaaa5a22ed231c25d5d4ee8e.
//
// Solidity: event AgentAccepted(uint256 indexed jobId, address indexed agent, uint256 indexed agentId)
func (_Job *JobFilterer) FilterAgentAccepted(opts *bind.FilterOpts, jobId []*big.Int, agent []common.Address, agentId []*big.Int) (*JobAgentAcceptedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _Job.contract.FilterLogs(opts, "AgentAccepted", jobIdRule, agentRule, agentIdRule)
	if err != nil {
		return nil, err
	}
	return &JobAgentAcceptedIterator{contract: _Job.contract, event: "AgentAccepted", logs: logs, sub: sub}, nil
}

// WatchAgentAccepted is a free log subscription operation binding the contract event 0xa08a9d61f61bc14cb9b5ef024d79b5dbbb271dfebaaa5a22ed231c25d5d4ee8e.
//
// Solidity: event AgentAccepted(uint256 indexed jobId, address indexed agent, uint256 indexed agentId)
func (_Job *JobFilterer) WatchAgentAccepted(opts *bind.WatchOpts, sink chan<- *JobAgentAccepted, jobId []*big.Int, agent []common.Address, agentId []*big.Int) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _Job.contract.WatchLogs(opts, "AgentAccepted", jobIdRule, agentRule, agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobAgentAccepted)
				if err := _Job.contract.UnpackLog(event, "AgentAccepted", log); err != nil {
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

// ParseAgentAccepted is a log parse operation binding the contract event 0xa08a9d61f61bc14cb9b5ef024d79b5dbbb271dfebaaa5a22ed231c25d5d4ee8e.
//
// Solidity: event AgentAccepted(uint256 indexed jobId, address indexed agent, uint256 indexed agentId)
func (_Job *JobFilterer) ParseAgentAccepted(log types.Log) (*JobAgentAccepted, error) {
	event := new(JobAgentAccepted)
	if err := _Job.contract.UnpackLog(event, "AgentAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobDisputeRaisedIterator is returned from FilterDisputeRaised and is used to iterate over the raw logs and unpacked data for DisputeRaised events raised by the Job contract.
type JobDisputeRaisedIterator struct {
	Event *JobDisputeRaised // Event containing the contract specifics and raw log

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
func (it *JobDisputeRaisedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobDisputeRaised)
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
		it.Event = new(JobDisputeRaised)
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
func (it *JobDisputeRaisedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobDisputeRaisedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobDisputeRaised represents a DisputeRaised event raised by the Job contract.
type JobDisputeRaised struct {
	JobId    *big.Int
	RaisedBy common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDisputeRaised is a free log retrieval operation binding the contract event 0x84a477df8a28a4276ca6dee4458a06c3015f30c477d9c949ede4e13ff8a552b4.
//
// Solidity: event DisputeRaised(uint256 indexed jobId, address indexed raisedBy)
func (_Job *JobFilterer) FilterDisputeRaised(opts *bind.FilterOpts, jobId []*big.Int, raisedBy []common.Address) (*JobDisputeRaisedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var raisedByRule []interface{}
	for _, raisedByItem := range raisedBy {
		raisedByRule = append(raisedByRule, raisedByItem)
	}

	logs, sub, err := _Job.contract.FilterLogs(opts, "DisputeRaised", jobIdRule, raisedByRule)
	if err != nil {
		return nil, err
	}
	return &JobDisputeRaisedIterator{contract: _Job.contract, event: "DisputeRaised", logs: logs, sub: sub}, nil
}

// WatchDisputeRaised is a free log subscription operation binding the contract event 0x84a477df8a28a4276ca6dee4458a06c3015f30c477d9c949ede4e13ff8a552b4.
//
// Solidity: event DisputeRaised(uint256 indexed jobId, address indexed raisedBy)
func (_Job *JobFilterer) WatchDisputeRaised(opts *bind.WatchOpts, sink chan<- *JobDisputeRaised, jobId []*big.Int, raisedBy []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var raisedByRule []interface{}
	for _, raisedByItem := range raisedBy {
		raisedByRule = append(raisedByRule, raisedByItem)
	}

	logs, sub, err := _Job.contract.WatchLogs(opts, "DisputeRaised", jobIdRule, raisedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobDisputeRaised)
				if err := _Job.contract.UnpackLog(event, "DisputeRaised", log); err != nil {
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

// ParseDisputeRaised is a log parse operation binding the contract event 0x84a477df8a28a4276ca6dee4458a06c3015f30c477d9c949ede4e13ff8a552b4.
//
// Solidity: event DisputeRaised(uint256 indexed jobId, address indexed raisedBy)
func (_Job *JobFilterer) ParseDisputeRaised(log types.Log) (*JobDisputeRaised, error) {
	event := new(JobDisputeRaised)
	if err := _Job.contract.UnpackLog(event, "DisputeRaised", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobDisputeResolvedIterator is returned from FilterDisputeResolved and is used to iterate over the raw logs and unpacked data for DisputeResolved events raised by the Job contract.
type JobDisputeResolvedIterator struct {
	Event *JobDisputeResolved // Event containing the contract specifics and raw log

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
func (it *JobDisputeResolvedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobDisputeResolved)
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
		it.Event = new(JobDisputeResolved)
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
func (it *JobDisputeResolvedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobDisputeResolvedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobDisputeResolved represents a DisputeResolved event raised by the Job contract.
type JobDisputeResolved struct {
	JobId         *big.Int
	Resolver      common.Address
	AgentFavoured bool
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDisputeResolved is a free log retrieval operation binding the contract event 0x8fdd4548a8481406b6e29c0d6f25e27cd72502f79f4adf409468502e7920dabc.
//
// Solidity: event DisputeResolved(uint256 indexed jobId, address indexed resolver, bool agentFavoured)
func (_Job *JobFilterer) FilterDisputeResolved(opts *bind.FilterOpts, jobId []*big.Int, resolver []common.Address) (*JobDisputeResolvedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var resolverRule []interface{}
	for _, resolverItem := range resolver {
		resolverRule = append(resolverRule, resolverItem)
	}

	logs, sub, err := _Job.contract.FilterLogs(opts, "DisputeResolved", jobIdRule, resolverRule)
	if err != nil {
		return nil, err
	}
	return &JobDisputeResolvedIterator{contract: _Job.contract, event: "DisputeResolved", logs: logs, sub: sub}, nil
}

// WatchDisputeResolved is a free log subscription operation binding the contract event 0x8fdd4548a8481406b6e29c0d6f25e27cd72502f79f4adf409468502e7920dabc.
//
// Solidity: event DisputeResolved(uint256 indexed jobId, address indexed resolver, bool agentFavoured)
func (_Job *JobFilterer) WatchDisputeResolved(opts *bind.WatchOpts, sink chan<- *JobDisputeResolved, jobId []*big.Int, resolver []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var resolverRule []interface{}
	for _, resolverItem := range resolver {
		resolverRule = append(resolverRule, resolverItem)
	}

	logs, sub, err := _Job.contract.WatchLogs(opts, "DisputeResolved", jobIdRule, resolverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobDisputeResolved)
				if err := _Job.contract.UnpackLog(event, "DisputeResolved", log); err != nil {
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

// ParseDisputeResolved is a log parse operation binding the contract event 0x8fdd4548a8481406b6e29c0d6f25e27cd72502f79f4adf409468502e7920dabc.
//
// Solidity: event DisputeResolved(uint256 indexed jobId, address indexed resolver, bool agentFavoured)
func (_Job *JobFilterer) ParseDisputeResolved(log types.Log) (*JobDisputeResolved, error) {
	event := new(JobDisputeResolved)
	if err := _Job.contract.UnpackLog(event, "DisputeResolved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobEvaluatorAssignedIterator is returned from FilterEvaluatorAssigned and is used to iterate over the raw logs and unpacked data for EvaluatorAssigned events raised by the Job contract.
type JobEvaluatorAssignedIterator struct {
	Event *JobEvaluatorAssigned // Event containing the contract specifics and raw log

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
func (it *JobEvaluatorAssignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobEvaluatorAssigned)
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
		it.Event = new(JobEvaluatorAssigned)
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
func (it *JobEvaluatorAssignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobEvaluatorAssignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobEvaluatorAssigned represents a EvaluatorAssigned event raised by the Job contract.
type JobEvaluatorAssigned struct {
	JobId     *big.Int
	Evaluator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEvaluatorAssigned is a free log retrieval operation binding the contract event 0x50a93d710505e6f207121334c60e2a4c6312fdbae71f879f5abee6488e20b131.
//
// Solidity: event EvaluatorAssigned(uint256 indexed jobId, address indexed evaluator)
func (_Job *JobFilterer) FilterEvaluatorAssigned(opts *bind.FilterOpts, jobId []*big.Int, evaluator []common.Address) (*JobEvaluatorAssignedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var evaluatorRule []interface{}
	for _, evaluatorItem := range evaluator {
		evaluatorRule = append(evaluatorRule, evaluatorItem)
	}

	logs, sub, err := _Job.contract.FilterLogs(opts, "EvaluatorAssigned", jobIdRule, evaluatorRule)
	if err != nil {
		return nil, err
	}
	return &JobEvaluatorAssignedIterator{contract: _Job.contract, event: "EvaluatorAssigned", logs: logs, sub: sub}, nil
}

// WatchEvaluatorAssigned is a free log subscription operation binding the contract event 0x50a93d710505e6f207121334c60e2a4c6312fdbae71f879f5abee6488e20b131.
//
// Solidity: event EvaluatorAssigned(uint256 indexed jobId, address indexed evaluator)
func (_Job *JobFilterer) WatchEvaluatorAssigned(opts *bind.WatchOpts, sink chan<- *JobEvaluatorAssigned, jobId []*big.Int, evaluator []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var evaluatorRule []interface{}
	for _, evaluatorItem := range evaluator {
		evaluatorRule = append(evaluatorRule, evaluatorItem)
	}

	logs, sub, err := _Job.contract.WatchLogs(opts, "EvaluatorAssigned", jobIdRule, evaluatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobEvaluatorAssigned)
				if err := _Job.contract.UnpackLog(event, "EvaluatorAssigned", log); err != nil {
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

// ParseEvaluatorAssigned is a log parse operation binding the contract event 0x50a93d710505e6f207121334c60e2a4c6312fdbae71f879f5abee6488e20b131.
//
// Solidity: event EvaluatorAssigned(uint256 indexed jobId, address indexed evaluator)
func (_Job *JobFilterer) ParseEvaluatorAssigned(log types.Log) (*JobEvaluatorAssigned, error) {
	event := new(JobEvaluatorAssigned)
	if err := _Job.contract.UnpackLog(event, "EvaluatorAssigned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Job contract.
type JobInitializedIterator struct {
	Event *JobInitialized // Event containing the contract specifics and raw log

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
func (it *JobInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobInitialized)
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
		it.Event = new(JobInitialized)
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
func (it *JobInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobInitialized represents a Initialized event raised by the Job contract.
type JobInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Job *JobFilterer) FilterInitialized(opts *bind.FilterOpts) (*JobInitializedIterator, error) {

	logs, sub, err := _Job.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &JobInitializedIterator{contract: _Job.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Job *JobFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *JobInitialized) (event.Subscription, error) {

	logs, sub, err := _Job.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobInitialized)
				if err := _Job.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Job *JobFilterer) ParseInitialized(log types.Log) (*JobInitialized, error) {
	event := new(JobInitialized)
	if err := _Job.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobJobCancelledIterator is returned from FilterJobCancelled and is used to iterate over the raw logs and unpacked data for JobCancelled events raised by the Job contract.
type JobJobCancelledIterator struct {
	Event *JobJobCancelled // Event containing the contract specifics and raw log

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
func (it *JobJobCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobJobCancelled)
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
		it.Event = new(JobJobCancelled)
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
func (it *JobJobCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobJobCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobJobCancelled represents a JobCancelled event raised by the Job contract.
type JobJobCancelled struct {
	JobId       *big.Int
	CancelledBy common.Address
	Reason      string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterJobCancelled is a free log retrieval operation binding the contract event 0x141454c3e63a845c54bcfebc9a4f9843c448cd8c0bb7077e5f70235e4f441356.
//
// Solidity: event JobCancelled(uint256 indexed jobId, address indexed cancelledBy, string reason)
func (_Job *JobFilterer) FilterJobCancelled(opts *bind.FilterOpts, jobId []*big.Int, cancelledBy []common.Address) (*JobJobCancelledIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var cancelledByRule []interface{}
	for _, cancelledByItem := range cancelledBy {
		cancelledByRule = append(cancelledByRule, cancelledByItem)
	}

	logs, sub, err := _Job.contract.FilterLogs(opts, "JobCancelled", jobIdRule, cancelledByRule)
	if err != nil {
		return nil, err
	}
	return &JobJobCancelledIterator{contract: _Job.contract, event: "JobCancelled", logs: logs, sub: sub}, nil
}

// WatchJobCancelled is a free log subscription operation binding the contract event 0x141454c3e63a845c54bcfebc9a4f9843c448cd8c0bb7077e5f70235e4f441356.
//
// Solidity: event JobCancelled(uint256 indexed jobId, address indexed cancelledBy, string reason)
func (_Job *JobFilterer) WatchJobCancelled(opts *bind.WatchOpts, sink chan<- *JobJobCancelled, jobId []*big.Int, cancelledBy []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var cancelledByRule []interface{}
	for _, cancelledByItem := range cancelledBy {
		cancelledByRule = append(cancelledByRule, cancelledByItem)
	}

	logs, sub, err := _Job.contract.WatchLogs(opts, "JobCancelled", jobIdRule, cancelledByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobJobCancelled)
				if err := _Job.contract.UnpackLog(event, "JobCancelled", log); err != nil {
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

// ParseJobCancelled is a log parse operation binding the contract event 0x141454c3e63a845c54bcfebc9a4f9843c448cd8c0bb7077e5f70235e4f441356.
//
// Solidity: event JobCancelled(uint256 indexed jobId, address indexed cancelledBy, string reason)
func (_Job *JobFilterer) ParseJobCancelled(log types.Log) (*JobJobCancelled, error) {
	event := new(JobJobCancelled)
	if err := _Job.contract.UnpackLog(event, "JobCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobJobCompletedIterator is returned from FilterJobCompleted and is used to iterate over the raw logs and unpacked data for JobCompleted events raised by the Job contract.
type JobJobCompletedIterator struct {
	Event *JobJobCompleted // Event containing the contract specifics and raw log

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
func (it *JobJobCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobJobCompleted)
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
		it.Event = new(JobJobCompleted)
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
func (it *JobJobCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobJobCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobJobCompleted represents a JobCompleted event raised by the Job contract.
type JobJobCompleted struct {
	JobId      *big.Int
	ReleasedBy common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterJobCompleted is a free log retrieval operation binding the contract event 0x4dbba472101e9f148a3a5ecbe793f0ee16a7efe4bf3f8cbfab5330e1642ef955.
//
// Solidity: event JobCompleted(uint256 indexed jobId, address indexed releasedBy, uint256 amount)
func (_Job *JobFilterer) FilterJobCompleted(opts *bind.FilterOpts, jobId []*big.Int, releasedBy []common.Address) (*JobJobCompletedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var releasedByRule []interface{}
	for _, releasedByItem := range releasedBy {
		releasedByRule = append(releasedByRule, releasedByItem)
	}

	logs, sub, err := _Job.contract.FilterLogs(opts, "JobCompleted", jobIdRule, releasedByRule)
	if err != nil {
		return nil, err
	}
	return &JobJobCompletedIterator{contract: _Job.contract, event: "JobCompleted", logs: logs, sub: sub}, nil
}

// WatchJobCompleted is a free log subscription operation binding the contract event 0x4dbba472101e9f148a3a5ecbe793f0ee16a7efe4bf3f8cbfab5330e1642ef955.
//
// Solidity: event JobCompleted(uint256 indexed jobId, address indexed releasedBy, uint256 amount)
func (_Job *JobFilterer) WatchJobCompleted(opts *bind.WatchOpts, sink chan<- *JobJobCompleted, jobId []*big.Int, releasedBy []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var releasedByRule []interface{}
	for _, releasedByItem := range releasedBy {
		releasedByRule = append(releasedByRule, releasedByItem)
	}

	logs, sub, err := _Job.contract.WatchLogs(opts, "JobCompleted", jobIdRule, releasedByRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobJobCompleted)
				if err := _Job.contract.UnpackLog(event, "JobCompleted", log); err != nil {
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

// ParseJobCompleted is a log parse operation binding the contract event 0x4dbba472101e9f148a3a5ecbe793f0ee16a7efe4bf3f8cbfab5330e1642ef955.
//
// Solidity: event JobCompleted(uint256 indexed jobId, address indexed releasedBy, uint256 amount)
func (_Job *JobFilterer) ParseJobCompleted(log types.Log) (*JobJobCompleted, error) {
	event := new(JobJobCompleted)
	if err := _Job.contract.UnpackLog(event, "JobCompleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobJobInitialisedIterator is returned from FilterJobInitialised and is used to iterate over the raw logs and unpacked data for JobInitialised events raised by the Job contract.
type JobJobInitialisedIterator struct {
	Event *JobJobInitialised // Event containing the contract specifics and raw log

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
func (it *JobJobInitialisedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobJobInitialised)
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
		it.Event = new(JobJobInitialised)
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
func (it *JobJobInitialisedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobJobInitialisedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobJobInitialised represents a JobInitialised event raised by the Job contract.
type JobJobInitialised struct {
	JobId     *big.Int
	Principal common.Address
	AgentId   *big.Int
	JobType   uint8
	Budget    *big.Int
	Token     common.Address
	Deadline  uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterJobInitialised is a free log retrieval operation binding the contract event 0x6396f9d24200dbcd88fb1965753f38f76c647a2a0a793a00875957937869ce48.
//
// Solidity: event JobInitialised(uint256 indexed jobId, address indexed principal, uint256 indexed agentId, uint8 jobType, uint256 budget, address token, uint64 deadline)
func (_Job *JobFilterer) FilterJobInitialised(opts *bind.FilterOpts, jobId []*big.Int, principal []common.Address, agentId []*big.Int) (*JobJobInitialisedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var principalRule []interface{}
	for _, principalItem := range principal {
		principalRule = append(principalRule, principalItem)
	}
	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _Job.contract.FilterLogs(opts, "JobInitialised", jobIdRule, principalRule, agentIdRule)
	if err != nil {
		return nil, err
	}
	return &JobJobInitialisedIterator{contract: _Job.contract, event: "JobInitialised", logs: logs, sub: sub}, nil
}

// WatchJobInitialised is a free log subscription operation binding the contract event 0x6396f9d24200dbcd88fb1965753f38f76c647a2a0a793a00875957937869ce48.
//
// Solidity: event JobInitialised(uint256 indexed jobId, address indexed principal, uint256 indexed agentId, uint8 jobType, uint256 budget, address token, uint64 deadline)
func (_Job *JobFilterer) WatchJobInitialised(opts *bind.WatchOpts, sink chan<- *JobJobInitialised, jobId []*big.Int, principal []common.Address, agentId []*big.Int) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var principalRule []interface{}
	for _, principalItem := range principal {
		principalRule = append(principalRule, principalItem)
	}
	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _Job.contract.WatchLogs(opts, "JobInitialised", jobIdRule, principalRule, agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobJobInitialised)
				if err := _Job.contract.UnpackLog(event, "JobInitialised", log); err != nil {
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

// ParseJobInitialised is a log parse operation binding the contract event 0x6396f9d24200dbcd88fb1965753f38f76c647a2a0a793a00875957937869ce48.
//
// Solidity: event JobInitialised(uint256 indexed jobId, address indexed principal, uint256 indexed agentId, uint8 jobType, uint256 budget, address token, uint64 deadline)
func (_Job *JobFilterer) ParseJobInitialised(log types.Log) (*JobJobInitialised, error) {
	event := new(JobJobInitialised)
	if err := _Job.contract.UnpackLog(event, "JobInitialised", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobResultApprovedIterator is returned from FilterResultApproved and is used to iterate over the raw logs and unpacked data for ResultApproved events raised by the Job contract.
type JobResultApprovedIterator struct {
	Event *JobResultApproved // Event containing the contract specifics and raw log

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
func (it *JobResultApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobResultApproved)
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
		it.Event = new(JobResultApproved)
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
func (it *JobResultApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobResultApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobResultApproved represents a ResultApproved event raised by the Job contract.
type JobResultApproved struct {
	JobId     *big.Int
	Evaluator common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterResultApproved is a free log retrieval operation binding the contract event 0xdfa2a378c6a00b39f4d2b4247ba44e90a894960a78603b2bf5af5de40dbe183d.
//
// Solidity: event ResultApproved(uint256 indexed jobId, address indexed evaluator)
func (_Job *JobFilterer) FilterResultApproved(opts *bind.FilterOpts, jobId []*big.Int, evaluator []common.Address) (*JobResultApprovedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var evaluatorRule []interface{}
	for _, evaluatorItem := range evaluator {
		evaluatorRule = append(evaluatorRule, evaluatorItem)
	}

	logs, sub, err := _Job.contract.FilterLogs(opts, "ResultApproved", jobIdRule, evaluatorRule)
	if err != nil {
		return nil, err
	}
	return &JobResultApprovedIterator{contract: _Job.contract, event: "ResultApproved", logs: logs, sub: sub}, nil
}

// WatchResultApproved is a free log subscription operation binding the contract event 0xdfa2a378c6a00b39f4d2b4247ba44e90a894960a78603b2bf5af5de40dbe183d.
//
// Solidity: event ResultApproved(uint256 indexed jobId, address indexed evaluator)
func (_Job *JobFilterer) WatchResultApproved(opts *bind.WatchOpts, sink chan<- *JobResultApproved, jobId []*big.Int, evaluator []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var evaluatorRule []interface{}
	for _, evaluatorItem := range evaluator {
		evaluatorRule = append(evaluatorRule, evaluatorItem)
	}

	logs, sub, err := _Job.contract.WatchLogs(opts, "ResultApproved", jobIdRule, evaluatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobResultApproved)
				if err := _Job.contract.UnpackLog(event, "ResultApproved", log); err != nil {
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

// ParseResultApproved is a log parse operation binding the contract event 0xdfa2a378c6a00b39f4d2b4247ba44e90a894960a78603b2bf5af5de40dbe183d.
//
// Solidity: event ResultApproved(uint256 indexed jobId, address indexed evaluator)
func (_Job *JobFilterer) ParseResultApproved(log types.Log) (*JobResultApproved, error) {
	event := new(JobResultApproved)
	if err := _Job.contract.UnpackLog(event, "ResultApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobResultRejectedIterator is returned from FilterResultRejected and is used to iterate over the raw logs and unpacked data for ResultRejected events raised by the Job contract.
type JobResultRejectedIterator struct {
	Event *JobResultRejected // Event containing the contract specifics and raw log

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
func (it *JobResultRejectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobResultRejected)
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
		it.Event = new(JobResultRejected)
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
func (it *JobResultRejectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobResultRejectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobResultRejected represents a ResultRejected event raised by the Job contract.
type JobResultRejected struct {
	JobId     *big.Int
	Evaluator common.Address
	Reason    string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterResultRejected is a free log retrieval operation binding the contract event 0x8e23b50499d65aea2724488d5fbc5924cafc75e6613cea60f4ec898b84123d90.
//
// Solidity: event ResultRejected(uint256 indexed jobId, address indexed evaluator, string reason)
func (_Job *JobFilterer) FilterResultRejected(opts *bind.FilterOpts, jobId []*big.Int, evaluator []common.Address) (*JobResultRejectedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var evaluatorRule []interface{}
	for _, evaluatorItem := range evaluator {
		evaluatorRule = append(evaluatorRule, evaluatorItem)
	}

	logs, sub, err := _Job.contract.FilterLogs(opts, "ResultRejected", jobIdRule, evaluatorRule)
	if err != nil {
		return nil, err
	}
	return &JobResultRejectedIterator{contract: _Job.contract, event: "ResultRejected", logs: logs, sub: sub}, nil
}

// WatchResultRejected is a free log subscription operation binding the contract event 0x8e23b50499d65aea2724488d5fbc5924cafc75e6613cea60f4ec898b84123d90.
//
// Solidity: event ResultRejected(uint256 indexed jobId, address indexed evaluator, string reason)
func (_Job *JobFilterer) WatchResultRejected(opts *bind.WatchOpts, sink chan<- *JobResultRejected, jobId []*big.Int, evaluator []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var evaluatorRule []interface{}
	for _, evaluatorItem := range evaluator {
		evaluatorRule = append(evaluatorRule, evaluatorItem)
	}

	logs, sub, err := _Job.contract.WatchLogs(opts, "ResultRejected", jobIdRule, evaluatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobResultRejected)
				if err := _Job.contract.UnpackLog(event, "ResultRejected", log); err != nil {
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

// ParseResultRejected is a log parse operation binding the contract event 0x8e23b50499d65aea2724488d5fbc5924cafc75e6613cea60f4ec898b84123d90.
//
// Solidity: event ResultRejected(uint256 indexed jobId, address indexed evaluator, string reason)
func (_Job *JobFilterer) ParseResultRejected(log types.Log) (*JobResultRejected, error) {
	event := new(JobResultRejected)
	if err := _Job.contract.UnpackLog(event, "ResultRejected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// JobResultSubmittedIterator is returned from FilterResultSubmitted and is used to iterate over the raw logs and unpacked data for ResultSubmitted events raised by the Job contract.
type JobResultSubmittedIterator struct {
	Event *JobResultSubmitted // Event containing the contract specifics and raw log

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
func (it *JobResultSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(JobResultSubmitted)
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
		it.Event = new(JobResultSubmitted)
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
func (it *JobResultSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *JobResultSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// JobResultSubmitted represents a ResultSubmitted event raised by the Job contract.
type JobResultSubmitted struct {
	JobId     *big.Int
	Agent     common.Address
	ResultURI string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterResultSubmitted is a free log retrieval operation binding the contract event 0xc06b551d984e333ac851ab20b1454a08d92740468f52ff54c0cd5270817f20a9.
//
// Solidity: event ResultSubmitted(uint256 indexed jobId, address indexed agent, string resultURI)
func (_Job *JobFilterer) FilterResultSubmitted(opts *bind.FilterOpts, jobId []*big.Int, agent []common.Address) (*JobResultSubmittedIterator, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _Job.contract.FilterLogs(opts, "ResultSubmitted", jobIdRule, agentRule)
	if err != nil {
		return nil, err
	}
	return &JobResultSubmittedIterator{contract: _Job.contract, event: "ResultSubmitted", logs: logs, sub: sub}, nil
}

// WatchResultSubmitted is a free log subscription operation binding the contract event 0xc06b551d984e333ac851ab20b1454a08d92740468f52ff54c0cd5270817f20a9.
//
// Solidity: event ResultSubmitted(uint256 indexed jobId, address indexed agent, string resultURI)
func (_Job *JobFilterer) WatchResultSubmitted(opts *bind.WatchOpts, sink chan<- *JobResultSubmitted, jobId []*big.Int, agent []common.Address) (event.Subscription, error) {

	var jobIdRule []interface{}
	for _, jobIdItem := range jobId {
		jobIdRule = append(jobIdRule, jobIdItem)
	}
	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _Job.contract.WatchLogs(opts, "ResultSubmitted", jobIdRule, agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(JobResultSubmitted)
				if err := _Job.contract.UnpackLog(event, "ResultSubmitted", log); err != nil {
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

// ParseResultSubmitted is a log parse operation binding the contract event 0xc06b551d984e333ac851ab20b1454a08d92740468f52ff54c0cd5270817f20a9.
//
// Solidity: event ResultSubmitted(uint256 indexed jobId, address indexed agent, string resultURI)
func (_Job *JobFilterer) ParseResultSubmitted(log types.Log) (*JobResultSubmitted, error) {
	event := new(JobResultSubmitted)
	if err := _Job.contract.UnpackLog(event, "ResultSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
