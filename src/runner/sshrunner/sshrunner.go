package sshrunner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runner"
	"time"

	"golang.org/x/crypto/ssh"
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
func (sr *SSHRunner) SyncExec(input runner.Input) (*runner.Output, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
		output       = &runner.Output{Status: runner.Success}
	)

	// get auth method
	if auth, err = sr.authMethods(); err != nil {
		// TODO warn this error
	}

	clientConfig = &ssh.ClientConfig{
		User:    sr.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", sr.Host, sr.Port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	defer client.Close()

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	defer session.Close()

	// excute command
	cmd := compositCommand(input)
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	output.ExecStart = time.Now()
	if err = session.Start(cmd); err != nil {
		goto SSHResult
	}

	if err = session.Wait(); err != nil {
		goto SSHResult
	}

SSHResult:
	output.ExecEnd = time.Now()
	output.StdOutput = string(stdout.Bytes())
	output.StdError = string(stderr.Bytes())
	if err != nil || output.StdError != "" {
		output.Status = runner.Fail
	}
	return output, err
}

// authMethods get auth methods
func (sr *SSHRunner) authMethods() ([]ssh.AuthMethod, error) {
	var (
		err         error
		authkey     []byte
		signer      ssh.Signer
		authMethods = make([]ssh.AuthMethod, 0)
	)
	authMethods = append(authMethods, ssh.Password(sr.Password))

	if authkey, err = ioutil.ReadFile(sr.SSHKeyPath); err != nil {
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
