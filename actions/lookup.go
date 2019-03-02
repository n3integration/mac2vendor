package actions

import (
	"fmt"
	"log"

	m2v "github.com/n3integration/mac2vendor"
	"gopkg.in/urfave/cli.v1"
)

var (
	mac   string
	quiet bool
)

func init() {
	log.SetFlags(log.LstdFlags)
	register(cli.Command{
		Name:    "lookup",
		Aliases: []string{"resolve"},
		Action:  lookupAction,
		Usage:   "lookup a mac address and resolve its vendor",
		Flags: []cli.Flag{
			cli.StringFlag{
				Destination: &mac,
				Name:        "mac",
				Usage:       "the mac address to resolve",
			},
			cli.BoolFlag{
				Destination: &quiet,
				Name:        "quiet",
				Usage:       "whether or not to run in quiet mode",
			},
		},
	})
}

func lookupAction(_ *cli.Context) error {
	vnd, err := m2v.Lookup(mac)
	if err != nil {
		return err
	}

	if quiet {
		fmt.Println(vnd)
	} else {
		fmt.Printf("   MAC: %s\n", mac)
		fmt.Printf("Vendor: %s\n", vnd)
	}

	return nil
}
