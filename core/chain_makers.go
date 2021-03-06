

package core

import (
	"fmt"
	"math/big"

	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/core/state"
	"github.com/gcchains/chain/core/vm"
	"github.com/gcchains/chain/database"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

// BlockGen creates blocks for testing.
// See GenerateChain for a detailed explanation.
type BlockGen struct {
	i           int
	parent      *types.Block
	chain       []*types.Block
	chainReader consensus.ChainReader
	header      *types.Header
	pubStateDB  *state.StateDB
	privStateDB *state.StateDB

	gasPool  *GasPool
	txs      []*types.Transaction
	receipts []*types.Receipt

	config *configs.ChainConfig
	engine consensus.Engine
}

// SetCoinbase sets the coinbase of the generated block.
// It can be called at most once.
func (b *BlockGen) SetCoinbase(addr common.Address) {
	if b.gasPool != nil {
		if len(b.txs) > 0 {
			panic("coinbase must be set before adding transactions")
		}
		panic("coinbase can only be set once")
	}
	b.header.Coinbase = addr
	b.gasPool = new(GasPool).AddGas(b.header.GasLimit)
}

// SetExtra sets the extra data field of the generated block.
func (b *BlockGen) SetExtra(data []byte) {
	b.header.Extra = data
}

// AddTx adds a transaction to the generated block. If no coinbase has
// been set, the block's coinbase is set to the zero address.
//
// AddTx panics if the transaction cannot be executed. In addition to
// the protocol-imposed limitations (gas limit, etc.), there are some
// further limitations on the content of transactions that can be
// added. Notably, contract code relying on the BLOCKHASH instruction
// will panic during execution.
func (b *BlockGen) AddTx(tx *types.Transaction) {
	b.AddTxWithChain(nil, tx)
}

// AddTxWithChain adds a transaction to the generated block. If no coinbase has
// been set, the block's coinbase is set to the zero address.
//
// AddTxWithChain panics if the transaction cannot be executed. In addition to
// the protocol-imposed limitations (gas limit, etc.), there are some
// further limitations on the content of transactions that can be
// added. If contract code relies on the BLOCKHASH instruction,
// the block in chain will be returned.
func (b *BlockGen) AddTxWithChain(bc *BlockChain, tx *types.Transaction) {
	if b.gasPool == nil {
		b.SetCoinbase(common.Address{})
	}
	b.pubStateDB.Prepare(tx.Hash(), common.Hash{}, len(b.txs))
	b.privStateDB.Prepare(tx.Hash(), common.Hash{}, len(b.txs))

	var remoteDB database.RemoteDatabase
	if bc != nil {
		remoteDB = bc.remoteDB
	}
	receipt, _, _, err := ApplyTransaction(b.config, bc, &b.header.Coinbase, b.gasPool, b.pubStateDB, b.privStateDB, remoteDB,
		b.header, tx, &b.header.GasUsed, vm.Config{}, nil) // Account manager is not required in thes scenario.
	if err != nil {
		panic(err)
	}
	b.txs = append(b.txs, tx)
	b.receipts = append(b.receipts, receipt)
}

// Number returns the block number of the block being generated.
func (b *BlockGen) Number() *big.Int {
	return new(big.Int).Set(b.header.Number)
}

// AddUncheckedReceipt forcefully adds a receipts to the block without a
// backing transaction.
//
// AddUncheckedReceipt will cause consensus failures when used during real
// chain processing. This is best used in conjunction with raw block insertion.
func (b *BlockGen) AddUncheckedReceipt(receipt *types.Receipt) {
	b.receipts = append(b.receipts, receipt)
}

// TxNonce returns the next valid transaction nonce for the
// account at addr. It panics if the account does not exist.
func (b *BlockGen) TxNonce(addr common.Address) uint64 {
	if !b.pubStateDB.Exist(addr) {
		panic("account does not exist")
	}
	return b.pubStateDB.GetNonce(addr)
}

// PrevBlock returns a previously generated block by number. It panics if
// num is greater or equal to the number of the block being generated.
// For index -1, PrevBlock returns the parent block given to GenerateChain.
func (b *BlockGen) PrevBlock(index int) *types.Block {
	if index >= b.i {
		panic("block index out of range")
	}
	if index == -1 {
		return b.parent
	}
	return b.chain[index]
}

// OffsetTime modifies the time instance of a block, implicitly changing its
// associated difficulty. It's useful to test scenarios where forking is not
// tied to chain length directly.
func (b *BlockGen) OffsetTime(milliseconds int64) {
	b.header.Time.Add(b.header.Time, new(big.Int).SetInt64(milliseconds))
	if b.header.Time.Cmp(b.parent.Header().Time) <= 0 {
		panic("block time out of range")
	}
}

// GenerateChain creates a chain of n blocks. The first block's
// parent will be the provided parent. db is used to store
// intermediate states and should contain the parent's state trie.
//
// The generator function is called with a new block generator for
// every block. Any transactions and uncles added to the generator
// become part of the block. If gen is nil, the blocks will be empty
// and their coinbase will be the zero address.
//
// Blocks created by GenerateChain do not contain valid proof of work
// values. Inserting them into BlockChain requires use of FakePow or
// a similar non-validating proof of work implementation.
func GenerateChain(config *configs.ChainConfig, parent *types.Block, engine consensus.Engine, db database.Database, remoteDB database.RemoteDatabase, n int, gen func(int, *BlockGen)) ([]*types.Block, []types.Receipts) {
	if config == nil {
		config = configs.TestChainConfig
	}
	blocks, receipts := make(types.Blocks, n), make([]types.Receipts, n)
	blockchain, _ := NewBlockChain(db, nil, config, engine, vm.Config{}, remoteDB, nil)
	defer blockchain.Stop()

	genblock := func(i int, parent *types.Block, pubStatedb *state.StateDB, privStateDB *state.StateDB) (*types.Block, types.Receipts) {
		// TODO(karalabe): This is needed for clique, which depends on multiple blocks.
		// It's nonetheless ugly to spin up a blockchain here. Get rid of this somehow.

		b := &BlockGen{i: i, parent: parent, chain: blocks, chainReader: blockchain, pubStateDB: pubStatedb, privStateDB: privStateDB, config: config, engine: engine}

		b.header = makeHeader(b.chainReader, parent, pubStatedb, b.engine)

		// Execute any user modifications to the block and finalize it
		if gen != nil {
			gen(i, b)
		}

		if b.engine != nil {
			block, _ := b.engine.Finalize(b.chainReader, b.header, pubStatedb, b.txs, []*types.Header{}, b.receipts)
			// Write state changes to db
			root, err := pubStatedb.Commit(true)
			if err != nil {
				panic(fmt.Sprintf("state write error: %v", err))
			}
			if err := pubStatedb.Database().TrieDB().Commit(root, false); err != nil {
				panic(fmt.Sprintf("trie write error: %v", err))
			}

			// WARN: TODO: to avoid impeach block validation failure, set same gasLimit as its parent
			block.RefHeader().GasLimit = parent.GasLimit()
			block.RefHeader().Extra = make([]byte, 65)

			return block, b.receipts
		}
		return nil, nil
	}
	// end of genblock() function

	for i := 0; i < n; i++ {
		pubStatedb, err := state.New(parent.StateRoot(), state.NewDatabase(db))
		if err != nil {
			panic(err)
		}
		privStateDB, err := state.New(GetPrivateStateRoot(db, parent.StateRoot()), state.NewDatabase(db))
		if err != nil {
			panic(err)
		}
		block, receipt := genblock(i, parent, pubStatedb, privStateDB)

		blocks[i] = block
		receipts[i] = receipt
		parent = block
	}
	return blocks, receipts
}

func makeHeader(chain consensus.ChainReader, parent *types.Block, state *state.StateDB, engine consensus.Engine) *types.Header {
	var time *big.Int
	if parent.Time() == nil {
		time = big.NewInt(10)
	} else {
		time = new(big.Int).Add(parent.Time(), big.NewInt(10)) // block time is fixed at 10 seconds
	}

	header := &types.Header{}

	header.Number = new(big.Int).Add(parent.Number(), common.Big1)
	header.ParentHash = parent.Hash()

	_ = engine.PrepareBlock(chain, header)

	header.StateRoot = state.IntermediateRoot(true)
	header.Coinbase = parent.Coinbase()
	header.GasLimit = CalcGasLimit(parent)

	header.Time = time

	return header
}

// makeHeaderChain creates a deterministic chain of headers rooted at parent.
func makeHeaderChain(parent *types.Header, n int, engine consensus.Engine, db database.Database, seed int) []*types.Header {
	blocks := makeBlockChain(types.NewBlockWithHeader(parent), n, engine, db, seed)
	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	return headers
}

// makeBlockChain creates a deterministic chain of blocks rooted at parent.
func makeBlockChain(parent *types.Block, n int, engine consensus.Engine, db database.Database, seed int) []*types.Block {
	blocks, _ := GenerateChain(configs.TestChainConfig, parent, engine, db, nil, n, func(i int, b *BlockGen) {
		b.SetCoinbase(common.Address{0: byte(seed), 19: byte(i)})
	})
	return blocks
}
