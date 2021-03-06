

package core

import (
	"runtime"
	"testing"
	"time"

	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus/dpos"
	"github.com/gcchains/chain/core/vm"
	"github.com/gcchains/chain/database"
	"github.com/gcchains/chain/types"
)

// Tests that simple header verification works, for both good and bad blocks.
func TestHeaderVerification(t *testing.T) {
	testdb := database.NewMemDatabase()
	genesis := DefaultGenesisBlock()
	genesisBlock := genesis.MustCommit(testdb)

	remoteDB := database.NewIpfsDbWithAdapter(database.NewFakeIpfsAdapter())

	config := configs.ChainConfigInfo().Dpos
	d := dpos.NewFaker(config, testdb)

	blocks, _ := GenerateChain(genesis.Config, genesisBlock, d, testdb, remoteDB, 8, nil)

	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	// Run the header checker for blocks one-by-one, checking for both valid and invalid nonces
	chain, _ := NewBlockChain(testdb, nil, genesis.Config, d, vm.Config{}, remoteDB, nil)
	defer chain.Stop()

	for i := 0; i < len(blocks); i++ {
		for j, valid := range []bool{true, false} {
			var results <-chan error

			if valid {
				engine := dpos.NewFaker(config, testdb)
				_, results = engine.VerifyHeaders(chain, []*types.Header{headers[i]}, []bool{true}, []*types.Header{headers[i]})
			} else {
				// engine := dpos.New(config, testdb)
				engine := dpos.NewFakeFailer(config, testdb, headers[i].Number.Uint64())
				_, results = engine.VerifyHeaders(chain, []*types.Header{headers[i]}, []bool{true}, []*types.Header{headers[i]})
			}
			// Wait for the verification result
			select {
			case result := <-results:
				if (result == nil) != valid {
					t.Errorf("test %d.%d: validity mismatch: have %v, want %v", i, j, result, valid)
				}
			case <-time.After(time.Second):
				t.Fatalf("test %d.%d: verification timeout", i, j)
			}
			// Make sure no more data is returned
			select {
			case result := <-results:
				t.Fatalf("test %d.%d: unexpected result returned: %v", i, j, result)
			case <-time.After(25 * time.Millisecond):
			}
		}
		chain.InsertChain(blocks[i : i+1])
	}
}

// Tests that concurrent header verification works, for both good and bad blocks.
func TestHeaderConcurrentVerification2(t *testing.T) {
	testHeaderConcurrentVerification(t, 2)
}
func TestHeaderConcurrentVerification8(t *testing.T) {
	testHeaderConcurrentVerification(t, 8)
}
func TestHeaderConcurrentVerification32(t *testing.T) {
	testHeaderConcurrentVerification(t, 32)
}

func testHeaderConcurrentVerification(t *testing.T, threads int) {
	// Create a simple chain to verify
	var (
		testdb = database.NewMemDatabase()
	)

	genesis := DefaultGenesisBlock()
	genesisBlock := genesis.MustCommit(testdb)

	remoteDB := database.NewIpfsDbWithAdapter(database.NewFakeIpfsAdapter())

	config := configs.ChainConfigInfo().Dpos
	d := dpos.NewDoNothingFaker(config, testdb)

	blocks, _ := GenerateChain(genesis.Config, genesisBlock, d, testdb, remoteDB, 8, nil)

	headers := make([]*types.Header, len(blocks))
	seals := make([]bool, len(blocks))

	for i, block := range blocks {
		headers[i] = block.Header()
		seals[i] = true
	}
	// Set the number of threads to verify on
	old := runtime.GOMAXPROCS(threads)
	defer runtime.GOMAXPROCS(old)

	// Run the header checker for the entire block chain at once both for a valid and
	// also an invalid chain (enough if one arbitrary block is invalid).
	for i, valid := range []bool{true, false} {
		var results <-chan error

		if valid {
			engine := dpos.NewDoNothingFaker(config, testdb)
			chain, _ := NewBlockChain(testdb, nil, genesis.Config, engine, vm.Config{}, remoteDB, nil)
			_, results = chain.engine.VerifyHeaders(chain, headers, seals, headers)
			chain.Stop()
		} else {
			engine := dpos.NewFakeFailer(config, testdb, uint64(len(headers)-1))
			chain, _ := NewBlockChain(testdb, nil, genesis.Config, engine, vm.Config{}, remoteDB, nil)
			_, results = chain.engine.VerifyHeaders(chain, headers, seals, headers)
			chain.Stop()
		}
		// Wait for all the verification results
		checks := make(map[int]error)
		for j := 0; j < len(blocks); j++ {
			select {
			case result := <-results:
				checks[j] = result

			case <-time.After(time.Second):
				t.Fatalf("test %d.%d: verification timeout", i, j)
			}
		}
		// Check nonce check validity
		for j := 0; j < len(blocks); j++ {
			want := valid || (j < len(blocks)-2) // We chose the last-but-one nonce in the chain to fail
			if (checks[j] == nil) != want {
				t.Errorf("test %d.%d: validity mismatch: have %v, want %v", i, j, checks[j], want)
			}
			if !want {
				// A few blocks after the first error may pass verification due to concurrent
				// workers. We don't care about those in this test, just that the correct block
				// errors out.
				break
			}
		}
		// Make sure no more data is returned
		select {
		case result := <-results:
			t.Fatalf("test %d: unexpected result returned: %v", i, result)
		case <-time.After(25 * time.Millisecond):
		}
	}
}

// Tests that aborting a header validation indeed prevents further checks from being
// run, as well as checks that no left-over goroutines are leaked.
func TestHeaderConcurrentAbortion2(t *testing.T)  { testHeaderConcurrentAbortion(t, 2) }
func TestHeaderConcurrentAbortion8(t *testing.T)  { testHeaderConcurrentAbortion(t, 8) }
func TestHeaderConcurrentAbortion32(t *testing.T) { testHeaderConcurrentAbortion(t, 32) }

func testHeaderConcurrentAbortion(t *testing.T, threads int) {
	// Create a simple chain to verify

	var (
		testdb = database.NewMemDatabase()
	)

	genesis := DefaultGenesisBlock()
	genesisBlock := genesis.MustCommit(testdb)

	remoteDB := database.NewIpfsDbWithAdapter(database.NewFakeIpfsAdapter())

	config := configs.ChainConfigInfo().Dpos
	d := dpos.NewFaker(config, testdb)

	blocks, _ := GenerateChain(genesis.Config, genesisBlock, d, testdb, remoteDB, 8, nil)

	headers := make([]*types.Header, len(blocks))
	seals := make([]bool, len(blocks))

	for i, block := range blocks {
		headers[i] = block.Header()
		seals[i] = true
	}
	// Set the number of threads to verify on
	old := runtime.GOMAXPROCS(threads)
	defer runtime.GOMAXPROCS(old)

	engine := dpos.NewFakeDelayer(config, testdb, time.Millisecond)
	// Start the verifications and immediately abort
	chain, _ := NewBlockChain(testdb, nil, genesis.Config, engine, vm.Config{}, remoteDB, nil)
	defer chain.Stop()

	abort, results := chain.engine.VerifyHeaders(chain, headers, seals, headers)
	close(abort)

	// Deplete the results channel
	verified := 0
	for depleted := false; !depleted; {
		select {
		case result := <-results:
			if result != nil {
				t.Errorf("header %d: validation failed: %v", verified, result)
			}
			verified++
		case <-time.After(50 * time.Millisecond):
			depleted = true
		}
	}
	// Check that abortion was honored by not processing too many POWs
	if verified > 3*threads {
		t.Errorf("verification count too large: have %d, want below %d", verified, 3*threads)
	}
}
