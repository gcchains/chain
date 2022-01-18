

package deploy

import (
	"context"
	"fmt"
	"testing"

	"github.com/gcchains/chain/api/gcclient"
	"github.com/gcchains/chain/commons/log"
	"github.com/ethereum/go-ethereum/common"
)

func TestMainnetInitContract(t *testing.T) {
	t.Skip("skip contract verify test")

	client, err := gcclient.Dial("http://localhost:8501")
	if err != nil {
		log.Fatal(err.Error())
		t.Skip("skip mainnet init verify contract")
	}
	ctx := context.Background()

	
	addr := common.HexToAddress("0xED111e097C77253bc0febBd7e85F5Af6D60A8196") // real

	code, err := client.CodeAt(ctx, addr, nil)
	if len(code) > 0 {
		fmt.Println("contract code exist")
	} else {
		fmt.Println("contract code not exist")
	}
	if err != nil {
		println("DeployContract failed")
		log.Fatal(err.Error())
	}
}

func TestDevInitContract(t *testing.T) {
	t.Skip("skip contract verify test")

	// client, err, _, _, fromAddress := config.Connect("password")
	client, err := gcclient.Dial("http://localhost:8501")
	if err != nil {
		log.Fatal(err.Error())
		t.Skip("skip dev init verify contract")
	}
	ctx := context.Background()

	
	devContractAddressMap := map[string]common.Address{
		"ContractProposer":   common.HexToAddress("0x310236762f31bf0f69f792bd9fb01b5c679aa3f1"),
		"ContractReward":     common.HexToAddress("0xd6E4BdC19b4D1744Cf16Ce90419EDE5e78751002"),
		"ContractAdmission":  common.HexToAddress("0x0DDf4057eeDFb81D58029Be49bab19bbc45bC500"),
		"ContractCampaign":   common.HexToAddress("0x82104907AA699b2982Fc46f38Fd8C915d03Cdb8d"),
		"ContractRpt":        common.HexToAddress("0x019cC01ff9d88529b9e58FF26bfc53B1E060e915"),
		"ContractRegister":   common.HexToAddress("0x1Aae743244a7A5116470df8BD398e7D562ae8881"),
		"ContractPdash":      common.HexToAddress("0xd81ab6B1e656551F90B2d8749261949fde97096D"),
		"ContractPdashProxy": common.HexToAddress("0x1791F193C2F374f49bCbc120750749b7AF17204e"),
	}

	for name, addr := range devContractAddressMap {
		fmt.Println("contract name:", name)
		fmt.Println("addr:", addr.Hex())
		code, err := client.CodeAt(ctx, addr, nil)
		if len(code) > 0 {
			fmt.Println("contract code exist")
		} else {
			fmt.Println("contract code not exist")
		}
		if err != nil {
			println("DeployContract failed:" + name)
			log.Fatal(err.Error())
		}
	}
}
