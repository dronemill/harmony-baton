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
		APIHost:      c.String("harmony-api"),
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

// Machines will show help for the contianers section of the app
func (b *Baton) Machines(c *cli.Context) {
	cli.ShowAppHelp(c)
}

// ContainersAdd will add a container
func (b *Baton) ContainersAdd(c *cli.Context) {
	if err := b.maestroConnect(c); err != nil {
		fmt.Printf("%s\n\n", err.Error())
		return
	}

	machineID := c.String("machine-id")
	if machineID == "" {
		fmt.Println("machine-id is required")
		return
	}

	name := c.String("name")
	if name == "" {
		fmt.Println("name is required")
		return
	}

	hostname := c.String("hostname")
	if hostname == "" {
		fmt.Println("hostname is required")
		return
	}

	image := c.String("image")
	if image == "" {
		fmt.Println("image is required")
		return
	}

	entryPoint := c.String("entry-point")

	cntr := &harmonyclient.Container{
		MachineID:  machineID,
		Name:       name,
		Hostname:   hostname,
		Image:      image,
		EntryPoint: entryPoint,
	}

	newCntr, err := b.Harmony.ContainersAdd(cntr)

	if err != nil {
		fmt.Printf("Error encountered while attempting to create new container: %s\n", err.Error())
		return
	}

	fmt.Printf("%s\n", newCntr.ID)
}

// ContainersList will list containers
func (b *Baton) ContainersList(c *cli.Context) {
	if err := b.maestroConnect(c); err != nil {
		fmt.Printf("%s\n\n", err.Error())
		return
	}

	ctrs, err := b.Harmony.Containers()

	if err != nil {
		fmt.Printf("GOT ERROR: %s\n", err.Error())
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Hostname", "MachineID", "CID"})

	for _, v := range *ctrs {
		r := []string{
			v.ID,
			v.Name,
			v.Hostname,
			v.MachineID,
			v.CID,
		}

		table.Append(r)
	}

	fmt.Println()
	table.SetBorder(false)
	table.Render() // Send output
}

// MachinesList will list machines
func (b *Baton) MachinesList(c *cli.Context) {
	if err := b.maestroConnect(c); err != nil {
		fmt.Printf("%s\n\n", err.Error())
		return
	}

	ctrs, err := b.Harmony.Machines()

	if err != nil {
		fmt.Printf("GOT ERROR: %s\n", err.Error())
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Hostname"})

	for _, v := range *ctrs {
		r := []string{
			v.ID,
			v.Name,
			v.Hostname,
		}

		table.Append(r)
	}

	fmt.Println()
	table.SetBorder(false)
	table.Render() // Send output
}
