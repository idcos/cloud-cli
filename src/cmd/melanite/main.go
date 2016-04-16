package main

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"util"

	"github.com/codegangsta/cli"
)

var (
	version  = "v0.1.0"
	confPath = ".melanite.ini"
)

func init() {
	u, _ := user.Current()
	confPath = path.Join(u.HomeDir, confPath)
}

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Name = "Melanite (CLI tool)"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "init",
			Usage: "init melanite's setting file",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("init") {
			createConfFile()
			return
		}

		cli.ShowAppHelp(c)
	}

	app.Run(os.Args)
}

func createConfFile() {
	if util.FileExist(confPath) {
		if !util.Confirm("config file is already existed, Do you want to overwrite it?(y or n)") {
			return
		}
	}

	f, err := os.Create(confPath)
	if err != nil {
		fmt.Println("create config file error: %s", err)
		return
	}

	defer f.Close()

	var defaultConfContent = `[Logger]
level=error
logFile=
logType=console
`
	f.Write([]byte(defaultConfContent))
}
