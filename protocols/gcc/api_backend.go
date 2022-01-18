

package gcc

import (
	"context"
	"errors"
	"math/big"

	"github.com/gcchains/chain/accounts"
	"github.com/gcchains/chain/api/rpc"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus/dpos"
	"github.com/gcchains/chain/core"
	"github.com/gcchains/chain/core/bloombits"
	"github.com/gcchains/chain/core/rawdb"
	"github.com/gcchains/chain/core/state"
	"github.com/gcchains/chain/core/vm"
	"github.com/gcchains/chain/database"
	"github.com/gcchains/chain/protocols/gcc/gasprice"
	"github.com/gcchains/chain/protocols/gcc/syncer"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/event"
)

var (
	errNilBlock             = errors.New("nil block")
	errInvalidProposersList = errors.New("invalid proposers list")
)

// APIBackend implements gccapi.Backend for full nodes
type APIBackend struct {
	gcc *gcchainService
	gpo *gasprice.Oracle
}

// ChainConfig returns the active chain configuration.
func (b *APIBackend) ChainConfig() *configs.ChainConfig {
	return b.gcc.chainConfig
}

func (b *APIBackend) CurrentBlock() *types.Block {
	return b.gcc.blockchain.CurrentBlock()
}

func (b *APIBackend) SetHead(number uint64) {
	b.gcc.blockchain.SetHead(number)
}

func (b *APIBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.gcc.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.gcc.blockchain.CurrentBlock().Header(), nil
	}
	return b.gcc.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *APIBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.gcc.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.gcc.blockchain.CurrentBlock(), nil
	}
	return b.gcc.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *APIBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber, isPrivate bool) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.gcc.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	var stateDb *state.StateDB
	if isPrivate {
		stateDb, err = b.gcc.BlockChain().StatePrivAt(header.StateRoot)
	} else {
		stateDb, err = b.gcc.BlockChain().StateAt(header.StateRoot)
	}
	return stateDb, header, err
}

func (b *APIBackend) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.gcc.blockchain.GetBlockByHash(hash), nil
}

func (b *APIBackend) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	if number := rawdb.ReadHeaderNumber(b.gcc.chainDb, hash); number != nil {
		return rawdb.ReadReceipts(b.gcc.chainDb, hash, *number), nil
	}
	return nil, nil
}

func (b *APIBackend) GetPrivateReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := core.ReadPrivateReceipt(txHash, b.ChainDb())
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func (b *APIBackend) GetLogs(ctx context.Context, hash common.Hash) ([][]*types.Log, error) {
	number := rawdb.ReadHeaderNumber(b.gcc.chainDb, hash)
	if number == nil {
		return nil, nil
	}
	receipts := rawdb.ReadReceipts(b.gcc.chainDb, hash, *number)
	if receipts == nil {
		return nil, nil
	}
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

func (b *APIBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.gcc.BlockChain(), nil)
	return vm.NewEVM(context, state, b.gcc.chainConfig, vmCfg), vmError, nil
}

func (b *APIBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.gcc.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *APIBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.gcc.BlockChain().SubscribeChainEvent(ch)
}

func (b *APIBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.gcc.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *APIBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.gcc.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *APIBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.gcc.BlockChain().SubscribeLogsEvent(ch)
}

func (b *APIBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.gcc.txPool.AddLocal(signedTx)
}

func (b *APIBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.gcc.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *APIBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.gcc.txPool.Get(hash)
}

func (b *APIBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.gcc.txPool.State().GetNonce(addr), nil
}

func (b *APIBackend) Stats() (pending int, queued int) {
	return b.gcc.txPool.Stats()
}

func (b *APIBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.gcc.TxPool().Content()
}

func (b *APIBackend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return b.gcc.TxPool().SubscribeNewTxsEvent(ch)
}

func (b *APIBackend) Downloader() syncer.Syncer {
	return b.gcc.Downloader()
}

func (b *APIBackend) ProtocolVersion() int {
	return b.gcc.CpcVersion()
}

func (b *APIBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *APIBackend) ChainDb() database.Database {
	return b.gcc.ChainDb()
}

func (b *APIBackend) EventMux() *event.TypeMux {
	return b.gcc.EventMux()
}

func (b *APIBackend) AccountManager() *accounts.Manager {
	return b.gcc.AccountManager()
}

func (b *APIBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.gcc.bloomIndexer.Sections()
	return configs.BloomBitsBlocks, sections
}

func (b *APIBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.gcc.bloomRequests)
	}
}

// RemoteDB returns remote database instance.
func (b *APIBackend) RemoteDB() database.RemoteDatabase {
	return b.gcc.RemoteDB()
}

// RNode returns current RNode information
func (b *APIBackend) RNode() ([]common.Address, uint64) {
	block := b.gcc.blockchain.CurrentBlock()
	bn := block.Number()
	api := b.gcc.engine.(*dpos.Dpos).APIs(b.gcc.blockchain)
	rNodes, _ := api[0].Service.(*dpos.API).GetRNodes()
	return rNodes, bn.Uint64()
}

// CurrentProposerIndex return current proposer index, (0,1,...,11)
func (b *APIBackend) CurrentProposerIndex() uint64 {
	block := b.gcc.blockchain.CurrentBlock()
	bn := block.Number()
	vl, tl := b.ViewLen(), b.TermLen()
	// be cautious vl*tl does not overflow
	proposerIndex := ((bn.Uint64() - 1) % (vl * tl)) % tl
	return proposerIndex
}

// CurrentView return current view, (0,1,2)
func (b *APIBackend) CurrentView() uint64 {
	block := b.gcc.blockchain.CurrentBlock()
	bn := block.Number()
	vl, tl := b.ViewLen(), b.TermLen()
	span := ((bn.Uint64() - 1) % (vl * tl)) / tl
	return span
}

// CurrentTerm return current term
func (b *APIBackend) CurrentTerm() uint64 {
	block := b.gcc.blockchain.CurrentBlock()
	bn := block.Number()
	vl, tl := b.ViewLen(), b.TermLen()
	term := (bn.Uint64() - 1) / (vl * tl)
	return term
}

// ViewLen return current ViewLen
func (b *APIBackend) ViewLen() uint64 {
	return b.gcc.chainConfig.Dpos.ViewLen
}

// TermLen return current TermLen
func (b *APIBackend) TermLen() uint64 {
	return b.gcc.chainConfig.Dpos.TermLen
}

// CommitteMember return current committe
func (b *APIBackend) CommitteMember() []common.Address {
	block := b.gcc.blockchain.CurrentBlock()
	return block.Header().Dpos.Proposers
}

func (b *APIBackend) CalcRptInfo(address common.Address, addresses []common.Address, blockNum uint64) int64 {
	return b.gcc.engine.(*dpos.Dpos).GetCalcRptInfo(address, addresses, blockNum)
}

func (b *APIBackend) BlockReward(blockNum rpc.BlockNumber) *big.Int {
	return b.gcc.engine.(*dpos.Dpos).GetBlockReward(uint64(blockNum))
}

func (b *APIBackend) ProposerOf(blockNum rpc.BlockNumber) (common.Address, error) {
	proposers, err := b.Proposers(blockNum)
	if err != nil {
		return common.Address{}, err
	}

	vl, tl := b.ViewLen(), b.TermLen()
	// be cautious vl*tl does not overflow
	view := ((uint64(blockNum) - 1) % (vl * tl)) % tl
	if len(proposers) > int(view) {
		return proposers[int(view)], nil
	}

	return common.Address{}, errInvalidProposersList
}

// Proposers returns block Proposers information
func (b *APIBackend) Proposers(blockNum rpc.BlockNumber) ([]common.Address, error) {
	block := b.gcc.BlockChain().GetBlockByNumber(uint64(blockNum))
	if block == nil {
		return []common.Address{}, errNilBlock
	}
	proposers := block.Dpos().Proposers
	return proposers, nil
}

// Validators returns current block Validators information
func (b *APIBackend) Validators(blockNr rpc.BlockNumber) ([]common.Address, error) {
	api := b.gcc.engine.(*dpos.Dpos).APIs(b.gcc.blockchain)
	return api[0].Service.(*dpos.API).GetValidators(blockNr)
}

func (b *APIBackend) SupportPrivateTx(ctx context.Context) (bool, error) {
	return types.SupportTxType(types.PrivateTx), nil
}
