

package profile

import (
	"path"
	"runtime"

	"github.com/gcchains/chain/cmd/gcchain/flags"
	"github.com/gcchains/chain/commons/log"
	"github.com/urfave/cli"
)

type profileConfig struct {
	memProfileRate          int
	blockProfileRate        int
	traceFileName           string
	cpuFileName             string
	blockingProfileFileName string
	memProfileFileName      string
	pprofAddr               string
}

func getProfileConfig(ctx *cli.Context) profileConfig {
	dirPath := ""
	if ctx.IsSet(flags.ProfileFlagName) {
		dirPath = ctx.String(flags.ProfileFlagName)
	}
	profileAddress := "localhost:8931"
	if ctx.IsSet(flags.ProfileAddressFlagName) {
		profileAddress = ctx.String(flags.ProfileAddressFlagName)
	}
	return profileConfig{
		memProfileRate:          runtime.MemProfileRate,
		blockProfileRate:        1,
		traceFileName:           path.Join(dirPath, "gcchain-trace.trace"),
		cpuFileName:             path.Join(dirPath, "gcchain-cpu.profile"),
		blockingProfileFileName: path.Join(dirPath, "gcchain-block.profile"),
		memProfileFileName:      path.Join(dirPath, "gcchain-heap.profile"),
		pprofAddr:               profileAddress,
	}
}

// Start profiling
func Start(ctx *cli.Context) error {
	log.Info("profile start")
	// start profiling, tracing
	cfg := getProfileConfig(ctx)

	

	// pprof server
	StartPProfWebUi(cfg.pprofAddr)
	return nil
}

// Stops all running profiles, flushing their output to the respective file.
func Stop() {
	
	log.Info("profile stop")
}
