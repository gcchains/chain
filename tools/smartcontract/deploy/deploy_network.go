

package deploy

import (
	"math/big"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/tools/smartcontract/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gcchains/chain/contracts/dpos/network"
)

func DeployNetwork(password string, nonce uint64) common.Address {
	client, err, privateKey, _, fromAddress := config.Connect(password)
	printBalance(client, fromAddress)
	// Launch contract deploy transaction.
	auth := newTransactor(privateKey, new(big.Int).SetUint64(nonce))
	contractAddress, tx, _, err := network.DeployNetwork(auth, client)
	if err != nil {
		log.Fatal(err.Error())
	}
	printTx(tx, err, client, contractAddress)
	return contractAddress
}
