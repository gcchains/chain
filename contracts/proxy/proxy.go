

package proxy

import (
	"context"
	"math/big"

	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate abigen --sol proxycontract/proxy.sol --pkg contract --out proxycontract/proxy.go
//need generate in dir:contracts/proxy/proxycontract.

type Backend interface {
	bind.ContractBackend
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	BalanceAt(ctx context.Context, address common.Address, blockNum *big.Int) (*big.Int, error)
}
