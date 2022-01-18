

package deploy

import (
	"math/big"

	"github.com/gcchains/chain/commons/log"
	campaign "github.com/gcchains/chain/contracts/dpos/campaign"
	"github.com/gcchains/chain/tools/smartcontract/config"
	"github.com/ethereum/go-ethereum/common"
)

func DeployCampaign(acAddr common.Address, rewardAddr common.Address, password string, nonce uint64) common.Address {
	client, err, privateKey, _, fromAddress := config.Connect(password)
	printBalance(client, fromAddress)

	// Launch contract deploy transaction.
	auth := newTransactor(privateKey, new(big.Int).SetUint64(nonce))
	contractAddress, tx, _, err := campaign.DeployCampaign(auth, client, acAddr, rewardAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	printTx(tx, err, client, contractAddress)
	return contractAddress
}

func UpdateCampaignParameters(password string, campaignContractAddr common.Address, nonce1 uint64, nonce2 uint64) {
	client, err, privateKey, _, _ := config.Connect(password)
	// get chain config
	cfg, err := client.ChainConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	auth := newTransactor(privateKey, new(big.Int).SetUint64(nonce1))
	campaign, _ := campaign.NewCampaign(campaignContractAddr, client)
	tx, err := campaign.UpdateTermLen(auth, new(big.Int).SetUint64(cfg.Dpos.TermLen))
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Info("updated term len", "txhash", tx.Hash().Hex(), "termLen", cfg.Dpos.TermLen)
}
