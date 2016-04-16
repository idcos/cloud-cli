package main

import (
	"fmt"
	"os"

	"util"

	"github.com/codegangsta/cli"
)

var (
	version  = "v0.1.0"
	confPath = "~/.melanite"
)

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
		}
	}

	app.Run(os.Args)
}

func createConfFile() {
	if util.FileExist(defaultConfPath) {
		fmt.Println("config file is already exist.")
	}

}
