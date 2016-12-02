package main

import (
	"config"
	"fmt"
	"os"
	"os/user"
	"path"

	"utils"

	"config/iniconf"

	"logger"

	"strings"

	"model/yamlrepo"

	"github.com/astaxie/beego/logs"
	"github.com/urfave/cli"
)

var (
	version  = "v0.6.5"
	confPath = ".melanite.ini"
	conf     *config.Config
	log      *logs.BeeLogger
	repo     *yamlrepo.YAMLRepo
)

func init() {
	u, _ := user.Current()
	confPath = path.Join(u.HomeDir, confPath)
}

func main() {
	app := cli.NewApp()
	app.Version = version
	app.EnableBashCompletion = true
	app.Name = "Melanite (CLI tool)"

	if err := initConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initListSubCmd(app)
	initExecSubCmd(app)
	initLoginSubCmd(app)
	initRcpSubCmd(app)

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		log.Error("%s\n", err)
		os.Exit(1)
	}
}

func initConfig() error {
	var err error
	if !utils.FileExist(confPath) {
		if !utils.Confirm("Do you want to create your config file?(y or n)") {
			return fmt.Errorf("You should create your init config file")
		}
		createConfFile()
	}

	load := iniconf.New(confPath)
	conf, err = load.Load()
	if err != nil {
		return err
	}

	if conf.DataSource.Conn == "" || conf.DataSource.Type == "" {
		return fmt.Errorf("You should set DataSource in your config file: %s", confPath)
	}

	// fmt.Printf("conf: %v\n", conf)
	// init log
	if strings.ToLower(conf.Logger.LogType) == "console" {
		log = logger.NewConsoleLogger(conf.Logger.Level)
	} else {
		log = logger.NewFileLogger(conf.Logger.LogFile, conf.Logger.Level)
	}

	// init repo
	repo, _ = yamlrepo.New(conf.DataSource.Conn)
	return nil
}

func createConfFile() {
	f, err := os.Create(confPath)
	if err != nil {
		fmt.Println("create config file error: %s", err)
		return
	}

	defer f.Close()

	var defaultConfContent = `[Main]
sync=true
concurrentNum=5
timeout=30

[Logger]
level=error
logFile=
logType=console

[DataSource]
type=yaml
conn=
`
	f.Write([]byte(defaultConfContent))
}
