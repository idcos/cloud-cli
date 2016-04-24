package util

import "testing"

func TestExecLinuxCmd(t *testing.T) {
	if _, err := ExecLinuxCmd("uptime"); err != nil {
		t.Errorf("exec command error: %v", err)
	}

	if _, err := ExecLinuxCmd("notfoundcmd"); err == nil {
		t.Errorf("exec command error: %v", err)
	}
}
