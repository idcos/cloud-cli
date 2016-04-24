package util

import (
	"os"
	"testing"
)

func TestMd5File(t *testing.T) {
	var filePath = "testfile"
	f, err := os.Create(filePath)
	if err != nil {
		t.Errorf("prepare test failed: %v", err)
	}

	defer os.Remove(filePath)
	f.WriteString("testtesttset")
	f.Close()

	if _, err = Md5File(filePath); err != nil {
		t.Errorf("test Md5File error: %v", err)
	}
}

func TestChkMd5Info(t *testing.T) {
	var filePath = "testfile"
	f, err := os.Create(filePath)
	if err != nil {
		t.Errorf("prepare test failed: %v", err)
	}

	defer os.Remove(filePath)
	f.WriteString("testtesttset")
	f.Close()

	md5Str, _ := Md5File(filePath)
	if err = ChkMd5Info(filePath, md5Str); err != nil {
		t.Errorf("test ChkMd5Info error: %v", err)
	}
}

func TestFileExist(t *testing.T) {
	var filePath = "testfile"
	f, err := os.Create(filePath)
	if err != nil {
		t.Errorf("prepare test failed: %v", err)
	}

	defer os.Remove(filePath)
	f.WriteString("testtesttset")
	f.Close()

	if FileExist(filePath) == false {
		t.Errorf("test file exist error")
	}

	if FileExist("a/b/c/d") == true {
		t.Errorf("test file not exist error")
	}

	if FileExist("/tmp") == true {
		t.Errorf("test file exist error: file is dir")
	}
}

func TestDirExist(t *testing.T) {
	var dirPath = "tmp"
	if err := os.MkdirAll(dirPath, 0666); err != nil {
		t.Errorf("prepare test error: %v", err)
		return
	}
	defer os.RemoveAll(dirPath)

	if DirExist(dirPath) == false {
		t.Errorf("test dir exist error")
	}

	if DirExist("tmp/a/b/c/d/e") == true {
		t.Errorf("test dir not exist error")
	}
}
