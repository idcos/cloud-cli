package main

import (
	"runner"
	"runner/sshrunner"

	"fmt"

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
				fmt.Println(err)
				cli.ShowCommandHelp(c, "exec")
				return
			}
			if err = execCmd(ep); err != nil {
				fmt.Println(err)
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
	var groups, err = repo.FilterNodeGroupsAndNodes(ep.GroupName, ep.NodeName)
	if err != nil {
		return err
	}

	if len(groups) == 0 {
		return ErrNoNodeToExec
	}

	// exec cmd on node
	for _, g := range groups {
		for _, n := range g.Nodes {
			fmt.Printf("start exec cmd(%s) on Group(%s)->Node(%s): >>>\n", ep.Cmd, g.Name, n.Name)
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
	}
	return nil
}

func displayExecResult(output *runner.Output, err error) {
	if err != nil {
		fmt.Printf("Command exec failed: %s\n", err)
	}

	if output == nil {
		return
	}
	fmt.Printf("START TIME: %s\n", output.ExecStart.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("END   TIME: %s\n", output.ExecEnd.Format("2006-01-02 15:04:05.000"))
	fmt.Printf("STDOUT >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n%s\n", output.StdOutput)
	if output.StdError != "" {
		fmt.Printf("STDERR >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n%s\n", output.StdError)
	}
}
