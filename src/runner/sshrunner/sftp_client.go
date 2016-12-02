package sshrunner

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"utils"

	"github.com/pkg/sftp"
	pb "gopkg.in/cheggaaa/pb.v1"
)

var (
	// ErrLocalPathIsFile remote path is directory
	ErrLocalPathIsFile = "local path cannot be a file when remote path is directory"
)

type SFTPClient struct {
	FileTransBuf int
	conn         *SSHClient
	client       *sftp.Client
}

func NewSFTPClient(sc *SSHClient, fileTransBuf int) *SFTPClient {
	return &SFTPClient{
		FileTransBuf: fileTransBuf,
		conn:         sc,
	}
}

func (sf *SFTPClient) Close() {
	if sf.client != nil {
		sf.client.Close()
	}
}

// Put transfer file/directory to remote server
func (sf *SFTPClient) Put(localPath, remotePath string, bar *pb.ProgressBar) error {
	var (
		err           error
		localFileInfo os.FileInfo
	)

	// create sftp client
	if err = sf.createClient(); err != nil {
		return err
	}
	defer sf.Close()

	localFileInfo, err = os.Stat(localPath)
	if err != nil {
		return err
	}

	if localFileInfo.IsDir() { // localPath is directory
		if string(localPath[len(localPath)-1]) == "/" {
			remotePath = path.Join(remotePath, path.Base(localPath))
		}
		return sf.putDir(localPath, remotePath, sf.FileTransBuf, bar)
	} else { // localPath is file
		return sf.putFile(localPath, remotePath, sf.FileTransBuf, bar)
	}
}

// Get transfer file/directory from remote server
func (sf *SFTPClient) Get(localPath, remotePath string, bar *pb.ProgressBar) error {
	var (
		err            error
		remoteFileInfo os.FileInfo
	)

	// create sftp client
	if err = sf.createClient(); err != nil {
		return err
	}
	defer sf.Close()

	if remoteFileInfo, err = sf.client.Stat(remotePath); err != nil {
		return err
	}

	if remoteFileInfo.IsDir() {
		localPath = path.Join(localPath, sf.conn.Host) // create dir by hostname
		if string(remotePath[len(remotePath)-1]) == "/" {
			localPath = path.Join(localPath, path.Base(remotePath))
		}
		os.MkdirAll(localPath, os.ModePerm)
		return sf.getDir(localPath, remotePath, sf.FileTransBuf, bar)
	} else {
		return sf.getFile(localPath, remotePath, sf.FileTransBuf, bar)
	}

	return err
}

// RemotePathSize remote path size with all file in it
func (sf *SFTPClient) RemotePathSize(remotePath string) (int64, error) {
	var size int64
	var err error

	// create sftp client
	if err = sf.createClient(); err != nil {
		return size, err
	}
	defer sf.Close()

	if !sf.isRemoteDirExisted(remotePath) {
		info, err := sf.client.Stat(remotePath)
		if err != nil {
			return size, err
		}
		return info.Size(), nil
	}

	w := sf.client.Walk(remotePath)
	for w.Step() {
		if err = w.Err(); err != nil {
			return size, err
		}

		_, err = filepath.Rel(remotePath, w.Path())
		if err != nil {
			return size, err
		}
		if !w.Stat().IsDir() {
			size += w.Stat().Size()
		}
	}
	return size, err
}

// IsRemoteDirExisted is remote dir existed
func (sf *SFTPClient) isRemoteDirExisted(remoteDir string) bool {
	remoteFileInfo, err := sf.client.Stat(remoteDir)
	// TODO error type is "not found file or directory"
	if err != nil {
		return false
	}

	return remoteFileInfo.IsDir()
}

// MkRemoteDirs create remote directories
func (sf *SFTPClient) mkRemoteDirs(remoteDir string) error {
	// create parent directory first
	var parentDir = path.Dir(remoteDir)
	if !sf.isRemoteDirExisted(remoteDir) {
		sf.mkRemoteDirs(parentDir)
		return sf.client.Mkdir(remoteDir)
	}
	return nil
}

// GetFile get file from remote server
func (sf *SFTPClient) getFile(localPath, remoteFile string, fileTransBuf int, bar *pb.ProgressBar) error {

	srcFile, err := sf.client.Open(remoteFile)
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

	var fSize int64
	if fi, err := srcFile.Stat(); err != nil {
		return err
	} else {
		fSize = fi.Size()
	}

	var bufSize = fileTransBuf
	buf := make([]byte, bufSize)
	if bar == nil {
		bar = utils.NewProgressBar(localPath, fSize)
		bar.Start()
	}

	for {
		nread, _ := srcFile.Read(buf)
		if nread == 0 {
			break
		}
		dstFile.Write(buf[:nread])

		bar.Add(nread)
	}
	if bar.Total == bar.Get() {
		bar.Finish()
	}

	return err
}

// GetDir get directory from remote server
func (sf *SFTPClient) getDir(localPath, remoteDir string, fileTransBuf int, bar *pb.ProgressBar) error {
	localFileInfo, err := os.Stat(localPath)
	// remotepath is directory, localPath existed and be a file, cause error
	if err == nil && !localFileInfo.IsDir() {
		return fmt.Errorf(ErrLocalPathIsFile)
	}

	w := sf.client.Walk(remoteDir)
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
			if err = sf.getFile(path.Join(localPath, relRemotePath), w.Path(), fileTransBuf, bar); err != nil {
				return err
			}
		}
	}

	return nil
}

// PutFile put file to remote server
func (sf *SFTPClient) putFile(localPath, remoteDir string, fileTransBuf int, bar *pb.ProgressBar) error {
	filename := path.Base(localPath)
	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create remote dir
	if err := sf.mkRemoteDirs(remoteDir); err != nil {
		return err
	}

	dstFile, err := sf.client.Create(path.Join(remoteDir, filename))
	if err != nil {
		return err
	}
	defer dstFile.Close()

	var fSize int64
	if fi, err := srcFile.Stat(); err != nil {
		return err
	} else {
		fSize = fi.Size()
	}

	var bufSize = fileTransBuf
	buf := make([]byte, bufSize)
	if bar == nil {
		bar = utils.NewProgressBar(localPath, fSize)
		bar.Start()
	}

	for {
		nread, _ := srcFile.Read(buf)
		if nread == 0 {
			break
		}
		dstFile.Write(buf[:nread])

		bar.Add(nread)
	}
	if bar.Total == bar.Get() {
		bar.Finish()
	}

	return nil
}

// PutDir put directories to remote server
func (sf *SFTPClient) putDir(localDir, remoteDir string, fileTransBuf int, bar *pb.ProgressBar) error {

	return filepath.Walk(localDir, func(localPath string, info os.FileInfo, err error) error {
		relPath, err := filepath.Rel(localDir, localPath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// if the remote directory is existed, then omit create it
			if err := sf.mkRemoteDirs(path.Join(remoteDir, relPath)); err != nil {
				return err
			}
			return nil
		} else {
			return sf.putFile(localPath, path.Join(remoteDir, path.Dir(relPath)), fileTransBuf, bar)
		}
	})
}

// createClient create ssh client
func (sf *SFTPClient) createClient() error {
	var err error

	if err = sf.conn.createClient(); err != nil {
		return err
	}

	sf.client, err = sftp.NewClient(sf.conn.client)
	return err
}
