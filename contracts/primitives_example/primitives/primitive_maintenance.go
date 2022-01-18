

package primitives

import (
	"math/big"

	"github.com/gcchains/chain/commons/log"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate abigen --sol ../contracts/primitive_contracts.sol --pkg contracts --out ../contracts/primitive_contracts.go

// Definitions of primitive contracts

type GetMaintenance struct {
	Backend RptPrimitiveBackend
}

func (c *GetMaintenance) RequiredGas(input []byte) uint64 {
	return configs.GetMaintenanceGas
}

func (c *GetMaintenance) Run(input []byte) ([]byte, error) {
	addr, number, err := extractRptPrimitivesArgs(input)
	if err != nil {
		log.Error("primitive_maintenance got error", "error", err)
		return common.LeftPadBytes(new(big.Int).Bytes(), 32), nil
	}
	log.Debug("primitive_maintenance", "address", addr.Hex(), "number", number)

	// TODO: @AC get gcchain Backend and read balance.
	maintenance, err := c.Backend.Maintenance(addr, number)
	if err != nil {
		log.Error("NewBasicCollector,error", "error", err)
		return common.LeftPadBytes(new(big.Int).Bytes(), 32), nil
	}
	ret := new(big.Int).SetInt64(int64(maintenance))
	return common.LeftPadBytes(ret.Bytes(), 32), nil
}
