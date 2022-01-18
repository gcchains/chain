package output

import (
	"testing"

	status "github.com/gcchains/chain/cmd/gcchain/campaign/common"
)

func TestLogOutput(t *testing.T) {
	output := NewLogOutput()
	output.Info("Info")
	output.Error("Error")
	output.Warn("Warn")

	// Status
	output.Status(&status.Status{
		Mining:   true,
		RNode:    true,
		Proposer: true,
	})
	output.Status(&status.Status{
		Mining:   true,
		RNode:    true,
		Proposer: false,
	})
}
