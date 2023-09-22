package commands

import (
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	return &cli.App{
		Name:  "account",
		Usage: "Account",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "space",
				Aliases: []string{"s"},
				Usage:   "app serve space",
			},
		},
		Commands: []*cli.Command{
			Serve(),
		},
	}
}
