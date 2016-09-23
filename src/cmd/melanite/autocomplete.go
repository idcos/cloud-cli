package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func completeGroups() {
	groups, _ := repo.FilterNodeGroups("*")
	for _, g := range groups {
		fmt.Println(g.Name)
	}
}

func completeNodes(gName string) {
	nodes, _ := repo.FilterNodes(gName, "*")
	for _, n := range nodes {
		fmt.Println(n.Name)
	}
}

func groupAndNodeComplete(c *cli.Context) {
	if isAutoComplete(c.String("group")) {
		completeGroups()
	} else if isAutoComplete(c.String("node")) {
		completeNodes(c.String("group"))
	}
}

func isAutoComplete(curStr string) bool {
	// --generate-bash-completion is global option for cli
	// ex. "node" is a multi-option, so framework cli will add [--generate-bash-completion] after -n
	if strings.Contains(curStr, "--generate-bash-completion") {
		return true
	}
	return false
}
