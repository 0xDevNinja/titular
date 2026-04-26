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

// CapitalFormationModuleMetaData contains all meta data concerning the CapitalFormationModule contract.
var CapitalFormationModuleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"MAX_TIERS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_TWAP_BLOCKS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tier\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"configs\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"agentToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payoutToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pair\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentIsToken0\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentAdmin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"configured\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"claimedBitmap\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configure\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payoutToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pair_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"creator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentAdmin_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"fdvThresholds\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"payoutsUsdc\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isClaimed\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tier\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"outstanding\",\"inputs\":[{\"name\":\"payoutToken\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"snapshot\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"snapshots\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"priceCumulative\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blockTimestampLast\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"exists\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tierAt\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tier\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"threshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"payout\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tierCount\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Configured\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agentToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"payoutToken\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"pair\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"agentAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"tierCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalDeposit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MilestoneClaimed\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tier\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"fdvUsd\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"payoutUsd\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Snapshotted\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"priceCumulative\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"blockTimestampLast\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"blockNumber\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FdvBelowThreshold\",\"inputs\":[{\"name\":\"fdv\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"threshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InsufficientFunding\",\"inputs\":[{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"actual\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidTierCount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PairTokenMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PairTooYoung\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PairUninitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SnapshotMissing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ThresholdsNotIncreasing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TierAlreadyClaimed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TierOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TwapWindowTooShort\",\"inputs\":[{\"name\":\"elapsedBlocks\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"required\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroPayout\",\"inputs\":[]}]",
}

// CapitalFormationModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use CapitalFormationModuleMetaData.ABI instead.
var CapitalFormationModuleABI = CapitalFormationModuleMetaData.ABI

// CapitalFormationModule is an auto generated Go binding around an Ethereum contract.
type CapitalFormationModule struct {
	CapitalFormationModuleCaller     // Read-only binding to the contract
	CapitalFormationModuleTransactor // Write-only binding to the contract
	CapitalFormationModuleFilterer   // Log filterer for contract events
}

// CapitalFormationModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type CapitalFormationModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CapitalFormationModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CapitalFormationModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CapitalFormationModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CapitalFormationModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CapitalFormationModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CapitalFormationModuleSession struct {
	Contract     *CapitalFormationModule // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// CapitalFormationModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CapitalFormationModuleCallerSession struct {
	Contract *CapitalFormationModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// CapitalFormationModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CapitalFormationModuleTransactorSession struct {
	Contract     *CapitalFormationModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// CapitalFormationModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type CapitalFormationModuleRaw struct {
	Contract *CapitalFormationModule // Generic contract binding to access the raw methods on
}

// CapitalFormationModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CapitalFormationModuleCallerRaw struct {
	Contract *CapitalFormationModuleCaller // Generic read-only contract binding to access the raw methods on
}

// CapitalFormationModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CapitalFormationModuleTransactorRaw struct {
	Contract *CapitalFormationModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCapitalFormationModule creates a new instance of CapitalFormationModule, bound to a specific deployed contract.
func NewCapitalFormationModule(address common.Address, backend bind.ContractBackend) (*CapitalFormationModule, error) {
	contract, err := bindCapitalFormationModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CapitalFormationModule{CapitalFormationModuleCaller: CapitalFormationModuleCaller{contract: contract}, CapitalFormationModuleTransactor: CapitalFormationModuleTransactor{contract: contract}, CapitalFormationModuleFilterer: CapitalFormationModuleFilterer{contract: contract}}, nil
}

// NewCapitalFormationModuleCaller creates a new read-only instance of CapitalFormationModule, bound to a specific deployed contract.
func NewCapitalFormationModuleCaller(address common.Address, caller bind.ContractCaller) (*CapitalFormationModuleCaller, error) {
	contract, err := bindCapitalFormationModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CapitalFormationModuleCaller{contract: contract}, nil
}

// NewCapitalFormationModuleTransactor creates a new write-only instance of CapitalFormationModule, bound to a specific deployed contract.
func NewCapitalFormationModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*CapitalFormationModuleTransactor, error) {
	contract, err := bindCapitalFormationModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CapitalFormationModuleTransactor{contract: contract}, nil
}

// NewCapitalFormationModuleFilterer creates a new log filterer instance of CapitalFormationModule, bound to a specific deployed contract.
func NewCapitalFormationModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*CapitalFormationModuleFilterer, error) {
	contract, err := bindCapitalFormationModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CapitalFormationModuleFilterer{contract: contract}, nil
}

// bindCapitalFormationModule binds a generic wrapper to an already deployed contract.
func bindCapitalFormationModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CapitalFormationModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CapitalFormationModule *CapitalFormationModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CapitalFormationModule.Contract.CapitalFormationModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CapitalFormationModule *CapitalFormationModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.CapitalFormationModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CapitalFormationModule *CapitalFormationModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.CapitalFormationModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CapitalFormationModule *CapitalFormationModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CapitalFormationModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CapitalFormationModule *CapitalFormationModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CapitalFormationModule *CapitalFormationModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.contract.Transact(opts, method, params...)
}

// MAXTIERS is a free data retrieval call binding the contract method 0x07c3d4af.
//
// Solidity: function MAX_TIERS() view returns(uint8)
func (_CapitalFormationModule *CapitalFormationModuleCaller) MAXTIERS(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _CapitalFormationModule.contract.Call(opts, &out, "MAX_TIERS")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MAXTIERS is a free data retrieval call binding the contract method 0x07c3d4af.
//
// Solidity: function MAX_TIERS() view returns(uint8)
func (_CapitalFormationModule *CapitalFormationModuleSession) MAXTIERS() (uint8, error) {
	return _CapitalFormationModule.Contract.MAXTIERS(&_CapitalFormationModule.CallOpts)
}

// MAXTIERS is a free data retrieval call binding the contract method 0x07c3d4af.
//
// Solidity: function MAX_TIERS() view returns(uint8)
func (_CapitalFormationModule *CapitalFormationModuleCallerSession) MAXTIERS() (uint8, error) {
	return _CapitalFormationModule.Contract.MAXTIERS(&_CapitalFormationModule.CallOpts)
}

// MINTWAPBLOCKS is a free data retrieval call binding the contract method 0xd6eb6da8.
//
// Solidity: function MIN_TWAP_BLOCKS() view returns(uint64)
func (_CapitalFormationModule *CapitalFormationModuleCaller) MINTWAPBLOCKS(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _CapitalFormationModule.contract.Call(opts, &out, "MIN_TWAP_BLOCKS")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MINTWAPBLOCKS is a free data retrieval call binding the contract method 0xd6eb6da8.
//
// Solidity: function MIN_TWAP_BLOCKS() view returns(uint64)
func (_CapitalFormationModule *CapitalFormationModuleSession) MINTWAPBLOCKS() (uint64, error) {
	return _CapitalFormationModule.Contract.MINTWAPBLOCKS(&_CapitalFormationModule.CallOpts)
}

// MINTWAPBLOCKS is a free data retrieval call binding the contract method 0xd6eb6da8.
//
// Solidity: function MIN_TWAP_BLOCKS() view returns(uint64)
func (_CapitalFormationModule *CapitalFormationModuleCallerSession) MINTWAPBLOCKS() (uint64, error) {
	return _CapitalFormationModule.Contract.MINTWAPBLOCKS(&_CapitalFormationModule.CallOpts)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(address agentToken, address payoutToken, address pair, bool agentIsToken0, address creator, address agentAdmin, bool configured, uint8 claimedBitmap)
func (_CapitalFormationModule *CapitalFormationModuleCaller) Configs(opts *bind.CallOpts, agent common.Address) (struct {
	AgentToken    common.Address
	PayoutToken   common.Address
	Pair          common.Address
	AgentIsToken0 bool
	Creator       common.Address
	AgentAdmin    common.Address
	Configured    bool
	ClaimedBitmap uint8
}, error) {
	var out []interface{}
	err := _CapitalFormationModule.contract.Call(opts, &out, "configs", agent)

	outstruct := new(struct {
		AgentToken    common.Address
		PayoutToken   common.Address
		Pair          common.Address
		AgentIsToken0 bool
		Creator       common.Address
		AgentAdmin    common.Address
		Configured    bool
		ClaimedBitmap uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AgentToken = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.PayoutToken = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Pair = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.AgentIsToken0 = *abi.ConvertType(out[3], new(bool)).(*bool)
	outstruct.Creator = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.AgentAdmin = *abi.ConvertType(out[5], new(common.Address)).(*common.Address)
	outstruct.Configured = *abi.ConvertType(out[6], new(bool)).(*bool)
	outstruct.ClaimedBitmap = *abi.ConvertType(out[7], new(uint8)).(*uint8)

	return *outstruct, err

}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(address agentToken, address payoutToken, address pair, bool agentIsToken0, address creator, address agentAdmin, bool configured, uint8 claimedBitmap)
func (_CapitalFormationModule *CapitalFormationModuleSession) Configs(agent common.Address) (struct {
	AgentToken    common.Address
	PayoutToken   common.Address
	Pair          common.Address
	AgentIsToken0 bool
	Creator       common.Address
	AgentAdmin    common.Address
	Configured    bool
	ClaimedBitmap uint8
}, error) {
	return _CapitalFormationModule.Contract.Configs(&_CapitalFormationModule.CallOpts, agent)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(address agentToken, address payoutToken, address pair, bool agentIsToken0, address creator, address agentAdmin, bool configured, uint8 claimedBitmap)
func (_CapitalFormationModule *CapitalFormationModuleCallerSession) Configs(agent common.Address) (struct {
	AgentToken    common.Address
	PayoutToken   common.Address
	Pair          common.Address
	AgentIsToken0 bool
	Creator       common.Address
	AgentAdmin    common.Address
	Configured    bool
	ClaimedBitmap uint8
}, error) {
	return _CapitalFormationModule.Contract.Configs(&_CapitalFormationModule.CallOpts, agent)
}

// IsClaimed is a free data retrieval call binding the contract method 0xec36aed0.
//
// Solidity: function isClaimed(address agent, uint8 tier) view returns(bool)
func (_CapitalFormationModule *CapitalFormationModuleCaller) IsClaimed(opts *bind.CallOpts, agent common.Address, tier uint8) (bool, error) {
	var out []interface{}
	err := _CapitalFormationModule.contract.Call(opts, &out, "isClaimed", agent, tier)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsClaimed is a free data retrieval call binding the contract method 0xec36aed0.
//
// Solidity: function isClaimed(address agent, uint8 tier) view returns(bool)
func (_CapitalFormationModule *CapitalFormationModuleSession) IsClaimed(agent common.Address, tier uint8) (bool, error) {
	return _CapitalFormationModule.Contract.IsClaimed(&_CapitalFormationModule.CallOpts, agent, tier)
}

// IsClaimed is a free data retrieval call binding the contract method 0xec36aed0.
//
// Solidity: function isClaimed(address agent, uint8 tier) view returns(bool)
func (_CapitalFormationModule *CapitalFormationModuleCallerSession) IsClaimed(agent common.Address, tier uint8) (bool, error) {
	return _CapitalFormationModule.Contract.IsClaimed(&_CapitalFormationModule.CallOpts, agent, tier)
}

// Outstanding is a free data retrieval call binding the contract method 0xde58bafa.
//
// Solidity: function outstanding(address payoutToken) view returns(uint256)
func (_CapitalFormationModule *CapitalFormationModuleCaller) Outstanding(opts *bind.CallOpts, payoutToken common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CapitalFormationModule.contract.Call(opts, &out, "outstanding", payoutToken)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Outstanding is a free data retrieval call binding the contract method 0xde58bafa.
//
// Solidity: function outstanding(address payoutToken) view returns(uint256)
func (_CapitalFormationModule *CapitalFormationModuleSession) Outstanding(payoutToken common.Address) (*big.Int, error) {
	return _CapitalFormationModule.Contract.Outstanding(&_CapitalFormationModule.CallOpts, payoutToken)
}

// Outstanding is a free data retrieval call binding the contract method 0xde58bafa.
//
// Solidity: function outstanding(address payoutToken) view returns(uint256)
func (_CapitalFormationModule *CapitalFormationModuleCallerSession) Outstanding(payoutToken common.Address) (*big.Int, error) {
	return _CapitalFormationModule.Contract.Outstanding(&_CapitalFormationModule.CallOpts, payoutToken)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CapitalFormationModule *CapitalFormationModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CapitalFormationModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CapitalFormationModule *CapitalFormationModuleSession) Owner() (common.Address, error) {
	return _CapitalFormationModule.Contract.Owner(&_CapitalFormationModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CapitalFormationModule *CapitalFormationModuleCallerSession) Owner() (common.Address, error) {
	return _CapitalFormationModule.Contract.Owner(&_CapitalFormationModule.CallOpts)
}

// Snapshots is a free data retrieval call binding the contract method 0x34b3081f.
//
// Solidity: function snapshots(address agent) view returns(uint256 priceCumulative, uint32 blockTimestampLast, uint64 blockNumber, bool exists)
func (_CapitalFormationModule *CapitalFormationModuleCaller) Snapshots(opts *bind.CallOpts, agent common.Address) (struct {
	PriceCumulative    *big.Int
	BlockTimestampLast uint32
	BlockNumber        uint64
	Exists             bool
}, error) {
	var out []interface{}
	err := _CapitalFormationModule.contract.Call(opts, &out, "snapshots", agent)

	outstruct := new(struct {
		PriceCumulative    *big.Int
		BlockTimestampLast uint32
		BlockNumber        uint64
		Exists             bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PriceCumulative = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.BlockTimestampLast = *abi.ConvertType(out[1], new(uint32)).(*uint32)
	outstruct.BlockNumber = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Exists = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Snapshots is a free data retrieval call binding the contract method 0x34b3081f.
//
// Solidity: function snapshots(address agent) view returns(uint256 priceCumulative, uint32 blockTimestampLast, uint64 blockNumber, bool exists)
func (_CapitalFormationModule *CapitalFormationModuleSession) Snapshots(agent common.Address) (struct {
	PriceCumulative    *big.Int
	BlockTimestampLast uint32
	BlockNumber        uint64
	Exists             bool
}, error) {
	return _CapitalFormationModule.Contract.Snapshots(&_CapitalFormationModule.CallOpts, agent)
}

// Snapshots is a free data retrieval call binding the contract method 0x34b3081f.
//
// Solidity: function snapshots(address agent) view returns(uint256 priceCumulative, uint32 blockTimestampLast, uint64 blockNumber, bool exists)
func (_CapitalFormationModule *CapitalFormationModuleCallerSession) Snapshots(agent common.Address) (struct {
	PriceCumulative    *big.Int
	BlockTimestampLast uint32
	BlockNumber        uint64
	Exists             bool
}, error) {
	return _CapitalFormationModule.Contract.Snapshots(&_CapitalFormationModule.CallOpts, agent)
}

// TierAt is a free data retrieval call binding the contract method 0x00eb4bc3.
//
// Solidity: function tierAt(address agent, uint8 tier) view returns(uint256 threshold, uint256 payout)
func (_CapitalFormationModule *CapitalFormationModuleCaller) TierAt(opts *bind.CallOpts, agent common.Address, tier uint8) (struct {
	Threshold *big.Int
	Payout    *big.Int
}, error) {
	var out []interface{}
	err := _CapitalFormationModule.contract.Call(opts, &out, "tierAt", agent, tier)

	outstruct := new(struct {
		Threshold *big.Int
		Payout    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Threshold = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Payout = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// TierAt is a free data retrieval call binding the contract method 0x00eb4bc3.
//
// Solidity: function tierAt(address agent, uint8 tier) view returns(uint256 threshold, uint256 payout)
func (_CapitalFormationModule *CapitalFormationModuleSession) TierAt(agent common.Address, tier uint8) (struct {
	Threshold *big.Int
	Payout    *big.Int
}, error) {
	return _CapitalFormationModule.Contract.TierAt(&_CapitalFormationModule.CallOpts, agent, tier)
}

// TierAt is a free data retrieval call binding the contract method 0x00eb4bc3.
//
// Solidity: function tierAt(address agent, uint8 tier) view returns(uint256 threshold, uint256 payout)
func (_CapitalFormationModule *CapitalFormationModuleCallerSession) TierAt(agent common.Address, tier uint8) (struct {
	Threshold *big.Int
	Payout    *big.Int
}, error) {
	return _CapitalFormationModule.Contract.TierAt(&_CapitalFormationModule.CallOpts, agent, tier)
}

// TierCount is a free data retrieval call binding the contract method 0x2ce13b3f.
//
// Solidity: function tierCount(address agent) view returns(uint256)
func (_CapitalFormationModule *CapitalFormationModuleCaller) TierCount(opts *bind.CallOpts, agent common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CapitalFormationModule.contract.Call(opts, &out, "tierCount", agent)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TierCount is a free data retrieval call binding the contract method 0x2ce13b3f.
//
// Solidity: function tierCount(address agent) view returns(uint256)
func (_CapitalFormationModule *CapitalFormationModuleSession) TierCount(agent common.Address) (*big.Int, error) {
	return _CapitalFormationModule.Contract.TierCount(&_CapitalFormationModule.CallOpts, agent)
}

// TierCount is a free data retrieval call binding the contract method 0x2ce13b3f.
//
// Solidity: function tierCount(address agent) view returns(uint256)
func (_CapitalFormationModule *CapitalFormationModuleCallerSession) TierCount(agent common.Address) (*big.Int, error) {
	return _CapitalFormationModule.Contract.TierCount(&_CapitalFormationModule.CallOpts, agent)
}

// Claim is a paid mutator transaction binding the contract method 0x00e93d31.
//
// Solidity: function claim(address agent, uint8 tier) returns()
func (_CapitalFormationModule *CapitalFormationModuleTransactor) Claim(opts *bind.TransactOpts, agent common.Address, tier uint8) (*types.Transaction, error) {
	return _CapitalFormationModule.contract.Transact(opts, "claim", agent, tier)
}

// Claim is a paid mutator transaction binding the contract method 0x00e93d31.
//
// Solidity: function claim(address agent, uint8 tier) returns()
func (_CapitalFormationModule *CapitalFormationModuleSession) Claim(agent common.Address, tier uint8) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.Claim(&_CapitalFormationModule.TransactOpts, agent, tier)
}

// Claim is a paid mutator transaction binding the contract method 0x00e93d31.
//
// Solidity: function claim(address agent, uint8 tier) returns()
func (_CapitalFormationModule *CapitalFormationModuleTransactorSession) Claim(agent common.Address, tier uint8) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.Claim(&_CapitalFormationModule.TransactOpts, agent, tier)
}

// Configure is a paid mutator transaction binding the contract method 0x3575b693.
//
// Solidity: function configure(address agent, address agentToken_, address payoutToken_, address pair_, address creator_, address agentAdmin_, uint256[] fdvThresholds, uint256[] payoutsUsdc) returns()
func (_CapitalFormationModule *CapitalFormationModuleTransactor) Configure(opts *bind.TransactOpts, agent common.Address, agentToken_ common.Address, payoutToken_ common.Address, pair_ common.Address, creator_ common.Address, agentAdmin_ common.Address, fdvThresholds []*big.Int, payoutsUsdc []*big.Int) (*types.Transaction, error) {
	return _CapitalFormationModule.contract.Transact(opts, "configure", agent, agentToken_, payoutToken_, pair_, creator_, agentAdmin_, fdvThresholds, payoutsUsdc)
}

// Configure is a paid mutator transaction binding the contract method 0x3575b693.
//
// Solidity: function configure(address agent, address agentToken_, address payoutToken_, address pair_, address creator_, address agentAdmin_, uint256[] fdvThresholds, uint256[] payoutsUsdc) returns()
func (_CapitalFormationModule *CapitalFormationModuleSession) Configure(agent common.Address, agentToken_ common.Address, payoutToken_ common.Address, pair_ common.Address, creator_ common.Address, agentAdmin_ common.Address, fdvThresholds []*big.Int, payoutsUsdc []*big.Int) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.Configure(&_CapitalFormationModule.TransactOpts, agent, agentToken_, payoutToken_, pair_, creator_, agentAdmin_, fdvThresholds, payoutsUsdc)
}

// Configure is a paid mutator transaction binding the contract method 0x3575b693.
//
// Solidity: function configure(address agent, address agentToken_, address payoutToken_, address pair_, address creator_, address agentAdmin_, uint256[] fdvThresholds, uint256[] payoutsUsdc) returns()
func (_CapitalFormationModule *CapitalFormationModuleTransactorSession) Configure(agent common.Address, agentToken_ common.Address, payoutToken_ common.Address, pair_ common.Address, creator_ common.Address, agentAdmin_ common.Address, fdvThresholds []*big.Int, payoutsUsdc []*big.Int) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.Configure(&_CapitalFormationModule.TransactOpts, agent, agentToken_, payoutToken_, pair_, creator_, agentAdmin_, fdvThresholds, payoutsUsdc)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CapitalFormationModule *CapitalFormationModuleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CapitalFormationModule.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CapitalFormationModule *CapitalFormationModuleSession) RenounceOwnership() (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.RenounceOwnership(&_CapitalFormationModule.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CapitalFormationModule *CapitalFormationModuleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.RenounceOwnership(&_CapitalFormationModule.TransactOpts)
}

// Snapshot is a paid mutator transaction binding the contract method 0x26512160.
//
// Solidity: function snapshot(address agent) returns()
func (_CapitalFormationModule *CapitalFormationModuleTransactor) Snapshot(opts *bind.TransactOpts, agent common.Address) (*types.Transaction, error) {
	return _CapitalFormationModule.contract.Transact(opts, "snapshot", agent)
}

// Snapshot is a paid mutator transaction binding the contract method 0x26512160.
//
// Solidity: function snapshot(address agent) returns()
func (_CapitalFormationModule *CapitalFormationModuleSession) Snapshot(agent common.Address) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.Snapshot(&_CapitalFormationModule.TransactOpts, agent)
}

// Snapshot is a paid mutator transaction binding the contract method 0x26512160.
//
// Solidity: function snapshot(address agent) returns()
func (_CapitalFormationModule *CapitalFormationModuleTransactorSession) Snapshot(agent common.Address) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.Snapshot(&_CapitalFormationModule.TransactOpts, agent)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CapitalFormationModule *CapitalFormationModuleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CapitalFormationModule.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CapitalFormationModule *CapitalFormationModuleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.TransferOwnership(&_CapitalFormationModule.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CapitalFormationModule *CapitalFormationModuleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CapitalFormationModule.Contract.TransferOwnership(&_CapitalFormationModule.TransactOpts, newOwner)
}

// CapitalFormationModuleConfiguredIterator is returned from FilterConfigured and is used to iterate over the raw logs and unpacked data for Configured events raised by the CapitalFormationModule contract.
type CapitalFormationModuleConfiguredIterator struct {
	Event *CapitalFormationModuleConfigured // Event containing the contract specifics and raw log

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
func (it *CapitalFormationModuleConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapitalFormationModuleConfigured)
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
		it.Event = new(CapitalFormationModuleConfigured)
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
func (it *CapitalFormationModuleConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CapitalFormationModuleConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CapitalFormationModuleConfigured represents a Configured event raised by the CapitalFormationModule contract.
type CapitalFormationModuleConfigured struct {
	Agent        common.Address
	AgentToken   common.Address
	PayoutToken  common.Address
	Pair         common.Address
	Creator      common.Address
	AgentAdmin   common.Address
	TierCount    *big.Int
	TotalDeposit *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterConfigured is a free log retrieval operation binding the contract event 0x987bfd904fb3c2fb18b4503a91cb3b90dd7fda6af766f422ab46988b183dc8a0.
//
// Solidity: event Configured(address indexed agent, address agentToken, address payoutToken, address pair, address creator, address agentAdmin, uint256 tierCount, uint256 totalDeposit)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) FilterConfigured(opts *bind.FilterOpts, agent []common.Address) (*CapitalFormationModuleConfiguredIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _CapitalFormationModule.contract.FilterLogs(opts, "Configured", agentRule)
	if err != nil {
		return nil, err
	}
	return &CapitalFormationModuleConfiguredIterator{contract: _CapitalFormationModule.contract, event: "Configured", logs: logs, sub: sub}, nil
}

// WatchConfigured is a free log subscription operation binding the contract event 0x987bfd904fb3c2fb18b4503a91cb3b90dd7fda6af766f422ab46988b183dc8a0.
//
// Solidity: event Configured(address indexed agent, address agentToken, address payoutToken, address pair, address creator, address agentAdmin, uint256 tierCount, uint256 totalDeposit)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) WatchConfigured(opts *bind.WatchOpts, sink chan<- *CapitalFormationModuleConfigured, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _CapitalFormationModule.contract.WatchLogs(opts, "Configured", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CapitalFormationModuleConfigured)
				if err := _CapitalFormationModule.contract.UnpackLog(event, "Configured", log); err != nil {
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

// ParseConfigured is a log parse operation binding the contract event 0x987bfd904fb3c2fb18b4503a91cb3b90dd7fda6af766f422ab46988b183dc8a0.
//
// Solidity: event Configured(address indexed agent, address agentToken, address payoutToken, address pair, address creator, address agentAdmin, uint256 tierCount, uint256 totalDeposit)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) ParseConfigured(log types.Log) (*CapitalFormationModuleConfigured, error) {
	event := new(CapitalFormationModuleConfigured)
	if err := _CapitalFormationModule.contract.UnpackLog(event, "Configured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CapitalFormationModuleMilestoneClaimedIterator is returned from FilterMilestoneClaimed and is used to iterate over the raw logs and unpacked data for MilestoneClaimed events raised by the CapitalFormationModule contract.
type CapitalFormationModuleMilestoneClaimedIterator struct {
	Event *CapitalFormationModuleMilestoneClaimed // Event containing the contract specifics and raw log

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
func (it *CapitalFormationModuleMilestoneClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapitalFormationModuleMilestoneClaimed)
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
		it.Event = new(CapitalFormationModuleMilestoneClaimed)
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
func (it *CapitalFormationModuleMilestoneClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CapitalFormationModuleMilestoneClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CapitalFormationModuleMilestoneClaimed represents a MilestoneClaimed event raised by the CapitalFormationModule contract.
type CapitalFormationModuleMilestoneClaimed struct {
	Agent     common.Address
	Tier      uint8
	FdvUsd    *big.Int
	PayoutUsd *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterMilestoneClaimed is a free log retrieval operation binding the contract event 0xe9fb2821e39309960329845dfae95b6d6fd4dab51c2ec700a3645a7ad8331da2.
//
// Solidity: event MilestoneClaimed(address indexed agent, uint8 indexed tier, uint256 fdvUsd, uint256 payoutUsd)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) FilterMilestoneClaimed(opts *bind.FilterOpts, agent []common.Address, tier []uint8) (*CapitalFormationModuleMilestoneClaimedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var tierRule []interface{}
	for _, tierItem := range tier {
		tierRule = append(tierRule, tierItem)
	}

	logs, sub, err := _CapitalFormationModule.contract.FilterLogs(opts, "MilestoneClaimed", agentRule, tierRule)
	if err != nil {
		return nil, err
	}
	return &CapitalFormationModuleMilestoneClaimedIterator{contract: _CapitalFormationModule.contract, event: "MilestoneClaimed", logs: logs, sub: sub}, nil
}

// WatchMilestoneClaimed is a free log subscription operation binding the contract event 0xe9fb2821e39309960329845dfae95b6d6fd4dab51c2ec700a3645a7ad8331da2.
//
// Solidity: event MilestoneClaimed(address indexed agent, uint8 indexed tier, uint256 fdvUsd, uint256 payoutUsd)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) WatchMilestoneClaimed(opts *bind.WatchOpts, sink chan<- *CapitalFormationModuleMilestoneClaimed, agent []common.Address, tier []uint8) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var tierRule []interface{}
	for _, tierItem := range tier {
		tierRule = append(tierRule, tierItem)
	}

	logs, sub, err := _CapitalFormationModule.contract.WatchLogs(opts, "MilestoneClaimed", agentRule, tierRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CapitalFormationModuleMilestoneClaimed)
				if err := _CapitalFormationModule.contract.UnpackLog(event, "MilestoneClaimed", log); err != nil {
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

// ParseMilestoneClaimed is a log parse operation binding the contract event 0xe9fb2821e39309960329845dfae95b6d6fd4dab51c2ec700a3645a7ad8331da2.
//
// Solidity: event MilestoneClaimed(address indexed agent, uint8 indexed tier, uint256 fdvUsd, uint256 payoutUsd)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) ParseMilestoneClaimed(log types.Log) (*CapitalFormationModuleMilestoneClaimed, error) {
	event := new(CapitalFormationModuleMilestoneClaimed)
	if err := _CapitalFormationModule.contract.UnpackLog(event, "MilestoneClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CapitalFormationModuleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CapitalFormationModule contract.
type CapitalFormationModuleOwnershipTransferredIterator struct {
	Event *CapitalFormationModuleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CapitalFormationModuleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapitalFormationModuleOwnershipTransferred)
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
		it.Event = new(CapitalFormationModuleOwnershipTransferred)
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
func (it *CapitalFormationModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CapitalFormationModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CapitalFormationModuleOwnershipTransferred represents a OwnershipTransferred event raised by the CapitalFormationModule contract.
type CapitalFormationModuleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CapitalFormationModuleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CapitalFormationModule.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CapitalFormationModuleOwnershipTransferredIterator{contract: _CapitalFormationModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CapitalFormationModuleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CapitalFormationModule.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CapitalFormationModuleOwnershipTransferred)
				if err := _CapitalFormationModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_CapitalFormationModule *CapitalFormationModuleFilterer) ParseOwnershipTransferred(log types.Log) (*CapitalFormationModuleOwnershipTransferred, error) {
	event := new(CapitalFormationModuleOwnershipTransferred)
	if err := _CapitalFormationModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CapitalFormationModuleSnapshottedIterator is returned from FilterSnapshotted and is used to iterate over the raw logs and unpacked data for Snapshotted events raised by the CapitalFormationModule contract.
type CapitalFormationModuleSnapshottedIterator struct {
	Event *CapitalFormationModuleSnapshotted // Event containing the contract specifics and raw log

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
func (it *CapitalFormationModuleSnapshottedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CapitalFormationModuleSnapshotted)
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
		it.Event = new(CapitalFormationModuleSnapshotted)
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
func (it *CapitalFormationModuleSnapshottedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CapitalFormationModuleSnapshottedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CapitalFormationModuleSnapshotted represents a Snapshotted event raised by the CapitalFormationModule contract.
type CapitalFormationModuleSnapshotted struct {
	Agent              common.Address
	PriceCumulative    *big.Int
	BlockTimestampLast uint32
	BlockNumber        uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterSnapshotted is a free log retrieval operation binding the contract event 0xa2302c5afe828dc0a841009f8570ce4231545a1b20297ee1f15d562f55f9fa80.
//
// Solidity: event Snapshotted(address indexed agent, uint256 priceCumulative, uint32 blockTimestampLast, uint64 blockNumber)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) FilterSnapshotted(opts *bind.FilterOpts, agent []common.Address) (*CapitalFormationModuleSnapshottedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _CapitalFormationModule.contract.FilterLogs(opts, "Snapshotted", agentRule)
	if err != nil {
		return nil, err
	}
	return &CapitalFormationModuleSnapshottedIterator{contract: _CapitalFormationModule.contract, event: "Snapshotted", logs: logs, sub: sub}, nil
}

// WatchSnapshotted is a free log subscription operation binding the contract event 0xa2302c5afe828dc0a841009f8570ce4231545a1b20297ee1f15d562f55f9fa80.
//
// Solidity: event Snapshotted(address indexed agent, uint256 priceCumulative, uint32 blockTimestampLast, uint64 blockNumber)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) WatchSnapshotted(opts *bind.WatchOpts, sink chan<- *CapitalFormationModuleSnapshotted, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _CapitalFormationModule.contract.WatchLogs(opts, "Snapshotted", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CapitalFormationModuleSnapshotted)
				if err := _CapitalFormationModule.contract.UnpackLog(event, "Snapshotted", log); err != nil {
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

// ParseSnapshotted is a log parse operation binding the contract event 0xa2302c5afe828dc0a841009f8570ce4231545a1b20297ee1f15d562f55f9fa80.
//
// Solidity: event Snapshotted(address indexed agent, uint256 priceCumulative, uint32 blockTimestampLast, uint64 blockNumber)
func (_CapitalFormationModule *CapitalFormationModuleFilterer) ParseSnapshotted(log types.Log) (*CapitalFormationModuleSnapshotted, error) {
	event := new(CapitalFormationModuleSnapshotted)
	if err := _CapitalFormationModule.contract.UnpackLog(event, "Snapshotted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
