package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/pkg/sftp"
)

var (
	// ErrFileExisted error message when file exist
	ErrFileExisted = "target file %s is existed"
	// ErrMd5Check error message when check md5
	ErrMd5Check = "md5 is not same: Origin MD5(%s)\tCurrent MD5(%s)"
	// ErrLocalPathIsFile remote path is directory
	ErrLocalPathIsFile = "local path cannot be a file when remote path is directory"
)

// Md5File generate md5 string
func Md5File(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5.Sum(content)), nil
}

// ChkMd5Info check md5 info is same or not
func ChkMd5Info(filepath string, md5Str string) error {
	if newMd5Str, err := Md5File(filepath); err != nil {
		return err
	} else if newMd5Str != md5Str {
		return fmt.Errorf(ErrMd5Check, newMd5Str, md5Str)
	}

	return nil
}

// FileExist 判断文件是否存在
func FileExist(filepath string) bool {
	fi, err := os.Stat(filepath)

	return (err == nil || os.IsExist(err)) && !fi.IsDir()
}

// DirExist 判断文件夹是否存在
func DirExist(dirpath string) bool {
	fi, err := os.Stat(dirpath)

	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

// IsDir is directory or not
func IsDir(filepath string) bool {
	fi, err := os.Stat(filepath)

	return err == nil && fi.IsDir()
}

// PutFile put file to remote server
func PutFile(sftpClient *sftp.Client, localPath, remoteDir string, fileTransBuf int, bar *pb.ProgressBar) error {
	filename := path.Base(localPath)
	srcFile, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create remote dir
	if err := MkRemoteDirs(sftpClient, remoteDir); err != nil {
		return err
	}

	dstFile, err := sftpClient.Create(path.Join(remoteDir, filename))
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
		bar = NewProgressBar(localPath, fSize)
		bar.Start()
	}

	var i int64
	for {
		i++
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
func PutDir(sftpClient *sftp.Client, localDir, remoteDir string, fileTransBuf int, bar *pb.ProgressBar) error {

	return filepath.Walk(localDir, func(localPath string, info os.FileInfo, err error) error {
		relPath, err := filepath.Rel(localDir, localPath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// if the remote directory is existed, then omit create it
			if err := MkRemoteDirs(sftpClient, path.Join(remoteDir, relPath)); err != nil {
				return err
			}
			return nil
		} else {
			return PutFile(sftpClient, localPath, path.Join(remoteDir, path.Dir(relPath)), fileTransBuf, bar)
		}
	})
}

// IsRemoteDirExisted is remote dir existed
func IsRemoteDirExisted(sftpClient *sftp.Client, remoteDir string) bool {
	remoteFileInfo, err := sftpClient.Stat(remoteDir)
	// TODO error type is "not found file or directory"
	if err != nil {
		return false
	}

	return remoteFileInfo.IsDir()
}

// MkRemoteDirs create remote directories
func MkRemoteDirs(sftpClient *sftp.Client, remoteDir string) error {
	// create parent directory first
	var parentDir = path.Dir(remoteDir)
	if !IsRemoteDirExisted(sftpClient, remoteDir) {
		MkRemoteDirs(sftpClient, parentDir)
		return sftpClient.Mkdir(remoteDir)
	}
	return nil
}

// GetFile get file from remote server
func GetFile(sftpClient *sftp.Client, localPath, remoteFile string, fileTransBuf int, bar *pb.ProgressBar) error {

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

	var fSize int64
	if fi, err := srcFile.Stat(); err != nil {
		return err
	} else {
		fSize = fi.Size()
	}

	var bufSize = fileTransBuf
	buf := make([]byte, bufSize)
	if bar == nil {
		bar = NewProgressBar(localPath, fSize)
		bar.Start()
	}

	var i int64
	for {
		i++
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
func GetDir(sftpClient *sftp.Client, localPath, remoteDir string, fileTransBuf int, bar *pb.ProgressBar) error {
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
			if err = GetFile(sftpClient, path.Join(localPath, relRemotePath), w.Path(), fileTransBuf, bar); err != nil {
				return err
			}
		}
	}

	return nil
}

// LocalDirSize directory size with all file in it
func LocalDirSize(dirPath string) (int64, error) {
	if !IsDir(dirPath) {
		info, err := os.Stat(dirPath)
		if err != nil {
			return 0, err
		}
		return info.Size(), nil
	}

	var size int64
	err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})

	return size, err
}

// RemoteDirSize directory size with all file in it
func RemoteDirSize(sftpClient *sftp.Client, dirPath string) (int64, error) {
	if !IsRemoteDirExisted(sftpClient, dirPath) {
		info, err := sftpClient.Stat(dirPath)
		if err != nil {
			return 0, err
		}
		return info.Size(), nil
	}

	var size int64
	var err error

	w := sftpClient.Walk(dirPath)
	for w.Step() {
		if err = w.Err(); err != nil {
			return size, err
		}

		_, err = filepath.Rel(dirPath, w.Path())
		if err != nil {
			return size, err
		}
		if !w.Stat().IsDir() {
			size += w.Stat().Size()
		}
	}
	return size, err
}
