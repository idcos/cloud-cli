package main

import (
	"fmt"
	"os"
	"utils"

	"github.com/urfave/cli"
)

func initPingSubCmd(app *cli.App) {
	pingSubCmd := cli.Command{
		Name:        "ping",
		Usage:       "ping <option>",
		Description: "ping nodes",
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
		Action: func(c *cli.Context) error {
			// 如果有 --generate-bash-completion 参数, 则不执行默认命令
			if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
				groupAndNodeComplete(c)
				return nil
			}

			var groupName = c.String("group")
			var nodeName = c.String("node")
			return pingNodes(groupName, nodeName)
		},
	}

	if app.Commands == nil {
		app.Commands = cli.Commands{pingSubCmd}
	} else {
		app.Commands = append(app.Commands, pingSubCmd)
	}
}

func pingNodes(groupName, nodeName string) error {
	var nodes, err = repo.FilterNodes(groupName, nodeName)
	if err != nil {
		return err
	}

	for _, n := range nodes {
		if pingTimes(n.Host, 30) {
			fmt.Printf("%-30s--------------------[%s]\n", n.Host, utils.FgBoldGreen("OK"))
		} else {
			fmt.Printf("%-30s--------------------[%s]\n", n.Host, utils.FgBoldRed("NG"))
		}
	}
	return nil
}

func pingTimes(host string, times int) bool {

	for i := 0; i < times; i++ {
		if utils.Ping(host, 2) {
			return true
		}
	}

	return false
}
