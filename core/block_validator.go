

package core

import (
	"fmt"

	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/core/state"
	"github.com/gcchains/chain/types"
)

// BlockValidator is responsible for validating block headers, uncles and
// processed state.
//
// BlockValidator implements Validator.
type BlockValidator struct {
	config *configs.ChainConfig // Chain configuration options
	bc     *BlockChain          // Canonical block chain
	engine consensus.Engine     // Consensus engine used for validating
}

// NewBlockValidator returns a new block validator which is safe for re-use
func NewBlockValidator(config *configs.ChainConfig, blockchain *BlockChain, engine consensus.Engine) *BlockValidator {
	validator := &BlockValidator{
		config: config,
		engine: engine,
		bc:     blockchain,
	}
	return validator
}

// ValidateBody validates the header's transaction root.
// The headers are assumed to be already validated at this point.
func (v *BlockValidator) ValidateBody(block *types.Block) error {
	// check whether the block's known, and if not, that it's linkable
	if v.bc.HasBlockAndState(block.Hash(), block.NumberU64()) {
		return ErrKnownBlock
	}

	if !v.bc.HasBlockAndState(block.ParentHash(), block.NumberU64()-1) {
		// we do not have the parent block
		if !v.bc.HasBlock(block.ParentHash(), block.NumberU64()-1) {
			return consensus.ErrUnknownAncestor
		}
		// we have the parent block but its state is pruned
		return consensus.ErrPrunedAncestor
	}

	// if a block already exists, but the state is missing.  we will also try to insert it.
	// header validity is known at this point, check the transactions validity
	header := block.Header()
	if hash := types.DeriveSha(block.Transactions()); hash != header.TxsRoot {
		return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash, header.TxsRoot)
	}
	return nil
}

// ValidateState validates the various changes that happen after a state
// transition, such as amount of used gas, the receipt roots and the state root
// itself. ValidateState returns a database batch if the validation was a success
// otherwise nil and an error is returned.
func (v *BlockValidator) ValidateState(block, parent *types.Block, statedb *state.StateDB, receipts types.Receipts, usedGas uint64) error {
	// this is a copy of the block header.
	header := block.Header()
	if block.GasUsed() != usedGas {
		return fmt.Errorf("invalid gas used (remote: %d local: %d)", block.GasUsed(), usedGas)
	}
	// Validate the received block's bloom with the one derived from the generated receipts.
	// For valid blocks this should always validate to true.
	rbloom := types.CreateBloom(receipts)
	if rbloom != header.LogsBloom {
		return fmt.Errorf("invalid bloom (remote: %x  local: %x)", header.LogsBloom, rbloom)
	}
	// Tre receipt Trie's root (R = (Tr [[H1, R1], ... [Hn, R1]]))
	receiptSha := types.DeriveSha(receipts)
	if receiptSha != header.ReceiptsRoot {
		return fmt.Errorf("invalid receipt root hash (remote: %x local: %x)", header.ReceiptsRoot, receiptSha)
	}
	// Validate the state root against the received state root and throw
	// an error if they don't match.
	if root := statedb.IntermediateRoot(true); header.StateRoot != root {
		return fmt.Errorf("invalid merkle root (remote: %x local: %x)", header.StateRoot, root)
	}
	return nil
}

// CalcGasLimit computes the gas limit of the next block after parent.
// This is miner strategy, not consensus protocol.
func CalcGasLimit(parent *types.Block) uint64 {
	// contrib = (parentGasUsed * 3 / 2) / 1024
	contrib := (parent.GasUsed() + parent.GasUsed()/2) / configs.GasLimitBoundDivisor

	// decay = parentGasLimit / 1024 -1
	decay := parent.GasLimit()/configs.GasLimitBoundDivisor - 1

	/*
		strategy: gasLimit of block-to-mine is set based on parent's
		gasUsed value.  if parentGasUsed > parentGasLimit * (2/3) then we
		increase it, otherwise lower it (or leave it unchanged if it's right
		at that usage) the amount increased/decreased depends on how far away
		from parentGasLimit * (2/3) parentGasUsed is.
	*/
	limit := parent.GasLimit() - decay + contrib
	if limit < configs.MinGasLimit {
		limit = configs.MinGasLimit
	}
	// however, if we're now below the target (TargetGasLimit) we increase the
	// limit as much as we can (parentGasLimit / 1024 -1)
	if limit < configs.TargetGasLimit {
		limit = parent.GasLimit() + decay
		if limit > configs.TargetGasLimit {
			limit = configs.TargetGasLimit
		}
	}
	if limit > configs.MaxGasLimit {
		limit = configs.MaxGasLimit
	}
	return limit
}
