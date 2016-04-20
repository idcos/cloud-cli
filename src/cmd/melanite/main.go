package main

import (
	"config"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"

	"util"

	"config/iniconf"

	"logger"

	"strings"

	"model/yamlrepo"

	"github.com/astaxie/beego/logs"
	"github.com/codegangsta/cli"
)

var (
	version  = "v0.1.0"
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
	app.Name = "Melanite (CLI tool)"

	if err := checkConfigFile(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initListSubCmd(app)

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func GetConfig() *config.Config {
	return conf
}

func GetLogger() *logs.BeeLogger {
	return log
}

func GetRepo() *yamlrepo.YAMLRepo {
	if repo == nil {
		repo, _ = yamlrepo.New(conf.DataSource.Conn)
	}
	return repo
}

func checkConfigFile() error {
	var err error
	if !util.FileExist(confPath) {
		if !util.Confirm("Do you want to create your config file?(y or n)") {
			return errors.New("You should create your init config file")
		}
		createConfFile()
	}

	load := iniconf.New(confPath)
	conf, err = load.Load()
	if err != nil {
		return err
	}

	if conf.DataSource.Conn == "" || conf.DataSource.Type == "" {
		return errors.New(fmt.Sprintf("You should set DataSource in your config file: %s", confPath))
	}

	// fmt.Printf("conf: %v\n", conf)
	if strings.ToLower(conf.Logger.LogType) == "console" {
		log = logger.NewConsoleLogger(conf.Logger.Level)
	} else {
		log = logger.NewFileLogger(conf.Logger.LogFile, conf.Logger.Level)
	}
	return nil
}

func createConfFile() {
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

[DataSource]
type=yaml
conn=
`
	f.Write([]byte(defaultConfContent))
}
