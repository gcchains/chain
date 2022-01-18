

package gcc

import (
	"math/big"
	"os"
	"os/user"
	"time"

	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/core"
	"github.com/gcchains/chain/private"
	"github.com/gcchains/chain/protocols/gcc/gasprice"
	"github.com/gcchains/chain/protocols/gcc/syncer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// DefaultConfig contains default settings.
var DefaultConfig = Config{
	NetworkId:     configs.MainnetNetworkId,
	DatabaseCache: 768,
	TrieCache:     256,
	TrieTimeout:   60 * time.Minute,
	GasPrice:      big.NewInt(18 * configs.Shannon),

	TxPool: core.DefaultTxPoolConfig,
	GPO: gasprice.Config{
		Blocks:     20,
		Percentile: 60,
	},
	PrivateTx: private.DefaultConfig(),
	SyncMode:  syncer.FullSync,
}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	// The genesis block, which is inserted if the database is empty.
	// If nil, the gcchain main net block is used.
	Genesis *core.Genesis `toml:",omitempty"`

	// Protocol options
	NetworkId uint64 // Network ID to use for selecting peers to connect to
	NoPruning bool   // TODO: remove it {AC}

	// Database options
	SkipBcVersionCheck bool `toml:"-"`
	DatabaseHandles    int  `toml:"-"`
	DatabaseCache      int
	TrieCache          int
	TrieTimeout        time.Duration

	// Mining-related options
	Gccbase      common.Address `toml:",omitempty"`
	MinerThreads int            `toml:",omitempty"`
	ExtraData    []byte         `toml:",omitempty"`
	GasPrice     *big.Int

	// Transaction pool options
	TxPool core.TxPoolConfig

	// Gas Price Oracle options
	GPO gasprice.Config

	// Enables tracking of SHA3 preimages in the VM
	EnablePreimageRecording bool

	// Miscellaneous options
	DocRoot string `toml:"-"`

	// Private Tx related configuration
	PrivateTx private.Config

	SyncMode syncer.SyncMode `toml:"-"`
}

type configMarshaling struct {
	ExtraData hexutil.Bytes
}
