package backend

import (
	"fmt"

	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/p2p"
)

const (
	// ProtocolName protocol name
	ProtocolName = "dpos"

	// ProtocolVersion protocol version
	ProtocolVersion = 65

	// ProtocolLength protocol length, max msg code
	ProtocolLength = 100
)

// Protocol messages belonging to gcc/01
const (
	// PbftMsgOutset is not a msg code, just used for msg code comparing
	PbftMsgOutset = 0x42

	// NewSignerMsg is a msg code used for network building
	NewSignerMsg = 0x42

	// those are messages for normal block verification
	PreprepareBlockMsg = 0x43
	PrepareHeaderMsg   = 0x44
	CommitHeaderMsg    = 0x45
	ValidateBlockMsg   = 0x46

	// those are messages for abnormal(impeachment) block verification
	PreprepareImpeachBlockMsg = 0x47
	PrepareImpeachHeaderMsg   = 0x48
	CommitImpeachHeaderMsg    = 0x49
	ValidateImpeachBlockMsg   = 0x50
)

// ProtocolMaxMsgSize Maximum cap on the size of a protocol message
const ProtocolMaxMsgSize = 10 * 1024 * 1024

type errCode int

const (
	// ErrMsgTooLarge is returned if msg if too large
	ErrMsgTooLarge = iota

	// ErrDecode is returned if decode failed
	ErrDecode

	// ErrInvalidMsgCode is returned if msg code is invalid
	ErrInvalidMsgCode

	// ErrProtocolVersionMismatch is returned if protocol version is not matched when handshaking
	ErrProtocolVersionMismatch

	// ErrNetworkIDMismatch is returned if networkid is not matched when handshaking
	ErrNetworkIDMismatch

	// ErrGenesisBlockMismatch is returned if genesis block is different from remote signer
	ErrGenesisBlockMismatch

	// ErrNoStatusMsg is returned if failed when reading status msg
	ErrNoStatusMsg

	// ErrExtraStatusMsg is returned if failed when extracting status msg
	ErrExtraStatusMsg

	// ErrSuspendedPeer is returned if remote signer is dead
	ErrSuspendedPeer
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

// XXX change once legacy code is out
var errorToString = map[int]string{
	ErrMsgTooLarge:             "Message too long",
	ErrDecode:                  "Invalid message",
	ErrInvalidMsgCode:          "Invalid message code",
	ErrProtocolVersionMismatch: "Protocol version mismatch",
	ErrNetworkIDMismatch:       "NetworkId mismatch",
	ErrGenesisBlockMismatch:    "Genesis block mismatch",
	ErrNoStatusMsg:             "No status message",
	ErrExtraStatusMsg:          "Extra status message",
	ErrSuspendedPeer:           "Suspended peer",
}

// SignerStatusData represents signer status when handshaking
type SignerStatusData struct {
	ProtocolVersion uint32
	Mac             string
	Sig             []byte
}

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

// IsSyncMsg checks if msg is a sync msg
func IsSyncMsg(msg p2p.Msg) bool {
	return msg.Code < PbftMsgOutset
}

// IsDposMsg checks if msg is a dpos related msg
func IsDposMsg(msg p2p.Msg) bool {
	return msg.Code >= PbftMsgOutset
}

// RecoverBlockFromMsg recovers a block from a p2p msg
func RecoverBlockFromMsg(msg p2p.Msg, p interface{}) (*types.Block, error) {
	// recover the block
	var block *types.Block
	if err := msg.Decode(&block); err != nil {
		return nil, errResp(ErrDecode, "%v: %v", msg, err)
	}
	block.ReceivedAt = msg.ReceivedAt
	block.ReceivedFrom = p

	return block, nil
}

// RecoverHeaderFromMsg recovers a header from a p2p msg
func RecoverHeaderFromMsg(msg p2p.Msg, p interface{}) (*types.Header, error) {
	// retrieve the header
	var header *types.Header
	if err := msg.Decode(&header); err != nil {
		return nil, errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	return header, nil
}
