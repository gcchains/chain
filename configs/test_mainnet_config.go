


package configs

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// testMainnet configuration
var (
	// contract
	testMainnetProxyContractRegister = common.HexToAddress("0xd4826921aa2dba7930117782ed113576ccebed93")

	TestMainnetContractAddressMap = map[string]common.Address{
		ContractRpt:       common.HexToAddress("0x7e9925bea4af2ebea16dd8ba9894d4503e1c0278"),
		ContractRnode:     common.HexToAddress("0x14826921aa2dba7930117782ed183576ccebed93"),
		ContractAdmission: common.HexToAddress("0xa5e0ea2a14d91031986c2f25f1e724beeeb61781"),
		ContractCampaign:  common.HexToAddress("0x126b6864749cde85a29afea57ffeae115b21b505"),
		ContractNetwork:   common.HexToAddress("0xe30361ffce7f560cb69b7d9254daeb35de8c1f84"),
	}

	// config
	testMainnetDefaultCandidates = []common.Address{
		common.HexToAddress("0x1e61732d0b1c1674151a01ac0bba824c5b1258fb"), // #1
		common.HexToAddress("0xaa6cf4f0338e04a40719dfa3c653efc1cd9e65c9"), // #2
		common.HexToAddress("0x7170f518ca82897375f009ddea319df08f31bcff"), // #3
		common.HexToAddress("0x1c61559aa727380e3fa516b6a7ae397b87e12384"), // #5
		common.HexToAddress("0xc5b481311bbcabb96ed0c835cee69b471419f49c"), // #4
		common.HexToAddress("0x6e7fdba0f15067a25a3cf1df90429e3c919411e3"), // #6

		common.HexToAddress("0x17e81a296f5b80d319d2f3018f2d5998530e79e4"), // #14
		common.HexToAddress("0x52e584b4fba8188eb7edcabb18e65161a99acc67"), // #15
		common.HexToAddress("0x030352bba31c0c7cec8669f64a26d96d5d179bdb"), // #16
		common.HexToAddress("0x1561ebb8a40114c1cf3cc0a628df5a1bd1663b26"), // #17
		common.HexToAddress("0xca8e011de1edea4929328bb86e35daa681c47ed0"), // #18
		common.HexToAddress("0xcc9cd266771b331fd424ea14dc30fc8511bec628"), // #19
	}
	testMainnetChainConfig = &ChainConfig{
		ChainID: big.NewInt(TestMainnetChainId),
		Dpos: &DposConfig{
			Period:                TestMainnetBlockPeriod,
			TermLen:               12,
			ViewLen:               3,
			FaultyNumber:          TestMainnetFaultyValidatorsNumber,
			MaxInitBlockNumber:    DefaultTestMainnetMaxInitBlockNumber,
			ProxyContractRegister: testMainnetProxyContractRegister,
			Contracts:             TestMainnetContractAddressMap,
			ImpeachTimeout:        time.Millisecond * TestMainnetBlockPeriod,
		},
	}
	testMainnetProposers = testMainnetDefaultCandidates

	testMainnetValidators = []common.Address{
		common.HexToAddress("0x0b2ee61452cc151565ed4b8eabff85c3f185c149"),
		common.HexToAddress("0x6a3678cac50b9266f82a1e1a12bd26edc1e743a3"),
		common.HexToAddress("0xc6bfd405a91a39fa06f3cf0f568c3a1a40c29882"),
		common.HexToAddress("0xaee4ec17edd59f5a2a0fe1fc786d217bea1ac3d9"),
		common.HexToAddress("0x17be125f3c60105b44e3242f5c5509d6c993ebb8"), // #11
		common.HexToAddress("0x30a36525ca16504939e944e81422bdac745dd050"), // #12
		common.HexToAddress("0x1341844d109c918f70d1ff4e621bc8da097b8d83"), // #13
	}

	testMainnetBootnodes = []string{
		"enode://d9bd60488c269f1324ed7811341f4c81ce41ed0531852ea265148a617f4cb99d58c95a979de059c5052b1d38eb4462b715f8d8db40f92cf3828d884265176cde@g1.mainnet.gc-server.com:30001",
		"enode://775cedbf2026c065b67fc80a38995c2999d5b3c113a80695115b2606ad4025dacae6947034189032453c0a71d8b49465f161e9c3b0e85a8cf9b4281e1cc198e4@g2.mainnet.gc-server.com:30002",
		"enode://121a1567059168e3cecc1d5c60217cd73dc5b299c3865f2f9eea621a93e7e8bb266f7ade0b1db28160686c740be6de4c1000ef441d7cd8c97388eca1790bc61a@g3.mainnet.gc-server.com:30003",
		"enode://fb56640fb3b8dec3473ea3906ef59b97c4f7911d86be27ed65908fa706d2fbe91800b7a221ba45127cbe4b5eb26291f7fbd3984cd1ab587d6fb53535ce4e0069@g4.mainnet.gc-server.com:30004",
	}

	defaultTestMainnetValidatorNodes = []string{
		"enode://1ddfb534019e6b446fb4465742f266d04fae662089e3dac6a4c841ad0fcf5529f8d049203079bb64e20d1a32fc84b584920839a2120cd5e8886744719452d936@s1.mainnet.gc-server.com:30017",
		"enode://f2a460e5d5008d0ba8ec0744c90df9d1fc01553e00025af483995a25d89e773de18a37972c38bdcf47917fc820738455b85675bb21b026a75768c68d5540d053@s2.mainnet.gc-server.com:30018",
		"enode://f3045792b9e9ad894cb36b488f1cf97065013cda9ef60f2b14840214683b32f3dadf61450a9f677457c1a4b75a5b456947f48f32b0019c7470cced9af1829993@s3.mainnet.gc-server.com:30019",
		"enode://1e14fce25a846bd5c91728fec7fb7071c92e2b9f8f4b710dbce79d8b6098592591ebeebfe6c59ee5bfd6f75387926f9342ae004d6ff8d2f97fc6d7e91e8f41be@s4.mainnet.gc-server.com:30020",
		"enode://00e5229f3792264032a335259671996da3714f90f8d19defd0abce4e2751527e644a76ae19b994f9b28b4d652826fa0766298260db6df70aa2def7461c50d662@s5.mainnet.gc-server.com:30021",
		"enode://269699f91013336e4ecf329aac4a4a6ee3957c7c7577996f9db821013e2e232ef8151e200cc2ab7ea9265121642b05b1cd21640d29e1e4228f6af737f353275c@s6.mainnet.gc-server.com:30022",
		"enode://ee4c7418336745ed8a54da5fd8b151ade53b0b2a23b8e1d5eecfae483d15f5ff9e440155c47311dc226c44d44dce0080a6246204ed992f1e37d7094df4289169@s7.mainnet.gc-server.com:30023",
	}
)
