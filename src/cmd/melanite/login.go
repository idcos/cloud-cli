package main

import (
	"fmt"
	"os"
	"runner/sshrunner"
	"util"

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
	// get node info for login
	var nodes, err = repo.FilterNodes(groupName, nodeName)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return fmt.Errorf("found No nodes")
	}

	var n = nodes[0]
	if len(nodes) > 1 {
		fmt.Printf("%-3s\t%-10s\t%-10s\n", "No.", "Name", "Host")
		fmt.Println(util.FgBoldBlue("=========================================================="))
		for index, n := range nodes {
			fmt.Printf("%-3d\t%-10s\t%-10s\n", index+1, n.Name, n.Host)
		}
		var loginNo = util.LoginNo(fmt.Sprintf("Please input the No.(%s<No.<%s) You want to login: ",
			util.FgBoldRed(0), util.FgBoldRed(len(nodes)+1)),
			1,
			len(nodes)+1)
		n = nodes[loginNo-1]
	}

	var runCmd = sshrunner.New(n.User, n.Password, n.KeyPath, n.Host, n.Port)
	runCmd.Login(n.Host, conf.Main.LoginShell)

	return nil
}
