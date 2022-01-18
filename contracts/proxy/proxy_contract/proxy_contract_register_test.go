

package contract

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/accounts/keystore"
	"github.com/gcchains/chain/api/gcclient"
	"github.com/gcchains/chain/commons/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestProxyContractRegister(t *testing.T) {
	t.Skip("we shall use a simulated backend.")

	client, privateKey, fromAddress, gasLimit, gasPrice, instance, ctx, contractAddress := deployProxyContract()

	registerProxyAndGet(t, privateKey, gasLimit, gasPrice, instance, ctx, client, fromAddress, contractAddress)
}

func deployProxyContract() (*gcclient.Client, *ecdsa.PrivateKey, common.Address, int, *big.Int, *ProxyContractRegister, context.Context, common.Address) {
	client, err := gcclient.Dial("http://localhost:8501")
	// local
	if err != nil {
		log.Fatal(err.Error())
	}
	file, _ := os.Open("../../../examples/gcchain/data/dd1/keystore/")
	keyPath, err := filepath.Abs(filepath.Dir(file.Name()))
	kst := keystore.NewKeyStore(keyPath, 2, 1)
	account := kst.Accounts()[0]
	account, key, err := kst.GetDecryptedKey(account, "password")
	privateKey := key.PrivateKey
	if err != nil {
		log.Fatal(err.Error())
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("from address:", fromAddress.Hex())
	bal, err := client.BalanceAt(context.Background(), fromAddress, nil)
	fmt.Println("bal:", bal)
	gasLimit := 3000000
	gasPrice, err := client.SuggestGasPrice(context.Background())
	fmt.Println("gasPrice:", gasPrice)
	if err != nil {
		log.Fatal(err.Error())
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(gasLimit)
	auth.GasPrice = gasPrice
	address, tx, instance, err := DeployProxyContractRegister(auth, client)
	if err != nil {
		log.Fatal(err.Error())
	}
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Println("receipt.Status:", receipt.Status)
	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())
	startTime := time.Now()
	fmt.Printf("TX start @:%s", time.Now())
	addressAfterMined, err := bind.WaitDeployed(ctx, client, tx)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("tx mining take time:%s\n", time.Since(startTime))
	if !bytes.Equal(address.Bytes(), addressAfterMined.Bytes()) {
		log.Fatalf("mined address :%s,before mined address:%s", addressAfterMined, address)
	}
	return client, privateKey, fromAddress, gasLimit, gasPrice, instance, ctx, address
}

func registerProxyAndGet(t *testing.T, privateKey *ecdsa.PrivateKey, gasLimit int, gasPrice *big.Int, instance *ProxyContractRegister, ctx context.Context, client *gcclient.Client, fromAddress, contractAddress common.Address) {

	// 2. register node2 public key on chain (claim campaign)
	fmt.Println("2.RegisterProxyContract")
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Value = big.NewInt(500)
	// in wei
	auth.GasLimit = uint64(gasLimit)
	// in units
	auth.GasPrice = gasPrice

	proxyAddr := common.HexToAddress("0x1a9fae75908752d1abf4dca45e1cac311c376290")
	realAddr := common.HexToAddress("0bd6a2b982f158e80a50beac02341871e2d016df")
	transaction, err := instance.RegisterProxyContract(auth, proxyAddr, realAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("TX:", transaction.Hash().Hex())
	startTime := time.Now()
	receipt, err := bind.WaitMined(ctx, client, transaction)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("tx mining take time:%s\n", time.Since(startTime))
	fmt.Println("receipt.Status:", receipt.Status)

	// 3. node1 get public key with other's(node2) address
	fmt.Println("3.GetRealContract")
	realAddrOnChain, err := instance.GetRealContract(nil, proxyAddr)

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("realAddr:", realAddrOnChain.Hex(),
		" of proxyAddr:", proxyAddr.Hex())

	if !bytes.Equal(realAddr.Bytes(), realAddrOnChain.Bytes()) {
		t.Errorf("expected enode:%v,got %v", realAddr, realAddrOnChain)
	}

	// register proxy contract address as proxy should be invalid
	proxyAddr = contractAddress
	realAddr = common.HexToAddress("0bd6a2b982fc58e80a20beac02349871e2d026df")
	transaction, err = instance.RegisterProxyContract(auth, proxyAddr, realAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("TX:", transaction.Hash().Hex())
	startTime = time.Now()
	receipt, err = bind.WaitMined(ctx, client, transaction)
	if err != nil {
		log.Fatalf("failed to deploy contact when mining :%v", err)
	}
	fmt.Printf("tx mining take time:%s\n", time.Since(startTime))
	fmt.Println("receipt.Status:", receipt.Status)

	// 4. GetRealContract
	fmt.Println("4.GetRealContract")
	realAddrOnChain, err = instance.GetRealContract(nil, proxyAddr)

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("realAddr:", realAddrOnChain.Hex(), "of proxyAddr:", proxyAddr.Hex())

	empty := common.Address{}
	if !bytes.Equal(empty.Bytes(), realAddrOnChain.Bytes()) {
		t.Errorf("expected enode:%v,got %v", empty, realAddrOnChain)
	}
}
