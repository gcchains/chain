

package core

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/gcchains/chain/configs"
	"github.com/naoina/toml"
)

func TestDefaultGenesisBlock_MarshalTOML(t *testing.T) {
	genesisblock := DefaultGenesisBlock()
	fmt.Println("==============toml=====================")
	err := toml.NewEncoder(os.Stdout).Encode(genesisblock)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("==============json=====================")
	ss, err := json.Marshal(genesisblock)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("genesisblock", string(ss))
}

func TestGenesisAccount_MarshalTOML(t *testing.T) {
	genesisblock := GenesisAccount{
		Nonce:      100,
		Code:       []byte("hhhhh"),
		PrivateKey: []byte("hhhh2222h"),
	}
	fmt.Println("===================================")
	err := toml.NewEncoder(os.Stdout).Encode(genesisblock)
	if err != nil {
		t.Error(err)
	}
}

func TestGenesisAccount_MarshalJson(t *testing.T) {
	genesisblock := GenesisAccount{
		Nonce:      100,
		Code:       []byte("hhhhh"),
		PrivateKey: []byte("hhhh2222h"),
	}
	fmt.Println("===================================")

	ss, err := json.Marshal(genesisblock)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("ss", string(ss))
}

func TestTestnetGenesisBlock_MarshalTOML(t *testing.T) {
	origMode := configs.GetRunMode()
	configs.SetRunMode(configs.Testnet)
	genesisblock := DefaultGenesisBlock()
	fmt.Println("==============toml=====================")
	err := toml.NewEncoder(os.Stdout).Encode(genesisblock)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("==============json=====================")
	ss, err := json.Marshal(genesisblock)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("genesisblock", string(ss))
	configs.SetRunMode(origMode)
}

func TestMainnetGenesisBlock_MarshalTOML(t *testing.T) {
	origMode := configs.GetRunMode()
	configs.SetRunMode(configs.Mainnet)
	genesisblock := DefaultGenesisBlock()
	fmt.Println("==============toml=====================")
	err := toml.NewEncoder(os.Stdout).Encode(genesisblock)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("==============json=====================")
	ss, err := json.Marshal(genesisblock)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("genesisblock", string(ss))
	configs.SetRunMode(origMode)
}
