package primitive_register

import (
	"context"
	"math/big"

	gcchain "/gcchain/chain"
	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/contracts/dpos/primitive_backend"
	"github.com/gcchains/chain/contracts/primitives_example/primitives"
	"github.com/gcchains/chain/core/vm"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

// this ContractAPI only use read contract can't Write or Event filtering
type ContractAPI interface {
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call gcchain.CallMsg, blockNumber *big.Int) ([]byte, error)
	PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error)
	PendingCallContract(ctx context.Context, call gcchain.CallMsg) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, call gcchain.CallMsg) (gas uint64, err error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	FilterLogs(ctx context.Context, query gcchain.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, query gcchain.FilterQuery, ch chan<- types.Log) (gcchain.Subscription, error)
}

func RegisterPrimitiveContracts() {
	chainClient := GetChainClient()
	for addr, c := range MakePrimitiveContracts(chainClient, chainClient) {
		err := vm.RegisterPrimitiveContract(addr, c)
		if err != nil {
			log.Fatal("register primitive contract error", "error", err, "addr", addr)
		}
	}
}

func GetChainClient() *primitive_backend.ApiClient {
	return &primitive_backend.ApiClient{ChainBackend: primitive_backend.GetApiBackendHolderInstance().ChainBackend, ContractBackend: primitive_backend.GetApiBackendHolderInstance().ContractBackend}
}

func MakePrimitiveContracts(contractClient bind.ContractBackend, chainClient *primitive_backend.ApiClient) map[common.Address]vm.PrimitiveContract {
	contracts := make(map[common.Address]vm.PrimitiveContract)

	// we start from 100 to reserve enough space for upstream primitive contracts.
	RptEvaluator, err := primitives.NewRptEvaluator(contractClient, chainClient)
	if err != nil {
		log.Fatal("s.RptEvaluator is file")
	}
	contracts[common.BytesToAddress([]byte{100})] = &primitives.GetRank{Backend: RptEvaluator}
	contracts[common.BytesToAddress([]byte{101})] = &primitives.GetMaintenance{Backend: RptEvaluator}
	contracts[common.BytesToAddress([]byte{102})] = &primitives.GetProxyCount{Backend: RptEvaluator}
	contracts[common.BytesToAddress([]byte{103})] = &primitives.GetUploadReward{Backend: RptEvaluator}
	contracts[common.BytesToAddress([]byte{104})] = &primitives.GetTxVolume{Backend: RptEvaluator}
	contracts[common.BytesToAddress([]byte{105})] = &primitives.IsProxy{Backend: RptEvaluator}
	return contracts
}
