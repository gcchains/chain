package primitive_register

import (
	"context"
	"math/big"

	gcchain "/gcchain/chain"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/contracts/dpos/primitives"
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
	for addr, c := range MakePrimitiveContracts() {
		err := vm.RegisterPrimitiveContract(addr, c)
		if err != nil {
			log.Fatal("register primitive contract error", "error", err, "addr", addr)
		}
	}
}

func MakePrimitiveContracts() map[common.Address]vm.PrimitiveContract {
	contracts := make(map[common.Address]vm.PrimitiveContract)

	contracts[common.BytesToAddress([]byte{106})] = &primitives.CpuPowValidate{}
	contracts[common.BytesToAddress([]byte{107})] = &primitives.MemPowValidate{}
	return contracts
}
