

package bind_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/gcchains/chain/accounts/abi/bind"
	"github.com/gcchains/chain/accounts/abi/bind/backends"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/core"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad3e901695e1b4b9ee17ae26c6668d313eac2f96dbcda3f291")

var waitDeployedTests = map[string]struct {
	code        string
	gas         uint64
	wantAddress common.Address
	wantErr     error
}{
	"successful deploy": {
		code:        `3060304052600a8060106000396000f360606040526008565b00`,
		gas:         3000000,
		wantAddress: common.HexToAddress("0x3a220f351252049d385b29beca14e27f204c296a"),
	},

	// TODO: fix this.
	// "empty code": {
	// 	code:        ``,
	// 	gas:         300000,
	// 	wantErr:     bind.ErrNoCodeAfterDeploy,
	// 	wantAddress: common.HexToAddress("0x3a220f351252049d385b29beca14e27f204c296a"),
	// },
}

func TestWaitDeployed(t *testing.T) {
	for name, test := range waitDeployedTests {
		backend := backends.NewDposSimulatedBackend(core.GenesisAlloc{
			crypto.PubkeyToAddress(testKey.PublicKey): {Balance: big.NewInt(10000000000)},
		})

		// Create the transaction.
		tx := types.NewContractCreation(0, big.NewInt(0), test.gas, big.NewInt(1), common.FromHex(test.code))
		tx, _ = types.SignTx(tx, types.NewCep1Signer(big.NewInt(configs.MainnetChainId)), testKey)

		// Wait for it to get mined in the background.
		var (
			err     error
			address common.Address
			mined   = make(chan struct{})
			ctx     = context.Background()
		)
		go func() {
			address, err = bind.WaitDeployed(ctx, backend, tx)
			close(mined)
		}()

		// Send and mine the transaction.
		backend.SendTransaction(ctx, tx)
		backend.Commit()

		select {
		case <-mined:
			if err != test.wantErr {
				t.Errorf("Test %q: error mismatch: got %q, want %q", name, err, test.wantErr)
			}
			if address != test.wantAddress {
				t.Errorf("Test %q: unexpected contract address %s", name, address.Hex())
			}
		case <-time.After(2 * time.Second):
			t.Errorf("Test %q: timeout", name)
		}
	}
}
