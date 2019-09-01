package main

import (
	"github.com/urfave/cli"

	"log"
	"os"
	shiphand "shiphand/app"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.1.0"
	app.Name = "shiphand"
	app.Usage = "Job orchestrator for Anchor CI"
	app.Flags = shiphand.GetOpts()
	app.Action = shiphand.AppMain

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
