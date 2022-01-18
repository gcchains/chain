

package dpos

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gcchains/chain/accounts/keystore"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/database"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	lru "github.com/hashicorp/golang-lru"
)

func Test_sigHash(t *testing.T) {
	tx1 := types.NewTransaction(0, common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"), big.NewInt(10), 50000, big.NewInt(10), nil)
	tx1, _ = tx1.WithSignature(types.HomesteadSigner{}, common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70818d094f8a8fae537ce25ed8cb5af9adac3f1411f69bd515bd2ba031522df09b97dd72b100"))
	newHeader := &types.Header{
		ParentHash:   common.HexToHash("0x83cafc574e1f51ba9dc0568fc617a08ea2429fb384059c972f1119fa1c8dd55"),
		Coinbase:     common.HexToAddress("0x8888f1F195AFa192CfeE861698584c030f4c9dB1"),
		StateRoot:    common.HexToHash("0xef1552a40b7165c3cd773806b9e0c165b75356e0314bf07061279c729f51e017"),
		TxsRoot:      common.HexToHash("0x5fe50b260da6301036625b850b5d6ced6d0a9f814c0688bc91ffb7b7a3a54b67"),
		ReceiptsRoot: common.HexToHash("0xbc37d79753ad738a6dac4921e57392f145d8887476de3f783d1a7edae9283e52"),
		Number:       big.NewInt(1),
		GasLimit:     uint64(3141592),
		GasUsed:      uint64(21000),
		Time:         big.NewInt(1426516743),
		Extra:        common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000000"),
		Dpos: types.DposSnap{
			Seal: types.HexToDposSig("0x0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
			Sigs: []types.DposSignature{},
			Proposers: []common.Address{
				common.HexToAddress("0xe94b7b6c5a0e526a4d97f9768ad1097bde25c62a"),
				common.HexToAddress("0xc053021cebd0730e3a18a058d7d1cb1204c4a092"),
				common.HexToAddress("0x1f3dd127de235f15ffb4fc0d71469d1339df6465"),
			},
			Validators: []common.Address{},
		},
	}

	type args struct {
		header *types.Header
	}
	tests := []struct {
		name     string
		args     args
		wantHash common.Hash
	}{
		{"sigHash", args{newHeader}, common.HexToHash("0x22c08daa37af74c531f057987c3c9400c5c528518cb784fac2c7d5dfaa337c9c")},
	}

	dposUtil := &defaultDposUtil{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHash := dposUtil.sigHash(tt.args.header); !reflect.DeepEqual(gotHash, tt.wantHash) {
				t.Errorf("sigHash(%v) = %v, want %v", tt.args.header, gotHash.Hex(), tt.wantHash.Hex())
			}
		})
	}
}

func getAccount(keyStoreFilePath string, passphrase string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, common.Address) {
	// Load account.
	file, err := os.Open(keyStoreFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	keyPath, err := filepath.Abs(filepath.Dir(file.Name()))
	if err != nil {
		log.Fatal(err.Error())
	}

	kst := keystore.NewKeyStore(keyPath, 2, 1)

	// Get account.
	account := kst.Accounts()[0]
	account, key, err := kst.GetDecryptedKey(account, passphrase)
	if err != nil {
		log.Fatal(err.Error())
	}

	privateKey := key.PrivateKey
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	return privateKey, publicKeyECDSA, fromAddress
}

func getTestAccount() (common.Address, *ecdsa.PrivateKey) {
	privateKey, _ := crypto.HexToECDSA("1ad9c8855b740a0b7ed4c221dbad0f3a83a49cad6b3fe8d5111ac83d38b6a19")
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress, privateKey
}

func Test_ecrecover(t *testing.T) {

	addr, privKey := getTestAccount()

	tx1 := types.NewTransaction(0, common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af512d87"), big.NewInt(10), 50000, big.NewInt(10), nil)
	tx1, _ = tx1.WithSignature(types.HomesteadSigner{}, common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8156f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100"))
	cache, _ := lru.NewARC(10)

	newHeader := &types.Header{
		ParentHash:   common.HexToHash("0x83cafc574e1f51ba9dc0568fc617a08ea2429fb381059c972f13b19fa1c8dd55"),
		Coinbase:     common.HexToAddress("0x8888f1F195AFa192CfeE860691584c030f4c9dB1"),
		StateRoot:    common.HexToHash("0xef1552a40b7165c3cd773806b9e0c165b75156e0314bf0706f279c729f51e017"),
		TxsRoot:      common.HexToHash("0x5fe50b260da6301036625b850b5d6ced6d0a9f814c0688bc91ffb7b7a3a54b67"),
		ReceiptsRoot: common.HexToHash("0xbc37d79753ad738a6dac4921e57392f145d8887476de3f783dfa1edae9283e52"),
		Number:       big.NewInt(1),
		GasLimit:     uint64(3141592),
		GasUsed:      uint64(21000),
		Time:         big.NewInt(1426516743),
		Extra:        hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000"),
		Dpos: types.DposSnap{
			Proposers: []common.Address{
				addr,
			},
			Sigs:       make([]types.DposSignature, 3),
			Validators: []common.Address{},
		},
	}

	dph := &defaultDposHelper{&defaultDposUtil{}}
	hashBytes := dph.sigHash(newHeader).Bytes()
	hashBytesWithState, _ := hashBytesWithState(hashBytes, consensus.Commit)
	proposerSig, _ := crypto.Sign(hashBytes, privKey)
	validatorSig, _ := crypto.Sign(hashBytesWithState, privKey)

	copy(newHeader.Dpos.Seal[:], proposerSig[:])
	copy(newHeader.Dpos.Sigs[0][:], validatorSig[:])
	copy(newHeader.Dpos.Sigs[1][:], validatorSig[:])
	copy(newHeader.Dpos.Sigs[2][:], validatorSig[:])

	sigs := &signatures{
		sigs: make(map[common.Address][]byte),
	}
	sigs.setSig(
		addr,
		proposerSig,
	)

	existingCache, _ := lru.NewARC(10)
	fmt.Println("newHeader.Hash():", newHeader.Hash().Hex())
	existingCache.Add(newHeader.Hash(), sigs)

	dposUtil := &defaultDposUtil{}
	

	noSignerSigHeader := types.CopyHeader(newHeader)
	noSignerSigHeader.Dpos.Seal = types.DposSignature{}
	copy(noSignerSigHeader.Dpos.Sigs[0][:], validatorSig[:])
	copy(noSignerSigHeader.Dpos.Sigs[1][:], validatorSig[:])
	copy(noSignerSigHeader.Dpos.Sigs[2][:], validatorSig[:])

	type args struct {
		header   *types.Header
		sigcache *lru.ARCCache
	}
	tests := []struct {
		name    string
		args    args
		want    common.Address
		want1   []common.Address
		wantErr bool
	}{
		{"leaderSigHeader already cached,success", args{newHeader, cache}, addr, []common.Address{addr, addr, addr}, false},
		{"no signers' signatures. fail", args{noSignerSigHeader, cache}, common.Address{}, []common.Address{addr, addr, addr}, false},
		{"success", args{newHeader, cache}, addr, []common.Address{addr, addr, addr}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := dposUtil.ecrecover(tt.args.header, tt.args.sigcache)
			if (err != nil) != tt.wantErr {
				t.Errorf("ecrecover(%v, %v) error = %v, wantErr %v", tt.args.header, tt.args.sigcache, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ecrecover got = %v, want %v", got.Hex(), tt.want.Hex())
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				gotAddrs := []string{}
				for _, addr := range got1 {
					gotAddrs = append(gotAddrs, addr.Hex())
				}
				wantAddrs := []string{}
				for _, addr := range tt.want1 {
					wantAddrs = append(wantAddrs, addr.Hex())
				}
				t.Errorf("ecrecover got1 = %v, want %v", strings.Join(gotAddrs, ","),
					strings.Join(wantAddrs, ","))
			}
		})
	}
}

type fakeDBForSignedBlocksRecord struct {
	m map[string][]byte
}

func newFakeDBForSignedBlocksRecord() *fakeDBForSignedBlocksRecord {
	return &fakeDBForSignedBlocksRecord{
		m: make(map[string][]byte),
	}
}

func (f *fakeDBForSignedBlocksRecord) Put(key []byte, value []byte) error {
	fmt.Println("excuting put, key", key, "value", value)
	f.m[string(key)] = value
	return nil
}

func (f *fakeDBForSignedBlocksRecord) Delete(key []byte) error {
	panic("not implemented")
}

func (f *fakeDBForSignedBlocksRecord) Get(key []byte) ([]byte, error) {
	fmt.Println("excuting get, key", key)
	value, ok := f.m[string(key)]
	if ok {
		return value, nil
	}
	return nil, errors.New("no value")
}

func (f *fakeDBForSignedBlocksRecord) Has(key []byte) (bool, error) {
	panic("not implemented")
}

func (f *fakeDBForSignedBlocksRecord) Close() {
	panic("not implemented")
}

func (f *fakeDBForSignedBlocksRecord) NewBatch() database.Batch {
	panic("not implemented")
}

func Test_newSignedBlocksRecord(t *testing.T) {
	db := newFakeDBForSignedBlocksRecord()
	fsbr := newSignedBlocksRecord(db)

	number, hash := generateNH()
	fmt.Println(number, hash)

	fsbr.markAsSigned(number, hash)
	if h, ok := fsbr.ifAlreadySigned(number); h != hash || !ok {
		t.Error("hh", "hash", h, "want", hash, "ok", ok)
	}

	// TODO: add more tests here

}

// generate random number and hash
func generateNH() (number uint64, hash common.Hash) {
	number = uint64(time.Now().UnixNano())
	hash = common.BytesToHash(numberToBytes(number))
	return
}
