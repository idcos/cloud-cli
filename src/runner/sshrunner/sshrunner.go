package sshrunner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runner"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// AuthTypePwd SSH login by password
	authTypePwd string = "PASSWORD"
	// AuthTypeKey SSH login by ssh-key
	authTypeKey string = "KEY"
)

var (
	// ErrInvalidAuthType error message for invalid auth type
	ErrInvalidAuthType = "Invalid auth type : %s\n"
)

// SSHRunner execute command by ssh
type SSHRunner struct {
	User       string
	Password   string
	SSHKeyPath string
	Host       string
	Port       int
	client     *ssh.Client
	session    *ssh.Session
}

func New(user, password, sshKeyPath, host string, port int) *SSHRunner {

	if port == 0 {
		port = 22
	}

	return &SSHRunner{
		User:       user,
		Password:   password,
		SSHKeyPath: sshKeyPath,
		Host:       host,
		Port:       port,
	}
}

// SyncExec execute command sync
func (sr *SSHRunner) SyncExec(input runner.Input) *runner.Output {
	var (
		err            error
		output         = &runner.Output{Status: runner.Fail}
		cmd            = compositCommand(input)
		stdout, stderr bytes.Buffer
		errChan        = make(chan error)
	)
	output.ExecStart = time.Now()

	// create session
	if err = sr.createSSHSession(); err != nil {
		goto SSHRunnerResult
	}
	defer sr.Close()

	sr.session.Stdout = &stdout
	sr.session.Stderr = &stderr

	go func(session *ssh.Session) {
		if err = session.Start(cmd); err != nil {
			errChan <- err
		}

		if err = session.Wait(); err != nil {
			errChan <- err
		}
		errChan <- nil
	}(sr.session)

	select {
	case err = <-errChan:
	case <-time.After(input.Timeout):
		err = fmt.Errorf("exec command(%s) on host(%s) TIMEOUT", input.Command, input.ExecHost)
		output.Status = runner.Timeout
	}

	output.StdOutput = string(stdout.Bytes())
	output.StdError = string(stderr.Bytes())

SSHRunnerResult:
	output.ExecEnd = time.Now()
	output.Err = err
	if output.Err == nil && output.StdError == "" {
		output.Status = runner.Success
	}
	return output
}

// ConcurrentExec execute command sync
func (sr *SSHRunner) ConcurrentExec(input runner.Input, outputChan chan *runner.ConcurrentOutput, limitChan chan int) {
	limitChan <- 1
	var output = sr.SyncExec(input)
	outputChan <- &runner.ConcurrentOutput{In: input, Out: output}
	<-limitChan
}

// Login login to remote server
func (sr *SSHRunner) Login(hostName, shell string) error {
	var err error

	// create session
	if err = sr.createSSHSession(); err != nil {
		return err
	}
	defer sr.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	// excute command
	sr.session.Stdout = os.Stdout
	sr.session.Stderr = os.Stderr
	sr.session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := sr.session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}
	if err := sr.session.Run(shell); err != nil {
		return err
	}
	return nil
}

func (sr *SSHRunner) Close() {
	if sr.session != nil {
		sr.session.Close()
	}

	if sr.client != nil {
		sr.client.Close()
	}
}

// createSSHSession create session for ssh use
func (sr *SSHRunner) createSSHSession() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		err          error
	)
	// get auth method
	auth, _ = authMethods(sr.Password, sr.SSHKeyPath)

	clientConfig = &ssh.ClientConfig{
		User:    sr.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", sr.Host, sr.Port)

	if sr.client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}

	// create session
	if sr.session, err = sr.client.NewSession(); err != nil {
		return err
	}

	return nil
}

// authMethods get auth methods
func authMethods(password, sshKeyPath string) ([]ssh.AuthMethod, error) {
	var (
		err         error
		authkey     []byte
		signer      ssh.Signer
		authMethods = make([]ssh.AuthMethod, 0)
	)
	authMethods = append(authMethods, ssh.Password(password))

	if authkey, err = ioutil.ReadFile(sshKeyPath); err != nil {
		return authMethods, err
	}

	if signer, err = ssh.ParsePrivateKey(authkey); err != nil {
		return authMethods, err
	}

	authMethods = append(authMethods, ssh.PublicKeys(signer))
	return authMethods, nil
}

func compositCommand(input runner.Input) string {
	return fmt.Sprintf(`su - '%s' -c '%s'`, input.ExecUser, input.Command)
}
