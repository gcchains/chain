

package dpos

import (
	"bytes"
	"encoding/binary"
	"sync"
	"time"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/database"

	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
)

func nanosecondToMillisecond(t int64) int64 {
	return t * int64(time.Nanosecond) / int64(time.Millisecond)
}

func millisecondToNanosecond(t int64) int64 {
	return t * int64(time.Millisecond) / int64(time.Nanosecond)
}

// signatures represents signatures of a block signed by validators
type signatures struct {
	lock sync.RWMutex
	sigs map[common.Address][]byte
}

// getSig gets addr's sig
func (s *signatures) getSig(addr common.Address) (sig []byte, ok bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	sig, ok = s.sigs[addr]
	return sig, ok
}

// setSig sets addr's sig
func (s *signatures) setSig(addr common.Address, sig []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.sigs[addr] = sig
}

type dposUtil interface {
	sigHash(header *types.Header) (hash common.Hash)
	ecrecover(header *types.Header, sigcache *lru.ARCCache) (common.Address, []common.Address, error)
}

type defaultDposUtil struct {
	lock sync.RWMutex
}

// sigHash returns the hash which is used as input for the proof-of-authority
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func (d *defaultDposUtil) sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	contentToHash := []interface{}{
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
		types.BlockNonce{},
	}
	rlp.Encode(hasher, contentToHash)

	hasher.Sum(hash[:0])
	return hash
}

// ecrecover extracts the gcchain account address from a signed header.
// the return value is (the_proposer_address, validators_committee_addresses, error)
func (d *defaultDposUtil) ecrecover(header *types.Header, sigcache *lru.ARCCache) (common.Address, []common.Address, error) {
	d.lock.Lock()
	defer d.lock.Unlock()

	hash := header.Hash()
	var proposer common.Address

	if !bytes.Equal(header.Dpos.Seal[:], new(types.DposSignature)[:]) {
		// Retrieve leader's signature
		proposerSig := header.Dpos.Seal

		// Recover the public key and the gcchain address of leader.
		proposerPubKey, err := crypto.Ecrecover(d.sigHash(header).Bytes(), proposerSig[:])
		if err != nil {
			return common.Address{}, []common.Address{}, err
		}
		copy(proposer[:], crypto.Keccak256(proposerPubKey[1:])[12:])

		// Cache proposer signature.
		if sigs, known := sigcache.Get(hash); known {
			sigs.(*signatures).setSig(proposer, proposerSig[:])
		} else {
			sigs := &signatures{
				sigs: make(map[common.Address][]byte),
			}
			sigs.setSig(proposer, proposerSig[:])
			sigcache.Add(hash, sigs)
		}
	}

	// Recover the public key and the gcchain address of signers one by one.
	var validators []common.Address
	for i := 0; i < len(header.Dpos.Sigs); i++ {
		signerSig := header.Dpos.Sigs[i]

		noSigner := bytes.Equal(signerSig[:], make([]byte, extraSeal))
		if !noSigner {

			// Recover it!
			hashToSign, err := hashBytesWithState(d.sigHash(header).Bytes(), consensus.Commit)
			signerPubkey, err := crypto.Ecrecover(hashToSign, signerSig[:])
			if err != nil {
				continue
			}

			var validator common.Address
			copy(validator[:], crypto.Keccak256(signerPubkey[1:])[12:])

			// Cache it!
			sigs, ok := sigcache.Get(hash)
			if ok {
				sigs.(*signatures).setSig(validator, signerSig[:])

			} else {
				sigs := &signatures{
					sigs: make(map[common.Address][]byte),
				}
				sigs.setSig(validator, signerSig[:])
				sigcache.Add(hash, sigs)
			}

			// Add signer to known signers
			validators = append(validators, validator)
		}
	}
	return proposer, validators, nil
}

const (
	maxSignedBlocksRecordInCache = 1024
)

type signedBlocksRecord struct {
	cache *lru.ARCCache
	db    database.Database
	lock  sync.RWMutex
}

func newSignedBlocksRecord(db database.Database) *signedBlocksRecord {
	cache, _ := lru.NewARC(maxSignedBlocksRecordInCache)
	return &signedBlocksRecord{
		db:    db,
		cache: cache,
	}
}

func (sbr *signedBlocksRecord) ifAlreadySigned(number uint64) (common.Hash, bool) {
	sbr.lock.RLock()
	defer sbr.lock.RUnlock()

	// retrieve from cache
	h, ok := sbr.cache.Get(number)
	if ok {
		hash := h.(common.Hash)
		return hash, ok
	}

	// retrieve from db
	hb, err := sbr.db.Get(numberToBytes(number))
	if err == nil {
		hash, ok := common.BytesToHash(hb), true
		return hash, ok
	}

	return common.Hash{}, false
}

func (sbr *signedBlocksRecord) markAsSigned(number uint64, hash common.Hash) (err error) {
	sbr.lock.Lock()
	defer sbr.lock.Unlock()

	// add to cache
	sbr.cache.Add(number, hash)

	// add to db
	err = sbr.db.Put(numberToBytes(number), hash.Bytes())

	return
}

func numberToBytes(number uint64) []byte {
	numberBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberBytes, number)
	return numberBytes
}

func hashBytesWithState(hash []byte, state consensus.State) (signHashBytes []byte, err error) {
	var (
		prepreparePrefix = "Prepare"
	)

	var bytesToSign []byte
	switch state {
	case consensus.Prepare, consensus.ImpeachPrepare:
		bytesToSign = append([]byte(prepreparePrefix), hash...)
	case consensus.Commit, consensus.ImpeachCommit:
		bytesToSign = hash
	default:
		log.Warn("unknown state when signing hash with state", "state", state)
		// TODO: add new error type here
		err = nil
	}

	var signHash common.Hash
	if len(bytesToSign) > len(hash) {
		hasher := sha3.NewKeccak256()
		hasher.Write(bytesToSign)
		hasher.Sum(signHash[:0])
	} else {
		signHash = common.BytesToHash(hash)
	}

	signHashBytes = signHash.Bytes()
	return
}
