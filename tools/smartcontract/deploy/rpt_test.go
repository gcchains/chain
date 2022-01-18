

package deploy

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/gcchains/chain/contracts/primitives_example/rpt"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/tools/smartcontract/config"
	"github.com/ethereum/go-ethereum/common"
)

func TestRpt(t *testing.T) {
	t.Skip("skip rpt integrate test")

	client, err, _, _, fromAddress := config.Connect("password")
	ctx := context.Background()
	printBalance(client, fromAddress)

	
	addr := common.HexToAddress("0x81104907aa691b2982fc46f38fd8c115d03cdb8d") 
	r, err := contracts.NewRpt(addr, client)

	code, err := client.CodeAt(ctx, addr, nil)
	if len(code) > 0 {
		fmt.Println("contract code exist")
	} else {
		fmt.Println("contract code not exist")
	}
	if err != nil {
		println("DeployRpt")
		log.Fatal(err.Error())
	}

	a, err := r.GetRpt(nil, common.HexToAddress("091e7baea6a6c7c4c1dfeb977efac326af152d85"), big.NewInt(0))
	if err != nil {
		println("GetRpt", "error:", err)
		log.Fatal(err.Error())
	}
	fmt.Println("rpt is :", a)

	b, err := r.GetRpt(nil, common.HexToAddress("0xE14B7b6C5A0e526A4D97f9768AD1097bdE25c61a"), big.NewInt(0))
	if err != nil {
		println("GetRpt", "error:", err)
		log.Fatal(err.Error())
	}
	fmt.Println("rpt is :", b)

	windowsize, err := r.Window(nil)
	if err != nil {
		log.Fatal("get windowzie is error")
	}
	fmt.Println("winodowsize is:", windowsize.Uint64())
}
