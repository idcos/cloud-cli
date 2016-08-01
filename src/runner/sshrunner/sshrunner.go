package sshrunner

import (
	"fmt"
	"runner"
	"time"
)

const (
	// AuthTypePwd SSH login by password
	authTypePwd string = "PASSWORD"
	// AuthTypeKey SSH login by ssh-key
	authTypeKey string = "KEY"
)

// SSHRunner execute command by ssh
type SSHRunner struct {
	client *SSHClient
}

func New(user, password, sshKeyPath, host string, port int) *SSHRunner {

	return &SSHRunner{
		client: NewSSHClient(user, password, sshKeyPath, host, port),
	}
}

// SyncExec execute command sync
func (sr *SSHRunner) SyncExec(input runner.ExecInput) *runner.ExecOutput {
	var cmd = compositCommand(input)
	var output = &runner.ExecOutput{Status: runner.Fail}

	output.ExecStart = time.Now()
	status, stdout, stderr, err := sr.client.ExecNointeractiveCmd(cmd, input.Timeout)
	output.ExecEnd = time.Now()

	output.Status = status
	output.Err = err
	output.StdOutput = string(stdout.Bytes())
	output.StdError = string(stderr.Bytes())

	return output
}

// ConcurrentExec execute command sync
func (sr *SSHRunner) ConcurrentExec(input runner.ExecInput, outputChan chan *runner.ConcurrentOutput, limitChan chan int) {
	limitChan <- 1
	var output = sr.SyncExec(input)
	outputChan <- &runner.ConcurrentOutput{In: input, Out: output}
	<-limitChan
}

// Login login to remote server
func (sr *SSHRunner) Login(shell string) error {
	return sr.client.ExecInteractiveCmd(shell)
}

// SyncPut copy file to remote server sync
func (sr *SSHRunner) SyncPut(input runner.RcpInput) *runner.RcpOutput {
	rcpStart := time.Now()
	err := sr.client.Put(input.SrcPath, input.DstPath)

	return &runner.RcpOutput{
		RcpStart: rcpStart,
		RcpEnd:   time.Now(),
		Err:      err,
	}
}

func compositCommand(input runner.ExecInput) string {
	return fmt.Sprintf(`su - '%s' -c '%s'`, input.ExecUser, input.Command)
}
