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

// BondingCurveMetaData contains all meta data concerning the BondingCurve contract.
var BondingCurveMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"quoteToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeRouter_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"virtualQuoteReserve_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"virtualAgentReserve_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"graduationThreshold_\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"initialAgentReserve_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BPS_DENOMINATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"FEE_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"agentToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buy\",\"inputs\":[{\"name\":\"minAgentOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"quoteIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"graduated\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"graduationThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"graduator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialAgentReserve\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"k\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pullForGraduation\",\"inputs\":[],\"outputs\":[{\"name\":\"quoteAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"agentAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"quoteBuy\",\"inputs\":[{\"name\":\"quoteIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"agentOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quoteSell\",\"inputs\":[{\"name\":\"agentIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"quoteOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quoteToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"realAgentReserve\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"realQuoteReserve\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestGraduation\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sell\",\"inputs\":[{\"name\":\"agentIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minQuoteOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFeeRouter\",\"inputs\":[{\"name\":\"next\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setGraduator\",\"inputs\":[{\"name\":\"next\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"virtualAgentReserve\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"virtualQuoteReserve\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Bought\",\"inputs\":[{\"name\":\"buyer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"quoteIn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"agentOut\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FeeRouterSet\",\"inputs\":[{\"name\":\"prev\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"next\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Graduated\",\"inputs\":[{\"name\":\"quoteReserve\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"agentReserve\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GraduatorSet\",\"inputs\":[{\"name\":\"prev\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"next\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Pulled\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"quoteAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"agentAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Sold\",\"inputs\":[{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"agentIn\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"quoteOut\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyGraduated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BelowThreshold\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExceedsGraduationThreshold\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GraduatorAlreadySet\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientOutput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidReserves\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGraduated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotGraduator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroValue\",\"inputs\":[]}]",
}

// BondingCurveABI is the input ABI used to generate the binding from.
// Deprecated: Use BondingCurveMetaData.ABI instead.
var BondingCurveABI = BondingCurveMetaData.ABI

// BondingCurve is an auto generated Go binding around an Ethereum contract.
type BondingCurve struct {
	BondingCurveCaller     // Read-only binding to the contract
	BondingCurveTransactor // Write-only binding to the contract
	BondingCurveFilterer   // Log filterer for contract events
}

// BondingCurveCaller is an auto generated read-only Go binding around an Ethereum contract.
type BondingCurveCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BondingCurveTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BondingCurveTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BondingCurveFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BondingCurveFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BondingCurveSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BondingCurveSession struct {
	Contract     *BondingCurve     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BondingCurveCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BondingCurveCallerSession struct {
	Contract *BondingCurveCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// BondingCurveTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BondingCurveTransactorSession struct {
	Contract     *BondingCurveTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// BondingCurveRaw is an auto generated low-level Go binding around an Ethereum contract.
type BondingCurveRaw struct {
	Contract *BondingCurve // Generic contract binding to access the raw methods on
}

// BondingCurveCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BondingCurveCallerRaw struct {
	Contract *BondingCurveCaller // Generic read-only contract binding to access the raw methods on
}

// BondingCurveTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BondingCurveTransactorRaw struct {
	Contract *BondingCurveTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBondingCurve creates a new instance of BondingCurve, bound to a specific deployed contract.
func NewBondingCurve(address common.Address, backend bind.ContractBackend) (*BondingCurve, error) {
	contract, err := bindBondingCurve(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BondingCurve{BondingCurveCaller: BondingCurveCaller{contract: contract}, BondingCurveTransactor: BondingCurveTransactor{contract: contract}, BondingCurveFilterer: BondingCurveFilterer{contract: contract}}, nil
}

// NewBondingCurveCaller creates a new read-only instance of BondingCurve, bound to a specific deployed contract.
func NewBondingCurveCaller(address common.Address, caller bind.ContractCaller) (*BondingCurveCaller, error) {
	contract, err := bindBondingCurve(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BondingCurveCaller{contract: contract}, nil
}

// NewBondingCurveTransactor creates a new write-only instance of BondingCurve, bound to a specific deployed contract.
func NewBondingCurveTransactor(address common.Address, transactor bind.ContractTransactor) (*BondingCurveTransactor, error) {
	contract, err := bindBondingCurve(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BondingCurveTransactor{contract: contract}, nil
}

// NewBondingCurveFilterer creates a new log filterer instance of BondingCurve, bound to a specific deployed contract.
func NewBondingCurveFilterer(address common.Address, filterer bind.ContractFilterer) (*BondingCurveFilterer, error) {
	contract, err := bindBondingCurve(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BondingCurveFilterer{contract: contract}, nil
}

// bindBondingCurve binds a generic wrapper to an already deployed contract.
func bindBondingCurve(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BondingCurveMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BondingCurve *BondingCurveRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BondingCurve.Contract.BondingCurveCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BondingCurve *BondingCurveRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BondingCurve.Contract.BondingCurveTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BondingCurve *BondingCurveRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BondingCurve.Contract.BondingCurveTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BondingCurve *BondingCurveCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BondingCurve.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BondingCurve *BondingCurveTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BondingCurve.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BondingCurve *BondingCurveTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BondingCurve.Contract.contract.Transact(opts, method, params...)
}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint16)
func (_BondingCurve *BondingCurveCaller) BPSDENOMINATOR(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "BPS_DENOMINATOR")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint16)
func (_BondingCurve *BondingCurveSession) BPSDENOMINATOR() (uint16, error) {
	return _BondingCurve.Contract.BPSDENOMINATOR(&_BondingCurve.CallOpts)
}

// BPSDENOMINATOR is a free data retrieval call binding the contract method 0xe1a45218.
//
// Solidity: function BPS_DENOMINATOR() view returns(uint16)
func (_BondingCurve *BondingCurveCallerSession) BPSDENOMINATOR() (uint16, error) {
	return _BondingCurve.Contract.BPSDENOMINATOR(&_BondingCurve.CallOpts)
}

// FEEBPS is a free data retrieval call binding the contract method 0xbf333f2c.
//
// Solidity: function FEE_BPS() view returns(uint16)
func (_BondingCurve *BondingCurveCaller) FEEBPS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "FEE_BPS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// FEEBPS is a free data retrieval call binding the contract method 0xbf333f2c.
//
// Solidity: function FEE_BPS() view returns(uint16)
func (_BondingCurve *BondingCurveSession) FEEBPS() (uint16, error) {
	return _BondingCurve.Contract.FEEBPS(&_BondingCurve.CallOpts)
}

// FEEBPS is a free data retrieval call binding the contract method 0xbf333f2c.
//
// Solidity: function FEE_BPS() view returns(uint16)
func (_BondingCurve *BondingCurveCallerSession) FEEBPS() (uint16, error) {
	return _BondingCurve.Contract.FEEBPS(&_BondingCurve.CallOpts)
}

// AgentToken is a free data retrieval call binding the contract method 0xb05707e9.
//
// Solidity: function agentToken() view returns(address)
func (_BondingCurve *BondingCurveCaller) AgentToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "agentToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentToken is a free data retrieval call binding the contract method 0xb05707e9.
//
// Solidity: function agentToken() view returns(address)
func (_BondingCurve *BondingCurveSession) AgentToken() (common.Address, error) {
	return _BondingCurve.Contract.AgentToken(&_BondingCurve.CallOpts)
}

// AgentToken is a free data retrieval call binding the contract method 0xb05707e9.
//
// Solidity: function agentToken() view returns(address)
func (_BondingCurve *BondingCurveCallerSession) AgentToken() (common.Address, error) {
	return _BondingCurve.Contract.AgentToken(&_BondingCurve.CallOpts)
}

// FeeRouter is a free data retrieval call binding the contract method 0xf29ebf61.
//
// Solidity: function feeRouter() view returns(address)
func (_BondingCurve *BondingCurveCaller) FeeRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "feeRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeRouter is a free data retrieval call binding the contract method 0xf29ebf61.
//
// Solidity: function feeRouter() view returns(address)
func (_BondingCurve *BondingCurveSession) FeeRouter() (common.Address, error) {
	return _BondingCurve.Contract.FeeRouter(&_BondingCurve.CallOpts)
}

// FeeRouter is a free data retrieval call binding the contract method 0xf29ebf61.
//
// Solidity: function feeRouter() view returns(address)
func (_BondingCurve *BondingCurveCallerSession) FeeRouter() (common.Address, error) {
	return _BondingCurve.Contract.FeeRouter(&_BondingCurve.CallOpts)
}

// Graduated is a free data retrieval call binding the contract method 0xe7c2b772.
//
// Solidity: function graduated() view returns(bool)
func (_BondingCurve *BondingCurveCaller) Graduated(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "graduated")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Graduated is a free data retrieval call binding the contract method 0xe7c2b772.
//
// Solidity: function graduated() view returns(bool)
func (_BondingCurve *BondingCurveSession) Graduated() (bool, error) {
	return _BondingCurve.Contract.Graduated(&_BondingCurve.CallOpts)
}

// Graduated is a free data retrieval call binding the contract method 0xe7c2b772.
//
// Solidity: function graduated() view returns(bool)
func (_BondingCurve *BondingCurveCallerSession) Graduated() (bool, error) {
	return _BondingCurve.Contract.Graduated(&_BondingCurve.CallOpts)
}

// GraduationThreshold is a free data retrieval call binding the contract method 0x8b0bc501.
//
// Solidity: function graduationThreshold() view returns(uint256)
func (_BondingCurve *BondingCurveCaller) GraduationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "graduationThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GraduationThreshold is a free data retrieval call binding the contract method 0x8b0bc501.
//
// Solidity: function graduationThreshold() view returns(uint256)
func (_BondingCurve *BondingCurveSession) GraduationThreshold() (*big.Int, error) {
	return _BondingCurve.Contract.GraduationThreshold(&_BondingCurve.CallOpts)
}

// GraduationThreshold is a free data retrieval call binding the contract method 0x8b0bc501.
//
// Solidity: function graduationThreshold() view returns(uint256)
func (_BondingCurve *BondingCurveCallerSession) GraduationThreshold() (*big.Int, error) {
	return _BondingCurve.Contract.GraduationThreshold(&_BondingCurve.CallOpts)
}

// Graduator is a free data retrieval call binding the contract method 0x0245a483.
//
// Solidity: function graduator() view returns(address)
func (_BondingCurve *BondingCurveCaller) Graduator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "graduator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Graduator is a free data retrieval call binding the contract method 0x0245a483.
//
// Solidity: function graduator() view returns(address)
func (_BondingCurve *BondingCurveSession) Graduator() (common.Address, error) {
	return _BondingCurve.Contract.Graduator(&_BondingCurve.CallOpts)
}

// Graduator is a free data retrieval call binding the contract method 0x0245a483.
//
// Solidity: function graduator() view returns(address)
func (_BondingCurve *BondingCurveCallerSession) Graduator() (common.Address, error) {
	return _BondingCurve.Contract.Graduator(&_BondingCurve.CallOpts)
}

// InitialAgentReserve is a free data retrieval call binding the contract method 0x18b722bb.
//
// Solidity: function initialAgentReserve() view returns(uint256)
func (_BondingCurve *BondingCurveCaller) InitialAgentReserve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "initialAgentReserve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialAgentReserve is a free data retrieval call binding the contract method 0x18b722bb.
//
// Solidity: function initialAgentReserve() view returns(uint256)
func (_BondingCurve *BondingCurveSession) InitialAgentReserve() (*big.Int, error) {
	return _BondingCurve.Contract.InitialAgentReserve(&_BondingCurve.CallOpts)
}

// InitialAgentReserve is a free data retrieval call binding the contract method 0x18b722bb.
//
// Solidity: function initialAgentReserve() view returns(uint256)
func (_BondingCurve *BondingCurveCallerSession) InitialAgentReserve() (*big.Int, error) {
	return _BondingCurve.Contract.InitialAgentReserve(&_BondingCurve.CallOpts)
}

// K is a free data retrieval call binding the contract method 0xb4f40c61.
//
// Solidity: function k() view returns(uint256)
func (_BondingCurve *BondingCurveCaller) K(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "k")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// K is a free data retrieval call binding the contract method 0xb4f40c61.
//
// Solidity: function k() view returns(uint256)
func (_BondingCurve *BondingCurveSession) K() (*big.Int, error) {
	return _BondingCurve.Contract.K(&_BondingCurve.CallOpts)
}

// K is a free data retrieval call binding the contract method 0xb4f40c61.
//
// Solidity: function k() view returns(uint256)
func (_BondingCurve *BondingCurveCallerSession) K() (*big.Int, error) {
	return _BondingCurve.Contract.K(&_BondingCurve.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BondingCurve *BondingCurveCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BondingCurve *BondingCurveSession) Owner() (common.Address, error) {
	return _BondingCurve.Contract.Owner(&_BondingCurve.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BondingCurve *BondingCurveCallerSession) Owner() (common.Address, error) {
	return _BondingCurve.Contract.Owner(&_BondingCurve.CallOpts)
}

// QuoteBuy is a free data retrieval call binding the contract method 0x4beb394c.
//
// Solidity: function quoteBuy(uint256 quoteIn) view returns(uint256 agentOut, uint256 fee)
func (_BondingCurve *BondingCurveCaller) QuoteBuy(opts *bind.CallOpts, quoteIn *big.Int) (struct {
	AgentOut *big.Int
	Fee      *big.Int
}, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "quoteBuy", quoteIn)

	outstruct := new(struct {
		AgentOut *big.Int
		Fee      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AgentOut = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// QuoteBuy is a free data retrieval call binding the contract method 0x4beb394c.
//
// Solidity: function quoteBuy(uint256 quoteIn) view returns(uint256 agentOut, uint256 fee)
func (_BondingCurve *BondingCurveSession) QuoteBuy(quoteIn *big.Int) (struct {
	AgentOut *big.Int
	Fee      *big.Int
}, error) {
	return _BondingCurve.Contract.QuoteBuy(&_BondingCurve.CallOpts, quoteIn)
}

// QuoteBuy is a free data retrieval call binding the contract method 0x4beb394c.
//
// Solidity: function quoteBuy(uint256 quoteIn) view returns(uint256 agentOut, uint256 fee)
func (_BondingCurve *BondingCurveCallerSession) QuoteBuy(quoteIn *big.Int) (struct {
	AgentOut *big.Int
	Fee      *big.Int
}, error) {
	return _BondingCurve.Contract.QuoteBuy(&_BondingCurve.CallOpts, quoteIn)
}

// QuoteSell is a free data retrieval call binding the contract method 0xa64190c4.
//
// Solidity: function quoteSell(uint256 agentIn) view returns(uint256 quoteOut, uint256 fee)
func (_BondingCurve *BondingCurveCaller) QuoteSell(opts *bind.CallOpts, agentIn *big.Int) (struct {
	QuoteOut *big.Int
	Fee      *big.Int
}, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "quoteSell", agentIn)

	outstruct := new(struct {
		QuoteOut *big.Int
		Fee      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.QuoteOut = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Fee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// QuoteSell is a free data retrieval call binding the contract method 0xa64190c4.
//
// Solidity: function quoteSell(uint256 agentIn) view returns(uint256 quoteOut, uint256 fee)
func (_BondingCurve *BondingCurveSession) QuoteSell(agentIn *big.Int) (struct {
	QuoteOut *big.Int
	Fee      *big.Int
}, error) {
	return _BondingCurve.Contract.QuoteSell(&_BondingCurve.CallOpts, agentIn)
}

// QuoteSell is a free data retrieval call binding the contract method 0xa64190c4.
//
// Solidity: function quoteSell(uint256 agentIn) view returns(uint256 quoteOut, uint256 fee)
func (_BondingCurve *BondingCurveCallerSession) QuoteSell(agentIn *big.Int) (struct {
	QuoteOut *big.Int
	Fee      *big.Int
}, error) {
	return _BondingCurve.Contract.QuoteSell(&_BondingCurve.CallOpts, agentIn)
}

// QuoteToken is a free data retrieval call binding the contract method 0x217a4b70.
//
// Solidity: function quoteToken() view returns(address)
func (_BondingCurve *BondingCurveCaller) QuoteToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "quoteToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// QuoteToken is a free data retrieval call binding the contract method 0x217a4b70.
//
// Solidity: function quoteToken() view returns(address)
func (_BondingCurve *BondingCurveSession) QuoteToken() (common.Address, error) {
	return _BondingCurve.Contract.QuoteToken(&_BondingCurve.CallOpts)
}

// QuoteToken is a free data retrieval call binding the contract method 0x217a4b70.
//
// Solidity: function quoteToken() view returns(address)
func (_BondingCurve *BondingCurveCallerSession) QuoteToken() (common.Address, error) {
	return _BondingCurve.Contract.QuoteToken(&_BondingCurve.CallOpts)
}

// RealAgentReserve is a free data retrieval call binding the contract method 0x272a273f.
//
// Solidity: function realAgentReserve() view returns(uint256)
func (_BondingCurve *BondingCurveCaller) RealAgentReserve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "realAgentReserve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RealAgentReserve is a free data retrieval call binding the contract method 0x272a273f.
//
// Solidity: function realAgentReserve() view returns(uint256)
func (_BondingCurve *BondingCurveSession) RealAgentReserve() (*big.Int, error) {
	return _BondingCurve.Contract.RealAgentReserve(&_BondingCurve.CallOpts)
}

// RealAgentReserve is a free data retrieval call binding the contract method 0x272a273f.
//
// Solidity: function realAgentReserve() view returns(uint256)
func (_BondingCurve *BondingCurveCallerSession) RealAgentReserve() (*big.Int, error) {
	return _BondingCurve.Contract.RealAgentReserve(&_BondingCurve.CallOpts)
}

// RealQuoteReserve is a free data retrieval call binding the contract method 0x4f1f58fd.
//
// Solidity: function realQuoteReserve() view returns(uint256)
func (_BondingCurve *BondingCurveCaller) RealQuoteReserve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "realQuoteReserve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RealQuoteReserve is a free data retrieval call binding the contract method 0x4f1f58fd.
//
// Solidity: function realQuoteReserve() view returns(uint256)
func (_BondingCurve *BondingCurveSession) RealQuoteReserve() (*big.Int, error) {
	return _BondingCurve.Contract.RealQuoteReserve(&_BondingCurve.CallOpts)
}

// RealQuoteReserve is a free data retrieval call binding the contract method 0x4f1f58fd.
//
// Solidity: function realQuoteReserve() view returns(uint256)
func (_BondingCurve *BondingCurveCallerSession) RealQuoteReserve() (*big.Int, error) {
	return _BondingCurve.Contract.RealQuoteReserve(&_BondingCurve.CallOpts)
}

// VirtualAgentReserve is a free data retrieval call binding the contract method 0x4b6db6cc.
//
// Solidity: function virtualAgentReserve() view returns(uint256)
func (_BondingCurve *BondingCurveCaller) VirtualAgentReserve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "virtualAgentReserve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VirtualAgentReserve is a free data retrieval call binding the contract method 0x4b6db6cc.
//
// Solidity: function virtualAgentReserve() view returns(uint256)
func (_BondingCurve *BondingCurveSession) VirtualAgentReserve() (*big.Int, error) {
	return _BondingCurve.Contract.VirtualAgentReserve(&_BondingCurve.CallOpts)
}

// VirtualAgentReserve is a free data retrieval call binding the contract method 0x4b6db6cc.
//
// Solidity: function virtualAgentReserve() view returns(uint256)
func (_BondingCurve *BondingCurveCallerSession) VirtualAgentReserve() (*big.Int, error) {
	return _BondingCurve.Contract.VirtualAgentReserve(&_BondingCurve.CallOpts)
}

// VirtualQuoteReserve is a free data retrieval call binding the contract method 0xf5596626.
//
// Solidity: function virtualQuoteReserve() view returns(uint256)
func (_BondingCurve *BondingCurveCaller) VirtualQuoteReserve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BondingCurve.contract.Call(opts, &out, "virtualQuoteReserve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VirtualQuoteReserve is a free data retrieval call binding the contract method 0xf5596626.
//
// Solidity: function virtualQuoteReserve() view returns(uint256)
func (_BondingCurve *BondingCurveSession) VirtualQuoteReserve() (*big.Int, error) {
	return _BondingCurve.Contract.VirtualQuoteReserve(&_BondingCurve.CallOpts)
}

// VirtualQuoteReserve is a free data retrieval call binding the contract method 0xf5596626.
//
// Solidity: function virtualQuoteReserve() view returns(uint256)
func (_BondingCurve *BondingCurveCallerSession) VirtualQuoteReserve() (*big.Int, error) {
	return _BondingCurve.Contract.VirtualQuoteReserve(&_BondingCurve.CallOpts)
}

// Buy is a paid mutator transaction binding the contract method 0xd6febde8.
//
// Solidity: function buy(uint256 minAgentOut, uint256 quoteIn) returns()
func (_BondingCurve *BondingCurveTransactor) Buy(opts *bind.TransactOpts, minAgentOut *big.Int, quoteIn *big.Int) (*types.Transaction, error) {
	return _BondingCurve.contract.Transact(opts, "buy", minAgentOut, quoteIn)
}

// Buy is a paid mutator transaction binding the contract method 0xd6febde8.
//
// Solidity: function buy(uint256 minAgentOut, uint256 quoteIn) returns()
func (_BondingCurve *BondingCurveSession) Buy(minAgentOut *big.Int, quoteIn *big.Int) (*types.Transaction, error) {
	return _BondingCurve.Contract.Buy(&_BondingCurve.TransactOpts, minAgentOut, quoteIn)
}

// Buy is a paid mutator transaction binding the contract method 0xd6febde8.
//
// Solidity: function buy(uint256 minAgentOut, uint256 quoteIn) returns()
func (_BondingCurve *BondingCurveTransactorSession) Buy(minAgentOut *big.Int, quoteIn *big.Int) (*types.Transaction, error) {
	return _BondingCurve.Contract.Buy(&_BondingCurve.TransactOpts, minAgentOut, quoteIn)
}

// PullForGraduation is a paid mutator transaction binding the contract method 0x61944511.
//
// Solidity: function pullForGraduation() returns(uint256 quoteAmount, uint256 agentAmount)
func (_BondingCurve *BondingCurveTransactor) PullForGraduation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BondingCurve.contract.Transact(opts, "pullForGraduation")
}

// PullForGraduation is a paid mutator transaction binding the contract method 0x61944511.
//
// Solidity: function pullForGraduation() returns(uint256 quoteAmount, uint256 agentAmount)
func (_BondingCurve *BondingCurveSession) PullForGraduation() (*types.Transaction, error) {
	return _BondingCurve.Contract.PullForGraduation(&_BondingCurve.TransactOpts)
}

// PullForGraduation is a paid mutator transaction binding the contract method 0x61944511.
//
// Solidity: function pullForGraduation() returns(uint256 quoteAmount, uint256 agentAmount)
func (_BondingCurve *BondingCurveTransactorSession) PullForGraduation() (*types.Transaction, error) {
	return _BondingCurve.Contract.PullForGraduation(&_BondingCurve.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BondingCurve *BondingCurveTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BondingCurve.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BondingCurve *BondingCurveSession) RenounceOwnership() (*types.Transaction, error) {
	return _BondingCurve.Contract.RenounceOwnership(&_BondingCurve.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BondingCurve *BondingCurveTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BondingCurve.Contract.RenounceOwnership(&_BondingCurve.TransactOpts)
}

// RequestGraduation is a paid mutator transaction binding the contract method 0xd82daa07.
//
// Solidity: function requestGraduation() returns()
func (_BondingCurve *BondingCurveTransactor) RequestGraduation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BondingCurve.contract.Transact(opts, "requestGraduation")
}

// RequestGraduation is a paid mutator transaction binding the contract method 0xd82daa07.
//
// Solidity: function requestGraduation() returns()
func (_BondingCurve *BondingCurveSession) RequestGraduation() (*types.Transaction, error) {
	return _BondingCurve.Contract.RequestGraduation(&_BondingCurve.TransactOpts)
}

// RequestGraduation is a paid mutator transaction binding the contract method 0xd82daa07.
//
// Solidity: function requestGraduation() returns()
func (_BondingCurve *BondingCurveTransactorSession) RequestGraduation() (*types.Transaction, error) {
	return _BondingCurve.Contract.RequestGraduation(&_BondingCurve.TransactOpts)
}

// Sell is a paid mutator transaction binding the contract method 0xd79875eb.
//
// Solidity: function sell(uint256 agentIn, uint256 minQuoteOut) returns()
func (_BondingCurve *BondingCurveTransactor) Sell(opts *bind.TransactOpts, agentIn *big.Int, minQuoteOut *big.Int) (*types.Transaction, error) {
	return _BondingCurve.contract.Transact(opts, "sell", agentIn, minQuoteOut)
}

// Sell is a paid mutator transaction binding the contract method 0xd79875eb.
//
// Solidity: function sell(uint256 agentIn, uint256 minQuoteOut) returns()
func (_BondingCurve *BondingCurveSession) Sell(agentIn *big.Int, minQuoteOut *big.Int) (*types.Transaction, error) {
	return _BondingCurve.Contract.Sell(&_BondingCurve.TransactOpts, agentIn, minQuoteOut)
}

// Sell is a paid mutator transaction binding the contract method 0xd79875eb.
//
// Solidity: function sell(uint256 agentIn, uint256 minQuoteOut) returns()
func (_BondingCurve *BondingCurveTransactorSession) Sell(agentIn *big.Int, minQuoteOut *big.Int) (*types.Transaction, error) {
	return _BondingCurve.Contract.Sell(&_BondingCurve.TransactOpts, agentIn, minQuoteOut)
}

// SetFeeRouter is a paid mutator transaction binding the contract method 0xc267ada4.
//
// Solidity: function setFeeRouter(address next) returns()
func (_BondingCurve *BondingCurveTransactor) SetFeeRouter(opts *bind.TransactOpts, next common.Address) (*types.Transaction, error) {
	return _BondingCurve.contract.Transact(opts, "setFeeRouter", next)
}

// SetFeeRouter is a paid mutator transaction binding the contract method 0xc267ada4.
//
// Solidity: function setFeeRouter(address next) returns()
func (_BondingCurve *BondingCurveSession) SetFeeRouter(next common.Address) (*types.Transaction, error) {
	return _BondingCurve.Contract.SetFeeRouter(&_BondingCurve.TransactOpts, next)
}

// SetFeeRouter is a paid mutator transaction binding the contract method 0xc267ada4.
//
// Solidity: function setFeeRouter(address next) returns()
func (_BondingCurve *BondingCurveTransactorSession) SetFeeRouter(next common.Address) (*types.Transaction, error) {
	return _BondingCurve.Contract.SetFeeRouter(&_BondingCurve.TransactOpts, next)
}

// SetGraduator is a paid mutator transaction binding the contract method 0xc169acde.
//
// Solidity: function setGraduator(address next) returns()
func (_BondingCurve *BondingCurveTransactor) SetGraduator(opts *bind.TransactOpts, next common.Address) (*types.Transaction, error) {
	return _BondingCurve.contract.Transact(opts, "setGraduator", next)
}

// SetGraduator is a paid mutator transaction binding the contract method 0xc169acde.
//
// Solidity: function setGraduator(address next) returns()
func (_BondingCurve *BondingCurveSession) SetGraduator(next common.Address) (*types.Transaction, error) {
	return _BondingCurve.Contract.SetGraduator(&_BondingCurve.TransactOpts, next)
}

// SetGraduator is a paid mutator transaction binding the contract method 0xc169acde.
//
// Solidity: function setGraduator(address next) returns()
func (_BondingCurve *BondingCurveTransactorSession) SetGraduator(next common.Address) (*types.Transaction, error) {
	return _BondingCurve.Contract.SetGraduator(&_BondingCurve.TransactOpts, next)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BondingCurve *BondingCurveTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BondingCurve.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BondingCurve *BondingCurveSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BondingCurve.Contract.TransferOwnership(&_BondingCurve.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BondingCurve *BondingCurveTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BondingCurve.Contract.TransferOwnership(&_BondingCurve.TransactOpts, newOwner)
}

// BondingCurveBoughtIterator is returned from FilterBought and is used to iterate over the raw logs and unpacked data for Bought events raised by the BondingCurve contract.
type BondingCurveBoughtIterator struct {
	Event *BondingCurveBought // Event containing the contract specifics and raw log

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
func (it *BondingCurveBoughtIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondingCurveBought)
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
		it.Event = new(BondingCurveBought)
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
func (it *BondingCurveBoughtIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondingCurveBoughtIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondingCurveBought represents a Bought event raised by the BondingCurve contract.
type BondingCurveBought struct {
	Buyer    common.Address
	QuoteIn  *big.Int
	AgentOut *big.Int
	Fee      *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBought is a free log retrieval operation binding the contract event 0xedba86fd2b22962d534e70ad9b0ff8730de46f636146f2bab6a72cbb1ebbcc53.
//
// Solidity: event Bought(address indexed buyer, uint256 quoteIn, uint256 agentOut, uint256 fee)
func (_BondingCurve *BondingCurveFilterer) FilterBought(opts *bind.FilterOpts, buyer []common.Address) (*BondingCurveBoughtIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _BondingCurve.contract.FilterLogs(opts, "Bought", buyerRule)
	if err != nil {
		return nil, err
	}
	return &BondingCurveBoughtIterator{contract: _BondingCurve.contract, event: "Bought", logs: logs, sub: sub}, nil
}

// WatchBought is a free log subscription operation binding the contract event 0xedba86fd2b22962d534e70ad9b0ff8730de46f636146f2bab6a72cbb1ebbcc53.
//
// Solidity: event Bought(address indexed buyer, uint256 quoteIn, uint256 agentOut, uint256 fee)
func (_BondingCurve *BondingCurveFilterer) WatchBought(opts *bind.WatchOpts, sink chan<- *BondingCurveBought, buyer []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _BondingCurve.contract.WatchLogs(opts, "Bought", buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondingCurveBought)
				if err := _BondingCurve.contract.UnpackLog(event, "Bought", log); err != nil {
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

// ParseBought is a log parse operation binding the contract event 0xedba86fd2b22962d534e70ad9b0ff8730de46f636146f2bab6a72cbb1ebbcc53.
//
// Solidity: event Bought(address indexed buyer, uint256 quoteIn, uint256 agentOut, uint256 fee)
func (_BondingCurve *BondingCurveFilterer) ParseBought(log types.Log) (*BondingCurveBought, error) {
	event := new(BondingCurveBought)
	if err := _BondingCurve.contract.UnpackLog(event, "Bought", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BondingCurveFeeRouterSetIterator is returned from FilterFeeRouterSet and is used to iterate over the raw logs and unpacked data for FeeRouterSet events raised by the BondingCurve contract.
type BondingCurveFeeRouterSetIterator struct {
	Event *BondingCurveFeeRouterSet // Event containing the contract specifics and raw log

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
func (it *BondingCurveFeeRouterSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondingCurveFeeRouterSet)
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
		it.Event = new(BondingCurveFeeRouterSet)
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
func (it *BondingCurveFeeRouterSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondingCurveFeeRouterSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondingCurveFeeRouterSet represents a FeeRouterSet event raised by the BondingCurve contract.
type BondingCurveFeeRouterSet struct {
	Prev common.Address
	Next common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterFeeRouterSet is a free log retrieval operation binding the contract event 0xe7f2ec24c1dc728e0b1df4cd4a80b136f2b519236115f081809e6b8b0becc0f6.
//
// Solidity: event FeeRouterSet(address indexed prev, address indexed next)
func (_BondingCurve *BondingCurveFilterer) FilterFeeRouterSet(opts *bind.FilterOpts, prev []common.Address, next []common.Address) (*BondingCurveFeeRouterSetIterator, error) {

	var prevRule []interface{}
	for _, prevItem := range prev {
		prevRule = append(prevRule, prevItem)
	}
	var nextRule []interface{}
	for _, nextItem := range next {
		nextRule = append(nextRule, nextItem)
	}

	logs, sub, err := _BondingCurve.contract.FilterLogs(opts, "FeeRouterSet", prevRule, nextRule)
	if err != nil {
		return nil, err
	}
	return &BondingCurveFeeRouterSetIterator{contract: _BondingCurve.contract, event: "FeeRouterSet", logs: logs, sub: sub}, nil
}

// WatchFeeRouterSet is a free log subscription operation binding the contract event 0xe7f2ec24c1dc728e0b1df4cd4a80b136f2b519236115f081809e6b8b0becc0f6.
//
// Solidity: event FeeRouterSet(address indexed prev, address indexed next)
func (_BondingCurve *BondingCurveFilterer) WatchFeeRouterSet(opts *bind.WatchOpts, sink chan<- *BondingCurveFeeRouterSet, prev []common.Address, next []common.Address) (event.Subscription, error) {

	var prevRule []interface{}
	for _, prevItem := range prev {
		prevRule = append(prevRule, prevItem)
	}
	var nextRule []interface{}
	for _, nextItem := range next {
		nextRule = append(nextRule, nextItem)
	}

	logs, sub, err := _BondingCurve.contract.WatchLogs(opts, "FeeRouterSet", prevRule, nextRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondingCurveFeeRouterSet)
				if err := _BondingCurve.contract.UnpackLog(event, "FeeRouterSet", log); err != nil {
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

// ParseFeeRouterSet is a log parse operation binding the contract event 0xe7f2ec24c1dc728e0b1df4cd4a80b136f2b519236115f081809e6b8b0becc0f6.
//
// Solidity: event FeeRouterSet(address indexed prev, address indexed next)
func (_BondingCurve *BondingCurveFilterer) ParseFeeRouterSet(log types.Log) (*BondingCurveFeeRouterSet, error) {
	event := new(BondingCurveFeeRouterSet)
	if err := _BondingCurve.contract.UnpackLog(event, "FeeRouterSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BondingCurveGraduatedIterator is returned from FilterGraduated and is used to iterate over the raw logs and unpacked data for Graduated events raised by the BondingCurve contract.
type BondingCurveGraduatedIterator struct {
	Event *BondingCurveGraduated // Event containing the contract specifics and raw log

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
func (it *BondingCurveGraduatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondingCurveGraduated)
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
		it.Event = new(BondingCurveGraduated)
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
func (it *BondingCurveGraduatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondingCurveGraduatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondingCurveGraduated represents a Graduated event raised by the BondingCurve contract.
type BondingCurveGraduated struct {
	QuoteReserve *big.Int
	AgentReserve *big.Int
	Timestamp    *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterGraduated is a free log retrieval operation binding the contract event 0x72a089bf72f8bdb633c01144c6cf486e8b100097b06bd326948141b7bd827d88.
//
// Solidity: event Graduated(uint256 quoteReserve, uint256 agentReserve, uint256 timestamp)
func (_BondingCurve *BondingCurveFilterer) FilterGraduated(opts *bind.FilterOpts) (*BondingCurveGraduatedIterator, error) {

	logs, sub, err := _BondingCurve.contract.FilterLogs(opts, "Graduated")
	if err != nil {
		return nil, err
	}
	return &BondingCurveGraduatedIterator{contract: _BondingCurve.contract, event: "Graduated", logs: logs, sub: sub}, nil
}

// WatchGraduated is a free log subscription operation binding the contract event 0x72a089bf72f8bdb633c01144c6cf486e8b100097b06bd326948141b7bd827d88.
//
// Solidity: event Graduated(uint256 quoteReserve, uint256 agentReserve, uint256 timestamp)
func (_BondingCurve *BondingCurveFilterer) WatchGraduated(opts *bind.WatchOpts, sink chan<- *BondingCurveGraduated) (event.Subscription, error) {

	logs, sub, err := _BondingCurve.contract.WatchLogs(opts, "Graduated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondingCurveGraduated)
				if err := _BondingCurve.contract.UnpackLog(event, "Graduated", log); err != nil {
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

// ParseGraduated is a log parse operation binding the contract event 0x72a089bf72f8bdb633c01144c6cf486e8b100097b06bd326948141b7bd827d88.
//
// Solidity: event Graduated(uint256 quoteReserve, uint256 agentReserve, uint256 timestamp)
func (_BondingCurve *BondingCurveFilterer) ParseGraduated(log types.Log) (*BondingCurveGraduated, error) {
	event := new(BondingCurveGraduated)
	if err := _BondingCurve.contract.UnpackLog(event, "Graduated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BondingCurveGraduatorSetIterator is returned from FilterGraduatorSet and is used to iterate over the raw logs and unpacked data for GraduatorSet events raised by the BondingCurve contract.
type BondingCurveGraduatorSetIterator struct {
	Event *BondingCurveGraduatorSet // Event containing the contract specifics and raw log

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
func (it *BondingCurveGraduatorSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondingCurveGraduatorSet)
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
		it.Event = new(BondingCurveGraduatorSet)
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
func (it *BondingCurveGraduatorSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondingCurveGraduatorSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondingCurveGraduatorSet represents a GraduatorSet event raised by the BondingCurve contract.
type BondingCurveGraduatorSet struct {
	Prev common.Address
	Next common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterGraduatorSet is a free log retrieval operation binding the contract event 0x7f72df928c4166945e33c99cb2d0977b39e009056cec0240694e7f65a0aedbe9.
//
// Solidity: event GraduatorSet(address indexed prev, address indexed next)
func (_BondingCurve *BondingCurveFilterer) FilterGraduatorSet(opts *bind.FilterOpts, prev []common.Address, next []common.Address) (*BondingCurveGraduatorSetIterator, error) {

	var prevRule []interface{}
	for _, prevItem := range prev {
		prevRule = append(prevRule, prevItem)
	}
	var nextRule []interface{}
	for _, nextItem := range next {
		nextRule = append(nextRule, nextItem)
	}

	logs, sub, err := _BondingCurve.contract.FilterLogs(opts, "GraduatorSet", prevRule, nextRule)
	if err != nil {
		return nil, err
	}
	return &BondingCurveGraduatorSetIterator{contract: _BondingCurve.contract, event: "GraduatorSet", logs: logs, sub: sub}, nil
}

// WatchGraduatorSet is a free log subscription operation binding the contract event 0x7f72df928c4166945e33c99cb2d0977b39e009056cec0240694e7f65a0aedbe9.
//
// Solidity: event GraduatorSet(address indexed prev, address indexed next)
func (_BondingCurve *BondingCurveFilterer) WatchGraduatorSet(opts *bind.WatchOpts, sink chan<- *BondingCurveGraduatorSet, prev []common.Address, next []common.Address) (event.Subscription, error) {

	var prevRule []interface{}
	for _, prevItem := range prev {
		prevRule = append(prevRule, prevItem)
	}
	var nextRule []interface{}
	for _, nextItem := range next {
		nextRule = append(nextRule, nextItem)
	}

	logs, sub, err := _BondingCurve.contract.WatchLogs(opts, "GraduatorSet", prevRule, nextRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondingCurveGraduatorSet)
				if err := _BondingCurve.contract.UnpackLog(event, "GraduatorSet", log); err != nil {
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

// ParseGraduatorSet is a log parse operation binding the contract event 0x7f72df928c4166945e33c99cb2d0977b39e009056cec0240694e7f65a0aedbe9.
//
// Solidity: event GraduatorSet(address indexed prev, address indexed next)
func (_BondingCurve *BondingCurveFilterer) ParseGraduatorSet(log types.Log) (*BondingCurveGraduatorSet, error) {
	event := new(BondingCurveGraduatorSet)
	if err := _BondingCurve.contract.UnpackLog(event, "GraduatorSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BondingCurveOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BondingCurve contract.
type BondingCurveOwnershipTransferredIterator struct {
	Event *BondingCurveOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BondingCurveOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondingCurveOwnershipTransferred)
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
		it.Event = new(BondingCurveOwnershipTransferred)
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
func (it *BondingCurveOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondingCurveOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondingCurveOwnershipTransferred represents a OwnershipTransferred event raised by the BondingCurve contract.
type BondingCurveOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BondingCurve *BondingCurveFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BondingCurveOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BondingCurve.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BondingCurveOwnershipTransferredIterator{contract: _BondingCurve.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BondingCurve *BondingCurveFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BondingCurveOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BondingCurve.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondingCurveOwnershipTransferred)
				if err := _BondingCurve.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BondingCurve *BondingCurveFilterer) ParseOwnershipTransferred(log types.Log) (*BondingCurveOwnershipTransferred, error) {
	event := new(BondingCurveOwnershipTransferred)
	if err := _BondingCurve.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BondingCurvePulledIterator is returned from FilterPulled and is used to iterate over the raw logs and unpacked data for Pulled events raised by the BondingCurve contract.
type BondingCurvePulledIterator struct {
	Event *BondingCurvePulled // Event containing the contract specifics and raw log

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
func (it *BondingCurvePulledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondingCurvePulled)
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
		it.Event = new(BondingCurvePulled)
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
func (it *BondingCurvePulledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondingCurvePulledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondingCurvePulled represents a Pulled event raised by the BondingCurve contract.
type BondingCurvePulled struct {
	To          common.Address
	QuoteAmount *big.Int
	AgentAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterPulled is a free log retrieval operation binding the contract event 0xbe8a6c37f88abb4f69e22e6409b95925635482005726063299ea716f1585b392.
//
// Solidity: event Pulled(address indexed to, uint256 quoteAmount, uint256 agentAmount)
func (_BondingCurve *BondingCurveFilterer) FilterPulled(opts *bind.FilterOpts, to []common.Address) (*BondingCurvePulledIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BondingCurve.contract.FilterLogs(opts, "Pulled", toRule)
	if err != nil {
		return nil, err
	}
	return &BondingCurvePulledIterator{contract: _BondingCurve.contract, event: "Pulled", logs: logs, sub: sub}, nil
}

// WatchPulled is a free log subscription operation binding the contract event 0xbe8a6c37f88abb4f69e22e6409b95925635482005726063299ea716f1585b392.
//
// Solidity: event Pulled(address indexed to, uint256 quoteAmount, uint256 agentAmount)
func (_BondingCurve *BondingCurveFilterer) WatchPulled(opts *bind.WatchOpts, sink chan<- *BondingCurvePulled, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _BondingCurve.contract.WatchLogs(opts, "Pulled", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondingCurvePulled)
				if err := _BondingCurve.contract.UnpackLog(event, "Pulled", log); err != nil {
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

// ParsePulled is a log parse operation binding the contract event 0xbe8a6c37f88abb4f69e22e6409b95925635482005726063299ea716f1585b392.
//
// Solidity: event Pulled(address indexed to, uint256 quoteAmount, uint256 agentAmount)
func (_BondingCurve *BondingCurveFilterer) ParsePulled(log types.Log) (*BondingCurvePulled, error) {
	event := new(BondingCurvePulled)
	if err := _BondingCurve.contract.UnpackLog(event, "Pulled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BondingCurveSoldIterator is returned from FilterSold and is used to iterate over the raw logs and unpacked data for Sold events raised by the BondingCurve contract.
type BondingCurveSoldIterator struct {
	Event *BondingCurveSold // Event containing the contract specifics and raw log

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
func (it *BondingCurveSoldIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BondingCurveSold)
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
		it.Event = new(BondingCurveSold)
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
func (it *BondingCurveSoldIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BondingCurveSoldIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BondingCurveSold represents a Sold event raised by the BondingCurve contract.
type BondingCurveSold struct {
	Seller   common.Address
	AgentIn  *big.Int
	QuoteOut *big.Int
	Fee      *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSold is a free log retrieval operation binding the contract event 0xe029f26dbcf8c42dd2f352c10214a7fc26773dc62482c6241334a0402ac09a80.
//
// Solidity: event Sold(address indexed seller, uint256 agentIn, uint256 quoteOut, uint256 fee)
func (_BondingCurve *BondingCurveFilterer) FilterSold(opts *bind.FilterOpts, seller []common.Address) (*BondingCurveSoldIterator, error) {

	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _BondingCurve.contract.FilterLogs(opts, "Sold", sellerRule)
	if err != nil {
		return nil, err
	}
	return &BondingCurveSoldIterator{contract: _BondingCurve.contract, event: "Sold", logs: logs, sub: sub}, nil
}

// WatchSold is a free log subscription operation binding the contract event 0xe029f26dbcf8c42dd2f352c10214a7fc26773dc62482c6241334a0402ac09a80.
//
// Solidity: event Sold(address indexed seller, uint256 agentIn, uint256 quoteOut, uint256 fee)
func (_BondingCurve *BondingCurveFilterer) WatchSold(opts *bind.WatchOpts, sink chan<- *BondingCurveSold, seller []common.Address) (event.Subscription, error) {

	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _BondingCurve.contract.WatchLogs(opts, "Sold", sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BondingCurveSold)
				if err := _BondingCurve.contract.UnpackLog(event, "Sold", log); err != nil {
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

// ParseSold is a log parse operation binding the contract event 0xe029f26dbcf8c42dd2f352c10214a7fc26773dc62482c6241334a0402ac09a80.
//
// Solidity: event Sold(address indexed seller, uint256 agentIn, uint256 quoteOut, uint256 fee)
func (_BondingCurve *BondingCurveFilterer) ParseSold(log types.Log) (*BondingCurveSold, error) {
	event := new(BondingCurveSold)
	if err := _BondingCurve.contract.UnpackLog(event, "Sold", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
