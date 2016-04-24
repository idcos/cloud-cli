package runner

import "time"

// Input input format for runner interface
type Input struct {
	// ExecUser username for exec command
	ExecUser string
	// ExecHost hostname for exec command
	ExecHost string
	// Command command for exec
	Command string
}

// Output output format for runner interface
type Output struct {
	// Status for exec result
	Status string
	// StatusCode code for exec result
	StatusCode int
	// StdError error output for exec result
	StdError string
	// StdOutput normal output for exec result
	StdOutput string
	// ExecStart start time when exec command
	ExecStart time.Time
	// ExecEnd end time when exec command
	ExecEnd time.Time
}

// IRunner runner interface
type IRunner interface {
	SyncExec(input Input) (Output, error)
}
