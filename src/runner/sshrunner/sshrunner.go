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
	client *SSHClient
}

func New(user, password, sshKeyPath, host string, port, fileTransBuf int) *SSHRunner {

	return &SSHRunner{
		client: NewSSHClient(user, password, sshKeyPath, host, port, fileTransBuf),
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
func (sr *SSHRunner) ConcurrentExec(input runner.ExecInput, outputChan chan *runner.ConcurrentExecOutput, limitChan chan int) {
	limitChan <- 1
	var output = sr.SyncExec(input)
	outputChan <- &runner.ConcurrentExecOutput{In: input, Out: output}
	<-limitChan
}

// Login login to remote server
func (sr *SSHRunner) Login(shell string) error {
	return sr.client.ExecInteractiveCmd(shell)
}

// SyncPut copy file to remote server sync
func (sr *SSHRunner) SyncPut(input runner.RcpInput) *runner.RcpOutput {
	rcpStart := time.Now()
	err := sr.client.Put(input.SrcPath, input.DstPath, nil)

	return &runner.RcpOutput{
		RcpStart: rcpStart,
		RcpEnd:   time.Now(),
		Err:      err,
	}
}

// SyncGet copy file from remote server sync
func (sr *SSHRunner) SyncGet(input runner.RcpInput) *runner.RcpOutput {
	rcpStart := time.Now()
	err := sr.client.Get(input.DstPath, input.SrcPath, nil)

	return &runner.RcpOutput{
		RcpStart: rcpStart,
		RcpEnd:   time.Now(),
		Err:      err,
	}
}

// ConcurrentGet copy file to remote server concurrency
func (sr *SSHRunner) ConcurrentGet(input runner.RcpInput, outputChan chan *runner.ConcurrentRcpOutput, limitChan chan int, pool *pb.Pool) {
	var (
		err     error
		bar     *pb.ProgressBar
		dirSize int64
	)

	limitChan <- 1
	rcpStart := time.Now()

	dirSize, err = utils.RemoteDirSize(sr.client.sftpClient, input.SrcPath)
	if err != nil {
		goto GetResult
	}

	bar = utils.NewProgressBar(input.RcpHost, dirSize)
	pool.Add(bar)
	err = sr.client.Get(input.DstPath, input.SrcPath, bar)

GetResult:
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
	var bar *pb.ProgressBar
	rcpStart := time.Now()

	dirSize, err := utils.LocalDirSize(input.SrcPath)
	if err != nil {
		goto PutResult
	}

	bar = utils.NewProgressBar(input.RcpHost, dirSize)
	pool.Add(bar)
	err = sr.client.Put(input.SrcPath, input.DstPath, bar)

PutResult:
	outputChan <- &runner.ConcurrentRcpOutput{In: input, Out: &runner.RcpOutput{
		RcpStart: rcpStart,
		RcpEnd:   time.Now(),
		Err:      err,
	}}
	<-limitChan
}

func compositCommand(input runner.ExecInput) string {
	return fmt.Sprintf(`su - '%s' -c '%s'`, input.ExecUser, input.Command)
}
