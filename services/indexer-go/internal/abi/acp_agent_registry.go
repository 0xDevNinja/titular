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

// AgentRegistryAgentInfo is an auto generated low-level Go binding around an user-defined struct.
type AgentRegistryAgentInfo struct {
	Controller      common.Address
	MetadataURI     string
	Capabilities    *big.Int
	ReputationScore *big.Int
	RegisteredAt    *big.Int
	Active          bool
}

// AgentRegistryMetaData contains all meta data concerning the AgentRegistry contract.
var AgentRegistryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PAUSER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SCORER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptControllerTransfer\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"agentsByController\",\"inputs\":[{\"name\":\"controller\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"ids\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelControllerTransfer\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAgent\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"info\",\"type\":\"tuple\",\"internalType\":\"structAgentRegistry.AgentInfo\",\"components\":[{\"name\":\"controller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"capabilities\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reputationScore\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"registeredAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingController\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postScore\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"delta\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proposeControllerTransfer\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"proposed\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"register\",\"inputs\":[{\"name\":\"controller\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"capabilities\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"scorerNonce\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setActive\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCapabilities\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"capabilities\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMetadata\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalAgents\",\"inputs\":[],\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ActiveStatusChanged\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"active\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"AgentRegistered\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"controller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"capabilities\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CapabilitiesUpdated\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"capabilities\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ControllerTransferAccepted\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"oldController\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newController\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ControllerTransferCancelled\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ControllerTransferProposed\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"proposed\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataUpdated\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"metadataURI\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ScorePosted\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"scorer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"delta\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"},{\"name\":\"newTotal\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"},{\"name\":\"nonce\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AgentInactive\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"AgentNotFound\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"EmptyMetadataURI\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidNonce\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expected\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"provided\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NoPendingTransfer\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotController\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"NotProposedController\",\"inputs\":[{\"name\":\"agentId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
}

// AgentRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use AgentRegistryMetaData.ABI instead.
var AgentRegistryABI = AgentRegistryMetaData.ABI

// AgentRegistry is an auto generated Go binding around an Ethereum contract.
type AgentRegistry struct {
	AgentRegistryCaller     // Read-only binding to the contract
	AgentRegistryTransactor // Write-only binding to the contract
	AgentRegistryFilterer   // Log filterer for contract events
}

// AgentRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AgentRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AgentRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgentRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgentRegistrySession struct {
	Contract     *AgentRegistry    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AgentRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgentRegistryCallerSession struct {
	Contract *AgentRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// AgentRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgentRegistryTransactorSession struct {
	Contract     *AgentRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// AgentRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AgentRegistryRaw struct {
	Contract *AgentRegistry // Generic contract binding to access the raw methods on
}

// AgentRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgentRegistryCallerRaw struct {
	Contract *AgentRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// AgentRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgentRegistryTransactorRaw struct {
	Contract *AgentRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAgentRegistry creates a new instance of AgentRegistry, bound to a specific deployed contract.
func NewAgentRegistry(address common.Address, backend bind.ContractBackend) (*AgentRegistry, error) {
	contract, err := bindAgentRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AgentRegistry{AgentRegistryCaller: AgentRegistryCaller{contract: contract}, AgentRegistryTransactor: AgentRegistryTransactor{contract: contract}, AgentRegistryFilterer: AgentRegistryFilterer{contract: contract}}, nil
}

// NewAgentRegistryCaller creates a new read-only instance of AgentRegistry, bound to a specific deployed contract.
func NewAgentRegistryCaller(address common.Address, caller bind.ContractCaller) (*AgentRegistryCaller, error) {
	contract, err := bindAgentRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryCaller{contract: contract}, nil
}

// NewAgentRegistryTransactor creates a new write-only instance of AgentRegistry, bound to a specific deployed contract.
func NewAgentRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AgentRegistryTransactor, error) {
	contract, err := bindAgentRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryTransactor{contract: contract}, nil
}

// NewAgentRegistryFilterer creates a new log filterer instance of AgentRegistry, bound to a specific deployed contract.
func NewAgentRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AgentRegistryFilterer, error) {
	contract, err := bindAgentRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryFilterer{contract: contract}, nil
}

// bindAgentRegistry binds a generic wrapper to an already deployed contract.
func bindAgentRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AgentRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentRegistry *AgentRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentRegistry.Contract.AgentRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentRegistry *AgentRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AgentRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentRegistry *AgentRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AgentRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentRegistry *AgentRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentRegistry *AgentRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentRegistry *AgentRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentRegistry.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistryCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistrySession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgentRegistry.Contract.DEFAULTADMINROLE(&_AgentRegistry.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistryCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _AgentRegistry.Contract.DEFAULTADMINROLE(&_AgentRegistry.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistryCaller) PAUSERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "PAUSER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistrySession) PAUSERROLE() ([32]byte, error) {
	return _AgentRegistry.Contract.PAUSERROLE(&_AgentRegistry.CallOpts)
}

// PAUSERROLE is a free data retrieval call binding the contract method 0xe63ab1e9.
//
// Solidity: function PAUSER_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistryCallerSession) PAUSERROLE() ([32]byte, error) {
	return _AgentRegistry.Contract.PAUSERROLE(&_AgentRegistry.CallOpts)
}

// SCORERROLE is a free data retrieval call binding the contract method 0xa4d6d495.
//
// Solidity: function SCORER_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistryCaller) SCORERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "SCORER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SCORERROLE is a free data retrieval call binding the contract method 0xa4d6d495.
//
// Solidity: function SCORER_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistrySession) SCORERROLE() ([32]byte, error) {
	return _AgentRegistry.Contract.SCORERROLE(&_AgentRegistry.CallOpts)
}

// SCORERROLE is a free data retrieval call binding the contract method 0xa4d6d495.
//
// Solidity: function SCORER_ROLE() view returns(bytes32)
func (_AgentRegistry *AgentRegistryCallerSession) SCORERROLE() ([32]byte, error) {
	return _AgentRegistry.Contract.SCORERROLE(&_AgentRegistry.CallOpts)
}

// AgentsByController is a free data retrieval call binding the contract method 0x53c0ada8.
//
// Solidity: function agentsByController(address controller) view returns(uint256[] ids)
func (_AgentRegistry *AgentRegistryCaller) AgentsByController(opts *bind.CallOpts, controller common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "agentsByController", controller)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// AgentsByController is a free data retrieval call binding the contract method 0x53c0ada8.
//
// Solidity: function agentsByController(address controller) view returns(uint256[] ids)
func (_AgentRegistry *AgentRegistrySession) AgentsByController(controller common.Address) ([]*big.Int, error) {
	return _AgentRegistry.Contract.AgentsByController(&_AgentRegistry.CallOpts, controller)
}

// AgentsByController is a free data retrieval call binding the contract method 0x53c0ada8.
//
// Solidity: function agentsByController(address controller) view returns(uint256[] ids)
func (_AgentRegistry *AgentRegistryCallerSession) AgentsByController(controller common.Address) ([]*big.Int, error) {
	return _AgentRegistry.Contract.AgentsByController(&_AgentRegistry.CallOpts, controller)
}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((address,string,uint256,int256,uint48,bool) info)
func (_AgentRegistry *AgentRegistryCaller) GetAgent(opts *bind.CallOpts, agentId *big.Int) (AgentRegistryAgentInfo, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "getAgent", agentId)

	if err != nil {
		return *new(AgentRegistryAgentInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(AgentRegistryAgentInfo)).(*AgentRegistryAgentInfo)

	return out0, err

}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((address,string,uint256,int256,uint48,bool) info)
func (_AgentRegistry *AgentRegistrySession) GetAgent(agentId *big.Int) (AgentRegistryAgentInfo, error) {
	return _AgentRegistry.Contract.GetAgent(&_AgentRegistry.CallOpts, agentId)
}

// GetAgent is a free data retrieval call binding the contract method 0x2de5aaf7.
//
// Solidity: function getAgent(uint256 agentId) view returns((address,string,uint256,int256,uint48,bool) info)
func (_AgentRegistry *AgentRegistryCallerSession) GetAgent(agentId *big.Int) (AgentRegistryAgentInfo, error) {
	return _AgentRegistry.Contract.GetAgent(&_AgentRegistry.CallOpts, agentId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentRegistry *AgentRegistryCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentRegistry *AgentRegistrySession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _AgentRegistry.Contract.GetRoleAdmin(&_AgentRegistry.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_AgentRegistry *AgentRegistryCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _AgentRegistry.Contract.GetRoleAdmin(&_AgentRegistry.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _AgentRegistry.Contract.HasRole(&_AgentRegistry.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _AgentRegistry.Contract.HasRole(&_AgentRegistry.CallOpts, role, account)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AgentRegistry *AgentRegistrySession) Paused() (bool, error) {
	return _AgentRegistry.Contract.Paused(&_AgentRegistry.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) Paused() (bool, error) {
	return _AgentRegistry.Contract.Paused(&_AgentRegistry.CallOpts)
}

// PendingController is a free data retrieval call binding the contract method 0x5c3bd345.
//
// Solidity: function pendingController(uint256 agentId) view returns(address proposed)
func (_AgentRegistry *AgentRegistryCaller) PendingController(opts *bind.CallOpts, agentId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "pendingController", agentId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingController is a free data retrieval call binding the contract method 0x5c3bd345.
//
// Solidity: function pendingController(uint256 agentId) view returns(address proposed)
func (_AgentRegistry *AgentRegistrySession) PendingController(agentId *big.Int) (common.Address, error) {
	return _AgentRegistry.Contract.PendingController(&_AgentRegistry.CallOpts, agentId)
}

// PendingController is a free data retrieval call binding the contract method 0x5c3bd345.
//
// Solidity: function pendingController(uint256 agentId) view returns(address proposed)
func (_AgentRegistry *AgentRegistryCallerSession) PendingController(agentId *big.Int) (common.Address, error) {
	return _AgentRegistry.Contract.PendingController(&_AgentRegistry.CallOpts, agentId)
}

// ScorerNonce is a free data retrieval call binding the contract method 0x2d462cea.
//
// Solidity: function scorerNonce(uint256 ) view returns(uint256)
func (_AgentRegistry *AgentRegistryCaller) ScorerNonce(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "scorerNonce", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ScorerNonce is a free data retrieval call binding the contract method 0x2d462cea.
//
// Solidity: function scorerNonce(uint256 ) view returns(uint256)
func (_AgentRegistry *AgentRegistrySession) ScorerNonce(arg0 *big.Int) (*big.Int, error) {
	return _AgentRegistry.Contract.ScorerNonce(&_AgentRegistry.CallOpts, arg0)
}

// ScorerNonce is a free data retrieval call binding the contract method 0x2d462cea.
//
// Solidity: function scorerNonce(uint256 ) view returns(uint256)
func (_AgentRegistry *AgentRegistryCallerSession) ScorerNonce(arg0 *big.Int) (*big.Int, error) {
	return _AgentRegistry.Contract.ScorerNonce(&_AgentRegistry.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentRegistry *AgentRegistryCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentRegistry *AgentRegistrySession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgentRegistry.Contract.SupportsInterface(&_AgentRegistry.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_AgentRegistry *AgentRegistryCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _AgentRegistry.Contract.SupportsInterface(&_AgentRegistry.CallOpts, interfaceId)
}

// TotalAgents is a free data retrieval call binding the contract method 0xc5053712.
//
// Solidity: function totalAgents() view returns(uint256 count)
func (_AgentRegistry *AgentRegistryCaller) TotalAgents(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentRegistry.contract.Call(opts, &out, "totalAgents")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAgents is a free data retrieval call binding the contract method 0xc5053712.
//
// Solidity: function totalAgents() view returns(uint256 count)
func (_AgentRegistry *AgentRegistrySession) TotalAgents() (*big.Int, error) {
	return _AgentRegistry.Contract.TotalAgents(&_AgentRegistry.CallOpts)
}

// TotalAgents is a free data retrieval call binding the contract method 0xc5053712.
//
// Solidity: function totalAgents() view returns(uint256 count)
func (_AgentRegistry *AgentRegistryCallerSession) TotalAgents() (*big.Int, error) {
	return _AgentRegistry.Contract.TotalAgents(&_AgentRegistry.CallOpts)
}

// AcceptControllerTransfer is a paid mutator transaction binding the contract method 0xf87ebe5e.
//
// Solidity: function acceptControllerTransfer(uint256 agentId) returns()
func (_AgentRegistry *AgentRegistryTransactor) AcceptControllerTransfer(opts *bind.TransactOpts, agentId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "acceptControllerTransfer", agentId)
}

// AcceptControllerTransfer is a paid mutator transaction binding the contract method 0xf87ebe5e.
//
// Solidity: function acceptControllerTransfer(uint256 agentId) returns()
func (_AgentRegistry *AgentRegistrySession) AcceptControllerTransfer(agentId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AcceptControllerTransfer(&_AgentRegistry.TransactOpts, agentId)
}

// AcceptControllerTransfer is a paid mutator transaction binding the contract method 0xf87ebe5e.
//
// Solidity: function acceptControllerTransfer(uint256 agentId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) AcceptControllerTransfer(agentId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.AcceptControllerTransfer(&_AgentRegistry.TransactOpts, agentId)
}

// CancelControllerTransfer is a paid mutator transaction binding the contract method 0x5d71e0a1.
//
// Solidity: function cancelControllerTransfer(uint256 agentId) returns()
func (_AgentRegistry *AgentRegistryTransactor) CancelControllerTransfer(opts *bind.TransactOpts, agentId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "cancelControllerTransfer", agentId)
}

// CancelControllerTransfer is a paid mutator transaction binding the contract method 0x5d71e0a1.
//
// Solidity: function cancelControllerTransfer(uint256 agentId) returns()
func (_AgentRegistry *AgentRegistrySession) CancelControllerTransfer(agentId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.CancelControllerTransfer(&_AgentRegistry.TransactOpts, agentId)
}

// CancelControllerTransfer is a paid mutator transaction binding the contract method 0x5d71e0a1.
//
// Solidity: function cancelControllerTransfer(uint256 agentId) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) CancelControllerTransfer(agentId *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.CancelControllerTransfer(&_AgentRegistry.TransactOpts, agentId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentRegistry *AgentRegistryTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentRegistry *AgentRegistrySession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.GrantRole(&_AgentRegistry.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.GrantRole(&_AgentRegistry.TransactOpts, role, account)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AgentRegistry *AgentRegistryTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AgentRegistry *AgentRegistrySession) Pause() (*types.Transaction, error) {
	return _AgentRegistry.Contract.Pause(&_AgentRegistry.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_AgentRegistry *AgentRegistryTransactorSession) Pause() (*types.Transaction, error) {
	return _AgentRegistry.Contract.Pause(&_AgentRegistry.TransactOpts)
}

// PostScore is a paid mutator transaction binding the contract method 0x9c641ce0.
//
// Solidity: function postScore(uint256 agentId, int256 delta, uint256 nonce) returns()
func (_AgentRegistry *AgentRegistryTransactor) PostScore(opts *bind.TransactOpts, agentId *big.Int, delta *big.Int, nonce *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "postScore", agentId, delta, nonce)
}

// PostScore is a paid mutator transaction binding the contract method 0x9c641ce0.
//
// Solidity: function postScore(uint256 agentId, int256 delta, uint256 nonce) returns()
func (_AgentRegistry *AgentRegistrySession) PostScore(agentId *big.Int, delta *big.Int, nonce *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.PostScore(&_AgentRegistry.TransactOpts, agentId, delta, nonce)
}

// PostScore is a paid mutator transaction binding the contract method 0x9c641ce0.
//
// Solidity: function postScore(uint256 agentId, int256 delta, uint256 nonce) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) PostScore(agentId *big.Int, delta *big.Int, nonce *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.PostScore(&_AgentRegistry.TransactOpts, agentId, delta, nonce)
}

// ProposeControllerTransfer is a paid mutator transaction binding the contract method 0x98310b5e.
//
// Solidity: function proposeControllerTransfer(uint256 agentId, address proposed) returns()
func (_AgentRegistry *AgentRegistryTransactor) ProposeControllerTransfer(opts *bind.TransactOpts, agentId *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "proposeControllerTransfer", agentId, proposed)
}

// ProposeControllerTransfer is a paid mutator transaction binding the contract method 0x98310b5e.
//
// Solidity: function proposeControllerTransfer(uint256 agentId, address proposed) returns()
func (_AgentRegistry *AgentRegistrySession) ProposeControllerTransfer(agentId *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.ProposeControllerTransfer(&_AgentRegistry.TransactOpts, agentId, proposed)
}

// ProposeControllerTransfer is a paid mutator transaction binding the contract method 0x98310b5e.
//
// Solidity: function proposeControllerTransfer(uint256 agentId, address proposed) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) ProposeControllerTransfer(agentId *big.Int, proposed common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.ProposeControllerTransfer(&_AgentRegistry.TransactOpts, agentId, proposed)
}

// Register is a paid mutator transaction binding the contract method 0xfc0d1b84.
//
// Solidity: function register(address controller, string metadataURI, uint256 capabilities) returns(uint256 agentId)
func (_AgentRegistry *AgentRegistryTransactor) Register(opts *bind.TransactOpts, controller common.Address, metadataURI string, capabilities *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "register", controller, metadataURI, capabilities)
}

// Register is a paid mutator transaction binding the contract method 0xfc0d1b84.
//
// Solidity: function register(address controller, string metadataURI, uint256 capabilities) returns(uint256 agentId)
func (_AgentRegistry *AgentRegistrySession) Register(controller common.Address, metadataURI string, capabilities *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.Register(&_AgentRegistry.TransactOpts, controller, metadataURI, capabilities)
}

// Register is a paid mutator transaction binding the contract method 0xfc0d1b84.
//
// Solidity: function register(address controller, string metadataURI, uint256 capabilities) returns(uint256 agentId)
func (_AgentRegistry *AgentRegistryTransactorSession) Register(controller common.Address, metadataURI string, capabilities *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.Register(&_AgentRegistry.TransactOpts, controller, metadataURI, capabilities)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentRegistry *AgentRegistryTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentRegistry *AgentRegistrySession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RenounceRole(&_AgentRegistry.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RenounceRole(&_AgentRegistry.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentRegistry *AgentRegistryTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentRegistry *AgentRegistrySession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RevokeRole(&_AgentRegistry.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _AgentRegistry.Contract.RevokeRole(&_AgentRegistry.TransactOpts, role, account)
}

// SetActive is a paid mutator transaction binding the contract method 0xe60a955d.
//
// Solidity: function setActive(uint256 agentId, bool active) returns()
func (_AgentRegistry *AgentRegistryTransactor) SetActive(opts *bind.TransactOpts, agentId *big.Int, active bool) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "setActive", agentId, active)
}

// SetActive is a paid mutator transaction binding the contract method 0xe60a955d.
//
// Solidity: function setActive(uint256 agentId, bool active) returns()
func (_AgentRegistry *AgentRegistrySession) SetActive(agentId *big.Int, active bool) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetActive(&_AgentRegistry.TransactOpts, agentId, active)
}

// SetActive is a paid mutator transaction binding the contract method 0xe60a955d.
//
// Solidity: function setActive(uint256 agentId, bool active) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) SetActive(agentId *big.Int, active bool) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetActive(&_AgentRegistry.TransactOpts, agentId, active)
}

// SetCapabilities is a paid mutator transaction binding the contract method 0xb3473f83.
//
// Solidity: function setCapabilities(uint256 agentId, uint256 capabilities) returns()
func (_AgentRegistry *AgentRegistryTransactor) SetCapabilities(opts *bind.TransactOpts, agentId *big.Int, capabilities *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "setCapabilities", agentId, capabilities)
}

// SetCapabilities is a paid mutator transaction binding the contract method 0xb3473f83.
//
// Solidity: function setCapabilities(uint256 agentId, uint256 capabilities) returns()
func (_AgentRegistry *AgentRegistrySession) SetCapabilities(agentId *big.Int, capabilities *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetCapabilities(&_AgentRegistry.TransactOpts, agentId, capabilities)
}

// SetCapabilities is a paid mutator transaction binding the contract method 0xb3473f83.
//
// Solidity: function setCapabilities(uint256 agentId, uint256 capabilities) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) SetCapabilities(agentId *big.Int, capabilities *big.Int) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetCapabilities(&_AgentRegistry.TransactOpts, agentId, capabilities)
}

// SetMetadata is a paid mutator transaction binding the contract method 0x593aa283.
//
// Solidity: function setMetadata(uint256 agentId, string metadataURI) returns()
func (_AgentRegistry *AgentRegistryTransactor) SetMetadata(opts *bind.TransactOpts, agentId *big.Int, metadataURI string) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "setMetadata", agentId, metadataURI)
}

// SetMetadata is a paid mutator transaction binding the contract method 0x593aa283.
//
// Solidity: function setMetadata(uint256 agentId, string metadataURI) returns()
func (_AgentRegistry *AgentRegistrySession) SetMetadata(agentId *big.Int, metadataURI string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetMetadata(&_AgentRegistry.TransactOpts, agentId, metadataURI)
}

// SetMetadata is a paid mutator transaction binding the contract method 0x593aa283.
//
// Solidity: function setMetadata(uint256 agentId, string metadataURI) returns()
func (_AgentRegistry *AgentRegistryTransactorSession) SetMetadata(agentId *big.Int, metadataURI string) (*types.Transaction, error) {
	return _AgentRegistry.Contract.SetMetadata(&_AgentRegistry.TransactOpts, agentId, metadataURI)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AgentRegistry *AgentRegistryTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentRegistry.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AgentRegistry *AgentRegistrySession) Unpause() (*types.Transaction, error) {
	return _AgentRegistry.Contract.Unpause(&_AgentRegistry.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_AgentRegistry *AgentRegistryTransactorSession) Unpause() (*types.Transaction, error) {
	return _AgentRegistry.Contract.Unpause(&_AgentRegistry.TransactOpts)
}

// AgentRegistryActiveStatusChangedIterator is returned from FilterActiveStatusChanged and is used to iterate over the raw logs and unpacked data for ActiveStatusChanged events raised by the AgentRegistry contract.
type AgentRegistryActiveStatusChangedIterator struct {
	Event *AgentRegistryActiveStatusChanged // Event containing the contract specifics and raw log

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
func (it *AgentRegistryActiveStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryActiveStatusChanged)
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
		it.Event = new(AgentRegistryActiveStatusChanged)
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
func (it *AgentRegistryActiveStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryActiveStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryActiveStatusChanged represents a ActiveStatusChanged event raised by the AgentRegistry contract.
type AgentRegistryActiveStatusChanged struct {
	AgentId *big.Int
	Active  bool
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterActiveStatusChanged is a free log retrieval operation binding the contract event 0x632239944907e2f939f78bb06a2d210ac62cc0fbe3bff1b6244b4883721df129.
//
// Solidity: event ActiveStatusChanged(uint256 indexed agentId, bool active)
func (_AgentRegistry *AgentRegistryFilterer) FilterActiveStatusChanged(opts *bind.FilterOpts, agentId []*big.Int) (*AgentRegistryActiveStatusChangedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "ActiveStatusChanged", agentIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryActiveStatusChangedIterator{contract: _AgentRegistry.contract, event: "ActiveStatusChanged", logs: logs, sub: sub}, nil
}

// WatchActiveStatusChanged is a free log subscription operation binding the contract event 0x632239944907e2f939f78bb06a2d210ac62cc0fbe3bff1b6244b4883721df129.
//
// Solidity: event ActiveStatusChanged(uint256 indexed agentId, bool active)
func (_AgentRegistry *AgentRegistryFilterer) WatchActiveStatusChanged(opts *bind.WatchOpts, sink chan<- *AgentRegistryActiveStatusChanged, agentId []*big.Int) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "ActiveStatusChanged", agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryActiveStatusChanged)
				if err := _AgentRegistry.contract.UnpackLog(event, "ActiveStatusChanged", log); err != nil {
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

// ParseActiveStatusChanged is a log parse operation binding the contract event 0x632239944907e2f939f78bb06a2d210ac62cc0fbe3bff1b6244b4883721df129.
//
// Solidity: event ActiveStatusChanged(uint256 indexed agentId, bool active)
func (_AgentRegistry *AgentRegistryFilterer) ParseActiveStatusChanged(log types.Log) (*AgentRegistryActiveStatusChanged, error) {
	event := new(AgentRegistryActiveStatusChanged)
	if err := _AgentRegistry.contract.UnpackLog(event, "ActiveStatusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryAgentRegisteredIterator is returned from FilterAgentRegistered and is used to iterate over the raw logs and unpacked data for AgentRegistered events raised by the AgentRegistry contract.
type AgentRegistryAgentRegisteredIterator struct {
	Event *AgentRegistryAgentRegistered // Event containing the contract specifics and raw log

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
func (it *AgentRegistryAgentRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryAgentRegistered)
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
		it.Event = new(AgentRegistryAgentRegistered)
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
func (it *AgentRegistryAgentRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryAgentRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryAgentRegistered represents a AgentRegistered event raised by the AgentRegistry contract.
type AgentRegistryAgentRegistered struct {
	AgentId      *big.Int
	Controller   common.Address
	MetadataURI  string
	Capabilities *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAgentRegistered is a free log retrieval operation binding the contract event 0xc29f819ac362ff9c94de06666235808451aafd8894a2dffb86a080a965efeae3.
//
// Solidity: event AgentRegistered(uint256 indexed agentId, address indexed controller, string metadataURI, uint256 capabilities)
func (_AgentRegistry *AgentRegistryFilterer) FilterAgentRegistered(opts *bind.FilterOpts, agentId []*big.Int, controller []common.Address) (*AgentRegistryAgentRegisteredIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "AgentRegistered", agentIdRule, controllerRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryAgentRegisteredIterator{contract: _AgentRegistry.contract, event: "AgentRegistered", logs: logs, sub: sub}, nil
}

// WatchAgentRegistered is a free log subscription operation binding the contract event 0xc29f819ac362ff9c94de06666235808451aafd8894a2dffb86a080a965efeae3.
//
// Solidity: event AgentRegistered(uint256 indexed agentId, address indexed controller, string metadataURI, uint256 capabilities)
func (_AgentRegistry *AgentRegistryFilterer) WatchAgentRegistered(opts *bind.WatchOpts, sink chan<- *AgentRegistryAgentRegistered, agentId []*big.Int, controller []common.Address) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var controllerRule []interface{}
	for _, controllerItem := range controller {
		controllerRule = append(controllerRule, controllerItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "AgentRegistered", agentIdRule, controllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryAgentRegistered)
				if err := _AgentRegistry.contract.UnpackLog(event, "AgentRegistered", log); err != nil {
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

// ParseAgentRegistered is a log parse operation binding the contract event 0xc29f819ac362ff9c94de06666235808451aafd8894a2dffb86a080a965efeae3.
//
// Solidity: event AgentRegistered(uint256 indexed agentId, address indexed controller, string metadataURI, uint256 capabilities)
func (_AgentRegistry *AgentRegistryFilterer) ParseAgentRegistered(log types.Log) (*AgentRegistryAgentRegistered, error) {
	event := new(AgentRegistryAgentRegistered)
	if err := _AgentRegistry.contract.UnpackLog(event, "AgentRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryCapabilitiesUpdatedIterator is returned from FilterCapabilitiesUpdated and is used to iterate over the raw logs and unpacked data for CapabilitiesUpdated events raised by the AgentRegistry contract.
type AgentRegistryCapabilitiesUpdatedIterator struct {
	Event *AgentRegistryCapabilitiesUpdated // Event containing the contract specifics and raw log

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
func (it *AgentRegistryCapabilitiesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryCapabilitiesUpdated)
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
		it.Event = new(AgentRegistryCapabilitiesUpdated)
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
func (it *AgentRegistryCapabilitiesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryCapabilitiesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryCapabilitiesUpdated represents a CapabilitiesUpdated event raised by the AgentRegistry contract.
type AgentRegistryCapabilitiesUpdated struct {
	AgentId      *big.Int
	Capabilities *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterCapabilitiesUpdated is a free log retrieval operation binding the contract event 0x2852a2d71dd37e385f378b5a9e262f9d2ed0cf51b7312a5f9229cc5cd39ea4c1.
//
// Solidity: event CapabilitiesUpdated(uint256 indexed agentId, uint256 capabilities)
func (_AgentRegistry *AgentRegistryFilterer) FilterCapabilitiesUpdated(opts *bind.FilterOpts, agentId []*big.Int) (*AgentRegistryCapabilitiesUpdatedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "CapabilitiesUpdated", agentIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryCapabilitiesUpdatedIterator{contract: _AgentRegistry.contract, event: "CapabilitiesUpdated", logs: logs, sub: sub}, nil
}

// WatchCapabilitiesUpdated is a free log subscription operation binding the contract event 0x2852a2d71dd37e385f378b5a9e262f9d2ed0cf51b7312a5f9229cc5cd39ea4c1.
//
// Solidity: event CapabilitiesUpdated(uint256 indexed agentId, uint256 capabilities)
func (_AgentRegistry *AgentRegistryFilterer) WatchCapabilitiesUpdated(opts *bind.WatchOpts, sink chan<- *AgentRegistryCapabilitiesUpdated, agentId []*big.Int) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "CapabilitiesUpdated", agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryCapabilitiesUpdated)
				if err := _AgentRegistry.contract.UnpackLog(event, "CapabilitiesUpdated", log); err != nil {
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

// ParseCapabilitiesUpdated is a log parse operation binding the contract event 0x2852a2d71dd37e385f378b5a9e262f9d2ed0cf51b7312a5f9229cc5cd39ea4c1.
//
// Solidity: event CapabilitiesUpdated(uint256 indexed agentId, uint256 capabilities)
func (_AgentRegistry *AgentRegistryFilterer) ParseCapabilitiesUpdated(log types.Log) (*AgentRegistryCapabilitiesUpdated, error) {
	event := new(AgentRegistryCapabilitiesUpdated)
	if err := _AgentRegistry.contract.UnpackLog(event, "CapabilitiesUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryControllerTransferAcceptedIterator is returned from FilterControllerTransferAccepted and is used to iterate over the raw logs and unpacked data for ControllerTransferAccepted events raised by the AgentRegistry contract.
type AgentRegistryControllerTransferAcceptedIterator struct {
	Event *AgentRegistryControllerTransferAccepted // Event containing the contract specifics and raw log

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
func (it *AgentRegistryControllerTransferAcceptedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryControllerTransferAccepted)
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
		it.Event = new(AgentRegistryControllerTransferAccepted)
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
func (it *AgentRegistryControllerTransferAcceptedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryControllerTransferAcceptedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryControllerTransferAccepted represents a ControllerTransferAccepted event raised by the AgentRegistry contract.
type AgentRegistryControllerTransferAccepted struct {
	AgentId       *big.Int
	OldController common.Address
	NewController common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterControllerTransferAccepted is a free log retrieval operation binding the contract event 0x969d6b1394c5147c28489f50960712e680e34e1e14ea24ff72bd52f315930c3a.
//
// Solidity: event ControllerTransferAccepted(uint256 indexed agentId, address indexed oldController, address indexed newController)
func (_AgentRegistry *AgentRegistryFilterer) FilterControllerTransferAccepted(opts *bind.FilterOpts, agentId []*big.Int, oldController []common.Address, newController []common.Address) (*AgentRegistryControllerTransferAcceptedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var oldControllerRule []interface{}
	for _, oldControllerItem := range oldController {
		oldControllerRule = append(oldControllerRule, oldControllerItem)
	}
	var newControllerRule []interface{}
	for _, newControllerItem := range newController {
		newControllerRule = append(newControllerRule, newControllerItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "ControllerTransferAccepted", agentIdRule, oldControllerRule, newControllerRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryControllerTransferAcceptedIterator{contract: _AgentRegistry.contract, event: "ControllerTransferAccepted", logs: logs, sub: sub}, nil
}

// WatchControllerTransferAccepted is a free log subscription operation binding the contract event 0x969d6b1394c5147c28489f50960712e680e34e1e14ea24ff72bd52f315930c3a.
//
// Solidity: event ControllerTransferAccepted(uint256 indexed agentId, address indexed oldController, address indexed newController)
func (_AgentRegistry *AgentRegistryFilterer) WatchControllerTransferAccepted(opts *bind.WatchOpts, sink chan<- *AgentRegistryControllerTransferAccepted, agentId []*big.Int, oldController []common.Address, newController []common.Address) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var oldControllerRule []interface{}
	for _, oldControllerItem := range oldController {
		oldControllerRule = append(oldControllerRule, oldControllerItem)
	}
	var newControllerRule []interface{}
	for _, newControllerItem := range newController {
		newControllerRule = append(newControllerRule, newControllerItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "ControllerTransferAccepted", agentIdRule, oldControllerRule, newControllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryControllerTransferAccepted)
				if err := _AgentRegistry.contract.UnpackLog(event, "ControllerTransferAccepted", log); err != nil {
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

// ParseControllerTransferAccepted is a log parse operation binding the contract event 0x969d6b1394c5147c28489f50960712e680e34e1e14ea24ff72bd52f315930c3a.
//
// Solidity: event ControllerTransferAccepted(uint256 indexed agentId, address indexed oldController, address indexed newController)
func (_AgentRegistry *AgentRegistryFilterer) ParseControllerTransferAccepted(log types.Log) (*AgentRegistryControllerTransferAccepted, error) {
	event := new(AgentRegistryControllerTransferAccepted)
	if err := _AgentRegistry.contract.UnpackLog(event, "ControllerTransferAccepted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryControllerTransferCancelledIterator is returned from FilterControllerTransferCancelled and is used to iterate over the raw logs and unpacked data for ControllerTransferCancelled events raised by the AgentRegistry contract.
type AgentRegistryControllerTransferCancelledIterator struct {
	Event *AgentRegistryControllerTransferCancelled // Event containing the contract specifics and raw log

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
func (it *AgentRegistryControllerTransferCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryControllerTransferCancelled)
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
		it.Event = new(AgentRegistryControllerTransferCancelled)
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
func (it *AgentRegistryControllerTransferCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryControllerTransferCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryControllerTransferCancelled represents a ControllerTransferCancelled event raised by the AgentRegistry contract.
type AgentRegistryControllerTransferCancelled struct {
	AgentId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterControllerTransferCancelled is a free log retrieval operation binding the contract event 0x5d4293db10304d02137eb59859f20bef23395587c83ae50ada99be07f7aa9aad.
//
// Solidity: event ControllerTransferCancelled(uint256 indexed agentId)
func (_AgentRegistry *AgentRegistryFilterer) FilterControllerTransferCancelled(opts *bind.FilterOpts, agentId []*big.Int) (*AgentRegistryControllerTransferCancelledIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "ControllerTransferCancelled", agentIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryControllerTransferCancelledIterator{contract: _AgentRegistry.contract, event: "ControllerTransferCancelled", logs: logs, sub: sub}, nil
}

// WatchControllerTransferCancelled is a free log subscription operation binding the contract event 0x5d4293db10304d02137eb59859f20bef23395587c83ae50ada99be07f7aa9aad.
//
// Solidity: event ControllerTransferCancelled(uint256 indexed agentId)
func (_AgentRegistry *AgentRegistryFilterer) WatchControllerTransferCancelled(opts *bind.WatchOpts, sink chan<- *AgentRegistryControllerTransferCancelled, agentId []*big.Int) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "ControllerTransferCancelled", agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryControllerTransferCancelled)
				if err := _AgentRegistry.contract.UnpackLog(event, "ControllerTransferCancelled", log); err != nil {
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

// ParseControllerTransferCancelled is a log parse operation binding the contract event 0x5d4293db10304d02137eb59859f20bef23395587c83ae50ada99be07f7aa9aad.
//
// Solidity: event ControllerTransferCancelled(uint256 indexed agentId)
func (_AgentRegistry *AgentRegistryFilterer) ParseControllerTransferCancelled(log types.Log) (*AgentRegistryControllerTransferCancelled, error) {
	event := new(AgentRegistryControllerTransferCancelled)
	if err := _AgentRegistry.contract.UnpackLog(event, "ControllerTransferCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryControllerTransferProposedIterator is returned from FilterControllerTransferProposed and is used to iterate over the raw logs and unpacked data for ControllerTransferProposed events raised by the AgentRegistry contract.
type AgentRegistryControllerTransferProposedIterator struct {
	Event *AgentRegistryControllerTransferProposed // Event containing the contract specifics and raw log

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
func (it *AgentRegistryControllerTransferProposedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryControllerTransferProposed)
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
		it.Event = new(AgentRegistryControllerTransferProposed)
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
func (it *AgentRegistryControllerTransferProposedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryControllerTransferProposedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryControllerTransferProposed represents a ControllerTransferProposed event raised by the AgentRegistry contract.
type AgentRegistryControllerTransferProposed struct {
	AgentId  *big.Int
	Proposed common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterControllerTransferProposed is a free log retrieval operation binding the contract event 0x41b8c1345e7590a4eb37c55b8c1cef180d9d1c694be1ae96d015bfe11a278e22.
//
// Solidity: event ControllerTransferProposed(uint256 indexed agentId, address indexed proposed)
func (_AgentRegistry *AgentRegistryFilterer) FilterControllerTransferProposed(opts *bind.FilterOpts, agentId []*big.Int, proposed []common.Address) (*AgentRegistryControllerTransferProposedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "ControllerTransferProposed", agentIdRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryControllerTransferProposedIterator{contract: _AgentRegistry.contract, event: "ControllerTransferProposed", logs: logs, sub: sub}, nil
}

// WatchControllerTransferProposed is a free log subscription operation binding the contract event 0x41b8c1345e7590a4eb37c55b8c1cef180d9d1c694be1ae96d015bfe11a278e22.
//
// Solidity: event ControllerTransferProposed(uint256 indexed agentId, address indexed proposed)
func (_AgentRegistry *AgentRegistryFilterer) WatchControllerTransferProposed(opts *bind.WatchOpts, sink chan<- *AgentRegistryControllerTransferProposed, agentId []*big.Int, proposed []common.Address) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var proposedRule []interface{}
	for _, proposedItem := range proposed {
		proposedRule = append(proposedRule, proposedItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "ControllerTransferProposed", agentIdRule, proposedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryControllerTransferProposed)
				if err := _AgentRegistry.contract.UnpackLog(event, "ControllerTransferProposed", log); err != nil {
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

// ParseControllerTransferProposed is a log parse operation binding the contract event 0x41b8c1345e7590a4eb37c55b8c1cef180d9d1c694be1ae96d015bfe11a278e22.
//
// Solidity: event ControllerTransferProposed(uint256 indexed agentId, address indexed proposed)
func (_AgentRegistry *AgentRegistryFilterer) ParseControllerTransferProposed(log types.Log) (*AgentRegistryControllerTransferProposed, error) {
	event := new(AgentRegistryControllerTransferProposed)
	if err := _AgentRegistry.contract.UnpackLog(event, "ControllerTransferProposed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryMetadataUpdatedIterator is returned from FilterMetadataUpdated and is used to iterate over the raw logs and unpacked data for MetadataUpdated events raised by the AgentRegistry contract.
type AgentRegistryMetadataUpdatedIterator struct {
	Event *AgentRegistryMetadataUpdated // Event containing the contract specifics and raw log

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
func (it *AgentRegistryMetadataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryMetadataUpdated)
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
		it.Event = new(AgentRegistryMetadataUpdated)
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
func (it *AgentRegistryMetadataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryMetadataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryMetadataUpdated represents a MetadataUpdated event raised by the AgentRegistry contract.
type AgentRegistryMetadataUpdated struct {
	AgentId     *big.Int
	MetadataURI string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdated is a free log retrieval operation binding the contract event 0x459157ba24c7ab9878b165ef465fa6ae2ab42bcd8445f576be378768b0c47309.
//
// Solidity: event MetadataUpdated(uint256 indexed agentId, string metadataURI)
func (_AgentRegistry *AgentRegistryFilterer) FilterMetadataUpdated(opts *bind.FilterOpts, agentId []*big.Int) (*AgentRegistryMetadataUpdatedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "MetadataUpdated", agentIdRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryMetadataUpdatedIterator{contract: _AgentRegistry.contract, event: "MetadataUpdated", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdated is a free log subscription operation binding the contract event 0x459157ba24c7ab9878b165ef465fa6ae2ab42bcd8445f576be378768b0c47309.
//
// Solidity: event MetadataUpdated(uint256 indexed agentId, string metadataURI)
func (_AgentRegistry *AgentRegistryFilterer) WatchMetadataUpdated(opts *bind.WatchOpts, sink chan<- *AgentRegistryMetadataUpdated, agentId []*big.Int) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "MetadataUpdated", agentIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryMetadataUpdated)
				if err := _AgentRegistry.contract.UnpackLog(event, "MetadataUpdated", log); err != nil {
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

// ParseMetadataUpdated is a log parse operation binding the contract event 0x459157ba24c7ab9878b165ef465fa6ae2ab42bcd8445f576be378768b0c47309.
//
// Solidity: event MetadataUpdated(uint256 indexed agentId, string metadataURI)
func (_AgentRegistry *AgentRegistryFilterer) ParseMetadataUpdated(log types.Log) (*AgentRegistryMetadataUpdated, error) {
	event := new(AgentRegistryMetadataUpdated)
	if err := _AgentRegistry.contract.UnpackLog(event, "MetadataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the AgentRegistry contract.
type AgentRegistryPausedIterator struct {
	Event *AgentRegistryPaused // Event containing the contract specifics and raw log

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
func (it *AgentRegistryPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryPaused)
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
		it.Event = new(AgentRegistryPaused)
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
func (it *AgentRegistryPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryPaused represents a Paused event raised by the AgentRegistry contract.
type AgentRegistryPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_AgentRegistry *AgentRegistryFilterer) FilterPaused(opts *bind.FilterOpts) (*AgentRegistryPausedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryPausedIterator{contract: _AgentRegistry.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_AgentRegistry *AgentRegistryFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *AgentRegistryPaused) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryPaused)
				if err := _AgentRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParsePaused(log types.Log) (*AgentRegistryPaused, error) {
	event := new(AgentRegistryPaused)
	if err := _AgentRegistry.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the AgentRegistry contract.
type AgentRegistryRoleAdminChangedIterator struct {
	Event *AgentRegistryRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *AgentRegistryRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryRoleAdminChanged)
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
		it.Event = new(AgentRegistryRoleAdminChanged)
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
func (it *AgentRegistryRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryRoleAdminChanged represents a RoleAdminChanged event raised by the AgentRegistry contract.
type AgentRegistryRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgentRegistry *AgentRegistryFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*AgentRegistryRoleAdminChangedIterator, error) {

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

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryRoleAdminChangedIterator{contract: _AgentRegistry.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_AgentRegistry *AgentRegistryFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *AgentRegistryRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryRoleAdminChanged)
				if err := _AgentRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseRoleAdminChanged(log types.Log) (*AgentRegistryRoleAdminChanged, error) {
	event := new(AgentRegistryRoleAdminChanged)
	if err := _AgentRegistry.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the AgentRegistry contract.
type AgentRegistryRoleGrantedIterator struct {
	Event *AgentRegistryRoleGranted // Event containing the contract specifics and raw log

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
func (it *AgentRegistryRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryRoleGranted)
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
		it.Event = new(AgentRegistryRoleGranted)
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
func (it *AgentRegistryRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryRoleGranted represents a RoleGranted event raised by the AgentRegistry contract.
type AgentRegistryRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentRegistry *AgentRegistryFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AgentRegistryRoleGrantedIterator, error) {

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

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryRoleGrantedIterator{contract: _AgentRegistry.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentRegistry *AgentRegistryFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *AgentRegistryRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryRoleGranted)
				if err := _AgentRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseRoleGranted(log types.Log) (*AgentRegistryRoleGranted, error) {
	event := new(AgentRegistryRoleGranted)
	if err := _AgentRegistry.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the AgentRegistry contract.
type AgentRegistryRoleRevokedIterator struct {
	Event *AgentRegistryRoleRevoked // Event containing the contract specifics and raw log

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
func (it *AgentRegistryRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryRoleRevoked)
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
		it.Event = new(AgentRegistryRoleRevoked)
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
func (it *AgentRegistryRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryRoleRevoked represents a RoleRevoked event raised by the AgentRegistry contract.
type AgentRegistryRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentRegistry *AgentRegistryFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AgentRegistryRoleRevokedIterator, error) {

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

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryRoleRevokedIterator{contract: _AgentRegistry.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_AgentRegistry *AgentRegistryFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *AgentRegistryRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryRoleRevoked)
				if err := _AgentRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseRoleRevoked(log types.Log) (*AgentRegistryRoleRevoked, error) {
	event := new(AgentRegistryRoleRevoked)
	if err := _AgentRegistry.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryScorePostedIterator is returned from FilterScorePosted and is used to iterate over the raw logs and unpacked data for ScorePosted events raised by the AgentRegistry contract.
type AgentRegistryScorePostedIterator struct {
	Event *AgentRegistryScorePosted // Event containing the contract specifics and raw log

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
func (it *AgentRegistryScorePostedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryScorePosted)
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
		it.Event = new(AgentRegistryScorePosted)
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
func (it *AgentRegistryScorePostedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryScorePostedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryScorePosted represents a ScorePosted event raised by the AgentRegistry contract.
type AgentRegistryScorePosted struct {
	AgentId  *big.Int
	Scorer   common.Address
	Delta    *big.Int
	NewTotal *big.Int
	Nonce    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterScorePosted is a free log retrieval operation binding the contract event 0xa19792167f6e5259706dc8bacefe76b932fe47d8538e98761711b78e49129e3a.
//
// Solidity: event ScorePosted(uint256 indexed agentId, address indexed scorer, int256 delta, int256 newTotal, uint256 nonce)
func (_AgentRegistry *AgentRegistryFilterer) FilterScorePosted(opts *bind.FilterOpts, agentId []*big.Int, scorer []common.Address) (*AgentRegistryScorePostedIterator, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var scorerRule []interface{}
	for _, scorerItem := range scorer {
		scorerRule = append(scorerRule, scorerItem)
	}

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "ScorePosted", agentIdRule, scorerRule)
	if err != nil {
		return nil, err
	}
	return &AgentRegistryScorePostedIterator{contract: _AgentRegistry.contract, event: "ScorePosted", logs: logs, sub: sub}, nil
}

// WatchScorePosted is a free log subscription operation binding the contract event 0xa19792167f6e5259706dc8bacefe76b932fe47d8538e98761711b78e49129e3a.
//
// Solidity: event ScorePosted(uint256 indexed agentId, address indexed scorer, int256 delta, int256 newTotal, uint256 nonce)
func (_AgentRegistry *AgentRegistryFilterer) WatchScorePosted(opts *bind.WatchOpts, sink chan<- *AgentRegistryScorePosted, agentId []*big.Int, scorer []common.Address) (event.Subscription, error) {

	var agentIdRule []interface{}
	for _, agentIdItem := range agentId {
		agentIdRule = append(agentIdRule, agentIdItem)
	}
	var scorerRule []interface{}
	for _, scorerItem := range scorer {
		scorerRule = append(scorerRule, scorerItem)
	}

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "ScorePosted", agentIdRule, scorerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryScorePosted)
				if err := _AgentRegistry.contract.UnpackLog(event, "ScorePosted", log); err != nil {
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

// ParseScorePosted is a log parse operation binding the contract event 0xa19792167f6e5259706dc8bacefe76b932fe47d8538e98761711b78e49129e3a.
//
// Solidity: event ScorePosted(uint256 indexed agentId, address indexed scorer, int256 delta, int256 newTotal, uint256 nonce)
func (_AgentRegistry *AgentRegistryFilterer) ParseScorePosted(log types.Log) (*AgentRegistryScorePosted, error) {
	event := new(AgentRegistryScorePosted)
	if err := _AgentRegistry.contract.UnpackLog(event, "ScorePosted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentRegistryUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the AgentRegistry contract.
type AgentRegistryUnpausedIterator struct {
	Event *AgentRegistryUnpaused // Event containing the contract specifics and raw log

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
func (it *AgentRegistryUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentRegistryUnpaused)
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
		it.Event = new(AgentRegistryUnpaused)
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
func (it *AgentRegistryUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentRegistryUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentRegistryUnpaused represents a Unpaused event raised by the AgentRegistry contract.
type AgentRegistryUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_AgentRegistry *AgentRegistryFilterer) FilterUnpaused(opts *bind.FilterOpts) (*AgentRegistryUnpausedIterator, error) {

	logs, sub, err := _AgentRegistry.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &AgentRegistryUnpausedIterator{contract: _AgentRegistry.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_AgentRegistry *AgentRegistryFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *AgentRegistryUnpaused) (event.Subscription, error) {

	logs, sub, err := _AgentRegistry.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentRegistryUnpaused)
				if err := _AgentRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_AgentRegistry *AgentRegistryFilterer) ParseUnpaused(log types.Log) (*AgentRegistryUnpaused, error) {
	event := new(AgentRegistryUnpaused)
	if err := _AgentRegistry.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
