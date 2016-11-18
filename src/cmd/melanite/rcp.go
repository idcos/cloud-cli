package main

import (
	"fmt"
	"model"
	"os"
	"runner"
	"runner/sshrunner"
	"time"

	"utils"

	"github.com/urfave/cli"
)

var (
	// ErrSrcRequired require src option
	ErrSrcRequired = fmt.Errorf("option --src is required")
	// ErrDstRequired require dst option
	ErrDstRequired = fmt.Errorf("option --dst is required")
	// ErrNoNodeToRcp no more node to execute
	ErrNoNodeToRcp = fmt.Errorf("found no node to copy file/directory")
)

type rcpParams struct {
	GroupName string
	NodeNames []string
	User      string
	Src       string
	Dst       string
	Yes       bool
}

func initRcpSubCmd(app *cli.App) {

	putSubCmd := cli.Command{
		Name:        "put",
		Usage:       "put <options>",
		Description: "copy file or directory to remote servers",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "src",
				Value: "",
				Usage: "source file or directory",
			},
			cli.StringFlag{
				Name:  "dst",
				Value: "",
				Usage: "destination *directory*",
			},
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
			cli.BoolFlag{
				Name:  "y,yes",
				Usage: "is confirm before excute command?",
			},
		},
		Action: func(c *cli.Context) error {
			// 如果有 --generate-bash-completion 参数, 则不执行默认命令
			if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
				groupAndNodeComplete(c)
				return nil
			}

			var rp, err = checkRcpParams(c)
			if err != nil {
				fmt.Println(utils.FgRed(err))
				cli.ShowCommandHelp(c, "put")
				return nil
			}

			return rcpCmd(rp, true)
		},
	}

	getSubCmd := cli.Command{
		Name:        "get",
		Usage:       "get <options>",
		Description: "copy file or directory from remote servers",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "src",
				Value: "",
				Usage: "source *directory*",
			},
			cli.StringFlag{
				Name:  "dst",
				Value: "",
				Usage: "destination file or directory",
			},
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
			cli.BoolFlag{
				Name:  "y,yes",
				Usage: "is confirm before excute command?",
			},
		},
		Action: func(c *cli.Context) error {
			// 如果有 --generate-bash-completion 参数, 则不执行默认命令
			if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
				groupAndNodeComplete(c)
				return nil
			}

			var rp, err = checkRcpParams(c)
			if err != nil {
				fmt.Println(utils.FgRed(err))
				cli.ShowCommandHelp(c, "get")
				return nil
			}

			return rcpCmd(rp, false)
		},
	}

	if app.Commands == nil {
		app.Commands = cli.Commands{putSubCmd, getSubCmd}
	} else {
		app.Commands = append(app.Commands, putSubCmd, getSubCmd)
	}
}

func checkRcpParams(c *cli.Context) (rcpParams, error) {
	var rp = rcpParams{
		GroupName: c.String("group"),
		NodeNames: c.StringSlice("node"),
		User:      c.String("user"),
		Src:       c.String("src"),
		Dst:       c.String("dst"),
		Yes:       c.Bool("yes"),
	}

	if rp.Dst == "" {
		return rp, ErrSrcRequired
	}

	if rp.Src == "" {
		return rp, ErrSrcRequired
	}

	return rp, nil
}

func rcpCmd(rp rcpParams, isPut bool) error {
	// TODO should use sshrunner from config

	// get node info for exec
	var nodes, _ = repo.FilterNodes(rp.GroupName, rp.NodeNames...)

	if len(nodes) == 0 {
		return ErrNoNodeToRcp
	}

	if !rp.Yes && !confirmRcp(nodes, rp.User, rp.Src, rp.Dst) {
		return nil
	}

	// exec cmd on node
	if conf.Main.Sync {
		// TODO copy file concurrency
		return nil
	} else {
		return syncRcp(nodes, rp, isPut)
	}
	return nil
}

func confirmRcp(nodes []model.Node, user, from, to string) bool {
	fmt.Printf("%-3s\t%-10s\t%-10s\n", "No.", "Name", "Host")
	fmt.Println("----------------------------------------------------------------------")
	for index, n := range nodes {
		fmt.Printf("%-3d\t%-10s\t%-10s\n", index+1, n.Name, n.Host)
	}

	fmt.Println()
	return utils.Confirm(fmt.Sprintf("You want to copy [%s] to [%s] by UESR(%s) at the above nodes, yes/no(y/n) ?",
		utils.FgBoldRed(from), utils.FgBoldRed(to), utils.FgBoldRed(user)))
}

func syncRcp(nodes []model.Node, rp rcpParams, isPut bool) error {
	var allOutputs = make([]*runner.RcpOutput, 0)
	var execStart = time.Now()
	for _, n := range nodes {
		fmt.Printf("%s(%s):\n", utils.FgBoldGreen(n.Name), utils.FgBoldGreen(n.Host))
		var sftpClient = sshrunner.New(n.User, n.Password, n.KeyPath, n.Host, n.Port)

		var input = runner.RcpInput{
			SrcPath: rp.Src,
			DstPath: rp.Dst,
			RcpHost: n.Host,
			RcpUser: rp.User,
		}

		// display result
		var output *runner.RcpOutput
		if isPut {
			output = sftpClient.SyncPut(input)
		} else {
			output = sftpClient.SyncGet(input)
		}
		displayRcpResult(output)
		allOutputs = append(allOutputs, output)
	}
	displayTotalRcpResult(allOutputs, execStart, time.Now())
	return nil
}

func displayRcpResult(output *runner.RcpOutput) {
	if output.Err != nil {
		fmt.Printf("copy file/directory failed: %s\n", utils.FgRed(output.Err))
	}

	fmt.Printf("time costs: %v\n", output.RcpEnd.Sub(output.RcpStart))
	fmt.Println(utils.FgBoldBlue("==========================================================\n"))
}

func displayTotalRcpResult(outputs []*runner.RcpOutput, rcpStart, rcpEnd time.Time) {
	var successCnt, failCnt int

	for _, output := range outputs {
		if output.Err != nil {
			failCnt += 1
		} else {
			successCnt += 1
		}
	}

	fmt.Printf("total time costs: %v\nRCP success nodes: %s | fail nodes: %s\n\n\n",
		rcpEnd.Sub(rcpStart),
		utils.FgBoldGreen(successCnt),
		utils.FgBoldRed(failCnt))
}
