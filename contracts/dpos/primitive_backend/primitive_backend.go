package primitive_backend

import (
	"context"
	"math/big"
	"sync"

	gcchain "/gcchain/chain"
	"github.com/gcchains/chain/api/rpc"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/core/state"
	"github.com/gcchains/chain/internal/gccapi"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// used to hold *APIBackend
type BackendHolder struct {
	ChainBackend    ChainAPIBackend
	ContractBackend ContractAPIbcakend
}

var apiBackendHolderInstance *BackendHolder
var onceApiBackendHoldCreation sync.Once
var onceApiBackendHoldInit sync.Once

func GetApiBackendHolderInstance() *BackendHolder {
	onceApiBackendHoldCreation.Do(func() {
		apiBackendHolderInstance = &BackendHolder{}
	})
	return apiBackendHolderInstance
}

func (rb *BackendHolder) Init(chainBackend ChainAPIBackend, contractBackend ContractAPIbcakend) {
	onceApiBackendHoldInit.Do(func() {
		log.Debug("init BackendHolder", "ChainBackend", chainBackend, "ContractBackend", contractBackend)
		rb.ChainBackend = chainBackend
		rb.ContractBackend = contractBackend
	})
}

func GetChainClient() *ApiClient {
	return &ApiClient{ChainBackend: GetApiBackendHolderInstance().ChainBackend, ContractBackend: GetApiBackendHolderInstance().ContractBackend}
}

type ChainClient struct {
	ChainAPIBackend
}

type ChainAPIBackend interface {
	StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber, isPrivate bool) (*state.StateDB, *types.Header, error)
	BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error)
	HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error)
}
type ContractAPIbcakend interface {
	Call(ctx context.Context, args gccapi.CallArgs, blockNr rpc.BlockNumber) (hexutil.Bytes, error)
}
type ApiClient struct {
	ChainBackend    ChainAPIBackend
	ContractBackend ContractAPIbcakend
}

// BalanceAt returns the wei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (cc *ApiClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	state, _, err := cc.ChainBackend.StateAndHeaderByNumber(ctx, rpc.BlockNumber(blockNumber.Uint64()), false)
	if state == nil || err != nil {
		return nil, err
	}
	return state.GetBalance(account), state.Error()
}

func (cc *ApiClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	state, _, err := cc.ChainBackend.StateAndHeaderByNumber(ctx, rpc.BlockNumber(blockNumber.Uint64()), false)
	if state == nil || err != nil {
		return 0, err
	}
	return state.GetNonce(account), state.Error()
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
func (cc *ApiClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	return cc.ChainBackend.BlockByNumber(ctx, rpc.BlockNumber(number.Uint64()))
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (cc *ApiClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	return cc.ChainBackend.HeaderByNumber(ctx, rpc.BlockNumber(number.Uint64()))
}

func (cc *ApiClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	blockNr := rpc.LatestBlockNumber
	state, _, err := cc.ChainBackend.StateAndHeaderByNumber(ctx, blockNr, false)
	if state == nil || err != nil {
		return nil, err
	}
	code := state.GetCode(account)
	return code, state.Error()
}

func (cc *ApiClient) CallContract(ctx context.Context, call gcchain.CallMsg, blockNumber *big.Int) ([]byte, error) {
	result, err := cc.ContractBackend.Call(ctx, toCallArg(call), rpc.LatestBlockNumber)
	if err != nil {
		log.Fatal("CallContract using PublicBlockChainAPI is error ", "error is ", err)
	}
	return result, err
}
func toCallArg(msg gcchain.CallMsg) gccapi.CallArgs {
	arg := gccapi.CallArgs{
		From: msg.From,
		To:   msg.To,
	}
	if len(msg.Data) > 0 {
		arg.Data = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg.Value = hexutil.Big(*msg.Value)
	}
	if msg.Gas != 0 {
		arg.Gas = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg.GasPrice = hexutil.Big(*msg.GasPrice)
	}
	return arg
}
func (cc *ApiClient) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	blockNr := rpc.PendingBlockNumber
	state, _, err := cc.ChainBackend.StateAndHeaderByNumber(ctx, blockNr, false)
	if state == nil || err != nil {
		return nil, err
	}
	code := state.GetCode(contract)
	return code, state.Error()
}

func (cc *ApiClient) PendingCallContract(ctx context.Context, call gcchain.CallMsg) ([]byte, error) {
	result, err := cc.ContractBackend.Call(ctx, toCallArg(call), rpc.PendingBlockNumber)
	if err != nil {
		log.Fatal("CallContract using PublicBlockChainAPI is error ", "error is ", err)
	}
	return result, err
}

func (cc *ApiClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	panic("that is fake PendingNonceAt please using RPC to call real function")
}

func (cc *ApiClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	panic("call the fake SuggestGasPrice,please using RPC to call real SuggestGasPrice")
}

func (cc *ApiClient) EstimateGas(ctx context.Context, call gcchain.CallMsg) (gas uint64, err error) {
	panic("that is fake PendingNonceAt please using RPC to call real function")
}
func (cc *ApiClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	panic("that is fake PendingNonceAt please using RPC to call real function")
}

func (cc *ApiClient) FilterLogs(ctx context.Context, query gcchain.FilterQuery) ([]types.Log, error) {
	panic("this is a fake FilterLogs,please use RPC call the real function")
}

func (cc *ApiClient) SubscribeFilterLogs(ctx context.Context, q gcchain.FilterQuery, ch chan<- types.Log) (gcchain.Subscription, error) {
	panic("this is a fake SubscribeFilterLogs,please use RPC call the real function")
}
