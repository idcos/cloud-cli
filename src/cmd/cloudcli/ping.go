package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"utils"

	fastping "github.com/tatsushid/go-fastping"
	"github.com/urfave/cli"
)

func initPingSubCmd(app *cli.App) {
	pingSubCmd := cli.Command{
		Name:        "ping",
		Usage:       "ping <option>",
		Description: "ping nodes",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "g,group",
				Value: "*",
				Usage: "list group and it's nodes",
			},
			cli.StringFlag{
				Name:  "n,node",
				Value: "*",
				Usage: "list nodes",
			},
		},
		Action: func(c *cli.Context) error {
			// 如果有 --generate-bash-completion 参数, 则不执行默认命令
			if os.Args[len(os.Args)-1] == "--generate-bash-completion" {
				groupAndNodeComplete(c)
				return nil
			}

			var groupName = c.String("group")
			var nodeName = c.String("node")
			return pingNodes(groupName, nodeName)
		},
	}

	if app.Commands == nil {
		app.Commands = cli.Commands{pingSubCmd}
	} else {
		app.Commands = append(app.Commands, pingSubCmd)
	}
}

func pingNodes(groupName, nodeName string) error {
	var nodes, err = repo.FilterNodes(groupName, nodeName)
	if err != nil {
		return err
	}

	for _, n := range nodes {
		if ping(n.Host) {
			fmt.Printf("%-30s--------------------[%s]\n", n.Host, utils.FgBoldGreen("OK"))
		} else {
			fmt.Printf("%-30s--------------------[%s]\n", n.Host, utils.FgBoldRed("NG"))
		}
	}
	return nil
}

func ping(host string) bool {
	p := fastping.NewPinger()
	p.Network("udp")
	p.MaxRTT = time.Second

	ra, err := net.ResolveIPAddr("ip4:icmp", host)
	if err != nil {
		fmt.Println(err)
		return false
	}
	p.AddIPAddr(ra)

	var pingResult = false
	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		pingResult = true
	}
	p.OnIdle = func() {}

	if err := p.Run(); err != nil {
		log.Error("ping error: %v\n", err)
		pingResult = false
	}
	<-p.Done()

	return pingResult
}
