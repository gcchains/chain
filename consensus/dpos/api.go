package dpos

import (
	"github.com/gcchains/chain/api/rpc"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain consensus.ChainReader
	dpos  *Dpos
	dh    *defaultDposHelper
}

// GetSnapshot retrieves the state Snapshot at a given block.
func (api *API) GetSnapshot(number rpc.BlockNumber) (*DposSnapshot, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == 0 || number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its Snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.dpos.dh.snapshot(api.dpos, api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetSnapshotAtHash retrieves the state Snapshot at a given block.
func (api *API) GetSnapshotAtHash(hash common.Hash) (*DposSnapshot, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	return api.dpos.dh.snapshot(api.dpos, api.chain, header.Number.Uint64(), header.Hash(), nil)
}

// GetProposers retrieves the Proposers at a given block.
func (api *API) GetProposers(number rpc.BlockNumber) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == 0 || number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its Proposers
	if header == nil {
		return nil, errUnknownBlock
	}
	return header.Dpos.Proposers, nil
}

// GetValidators retrieves the Validators at a given block.
func (api *API) GetValidators(number rpc.BlockNumber) ([]common.Address, error) {
	return api.dpos.ValidatorsOf(uint64(number))
}

// GetRNodes retrieves current RNodes.
func (api *API) GetRNodes() ([]common.Address, error) {
	return api.dpos.GetRNodes()
}
