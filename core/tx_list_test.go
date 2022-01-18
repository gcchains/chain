// +build !race



package core

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Tests that transactions can be added to strict lists and list contents and
// nonce boundaries are correctly maintained.
func TestStrictTxListAdd(t *testing.T) {
	// Generate a list of transactions to insert
	key, _ := crypto.GenerateKey()

	txs := make(types.Transactions, 1024)
	for i := 0; i < len(txs); i++ {
		txs[i] = transaction(uint64(i), 0, key)
	}
	// Insert the transactions in a random order
	list := newTxList(true)
	for _, v := range rand.Perm(len(txs)) {
		list.Add(txs[v], DeprecatedDefaultTxPoolConfig.PriceBump)
	}
	// Verify internal state
	if len(list.txs.items) != len(txs) {
		t.Errorf("transaction count mismatch: have %d, want %d", len(list.txs.items), len(txs))
	}
	for i, tx := range txs {
		if list.txs.items[tx.Nonce()].Tx() != tx {
			t.Errorf("item %d: transaction mismatch: have %v, want %v", i, list.txs.items[tx.Nonce()], tx)
		}
	}
}

func generateTxSortedMapTestcase1() (*txSortedMap, time.Time, []types.Transactions) {
	// create a basetime
	baseTime, _ := time.Parse(time.RFC3339, "2019-01-01T00:00:00+00:00")

	var (
		timeIter = baseTime
		argTime  = baseTime.Add(2 * time.Minute)
		size     = 20               // transaction map size
		timeGap  = 10 * time.Second // time gap between txs
		expected = argTime.Sub(baseTime) / timeGap
	)

	// generate a series of transactions with updateTime increased 10 seconds one by one.
	items := make(map[uint64]*TimedTransaction)
	index := make(nonceHeap, size)
	for i := uint64(0); i < uint64(size); i++ {
		items[i] = &TimedTransaction{
			Transaction: types.NewTransaction(i, common.Address{}, big.NewInt(0), uint64(0), big.NewInt(0), []byte{}),
			updateTime:  timeIter,
		}
		index = append(index, i)
		timeIter = timeIter.Add(timeGap)
	}

	// expected result is the first 6 txs in one txs slice
	result := func() []types.Transactions {
		var results []types.Transactions
		var result []*types.Transaction
		for i := 0; i < int(expected); i++ {
			result = append(
				result,
				types.NewTransaction(uint64(i), common.Address{}, big.NewInt(0), uint64(0), big.NewInt(0), []byte{}),
			)
		}
		results = append(results, result)
		return results
	}()

	return &txSortedMap{
		items: items,
		index: &index,
	}, argTime, result
}

func Test_txSortedMap_AllBefore(t *testing.T) {
	type fields struct {
		items map[uint64]*TimedTransaction
		index *nonceHeap
		cache types.Transactions
	}
	type args struct {
		t time.Time
	}

	fakeTxSortedMap1, argTime1, result1 := generateTxSortedMapTestcase1()

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantResults []types.Transactions
	}{
		// TODO: Add test cases.
		{
			"test1",
			fields{
				fakeTxSortedMap1.items,
				fakeTxSortedMap1.index,
				fakeTxSortedMap1.cache,
			},
			args{
				argTime1,
			},
			result1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &txSortedMap{
				items: tt.fields.items,
				index: tt.fields.index,
				cache: tt.fields.cache,
			}
			if gotResults := m.AllBefore(tt.args.t); !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("txSortedMap.AllBefore() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}
