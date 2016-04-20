package main

import "github.com/codegangsta/cli"

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
			if groupName == "" {
				listGroups(groupName)
			} else {
				listNodes(groupName, nodeName)
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

	return nil
}

func listNodes(groupName, nodeName string) error {
	return nil
}
