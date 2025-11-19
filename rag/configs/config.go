package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

type LoggerConfig struct {
	Level        string `toml:"level"`
	Format       string `toml:"format"`
	Output       string `toml:"output"`
	LogDir       string `toml:"log_dir"`
	LogFile      string `toml:"log_file"`
	MaxSize      int    `toml:"max_size"`
	MaxBackups   int    `toml:"max_backups"`
	MaxAge       int    `toml:"max_age"`
	Compress     bool   `toml:"compress"`
	EnableCaller bool   `toml:"enable_caller"`
	CallerSkip   int    `toml:"caller_skip"`
}

func (c *LoggerConfig) GetLogFilePath() string {
	return filepath.Join(c.LogDir, c.LogFile)
}

func (c *LoggerConfig) EnsureLogDir() error {
	return os.MkdirAll(c.LogDir, 0755)
}

type MySQLConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	DBName    string `toml:"dbname"`
	Charset   string
	ParseTime bool `toml:"parseTime"`
	Loc       string
}

type LLMConfig struct {
	ApiKey  string `toml:"api_key"`
	BaseURL string `toml:"base_url"`
}

type WeaviateConfig struct {
	Host   string `toml:"host"`
	Scheme string `toml:"scheme"`
	ApiKey string `toml:"api_key"`
}

type Config struct {
	MySQL    MySQLConfig    `toml:"mysql"`
	LLM      LLMConfig      `toml:"llm"`
	Log      LoggerConfig   `toml:"log"`
	Weaviate WeaviateConfig `toml:"weaviate"`
}

var (
	config     *Config
	configOnce sync.Once
)

// InitConfig 读取配置文件，只初始化一次
func InitConfig(configPath string) error {
	var err error
	configOnce.Do(func() {
		var c Config
		if _, e := toml.DecodeFile(configPath, &c); e != nil {
			err = fmt.Errorf("failed to decode config: %w", e)
			return
		}
		config = &c
	})
	return err
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return config
}
