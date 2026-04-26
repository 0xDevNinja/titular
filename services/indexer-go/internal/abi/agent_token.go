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

// AgentTokenMetaData contains all meta data concerning the AgentToken contract.
var AgentTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BPS_DENOMINATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TOTAL_SUPPLY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TRADE_TAX_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bondingCurve\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"creator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"name_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"creator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeRouter_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bondingCurve_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AgentTokenInitialized\",\"inputs\":[{\"name\":\"creator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"feeRouter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"bondingCurve\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TaxCollected\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"taxAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ERC20InsufficientAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"allowance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InsufficientBalance\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC20InvalidSpender\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"EmptyMetadata\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// AgentTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use AgentTokenMetaData.ABI instead.
var AgentTokenABI = AgentTokenMetaData.ABI

// AgentToken is an auto generated Go binding around an Ethereum contract.
type AgentToken struct {
	AgentTokenCaller     // Read-only binding to the contract
	AgentTokenTransactor // Write-only binding to the contract
	AgentTokenFilterer   // Log filterer for contract events
}

// AgentTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type AgentTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AgentTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgentTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgentTokenSession struct {
	Contract     *AgentToken       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AgentTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgentTokenCallerSession struct {
	Contract *AgentTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AgentTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgentTokenTransactorSession struct {
	Contract     *AgentTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AgentTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type AgentTokenRaw struct {
	Contract *AgentToken // Generic contract binding to access the raw methods on
}

// AgentTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgentTokenCallerRaw struct {
	Contract *AgentTokenCaller // Generic read-only contract binding to access the raw methods on
}

// AgentTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgentTokenTransactorRaw struct {
	Contract *AgentTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAgentToken creates a new instance of AgentToken, bound to a specific deployed contract.
func NewAgentToken(address common.Address, backend bind.ContractBackend) (*AgentToken, error) {
	contract, err := bindAgentToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AgentToken{AgentTokenCaller: AgentTokenCaller{contract: contract}, AgentTokenTransactor: AgentTokenTransactor{contract: contract}, AgentTokenFilterer: AgentTokenFilterer{contract: contract}}, nil
}

// NewAgentTokenCaller creates a new read-only instance of AgentToken, bound to a specific deployed contract.
func NewAgentTokenCaller(address common.Address, caller bind.ContractCaller) (*AgentTokenCaller, error) {
	contract, err := bindAgentToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgentTokenCaller{contract: contract}, nil
}

// NewAgentTokenTransactor creates a new write-only instance of AgentToken, bound to a specific deployed contract.
func NewAgentTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*AgentTokenTransactor, error) {
	contract, err := bindAgentToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgentTokenTransactor{contract: contract}, nil
}

// NewAgentTokenFilterer creates a new log filterer instance of AgentToken, bound to a specific deployed contract.
func NewAgentTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*AgentTokenFilterer, error) {
	contract, err := bindAgentToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgentTokenFilterer{contract: contract}, nil
}

// bindAgentToken binds a generic wrapper to an already deployed contract.
func bindAgentToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AgentTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentToken *AgentTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentToken.Contract.AgentTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentToken *AgentTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentToken.Contract.AgentTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentToken *AgentTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentToken.Contract.AgentTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentToken *AgentTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentToken *AgentTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentToken *AgentTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentToken.Contract.contract.Transact(opts, method, params...)
}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint16)
func (_AgentToken *AgentTokenCaller) BPSDENOMINATOR(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "BPS_DENOMINATOR")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint16)
func (_AgentToken *AgentTokenSession) BPSDENOMINATOR() (uint16, error) {
	return _AgentToken.Contract.BPSDENOMINATOR(&_AgentToken.CallOpts)
}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint16)
func (_AgentToken *AgentTokenCallerSession) BPSDENOMINATOR() (uint16, error) {
	return _AgentToken.Contract.BPSDENOMINATOR(&_AgentToken.CallOpts)
}

// TOTALSUPPLY is a free data retrieval call binding the contract method 0x902d55a5.
//
// Solidity: function TOTAL_SUPPLY() view returns(uint256)
func (_AgentToken *AgentTokenCaller) TOTALSUPPLY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "TOTAL_SUPPLY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TOTALSUPPLY is a free data retrieval call binding the contract method 0x902d55a5.
//
// Solidity: function TOTAL_SUPPLY() view returns(uint256)
func (_AgentToken *AgentTokenSession) TOTALSUPPLY() (*big.Int, error) {
	return _AgentToken.Contract.TOTALSUPPLY(&_AgentToken.CallOpts)
}

// TOTALSUPPLY is a free data retrieval call binding the contract method 0x902d55a5.
//
// Solidity: function TOTAL_SUPPLY() view returns(uint256)
func (_AgentToken *AgentTokenCallerSession) TOTALSUPPLY() (*big.Int, error) {
	return _AgentToken.Contract.TOTALSUPPLY(&_AgentToken.CallOpts)
}

// TRADETAXBPS is a free data retrieval call binding the contract method 0x2185ec77.
//
// Solidity: function TRADE_TAX_BPS() view returns(uint16)
func (_AgentToken *AgentTokenCaller) TRADETAXBPS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "TRADE_TAX_BPS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// TRADETAXBPS is a free data retrieval call binding the contract method 0x2185ec77.
//
// Solidity: function TRADE_TAX_BPS() view returns(uint16)
func (_AgentToken *AgentTokenSession) TRADETAXBPS() (uint16, error) {
	return _AgentToken.Contract.TRADETAXBPS(&_AgentToken.CallOpts)
}

// TRADETAXBPS is a free data retrieval call binding the contract method 0x2185ec77.
//
// Solidity: function TRADE_TAX_BPS() view returns(uint16)
func (_AgentToken *AgentTokenCallerSession) TRADETAXBPS() (uint16, error) {
	return _AgentToken.Contract.TRADETAXBPS(&_AgentToken.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AgentToken *AgentTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AgentToken *AgentTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _AgentToken.Contract.Allowance(&_AgentToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_AgentToken *AgentTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _AgentToken.Contract.Allowance(&_AgentToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AgentToken *AgentTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AgentToken *AgentTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _AgentToken.Contract.BalanceOf(&_AgentToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_AgentToken *AgentTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _AgentToken.Contract.BalanceOf(&_AgentToken.CallOpts, account)
}

// BondingCurve is a free data retrieval call binding the contract method 0xeff1d50e.
//
// Solidity: function bondingCurve() view returns(address)
func (_AgentToken *AgentTokenCaller) BondingCurve(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "bondingCurve")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BondingCurve is a free data retrieval call binding the contract method 0xeff1d50e.
//
// Solidity: function bondingCurve() view returns(address)
func (_AgentToken *AgentTokenSession) BondingCurve() (common.Address, error) {
	return _AgentToken.Contract.BondingCurve(&_AgentToken.CallOpts)
}

// BondingCurve is a free data retrieval call binding the contract method 0xeff1d50e.
//
// Solidity: function bondingCurve() view returns(address)
func (_AgentToken *AgentTokenCallerSession) BondingCurve() (common.Address, error) {
	return _AgentToken.Contract.BondingCurve(&_AgentToken.CallOpts)
}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_AgentToken *AgentTokenCaller) Creator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "creator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_AgentToken *AgentTokenSession) Creator() (common.Address, error) {
	return _AgentToken.Contract.Creator(&_AgentToken.CallOpts)
}

// Creator is a free data retrieval call binding the contract method 0x02d05d3f.
//
// Solidity: function creator() view returns(address)
func (_AgentToken *AgentTokenCallerSession) Creator() (common.Address, error) {
	return _AgentToken.Contract.Creator(&_AgentToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AgentToken *AgentTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AgentToken *AgentTokenSession) Decimals() (uint8, error) {
	return _AgentToken.Contract.Decimals(&_AgentToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_AgentToken *AgentTokenCallerSession) Decimals() (uint8, error) {
	return _AgentToken.Contract.Decimals(&_AgentToken.CallOpts)
}

// FeeRouter is a free data retrieval call binding the contract method 0xf29ebf61.
//
// Solidity: function feeRouter() view returns(address)
func (_AgentToken *AgentTokenCaller) FeeRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "feeRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeRouter is a free data retrieval call binding the contract method 0xf29ebf61.
//
// Solidity: function feeRouter() view returns(address)
func (_AgentToken *AgentTokenSession) FeeRouter() (common.Address, error) {
	return _AgentToken.Contract.FeeRouter(&_AgentToken.CallOpts)
}

// FeeRouter is a free data retrieval call binding the contract method 0xf29ebf61.
//
// Solidity: function feeRouter() view returns(address)
func (_AgentToken *AgentTokenCallerSession) FeeRouter() (common.Address, error) {
	return _AgentToken.Contract.FeeRouter(&_AgentToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AgentToken *AgentTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AgentToken *AgentTokenSession) Name() (string, error) {
	return _AgentToken.Contract.Name(&_AgentToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_AgentToken *AgentTokenCallerSession) Name() (string, error) {
	return _AgentToken.Contract.Name(&_AgentToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentToken *AgentTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentToken *AgentTokenSession) Owner() (common.Address, error) {
	return _AgentToken.Contract.Owner(&_AgentToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentToken *AgentTokenCallerSession) Owner() (common.Address, error) {
	return _AgentToken.Contract.Owner(&_AgentToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AgentToken *AgentTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AgentToken *AgentTokenSession) Symbol() (string, error) {
	return _AgentToken.Contract.Symbol(&_AgentToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_AgentToken *AgentTokenCallerSession) Symbol() (string, error) {
	return _AgentToken.Contract.Symbol(&_AgentToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AgentToken *AgentTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AgentToken *AgentTokenSession) TotalSupply() (*big.Int, error) {
	return _AgentToken.Contract.TotalSupply(&_AgentToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_AgentToken *AgentTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _AgentToken.Contract.TotalSupply(&_AgentToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_AgentToken *AgentTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_AgentToken *AgentTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentToken.Contract.Approve(&_AgentToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_AgentToken *AgentTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentToken.Contract.Approve(&_AgentToken.TransactOpts, spender, value)
}

// Initialize is a paid mutator transaction binding the contract method 0xdb0ed6a0.
//
// Solidity: function initialize(string name_, string symbol_, address creator_, address feeRouter_, address bondingCurve_) returns()
func (_AgentToken *AgentTokenTransactor) Initialize(opts *bind.TransactOpts, name_ string, symbol_ string, creator_ common.Address, feeRouter_ common.Address, bondingCurve_ common.Address) (*types.Transaction, error) {
	return _AgentToken.contract.Transact(opts, "initialize", name_, symbol_, creator_, feeRouter_, bondingCurve_)
}

// Initialize is a paid mutator transaction binding the contract method 0xdb0ed6a0.
//
// Solidity: function initialize(string name_, string symbol_, address creator_, address feeRouter_, address bondingCurve_) returns()
func (_AgentToken *AgentTokenSession) Initialize(name_ string, symbol_ string, creator_ common.Address, feeRouter_ common.Address, bondingCurve_ common.Address) (*types.Transaction, error) {
	return _AgentToken.Contract.Initialize(&_AgentToken.TransactOpts, name_, symbol_, creator_, feeRouter_, bondingCurve_)
}

// Initialize is a paid mutator transaction binding the contract method 0xdb0ed6a0.
//
// Solidity: function initialize(string name_, string symbol_, address creator_, address feeRouter_, address bondingCurve_) returns()
func (_AgentToken *AgentTokenTransactorSession) Initialize(name_ string, symbol_ string, creator_ common.Address, feeRouter_ common.Address, bondingCurve_ common.Address) (*types.Transaction, error) {
	return _AgentToken.Contract.Initialize(&_AgentToken.TransactOpts, name_, symbol_, creator_, feeRouter_, bondingCurve_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AgentToken *AgentTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AgentToken *AgentTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _AgentToken.Contract.RenounceOwnership(&_AgentToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AgentToken *AgentTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AgentToken.Contract.RenounceOwnership(&_AgentToken.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_AgentToken *AgentTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_AgentToken *AgentTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentToken.Contract.Transfer(&_AgentToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_AgentToken *AgentTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentToken.Contract.Transfer(&_AgentToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_AgentToken *AgentTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_AgentToken *AgentTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentToken.Contract.TransferFrom(&_AgentToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_AgentToken *AgentTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _AgentToken.Contract.TransferFrom(&_AgentToken.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AgentToken *AgentTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AgentToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AgentToken *AgentTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AgentToken.Contract.TransferOwnership(&_AgentToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AgentToken *AgentTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AgentToken.Contract.TransferOwnership(&_AgentToken.TransactOpts, newOwner)
}

// AgentTokenAgentTokenInitializedIterator is returned from FilterAgentTokenInitialized and is used to iterate over the raw logs and unpacked data for AgentTokenInitialized events raised by the AgentToken contract.
type AgentTokenAgentTokenInitializedIterator struct {
	Event *AgentTokenAgentTokenInitialized // Event containing the contract specifics and raw log

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
func (it *AgentTokenAgentTokenInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentTokenAgentTokenInitialized)
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
		it.Event = new(AgentTokenAgentTokenInitialized)
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
func (it *AgentTokenAgentTokenInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentTokenAgentTokenInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentTokenAgentTokenInitialized represents a AgentTokenInitialized event raised by the AgentToken contract.
type AgentTokenAgentTokenInitialized struct {
	Creator      common.Address
	FeeRouter    common.Address
	BondingCurve common.Address
	Name         string
	Symbol       string
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAgentTokenInitialized is a free log retrieval operation binding the contract event 0x31f205857a54813b55cd88e4a445625e2b138a25179d428f10a260e4dff33c74.
//
// Solidity: event AgentTokenInitialized(address indexed creator, address indexed feeRouter, address indexed bondingCurve, string name, string symbol)
func (_AgentToken *AgentTokenFilterer) FilterAgentTokenInitialized(opts *bind.FilterOpts, creator []common.Address, feeRouter []common.Address, bondingCurve []common.Address) (*AgentTokenAgentTokenInitializedIterator, error) {

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var feeRouterRule []interface{}
	for _, feeRouterItem := range feeRouter {
		feeRouterRule = append(feeRouterRule, feeRouterItem)
	}
	var bondingCurveRule []interface{}
	for _, bondingCurveItem := range bondingCurve {
		bondingCurveRule = append(bondingCurveRule, bondingCurveItem)
	}

	logs, sub, err := _AgentToken.contract.FilterLogs(opts, "AgentTokenInitialized", creatorRule, feeRouterRule, bondingCurveRule)
	if err != nil {
		return nil, err
	}
	return &AgentTokenAgentTokenInitializedIterator{contract: _AgentToken.contract, event: "AgentTokenInitialized", logs: logs, sub: sub}, nil
}

// WatchAgentTokenInitialized is a free log subscription operation binding the contract event 0x31f205857a54813b55cd88e4a445625e2b138a25179d428f10a260e4dff33c74.
//
// Solidity: event AgentTokenInitialized(address indexed creator, address indexed feeRouter, address indexed bondingCurve, string name, string symbol)
func (_AgentToken *AgentTokenFilterer) WatchAgentTokenInitialized(opts *bind.WatchOpts, sink chan<- *AgentTokenAgentTokenInitialized, creator []common.Address, feeRouter []common.Address, bondingCurve []common.Address) (event.Subscription, error) {

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var feeRouterRule []interface{}
	for _, feeRouterItem := range feeRouter {
		feeRouterRule = append(feeRouterRule, feeRouterItem)
	}
	var bondingCurveRule []interface{}
	for _, bondingCurveItem := range bondingCurve {
		bondingCurveRule = append(bondingCurveRule, bondingCurveItem)
	}

	logs, sub, err := _AgentToken.contract.WatchLogs(opts, "AgentTokenInitialized", creatorRule, feeRouterRule, bondingCurveRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentTokenAgentTokenInitialized)
				if err := _AgentToken.contract.UnpackLog(event, "AgentTokenInitialized", log); err != nil {
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

// ParseAgentTokenInitialized is a log parse operation binding the contract event 0x31f205857a54813b55cd88e4a445625e2b138a25179d428f10a260e4dff33c74.
//
// Solidity: event AgentTokenInitialized(address indexed creator, address indexed feeRouter, address indexed bondingCurve, string name, string symbol)
func (_AgentToken *AgentTokenFilterer) ParseAgentTokenInitialized(log types.Log) (*AgentTokenAgentTokenInitialized, error) {
	event := new(AgentTokenAgentTokenInitialized)
	if err := _AgentToken.contract.UnpackLog(event, "AgentTokenInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the AgentToken contract.
type AgentTokenApprovalIterator struct {
	Event *AgentTokenApproval // Event containing the contract specifics and raw log

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
func (it *AgentTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentTokenApproval)
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
		it.Event = new(AgentTokenApproval)
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
func (it *AgentTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentTokenApproval represents a Approval event raised by the AgentToken contract.
type AgentTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AgentToken *AgentTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*AgentTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _AgentToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &AgentTokenApprovalIterator{contract: _AgentToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AgentToken *AgentTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *AgentTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _AgentToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentTokenApproval)
				if err := _AgentToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_AgentToken *AgentTokenFilterer) ParseApproval(log types.Log) (*AgentTokenApproval, error) {
	event := new(AgentTokenApproval)
	if err := _AgentToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentTokenInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AgentToken contract.
type AgentTokenInitializedIterator struct {
	Event *AgentTokenInitialized // Event containing the contract specifics and raw log

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
func (it *AgentTokenInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentTokenInitialized)
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
		it.Event = new(AgentTokenInitialized)
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
func (it *AgentTokenInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentTokenInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentTokenInitialized represents a Initialized event raised by the AgentToken contract.
type AgentTokenInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AgentToken *AgentTokenFilterer) FilterInitialized(opts *bind.FilterOpts) (*AgentTokenInitializedIterator, error) {

	logs, sub, err := _AgentToken.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AgentTokenInitializedIterator{contract: _AgentToken.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AgentToken *AgentTokenFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AgentTokenInitialized) (event.Subscription, error) {

	logs, sub, err := _AgentToken.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentTokenInitialized)
				if err := _AgentToken.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_AgentToken *AgentTokenFilterer) ParseInitialized(log types.Log) (*AgentTokenInitialized, error) {
	event := new(AgentTokenInitialized)
	if err := _AgentToken.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AgentToken contract.
type AgentTokenOwnershipTransferredIterator struct {
	Event *AgentTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *AgentTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentTokenOwnershipTransferred)
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
		it.Event = new(AgentTokenOwnershipTransferred)
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
func (it *AgentTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentTokenOwnershipTransferred represents a OwnershipTransferred event raised by the AgentToken contract.
type AgentTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AgentToken *AgentTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AgentTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AgentToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AgentTokenOwnershipTransferredIterator{contract: _AgentToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AgentToken *AgentTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AgentTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AgentToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentTokenOwnershipTransferred)
				if err := _AgentToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_AgentToken *AgentTokenFilterer) ParseOwnershipTransferred(log types.Log) (*AgentTokenOwnershipTransferred, error) {
	event := new(AgentTokenOwnershipTransferred)
	if err := _AgentToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentTokenTaxCollectedIterator is returned from FilterTaxCollected and is used to iterate over the raw logs and unpacked data for TaxCollected events raised by the AgentToken contract.
type AgentTokenTaxCollectedIterator struct {
	Event *AgentTokenTaxCollected // Event containing the contract specifics and raw log

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
func (it *AgentTokenTaxCollectedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentTokenTaxCollected)
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
		it.Event = new(AgentTokenTaxCollected)
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
func (it *AgentTokenTaxCollectedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentTokenTaxCollectedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentTokenTaxCollected represents a TaxCollected event raised by the AgentToken contract.
type AgentTokenTaxCollected struct {
	From      common.Address
	To        common.Address
	TaxAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTaxCollected is a free log retrieval operation binding the contract event 0x5d37fd68fe66745a199f8c603e00ae02183f4aabb8ec0089589b0b40c4ead5e1.
//
// Solidity: event TaxCollected(address indexed from, address indexed to, uint256 taxAmount)
func (_AgentToken *AgentTokenFilterer) FilterTaxCollected(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AgentTokenTaxCollectedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AgentToken.contract.FilterLogs(opts, "TaxCollected", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AgentTokenTaxCollectedIterator{contract: _AgentToken.contract, event: "TaxCollected", logs: logs, sub: sub}, nil
}

// WatchTaxCollected is a free log subscription operation binding the contract event 0x5d37fd68fe66745a199f8c603e00ae02183f4aabb8ec0089589b0b40c4ead5e1.
//
// Solidity: event TaxCollected(address indexed from, address indexed to, uint256 taxAmount)
func (_AgentToken *AgentTokenFilterer) WatchTaxCollected(opts *bind.WatchOpts, sink chan<- *AgentTokenTaxCollected, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AgentToken.contract.WatchLogs(opts, "TaxCollected", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentTokenTaxCollected)
				if err := _AgentToken.contract.UnpackLog(event, "TaxCollected", log); err != nil {
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

// ParseTaxCollected is a log parse operation binding the contract event 0x5d37fd68fe66745a199f8c603e00ae02183f4aabb8ec0089589b0b40c4ead5e1.
//
// Solidity: event TaxCollected(address indexed from, address indexed to, uint256 taxAmount)
func (_AgentToken *AgentTokenFilterer) ParseTaxCollected(log types.Log) (*AgentTokenTaxCollected, error) {
	event := new(AgentTokenTaxCollected)
	if err := _AgentToken.contract.UnpackLog(event, "TaxCollected", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the AgentToken contract.
type AgentTokenTransferIterator struct {
	Event *AgentTokenTransfer // Event containing the contract specifics and raw log

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
func (it *AgentTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentTokenTransfer)
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
		it.Event = new(AgentTokenTransfer)
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
func (it *AgentTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentTokenTransfer represents a Transfer event raised by the AgentToken contract.
type AgentTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AgentToken *AgentTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*AgentTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AgentToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &AgentTokenTransferIterator{contract: _AgentToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AgentToken *AgentTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AgentTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _AgentToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentTokenTransfer)
				if err := _AgentToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_AgentToken *AgentTokenFilterer) ParseTransfer(log types.Log) (*AgentTokenTransfer, error) {
	event := new(AgentTokenTransfer)
	if err := _AgentToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
