

package tests

import (
	"math/big"
	"testing"

	"github.com/gcchains/chain/configs"
)

func TestTransaction(t *testing.T) {
	t.Parallel()

	txt := new(testMatcher)
	txt.config(`^Homestead/`, configs.ChainConfig{})
	txt.config(`^EIP155/`, configs.ChainConfig{
		ChainID: big.NewInt(1),
	})
	txt.config(`^Byzantium/`, configs.ChainConfig{})

	txt.walk(t, transactionTestDir, func(t *testing.T, name string, test *TransactionTest) {
		cfg := txt.findConfig(name)
		if err := txt.checkFailure(t, name, test.Run(cfg)); err != nil {
			t.Error(err)
		}
	})
}
