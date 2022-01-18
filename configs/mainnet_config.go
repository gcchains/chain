


package configs

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Mainnet configuration
var (
	// contract
	mainnetProxyContractRegister = common.HexToAddress("0xd4826927aa2dba7230117782ed183526ccebed93")

	MainnetContractAddressMap = map[string]common.Address{
		ContractRpt:       common.HexToAddress("0x16cb35DD47421295215b01a41a8e224E6eb39235"),
		ContractRnode:     common.HexToAddress("0x26130DA5aA1851313a7225D3735BED76029560DA"),
		ContractAdmission: common.HexToAddress("0xB3178aa5f6B5ABDc534e5bDEEc70B7e36BBDa4e2"),
		ContractCampaign:  common.HexToAddress("0x1A186bE66Dd20c1699Add34A49A3019a93a7Fcd1"),
		ContractNetwork:   common.HexToAddress("0xFE4e9816C4B05D0be1fe1fb951FfAB44e3309118"),
	}

	// config
	mainnetDefaultCandidates = []common.Address{
		common.HexToAddress("0x50f8c76f6d1442c54905c74145ae163132b9f4ae"), // #1
		common.HexToAddress("0x8ab63651e1ce7eed40af33276011a5e2e1a133a2"), // #2
		common.HexToAddress("0x501f6cf7b2417671d770998e3b785474871fef1d"), // #3
		common.HexToAddress("0x9508e430ce672150bcf6bef9c4c0adf303b21c5c"), // #4
		common.HexToAddress("0x049295e2e125cec28ddeeb63507e654b6d117423"), // #5
		common.HexToAddress("0x8c65cb8561c4945d4b419af9066874118f19ba43"), // #6
		common.HexToAddress("0x1d1f1d14f303b746121564e3295e2957e74ea1d2"), // #14
		common.HexToAddress("0x13ae2b32ef3fad80707d4da0f49c675d3efc717c"), // #15
		common.HexToAddress("0x5a51d5ef67c047b5d748724f7401781ed0af61ed"), // #16
		common.HexToAddress("0x1f071085dfdfa4a65f8870e50348f277d6fcd91c"), // #17
		common.HexToAddress("0xcb6fb6a101d6c126f80053fe17ca41188e24fe2f"), // #18
		common.HexToAddress("0xfaf2a1cdc4da310b52ad7d8d16e8c1bd5d4c0bd0"), // #19
	}
	mainnetChainConfig = &ChainConfig{
		ChainID: big.NewInt(MainnetChainId),
		Dpos: &DposConfig{
			Period:                MainnetBlockPeriod,
			TermLen:               12,
			ViewLen:               3,
			FaultyNumber:          MainnetFaultyValidatorsNumber,
			MaxInitBlockNumber:    DefaultMainnetMaxInitBlockNumber,
			ProxyContractRegister: mainnetProxyContractRegister,
			Contracts:             MainnetContractAddressMap,
			ImpeachTimeout:        time.Millisecond * MainnetBlockPeriod,
		},
	}
	mainnetProposers = mainnetDefaultCandidates

	mainnetValidators = []common.Address{
		common.HexToAddress("0x2effd798190059ed313fa7d01483bcfa7ea637be"), //#7
		common.HexToAddress("0x1f11dc7132c31dd26dfac3754b3a7ea0da1ea351"), //#8
		common.HexToAddress("0x1a9463aa4b1157681421c72e052d2cf8b6498a38"), //#9
		common.HexToAddress("0xd08975bb1c17c8139cf5107e1fd896d46b7a841a"), //#10
		common.HexToAddress("0xce1eb5797e457cf199f0bf2c994c46614bfb4feb"), // #11
		common.HexToAddress("0x15b5a0709ae1cf751377bf9e26b8340c7bb9112b"), // #12
		common.HexToAddress("0x1e693ea09b1593bd3c715186806121f5b69371b6"), // #13
	}

	mainnetBootnodes = []string{
		"enode://fac8239134ee0ae6a1d9a0035fd83157011ac27e6ed0c8b00bdd422e7f100b8333660255aac4a21c54df9548fdbeaa81e6216338d41809a9218118dd62e08764@g1.mainnet.gc-server.com:30010",
		"enode://c009b708a911bc610204acf46cea61358ec2113a2cf9e62e93b18139a40763af9853238b3162c62d07aeea0d8d5c6d91b08b88747a0bf4f737bb9a1230b3561b@g2.mainnet.gc-server.com:30010",
		"enode://2258c9581e453fa0d4f75529824811b8db7e168611cda737f24b0c3f44ae1242d924cc600c1f78fc7d077c80e9a3f124a47f4fa221327e741aae3a515493e059@g3.mainnet.gc-server.com:30010",
	}

	defaultMainnetValidatorNodes = []string{
		"enode://503dd5cd4a532635f484c3ea9abb906a3a286cc79112f670595f021da9a1570f1e984e9d3e8374084d8b6d136798ba771a01a76f80a6acd67ff824864ba67e0a@s1.mainnet.gc-server.com:30010",
		"enode://f64d08fac6acd07be96df3bcfa3f73128a128a374ccdb1d88b565c43662ab7148ab85548f89b319b5b964491e3e6f677a468c03111398c653183c310c3c22de4@s2.mainnet.gc-server.com:30010",
		"enode://196f9d87497b3e5de234eb6b35cdbafbf9851b006bac1246eaf71868464306fdeb971b5a182185ece04e747cc3e219376c72eb8df30b5d74918f7a7105ff6aaf@s3.mainnet.gc-server.com:30010",
		"enode://9d570a1366a15b45597a48f0ea8e0c7d2a5519fd23106d53629d69156176a58dbdfca57bdfb4dd205f02aaf81b6f02bbb70699adff6d56849ce140ac694021bf@s4.mainnet.gc-server.com:30010",
		"enode://6bc0cadbf3f60876a9b7da3f6f5bad406d484fffc414f993a1dfa6535c19b457c30287cd7e4ac43980efafd52499ebac151752b2f889f18d3241086a32f7f678@s5.mainnet.gc-server.com:30010",
		"enode://80945a6705944d157515a83258b4dedf939cf42ed5cdc240081ea4587c2227841e20b4d078c5923e01c3d8cb746a2569e8ac25933ab9f4f54c6e2a17af318005@s6.mainnet.gc-server.com:30010",
		"enode://01e620180eb2c0fc3d15e0de246fd897781a72ee1f7431a05dc082b2a3f0488eb297ce8d59e4c10d362e02152d300bbc3f7cd42306e8c7deab62617be1debdd4@s7.mainnet.gc-server.com:30010",
	}
)
