

package primitives

import (
	"math/big"

	"github.com/gcchains/chain/commons/log"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate abigen --sol ../contracts/primitive_contracts.sol --pkg contracts --out ../contracts/primitive_contracts.go

// Definitions of primitive contracts
//

// GetProxyCount returns the count of transactions processed by the proxy specified by given account address
type GetProxyCount struct {
	Backend RptPrimitiveBackend
}

func (c *GetProxyCount) RequiredGas(input []byte) uint64 {
	return configs.IsProxyGas
}

func (c *GetProxyCount) Run(input []byte) ([]byte, error) {
	addr, number, err := extractRptPrimitivesArgs(input)
	if err != nil {
		log.Warnf("primitive_proxy_count got error %v", err)
		return common.LeftPadBytes(new(big.Int).Bytes(), 32), nil
	}
	log.Debug("primitive_proxy_count", "address", addr.Hex(), "block number", number)

	// TODO: @AC get gcchain Backend and read balance.
	_, proxyCount, err := c.Backend.ProxyInfo(addr, number)
	if err != nil {
		log.Warn("NewBasicCollector,error", "error", err, "address", addr.Hex())
		return common.LeftPadBytes(new(big.Int).Bytes(), 32), nil
	}
	ret := new(big.Int).SetInt64(int64(proxyCount))
	return common.LeftPadBytes(ret.Bytes(), 32), nil
}

type IsProxy struct {
	Backend RptPrimitiveBackend
}

func (c *IsProxy) RequiredGas(input []byte) uint64 {
	return configs.IsProxyGas
}

func (c *IsProxy) Run(input []byte) ([]byte, error) {
	addr, number, err := extractRptPrimitivesArgs(input)
	if err != nil {
		log.Error("primitive_is_proxy got error", "error", err)
		return common.LeftPadBytes(new(big.Int).Bytes(), 32), nil
	}
	log.Debug("primitive_is_proxy", "address", addr.Hex(), "number", number)

	isProxy, _, err := c.Backend.ProxyInfo(addr, number)
	if err != nil {
		log.Error("NewBasicCollector,error", "error", err, "address", addr.Hex())
		ret := new(big.Int).SetInt64(int64(0))
		return common.LeftPadBytes(ret.Bytes(), 32), nil
	}
	ret := new(big.Int).SetInt64(int64(isProxy))
	return common.LeftPadBytes(ret.Bytes(), 32), nil
}
