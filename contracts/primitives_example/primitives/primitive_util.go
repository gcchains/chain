package primitives

import (
	"math/big"
	"sort"
)

func calcuateRank(myBalance *big.Int, balances []float64) int64 {

	index := sort.SearchFloat64s(balances, float64(myBalance.Uint64()))
	blockNumber := configs.ChainConfigInfo().Dpos.ViewLen * configs.ChainConfigInfo().Dpos.TermLen
	rank := int64((uint64(index) / blockNumber) * 100) // solidity can't use float
	return rank
}
