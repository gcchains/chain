

package keystore

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

const (
	veryLightScryptN = 2
	veryLightScryptP = 1
)

// Tests that a json key file can be decrypted and encrypted in multiple rounds.
func TestKeyEncryptDecrypt(t *testing.T) {
	keyjson, err := ioutil.ReadFile("testdata/very-light-scrypt.json")
	if err != nil {
		t.Fatal(err)
	}
	password := ""
	address := common.HexToAddress("45dea0fb0bba44f5fcf290bba72fd57d7117cbb8")

	// Do a few rounds of decryption and encryption
	for i := 0; i < 3; i++ {
		// Try a bad password first
		if _, err := DecryptKey(keyjson, password+"bad"); err == nil {
			t.Errorf("test %d: json key decrypted with bad password", i)
		}
		// Decrypt with the correct password
		key, err := DecryptKey(keyjson, password)
		if err != nil {
			t.Fatalf("test %d: json key failed to decrypt: %v", i, err)
		}
		if key.Address != address {
			t.Errorf("test %d: key address mismatch: have %x, want %x", i, key.Address, address)
		}
		// Recrypt with a new password and start over
		password += "new data appended"
		if keyjson, err = EncryptKey(key, password, veryLightScryptN, veryLightScryptP); err != nil {
			t.Errorf("test %d: failed to recrypt key %v", i, err)
		}
	}
}

func TestMergeBytes(t *testing.T) {
	b1, b2 := []byte("bbb1"), []byte("aaaccc2")
	b3 := mergeBytes(b1, b2)
	expected := []byte("bbb1" + "aaaccc2")
	if !reflect.DeepEqual(b3, expected) {
		t.Errorf("merge bytes error: have %x, want %x", b3, expected)
	}
}

func TestDecryptKeyTestnet(t *testing.T) {
	address := common.HexToAddress("2a15146f434c0105cfae639de2ac3bb543539b24")
	keyjson, err := ioutil.ReadFile("../../examples/gcchain/conf-testnet/keys/key1")
	if err != nil {
		t.Skip("file not found.skip.")
	}
	password := "123456!"
	key, err := DecryptKey(keyjson, password)
	if err != nil {
		t.Fatalf("test %d: json key failed to decrypt: %v", 1, err)
	}
	if key.Address != address {
		t.Errorf("test %d: key address mismatch: have %x, want %x", 1, key.Address, address)
	}

}

func TestDecryptKeyDev(t *testing.T) {
	address := common.HexToAddress("e94b7b6c5a1e526a4d97f9268ad6097bde25c62a")
	keyjson, err := ioutil.ReadFile("../../examples/gcchain/conf-dev/keys/key1")
	if err != nil {
		t.Fatal(err)
	}
	password := "password"
	key, err := DecryptKey(keyjson, password)
	if err != nil {
		t.Fatalf("test %d: json key failed to decrypt: %v", 1, err)
	}
	if key.Address != address {
		t.Errorf("test %d: key address mismatch: have %x, want %x", 1, key.Address, address)
	}
}
