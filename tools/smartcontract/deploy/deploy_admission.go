

package deploy

import (
	"math/big"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/contracts/dpos/admission"
	"github.com/gcchains/chain/tools/smartcontract/config"
	"github.com/ethereum/go-ethereum/common"
)

func DeployAdmission(password string, nonce uint64) common.Address {
	client, err, privateKey, _, fromAddress := config.Connect(password)
	printBalance(client, fromAddress)
	// Launch contract deploy transaction.
	// auth := newAuth(client, privateKey, fromAddress)
	auth := newTransactor(privateKey, new(big.Int).SetUint64(nonce))
	contractAddress, tx, _, err := admission.DeployAdmission(auth, client,
		new(big.Int).SetUint64(config.DefaultCPUDifficulty),
		new(big.Int).SetUint64(config.DefaultMemoryDifficulty),
		new(big.Int).SetUint64(config.DefaultCpuWorkTimeout),
		new(big.Int).SetUint64(config.DefaultMemoryWorkTimeout))
	if err != nil {
		log.Fatal(err.Error())
	}
	printTx(tx, err, client, contractAddress)
	return contractAddress
}
