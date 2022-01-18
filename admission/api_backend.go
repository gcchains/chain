

package admission

import (
	"github.com/gcchains/chain/accounts/keystore"
	"github.com/gcchains/chain/api/gcclient"
	"github.com/gcchains/chain/api/rpc"
	"github.com/gcchains/chain/consensus"
	contracts "github.com/gcchains/chain/contracts/dpos/campaign/tests"
	"github.com/ethereum/go-ethereum/common"
)

type AdmissionApiBackend struct {
	admissionControl *AdmissionControl
}

func NewAdmissionApiBackend(chain consensus.ChainReader, address common.Address, admissionContractAddr common.Address,
	campaignContractAddr common.Address, rNodeContractAddr common.Address, networkContractAddr common.Address) ApiBackend {
	return &AdmissionApiBackend{
		admissionControl: NewAdmissionControl(chain, address, admissionContractAddr, campaignContractAddr, rNodeContractAddr, networkContractAddr),
	}
}

// APIs returns the collection of RPC services the admission package offers.
func (b *AdmissionApiBackend) Apis() []rpc.API {
	return []rpc.API{
		{
			Namespace: "admission",
			Version:   "1.2",
			Service:   b,
			Public:    false,
		},
	}
}

func (b *AdmissionApiBackend) IgnoreNetworkCheck() {
	b.admissionControl.IgnoreNetworkCheck()
}

func (b *AdmissionApiBackend) CheckNetworkStatus() bool {
	return b.admissionControl.CheckNetworkStatus()
}

func (b *AdmissionApiBackend) FundForRNode() error {
	return b.admissionControl.FundForRNode()
}

func (b *AdmissionApiBackend) IsRNode() (bool, error) {
	return b.admissionControl.IsRNode()
}

func (b *AdmissionApiBackend) Campaign(terms uint64) error {
	return b.admissionControl.Campaign(terms)
}

func (b *AdmissionApiBackend) Abort() {
	b.admissionControl.Abort()
}

func (b *AdmissionApiBackend) GetStatus() (workStatus, error) {
	return b.admissionControl.GetStatus()
}

func (b *AdmissionApiBackend) GetResult() map[string]Result {
	return b.admissionControl.GetResult()
}

func (b *AdmissionApiBackend) SetAdmissionKey(key *keystore.Key) {
	b.admissionControl.SetAdmissionKey(key)
}

func (b *AdmissionApiBackend) AdmissionKey() *keystore.Key {
	return b.admissionControl.key
}

// RegisterInProcHandler registers the rpc.Server, handles RPC request to process the API requests in process
func (b *AdmissionApiBackend) RegisterInProcHandler(localRPCServer *rpc.Server) {
	client := rpc.DialInProc(localRPCServer)
	b.admissionControl.setClientBackend(gcclient.NewClient(client))
}

func (b *AdmissionApiBackend) SetContractBackend(contractBackend contracts.Backend) {
	b.admissionControl.SetSimulateBackend(contractBackend)
}
