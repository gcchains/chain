

package config

import (
	"crypto/ecdsa"
	"os"
	"path/filepath"

	"github.com/gcchains/chain/accounts/keystore"
	"github.com/gcchains/chain/api/gcclient"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/configs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	endPoint         = "http://localhost:8501"
	keyStoreFilePath = "./chain/examples/gcchain/data/data1/keystore/"
	// DefaultCPUDifficulty = uint64(19) // 1 cpu
	DefaultCPUDifficulty = uint64(12) // 1 cpu
	// DefaultMemoryDifficulty  = uint64(9) // 16 G
	DefaultMemoryDifficulty  = uint64(6) // 2 G
	DefaultCpuWorkTimeout    = uint64(5)
	DefaultMemoryWorkTimeout = uint64(5)
)

// overwrite from environment variables
func init() {
	if val := os.Getenv("gcchain_KEYSTORE_FILEPATH"); val != "" {
		keyStoreFilePath = val
	}
}

func SetConfig(ep, ksPath string) {
	endPoint = ep
	keyStoreFilePath = ksPath
}

func Connect(password string) (*gcclient.Client, error, *ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address) {
	ep, err := configs.ResolveUrl(endPoint)
	if err != nil {
		log.Fatal("unknown endpoint", "endpoint", endPoint, "err", err)
	}
	// Create client.
	client, err := gcclient.Dial(ep)
	if err != nil {
		log.Fatal(err.Error())
	}

	chainConfig, err := client.ChainConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	chainId, runMode := chainConfig.ChainID.Uint64(), configs.Mainnet
	switch chainId {
	case configs.DevChainId:
		runMode = configs.Dev
	case configs.MainnetChainId:
		runMode = configs.Mainnet
	case configs.TestMainnetChainId:
		runMode = configs.TestMainnet
	case configs.TestnetChainId:
		runMode = configs.Testnet
	default:
		log.Fatal("unknown chain id")
	}
	configs.SetRunMode(runMode)

	// Open keystore file.
	file, err := os.Open(keyStoreFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	keyPath, err := filepath.Abs(filepath.Dir(file.Name()))
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create keystore and get account.
	kst := keystore.NewKeyStore(keyPath, 2, 1)
	account := kst.Accounts()[0]
	account, key, err := kst.GetDecryptedKey(account, password)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Get private and public keys.
	privateKey := key.PrivateKey
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	// Get contractAddress.
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// fmt.Println("from contractAddress:", fromAddress.Hex()) // 0xe94b7b6c5a0e526a4d97f9768ad6097bde25c62a

	return client, err, privateKey, publicKeyECDSA, fromAddress
}
