package backend

import (
	"errors"
	"sync"
	"time"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/database"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	maxPendingBlocks = 16
)

var (
	// ErrUnknownHandlerMode is returned if in an unknown mode
	ErrUnknownHandlerMode = errors.New("Unknown dpos handler mode")

	// ErrFailToAddPendingBlock is returned if failed to add block to pending
	ErrFailToAddPendingBlock = errors.New("fail to add pending block")

	// ErrNotSigner is returned if i am not a signer when handshaking
	// with remote signer
	ErrNotSigner = errors.New("local peer is not in the GC committees")
)

// Handler implements PbftHandler
type Handler struct {
	mode   HandlerMode
	config *configs.DposConfig

	available bool
	coinbase  common.Address
	lock      sync.RWMutex

	dialer *Dialer
	lbft   *LBFT
	dpos   DposService
	fsm    ConsensusStateMachine

	knownBlocks           *RecentBlocks
	unknownAncestorBlocks *RecentBlocks
	pendingBlockCh        chan *types.Block
	pendingImpeachBlockCh chan *types.Block
	quitCh                chan struct{}

	broadcastRecord   *broadcastRecord
	impeachmentRecord *impeachmentRecord
}

// NewHandler creates a new Handler
func NewHandler(config *configs.DposConfig, coinbase common.Address, db database.Database) *Handler {

	h := &Handler{
		config:                config,
		available:             false,
		coinbase:              coinbase,
		dialer:                NewDialer(),
		knownBlocks:           NewRecentBlocks(db),
		unknownAncestorBlocks: NewRecentBlocks(db),
		pendingBlockCh:        make(chan *types.Block),
		pendingImpeachBlockCh: make(chan *types.Block),
		quitCh:                make(chan struct{}),
		broadcastRecord:       newBroadcastRecord(),
		impeachmentRecord:     newImpeachmentRecord(),
	}

	// h.mode = LBFTMode
	h.mode = LBFT2Mode

	return h
}

// Start starts handler
func (h *Handler) Start() {

	// dial default validators
	go h.dialer.dialAllRemoteValidators(0)
	go h.dialer.KeepConnection()

	// broadcast mined pending block loop
	go h.PendingBlockBroadcastLoop()

	// broadcast impeachment block loop
	go h.PendingImpeachBlockBroadcastLoop()

	// unknown ancestor block handler
	go h.procUnknownAncestorsLoop()
}

// Stop stops all
func (h *Handler) Stop() {
	h.dialer.Stop()

	close(h.quitCh)
	h.quitCh = make(chan struct{})

	return
}

// GetProtocol returns handler protocol
func (h *Handler) GetProtocol() consensus.Protocol {
	return h
}

// NodeInfo returns node status
func (h *Handler) NodeInfo() interface{} {
	return h.dpos.Status()
	// TODO: fix this
	// Identity, Number, Hash, State
}

// Name returns protocol name
func (h *Handler) Name() string {
	return ProtocolName
}

// Version returns protocol version
func (h *Handler) Version() uint {
	return ProtocolVersion
}

// Length returns protocol max msg code
func (h *Handler) Length() uint64 {
	return ProtocolLength
}

// AddPeer adds a p2p peer to local peer set
func (h *Handler) AddPeer(version int, p *p2p.Peer, rw p2p.MsgReadWriter) (string, bool, bool, error) {
	blk := h.dpos.GetCurrentBlock()
	if blk == nil {
		log.Error("current block is nil", "block", blk)
		return "", false, false, errNilBlock
	}

	var (
		number      = blk.NumberU64()
		currentTerm = h.dpos.TermOf(number)
		futureTerm  = h.dpos.FutureTermOf(number)
	)

	mac, sig, err := h.dpos.GetMac()
	if err != nil {
		log.Fatal("err when get message authentication code", "err", err)
	}

	return h.dialer.AddPeer(version, p, rw, mac, sig, currentTerm, futureTerm)
}

// RemovePeer removes a p2p peer with its addr
func (h *Handler) RemovePeer(addr string) {

	log.Debug("removing dpos peer", "addr", addr)

	_ = h.dialer.removeRemoteProposers(addr)
	_ = h.dialer.removeRemoteValidators(addr)
}

// HandleMsg handles a msg of peer with id "addr"
func (h *Handler) HandleMsg(addr string, version int, p *p2p.Peer, rw p2p.MsgReadWriter, msg p2p.Msg) (string, error) {

	remoteValidator, isV := h.dialer.getValidator(addr)
	remoteProposer, isP := h.dialer.getProposer(addr)

	if isV {
		return addr, h.handleMsg(remoteValidator.RemoteSigner, msg)
	} else if isP {
		return addr, h.handleMsg(remoteProposer.RemoteSigner, msg)
	} else {
		// TODO: the remote proposer is not in current proposer list, fix this
		log.Debug("handling remote proposer connection msg", "remote addr", p.RemoteAddr().String(), "local addr", p.LocalAddr().String(), "coinbase", addr, "msgcode", msg.Code)
		return "", nil
	}
}

func (h *Handler) handleMsg(p *RemoteSigner, msg p2p.Msg) error {
	if msg.Code == NewSignerMsg {
		log.Debug("received NewSignerMsg", "coinbase", p.Coinbase().Hex(), "remote addr", p.RemoteAddr().String(), "local addr", p.LocalAddr().String())
		return nil
	}

	switch h.mode {
	case LBFTMode:
		return h.handleLBFTMsg(msg, p)
	case LBFT2Mode:
		return h.handleLBFT2Msg(msg, p)
	default:
		return ErrUnknownHandlerMode
	}
}

// SetServer sets dialer.server
func (h *Handler) SetServer(server *p2p.Server) {
	h.dialer.SetServer(server)
}

// SetDposService sets dpos service to handler
func (h *Handler) SetDposService(dpos DposService) {
	h.dpos = dpos
	h.dialer.SetDposService(dpos)
}

// SetDposStateMachine sets dpos state machine
func (h *Handler) SetDposStateMachine(fsm ConsensusStateMachine) {
	h.fsm = fsm
}

// Coinbase returns handler.signer
func (h *Handler) Coinbase() common.Address {
	h.lock.Lock()
	defer h.lock.Unlock()

	return h.coinbase
}

// SetCoinbase sets coinbase of handler
func (h *Handler) SetCoinbase(coinbase common.Address) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.coinbase != coinbase {
		h.coinbase = coinbase
	}
}

// SetAvailable sets available
func (h *Handler) SetAvailable() {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.available = true
}

// Available returns if handler is available
func (h *Handler) Available() bool {
	h.lock.RLock()
	defer h.lock.RUnlock()

	return h.available
}

func (h *Handler) procUnknownAncestorsLoop() {
	for {
		for _, bi := range h.unknownAncestorBlocks.GetBlockIdentifiers() {

			// if less than current number, drop it!
			blk := h.dpos.GetCurrentBlock()
			if blk == nil {
				continue
			}

			if bi.number <= blk.NumberU64() {

				h.unknownAncestorBlocks.RemoveBlock(bi)
				log.Debug("unknown ancestor block's number is less than current number, drop it!", "number", bi.number, "hash", bi.hash.Hex())

				continue
			}

			// handle this unknown ancestor block!
			block, err := h.unknownAncestorBlocks.GetBlock(bi)
			if block != nil && err == nil {
				var msg p2p.Msg
				size, r, err := rlp.EncodeToReader(block)
				if err != nil {
					log.Warn("failed to encode unknown ancestor block", "err", err)
					continue
				}

				if block.Impeachment() {
					// impeach block
					msg = p2p.Msg{Code: PreprepareImpeachBlockMsg, Size: uint32(size), Payload: r}

				} else {
					// not impeach block
					msg = p2p.Msg{Code: PreprepareBlockMsg, Size: uint32(size), Payload: r}

				}

				go h.handleLBFT2Msg(msg, nil)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (h *Handler) PeerInfos() ([]*PeerInfo, error) {
	return h.dialer.PeerInfos()
}
