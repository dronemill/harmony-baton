package main

import (
	"fmt"
	"os"
	"regexp"

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

//
// CONTAINERS
//

// Containers will show help for the contianers section of the app
func (b *Baton) Containers(c *cli.Context) {
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
		Enabled:    c.Bool("start"),
	}

	newCntr, err := b.Harmony.ContainersAdd(cntr)

	if err != nil {
		fmt.Printf("Error encountered while attempting to create new container: %s\n", err.Error())
		return
	}

	fmt.Printf("%s\n", newCntr.ID)
}

// ContainersShow will show containers
func (b *Baton) ContainersShow(c *cli.Context) {
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
		cid := v.CID
		if len(cid) > 12 {
			cid = v.CID[0:12]
		}
		r := []string{
			v.ID,
			v.Name,
			v.Hostname,
			v.MachineID,
			cid,
		}

		table.Append(r)
	}

	fmt.Println()
	table.SetBorder(false)
	table.Render() // Send output
}

// ContainersStart will start a container
func (b *Baton) ContainersStart(c *cli.Context) {
	b.containersSetEnabled(c, true)
}

// ContainersStop will start a container
func (b *Baton) ContainersStop(c *cli.Context) {
	b.containersSetEnabled(c, false)
}

// ContainersStop will stop a container
func (b *Baton) containersSetEnabled(c *cli.Context, enabled bool) {
	if len(c.Args()) == 0 {
		fmt.Print("ContainerID or Name is required\n\n")
		cli.ShowAppHelp(c)
		return
	}

	if err := b.maestroConnect(c); err != nil {
		fmt.Printf("%s\n\n", err.Error())
		return
	}

	container := b.findContainer(c, c.Args()[0])
	containerID := container.ID

	if err := b.Harmony.ContainersEnabledUpdate(containerID, enabled); err != nil {
		fmt.Printf("%s\n\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("OK\n")
}

// findContainer by its name or ID
func (b *Baton) findContainer(c *cli.Context, containerID string) (container *harmonyclient.Container) {
	var err error
	if matched, _ := regexp.MatchString("^[0-9]*$", containerID); matched {
		container, err = b.Harmony.Container(containerID)
	} else {
		container, err = b.Harmony.ContainerByName(containerID)
	}

	if err != nil {
		fmt.Printf("GOT ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	if container == nil {
		fmt.Printf("[404] Container not found [%s]\n", containerID)
		os.Exit(1)
	}

	return container
}

//
// MACHINES
//

// Machines will show help for the contianers section of the app
func (b *Baton) Machines(c *cli.Context) {
	cli.ShowAppHelp(c)
}

// MachinesShow will show machines
func (b *Baton) MachinesShow(c *cli.Context) {
	if err := b.maestroConnect(c); err != nil {
		fmt.Printf("%s\n\n", err.Error())
		return
	}

	if len(c.Args()) == 0 {
		b.showMachines(c)
		return
	}

	machineID := c.Args()[0]
	if matched, _ := regexp.MatchString("^[0-9]*$", machineID); matched {
		b.showMachineByID(c, machineID)
	} else {
		b.showMachineByName(c, machineID)
	}
}

// showMachines is the command processor for showing all machines
func (b *Baton) showMachines(c *cli.Context) {
	machines, err := b.Harmony.Machines()

	if err != nil {
		fmt.Printf("GOT ERROR: %s\n", err.Error())
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Hostname"})

	for _, v := range *machines {
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

// showMachineByID will show a machine by its id
func (b *Baton) showMachineByID(c *cli.Context, ID string) {
	m, err := b.Harmony.Machine(ID)

	if err != nil {
		fmt.Printf("GOT ERROR: %s\n", err.Error())
		return
	}

	if m == nil {
		fmt.Printf("ERROR: machineID '%s' not found\n", ID)
		return
	}

	b.renderMachine(c, m)
}

// showMachineByName will show a machine by its name
func (b *Baton) showMachineByName(c *cli.Context, name string) {
	m, err := b.Harmony.MachineByName(name)

	if err != nil {
		fmt.Printf("GOT ERROR: %s\n", err.Error())
		return
	}

	if m == nil {
		fmt.Printf("ERROR: machine with name '%s' not found\n", name)
		return
	}

	b.renderMachine(c, m)

}

// renderMachine will output a formatted machine
func (b *Baton) renderMachine(c *cli.Context, m *harmonyclient.Machine) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Hostname"})

	r := []string{
		m.ID,
		m.Name,
		m.Hostname,
	}
	table.Append(r)

	fmt.Println()
	table.SetBorder(false)
	table.Render() // Send output

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Containers"})
	for _, v := range m.ContainerIDs {
		r := []string{
			v,
		}

		table.Append(r)
	}

	fmt.Println()
	table.SetBorder(false)
	table.Render() // Send output

}
