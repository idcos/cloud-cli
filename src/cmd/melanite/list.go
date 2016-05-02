package main

import (
	"fmt"
	"model"

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
			cli.BoolFlag{
				Name:  "a,all",
				Usage: "is list all info about node?",
			},
		},
		Action: func(c *cli.Context) {
			var groupName = c.String("group")
			var nodeName = c.String("node")
			var isDisplayAll = c.Bool("all")
			if err := listNodes(groupName, nodeName, isDisplayAll); err != nil {
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

func listNodes(groupName, nodeName string, isDisplayAll bool) error {
	repo := GetRepo()

	var groups, err = repo.FilterNodeGroupsAndNodes(groupName, nodeName)
	if err != nil {
		return err
	}

	if isDisplayAll {
		displayDetailInfo(groups)
	} else {
		displaySimpleInfo(groups)
	}

	return nil
}

func displaySimpleInfo(groups []model.NodeGroup) {
	for _, g := range groups {
		fmt.Printf("Group(%s) Nodes: >>>\n", util.FgBoldGreen(g.Name))
		fmt.Printf("%-3s\t%-10s\t%-10s\n", "No.", "Name", "IP")
		fmt.Println(util.FgBoldBlue("=========================================================="))
		for index, n := range g.Nodes {
			fmt.Printf("%-3d\t%-10s\t%-10s\n", index+1, n.Name, n.Host)
		}
		fmt.Println("")
	}
}

func displayDetailInfo(groups []model.NodeGroup) {
	for _, g := range groups {
		fmt.Printf("Group(%s) Nodes: >>>\n", util.FgBoldGreen(g.Name))
		fmt.Printf("%-3s\t%-10s\t%-30s\t%-5s\t%-8s\t%-15s\t%-20s\n", "No.", "Name", "IP", "Port", "User", "Password", "KeyPath")
		fmt.Println(util.FgBoldBlue("========================================================================================================"))
		for index, n := range g.Nodes {
			fmt.Printf("%-3d\t%-10s\t%-30s\t%-5d\t%-8s\t%-15s\t%-20s\n", index+1, n.Name, n.Host, n.Port, n.User, n.Password, n.KeyPath)
		}
		fmt.Println("")
	}
}
