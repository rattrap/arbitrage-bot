// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// ITickLensPopulatedTick is an auto generated low-level Go binding around an user-defined struct.
type ITickLensPopulatedTick struct {
	Tick           *big.Int
	LiquidityNet   *big.Int
	LiquidityGross *big.Int
}

// TickLensMetaData contains all meta data concerning the TickLens contract.
var TickLensMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"pool\",\"type\":\"address\"},{\"internalType\":\"int16\",\"name\":\"tickBitmapIndex\",\"type\":\"int16\"}],\"name\":\"getPopulatedTicksInWord\",\"outputs\":[{\"components\":[{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"},{\"internalType\":\"int128\",\"name\":\"liquidityNet\",\"type\":\"int128\"},{\"internalType\":\"uint128\",\"name\":\"liquidityGross\",\"type\":\"uint128\"}],\"internalType\":\"structITickLens.PopulatedTick[]\",\"name\":\"populatedTicks\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TickLensABI is the input ABI used to generate the binding from.
// Deprecated: Use TickLensMetaData.ABI instead.
var TickLensABI = TickLensMetaData.ABI

// TickLens is an auto generated Go binding around an Ethereum contract.
type TickLens struct {
	TickLensCaller     // Read-only binding to the contract
	TickLensTransactor // Write-only binding to the contract
	TickLensFilterer   // Log filterer for contract events
}

// TickLensCaller is an auto generated read-only Go binding around an Ethereum contract.
type TickLensCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TickLensTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TickLensTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TickLensFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TickLensFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TickLensSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TickLensSession struct {
	Contract     *TickLens         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TickLensCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TickLensCallerSession struct {
	Contract *TickLensCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TickLensTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TickLensTransactorSession struct {
	Contract     *TickLensTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TickLensRaw is an auto generated low-level Go binding around an Ethereum contract.
type TickLensRaw struct {
	Contract *TickLens // Generic contract binding to access the raw methods on
}

// TickLensCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TickLensCallerRaw struct {
	Contract *TickLensCaller // Generic read-only contract binding to access the raw methods on
}

// TickLensTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TickLensTransactorRaw struct {
	Contract *TickLensTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTickLens creates a new instance of TickLens, bound to a specific deployed contract.
func NewTickLens(address common.Address, backend bind.ContractBackend) (*TickLens, error) {
	contract, err := bindTickLens(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TickLens{TickLensCaller: TickLensCaller{contract: contract}, TickLensTransactor: TickLensTransactor{contract: contract}, TickLensFilterer: TickLensFilterer{contract: contract}}, nil
}

// NewTickLensCaller creates a new read-only instance of TickLens, bound to a specific deployed contract.
func NewTickLensCaller(address common.Address, caller bind.ContractCaller) (*TickLensCaller, error) {
	contract, err := bindTickLens(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TickLensCaller{contract: contract}, nil
}

// NewTickLensTransactor creates a new write-only instance of TickLens, bound to a specific deployed contract.
func NewTickLensTransactor(address common.Address, transactor bind.ContractTransactor) (*TickLensTransactor, error) {
	contract, err := bindTickLens(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TickLensTransactor{contract: contract}, nil
}

// NewTickLensFilterer creates a new log filterer instance of TickLens, bound to a specific deployed contract.
func NewTickLensFilterer(address common.Address, filterer bind.ContractFilterer) (*TickLensFilterer, error) {
	contract, err := bindTickLens(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TickLensFilterer{contract: contract}, nil
}

// bindTickLens binds a generic wrapper to an already deployed contract.
func bindTickLens(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TickLensMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TickLens *TickLensRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TickLens.Contract.TickLensCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TickLens *TickLensRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TickLens.Contract.TickLensTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TickLens *TickLensRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TickLens.Contract.TickLensTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TickLens *TickLensCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TickLens.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TickLens *TickLensTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TickLens.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TickLens *TickLensTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TickLens.Contract.contract.Transact(opts, method, params...)
}

// GetPopulatedTicksInWord is a free data retrieval call binding the contract method 0x351fb478.
//
// Solidity: function getPopulatedTicksInWord(address pool, int16 tickBitmapIndex) view returns((int24,int128,uint128)[] populatedTicks)
func (_TickLens *TickLensCaller) GetPopulatedTicksInWord(opts *bind.CallOpts, pool common.Address, tickBitmapIndex int16) ([]ITickLensPopulatedTick, error) {
	var out []interface{}
	err := _TickLens.contract.Call(opts, &out, "getPopulatedTicksInWord", pool, tickBitmapIndex)

	if err != nil {
		return *new([]ITickLensPopulatedTick), err
	}

	out0 := *abi.ConvertType(out[0], new([]ITickLensPopulatedTick)).(*[]ITickLensPopulatedTick)

	return out0, err

}

// GetPopulatedTicksInWord is a free data retrieval call binding the contract method 0x351fb478.
//
// Solidity: function getPopulatedTicksInWord(address pool, int16 tickBitmapIndex) view returns((int24,int128,uint128)[] populatedTicks)
func (_TickLens *TickLensSession) GetPopulatedTicksInWord(pool common.Address, tickBitmapIndex int16) ([]ITickLensPopulatedTick, error) {
	return _TickLens.Contract.GetPopulatedTicksInWord(&_TickLens.CallOpts, pool, tickBitmapIndex)
}

// GetPopulatedTicksInWord is a free data retrieval call binding the contract method 0x351fb478.
//
// Solidity: function getPopulatedTicksInWord(address pool, int16 tickBitmapIndex) view returns((int24,int128,uint128)[] populatedTicks)
func (_TickLens *TickLensCallerSession) GetPopulatedTicksInWord(pool common.Address, tickBitmapIndex int16) ([]ITickLensPopulatedTick, error) {
	return _TickLens.Contract.GetPopulatedTicksInWord(&_TickLens.CallOpts, pool, tickBitmapIndex)
}
