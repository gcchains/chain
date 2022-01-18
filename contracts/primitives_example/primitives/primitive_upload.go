

package primitives

import (
	"math/big"

	"github.com/gcchains/chain/commons/log"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate abigen --sol ../contracts/primitive_contracts.sol --pkg contracts --out ../contracts/primitive_contracts.go

// Definitions of primitive contracts

type GetUploadReward struct {
	Backend RptPrimitiveBackend
}

func (c *GetUploadReward) RequiredGas(input []byte) uint64 {
	return configs.GetUploadRewardGas
}

func (c *GetUploadReward) Run(input []byte) ([]byte, error) {
	addr, number, err := extractRptPrimitivesArgs(input)
	if err != nil {
		log.Error("primitive_uploadreward got error ", "error", err)
		return common.LeftPadBytes(new(big.Int).Bytes(), 32), nil
	}
	log.Debug("primitive_uploadreward, address", "addr", addr.Hex(), "number", number)

	// TODO: @AC get gcchain backend and read balance.
	uploadReward, err := c.Backend.UploadCount(addr, number)
	if err != nil {
		log.Error("NewBasicCollector,error", "error", err)
	}
	ret := new(big.Int).SetInt64(int64(uploadReward))
	return common.LeftPadBytes(ret.Bytes(), 32), nil
}
