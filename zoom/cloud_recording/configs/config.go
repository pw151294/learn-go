package configs

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Server    ServerConfig    `toml:"server"`
	MinIO     MinIOConfig     `toml:"minio"`
	ES        ESConfig        `toml:"es"`
	CDN       CDNConfig       `toml:"cdn"`
	Auth      AuthConfig      `toml:"auth"`
	Lifecycle LifecycleConfig `toml:"lifecycle"`
}

type ServerConfig struct {
	Port int    `toml:"port"`
	Host string `toml:"host"`
}

type MinIOConfig struct {
	Endpoint        string `toml:"endpoint"`
	AccessKeyID     string `toml:"access_key_id"`
	SecretAccessKey string `toml:"secret_access_key"`
	UseSSL          bool   `toml:"use_ssl"`
	HotBucket       string `toml:"hot_bucket"`  // rec-hot
	WarmBucket      string `toml:"warm_bucket"` // rec-warm
	ColdBucket      string `toml:"cold_bucket"` // rec-cold
}

type ESConfig struct {
	URL       string `toml:"url"`
	IndexName string `toml:"index_name"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
}

type CDNConfig struct {
	Enabled bool   `toml:"enabled"`
	BaseURL string `toml:"base_url"`
}

type AuthConfig struct {
	SignSecret    string `toml:"sign_secret"`
	DefaultExpiry int    `toml:"default_expiry"` // 秒，默认 3600
}

type LifecycleConfig struct {
	HotDays       int `toml:"hot_days"`       // 热存储天数，默认 7
	WarmDays      int `toml:"warm_days"`      // 温存储天数，默认 30
	RetentionDays int `toml:"retention_days"` // 保留天数，默认 365
	ScanInterval  int `toml:"scan_interval"`  // 扫描间隔（秒），默认 3600
}

// Load 从文件加载配置
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	cfg := &Config{}
	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	return cfg, nil
}
