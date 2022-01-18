

package admission

import (
	"sync"

	"github.com/gcchains/chain/accounts/keystore"
	"github.com/gcchains/chain/api/rpc"
	contracts "github.com/gcchains/chain/contracts/dpos/campaign/tests"
)

// ApiBackend interface provides the common JSON-RPC API.
type ApiBackend interface {
	// APIs returns the collection of RPC services the admission package offers.
	Apis() []rpc.API

	// IsRNode returns true or false indicate whether the node is RNode
	IsRNode() (bool, error)

	// FundForRNode sends money to reward contract to become RNode
	FundForRNode() error

	// Campaign starts running all the proof work to generate the campaign information and waits all proof work done, send msg
	Campaign(times uint64) error

	// Abort cancels all the proof work associated to the workType.
	Abort()

	// GetStatus gets status of campaign
	GetStatus() (workStatus, error)

	// getResult returns the work proof result
	GetResult() map[string]Result

	// SetAdmissionKey sets the key for admission control to participate campaign
	SetAdmissionKey(key *keystore.Key)

	// AdmissionKey returns keystore key
	AdmissionKey() *keystore.Key

	// RegisterInProcHandler registers the rpc.Server, handles RPC request to process the API requests in process
	RegisterInProcHandler(localRPCServer *rpc.Server)

	SetContractBackend(contractBackend contracts.Backend)

	// IgnoreNetworkCheck tells ac backend to ignore network status check
	IgnoreNetworkCheck()
}

// ProofWork represent a proof work
type ProofWork interface {
	// prove starts memory/cpu/... POW work.
	prove(abort <-chan interface{}, wg *sync.WaitGroup)

	// error returns err if proof work is abnormal
	error() error

	// result returns the work proof result
	result() Result
}
