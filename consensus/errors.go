

package consensus

import (
	"errors"
)

var (
	// ErrUnknownAncestor is returned when validating a block requires an ancestor
	// that is unknown.
	ErrUnknownAncestor = errors.New("unknown ancestor")

	// ErrPrunedAncestor is returned when validating a block requires an ancestor
	// that is known, but the state of which is not available.
	ErrPrunedAncestor = errors.New("pruned ancestor")

	// ErrFutureBlock is returned when a block's timestamp is in the future according
	// to the current node.
	ErrFutureBlock = errors.New("block in the future")

	// ErrInvalidTimestamp is returned when a block's timestamp is larger than parent's
	// timestamp + period + timeout.
	ErrInvalidTimestamp = errors.New("invalid timestamp")

	// ErrInvalidNumber is returned if a block's number doesn't equal it's parent's
	// plus one.
	ErrInvalidNumber = errors.New("invalid block number")

	// ErrNotEnoughSigs is returned if there is not enough signatures for a block.
	ErrNotEnoughSigs = errors.New("not enough signatures in block")

	// ErrUnauthorized is returned if a header is signed by a non-authorized entity.
	ErrUnauthorized = errors.New("unauthorized leader")

	// ErrNotInProposerCommittee is returned  if the account is not in proposer committee.
	ErrNotInProposerCommittee = errors.New("not in proposer committee")

	// ErrUnknownLbftState is returned if committee handler's state is unknown
	ErrUnknownLbftState = errors.New("unknown lbft state")

	// ErrInvalidSigners is returned if a block contains an invalid extra sigers bytes.
	ErrInvalidSigners = errors.New("invalid signer list on checkpoint block")

	// ErrInvalidNormalCoinbase is returned if a normal block's coinbase is 0x00.
	ErrInvalidNormalCoinbase = errors.New("invalid normal coinbase, it's 0x00")

	// --- those are invalid impeach block errors ---

	// ErrInvalidImpeachCoinbase is returned if an impeach block's coinbase is not 0x00.
	ErrInvalidImpeachCoinbase = errors.New("invalid impeach coinbase")

	// ErrInvalidImpeachStateRoot is returned if an impeach block's StateRoot is not equal to parents'.
	ErrInvalidImpeachStateRoot = errors.New("invalid impeach state root")

	// ErrInvalidImpeachTxsRoot is returned if an impeach block's TxsRoot is not valid.
	ErrInvalidImpeachTxsRoot = errors.New("invalid impeach txs root")

	// ErrInvalidImpeachReceiptsRoot is returned if an impeach block's ReceiptsRoot is not valid.
	ErrInvalidImpeachReceiptsRoot = errors.New("invalid impeach receipts root")

	// ErrInvalidImpeachLogsBloom is returned if an impeach block's LogsBloom is not valid.
	ErrInvalidImpeachLogsBloom = errors.New("invalid impeach LogsBloom")

	// ErrInvalidImpeachGasLimit is returned if an impeach block's GasLimit is not valid.
	ErrInvalidImpeachGasLimit = errors.New("invalid impeach GasLimit")

	// ErrInvalidImpeachGasUsed is returned if an impeach block's GasUsed is not valid.
	ErrInvalidImpeachGasUsed = errors.New("invalid impeach GasUsed")

	// ErrInvalidImpeachTimestamp is returned if an impeach block's Timestamp is not valid.
	ErrInvalidImpeachTimestamp = errors.New("invalid impeach Timestamp")

	// ErrInvalidImpeachExtra is returned if an impeach block's Extra is not valid.
	ErrInvalidImpeachExtra = errors.New("invalid impeach Extra")

	// ErrInvalidImpeachDposSnap is returned if an impeach block's DposSnap is not valid.
	ErrInvalidImpeachDposSnap = errors.New("invalid impeach DposSnap")

	// ErrInvalidImpeachTxs is returned if an impeach block contrains txs
	ErrInvalidImpeachTxs = errors.New("invalid impeach txs")

	// --- those are invalid impeach block errors ---

	// ErrorInvalidValidatorsList is returned if the validators list is invalid
	ErrorInvalidValidatorsList = errors.New("invalid validators list")
)
