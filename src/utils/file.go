package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	// ErrFileExisted error message when file exist
	ErrFileExisted = "target file %s is existed"
	// ErrMd5Check error message when check md5
	ErrMd5Check = "md5 is not same: Origin MD5(%s)\tCurrent MD5(%s)"
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

// LocalPathSize directory size with all file in it
func LocalPathSize(dirPath string) (int64, error) {
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
