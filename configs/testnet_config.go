


package configs

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)


var (
	// contract
	testnetProxyContractRegister = common.HexToAddress("0x36a8ac0cad2150e036de638aa492042eeb823c6b")

	TestnetContractAddressMap = map[string]common.Address{
		ContractAdmission: common.HexToAddress("0x82102c2A09DEe47D1DDcf398AeF7877F99d787c4"),
		ContractCampaign:  common.HexToAddress("0x3B11cA41A28571e22e342299dC308f08EDD7F011"),
		ContractRpt:       common.HexToAddress("0x7a174062c3C8551649A86AE1b8a84282D901C2C3"),
		ContractRnode:     common.HexToAddress("0xF0f87e064C76674fE7c4dDceE3603AFC37998658"),
	}

	// config
	testnetDefaultCandidates = []common.Address{
		common.HexToAddress("0x3a15146f433c0205cfae639de2ac3bb543539b24"), // #1
		common.HexToAddress("0xb436e2feff376c30beb9d89e825281baa3956d4c"), // #2
		common.HexToAddress("0xf675b1e939892cad974a17da6bcd13ceae3b73dd"), // #3
		common.HexToAddress("0x37a99234187e95f28f8f69d44437fb16c465071c"), // #4
		common.HexToAddress("0x3326d5248928b83f63a80e424a1c6d39cb334624"), // #5
		common.HexToAddress("0x2661177788fe63888e93cf38b5e4e31306301170"), // #6
	}
	testnetChainConfig = &ChainConfig{
		ChainID: big.NewInt(TestnetChainId),
		Dpos: &DposConfig{
			Period:                TestnetBlockPeriod,
			TermLen:               4,
			ViewLen:               3,
			FaultyNumber:          TestnetFaultyValidatorsNumber,
			MaxInitBlockNumber:    DefaultTestnetMaxInitBlockNumber,
			ProxyContractRegister: testnetProxyContractRegister,
			Contracts:             TestnetContractAddressMap,
			ImpeachTimeout:        time.Millisecond * TestnetBlockPeriod * 2,
		},
	}

	testnetProposers = testnetDefaultCandidates[0:4]

	testnetValidators = []common.Address{
		common.HexToAddress("0x177b2a835f27a3989dfca814b37d08c54e1de883"),
		common.HexToAddress("0x332062f84f982050c820b5ec98631825d003ec8e"),
		common.HexToAddress("0x2da372d6326573aa5e1863ba3fa724a231c477d3"),
		common.HexToAddress("0x38e86c815665de503a210ff4b8e8572b8c201309"),
	}

	testnetBootnodes = []string{
		"enode://18c444f813e3fbef9848748306a4a4b2fa2d90090a31e59c1dcdfa55a7435a18abaabffcd205fa976a2f4f9b1832ffd361b1e53bcef6b052823dd442b1722bf8@127.0.0.1:31000",
		"enode://9eedb4aa96949a2db1307a5360e604f5149a2933b82f70c7ac3080362db170a17513de101e39d36634994a22003c1f77980699d72636d9f747ece888e0c98395@127.0.0.1:31000",
	}

	
	defaultTestnetValidatorNodes = []string{
		
		"enode://3f11492af45df3c06fbdc43534a66615baa58dc58a4918a3b693bf67c5baad4098ea5e03a63a26ed53890865b8aa30550727ebff32b6823b72ad5c9dd28b4313@127.0.0.1:30017",
	
		"enode://f22094e4153d73d304d0e362704ecfec5fa928e56b62703b599a66e445c7bfa3b7a525118dc7e41fdf98263e428bda4d98cb3405f50ae539add6fc5aa87c98b9@127.0.0.1:30018",
		
		"enode://33925c2c99a2bc8ebb05d0946ee673d18fb1e2905b3e1e7ea4c840dd6cac5fc369ac54d1c791b158dbba3734494422fb01103ac4384f932d214aba69e0b43518@127.0.0.1:30019",
		
		"enode://d4175e2c796a6591e52e788483261bb54cfc337e0ba881f033cafd1333ea94d22f3c84fa8652a343b2cb155b6443d3494a9010a1b993e63841374e9311382513@127.0.0.1:30020",
	}
)
