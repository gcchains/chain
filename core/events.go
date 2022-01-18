

package core

import (
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

// NewTxsEvent is posted when a batch of transactions enter the transaction pool.
type NewTxsEvent struct {
	Txs            []*types.Transaction
	ForceBroadcast bool
}

// PendingLogsEvent is posted pre mining and notifies of pending logs.
type PendingLogsEvent struct {
	Logs []*types.Log
}

// PendingStateEvent is posted pre mining and notifies of pending state changes.
type PendingStateEvent struct{}

// NewMinedBlockEvent is posted when a block has been imported.
type NewMinedBlockEvent struct{ Block *types.Block }

// RemovedLogsEvent is posted when a reorg happens
type RemovedLogsEvent struct{ Logs []*types.Log }

type ChainEvent struct {
	Block *types.Block
	Hash  common.Hash
	Logs  []*types.Log
}

type ChainSideEvent struct {
	Block *types.Block
}

type ChainHeadEvent struct {
	Block *types.Block
}

type ChainLatestEvent struct {
	Block *types.Block
}

type InsertionStartEvent struct{}
type InsertionDoneEvent struct{}
