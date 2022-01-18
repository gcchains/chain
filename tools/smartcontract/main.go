

package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/tools/smartcontract/config"
	"github.com/gcchains/chain/tools/smartcontract/deploy"
	"github.com/ethereum/go-ethereum/common"
)

var (
	proxyCampaignContractAddress common.Address
)

func main() {
	log.Info("cmdline args", "args", os.Args)
	if len(os.Args) != 4 {
		fmt.Println("Usage: smartcontract <endpoint> <keystore path> <password>")
		return
	}
	config.SetConfig(os.Args[1], os.Args[2])

	password := os.Args[3]
	log.Info("start deploy contract")

	var wg1 sync.WaitGroup
	wg1.Add(1)
	go deployRNodeAdmissionAndCampaign(password, &wg1)
	wg1.Wait()

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go deployRpt(password, &wg2)

	wg2.Wait()

	var wg3 sync.WaitGroup
	wg3.Add(1)
	go deployNetwork(password, &wg3)
	wg3.Wait()

	deploy.UpdateCampaignParameters(password, proxyCampaignContractAddress, 5, 6)

	fmt.Println("======== init contract deploy completed=========")
}

func deployRNodeAdmissionAndCampaign(password string, wg *sync.WaitGroup) {
	// 1
	title := "[1.DeployRNode]"
	rnodeAddress := deploy.DeployRNode(password, 0)
	deploy.PrintContract(title, rnodeAddress)

	// 2
	title = "[2.DeployAdmission]"
	admissionAddress := deploy.DeployAdmission(password, 1)
	deploy.PrintContract(title, admissionAddress)

	// 3
	title = "[3.DeployCampaign]"
	campaignAddress := deploy.DeployCampaign(admissionAddress, rnodeAddress, password, 2)
	deploy.PrintContract(title, campaignAddress)
	proxyCampaignContractAddress = campaignAddress
	wg.Done()
}

func deployRpt(password string, wg *sync.WaitGroup) {
	// 4
	title := "[4.DeployRpt]"
	rptAddress := deploy.DeployRpt(password, 3)
	deploy.PrintContract(title, rptAddress)
	wg.Done()
}

func deployNetwork(password string, wg *sync.WaitGroup) {
	//5
	title := "[5.DeployNetwork]"
	networkAddress := deploy.DeployNetwork(password, 4)
	deploy.PrintContract(title, networkAddress)
	wg.Done()
}
