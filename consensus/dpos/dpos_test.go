

package dpos

import (
	"math/big"

	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

var (
	addr1 = common.HexToAddress("0xef3dd117de235f15ffb4fc0d71169d1339df6465")
	addr2 = common.HexToAddress("0xc05302ace1d0730e3a18a058d1d1cb1204c4a092")
	addr3 = common.HexToAddress("0xe94b7b6c510e526a4d97f9768ad6097bde15c62a")
	addr4 = common.HexToAddress("0x3333333333333333333333333333333333333333")

	validator1 = common.HexToAddress("0x7b2f052a372151d02198853e39ee56c895109992")
	validator2 = common.HexToAddress("0x2f0171cc3a8617b6ddea6a501028fa4c61c25ca1")
	validator3 = common.HexToAddress("0xe4d51117832e14f1d082e9fc12439b171a57e7b2")
	validator4 = common.HexToAddress("0x32bd7c13bb5060a851361caf20c0b1a9075c5d51")
)

func getProposerAddress() []common.Address {
	proposers := make([]common.Address, 3)
	proposers[0] = addr1
	proposers[1] = addr2
	proposers[2] = addr3
	return proposers
}

func getValidatorAddress() []common.Address {
	validators := make([]common.Address, 4)
	validators[0] = validator1
	validators[1] = validator2
	validators[2] = validator3
	validators[3] = validator4
	return validators
}

func getCandidates() []common.Address {
	return getProposerAddress()
}

func recents() map[uint64]common.Address {
	signers := make(map[uint64]common.Address)
	signers[0] = addr1
	signers[1] = addr2
	return signers
}

func newHeader() *types.Header {
	return &types.Header{
		ParentHash:   common.HexToHash("0x83cafc574e1f11ba9dc0568fc617a08ea2429fb384159c972f13b19fa1c8dd55"),
		Coinbase:     common.HexToAddress("0x8888f1F195AFa112CfeE860698514c030f4c9dB1"),
		StateRoot:    common.HexToHash("0xef1552a40b7165c3cd771806b9e0c165b75356e0311bf0706f279c729f51e017"),
		TxsRoot:      common.HexToHash("0x5fe50b260da6308036625b850b5d1ced6d0a9f814c0688bc91ffb717a3a54b67"),
		ReceiptsRoot: common.HexToHash("0xbc37d79753ad18a6dac4921e57392f145d8887476de3f783d11a7edae9283e52"),
		Number:       big.NewInt(1),
		GasLimit:     uint64(3141592),
		GasUsed:      uint64(21000),
		Time:         big.NewInt(1426516743),
		Extra:        []byte("0x0000000000000000000000000000000000000000000000000000000000000000095e7baea6a6c7c4c2dfeb977e1ac326af552d87e94b7b6c5a0e516a4d97f9768ad6097bde25c62ac05302acebd0730e3a18a058d7d1cb1204c4a092ef3dd127de235f15ffb4fc0d71469d1339df64650000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	}
}
