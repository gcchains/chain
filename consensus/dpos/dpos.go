// Package dpos implements the dpos consensus engine.
package dpos

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gcchains/chain/accounts"
	"github.com/gcchains/chain/admission"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/consensus/dpos/backend"
	"github.com/gcchains/chain/consensus/dpos/campaign"
	"github.com/gcchains/chain/consensus/dpos/rnode"
	"github.com/gcchains/chain/consensus/dpos/rpt"
	"github.com/gcchains/chain/database"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
	lru "github.com/hashicorp/golang-lru"
)

var (
	errDposProtocolNotWorking = errors.New("Dpos Protocol is not working")
)

// BroadcastBlockFn broadcasts a block to normal peers(not pbft replicas)
type BroadcastBlockFn func(block *types.Block, prop bool)

// SyncFromPeerFn tries to sync blocks from given peer
type SyncFromPeerFn func(p *p2p.Peer)

// SyncFromBestPeerFn tries to sync blocks from best peer
type SyncFromBestPeerFn func()

const (
	inMemorySnapshots  = 100 // Number of recent vote snapshots to keep in memory
	inMemorySignatures = 100 // Number of recent block signatures to keep in memory
)

// Mode defines the type a dpos engine makes.
type Mode uint

// DposMode
const (
	NormalMode Mode = iota
	FakeMode
	DoNothingFakeMode
	PbftFakeMode
)

// Dpos is the proof-of-reputation consensus engine proposed to support the
// gcchain testnet.
type Dpos struct {
	dh     dposHelper
	db     database.Database   // Database to store and retrieve Snapshot checkpoints
	config *configs.DposConfig // Consensus engine configuration parameters

	recentSnaps *lru.ARCCache // Snapshots for recent block to speed up reorgs
	finalSigs   *lru.ARCCache // Final signatures of recent blocks to speed up mining
	prepareSigs *lru.ARCCache // The signatures of recent blocks for 'prepared' state

	signedBlocks *signedBlocksRecord // Record signed blocks.

	currentSnap     *DposSnapshot // Current snapshot
	currentSnapLock sync.RWMutex

	coinbase     common.Address // Coinbase of the miner(proposer or validator)
	signFn       backend.SignFn // Sign function to authorize hashes with
	coinbaseLock sync.RWMutex   // Protects the signer fields

	handler *backend.Handler

	isMiner     bool
	isMinerLock sync.RWMutex

	isValidator int32

	mode      Mode // used for test, always accept a block.
	fakeFail  uint64
	fakeDelay time.Duration // Time delay to sleep for before returning from verify
	modeLock  sync.RWMutex

	pbftState consensus.State
	stateLock sync.RWMutex

	ac admission.ApiBackend

	rNodeBackend    rnode.RNodeService
	rptBackend      rpt.RptService
	campaignBackend campaign.CandidateService

	chain consensus.ChainReadWriter

	pmBroadcastBlockFn   BroadcastBlockFn
	pmSyncFromPeerFn     SyncFromPeerFn
	pmSyncFromBestPeerFn SyncFromBestPeerFn

	quitSync chan struct{}

	lastCampaignTerm uint64 // the last term which the node has participated in campaign
	isToCampaign     int32  // indicate whether or not participate campaign, only elected proposer node can do mining
	// indicate whether the miner is running, there is a case that the dpos is running mining while campaign is stop,
	// it is by design and actually it does not generate any block in this case.
	runningMiner         int32
	validatorInitialized int32
}

// SignHash signs a hash msg with dpos coinbase account
func (d *Dpos) SignHash(hash []byte) ([]byte, error) {
	d.coinbaseLock.Lock()
	defer d.coinbaseLock.Unlock()

	var (
		coinbase = d.coinbase
		account  = accounts.Account{Address: coinbase}
	)

	return d.signFn(account, hash)
}

// IsMiner returns if local coinbase is a miner(proposer or validator)
func (d *Dpos) IsMiner() bool {
	d.isMinerLock.RLock()
	defer d.isMinerLock.RUnlock()

	return d.isMiner
}

// SetAsMiner sets local coinbase as a miner
func (d *Dpos) SetAsMiner(isMiner bool) {
	d.isMinerLock.Lock()
	defer d.isMinerLock.Unlock()

	d.isMiner = isMiner
}

// IsValidator returns if the node is running as a validator
func (d *Dpos) IsValidator() bool {
	return atomic.LoadInt32(&d.isValidator) == 1
}

// SetAsValidator sets the consensus engine working as a validator
func (d *Dpos) SetAsValidator(isValidator bool) {
	if isValidator {
		atomic.StoreInt32(&d.isValidator, 1)
	} else {
		atomic.StoreInt32(&d.isValidator, 0)
	}
}

// IsToCampaign returns if it is time to campaign
func (d *Dpos) IsToCampaign() bool {
	return atomic.LoadInt32(&d.isToCampaign) > 0
}

// SetToCampaign sets isToCampaign as true
func (d *Dpos) SetToCampaign(isToCampaign bool) {
	if isToCampaign {
		atomic.StoreInt32(&d.isToCampaign, 1)
	} else {
		atomic.StoreInt32(&d.isToCampaign, 0)
	}
}

// Mode returns dpos mode
func (d *Dpos) Mode() Mode {
	d.modeLock.RLock()
	defer d.modeLock.RUnlock()

	return d.mode
}

// CurrentSnap returns current dpos snapshot
func (d *Dpos) CurrentSnap() *DposSnapshot {
	d.currentSnapLock.RLock()
	defer d.currentSnapLock.RUnlock()

	return d.currentSnap
}

// SetCurrentSnap sets current dpos snapshot
func (d *Dpos) SetCurrentSnap(snap *DposSnapshot) {
	d.currentSnapLock.Lock()
	defer d.currentSnapLock.Unlock()

	d.currentSnap = snap
}

// New creates a Dpos proof-of-reputation consensus engine with the initial
// signers set to the ones provided by the user.
func New(config *configs.DposConfig, db database.Database) *Dpos {

	// Set any missing consensus parameters to their defaults
	conf := *config
	if conf.TermLen == 0 || conf.ViewLen == 0 {
		log.Fatal("wrong term length or view length configuration", "term length", conf.TermLen, "view length", conf.ViewLen)
		return nil
	}

	// Allocate the Snapshot caches and create the engine
	recentSnaps, _ := lru.NewARC(inMemorySnapshots)
	finalSigs, _ := lru.NewARC(inMemorySignatures)
	preparedSigs, _ := lru.NewARC(inMemorySignatures)

	signedBlocks := newSignedBlocksRecord(db)

	return &Dpos{
		dh:           &defaultDposHelper{&defaultDposUtil{}},
		config:       &conf,
		handler:      backend.NewHandler(&conf, common.Address{}, db),
		db:           db,
		recentSnaps:  recentSnaps,
		finalSigs:    finalSigs,
		prepareSigs:  preparedSigs,
		signedBlocks: signedBlocks,
	}
}

// NewFaker creates a new fake dpos
func NewFaker(config *configs.DposConfig, db database.Database) *Dpos {
	d := New(config, db)
	d.mode = FakeMode
	return d
}

// NewDoNothingFaker creates a new fake dpos, do nothing when verifying blocks
func NewDoNothingFaker(config *configs.DposConfig, db database.Database) *Dpos {
	d := New(config, db)
	d.mode = DoNothingFakeMode
	return d
}

// NewFakeFailer creates a new fake dpos, always fails when verifying blocks
func NewFakeFailer(config *configs.DposConfig, db database.Database, fail uint64) *Dpos {
	d := NewDoNothingFaker(config, db)
	d.fakeFail = fail
	return d
}

// NewFakeDelayer creates a new fake dpos, delays when verifying blocks
func NewFakeDelayer(config *configs.DposConfig, db database.Database, delay time.Duration) *Dpos {
	d := NewFaker(config, db)
	d.fakeDelay = delay
	return d
}

// NewPbftFaker creates a new fake dpos to work with pbft, not in use now
func NewPbftFaker(config *configs.DposConfig, db database.Database) *Dpos {
	d := New(config, db)
	d.mode = PbftFakeMode
	return d
}

// SetHandler sets dpos.handler
func (d *Dpos) SetHandler(handler *backend.Handler) error {
	d.handler = handler
	return nil
}

// IfSigned checks if already signed a block
func (d *Dpos) IfSigned(number uint64) (common.Hash, bool) {
	return d.signedBlocks.ifAlreadySigned(number)
}

// MarkAsSigned marks signed a hash as signed
func (d *Dpos) MarkAsSigned(number uint64, hash common.Hash) error {
	return d.signedBlocks.markAsSigned(number, hash)
}

// SetChain is called by test file to assign the value of Dpos.chain, as well as DPor.currentSnapshot
func (d *Dpos) SetChain(blockchain consensus.ChainReadWriter) {
	d.chain = blockchain

	header := d.chain.CurrentHeader()
	number := header.Number.Uint64()
	hash := header.Hash()

	snap, _ := d.dh.snapshot(d, d.chain, number, hash, nil)
	d.SetCurrentSnap(snap)
}

func (d *Dpos) SetupAsValidator(blockchain consensus.ChainReadWriter, server *p2p.Server, pmBroadcastBlockFn BroadcastBlockFn, pmSyncFromPeerFn SyncFromPeerFn, pmSyncFromBestPeerFn SyncFromBestPeerFn) {
	initialized := atomic.LoadInt32(&d.validatorInitialized) > 0
	// avoid launch handler twice
	if initialized {
		return
	}
	atomic.StoreInt32(&d.validatorInitialized, 1)

	d.initMinerAndValidator(blockchain, server, pmBroadcastBlockFn, pmSyncFromPeerFn, pmSyncFromBestPeerFn)
	return
}

// StartMining starts to create a handler and start it.
func (d *Dpos) StartMining(blockchain consensus.ChainReadWriter, server *p2p.Server, pmBroadcastBlockFn BroadcastBlockFn, pmSyncFromPeerFn SyncFromPeerFn, pmSyncFromBestPeerFn SyncFromBestPeerFn) {
	running := atomic.LoadInt32(&d.runningMiner) > 0
	// avoid launch handler twice
	if running {
		return
	}
	atomic.StoreInt32(&d.runningMiner, 1)

	d.initMinerAndValidator(blockchain, server, pmBroadcastBlockFn, pmSyncFromPeerFn, pmSyncFromBestPeerFn)
	return
}

func (d *Dpos) initMinerAndValidator(blockchain consensus.ChainReadWriter, server *p2p.Server, pmBroadcastBlockFn BroadcastBlockFn, pmSyncFromPeerFn SyncFromPeerFn, pmSyncFromBestPeerFn SyncFromBestPeerFn) {
	d.chain = blockchain

	d.pmBroadcastBlockFn = pmBroadcastBlockFn
	d.pmSyncFromPeerFn = pmSyncFromPeerFn
	d.pmSyncFromBestPeerFn = pmSyncFromBestPeerFn

	var (
		faulty  = d.config.FaultyNumber
		handler = d.handler
	)

	if d.IsValidator() {
		fsm := backend.NewLBFT2(faulty, d, handler.ReceiveImpeachBlock, handler.ReceiveFailbackImpeachBlock, d.db)
		handler.SetDposStateMachine(fsm)
	}

	handler.SetServer(server)
	handler.SetDposService(d)
	handler.SetAvailable()

	d.handler = handler

	log.Debug("set dpos handler available!")

	var (
		header = d.chain.CurrentHeader()
		hash   = header.Hash()
		number = header.Number.Uint64()
	)

	snap, _ := d.dh.snapshot(d, d.chain, number, hash, nil)
	d.SetCurrentSnap(snap)

	go d.handler.Start()

	return
}

// StopMining stops dpos engine
func (d *Dpos) StopMining() {
	running := atomic.LoadInt32(&d.runningMiner) > 0
	// avoid close twice
	if !running {
		return
	}
	atomic.StoreInt32(&d.runningMiner, 0)

	d.handler.Stop()
	return
}

// Coinbase returns current coinbase
func (d *Dpos) Coinbase() common.Address {
	d.coinbaseLock.RLock()
	defer d.coinbaseLock.RUnlock()

	return d.coinbase
}

// Protocol returns Dpos p2p protocol
func (d *Dpos) Protocol() consensus.Protocol {
	return d.handler.GetProtocol()
}

// PbftStatus returns current state of dpos
func (d *Dpos) PbftStatus() *consensus.PbftStatus {
	state := d.State()
	head := d.chain.CurrentHeader()
	return &consensus.PbftStatus{
		State: state,
		Head:  head,
	}
}

// HandleMinedBlock receives a block to add to handler's pending block channel
func (d *Dpos) HandleMinedBlock(block *types.Block) error {
	return d.handler.ReceiveMinedPendingBlock(block)
}

// ImpeachTimeout returns impeach time out
func (d *Dpos) ImpeachTimeout() time.Duration {
	return d.config.ImpeachTimeout
}

// SetupAdmission setups admission backend
func (d *Dpos) SetupAdmission(ac admission.ApiBackend) {
	d.ac = ac
}

func (d *Dpos) SetRptBackend(rptContract common.Address, backend backend.ClientBackend) {
	d.rptBackend, _ = rpt.NewRptService(rptContract, backend)
}

func (d *Dpos) GetRptBackend() rpt.RptService {
	return d.rptBackend
}

func (d *Dpos) SetCampaignBackend(campaignContract common.Address, backend backend.ClientBackend) {
	d.campaignBackend, _ = campaign.NewCampaignService(campaignContract, backend)
}

func (d *Dpos) GetCandidateBackend() campaign.CandidateService {
	return d.campaignBackend
}

func (d *Dpos) SetRNodeBackend(rNodeContract common.Address, backend backend.ClientBackend) {
	d.rNodeBackend, _ = rnode.NewRNodeService(rNodeContract, backend)
}

func (d *Dpos) GetRNodes() ([]common.Address, error) {
	if d.rNodeBackend != nil {
		return d.rNodeBackend.GetRNodes()
	}

	return []common.Address{}, nil
}
