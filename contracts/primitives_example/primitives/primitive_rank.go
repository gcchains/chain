

package primitives

import (
	"math/big"

	"github.com/gcchains/chain/commons/log"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate abigen --sol ../contracts/primitive_contracts.sol --pkg contracts --out ../contracts/primitive_contracts.go

// Definitions of primitive contracts

type GetRank struct {
	Backend RptPrimitiveBackend
}

func (c *GetRank) RequiredGas(input []byte) uint64 {
	return configs.GetRankGas
}

func (c *GetRank) Run(input []byte) ([]byte, error) {
	addr, number, err := extractRptPrimitivesArgs(input)
	if err != nil {
		log.Error("primitive_rank got error", "error", err)
		return common.LeftPadBytes(new(big.Int).Bytes(), 32), nil
	}
	log.Debug("primitive_rank, address", "addr", addr.Hex(), "number", number)

	// TODO: @AC get gcchain Backend and read balance.
	coinAge, err := c.Backend.Rank(addr, number)
	if err != nil {
		log.Error("NewBasicCollector,error", "error", err, "addr", addr.Hex())
	}
	ret := new(big.Int).SetInt64(int64(coinAge))
	return common.LeftPadBytes(ret.Bytes(), 32), nil
}
