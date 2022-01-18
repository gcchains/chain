

package primitives

import (
	"math/big"

	"github.com/gcchains/chain/commons/log"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate abigen --sol ../contracts/primitive_contracts.sol --pkg contracts --out ../contracts/primitive_contracts.go

// Definitions of primitive contracts
//

type GetTxVolume struct {
	Backend RptPrimitiveBackend
}

func (c *GetTxVolume) RequiredGas(input []byte) uint64 {
	return configs.GetTxVolumeGas
}

func (c *GetTxVolume) Run(input []byte) ([]byte, error) {
	addr, number, err := extractRptPrimitivesArgs(input)
	if err != nil {
		log.Error("primitive_txvolume got error", "error", err)
		return common.LeftPadBytes(new(big.Int).Bytes(), 32), nil
	}
	log.Debug("primitive_txvolume", "addr", addr, "block number", number)

	// TODO: @AC get gcchain backend and read balance.
	txVolume, err := c.Backend.TxVolume(addr, number)
	if err != nil {
		log.Error("NewBasicCollector,error", "error", err)
	}
	ret := new(big.Int).SetInt64(int64(txVolume))
	return common.LeftPadBytes(ret.Bytes(), 32), nil
}
