

package flags

import (
	"github.com/gcchains/chain/commons/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

var flagMap = make(map[string]cli.Flag)

func init() {
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version, v",
		Usage: "Print the version",
	}

	cli.HelpFlag = cli.BoolFlag{
		Name:  "help, h",
		Usage: "Show help",
	}

	Register(GeneralFlags...)
}

func Register(flags ...cli.Flag) {
	for _, flag := range flags {
		if _, ok := flagMap[flag.GetName()]; ok {
			log.Fatalf("Flag already exists: %v", flag.GetName())
		}
		flagMap[flag.GetName()] = flag
	}
}

func GetByName(name string) cli.Flag {
	flag, ok := flagMap[name]
	if !ok {
		log.Fatalf("Flag does not exist: %v", name)
	}
	return flag
}

// begin flags

const (
	KeystorePath = "keystore"
	Endpoint     = "endpoint"
	ContractAddr = "contractaddr"
)

var GeneralFlags = []cli.Flag{
	cli.StringFlag{
		Name:  KeystorePath,
		Usage: "Keystore file path for contract admin",
	},
	cli.StringFlag{
		Name:  Endpoint,
		Usage: "Endpoint to interact with",
	},
	cli.StringFlag{
		Name:  ContractAddr,
		Usage: "Contract address",
	},
}

func GetContractAddress(ctx *cli.Context) common.Address {
	if !ctx.IsSet(ContractAddr) {
		log.Fatal("contract address must be provided!")
	}

	contractAddr := common.HexToAddress(ctx.String(ContractAddr))
	return contractAddr
}

func GetEndpoint(ctx *cli.Context) string {
	if !ctx.IsSet(Endpoint) {
		log.Fatal("endpoint must be provided!")
	}

	endpoint := ctx.String(Endpoint)
	return endpoint
}

func GetKeystorePath(ctx *cli.Context) string {
	if !ctx.IsSet(KeystorePath) {
		log.Fatal("keystore path must be provided!")
	}

	keystorePath := ctx.String(KeystorePath)
	return keystorePath
}
