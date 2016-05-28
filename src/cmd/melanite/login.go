package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func initLoginSubCmd(app *cli.App) {

	loginSubCmd := cli.Command{
		Name:        "login",
		Usage:       "login <options>",
		Description: "login to one node",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "g,group",
				Value: "*",
				Usage: "the group's node you want to login",
			},
			cli.StringFlag{
				Name:  "n,node",
				Value: "*",
				Usage: "the node you want to login",
			},
		},
		Action: func(c *cli.Context) error {
			// 如果有 --generate-bash-completion 参数, 则不执行默认命令
			if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
				bashComplete(c)
				return nil
			}

			var groupName = c.String("group")
			var nodeName = c.String("node")
			var err error
			if err = loginNode(groupName, nodeName); err != nil {
				fmt.Println(err)
			}
			return err
		},
	}

	if app.Commands == nil {
		app.Commands = cli.Commands{loginSubCmd}
	} else {
		app.Commands = append(app.Commands, loginSubCmd)
	}
}

func loginNode(groupName, nodeName string) error {

	return nil
}
