package manager

import (
	"context"
	"testing"

	out "github.com/gcchains/chain/cmd/gcchain/campaign/output"
)

func TestManager(t *testing.T) {
	t.Skip()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	endPoint := "http://127.0.0.1:8503"
	kspath := "/Users/steve/.gcchain/keystore/"
	password := "/Users/steve/.gcchain/password"
	output := out.NewLogOutput()

	manager, _ := NewConsole(&ctx, endPoint, kspath, password, &output)

	manager.GetStatus()
}
