

// Package gcc implements the gcchain protocol.
package gcc

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/gcchains/chain/accounts"
	"github.com/gcchains/chain/admission"
	"github.com/gcchains/chain/api/rpc"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus"
	"github.com/gcchains/chain/consensus/dpos"
	"github.com/gcchains/chain/contracts/dpos/primitive_backend"
	"github.com/gcchains/chain/core"
	"github.com/gcchains/chain/core/bloombits"
	"github.com/gcchains/chain/core/rawdb"
	"github.com/gcchains/chain/core/vm"
	"github.com/gcchains/chain/database"
	"github.com/gcchains/chain/internal/gccapi"
	"github.com/gcchains/chain/miner"
	"github.com/gcchains/chain/node"
	"github.com/gcchains/chain/private"
	"github.com/gcchains/chain/protocols/gcc/filters"
	"github.com/gcchains/chain/protocols/gcc/gasprice"
	"github.com/gcchains/chain/protocols/gcc/syncer"
	"github.com/gcchains/chain/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/p2p"
)

var (
	errForbidValidatorMining = errors.New("Validator is forbidden to mine.")
	errNotAdmissionKey       = errors.New("Admission key is missing, need to run with --mine flag.")
)

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

// gcchainService implements the gcchainService full node service.
type gcchainService struct {
	config      *Config
	chainConfig *configs.ChainConfig

	// Channel for shutting down the service
	shutdownChan chan bool // Channel for shutting down the gcchain

	// Handlers
	txPool          *core.TxPool
	blockchain      *core.BlockChain
	protocolManager *ProtocolManager
	lesServer       LesServer

	server *p2p.Server

	// DB interfaces
	chainDb database.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // LogsBloom indexer operating during block imports

	// chain service backend
	APIBackend          *APIBackend
	AdmissionApiBackend admission.ApiBackend

	miner    *miner.Miner
	gasPrice *big.Int
	coinbase common.Address

	networkID     uint64
	netRPCService *gccapi.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and coinbase)

	remoteDB database.RemoteDatabase // remoteDB represents an remote distributed database.
}

func (s *gcchainService) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// New creates a new gcchainService object (including the
// initialisation of the common gcchainService object)
func New(ctx *node.ServiceContext, config *Config) (*gcchainService, error) {

	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*configs.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	var remoteDB database.RemoteDatabase
	switch config.PrivateTx.RemoteDBType {
	case private.IPFS:
		remoteDB = database.NewIpfsDB(config.PrivateTx.RemoteDBParams)
		log.Info("Initialize remote database", "database", "IPFS")
	case private.Dummy:
		remoteDB = new(database.DummyDatabase)
		log.Info("Initialize remote database", "database", "Dummy")
	default:
		remoteDB = database.NewIpfsDB(private.DefaultIpfsUrl)
		log.Info("Initialize remote database", "database", "IPFS")
	}

	gcc := &gcchainService{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		shutdownChan:   make(chan bool),
		networkID:      config.NetworkId,
		gasPrice:       config.GasPrice,
		coinbase:       config.Gccbase,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, configs.BloomBitsBlocks),
		remoteDB:       remoteDB,
	}

	gcc.engine = gcc.CreateConsensusEngine(ctx, chainConfig, chainDb)
	if gcc.engine == nil {
		return nil, errBadEngine
	}
	gcc.APIBackend = &APIBackend{gcc, nil}

	// gas related
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	gcc.APIBackend.gpo = gasprice.NewOracle(gcc.APIBackend, gpoParams)

	contractAddrs := configs.ChainConfigInfo().Dpos.Contracts

	contractClient := gccapi.NewPublicBlockChainAPI(gcc.APIBackend)
	primitive_backend.GetApiBackendHolderInstance().Init(gcc.APIBackend, contractClient)
	if dpos, ok := gcc.engine.(*dpos.Dpos); ok {
		dpos.SetCampaignBackend(contractAddrs[configs.ContractCampaign], primitive_backend.GetChainClient())
		dpos.SetRptBackend(contractAddrs[configs.ContractRpt], primitive_backend.GetChainClient())
		dpos.SetRNodeBackend(contractAddrs[configs.ContractRnode], primitive_backend.GetChainClient())
	}

	log.Info("Initialising gcchain protocol", "versions", ProtocolVersions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		bcVersion := rawdb.ReadDatabaseVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run gcchain upgradedb.\n", bcVersion, core.BlockChainVersion)
		}
		rawdb.WriteDatabaseVersion(chainDb, core.BlockChainVersion)
	}

	var (
		vmConfig    = vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}
		cacheConfig = &core.CacheConfig{Disabled: config.NoPruning, TrieNodeLimit: config.TrieCache, TrieTimeLimit: config.TrieTimeout}
	)
	gcc.blockchain, err = core.NewBlockChain(chainDb, cacheConfig, gcc.chainConfig, gcc.engine, vmConfig, remoteDB, ctx.AccountManager)
	if err != nil {
		return nil, err
	}
	gcc.blockchain.SetTypeMux(gcc.eventMux)
	gcc.blockchain.SetSyncMode(config.SyncMode)

	// admission must initialize after blockchain has been initialized
	gcc.AdmissionApiBackend = admission.NewAdmissionApiBackend(gcc.blockchain, gcc.coinbase,
		contractAddrs[configs.ContractAdmission],
		contractAddrs[configs.ContractCampaign],
		contractAddrs[configs.ContractRnode],
		contractAddrs[configs.ContractNetwork])

	if dpos, ok := gcc.engine.(*dpos.Dpos); ok {
		dpos.SetupAdmission(gcc.AdmissionApiBackend)
		dpos.SetChain(gcc.blockchain)
	}

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*configs.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		gcc.blockchain.SetHead(compat.RewindTo)
		rawdb.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	gcc.bloomIndexer.Start(gcc.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	gcc.txPool = core.NewTxPool(config.TxPool, gcc.chainConfig, gcc.blockchain)

	if gcc.protocolManager, err = NewProtocolManager(gcc.chainConfig, config.NetworkId, gcc.eventMux, gcc.txPool, gcc.engine, gcc.blockchain, chainDb, gcc.coinbase, config.SyncMode); err != nil {
		return nil, err
	}

	gcc.miner = miner.New(gcc, gcc.chainConfig, gcc.EventMux(), gcc.engine)

	return gcc, nil
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (database.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// SetAsMiner sets dpos engine as miner
func (s *gcchainService) SetAsMiner(isMiner bool) {
	if dpos, ok := s.engine.(*dpos.Dpos); ok {
		if dpos.IsValidator() {
			return // not execute miner-related operations if node is validator
		}
		dpos.SetAsMiner(isMiner)
	}
}

// SetAsValidator sets the node as validator and it cannot set back to false once it is set to true.
func (s *gcchainService) SetAsValidator() {
	s.engine.(*dpos.Dpos).SetAsValidator(true)
}

// CreateConsensusEngine creates the required type of consensus engine instance for an gcchain service
func (s *gcchainService) CreateConsensusEngine(ctx *node.ServiceContext, chainConfig *configs.ChainConfig,
	db database.Database) consensus.Engine {
	eb, err := s.Coinbase()
	if err != nil {
		log.Debug("coinbase is not set, but is allowed for non-miner node", "error", err)
	}
	// If Dpos is requested, set it up
	if chainConfig.Dpos != nil {
		dpos := dpos.New(chainConfig.Dpos, db)
		if eb != (common.Address{}) {
			wallet, err := s.accountManager.Find(accounts.Account{Address: eb})
			if wallet == nil || err != nil {
				log.Error("Coinbase account unavailable locally", "err", err)
				return nil
			}
			dpos.Authorize(eb, wallet.SignHash)
		}
		return dpos
	}
	return nil
}

// APIs return the collection of RPC services the gcc package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *gcchainService) APIs() []rpc.API {
	apis := gccapi.GetAPIs(s.APIBackend)

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	// Append any APIs exposed explicitly by the admission control
	apis = append(apis, s.AdmissionApiBackend.Apis()...)

	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicgcchainAPI(s),
			Public:    true,
		},
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		},
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   syncer.NewPublicDownloaderAPI(s.protocolManager.syncer, s.eventMux),
			Public:    true,
		},
		{
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.APIBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *gcchainService) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *gcchainService) Coinbase() (coinbase common.Address, err error) {
	s.lock.RLock()
	coinbase = s.coinbase
	s.lock.RUnlock()

	if coinbase != (common.Address{}) {
		return coinbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accs := wallets[0].Accounts(); len(accs) > 0 {
			coinbase = accs[0].Address

			s.lock.Lock()
			s.coinbase = coinbase
			s.lock.Unlock()

			log.Info("Coinbase automatically configured", "address", coinbase)
			return coinbase, nil
		}
	}
	return common.Address{}, fmt.Errorf("coinbase must be explicitly specified")
}

// SetCoinbase sets the mining reward address.
func (s *gcchainService) SetCoinbase(coinbase common.Address) {
	s.lock.Lock()
	s.coinbase = coinbase
	s.lock.Unlock()

	s.miner.SetCoinbase(coinbase)
}

func (s *gcchainService) StartMining(local bool) error {
	if s.IsMining() {
		return nil
	}

	coinbase, err := s.Coinbase()
	if err != nil {
		log.Error("Cannot start mining without coinbase", "err", err)
		return fmt.Errorf("coinbase missing: %v", err)
	}

	// post-requisite: miner.isMining == true && dpos.IsMiner() == true && dpos.isToCampaign == true
	if dpos, ok := s.engine.(*dpos.Dpos); ok {
		if dpos.IsValidator() {
			return errForbidValidatorMining
		}
		if s.AdmissionApiBackend.AdmissionKey() == nil {
			return errNotAdmissionKey
		}

		if dpos.Coinbase() != coinbase {
			wallet, err := s.accountManager.Find(accounts.Account{Address: coinbase})
			if wallet == nil || err != nil {
				log.Error("Etherbase account unavailable locally", "err", err)
				return nil
			}
			dpos.Authorize(coinbase, wallet.SignHash)
		}

		log.Debug("server.nodeid", "enode", s.server.NodeInfo().Enode)

		dpos.SetToCampaign(true)

		// make sure dpos.StartMining start once
		dpos.SetAsMiner(true)
		go dpos.StartMining(s.blockchain, s.server, s.protocolManager.BroadcastBlock, s.protocolManager.SyncFromPeer, s.protocolManager.SyncFromBestPeer)
		log.Info("start participating campaign", "campaign", dpos.IsToCampaign())
	}

	// Propagate the initial price point to the transaction pool
	s.lock.RLock()
	price := s.gasPrice
	s.lock.RUnlock()

	s.txPool.SetGasPrice(price)

	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so none will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)
	}

	// make sure miner.Start() start once
	if !s.miner.IsMining() {
		go s.miner.Start(coinbase)
	}
	return nil
}

func (s *gcchainService) SetupValidator() error {
	if dpos, ok := s.engine.(*dpos.Dpos); ok {
		dpos.SetAsValidator(true)
		dpos.SetupAsValidator(s.blockchain, s.server, s.protocolManager.BroadcastBlock, s.protocolManager.SyncFromPeer, s.protocolManager.SyncFromBestPeer)
	}
	return nil
}

func (s *gcchainService) StopMining() {
	if !s.IsMining() {
		return
	}

	if dpos, ok := s.engine.(*dpos.Dpos); ok {
		// for dpos, keep miner mining, just stop participating campaign
		dpos.SetToCampaign(false)
		log.Info("stopped participating campaign", "campaign", dpos.IsToCampaign())
	} else {
		s.miner.Stop()
	}
}

func (s *gcchainService) IsMining() bool {
	// post-requisite: miner.isMining == true && dpos.IsMiner() == true && dpos.isToCampaign == true
	if dpos, ok := s.engine.(*dpos.Dpos); ok {
		return s.miner.IsMining() && dpos.IsMiner() && dpos.IsToCampaign()
	} else {
		return s.miner.IsMining()
	}
}

func (s *gcchainService) Miner() *miner.Miner { return s.miner }

func (s *gcchainService) AccountManager() *accounts.Manager { return s.accountManager }
func (s *gcchainService) BlockChain() *core.BlockChain      { return s.blockchain }
func (s *gcchainService) TxPool() *core.TxPool              { return s.txPool }
func (s *gcchainService) EventMux() *event.TypeMux          { return s.eventMux }
func (s *gcchainService) Engine() consensus.Engine          { return s.engine }
func (s *gcchainService) ChainDb() database.Database        { return s.chainDb }
func (s *gcchainService) IsListening() bool                 { return true }                                           // Always listening
func (s *gcchainService) CpcVersion() int                   { return int(s.protocolManager.SubProtocols[0].Version) } // the first protocol is the latest version.
func (s *gcchainService) NetVersion() uint64                { return s.networkID }
func (s *gcchainService) Downloader() syncer.Syncer {
	return s.protocolManager.syncer
}
func (s *gcchainService) RemoteDB() database.RemoteDatabase { return s.remoteDB }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *gcchainService) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// start implements node.service, starting all internal goroutines needed by the
// gcchain protocol implementation.
func (s *gcchainService) Start(srvr *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = gccapi.NewPublicNetAPI(srvr, s.NetVersion())

	s.server = srvr

	log.Info("gcchainService started")

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers

	// start the networking layer and the light server if requested
	// by this time, the p2p has already started.  we are only starting the upper layer handling.
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}

	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// gcchain protocol.
func (s *gcchainService) Stop() error {
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
