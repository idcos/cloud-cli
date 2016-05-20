package runner

import "time"

type OutputStaus string

const (
	Fail    OutputStaus = "fail"
	Success OutputStaus = "success"
	Timeout OutputStaus = "timeout"
)

// Input input format for runner interface
type Input struct {
	// ExecUser username for exec command
	ExecUser string
	// ExecHost hostname for exec command
	ExecHost string
	// Command command for exec
	Command string
	// Timeout command exec timeout
	Timeout time.Duration
}

// Output output format for runner interface
type Output struct {
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

//
type ConcurrentOutput struct {
	In  Input
	Out *Output
}

// IRunner runner interface
type IRunner interface {
	// exec command sync
	SyncExec(input Input) *Output
	// exec command concurrency
	ConcurrentExec(input Input, outputChan chan *ConcurrentOutput, limitChan chan int)
}
