package backend

import (
	"time"

	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	defaultWaitDialTimes = 20
)

func waitForEnoughValidators(h *Handler, term uint64, quitCh chan struct{}) (validators map[common.Address]*RemoteValidator) {
	for i := 0; i < defaultWaitDialTimes; i++ {
		select {
		case <-quitCh:
			return

		default:

			if validators, enough := h.dialer.EnoughValidatorsOfTerm(term); enough {
				return validators
			}

			time.Sleep((h.dpos.Period() + h.dpos.ImpeachTimeout()) / defaultWaitDialTimes)
		}
	}
	return
}

func waitForEnoughImpeachValidators(h *Handler, term uint64, quitCh chan struct{}) (validators map[common.Address]*RemoteValidator) {
	for i := 0; i < defaultWaitDialTimes; i++ {
		select {
		case <-quitCh:
			return

		default:

			if validators, enough := h.dialer.EnoughImpeachValidatorsOfTerm(term); enough {
				return validators
			}

			time.Sleep((h.dpos.Period() + h.dpos.ImpeachTimeout()) / defaultWaitDialTimes)
		}
	}
	return
}

// waitForEnoughPreprepareValidators is used by proposer when trying to broadcast a preprepare block msg
func waitForEnoughPreprepareValidators(h *Handler, term uint64, quitCh chan struct{}, deadline time.Time) (validators map[common.Address]*RemoteValidator) {
	for i := 0; i < defaultWaitDialTimes; i++ {
		select {
		case <-quitCh:
			return

		default:

			if validators, enough := h.dialer.EnoughImpeachValidatorsOfTerm(term); enough || time.Now().After(deadline) {
				return validators
			}

			time.Sleep(500 * time.Millisecond)
		}
	}
	return
}

// ProposerBroadcastPreprepareBlock broadcasts generated block to validators
func (h *Handler) ProposerBroadcastPreprepareBlock(block *types.Block) {

	log.Debug("broadcasting preprepare block", "number", block.NumberU64(), "hash", block.Hash().Hex())

	term := h.dpos.TermOf(block.NumberU64())
	deadline := block.Timestamp().Add(h.dpos.BlockDelay()).Add(-1 * time.Second)
	validators := waitForEnoughPreprepareValidators(h, term, h.quitCh, deadline)

	for _, peer := range validators {
		peer.AsyncSendPreprepareBlock(block)
	}
}

// BroadcastPreprepareBlock broadcasts generated block to validators
func (h *Handler) BroadcastPreprepareBlock(block *types.Block) {

	log.Debug("broadcasting preprepare block", "number", block.NumberU64(), "hash", block.Hash().Hex())

	term := h.dpos.TermOf(block.NumberU64())
	validators := waitForEnoughValidators(h, term, h.quitCh)

	for _, peer := range validators {
		peer.AsyncSendPreprepareBlock(block)
	}
}

// BroadcastPreprepareImpeachBlock broadcasts generated impeach block to validators
func (h *Handler) BroadcastPreprepareImpeachBlock(block *types.Block) {

	log.Debug("broadcasting preprepare impeach block", "number", block.NumberU64(), "hash", block.Hash().Hex())

	term := h.dpos.TermOf(block.NumberU64())
	validators := waitForEnoughImpeachValidators(h, term, h.quitCh)

	for _, peer := range validators {
		peer.AsyncSendPreprepareImpeachBlock(block)
	}
}

// BroadcastPrepareHeader broadcasts signed prepare header to remote validators
func (h *Handler) BroadcastPrepareHeader(header *types.Header) {

	log.Debug("broadcasting prepare header", "number", header.Number.Uint64(), "hash", header.Hash().Hex())

	term := h.dpos.TermOf(header.Number.Uint64())
	validators := waitForEnoughValidators(h, term, h.quitCh)

	for _, peer := range validators {
		peer.AsyncSendPrepareHeader(header)
	}
}

// BroadcastPrepareImpeachHeader broadcasts signed impeach prepare header to remote validators
func (h *Handler) BroadcastPrepareImpeachHeader(header *types.Header) {

	log.Debug("broadcasting prepare impeach header", "number", header.Number.Uint64(), "hash", header.Hash().Hex())

	term := h.dpos.TermOf(header.Number.Uint64())
	validators := waitForEnoughImpeachValidators(h, term, h.quitCh)

	for _, peer := range validators {
		peer.AsyncSendPrepareImpeachHeader(header)
	}
}

// BroadcastCommitHeader broadcasts signed commit header to remote validators
func (h *Handler) BroadcastCommitHeader(header *types.Header) {

	log.Debug("broadcasting commit header", "number", header.Number.Uint64(), "hash", header.Hash().Hex())

	term := h.dpos.TermOf(header.Number.Uint64())
	validators := waitForEnoughValidators(h, term, h.quitCh)

	for _, peer := range validators {
		peer.AsyncSendCommitHeader(header)
	}
}

// BroadcastCommitImpeachHeader broadcasts signed impeach commit header to remote validators
func (h *Handler) BroadcastCommitImpeachHeader(header *types.Header) {

	log.Debug("broadcasting commit impeach header", "number", header.Number.Uint64(), "hash", header.Hash().Hex())

	term := h.dpos.TermOf(header.Number.Uint64())
	validators := waitForEnoughImpeachValidators(h, term, h.quitCh)

	for _, peer := range validators {
		peer.AsyncSendCommitImpeachHeader(header)
	}
}

// BroadcastValidateBlock broadcasts validate block to validators
func (h *Handler) BroadcastValidateBlock(block *types.Block) {

	log.Debug("broadcasting validate block", "number", block.NumberU64(), "hash", block.Hash().Hex())

	term := h.dpos.TermOf(block.NumberU64())
	validators := waitForEnoughValidators(h, term, h.quitCh)

	for _, peer := range validators {
		peer.AsyncSendValidateBlock(block)
	}
}

// BroadcastValidateImpeachBlock broadcasts validate impeach block to validators
func (h *Handler) BroadcastValidateImpeachBlock(block *types.Block) {

	log.Debug("broadcasting validate impeach block", "number", block.NumberU64(), "hash", block.Hash().Hex())

	term := h.dpos.TermOf(block.NumberU64())
	validators := waitForEnoughImpeachValidators(h, term, h.quitCh)

	for _, peer := range validators {
		peer.AsyncSendImpeachValidateBlock(block)
	}
}

// PendingBlockBroadcastLoop loops to broadcast blocks
func (h *Handler) PendingBlockBroadcastLoop() {

	for {
		select {
		case pendingBlock := <-h.pendingBlockCh:

			// broadcast mined pending block to remote signers
			go h.ProposerBroadcastPreprepareBlock(pendingBlock)

		case <-h.quitCh:
			return
		}
	}
}

// PendingImpeachBlockBroadcastLoop loops to broadcasts pending impeachment block
func (h *Handler) PendingImpeachBlockBroadcastLoop() {
	for {
		select {
		case impeachBlock := <-h.pendingImpeachBlockCh:

			if h.mode == LBFT2Mode {
				size, r, err := rlp.EncodeToReader(impeachBlock)
				if err != nil {
					log.Warn("failed to encode composed impeach block", "err", err)
					continue
				}
				msg := p2p.Msg{Code: PreprepareImpeachBlockMsg, Size: uint32(size), Payload: r}

				var (
					number = impeachBlock.NumberU64()
					hash   = impeachBlock.Hash()
				)

				// only validators can impeach
				isValidator, err := h.dpos.VerifyValidatorOf(h.Coinbase(), h.dpos.TermOf(number))

				if !h.impeachmentRecord.ifImpeached(number, hash) && isValidator && err == nil {

					// handle the impeach block
					go h.handleLBFT2Msg(msg, nil)

					h.impeachmentRecord.markAsImpeached(number, hash)
				}
			}

		case <-h.quitCh:
			return
		}
	}
}
