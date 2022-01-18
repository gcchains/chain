

package dpos

import (
	"math"
	"reflect"
	"time"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

type dposHelper interface {
	dposUtil

	verifyHeader(d *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header,
		verifySigs bool, verifyProposers bool) error

	snapshot(d *Dpos, chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) (*DposSnapshot, error)

	verifyBasic(d *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error
	verifySeal(d *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error
	verifySignatures(d *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error
	verifyProposers(dpos *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error

	signHeader(d *Dpos, chain consensus.ChainReader, header *types.Header, state consensus.State) error
	validateBlock(d *Dpos, chain consensus.ChainReader, block *types.Block, verifySigs bool, verifyProposers bool) error
}

type defaultDposHelper struct {
	dposUtil
}

// validateBlock checks basic fields in a block, this is called only by validators
func (dh *defaultDposHelper) validateBlock(c *Dpos, chain consensus.ChainReader, block *types.Block, verifySigs bool, verifyProposers bool) error {

	// verify the `validators` field in the header is empty
	if len(block.Header().Dpos.Validators) != 0 {
		return consensus.ErrorInvalidValidatorsList
	}

	// verify the block header according to Dpos Protocol
	if err := dh.verifyHeader(c, chain, block.Header(), nil, block.RefHeader(), verifySigs, verifyProposers); err != nil {
		return err
	}

	// validate transactions in the block
	if err := chain.ValidateBlockBody(block); err != nil {
		return err
	}

	// all is well!
	return nil
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (dh *defaultDposHelper) verifyHeader(dpos *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header,
	refHeader *types.Header, verifySigs bool, verifyProposers bool) error {

	var (
		number    = header.Number.Uint64()
		isImpeach = header.Impeachment()
	)

	if number == 0 {
		return nil
	}

	err := dh.verifyBasic(dpos, chain, header, parents, refHeader)
	if err != nil {
		return err
	}

	// verify dpos seal, genesis block not need this check
	if verifyProposers {
		if isImpeach {
			if err := dh.verifyDposSnapImpeach(dpos, chain, header, parents, refHeader); err != nil {
				log.Warn("verifying dpos snap of impeach failed", "error", err, "hash", header.Hash().Hex())
				return err
			}

		} else {
			// verify proposers
			if err := dh.verifyProposers(dpos, chain, header, parents, refHeader); err != nil {
				log.Warn("verifying proposers failed", "error", err, "hash", header.Hash().Hex())
				return err
			}

			// verify proposer's seal
			if err := dh.verifySeal(dpos, chain, header, parents, refHeader); err != nil {
				log.Warn("verifying seal failed", "error", err, "hash", header.Hash().Hex())
				return err
			}
		}
	}

	// verify dpos signatures if required
	if verifySigs {
		if err := dh.verifySignatures(dpos, chain, header, parents, refHeader); err != nil {
			log.Debug("verifying validator signatures failed", "error", err, "hash", header.Hash().Hex())
			return err
		}
	}

	return nil
}

// verifyBasic verifies basic fields of the header, i.e. Number, Hash, Coinbase, Time
func (dh *defaultDposHelper) verifyBasic(dpos *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error {

	// if nil number, return error
	if header.Number == nil {
		return errUnknownBlock
	}

	// switch dpos mode, test related
	switch dpos.Mode() {
	case DoNothingFakeMode:
		// do nothing
	case FakeMode:
		return nil
	case PbftFakeMode:
		return nil
	}

	var (
		number    = header.Number.Uint64()
		hash      = header.Hash()
		isImpeach = header.Impeachment()
	)

	if number == 0 {
		return nil
	}

	// Ensure the block's parent is valid
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		blk := chain.GetBlock(header.ParentHash, number-1)
		if blk != nil {
			parent = blk.Header()
			log.Debug("dpos_helper get block", "blk", blk.NumberU64(), "parent_is_nil", parent == nil)
		}
		if parent == nil {
			parent = chain.GetHeaderByNumber(number - 1)
			log.Debug("dpos_helper get block(blk is nil)", "parent_is_nil", parent == nil)
		}
	}

	// Ensure that the block's parent is valid
	if parent == nil {
		log.Debug("parent is nil when verifying the header", "number", number, "hash", hash)
		return consensus.ErrUnknownAncestor
	}
	if parent.Number.Uint64() != number-1 {
		log.Debug("parent's number is not equal to header.number when verifying the header", "number", number, "hash", hash, "parent.number", parent.Number.Uint64())
		return consensus.ErrUnknownAncestor
	}
	if parent.Hash() != header.ParentHash {
		log.Debug("parent's hash is not equal to header.parentHash when verifying the header", "number", number, "hash", hash, "parent.hash", parent.Hash().Hex(), "header.parentHash", header.ParentHash.Hex())
		return consensus.ErrUnknownAncestor
	}

	// If timestamp is in a valid field, wait for it, otherwise, return invalid timestamp.
	log.Debug("timestamp related values", "parent timestamp", parent.Timestamp(), "block timestamp", header.Timestamp(), "period", dpos.config.PeriodDuration(), "timeout", dpos.config.ImpeachTimeout)

	// Ensure that the block's timestamp is valid
	if dpos.Mode() == NormalMode && number > dpos.config.MaxInitBlockNumber && !isImpeach {

		if header.Timestamp().Before(parent.Timestamp().Add(dpos.config.PeriodDuration())) {
			return ErrInvalidTimestamp
		}
		if header.Timestamp().After(parent.Timestamp().Add(dpos.config.PeriodDuration()).Add(dpos.config.ImpeachTimeout)) {
			return ErrInvalidTimestamp
		}
	}

	// Ensure that the block's gasLimit is valid
	if header.GasLimit > configs.MaxGasLimit || header.GasLimit < configs.MinGasLimit || header.GasUsed > header.GasLimit {
		return ErrInvalidGasLimit
	}

	if isImpeach {
		return dh.verifyBasicImpeach(dpos, chain, header, parent)
	}

	// Delay to verify it!
	delay := header.Timestamp().Sub(time.Now())
	log.Debug("delaying to verify the block", "delay", delay)
	<-time.After(delay)

	return nil
}

// verifyBasicImpeach verifies basic fields of an impeach header, i.e. Number, Hash, Coinbase, Time
func (dh *defaultDposHelper) verifyBasicImpeach(dpos *Dpos, chain consensus.ChainReader, header *types.Header, parent *types.Header) error {

	expectedImpeachBlock := types.NewBlock(header, []*types.Transaction{}, []*types.Receipt{})
	expectedImpeachBlock.RefHeader().Extra = make([]byte, extraSeal)

	if header.StateRoot != parent.StateRoot {
		return consensus.ErrInvalidImpeachStateRoot
	}

	if header.TxsRoot != expectedImpeachBlock.TxsRoot() {
		return consensus.ErrInvalidImpeachTxsRoot
	}

	if header.ReceiptsRoot != expectedImpeachBlock.ReceiptsRoot() {
		return consensus.ErrInvalidImpeachReceiptsRoot
	}

	if header.LogsBloom != expectedImpeachBlock.LogsBloom() {
		return consensus.ErrInvalidImpeachLogsBloom
	}

	if header.GasLimit != parent.GasLimit {
		return consensus.ErrInvalidImpeachGasLimit
	}

	if header.GasUsed != 0 {
		return consensus.ErrInvalidImpeachGasUsed
	}

	if len(header.Extra) != len(expectedImpeachBlock.Extra()) {
		return consensus.ErrInvalidImpeachExtra
	}

	for i, x := range header.Extra {
		if x != expectedImpeachBlock.Extra()[i] {
			return consensus.ErrInvalidImpeachExtra
		}
	}

	return nil
}

// ProposersImpeach verifies dpos snap fields of an impeach header
func (dh *defaultDposHelper) verifyDposSnapImpeach(dpos *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error {

	// Ensure the block's parent is valid
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		blk := chain.GetBlock(header.ParentHash, header.Number.Uint64()-1)
		if blk != nil {
			parent = blk.Header()
		}
	}

	parentHeader := parent
	if parentHeader == nil {
		return consensus.ErrUnknownAncestor
	}

	expectedImpeachBlock, err := dpos.CreateImpeachBlockAt(parentHeader)
	if err != nil {
		return err
	}

	if header.Dpos.Seal != expectedImpeachBlock.Header().Dpos.Seal {
		return consensus.ErrInvalidImpeachDposSnap
	}

	if len(header.Dpos.Proposers) != len(expectedImpeachBlock.Header().Dpos.Proposers) {
		return consensus.ErrInvalidImpeachDposSnap
	}

	for i, x := range header.Dpos.Proposers {
		if x != expectedImpeachBlock.Header().Dpos.Proposers[i] {
			return consensus.ErrInvalidImpeachDposSnap
		}
	}

	if len(header.Dpos.Validators) != len(expectedImpeachBlock.Header().Dpos.Validators) {
		return consensus.ErrInvalidImpeachDposSnap
	}

	for i, x := range header.Dpos.Validators {
		if x != expectedImpeachBlock.Header().Dpos.Validators[i] {
			return consensus.ErrInvalidImpeachDposSnap
		}
	}

	return nil
}

// verifyProposers verifies dpos proposers
func (dh *defaultDposHelper) verifyProposers(dpos *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error {

	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}

	// Retrieve the Snapshot needed to verify this header and cache it
	snap, err := dh.snapshot(dpos, chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	// Check proposers
	proposers := snap.ProposersOf(number)
	if !reflect.DeepEqual(header.Dpos.Proposers, proposers) {
		if dpos.Mode() == NormalMode {
			log.Debug("err: invalid proposer list")
			log.Debug("~~~~~~~~~~~~~~~~~~~~~~~~")
			log.Debug("proposers in block dpos snap:")
			for round, signer := range header.Dpos.Proposers {
				log.Debug("proposer", "addr", signer.Hex(), "idx", round)
			}

			log.Debug("~~~~~~~~~~~~~~~~~~~~~~~~")
			log.Debug("proposers in snapshot:")
			for round, signer := range proposers {
				log.Debug("validator", "addr", signer.Hex(), "idx", round)
			}

			log.Debug("~~~~~~~~~~~~~~~~~~~~~~~~")
			log.Debug("recent proposers: ")
			for i := snap.TermOf(number); i < snap.TermOf(number)+5; i++ {
				log.Debug("----------------------")
				log.Debug("proposers in snapshot of:", "term idx", i)
				for _, s := range snap.getRecentProposers(i) {
					log.Debug("signer", "s", s.Hex())
				}
			}

			return consensus.ErrInvalidSigners
		}
	}

	return nil
}

// Snapshot retrieves the authorization Snapshot at a given point in time.
// @param chainSeg  the segment of a chain, composed by ancestors and the block(specified by parameter [number] and [hash])
// in the order of ascending block number.
func (dh *defaultDposHelper) snapshot(dpos *Dpos, chain consensus.ChainReader, number uint64, hash common.Hash, chainSeg []*types.Header) (*DposSnapshot, error) {
	// Search for a Snapshot in memory or on disk for checkpoints
	var (
		headers []*types.Header
		snap    *DposSnapshot
	)

	log.Debug("defaultDposHelper snapshot", "number", number, "hash", hash.Hex(), "len(parent and itself)", len(chainSeg))

	// return early if already know it!
	if dpos.CurrentSnap() != nil && dpos.CurrentSnap().hash() == hash {
		return dpos.CurrentSnap(), nil
	}

	numberIter := number
	for snap == nil {

		// If an in-memory snapshot can be found, use it
		if got, ok := dpos.recentSnaps.Get(hash); ok {
			snap = got.(*DposSnapshot)
			break
		}

		// If an on-disk checkpoint Snapshot can be found, use that
		log.Debug("loading snapshot", "number", numberIter, "hash", hash.Hex())
		s, err := loadSnapshot(dpos.config, dpos.db, hash)
		if err == nil {
			log.Debug("Loaded checkpoint Snapshot from disk", "number", numberIter, "hash", hash.Hex())
			snap = s
			break
		} else {
			log.Debug("loading snapshot fails", "error", err)
		}

		// If we're at block zero, make a Snapshot
		if numberIter == 0 {
			// Retrieve genesis block and verify it
			genesis := chain.GetHeaderByNumber(0)
			if err := dpos.dh.verifyHeader(dpos, chain, genesis, nil, nil, false, false); err != nil {
				return nil, err
			}

			var proposers []common.Address
			var validators []common.Address
			if dpos.Mode() == FakeMode || dpos.Mode() == DoNothingFakeMode {
				// do nothing when test,empty proposers assigned
			} else {
				// Create a snapshot from the genesis block
				proposers = genesis.Dpos.CopyProposers()
				validators = genesis.Dpos.CopyValidators()
			}
			snap = newSnapshot(dpos.config, 0, genesis.Hash(), proposers, validators, FakeMode)
			if err := snap.store(dpos.db); err != nil {
				return nil, err
			}
			log.Debug("Stored genesis voting Snapshot to disk")
			break
		}

		// No Snapshot for this header, gather the header and move backward
		var header *types.Header
		if len(chainSeg) > 0 {
			// If we have explicit chainSeg, pick from there (enforced)
			header = chainSeg[len(chainSeg)-1]
			if header.Hash() != hash || header.Number.Uint64() != numberIter {
				return nil, consensus.ErrUnknownAncestor
			}
			chainSeg = chainSeg[:len(chainSeg)-1]
		} else {
			// No explicit chainSeg (or no more left), reach out to the database
			header = chain.GetHeader(hash, numberIter)
			if header == nil {
				return nil, consensus.ErrUnknownAncestor
			}
		}

		headers = append(headers, header)
		numberIter, hash = numberIter-1, header.ParentHash
	}

	// Previous Snapshot found, apply any pending headers on top of it
	for i := 0; i < len(headers)/2; i++ {
		headers[i], headers[len(headers)-1-i] = headers[len(headers)-1-i], headers[i]
	}

	var (
		candidateService = dpos.GetCandidateBackend()
		rptService       = dpos.GetRptBackend()
	)

	var timeToUpdateCommittee bool
	_, headNumber := chain.KnownHead()

	log.Debug("known chain head", "number", headNumber)

	timeToUpdateCommittee = dpos.IsMiner() || dpos.IsValidator()
	startBlockNumberOfRptCalculate := float64(int(headNumber) - configs.DefaultFullSyncPivot)
	timeToUpdateRpts := float64(snap.number()) > math.Max(0., startBlockNumberOfRptCalculate)
	timeToUpdateCommittee = timeToUpdateCommittee && timeToUpdateRpts

	log.Debug("now apply a batch of headers to get a new snap")

	applyStartTime := time.Now()

	// Apply headers to the snapshot and updates RPTs
	newSnap, err := snap.apply(headers, timeToUpdateCommittee, candidateService, rptService)
	if err != nil {
		return nil, err
	}

	log.Debug("now created a new snap", "number", newSnap.number(), "hash", newSnap.hash().Hex(), "apply elapsed", common.PrettyDuration(time.Now().Sub(applyStartTime)))

	// Save to cache
	dpos.recentSnaps.Add(newSnap.hash(), newSnap)

	// If we've generated a new checkpoint Snapshot, save to disk
	if err = newSnap.store(dpos.db); err != nil {
		log.Warn("failed to store dpos snapshot", "error", err)
		return nil, err
	}
	log.Debug("Stored snapshot to disk", "number", newSnap.number(), "hash", newSnap.hash().Hex())

	if dpos.CurrentSnap() == nil || (dpos.CurrentSnap() != nil && newSnap.number() >= dpos.CurrentSnap().number()) {
		dpos.SetCurrentSnap(newSnap)
	}

	return newSnap, err
}

// verifySeal checks whether the dpos seal is signature of a correct proposer.
// The method accepts an optional list of parent headers that aren't yet part of the local blockchain to generate
// the snapshots from.
func (dh *defaultDposHelper) verifySeal(dpos *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error {
	hash := header.Hash()
	number := header.Number.Uint64()

	// Verifying the genesis block is not supported
	if number == 0 {
		return errUnknownBlock
	}

	// Fake Dpos doesn't do seal check
	if dpos.Mode() == FakeMode || dpos.Mode() == DoNothingFakeMode {
		time.Sleep(dpos.fakeDelay)
		if dpos.fakeFail == number {
			return errFakerFail
		}
		return nil
	}

	// Resolve the authorization key and check against signers
	proposer, _, err := dh.ecrecover(header, dpos.finalSigs)
	if err != nil {
		return err
	}

	// Retrieve the Snapshot needed to verify this header and cache it
	snap, err := dh.snapshot(dpos, chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	// Some debug infos here
	log.Debug("--------dpos.verifySeal--------")
	log.Debug("hash", "hash", hash.Hex())
	log.Debug("number", "number", number)
	if chain.CurrentBlock() != nil {
		log.Debug("current header", "number", chain.CurrentBlock().NumberU64())
	}
	log.Debug("proposer", "address", proposer.Hex())

	// Check if the proposer is right proposer
	ok, err := snap.IsProposerOf(proposer, number)
	if err != nil {
		return err
	}
	// If proposer is a wrong leader, return err
	if !ok {
		return consensus.ErrUnauthorized
	}

	return nil
}

// verifySignatures verifies whether the signatures of the header is signed by correct validator committee
func (dh *defaultDposHelper) verifySignatures(dpos *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error {
	var (
		number    = header.Number.Uint64()
		hash      = header.Hash()
		isImpeach = header.Impeachment()
	)

	// Verifying the genesis block is not supported
	if number == 0 {
		return errUnknownBlock
	}

	// Fake Dpos doesn't do seal check
	if dpos.Mode() == FakeMode || dpos.Mode() == DoNothingFakeMode {
		time.Sleep(dpos.fakeDelay)
		if dpos.fakeFail == number {
			return errFakerFail
		}
		return nil
	}

	// Resolve the authorization keys
	proposer, validators, err := dh.ecrecover(header, dpos.finalSigs)
	if err != nil {
		return err
	}

	// Retrieve the Snapshot needed to verify this header and cache it
	snap, err := dh.snapshot(dpos, chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	expectValidators := snap.ValidatorsOf(number)

	// Some debug infos here
	log.Debug("--------dpos.verifySigs--------")
	log.Debug("hash", "hash", hash.Hex())
	log.Debug("number", "number", number)
	if chain.CurrentBlock() != nil {
		log.Debug("current header", "number", chain.CurrentBlock().NumberU64())
	}
	log.Debug("proposer", "address", proposer.Hex())

	defaultValidators, _ := dpos.ValidatorsOf(chain.CurrentHeader().Number.Uint64())
	log.Debug("number of validators", "count", len(defaultValidators))

	log.Debug("validators recovered from header: ")
	for idx, validator := range validators {
		log.Debug("validator", "addr", validator.Hex(), "idx", idx)
	}
	log.Debug("validators in snapshot: ")
	for idx, signer := range expectValidators {
		log.Debug("validator", "addr", signer.Hex(), "idx", idx)
	}

	count := 0
	for _, v := range validators {
		for _, ev := range expectValidators {
			if v == ev {
				count++
			}
		}
	}

	// if not reached to 2f + 1, the validation fails
	if !isImpeach {
		if !dpos.config.Certificate(uint64(count)) {
			return consensus.ErrNotEnoughSigs
		}
	} else {
		if !dpos.config.ImpeachCertificate(uint64(count)) {
			return consensus.ErrNotEnoughSigs
		}
	}

	// pass
	return nil
}

// signHeader signs the given refHeader if self is in the committee
func (dh *defaultDposHelper) signHeader(dpos *Dpos, chain consensus.ChainReader, header *types.Header, state consensus.State) error {
	hash := header.Hash()
	number := header.Number.Uint64()

	// Retrieve the Snapshot needed to verify this header and cache it
	snap, err := dh.snapshot(dpos, chain, number-1, header.ParentHash, nil)
	if err != nil {
		log.Warn("getting dpos snapshot failed", "error", err)
		return err
	}

	var s interface{}
	var ok bool
	// Retrieve signatures of the block in cache
	if state == consensus.Commit || state == consensus.ImpeachCommit {
		s, ok = dpos.finalSigs.Get(hash) // check if it needs a lock
		if !ok || s == nil {
			s = &signatures{
				sigs: make(map[common.Address][]byte),
			}
			dpos.finalSigs.Add(hash, s)
		}
	} else if state == consensus.Prepare || state == consensus.ImpeachPrepare {
		s, ok = dpos.prepareSigs.Get(hash)
		if !ok || s == nil {
			s = &signatures{
				sigs: make(map[common.Address][]byte),
			}
			dpos.prepareSigs.Add(hash, s)
		}
	} else {
		log.Warn("the state is unexpected for signing header", "state", state)
		return errInvalidStateForSign
	}

	// Copy all signatures to allSigs
	allSigs := make([]types.DposSignature, dpos.config.ValidatorsLen())
	validators := snap.ValidatorsOf(number)
	if dpos.config.ValidatorsLen() != uint64(len(validators)) {
		log.Warn("validator committee length not equal to validators length", "config.ValidatorsLen", dpos.config.ValidatorsLen(), "validatorLen", len(validators))
	}

	// fulfill all known validator signatures to dpos.sigs to accumulate
	for signPos, signer := range snap.ValidatorsOf(number) {
		if sigHash, ok := s.(*signatures).getSig(signer); ok {
			copy(allSigs[signPos][:], sigHash)
		}
	}
	header.Dpos.Sigs = allSigs

	// Sign the block if self is in the committee
	if snap.IsValidatorOf(dpos.Coinbase(), number) {
		// NOTE: sign a block only once
		if signedHash, signed := dpos.IfSigned(number); signed && signedHash != header.Hash() && state != consensus.ImpeachPrepare && state != consensus.ImpeachCommit {
			return errMultiBlocksInOneHeight
		}

		// get hash with state
		hashToSign, err := hashBytesWithState(dpos.dh.sigHash(header).Bytes(), state)
		if err != nil {
			log.Warn("failed to get hash bytes with state", "number", number, "hash", hash.Hex(), "state", state)
			return err
		}

		// Sign it
		sighash, err := dpos.SignHash(hashToSign)
		if err != nil {
			log.Warn("signing block header failed", "error", err)
			return err
		}

		// if the sigs length is wrong, reset it with correct ValidatorsLen
		if len(header.Dpos.Sigs) != int(snap.config.ValidatorsLen()) {
			header.Dpos.Sigs = make([]types.DposSignature, snap.config.ValidatorsLen())
		}

		// mark as signed
		err = dpos.MarkAsSigned(number, hash)
		if err != nil {
			return err
		}

		// Copy signer's signature to the right position in the allSigs
		sigPos, _ := snap.ValidatorViewOf(dpos.Coinbase(), number)
		copy(header.Dpos.Sigs[sigPos][:], sighash)

		// Record new sig to signature cache
		s.(*signatures).setSig(dpos.Coinbase(), sighash)

		return nil
	}
	log.Warn("signing block failed", "error", errValidatorNotInCommittee)
	return errValidatorNotInCommittee
}
