

package gcc

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
	"sync"

	"github.com/gcchains/chain/api/rpc"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus/dpos"
	"github.com/gcchains/chain/consensus/dpos/backend"
	"github.com/gcchains/chain/core"
	"github.com/gcchains/chain/core/rawdb"
	"github.com/gcchains/chain/core/state"
	"github.com/gcchains/chain/internal/gccapi"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

// PublicgcchainAPI provides an API to access gcchain full node-related
// information.
type PublicgcchainAPI struct {
	c *gcchainService
}

// NewPublicgcchainAPI creates a new gcchain protocol API for full nodes.
func NewPublicgcchainAPI(e *gcchainService) *PublicgcchainAPI {
	return &PublicgcchainAPI{e}
}

// Coinbase is the address that mining rewards will be send to
func (api *PublicgcchainAPI) Gccbase() (common.Address, error) {
	return api.c.Coinbase()
}

// Coinbase is the address that mining rewards will be send to (alias for Coinbase)
func (api *PublicgcchainAPI) Coinbase() (common.Address, error) {
	return api.Gccbase()
}

// PublicMinerAPI provides an API to control the miner.
// It offers only methods that operate on data that pose no security risk when it is publicly accessible.
type PublicMinerAPI struct {
	c *gcchainService
}

// NewPublicMinerAPI create a new PublicMinerAPI instance.
func NewPublicMinerAPI(e *gcchainService) *PublicMinerAPI {
	return &PublicMinerAPI{e}
}

// Mining returns an indication if this node is currently mining.
func (api *PublicMinerAPI) Mining() bool {
	return api.c.IsMining()
}

// PrivateMinerAPI provides private RPC methods to control the miner.
// These methods can be abused by external users and must be considered insecure for use by untrusted users.
type PrivateMinerAPI struct {
	c    *gcchainService
	lock sync.Mutex // Protected mining starting and stopping
}

// NewPrivateMinerAPI create a new RPC service which controls the miner of this node.
func NewPrivateMinerAPI(e *gcchainService) *PrivateMinerAPI {
	return &PrivateMinerAPI{c: e}
}

// Start the miner with the given number of threads. If threads is nil the number
// of workers started is equal to the number of logical CPUs that are usable by
// this process. If mining is already running, this method adjust the number of
// threads allowed to use.
func (api *PrivateMinerAPI) Start(threads *int) error {
	// make sure the api executes in sequence(no parallel)
	api.lock.Lock()
	defer api.lock.Unlock()

	if api.c.IsMining() {
		return nil
	}

	// Set the number of threads if the seal engine supports it
	if threads == nil {
		threads = new(int)
	} else if *threads == 0 {
		*threads = -1 // Disable the miner from within
	}
	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := api.c.engine.(threaded); ok {
		log.Info("Updated mining threads", "threads", *threads)
		th.SetThreads(*threads)
	}
	// Start the miner and return
	// TODO: @ac pass the concrete client argument in the function call
	return api.c.StartMining(true)
}

// Stop the miner
func (api *PrivateMinerAPI) Stop() bool {
	// make sure the api executes in sequence(no parallel)
	api.lock.Lock()
	defer api.lock.Unlock()

	if !api.c.IsMining() {
		return true
	}

	type threaded interface {
		SetThreads(threads int)
	}
	if th, ok := api.c.engine.(threaded); ok {
		th.SetThreads(-1)
	}
	api.c.StopMining()
	return true
}

// SetExtra sets the extra data string that is included when this miner mines a block.
func (api *PrivateMinerAPI) SetExtra(extra string) (bool, error) {
	// make sure the api executes in sequence(no parallel)
	api.lock.Lock()
	defer api.lock.Unlock()

	if err := api.c.Miner().SetExtra([]byte(extra)); err != nil {
		return false, err
	}
	return true, nil
}

// SetGasPrice sets the minimum accepted gas price for the miner.
func (api *PrivateMinerAPI) SetGasPrice(gasPrice hexutil.Big) bool {
	// make sure the api executes in sequence(no parallel)
	api.lock.Lock()
	defer api.lock.Unlock()

	api.c.lock.Lock()
	api.c.gasPrice = (*big.Int)(&gasPrice)
	api.c.lock.Unlock()

	api.c.txPool.SetGasPrice((*big.Int)(&gasPrice))
	return true
}

func (api *PrivateMinerAPI) SetCoinbase(coinbase common.Address) bool {
	// make sure the api executes in sequence(no parallel)
	api.lock.Lock()
	defer api.lock.Unlock()

	api.c.SetCoinbase(coinbase)
	return true
}

// PrivateAdminAPI is the collection of gcchain full node-related APIs
// exposed over the private admin endpoint.
type PrivateAdminAPI struct {
	gcc *gcchainService
}

// NewPrivateAdminAPI creates a new API definition for the full node private
// admin methods of the gcchain service.
func NewPrivateAdminAPI(eth *gcchainService) *PrivateAdminAPI {
	return &PrivateAdminAPI{gcc: eth}
}

// ExportChain exports the current blockchain into a local file.
func (api *PrivateAdminAPI) ExportChain(file string) (bool, error) {
	// Make sure we can create the file to export into
	out, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return false, err
	}
	defer out.Close()

	var writer io.Writer = out
	if strings.HasSuffix(file, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}

	// Export the blockchain
	if err := api.gcc.BlockChain().Export(writer); err != nil {
		return false, err
	}
	return true, nil
}

func hasAllBlocks(chain *core.BlockChain, bs []*types.Block) bool {
	for _, b := range bs {
		if !chain.HasBlock(b.Hash(), b.NumberU64()) {
			return false
		}
	}

	return true
}

// ImportChain imports a blockchain from a local file.
func (api *PrivateAdminAPI) ImportChain(file string) (bool, error) {
	// Make sure the can access the file to import
	in, err := os.Open(file)
	if err != nil {
		return false, err
	}
	defer in.Close()

	var reader io.Reader = in
	if strings.HasSuffix(file, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return false, err
		}
	}

	// Run actual the import in pre-configured batches
	stream := rlp.NewStream(reader, 0)

	blocks, index := make([]*types.Block, 0, 2500), 0
	for batch := 0; ; batch++ {
		// Load a batch of blocks from the input file
		for len(blocks) < cap(blocks) {
			block := new(types.Block)
			if err := stream.Decode(block); err == io.EOF {
				break
			} else if err != nil {
				return false, fmt.Errorf("block %d: failed to parse: %v", index, err)
			}
			blocks = append(blocks, block)
			index++
		}
		if len(blocks) == 0 {
			break
		}

		if hasAllBlocks(api.gcc.BlockChain(), blocks) {
			blocks = blocks[:0]
			continue
		}
		// Import the batch and reset the buffer
		if _, err := api.gcc.BlockChain().InsertChain(blocks); err != nil {
			return false, fmt.Errorf("batch %d: failed to insert: %v", batch, err)
		}
		blocks = blocks[:0]
	}
	return true, nil
}

// PublicDebugAPI is the collection of gcchain full node APIs exposed
// over the public debugging endpoint.
type PublicDebugAPI struct {
	gcc *gcchainService
}

// NewPublicDebugAPI creates a new API definition for the full node-
// related public debug methods of the gcchain service.
func NewPublicDebugAPI(eth *gcchainService) *PublicDebugAPI {
	return &PublicDebugAPI{gcc: eth}
}

// DumpBlock retrieves the entire state of the database at a given block.
func (api *PublicDebugAPI) DumpBlock(blockNr rpc.BlockNumber) (state.Dump, error) {
	if blockNr == rpc.PendingBlockNumber {
		// If we're dumping the pending state, we need to request
		// both the pending block as well as the pending state from
		// the miner and operate on those
		_, stateDb := api.gcc.miner.Pending()
		return stateDb.RawDump(), nil
	}
	var block *types.Block
	if blockNr == rpc.LatestBlockNumber {
		block = api.gcc.blockchain.CurrentBlock()
	} else {
		block = api.gcc.blockchain.GetBlockByNumber(uint64(blockNr))
	}
	if block == nil {
		return state.Dump{}, fmt.Errorf("block #%d not found", blockNr)
	}
	stateDb, err := api.gcc.BlockChain().StateAt(block.StateRoot())
	if err != nil {
		return state.Dump{}, err
	}
	return stateDb.RawDump(), nil
}

// DposPeers retrieves all the information we know about each individual peer
// at the dpos protocol level
func (api *PublicDebugAPI) DposPeers() ([]*backend.PeerInfo, error) {
	return api.gcc.Engine().(*dpos.Dpos).PeerInfos()
}

// PrivateDebugAPI is the collection of gcchain full node APIs exposed over
// the private debugging endpoint.
type PrivateDebugAPI struct {
	config *configs.ChainConfig
	gcc    *gcchainService
}

// NewPrivateDebugAPI creates a new API definition for the full node-related
// private debug methods of the gcchain service.
func NewPrivateDebugAPI(config *configs.ChainConfig, service *gcchainService) *PrivateDebugAPI {
	return &PrivateDebugAPI{config: config, gcc: service}
}

// Preimage is a debug API function that returns the preimage for a sha3 hash, if known.
func (api *PrivateDebugAPI) Preimage(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	if preimage := rawdb.ReadPreimage(api.gcc.ChainDb(), hash); preimage != nil {
		return preimage, nil
	}
	return nil, errors.New("unknown preimage")
}

// BadBlockArgs represents the entries in the list returned when bad blocks are queried.
type BadBlockArgs struct {
	Hash  common.Hash            `json:"hash"`
	Block map[string]interface{} `json:"block"`
	RLP   string                 `json:"rlp"`
}

// GetBadBlocks returns a list of the last 'bad blocks' that the client has seen on the network
// and returns them as a JSON list of block-hashes
func (api *PrivateDebugAPI) GetBadBlocks(ctx context.Context) ([]*BadBlockArgs, error) {
	blocks := api.gcc.BlockChain().BadBlocks()
	results := make([]*BadBlockArgs, len(blocks))

	var err error
	for i, block := range blocks {
		results[i] = &BadBlockArgs{
			Hash: block.Hash(),
		}
		if rlpBytes, err := rlp.EncodeToBytes(block); err != nil {
			results[i].RLP = err.Error() // Hacky, but hey, it works
		} else {
			results[i].RLP = fmt.Sprintf("0x%x", rlpBytes)
		}
		if results[i].Block, err = gccapi.RPCMarshalBlock(block, true, true); err != nil {
			results[i].Block = map[string]interface{}{"error": err.Error()}
		}
	}
	return results, nil
}

// StorageRangeResult is the result of a debug_storageRangeAt API call.
type StorageRangeResult struct {
	Storage storageMap   `json:"storage"`
	NextKey *common.Hash `json:"nextKey"` // nil if Storage includes the last key in the trie.
}

type storageMap map[common.Hash]storageEntry

type storageEntry struct {
	Key   *common.Hash `json:"key"`
	Value common.Hash  `json:"value"`
}

// StorageRangeAt returns the storage at the given block height and transaction index.
func (api *PrivateDebugAPI) StorageRangeAt(ctx context.Context, blockHash common.Hash, txIndex int, contractAddress common.Address, keyStart hexutil.Bytes, maxResult int) (StorageRangeResult, error) {
	_, _, statedb, err := api.computeTxEnv(blockHash, txIndex, 0)
	if err != nil {
		return StorageRangeResult{}, err
	}
	st := statedb.StorageTrie(contractAddress)
	if st == nil {
		return StorageRangeResult{}, fmt.Errorf("account %x doesn't exist", contractAddress)
	}
	return storageRangeAt(st, keyStart, maxResult)
}

func storageRangeAt(st state.Trie, start []byte, maxResult int) (StorageRangeResult, error) {
	it := trie.NewIterator(st.NodeIterator(start))
	result := StorageRangeResult{Storage: storageMap{}}
	for i := 0; i < maxResult && it.Next(); i++ {
		_, content, _, err := rlp.Split(it.Value)
		if err != nil {
			return StorageRangeResult{}, err
		}
		e := storageEntry{Value: common.BytesToHash(content)}
		if preimage := st.GetKey(it.Key); preimage != nil {
			preimage := common.BytesToHash(preimage)
			e.Key = &preimage
		}
		result.Storage[common.BytesToHash(it.Key)] = e
	}
	// Add the 'next key' so clients can continue downloading.
	if it.Next() {
		next := common.BytesToHash(it.Key)
		result.NextKey = &next
	}
	return result, nil
}

// GetModifiedAccountsByNumber returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *PrivateDebugAPI) GetModifiedAccountsByNumber(startNum uint64, endNum *uint64) ([]common.Address, error) {
	var startBlock, endBlock *types.Block

	startBlock = api.gcc.blockchain.GetBlockByNumber(startNum)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startNum)
	}

	if endNum == nil {
		endBlock = startBlock
		startBlock = api.gcc.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.gcc.blockchain.GetBlockByNumber(*endNum)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %d not found", *endNum)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

// GetModifiedAccountsByHash returns all accounts that have changed between the
// two blocks specified. A change is defined as a difference in nonce, balance,
// code hash, or storage hash.
//
// With one parameter, returns the list of accounts modified in the specified block.
func (api *PrivateDebugAPI) GetModifiedAccountsByHash(startHash common.Hash, endHash *common.Hash) ([]common.Address, error) {
	var startBlock, endBlock *types.Block
	startBlock = api.gcc.blockchain.GetBlockByHash(startHash)
	if startBlock == nil {
		return nil, fmt.Errorf("start block %x not found", startHash)
	}

	if endHash == nil {
		endBlock = startBlock
		startBlock = api.gcc.blockchain.GetBlockByHash(startBlock.ParentHash())
		if startBlock == nil {
			return nil, fmt.Errorf("block %x has no parent", endBlock.Number())
		}
	} else {
		endBlock = api.gcc.blockchain.GetBlockByHash(*endHash)
		if endBlock == nil {
			return nil, fmt.Errorf("end block %x not found", *endHash)
		}
	}
	return api.getModifiedAccounts(startBlock, endBlock)
}

func (api *PrivateDebugAPI) getModifiedAccounts(startBlock, endBlock *types.Block) ([]common.Address, error) {
	if startBlock.Number().Uint64() >= endBlock.Number().Uint64() {
		return nil, fmt.Errorf("start block height (%d) must be less than end block height (%d)", startBlock.Number().Uint64(), endBlock.Number().Uint64())
	}

	oldTrie, err := trie.NewSecure(startBlock.StateRoot(), trie.NewDatabase(api.gcc.chainDb), 0)
	if err != nil {
		return nil, err
	}
	newTrie, err := trie.NewSecure(endBlock.StateRoot(), trie.NewDatabase(api.gcc.chainDb), 0)
	if err != nil {
		return nil, err
	}

	diff, _ := trie.NewDifferenceIterator(oldTrie.NodeIterator([]byte{}), newTrie.NodeIterator([]byte{}))
	iter := trie.NewIterator(diff)

	var dirty []common.Address
	for iter.Next() {
		key := newTrie.GetKey(iter.Key)
		if key == nil {
			return nil, fmt.Errorf("no preimage found for hash %x", iter.Key)
		}
		dirty = append(dirty, common.BytesToAddress(key))
	}
	return dirty, nil
}
