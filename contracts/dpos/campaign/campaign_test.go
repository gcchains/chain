

package campaign_test

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/gcchains/chain/contracts/dpos/network"

	"github.com/gcchains/chain/types"

	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/accounts/abi/bind/backends"
	admission2 "github.com/gcchains/chain/admission"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/contracts/dpos/admission"
	campaign "github.com/gcchains/chain/contracts/dpos/campaign"
	contracts "github.com/gcchains/chain/contracts/dpos/campaign/tests"
	rnode "github.com/gcchains/chain/contracts/dpos/rnode"
	"github.com/gcchains/chain/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	key, _      = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c1668d313eac2f96dbcd13f291")
	addr        = crypto.PubkeyToAddress(key.PublicKey)
	numPerRound = 12
	initBalance = new(big.Int).Mul(big.NewInt(1000000), big.NewInt(configs.Gcc))
	version     = new(big.Int).SetInt64(1)
)

func deploy(prvKey *ecdsa.PrivateKey, amount *big.Int, backend *backends.SimulatedBackend) (campaignAddr common.Address, admissionAddr common.Address, rNodeAddr common.Address, networkAddr common.Address, err error) {
	deployTransactor := bind.NewKeyedTransactor(prvKey)
	addrRNode, _, _, err := rnode.DeployRnode(deployTransactor, backend)
	acAddr, _, _, err := admission.DeployAdmission(deployTransactor, backend, big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1))
	networkAddr, _, _, _ = network.DeployNetwork(deployTransactor, backend)
	addr, _, _, err := campaign.DeployCampaign(deployTransactor, backend, acAddr, addrRNode)
	if err != nil {
		return common.Address{}, common.Address{}, common.Address{}, common.Address{}, err
	}
	backend.Commit()
	return addr, acAddr, addrRNode, networkAddr, nil
}

func fundToCampaign(prvKey *ecdsa.PrivateKey, rNodeAddr common.Address, backend *backends.SimulatedBackend) error {
	transactOpts := bind.NewKeyedTransactor(prvKey)
	rNodeContract, err := rnode.NewRnode(rNodeAddr, backend)
	if err != nil {
		return err
	}

	transactOpts.Value = new(big.Int).Mul(big.NewInt(210000), big.NewInt(configs.Gcc))
	_, err = rNodeContract.JoinRnode(transactOpts, big.NewInt(1))
	if err != nil {
		return err
	}

	backend.Commit()

	isRNode, _ := rNodeContract.IsRnode(nil, transactOpts.From)
	if !isRNode {
		return errors.New("The address is not RNode")
	}

	return nil
}

func TestDeployCampaign(t *testing.T) {
	contractBackend := backends.NewDposSimulatedBackend(core.GenesisAlloc{addr: {Balance: initBalance}})
	contractAddr, _, rNodeAddr, _, err := deploy(key, big.NewInt(0), contractBackend)
	checkError(t, "deploy contract: expected no error, got %v", err)

	err = fundToCampaign(key, rNodeAddr, contractBackend)
	checkError(t, "encounter error when fund to become campaign", err)

	transactOpts := bind.NewKeyedTransactor(key)
	campaign, err := contracts.NewCampaignWrapper(transactOpts, contractAddr, contractBackend)
	checkError(t, "can't deploy root registry: %v", err)
	_ = contractAddr
	contractBackend.Commit()

	// maximumNoc
	maximumNoc, err := campaign.MaximumNoc()
	checkError(t, "maximumNoc error: %v", err)
	t.Log("maximumNoc:", maximumNoc)

	// viewIdx
	viewIdx, err := campaign.TermIdx()
	checkError(t, "viewIdx error: %v", err)
	t.Log("viewIdx:", viewIdx)

	//minimumNoc
	minimumNoc, err := campaign.MinNoc()
	checkError(t, "minimumNoc error: %v", err)
	t.Log("minimumNoc:", minimumNoc)

	// test contract map variable call.
	numOfCampaign, startViewIdx, endViewIdx, err := campaign.CandidateInfoOf(addr)
	checkError(t, "CandidateInfoOf error: %v", err)
	t.Log("candidate info of", addr.Hex(), ":", numOfCampaign, startViewIdx, endViewIdx)

	verifyCandidates(campaign, t, big.NewInt(0), 0)
}

func checkError(t *testing.T, msg string, err error) {
	if err != nil {
		t.Fatalf(msg, err)
	}
}

func TestClaimCampaign(t *testing.T) {
	contractBackend := backends.NewDposSimulatedBackend(core.GenesisAlloc{addr: {Balance: initBalance}})

	campaignAddr, acAddr, rNodeAddr, networkAddr, err := deploy(key, big.NewInt(0), contractBackend)
	checkError(t, "deploy contract: expected no error, got %v", err)
	contractBackend.Commit()

	if err := fundToCampaign(key, rNodeAddr, contractBackend); err != nil {
		t.Error(err)
	}

	transactOpts := bind.NewKeyedTransactor(key)

	campaign, err := campaign.NewCampaign(campaignAddr, contractBackend)
	checkError(t, "can't deploy root registry: %v", err)
	_ = campaignAddr

	// compute cpu&memory pow
	ac := admission2.NewAdmissionControl(contractBackend.Blockchain(), addr, acAddr, campaignAddr, rNodeAddr, networkAddr)
	ac.SetSimulateBackend(contractBackend)
	configs.ChainConfigInfo().Dpos.Contracts[configs.ContractAdmission] = acAddr
	ac.Campaign(1)
	<-ac.DoneCh() // wait for done
	results := ac.GetResult()
	cpuBlockNum := results[admission2.Cpu].BlockNumber
	cpuNonce := results[admission2.Cpu].Nonce
	memBlockNum := results[admission2.Memory].BlockNumber
	memNonce := results[admission2.Memory].Nonce

	t.Log("cpuBlocksNum:", cpuBlockNum, "cpuNonce", cpuNonce, "memBlockNum", memBlockNum, "MemNonce", memNonce)
	ac1, _ := admission.NewAdmission(acAddr, contractBackend)
	result, _ := ac1.Verify(nil, cpuNonce, big.NewInt(cpuBlockNum), memNonce, big.NewInt(cpuBlockNum), transactOpts.From)
	if result == false {
		t.Error("Admission verify error")
	}

	transactOpts.GasLimit = uint64(50000000)
	// ClaimCampaign 1st time
	tx, err := campaign.ClaimCampaign(transactOpts, big.NewInt(1), cpuNonce, big.NewInt(cpuBlockNum), memNonce, big.NewInt(memBlockNum), version)
	checkError(t, "ClaimCampaign error:", err)
	t.Log("ClaimCampaign tx:", tx.Hash().Hex())
	contractBackend.Commit()

	receipt, _ := contractBackend.TransactionReceipt(context.Background(), tx.Hash())
	if receipt.Status == types.ReceiptStatusFailed {
		t.Fatal("the transaction should success but it failed!")
	}

	// verify result
	numOfCampaign, startViewIdx, endViewIdx, err := campaign.CandidateInfoOf(nil, addr)
	checkError(t, "CandidateInfoOf error: %v", err)
	t.Log("candidate info of", addr.Hex(), ":", numOfCampaign, startViewIdx, endViewIdx)
	assertCampaign(1, numOfCampaign, t)

	termIdx1, err := campaign.TermIdx(nil)
	if termIdx1.Uint64() != uint64(0) {
		t.Fatal("termIdx is not correct", "expected", 0, "actual", termIdx1)
	}

	// go forward to next term
	numPerRnd, _ := campaign.NumPerRound(nil)
	t.Log("numPerRnd", numPerRnd)
	for i := int64(0); i < numPerRnd.Int64(); i++ {
		contractBackend.Commit()
	}

	// ac again
	ac.Campaign(1)
	<-ac.DoneCh() // wait for done
	results = ac.GetResult()
	cpuBlockNum = results[admission2.Cpu].BlockNumber
	cpuNonce = results[admission2.Cpu].Nonce
	memBlockNum = results[admission2.Memory].BlockNumber
	memNonce = results[admission2.Memory].Nonce

	t.Log("cpuBlocksNum:", cpuBlockNum, "cpuNonce", cpuNonce, "memBlockNum", memBlockNum, "MemNonce", memNonce)
	result, _ = ac1.Verify(nil, cpuNonce, big.NewInt(cpuBlockNum), memNonce, big.NewInt(cpuBlockNum), transactOpts.From)
	if result == false {
		t.Error("Admission verify error")
	}

	tx, err = campaign.ClaimCampaign(transactOpts, big.NewInt(1), cpuNonce, big.NewInt(cpuBlockNum), memNonce, big.NewInt(memBlockNum), version)

	checkError(t, "ClaimCampaign error: %v", err)
	fmt.Println("ClaimCampaign tx:", tx.Hash().Hex())
	contractBackend.Commit()

	receipt, _ = contractBackend.TransactionReceipt(context.Background(), tx.Hash())
	if receipt.Status == types.ReceiptStatusFailed {
		t.Log("gas", receipt.GasUsed, ", cumulative gas", receipt.CumulativeGasUsed)
		t.Fatal("the transaction should success but it failed!")
	}

	// termIdx should be updated to current
	termIdx2, err := campaign.TermIdx(nil)
	if termIdx2.Uint64() != uint64(1) {
		t.Fatal("termIdx is not correct", "expected", 1, "actual", termIdx2)
	}

	// test contract map variable call.
	numOfCampaign, startViewIdx, endViewIdx, err = campaign.CandidateInfoOf(nil, addr)
	checkError(t, "CandidateInfoOf error: %v", err)
	fmt.Println("candidate info of", addr.Hex(), ":", numOfCampaign, startViewIdx, endViewIdx)
	// the second claim of campaign does not take effect as the previous campaign is not finished
	assertCampaign(2, numOfCampaign, t)

	// get candidates by view index
	candidates, err := campaign.CandidatesOf(nil, startViewIdx)
	checkError(t, "CandidatesOf error: %v", err)
	fmt.Println("len(candidates):", len(candidates))
	if len(candidates) != 1 {
		t.Fatal("len(candidates) != 1")
	}
	fmt.Println("candidates of first view:")
	for i := 0; i < len(candidates); i++ {
		fmt.Println("number", i, candidates[i].Hex())
	}
	printBalance(contractBackend)
}

func TestClaimCampaignOnSecondTerm(t *testing.T) {
	contractBackend := backends.NewDposSimulatedBackend(core.GenesisAlloc{addr: {Balance: initBalance}})
	printBalance(contractBackend)

	fmt.Println("deploy Campaign")
	campaignAddr, acAddr, rNodeAddr, networkAddr, err := deploy(key, big.NewInt(0), contractBackend)
	fmt.Println("campaignAddr:", campaignAddr)
	checkError(t, "deploy contract: expected no error, got %v", err)
	contractBackend.Commit()
	printBalance(contractBackend)

	fundToCampaign(key, rNodeAddr, contractBackend)

	fmt.Println("load Campaign")
	transactOpts := bind.NewKeyedTransactor(key)
	campaign, err := contracts.NewCampaignWrapper(transactOpts, campaignAddr, contractBackend)
	checkError(t, "can't deploy root registry: %v", err)
	_ = campaignAddr
	printBalance(contractBackend)

	// setup TransactOpts
	campaign.TransactOpts = *bind.NewKeyedTransactor(key)
	campaign.TransactOpts.Value = big.NewInt(40)
	campaign.TransactOpts.GasLimit = 1000000

	// compute cpu&memory pow
	ac := admission2.NewAdmissionControl(contractBackend.Blockchain(), addr, acAddr, campaignAddr, rNodeAddr, networkAddr)
	ac.SetSimulateBackend(contractBackend)
	configs.ChainConfigInfo().Dpos.Contracts[configs.ContractAdmission] = acAddr
	ac.Campaign(2)
	<-ac.DoneCh() // wait for done
	results := ac.GetResult()
	cpuBlockNum := results[admission2.Cpu].BlockNumber
	cpuNonce := results[admission2.Cpu].Nonce
	memBlockNum := results[admission2.Memory].BlockNumber
	memNonce := results[admission2.Memory].Nonce

	rNodeContract, err := rnode.NewRnode(rNodeAddr, contractBackend)
	isCan, _ := rNodeContract.IsRnode(&bind.CallOpts{From: transactOpts.From}, transactOpts.From)
	_ = isCan

	if !isCan {
		t.Error("This node is not Rnode!")
	}

	campaign.TransactOpts.Value = big.NewInt(0)
	tx, err := campaign.ClaimCampaign(big.NewInt(2), cpuNonce, big.NewInt(cpuBlockNum), memNonce, big.NewInt(memBlockNum), version)
	checkError(t, "ClaimCampaign error:", err)
	contractBackend.Commit()
	receipt, _ := contractBackend.TransactionReceipt(context.Background(), tx.Hash())
	if receipt.Status == types.ReceiptStatusFailed {
		t.Fatal("the transaction should success but it failed!")
	}

	// wait for view change
	waitForViewChange(contractBackend, 3)

	// get candidates by start view index
	verifyCandidates(campaign, t, big.NewInt(1), 1)
	verifyCandidates(campaign, t, big.NewInt(4), 0)
}

func TestClaimAndViewChangeThenQuitCampaign(t *testing.T) {
	contractBackend := backends.NewDposSimulatedBackend(core.GenesisAlloc{addr: {Balance: initBalance}})
	printBalance(contractBackend)

	fmt.Println("deploy Campaign")
	campaignAddr, acAddr, rNodeAddr, networkAddr, err := deploy(key, big.NewInt(0), contractBackend)
	fmt.Println("campaignAddr:", campaignAddr)
	checkError(t, "deploy contract: expected no error, got %v", err)
	contractBackend.Commit()
	printBalance(contractBackend)

	fundToCampaign(key, rNodeAddr, contractBackend)

	fmt.Println("load Campaign")
	transactOpts := bind.NewKeyedTransactor(key)
	campaign, err := contracts.NewCampaignWrapper(transactOpts, campaignAddr, contractBackend)
	checkError(t, "can't deploy root registry: %v", err)
	_ = campaignAddr
	printBalance(contractBackend)

	// setup TransactOpts
	campaign.TransactOpts = *bind.NewKeyedTransactor(key)
	campaign.TransactOpts.GasLimit = 1000000

	// compute cpu&memory pow
	ac := admission2.NewAdmissionControl(contractBackend.Blockchain(), addr, acAddr, campaignAddr, rNodeAddr, networkAddr)
	ac.SetSimulateBackend(contractBackend)
	configs.ChainConfigInfo().Dpos.Contracts[configs.ContractAdmission] = acAddr
	ac.Campaign(1)
	<-ac.DoneCh() // wait for done
	results := ac.GetResult()
	cpuBlockNum := results[admission2.Cpu].BlockNumber
	cpuNonce := results[admission2.Cpu].Nonce
	memBlockNum := results[admission2.Memory].BlockNumber
	memNonce := results[admission2.Memory].Nonce

	campaign.TransactOpts.Value = big.NewInt(0)
	tx, err := campaign.ClaimCampaign(big.NewInt(2), cpuNonce, big.NewInt(cpuBlockNum), memNonce, big.NewInt(memBlockNum), version)
	checkError(t, "ClaimCampaign error:", err)
	fmt.Println("ClaimCampaign tx:", tx.Hash().Hex())
	contractBackend.Commit()
	printBalance(contractBackend)

	// test contract map variable call.
	numOfCampaign, startViewIdx, endViewIdx, err := campaign.CandidateInfoOf(addr)
	checkError(t, "CandidateInfoOf error: %v", err)
	fmt.Println("candidate info of", addr.Hex(), ":", numOfCampaign, startViewIdx, endViewIdx)
	assertCampaign(2, numOfCampaign, t)

	// get candidates by view index
	verifyCandidates(campaign, t, startViewIdx, 1)
	contractBackend.Commit()

	// wait for view change
	waitForViewChange(contractBackend, 2)
}

func waitForViewChange(contractBackend *backends.SimulatedBackend, viewIdx int) {
	i := 0
	for i < viewIdx*numPerRound {
		contractBackend.Commit()
		i++
	}
}

func verifyCandidates(campaign *contracts.CampaignWrapper, t *testing.T, viewIdx *big.Int, candidateLengh int) {
	candidates, err := campaign.CandidatesOf(viewIdx)
	checkError(t, "CandidatesOf error: %v", err)
	fmt.Println("len(candidates):", len(candidates))
	if len(candidates) != candidateLengh {
		t.Fatal("len(candidates) != ", candidateLengh)
	}
}

func printBalance(contractBackend *backends.SimulatedBackend) {
	addrBalance, _ := contractBackend.BalanceAt(context.Background(), addr, nil)
	fmt.Println("==== addrBalance ==== ", addrBalance)
}

func assertCampaign(expectNum int64, numOfCampaign *big.Int, t *testing.T) {
	if numOfCampaign.Cmp(big.NewInt(expectNum)) != 0 {
		t.Fatal("unexpected numOfCampaign:", numOfCampaign)
	}
}
