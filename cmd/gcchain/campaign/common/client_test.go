package common

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/gcchains/chain/api/gcclient"
	"github.com/gcchains/chain/api/rpc"
	"github.com/gcchains/chain/configs"
	"github.com/ethereum/go-ethereum/common"
)

var skip bool

func init() {
	skip = true
}

func buildClient(ctx *context.Context, t *testing.T) (*gcclient.Client, *ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address) {
	endPoint := "http://127.0.0.1:8523"
	keyStoreFilePath := "~/.gcchain/keystore/"
	password := "password"
	client, privateKey, publicKeyECDSA, fromAddress, err := NewCpcClient(endPoint, keyStoreFilePath, password)
	if err != nil {
		t.Log(err.Error())
	}
	return client, privateKey, publicKeyECDSA, fromAddress
}

func checkError(t *testing.T, err error) {
	if err != nil {
		t.Log(err)
	}
}

func TestNewCpcClient(t *testing.T) {
	if skip {
		t.Skip()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, _, _, fromAddress := buildClient(&ctx, t)
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		t.Log(err)
	}
	t.Log("Balance", new(big.Int).Div(balance, big.NewInt(configs.Gcc)))
}

func TestContractExist(t *testing.T) {
	if skip {
		t.Skip()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, _, _, _ := buildClient(&ctx, t)
	contracts := map[string]common.Address{
		"ContractReward":    common.HexToAddress("0x94576e35a55D6BbF3bB45120bC831a668557eF42"),
		"ContractCampaign":  common.HexToAddress("0x2404Bf355428523F8e52E68Df00A0521e413F98E"),
		"ContractAdmission": common.HexToAddress("0x6f01875F462CBBc956CB9C0396dE6053A31C9C99"),
	}
	for name, addr := range contracts {
		code, err := client.CodeAt(ctx, addr, nil)
		if len(code) > 0 {
			t.Log("contract " + name + " code exist")
		} else {
			t.Log("contract " + name + " code not exist")
		}
		if err != nil {
			t.Error("DeployContract failed: " + name)
		}
	}
}

func TestStatus(t *testing.T) {
	if skip {
		t.Skip()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Miner
	endPoint := "http://127.0.0.1:8503"
	c, err := rpc.DialContext(ctx, endPoint)
	checkError(t, err)
	var result bool
	c.CallContext(ctx, &result, "eth_mining")
	t.Log("Is Mining:", result)

	// RNode
	// see TestReward.

	// Proposer
	client, _, _, fromAddress := buildClient(&ctx, t)
	head, err := client.HeaderByNumber(ctx, nil)
	checkError(t, err)
	paddrs := head.Dpos.Proposers
	for index, addr := range paddrs {
		t.Log(index, addr.Hex())
		if fromAddress == addr {
			t.Log("Is Proposers")
		}
	}
}

func TestMining(t *testing.T) {
	if skip {
		t.Skip()
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	endPoint := "http://127.0.0.1:8503"
	c, err := rpc.DialContext(ctx, endPoint)
	checkError(t, err)

	// Start Mining
	err = c.CallContext(ctx, nil, "miner_start", 1)
	checkError(t, err)

	// check
	var result bool
	c.CallContext(ctx, &result, "eth_mining")
	t.Log("Is Mining:", result)

	// Stop Mining
	err = c.CallContext(ctx, nil, "miner_stop")
	checkError(t, err)

	// check
	c.CallContext(ctx, &result, "eth_mining")
	t.Log("Is Mining:", result)
	if result == true {
		t.Error("Expect false but true. Stop mining failed.")
	}
}
