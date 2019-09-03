package app

// TODO: Decouple this from the application, investigate Traefik's CLI args

import (
	"github.com/urfave/cli"
)

func GetOpts() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "run, r",
			Usage: "For running a local CI definition",
		},
	}
}
