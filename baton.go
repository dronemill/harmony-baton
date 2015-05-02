package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/dronemill/harmony-client-go"
	"github.com/olekukonko/tablewriter"
)

// Baton is the main app contianer
type Baton struct {
	// the cli app
	App *cli.App

	// hamrony client
	Harmony *harmonyclient.Client
}

// maestroConnect will get a connected maestro client
func (b *Baton) maestroConnect(c *cli.Context) error {
	config := harmonyclient.Config{
		APIHost:      c.String("maestro"),
		APIVersion:   "v1",
		APIVerifySSL: !c.Bool("noverifyssl"),
	}

	var err error
	b.Harmony, err = harmonyclient.NewHarmonyClient(config)

	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed connecting to the maestro: %s", err.Error()))
	}

	return nil
}

// Containers will show help for the contianers section of the app
func (b *Baton) Containers(c *cli.Context) {
	cli.ShowAppHelp(c)
}

// ContainersAdd will add a container
func (b *Baton) ContainersAdd(c *cli.Context) {

}

// ContainersList will list containers
func (b *Baton) ContainersList(c *cli.Context) {
	ctrs, err := b.Harmony.Containers()

	if err != nil {
		fmt.Printf("GOT ERROR: %s\n", err.Error())
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "MachineID", "Name", "Hostname", "CID"})

	for _, v := range *ctrs {
		r := []string{
			v.ID,
			v.MachineID,
			v.Name,
			v.Hostname,
			v.CID,
		}

		table.Append(r)
	}

	fmt.Println()
	table.SetBorder(false)
	table.Render() // Send output
}
