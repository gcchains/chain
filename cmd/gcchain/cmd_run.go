

package main

import (
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gcchains/chain/accounts"
	"github.com/gcchains/chain/accounts/keystore"
	"github.com/gcchains/chain/api/gcclient"
	"github.com/gcchains/chain/cmd/gcchain/flags"
	"github.com/gcchains/chain/commons/chainmetrics"
	"github.com/gcchains/chain/commons/log"
	times "github.com/gcchains/chain/commons/time"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/consensus/dpos"
	"github.com/gcchains/chain/contracts/dpos/primitive_register"
	"github.com/gcchains/chain/internal/profile"
	"github.com/gcchains/chain/node"
	"github.com/gcchains/chain/protocols/gcc"
	"github.com/urfave/cli"
)

var runCommand cli.Command

// BusyWarning is a warning msg
const (
	BusyWarning = `


################################################################################
################################################################################

	You are either a CURRENT or a FUTURE PROPOSER.
	Please use =gcchain campaign =command to quit mining first.

################################################################################
################################################################################


	`

	InvalidSystemClockWarn = `


################################################################################
################################################################################

	System clock need to be synchronized.
		
################################################################################
################################################################################


	`
)

func init() {
	runFlags := append([]cli.Flag(nil), flags.RpcFlags...)
	runFlags = append(runFlags, flags.GeneralFlags...)
	runFlags = append(runFlags, flags.NodeFlags...)
	runFlags = append(runFlags, flags.MinerFlags...)
	runFlags = append(runFlags, flags.SyncFlags...)
	runFlags = append(runFlags, flags.P2pFlags...)
	runFlags = append(runFlags, flags.AccountFlags...)
	runFlags = append(runFlags, flags.ChainFlags...)
	runFlags = append(runFlags, flags.LogFlags...)
	runCommand = cli.Command{
		Action: run,
		Name:   "run",
		Flags:  runFlags,
		Usage:  "Run a gcchain node",
		Before: func(ctx *cli.Context) error {
			return nil
		},
		After: func(ctx *cli.Context) error {
			if ctx.IsSet(flags.ProfileFlagName) {
				profile.Stop()
			}
			log.Info("Exit \"gcchain run\" command")
			return nil
		},
	}
}

func run(ctx *cli.Context) error {
	isRunningChain = true
	err := times.InvalidSystemClock()
	if err != nil {
		log.Warn(InvalidSystemClockWarn)
		log.Fatalf("System clock need to be synchronized.there is more than %v seconds gap between ntp and this server", times.MaxGapDuration)
	}

	if ctx.IsSet(flags.MineFlagName) && ctx.IsSet(flags.ValidatorFlagName) {
		log.Fatalf("A node cannot be both miner and validator.")
	}

	n := createNode(ctx)
	bootstrap(ctx, n)
	n.Wait()
	return nil
}

// Register chain services for a *full* node.
func registerChainService(cfg *gcc.Config, n *node.Node, cliCtx *cli.Context) {
	err := n.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		fullNode, err := gcc.New(ctx, cfg)
		primitive_register.RegisterPrimitiveContracts()

		if cliCtx.Bool(flags.MineFlagName) {
			fullNode.SetAsMiner(true)
		}

		if cliCtx.Bool(flags.ValidatorFlagName) {
			fullNode.SetAsValidator()
		}
		return fullNode, err
	})
	if err != nil {
		log.Fatalf("Failed to register the chain service: %v", err)
	}
}

// Creates a node with chain services registered
func createNode(ctx *cli.Context) *node.Node {
	cfg, n := newConfigNode(ctx)
	registerChainService(&cfg.Gcc, n, ctx)
	return n
}

// Starts up the node
func startNode(n *node.Node) {
	// launch the node itself
	if err := n.Start(); err != nil {
		log.Fatalf("Error starting protocol n: %v", err)
	}
}

// makePasswordList reads password lines from the file specified by the global --password flag.
func makePasswordList(ctx *cli.Context) []string {
	path := ctx.String(flags.PasswordFlagName)
	if path == "" {
		return nil
	}
	text, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read password file: %v", err)
	}
	lines := strings.Split(string(text), "\n")
	// Sanitise DOS line endings.
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines
}

func unlockAccounts(ctx *cli.Context, n *node.Node) *keystore.Key {
	ks := n.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
	passwords := makePasswordList(ctx)
	unlock := ctx.String("unlock")
	unlocks := strings.FieldsFunc(unlock, func(c rune) bool { return c == ',' })
	for i, account := range unlocks {
		// log.Infof("%v, %v\n", i, account)
		if i < len(passwords) {
			_, key, err := unlockAccountWithPassword(ks, account, passwords[i])
			if err != nil {
				log.Error("unlock account error", "err", err)
				return nil
			}
			return key
		} else {
			_, _, key, err := unlockAccountWithPrompt(ks, account)
			if err != nil {
				log.Error("unlock account error", "err", err)
				return nil
			}
			return key
		}
	}
	return nil
}

// TODO @chengxin @xumx please be sure about the underlying logic.
// cf. those in the keystore package.
func handleWallet(n *node.Node) {
	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent, 16)
	n.AccountManager().Subscribe(events)

	go func() {
		// Create a chain state reader for self-derivation
		rpcClient, err := n.Attach()
		if err != nil {
			log.Fatalf("Failed to attach to self: %v", err)
		}
		stateReader := gcclient.NewClient(rpcClient)

		// Open any wallets already attached
		for _, wallet := range n.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				if event.Wallet.URL().Scheme == "ledger" {
					event.Wallet.SelfDerive(accounts.DefaultLedgerBaseDerivationPath, stateReader)
				} else {
					event.Wallet.SelfDerive(accounts.DefaultBaseDerivationPath, stateReader)
				}

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()
}

func setupMining(ctx *cli.Context, n *node.Node, key *keystore.Key) {
	var gcchainService *gcc.gcchainService
	// gcchainService will point to the real gcchain service in n.services
	if err := n.Service(&gcchainService); err != nil {
		log.Fatalf("gcchain service not running: %v", err)
	}

	gcchainService.AdmissionApiBackend.SetAdmissionKey(key)
	if configs.IgnoreNetworkStatusCheck {
		gcchainService.AdmissionApiBackend.IgnoreNetworkCheck()
	}

	if ctx.Bool(flags.MineFlagName) {
		if err := gcchainService.StartMining(true); err != nil {
			log.Fatalf("Failed to start mining: %v", err)
		}
	}

	if ctx.Bool(flags.ValidatorFlagName) {
		if err := gcchainService.SetupValidator(); err != nil {
			log.Fatalf("Failed to setup validator: %v", err)
		}
	}
}

func handleInterrupt(n *node.Node) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)

	var gcchainService *gcc.gcchainService
	if err := n.Service(&gcchainService); err != nil {
		log.Fatalf("gcchain service not running: %v", err)
	}

WaitSignal:
	signal := <-sigc
	log.Info("Got interrupt")
	// Warn to not to stop if local coinbase is a current or future proposer!
	if coinbase, err := gcchainService.Coinbase(); err == nil && signal != syscall.SIGTERM {
		if gcchainService.Engine().(*dpos.Dpos).IsCurrentOrFutureProposer(coinbase) {
			log.Warn(BusyWarning)
			goto WaitSignal
		}
	}

	log.Info("Shutting down...")

	go n.Stop()
	for i := 10; i > 0; i-- {
		<-sigc
		if i > 1 {
			log.Warn("Already shutting down, interrupt more to panic.", "times", i-1)
		}
	}
}

func bootstrap(ctx *cli.Context, n *node.Node) {
	// start profiling
	if ctx.IsSet(flags.ProfileFlagName) {
		if err := profile.Start(ctx); err != nil {
			log.Fatalf("start profiling failed: %v\n", err)
		}
	}

	// init metrics
	if ctx.IsSet(flags.MetricGatewayFlagName) {
		chainmetrics.InitMetrics(ctx.String(flags.PortFlagName), ctx.String(flags.MetricGatewayFlagName))
	}

	startNode(n)
	key := unlockAccounts(ctx, n)
	handleWallet(n)
	setupMining(ctx, n, key)
	// handle user interrupt
	go handleInterrupt(n)
}
