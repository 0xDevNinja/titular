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

// LaunchpadFactoryAgent is an auto generated low-level Go binding around an user-defined struct.
type LaunchpadFactoryAgent struct {
	Token     common.Address
	Curve     common.Address
	LpLock    common.Address
	Pair      common.Address
	Creator   common.Address
	CreatedAt uint64
	Modules   [][32]byte
}

// LaunchpadFactoryLaunchParams is an auto generated low-level Go binding around an user-defined struct.
type LaunchpadFactoryLaunchParams struct {
	Name           string
	Symbol         string
	ImageURI       string
	SoulURI        string
	EnabledModules [][32]byte
	ModuleData     [][]byte
}

// LaunchpadFactoryMetaData contains all meta data concerning the LaunchpadFactory contract.
var LaunchpadFactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"tituToken_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"treasury_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentTokenImpl_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bondingCurveImpl_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"feeRouter_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"graduator_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"uniV2Factory_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"uniV2Router_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"owner_\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"GRADUATION_THRESHOLD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"INITIAL_AGENT_RESERVE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_NAME_LEN\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_SYMBOL_LEN\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MODULE_ID_AIRDROP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MODULE_ID_ANTI_SNIPER\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MODULE_ID_CAPITAL_FORMATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MODULE_ID_EXISTING_TOKEN\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MODULE_ID_LAUNCH_RADAR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MODULE_ID_PRE_BUY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MODULE_ID_SIXTY_DAYS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VIRTUAL_AGENT_RESERVE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VIRTUAL_QUOTE_RESERVE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"agentCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"agentTokenImpl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"bondingCurveImpl\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAgent\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structLaunchpadFactory.Agent\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"curve\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"lpLock\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"pair\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"createdAt\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"modules\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"graduator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractGraduator\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"launchAgent\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structLaunchpadFactory.LaunchParams\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"imageURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"soulURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"enabledModules\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"},{\"name\":\"moduleData\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}]}],\"outputs\":[{\"name\":\"agentToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"bondingCurve\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"modules\",\"inputs\":[{\"name\":\"moduleId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"module\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setModule\",\"inputs\":[{\"name\":\"moduleId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"module\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tituToken\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"treasury\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"uniV2Factory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUniswapV2Factory\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"uniV2Router\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIUniswapV2Router02\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AgentLaunched\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"curve\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"creator\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"lpLock\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"pair\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"modules\",\"type\":\"bytes32[]\",\"indexed\":false,\"internalType\":\"bytes32[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ModuleSet\",\"inputs\":[{\"name\":\"moduleId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"prev\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"module\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"DuplicateModule\",\"inputs\":[{\"name\":\"moduleId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"EmptyName\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptySymbol\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedDeployment\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[{\"name\":\"balance\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"needed\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidModuleId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LengthMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NameTooLong\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SymbolTooLong\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnknownAgent\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnknownModule\",\"inputs\":[{\"name\":\"moduleId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// LaunchpadFactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use LaunchpadFactoryMetaData.ABI instead.
var LaunchpadFactoryABI = LaunchpadFactoryMetaData.ABI

// LaunchpadFactory is an auto generated Go binding around an Ethereum contract.
type LaunchpadFactory struct {
	LaunchpadFactoryCaller     // Read-only binding to the contract
	LaunchpadFactoryTransactor // Write-only binding to the contract
	LaunchpadFactoryFilterer   // Log filterer for contract events
}

// LaunchpadFactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type LaunchpadFactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LaunchpadFactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LaunchpadFactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LaunchpadFactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LaunchpadFactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LaunchpadFactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LaunchpadFactorySession struct {
	Contract     *LaunchpadFactory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LaunchpadFactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LaunchpadFactoryCallerSession struct {
	Contract *LaunchpadFactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// LaunchpadFactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LaunchpadFactoryTransactorSession struct {
	Contract     *LaunchpadFactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// LaunchpadFactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type LaunchpadFactoryRaw struct {
	Contract *LaunchpadFactory // Generic contract binding to access the raw methods on
}

// LaunchpadFactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LaunchpadFactoryCallerRaw struct {
	Contract *LaunchpadFactoryCaller // Generic read-only contract binding to access the raw methods on
}

// LaunchpadFactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LaunchpadFactoryTransactorRaw struct {
	Contract *LaunchpadFactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLaunchpadFactory creates a new instance of LaunchpadFactory, bound to a specific deployed contract.
func NewLaunchpadFactory(address common.Address, backend bind.ContractBackend) (*LaunchpadFactory, error) {
	contract, err := bindLaunchpadFactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LaunchpadFactory{LaunchpadFactoryCaller: LaunchpadFactoryCaller{contract: contract}, LaunchpadFactoryTransactor: LaunchpadFactoryTransactor{contract: contract}, LaunchpadFactoryFilterer: LaunchpadFactoryFilterer{contract: contract}}, nil
}

// NewLaunchpadFactoryCaller creates a new read-only instance of LaunchpadFactory, bound to a specific deployed contract.
func NewLaunchpadFactoryCaller(address common.Address, caller bind.ContractCaller) (*LaunchpadFactoryCaller, error) {
	contract, err := bindLaunchpadFactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LaunchpadFactoryCaller{contract: contract}, nil
}

// NewLaunchpadFactoryTransactor creates a new write-only instance of LaunchpadFactory, bound to a specific deployed contract.
func NewLaunchpadFactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*LaunchpadFactoryTransactor, error) {
	contract, err := bindLaunchpadFactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LaunchpadFactoryTransactor{contract: contract}, nil
}

// NewLaunchpadFactoryFilterer creates a new log filterer instance of LaunchpadFactory, bound to a specific deployed contract.
func NewLaunchpadFactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*LaunchpadFactoryFilterer, error) {
	contract, err := bindLaunchpadFactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LaunchpadFactoryFilterer{contract: contract}, nil
}

// bindLaunchpadFactory binds a generic wrapper to an already deployed contract.
func bindLaunchpadFactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LaunchpadFactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LaunchpadFactory *LaunchpadFactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LaunchpadFactory.Contract.LaunchpadFactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LaunchpadFactory *LaunchpadFactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.LaunchpadFactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LaunchpadFactory *LaunchpadFactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.LaunchpadFactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LaunchpadFactory *LaunchpadFactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LaunchpadFactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LaunchpadFactory *LaunchpadFactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LaunchpadFactory *LaunchpadFactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.contract.Transact(opts, method, params...)
}

// GRADUATIONTHRESHOLD is a free data retrieval call binding the contract method 0xfcfc0c09.
//
// Solidity: function GRADUATION_THRESHOLD() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCaller) GRADUATIONTHRESHOLD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "GRADUATION_THRESHOLD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GRADUATIONTHRESHOLD is a free data retrieval call binding the contract method 0xfcfc0c09.
//
// Solidity: function GRADUATION_THRESHOLD() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactorySession) GRADUATIONTHRESHOLD() (*big.Int, error) {
	return _LaunchpadFactory.Contract.GRADUATIONTHRESHOLD(&_LaunchpadFactory.CallOpts)
}

// GRADUATIONTHRESHOLD is a free data retrieval call binding the contract method 0xfcfc0c09.
//
// Solidity: function GRADUATION_THRESHOLD() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) GRADUATIONTHRESHOLD() (*big.Int, error) {
	return _LaunchpadFactory.Contract.GRADUATIONTHRESHOLD(&_LaunchpadFactory.CallOpts)
}

// INITIALAGENTRESERVE is a free data retrieval call binding the contract method 0x6da031f6.
//
// Solidity: function INITIAL_AGENT_RESERVE() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCaller) INITIALAGENTRESERVE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "INITIAL_AGENT_RESERVE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITIALAGENTRESERVE is a free data retrieval call binding the contract method 0x6da031f6.
//
// Solidity: function INITIAL_AGENT_RESERVE() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactorySession) INITIALAGENTRESERVE() (*big.Int, error) {
	return _LaunchpadFactory.Contract.INITIALAGENTRESERVE(&_LaunchpadFactory.CallOpts)
}

// INITIALAGENTRESERVE is a free data retrieval call binding the contract method 0x6da031f6.
//
// Solidity: function INITIAL_AGENT_RESERVE() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) INITIALAGENTRESERVE() (*big.Int, error) {
	return _LaunchpadFactory.Contract.INITIALAGENTRESERVE(&_LaunchpadFactory.CallOpts)
}

// MAXNAMELEN is a free data retrieval call binding the contract method 0x0a63d8bb.
//
// Solidity: function MAX_NAME_LEN() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCaller) MAXNAMELEN(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "MAX_NAME_LEN")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXNAMELEN is a free data retrieval call binding the contract method 0x0a63d8bb.
//
// Solidity: function MAX_NAME_LEN() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactorySession) MAXNAMELEN() (*big.Int, error) {
	return _LaunchpadFactory.Contract.MAXNAMELEN(&_LaunchpadFactory.CallOpts)
}

// MAXNAMELEN is a free data retrieval call binding the contract method 0x0a63d8bb.
//
// Solidity: function MAX_NAME_LEN() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) MAXNAMELEN() (*big.Int, error) {
	return _LaunchpadFactory.Contract.MAXNAMELEN(&_LaunchpadFactory.CallOpts)
}

// MAXSYMBOLLEN is a free data retrieval call binding the contract method 0x88a512b7.
//
// Solidity: function MAX_SYMBOL_LEN() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCaller) MAXSYMBOLLEN(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "MAX_SYMBOL_LEN")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXSYMBOLLEN is a free data retrieval call binding the contract method 0x88a512b7.
//
// Solidity: function MAX_SYMBOL_LEN() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactorySession) MAXSYMBOLLEN() (*big.Int, error) {
	return _LaunchpadFactory.Contract.MAXSYMBOLLEN(&_LaunchpadFactory.CallOpts)
}

// MAXSYMBOLLEN is a free data retrieval call binding the contract method 0x88a512b7.
//
// Solidity: function MAX_SYMBOL_LEN() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) MAXSYMBOLLEN() (*big.Int, error) {
	return _LaunchpadFactory.Contract.MAXSYMBOLLEN(&_LaunchpadFactory.CallOpts)
}

// MODULEIDAIRDROP is a free data retrieval call binding the contract method 0xae94763f.
//
// Solidity: function MODULE_ID_AIRDROP() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCaller) MODULEIDAIRDROP(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "MODULE_ID_AIRDROP")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MODULEIDAIRDROP is a free data retrieval call binding the contract method 0xae94763f.
//
// Solidity: function MODULE_ID_AIRDROP() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactorySession) MODULEIDAIRDROP() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDAIRDROP(&_LaunchpadFactory.CallOpts)
}

// MODULEIDAIRDROP is a free data retrieval call binding the contract method 0xae94763f.
//
// Solidity: function MODULE_ID_AIRDROP() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) MODULEIDAIRDROP() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDAIRDROP(&_LaunchpadFactory.CallOpts)
}

// MODULEIDANTISNIPER is a free data retrieval call binding the contract method 0x6828f47b.
//
// Solidity: function MODULE_ID_ANTI_SNIPER() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCaller) MODULEIDANTISNIPER(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "MODULE_ID_ANTI_SNIPER")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MODULEIDANTISNIPER is a free data retrieval call binding the contract method 0x6828f47b.
//
// Solidity: function MODULE_ID_ANTI_SNIPER() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactorySession) MODULEIDANTISNIPER() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDANTISNIPER(&_LaunchpadFactory.CallOpts)
}

// MODULEIDANTISNIPER is a free data retrieval call binding the contract method 0x6828f47b.
//
// Solidity: function MODULE_ID_ANTI_SNIPER() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) MODULEIDANTISNIPER() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDANTISNIPER(&_LaunchpadFactory.CallOpts)
}

// MODULEIDCAPITALFORMATION is a free data retrieval call binding the contract method 0x804f5b9c.
//
// Solidity: function MODULE_ID_CAPITAL_FORMATION() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCaller) MODULEIDCAPITALFORMATION(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "MODULE_ID_CAPITAL_FORMATION")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MODULEIDCAPITALFORMATION is a free data retrieval call binding the contract method 0x804f5b9c.
//
// Solidity: function MODULE_ID_CAPITAL_FORMATION() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactorySession) MODULEIDCAPITALFORMATION() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDCAPITALFORMATION(&_LaunchpadFactory.CallOpts)
}

// MODULEIDCAPITALFORMATION is a free data retrieval call binding the contract method 0x804f5b9c.
//
// Solidity: function MODULE_ID_CAPITAL_FORMATION() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) MODULEIDCAPITALFORMATION() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDCAPITALFORMATION(&_LaunchpadFactory.CallOpts)
}

// MODULEIDEXISTINGTOKEN is a free data retrieval call binding the contract method 0xe64180c0.
//
// Solidity: function MODULE_ID_EXISTING_TOKEN() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCaller) MODULEIDEXISTINGTOKEN(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "MODULE_ID_EXISTING_TOKEN")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MODULEIDEXISTINGTOKEN is a free data retrieval call binding the contract method 0xe64180c0.
//
// Solidity: function MODULE_ID_EXISTING_TOKEN() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactorySession) MODULEIDEXISTINGTOKEN() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDEXISTINGTOKEN(&_LaunchpadFactory.CallOpts)
}

// MODULEIDEXISTINGTOKEN is a free data retrieval call binding the contract method 0xe64180c0.
//
// Solidity: function MODULE_ID_EXISTING_TOKEN() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) MODULEIDEXISTINGTOKEN() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDEXISTINGTOKEN(&_LaunchpadFactory.CallOpts)
}

// MODULEIDLAUNCHRADAR is a free data retrieval call binding the contract method 0xeffe92a6.
//
// Solidity: function MODULE_ID_LAUNCH_RADAR() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCaller) MODULEIDLAUNCHRADAR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "MODULE_ID_LAUNCH_RADAR")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MODULEIDLAUNCHRADAR is a free data retrieval call binding the contract method 0xeffe92a6.
//
// Solidity: function MODULE_ID_LAUNCH_RADAR() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactorySession) MODULEIDLAUNCHRADAR() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDLAUNCHRADAR(&_LaunchpadFactory.CallOpts)
}

// MODULEIDLAUNCHRADAR is a free data retrieval call binding the contract method 0xeffe92a6.
//
// Solidity: function MODULE_ID_LAUNCH_RADAR() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) MODULEIDLAUNCHRADAR() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDLAUNCHRADAR(&_LaunchpadFactory.CallOpts)
}

// MODULEIDPREBUY is a free data retrieval call binding the contract method 0xda439155.
//
// Solidity: function MODULE_ID_PRE_BUY() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCaller) MODULEIDPREBUY(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "MODULE_ID_PRE_BUY")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MODULEIDPREBUY is a free data retrieval call binding the contract method 0xda439155.
//
// Solidity: function MODULE_ID_PRE_BUY() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactorySession) MODULEIDPREBUY() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDPREBUY(&_LaunchpadFactory.CallOpts)
}

// MODULEIDPREBUY is a free data retrieval call binding the contract method 0xda439155.
//
// Solidity: function MODULE_ID_PRE_BUY() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) MODULEIDPREBUY() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDPREBUY(&_LaunchpadFactory.CallOpts)
}

// MODULEIDSIXTYDAYS is a free data retrieval call binding the contract method 0x178e8777.
//
// Solidity: function MODULE_ID_SIXTY_DAYS() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCaller) MODULEIDSIXTYDAYS(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "MODULE_ID_SIXTY_DAYS")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// MODULEIDSIXTYDAYS is a free data retrieval call binding the contract method 0x178e8777.
//
// Solidity: function MODULE_ID_SIXTY_DAYS() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactorySession) MODULEIDSIXTYDAYS() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDSIXTYDAYS(&_LaunchpadFactory.CallOpts)
}

// MODULEIDSIXTYDAYS is a free data retrieval call binding the contract method 0x178e8777.
//
// Solidity: function MODULE_ID_SIXTY_DAYS() view returns(bytes32)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) MODULEIDSIXTYDAYS() ([32]byte, error) {
	return _LaunchpadFactory.Contract.MODULEIDSIXTYDAYS(&_LaunchpadFactory.CallOpts)
}

// VIRTUALAGENTRESERVE is a free data retrieval call binding the contract method 0xfeb8405f.
//
// Solidity: function VIRTUAL_AGENT_RESERVE() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCaller) VIRTUALAGENTRESERVE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "VIRTUAL_AGENT_RESERVE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VIRTUALAGENTRESERVE is a free data retrieval call binding the contract method 0xfeb8405f.
//
// Solidity: function VIRTUAL_AGENT_RESERVE() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactorySession) VIRTUALAGENTRESERVE() (*big.Int, error) {
	return _LaunchpadFactory.Contract.VIRTUALAGENTRESERVE(&_LaunchpadFactory.CallOpts)
}

// VIRTUALAGENTRESERVE is a free data retrieval call binding the contract method 0xfeb8405f.
//
// Solidity: function VIRTUAL_AGENT_RESERVE() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) VIRTUALAGENTRESERVE() (*big.Int, error) {
	return _LaunchpadFactory.Contract.VIRTUALAGENTRESERVE(&_LaunchpadFactory.CallOpts)
}

// VIRTUALQUOTERESERVE is a free data retrieval call binding the contract method 0xc0dfc9ea.
//
// Solidity: function VIRTUAL_QUOTE_RESERVE() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCaller) VIRTUALQUOTERESERVE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "VIRTUAL_QUOTE_RESERVE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VIRTUALQUOTERESERVE is a free data retrieval call binding the contract method 0xc0dfc9ea.
//
// Solidity: function VIRTUAL_QUOTE_RESERVE() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactorySession) VIRTUALQUOTERESERVE() (*big.Int, error) {
	return _LaunchpadFactory.Contract.VIRTUALQUOTERESERVE(&_LaunchpadFactory.CallOpts)
}

// VIRTUALQUOTERESERVE is a free data retrieval call binding the contract method 0xc0dfc9ea.
//
// Solidity: function VIRTUAL_QUOTE_RESERVE() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) VIRTUALQUOTERESERVE() (*big.Int, error) {
	return _LaunchpadFactory.Contract.VIRTUALQUOTERESERVE(&_LaunchpadFactory.CallOpts)
}

// AgentCount is a free data retrieval call binding the contract method 0xb7dc1284.
//
// Solidity: function agentCount() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCaller) AgentCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "agentCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AgentCount is a free data retrieval call binding the contract method 0xb7dc1284.
//
// Solidity: function agentCount() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactorySession) AgentCount() (*big.Int, error) {
	return _LaunchpadFactory.Contract.AgentCount(&_LaunchpadFactory.CallOpts)
}

// AgentCount is a free data retrieval call binding the contract method 0xb7dc1284.
//
// Solidity: function agentCount() view returns(uint256)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) AgentCount() (*big.Int, error) {
	return _LaunchpadFactory.Contract.AgentCount(&_LaunchpadFactory.CallOpts)
}

// AgentTokenImpl is a free data retrieval call binding the contract method 0x6d32fa52.
//
// Solidity: function agentTokenImpl() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCaller) AgentTokenImpl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "agentTokenImpl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentTokenImpl is a free data retrieval call binding the contract method 0x6d32fa52.
//
// Solidity: function agentTokenImpl() view returns(address)
func (_LaunchpadFactory *LaunchpadFactorySession) AgentTokenImpl() (common.Address, error) {
	return _LaunchpadFactory.Contract.AgentTokenImpl(&_LaunchpadFactory.CallOpts)
}

// AgentTokenImpl is a free data retrieval call binding the contract method 0x6d32fa52.
//
// Solidity: function agentTokenImpl() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) AgentTokenImpl() (common.Address, error) {
	return _LaunchpadFactory.Contract.AgentTokenImpl(&_LaunchpadFactory.CallOpts)
}

// BondingCurveImpl is a free data retrieval call binding the contract method 0x2a4fb36a.
//
// Solidity: function bondingCurveImpl() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCaller) BondingCurveImpl(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "bondingCurveImpl")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BondingCurveImpl is a free data retrieval call binding the contract method 0x2a4fb36a.
//
// Solidity: function bondingCurveImpl() view returns(address)
func (_LaunchpadFactory *LaunchpadFactorySession) BondingCurveImpl() (common.Address, error) {
	return _LaunchpadFactory.Contract.BondingCurveImpl(&_LaunchpadFactory.CallOpts)
}

// BondingCurveImpl is a free data retrieval call binding the contract method 0x2a4fb36a.
//
// Solidity: function bondingCurveImpl() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) BondingCurveImpl() (common.Address, error) {
	return _LaunchpadFactory.Contract.BondingCurveImpl(&_LaunchpadFactory.CallOpts)
}

// FeeRouter is a free data retrieval call binding the contract method 0xf29ebf61.
//
// Solidity: function feeRouter() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCaller) FeeRouter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "feeRouter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeRouter is a free data retrieval call binding the contract method 0xf29ebf61.
//
// Solidity: function feeRouter() view returns(address)
func (_LaunchpadFactory *LaunchpadFactorySession) FeeRouter() (common.Address, error) {
	return _LaunchpadFactory.Contract.FeeRouter(&_LaunchpadFactory.CallOpts)
}

// FeeRouter is a free data retrieval call binding the contract method 0xf29ebf61.
//
// Solidity: function feeRouter() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) FeeRouter() (common.Address, error) {
	return _LaunchpadFactory.Contract.FeeRouter(&_LaunchpadFactory.CallOpts)
}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((address,address,address,address,address,uint64,bytes32[]))
func (_LaunchpadFactory *LaunchpadFactoryCaller) GetAgent(opts *bind.CallOpts, agentId *big.Int) (LaunchpadFactoryAgent, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "getAgent", agentId)

	if err != nil {
		return *new(LaunchpadFactoryAgent), err
	}

	out0 := *abi.ConvertType(out[0], new(LaunchpadFactoryAgent)).(*LaunchpadFactoryAgent)

	return out0, err

}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((address,address,address,address,address,uint64,bytes32[]))
func (_LaunchpadFactory *LaunchpadFactorySession) GetAgent(agentId *big.Int) (LaunchpadFactoryAgent, error) {
	return _LaunchpadFactory.Contract.GetAgent(&_LaunchpadFactory.CallOpts, agentId)
}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((address,address,address,address,address,uint64,bytes32[]))
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) GetAgent(agentId *big.Int) (LaunchpadFactoryAgent, error) {
	return _LaunchpadFactory.Contract.GetAgent(&_LaunchpadFactory.CallOpts, agentId)
}

// Graduator is a free data retrieval call binding the contract method 0x0245a483.
//
// Solidity: function graduator() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCaller) Graduator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "graduator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Graduator is a free data retrieval call binding the contract method 0x0245a483.
//
// Solidity: function graduator() view returns(address)
func (_LaunchpadFactory *LaunchpadFactorySession) Graduator() (common.Address, error) {
	return _LaunchpadFactory.Contract.Graduator(&_LaunchpadFactory.CallOpts)
}

// Graduator is a free data retrieval call binding the contract method 0x0245a483.
//
// Solidity: function graduator() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) Graduator() (common.Address, error) {
	return _LaunchpadFactory.Contract.Graduator(&_LaunchpadFactory.CallOpts)
}

// Modules is a free data retrieval call binding the contract method 0xb0b6cc1a.
//
// Solidity: function modules(bytes32 moduleId) view returns(address module)
func (_LaunchpadFactory *LaunchpadFactoryCaller) Modules(opts *bind.CallOpts, moduleId [32]byte) (common.Address, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "modules", moduleId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Modules is a free data retrieval call binding the contract method 0xb0b6cc1a.
//
// Solidity: function modules(bytes32 moduleId) view returns(address module)
func (_LaunchpadFactory *LaunchpadFactorySession) Modules(moduleId [32]byte) (common.Address, error) {
	return _LaunchpadFactory.Contract.Modules(&_LaunchpadFactory.CallOpts, moduleId)
}

// Modules is a free data retrieval call binding the contract method 0xb0b6cc1a.
//
// Solidity: function modules(bytes32 moduleId) view returns(address module)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) Modules(moduleId [32]byte) (common.Address, error) {
	return _LaunchpadFactory.Contract.Modules(&_LaunchpadFactory.CallOpts, moduleId)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LaunchpadFactory *LaunchpadFactorySession) Owner() (common.Address, error) {
	return _LaunchpadFactory.Contract.Owner(&_LaunchpadFactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) Owner() (common.Address, error) {
	return _LaunchpadFactory.Contract.Owner(&_LaunchpadFactory.CallOpts)
}

// TituToken is a free data retrieval call binding the contract method 0x7862355c.
//
// Solidity: function tituToken() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCaller) TituToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "tituToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TituToken is a free data retrieval call binding the contract method 0x7862355c.
//
// Solidity: function tituToken() view returns(address)
func (_LaunchpadFactory *LaunchpadFactorySession) TituToken() (common.Address, error) {
	return _LaunchpadFactory.Contract.TituToken(&_LaunchpadFactory.CallOpts)
}

// TituToken is a free data retrieval call binding the contract method 0x7862355c.
//
// Solidity: function tituToken() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) TituToken() (common.Address, error) {
	return _LaunchpadFactory.Contract.TituToken(&_LaunchpadFactory.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCaller) Treasury(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "treasury")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_LaunchpadFactory *LaunchpadFactorySession) Treasury() (common.Address, error) {
	return _LaunchpadFactory.Contract.Treasury(&_LaunchpadFactory.CallOpts)
}

// Treasury is a free data retrieval call binding the contract method 0x61d027b3.
//
// Solidity: function treasury() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) Treasury() (common.Address, error) {
	return _LaunchpadFactory.Contract.Treasury(&_LaunchpadFactory.CallOpts)
}

// UniV2Factory is a free data retrieval call binding the contract method 0x3da04b87.
//
// Solidity: function uniV2Factory() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCaller) UniV2Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "uniV2Factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UniV2Factory is a free data retrieval call binding the contract method 0x3da04b87.
//
// Solidity: function uniV2Factory() view returns(address)
func (_LaunchpadFactory *LaunchpadFactorySession) UniV2Factory() (common.Address, error) {
	return _LaunchpadFactory.Contract.UniV2Factory(&_LaunchpadFactory.CallOpts)
}

// UniV2Factory is a free data retrieval call binding the contract method 0x3da04b87.
//
// Solidity: function uniV2Factory() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) UniV2Factory() (common.Address, error) {
	return _LaunchpadFactory.Contract.UniV2Factory(&_LaunchpadFactory.CallOpts)
}

// UniV2Router is a free data retrieval call binding the contract method 0x958c2e52.
//
// Solidity: function uniV2Router() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCaller) UniV2Router(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LaunchpadFactory.contract.Call(opts, &out, "uniV2Router")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UniV2Router is a free data retrieval call binding the contract method 0x958c2e52.
//
// Solidity: function uniV2Router() view returns(address)
func (_LaunchpadFactory *LaunchpadFactorySession) UniV2Router() (common.Address, error) {
	return _LaunchpadFactory.Contract.UniV2Router(&_LaunchpadFactory.CallOpts)
}

// UniV2Router is a free data retrieval call binding the contract method 0x958c2e52.
//
// Solidity: function uniV2Router() view returns(address)
func (_LaunchpadFactory *LaunchpadFactoryCallerSession) UniV2Router() (common.Address, error) {
	return _LaunchpadFactory.Contract.UniV2Router(&_LaunchpadFactory.CallOpts)
}

// LaunchAgent is a paid mutator transaction binding the contract method 0x8693c5c8.
//
// Solidity: function launchAgent((string,string,string,string,bytes32[],bytes[]) params) returns(address agentToken, address bondingCurve, uint256 agentId)
func (_LaunchpadFactory *LaunchpadFactoryTransactor) LaunchAgent(opts *bind.TransactOpts, params LaunchpadFactoryLaunchParams) (*types.Transaction, error) {
	return _LaunchpadFactory.contract.Transact(opts, "launchAgent", params)
}

// LaunchAgent is a paid mutator transaction binding the contract method 0x8693c5c8.
//
// Solidity: function launchAgent((string,string,string,string,bytes32[],bytes[]) params) returns(address agentToken, address bondingCurve, uint256 agentId)
func (_LaunchpadFactory *LaunchpadFactorySession) LaunchAgent(params LaunchpadFactoryLaunchParams) (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.LaunchAgent(&_LaunchpadFactory.TransactOpts, params)
}

// LaunchAgent is a paid mutator transaction binding the contract method 0x8693c5c8.
//
// Solidity: function launchAgent((string,string,string,string,bytes32[],bytes[]) params) returns(address agentToken, address bondingCurve, uint256 agentId)
func (_LaunchpadFactory *LaunchpadFactoryTransactorSession) LaunchAgent(params LaunchpadFactoryLaunchParams) (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.LaunchAgent(&_LaunchpadFactory.TransactOpts, params)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LaunchpadFactory *LaunchpadFactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LaunchpadFactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LaunchpadFactory *LaunchpadFactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.RenounceOwnership(&_LaunchpadFactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LaunchpadFactory *LaunchpadFactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.RenounceOwnership(&_LaunchpadFactory.TransactOpts)
}

// SetModule is a paid mutator transaction binding the contract method 0x541cd468.
//
// Solidity: function setModule(bytes32 moduleId, address module) returns()
func (_LaunchpadFactory *LaunchpadFactoryTransactor) SetModule(opts *bind.TransactOpts, moduleId [32]byte, module common.Address) (*types.Transaction, error) {
	return _LaunchpadFactory.contract.Transact(opts, "setModule", moduleId, module)
}

// SetModule is a paid mutator transaction binding the contract method 0x541cd468.
//
// Solidity: function setModule(bytes32 moduleId, address module) returns()
func (_LaunchpadFactory *LaunchpadFactorySession) SetModule(moduleId [32]byte, module common.Address) (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.SetModule(&_LaunchpadFactory.TransactOpts, moduleId, module)
}

// SetModule is a paid mutator transaction binding the contract method 0x541cd468.
//
// Solidity: function setModule(bytes32 moduleId, address module) returns()
func (_LaunchpadFactory *LaunchpadFactoryTransactorSession) SetModule(moduleId [32]byte, module common.Address) (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.SetModule(&_LaunchpadFactory.TransactOpts, moduleId, module)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LaunchpadFactory *LaunchpadFactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _LaunchpadFactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LaunchpadFactory *LaunchpadFactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.TransferOwnership(&_LaunchpadFactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LaunchpadFactory *LaunchpadFactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _LaunchpadFactory.Contract.TransferOwnership(&_LaunchpadFactory.TransactOpts, newOwner)
}

// LaunchpadFactoryAgentLaunchedIterator is returned from FilterAgentLaunched and is used to iterate over the raw logs and unpacked data for AgentLaunched events raised by the LaunchpadFactory contract.
type LaunchpadFactoryAgentLaunchedIterator struct {
	Event *LaunchpadFactoryAgentLaunched // Event containing the contract specifics and raw log

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
func (it *LaunchpadFactoryAgentLaunchedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LaunchpadFactoryAgentLaunched)
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
		it.Event = new(LaunchpadFactoryAgentLaunched)
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
func (it *LaunchpadFactoryAgentLaunchedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LaunchpadFactoryAgentLaunchedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LaunchpadFactoryAgentLaunched represents a AgentLaunched event raised by the LaunchpadFactory contract.
type LaunchpadFactoryAgentLaunched struct {
	AgentId *big.Int
	Token   common.Address
	Curve   common.Address
	Creator common.Address
	LpLock  common.Address
	Pair    common.Address
	Modules [][32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterAgentLaunched is a free log retrieval operation binding the contract event 0xc480e8abcc9946f900ed3b50d1a70c7afdc2c4f8254ed2a347a425f348dc41f8.
//
// Solidity: event AgentLaunched(uint256 indexed agentId, address indexed token, address indexed curve, address creator, address lpLock, address pair, bytes32[] modules)
func (_LaunchpadFactory *LaunchpadFactoryFilterer) FilterAgentLaunched(opts *bind.FilterOpts, agentId []*big.Int, token []common.Address, curve []common.Address) (*LaunchpadFactoryAgentLaunchedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var curveRule []interface{}
	for _, curveItem := range curve {
		curveRule = append(curveRule, curveItem)
	}

	logs, sub, err := _LaunchpadFactory.contract.FilterLogs(opts, "AgentLaunched", agentIdRule, tokenRule, curveRule)
	if err != nil {
		return nil, err
	}
	return &LaunchpadFactoryAgentLaunchedIterator{contract: _LaunchpadFactory.contract, event: "AgentLaunched", logs: logs, sub: sub}, nil
}

// WatchAgentLaunched is a free log subscription operation binding the contract event 0xc480e8abcc9946f900ed3b50d1a70c7afdc2c4f8254ed2a347a425f348dc41f8.
//
// Solidity: event AgentLaunched(uint256 indexed agentId, address indexed token, address indexed curve, address creator, address lpLock, address pair, bytes32[] modules)
func (_LaunchpadFactory *LaunchpadFactoryFilterer) WatchAgentLaunched(opts *bind.WatchOpts, sink chan<- *LaunchpadFactoryAgentLaunched, agentId []*big.Int, token []common.Address, curve []common.Address) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var curveRule []interface{}
	for _, curveItem := range curve {
		curveRule = append(curveRule, curveItem)
	}

	logs, sub, err := _LaunchpadFactory.contract.WatchLogs(opts, "AgentLaunched", agentIdRule, tokenRule, curveRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LaunchpadFactoryAgentLaunched)
				if err := _LaunchpadFactory.contract.UnpackLog(event, "AgentLaunched", log); err != nil {
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

// ParseAgentLaunched is a log parse operation binding the contract event 0xc480e8abcc9946f900ed3b50d1a70c7afdc2c4f8254ed2a347a425f348dc41f8.
//
// Solidity: event AgentLaunched(uint256 indexed agentId, address indexed token, address indexed curve, address creator, address lpLock, address pair, bytes32[] modules)
func (_LaunchpadFactory *LaunchpadFactoryFilterer) ParseAgentLaunched(log types.Log) (*LaunchpadFactoryAgentLaunched, error) {
	event := new(LaunchpadFactoryAgentLaunched)
	if err := _LaunchpadFactory.contract.UnpackLog(event, "AgentLaunched", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LaunchpadFactoryModuleSetIterator is returned from FilterModuleSet and is used to iterate over the raw logs and unpacked data for ModuleSet events raised by the LaunchpadFactory contract.
type LaunchpadFactoryModuleSetIterator struct {
	Event *LaunchpadFactoryModuleSet // Event containing the contract specifics and raw log

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
func (it *LaunchpadFactoryModuleSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LaunchpadFactoryModuleSet)
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
		it.Event = new(LaunchpadFactoryModuleSet)
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
func (it *LaunchpadFactoryModuleSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LaunchpadFactoryModuleSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LaunchpadFactoryModuleSet represents a ModuleSet event raised by the LaunchpadFactory contract.
type LaunchpadFactoryModuleSet struct {
	ModuleId [32]byte
	Prev     common.Address
	Module   common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterModuleSet is a free log retrieval operation binding the contract event 0x805f6cfd1882ecfb5ba14cc1b701b7242cf95893badde12b29ffd0236bef0422.
//
// Solidity: event ModuleSet(bytes32 indexed moduleId, address indexed prev, address indexed module)
func (_LaunchpadFactory *LaunchpadFactoryFilterer) FilterModuleSet(opts *bind.FilterOpts, moduleId [][32]byte, prev []common.Address, module []common.Address) (*LaunchpadFactoryModuleSetIterator, error) {

	var moduleIdRule []interface{}
	for _, moduleIdItem := range moduleId {
		moduleIdRule = append(moduleIdRule, moduleIdItem)
	}
	var prevRule []interface{}
	for _, prevItem := range prev {
		prevRule = append(prevRule, prevItem)
	}
	var moduleRule []interface{}
	for _, moduleItem := range module {
		moduleRule = append(moduleRule, moduleItem)
	}

	logs, sub, err := _LaunchpadFactory.contract.FilterLogs(opts, "ModuleSet", moduleIdRule, prevRule, moduleRule)
	if err != nil {
		return nil, err
	}
	return &LaunchpadFactoryModuleSetIterator{contract: _LaunchpadFactory.contract, event: "ModuleSet", logs: logs, sub: sub}, nil
}

// WatchModuleSet is a free log subscription operation binding the contract event 0x805f6cfd1882ecfb5ba14cc1b701b7242cf95893badde12b29ffd0236bef0422.
//
// Solidity: event ModuleSet(bytes32 indexed moduleId, address indexed prev, address indexed module)
func (_LaunchpadFactory *LaunchpadFactoryFilterer) WatchModuleSet(opts *bind.WatchOpts, sink chan<- *LaunchpadFactoryModuleSet, moduleId [][32]byte, prev []common.Address, module []common.Address) (event.Subscription, error) {

	var moduleIdRule []interface{}
	for _, moduleIdItem := range moduleId {
		moduleIdRule = append(moduleIdRule, moduleIdItem)
	}
	var prevRule []interface{}
	for _, prevItem := range prev {
		prevRule = append(prevRule, prevItem)
	}
	var moduleRule []interface{}
	for _, moduleItem := range module {
		moduleRule = append(moduleRule, moduleItem)
	}

	logs, sub, err := _LaunchpadFactory.contract.WatchLogs(opts, "ModuleSet", moduleIdRule, prevRule, moduleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LaunchpadFactoryModuleSet)
				if err := _LaunchpadFactory.contract.UnpackLog(event, "ModuleSet", log); err != nil {
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

// ParseModuleSet is a log parse operation binding the contract event 0x805f6cfd1882ecfb5ba14cc1b701b7242cf95893badde12b29ffd0236bef0422.
//
// Solidity: event ModuleSet(bytes32 indexed moduleId, address indexed prev, address indexed module)
func (_LaunchpadFactory *LaunchpadFactoryFilterer) ParseModuleSet(log types.Log) (*LaunchpadFactoryModuleSet, error) {
	event := new(LaunchpadFactoryModuleSet)
	if err := _LaunchpadFactory.contract.UnpackLog(event, "ModuleSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LaunchpadFactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the LaunchpadFactory contract.
type LaunchpadFactoryOwnershipTransferredIterator struct {
	Event *LaunchpadFactoryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *LaunchpadFactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LaunchpadFactoryOwnershipTransferred)
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
		it.Event = new(LaunchpadFactoryOwnershipTransferred)
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
func (it *LaunchpadFactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LaunchpadFactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LaunchpadFactoryOwnershipTransferred represents a OwnershipTransferred event raised by the LaunchpadFactory contract.
type LaunchpadFactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LaunchpadFactory *LaunchpadFactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*LaunchpadFactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LaunchpadFactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &LaunchpadFactoryOwnershipTransferredIterator{contract: _LaunchpadFactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LaunchpadFactory *LaunchpadFactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LaunchpadFactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LaunchpadFactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LaunchpadFactoryOwnershipTransferred)
				if err := _LaunchpadFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_LaunchpadFactory *LaunchpadFactoryFilterer) ParseOwnershipTransferred(log types.Log) (*LaunchpadFactoryOwnershipTransferred, error) {
	event := new(LaunchpadFactoryOwnershipTransferred)
	if err := _LaunchpadFactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
