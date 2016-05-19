package config

// Loader 定义统一的配置加载接口
type Loader interface {
	Load() (*Config, error)
	Save(*Config) error
}

// Config config 数据结构体
type Config struct {
	Main struct {
		Sync bool `ini:"sync"`
	}
	Logger struct {
		Level   string `ini:"level"`
		LogFile string `ini:"logFile"`
		LogType string `ini:"logType"`
	}
	DataSource struct {
		Type string `ini:"type"`
		Conn string `ini:"conn"`
	}
}
