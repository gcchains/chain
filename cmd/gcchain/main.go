

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/gcchains/chain/cmd/gcchain/flags"
	"github.com/gcchains/chain/configs"
	"github.com/urfave/cli"
)

func newApp() *cli.App {
	app := cli.NewApp()
	// the executable name
	app.Name = filepath.Base(os.Args[0])
	app.Authors = []cli.Author{
		{
			Name:  "The gcchain authors",
			Email: "info@gcchain.io",
		},
	}
	app.Version = configs.Version
	app.Copyright = "GCOIN foundation"
	app.Usage = "Executable for the gcchain blockchain network"
	/
	app.Action = cli.ShowAppHelp

	app.Commands = []cli.Command{
		accountCommand,
		runCommand,
		dumpConfigCommand,
		chainCommand,
		campaignCommand,
	}

	// global flags
	app.Flags = append(app.Flags, flags.ConfigFileFlag)

	// maintain order
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}

func main() {
	if err := newApp().Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
