package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

var app *cli.App

func main() {
	app = cli.NewApp()
	app.Name = "Baton"
	app.Usage = "The instrument of the Maestro."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "maestro",
			Value: "127.0.0.1:4774",
			Usage: "the maestro to connect to",
		},
	}

	app.Action = func(c *cli.Context) {
		fmt.Printf("Maestro: %s\n", c.String("maestro"))
	}

	app.Run(os.Args)
}
