package main

import (
	"log"
	"os"

	"github.com/n3integration/mac2vendor/actions"
	"gopkg.in/urfave/cli.v1"
)

var Version = "snapshot"

func main() {
	app := cli.NewApp()
	app.Name = "mac2vnd"
	app.Version = Version
	app.Usage = "mac address to vendor resolution utilities"
	app.Commands = actions.GetCommands()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
