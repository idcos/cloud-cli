package main

import (
	"model"
	"runner"
	"runner/sshrunner"

	"fmt"

	"util"

	"github.com/codegangsta/cli"
)

var (
	// ErrCmdRequired require cmd option
	ErrCmdRequired = fmt.Errorf("option -c/--cmd is required")
	// ErrNoNodeToExec no more node to execute
	ErrNoNodeToExec = fmt.Errorf("found no node to execute")
)

type execParams struct {
	GroupName string
	NodeName  string
	User      string
	Cmd       string
}

func initExecSubCmd(app *cli.App) {
	execSubCmd := cli.Command{
		Name:        "exec",
		Usage:       "exec <options>",
		Description: "exec command on groups or nodes",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "g,group",
				Value: "*",
				Usage: "exec command on group",
			},
			cli.StringFlag{
				Name:  "n,node",
				Value: "",
				Usage: "exec command on node",
			},
			cli.StringFlag{
				Name:  "u,user",
				Value: "root",
				Usage: "user who exec the command",
			},
			cli.StringFlag{
				Name:  "c,cmd",
				Value: "",
				Usage: "command for exec",
			},
		},
		Action: func(c *cli.Context) {
			var ep, err = checkExecParams(c)
			if err != nil {
				fmt.Println(util.FgRed(err))
				cli.ShowCommandHelp(c, "exec")
				return
			}
			if err = execCmd(ep); err != nil {
				fmt.Println(util.FgRed(err))
			}
		},
	}

	if app.Commands == nil {
		app.Commands = cli.Commands{execSubCmd}
	} else {
		app.Commands = append(app.Commands, execSubCmd)
	}
}

func checkExecParams(c *cli.Context) (execParams, error) {
	var ep = execParams{
		GroupName: c.String("group"),
		NodeName:  c.String("node"),
		User:      c.String("user"),
		Cmd:       c.String("cmd"),
	}

	if ep.Cmd == "" {
		return ep, ErrCmdRequired
	}

	return ep, nil
}

func execCmd(ep execParams) error {
	// TODO should use sshrunner from config

	// get node info for exec
	repo := GetRepo()
	var nodes, err = repo.FilterNodes(ep.GroupName, ep.NodeName)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return ErrNoNodeToExec
	}

	if !confirmExec(nodes, ep.User, ep.Cmd) {
		return nil
	}

	// exec cmd on node
	for _, n := range nodes {
		fmt.Printf("Start to excute \"%s\" on %s(%s):\n", util.FgBoldGreen(ep.Cmd), util.FgBoldGreen(n.Name), util.FgBoldGreen(n.Host))
		var runCmd = sshrunner.New(n.User, n.Password, n.KeyPath, n.Host, n.Port)
		var input = runner.Input{
			ExecHost: n.Host,
			ExecUser: ep.User,
			Command:  ep.Cmd,
		}

		// display result
		output, err := runCmd.SyncExec(input)
		displayExecResult(output, err)
	}
	return nil
}

func displayExecResult(output *runner.Output, err error) {
	if err != nil {
		fmt.Printf("Command exec failed: %s\n", util.FgRed(err))
	}

	if output == nil {
		return
	}

	fmt.Printf(">>>>>>>>>>>>>>>>>>>> STDOUT >>>>>>>>>>>>>>>>>>>>\n%s\n", output.StdOutput)
	if output.StdError != "" {
		fmt.Printf(">>>>>>>>>>>>>>>>>>>> STDERR >>>>>>>>>>>>>>>>>>>>\n%s\n", output.StdError)
	}
	fmt.Printf("time costs: %v\n", output.ExecEnd.Sub(output.ExecStart))
	fmt.Println(util.FgBoldBlue("==========================================================\n"))
}

func confirmExec(nodes []model.Node, user, cmd string) bool {
	fmt.Printf("%-3s\t%-10s\t%-10s\n", "No.", "Name", "IP")
	for index, n := range nodes {
		fmt.Printf("%-3d\t%-10s\t%-10s\n", index+1, n.Name, n.Host)
	}

	fmt.Println()
	return util.Confirm(fmt.Sprintf("You want to exec COMMAND(%s) by UESR(%s) at the above nodes, yes/no(y/n) ?",
		util.FgBoldRed(cmd), util.FgBoldRed(user)))
}
