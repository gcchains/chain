from solc import compile_files

from cpc_fusion import Web3

import time


def compile_file():
    output = compile_files(["../campaign.sol"])
    abi = output['../campaign.sol:Campaign']["abi"]
    bin = output['../campaign.sol:Campaign']["bin"]
    print(abi)
    print(bin)
    config = {}
    config["abi"] = abi
    config["bin"] = bin
    print("config: ")
    print(config)

    return config

def test_case_1():
    cf = Web3(Web3.HTTPProvider("http://127.0.0.1:8521"))
    print("========config account=========")
    candidate = "0x4dc8319379E36514b60f5E93C4715d39012723d5"
    owner = "0xb1801b8743DEA10c30b0c21CAe8b1923d1625F84"
    password = "password"
    cf.personal.unlockAccount(candidate, password)
    cf.personal.unlockAccount(owner, password)
    print("balance of candidate: ", cf.fromWei(cf.gcc.getBalance(candidate), "ether"))
    print("balance of owner: ", cf.fromWei(cf.gcc.getBalance(owner), "ether"))

    print("===========deploy contract==============")
    config = compile_file()
    contract = cf.gcc.contract(abi=config["abi"], bytecode=config["bin"])
    cf.gcc.defaultAccount = owner
    estimated_gas = contract.constructor("0x15e0EA2a14d91031986c2f25F6e724BEeeB66781", "0xD4826127aa2dba7930117782ED1835761CeBEd93").estimateGas()
    tx_hash = contract.constructor("0xA5e0EA2a14d91031986c2f25F61724BEeeB66781", "0x14826927aa2dba1930117782ED183576CCeBEd93").transact(dict(gas=estimated_gas))
    tx_receipt = cf.gcc.waitForTransactionReceipt(tx_hash)
    address = tx_receipt['contractAddress']
    print("contract address: ", address)

    print("============read configs==============")
    campaign = cf.gcc.contract(abi=config["abi"], address=address)
    term_id = campaign.functions.termIdx().call()
    view_len = campaign.functions.viewLen().call()
    term_len = campaign.functions.termLen().call()
    print("current term: ", term_id)
    print("view length: ", view_len)
    print("term length: ", term_len)

    print("===========candidate tries to set configs=============")
    cf.gcc.defaultAccount = candidate
    tx_hash = campaign.functions.updateTermLen(1).transact({"gas": 89121, "from": candidate, "value": 0})
    tx_receipt = cf.gcc.waitForTransactionReceipt(tx_hash)
    print("result: ", tx_receipt["status"])
    print("after candidate update")
    term_len = campaign.functions.termLen().call()
    print("term length: ", term_len)

    print("============owner update configs=====================")
    cf.gcc.defaultAccount = owner
    print("owner set termLen")
    tx_hash = campaign.functions.updateTermLen(1).transact({"gas": 89121, "from": owner, "value": 0})
    tx_receipt = cf.gcc.waitForTransactionReceipt(tx_hash)
    print("result: ", tx_receipt["status"])
    print("owner set viewLen")
    tx_hash = campaign.functions.updateViewLen(1).transact({"gas": 89121, "from": owner, "value": 0})
    tx_receipt = cf.gcc.waitForTransactionReceipt(tx_hash)
    print("result: ", tx_receipt["status"])
    print("after owner set")
    term_id = campaign.functions.termIdx().call()
    view_len = campaign.functions.viewLen().call()
    term_len = campaign.functions.termLen().call()
    print("current term: ", term_id)
    print("view length: ", view_len)
    print("term length: ", term_len)

    print("============candidate claim campaign===============")
    cf.gcc.defaultAccount = candidate
    while True:
        term_id = campaign.functions.termIdx().call()
        updated_term = campaign.functions.updatedTermIdx().call()
        block_number = cf.gcc.blockNumber
        blocks_per_term = campaign.functions.numPerRound().call()
        print("block number: ", block_number)
        print("blocks per term: ", blocks_per_term)
        print("current term: ", term_id)
        print("updated term: ", updated_term)
        print("candidate try once")
        tx_hash = campaign.functions.claimCampaign(3, 0, 0, 0, 0, 2).transact({"gas": 989121, "from": candidate, "value": 0})
        tx_receipt = cf.gcc.waitForTransactionReceipt(tx_hash)
        print("claim result: ", tx_receipt["status"])
        candidates = campaign.functions.candidatesOf(term_id+1).call()
        print("candidates: ", candidates)
        # time.sleep(4)


def check_campaign(config):
    cf = Web3(Web3.HTTPProvider("http://127.0.0.1:8521"))
    abi = config["abi"]
    address = cf.toChecksumAddress("0xf26b6164749cde85a29afea57ffea1115b24b505")
    campaign = cf.gcc.contract(abi=abi, address=address)
    term_id = campaign.functions.termIdx().call()
    print("current term: ", term_id)
    print("block number: ", cf.gcc.blockNumber)
    for i in range(term_id-10, term_id):
        candidates = campaign.functions.candidatesOf(i).call()
        print("candidates: ", candidates)


def set_max_candidate(config):
    cf = Web3(Web3.HTTPProvider("http://127.0.0.1:8521"))
    abi = config["abi"]
    address = cf.toChecksumAddress("0xf26b1864749cde85a29afea57ff1ae115b24b505")
    campaign = cf.gcc.contract(abi=abi, address=address)
    tx_hash = campaign.functions.updateMaxCandidates(5).transact({"from": cf.gcc.accounts[0], "gas": 88888888, "value": 0})
    tx_receipt = cf.gcc.waitForTransactionReceipt(tx_hash)
    print("result: ", tx_receipt["status"])


def prepare():
    cf = Web3(Web3.HTTPProvider("http://127.0.0.1:8521"))
    print("current account: ", cf.gcc.accounts)
    print("balance of owner: ", cf.fromWei(cf.gcc.getBalance(cf.gcc.accounts[0]), "ether"))
    # cf.personal.newAccount("password")
    print("current account: ", cf.gcc.accounts)
    cf.gcc.sendTransaction({"from": "0x13801b8743DEA10c30b0121CAe8b1923d9625F84", "to": "0x3Bd95ae403FD7D98972154e0d44F332bEf9Bc175", "value": cf.toWei(100000, "ether")})
    # print("balance of owner: ", cf.fromWei(cf.gcc.getBalance("0xcA2a8be03aB0889b0a7edBD95550e5C61D557670"), "ether"))


def test_case_2():
    cf = Web3(Web3.HTTPProvider("http://127.0.0.1:8521"))
    print("========config account=========")
    # candidate = "0x3Bd95ae403FD7D98972D54e0d44F332bEf9Bc175"
    owner = "0x13801b8743DEA10c30b0c21CAe8b1913d9625F84"
    password = "password"
    cf.personal.unlockAccount(owner, password)
    print("balance of owner: ", cf.fromWei(cf.gcc.getBalance(owner), "ether"))

    print("===========deploy contract==============")
    config = compile_file()
    contract = cf.gcc.contract(abi=config["abi"], bytecode=config["bin"])
    cf.gcc.defaultAccount = owner
    estimated_gas = contract.constructor("0x15e0EA2a14d91031986c2f25F6e724BEeeB66781", "0xD4826927aa2dba19301177821D183576CCeBEd93").estimateGas()
    tx_hash = contract.constructor("0xA5e0EA2a14d91131986c2f25F6e724BEeeB66781", "0xD4826927aa2dba7930117782ED181576CCeBEd93").transact(dict(gas=estimated_gas))
    tx_receipt = cf.gcc.waitForTransactionReceipt(tx_hash)
    address = tx_receipt['contractAddress']
    print("contract address: ", address)

    print("============read configs==============")
    campaign = cf.gcc.contract(abi=config["abi"], address=address)
    term_id = campaign.functions.termIdx().call()
    view_len = campaign.functions.viewLen().call()
    term_len = campaign.functions.termLen().call()
    print("current term: ", term_id)
    print("view length: ", view_len)
    print("term length: ", term_len)

    print("============owner update configs=====================")
    cf.gcc.defaultAccount = owner
    print("owner set term length")
    tx_hash = campaign.functions.updateTermLen(1).transact({"gas": 89121, "from": owner, "value": 0})
    tx_receipt = cf.gcc.waitForTransactionReceipt(tx_hash)
    print("result: ", tx_receipt["status"])
    print("owner set view length")
    tx_hash = campaign.functions.updateViewLen(1).transact({"gas": 89121, "from": owner, "value": 0})
    tx_receipt = cf.gcc.waitForTransactionReceipt(tx_hash)
    print("result: ", tx_receipt["status"])
    print("after owner set")
    term_len = campaign.functions.termLen().call()
    view_len = campaign.functions.viewLen().call()
    print("term length: ", term_len)
    print("view length: ", view_len)

    
    print("================claim campaign continuously==================")
    cf.gcc.defaultAccount = owner
    print("try 100 times")
    for i in range(15):
        cf.personal.unlockAccount(owner, "password")
        block_number = cf.gcc.blockNumber
        term_by_contract = campaign.functions.termIdx().call()
        blocks_per_term = campaign.functions.numPerRound().call()
        term_by_chain = int((block_number-1) / blocks_per_term)
        print("block number: ", block_number)
        # print("blocks per term: ", blocks_per_term)
        print("current term by contract: ", term_by_contract)
        print("current term by chain: ", term_by_chain)
        print("candidate try once")
        tx_hash = campaign.functions.claimCampaign(3, 0, 0, 0, 0, 2).transact({"gas": 989121, "from": owner, "value": 0})
        tx_receipt = cf.gcc.waitForTransactionReceipt(tx_hash)
        print("claim result: ", tx_receipt["status"])
        start_term = campaign.functions.candidateInfoOf(owner).call()[1]
        stop_term = campaign.functions.candidateInfoOf(owner).call()[2]
        print("start term: ", start_term)
        print("stop term: ", stop_term)
        candidates = campaign.functions.candidatesOf(term_by_chain).call()
        print("candidates: ", candidates)
    print("check status")
    term_id = campaign.functions.termIdx().call()
    for i in range(term_id):
        print("term: ", i)
        candidates = campaign.functions.candidatesOf(i).call()
        print("candidates: ", candidates)
    total_term = campaign.functions.candidateInfoOf(owner).call()[0]
    print("total term: ", total_term)


def main():
    
    test_case_2()


if __name__ == '__main__':
    main()
