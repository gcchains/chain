

package dpos

import (
	"bytes"
	"errors"
	"math/big"
	"time"

	"github.com/gcchains/chain/api/rpc"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/consensus/dpos/backend"
	"github.com/gcchains/chain/core/state"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

// Dpos proof-of-reputation protocol constants.
const (
	defaultCampaignTerms = uint64(3) // Default number of terms to campaign for proposer committee.

	extraVanity = 32 // Fixed number of extra-data prefix bytes reserved for signer vanity
	extraSeal   = 65 // Fixed number of extra-data suffix bytes reserved for signer seal
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")

	// errInvalidCheckpointBeneficiary is returned if a checkpoint/epoch transition
	// block has a beneficiary set to non-zeroes.
	errInvalidCheckpointBeneficiary = errors.New("beneficiary in checkpoint block non-zero")

	// errMissingVanity is returned if a block's extra-data section is shorter than
	// 32 bytes, which is required to store the signer vanity.
	errMissingVanity = errors.New("extra-data 32 byte vanity prefix missing")

	// errMissingSignature is returned if a block's extra-data section doesn't seem
	// to contain a 65 byte secp256k1 signature.
	errMissingSignature = errors.New("extra-data 65 byte suffix signature missing")

	// ErrInvalidTimestamp is returned if the timestamp of a block is lower than
	// the previous block's timestamp + the minimum block period.
	ErrInvalidTimestamp = errors.New("invalid timestamp")

	// ErrInvalidGasLimit is returned if the gasLimit of a block is invalid
	ErrInvalidGasLimit = errors.New("invalid gas limit for the block")

	// errInvalidChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidChain = errors.New("invalid voting chain")

	// --- new error types ---

	// errMultiBlocksInOneHeight is returned if there is multi blocks in one height in the chain.
	errMultiBlocksInOneHeight = errors.New("multi blocks in one height")

	// errInvalidValidatorSigs is returned if the dpos sigs are not sigend by correct validator committtee.
	errInvalidValidatorSigs = errors.New("invalid validator signatures")

	// errNoSigsInCache is returned if the cache is unable to store and return sigs.
	errNoSigsInCache = errors.New("signatures not found in cache")

	errFakerFail = errors.New("error fake fail")

	// --- our new error types ---

	// errVerifyUncleNotAllowed is returned when verify uncle block.
	errVerifyUncleNotAllowed = errors.New("uncles not allowed")

	// errWaitTransactions is returned if an empty block is attempted to be sealed
	// on an instant chain (0 second period). It's important to refuse these as the
	// block reward is zero, so an empty block just bloats the chain... fast.
	errWaitTransactions = errors.New("waiting for transactions")

	errInvalidStateForSign = errors.New("the state is unexpected for signing header")
)

// Author implements consensus.Engine, returning the gcchain address recovered
// from the signature in the header's extra-data section.
func (d *Dpos) Author(header *types.Header) (common.Address, error) {
	proposer, _, err := d.dh.ecrecover(header, d.finalSigs)
	return proposer, err
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (d *Dpos) VerifyHeader(chain consensus.ChainReader, header *types.Header, verifySigs bool, refHeader *types.Header) error {
	return d.dh.verifyHeader(d, chain, header, nil, refHeader, verifySigs, false)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers. The
// method returns a quit channel to abort the operations and a results channel to
// retrieve the async verifications (the order is that of the input slice).
func (d *Dpos) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, verifySigs []bool, refHeaders []*types.Header) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for i, header := range headers {
			err := d.dh.verifyHeader(d, chain, header, headers[:i], refHeaders[i], verifySigs[i], false)

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// VerifySeal implements consensus.Engine, checking whether the signature contained
// in the header satisfies the consensus protocol requirements.
func (d *Dpos) VerifySeal(chain consensus.ChainReader, header *types.Header, refHeader *types.Header) error {
	return d.dh.verifySeal(d, chain, header, nil, refHeader)
}

// VerifySigs checks if header has enough signatures of validators.
func (d *Dpos) VerifySigs(chain consensus.ChainReader, header *types.Header, refHeader *types.Header) error {
	return d.dh.verifySignatures(d, chain, header, nil, refHeader)
}

// PrepareBlock implements consensus.Engine, preparing all the consensus fields of the
// header for running the transactions on top.
func (d *Dpos) PrepareBlock(chain consensus.ChainReader, header *types.Header) error {
	number := header.Number.Uint64()

	// Create a snapshot
	snap, err := d.dh.snapshot(d, chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}

	// Ensure the extra data has all its components
	if len(header.Extra) < extraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, extraVanity-len(header.Extra))...)
	}
	header.Extra = header.Extra[:extraVanity]

	for _, proposer := range snap.ProposersOf(number) {
		header.Dpos.Proposers = append(header.Dpos.Proposers, proposer)
	}

	log.Debug("prepare a block", "number", header.Number.Uint64(), "proposers", header.Dpos.ProposersFormatText(),
		"validators", header.Dpos.ValidatorsFormatText())

	// Set correct signatures size
	header.Dpos.Sigs = make([]types.DposSignature, d.config.ValidatorsLen())

	// Ensure the timestamp has the correct delay
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		log.Warn("consensus.ErrUnknownAncestor 4", "number", number, "parentHash", header.ParentHash.Hex())
		return consensus.ErrUnknownAncestor
	}

	header.SetTimestamp(parent.Timestamp().Add(d.config.PeriodDuration()))
	if header.Timestamp().Before(time.Now()) {
		header.SetTimestamp(time.Now())
	}
	return nil
}

// TryCampaign tries to start campaign
func (d *Dpos) TryCampaign() {
	if d.ac == nil {
		// it is not able to campaign in the situation
		log.Debug("It is not able to campaign in the situation")
		return
	}

	snap := d.CurrentSnap()
	if snap != nil {
		isV := snap.IsValidatorOf(d.coinbase, snap.Number)
		log.Debug("check if participate campaign", "isToCampaign", d.IsToCampaign(), "isStartCampaign", snap.isStartCampaign(), "number", snap.number(), "isValidator", isV)

		if d.IsToCampaign() && snap.isAboutToCampaign() && !isV {
			// make sure it is a RNode, because only RNode has permission to participate campaign
			isRNode, err := d.ac.IsRNode()
			if err != nil {
				log.Debug("encounter error when invoke IsRNode()", "error", err)
				return
			}

			if !isRNode {
				log.Info("It is not RNode, cannot participate campaign")
				if err := d.ac.FundForRNode(); err != nil {
					log.Debug("failed to FundForRNode", "detail", err)
					return
				}
				log.Info("already send money to become RNode")
			}
		}

		if d.IsToCampaign() && snap.isStartCampaign() && !isV {
			newTerm := d.CurrentSnap().TermOf(snap.Number)
			if newTerm > d.lastCampaignTerm+defaultCampaignTerms-1 {
				d.lastCampaignTerm = newTerm
				log.Info("campaign for proposer committee", "eleTerm", newTerm)
				if d.ac != nil {
					d.ac.Campaign(defaultCampaignTerms)
				}
			}
		}
	}
}

// GetBlockReward returns block reward according to block number
func (d *Dpos) GetBlockReward(blockNum uint64) *big.Int {
	reward := getBlockReward(new(big.Int).SetUint64(blockNum))
	return reward
}

func getBlockReward(number *big.Int) *big.Int {
	var amount *big.Int
	if number.Cmp(configs.Cep1LastBlockY1()) <= 0 {
		amount = configs.Cep1BlockRewardY1()
	} else if number.Cmp(configs.Cep1LastBlockY2()) <= 0 {
		amount = configs.Cep1BlockRewardY2()
	} else if number.Cmp(configs.Cep1LastBlockY3()) <= 0 {
		amount = configs.Cep1BlockRewardY3()
	} else if number.Cmp(configs.Cep1LastBlockY4()) <= 0 {
		amount = configs.Cep1BlockRewardY4()
	} else if number.Cmp(configs.Cep1LastBlockY5()) <= 0 {
		amount = configs.Cep1BlockRewardY5()
	} else {
		amount = big.NewInt(0)
	}
	return amount
}

func addCoinbaseReward(coinbase common.Address, state *state.StateDB, number *big.Int) {
	amount := getBlockReward(number)
	state.AddBalance(coinbase, amount)
}

// Finalize implements consensus.Engine, ensuring no uncles are set, nor block
// rewards given, and returns the final block.
func (d *Dpos) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {

	if (header.Coinbase != common.Address{}) {
		addCoinbaseReward(header.Coinbase, state, header.Number)
	}

	// last step
	header.StateRoot = state.IntermediateRoot(true)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, receipts), nil
}

// Authorize injects a private key into the consensus engine to mint new blocks
// with.
func (d *Dpos) Authorize(signer common.Address, signFn backend.SignFn) {
	d.coinbaseLock.Lock()
	d.coinbase = signer
	d.signFn = signFn
	d.coinbaseLock.Unlock()

	if d.handler == nil {
		d.handler = backend.NewHandler(d.config, d.Coinbase(), d.db)
	}
	if d.handler.Coinbase() != signer {
		d.handler.SetCoinbase(signer)
	}
}

// Seal implements consensus.Engine, attempting to create a sealed block using
// the local signing credentials.
// NB please populate the correct field values.  we are now removing some fields such as nonce.
func (d *Dpos) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	var (
		header = block.Header()
		number = header.Number.Uint64()

		coinbase = d.Coinbase()
		signFn   = d.SignHash
	)

	// Sealing the genesis block is not supported
	if number == 0 {
		return nil, errUnknownBlock
	}

	// For 0-period chains, refuse to seal empty blocks (no reward but would spin sealing)
	if d.config.Period == 0 && len(block.Transactions()) == 0 {
		return nil, errWaitTransactions
	}

	// Bail out if we're unauthorized to sign a block
	snap, err := d.dh.snapshot(d, chain, number-1, header.ParentHash, nil)
	if err != nil {
		return nil, err
	}

	ok, err := snap.IsProposerOf(coinbase, number)
	if err != nil {
		if err == errProposerNotInCommittee {
			return nil, consensus.ErrNotInProposerCommittee
		}

		log.Debug("Error occurs when seal block", "error", err)
		return nil, err

	}
	if !ok {
		return nil, consensus.ErrUnauthorized
	}

	// Sweet, the protocol permits us to sign the block, wait for our time
	delay := header.Timestamp().Sub(time.Now())
	log.Debug("Waiting for slot to sign and propagate", "delay", delay)

	select {
	case <-stop:
		log.Warn("Quit block sealing", "number", block.NumberU64(), "hash", block.Hash().Hex())

		return nil, nil
	case <-time.After(delay):
		log.Debug("wait for seal", "delay", delay)
	}

	// Proposer seals the block with signature
	sighash, err := signFn(d.dh.sigHash(header).Bytes())
	if err != nil {
		return nil, err
	}
	copy(header.Dpos.Seal[:], sighash)

	// Create a signature space for validators
	header.Dpos.Sigs = make([]types.DposSignature, len(header.Dpos.Validators))
	log.Debug("sealed the block", "hash", header.Hash().Hex(), "number", header.Number)

	// Update dpos current snapshot
	d.SetCurrentSnap(snap)

	return block.WithSeal(header), nil
}

// CanMakeBlock checks if the given coinbase is ready to propose a block
func (d *Dpos) CanMakeBlock(chain consensus.ChainReader, coinbase common.Address, parent *types.Header) bool {
	number := parent.Number.Uint64()
	// Bail out if we're unauthorized to sign a block
	snap, err := d.dh.snapshot(d, chain, number, parent.Hash(), nil)
	if err != nil {
		log.Debug("Error occurs when create a snapshot", "error", err)
		return false
	}

	log.Debug("created an snapshot")

	// check if it is the in-charge proposer for next block
	ok, err := snap.IsProposerOf(coinbase, number+1)
	if err != nil {
		log.Debug("it is not proposer", "msg", err)
		return false
	}
	log.Debug("now can finished CanMakeBlock call", "ok", ok)
	return ok
}

// APIs implements consensus.Engine, returning the user facing RPC API to allow
// controlling the signer voting.
func (d *Dpos) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "dpos",
		Version:   "1.0",
		Service:   &API{chain: chain, dpos: d},
		Public:    false,
	}}
}

// State returns current pbft phrase, one of (PrePrepare, Prepare, Commit).
func (d *Dpos) State() consensus.State {
	d.stateLock.Lock()
	defer d.stateLock.Unlock()
	return d.pbftState
}

// GetCalcRptInfo get the rpt value of an address at specific block number
func (d *Dpos) GetCalcRptInfo(address common.Address, addresses []common.Address, blockNum uint64) int64 {
	rptService := d.GetRptBackend()
	if rptService == nil {
		log.Fatal("dpos rpt service is nil")
	}
	rp := rptService.CalcRptInfo(address, addresses, blockNum)
	return rp.Rpt
}
