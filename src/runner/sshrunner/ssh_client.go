package sshrunner

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runner"
	"time"

	"github.com/pkg/sftp"

	"path"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	ErrLocalPathIsFile = "local path cannot be a file when remote path is directory"
)

type SSHClient struct {
	User       string
	Password   string
	SSHKeyPath string
	Host       string
	Port       int
	client     *ssh.Client
	session    *ssh.Session
	sftpClient *sftp.Client
}

func NewSSHClient(user, password, sshKeyPath, host string, port int) *SSHClient {
	if port == 0 {
		port = 22
	}

	return &SSHClient{
		User:       user,
		Password:   password,
		SSHKeyPath: sshKeyPath,
		Host:       host,
		Port:       port,
	}
}

// Close release resources
func (sc *SSHClient) Close() {
	if sc.session != nil {
		sc.session.Close()
	}

	if sc.client != nil {
		sc.client.Close()
	}
}

// ExecNointeractiveCmd exec command without interactive
func (sc *SSHClient) ExecNointeractiveCmd(cmd string, timeout time.Duration) (status runner.OutputStaus, stdout, stderr *bytes.Buffer, err error) {
	status = runner.Fail
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	var errChan = make(chan error)

	// create session
	if err = sc.createSession(); err != nil {
		status = runner.Timeout
		return
	}
	defer sc.Close()

	sc.session.Stdout = stdout
	sc.session.Stderr = stderr

	go func(session *ssh.Session) {
		if err = session.Start(cmd); err != nil {
			errChan <- err
		}

		if err = session.Wait(); err != nil {
			errChan <- err
		}
		errChan <- nil
	}(sc.session)

	select {
	case err = <-errChan:
	case <-time.After(timeout):
		err = fmt.Errorf("exec command(%s) on host(%s) TIMEOUT", cmd, sc.Host)
		status = runner.Timeout
	}

	if err == nil {
		status = runner.Success
	}

	return
}

// ExecInteractiveCmd exec command with interactive
func (sc *SSHClient) ExecInteractiveCmd(cmd string) error {
	var err error

	// create session
	if err = sc.createSession(); err != nil {
		return err
	}
	defer sc.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	// excute command
	sc.session.Stdout = os.Stdout
	sc.session.Stderr = os.Stderr
	sc.session.Stdin = os.Stdin

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
	if err := sc.session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}
	if err := sc.session.Run(cmd); err != nil {
		return err
	}
	return nil
}

// Put transfer file/directory to remote server
func (sc *SSHClient) Put(localPath, remotePath string) error {
	var (
		err           error
		localFileInfo os.FileInfo
	)

	// create client
	if err = sc.createClient(); err != nil {
		return err
	}
	sc.sftpClient, err = sftp.NewClient(sc.client)
	if err != nil {
		return err
	}
	defer sc.sftpClient.Close()

	localFileInfo, err = os.Stat(localPath)
	if err != nil {
		return err
	}

	if localFileInfo.IsDir() { // localPath is directory
		return putDir(sc.sftpClient, localPath, remotePath)
	} else { // localPath is file
		return putFile(sc.sftpClient, localPath, remotePath)
	}
}

// Get transfer file/directory from remote server
func (sc *SSHClient) Get(localPath, remotePath string) error {
	var (
		err            error
		remoteFileInfo os.FileInfo
	)

	// create client
	if err = sc.createClient(); err != nil {
		return err
	}
	sc.sftpClient, err = sftp.NewClient(sc.client)
	if err != nil {
		return err
	}
	defer sc.sftpClient.Close()

	if remoteFileInfo, err = sc.sftpClient.Stat(remotePath); err != nil {
		return err
	}

	if remoteFileInfo.IsDir() {
		return getDir(sc.sftpClient, localPath, remotePath)
	} else {
		return getFile(sc.sftpClient, localPath, remotePath)
	}

	return err
}

// createClient create ssh client
func (sc *SSHClient) createClient() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		err          error
	)
	// get auth method
	auth, _ = authMethods(sc.Password, sc.SSHKeyPath)

	clientConfig = &ssh.ClientConfig{
		User:    sc.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", sc.Host, sc.Port)

	if sc.client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	return nil
}

// createSession create session for ssh use
func (sc *SSHClient) createSession() error {
	var err error

	// create client
	if err = sc.createClient(); err != nil {
		return err
	}

	// create session
	if sc.session, err = sc.client.NewSession(); err != nil {
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

func putFile(sftpClient *sftp.Client, localPath, remoteDir string) error {
	filename := path.Base(localPath)
	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create remote dir
	if err := mkRemoteDirs(sftpClient, remoteDir); err != nil {
		return err
	}

	dstFile, err := sftpClient.Create(path.Join(remoteDir, filename))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		nread, _ := srcFile.Read(buf)
		if nread == 0 {
			break
		}
		dstFile.Write(buf[:nread])
	}

	return nil
}

func putDir(sftpClient *sftp.Client, localDir, remoteDir string) error {

	return filepath.Walk(localDir, func(localPath string, info os.FileInfo, err error) error {
		relPath, err := filepath.Rel(localDir, localPath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// if the remote directory is existed, then omit create it
			if err := mkRemoteDirs(sftpClient, path.Join(remoteDir, relPath)); err != nil {
				return err
			}
			return nil
		} else {
			return putFile(sftpClient, localPath, path.Join(remoteDir, path.Dir(relPath)))
		}
	})
}

func isRemoteDirExisted(sftpClient *sftp.Client, remoteDir string) bool {
	remoteFileInfo, err := sftpClient.Stat(remoteDir)
	// TODO error type is "not found file or directory"
	if err != nil {
		return false
	}

	return remoteFileInfo.IsDir()
}

func mkRemoteDirs(sftpClient *sftp.Client, remoteDir string) error {
	// create parent directory first
	var parentDir = path.Dir(remoteDir)
	if !isRemoteDirExisted(sftpClient, remoteDir) {
		mkRemoteDirs(sftpClient, parentDir)
		return sftpClient.Mkdir(remoteDir)
	}
	return nil
}

func getFile(sftpClient *sftp.Client, localPath, remoteFile string) error {

	srcFile, err := sftpClient.Open(remoteFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// localPath is directory, then localFile's name == remoteFile's name
	localFileInfo, err := os.Stat(localPath)
	if err == nil && localFileInfo.IsDir() {
		localPath = path.Join(localPath, path.Base(remoteFile))
	}

	dstFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = srcFile.WriteTo(dstFile)
	return err
}

func getDir(sftpClient *sftp.Client, localPath, remoteDir string) error {
	localFileInfo, err := os.Stat(localPath)
	// remotepath is directory, localPath existed and be a file, cause error
	if err == nil && !localFileInfo.IsDir() {
		return fmt.Errorf(ErrLocalPathIsFile)
	}

	w := sftpClient.Walk(remoteDir)
	for w.Step() {
		if err = w.Err(); err != nil {
			return err
		}

		relRemotePath, err := filepath.Rel(remoteDir, w.Path())
		if err != nil {
			return err
		}
		if w.Stat().IsDir() {
			if err = os.MkdirAll(path.Join(localPath, relRemotePath), os.ModePerm); err != nil {
				return err
			}
		} else {
			if err = getFile(sftpClient, path.Join(localPath, relRemotePath), w.Path()); err != nil {
				return err
			}
		}
	}

	return nil
}
