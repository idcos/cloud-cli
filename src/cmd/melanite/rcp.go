package main

import (
	"fmt"
	"os"
	"util"

	"github.com/urfave/cli"
)

var (
	// ErrSrcRequired require src option
	ErrSrcRequired = fmt.Errorf("option --src is required")
	// ErrDstRequired require dst option
	ErrDstRequired = fmt.Errorf("option --dst is required")
)

type rcpParams struct {
	GroupName string
	NodeNames []string
	User      string
	Src       string
	Dst       string
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
		},
		Action: func(c *cli.Context) error {
			// 如果有 --generate-bash-completion 参数, 则不执行默认命令
			if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
				groupAndNodeComplete(c)
				return nil
			}
			return nil
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
				Usage: "source file or directory",
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
		},
		Action: func(c *cli.Context) error {
			// 如果有 --generate-bash-completion 参数, 则不执行默认命令
			if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
				groupAndNodeComplete(c)
				return nil
			}

			var rp, err = checkRcpParams(c)
			if err != nil {
				fmt.Println(util.FgRed(err))
				cli.ShowCommandHelp(c, "exec")
				return err
			}
			if err = rcpCmd(rp); err != nil {
				fmt.Println(util.FgRed(err))
			}

			return nil
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
	}

	if rp.Dst == "" {
		return rp, ErrSrcRequired
	}

	if rp.Src == "" {
		return rp, ErrSrcRequired
	}

	return rp, nil
}

func rcpCmd(rp rcpParams) error {

	return nil
}
