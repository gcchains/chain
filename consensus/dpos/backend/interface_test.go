package backend

import (
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

func TestValidMacSig(t *testing.T) {
	type args struct {
		mac string
		sig []byte
	}

	var (
		key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b2b9ee17ae16c6668d313eac2f26dbcda3f291")
		addr   = crypto.PubkeyToAddress(key.PublicKey)
		mac    = "gcchain|" + time.Now().Format(time.RFC3339)                   // correct
		mac2   = "gcchain|" + time.Now().Add(-time.Minute).Format(time.RFC3339) // wrong time
	)

	getSig := func(mac string) []byte {
		var hash common.Hash
		hasher := sha3.NewKeccak256()
		hasher.Write([]byte(mac))
		hasher.Sum(hash[:0])
		sig, _ := crypto.Sign(hash.Bytes(), key)
		return sig
	}

	tests := []struct {
		name       string
		args       args
		wantValid  bool
		wantSigner common.Address
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			"test1", // succeed
			args{
				mac,
				getSig(mac),
			},
			true,
			addr,
			false,
		},
		{
			"test2", // wrong timestamp
			args{
				mac2,
				getSig(mac2),
			},
			false,
			common.Address{},
			false,
		},
		{
			"test3", // wrong signature
			args{
				mac,
				[]byte{},
			},
			false,
			common.Address{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValid, gotSigner, err := ValidMacSig(tt.args.mac, tt.args.sig)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidMacSig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotValid != tt.wantValid {
				t.Errorf("ValidMacSig() gotValid = %v, want %v", gotValid, tt.wantValid)
			}
			if !reflect.DeepEqual(gotSigner, tt.wantSigner) {
				t.Errorf("ValidMacSig() gotSigner = %v, want %v", gotSigner, tt.wantSigner)
			}
		})
	}
}

func TestIsCheckPoint(t *testing.T) {
	type args struct {
		number  uint64
		termLen uint64
		viewLen uint64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			"succeed",
			args{
				12,
				4,
				3,
			},
			true,
		},
		{
			"fail",
			args{
				11,
				4,
				3,
			},
			false,
		},
		{
			"edge1",
			args{
				0,
				4,
				3,
			},
			false,
		},

		{
			"edge2",
			args{
				12,
				4,
				0,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCheckPoint(tt.args.number, tt.args.termLen, tt.args.viewLen); got != tt.want {
				t.Errorf("IsCheckPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
