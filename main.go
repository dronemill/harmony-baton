package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var baton = Baton{}

var flagsBase = []cli.Flag{
	cli.StringFlag{
		Name:  "harmony-api",
		Value: "http://harmony.dev:4774",
		Usage: "the maestro to connect to",
	},
	cli.BoolFlag{
		Name:  "noverifyssl",
		Usage: "do not verify api ssl connections",
	},
}

var flagsContianersAdd = append(flagsBase,
	cli.StringFlag{
		Name:  "machine-id",
		Usage: "the machine to run the container on",
	},
	cli.StringFlag{
		Name:  "name",
		Usage: "the contianer name",
	},
	cli.StringFlag{
		Name:  "hostname",
		Usage: "the container hostname",
	},
	cli.StringFlag{
		Name:  "image",
		Usage: "the container image",
	},
	cli.StringFlag{
		Name:  "entry-point",
		Usage: "the container entry-point",
	},
	cli.BoolFlag{
		Name:  "start",
		Usage: "start container after being created",
	},
)

func main() {
	baton.App = cli.NewApp()
	baton.App.Name = "Baton"
	baton.App.Usage = "The instrument of the Maestro."

	baton.App.Flags = flagsBase

	baton.App.Commands = []cli.Command{
		{
			Name:   "containers",
			Usage:  "manage containers",
			Action: baton.Containers,
			Subcommands: []cli.Command{
				{
					Name:   "add",
					Usage:  "add a new container. Returns new container id",
					Flags:  flagsContianersAdd,
					Action: baton.ContainersAdd,
				},
				{
					Name:   "show",
					Usage:  "containers show",
					Flags:  flagsBase,
					Action: baton.ContainersShow,
				},
				{
					Name:   "start",
					Usage:  "containers start <containerID|name>",
					Flags:  flagsBase,
					Action: baton.ContainersStart,
				},
				{
					Name:   "stop",
					Usage:  "containers stop <containerID|name>",
					Flags:  flagsBase,
					Action: baton.ContainersStop,
				},
			},
		},
		{
			Name:   "machines",
			Usage:  "manage machines",
			Action: baton.Machines,
			Subcommands: []cli.Command{
				{
					Name:   "show",
					Usage:  "machines show [machineID | name] [containers]",
					Flags:  flagsBase,
					Action: baton.MachinesShow,
				},
			},
		},
	}

	baton.App.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}

	baton.App.Run(os.Args)
}
