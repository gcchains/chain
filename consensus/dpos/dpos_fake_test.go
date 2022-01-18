

package dpos

import (
	"errors"
	"math/big"
	"time"

	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

type FakeReader struct {
	consensus.ChainReader
}

func (*FakeReader) Config() *configs.ChainConfig {
	// TODO @hmw populate this config
	return &configs.ChainConfig{}
}

func (*FakeReader) GetHeaderByNumber(number uint64) *types.Header {
	return &types.Header{Number: big.NewInt(0), Time: big.NewInt(0).Sub(big.NewInt(time.Now().Unix()), big.NewInt(100))}
}

type fakeDposUtil struct {
	dposUtil
	success bool
}

type fakeDposHelper struct {
	dposUtil
	verifySuccess   bool
	snapshotSuccess bool
}

func (f *fakeDposHelper) verifySignatures(d *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error {
	panic("implement me")
}

func (f *fakeDposHelper) verifyHeader(d *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header, verifySigs bool, verifyProposers bool) error {
	if f.verifySuccess {
		return nil
	}

	return errors.New("verify Header")
}

func (*fakeDposHelper) verifyBasic(c *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error {
	panic("implement me")
}

func (*fakeDposHelper) verifyProposers(c *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error {
	panic("implement me")
}

func (*fakeDposHelper) validateBlock(d *Dpos, chain consensus.ChainReader, block *types.Block, verifySigs bool, verifyProposers bool) error {
	return nil
}

func (f *fakeDposHelper) snapshot(c *Dpos, chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) (*DposSnapshot, error) {
	if f.snapshotSuccess {
		return &DposSnapshot{}, nil
	}

	return nil, errors.New("err")
}

func (*fakeDposHelper) verifySeal(c *Dpos, chain consensus.ChainReader, header *types.Header, parents []*types.Header, refHeader *types.Header) error {
	panic("implement me")
}

func (*fakeDposHelper) signHeader(d *Dpos, chain consensus.ChainReader, header *types.Header, state consensus.State) error {

	return nil
}
