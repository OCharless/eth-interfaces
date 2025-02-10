// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Disperse

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

// DisperseMetaData contains all meta data concerning the Disperse contract.
var DisperseMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"disperseEther\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"disperseToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"recipients\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"disperseTokenSimple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DisperseABI is the input ABI used to generate the binding from.
// Deprecated: Use DisperseMetaData.ABI instead.
var DisperseABI = DisperseMetaData.ABI

// Disperse is an auto generated Go binding around an Ethereum contract.
type Disperse struct {
	DisperseCaller     // Read-only binding to the contract
	DisperseTransactor // Write-only binding to the contract
	DisperseFilterer   // Log filterer for contract events
}

// DisperseCaller is an auto generated read-only Go binding around an Ethereum contract.
type DisperseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DisperseTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DisperseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DisperseFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DisperseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DisperseSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DisperseSession struct {
	Contract     *Disperse         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DisperseCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DisperseCallerSession struct {
	Contract *DisperseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// DisperseTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DisperseTransactorSession struct {
	Contract     *DisperseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DisperseRaw is an auto generated low-level Go binding around an Ethereum contract.
type DisperseRaw struct {
	Contract *Disperse // Generic contract binding to access the raw methods on
}

// DisperseCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DisperseCallerRaw struct {
	Contract *DisperseCaller // Generic read-only contract binding to access the raw methods on
}

// DisperseTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DisperseTransactorRaw struct {
	Contract *DisperseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDisperse creates a new instance of Disperse, bound to a specific deployed contract.
func NewDisperse(address common.Address, backend bind.ContractBackend) (*Disperse, error) {
	contract, err := bindDisperse(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Disperse{DisperseCaller: DisperseCaller{contract: contract}, DisperseTransactor: DisperseTransactor{contract: contract}, DisperseFilterer: DisperseFilterer{contract: contract}}, nil
}

// NewDisperseCaller creates a new read-only instance of Disperse, bound to a specific deployed contract.
func NewDisperseCaller(address common.Address, caller bind.ContractCaller) (*DisperseCaller, error) {
	contract, err := bindDisperse(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DisperseCaller{contract: contract}, nil
}

// NewDisperseTransactor creates a new write-only instance of Disperse, bound to a specific deployed contract.
func NewDisperseTransactor(address common.Address, transactor bind.ContractTransactor) (*DisperseTransactor, error) {
	contract, err := bindDisperse(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DisperseTransactor{contract: contract}, nil
}

// NewDisperseFilterer creates a new log filterer instance of Disperse, bound to a specific deployed contract.
func NewDisperseFilterer(address common.Address, filterer bind.ContractFilterer) (*DisperseFilterer, error) {
	contract, err := bindDisperse(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DisperseFilterer{contract: contract}, nil
}

// bindDisperse binds a generic wrapper to an already deployed contract.
func bindDisperse(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DisperseMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Disperse *DisperseRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Disperse.Contract.DisperseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Disperse *DisperseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Disperse.Contract.DisperseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Disperse *DisperseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Disperse.Contract.DisperseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Disperse *DisperseCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Disperse.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Disperse *DisperseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Disperse.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Disperse *DisperseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Disperse.Contract.contract.Transact(opts, method, params...)
}

// DisperseEther is a paid mutator transaction binding the contract method 0xe63d38ed.
//
// Solidity: function disperseEther(address[] recipients, uint256[] values) payable returns()
func (_Disperse *DisperseTransactor) DisperseEther(opts *bind.TransactOpts, recipients []common.Address, values []*big.Int) (*types.Transaction, error) {
	return _Disperse.contract.Transact(opts, "disperseEther", recipients, values)
}

// DisperseEther is a paid mutator transaction binding the contract method 0xe63d38ed.
//
// Solidity: function disperseEther(address[] recipients, uint256[] values) payable returns()
func (_Disperse *DisperseSession) DisperseEther(recipients []common.Address, values []*big.Int) (*types.Transaction, error) {
	return _Disperse.Contract.DisperseEther(&_Disperse.TransactOpts, recipients, values)
}

// DisperseEther is a paid mutator transaction binding the contract method 0xe63d38ed.
//
// Solidity: function disperseEther(address[] recipients, uint256[] values) payable returns()
func (_Disperse *DisperseTransactorSession) DisperseEther(recipients []common.Address, values []*big.Int) (*types.Transaction, error) {
	return _Disperse.Contract.DisperseEther(&_Disperse.TransactOpts, recipients, values)
}

// DisperseToken is a paid mutator transaction binding the contract method 0xc73a2d60.
//
// Solidity: function disperseToken(address token, address[] recipients, uint256[] values) returns()
func (_Disperse *DisperseTransactor) DisperseToken(opts *bind.TransactOpts, token common.Address, recipients []common.Address, values []*big.Int) (*types.Transaction, error) {
	return _Disperse.contract.Transact(opts, "disperseToken", token, recipients, values)
}

// DisperseToken is a paid mutator transaction binding the contract method 0xc73a2d60.
//
// Solidity: function disperseToken(address token, address[] recipients, uint256[] values) returns()
func (_Disperse *DisperseSession) DisperseToken(token common.Address, recipients []common.Address, values []*big.Int) (*types.Transaction, error) {
	return _Disperse.Contract.DisperseToken(&_Disperse.TransactOpts, token, recipients, values)
}

// DisperseToken is a paid mutator transaction binding the contract method 0xc73a2d60.
//
// Solidity: function disperseToken(address token, address[] recipients, uint256[] values) returns()
func (_Disperse *DisperseTransactorSession) DisperseToken(token common.Address, recipients []common.Address, values []*big.Int) (*types.Transaction, error) {
	return _Disperse.Contract.DisperseToken(&_Disperse.TransactOpts, token, recipients, values)
}

// DisperseTokenSimple is a paid mutator transaction binding the contract method 0x51ba162c.
//
// Solidity: function disperseTokenSimple(address token, address[] recipients, uint256[] values) returns()
func (_Disperse *DisperseTransactor) DisperseTokenSimple(opts *bind.TransactOpts, token common.Address, recipients []common.Address, values []*big.Int) (*types.Transaction, error) {
	return _Disperse.contract.Transact(opts, "disperseTokenSimple", token, recipients, values)
}

// DisperseTokenSimple is a paid mutator transaction binding the contract method 0x51ba162c.
//
// Solidity: function disperseTokenSimple(address token, address[] recipients, uint256[] values) returns()
func (_Disperse *DisperseSession) DisperseTokenSimple(token common.Address, recipients []common.Address, values []*big.Int) (*types.Transaction, error) {
	return _Disperse.Contract.DisperseTokenSimple(&_Disperse.TransactOpts, token, recipients, values)
}

// DisperseTokenSimple is a paid mutator transaction binding the contract method 0x51ba162c.
//
// Solidity: function disperseTokenSimple(address token, address[] recipients, uint256[] values) returns()
func (_Disperse *DisperseTransactorSession) DisperseTokenSimple(token common.Address, recipients []common.Address, values []*big.Int) (*types.Transaction, error) {
	return _Disperse.Contract.DisperseTokenSimple(&_Disperse.TransactOpts, token, recipients, values)
}
