package iniconf

import "testing"

var iniContent = `
[Main]
sync=false
concurrentNum=5
timeout=300
loginShell=/bin/bash

[Logger]
level=debug
logFile=
logType=console

[DataSource]
type=yaml
#conn=~/tmp/hf-project.yaml
conn=~/tmp/docker_test.yaml

[CmdAlias]
-: puppet agent -t
-: ls /
`

func TestCmdAlias(t *testing.T) {
	var ini = NewContent([]byte(iniContent))
	var config, err = ini.Load()
	if err != nil {
		t.Errorf("load ini file error: %v", err)
	}

	if config.CmdAlias["#1"] != "puppet agent -t" ||
		config.CmdAlias["#2"] != "ls /" {
		t.Errorf("parse CmdAlias error")
	}
}
