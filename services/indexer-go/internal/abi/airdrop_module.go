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

// AirdropModuleMetaData contains all meta data concerning the AirdropModule contract.
var AirdropModuleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_ALLOCATION_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proof\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimed\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configs\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"allocation\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"snapshotBlock\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"deadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"agentAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"configured\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configure\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"snapshotBlock\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"allocation\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"deadline\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sweep\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Claimed\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Configured\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"merkleRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"snapshotBlock\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"allocation\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"deadline\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"agentAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Swept\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AllocationOverCap\",\"inputs\":[{\"name\":\"allocation\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"cap\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"AlreadyClaimed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DeadlineInPast\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DeadlineNotPassed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DeadlinePassed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"have\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"need\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAgentAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NothingToSweep\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAllocation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroRoot\",\"inputs\":[]}]",
}

// AirdropModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use AirdropModuleMetaData.ABI instead.
var AirdropModuleABI = AirdropModuleMetaData.ABI

// AirdropModule is an auto generated Go binding around an Ethereum contract.
type AirdropModule struct {
	AirdropModuleCaller     // Read-only binding to the contract
	AirdropModuleTransactor // Write-only binding to the contract
	AirdropModuleFilterer   // Log filterer for contract events
}

// AirdropModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type AirdropModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AirdropModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AirdropModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AirdropModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AirdropModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AirdropModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AirdropModuleSession struct {
	Contract     *AirdropModule    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AirdropModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AirdropModuleCallerSession struct {
	Contract *AirdropModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// AirdropModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AirdropModuleTransactorSession struct {
	Contract     *AirdropModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// AirdropModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type AirdropModuleRaw struct {
	Contract *AirdropModule // Generic contract binding to access the raw methods on
}

// AirdropModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AirdropModuleCallerRaw struct {
	Contract *AirdropModuleCaller // Generic read-only contract binding to access the raw methods on
}

// AirdropModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AirdropModuleTransactorRaw struct {
	Contract *AirdropModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAirdropModule creates a new instance of AirdropModule, bound to a specific deployed contract.
func NewAirdropModule(address common.Address, backend bind.ContractBackend) (*AirdropModule, error) {
	contract, err := bindAirdropModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AirdropModule{AirdropModuleCaller: AirdropModuleCaller{contract: contract}, AirdropModuleTransactor: AirdropModuleTransactor{contract: contract}, AirdropModuleFilterer: AirdropModuleFilterer{contract: contract}}, nil
}

// NewAirdropModuleCaller creates a new read-only instance of AirdropModule, bound to a specific deployed contract.
func NewAirdropModuleCaller(address common.Address, caller bind.ContractCaller) (*AirdropModuleCaller, error) {
	contract, err := bindAirdropModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AirdropModuleCaller{contract: contract}, nil
}

// NewAirdropModuleTransactor creates a new write-only instance of AirdropModule, bound to a specific deployed contract.
func NewAirdropModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*AirdropModuleTransactor, error) {
	contract, err := bindAirdropModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AirdropModuleTransactor{contract: contract}, nil
}

// NewAirdropModuleFilterer creates a new log filterer instance of AirdropModule, bound to a specific deployed contract.
func NewAirdropModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*AirdropModuleFilterer, error) {
	contract, err := bindAirdropModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AirdropModuleFilterer{contract: contract}, nil
}

// bindAirdropModule binds a generic wrapper to an already deployed contract.
func bindAirdropModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AirdropModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AirdropModule *AirdropModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AirdropModule.Contract.AirdropModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AirdropModule *AirdropModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AirdropModule.Contract.AirdropModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AirdropModule *AirdropModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AirdropModule.Contract.AirdropModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AirdropModule *AirdropModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AirdropModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AirdropModule *AirdropModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AirdropModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AirdropModule *AirdropModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AirdropModule.Contract.contract.Transact(opts, method, params...)
}

// MAXALLOCATIONBPS is a free data retrieval call binding the contract method 0xb308a8e8.
//
// Solidity: function MAX_ALLOCATION_BPS() view returns(uint16)
func (_AirdropModule *AirdropModuleCaller) MAXALLOCATIONBPS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _AirdropModule.contract.Call(opts, &out, "MAX_ALLOCATION_BPS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MAXALLOCATIONBPS is a free data retrieval call binding the contract method 0xb308a8e8.
//
// Solidity: function MAX_ALLOCATION_BPS() view returns(uint16)
func (_AirdropModule *AirdropModuleSession) MAXALLOCATIONBPS() (uint16, error) {
	return _AirdropModule.Contract.MAXALLOCATIONBPS(&_AirdropModule.CallOpts)
}

// MAXALLOCATIONBPS is a free data retrieval call binding the contract method 0xb308a8e8.
//
// Solidity: function MAX_ALLOCATION_BPS() view returns(uint16)
func (_AirdropModule *AirdropModuleCallerSession) MAXALLOCATIONBPS() (uint16, error) {
	return _AirdropModule.Contract.MAXALLOCATIONBPS(&_AirdropModule.CallOpts)
}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint16)
func (_AirdropModule *AirdropModuleCaller) MAXBPS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _AirdropModule.contract.Call(opts, &out, "MAX_BPS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint16)
func (_AirdropModule *AirdropModuleSession) MAXBPS() (uint16, error) {
	return _AirdropModule.Contract.MAXBPS(&_AirdropModule.CallOpts)
}

// MAXBPS is a free data retrieval call binding the contract method 0xfd967f47.
//
// Solidity: function MAX_BPS() view returns(uint16)
func (_AirdropModule *AirdropModuleCallerSession) MAXBPS() (uint16, error) {
	return _AirdropModule.Contract.MAXBPS(&_AirdropModule.CallOpts)
}

// Claimed is a free data retrieval call binding the contract method 0x0c9cbf0e.
//
// Solidity: function claimed(address agent, address recipient) view returns(bool)
func (_AirdropModule *AirdropModuleCaller) Claimed(opts *bind.CallOpts, agent common.Address, recipient common.Address) (bool, error) {
	var out []interface{}
	err := _AirdropModule.contract.Call(opts, &out, "claimed", agent, recipient)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Claimed is a free data retrieval call binding the contract method 0x0c9cbf0e.
//
// Solidity: function claimed(address agent, address recipient) view returns(bool)
func (_AirdropModule *AirdropModuleSession) Claimed(agent common.Address, recipient common.Address) (bool, error) {
	return _AirdropModule.Contract.Claimed(&_AirdropModule.CallOpts, agent, recipient)
}

// Claimed is a free data retrieval call binding the contract method 0x0c9cbf0e.
//
// Solidity: function claimed(address agent, address recipient) view returns(bool)
func (_AirdropModule *AirdropModuleCallerSession) Claimed(agent common.Address, recipient common.Address) (bool, error) {
	return _AirdropModule.Contract.Claimed(&_AirdropModule.CallOpts, agent, recipient)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(bytes32 merkleRoot, uint128 allocation, uint64 snapshotBlock, uint64 deadline, address agentAdmin, bool configured)
func (_AirdropModule *AirdropModuleCaller) Configs(opts *bind.CallOpts, agent common.Address) (struct {
	MerkleRoot    [32]byte
	Allocation    *big.Int
	SnapshotBlock uint64
	Deadline      uint64
	AgentAdmin    common.Address
	Configured    bool
}, error) {
	var out []interface{}
	err := _AirdropModule.contract.Call(opts, &out, "configs", agent)

	outstruct := new(struct {
		MerkleRoot    [32]byte
		Allocation    *big.Int
		SnapshotBlock uint64
		Deadline      uint64
		AgentAdmin    common.Address
		Configured    bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MerkleRoot = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Allocation = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.SnapshotBlock = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Deadline = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.AgentAdmin = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Configured = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(bytes32 merkleRoot, uint128 allocation, uint64 snapshotBlock, uint64 deadline, address agentAdmin, bool configured)
func (_AirdropModule *AirdropModuleSession) Configs(agent common.Address) (struct {
	MerkleRoot    [32]byte
	Allocation    *big.Int
	SnapshotBlock uint64
	Deadline      uint64
	AgentAdmin    common.Address
	Configured    bool
}, error) {
	return _AirdropModule.Contract.Configs(&_AirdropModule.CallOpts, agent)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(bytes32 merkleRoot, uint128 allocation, uint64 snapshotBlock, uint64 deadline, address agentAdmin, bool configured)
func (_AirdropModule *AirdropModuleCallerSession) Configs(agent common.Address) (struct {
	MerkleRoot    [32]byte
	Allocation    *big.Int
	SnapshotBlock uint64
	Deadline      uint64
	AgentAdmin    common.Address
	Configured    bool
}, error) {
	return _AirdropModule.Contract.Configs(&_AirdropModule.CallOpts, agent)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AirdropModule *AirdropModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AirdropModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AirdropModule *AirdropModuleSession) Owner() (common.Address, error) {
	return _AirdropModule.Contract.Owner(&_AirdropModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AirdropModule *AirdropModuleCallerSession) Owner() (common.Address, error) {
	return _AirdropModule.Contract.Owner(&_AirdropModule.CallOpts)
}

// Claim is a paid mutator transaction binding the contract method 0xfabed412.
//
// Solidity: function claim(address agent, address recipient, uint256 amount, bytes32[] proof) returns()
func (_AirdropModule *AirdropModuleTransactor) Claim(opts *bind.TransactOpts, agent common.Address, recipient common.Address, amount *big.Int, proof [][32]byte) (*types.Transaction, error) {
	return _AirdropModule.contract.Transact(opts, "claim", agent, recipient, amount, proof)
}

// Claim is a paid mutator transaction binding the contract method 0xfabed412.
//
// Solidity: function claim(address agent, address recipient, uint256 amount, bytes32[] proof) returns()
func (_AirdropModule *AirdropModuleSession) Claim(agent common.Address, recipient common.Address, amount *big.Int, proof [][32]byte) (*types.Transaction, error) {
	return _AirdropModule.Contract.Claim(&_AirdropModule.TransactOpts, agent, recipient, amount, proof)
}

// Claim is a paid mutator transaction binding the contract method 0xfabed412.
//
// Solidity: function claim(address agent, address recipient, uint256 amount, bytes32[] proof) returns()
func (_AirdropModule *AirdropModuleTransactorSession) Claim(agent common.Address, recipient common.Address, amount *big.Int, proof [][32]byte) (*types.Transaction, error) {
	return _AirdropModule.Contract.Claim(&_AirdropModule.TransactOpts, agent, recipient, amount, proof)
}

// Configure is a paid mutator transaction binding the contract method 0x7d8f5311.
//
// Solidity: function configure(address agent, bytes32 root, uint64 snapshotBlock, uint128 allocation, uint64 deadline, address admin) returns()
func (_AirdropModule *AirdropModuleTransactor) Configure(opts *bind.TransactOpts, agent common.Address, root [32]byte, snapshotBlock uint64, allocation *big.Int, deadline uint64, admin common.Address) (*types.Transaction, error) {
	return _AirdropModule.contract.Transact(opts, "configure", agent, root, snapshotBlock, allocation, deadline, admin)
}

// Configure is a paid mutator transaction binding the contract method 0x7d8f5311.
//
// Solidity: function configure(address agent, bytes32 root, uint64 snapshotBlock, uint128 allocation, uint64 deadline, address admin) returns()
func (_AirdropModule *AirdropModuleSession) Configure(agent common.Address, root [32]byte, snapshotBlock uint64, allocation *big.Int, deadline uint64, admin common.Address) (*types.Transaction, error) {
	return _AirdropModule.Contract.Configure(&_AirdropModule.TransactOpts, agent, root, snapshotBlock, allocation, deadline, admin)
}

// Configure is a paid mutator transaction binding the contract method 0x7d8f5311.
//
// Solidity: function configure(address agent, bytes32 root, uint64 snapshotBlock, uint128 allocation, uint64 deadline, address admin) returns()
func (_AirdropModule *AirdropModuleTransactorSession) Configure(agent common.Address, root [32]byte, snapshotBlock uint64, allocation *big.Int, deadline uint64, admin common.Address) (*types.Transaction, error) {
	return _AirdropModule.Contract.Configure(&_AirdropModule.TransactOpts, agent, root, snapshotBlock, allocation, deadline, admin)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AirdropModule *AirdropModuleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AirdropModule.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AirdropModule *AirdropModuleSession) RenounceOwnership() (*types.Transaction, error) {
	return _AirdropModule.Contract.RenounceOwnership(&_AirdropModule.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AirdropModule *AirdropModuleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AirdropModule.Contract.RenounceOwnership(&_AirdropModule.TransactOpts)
}

// Sweep is a paid mutator transaction binding the contract method 0x01681a62.
//
// Solidity: function sweep(address agent) returns()
func (_AirdropModule *AirdropModuleTransactor) Sweep(opts *bind.TransactOpts, agent common.Address) (*types.Transaction, error) {
	return _AirdropModule.contract.Transact(opts, "sweep", agent)
}

// Sweep is a paid mutator transaction binding the contract method 0x01681a62.
//
// Solidity: function sweep(address agent) returns()
func (_AirdropModule *AirdropModuleSession) Sweep(agent common.Address) (*types.Transaction, error) {
	return _AirdropModule.Contract.Sweep(&_AirdropModule.TransactOpts, agent)
}

// Sweep is a paid mutator transaction binding the contract method 0x01681a62.
//
// Solidity: function sweep(address agent) returns()
func (_AirdropModule *AirdropModuleTransactorSession) Sweep(agent common.Address) (*types.Transaction, error) {
	return _AirdropModule.Contract.Sweep(&_AirdropModule.TransactOpts, agent)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AirdropModule *AirdropModuleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AirdropModule.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AirdropModule *AirdropModuleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AirdropModule.Contract.TransferOwnership(&_AirdropModule.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AirdropModule *AirdropModuleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AirdropModule.Contract.TransferOwnership(&_AirdropModule.TransactOpts, newOwner)
}

// AirdropModuleClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the AirdropModule contract.
type AirdropModuleClaimedIterator struct {
	Event *AirdropModuleClaimed // Event containing the contract specifics and raw log

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
func (it *AirdropModuleClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AirdropModuleClaimed)
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
		it.Event = new(AirdropModuleClaimed)
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
func (it *AirdropModuleClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AirdropModuleClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AirdropModuleClaimed represents a Claimed event raised by the AirdropModule contract.
type AirdropModuleClaimed struct {
	Agent     common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0xf7a40077ff7a04c7e61f6f26fb13774259ddf1b6bce9ecf26a8276cdd3992683.
//
// Solidity: event Claimed(address indexed agent, address indexed recipient, uint256 amount)
func (_AirdropModule *AirdropModuleFilterer) FilterClaimed(opts *bind.FilterOpts, agent []common.Address, recipient []common.Address) (*AirdropModuleClaimedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AirdropModule.contract.FilterLogs(opts, "Claimed", agentRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &AirdropModuleClaimedIterator{contract: _AirdropModule.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0xf7a40077ff7a04c7e61f6f26fb13774259ddf1b6bce9ecf26a8276cdd3992683.
//
// Solidity: event Claimed(address indexed agent, address indexed recipient, uint256 amount)
func (_AirdropModule *AirdropModuleFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *AirdropModuleClaimed, agent []common.Address, recipient []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _AirdropModule.contract.WatchLogs(opts, "Claimed", agentRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AirdropModuleClaimed)
				if err := _AirdropModule.contract.UnpackLog(event, "Claimed", log); err != nil {
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

// ParseClaimed is a log parse operation binding the contract event 0xf7a40077ff7a04c7e61f6f26fb13774259ddf1b6bce9ecf26a8276cdd3992683.
//
// Solidity: event Claimed(address indexed agent, address indexed recipient, uint256 amount)
func (_AirdropModule *AirdropModuleFilterer) ParseClaimed(log types.Log) (*AirdropModuleClaimed, error) {
	event := new(AirdropModuleClaimed)
	if err := _AirdropModule.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AirdropModuleConfiguredIterator is returned from FilterConfigured and is used to iterate over the raw logs and unpacked data for Configured events raised by the AirdropModule contract.
type AirdropModuleConfiguredIterator struct {
	Event *AirdropModuleConfigured // Event containing the contract specifics and raw log

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
func (it *AirdropModuleConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AirdropModuleConfigured)
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
		it.Event = new(AirdropModuleConfigured)
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
func (it *AirdropModuleConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AirdropModuleConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AirdropModuleConfigured represents a Configured event raised by the AirdropModule contract.
type AirdropModuleConfigured struct {
	Agent         common.Address
	MerkleRoot    [32]byte
	SnapshotBlock uint64
	Allocation    *big.Int
	Deadline      uint64
	AgentAdmin    common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterConfigured is a free log retrieval operation binding the contract event 0x00316ebf58d741c0423f6ba06de7853745eade101944eae2b70b0d2c8e7b9863.
//
// Solidity: event Configured(address indexed agent, bytes32 merkleRoot, uint64 snapshotBlock, uint128 allocation, uint64 deadline, address agentAdmin)
func (_AirdropModule *AirdropModuleFilterer) FilterConfigured(opts *bind.FilterOpts, agent []common.Address) (*AirdropModuleConfiguredIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _AirdropModule.contract.FilterLogs(opts, "Configured", agentRule)
	if err != nil {
		return nil, err
	}
	return &AirdropModuleConfiguredIterator{contract: _AirdropModule.contract, event: "Configured", logs: logs, sub: sub}, nil
}

// WatchConfigured is a free log subscription operation binding the contract event 0x00316ebf58d741c0423f6ba06de7853745eade101944eae2b70b0d2c8e7b9863.
//
// Solidity: event Configured(address indexed agent, bytes32 merkleRoot, uint64 snapshotBlock, uint128 allocation, uint64 deadline, address agentAdmin)
func (_AirdropModule *AirdropModuleFilterer) WatchConfigured(opts *bind.WatchOpts, sink chan<- *AirdropModuleConfigured, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _AirdropModule.contract.WatchLogs(opts, "Configured", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AirdropModuleConfigured)
				if err := _AirdropModule.contract.UnpackLog(event, "Configured", log); err != nil {
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

// ParseConfigured is a log parse operation binding the contract event 0x00316ebf58d741c0423f6ba06de7853745eade101944eae2b70b0d2c8e7b9863.
//
// Solidity: event Configured(address indexed agent, bytes32 merkleRoot, uint64 snapshotBlock, uint128 allocation, uint64 deadline, address agentAdmin)
func (_AirdropModule *AirdropModuleFilterer) ParseConfigured(log types.Log) (*AirdropModuleConfigured, error) {
	event := new(AirdropModuleConfigured)
	if err := _AirdropModule.contract.UnpackLog(event, "Configured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AirdropModuleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AirdropModule contract.
type AirdropModuleOwnershipTransferredIterator struct {
	Event *AirdropModuleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AirdropModuleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AirdropModuleOwnershipTransferred)
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
		it.Event = new(AirdropModuleOwnershipTransferred)
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
func (it *AirdropModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AirdropModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AirdropModuleOwnershipTransferred represents a OwnershipTransferred event raised by the AirdropModule contract.
type AirdropModuleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AirdropModule *AirdropModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AirdropModuleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AirdropModule.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AirdropModuleOwnershipTransferredIterator{contract: _AirdropModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AirdropModule *AirdropModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AirdropModuleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AirdropModule.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AirdropModuleOwnershipTransferred)
				if err := _AirdropModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_AirdropModule *AirdropModuleFilterer) ParseOwnershipTransferred(log types.Log) (*AirdropModuleOwnershipTransferred, error) {
	event := new(AirdropModuleOwnershipTransferred)
	if err := _AirdropModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AirdropModuleSweptIterator is returned from FilterSwept and is used to iterate over the raw logs and unpacked data for Swept events raised by the AirdropModule contract.
type AirdropModuleSweptIterator struct {
	Event *AirdropModuleSwept // Event containing the contract specifics and raw log

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
func (it *AirdropModuleSweptIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AirdropModuleSwept)
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
		it.Event = new(AirdropModuleSwept)
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
func (it *AirdropModuleSweptIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AirdropModuleSweptIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AirdropModuleSwept represents a Swept event raised by the AirdropModule contract.
type AirdropModuleSwept struct {
	Agent  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterSwept is a free log retrieval operation binding the contract event 0xc36b5179cb9c303b200074996eab2b3473eac370fdd7eba3bec636fe35109696.
//
// Solidity: event Swept(address indexed agent, uint256 amount)
func (_AirdropModule *AirdropModuleFilterer) FilterSwept(opts *bind.FilterOpts, agent []common.Address) (*AirdropModuleSweptIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _AirdropModule.contract.FilterLogs(opts, "Swept", agentRule)
	if err != nil {
		return nil, err
	}
	return &AirdropModuleSweptIterator{contract: _AirdropModule.contract, event: "Swept", logs: logs, sub: sub}, nil
}

// WatchSwept is a free log subscription operation binding the contract event 0xc36b5179cb9c303b200074996eab2b3473eac370fdd7eba3bec636fe35109696.
//
// Solidity: event Swept(address indexed agent, uint256 amount)
func (_AirdropModule *AirdropModuleFilterer) WatchSwept(opts *bind.WatchOpts, sink chan<- *AirdropModuleSwept, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _AirdropModule.contract.WatchLogs(opts, "Swept", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AirdropModuleSwept)
				if err := _AirdropModule.contract.UnpackLog(event, "Swept", log); err != nil {
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

// ParseSwept is a log parse operation binding the contract event 0xc36b5179cb9c303b200074996eab2b3473eac370fdd7eba3bec636fe35109696.
//
// Solidity: event Swept(address indexed agent, uint256 amount)
func (_AirdropModule *AirdropModuleFilterer) ParseSwept(log types.Log) (*AirdropModuleSwept, error) {
	event := new(AirdropModuleSwept)
	if err := _AirdropModule.contract.UnpackLog(event, "Swept", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
