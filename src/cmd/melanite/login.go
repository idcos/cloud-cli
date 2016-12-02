package main

import (
	"bytes"
	"fmt"
	"model"
	"os"
	"os/exec"
	"runner/sshrunner"
	"strings"

	"utils"

	"io/ioutil"

	"github.com/urfave/cli"
)

var loginOptions = []string{"-g", "--group", "-n", "--node"}

func initLoginSubCmd(app *cli.App) {
	loginSubCmd := cli.Command{
		Name:        "login",
		Usage:       "login <options>",
		Description: "login to one node",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "g,group",
				Value: "*",
				Usage: "the group's node you want to login",
			},
			cli.StringFlag{
				Name:  "n,node",
				Value: "*",
				Usage: "the node you want to login",
			},
		},
		BashComplete: func(c *cli.Context) {
			for _, opt := range loginOptions {
				fmt.Println(opt)
			}
		},
		Action: func(c *cli.Context) error {
			// 如果有 --generate-bash-completion 参数, 则不执行默认命令
			if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
				groupAndNodeComplete(c)
				return nil
			}

			var groupName = c.String("group")
			var nodeName = c.String("node")
			return loginNode(groupName, nodeName)
		},
	}

	if app.Commands == nil {
		app.Commands = cli.Commands{loginSubCmd}
	} else {
		app.Commands = append(app.Commands, loginSubCmd)
	}
}

func loginNode(groupName, nodeName string) error {
	// get node info for login
	var (
		nodes []model.Node
		index int
		err   error
	)
	nodes, err = repo.FilterNodes(groupName, nodeName)
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return fmt.Errorf("found No nodes")
	}

	index, err = selectNodeByVim(nodes)
	if err != nil {
		return err
	}

	var runCmd = sshrunner.New(nodes[index].User,
		nodes[index].Password,
		nodes[index].KeyPath,
		nodes[index].Host,
		nodes[index].Port,
		conf.Main.FileTransBuf)

	if conf.Main.LoginShell == "" {
		return fmt.Errorf("You must set loginShell in .melanite.ini")
	}

	return runCmd.Login(conf.Main.LoginShell)
}

// this function has refered https://github.com/jpalardy/warp, thanks jpalardy
func selectNodeByVim(nodes []model.Node) (int, error) {
	if len(nodes) == 1 {
		return 0, nil
	}

	var (
		in   *bytes.Buffer
		name string
		host string
		err  error
	)

	cmd := exec.Command("vim",
		"-c", "setlocal noreadonly",
		"-c", "setlocal cursorline",
		"-c", "setlocal number",
		"-c", "nnoremap <buffer> <CR> V:w! ~/.picked<CR>:q!<CR>",
		"--noplugin", "-R",
		"-")

	in = bytes.NewBuffer(nil)

	cmd.Stdin = in
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	in.WriteString(fmt.Sprintf("%-10s\t%-10s\n", "Name", "Host"))
	in.WriteString("==========================================================\n")
	for _, n := range nodes {
		in.WriteString(fmt.Sprintf("%-10s\t%-10s\n", n.Name, n.Host))
	}

	err = cmd.Run()
	if err != nil {
		return 0, err
	}

	// get node and host which is selected in vim
	name, host, err = getNodeNameAndHost()
	if err != nil {
		return 0, err
	}

	for index, n := range nodes {
		if n.Name == name && n.Host == host {
			return index, nil
		}
	}

	// clear the file ~/.picked
	return 0, fmt.Errorf("Not fount the node you select, name: %s\thost: %s\n", name, host)
}

func getNodeNameAndHost() (string, string, error) {
	var (
		filePath string
		content  []byte
		err      error
	)

	filePath, err = utils.ConvertHomeDir("~/.picked")
	if err != nil {
		return "", "", err
	}

	content, err = ioutil.ReadFile(filePath)
	if err != nil {
		return "", "", err
	}

	var result = strings.Split(string(content), "\t")
	if len(result) != 2 {
		return "", "", fmt.Errorf("wrong content in tempfile ~/.picked")
	}

	return utils.Trim(result[0], " ", "\t", "\n"), utils.Trim(result[1], " ", "\t", "\n"), nil
}
