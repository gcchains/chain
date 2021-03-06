

package core

import (
	"fmt"
	"math/big"

	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus/dpos"
	"github.com/gcchains/chain/core/vm"
	"github.com/gcchains/chain/database"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func ExampleGenerateChain() {
	var (
		key1, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _  = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		key3, _  = crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
		addr1    = crypto.PubkeyToAddress(key1.PublicKey)
		addr2    = crypto.PubkeyToAddress(key2.PublicKey)
		addr3    = crypto.PubkeyToAddress(key3.PublicKey)
		db       = database.NewMemDatabase()
		remoteDB = database.NewIpfsDbWithAdapter(database.NewFakeIpfsAdapter())
	)

	type configBackup struct {
		Cep1LastBlockY1 *big.Int
		Cep1LastBlockY2 *big.Int
		Cep1LastBlockY3 *big.Int
		Cep1LastBlockY4 *big.Int
		Cep1LastBlockY5 *big.Int
	}
	bak := configBackup{
		Cep1LastBlockY1: configs.Cep1LastBlockY1(),
		Cep1LastBlockY2: configs.Cep1LastBlockY2(),
		Cep1LastBlockY3: configs.Cep1LastBlockY3(),
		Cep1LastBlockY4: configs.Cep1LastBlockY4(),
		Cep1LastBlockY5: configs.Cep1LastBlockY5(),
	}

	configs.TestOnly_SetCep1LastBlockY1(big.NewInt(1))
	configs.TestOnly_SetCep1LastBlockY2(big.NewInt(2))
	configs.TestOnly_SetCep1LastBlockY3(big.NewInt(3))
	configs.TestOnly_SetCep1LastBlockY4(big.NewInt(4))
	configs.TestOnly_SetCep1LastBlockY5(big.NewInt(5))

	// recover configs.Cep1LastBlockYx
	defer func() {
		configs.TestOnly_SetCep1LastBlockY1(bak.Cep1LastBlockY1)
		configs.TestOnly_SetCep1LastBlockY2(bak.Cep1LastBlockY2)
		configs.TestOnly_SetCep1LastBlockY3(bak.Cep1LastBlockY3)
		configs.TestOnly_SetCep1LastBlockY4(bak.Cep1LastBlockY4)
		configs.TestOnly_SetCep1LastBlockY5(bak.Cep1LastBlockY5)

	}()

	// Ensure that key1 has some funds in the genesis block.
	gspec := DefaultGenesisBlock()
	gspec.Alloc = GenesisAlloc{addr1: {Balance: big.NewInt(1000000)}}
	// Config: &configs.ChainConfig{HomesteadBlock: new(big.Int)},
	genesis := gspec.MustCommit(db)

	engine := dpos.NewFaker(configs.ChainConfigInfo().Dpos, db)
	// This call generates a chain of 5 blocks. The function runs for
	// each block and adds different features to gen based on the
	// block index.
	signer := types.HomesteadSigner{}
	n := 5
	chain, _ := GenerateChain(gspec.Config, genesis, engine, db, remoteDB, n, func(i int, gen *BlockGen) {
		switch i {
		case 0:
			gen.SetCoinbase(addr1)
			// In block 1, addr1 sends addr2 some ether.
			tx, _ := types.SignTx(types.NewTransaction(gen.TxNonce(addr1), addr2, big.NewInt(10000), configs.TxGas, nil, nil), signer, key1)
			gen.AddTx(tx)
		case 1:
			gen.SetCoinbase(addr2)
			// In block 2, addr1 sends some more ether to addr2.
			// addr2 passes it on to addr3.
			tx1, _ := types.SignTx(types.NewTransaction(gen.TxNonce(addr1), addr2, big.NewInt(1000), configs.TxGas, nil, nil), signer, key1)
			tx2, _ := types.SignTx(types.NewTransaction(gen.TxNonce(addr2), addr3, big.NewInt(1000), configs.TxGas, nil, nil), signer, key2)
			gen.AddTx(tx1)
			gen.AddTx(tx2)
		case 2:
			// Block 3 is empty but was *mined* by addr3.
			gen.SetCoinbase(addr3)
		}
	})

	// Import the chain. This runs all block validation rules.
	blockchain, _ := NewBlockChain(db, nil, gspec.Config, engine, vm.Config{}, remoteDB, nil)
	defer blockchain.Stop()

	if i, err := blockchain.InsertChain(chain); err != nil {
		fmt.Printf("insert error (block %d): %v\n", chain[i].NumberU64(), err)
		return
	}

	state, _ := blockchain.State()
	fmt.Printf("last block: #%d\n", blockchain.CurrentBlock().Number())
	fmt.Println("balance of addr1:", new(big.Int).Sub(state.GetBalance(addr1), new(big.Int).Mul(big.NewInt(1265), big.NewInt(1e+16))))
	fmt.Println("balance of addr2:", new(big.Int).Sub(state.GetBalance(addr2), new(big.Int).Mul(big.NewInt(951), big.NewInt(1e+16))))

	sub := new(big.Int).Mul(big.NewInt(713+539+403), big.NewInt(1e+16))
	balanceAddr3 := new(big.Int).Sub(state.GetBalance(addr3), sub)
	fmt.Println("balance of addr3 (adjusted):", balanceAddr3)
	// Output:
	// last block: #5
	// balance of addr1: 989000
	// balance of addr2: 10000
	// balance of addr3 (adjusted): 1000
}
