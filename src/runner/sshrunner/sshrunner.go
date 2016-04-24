package sshrunner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runner"
	"strings"
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
	LoginType  string
	Host       string
	Port       int
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
		output       = &runner.Output{}
	)

	// get auth method
	if auth, err = sr.authMethods(); err != nil {
		return nil, err
	}

	clientConfig = &ssh.ClientConfig{
		User: sr.User,
		Auth: auth,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%s", sr.Host, sr.Port)

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
		return nil, err
	}

	if err = session.Wait(); err != nil {
		return nil, err
	}
	output.ExecEnd = time.Now()
	output.StdOutput = string(stdout.Bytes())
	output.StdError = string(stderr.Bytes())

	return output, err
}

// authMethods get auth methods
func (sr *SSHRunner) authMethods() ([]ssh.AuthMethod, error) {
	if strings.ToUpper(sr.LoginType) == authTypePwd {
		return []ssh.AuthMethod{ssh.Password(sr.Password)}, nil
	}

	if strings.ToUpper(sr.LoginType) == authTypeKey {
		var (
			err     error
			authkey []byte
			signer  ssh.Signer
		)

		if authkey, err = ioutil.ReadFile(sr.SSHKeyPath); err != nil {
			return nil, err
		}

		if signer, err = ssh.ParsePrivateKey(authkey); err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
	}

	return nil, fmt.Errorf(ErrInvalidAuthType, sr.LoginType)
}

func compositCommand(input runner.Input) string {
	return fmt.Sprintf(`su - '%s' -c '%s'`, input.ExecUser, input.Command)
}
