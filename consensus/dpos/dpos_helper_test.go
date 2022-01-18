

package dpos

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	lru "github.com/hashicorp/golang-lru"
)

func Test_dposHelper_verifyHeader(t *testing.T) {
	t.Skip("need to redesign the unittests for dposHelper")

	dh := &defaultDposHelper{}

	extraErr1 := "0x00000000000000000000000000000000"
	fmt.Println("extraErr1:", extraErr1)

	rightExtra := "0x0000000000000000000000000000000000000000000000000000000000000000"
	rightSeal := "0xc9efd3956760d72613081c50294ad582d1e36bea45878f3570cc9e1525b997472120d0ef25f88c3b64122b967bd5063633b744bc4e3ae3afc316bb1e5c7edc1d00"
	rightAddr := "0xe94b7b1c5a0e526a4d97f1768ad6097bde25c61a"

	time1 := big.NewInt(time.Now().Unix() + 100)
	time := big.NewInt(time.Now().Unix() - 100)

	type args struct {
		c         *Dpos
		chain     consensus.ChainReader
		header    *types.Header
		parents   []*types.Header
		refHeader *types.Header
	}
	tests := []struct {
		name    string
		d       *defaultDposHelper
		args    args
		wantErr bool
	}{
		{"header.Number is nil", dh, args{header: &types.Header{Number: nil, Time: time1}}, true},

		{"header.Time error", dh, args{header: &types.Header{Number: big.NewInt(6),
			Time: time1}}, true},

		{"errInvalidCheckpointBeneficiary", dh,
			args{header: &types.Header{Number: big.NewInt(6), Time: time, Coinbase: common.HexToAddress("aaaaa")},
				c: &Dpos{config: &configs.DposConfig{TermLen: 3}}}, true},

		{"header.Extra error1", dh,
			args{
				header: &types.Header{
					Number: big.NewInt(5), Time: time, Extra: hexutil.MustDecode(string(extraErr1))},
				c: &Dpos{config: &configs.DposConfig{TermLen: 3}}}, true},

		{"errInvalidDifficulty", dh,
			args{
				header: &types.Header{
					Number: big.NewInt(7),
					Time:   time,
					Extra:  hexutil.MustDecode(string(rightExtra)),
					Dpos: types.DposSnap{
						Seal: types.HexToDposSig(rightSeal),
						Proposers: []common.Address{
							common.HexToAddress(rightAddr),
						},
					}},
				c: &Dpos{config: &configs.DposConfig{TermLen: 3}}}, true},

		{"success", dh,
			args{
				header: &types.Header{
					Number: big.NewInt(0),
					Time:   time,
					Extra:  hexutil.MustDecode(string(rightExtra)),
					Dpos: types.DposSnap{
						Seal: types.HexToDposSig(rightSeal),
						Proposers: []common.Address{
							common.HexToAddress(rightAddr),
						},
					},
				},
				c:       &Dpos{config: &configs.DposConfig{TermLen: 3}, dh: &defaultDposHelper{}},
				chain:   &FakeReader{},
				parents: []*types.Header{},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dh := &defaultDposHelper{}
			if err := dh.verifyHeader(tt.args.c, tt.args.chain, tt.args.header, tt.args.parents, tt.args.refHeader, true, false); (err != nil) != tt.wantErr {
				t.Errorf("defaultDposHelper.verifyHeader(%v, %v, %v, %v, %v) error = %v, wantErr %v", tt.args.c, tt.args.chain, tt.args.header, tt.args.parents, tt.args.refHeader, err, tt.wantErr)
			}
		})
	}
}

func Test_dposHelper_verifyCascadingFields(t *testing.T) {
	t.Skip("need to redesign the unittests for dposHelper")

	recents, _ := lru.NewARC(10)
	rightExtra := "0x0000000000000000000000000000000000000000000000000000000000000000"
	seal := "0xc1efd3956760d72613081c50294ad582d0e36bea45878f3510cc9e8525b997472120d0ef25f88c3b64122b167bd5063633b744bc4e3ae3afc316bb1e5c7edc1d00"
	proposer := "0x194b7b6c5a0e526a4d97f9718ad6091bde25c62a"
	time1 := big.NewInt(time.Now().Unix() - 100)
	time2 := big.NewInt(time.Now().Unix() + 100)
	header := &types.Header{Number: big.NewInt(0), Time: time1}
	parentHash := header.Hash()
	recents.Add(parentHash, &DposSnapshot{config: &configs.DposConfig{Period: 3, ViewLen: 3, TermLen: 3}, RecentProposers: make(map[uint64][]common.Address)})
	chain := &FakeReader{}
	type args struct {
		d         *Dpos
		chain     consensus.ChainReader
		header    *types.Header
		parents   []*types.Header
		refHeader *types.Header
	}
	tests := []struct {
		name    string
		d       *defaultDposHelper
		args    args
		wantErr bool
	}{
		{"success when block 0", &defaultDposHelper{},
			args{d: &Dpos{recentSnaps: recents, config: &configs.DposConfig{Period: 3, ViewLen: 3, TermLen: 4}},
				header: &types.Header{Number: big.NewInt(0), ParentHash: parentHash}, chain: chain}, false},
		{"fail with parent block", &defaultDposHelper{},
			args{d: &Dpos{recentSnaps: recents, config: &configs.DposConfig{Period: 3, ViewLen: 3, TermLen: 4}},
				header:  &types.Header{Number: big.NewInt(1), ParentHash: parentHash, Time: time1},
				parents: []*types.Header{header}, chain: chain}, true},
		{"errInvalidSigners", &defaultDposHelper{},
			args{d: &Dpos{recentSnaps: recents, config: &configs.DposConfig{Period: 3, ViewLen: 3, TermLen: 4}, dh: &defaultDposHelper{}},
				header: &types.Header{Number: big.NewInt(1), ParentHash: parentHash, Time: time2,
					Extra: hexutil.MustDecode(rightExtra), Dpos: types.DposSnap{Seal: types.HexToDposSig(seal),
						Proposers: []common.Address{common.HexToAddress(proposer)}},
				},
				parents: []*types.Header{header}, chain: chain}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultDposHelper{}
			if err := d.verifyProposers(tt.args.d, tt.args.chain, tt.args.header, tt.args.parents, tt.args.refHeader); (err != nil) != tt.wantErr {
				t.Errorf("defaultDposHelper.verifyProposers(%v, %v, %v, %v, %v) error = %v, wantErr %v", tt.args.d, tt.args.chain, tt.args.header, tt.args.parents, tt.args.refHeader, err, tt.wantErr)
			}
		})
	}
}

func Test_dposHelper_snapshot(t *testing.T) {
	type args struct {
		c       *Dpos
		chain   consensus.ChainReader
		number  uint64
		hash    common.Hash
		parents []*types.Header
	}
	tests := []struct {
		name    string
		d       *defaultDposHelper
		args    args
		want    *DposSnapshot
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &defaultDposHelper{}
			got, err := d.snapshot(tt.args.c, tt.args.chain, tt.args.number, tt.args.hash, tt.args.parents)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultDposHelper.Snapshot(%v, %v, %v, %v, %v) error = %v, wantErr %v", tt.args.c, tt.args.chain, tt.args.number, tt.args.hash, tt.args.parents, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultDposHelper.Snapshot(%v, %v, %v, %v, %v) = %v, want %v", tt.args.c, tt.args.chain, tt.args.number, tt.args.hash, tt.args.parents, got, tt.want)
			}
		})
	}
}

func Test_dposHelper_verifySeal(t *testing.T) {

	rightExtra := "0x0000000000000000000000000000000000000000000000000000000000000000"
	rightAddr := "0xe94b7b6c1a0e526a4d97f9761ad6097bde25c62a"
	rightSeal := "0x19efd3956760d72613081c50294ad582d0e36bea45878f3570cc9e8525b997472121d0ef25f88c3b64122b967bd5063633b714bc4e3ae3afc116bb4e5c7edc1d00"

	time1 := big.NewInt(time.Now().Unix() - 100)

	header := &types.Header{Number: big.NewInt(0), Time: time1}
	parentHash := header.Hash()
	recents, _ := lru.NewARC(10)
	recents.Add(parentHash, &DposSnapshot{})
	
	type args struct {
		c         *Dpos
		chain     consensus.ChainReader
		header    *types.Header
		parents   []*types.Header
		refHeader *types.Header
	}
	tests := []struct {
		name    string
		d       *defaultDposHelper
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"fail when block number is 0", &defaultDposHelper{},
			args{
				c: &Dpos{
					config:      &configs.DposConfig{Period: 3},
					db:          &fakeDb{1},
					recentSnaps: recents, dh: &defaultDposHelper{}},
				chain: &FakeReader{},
				header: &types.Header{
					Number: big.NewInt(0),
					Time:   time1,
					Extra:  hexutil.MustDecode(string(rightExtra)),
					Dpos: types.DposSnap{
						Proposers: []common.Address{
							common.HexToAddress(rightAddr),
						},
						Seal: types.HexToDposSig(rightSeal),
					},
				}},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.verifySeal(tt.args.c, tt.args.chain, tt.args.header, tt.args.parents, tt.args.refHeader); (err != nil) != tt.wantErr {
				t.Errorf("defaultDposHelper.verifySeal(%v, %v, %v, %v, %v) error = %v, wantErr %v", tt.args.c, tt.args.chain, tt.args.header, tt.args.parents, tt.args.refHeader, err, tt.wantErr)
			}
		})
	}
}
