package runner

import "time"

type OutputStaus string

const (
	Fail    OutputStaus = "fail"
	Success OutputStaus = "success"
	Timeout OutputStaus = "timeout"
)

// ExecInput input format for runner interface
type ExecInput struct {
	// ExecUser username for exec command
	ExecUser string
	// ExecHost hostname for exec command
	ExecHost string
	// Command command for exec
	Command string
	// Timeout command exec timeout
	Timeout time.Duration
}

// ExecOutput output format for runner interface
type ExecOutput struct {
	// Status for exec result
	Status OutputStaus
	// StdError error output for exec result
	StdError string
	// StdOutput normal output for exec result
	StdOutput string
	// ExecStart start time when exec command
	ExecStart time.Time
	// ExecEnd end time when exec command
	ExecEnd time.Time
	// Err error info about exec command
	Err error
}

type RcpInput struct {
	// LocalPath local path for file or directory
	SrcPath string
	// RemotePath remote path for file or directory
	DstPath string
	// RcpHost remote host
	RcpHost string
	// RcpUser remote user
	RcpUser string
}

type RcpOutput struct {
	// RcpStart start time when exec command
	RcpStart time.Time
	// RcpEnd end time when exec command
	RcpEnd time.Time
	// Err error info about exec command
	Err error
}

// ConcurrentOutput output for concurrent exec command
type ConcurrentOutput struct {
	// In for confirm the Out is from which node
	In ExecInput
	// Out concurrent exec output
	Out *ExecOutput
}

// IRunner runner interface
type IRunner interface {
	// SyncExec exec command sync
	SyncExec(input ExecInput) *ExecOutput
	// ConcurrentExec exec command concurrency
	ConcurrentExec(input ExecInput, outputChan chan *ConcurrentOutput, limitChan chan int)
	// Login login to remote server
	Login(shell string) error
	// SyncPut copy file to remote server sync
	SyncPut(input RcpInput) *RcpOutput
	// SyncGet copy file from remote server sync
	SyncGet(input RcpInput) *RcpOutput
	// ConcurrentRcp copy file to remote server concurrency
}
