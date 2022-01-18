package backend

import (
	"testing"

	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
)

func TestHandler_BroadcastMinedBlock(t *testing.T) {
	type fields struct {
		mode           HandlerMode
		config         *configs.DposConfig
		available      bool
		coinbase       common.Address
		dialer         *Dialer
		snap           *consensus.PbftStatus
		dpos           DposService
		knownBlocks    *RecentBlocks
		pendingBlockCh chan *types.Block
		quitSync       chan struct{}
	}
	type args struct {
		block *types.Block
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				mode:           tt.fields.mode,
				config:         tt.fields.config,
				available:      tt.fields.available,
				coinbase:       tt.fields.coinbase,
				dialer:         tt.fields.dialer,
				dpos:           tt.fields.dpos,
				knownBlocks:    tt.fields.knownBlocks,
				pendingBlockCh: tt.fields.pendingBlockCh,
				quitCh:         tt.fields.quitSync,
			}
			h.BroadcastPreprepareBlock(tt.args.block)
		})
	}
}
