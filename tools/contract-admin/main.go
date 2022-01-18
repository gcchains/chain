package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/tools/contract-admin/admission"
	"github.com/gcchains/chain/tools/contract-admin/campaign"
	"github.com/gcchains/chain/tools/contract-admin/network"
	"github.com/gcchains/chain/tools/contract-admin/rnode"
	"github.com/gcchains/chain/tools/contract-admin/rpt"
	"github.com/urfave/cli"
)

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Authors = []cli.Author{
		{
			Name:  "The gcchain authors",
			Email: "info@gcchain.io",
		},
	}
	app.Version = configs.Version
	app.Copyright = "LGPL"
	app.Usage = "Executable for the gcchain official contract admin"

	app.Action = cli.ShowAppHelp

	app.Commands = []cli.Command{
		admission.AdmissionCommand,
		campaign.CampaignCommand,
		network.NetworkCommand,
		rnode.RnodeCommand,
		rpt.RptCommand,
	}

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
