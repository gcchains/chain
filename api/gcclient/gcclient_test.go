

package gcclient_test

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"

	gcchain "/gcchain/chain"
	"github.com/gcchains/chain/api/gcclient"
	"github.com/ethereum/go-ethereum/common"
)

// Verify that Client implements the ethereum interfaces.
var (
	_ = gcchain.ChainReader(&gcclient.Client{})
	_ = gcchain.TransactionReader(&gcclient.Client{})
	_ = gcchain.ChainStateReader(&gcclient.Client{})
	_ = gcchain.ChainSyncReader(&gcclient.Client{})
	_ = gcchain.ContractCaller(&gcclient.Client{})
	_ = gcchain.GasEstimator(&gcclient.Client{})
	_ = gcchain.GasPricer(&gcclient.Client{})
	_ = gcchain.LogFilterer(&gcclient.Client{})
	_ = gcchain.PendingStateReader(&gcclient.Client{})
	// _ = ethereum.PendingStateEventer(&Client{})
	_ = gcchain.PendingContractCaller(&gcclient.Client{})
)

func TestGetRNodes(t *testing.T) {
	t.Skip("Must start chain to test")
	client, err := gcclient.Dial("http://localhost:8501")
	// local
	if err != nil {
		log.Fatal(err.Error())
	}
	rnodes, err := client.GetRNodes(context.Background())
	fmt.Println(rnodes, err)
	fmt.Println("rpt is :", "addr", rnodes[0].Address, "rpt", rnodes[0].Rpt, "status", rnodes[0].Status)

	if rnodes[0].Rpt == 0 {
		t.Errorf("GetRNodes failed")
	}
}

func TestGetCurrentTerm(t *testing.T) {
	t.Skip("must start chain to test")
	client, err := gcclient.Dial("http://localhost:8501")
	if err != nil {
		log.Fatal(err.Error())
	}
	currentTerm, err := client.GetCurrentTerm(context.Background())
	fmt.Println("currentTerm", currentTerm)

	if err != nil {
		t.Errorf("GetCurrentTerm failed")
	}
}

func TestGetCurrentView(t *testing.T) {
	t.Skip("must start chain to test")
	client, err := gcclient.Dial("http://localhost:8501")
	// local
	if err != nil {
		log.Fatal(err.Error())
	}
	currentView, err := client.GetCurrentView(context.Background())
	fmt.Println("currentTerm", currentView)

	if err != nil {
		t.Errorf("GetCurrentView failed")
	}
}

func TestGetBlockGenerationInfoList(t *testing.T) {
	t.Skip("must start chain to test")
	client, err := gcclient.Dial("http://localhost:8501")
	// local
	if err != nil {
		log.Fatal(err.Error())
	}
	blockGenInfo, err := client.GetBlockGenerationInfo(context.Background())
	fmt.Println("committee is:", blockGenInfo, len(blockGenInfo.Proposers))
	fmt.Println("blockGenInfo ", "ProposerIndex:", blockGenInfo.ProposerIndex, "Term :", blockGenInfo.Term, "BlockNumber :", blockGenInfo.BlockNumber, "Proposer", blockGenInfo.Proposer, "Porposers", blockGenInfo.Proposers)
	if len(blockGenInfo.Proposers) != 4 {
		t.Errorf("GetBlockGenerationInfoList failed")
	}
}

func TestGetCommitteeNumber(t *testing.T) {
	t.Skip("must start chain to test")
	client, err := gcclient.Dial("http://localhost:8501")
	if err != nil {
		log.Fatal(err.Error())
	}
	committeesNum, err := client.GetCommitteeNumber(context.Background())
	fmt.Println("committees is :", committeesNum)

	if committeesNum != 4 {
		t.Errorf("GetCommittee failed")
	}
}

func TestGetBlockReward(t *testing.T) {
	t.Skip("must start chain to test")
	client, err := gcclient.Dial("http://localhost:8501")
	if err != nil {
		log.Fatal(err.Error())
	}
	reward := client.GetBlockReward(context.Background(), new(big.Int).SetInt64(100))
	fmt.Println("block number", 100, "block reward", reward.Uint64())
}

func TestClient_BlockByNumber(t *testing.T) {
	t.Skip("must start chain to test")
	client, err := gcclient.Dial("http://localhost:8501")
	// local
	if err != nil {
		log.Fatal(err.Error())
	}
	block, err := client.BlockByNumber(context.Background(), big.NewInt(138))
	if err != nil {
		log.Fatal("BlockByNumber is error: ", err)
	}
	Number := block.Number()
	tx := block.Transactions()
	fmt.Println("block Transactions is :", tx)
	fmt.Println("the blcok number is :", Number)

}

func TestClient_ChainConfig(t *testing.T) {
	t.Skip("must start chain to test")
	client, err := gcclient.Dial("http://localhost:8501")
	// local
	if err != nil {
		log.Fatal(err.Error())
	}
	cfg, err := client.ChainConfig()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("chain config", "viewLen", cfg.Dpos.ViewLen, "termLen", cfg.Dpos.TermLen)
}

func TestGetProposerByBlock(t *testing.T) {
	t.Skip("must start chain to test")
	client, err := gcclient.Dial("http://localhost:8501")
	// local
	if err != nil {
		log.Fatal(err.Error())
	}
	cfg, err := client.GetProposerByBlock(context.Background(), big.NewInt(1))
	if err != nil {
		t.Fatal(err)
	}
	if cfg != common.HexToAddress("0xc05302acebd0730e3a18a048d7d1cb1202c4a092") {
		t.Fatal("wrong Proposer ", "we want ", common.HexToAddress("0xc05302acebd0730e3a18a048d7d1cb1202c4a092"), "we get ", cfg)
	}
}
