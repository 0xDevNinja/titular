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

// FeeRouterMetaData contains all meta data concerning the FeeRouter contract.
var FeeRouterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"treasury_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BPS_DENOMINATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"DEFAULT_CREATOR_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"clearRoute\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"distribute\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"contractIERC20\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"distributeNative\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getSplit\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"creatorAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"treasuryAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"routes\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"creatorBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"configured\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setRoute\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"creatorBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Distributed\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"creatorAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"treasuryAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouteCleared\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RouteSet\",\"inputs\":[{\"name\":\"agent\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"creatorBps\",\"type\":\"uint16\",\"indexed\":false,\"internalType\":\"uint16\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"BpsTooHigh\",\"inputs\":[{\"name\":\"creatorBps\",\"type\":\"uint16\",\"internalType\":\"uint16\"}]},{\"type\":\"error\",\"name\":\"InsufficientValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NativeTransferFailed\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// FeeRouterABI is the input ABI used to generate the binding from.
// Deprecated: Use FeeRouterMetaData.ABI instead.
var FeeRouterABI = FeeRouterMetaData.ABI

// FeeRouter is an auto generated Go binding around an Ethereum contract.
type FeeRouter struct {
	FeeRouterCaller     // Read-only binding to the contract
	FeeRouterTransactor // Write-only binding to the contract
	FeeRouterFilterer   // Log filterer for contract events
}

// FeeRouterCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeeRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeRouterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeeRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeRouterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeeRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeRouterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeeRouterSession struct {
	Contract     *FeeRouter        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeeRouterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeeRouterCallerSession struct {
	Contract *FeeRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// FeeRouterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeeRouterTransactorSession struct {
	Contract     *FeeRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// FeeRouterRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeeRouterRaw struct {
	Contract *FeeRouter // Generic contract binding to access the raw methods on
}

// FeeRouterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeeRouterCallerRaw struct {
	Contract *FeeRouterCaller // Generic read-only contract binding to access the raw methods on
}

// FeeRouterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeeRouterTransactorRaw struct {
	Contract *FeeRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeeRouter creates a new instance of FeeRouter, bound to a specific deployed contract.
func NewFeeRouter(address common.Address, backend bind.ContractBackend) (*FeeRouter, error) {
	contract, err := bindFeeRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeRouter{FeeRouterCaller: FeeRouterCaller{contract: contract}, FeeRouterTransactor: FeeRouterTransactor{contract: contract}, FeeRouterFilterer: FeeRouterFilterer{contract: contract}}, nil
}

// NewFeeRouterCaller creates a new read-only instance of FeeRouter, bound to a specific deployed contract.
func NewFeeRouterCaller(address common.Address, caller bind.ContractCaller) (*FeeRouterCaller, error) {
	contract, err := bindFeeRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeRouterCaller{contract: contract}, nil
}

// NewFeeRouterTransactor creates a new write-only instance of FeeRouter, bound to a specific deployed contract.
func NewFeeRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeRouterTransactor, error) {
	contract, err := bindFeeRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeRouterTransactor{contract: contract}, nil
}

// NewFeeRouterFilterer creates a new log filterer instance of FeeRouter, bound to a specific deployed contract.
func NewFeeRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeRouterFilterer, error) {
	contract, err := bindFeeRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeRouterFilterer{contract: contract}, nil
}

// bindFeeRouter binds a generic wrapper to an already deployed contract.
func bindFeeRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeRouter *FeeRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeRouter.Contract.FeeRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeRouter *FeeRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeRouter.Contract.FeeRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeRouter *FeeRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeRouter.Contract.FeeRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeRouter *FeeRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeRouter *FeeRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeRouter *FeeRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeRouter.Contract.contract.Transact(opts, method, params...)
}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint16)
func (_FeeRouter *FeeRouterCaller) BPSDENOMINATOR(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _FeeRouter.contract.Call(opts, &out, "BPS_DENOMINATOR")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint16)
func (_FeeRouter *FeeRouterSession) BPSDENOMINATOR() (uint16, error) {
	return _FeeRouter.Contract.BPSDENOMINATOR(&_FeeRouter.CallOpts)
}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint16)
func (_FeeRouter *FeeRouterCallerSession) BPSDENOMINATOR() (uint16, error) {
	return _FeeRouter.Contract.BPSDENOMINATOR(&_FeeRouter.CallOpts)
}

// DEFAULTCREATORBPS is a free data retrieval call binding the contract method 0xe52e3820.
//
// Solidity: function DEFAULT_CREATOR_BPS() view returns(uint16)
func (_FeeRouter *FeeRouterCaller) DEFAULTCREATORBPS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _FeeRouter.contract.Call(opts, &out, "DEFAULT_CREATOR_BPS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// DEFAULTCREATORBPS is a free data retrieval call binding the contract method 0xe52e3820.
//
// Solidity: function DEFAULT_CREATOR_BPS() view returns(uint16)
func (_FeeRouter *FeeRouterSession) DEFAULTCREATORBPS() (uint16, error) {
	return _FeeRouter.Contract.DEFAULTCREATORBPS(&_FeeRouter.CallOpts)
}

// DEFAULTCREATORBPS is a free data retrieval call binding the contract method 0xe52e3820.
//
// Solidity: function DEFAULT_CREATOR_BPS() view returns(uint16)
func (_FeeRouter *FeeRouterCallerSession) DEFAULTCREATORBPS() (uint16, error) {
	return _FeeRouter.Contract.DEFAULTCREATORBPS(&_FeeRouter.CallOpts)
}

// GetSplit is a free data retrieval call binding the contract method 0x9269df26.
//
// Solidity: function getSplit(address agent, uint256 amount) view returns(uint256 creatorAmount, uint256 treasuryAmount)
func (_FeeRouter *FeeRouterCaller) GetSplit(opts *bind.CallOpts, agent common.Address, amount *big.Int) (struct {
	CreatorAmount  *big.Int
	TreasuryAmount *big.Int
}, error) {
	var out []interface{}
	err := _FeeRouter.contract.Call(opts, &out, "getSplit", agent, amount)

	outstruct := new(struct {
		CreatorAmount  *big.Int
		TreasuryAmount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CreatorAmount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TreasuryAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetSplit is a free data retrieval call binding the contract method 0x9269df26.
//
// Solidity: function getSplit(address agent, uint256 amount) view returns(uint256 creatorAmount, uint256 treasuryAmount)
func (_FeeRouter *FeeRouterSession) GetSplit(agent common.Address, amount *big.Int) (struct {
	CreatorAmount  *big.Int
	TreasuryAmount *big.Int
}, error) {
	return _FeeRouter.Contract.GetSplit(&_FeeRouter.CallOpts, agent, amount)
}

// GetSplit is a free data retrieval call binding the contract method 0x9269df26.
//
// Solidity: function getSplit(address agent, uint256 amount) view returns(uint256 creatorAmount, uint256 treasuryAmount)
func (_FeeRouter *FeeRouterCallerSession) GetSplit(agent common.Address, amount *big.Int) (struct {
	CreatorAmount  *big.Int
	TreasuryAmount *big.Int
}, error) {
	return _FeeRouter.Contract.GetSplit(&_FeeRouter.CallOpts, agent, amount)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeeRouter *FeeRouterCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeRouter.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeeRouter *FeeRouterSession) Owner() (common.Address, error) {
	return _FeeRouter.Contract.Owner(&_FeeRouter.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeeRouter *FeeRouterCallerSession) Owner() (common.Address, error) {
	return _FeeRouter.Contract.Owner(&_FeeRouter.CallOpts)
}

// Routes is a free data retrieval call binding the contract method 0xd7409659.
//
// Solidity: function routes(address agent) view returns(address creator, uint16 creatorBps, bool configured)
func (_FeeRouter *FeeRouterCaller) Routes(opts *bind.CallOpts, agent common.Address) (struct {
	Creator    common.Address
	CreatorBps uint16
	Configured bool
}, error) {
	var out []interface{}
	err := _FeeRouter.contract.Call(opts, &out, "routes", agent)

	outstruct := new(struct {
		Creator    common.Address
		CreatorBps uint16
		Configured bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Creator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.CreatorBps = *abi.ConvertType(out[1], new(uint16)).(*uint16)
	outstruct.Configured = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// Routes is a free data retrieval call binding the contract method 0xd7409659.
//
// Solidity: function routes(address agent) view returns(address creator, uint16 creatorBps, bool configured)
func (_FeeRouter *FeeRouterSession) Routes(agent common.Address) (struct {
	Creator    common.Address
	CreatorBps uint16
	Configured bool
}, error) {
	return _FeeRouter.Contract.Routes(&_FeeRouter.CallOpts, agent)
}

// Routes is a free data retrieval call binding the contract method 0xd7409659.
//
// Solidity: function routes(address agent) view returns(address creator, uint16 creatorBps, bool configured)
func (_FeeRouter *FeeRouterCallerSession) Routes(agent common.Address) (struct {
	Creator    common.Address
	CreatorBps uint16
	Configured bool
}, error) {
	return _FeeRouter.Contract.Routes(&_FeeRouter.CallOpts, agent)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_FeeRouter *FeeRouterCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeRouter.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_FeeRouter *FeeRouterSession) Treasury() (common.Address, error) {
	return _FeeRouter.Contract.Treasury(&_FeeRouter.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_FeeRouter *FeeRouterCallerSession) Treasury() (common.Address, error) {
	return _FeeRouter.Contract.Treasury(&_FeeRouter.CallOpts)
}

// ClearRoute is a paid mutator transaction binding the contract method 0x8ebd2d29.
//
// Solidity: function clearRoute(address agent) returns()
func (_FeeRouter *FeeRouterTransactor) ClearRoute(opts *bind.TransactOpts, agent common.Address) (*types.Transaction, error) {
	return _FeeRouter.contract.Transact(opts, "clearRoute", agent)
}

// ClearRoute is a paid mutator transaction binding the contract method 0x8ebd2d29.
//
// Solidity: function clearRoute(address agent) returns()
func (_FeeRouter *FeeRouterSession) ClearRoute(agent common.Address) (*types.Transaction, error) {
	return _FeeRouter.Contract.ClearRoute(&_FeeRouter.TransactOpts, agent)
}

// ClearRoute is a paid mutator transaction binding the contract method 0x8ebd2d29.
//
// Solidity: function clearRoute(address agent) returns()
func (_FeeRouter *FeeRouterTransactorSession) ClearRoute(agent common.Address) (*types.Transaction, error) {
	return _FeeRouter.Contract.ClearRoute(&_FeeRouter.TransactOpts, agent)
}

// Distribute is a paid mutator transaction binding the contract method 0xa16a9183.
//
// Solidity: function distribute(address agent, address token, uint256 amount) returns()
func (_FeeRouter *FeeRouterTransactor) Distribute(opts *bind.TransactOpts, agent common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeRouter.contract.Transact(opts, "distribute", agent, token, amount)
}

// Distribute is a paid mutator transaction binding the contract method 0xa16a9183.
//
// Solidity: function distribute(address agent, address token, uint256 amount) returns()
func (_FeeRouter *FeeRouterSession) Distribute(agent common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeRouter.Contract.Distribute(&_FeeRouter.TransactOpts, agent, token, amount)
}

// Distribute is a paid mutator transaction binding the contract method 0xa16a9183.
//
// Solidity: function distribute(address agent, address token, uint256 amount) returns()
func (_FeeRouter *FeeRouterTransactorSession) Distribute(agent common.Address, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeRouter.Contract.Distribute(&_FeeRouter.TransactOpts, agent, token, amount)
}

// DistributeNative is a paid mutator transaction binding the contract method 0x9d71020e.
//
// Solidity: function distributeNative(address agent) payable returns()
func (_FeeRouter *FeeRouterTransactor) DistributeNative(opts *bind.TransactOpts, agent common.Address) (*types.Transaction, error) {
	return _FeeRouter.contract.Transact(opts, "distributeNative", agent)
}

// DistributeNative is a paid mutator transaction binding the contract method 0x9d71020e.
//
// Solidity: function distributeNative(address agent) payable returns()
func (_FeeRouter *FeeRouterSession) DistributeNative(agent common.Address) (*types.Transaction, error) {
	return _FeeRouter.Contract.DistributeNative(&_FeeRouter.TransactOpts, agent)
}

// DistributeNative is a paid mutator transaction binding the contract method 0x9d71020e.
//
// Solidity: function distributeNative(address agent) payable returns()
func (_FeeRouter *FeeRouterTransactorSession) DistributeNative(agent common.Address) (*types.Transaction, error) {
	return _FeeRouter.Contract.DistributeNative(&_FeeRouter.TransactOpts, agent)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FeeRouter *FeeRouterTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeRouter.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FeeRouter *FeeRouterSession) RenounceOwnership() (*types.Transaction, error) {
	return _FeeRouter.Contract.RenounceOwnership(&_FeeRouter.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FeeRouter *FeeRouterTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FeeRouter.Contract.RenounceOwnership(&_FeeRouter.TransactOpts)
}

// SetRoute is a paid mutator transaction binding the contract method 0x9e66bb6f.
//
// Solidity: function setRoute(address agent, address creator, uint16 creatorBps) returns()
func (_FeeRouter *FeeRouterTransactor) SetRoute(opts *bind.TransactOpts, agent common.Address, creator common.Address, creatorBps uint16) (*types.Transaction, error) {
	return _FeeRouter.contract.Transact(opts, "setRoute", agent, creator, creatorBps)
}

// SetRoute is a paid mutator transaction binding the contract method 0x9e66bb6f.
//
// Solidity: function setRoute(address agent, address creator, uint16 creatorBps) returns()
func (_FeeRouter *FeeRouterSession) SetRoute(agent common.Address, creator common.Address, creatorBps uint16) (*types.Transaction, error) {
	return _FeeRouter.Contract.SetRoute(&_FeeRouter.TransactOpts, agent, creator, creatorBps)
}

// SetRoute is a paid mutator transaction binding the contract method 0x9e66bb6f.
//
// Solidity: function setRoute(address agent, address creator, uint16 creatorBps) returns()
func (_FeeRouter *FeeRouterTransactorSession) SetRoute(agent common.Address, creator common.Address, creatorBps uint16) (*types.Transaction, error) {
	return _FeeRouter.Contract.SetRoute(&_FeeRouter.TransactOpts, agent, creator, creatorBps)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FeeRouter *FeeRouterTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FeeRouter.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FeeRouter *FeeRouterSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FeeRouter.Contract.TransferOwnership(&_FeeRouter.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FeeRouter *FeeRouterTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FeeRouter.Contract.TransferOwnership(&_FeeRouter.TransactOpts, newOwner)
}

// FeeRouterDistributedIterator is returned from FilterDistributed and is used to iterate over the raw logs and unpacked data for Distributed events raised by the FeeRouter contract.
type FeeRouterDistributedIterator struct {
	Event *FeeRouterDistributed // Event containing the contract specifics and raw log

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
func (it *FeeRouterDistributedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeRouterDistributed)
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
		it.Event = new(FeeRouterDistributed)
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
func (it *FeeRouterDistributedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeRouterDistributedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeRouterDistributed represents a Distributed event raised by the FeeRouter contract.
type FeeRouterDistributed struct {
	Agent          common.Address
	Token          common.Address
	CreatorAmount  *big.Int
	TreasuryAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterDistributed is a free log retrieval operation binding the contract event 0xf57b0de383850a0698b76f6f85ebcc94687f139a1944025cd7c54954304c20d9.
//
// Solidity: event Distributed(address indexed agent, address indexed token, uint256 creatorAmount, uint256 treasuryAmount)
func (_FeeRouter *FeeRouterFilterer) FilterDistributed(opts *bind.FilterOpts, agent []common.Address, token []common.Address) (*FeeRouterDistributedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeRouter.contract.FilterLogs(opts, "Distributed", agentRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &FeeRouterDistributedIterator{contract: _FeeRouter.contract, event: "Distributed", logs: logs, sub: sub}, nil
}

// WatchDistributed is a free log subscription operation binding the contract event 0xf57b0de383850a0698b76f6f85ebcc94687f139a1944025cd7c54954304c20d9.
//
// Solidity: event Distributed(address indexed agent, address indexed token, uint256 creatorAmount, uint256 treasuryAmount)
func (_FeeRouter *FeeRouterFilterer) WatchDistributed(opts *bind.WatchOpts, sink chan<- *FeeRouterDistributed, agent []common.Address, token []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _FeeRouter.contract.WatchLogs(opts, "Distributed", agentRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeRouterDistributed)
				if err := _FeeRouter.contract.UnpackLog(event, "Distributed", log); err != nil {
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

// ParseDistributed is a log parse operation binding the contract event 0xf57b0de383850a0698b76f6f85ebcc94687f139a1944025cd7c54954304c20d9.
//
// Solidity: event Distributed(address indexed agent, address indexed token, uint256 creatorAmount, uint256 treasuryAmount)
func (_FeeRouter *FeeRouterFilterer) ParseDistributed(log types.Log) (*FeeRouterDistributed, error) {
	event := new(FeeRouterDistributed)
	if err := _FeeRouter.contract.UnpackLog(event, "Distributed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeRouterOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FeeRouter contract.
type FeeRouterOwnershipTransferredIterator struct {
	Event *FeeRouterOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FeeRouterOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeRouterOwnershipTransferred)
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
		it.Event = new(FeeRouterOwnershipTransferred)
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
func (it *FeeRouterOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeRouterOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeRouterOwnershipTransferred represents a OwnershipTransferred event raised by the FeeRouter contract.
type FeeRouterOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FeeRouter *FeeRouterFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FeeRouterOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FeeRouter.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FeeRouterOwnershipTransferredIterator{contract: _FeeRouter.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FeeRouter *FeeRouterFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeeRouterOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FeeRouter.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeRouterOwnershipTransferred)
				if err := _FeeRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_FeeRouter *FeeRouterFilterer) ParseOwnershipTransferred(log types.Log) (*FeeRouterOwnershipTransferred, error) {
	event := new(FeeRouterOwnershipTransferred)
	if err := _FeeRouter.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeRouterRouteClearedIterator is returned from FilterRouteCleared and is used to iterate over the raw logs and unpacked data for RouteCleared events raised by the FeeRouter contract.
type FeeRouterRouteClearedIterator struct {
	Event *FeeRouterRouteCleared // Event containing the contract specifics and raw log

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
func (it *FeeRouterRouteClearedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeRouterRouteCleared)
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
		it.Event = new(FeeRouterRouteCleared)
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
func (it *FeeRouterRouteClearedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeRouterRouteClearedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeRouterRouteCleared represents a RouteCleared event raised by the FeeRouter contract.
type FeeRouterRouteCleared struct {
	Agent common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRouteCleared is a free log retrieval operation binding the contract event 0xf13e05d9cb53ed68362bc7ee84fb0c6c651d6493c29a2a2847c4926aeed3258b.
//
// Solidity: event RouteCleared(address indexed agent)
func (_FeeRouter *FeeRouterFilterer) FilterRouteCleared(opts *bind.FilterOpts, agent []common.Address) (*FeeRouterRouteClearedIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _FeeRouter.contract.FilterLogs(opts, "RouteCleared", agentRule)
	if err != nil {
		return nil, err
	}
	return &FeeRouterRouteClearedIterator{contract: _FeeRouter.contract, event: "RouteCleared", logs: logs, sub: sub}, nil
}

// WatchRouteCleared is a free log subscription operation binding the contract event 0xf13e05d9cb53ed68362bc7ee84fb0c6c651d6493c29a2a2847c4926aeed3258b.
//
// Solidity: event RouteCleared(address indexed agent)
func (_FeeRouter *FeeRouterFilterer) WatchRouteCleared(opts *bind.WatchOpts, sink chan<- *FeeRouterRouteCleared, agent []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}

	logs, sub, err := _FeeRouter.contract.WatchLogs(opts, "RouteCleared", agentRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeRouterRouteCleared)
				if err := _FeeRouter.contract.UnpackLog(event, "RouteCleared", log); err != nil {
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

// ParseRouteCleared is a log parse operation binding the contract event 0xf13e05d9cb53ed68362bc7ee84fb0c6c651d6493c29a2a2847c4926aeed3258b.
//
// Solidity: event RouteCleared(address indexed agent)
func (_FeeRouter *FeeRouterFilterer) ParseRouteCleared(log types.Log) (*FeeRouterRouteCleared, error) {
	event := new(FeeRouterRouteCleared)
	if err := _FeeRouter.contract.UnpackLog(event, "RouteCleared", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeRouterRouteSetIterator is returned from FilterRouteSet and is used to iterate over the raw logs and unpacked data for RouteSet events raised by the FeeRouter contract.
type FeeRouterRouteSetIterator struct {
	Event *FeeRouterRouteSet // Event containing the contract specifics and raw log

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
func (it *FeeRouterRouteSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeRouterRouteSet)
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
		it.Event = new(FeeRouterRouteSet)
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
func (it *FeeRouterRouteSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeRouterRouteSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeRouterRouteSet represents a RouteSet event raised by the FeeRouter contract.
type FeeRouterRouteSet struct {
	Agent      common.Address
	Creator    common.Address
	CreatorBps uint16
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRouteSet is a free log retrieval operation binding the contract event 0x573e30a0b37c2f487eb2b7f0240ea1db1cfcc6b21d6734e1246e64d3fec1ade2.
//
// Solidity: event RouteSet(address indexed agent, address indexed creator, uint16 creatorBps)
func (_FeeRouter *FeeRouterFilterer) FilterRouteSet(opts *bind.FilterOpts, agent []common.Address, creator []common.Address) (*FeeRouterRouteSetIterator, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _FeeRouter.contract.FilterLogs(opts, "RouteSet", agentRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &FeeRouterRouteSetIterator{contract: _FeeRouter.contract, event: "RouteSet", logs: logs, sub: sub}, nil
}

// WatchRouteSet is a free log subscription operation binding the contract event 0x573e30a0b37c2f487eb2b7f0240ea1db1cfcc6b21d6734e1246e64d3fec1ade2.
//
// Solidity: event RouteSet(address indexed agent, address indexed creator, uint16 creatorBps)
func (_FeeRouter *FeeRouterFilterer) WatchRouteSet(opts *bind.WatchOpts, sink chan<- *FeeRouterRouteSet, agent []common.Address, creator []common.Address) (event.Subscription, error) {

	var agentRule []interface{}
	for _, agentItem := range agent {
		agentRule = append(agentRule, agentItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _FeeRouter.contract.WatchLogs(opts, "RouteSet", agentRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeRouterRouteSet)
				if err := _FeeRouter.contract.UnpackLog(event, "RouteSet", log); err != nil {
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

// ParseRouteSet is a log parse operation binding the contract event 0x573e30a0b37c2f487eb2b7f0240ea1db1cfcc6b21d6734e1246e64d3fec1ade2.
//
// Solidity: event RouteSet(address indexed agent, address indexed creator, uint16 creatorBps)
func (_FeeRouter *FeeRouterFilterer) ParseRouteSet(log types.Log) (*FeeRouterRouteSet, error) {
	event := new(FeeRouterRouteSet)
	if err := _FeeRouter.contract.UnpackLog(event, "RouteSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
