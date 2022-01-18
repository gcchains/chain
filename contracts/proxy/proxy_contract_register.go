

package proxy

import (
	"math/big"

	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/contracts/proxy/proxy_contract"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate abigen --sol contract/proxyContractRegister.sol --pkg contract --out contract/proxyContractRegister.go

type ProxyContractRegister struct {
	*contract.ProxyContractRegisterSession
	contractBackend bind.ContractBackend
}

func NewProxyContractRegister(transactOpts *bind.TransactOpts, contractAddr common.Address, contractBackend Backend) (*ProxyContractRegister, error) {
	c, err := contract.NewProxyContractRegister(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &ProxyContractRegister{
		&contract.ProxyContractRegisterSession{
			Contract:     c,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

func DeployProxyContractRegister(transactOpts *bind.TransactOpts, contractBackend Backend) (common.Address, *types.Transaction, *ProxyContractRegister, error) {
	contractAddr, tx, _, err := contract.DeployProxyContractRegister(transactOpts, contractBackend)
	if err != nil {
		return contractAddr, tx, nil, err
	}
	register, err := NewProxyContractRegister(transactOpts, contractAddr, contractBackend)
	if err != nil {
		return contractAddr, tx, nil, err
	}

	return contractAddr, tx, register, err
}

func (self *ProxyContractRegister) GetRealContract(addr common.Address) (common.Address, error) {
	realAddress, err := self.Contract.GetRealContract(&self.CallOpts, addr)
	if err != nil {
		return common.Address{}, err
	}
	log.Info("address:%v,realAddress:%v", addr, realAddress.Hex())
	return realAddress, err
}

func (self *ProxyContractRegister) RegisterPublicKey(proxyAddress, realAddress common.Address) (*types.Transaction, error) {
	self.TransactOpts.GasLimit = 300000
	self.TransactOpts.Value = big.NewInt(500)
	return self.Contract.RegisterProxyContract(&self.TransactOpts, proxyAddress, realAddress)
}
