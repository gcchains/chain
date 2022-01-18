

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/gcchains/chain/api/gcclient"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli"
)

const (
	ENDPOINT = "endpoint"
)

func main() {
	app := cli.NewApp()
	app.Usage = "executable command tool for easy test"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  ENDPOINT + ",ep",
			Value: "http://127.0.0.1:8521",
			Usage: "endpoint like http://127.0.0.1:8521",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "get_block_number",
			Aliases: []string{"bn"},
			Usage:   "get block number",
			Action: func(c *cli.Context) error {
				endpoint := c.GlobalString(ENDPOINT)
				client, err := gcclient.Dial(endpoint)
				if err != nil {
					log.Fatal(err.Error())
				}
				blockNumber := client.GetBlockNumber()
				fmt.Println(blockNumber)
				return nil
			},
		},
		{
			Name:    "balance",
			Aliases: []string{"bal"},
			Usage:   "get balance by address. \n\t\texample: ./testtool bal ${addr}",
			Action: func(c *cli.Context) error {
				endpoint := c.GlobalString(ENDPOINT)
				client, err := gcclient.Dial(endpoint)
				if err != nil {
					log.Fatal(err.Error())
				}
				balance, _ := client.BalanceAt(context.Background(), common.HexToAddress(c.Args().First()), nil)
				fmt.Println(balance)
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
