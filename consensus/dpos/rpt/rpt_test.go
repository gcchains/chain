

package rpt_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus/dpos"
	"github.com/gcchains/chain/consensus/dpos/rpt"
	"github.com/gcchains/chain/core"
	"github.com/gcchains/chain/core/vm"
	"github.com/gcchains/chain/database"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

// --- --- --- --- --- --- --- --- --- --- --- below are bench tests for new rpt calculation method --- --- --- --- --- --- --- --- --- --- ---

func newBlockchain(n int) *core.BlockChain {
	db := database.NewMemDatabase()
	remoteDB := database.NewIpfsDbWithAdapter(database.NewFakeIpfsAdapter())
	gspec := core.DefaultGenesisBlock()
	genesis := gspec.MustCommit(db)
	config := gspec.Config
	dposConfig := config.Dpos
	dposFakeEngine := dpos.NewFaker(dposConfig, db)
	blocks, _ := core.GenerateChain(config, genesis, dposFakeEngine, db, remoteDB, n, nil)
	blockchain, _ := core.NewBlockChain(db, nil, gspec.Config, dposFakeEngine, vm.Config{}, remoteDB, nil)
	_, _ = blockchain.InsertChain(blocks)
	return blockchain
}

func newBlockchainWithBalances(n int, accounts []common.Address) *core.BlockChain {
	db := database.NewMemDatabase()
	remoteDB := database.NewIpfsDbWithAdapter(database.NewFakeIpfsAdapter())
	gspec := core.DefaultGenesisBlock()

	for i, addr := range accounts {
		gspec.Alloc[addr] = core.GenesisAccount{Balance: new(big.Int).Mul(big.NewInt(int64(i/3)), big.NewInt(configs.Gcc))}

		log.Info("alloc", "addr", addr.Hex(), "i", i, "bal", i/3)
	}

	genesis := gspec.MustCommit(db)
	config := gspec.Config
	dposConfig := config.Dpos
	dposFakeEngine := dpos.NewFaker(dposConfig, db)
	blocks, _ := core.GenerateChain(config, genesis, dposFakeEngine, db, remoteDB, n, nil)
	blockchain, _ := core.NewBlockChain(db, nil, gspec.Config, dposFakeEngine, vm.Config{}, remoteDB, nil)
	_, _ = blockchain.InsertChain(blocks)
	return blockchain
}

type fakeChainBackendForRptCollector struct {
	blockchain *core.BlockChain
}

func newFakeChainBackendForRptCollector(n int) *fakeChainBackendForRptCollector {
	bc := newBlockchain(n)
	return &fakeChainBackendForRptCollector{
		blockchain: bc,
	}
}

func newFakeChainBackendForRptCollectorWithBalances(n int, accounts []common.Address) *fakeChainBackendForRptCollector {
	bc := newBlockchainWithBalances(n, accounts)
	return &fakeChainBackendForRptCollector{
		blockchain: bc,
	}
}

func (fc *fakeChainBackendForRptCollector) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return fc.blockchain.GetBlockByNumber(number.Uint64()).Header(), nil
}

func (fc *fakeChainBackendForRptCollector) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	header, err := fc.HeaderByNumber(ctx, blockNumber)
	state, err := fc.blockchain.StateAt(header.StateRoot)
	if state == nil || err != nil {
		return nil, err
	}
	return state.GetBalance(account), state.Error()
}

func (fc *fakeChainBackendForRptCollector) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	header, err := fc.HeaderByNumber(ctx, blockNumber)
	state, err := fc.blockchain.StateAt(header.StateRoot)
	if state == nil || err != nil {
		return 0, err
	}
	return state.GetNonce(account), state.Error()
}

func generateABatchAccounts(n int) []common.Address {
	var addresses []common.Address
	for i := 1; i < n; i++ {
		addresses = append(addresses, common.HexToAddress("0x"+fmt.Sprintf("%040x", i)))
	}
	return addresses
}

func BenchmarkRptOf_10a(b *testing.B) {
	benchRptOf(b, 10)
}

func BenchmarkRptOf_100a(b *testing.B) {
	benchRptOf(b, 100)
}

func BenchmarkRptOf_200a(b *testing.B) {
	benchRptOf(b, 200)
}

func BenchmarkRptOf_1000a(b *testing.B) {
	benchRptOf(b, 1000)
}

func benchRptOf(b *testing.B, numAccount int) {
	b.ReportAllocs()
	b.ResetTimer()

	addrs := generateABatchAccounts(numAccount)
	fc := newFakeChainBackendForRptCollector(1000)

	rptCollector := rpt.NewRptCollectorImpl6(nil, fc)
	for i, addr := range addrs {
		rpt := rptCollector.RptOf(addr, addrs, 500)
		b.Log("idx", i, "rpt", rpt.Rpt, "addr", addr.Hex())
	}
}

func TestRptOf4(t *testing.T) {

	numAccount := 30
	numBlocks := 1000
	accounts := generateABatchAccounts(numAccount)
	fc := newFakeChainBackendForRptCollectorWithBalances(numBlocks, accounts)

	rptCollector := rpt.NewRptCollectorImpl6(nil, fc)
	for i, addr := range accounts {
		rpt := rptCollector.RptOf(addr, accounts, 500)
		t.Log("idx", i, "rpt", rpt.Rpt, "addr", addr.Hex())
	}

}
func TestRptOf5(t *testing.T) {

	numAccount := 30
	numBlocks := 1000
	accounts := generateABatchAccounts(numAccount)
	fc := newFakeChainBackendForRptCollectorWithBalances(numBlocks, accounts)

	rptCollector := rpt.NewRptCollectorImpl6(nil, fc)
	for i, addr := range accounts {
		rpt := rptCollector.RptOf(addr, accounts, 500)
		t.Log("idx", i, "rpt", rpt.Rpt, "addr", addr.Hex())
	}

}

func Test_PctCount(t *testing.T) {
	type args struct {
		pct   int
		total int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				pct:   50,
				total: 8,
			},
			want: 4,
		},
		{
			name: "2",
			args: args{
				pct:   10,
				total: 8,
			},
			want: 0,
		},
		{
			name: "3",
			args: args{
				pct:   0,
				total: 8,
			},
			want: 0,
		},
		{
			name: "4",
			args: args{
				pct:   100,
				total: 8,
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rpt.PctCount(tt.args.pct, tt.args.total); got != tt.want {
				t.Errorf("pctCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
