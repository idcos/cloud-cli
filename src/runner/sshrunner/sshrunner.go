package sshrunner

import (
	"fmt"
	"runner"
	"time"
	"utils"

	pb "gopkg.in/cheggaaa/pb.v1"
)

const (
	// AuthTypePwd SSH login by password
	authTypePwd string = "PASSWORD"
	// AuthTypeKey SSH login by ssh-key
	authTypeKey string = "KEY"
)

// SSHRunner execute command by ssh
type SSHRunner struct {
	sshClient  *SSHClient
	sftpClient *SFTPClient
}

func New(user, password, sshKeyPath, host string, port, fileTransBuf int) *SSHRunner {
	sshClient := NewSSHClient(user, password, sshKeyPath, host, port)
	sftpClient := NewSFTPClient(sshClient, fileTransBuf)

	return &SSHRunner{
		sshClient:  sshClient,
		sftpClient: sftpClient,
	}
}

// SyncExec execute command sync
func (sr *SSHRunner) SyncExec(input runner.ExecInput) *runner.ExecOutput {
	var cmd = compositCommand(input)
	var output = &runner.ExecOutput{Status: runner.Fail}

	output.ExecStart = time.Now()
	status, stdout, stderr, err := sr.sshClient.ExecNointeractiveCmd(cmd, input.Timeout)
	output.ExecEnd = time.Now()

	output.Status = status
	output.Err = err
	output.StdOutput = string(stdout.Bytes())
	output.StdError = string(stderr.Bytes())

	return output
}

// ConcurrentExec execute command sync
func (sr *SSHRunner) ConcurrentExec(input runner.ExecInput, outputChan chan *runner.ConcurrentExecOutput, limitChan chan int) {
	limitChan <- 1
	var output = sr.SyncExec(input)
	outputChan <- &runner.ConcurrentExecOutput{In: input, Out: output}
	<-limitChan
}

// Login login to remote server
func (sr *SSHRunner) Login(shell string) error {
	return sr.sshClient.ExecInteractiveCmd(shell)
}

// SyncPut copy file to remote server sync
func (sr *SSHRunner) SyncPut(input runner.RcpInput) *runner.RcpOutput {
	rcpStart := time.Now()
	err := sr.sftpClient.Put(input.SrcPath, input.DstPath, nil)

	return &runner.RcpOutput{
		RcpStart: rcpStart,
		RcpEnd:   time.Now(),
		Err:      err,
	}
}

// SyncGet copy file from remote server sync
func (sr *SSHRunner) SyncGet(input runner.RcpInput) *runner.RcpOutput {
	rcpStart := time.Now()
	err := sr.sftpClient.Get(input.DstPath, input.SrcPath, nil)

	return &runner.RcpOutput{
		RcpStart: rcpStart,
		RcpEnd:   time.Now(),
		Err:      err,
	}
}

// ConcurrentGet copy file to remote server concurrency
func (sr *SSHRunner) ConcurrentGet(input runner.RcpInput, outputChan chan *runner.ConcurrentRcpOutput, limitChan chan int, pool *pb.Pool) {
	limitChan <- 1
	rcpStart := time.Now()

	bar := utils.NewProgressBar(input.RcpHost, input.RcpSize)
	pool.Add(bar)
	err := sr.sftpClient.Get(input.DstPath, input.SrcPath, bar)

	outputChan <- &runner.ConcurrentRcpOutput{In: input, Out: &runner.RcpOutput{
		RcpStart: rcpStart,
		RcpEnd:   time.Now(),
		Err:      err,
	}}
	<-limitChan
}

// ConcurrentPut copy file from remote server concurrency
func (sr *SSHRunner) ConcurrentPut(input runner.RcpInput, outputChan chan *runner.ConcurrentRcpOutput, limitChan chan int, pool *pb.Pool) {
	limitChan <- 1
	rcpStart := time.Now()

	bar := utils.NewProgressBar(input.RcpHost, input.RcpSize)
	pool.Add(bar)
	err := sr.sftpClient.Put(input.SrcPath, input.DstPath, bar)

	outputChan <- &runner.ConcurrentRcpOutput{In: input, Out: &runner.RcpOutput{
		RcpStart: rcpStart,
		RcpEnd:   time.Now(),
		Err:      err,
	}}
	<-limitChan
}

func (sr *SSHRunner) RemotePathSize(input runner.RcpInput) (int64, error) {
	return sr.sftpClient.RemotePathSize(input.SrcPath)
}

func compositCommand(input runner.ExecInput) string {
	return fmt.Sprintf(`su - '%s' -c '%s'`, input.ExecUser, input.Command)
}
