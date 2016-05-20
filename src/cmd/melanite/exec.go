package main

import (
	"model"
	"os"
	"runner"
	"runner/sshrunner"
	"time"

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
	NodeNames []string
	User      string
	Cmd       string
	Yes       bool
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
				Usage: "exec command on one group",
			},
			cli.StringSliceFlag{
				Name:  "n,node",
				Value: &cli.StringSlice{},
				Usage: "exec command on one or more nodes",
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
			cli.BoolFlag{
				Name:  "y,yes",
				Usage: "is confirm before excute command?",
			},
		},
		BashComplete: func(c *cli.Context) {
			bashComplete(c)
		},
		Action: func(c *cli.Context) {
			// 如果有 --generate-bash-completion 参数, 则不执行默认命令
			if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
				bashComplete(c)
				return
			}
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
		NodeNames: c.StringSlice("node"),
		User:      c.String("user"),
		Cmd:       c.String("cmd"),
		Yes:       c.Bool("yes"),
	}

	if ep.Cmd == "" {
		return ep, ErrCmdRequired
	}

	return ep, nil
}

func execCmd(ep execParams) error {
	// TODO should use sshrunner from config

	// get node info for exec
	var nodes, _ = repo.FilterNodes(ep.GroupName, ep.NodeNames...)

	if len(nodes) == 0 {
		return ErrNoNodeToExec
	}

	if !ep.Yes && !confirmExec(nodes, ep.User, ep.Cmd) {
		return nil
	}

	// exec cmd on node
	if conf.Main.Sync {
		return syncExecCmd(nodes, ep)
	} else {
		return concurrentExecCmd(nodes, ep)
	}
}

func syncExecCmd(nodes []model.Node, ep execParams) error {
	var allOutputs = make([]*runner.Output, 0)
	var execStart = time.Now()
	for _, n := range nodes {
		fmt.Printf("EXCUTE \"%s\" on %s(%s):\n", util.FgBoldGreen(ep.Cmd), util.FgBoldGreen(n.Name), util.FgBoldGreen(n.Host))
		var runCmd = sshrunner.New(n.User, n.Password, n.KeyPath, n.Host, n.Port)
		var input = runner.Input{
			ExecHost: n.Host,
			ExecUser: ep.User,
			Command:  ep.Cmd,
			Timeout:  time.Duration(conf.Main.Timeout) * time.Second,
		}

		// display result
		output := runCmd.SyncExec(input)
		displayExecResult(output)
		allOutputs = append(allOutputs, output)
	}
	displayTotalExecResult(allOutputs, execStart, time.Now())
	return nil
}

func concurrentExecCmd(nodes []model.Node, ep execParams) error {
	var allOutputs = make([]*runner.Output, 0)
	var concurrentLimitChan = make(chan int, conf.Main.ConcurrentNum)
	var outputChan = make(chan *runner.ConcurrentOutput)

	var execStart = time.Now()
	for _, n := range nodes {
		var runCmd = sshrunner.New(n.User, n.Password, n.KeyPath, n.Host, n.Port)
		var input = runner.Input{
			ExecHost: n.Host,
			ExecUser: ep.User,
			Command:  ep.Cmd,
			Timeout:  time.Duration(conf.Main.Timeout) * time.Second,
		}

		// exec comamnd
		go runCmd.ConcurrentExec(input, outputChan, concurrentLimitChan)
	}

	var totalCnt = len(nodes)
	for ch := range outputChan {
		totalCnt -= 1
		fmt.Printf("EXCUTE \"%s\" on %s(%s):\n", util.FgBoldGreen(ep.Cmd), util.FgBoldGreen(ch.In.ExecUser), util.FgBoldGreen(ch.In.ExecHost))
		displayExecResult(ch.Out)
		allOutputs = append(allOutputs, ch.Out)

		if totalCnt == 0 {
			close(outputChan)
		}
	}

	displayTotalExecResult(allOutputs, execStart, time.Now())
	return nil
}

func displayExecResult(output *runner.Output) {
	if output.Err != nil {
		fmt.Printf("Command exec failed: %s\n", util.FgRed(output.Err))
	}

	if output != nil {
		fmt.Printf(">>>>>>>>>>>>>>>>>>>> STDOUT >>>>>>>>>>>>>>>>>>>>\n%s\n", output.StdOutput)
		if output.Status == runner.Fail {
			fmt.Printf(">>>>>>>>>>>>>>>>>>>> STDERR >>>>>>>>>>>>>>>>>>>>\n%s\n", output.StdError)
		}
		fmt.Printf("time costs: %v\n", output.ExecEnd.Sub(output.ExecStart))
	}
	fmt.Println(util.FgBoldBlue("==========================================================\n"))
}

func displayTotalExecResult(outputs []*runner.Output, execStart, execEnd time.Time) {
	var successCnt, failCnt, timeoutCnt int

	for _, output := range outputs {
		switch output.Status {
		case runner.Success:
			successCnt += 1
		case runner.Fail:
			failCnt += 1
		case runner.Timeout:
			timeoutCnt += 1
		}
	}

	fmt.Printf("total time costs: %v\nEXEC success nodes: %s | fail nodes: %s | timeout nodes: %s\n\n\n",
		execEnd.Sub(execStart),
		util.FgBoldGreen(successCnt),
		util.FgBoldRed(failCnt),
		util.FgBoldYellow(timeoutCnt))
}

func confirmExec(nodes []model.Node, user, cmd string) bool {
	fmt.Printf("%-3s\t%-10s\t%-10s\n", "No.", "Name", "Host")
	fmt.Println("----------------------------------------------------------------------")
	for index, n := range nodes {
		fmt.Printf("%-3d\t%-10s\t%-10s\n", index+1, n.Name, n.Host)
	}

	fmt.Println()
	return util.Confirm(fmt.Sprintf("You want to exec COMMAND(%s) by UESR(%s) at the above nodes, yes/no(y/n) ?",
		util.FgBoldRed(cmd), util.FgBoldRed(user)))
}
