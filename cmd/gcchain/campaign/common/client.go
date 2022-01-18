package common

import (
	"crypto/ecdsa"
	"os"
	"path/filepath"

	"github.com/gcchains/chain/accounts/keystore"
	"github.com/gcchains/chain/api/gcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// NewCpcClient new a gcc.client
func NewCpcClient(ep string, kspath string, password string) (*gcclient.Client, *ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address, error) {
	// Create client.
	client, err := gcclient.Dial(ep)
	if err != nil {
		return nil, nil, nil, [20]byte{}, err
	}
	// Open keystore file.
	file, err := os.Open(kspath)
	if err != nil {
		return nil, nil, nil, [20]byte{}, err
	}
	keyPath, err := filepath.Abs(filepath.Dir(file.Name()))
	if err != nil {
		return nil, nil, nil, [20]byte{}, err
	}
	// Create keystore and get account.
	kst := keystore.NewKeyStore(keyPath, 2, 1)
	account := kst.Accounts()[0]
	account, key, err := kst.GetDecryptedKey(account, password)
	if err != nil {
		return nil, nil, nil, [20]byte{}, err
	}
	// Get private and public keys.
	privateKey := key.PrivateKey
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, nil, [20]byte{}, err
	}

	// Get contractAddress.
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	// fmt.Println("from contractAddress:", fromAddress.Hex()) 
	return client, privateKey, publicKeyECDSA, fromAddress, err
}
