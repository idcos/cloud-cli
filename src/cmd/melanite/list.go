package main

import (
	"fmt"

	"util"

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
				Value: "*",
				Usage: "list group and it's nodes",
			},
			cli.StringFlag{
				Name:  "n,node",
				Value: "*",
				Usage: "list nodes",
			},
		},
		Action: func(c *cli.Context) {
			var groupName = c.String("group")
			var nodeName = c.String("node")
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

func listNodes(groupName, nodeName string) error {
	repo := GetRepo()

	var groups, err = repo.FilterNodeGroupsAndNodes(groupName, nodeName)
	if err != nil {
		return err
	}

	for _, g := range groups {
		fmt.Printf("Group(%s) Nodes: >>>\n", util.FgBoldGreen(g.Name))
		fmt.Printf("%-3s\t%-10s\t%-10s\n", "No.", "Name", "IP")

		for index, n := range g.Nodes {
			fmt.Printf("%-3d\t%-10s\t%-10s\n", index+1, n.Name, n.Host)
		}

		fmt.Println()
	}

	return nil
}
