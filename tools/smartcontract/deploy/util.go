

package deploy

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"

	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/api/gcclient"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

func printTx(tx *types.Transaction, err error, client *gcclient.Client, contractAddress common.Address) context.Context {
	ctx := context.Background()
	
	addressAfterMined, err := bind.WaitDeployed(ctx, client, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	// fmt.Printf("tx mining take time:%s\n", time.Since(startTime))
	if !bytes.Equal(contractAddress.Bytes(), addressAfterMined.Bytes()) {
		log.Fatalf("mined contractAddress :%s,before mined contractAddress:%s", addressAfterMined, contractAddress)
	}
	return ctx
}

func printBalance(client *gcclient.Client, fromAddress common.Address) {
	// Check balance.
	bal, _ := client.BalanceAt(context.Background(), fromAddress, nil)
	_ = bal
	// fmt.Println("balance:", bal)
}

func PrintContract(title string, address common.Address) {
	fmt.Println("================================================================")
	fmt.Printf(title+" Contract Address: 0x%x\n", address)
	fmt.Println("================================================================")
}

func FormatPrint(msg string) {
	fmt.Println("\n\n================================================================")
	fmt.Println(msg)
	fmt.Println("================================================================")
}

type nonceCounter struct {
	nonce uint64
	lock  sync.RWMutex
}

var nonceInstance *nonceCounter
var once sync.Once
var needInit = true

func GetNonceInstance(init uint64) *nonceCounter {
	once.Do(func() {
		nonceInstance = &nonceCounter{nonce: init}
		needInit = false
	})
	return nonceInstance
}

func (p *nonceCounter) GetNonce() uint64 {
	p.lock.RLock()
	defer p.lock.RUnlock()
	orig := p.nonce
	p.nonce = p.nonce + 1
	return orig
}

func newAuth(client *gcclient.Client, privateKey *ecdsa.PrivateKey, fromAddress common.Address) *bind.TransactOpts {
	auth := bind.NewKeyedTransactor(privateKey)
	newNonce := uint64(0)

	if needInit {
		blockNumber := client.GetBlockNumber()

		initNonce, err := client.NonceAt(context.Background(), fromAddress, blockNumber)
		if err != nil {
			fmt.Println("get nonce failed", err)
		}
		GetNonceInstance(initNonce)
	}
	newNonce = GetNonceInstance(0).GetNonce()
	fmt.Println("newNonce:", newNonce)
	auth.Nonce = new(big.Int).SetUint64(newNonce)
	return auth

}

func newTransactor(privateKey *ecdsa.PrivateKey, nonce *big.Int) *bind.TransactOpts {
	auth := bind.NewKeyedTransactor(privateKey)
	if nonce.Cmp(big.NewInt(-1)) > 0 {
		auth.Nonce = nonce
	}
	return auth

}
