package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/astaxie/beego/logs"
)

const (
	// HomeDirFlag 当前用户家目录标识符
	HomeDirFlag = "~"

	// FileLog output log to file
	FileLog = "file"
	// ConsoleLog output log to console
	ConsoleLog = "console"
)

func selectLevel(level string) uint {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return 7
	case "warn":
		return 4
	case "error":
		return 3
	default:
		return 6 // default level: info
	}
}

// 将~转化为用户家目录
func rel2Abs(raw string) (string, error) {
	raw = strings.TrimSpace(raw)

	if !strings.HasPrefix(raw, HomeDirFlag) {
		return raw, nil
	}
	user, err := user.Current()
	if err != nil {
		return raw, err
	}
	return strings.Replace(raw, HomeDirFlag, user.HomeDir, 1), nil
}

// NewFileLogger output to file
func NewFileLogger(logFilePath, level string) *logs.BeeLogger {
	filename := strings.TrimSpace(logFilePath)
	if strings.HasPrefix(filename, HomeDirFlag) {
		filename, _ = rel2Abs(filename)
	}

	var logConf struct {
		FileName string `json:"filename"`
		Level    uint   `json:"level"`
	}
	logConf.FileName = filename
	logConf.Level = selectLevel(level)

	if err := os.MkdirAll(path.Dir(filename), os.ModePerm); err != nil {
		fmt.Printf("MkdirAll err: %s\n", err)
	}

	log := logs.NewLogger(1000)
	log.EnableFuncCallDepth(true) // 输出文件名和行号
	log.SetLogFuncCallDepth(3)

	logData, _ := json.Marshal(logConf)
	if err := log.SetLogger("file", string(logData)); err != nil {
		fmt.Printf("SetLogger err: %s\n", err)
	}

	// 尝试重置日志文件权限为0666
	os.Chmod(filename, 0666) // 不处理error

	return log
}

// NewConsoleLogger output to terminal
func NewConsoleLogger(level string) *logs.BeeLogger {
	var logConf struct {
		Level uint `json:"level"`
	}
	logConf.Level = selectLevel(level)

	log := logs.NewLogger(1000)
	log.EnableFuncCallDepth(true) // 输出文件名和行号
	log.SetLogFuncCallDepth(3)

	logData, _ := json.Marshal(logConf)
	if err := log.SetLogger("console", string(logData)); err != nil {
		fmt.Printf("SetLogger err: %s\n", err)
	}

	return log
}
