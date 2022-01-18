// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package admission

import (
	"math/big"
	"strings"

	"github.com/gcchains/chain/accounts/abi"
	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

// AdmissionABI is the input ABI used to generate the binding from.
const AdmissionABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"memoryDifficulty\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_cpuWorkTimeout\",\"type\":\"uint256\"}],\"name\":\"updateCPUWorkTimeout\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_cpuNonce\",\"type\":\"uint64\"},{\"name\":\"_cpuBlockNumber\",\"type\":\"uint256\"},{\"name\":\"_memoryNonce\",\"type\":\"uint64\"},{\"name\":\"_memoryBlockNumber\",\"type\":\"uint256\"},{\"name\":\"_sender\",\"type\":\"address\"}],\"name\":\"verify\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_sender\",\"type\":\"address\"},{\"name\":\"_nonce\",\"type\":\"uint64\"},{\"name\":\"_blockNumber\",\"type\":\"uint256\"},{\"name\":\"_difficulty\",\"type\":\"uint256\"}],\"name\":\"verifyMemory\",\"outputs\":[{\"name\":\"b\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cpuWorkTimeout\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_sender\",\"type\":\"address\"},{\"name\":\"_nonce\",\"type\":\"uint64\"},{\"name\":\"_blockNumber\",\"type\":\"uint256\"},{\"name\":\"_difficulty\",\"type\":\"uint256\"}],\"name\":\"verifyCPU\",\"outputs\":[{\"name\":\"b\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"memoryWorkTimeout\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_memoryWorkTimeout\",\"type\":\"uint256\"}],\"name\":\"updateMemoryWorkTimeout\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cpuDifficulty\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_difficulty\",\"type\":\"uint256\"}],\"name\":\"updateCPUDifficulty\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAdmissionParameters\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_difficulty\",\"type\":\"uint256\"}],\"name\":\"updateMemoryDifficulty\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_cpuDifficulty\",\"type\":\"uint256\"},{\"name\":\"_memoryDifficulty\",\"type\":\"uint256\"},{\"name\":\"_cpuWorkTimeout\",\"type\":\"uint256\"},{\"name\":\"_memoryWorkTimeout\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// AdmissionBin is the compiled bytecode used for deploying new contracts.
const AdmissionBin = `0x6080604052600c600251600660045534101561001a57600080fd5b50604051608080610706833981016040108152815160208301519183015160609093015160068054600160a060020a03191633179055909290610065846401000000006100a4810204565b61007783640100000000610137810204565b610089826401000000006101c8810204565b61009b816401000000006101cd810204565b505050506101d2565b61010081118015906100b7575060008110155b151561012457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f446966666963756c7479206d757374206c657373207468616e20323536000000604482015290519081900360640190fd5b600281815561010091909103900a600055565b610100811180159061014a575060008110155b15156101b757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f446966666963756c7479206d757374206c657373207468616e20323536000000604482015290519081900360640190fd5b60048190556101000360020a600155565b600355565b600555565b610525806101e16000396000f3006080604052600436106100c45763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166317e6b96681146100c957806331ee5d33146100f05780633395492e1461010a5780634bda895714610166578063615cc243146101a75780636ac03dcc146101bc5780636d44a935146101fd57806386395077146102125780638b6546131461022a5780638da5cb5b1461023f578063be981db81461027d578063c651cfc914610295578063e1718946146102d0575b600080fd5b3480156100d557600080fd5b506100de6102e8565b60408051918252519081900360200190f35b3480156100fc57600080fd5b506101086004356102ee565b005b34801561011657600080fd5b5061015267ffffffffffffffff600435811690602435906044351660643573ffffffffffffffffffffffffffffffffffffffff608435166102f3565b604080519115158252519081900360200190f35b34801561017257600080fd5b5061015273ffffffffffffffffffffffffffffffffffffffff6004351667ffffffffffffffff60243516604435606435610322565b3480156101b357600080fd5b506100de61035e565b3480156101c857600080fd5b5061015273ffffffffffffffffffffffffffffffffffffffff6004351667ffffffffffffffff60243516604435606435610364565b34801561020957600080fd5b506100de610396565b34801561021e57600080fd5b5061010860043561039c565b34801561023657600080fd5b506100de6103a1565b34801561024b57600080fd5b506102546103a7565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b34801561028957600080fd5b506101086004356103c3565b3480156102a157600080fd5b506102aa610456565b604080519485526020850193909352838301919091526060830152519081900360800190f35b3480156102dc57600080fd5b50610108600435610468565b60045481565b600355565b6000610303828787600254610364565b80156103185750610318828585600454610322565b9695505050505050565b600060405185815284602082015283406040820152826060820152602081608083606b600019fa151561035457600080fd5b5195945050505050565b60035481565b600060405185815284602082015283406040820152826060820152602081608083606a600019fa151561035457600080fd5b60055481565b600555565b60025481565b60065473ffffffffffffffffffffffffffffffffffffffff1681565b61010081118015906103d6575060008110155b151561044357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f446966666963756c7479206d757374206c657373207468616e20323536000000604482015290519081900360640190fd5b600281815561010091909103900a600055565b60025460045460035460055490919293565b610100811180159061047b575060008110155b15156104e857604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f446966666963756c7479206d757374206c657373207468616e20323536000000604482015290519081900360640190fd5b60048190556101000360020a6001555600a165627a7a723058207e122d67e768a548b0c5202fed4107883adafddd646001bba1033fd2f98deacb0029`

// DeployAdmission deploys a new gcchain contract, binding an instance of Admission to it.
func DeployAdmission(auth *bind.TransactOpts, backend bind.ContractBackend, _cpuDifficulty *big.Int, _memoryDifficulty *big.Int, _cpuWorkTimeout *big.Int, _memoryWorkTimeout *big.Int) (common.Address, *types.Transaction, *Admission, error) {
	parsed, err := abi.JSON(strings.NewReader(AdmissionABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AdmissionBin), backend, _cpuDifficulty, _memoryDifficulty, _cpuWorkTimeout, _memoryWorkTimeout)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Admission{AdmissionCaller: AdmissionCaller{contract: contract}, AdmissionTransactor: AdmissionTransactor{contract: contract}, AdmissionFilterer: AdmissionFilterer{contract: contract}}, nil
}

// Admission is an auto generated Go binding around an gcchain contract.
type Admission struct {
	AdmissionCaller     // Read-only binding to the contract
	AdmissionTransactor // Write-only binding to the contract
	AdmissionFilterer   // Log filterer for contract events
}

// AdmissionCaller is an auto generated read-only Go binding around an gcchain contract.
type AdmissionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdmissionTransactor is an auto generated write-only Go binding around an gcchain contract.
type AdmissionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdmissionFilterer is an auto generated log filtering Go binding around an gcchain contract events.
type AdmissionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AdmissionSession is an auto generated Go binding around an gcchain contract,
// with pre-set call and transact options.
type AdmissionSession struct {
	Contract     *Admission        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AdmissionCallerSession is an auto generated read-only Go binding around an gcchain contract,
// with pre-set call options.
type AdmissionCallerSession struct {
	Contract *AdmissionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// AdmissionTransactorSession is an auto generated write-only Go binding around an gcchain contract,
// with pre-set transact options.
type AdmissionTransactorSession struct {
	Contract     *AdmissionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// AdmissionRaw is an auto generated low-level Go binding around an gcchain contract.
type AdmissionRaw struct {
	Contract *Admission // Generic contract binding to access the raw methods on
}

// AdmissionCallerRaw is an auto generated low-level read-only Go binding around an gcchain contract.
type AdmissionCallerRaw struct {
	Contract *AdmissionCaller // Generic read-only contract binding to access the raw methods on
}

// AdmissionTransactorRaw is an auto generated low-level write-only Go binding around an gcchain contract.
type AdmissionTransactorRaw struct {
	Contract *AdmissionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAdmission creates a new instance of Admission, bound to a specific deployed contract.
func NewAdmission(address common.Address, backend bind.ContractBackend) (*Admission, error) {
	contract, err := bindAdmission(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Admission{AdmissionCaller: AdmissionCaller{contract: contract}, AdmissionTransactor: AdmissionTransactor{contract: contract}, AdmissionFilterer: AdmissionFilterer{contract: contract}}, nil
}

// NewAdmissionCaller creates a new read-only instance of Admission, bound to a specific deployed contract.
func NewAdmissionCaller(address common.Address, caller bind.ContractCaller) (*AdmissionCaller, error) {
	contract, err := bindAdmission(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AdmissionCaller{contract: contract}, nil
}

// NewAdmissionTransactor creates a new write-only instance of Admission, bound to a specific deployed contract.
func NewAdmissionTransactor(address common.Address, transactor bind.ContractTransactor) (*AdmissionTransactor, error) {
	contract, err := bindAdmission(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AdmissionTransactor{contract: contract}, nil
}

// NewAdmissionFilterer creates a new log filterer instance of Admission, bound to a specific deployed contract.
func NewAdmissionFilterer(address common.Address, filterer bind.ContractFilterer) (*AdmissionFilterer, error) {
	contract, err := bindAdmission(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AdmissionFilterer{contract: contract}, nil
}

// bindAdmission binds a generic wrapper to an already deployed contract.
func bindAdmission(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AdmissionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Admission *AdmissionRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Admission.Contract.AdmissionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Admission *AdmissionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Admission.Contract.AdmissionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Admission *AdmissionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Admission.Contract.AdmissionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Admission *AdmissionCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Admission.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Admission *AdmissionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Admission.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Admission *AdmissionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Admission.Contract.contract.Transact(opts, method, params...)
}

// CpuDifficulty is a free data retrieval call binding the contract method 0x8b654613.
//
// Solidity: function cpuDifficulty() constant returns(uint256)
func (_Admission *AdmissionCaller) CpuDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Admission.contract.Call(opts, out, "cpuDifficulty")
	return *ret0, err
}

// CpuDifficulty is a free data retrieval call binding the contract method 0x8b654613.
//
// Solidity: function cpuDifficulty() constant returns(uint256)
func (_Admission *AdmissionSession) CpuDifficulty() (*big.Int, error) {
	return _Admission.Contract.CpuDifficulty(&_Admission.CallOpts)
}

// CpuDifficulty is a free data retrieval call binding the contract method 0x8b654613.
//
// Solidity: function cpuDifficulty() constant returns(uint256)
func (_Admission *AdmissionCallerSession) CpuDifficulty() (*big.Int, error) {
	return _Admission.Contract.CpuDifficulty(&_Admission.CallOpts)
}

// CpuWorkTimeout is a free data retrieval call binding the contract method 0x615cc243.
//
// Solidity: function cpuWorkTimeout() constant returns(uint256)
func (_Admission *AdmissionCaller) CpuWorkTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Admission.contract.Call(opts, out, "cpuWorkTimeout")
	return *ret0, err
}

// CpuWorkTimeout is a free data retrieval call binding the contract method 0x615cc243.
//
// Solidity: function cpuWorkTimeout() constant returns(uint256)
func (_Admission *AdmissionSession) CpuWorkTimeout() (*big.Int, error) {
	return _Admission.Contract.CpuWorkTimeout(&_Admission.CallOpts)
}

// CpuWorkTimeout is a free data retrieval call binding the contract method 0x615cc243.
//
// Solidity: function cpuWorkTimeout() constant returns(uint256)
func (_Admission *AdmissionCallerSession) CpuWorkTimeout() (*big.Int, error) {
	return _Admission.Contract.CpuWorkTimeout(&_Admission.CallOpts)
}

// GetAdmissionParameters is a free data retrieval call binding the contract method 0xc651cfc9.
//
// Solidity: function getAdmissionParameters() constant returns(uint256, uint256, uint256, uint256)
func (_Admission *AdmissionCaller) GetAdmissionParameters(opts *bind.CallOpts) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
		ret3 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _Admission.contract.Call(opts, out, "getAdmissionParameters")
	return *ret0, *ret1, *ret2, *ret3, err
}

// GetAdmissionParameters is a free data retrieval call binding the contract method 0xc651cfc9.
//
// Solidity: function getAdmissionParameters() constant returns(uint256, uint256, uint256, uint256)
func (_Admission *AdmissionSession) GetAdmissionParameters() (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _Admission.Contract.GetAdmissionParameters(&_Admission.CallOpts)
}

// GetAdmissionParameters is a free data retrieval call binding the contract method 0xc651cfc9.
//
// Solidity: function getAdmissionParameters() constant returns(uint256, uint256, uint256, uint256)
func (_Admission *AdmissionCallerSession) GetAdmissionParameters() (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _Admission.Contract.GetAdmissionParameters(&_Admission.CallOpts)
}

// MemoryDifficulty is a free data retrieval call binding the contract method 0x17e6b966.
//
// Solidity: function memoryDifficulty() constant returns(uint256)
func (_Admission *AdmissionCaller) MemoryDifficulty(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Admission.contract.Call(opts, out, "memoryDifficulty")
	return *ret0, err
}

// MemoryDifficulty is a free data retrieval call binding the contract method 0x17e6b966.
//
// Solidity: function memoryDifficulty() constant returns(uint256)
func (_Admission *AdmissionSession) MemoryDifficulty() (*big.Int, error) {
	return _Admission.Contract.MemoryDifficulty(&_Admission.CallOpts)
}

// MemoryDifficulty is a free data retrieval call binding the contract method 0x17e6b966.
//
// Solidity: function memoryDifficulty() constant returns(uint256)
func (_Admission *AdmissionCallerSession) MemoryDifficulty() (*big.Int, error) {
	return _Admission.Contract.MemoryDifficulty(&_Admission.CallOpts)
}

// MemoryWorkTimeout is a free data retrieval call binding the contract method 0x6d44a935.
//
// Solidity: function memoryWorkTimeout() constant returns(uint256)
func (_Admission *AdmissionCaller) MemoryWorkTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Admission.contract.Call(opts, out, "memoryWorkTimeout")
	return *ret0, err
}

// MemoryWorkTimeout is a free data retrieval call binding the contract method 0x6d44a935.
//
// Solidity: function memoryWorkTimeout() constant returns(uint256)
func (_Admission *AdmissionSession) MemoryWorkTimeout() (*big.Int, error) {
	return _Admission.Contract.MemoryWorkTimeout(&_Admission.CallOpts)
}

// MemoryWorkTimeout is a free data retrieval call binding the contract method 0x6d44a935.
//
// Solidity: function memoryWorkTimeout() constant returns(uint256)
func (_Admission *AdmissionCallerSession) MemoryWorkTimeout() (*big.Int, error) {
	return _Admission.Contract.MemoryWorkTimeout(&_Admission.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Admission *AdmissionCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Admission.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Admission *AdmissionSession) Owner() (common.Address, error) {
	return _Admission.Contract.Owner(&_Admission.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Admission *AdmissionCallerSession) Owner() (common.Address, error) {
	return _Admission.Contract.Owner(&_Admission.CallOpts)
}

// Verify is a free data retrieval call binding the contract method 0x3395492e.
//
// Solidity: function verify(_cpuNonce uint64, _cpuBlockNumber uint256, _memoryNonce uint64, _memoryBlockNumber uint256, _sender address) constant returns(bool)
func (_Admission *AdmissionCaller) Verify(opts *bind.CallOpts, _cpuNonce uint64, _cpuBlockNumber *big.Int, _memoryNonce uint64, _memoryBlockNumber *big.Int, _sender common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Admission.contract.Call(opts, out, "verify", _cpuNonce, _cpuBlockNumber, _memoryNonce, _memoryBlockNumber, _sender)
	return *ret0, err
}

// Verify is a free data retrieval call binding the contract method 0x3395492e.
//
// Solidity: function verify(_cpuNonce uint64, _cpuBlockNumber uint256, _memoryNonce uint64, _memoryBlockNumber uint256, _sender address) constant returns(bool)
func (_Admission *AdmissionSession) Verify(_cpuNonce uint64, _cpuBlockNumber *big.Int, _memoryNonce uint64, _memoryBlockNumber *big.Int, _sender common.Address) (bool, error) {
	return _Admission.Contract.Verify(&_Admission.CallOpts, _cpuNonce, _cpuBlockNumber, _memoryNonce, _memoryBlockNumber, _sender)
}

// Verify is a free data retrieval call binding the contract method 0x3395492e.
//
// Solidity: function verify(_cpuNonce uint64, _cpuBlockNumber uint256, _memoryNonce uint64, _memoryBlockNumber uint256, _sender address) constant returns(bool)
func (_Admission *AdmissionCallerSession) Verify(_cpuNonce uint64, _cpuBlockNumber *big.Int, _memoryNonce uint64, _memoryBlockNumber *big.Int, _sender common.Address) (bool, error) {
	return _Admission.Contract.Verify(&_Admission.CallOpts, _cpuNonce, _cpuBlockNumber, _memoryNonce, _memoryBlockNumber, _sender)
}

// VerifyCPU is a free data retrieval call binding the contract method 0x6ac03dcc.
//
// Solidity: function verifyCPU(_sender address, _nonce uint64, _blockNumber uint256, _difficulty uint256) constant returns(b bool)
func (_Admission *AdmissionCaller) VerifyCPU(opts *bind.CallOpts, _sender common.Address, _nonce uint64, _blockNumber *big.Int, _difficulty *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Admission.contract.Call(opts, out, "verifyCPU", _sender, _nonce, _blockNumber, _difficulty)
	return *ret0, err
}

// VerifyCPU is a free data retrieval call binding the contract method 0x6ac03dcc.
//
// Solidity: function verifyCPU(_sender address, _nonce uint64, _blockNumber uint256, _difficulty uint256) constant returns(b bool)
func (_Admission *AdmissionSession) VerifyCPU(_sender common.Address, _nonce uint64, _blockNumber *big.Int, _difficulty *big.Int) (bool, error) {
	return _Admission.Contract.VerifyCPU(&_Admission.CallOpts, _sender, _nonce, _blockNumber, _difficulty)
}

// VerifyCPU is a free data retrieval call binding the contract method 0x6ac03dcc.
//
// Solidity: function verifyCPU(_sender address, _nonce uint64, _blockNumber uint256, _difficulty uint256) constant returns(b bool)
func (_Admission *AdmissionCallerSession) VerifyCPU(_sender common.Address, _nonce uint64, _blockNumber *big.Int, _difficulty *big.Int) (bool, error) {
	return _Admission.Contract.VerifyCPU(&_Admission.CallOpts, _sender, _nonce, _blockNumber, _difficulty)
}

// VerifyMemory is a free data retrieval call binding the contract method 0x4bda8957.
//
// Solidity: function verifyMemory(_sender address, _nonce uint64, _blockNumber uint256, _difficulty uint256) constant returns(b bool)
func (_Admission *AdmissionCaller) VerifyMemory(opts *bind.CallOpts, _sender common.Address, _nonce uint64, _blockNumber *big.Int, _difficulty *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Admission.contract.Call(opts, out, "verifyMemory", _sender, _nonce, _blockNumber, _difficulty)
	return *ret0, err
}

// VerifyMemory is a free data retrieval call binding the contract method 0x4bda8957.
//
// Solidity: function verifyMemory(_sender address, _nonce uint64, _blockNumber uint256, _difficulty uint256) constant returns(b bool)
func (_Admission *AdmissionSession) VerifyMemory(_sender common.Address, _nonce uint64, _blockNumber *big.Int, _difficulty *big.Int) (bool, error) {
	return _Admission.Contract.VerifyMemory(&_Admission.CallOpts, _sender, _nonce, _blockNumber, _difficulty)
}

// VerifyMemory is a free data retrieval call binding the contract method 0x4bda8957.
//
// Solidity: function verifyMemory(_sender address, _nonce uint64, _blockNumber uint256, _difficulty uint256) constant returns(b bool)
func (_Admission *AdmissionCallerSession) VerifyMemory(_sender common.Address, _nonce uint64, _blockNumber *big.Int, _difficulty *big.Int) (bool, error) {
	return _Admission.Contract.VerifyMemory(&_Admission.CallOpts, _sender, _nonce, _blockNumber, _difficulty)
}

// UpdateCPUDifficulty is a paid mutator transaction binding the contract method 0xbe981db8.
//
// Solidity: function updateCPUDifficulty(_difficulty uint256) returns()
func (_Admission *AdmissionTransactor) UpdateCPUDifficulty(opts *bind.TransactOpts, _difficulty *big.Int) (*types.Transaction, error) {
	return _Admission.contract.Transact(opts, "updateCPUDifficulty", _difficulty)
}

// UpdateCPUDifficulty is a paid mutator transaction binding the contract method 0xbe981db8.
//
// Solidity: function updateCPUDifficulty(_difficulty uint256) returns()
func (_Admission *AdmissionSession) UpdateCPUDifficulty(_difficulty *big.Int) (*types.Transaction, error) {
	return _Admission.Contract.UpdateCPUDifficulty(&_Admission.TransactOpts, _difficulty)
}

// UpdateCPUDifficulty is a paid mutator transaction binding the contract method 0xbe981db8.
//
// Solidity: function updateCPUDifficulty(_difficulty uint256) returns()
func (_Admission *AdmissionTransactorSession) UpdateCPUDifficulty(_difficulty *big.Int) (*types.Transaction, error) {
	return _Admission.Contract.UpdateCPUDifficulty(&_Admission.TransactOpts, _difficulty)
}

// UpdateCPUWorkTimeout is a paid mutator transaction binding the contract method 0x31ee5d33.
//
// Solidity: function updateCPUWorkTimeout(_cpuWorkTimeout uint256) returns()
func (_Admission *AdmissionTransactor) UpdateCPUWorkTimeout(opts *bind.TransactOpts, _cpuWorkTimeout *big.Int) (*types.Transaction, error) {
	return _Admission.contract.Transact(opts, "updateCPUWorkTimeout", _cpuWorkTimeout)
}

// UpdateCPUWorkTimeout is a paid mutator transaction binding the contract method 0x31ee5d33.
//
// Solidity: function updateCPUWorkTimeout(_cpuWorkTimeout uint256) returns()
func (_Admission *AdmissionSession) UpdateCPUWorkTimeout(_cpuWorkTimeout *big.Int) (*types.Transaction, error) {
	return _Admission.Contract.UpdateCPUWorkTimeout(&_Admission.TransactOpts, _cpuWorkTimeout)
}

// UpdateCPUWorkTimeout is a paid mutator transaction binding the contract method 0x31ee5d33.
//
// Solidity: function updateCPUWorkTimeout(_cpuWorkTimeout uint256) returns()
func (_Admission *AdmissionTransactorSession) UpdateCPUWorkTimeout(_cpuWorkTimeout *big.Int) (*types.Transaction, error) {
	return _Admission.Contract.UpdateCPUWorkTimeout(&_Admission.TransactOpts, _cpuWorkTimeout)
}

// UpdateMemoryDifficulty is a paid mutator transaction binding the contract method 0xe1718946.
//
// Solidity: function updateMemoryDifficulty(_difficulty uint256) returns()
func (_Admission *AdmissionTransactor) UpdateMemoryDifficulty(opts *bind.TransactOpts, _difficulty *big.Int) (*types.Transaction, error) {
	return _Admission.contract.Transact(opts, "updateMemoryDifficulty", _difficulty)
}

// UpdateMemoryDifficulty is a paid mutator transaction binding the contract method 0xe1718946.
//
// Solidity: function updateMemoryDifficulty(_difficulty uint256) returns()
func (_Admission *AdmissionSession) UpdateMemoryDifficulty(_difficulty *big.Int) (*types.Transaction, error) {
	return _Admission.Contract.UpdateMemoryDifficulty(&_Admission.TransactOpts, _difficulty)
}

// UpdateMemoryDifficulty is a paid mutator transaction binding the contract method 0xe1718946.
//
// Solidity: function updateMemoryDifficulty(_difficulty uint256) returns()
func (_Admission *AdmissionTransactorSession) UpdateMemoryDifficulty(_difficulty *big.Int) (*types.Transaction, error) {
	return _Admission.Contract.UpdateMemoryDifficulty(&_Admission.TransactOpts, _difficulty)
}

// UpdateMemoryWorkTimeout is a paid mutator transaction binding the contract method 0x86395077.
//
// Solidity: function updateMemoryWorkTimeout(_memoryWorkTimeout uint256) returns()
func (_Admission *AdmissionTransactor) UpdateMemoryWorkTimeout(opts *bind.TransactOpts, _memoryWorkTimeout *big.Int) (*types.Transaction, error) {
	return _Admission.contract.Transact(opts, "updateMemoryWorkTimeout", _memoryWorkTimeout)
}

// UpdateMemoryWorkTimeout is a paid mutator transaction binding the contract method 0x86395077.
//
// Solidity: function updateMemoryWorkTimeout(_memoryWorkTimeout uint256) returns()
func (_Admission *AdmissionSession) UpdateMemoryWorkTimeout(_memoryWorkTimeout *big.Int) (*types.Transaction, error) {
	return _Admission.Contract.UpdateMemoryWorkTimeout(&_Admission.TransactOpts, _memoryWorkTimeout)
}

// UpdateMemoryWorkTimeout is a paid mutator transaction binding the contract method 0x86395077.
//
// Solidity: function updateMemoryWorkTimeout(_memoryWorkTimeout uint256) returns()
func (_Admission *AdmissionTransactorSession) UpdateMemoryWorkTimeout(_memoryWorkTimeout *big.Int) (*types.Transaction, error) {
	return _Admission.Contract.UpdateMemoryWorkTimeout(&_Admission.TransactOpts, _memoryWorkTimeout)
}
