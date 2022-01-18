

package pdash

import (
	"context"
	"fmt"
	"math/big"

	"github.com/gcchains/chain/accounts/abi/bind"
	pdash "github.com/gcchains/chain/contracts/pdash/pdash_contract"
	"github.com/gcchains/chain/contracts/proxy"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

type RegisterBackend interface {
	bind.ContractBackend
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	BalanceAt(ctx context.Context, address common.Address, blockNum *big.Int) (*big.Int, error)
}

type RegisterWrapper struct {
	*pdash.RegisterSession
	contractBackend bind.ContractBackend
}

func NewRegisterWrapper(transactOpts *bind.TransactOpts, contractAddr common.Address, contractBackend bind.ContractBackend) (*RegisterWrapper, error) {
	c, err := pdash.NewRegister(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &RegisterWrapper{
		&pdash.RegisterSession{
			Contract:     c,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

func DeployRegisterAndReturnWrapper(transactOpts *bind.TransactOpts, contractBackend RegisterBackend) (common.Address, *RegisterWrapper, *proxy.ProxyContractRegister, error) {
	proxyaddress, _, proxyContractRegister, err := proxy.DeployProxyContractRegister(transactOpts, contractBackend)
	contractAddr, _, _, err := pdash.DeployRegister(transactOpts, contractBackend, proxyaddress)
	if err != nil {
		return contractAddr, nil, nil, err
	}
	register, err := NewRegisterWrapper(transactOpts, contractAddr, contractBackend)
	if err != nil {
		return contractAddr, nil, nil, err
	}

	return contractAddr, register, proxyContractRegister, err

}

func (self *RegisterWrapper) ClaimRegister(transactOpts *bind.TransactOpts, fileName string, fileHash [32]byte, fileSize *big.Int) (*types.Transaction, error) {
	fmt.Println("ClainRegister is called:")
	return self.Contract.ClaimRegister(transactOpts, fileName, fileHash, fileSize)
}

func (self *RegisterWrapper) GetUploadCount(CallOpts *bind.CallOpts, address common.Address) (*big.Int, error) {
	fmt.Println("GetUploadCount is call:")
	return self.Contract.GetUploadCount(CallOpts, address)
}
