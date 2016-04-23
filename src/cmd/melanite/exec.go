package main

import (
	"errors"

	"fmt"

	"github.com/codegangsta/cli"
)

var (
	ErrGroupORNodeRequired = errors.New("option -g/--group or -n/--node is required")
	ErrOnlyGroupOROnlyNode = errors.New("option -g/--group and -n/--node couldn't exist at same time")
	ErrCmdRequired         = errors.New("option -c/--cmd is required")
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
				Value: "",
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

	if ep.GroupName == "" && ep.NodeName == "" {
		return ep, ErrGroupORNodeRequired
	}

	if ep.GroupName != "" && ep.NodeName != "" {
		return ep, ErrOnlyGroupOROnlyNode
	}

	if ep.Cmd == "" {
		return ep, ErrCmdRequired
	}

	return ep, nil
}

func execCmd(ep execParams) error {

	return nil
}
