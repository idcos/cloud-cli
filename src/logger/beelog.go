package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"utils"

	"github.com/astaxie/beego/logs"
)

const (
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

// NewFileLogger output to file
func NewFileLogger(logFilePath, level string) *logs.BeeLogger {
	filename, _ := utils.ConvertHomeDir(logFilePath)

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
