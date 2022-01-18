

package tests

import (
	"fmt"
	"math/big"

	"github.com/gcchains/chain/configs"
)

// Forks table defines supported forks and their chain config.
var Forks = map[string]*configs.ChainConfig{
	// TODO: @AC confirm the real name of the initial phase(the first release of gcchain)
	"cep1": {
		ChainID: big.NewInt(45),
	},
}

// UnsupportedForkError is returned when a test requests a fork that isn't implemented.
type UnsupportedForkError struct {
	Name string
}

func (e UnsupportedForkError) Error() string {
	return fmt.Sprintf("unsupported fork %q", e.Name)
}
