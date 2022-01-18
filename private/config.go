

package private

const (
	DefaultIpfsUrl = "127.0.0.1:5001"
	Dummy          = "dummy"
	IPFS           = "ipfs"
	Swarm          = "swarm"
)

var (
	// Indicate whether support private transaction functionality
	SupportPrivateTx string
)

type Config struct {
	RemoteDBParams string
	RemoteDBType   string
}

func DefaultConfig() Config {
	return Config{
		RemoteDBType:   Dummy,
		RemoteDBParams: DefaultIpfsUrl,
	}
}
