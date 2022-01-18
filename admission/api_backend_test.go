package admission_test

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"testing"

	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/accounts/abi/bind/backends"
	"github.com/gcchains/chain/accounts/keystore"
	"github.com/gcchains/chain/admission"
	"github.com/gcchains/chain/api/rpc"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/consensus/dpos"
	acContracts "github.com/gcchains/chain/contracts/dpos/admission"
	campaign "github.com/gcchains/chain/contracts/dpos/campaign"
	"github.com/gcchains/chain/contracts/dpos/network"
	rnode "github.com/gcchains/chain/contracts/dpos/rnode"
	"github.com/gcchains/chain/core"
	"github.com/gcchains/chain/core/vm"
	"github.com/gcchains/chain/database"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	key      *keystore.Key
	addr     common.Address
	key1, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae17c6668d313eac2f96dbcda3f391")
	addr1    = crypto.PubkeyToAddress(key1.PublicKey)
	keyPath  = "../examples/gcchain/conf-dev/keys"
	password = "password"
	ks       *keystore.KeyStore

	alloc = core.GenesisAlloc{
		addr:  {Balance: big.NewInt(1000000000)},
		addr1: {Balance: big.NewInt(1000000000)},
	}
	// gspec = core.Genesis{Config: params.TestChainConfig, Alloc: alloc}
	gspec   = core.Genesis{Config: configs.ChainConfigInfo(), Alloc: alloc}
	testdb  = database.NewMemDatabase()
	genesis = gspec.MustCommit(testdb)
)

func deployAdmission(prvKey *ecdsa.PrivateKey, cpuDifficulty, memoryDifficulty, cpuWorkTimeout, memoryWorkTimeout uint64, backend *backends.SimulatedBackend) (common.Address, error) {
	deployTransactor := bind.NewKeyedTransactor(prvKey)
	acAddr, _, _, err := acContracts.DeployAdmission(deployTransactor, backend, new(big.Int).SetUint64(cpuDifficulty), new(big.Int).SetUint64(memoryDifficulty), new(big.Int).SetUint64(cpuWorkTimeout), new(big.Int).SetUint64(memoryWorkTimeout))
	if err != nil {
		return common.Address{}, err
	}
	backend.Commit()
	return acAddr, nil
}

func deployRnode(prvKey *ecdsa.PrivateKey, backend *backends.SimulatedBackend) (common.Address, *rnode.Rnode, error) {
	deployTransactor := bind.NewKeyedTransactor(prvKey)
	rNodeAddr, _, rNodeContract, err := rnode.DeployRnode(deployTransactor, backend)
	if err != nil {
		return common.Address{}, nil, err
	}
	backend.Commit()
	return rNodeAddr, rNodeContract, nil
}

func deployNetwork(prvKey *ecdsa.PrivateKey, backend *backends.SimulatedBackend) (common.Address, *network.Network, error) {
	deployTransactor := bind.NewKeyedTransactor(prvKey)
	networkAddr, _, networkContract, err := network.DeployNetwork(deployTransactor, backend)
	if err != nil {
		return common.Address{}, nil, err
	}
	backend.Commit()
	return networkAddr, networkContract, nil
}

func deployCampaign(prvKey *ecdsa.PrivateKey, backend *backends.SimulatedBackend, admissionContract common.Address, rNodeContract common.Address) (common.Address, error) {
	deployTransactor := bind.NewKeyedTransactor(prvKey)
	campaignAddr, _, _, err := campaign.DeployCampaign(deployTransactor, backend, admissionContract, rNodeContract)
	if err != nil {
		return common.Address{}, err
	}
	backend.Commit()
	return campaignAddr, nil
}

func deployRequiredContracts(t *testing.T) (*backends.SimulatedBackend, common.Address, common.Address, *rnode.Rnode, common.Address, common.Address) {
	var (
		cpuDifficulty     uint64 = 5
		memDifficulty     uint64 = 5
		cpuWorkTimeout    uint64 = 5
		memoryWorkTimeout uint64 = 5
	)

	contractBackend := backends.NewDposSimulatedBackend(core.GenesisAlloc{addr: {Balance: new(big.Int).Mul(big.NewInt(1000000), big.NewInt(configs.Gcc))}})
	admissionAddr, err := deployAdmission(key.PrivateKey, cpuDifficulty, memDifficulty, cpuWorkTimeout, memoryWorkTimeout, contractBackend)
	if err != nil {
		t.Fatal("deploy error", "error", err)
	}

	rNodeAddr, rNodeContract, err := deployRnode(key.PrivateKey, contractBackend)
	if err != nil {
		t.Fatal("deploy error", "error", err)
	}

	networkAddr, _, err := deployNetwork(key.PrivateKey, contractBackend)

	campaignAddr, err := deployCampaign(key.PrivateKey, contractBackend, admissionAddr, rNodeAddr)
	if err != nil {
		t.Fatal("deploy error", "error", err)
	}
	return contractBackend, admissionAddr, rNodeAddr, rNodeContract, campaignAddr, networkAddr
}

func init() {
	ks = keystore.NewKeyStore(keyPath, 2, 1)
	account := ks.Accounts()[0]
	account, k, err := ks.GetDecryptedKey(account, password)
	key = k
	if err != nil {
		new(testing.T).Fatal(err.Error())
	}

	ks.Unlock(account, password)
	addr = crypto.PubkeyToAddress(key.PrivateKey.PublicKey)
}

func newDummyChain() consensus.ChainReader {
	_, _ = core.GenerateChain(configs.TestChainConfig, genesis, dpos.New(configs.ChainConfigInfo().Dpos, testdb), testdb, nil, 8, nil)
	chain, _ := core.NewBlockChain(testdb, nil, configs.TestChainConfig, dpos.New(configs.ChainConfigInfo().Dpos, testdb), vm.Config{}, nil, nil)

	return chain
}

// newAC return a new AdmissionControl instance
func newAcApiBackend(chain consensus.ChainReader, admissionContractAddr common.Address, campaignContractAddr common.Address, rNodeContractAddr common.Address, networkContractAddr common.Address) admission.ApiBackend {
	return admission.NewAdmissionApiBackend(chain, addr, admissionContractAddr, campaignContractAddr, rNodeContractAddr, networkContractAddr)
}

// TestAPIs test apis
func TestApis(t *testing.T) {
	ac := newAcApiBackend(newDummyChain(), common.Address{}, common.Address{}, common.Address{}, common.Address{})
	apis := ac.Apis()

	wantApis := []rpc.API{
		{
			Namespace: "admission",
			Version:   "1.0",
			Service:   ac,
			Public:    false,
		},
	}
	if !reflect.DeepEqual(apis, wantApis) {
		t.Fatalf("ac.APIs want(%+v), but(%+v)", wantApis, apis)
	}
}

// TestCampaign tests campaign, check status, abort and check status
func TestAdmissionApiBackend_Campaign(t *testing.T) {
	contractBackend, admissionAddr, rNodeAddr, _, campaignAddr, networkAddr := deployRequiredContracts(t)

	ac := newAcApiBackend(contractBackend.Blockchain(), admissionAddr, campaignAddr, rNodeAddr, networkAddr)
	status, err := ac.GetStatus()
	var wantErr error
	if status != admission.AcIdle || !reflect.DeepEqual(err, wantErr) {
		t.Fatalf("Before starting campaign: GetStatus, want(status:%d, err:%v), but(status:%d, err:%v)\n", admission.AcIdle, wantErr, status, err)
	}
	ac.SetContractBackend(contractBackend)
	ac.SetAdmissionKey(key)
	ac.FundForRNode()
	contractBackend.Commit()
	ac.Campaign(1)
	status, err = ac.GetStatus()
	if status != admission.AcRunning || !reflect.DeepEqual(err, wantErr) {
		t.Fatalf("Started compaign: GetStatus, want(status:%d, err:%v), but(status:%d, err:%v)\n", admission.AcRunning, wantErr, status, err)
	}

	ac.Abort()
	status, err = ac.GetStatus()
	wantErr = admission.ErrPowAbort
	if status != admission.AcIdle || !reflect.DeepEqual(err, wantErr) {
		t.Fatalf("Aborted campaign: GetStatus, want(status:%d, err:%v), but(status:%d, err:%v)\n", admission.AcIdle, wantErr, status, err)
	}
}

// TestIsRNode returns a bool value indicating if the current node is RNode
func TestAdmissionApiBackend_IsRNode(t *testing.T) {
	contractBackend, admissionAddr, rewardAddr, _, campaignAddr, networkAddr := deployRequiredContracts(t)
	ac := newAcApiBackend(contractBackend.Blockchain(), admissionAddr, campaignAddr, rewardAddr, networkAddr)
	ac.SetContractBackend(contractBackend)
	ac.SetAdmissionKey(key)
	isRNode, _ := ac.IsRNode()
	if isRNode {
		t.Error("IsRNode() should return false before the node deposit enough money to reward")
	}

	err := ac.FundForRNode()
	if err != nil {
		t.Error("encounter error when funding money for RNode", "error", err)
	}
	contractBackend.Commit()

	isRNode, _ = ac.IsRNode()
	if !isRNode {
		t.Error("fund for rnode failed")
	}
}

func TestAdmissionApiBackend_FundForRNode(t *testing.T) {
	contractBackend, admissionAddr, rNodeAddr, _, campaignAddr, networkAddr := deployRequiredContracts(t)
	ac := newAcApiBackend(contractBackend.Blockchain(), admissionAddr, campaignAddr, rNodeAddr, networkAddr)
	ac.SetContractBackend(contractBackend)
	ac.SetAdmissionKey(key)

	isRNode, _ := ac.IsRNode()
	if isRNode {
		t.Error("IsRNode() should return false before the node deposit enough money to reward")
	}

	err := ac.FundForRNode()
	if err != nil {
		t.Error("encounter error when funding money for RNode", "error", err)
	}
	contractBackend.Commit()

	isRNode, _ = ac.IsRNode()
	if !isRNode {
		t.Error("IsRNode() should return true after the new round started")
	}

	b, _ := contractBackend.BalanceAt(context.Background(), rNodeAddr, nil)
	t.Log("the balance of reward", b.String())
	if b.String() != "200000000000000000000000" {
		t.Error("reward balance is not correct")
	}

	b, _ = contractBackend.BalanceAt(context.Background(), addr, nil)
	t.Log("the balance of the current node", b.String())
	tmp := new(big.Int).Div(b, big.NewInt(configs.Gcc))
	if new(big.Int).Sub(tmp, big.NewInt(800000)).CmpAbs(big.NewInt(1)) < 0 {
		t.Log("the remain balance should be approximately 800000 CPC", "actual", tmp)
	}

	// re-fund for RNode, expect it is not redundantly sending money to reward contract
	err = ac.FundForRNode()
	if err != nil {
		t.Error("encounter error when funding money for RNode", "error", err)
	}
	contractBackend.Commit()
	b2, _ := contractBackend.BalanceAt(context.Background(), addr, nil)
	t.Log("the balance of the current node", b2.String())
	if b2.Cmp(b) != 0 {
		t.Error("the balance should not change because it is already RNode and not need to send money to reward contract")
	}
}
