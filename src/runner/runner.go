package runner

import "time"

type RunnerInput struct {
	ExecUser string
	ExecHost string
	Command  string
}

type RunnerOutput struct {
	Status     string
	StatusCode int
	StdError   string
	StdOutput  string
	ExecStart  time.Time
	ExecEnd    time.Time
}

type Runner interface {
	SyncExec(input RunnerInput) (RunnerOutput, error)
}
