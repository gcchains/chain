// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rnode

import (
	"math/big"
	"strings"

	gcchain "/gcchain/chain"
	"github.com/gcchains/chain/accounts/abi"
	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

// RnodeABI is the input ABI used to generate the binding from.
const RnodeABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getRnodeNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_period\",\"type\":\"uint256\"}],\"name\":\"setPeriod\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"quitRnode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isContract\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"enabled\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"enableContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"refundAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"Participants\",\"outputs\":[{\"name\":\"lockedDeposit\",\"type\":\"uint256\"},{\"name\":\"lockedTime\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"setSupportedVersion\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"disableContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"threshold\",\"type\":\"uint256\"}],\"name\":\"setRnodeThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"isRnode\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"version\",\"type\":\"uint256\"}],\"name\":\"joinRnode\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rnodeThreshold\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"supportedVersion\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRnodes\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"period\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"investor\",\"type\":\"address\"}],\"name\":\"refund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"who\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"lockedDeposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"lockedTime\",\"type\":\"uint256\"}],\"name\":\"NewRnode\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"who\",\"type\":\"address\"}],\"name\":\"RnodeQuit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"who\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ownerRefund\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"numOfInvestor\",\"type\":\"uint256\"}],\"name\":\"ownerRefundAll\",\"type\":\"event\"}]"

// RnodeBin is the compiled bytecode used for deploying new contracts.
const RnodeBin = `0x60806040526117706001908155692a5a058fc295ed00000060025560038190556007805460ff1916909117905534801561003857600080fd5b5060008054600160a060020a03191633179055610be88061005a6000396000f3006080604052600436106100fb5763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416630b443f4281146101005780630f3a9f6514610127578063113c8498146101415780631627905514610156578063238dafe01461018b578063367edd32146101a057806338e771ab146101b5578063595aa13d146101ca5780635f86d4ca14610204578063894ba8331461021c578063975dd4b214610231578063a8f0769714610249578063aae80f781461026a578063b7b3e9da14610275578063d5601e9f1461028a578063e508bb851461029f578063ef78d4fd14610304578063fa89401a14610319575b600080fd5b34801561010c57600080fd5b5061011561033a565b60408051918252519081900360200190f35b34801561013357600080fd5b5061013f600435610341565b005b34801561014d57600080fd5b5061013f61036d565b34801561016257600080fd5b50610177600160a060020a036004351661043e565b604080519115158252519081900360200190f35b34801561019757600080fd5b50610177610446565b3480156101ac57600080fd5b5061013f61044f565b3480156101c157600080fd5b5061013f610475565b3480156101d657600080fd5b506101eb600160a060020a03600435166105c6565b6040805192835260208301919091528051918290030190f35b34801561021057600080fd5b5061013f6004356105df565b34801561022857600080fd5b5061013f6105fb565b34801561023d57600080fd5b5061013f60043561061e565b34801561025557600080fd5b50610177600160a060020a0360043516610651565b61013f60043561066a565b34801561028157600080fd5b506101156107f4565b34801561029657600080fd5b506101156107fa565b3480156102ab57600080fd5b506102b4610800565b60408051602080825283518183015283519192839290830191858101910280838360005b838110156102f05781810151838201526020016102d8565b505050509050019250505060405180910390f35b34801561031057600080fd5b50610115610811565b34801561032557600080fd5b5061013f600160a060020a0360043516610817565b6006545b90565b600054600160a060020a0316331461035857600080fd5b6201518081111561036857600080fd5b600155565b61037e60043363ffffffff61090516565b151561038957600080fd5b600180543360009081526008602052604090209091015442910111156103ae57600080fd5b3360008181526008602052604080822054905181156108fc0292818181858888f193505050501580156103e5573d6000803e3d6000fd5b50336000818152600860205260408120556104089060049063ffffffff61092416565b506040805133815290517f602a2a9c94f70293aa2be9077f0b2dc89d388bc293fdbcd968274f43494c380d9181900360200190a1565b6000903b1190565b60075460ff1681565b600054600160a060020a0316331461046657600080fd5b6007805460ff19166001179055565b60008054819081908190600160a060020a0316331461049357600080fd5b6006549350600092505b83831015610583576006805460009081106104b457fe5b6000918252602080832090910154600160a060020a0316808352600890915260408083205490519194509250839183156108fc02918491818181858888f19350505050158015610508573d6000803e3d6000fd5b50600160a060020a03821660009081526008602052604081205561053360048363ffffffff61092416565b5060408051600160a060020a03841681526020810183905281517f3914ba80eb00486e7a58b91fb4795283df0c5b507eea9cf7c77cce26cc70d25c929181900390910190a160019092019161049d565b6006541561058d57fe5b6040805185815290517fb65ebb6b17695b3a5612c7a0f6f60e649c02ba24b36b546b8d037e98215fdb8d9181900360200190a150505050565b6008602052600090815260409020805460019091015482565b600054600160a060020a031633146105f657600080fd5b600355565b600054600160a060020a0316331461061257600080fd5b6007805460ff19169055565b600054600160a060020a0316331461063557600080fd5b692a5a058fc295ed00000081101561064c57600080fd5b600255565b600061066460048363ffffffff61090516565b92915050565b60075460ff16151561067b57600080fd5b60035481101561068a57600080fd5b6106933361043e565b1561072557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f706c65617365206e6f742075736520636f6e74726163742063616c6c2074686960448201527f732066756e6374696f6e00000000000000000000000000000000000000000000606482015290519081900360840190fd5b61073660043363ffffffff61090516565b1561074057600080fd5b60025434101561074f57600080fd5b3360009081526008602052604090205461076f903463ffffffff610a6a16565b3360008181526008602052604090209182554260019092019190915561079d9060049063ffffffff610a8016565b5033600081815260086020908152604091829020805460019091015483519485529184015282820152517f586bfaa7a657ad9313326c9269639546950d589bd479b3d6928be469d6dc29039181900360600190a150565b60025481565b60035481565b606061080c6004610b0f565b905090565b60015481565b60008054600160a060020a0316331461082f57600080fd5b61084060048363ffffffff61090516565b151561084b57600080fd5b50600160a060020a03811660008181526008602052604080822054905190929183156108fc02918491818181858888f19350505050158015610891573d6000803e3d6000fd5b50600160a060020a0382166000908152600860205260408120556108bc60048363ffffffff61092416565b5060408051600160a060020a03841681526020810183905281517f3914ba80eb00486e7a58b91fb4795283df0c5b507eea9cf7c77cce26cc70d25c929181900390910190a15050565b600160a060020a03166000908152602091909152604090205460ff1690565b600160a060020a03811660009081526020839052604081205481908190819060ff1615156109555760009350610a61565b600160a060020a038516600090815260208781526040808320805460ff1916905560028901805460018b019093529220549094509250600019840184811061099957fe5b600091825260209091200154600287018054600160a060020a0390921692508291849081106109c457fe5b6000918252602080832091909101805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a03948516179055918316815260018801909152604090208290556002860180546000198501908110610a2057fe5b6000918252602090912001805473ffffffffffffffffffffffffffffffffffffffff1916905560028601805490610a5b906000198301610b75565b50600193505b50505092915050565b600082820183811015610a7957fe5b9392505050565b600160a060020a03811660009081526020839052604081205460ff1615610aa957506000610664565b50600160a060020a0316600081815260208381526040808320805460ff19166001908117909155600286018054968201845291842086905585810182559083529120909201805473ffffffffffffffffffffffffffffffffffffffff1916909117905590565b606081600201805480602002602001604051908101604052809291908181526020018280548015610b6957602002820191906000526020600020905b8154600160a060020a03168152600190910190602001808311610b4b575b50505050509050919050565b815481835581811115610b9957600083815260209020610b99918101908301610b9e565b505050565b61033e91905b80821115610bb85760008155600101610ba4565b50905600a165627a7a723058206dd2e368d6f0c7701b45b4d92495e1edfa972b0b9e6ad7e7a11b0f4d9c1f03a00029`

// DeployRnode deploys a new gcchain contract, binding an instance of Rnode to it.
func DeployRnode(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Rnode, error) {
	parsed, err := abi.JSON(strings.NewReader(RnodeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RnodeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Rnode{RnodeCaller: RnodeCaller{contract: contract}, RnodeTransactor: RnodeTransactor{contract: contract}, RnodeFilterer: RnodeFilterer{contract: contract}}, nil
}

// Rnode is an auto generated Go binding around an gcchain contract.
type Rnode struct {
	RnodeCaller     // Read-only binding to the contract
	RnodeTransactor // Write-only binding to the contract
	RnodeFilterer   // Log filterer for contract events
}

// RnodeCaller is an auto generated read-only Go binding around an gcchain contract.
type RnodeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RnodeTransactor is an auto generated write-only Go binding around an gcchain contract.
type RnodeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RnodeFilterer is an auto generated log filtering Go binding around an gcchain contract events.
type RnodeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RnodeSession is an auto generated Go binding around an gcchain contract,
// with pre-set call and transact options.
type RnodeSession struct {
	Contract     *Rnode            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RnodeCallerSession is an auto generated read-only Go binding around an gcchain contract,
// with pre-set call options.
type RnodeCallerSession struct {
	Contract *RnodeCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RnodeTransactorSession is an auto generated write-only Go binding around an gcchain contract,
// with pre-set transact options.
type RnodeTransactorSession struct {
	Contract     *RnodeTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RnodeRaw is an auto generated low-level Go binding around an gcchain contract.
type RnodeRaw struct {
	Contract *Rnode // Generic contract binding to access the raw methods on
}

// RnodeCallerRaw is an auto generated low-level read-only Go binding around an gcchain contract.
type RnodeCallerRaw struct {
	Contract *RnodeCaller // Generic read-only contract binding to access the raw methods on
}

// RnodeTransactorRaw is an auto generated low-level write-only Go binding around an gcchain contract.
type RnodeTransactorRaw struct {
	Contract *RnodeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRnode creates a new instance of Rnode, bound to a specific deployed contract.
func NewRnode(address common.Address, backend bind.ContractBackend) (*Rnode, error) {
	contract, err := bindRnode(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rnode{RnodeCaller: RnodeCaller{contract: contract}, RnodeTransactor: RnodeTransactor{contract: contract}, RnodeFilterer: RnodeFilterer{contract: contract}}, nil
}

// NewRnodeCaller creates a new read-only instance of Rnode, bound to a specific deployed contract.
func NewRnodeCaller(address common.Address, caller bind.ContractCaller) (*RnodeCaller, error) {
	contract, err := bindRnode(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RnodeCaller{contract: contract}, nil
}

// NewRnodeTransactor creates a new write-only instance of Rnode, bound to a specific deployed contract.
func NewRnodeTransactor(address common.Address, transactor bind.ContractTransactor) (*RnodeTransactor, error) {
	contract, err := bindRnode(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RnodeTransactor{contract: contract}, nil
}

// NewRnodeFilterer creates a new log filterer instance of Rnode, bound to a specific deployed contract.
func NewRnodeFilterer(address common.Address, filterer bind.ContractFilterer) (*RnodeFilterer, error) {
	contract, err := bindRnode(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RnodeFilterer{contract: contract}, nil
}

// bindRnode binds a generic wrapper to an already deployed contract.
func bindRnode(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RnodeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rnode *RnodeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Rnode.Contract.RnodeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rnode *RnodeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rnode.Contract.RnodeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rnode *RnodeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rnode.Contract.RnodeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rnode *RnodeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Rnode.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rnode *RnodeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rnode.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rnode *RnodeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rnode.Contract.contract.Transact(opts, method, params...)
}

// Participants is a free data retrieval call binding the contract method 0x595aa13d.
//
// Solidity: function Participants( address) constant returns(lockedDeposit uint256, lockedTime uint256)
func (_Rnode *RnodeCaller) Participants(opts *bind.CallOpts, arg0 common.Address) (struct {
	LockedDeposit *big.Int
	LockedTime    *big.Int
}, error) {
	ret := new(struct {
		LockedDeposit *big.Int
		LockedTime    *big.Int
	})
	out := ret
	err := _Rnode.contract.Call(opts, out, "Participants", arg0)
	return *ret, err
}

// Participants is a free data retrieval call binding the contract method 0x595aa13d.
//
// Solidity: function Participants( address) constant returns(lockedDeposit uint256, lockedTime uint256)
func (_Rnode *RnodeSession) Participants(arg0 common.Address) (struct {
	LockedDeposit *big.Int
	LockedTime    *big.Int
}, error) {
	return _Rnode.Contract.Participants(&_Rnode.CallOpts, arg0)
}

// Participants is a free data retrieval call binding the contract method 0x595aa13d.
//
// Solidity: function Participants( address) constant returns(lockedDeposit uint256, lockedTime uint256)
func (_Rnode *RnodeCallerSession) Participants(arg0 common.Address) (struct {
	LockedDeposit *big.Int
	LockedTime    *big.Int
}, error) {
	return _Rnode.Contract.Participants(&_Rnode.CallOpts, arg0)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(bool)
func (_Rnode *RnodeCaller) Enabled(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rnode.contract.Call(opts, out, "enabled")
	return *ret0, err
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(bool)
func (_Rnode *RnodeSession) Enabled() (bool, error) {
	return _Rnode.Contract.Enabled(&_Rnode.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(bool)
func (_Rnode *RnodeCallerSession) Enabled() (bool, error) {
	return _Rnode.Contract.Enabled(&_Rnode.CallOpts)
}

// GetRnodeNum is a free data retrieval call binding the contract method 0x0b443f42.
//
// Solidity: function getRnodeNum() constant returns(uint256)
func (_Rnode *RnodeCaller) GetRnodeNum(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rnode.contract.Call(opts, out, "getRnodeNum")
	return *ret0, err
}

// GetRnodeNum is a free data retrieval call binding the contract method 0x0b443f42.
//
// Solidity: function getRnodeNum() constant returns(uint256)
func (_Rnode *RnodeSession) GetRnodeNum() (*big.Int, error) {
	return _Rnode.Contract.GetRnodeNum(&_Rnode.CallOpts)
}

// GetRnodeNum is a free data retrieval call binding the contract method 0x0b443f42.
//
// Solidity: function getRnodeNum() constant returns(uint256)
func (_Rnode *RnodeCallerSession) GetRnodeNum() (*big.Int, error) {
	return _Rnode.Contract.GetRnodeNum(&_Rnode.CallOpts)
}

// GetRnodes is a free data retrieval call binding the contract method 0xe508bb85.
//
// Solidity: function getRnodes() constant returns(address[])
func (_Rnode *RnodeCaller) GetRnodes(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Rnode.contract.Call(opts, out, "getRnodes")
	return *ret0, err
}

// GetRnodes is a free data retrieval call binding the contract method 0xe508bb85.
//
// Solidity: function getRnodes() constant returns(address[])
func (_Rnode *RnodeSession) GetRnodes() ([]common.Address, error) {
	return _Rnode.Contract.GetRnodes(&_Rnode.CallOpts)
}

// GetRnodes is a free data retrieval call binding the contract method 0xe508bb85.
//
// Solidity: function getRnodes() constant returns(address[])
func (_Rnode *RnodeCallerSession) GetRnodes() ([]common.Address, error) {
	return _Rnode.Contract.GetRnodes(&_Rnode.CallOpts)
}

// IsContract is a free data retrieval call binding the contract method 0x16279055.
//
// Solidity: function isContract(addr address) constant returns(bool)
func (_Rnode *RnodeCaller) IsContract(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rnode.contract.Call(opts, out, "isContract", addr)
	return *ret0, err
}

// IsContract is a free data retrieval call binding the contract method 0x16279055.
//
// Solidity: function isContract(addr address) constant returns(bool)
func (_Rnode *RnodeSession) IsContract(addr common.Address) (bool, error) {
	return _Rnode.Contract.IsContract(&_Rnode.CallOpts, addr)
}

// IsContract is a free data retrieval call binding the contract method 0x16279055.
//
// Solidity: function isContract(addr address) constant returns(bool)
func (_Rnode *RnodeCallerSession) IsContract(addr common.Address) (bool, error) {
	return _Rnode.Contract.IsContract(&_Rnode.CallOpts, addr)
}

// IsRnode is a free data retrieval call binding the contract method 0xa8f07697.
//
// Solidity: function isRnode(addr address) constant returns(bool)
func (_Rnode *RnodeCaller) IsRnode(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rnode.contract.Call(opts, out, "isRnode", addr)
	return *ret0, err
}

// IsRnode is a free data retrieval call binding the contract method 0xa8f07697.
//
// Solidity: function isRnode(addr address) constant returns(bool)
func (_Rnode *RnodeSession) IsRnode(addr common.Address) (bool, error) {
	return _Rnode.Contract.IsRnode(&_Rnode.CallOpts, addr)
}

// IsRnode is a free data retrieval call binding the contract method 0xa8f07697.
//
// Solidity: function isRnode(addr address) constant returns(bool)
func (_Rnode *RnodeCallerSession) IsRnode(addr common.Address) (bool, error) {
	return _Rnode.Contract.IsRnode(&_Rnode.CallOpts, addr)
}

// Period is a free data retrieval call binding the contract method 0xef78d4fd.
//
// Solidity: function period() constant returns(uint256)
func (_Rnode *RnodeCaller) Period(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rnode.contract.Call(opts, out, "period")
	return *ret0, err
}

// Period is a free data retrieval call binding the contract method 0xef78d4fd.
//
// Solidity: function period() constant returns(uint256)
func (_Rnode *RnodeSession) Period() (*big.Int, error) {
	return _Rnode.Contract.Period(&_Rnode.CallOpts)
}

// Period is a free data retrieval call binding the contract method 0xef78d4fd.
//
// Solidity: function period() constant returns(uint256)
func (_Rnode *RnodeCallerSession) Period() (*big.Int, error) {
	return _Rnode.Contract.Period(&_Rnode.CallOpts)
}

// RnodeThreshold is a free data retrieval call binding the contract method 0xb7b3e9da.
//
// Solidity: function rnodeThreshold() constant returns(uint256)
func (_Rnode *RnodeCaller) RnodeThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rnode.contract.Call(opts, out, "rnodeThreshold")
	return *ret0, err
}

// RnodeThreshold is a free data retrieval call binding the contract method 0xb7b3e9da.
//
// Solidity: function rnodeThreshold() constant returns(uint256)
func (_Rnode *RnodeSession) RnodeThreshold() (*big.Int, error) {
	return _Rnode.Contract.RnodeThreshold(&_Rnode.CallOpts)
}

// RnodeThreshold is a free data retrieval call binding the contract method 0xb7b3e9da.
//
// Solidity: function rnodeThreshold() constant returns(uint256)
func (_Rnode *RnodeCallerSession) RnodeThreshold() (*big.Int, error) {
	return _Rnode.Contract.RnodeThreshold(&_Rnode.CallOpts)
}

// SupportedVersion is a free data retrieval call binding the contract method 0xd5601e9f.
//
// Solidity: function supportedVersion() constant returns(uint256)
func (_Rnode *RnodeCaller) SupportedVersion(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rnode.contract.Call(opts, out, "supportedVersion")
	return *ret0, err
}

// SupportedVersion is a free data retrieval call binding the contract method 0xd5601e9f.
//
// Solidity: function supportedVersion() constant returns(uint256)
func (_Rnode *RnodeSession) SupportedVersion() (*big.Int, error) {
	return _Rnode.Contract.SupportedVersion(&_Rnode.CallOpts)
}

// SupportedVersion is a free data retrieval call binding the contract method 0xd5601e9f.
//
// Solidity: function supportedVersion() constant returns(uint256)
func (_Rnode *RnodeCallerSession) SupportedVersion() (*big.Int, error) {
	return _Rnode.Contract.SupportedVersion(&_Rnode.CallOpts)
}

// DisableContract is a paid mutator transaction binding the contract method 0x894ba833.
//
// Solidity: function disableContract() returns()
func (_Rnode *RnodeTransactor) DisableContract(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rnode.contract.Transact(opts, "disableContract")
}

// DisableContract is a paid mutator transaction binding the contract method 0x894ba833.
//
// Solidity: function disableContract() returns()
func (_Rnode *RnodeSession) DisableContract() (*types.Transaction, error) {
	return _Rnode.Contract.DisableContract(&_Rnode.TransactOpts)
}

// DisableContract is a paid mutator transaction binding the contract method 0x894ba833.
//
// Solidity: function disableContract() returns()
func (_Rnode *RnodeTransactorSession) DisableContract() (*types.Transaction, error) {
	return _Rnode.Contract.DisableContract(&_Rnode.TransactOpts)
}

// EnableContract is a paid mutator transaction binding the contract method 0x367edd32.
//
// Solidity: function enableContract() returns()
func (_Rnode *RnodeTransactor) EnableContract(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rnode.contract.Transact(opts, "enableContract")
}

// EnableContract is a paid mutator transaction binding the contract method 0x367edd32.
//
// Solidity: function enableContract() returns()
func (_Rnode *RnodeSession) EnableContract() (*types.Transaction, error) {
	return _Rnode.Contract.EnableContract(&_Rnode.TransactOpts)
}

// EnableContract is a paid mutator transaction binding the contract method 0x367edd32.
//
// Solidity: function enableContract() returns()
func (_Rnode *RnodeTransactorSession) EnableContract() (*types.Transaction, error) {
	return _Rnode.Contract.EnableContract(&_Rnode.TransactOpts)
}

// JoinRnode is a paid mutator transaction binding the contract method 0xaae80f78.
//
// Solidity: function joinRnode(version uint256) returns()
func (_Rnode *RnodeTransactor) JoinRnode(opts *bind.TransactOpts, version *big.Int) (*types.Transaction, error) {
	return _Rnode.contract.Transact(opts, "joinRnode", version)
}

// JoinRnode is a paid mutator transaction binding the contract method 0xaae80f78.
//
// Solidity: function joinRnode(version uint256) returns()
func (_Rnode *RnodeSession) JoinRnode(version *big.Int) (*types.Transaction, error) {
	return _Rnode.Contract.JoinRnode(&_Rnode.TransactOpts, version)
}

// JoinRnode is a paid mutator transaction binding the contract method 0xaae80f78.
//
// Solidity: function joinRnode(version uint256) returns()
func (_Rnode *RnodeTransactorSession) JoinRnode(version *big.Int) (*types.Transaction, error) {
	return _Rnode.Contract.JoinRnode(&_Rnode.TransactOpts, version)
}

// QuitRnode is a paid mutator transaction binding the contract method 0x113c8498.
//
// Solidity: function quitRnode() returns()
func (_Rnode *RnodeTransactor) QuitRnode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rnode.contract.Transact(opts, "quitRnode")
}

// QuitRnode is a paid mutator transaction binding the contract method 0x113c8498.
//
// Solidity: function quitRnode() returns()
func (_Rnode *RnodeSession) QuitRnode() (*types.Transaction, error) {
	return _Rnode.Contract.QuitRnode(&_Rnode.TransactOpts)
}

// QuitRnode is a paid mutator transaction binding the contract method 0x113c8498.
//
// Solidity: function quitRnode() returns()
func (_Rnode *RnodeTransactorSession) QuitRnode() (*types.Transaction, error) {
	return _Rnode.Contract.QuitRnode(&_Rnode.TransactOpts)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(investor address) returns()
func (_Rnode *RnodeTransactor) Refund(opts *bind.TransactOpts, investor common.Address) (*types.Transaction, error) {
	return _Rnode.contract.Transact(opts, "refund", investor)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(investor address) returns()
func (_Rnode *RnodeSession) Refund(investor common.Address) (*types.Transaction, error) {
	return _Rnode.Contract.Refund(&_Rnode.TransactOpts, investor)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(investor address) returns()
func (_Rnode *RnodeTransactorSession) Refund(investor common.Address) (*types.Transaction, error) {
	return _Rnode.Contract.Refund(&_Rnode.TransactOpts, investor)
}

// RefundAll is a paid mutator transaction binding the contract method 0x38e771ab.
//
// Solidity: function refundAll() returns()
func (_Rnode *RnodeTransactor) RefundAll(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rnode.contract.Transact(opts, "refundAll")
}

// RefundAll is a paid mutator transaction binding the contract method 0x38e771ab.
//
// Solidity: function refundAll() returns()
func (_Rnode *RnodeSession) RefundAll() (*types.Transaction, error) {
	return _Rnode.Contract.RefundAll(&_Rnode.TransactOpts)
}

// RefundAll is a paid mutator transaction binding the contract method 0x38e771ab.
//
// Solidity: function refundAll() returns()
func (_Rnode *RnodeTransactorSession) RefundAll() (*types.Transaction, error) {
	return _Rnode.Contract.RefundAll(&_Rnode.TransactOpts)
}

// SetPeriod is a paid mutator transaction binding the contract method 0x0f3a9f65.
//
// Solidity: function setPeriod(_period uint256) returns()
func (_Rnode *RnodeTransactor) SetPeriod(opts *bind.TransactOpts, _period *big.Int) (*types.Transaction, error) {
	return _Rnode.contract.Transact(opts, "setPeriod", _period)
}

// SetPeriod is a paid mutator transaction binding the contract method 0x0f3a9f65.
//
// Solidity: function setPeriod(_period uint256) returns()
func (_Rnode *RnodeSession) SetPeriod(_period *big.Int) (*types.Transaction, error) {
	return _Rnode.Contract.SetPeriod(&_Rnode.TransactOpts, _period)
}

// SetPeriod is a paid mutator transaction binding the contract method 0x0f3a9f65.
//
// Solidity: function setPeriod(_period uint256) returns()
func (_Rnode *RnodeTransactorSession) SetPeriod(_period *big.Int) (*types.Transaction, error) {
	return _Rnode.Contract.SetPeriod(&_Rnode.TransactOpts, _period)
}

// SetRnodeThreshold is a paid mutator transaction binding the contract method 0x975dd4b2.
//
// Solidity: function setRnodeThreshold(threshold uint256) returns()
func (_Rnode *RnodeTransactor) SetRnodeThreshold(opts *bind.TransactOpts, threshold *big.Int) (*types.Transaction, error) {
	return _Rnode.contract.Transact(opts, "setRnodeThreshold", threshold)
}

// SetRnodeThreshold is a paid mutator transaction binding the contract method 0x975dd4b2.
//
// Solidity: function setRnodeThreshold(threshold uint256) returns()
func (_Rnode *RnodeSession) SetRnodeThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Rnode.Contract.SetRnodeThreshold(&_Rnode.TransactOpts, threshold)
}

// SetRnodeThreshold is a paid mutator transaction binding the contract method 0x975dd4b2.
//
// Solidity: function setRnodeThreshold(threshold uint256) returns()
func (_Rnode *RnodeTransactorSession) SetRnodeThreshold(threshold *big.Int) (*types.Transaction, error) {
	return _Rnode.Contract.SetRnodeThreshold(&_Rnode.TransactOpts, threshold)
}

// SetSupportedVersion is a paid mutator transaction binding the contract method 0x5f86d4ca.
//
// Solidity: function setSupportedVersion(version uint256) returns()
func (_Rnode *RnodeTransactor) SetSupportedVersion(opts *bind.TransactOpts, version *big.Int) (*types.Transaction, error) {
	return _Rnode.contract.Transact(opts, "setSupportedVersion", version)
}

// SetSupportedVersion is a paid mutator transaction binding the contract method 0x5f86d4ca.
//
// Solidity: function setSupportedVersion(version uint256) returns()
func (_Rnode *RnodeSession) SetSupportedVersion(version *big.Int) (*types.Transaction, error) {
	return _Rnode.Contract.SetSupportedVersion(&_Rnode.TransactOpts, version)
}

// SetSupportedVersion is a paid mutator transaction binding the contract method 0x5f86d4ca.
//
// Solidity: function setSupportedVersion(version uint256) returns()
func (_Rnode *RnodeTransactorSession) SetSupportedVersion(version *big.Int) (*types.Transaction, error) {
	return _Rnode.Contract.SetSupportedVersion(&_Rnode.TransactOpts, version)
}

// RnodeNewRnodeIterator is returned from FilterNewRnode and is used to iterate over the raw logs and unpacked data for NewRnode events raised by the Rnode contract.
type RnodeNewRnodeIterator struct {
	Event *RnodeNewRnode // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log       // Log channel receiving the found contract events
	sub  gcchain.Subscription // Subscription for errors, completion and termination
	done bool                 // Whether the subscription completed delivering logs
	fail error                // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RnodeNewRnodeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RnodeNewRnode)
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
		it.Event = new(RnodeNewRnode)
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
func (it *RnodeNewRnodeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RnodeNewRnodeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RnodeNewRnode represents a NewRnode event raised by the Rnode contract.
type RnodeNewRnode struct {
	Who           common.Address
	LockedDeposit *big.Int
	LockedTime    *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNewRnode is a free log retrieval operation binding the contract event 0x586bfaa7a657ad9313326c9269639546950d589bd479b3d6928be469d6dc2903.
//
// Solidity: e NewRnode(who address, lockedDeposit uint256, lockedTime uint256)
func (_Rnode *RnodeFilterer) FilterNewRnode(opts *bind.FilterOpts) (*RnodeNewRnodeIterator, error) {

	logs, sub, err := _Rnode.contract.FilterLogs(opts, "NewRnode")
	if err != nil {
		return nil, err
	}
	return &RnodeNewRnodeIterator{contract: _Rnode.contract, event: "NewRnode", logs: logs, sub: sub}, nil
}

// WatchNewRnode is a free log subscription operation binding the contract event 0x586bfaa7a657ad9313326c9269639546950d589bd479b3d6928be469d6dc2903.
//
// Solidity: e NewRnode(who address, lockedDeposit uint256, lockedTime uint256)
func (_Rnode *RnodeFilterer) WatchNewRnode(opts *bind.WatchOpts, sink chan<- *RnodeNewRnode) (event.Subscription, error) {

	logs, sub, err := _Rnode.contract.WatchLogs(opts, "NewRnode")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RnodeNewRnode)
				if err := _Rnode.contract.UnpackLog(event, "NewRnode", log); err != nil {
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

// RnodeRnodeQuitIterator is returned from FilterRnodeQuit and is used to iterate over the raw logs and unpacked data for RnodeQuit events raised by the Rnode contract.
type RnodeRnodeQuitIterator struct {
	Event *RnodeRnodeQuit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log       // Log channel receiving the found contract events
	sub  gcchain.Subscription // Subscription for errors, completion and termination
	done bool                 // Whether the subscription completed delivering logs
	fail error                // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RnodeRnodeQuitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RnodeRnodeQuit)
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
		it.Event = new(RnodeRnodeQuit)
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
func (it *RnodeRnodeQuitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RnodeRnodeQuitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RnodeRnodeQuit represents a RnodeQuit event raised by the Rnode contract.
type RnodeRnodeQuit struct {
	Who common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterRnodeQuit is a free log retrieval operation binding the contract event 0x602a2a9c94f70293aa2be9077f0b2dc89d388bc293fdbcd968274f43494c380d.
//
// Solidity: e RnodeQuit(who address)
func (_Rnode *RnodeFilterer) FilterRnodeQuit(opts *bind.FilterOpts) (*RnodeRnodeQuitIterator, error) {

	logs, sub, err := _Rnode.contract.FilterLogs(opts, "RnodeQuit")
	if err != nil {
		return nil, err
	}
	return &RnodeRnodeQuitIterator{contract: _Rnode.contract, event: "RnodeQuit", logs: logs, sub: sub}, nil
}

// WatchRnodeQuit is a free log subscription operation binding the contract event 0x602a2a9c94f70293aa2be9077f0b2dc89d388bc293fdbcd968274f43494c380d.
//
// Solidity: e RnodeQuit(who address)
func (_Rnode *RnodeFilterer) WatchRnodeQuit(opts *bind.WatchOpts, sink chan<- *RnodeRnodeQuit) (event.Subscription, error) {

	logs, sub, err := _Rnode.contract.WatchLogs(opts, "RnodeQuit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RnodeRnodeQuit)
				if err := _Rnode.contract.UnpackLog(event, "RnodeQuit", log); err != nil {
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

// RnodeOwnerRefundIterator is returned from FilterOwnerRefund and is used to iterate over the raw logs and unpacked data for OwnerRefund events raised by the Rnode contract.
type RnodeOwnerRefundIterator struct {
	Event *RnodeOwnerRefund // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log       // Log channel receiving the found contract events
	sub  gcchain.Subscription // Subscription for errors, completion and termination
	done bool                 // Whether the subscription completed delivering logs
	fail error                // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RnodeOwnerRefundIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RnodeOwnerRefund)
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
		it.Event = new(RnodeOwnerRefund)
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
func (it *RnodeOwnerRefundIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RnodeOwnerRefundIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RnodeOwnerRefund represents a OwnerRefund event raised by the Rnode contract.
type RnodeOwnerRefund struct {
	Who    common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOwnerRefund is a free log retrieval operation binding the contract event 0x3914ba80eb00486e7a58b91fb4795283df0c5b507eea9cf7c77cce26cc70d25c.
//
// Solidity: e ownerRefund(who address, amount uint256)
func (_Rnode *RnodeFilterer) FilterOwnerRefund(opts *bind.FilterOpts) (*RnodeOwnerRefundIterator, error) {

	logs, sub, err := _Rnode.contract.FilterLogs(opts, "ownerRefund")
	if err != nil {
		return nil, err
	}
	return &RnodeOwnerRefundIterator{contract: _Rnode.contract, event: "ownerRefund", logs: logs, sub: sub}, nil
}

// WatchOwnerRefund is a free log subscription operation binding the contract event 0x3914ba80eb00486e7a58b91fb4795283df0c5b507eea9cf7c77cce26cc70d25c.
//
// Solidity: e ownerRefund(who address, amount uint256)
func (_Rnode *RnodeFilterer) WatchOwnerRefund(opts *bind.WatchOpts, sink chan<- *RnodeOwnerRefund) (event.Subscription, error) {

	logs, sub, err := _Rnode.contract.WatchLogs(opts, "ownerRefund")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RnodeOwnerRefund)
				if err := _Rnode.contract.UnpackLog(event, "ownerRefund", log); err != nil {
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

// RnodeOwnerRefundAllIterator is returned from FilterOwnerRefundAll and is used to iterate over the raw logs and unpacked data for OwnerRefundAll events raised by the Rnode contract.
type RnodeOwnerRefundAllIterator struct {
	Event *RnodeOwnerRefundAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log       // Log channel receiving the found contract events
	sub  gcchain.Subscription // Subscription for errors, completion and termination
	done bool                 // Whether the subscription completed delivering logs
	fail error                // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *RnodeOwnerRefundAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RnodeOwnerRefundAll)
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
		it.Event = new(RnodeOwnerRefundAll)
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
func (it *RnodeOwnerRefundAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RnodeOwnerRefundAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RnodeOwnerRefundAll represents a OwnerRefundAll event raised by the Rnode contract.
type RnodeOwnerRefundAll struct {
	NumOfInvestor *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnerRefundAll is a free log retrieval operation binding the contract event 0xb65ebb6b17695b3a5612c7a0f6f60e649c02ba24b36b546b8d037e98215fdb8d.
//
// Solidity: e ownerRefundAll(numOfInvestor uint256)
func (_Rnode *RnodeFilterer) FilterOwnerRefundAll(opts *bind.FilterOpts) (*RnodeOwnerRefundAllIterator, error) {

	logs, sub, err := _Rnode.contract.FilterLogs(opts, "ownerRefundAll")
	if err != nil {
		return nil, err
	}
	return &RnodeOwnerRefundAllIterator{contract: _Rnode.contract, event: "ownerRefundAll", logs: logs, sub: sub}, nil
}

// WatchOwnerRefundAll is a free log subscription operation binding the contract event 0xb65ebb6b17695b3a5612c7a0f6f60e649c02ba24b36b546b8d037e98215fdb8d.
//
// Solidity: e ownerRefundAll(numOfInvestor uint256)
func (_Rnode *RnodeFilterer) WatchOwnerRefundAll(opts *bind.WatchOpts, sink chan<- *RnodeOwnerRefundAll) (event.Subscription, error) {

	logs, sub, err := _Rnode.contract.WatchLogs(opts, "ownerRefundAll")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RnodeOwnerRefundAll)
				if err := _Rnode.contract.UnpackLog(event, "ownerRefundAll", log); err != nil {
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

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a72305820b0c22c1f89a0edecb1f0595512d9688daf9d1af7bd246ce0b880bc0562666f260029`

// DeploySafeMath deploys a new gcchain contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an gcchain contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an gcchain contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an gcchain contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an gcchain contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an gcchain contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an gcchain contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an gcchain contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an gcchain contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an gcchain contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an gcchain contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}

// SetABI is the input ABI used to generate the binding from.
const SetABI = "[]"

// SetBin is the compiled bytecode used for deploying new contracts.
const SetBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a723058209dd4f9706744c33e1b4a8afbef333e72f055520a649897eaf114fa346462cd0e0029`

// DeploySet deploys a new gcchain contract, binding an instance of Set to it.
func DeploySet(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Set, error) {
	parsed, err := abi.JSON(strings.NewReader(SetABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SetBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Set{SetCaller: SetCaller{contract: contract}, SetTransactor: SetTransactor{contract: contract}, SetFilterer: SetFilterer{contract: contract}}, nil
}

// Set is an auto generated Go binding around an gcchain contract.
type Set struct {
	SetCaller     // Read-only binding to the contract
	SetTransactor // Write-only binding to the contract
	SetFilterer   // Log filterer for contract events
}

// SetCaller is an auto generated read-only Go binding around an gcchain contract.
type SetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SetTransactor is an auto generated write-only Go binding around an gcchain contract.
type SetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SetFilterer is an auto generated log filtering Go binding around an gcchain contract events.
type SetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SetSession is an auto generated Go binding around an gcchain contract,
// with pre-set call and transact options.
type SetSession struct {
	Contract     *Set              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SetCallerSession is an auto generated read-only Go binding around an gcchain contract,
// with pre-set call options.
type SetCallerSession struct {
	Contract *SetCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// SetTransactorSession is an auto generated write-only Go binding around an gcchain contract,
// with pre-set transact options.
type SetTransactorSession struct {
	Contract     *SetTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SetRaw is an auto generated low-level Go binding around an gcchain contract.
type SetRaw struct {
	Contract *Set // Generic contract binding to access the raw methods on
}

// SetCallerRaw is an auto generated low-level read-only Go binding around an gcchain contract.
type SetCallerRaw struct {
	Contract *SetCaller // Generic read-only contract binding to access the raw methods on
}

// SetTransactorRaw is an auto generated low-level write-only Go binding around an gcchain contract.
type SetTransactorRaw struct {
	Contract *SetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSet creates a new instance of Set, bound to a specific deployed contract.
func NewSet(address common.Address, backend bind.ContractBackend) (*Set, error) {
	contract, err := bindSet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Set{SetCaller: SetCaller{contract: contract}, SetTransactor: SetTransactor{contract: contract}, SetFilterer: SetFilterer{contract: contract}}, nil
}

// NewSetCaller creates a new read-only instance of Set, bound to a specific deployed contract.
func NewSetCaller(address common.Address, caller bind.ContractCaller) (*SetCaller, error) {
	contract, err := bindSet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SetCaller{contract: contract}, nil
}

// NewSetTransactor creates a new write-only instance of Set, bound to a specific deployed contract.
func NewSetTransactor(address common.Address, transactor bind.ContractTransactor) (*SetTransactor, error) {
	contract, err := bindSet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SetTransactor{contract: contract}, nil
}

// NewSetFilterer creates a new log filterer instance of Set, bound to a specific deployed contract.
func NewSetFilterer(address common.Address, filterer bind.ContractFilterer) (*SetFilterer, error) {
	contract, err := bindSet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SetFilterer{contract: contract}, nil
}

// bindSet binds a generic wrapper to an already deployed contract.
func bindSet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SetABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Set *SetRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Set.Contract.SetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Set *SetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Set.Contract.SetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Set *SetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Set.Contract.SetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Set *SetCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Set.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Set *SetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Set.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Set *SetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Set.Contract.contract.Transact(opts, method, params...)
}
