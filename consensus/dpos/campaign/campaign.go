

package campaign

import (
	"math/big"

	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/commons/log"
	campaignContract "github.com/gcchains/chain/contracts/dpos/campaign"
	"github.com/ethereum/go-ethereum/common"
)

// CandidateService provides methods to obtain all candidates from campaign contract
type CandidateService interface {
	CandidatesOf(term uint64) ([]common.Address, error)
}

// CandidateServiceImpl is the default candidate list collector
type CandidateServiceImpl struct {
	client   bind.ContractBackend
	contract common.Address
}

// NewCampaignService creates a concrete candidate service instance.
func NewCampaignService(campaignContract common.Address, backend bind.ContractBackend) (CandidateService, error) {

	rs := &CandidateServiceImpl{
		contract: campaignContract,
		client:   backend,
	}
	return rs, nil
}

// CandidatesOf implements CandidateService
func (rs *CandidateServiceImpl) CandidatesOf(term uint64) ([]common.Address, error) {

	// new campaign contract instance
	contractInstance, err := campaignContract.NewCampaign(rs.contract, rs.client)
	if err != nil {
		log.Debug("error when create campaign instance", "err", err)
		return nil, err
	}

	// candidates from new campaign contract
	cds, err := contractInstance.CandidatesOf(nil, new(big.Int).SetUint64(term))
	if err != nil {
		log.Debug("error when read candidates from campaign", "err", err)
		return nil, err
	}

	log.Debug("now read candidates from campaign contract", "len", len(cds), "contract addr", rs.contract.Hex())
	return cds, nil
}
