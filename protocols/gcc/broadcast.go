

package gcc

import (
	"time"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/consensus/dpos"
	"github.com/gcchains/chain/core"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

// BroadcastBlock will either propagate a block to a subset of it's peers, or
// will only announce it's availability (depending what's requested).
func (pm *ProtocolManager) BroadcastBlock(block *types.Block, propagate bool) {

	hash := block.Hash()
	peers := pm.peers.PeersWithoutBlock(hash)

	// If propagation is requested, send to a subset of the peer
	if propagate {
		if parent := pm.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1); parent == nil {
			log.Warn("Propagating dangling block", "number", block.Number(), "hash", hash.Hex())
			return
		}

		// Send the block to a subset of our peers
		// transfer := peers[:int(math.Sqrt(float64(len(peers))))]
		transfer := peers[:]

		for _, peer := range transfer {
			peer.AsyncSendNewBlock(block)
		}

		log.Debug("Propagated block", "number", block.NumberU64(), "hash", hash.Hex(), "recipients", len(transfer), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
		return
	}

	for _, peer := range peers {
		peer.AsyncSendNewBlockHash(block)
	}

	log.Debug("Propagated block hash and number announcement", "number", block.NumberU64(), "hash", hash.Hex(), "recipients", len(peers), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
	return
}

// BroadcastTxs will propagate a batch of transactions to all peers which are not known to
// already have the given transaction.
func (pm *ProtocolManager) BroadcastTxs(txs types.Transactions, force bool) {
	var txset = make(map[*peer]types.Transactions)

	// Broadcast transactions to a batch of peers not knowing about it
	txsLength := txs.Len()
	if force && txsLength > 0 {
		log.Debug("now rebroadcast waiting txs", "len", txsLength)
	}
	for _, tx := range txs {
		peers := pm.peers.PeersWithoutTx(tx.Hash())

		// if force, broadcast to all peers
		if force {
			peers = pm.peers.AllPeers()
		}

		for _, peer := range peers {
			txset[peer] = append(txset[peer], tx)
		}
		log.Debug("Broadcast transaction", "hash", tx.Hash().Hex(), "recipients", len(peers))
	}
	// FIXME include this again: peers = peers[:int(math.Sqrt(float64(len(peers))))]
	for peer, txs := range txset {
		peer.AsyncSendTransactions(txs)
	}
}

// Mined broadcast loop
func (pm *ProtocolManager) minedBroadcastLoop() {

	// automatically stops if unsubscribe
	for obj := range pm.minedBlockSub.Chan() {
		switch ev := obj.Data.(type) {
		case core.NewMinedBlockEvent:

			if pm.chainconfig.Dpos != nil && pm.engine.(*dpos.Dpos).Mode() == dpos.NormalMode {

				log.Debug("handling mined block with dpos handler", "number", ev.Block.NumberU64())

				// broadcast mined block with dpos handler
				pm.engine.(*dpos.Dpos).HandleMinedBlock(ev.Block)
			} else {
				pm.BroadcastBlock(ev.Block, true)
			}

		}
	}

}

func (pm *ProtocolManager) txBroadcastLoop() {
	for {
		select {
		case event := <-pm.txsCh:
			pm.BroadcastTxs(event.Txs, event.ForceBroadcast)

		// Err() channel will be closed when unsubscribing.
		case <-pm.txsSub.Err():
			return
		}
	}
}
