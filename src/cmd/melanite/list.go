package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func initListSubCmd(app *cli.App) {
	listSubCmd := cli.Command{
		Name:        "list",
		Usage:       "list <options>",
		Description: "list groups and nodes",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "g,group",
				Value: "",
				Usage: "list group and it's nodes",
			},
			cli.StringFlag{
				Name:  "n,node",
				Value: "",
				Usage: "list nodes",
			},
		},
		Action: func(c *cli.Context) {
			var groupName = c.String("group")
			var nodeName = c.String("node")
			if groupName != "" && nodeName == "" {
				if err := listGroups(groupName); err != nil {
					fmt.Println(err)
				}
				return
			}
			if err := listNodes(groupName, nodeName); err != nil {
				fmt.Println(err)
			}
		},
	}

	if app.Commands == nil {
		app.Commands = cli.Commands{listSubCmd}
	} else {
		app.Commands = append(app.Commands, listSubCmd)
	}
}

func listGroups(groupName string) error {
	repo := GetRepo()
	log := GetLogger()

	var groups, err = repo.GetNodeGroups(groupName)
	log.Debug("groups: %v", groups)
	if err != nil {
		return err
	}

	fmt.Printf("Groups: >>>\n")
	for index, g := range groups {
		fmt.Printf("[%d] Name: %-10s\n", index+1, g.Name)
	}

	return nil
}

func listNodes(groupName, nodeName string) error {
	repo := GetRepo()

	var nodes, err = repo.GetNodesByGroupName(groupName, nodeName)
	if err != nil {
		return err
	}

	fmt.Printf("Nodes: >>>\n")
	fmt.Printf("%-3s\t%-10s\t%-10s\n", "No.", "Name", "IP")
	for index, n := range nodes {
		fmt.Printf("%-3d\t%-10s\t%-10s\n", index+1, n.Name, n.IP)
	}

	return nil
}
