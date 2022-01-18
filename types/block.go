

// Package types contains data types related to Ethereum consensus.
package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"sort"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/gcchains/chain/commons/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	DposSigLength = 65
)

var (
	EmptyRootHash = DeriveSha(Transactions{})
)

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [8]byte

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockNonce", input, n[:])
}

//go:generate gencodec -type Header -formats json,toml -field-override headerMarshaling -out gen_header_json.go
const (
	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
)

// Header represents a block header in the Ethereum blockchain.
type Header struct {
	ParentHash   common.Hash    `json:"parentHash"       gencodec:"required"`
	Coinbase     common.Address `json:"miner"            gencodec:"required"`
	StateRoot    common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxsRoot      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptsRoot common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	LogsBloom    Bloom          `json:"logsBloom"        gencodec:"required"`
	Number       *big.Int       `json:"number"           gencodec:"required"`
	GasLimit     uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed      uint64         `json:"gasUsed"          gencodec:"required"`
	Time         *big.Int       `json:"timestamp"        gencodec:"required"` // this is a value with accuracy of Millisecond
	Extra        []byte         `json:"extraData"        gencodec:"required"`
	Dpos         DposSnap       `json:"dpos"             gencodec:"required"`
}

type DposSignature [DposSigLength]byte

// MarshalText encodes n as a hex string with 0x prefix.
func (n DposSignature) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *DposSignature) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("DposSignature", input, n[:])
}

func HexToDposSig(s string) DposSignature {
	var a DposSignature
	copy(a[:], common.FromHex(s))
	return a
}

func (d *DposSignature) IsEmpty() bool {
	return bytes.Equal(d[:], bytes.Repeat([]byte{0x00}, DposSigLength))
}

type DposSnap struct {
	Seal       DposSignature    `json:"seal"`       // the signature of the block's proposer
	Sigs       []DposSignature  `json:"sigs"`       // the signatures of validators to endorse the block
	Proposers  []common.Address `json:"proposers"`  // current proposers committee
	Validators []common.Address `json:"validators"` // updated validator committee in next epoch if it is not nil. Keep the same to current if it is nil.
}

func (d *DposSnap) SigsFormatText() string {
	items := make([]string, len(d.Sigs))
	for idx, sig := range d.Sigs {
		items[idx] = fmt.Sprintf("#%d %s", idx, common.Bytes2Hex(sig[:]))
	}
	return strings.Join(items, ",")
}

func (d *DposSnap) ProposersFormatText() string {
	items := make([]string, len(d.Proposers))
	for idx, p := range d.Proposers {
		items[idx] = fmt.Sprintf("[#%d %s]", idx, p.Hex())
	}
	return fmt.Sprintf("[%s]", strings.Join(items, ","))
}

func (d *DposSnap) ValidatorsFormatText() string {
	items := make([]string, len(d.Validators))
	for idx, v := range d.Validators {
		items[idx] = fmt.Sprintf("[#%d %s]", idx, v.Hex())
	}
	return fmt.Sprintf("[%s]", strings.Join(items, ","))
}

// field type overrides for gencodec
type headerMarshaling struct {
	Number   *hexutil.Big
	GasLimit hexutil.Uint64
	GasUsed  hexutil.Uint64
	Time     *hexutil.Big
	Extra    hexutil.Bytes
	Hash     common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
	Dpos     DposSnap
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.Hash {
	return sigHash(h)
}

// sigHash returns hash of header
func sigHash(header *Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()
	err := rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.Coinbase,
		header.StateRoot,
		header.TxsRoot,
		header.ReceiptsRoot,
		header.LogsBloom,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Dpos.Proposers,
		header.Dpos.Validators,
		header.Extra,
		common.Hash{},
		BlockNonce{},
	})
	if err != nil {
		log.Error("invalid hash encoding", "error", err)
		return common.Hash{}
	}
	hasher.Sum(hash[:0])
	return hash
}

// HashNoNonce returns the hash which is used as input for the proof-of-work search.
func (h *Header) HashNoNonce() common.Hash {
	return rlpHash([]interface{}{
		h.ParentHash,
		h.Coinbase,
		h.StateRoot,
		h.TxsRoot,
		h.ReceiptsRoot,
		h.LogsBloom,
		h.Number,
		h.GasLimit,
		h.GasUsed,
		h.Time,
		h.Dpos.Proposers,
		h.Dpos.Validators,
		h.Extra,
	})
}

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (h *Header) Size() common.StorageSize {
	dposSize := common.StorageSize(len(h.Dpos.Proposers))*common.StorageSize(unsafe.Sizeof(common.Address{})) +
		common.StorageSize(len(h.Dpos.Sigs))*common.StorageSize(unsafe.Sizeof(DposSignature{})) +
		common.StorageSize(len(h.Dpos.Validators))*common.StorageSize(unsafe.Sizeof(common.Address{})) +
		common.StorageSize(unsafe.Sizeof(h.Dpos.Seal))

	return common.StorageSize(unsafe.Sizeof(*h)) + common.StorageSize(len(h.Extra)+(h.Number.BitLen()+h.Time.BitLen())/8) + dposSize
}

func (h *Header) Timestamp() time.Time {
	return time.Unix(0, h.Time.Int64()*int64(time.Millisecond)/int64(time.Nanosecond))
}

func (h *Header) Impeachment() bool {
	return h.Coinbase == common.Address{}
}

func (h *Header) SetTimestamp(t time.Time) {
	timestamp := int64(t.UnixNano()) * int64(time.Nanosecond) / int64(time.Millisecond)
	if h.Time == nil {
		h.Time = new(big.Int)
	}
	h.Time.SetInt64(timestamp)
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions and uncles) together.
type Body struct {
	Transactions []*Transaction
}

// Block represents an entire block in the Ethereum blockchain.
type Block struct {
	header       *Header
	transactions Transactions

	// caches
	hash atomic.Value
	size atomic.Value

	// Td is used by package core to store the total difficulty
	// of the chain up to and including the block.
	td *big.Int

	// These fields are used to track inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
}

// DeprecatedTd is an old relic for extracting the TD of a block. It is in the
// code solely to facilitate upgrading the database from the old format to the
// new, after which it should be deleted. Do not use!
func (b *Block) DeprecatedTd() *big.Int {
	return b.td
}

// [deprecated by eth/63]
// StorageBlock defines the RLP encoding of a Block stored in the
// state database. The StorageBlock encoding contains fields that
// would otherwise need to be recomputed.
type StorageBlock Block

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Header *Header
	Txs    []*Transaction
}

// [deprecated by eth/63]
// "storage" block encoding. used for database.
type storageblock struct {
	Header *Header
	Txs    []*Transaction
	TD     *big.Int
}

// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the
// block.
//
// The values of TxsRoot, UncleHash, ReceiptsRoot and LogsBloom in header
// are ignored and set to values derived from the given txs, uncles
// and receipts.
func NewBlock(header *Header, txs []*Transaction, receipts []*Receipt) *Block {
	b := &Block{header: CopyHeader(header), td: new(big.Int)}

	// TODO: panic if len(txs) != len(receipts)
	if len(txs) == 0 {
		b.header.TxsRoot = EmptyRootHash
	} else {
		b.header.TxsRoot = DeriveSha(Transactions(txs))
		b.transactions = make(Transactions, len(txs))
		copy(b.transactions, txs)
	}

	if len(receipts) == 0 {
		b.header.ReceiptsRoot = EmptyRootHash
	} else {
		b.header.ReceiptsRoot = DeriveSha(Receipts(receipts))
		b.header.LogsBloom = CreateBloom(receipts)
	}

	return b
}

// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: CopyHeader(header)}
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	cpy := *h
	if cpy.Time = new(big.Int); h.Time != nil {
		cpy.Time.Set(h.Time)
	}
	if cpy.Number = new(big.Int); h.Number != nil {
		cpy.Number.Set(h.Number)
	}
	if len(h.Extra) > 0 {
		cpy.Extra = make([]byte, len(h.Extra))
		copy(cpy.Extra, h.Extra)
	}
	cpy.Dpos = *CopyDposSnap(&h.Dpos)
	return &cpy
}

func CopyDposSnap(d *DposSnap) *DposSnap {
	// copy DposSnap
	cpy := new(DposSnap)
	// copy DposSnap.Proposers
	cpy.Proposers = d.CopyProposers()
	// copy DposSnap.Sigs
	cpy.Sigs = make([]DposSignature, len(d.Sigs))
	for i := 0; i < len(d.Sigs); i++ {
		copy(cpy.Sigs[i][:], d.Sigs[i][:])
	}
	// copy DposSnap.Seal
	copy(cpy.Seal[:], d.Seal[:])
	// copy DposSnap.Validators
	cpy.Validators = d.CopyValidators()
	return cpy
}

func (d *DposSnap) CopyProposers() []common.Address {
	proposers := make([]common.Address, len(d.Proposers))
	for i := 0; i < len(d.Proposers); i++ {
		copy(proposers[i][:], d.Proposers[i][:])
	}
	return proposers
}

func (d *DposSnap) CopyValidators() []common.Address {
	validators := make([]common.Address, len(d.Validators))
	for i := 0; i < len(d.Validators); i++ {
		copy(validators[i][:], d.Validators[i][:])
	}
	return validators
}

// DecodeRLP decodes the Ethereum
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb extblock
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.header, b.transactions = eb.Header, eb.Txs
	b.size.Store(common.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extblock{
		Header: b.header,
		Txs:    b.transactions,
	})
}

// [deprecated by eth/63]
func (b *StorageBlock) DecodeRLP(s *rlp.Stream) error {
	var sb storageblock
	if err := s.Decode(&sb); err != nil {
		return err
	}
	b.header, b.transactions, b.td = sb.Header, sb.Txs, sb.TD
	return nil
}

// TODO: copies

func (b *Block) Transactions() Transactions { return b.transactions }

func (b *Block) Transaction(hash common.Hash) *Transaction {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}

func (b *Block) Number() *big.Int { return new(big.Int).Set(b.header.Number) }
func (b *Block) GasLimit() uint64 { return b.header.GasLimit }
func (b *Block) GasUsed() uint64  { return b.header.GasUsed }

func (b *Block) Time() *big.Int           { return new(big.Int).Set(b.header.Time) }
func (b *Block) Timestamp() time.Time     { return b.Header().Timestamp() }
func (b *Block) SetTimestamp(t time.Time) { b.RefHeader().SetTimestamp(t) }
func (b *Block) Impeachment() bool        { return b.Header().Impeachment() }

func (b *Block) NumberU64() uint64         { return b.header.Number.Uint64() }
func (b *Block) LogsBloom() Bloom          { return b.header.LogsBloom }
func (b *Block) Coinbase() common.Address  { return b.header.Coinbase }
func (b *Block) StateRoot() common.Hash    { return b.header.StateRoot }
func (b *Block) ParentHash() common.Hash   { return b.header.ParentHash }
func (b *Block) TxsRoot() common.Hash      { return b.header.TxsRoot }
func (b *Block) ReceiptsRoot() common.Hash { return b.header.ReceiptsRoot }
func (b *Block) Extra() []byte             { return common.CopyBytes(b.header.Extra) }
func (b *Block) Dpos() DposSnap            { return b.header.Dpos }

func (b *Block) RefHeader() *Header { return b.header } // TODO: fix it.
func (b *Block) Header() *Header    { return CopyHeader(b.header) }

// Body returns the non-header content of the block.
func (b *Block) Body() *Body { return &Body{b.transactions} }

func (b *Block) HashNoNonce() common.Hash {
	return b.header.HashNoNonce()
}

// Size returns the true RLP encoded storage size of the block, either by encoding
// and returning it, or returning a previsouly cached value.
func (b *Block) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block) WithSeal(header *Header) *Block {
	cpy := *header

	return &Block{
		header:       &cpy,
		transactions: b.transactions,
	}
}

// WithBody returns a new block with the given transaction and uncle contents.
func (b *Block) WithBody(transactions []*Transaction) *Block {
	block := &Block{
		header:       CopyHeader(b.header),
		transactions: make([]*Transaction, len(transactions)),
	}
	copy(block.transactions, transactions)
	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

type Blocks []*Block

type BlockBy func(b1, b2 *Block) bool

func (self BlockBy) Sort(blocks Blocks) {
	bs := blockSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(bs)
}

type blockSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}

func (self blockSorter) Len() int { return len(self.blocks) }
func (self blockSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}
func (self blockSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }

func Number(b1, b2 *Block) bool { return b1.header.Number.Cmp(b2.header.Number) < 0 }
