

package contracts

import (
	"context"
	"fmt"
	"math/big"

	"github.com/gcchains/chain/accounts/abi/bind"
	campaign "github.com/gcchains/chain/contracts/dpos/campaign"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate abigen --sol rpt.sol --pkg contracts --out rpt.go

//go:generate abigen --sol campaign/campaign.sol --pkg campaign --out campaign/campaign.go

// Backend wraps all methods required for campaign operation.
type Backend interface {
	bind.ContractBackend
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	// BalanceAt gets balance of specified account at the block specified by parameter blockNum. If blockNum is nil, will be the latest one.
	BalanceAt(ctx context.Context, address common.Address, blockNum *big.Int) (*big.Int, error)
}

type CampaignWrapper struct {
	*campaign.CampaignSession
	contractBackend bind.ContractBackend
}

func NewCampaignWrapper(transactOpts *bind.TransactOpts, contractAddr common.Address, contractBackend Backend) (*CampaignWrapper, error) {
	c, err := campaign.NewCampaign(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &CampaignWrapper{
		&campaign.CampaignSession{
			Contract:     c,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

func DeployCampaign(transactOpts *bind.TransactOpts, contractBackend Backend, admissionContractAddr common.Address, rewardContractAddr common.Address) (common.Address, *CampaignWrapper, error) {
	contractAddr, _, _, err := campaign.DeployCampaign(transactOpts, contractBackend, admissionContractAddr, rewardContractAddr)
	if err != nil {
		return contractAddr, nil, err
	}
	campaign, err := NewCampaignWrapper(transactOpts, contractAddr, contractBackend)
	if err != nil {
		return contractAddr, nil, err
	}

	return contractAddr, campaign, err
}

func (self *CampaignWrapper) MaximumNoc() (*big.Int, error) {
	fmt.Println("MaximumNoc is called")
	return self.Contract.MaxNoc(nil)
}

func (self *CampaignWrapper) ClaimCampaign(numOfCampaign *big.Int, cpuNonce uint64, cpuBlockNumber *big.Int,
	memoryNonce uint64, memoryBlockNumber *big.Int, version *big.Int) (*types.Transaction, error) {
	return self.Contract.ClaimCampaign(&self.TransactOpts, numOfCampaign, cpuNonce, cpuBlockNumber, memoryNonce, memoryBlockNumber, version)
}
