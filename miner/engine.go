
package miner

import (
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gcchains/chain/accounts"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/core"
	"github.com/gcchains/chain/core/state"
	"github.com/gcchains/chain/core/vm"
	"github.com/gcchains/chain/database"
	"github.com/gcchains/chain/protocols/gcc/syncer"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
)

// delayBeforeSeal returns 10% of the given delay duration
func delayBeforeSeal(d time.Duration) time.Duration {
	return d * 1 / 10
}

const (
	resultQueueSize = 10
	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096
	// chainLatestChanSize is the size of channel listening to ChainLatestEvent.
	chainLatestChanSize = 10
	// chainSideChanSize is the size of channel listening to ChainSideEvent.
	chainSideChanSize = 10
)

// Worker can register itself with the engine
type Worker interface {
	// retrieve the channel to pass work to a worker
	Work() chan<- *Work
	// retrieve result from a worker
	SetReturnCh(chan<- *Result)
	Stop()
	Start()
}

// Work is the workers current environment and holds
// all of the current state information
type Work struct {
	config *configs.ChainConfig // manifest which chain we are on
	signer types.Signer         // the signer, e.g., cep1signer to recover the transaction sender

	privState *state.StateDB          // apply public state changes here
	pubState  *state.StateDB          // apply private state changes here
	remoteDB  database.RemoteDatabase // ipfs database used for private tx processing
	tcount    int                     // tx count in cycle
	gasPool   *core.GasPool           // available gas used to pack transactions

	Block *types.Block // the new block

	header       *types.Header
	txs          []*types.Transaction
	pubReceipts  []*types.Receipt
	privReceipts []*types.Receipt

	accm      *accounts.Manager
	createdAt time.Time
}

type Result struct {
	Work  *Work
	Block *types.Block
}

// engine is the main object which takes care of applying messages to the new state
type engine struct {
	mu sync.Mutex

	config *configs.ChainConfig
	cons   consensus.Engine

	// update loop
	mux            *event.TypeMux
	txsCh          chan core.NewTxsEvent // new transactions enter the transaction pool
	txsSub         event.Subscription
	chainLatestCh  chan core.ChainLatestEvent // a new latest block has been inserted into the chain
	chainLatestSub event.Subscription
	chainSideCh    chan core.ChainSideEvent // a side block has been inserted
	chainSideSub   event.Subscription
	quitCh         chan struct{}

	workers map[Worker]struct{} // set of workers
	recv    chan *Result        // the channel that receives the result from workers

	backend Backend          // gcchain service backend
	chain   *core.BlockChain // pointer cuz we have only one canonical chain in the whole system
	proc    core.Validator
	chainDb database.Database

	coinbase common.Address
	extra    []byte

	currentMu   sync.RWMutex
	currentWork *Work

	snapshotMu    sync.RWMutex
	snapshotBlock *types.Block
	snapshotState *state.StateDB

	// atomic status counters
	mining    int32
	atWork    int32
	lastBlock uint64
}

func newEngine(config *configs.ChainConfig, cons consensus.Engine, coinbase common.Address, backend Backend, mux *event.TypeMux) *engine {
	e := &engine{
		config:        config,
		cons:          cons,
		backend:       backend,
		mux:           mux,
		txsCh:         make(chan core.NewTxsEvent, txChanSize),
		chainLatestCh: make(chan core.ChainLatestEvent, chainLatestChanSize),
		chainSideCh:   make(chan core.ChainSideEvent, chainSideChanSize),
		quitCh:        make(chan struct{}),
		chainDb:       backend.ChainDb(),
		recv:          make(chan *Result, resultQueueSize),
		chain:         backend.BlockChain(),
		proc:          backend.BlockChain().Validator(), // processor validator lock
		coinbase:      coinbase,
		workers:       make(map[Worker]struct{}),
	}

	// initially commit new work to make pending block and snapshot availableklk
	if backend.BlockChain().SyncMode() == syncer.FullSync {
		go e.update()
		e.commitNewWork()
	}
	return e
}

func (e *engine) setCoinbase(addr common.Address) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.coinbase = addr
}

func (e *engine) setExtra(extra []byte) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.extra = extra
}

func (e *engine) pending() (*types.Block, *state.StateDB) {
	if atomic.LoadInt32(&e.mining) == 0 {
		// return a snapshot to avoid contention on currentMu mutex
		e.snapshotMu.RLock()
		defer e.snapshotMu.RUnlock()
		return e.snapshotBlock, e.snapshotState.Copy()
	}

	e.currentMu.RLock()
	defer e.currentMu.RUnlock()
	return e.currentWork.Block, e.currentWork.pubState.Copy()
}

func (e *engine) pendingBlock() *types.Block {
	if atomic.LoadInt32(&e.mining) == 0 {
		// return a snapshot to avoid contention on currentMu mutex
		e.snapshotMu.RLock()
		defer e.snapshotMu.RUnlock()
		return e.snapshotBlock
	}

	e.currentMu.RLock()
	defer e.currentMu.RUnlock()
	return e.currentWork.Block
}

func (e *engine) start() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if atomic.LoadInt32(&e.mining) == 0 {
		// spin up workers
		for worker := range e.workers {
			worker.Start()
		}

		go e.wait()
	}

	atomic.StoreInt32(&e.mining, 1)
}

func (e *engine) stop() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if atomic.LoadInt32(&e.mining) == 1 {
		for worker := range e.workers {
			worker.Stop()
		}
		close(e.quitCh)
		e.quitCh = make(chan struct{})
	}
	atomic.StoreInt32(&e.mining, 0)
	atomic.StoreInt32(&e.atWork, 0)
}

func (e *engine) register(worker Worker) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.workers[worker] = struct{}{}
	worker.SetReturnCh(e.recv)
}

func (e *engine) unregister(worker Worker) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.workers, worker)
	worker.Stop()
}

// update dispatches blocks.
func (e *engine) update() {
	// Subscribe NewTxsEvent for tx pool
	e.txsSub = e.backend.TxPool().SubscribeNewTxsEvent(e.txsCh)
	// Subscribe events for blockchain
	e.chainLatestSub = e.backend.BlockChain().SubscribeChainLatestEvent(e.chainLatestCh)
	e.chainSideSub = e.backend.BlockChain().SubscribeChainSideEvent(e.chainSideCh)

	defer e.chainSideSub.Unsubscribe()
	defer e.chainLatestSub.Unsubscribe()
	defer e.txsSub.Unsubscribe()

	for {
		// A real event arrived, process interesting content

		select {
		// a new block has been inserted.  we start to mine based on this new tip.
		case <-e.chainLatestCh:

			log.Debug("now to commit new work", "now", time.Now())

			// commitNewWork must run no matter if it is mining, because pending block needs to be updated by commitNewWork
			e.commitNewWork()

			if atomic.LoadInt32(&e.mining) == 1 {
				// checks and tries to campaign if needed
				e.cons.TryCampaign()
			}

		// handle chainsideevent
		// we don't have uncle blocks
		case ev := <-e.chainSideCh:
			log.Warn("Got unexpected uncle block ", "hash", ev.Block.Hash().Hex())

		// Handle NewTxsEvent
		// add to the work (transaction set).  it's mainly for api use, e.g., pending block.
		case ev := <-e.txsCh:
			// Apply transactions to the pending state if we're *not* mining.
			// Note all transactions received may be compatible with transactions
			// already included in the current mining block. These transactions will
			// be automatically eliminated.
			if atomic.LoadInt32(&e.mining) == 0 {
				// critical section for current work
				e.currentMu.Lock()

				txs := make(map[common.Address]types.Transactions)
				for _, tx := range ev.Txs {
					// get the sender account
					acc, _ := types.Sender(e.currentWork.signer, tx)
					txs[acc] = append(txs[acc], tx)
				}
				txset := types.NewTransactionsByPriceAndNonce(e.currentWork.signer, txs)
				e.currentWork.commitTransactions(e.mux, txset, e.chain, e.coinbase, time.Now().Add(time.Second*10))
				e.updateSnapshot()
				e.currentMu.Unlock()
			}
		// System stopped
		case err := <-e.txsSub.Err():
			if err == nil {
				log.Info("system is stopped")
				return
			}
			log.Warn("txsSub got error", "error", err)
		case err := <-e.chainLatestSub.Err():
			if err == nil {
				log.Info("system is stopped")
				return
			}
			log.Warn("chainLatestSub got error", "error", err)
		case err := <-e.chainSideSub.Err():
			if err == nil {
				log.Info("system is stopped")
				return

			}
			log.Warn("chainSideSub got error", "error", err)
		}
	}

}

// wait handles mined blocks.
func (e *engine) wait() {
	for {
		select {
		case result := <-e.recv:
			atomic.AddInt32(&e.atWork, -1)

			if result == nil {
				continue
			}
			block := result.Block

			if block.NumberU64() <= e.lastBlock {
				// ignore the same height new block or old block
				log.Warn("detected duplicate blocks with same block number", "number", block.NumberU64())
				continue
			}
			e.lastBlock = block.NumberU64()

			// broadcast the block and announce chain insertion event
			_ = e.mux.Post(core.NewMinedBlockEvent{Block: block})

		case <-e.quitCh:
			log.Info("goroutine wait() quit")
			return
		}
	}
}

// push sends a new work task to currently live miner workers.
func (e *engine) push(work *Work) {
	if atomic.LoadInt32(&e.mining) != 1 {
		return
	}
	for worker := range e.workers {
		atomic.AddInt32(&e.atWork, 1)
		if ch := worker.Work(); ch != nil {
			ch <- work
		}
	}
}

// makeCurrentWork creates a new environment for the current cycle.
func (e *engine) makeCurrentWork(parent *types.Block, header *types.Header) error {
	pubState, err := e.chain.StateAt(parent.StateRoot())
	if err != nil {
		return err
	}

	privState, err := e.chain.StatePrivAt(parent.StateRoot())
	if err != nil {
		return err
	}

	work := &Work{
		config:    e.config,
		signer:    types.NewCep1Signer(e.config.ChainID),
		pubState:  pubState,
		privState: privState,
		header:    header,
		createdAt: time.Now(),
		remoteDB:  e.chain.RemoteDB(),
		accm:      e.backend.AccountManager(),
	}

	// Keep track of transactions which return errors so they can be removed
	work.tcount = 0
	e.currentWork = work
	return nil
}

// commitNewWork creates a new block. Calling this function multiple times will abort the previous work on workers.
func (e *engine) commitNewWork() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.currentMu.Lock()
	defer e.currentMu.Unlock()

	parent := e.chain.CurrentBlock() // the head of the blockchain
	tstart := time.Now()
	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     big.NewInt(0).SetUint64(num.Uint64() + 1),
		GasLimit:   core.CalcGasLimit(parent),
		Extra:      e.extra,
	}
	// Only set the coinbase if we are mining (avoid spurious block rewards)
	if atomic.LoadInt32(&e.mining) == 1 {
		header.Coinbase = e.coinbase
	}
	if err := e.cons.PrepareBlock(e.chain, header); err != nil {
		log.Error("Failed to prepare header for mining", "err", err)
		return
	}

	// delay to add more txs to tx pool
	delay := header.Timestamp().Sub(time.Now())
	delay = delayBeforeSeal(delay)
	log.Debug("Waiting for slot to seal", "delay", delay)

	select {
	case <-time.After(delay):
		log.Debug("now to make work and try to seal", "delay", delay)
	}

	err := e.makeCurrentWork(parent, header)
	if err != nil {
		log.Error("Failed to create mining context", "err", err)
		return
	}
	// create the current work task and check any fork transitions needed
	// note, there is no transaction in this block
	work := e.currentWork

	// we now populate the work with pending transactions
	pending, err := e.backend.TxPool().Pending()
	if err != nil {
		log.Error("Failed to fetch pending transactions", "err", err)
		return
	}
	txs := types.NewTransactionsByPriceAndNonce(e.currentWork.signer, pending)

	// break early at header.timestamp - delayBeforeSeal
	// timeline  ------------------------------------------
	//             |        9/10              |    1/10   |
	//           now               commitTxsBreakTime   header.timestamp
	//
	// timeline  ------------------------------------------
	log.Debug("timelog header.timestamp and now", "header.timestamp", header.Timestamp(), "now", time.Now(), "delay", header.Timestamp().Sub(time.Now()))
	delay = header.Timestamp().Sub(time.Now())
	delay = delayBeforeSeal(delay)
	commitTxsBreakTime := header.Timestamp()
	if delay > 0 {
		commitTxsBreakTime = header.Timestamp().Add(-delay)
	}

	log.Debug("timelog before commit txs", "header.timestamp", header.Timestamp(), "now", time.Now(), "delay", header.Timestamp().Sub(time.Now()), "commitTxsBreakTime", commitTxsBreakTime)

	work.commitTransactions(e.mux, txs, e.chain, e.coinbase, commitTxsBreakTime)

	log.Debug("timelog after commit txs", "header.timestamp", header.Timestamp(), "now", time.Now(), "delay", header.Timestamp().Sub(time.Now()))

	// Create the new block to seal with the consensus engine. Private tx's receipts are not involved computing block's
	// receipts hash and receipts bloom as they are private and not guaranteeing identical in different nodes.
	// Finalize will reward the coinbase.
	if work.Block, err = e.cons.Finalize(e.chain, header, work.pubState, work.txs, []*types.Header{}, work.pubReceipts); err != nil {
		log.Error("Failed to finalize block for sealing", "err", err)
		return
	}

	log.Debug("timelog after finalize", "header.timestamp", header.Timestamp(), "now", time.Now(), "delay", header.Timestamp().Sub(time.Now()))

	// We only care about logging if we're actually mining.
	if atomic.LoadInt32(&e.mining) == 1 {
		// only seal and broadcast the block when it is mining proposer
		if e.cons.CanMakeBlock(e.chain, e.coinbase, parent.Header()) {
			log.Debug("timelog pushing", "header.timestamp", header.Timestamp(), "now", time.Now(), "delay", header.Timestamp().Sub(time.Now()))
			e.push(work)
			log.Info("Commit new mining work", "number", work.Block.Number(), "hash", work.Block.Hash().Hex(), "txs", work.tcount, "elapsed", common.PrettyDuration(time.Since(tstart)))
		}
	}
	e.updateSnapshot()
}

func (e *engine) updateSnapshot() {
	e.snapshotMu.Lock()
	defer e.snapshotMu.Unlock()

	e.snapshotBlock = types.NewBlock(
		e.currentWork.header,
		e.currentWork.txs,
		e.currentWork.pubReceipts,
	)
	e.snapshotState = e.currentWork.pubState.Copy()
	// TODO: if need to snapshot private state?
}

// transactions are applied in ascending nonce order of each account.
func (w *Work) commitTransactions(mux *event.TypeMux, txs *types.TransactionsByPriceAndNonce, bc *core.BlockChain, coinbase common.Address, breakTimer time.Time) {
	if w.gasPool == nil {
		w.gasPool = new(core.GasPool).AddGas(w.header.GasLimit)
	}

	var coalescedLogs []*types.Log

	for {

		// If break timer is up, break now
		if time.Now().After(breakTimer) {
			break
		}

		// If we don't have enough gas for any further transactions then we're done
		if w.gasPool.Gas() < configs.TxGas {
			log.Debug("Not enough gas for further transactions", "have", w.gasPool, "want", configs.TxGas)
			break
		}
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			break
		}
		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		//
		from, _ := types.Sender(w.signer, tx)

		// Start executing the transaction
		w.pubState.Prepare(tx.Hash(), common.Hash{}, w.tcount)
		w.privState.Prepare(tx.Hash(), common.Hash{}, w.tcount)

		err, logs := w.commitTransaction(tx, bc, coinbase, w.gasPool)
		switch err {
		case core.ErrGasLimitReached:
			// Pop the current out-of-gas transaction without shifting in the next from the account
			log.Debug("Gas limit exceeded for current block", "sender", from)
			// we remove all the transactions associated with the account.
			txs.Pop()

		case core.ErrNonceTooLow:
			// New head notification data race between the transaction pool and miner, shift
			log.Debug("Skipping transaction with low nonce", "sender", from, "nonce", tx.Nonce())
			txs.Shift()

		case core.ErrNonceTooHigh:
			// Reorg notification data race between the transaction pool and miner, skip account =
			log.Debug("Skipping account with high nonce", "sender", from, "nonce", tx.Nonce())
			txs.Pop()

		case nil:
			// Everything ok, collect the logs and shift in the next transaction from the same account
			coalescedLogs = append(coalescedLogs, logs...)
			w.tcount++
			txs.Shift()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			log.Debug("Transaction failed, account skipped", "hash", tx.Hash(), "err", err)
			txs.Shift()
		}
	}

	if len(coalescedLogs) > 0 || w.tcount > 0 {
		// make a copy, the state caches the logs and these logs get "upgraded" from pending to mined
		// logs by filling in the block hash when the block was mined by the local miner. This can
		// cause a race condition if a log was "upgraded" before the PendingLogsEvent is processed.
		cpy := make([]*types.Log, len(coalescedLogs))
		for i, l := range coalescedLogs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		go func(logs []*types.Log, tcount int) {
			if len(logs) > 0 {
				// do we need pending log event?
				mux.Post(core.PendingLogsEvent{Logs: logs})
			}
			if tcount > 0 {
				mux.Post(core.PendingStateEvent{})
			}
		}(cpy, w.tcount)
	}
}

func (w *Work) commitTransaction(tx *types.Transaction, bc *core.BlockChain, coinbase common.Address, gp *core.GasPool) (error, []*types.Log) {
	snap := w.pubState.Snapshot()
	snapPriv := w.privState.Snapshot()

	pubReceipt, privReceipt, _, err := core.ApplyTransaction(w.config, bc, &coinbase, gp, w.pubState, w.privState, w.remoteDB,
		w.header, tx, &w.header.GasUsed, vm.Config{}, w.accm)
	if err != nil {
		w.pubState.RevertToSnapshot(snap)
		w.privState.RevertToSnapshot(snapPriv)
		return err, nil
	}
	w.txs = append(w.txs, tx)
	w.pubReceipts = append(w.pubReceipts, pubReceipt)
	if privReceipt != nil {
		w.privReceipts = append(w.privReceipts, privReceipt)
	}

	// TODO: consider whether append private logs to returned logs together.
	return nil, pubReceipt.Logs
}
