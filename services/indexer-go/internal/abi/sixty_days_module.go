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

// SixtyDaysModuleMetaData contains all meta data concerning the SixtyDaysModule contract.
var SixtyDaysModuleMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"WINDOW\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"accrueEscrow\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"accrueRefundPool\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimRefund\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimed\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"commit\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"configs\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"startTime\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"windowEnd\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"escrowToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bondingCurve\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"phase\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"configured\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"escrowAmountAccumulated\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"refundPool\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalContributed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"configure\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"escrowToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bondingCurve_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"creator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"startTime_\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"contributed\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"phaseOf\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"phase\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"previewRefund\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordContribution\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"refund\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Committed\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Configured\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"escrowToken\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"bondingCurve\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"startTime\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"windowEnd\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Contributed\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalContributed\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EscrowAccrued\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalEscrow\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RefundClaimed\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RefundPoolToppedUp\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalPool\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Refunded\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"escrowAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyClaimed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AlreadyConfigured\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoContribution\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotBondingCurve\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotCreator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInRefundPhase\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOpen\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotRefunding\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"WindowNotElapsed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]}]",
}

// SixtyDaysModuleABI is the input ABI used to generate the binding from.
// Deprecated: Use SixtyDaysModuleMetaData.ABI instead.
var SixtyDaysModuleABI = SixtyDaysModuleMetaData.ABI

// SixtyDaysModule is an auto generated Go binding around an Ethereum contract.
type SixtyDaysModule struct {
	SixtyDaysModuleCaller     // Read-only binding to the contract
	SixtyDaysModuleTransactor // Write-only binding to the contract
	SixtyDaysModuleFilterer   // Log filterer for contract events
}

// SixtyDaysModuleCaller is an auto generated read-only Go binding around an Ethereum contract.
type SixtyDaysModuleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SixtyDaysModuleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SixtyDaysModuleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SixtyDaysModuleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SixtyDaysModuleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SixtyDaysModuleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SixtyDaysModuleSession struct {
	Contract     *SixtyDaysModule  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SixtyDaysModuleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SixtyDaysModuleCallerSession struct {
	Contract *SixtyDaysModuleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// SixtyDaysModuleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SixtyDaysModuleTransactorSession struct {
	Contract     *SixtyDaysModuleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// SixtyDaysModuleRaw is an auto generated low-level Go binding around an Ethereum contract.
type SixtyDaysModuleRaw struct {
	Contract *SixtyDaysModule // Generic contract binding to access the raw methods on
}

// SixtyDaysModuleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SixtyDaysModuleCallerRaw struct {
	Contract *SixtyDaysModuleCaller // Generic read-only contract binding to access the raw methods on
}

// SixtyDaysModuleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SixtyDaysModuleTransactorRaw struct {
	Contract *SixtyDaysModuleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSixtyDaysModule creates a new instance of SixtyDaysModule, bound to a specific deployed contract.
func NewSixtyDaysModule(address common.Address, backend bind.ContractBackend) (*SixtyDaysModule, error) {
	contract, err := bindSixtyDaysModule(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModule{SixtyDaysModuleCaller: SixtyDaysModuleCaller{contract: contract}, SixtyDaysModuleTransactor: SixtyDaysModuleTransactor{contract: contract}, SixtyDaysModuleFilterer: SixtyDaysModuleFilterer{contract: contract}}, nil
}

// NewSixtyDaysModuleCaller creates a new read-only instance of SixtyDaysModule, bound to a specific deployed contract.
func NewSixtyDaysModuleCaller(address common.Address, caller bind.ContractCaller) (*SixtyDaysModuleCaller, error) {
	contract, err := bindSixtyDaysModule(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleCaller{contract: contract}, nil
}

// NewSixtyDaysModuleTransactor creates a new write-only instance of SixtyDaysModule, bound to a specific deployed contract.
func NewSixtyDaysModuleTransactor(address common.Address, transactor bind.ContractTransactor) (*SixtyDaysModuleTransactor, error) {
	contract, err := bindSixtyDaysModule(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleTransactor{contract: contract}, nil
}

// NewSixtyDaysModuleFilterer creates a new log filterer instance of SixtyDaysModule, bound to a specific deployed contract.
func NewSixtyDaysModuleFilterer(address common.Address, filterer bind.ContractFilterer) (*SixtyDaysModuleFilterer, error) {
	contract, err := bindSixtyDaysModule(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleFilterer{contract: contract}, nil
}

// bindSixtyDaysModule binds a generic wrapper to an already deployed contract.
func bindSixtyDaysModule(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SixtyDaysModuleMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SixtyDaysModule *SixtyDaysModuleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SixtyDaysModule.Contract.SixtyDaysModuleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SixtyDaysModule *SixtyDaysModuleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.SixtyDaysModuleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SixtyDaysModule *SixtyDaysModuleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.SixtyDaysModuleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SixtyDaysModule *SixtyDaysModuleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SixtyDaysModule.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SixtyDaysModule *SixtyDaysModuleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SixtyDaysModule *SixtyDaysModuleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.contract.Transact(opts, method, params...)
}

// WINDOW is a free data retrieval call binding the contract method 0xa53e5412.
//
// Solidity: function WINDOW() view returns(uint64)
func (_SixtyDaysModule *SixtyDaysModuleCaller) WINDOW(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _SixtyDaysModule.contract.Call(opts, &out, "WINDOW")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// WINDOW is a free data retrieval call binding the contract method 0xa53e5412.
//
// Solidity: function WINDOW() view returns(uint64)
func (_SixtyDaysModule *SixtyDaysModuleSession) WINDOW() (uint64, error) {
	return _SixtyDaysModule.Contract.WINDOW(&_SixtyDaysModule.CallOpts)
}

// WINDOW is a free data retrieval call binding the contract method 0xa53e5412.
//
// Solidity: function WINDOW() view returns(uint64)
func (_SixtyDaysModule *SixtyDaysModuleCallerSession) WINDOW() (uint64, error) {
	return _SixtyDaysModule.Contract.WINDOW(&_SixtyDaysModule.CallOpts)
}

// Claimed is a free data retrieval call binding the contract method 0x0c9cbf0e.
//
// Solidity: function claimed(address agent, address user) view returns(bool)
func (_SixtyDaysModule *SixtyDaysModuleCaller) Claimed(opts *bind.CallOpts, agent common.Address, user common.Address) (bool, error) {
	var out []interface{}
	err := _SixtyDaysModule.contract.Call(opts, &out, "claimed", agent, user)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Claimed is a free data retrieval call binding the contract method 0x0c9cbf0e.
//
// Solidity: function claimed(address agent, address user) view returns(bool)
func (_SixtyDaysModule *SixtyDaysModuleSession) Claimed(agent common.Address, user common.Address) (bool, error) {
	return _SixtyDaysModule.Contract.Claimed(&_SixtyDaysModule.CallOpts, agent, user)
}

// Claimed is a free data retrieval call binding the contract method 0x0c9cbf0e.
//
// Solidity: function claimed(address agent, address user) view returns(bool)
func (_SixtyDaysModule *SixtyDaysModuleCallerSession) Claimed(agent common.Address, user common.Address) (bool, error) {
	return _SixtyDaysModule.Contract.Claimed(&_SixtyDaysModule.CallOpts, agent, user)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(uint64 startTime, uint64 windowEnd, address escrowToken, address bondingCurve, address creator, uint8 phase, bool configured, uint256 escrowAmountAccumulated, uint256 refundPool, uint256 totalContributed)
func (_SixtyDaysModule *SixtyDaysModuleCaller) Configs(opts *bind.CallOpts, agent common.Address) (struct {
	StartTime               uint64
	WindowEnd               uint64
	EscrowToken             common.Address
	BondingCurve            common.Address
	Creator                 common.Address
	Phase                   uint8
	Configured              bool
	EscrowAmountAccumulated *big.Int
	RefundPool              *big.Int
	TotalContributed        *big.Int
}, error) {
	var out []interface{}
	err := _SixtyDaysModule.contract.Call(opts, &out, "configs", agent)

	outstruct := new(struct {
		StartTime               uint64
		WindowEnd               uint64
		EscrowToken             common.Address
		BondingCurve            common.Address
		Creator                 common.Address
		Phase                   uint8
		Configured              bool
		EscrowAmountAccumulated *big.Int
		RefundPool              *big.Int
		TotalContributed        *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.StartTime = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.WindowEnd = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.EscrowToken = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.BondingCurve = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Creator = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.Phase = *abi.ConvertType(out[5], new(uint8)).(*uint8)
	outstruct.Configured = *abi.ConvertType(out[6], new(bool)).(*bool)
	outstruct.EscrowAmountAccumulated = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.RefundPool = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.TotalContributed = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(uint64 startTime, uint64 windowEnd, address escrowToken, address bondingCurve, address creator, uint8 phase, bool configured, uint256 escrowAmountAccumulated, uint256 refundPool, uint256 totalContributed)
func (_SixtyDaysModule *SixtyDaysModuleSession) Configs(agent common.Address) (struct {
	StartTime               uint64
	WindowEnd               uint64
	EscrowToken             common.Address
	BondingCurve            common.Address
	Creator                 common.Address
	Phase                   uint8
	Configured              bool
	EscrowAmountAccumulated *big.Int
	RefundPool              *big.Int
	TotalContributed        *big.Int
}, error) {
	return _SixtyDaysModule.Contract.Configs(&_SixtyDaysModule.CallOpts, agent)
}

// Configs is a free data retrieval call binding the contract method 0xfce89878.
//
// Solidity: function configs(address agent) view returns(uint64 startTime, uint64 windowEnd, address escrowToken, address bondingCurve, address creator, uint8 phase, bool configured, uint256 escrowAmountAccumulated, uint256 refundPool, uint256 totalContributed)
func (_SixtyDaysModule *SixtyDaysModuleCallerSession) Configs(agent common.Address) (struct {
	StartTime               uint64
	WindowEnd               uint64
	EscrowToken             common.Address
	BondingCurve            common.Address
	Creator                 common.Address
	Phase                   uint8
	Configured              bool
	EscrowAmountAccumulated *big.Int
	RefundPool              *big.Int
	TotalContributed        *big.Int
}, error) {
	return _SixtyDaysModule.Contract.Configs(&_SixtyDaysModule.CallOpts, agent)
}

// Contributed is a free data retrieval call binding the contract method 0x108ec44e.
//
// Solidity: function contributed(address agent, address user) view returns(uint256)
func (_SixtyDaysModule *SixtyDaysModuleCaller) Contributed(opts *bind.CallOpts, agent common.Address, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SixtyDaysModule.contract.Call(opts, &out, "contributed", agent, user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Contributed is a free data retrieval call binding the contract method 0x108ec44e.
//
// Solidity: function contributed(address agent, address user) view returns(uint256)
func (_SixtyDaysModule *SixtyDaysModuleSession) Contributed(agent common.Address, user common.Address) (*big.Int, error) {
	return _SixtyDaysModule.Contract.Contributed(&_SixtyDaysModule.CallOpts, agent, user)
}

// Contributed is a free data retrieval call binding the contract method 0x108ec44e.
//
// Solidity: function contributed(address agent, address user) view returns(uint256)
func (_SixtyDaysModule *SixtyDaysModuleCallerSession) Contributed(agent common.Address, user common.Address) (*big.Int, error) {
	return _SixtyDaysModule.Contract.Contributed(&_SixtyDaysModule.CallOpts, agent, user)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SixtyDaysModule *SixtyDaysModuleCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SixtyDaysModule.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SixtyDaysModule *SixtyDaysModuleSession) Owner() (common.Address, error) {
	return _SixtyDaysModule.Contract.Owner(&_SixtyDaysModule.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SixtyDaysModule *SixtyDaysModuleCallerSession) Owner() (common.Address, error) {
	return _SixtyDaysModule.Contract.Owner(&_SixtyDaysModule.CallOpts)
}

// PhaseOf is a free data retrieval call binding the contract method 0xaab94cef.
//
// Solidity: function phaseOf(address agent) view returns(uint8 phase)
func (_SixtyDaysModule *SixtyDaysModuleCaller) PhaseOf(opts *bind.CallOpts, agent common.Address) (uint8, error) {
	var out []interface{}
	err := _SixtyDaysModule.contract.Call(opts, &out, "phaseOf", agent)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// PhaseOf is a free data retrieval call binding the contract method 0xaab94cef.
//
// Solidity: function phaseOf(address agent) view returns(uint8 phase)
func (_SixtyDaysModule *SixtyDaysModuleSession) PhaseOf(agent common.Address) (uint8, error) {
	return _SixtyDaysModule.Contract.PhaseOf(&_SixtyDaysModule.CallOpts, agent)
}

// PhaseOf is a free data retrieval call binding the contract method 0xaab94cef.
//
// Solidity: function phaseOf(address agent) view returns(uint8 phase)
func (_SixtyDaysModule *SixtyDaysModuleCallerSession) PhaseOf(agent common.Address) (uint8, error) {
	return _SixtyDaysModule.Contract.PhaseOf(&_SixtyDaysModule.CallOpts, agent)
}

// PreviewRefund is a free data retrieval call binding the contract method 0x94187bf0.
//
// Solidity: function previewRefund(address agent, address user) view returns(uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleCaller) PreviewRefund(opts *bind.CallOpts, agent common.Address, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SixtyDaysModule.contract.Call(opts, &out, "previewRefund", agent, user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewRefund is a free data retrieval call binding the contract method 0x94187bf0.
//
// Solidity: function previewRefund(address agent, address user) view returns(uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleSession) PreviewRefund(agent common.Address, user common.Address) (*big.Int, error) {
	return _SixtyDaysModule.Contract.PreviewRefund(&_SixtyDaysModule.CallOpts, agent, user)
}

// PreviewRefund is a free data retrieval call binding the contract method 0x94187bf0.
//
// Solidity: function previewRefund(address agent, address user) view returns(uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleCallerSession) PreviewRefund(agent common.Address, user common.Address) (*big.Int, error) {
	return _SixtyDaysModule.Contract.PreviewRefund(&_SixtyDaysModule.CallOpts, agent, user)
}

// AccrueEscrow is a paid mutator transaction binding the contract method 0x50408873.
//
// Solidity: function accrueEscrow(address agent, uint256 amount) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactor) AccrueEscrow(opts *bind.TransactOpts, agent common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SixtyDaysModule.contract.Transact(opts, "accrueEscrow", agent, amount)
}

// AccrueEscrow is a paid mutator transaction binding the contract method 0x50408873.
//
// Solidity: function accrueEscrow(address agent, uint256 amount) returns()
func (_SixtyDaysModule *SixtyDaysModuleSession) AccrueEscrow(agent common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.AccrueEscrow(&_SixtyDaysModule.TransactOpts, agent, amount)
}

// AccrueEscrow is a paid mutator transaction binding the contract method 0x50408873.
//
// Solidity: function accrueEscrow(address agent, uint256 amount) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactorSession) AccrueEscrow(agent common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.AccrueEscrow(&_SixtyDaysModule.TransactOpts, agent, amount)
}

// AccrueRefundPool is a paid mutator transaction binding the contract method 0x67eb4cc5.
//
// Solidity: function accrueRefundPool(address agent, uint256 amount) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactor) AccrueRefundPool(opts *bind.TransactOpts, agent common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SixtyDaysModule.contract.Transact(opts, "accrueRefundPool", agent, amount)
}

// AccrueRefundPool is a paid mutator transaction binding the contract method 0x67eb4cc5.
//
// Solidity: function accrueRefundPool(address agent, uint256 amount) returns()
func (_SixtyDaysModule *SixtyDaysModuleSession) AccrueRefundPool(agent common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.AccrueRefundPool(&_SixtyDaysModule.TransactOpts, agent, amount)
}

// AccrueRefundPool is a paid mutator transaction binding the contract method 0x67eb4cc5.
//
// Solidity: function accrueRefundPool(address agent, uint256 amount) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactorSession) AccrueRefundPool(agent common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.AccrueRefundPool(&_SixtyDaysModule.TransactOpts, agent, amount)
}

// ClaimRefund is a paid mutator transaction binding the contract method 0xbffa55d5.
//
// Solidity: function claimRefund(address agent) returns(uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleTransactor) ClaimRefund(opts *bind.TransactOpts, agent common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.contract.Transact(opts, "claimRefund", agent)
}

// ClaimRefund is a paid mutator transaction binding the contract method 0xbffa55d5.
//
// Solidity: function claimRefund(address agent) returns(uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleSession) ClaimRefund(agent common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.ClaimRefund(&_SixtyDaysModule.TransactOpts, agent)
}

// ClaimRefund is a paid mutator transaction binding the contract method 0xbffa55d5.
//
// Solidity: function claimRefund(address agent) returns(uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleTransactorSession) ClaimRefund(agent common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.ClaimRefund(&_SixtyDaysModule.TransactOpts, agent)
}

// Commit is a paid mutator transaction binding the contract method 0x369e8c1d.
//
// Solidity: function commit(address agent) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactor) Commit(opts *bind.TransactOpts, agent common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.contract.Transact(opts, "commit", agent)
}

// Commit is a paid mutator transaction binding the contract method 0x369e8c1d.
//
// Solidity: function commit(address agent) returns()
func (_SixtyDaysModule *SixtyDaysModuleSession) Commit(agent common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.Commit(&_SixtyDaysModule.TransactOpts, agent)
}

// Commit is a paid mutator transaction binding the contract method 0x369e8c1d.
//
// Solidity: function commit(address agent) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactorSession) Commit(agent common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.Commit(&_SixtyDaysModule.TransactOpts, agent)
}

// Configure is a paid mutator transaction binding the contract method 0xdb2ad89d.
//
// Solidity: function configure(address agent, address escrowToken_, address bondingCurve_, address creator_, uint64 startTime_) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactor) Configure(opts *bind.TransactOpts, agent common.Address, escrowToken_ common.Address, bondingCurve_ common.Address, creator_ common.Address, startTime_ uint64) (*types.Transaction, error) {
	return _SixtyDaysModule.contract.Transact(opts, "configure", agent, escrowToken_, bondingCurve_, creator_, startTime_)
}

// Configure is a paid mutator transaction binding the contract method 0xdb2ad89d.
//
// Solidity: function configure(address agent, address escrowToken_, address bondingCurve_, address creator_, uint64 startTime_) returns()
func (_SixtyDaysModule *SixtyDaysModuleSession) Configure(agent common.Address, escrowToken_ common.Address, bondingCurve_ common.Address, creator_ common.Address, startTime_ uint64) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.Configure(&_SixtyDaysModule.TransactOpts, agent, escrowToken_, bondingCurve_, creator_, startTime_)
}

// Configure is a paid mutator transaction binding the contract method 0xdb2ad89d.
//
// Solidity: function configure(address agent, address escrowToken_, address bondingCurve_, address creator_, uint64 startTime_) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactorSession) Configure(agent common.Address, escrowToken_ common.Address, bondingCurve_ common.Address, creator_ common.Address, startTime_ uint64) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.Configure(&_SixtyDaysModule.TransactOpts, agent, escrowToken_, bondingCurve_, creator_, startTime_)
}

// RecordContribution is a paid mutator transaction binding the contract method 0xfadbf92e.
//
// Solidity: function recordContribution(address agent, address user, uint256 amount) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactor) RecordContribution(opts *bind.TransactOpts, agent common.Address, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SixtyDaysModule.contract.Transact(opts, "recordContribution", agent, user, amount)
}

// RecordContribution is a paid mutator transaction binding the contract method 0xfadbf92e.
//
// Solidity: function recordContribution(address agent, address user, uint256 amount) returns()
func (_SixtyDaysModule *SixtyDaysModuleSession) RecordContribution(agent common.Address, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.RecordContribution(&_SixtyDaysModule.TransactOpts, agent, user, amount)
}

// RecordContribution is a paid mutator transaction binding the contract method 0xfadbf92e.
//
// Solidity: function recordContribution(address agent, address user, uint256 amount) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactorSession) RecordContribution(agent common.Address, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.RecordContribution(&_SixtyDaysModule.TransactOpts, agent, user, amount)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address agent) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactor) Refund(opts *bind.TransactOpts, agent common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.contract.Transact(opts, "refund", agent)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address agent) returns()
func (_SixtyDaysModule *SixtyDaysModuleSession) Refund(agent common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.Refund(&_SixtyDaysModule.TransactOpts, agent)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address agent) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactorSession) Refund(agent common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.Refund(&_SixtyDaysModule.TransactOpts, agent)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SixtyDaysModule.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SixtyDaysModule *SixtyDaysModuleSession) RenounceOwnership() (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.RenounceOwnership(&_SixtyDaysModule.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.RenounceOwnership(&_SixtyDaysModule.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SixtyDaysModule *SixtyDaysModuleSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.TransferOwnership(&_SixtyDaysModule.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SixtyDaysModule *SixtyDaysModuleTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SixtyDaysModule.Contract.TransferOwnership(&_SixtyDaysModule.TransactOpts, newOwner)
}

// SixtyDaysModuleCommittedIterator is returned from FilterCommitted and is used to iterate over the raw logs and unpacked data for Committed events raised by the SixtyDaysModule contract.
type SixtyDaysModuleCommittedIterator struct {
	Event *SixtyDaysModuleCommitted // Event containing the contract specifics and raw log

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
func (it *SixtyDaysModuleCommittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SixtyDaysModuleCommitted)
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
		it.Event = new(SixtyDaysModuleCommitted)
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
func (it *SixtyDaysModuleCommittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SixtyDaysModuleCommittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SixtyDaysModuleCommitted represents a Committed event raised by the SixtyDaysModule contract.
type SixtyDaysModuleCommitted struct {
	Agent   common.Address
	Creator common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterCommitted is a free log retrieval operation binding the contract event 0xa5e4e34711d102d6aaadc9166b6c0ca95c3c5c1e55a45836d9d1d4f7ccb329e1.
//
// Solidity: event Committed(address indexed agent, address indexed creator, uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) FilterCommitted(opts *bind.FilterOpts, agent []common.Address, creator []common.Address) (*SixtyDaysModuleCommittedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.FilterLogs(opts, "Committed", agentRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleCommittedIterator{contract: _SixtyDaysModule.contract, event: "Committed", logs: logs, sub: sub}, nil
}

// WatchCommitted is a free log subscription operation binding the contract event 0xa5e4e34711d102d6aaadc9166b6c0ca95c3c5c1e55a45836d9d1d4f7ccb329e1.
//
// Solidity: event Committed(address indexed agent, address indexed creator, uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) WatchCommitted(opts *bind.WatchOpts, sink chan<- *SixtyDaysModuleCommitted, agent []common.Address, creator []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.WatchLogs(opts, "Committed", agentRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SixtyDaysModuleCommitted)
				if err := _SixtyDaysModule.contract.UnpackLog(event, "Committed", log); err != nil {
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

// ParseCommitted is a log parse operation binding the contract event 0xa5e4e34711d102d6aaadc9166b6c0ca95c3c5c1e55a45836d9d1d4f7ccb329e1.
//
// Solidity: event Committed(address indexed agent, address indexed creator, uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) ParseCommitted(log types.Log) (*SixtyDaysModuleCommitted, error) {
	event := new(SixtyDaysModuleCommitted)
	if err := _SixtyDaysModule.contract.UnpackLog(event, "Committed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SixtyDaysModuleConfiguredIterator is returned from FilterConfigured and is used to iterate over the raw logs and unpacked data for Configured events raised by the SixtyDaysModule contract.
type SixtyDaysModuleConfiguredIterator struct {
	Event *SixtyDaysModuleConfigured // Event containing the contract specifics and raw log

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
func (it *SixtyDaysModuleConfiguredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SixtyDaysModuleConfigured)
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
		it.Event = new(SixtyDaysModuleConfigured)
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
func (it *SixtyDaysModuleConfiguredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SixtyDaysModuleConfiguredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SixtyDaysModuleConfigured represents a Configured event raised by the SixtyDaysModule contract.
type SixtyDaysModuleConfigured struct {
	Agent        common.Address
	EscrowToken  common.Address
	BondingCurve common.Address
	Creator      common.Address
	StartTime    uint64
	WindowEnd    uint64
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterConfigured is a free log retrieval operation binding the contract event 0x63d7d3f77acbd3002248a8835879046bbbce05b2af2b8baa23d42a1b778b1aea.
//
// Solidity: event Configured(address indexed agent, address indexed escrowToken, address indexed bondingCurve, address creator, uint64 startTime, uint64 windowEnd)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) FilterConfigured(opts *bind.FilterOpts, agent []common.Address, escrowToken []common.Address, bondingCurve []common.Address) (*SixtyDaysModuleConfiguredIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var escrowTokenRule []interface{}
	for _, escrowTokenItem := range escrowToken {
		escrowTokenRule = append(escrowTokenRule, escrowTokenItem)
	}
	var bondingCurveRule []interface{}
	for _, bondingCurveItem := range bondingCurve {
		bondingCurveRule = append(bondingCurveRule, bondingCurveItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.FilterLogs(opts, "Configured", agentRule, escrowTokenRule, bondingCurveRule)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleConfiguredIterator{contract: _SixtyDaysModule.contract, event: "Configured", logs: logs, sub: sub}, nil
}

// WatchConfigured is a free log subscription operation binding the contract event 0x63d7d3f77acbd3002248a8835879046bbbce05b2af2b8baa23d42a1b778b1aea.
//
// Solidity: event Configured(address indexed agent, address indexed escrowToken, address indexed bondingCurve, address creator, uint64 startTime, uint64 windowEnd)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) WatchConfigured(opts *bind.WatchOpts, sink chan<- *SixtyDaysModuleConfigured, agent []common.Address, escrowToken []common.Address, bondingCurve []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var escrowTokenRule []interface{}
	for _, escrowTokenItem := range escrowToken {
		escrowTokenRule = append(escrowTokenRule, escrowTokenItem)
	}
	var bondingCurveRule []interface{}
	for _, bondingCurveItem := range bondingCurve {
		bondingCurveRule = append(bondingCurveRule, bondingCurveItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.WatchLogs(opts, "Configured", agentRule, escrowTokenRule, bondingCurveRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SixtyDaysModuleConfigured)
				if err := _SixtyDaysModule.contract.UnpackLog(event, "Configured", log); err != nil {
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

// ParseConfigured is a log parse operation binding the contract event 0x63d7d3f77acbd3002248a8835879046bbbce05b2af2b8baa23d42a1b778b1aea.
//
// Solidity: event Configured(address indexed agent, address indexed escrowToken, address indexed bondingCurve, address creator, uint64 startTime, uint64 windowEnd)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) ParseConfigured(log types.Log) (*SixtyDaysModuleConfigured, error) {
	event := new(SixtyDaysModuleConfigured)
	if err := _SixtyDaysModule.contract.UnpackLog(event, "Configured", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SixtyDaysModuleContributedIterator is returned from FilterContributed and is used to iterate over the raw logs and unpacked data for Contributed events raised by the SixtyDaysModule contract.
type SixtyDaysModuleContributedIterator struct {
	Event *SixtyDaysModuleContributed // Event containing the contract specifics and raw log

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
func (it *SixtyDaysModuleContributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SixtyDaysModuleContributed)
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
		it.Event = new(SixtyDaysModuleContributed)
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
func (it *SixtyDaysModuleContributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SixtyDaysModuleContributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SixtyDaysModuleContributed represents a Contributed event raised by the SixtyDaysModule contract.
type SixtyDaysModuleContributed struct {
	Agent            common.Address
	User             common.Address
	Amount           *big.Int
	TotalContributed *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterContributed is a free log retrieval operation binding the contract event 0x44a1ddc51ace1542c482b6e70a4f0300716de6ac2f749e7f9d44e590806e3fd4.
//
// Solidity: event Contributed(address indexed agent, address indexed user, uint256 amount, uint256 totalContributed)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) FilterContributed(opts *bind.FilterOpts, agent []common.Address, user []common.Address) (*SixtyDaysModuleContributedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.FilterLogs(opts, "Contributed", agentRule, userRule)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleContributedIterator{contract: _SixtyDaysModule.contract, event: "Contributed", logs: logs, sub: sub}, nil
}

// WatchContributed is a free log subscription operation binding the contract event 0x44a1ddc51ace1542c482b6e70a4f0300716de6ac2f749e7f9d44e590806e3fd4.
//
// Solidity: event Contributed(address indexed agent, address indexed user, uint256 amount, uint256 totalContributed)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) WatchContributed(opts *bind.WatchOpts, sink chan<- *SixtyDaysModuleContributed, agent []common.Address, user []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.WatchLogs(opts, "Contributed", agentRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SixtyDaysModuleContributed)
				if err := _SixtyDaysModule.contract.UnpackLog(event, "Contributed", log); err != nil {
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

// ParseContributed is a log parse operation binding the contract event 0x44a1ddc51ace1542c482b6e70a4f0300716de6ac2f749e7f9d44e590806e3fd4.
//
// Solidity: event Contributed(address indexed agent, address indexed user, uint256 amount, uint256 totalContributed)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) ParseContributed(log types.Log) (*SixtyDaysModuleContributed, error) {
	event := new(SixtyDaysModuleContributed)
	if err := _SixtyDaysModule.contract.UnpackLog(event, "Contributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SixtyDaysModuleEscrowAccruedIterator is returned from FilterEscrowAccrued and is used to iterate over the raw logs and unpacked data for EscrowAccrued events raised by the SixtyDaysModule contract.
type SixtyDaysModuleEscrowAccruedIterator struct {
	Event *SixtyDaysModuleEscrowAccrued // Event containing the contract specifics and raw log

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
func (it *SixtyDaysModuleEscrowAccruedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SixtyDaysModuleEscrowAccrued)
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
		it.Event = new(SixtyDaysModuleEscrowAccrued)
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
func (it *SixtyDaysModuleEscrowAccruedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SixtyDaysModuleEscrowAccruedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SixtyDaysModuleEscrowAccrued represents a EscrowAccrued event raised by the SixtyDaysModule contract.
type SixtyDaysModuleEscrowAccrued struct {
	Agent       common.Address
	Amount      *big.Int
	TotalEscrow *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterEscrowAccrued is a free log retrieval operation binding the contract event 0xd673bd2736331a88f60c0a1effedd3d4a3e105a0aad13bca58eb9825fea6d372.
//
// Solidity: event EscrowAccrued(address indexed agent, uint256 amount, uint256 totalEscrow)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) FilterEscrowAccrued(opts *bind.FilterOpts, agent []common.Address) (*SixtyDaysModuleEscrowAccruedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.FilterLogs(opts, "EscrowAccrued", agentRule)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleEscrowAccruedIterator{contract: _SixtyDaysModule.contract, event: "EscrowAccrued", logs: logs, sub: sub}, nil
}

// WatchEscrowAccrued is a free log subscription operation binding the contract event 0xd673bd2736331a88f60c0a1effedd3d4a3e105a0aad13bca58eb9825fea6d372.
//
// Solidity: event EscrowAccrued(address indexed agent, uint256 amount, uint256 totalEscrow)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) WatchEscrowAccrued(opts *bind.WatchOpts, sink chan<- *SixtyDaysModuleEscrowAccrued, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.WatchLogs(opts, "EscrowAccrued", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SixtyDaysModuleEscrowAccrued)
				if err := _SixtyDaysModule.contract.UnpackLog(event, "EscrowAccrued", log); err != nil {
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

// ParseEscrowAccrued is a log parse operation binding the contract event 0xd673bd2736331a88f60c0a1effedd3d4a3e105a0aad13bca58eb9825fea6d372.
//
// Solidity: event EscrowAccrued(address indexed agent, uint256 amount, uint256 totalEscrow)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) ParseEscrowAccrued(log types.Log) (*SixtyDaysModuleEscrowAccrued, error) {
	event := new(SixtyDaysModuleEscrowAccrued)
	if err := _SixtyDaysModule.contract.UnpackLog(event, "EscrowAccrued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SixtyDaysModuleOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SixtyDaysModule contract.
type SixtyDaysModuleOwnershipTransferredIterator struct {
	Event *SixtyDaysModuleOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SixtyDaysModuleOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SixtyDaysModuleOwnershipTransferred)
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
		it.Event = new(SixtyDaysModuleOwnershipTransferred)
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
func (it *SixtyDaysModuleOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SixtyDaysModuleOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SixtyDaysModuleOwnershipTransferred represents a OwnershipTransferred event raised by the SixtyDaysModule contract.
type SixtyDaysModuleOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SixtyDaysModuleOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleOwnershipTransferredIterator{contract: _SixtyDaysModule.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SixtyDaysModuleOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SixtyDaysModuleOwnershipTransferred)
				if err := _SixtyDaysModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SixtyDaysModule *SixtyDaysModuleFilterer) ParseOwnershipTransferred(log types.Log) (*SixtyDaysModuleOwnershipTransferred, error) {
	event := new(SixtyDaysModuleOwnershipTransferred)
	if err := _SixtyDaysModule.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SixtyDaysModuleRefundClaimedIterator is returned from FilterRefundClaimed and is used to iterate over the raw logs and unpacked data for RefundClaimed events raised by the SixtyDaysModule contract.
type SixtyDaysModuleRefundClaimedIterator struct {
	Event *SixtyDaysModuleRefundClaimed // Event containing the contract specifics and raw log

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
func (it *SixtyDaysModuleRefundClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SixtyDaysModuleRefundClaimed)
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
		it.Event = new(SixtyDaysModuleRefundClaimed)
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
func (it *SixtyDaysModuleRefundClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SixtyDaysModuleRefundClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SixtyDaysModuleRefundClaimed represents a RefundClaimed event raised by the SixtyDaysModule contract.
type SixtyDaysModuleRefundClaimed struct {
	Agent  common.Address
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRefundClaimed is a free log retrieval operation binding the contract event 0x39bed68a008a68cbf907d7ff6bc3629912af6516cb837cfa3f871ad9f2b8a944.
//
// Solidity: event RefundClaimed(address indexed agent, address indexed user, uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) FilterRefundClaimed(opts *bind.FilterOpts, agent []common.Address, user []common.Address) (*SixtyDaysModuleRefundClaimedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.FilterLogs(opts, "RefundClaimed", agentRule, userRule)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleRefundClaimedIterator{contract: _SixtyDaysModule.contract, event: "RefundClaimed", logs: logs, sub: sub}, nil
}

// WatchRefundClaimed is a free log subscription operation binding the contract event 0x39bed68a008a68cbf907d7ff6bc3629912af6516cb837cfa3f871ad9f2b8a944.
//
// Solidity: event RefundClaimed(address indexed agent, address indexed user, uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) WatchRefundClaimed(opts *bind.WatchOpts, sink chan<- *SixtyDaysModuleRefundClaimed, agent []common.Address, user []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.WatchLogs(opts, "RefundClaimed", agentRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SixtyDaysModuleRefundClaimed)
				if err := _SixtyDaysModule.contract.UnpackLog(event, "RefundClaimed", log); err != nil {
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

// ParseRefundClaimed is a log parse operation binding the contract event 0x39bed68a008a68cbf907d7ff6bc3629912af6516cb837cfa3f871ad9f2b8a944.
//
// Solidity: event RefundClaimed(address indexed agent, address indexed user, uint256 amount)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) ParseRefundClaimed(log types.Log) (*SixtyDaysModuleRefundClaimed, error) {
	event := new(SixtyDaysModuleRefundClaimed)
	if err := _SixtyDaysModule.contract.UnpackLog(event, "RefundClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SixtyDaysModuleRefundPoolToppedUpIterator is returned from FilterRefundPoolToppedUp and is used to iterate over the raw logs and unpacked data for RefundPoolToppedUp events raised by the SixtyDaysModule contract.
type SixtyDaysModuleRefundPoolToppedUpIterator struct {
	Event *SixtyDaysModuleRefundPoolToppedUp // Event containing the contract specifics and raw log

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
func (it *SixtyDaysModuleRefundPoolToppedUpIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SixtyDaysModuleRefundPoolToppedUp)
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
		it.Event = new(SixtyDaysModuleRefundPoolToppedUp)
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
func (it *SixtyDaysModuleRefundPoolToppedUpIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SixtyDaysModuleRefundPoolToppedUpIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SixtyDaysModuleRefundPoolToppedUp represents a RefundPoolToppedUp event raised by the SixtyDaysModule contract.
type SixtyDaysModuleRefundPoolToppedUp struct {
	Agent     common.Address
	Amount    *big.Int
	TotalPool *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRefundPoolToppedUp is a free log retrieval operation binding the contract event 0x47ef087480a7d791bbc5e0550b0ca3f435515cc2548a7653c5bb5ce386ddf2d5.
//
// Solidity: event RefundPoolToppedUp(address indexed agent, uint256 amount, uint256 totalPool)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) FilterRefundPoolToppedUp(opts *bind.FilterOpts, agent []common.Address) (*SixtyDaysModuleRefundPoolToppedUpIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.FilterLogs(opts, "RefundPoolToppedUp", agentRule)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleRefundPoolToppedUpIterator{contract: _SixtyDaysModule.contract, event: "RefundPoolToppedUp", logs: logs, sub: sub}, nil
}

// WatchRefundPoolToppedUp is a free log subscription operation binding the contract event 0x47ef087480a7d791bbc5e0550b0ca3f435515cc2548a7653c5bb5ce386ddf2d5.
//
// Solidity: event RefundPoolToppedUp(address indexed agent, uint256 amount, uint256 totalPool)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) WatchRefundPoolToppedUp(opts *bind.WatchOpts, sink chan<- *SixtyDaysModuleRefundPoolToppedUp, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.WatchLogs(opts, "RefundPoolToppedUp", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SixtyDaysModuleRefundPoolToppedUp)
				if err := _SixtyDaysModule.contract.UnpackLog(event, "RefundPoolToppedUp", log); err != nil {
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

// ParseRefundPoolToppedUp is a log parse operation binding the contract event 0x47ef087480a7d791bbc5e0550b0ca3f435515cc2548a7653c5bb5ce386ddf2d5.
//
// Solidity: event RefundPoolToppedUp(address indexed agent, uint256 amount, uint256 totalPool)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) ParseRefundPoolToppedUp(log types.Log) (*SixtyDaysModuleRefundPoolToppedUp, error) {
	event := new(SixtyDaysModuleRefundPoolToppedUp)
	if err := _SixtyDaysModule.contract.UnpackLog(event, "RefundPoolToppedUp", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SixtyDaysModuleRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the SixtyDaysModule contract.
type SixtyDaysModuleRefundedIterator struct {
	Event *SixtyDaysModuleRefunded // Event containing the contract specifics and raw log

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
func (it *SixtyDaysModuleRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SixtyDaysModuleRefunded)
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
		it.Event = new(SixtyDaysModuleRefunded)
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
func (it *SixtyDaysModuleRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SixtyDaysModuleRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SixtyDaysModuleRefunded represents a Refunded event raised by the SixtyDaysModule contract.
type SixtyDaysModuleRefunded struct {
	Agent        common.Address
	Creator      common.Address
	EscrowAmount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0xec1e5ed733e00f1a00915d56caef57b4f52312dde4f9b3165f213319a0da156b.
//
// Solidity: event Refunded(address indexed agent, address indexed creator, uint256 escrowAmount)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) FilterRefunded(opts *bind.FilterOpts, agent []common.Address, creator []common.Address) (*SixtyDaysModuleRefundedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.FilterLogs(opts, "Refunded", agentRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &SixtyDaysModuleRefundedIterator{contract: _SixtyDaysModule.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0xec1e5ed733e00f1a00915d56caef57b4f52312dde4f9b3165f213319a0da156b.
//
// Solidity: event Refunded(address indexed agent, address indexed creator, uint256 escrowAmount)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *SixtyDaysModuleRefunded, agent []common.Address, creator []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SixtyDaysModule.contract.WatchLogs(opts, "Refunded", agentRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SixtyDaysModuleRefunded)
				if err := _SixtyDaysModule.contract.UnpackLog(event, "Refunded", log); err != nil {
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

// ParseRefunded is a log parse operation binding the contract event 0xec1e5ed733e00f1a00915d56caef57b4f52312dde4f9b3165f213319a0da156b.
//
// Solidity: event Refunded(address indexed agent, address indexed creator, uint256 escrowAmount)
func (_SixtyDaysModule *SixtyDaysModuleFilterer) ParseRefunded(log types.Log) (*SixtyDaysModuleRefunded, error) {
	event := new(SixtyDaysModuleRefunded)
	if err := _SixtyDaysModule.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
