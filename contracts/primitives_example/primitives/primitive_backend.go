

package primitives

// this package collects all reputation calculation related information,
// then calculates the reputations of candidates.

import (
	"context"
	"math/big"
	"sort"

	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/api/rpc"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/contracts/dpos/campaign"
	"github.com/gcchains/chain/contracts/dpos/primitive_backend"
	pdash "github.com/gcchains/chain/contracts/pdash/pdash_contract"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate abigen --sol contracts/primitive_contracts_inst.sol --pkg contracts --out contracts/primitive_contracts_inst.go

const (
	Created = iota
	Finished
)

const (
	defaultRank = 100 // 100 represent give the address a default rank
)

type RptPrimitiveBackend interface {
	// Rank returns the rank for given account address at the given block number.
	Rank(address common.Address, number uint64) (int64, error)

	// TxVolume returns the transaction volumn for given account address at the given block number.
	TxVolume(address common.Address, number uint64) (int64, error)

	// Maintenance returns the maintenance score for given account address at the given block number.
	Maintenance(address common.Address, number uint64) (int64, error)

	// UploadCount returns the upload score for given account address at the given block number.
	UploadCount(address common.Address, number uint64) (int64, error)

	// ProxyInfo returns a value indicating whether the given address is proxy and the count of transactions processed
	// by the proxy at the given block number.
	ProxyInfo(address common.Address, number uint64) (isProxy int64, proxyCount int64, err error)
}

type RptEvaluator struct {
	ContractClient bind.ContractBackend
	ChainClient    *primitive_backend.ApiClient
}

func NewRptEvaluator(contractClient bind.ContractBackend, chainClient *primitive_backend.ApiClient) (*RptEvaluator, error) {
	bc := &RptEvaluator{
		ContractClient: contractClient,
		ChainClient:    chainClient,
	}
	return bc, nil
}

func getBalanceAt(ctx context.Context, apiBackend primitive_backend.ChainAPIBackend, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	state, _, err := apiBackend.StateAndHeaderByNumber(ctx, rpc.BlockNumber(blockNumber.Uint64()), false)
	if state == nil || err != nil {
		return common.Big0, err
	}
	return state.GetBalance(account), state.Error()
}

// Rank is the func to get rank to rpt
func (re *RptEvaluator) Rank(address common.Address, number uint64) (int64, error) {
	var balances []float64
	myBalance, err := getBalanceAt(context.Background(), re.ChainClient.ChainBackend, address, big.NewInt(int64(number)))
	if err != nil {
		log.Warn("error with getReputationnode", "error", err)
		return defaultRank, err
	}
	contractAddress := configs.ChainConfigInfo().Dpos.Contracts[configs.ContractCampaign]
	intance, err := campaign.NewCampaign(contractAddress, re.ContractClient)
	if err != nil {
		log.Error("NewCampaign error", "error", err, "contractAddress", contractAddress.Hex())
		return defaultRank, err
	}
	// get the rnode in that block
	term := (number - 1) / (configs.ChainConfigInfo().Dpos.TermLen * configs.ChainConfigInfo().Dpos.ViewLen)
	rNodeAddress, err := intance.CandidatesOf(nil, big.NewInt(int64(term)))
	if err != nil || rNodeAddress == nil {
		log.Error("CandidatesOf error", "error", err, "contractAddress", contractAddress.Hex())
		return defaultRank, err
	}
	for _, committee := range rNodeAddress {
		balance, err := getBalanceAt(context.Background(), re.ChainClient.ChainBackend, committee, big.NewInt(int64(number)))
		if err != nil {
			log.Error("error with bc.BalanceAt", "error", err, "contractAddress", contractAddress.Hex())
			return defaultRank, err
		}
		balances = append(balances, float64(balance.Uint64()))
	}
	var rank int64
	sort.Sort(sort.Reverse(sort.Float64Slice(balances)))
	index := sort.SearchFloat64s(balances, float64(myBalance.Uint64()))
	rank = int64(float64(index) / float64(len(rNodeAddress)) * 100) // solidity can't use float,so we magnify rank 100 times
	return rank, err
}

// TxVolume is the func to get txVolume to rpt
func (re *RptEvaluator) TxVolume(address common.Address, number uint64) (int64, error) {
	block, err := re.ChainClient.BlockByNumber(context.Background(), big.NewInt(int64(number)))
	if block == nil || err != nil {
		log.Error("error with bc.getTxVolume", "error", err)
		return 0, err
	}
	txvs := int64(0)
	signer := types.NewCep1Signer(configs.ChainConfigInfo().ChainID)
	txs := block.Transactions()
	for _, tx := range txs {
		sender, err := signer.Sender(tx)
		if err != nil {
			continue
		}
		if sender == address {
			txvs += 1
		}
	}
	return txvs, nil
}

// leader:0,committee:1,rNode:2,nil:3
func (re *RptEvaluator) Maintenance(address common.Address, number uint64) (int64, error) {
	ld := int64(2)
	if configs.ChainConfigInfo().ChainID.Uint64() == uint64(4) {
		return 0, nil
	}
	block, err := re.ChainClient.BlockByNumber(context.Background(), big.NewInt(int64(number)))
	if block == nil || err != nil {
		log.Error("error with bc.getIfLeader", "error", err)
		return 0, err
	}
	header := block.Header()
	leader := header.Coinbase

	log.Debug("leader.Hex is ", "hex", leader.Hex())

	if leader == address {
		ld = 0
	} else {
		for _, committee := range header.Dpos.Proposers {
			if address == committee {
				ld = 1
			}
		}
	}
	return ld, nil
}

// UploadCount is the func to get uploadnumber to rpt
func (re *RptEvaluator) UploadCount(address common.Address, number uint64) (int64, error) {
	uploadNumber := int64(0)
	contractAddress := configs.ChainConfigInfo().Dpos.Contracts[configs.ContractRegister]
	upload, err := pdash.NewRegister(contractAddress, re.ContractClient)
	if err != nil {
		log.Error("NewRegister error", "error", err, "address", address.Hex(), "contractAddress", contractAddress.Hex())
		return uploadNumber, err
	}
	fileNumber, err := upload.GetUploadCount(nil, address)
	if err != nil {
		log.Error("GetUploadCount error", "error", err, "address", address.Hex(), "contractAddress", contractAddress.Hex())
		return uploadNumber, err
	}
	return fileNumber.Int64(), err
}

// ProxyInfo func return the node is proxy or not
func (re *RptEvaluator) ProxyInfo(address common.Address, number uint64) (int64, int64, error) {
	isProxy := int64(0)
	contractAddress := configs.ChainConfigInfo().Dpos.Contracts[configs.ContractPdashProxy]
	proxyInstance, err := pdash.NewPdashProxy(contractAddress, re.ContractClient)

	if err != nil {
		log.Error("NewPdashProxy error", "error", err, "address", address.Hex(), "contractAddress", contractAddress.Hex())
		return 0, 0, err
	}

	proxy, err := proxyInstance.GetProxyFileNumber(nil, address)
	if err != nil {
		log.Error("GetProxyFileNumber error", "error", err, "address", address.Hex(), "contractAddress", contractAddress.Hex())
		return 0, 0, err
	}

	if proxy.Uint64() == 0 {
		return isProxy, 0, nil
	}

	proxyCount, err := proxyInstance.GetProxyFileNumberInBlock(nil, address, big.NewInt(int64(number)))
	if err != nil {
		log.Error("GetProxyFileNumber error", "error", err, "address", address.Hex(), "contractAddress", contractAddress.Hex())
		return 0, 0, err
	}

	return isProxy, proxyCount.Int64(), err
}

// func (re *RptEvaluator) CommitteeMember(header *types.Header) []common.Address {
//	committee := make([]common.Address, len(header.Dpos.Proposers))
//	for i := 0; i < len(committee); i++ {
//		copy(committee[i][:], header.Dpos.Proposers[i][:])
//	}
//	return committee
// }
