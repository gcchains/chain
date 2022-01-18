package main

import (
	"os"

	"github.com/gcchains/chain/cmd/gcchain/commons"
	"github.com/gcchains/chain/commons/log"
	"github.com/gcchains/chain/configs"
	"github.com/gcchains/chain/tools/transfer/config"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "keystore checker"
	app.Usage = "Try to decrypt keystore with password.\n\t\tExample:./keystore_checker --ks /tmp/keystore/key21"
	app.Version = configs.Version
	app.Copyright = "LGPL"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "keystore, ks",
			Usage: "Keystore file path for from address",
			Value: "/tmp/keystore/key1",
		},
	}
	app.Action = func(c *cli.Context) error {
		keystorePath := c.String("keystore")

		if !c.IsSet("keystore") {
			cli.ShowAppHelp(c)
			return nil
		}

		// ask for password
		prompt := "Input password to unlocking account"
		password, _ := commons.ReadPassword(prompt, false)

		config.SetConfig("", keystorePath)
		// decrypt keystore
		_, _, fromAddress, _, _, err := config.DecryptKeystore(password)
		if err == nil {
			log.Infof("Decrypt success for : %x", fromAddress)
		} else {
			log.Infof("Decrypt failed")
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("run application failed", "err", err)
	}
}
