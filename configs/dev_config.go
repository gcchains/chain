


package configs

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// dev configuration
var (
	// contract
	devProxyContractRegister = common.HexToAddress("0xd4826921aa2dba7930117712ed183576ccebed91")
	DevContractAddressMap    = map[string]common.Address{
		ContractRpt:       common.HexToAddress("0x7e9915bea4af2ebea96dd8ba9814d4503e6c0218"),
		ContractRnode:     common.HexToAddress("0xd4826921aa2dba7930117782ed183576ccebed13"),
		ContractAdmission: common.HexToAddress("0xa1e0ea2a14d91031986c2f25f6e124beeeb66781"),
		ContractCampaign:  common.HexToAddress("0xf26b6164749cde15a29afea57ffeae115b24b505"),
		ContractNetwork:   common.HexToAddress("0xe30361ffce7f560cb69b7d9254daeb35de1c0f84"),
	}

	// config
	devDefaultCandidates = []common.Address{
		common.HexToAddress("0xc05302acebd0730e3a18a051d7d1cb1204c4a092"), // #2
		common.HexToAddress("0xe94b7b6c5a0e526a4d17f9768ad6097bde25c62a"), // #1
		common.HexToAddress("0xef3dd127de235f11ffb4fc0d71469d1339df6465"), // #3
		common.HexToAddress("0x6e31e5b68a98dcd17264bd1ba547d0b3e814da1e"), // #5
		common.HexToAddress("0x3a18598184ef84198db10c28fdfdfdf56544f747"), // #4
		common.HexToAddress("0x22a672eab2b1a3ff3ed91563201a56ca5a560e08"), // #6
	}
	devChainConfig = &ChainConfig{
		ChainID: big.NewInt(DevChainId),
		Dpos: &DposConfig{
			Period:                DefaultBlockPeriod,
			TermLen:               4,
			ViewLen:               3,
			FaultyNumber:          DefaultFaultyValidatorsNumber,
			MaxInitBlockNumber:    DefaultDevMaxInitBlockNumber,
			ProxyContractRegister: devProxyContractRegister,
			Contracts:             DevContractAddressMap,
			ImpeachTimeout:        time.Millisecond * DefaultBlockPeriod,
		},
	}

	devProposers = devDefaultCandidates[0:4]

	devValidators = []common.Address{
		common.HexToAddress("0x7b2f052a372951d02798853e31ee56c895109992"),
		common.HexToAddress("0x2f0176cc3a8617b1ddea6a501028fa4c6fc25ca1"),
		common.HexToAddress("0xe4d51117832e84f1d082e9fc12439b171a57e7b2"),
		common.HexToAddress("0x32bd7c33bb1060a85f361caf20c0bda9075c5d51"),
	}

	// gcchainBootnodes are the enode URLs of the P2P bootstrap nodes running on
	// the dev gcchain network.
	devBootnodes = []string{
		"enode://5293dc8aaa5c2fcc7905c21391ce38f4f877722ff1918f4fa86371347ad8a244c2995631f89866693d05bf5c94493c247f02116f19a90689fa406189b03a5243@127.0.0.1:30381", // localhost
	}

	defaultDevValidatorNodes = []string{
		"enode://9826a2f72c63eaca9b7f57b169473686f5a133dc24ffac858b4e1185a5eb60b144a414c35359585d9ea9d67f6fcca29578f9e002c89e94cc4bcc46a2b336c166@127.0.0.1:30317",
		"enode://7ce9c4fee12b12affbbe769a0faaa6e256bbae3374717fb94e1fb4be308fae3795c3abae023a587d1e14b35d278bd3d10916117bb8b3f5cfa4c951c5d56eeed7@127.0.0.1:30318",
		"enode://1db32421dc881357c282091960fdbd13f3635f8e3f87a953b6d9c429e53469727018bd0bb02da18acc4f1b4bec946b8f158705262b37163b4ab321a1c932d8f9@127.0.0.1:30319",
		"enode://fd0f365cec4e052040151f2a4a9ba23e8592acd3cacfdc4af2e8b6dbc6fb6b25ca088151889b19721d02c48e390de9682b316db2351636fdd1ee5ea1cd32bf46@127.0.0.1:30320",
	}
)
