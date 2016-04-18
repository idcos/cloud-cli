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
		},
		Action: func(c *cli.Context) {
			var groupName = c.String("group")
			if groupName == "" {
				listGroups()
			} else {
				listNodes(groupName)
			}
		},
	}
}

func listGroups() error {
	return nil
}

func listNodes(groupName string) error {
	return nil
}
