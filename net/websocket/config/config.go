package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
	"time"
)

var WskCfg *WskConfig

type WskConfig struct {
	ServerPort        int
	ServerIp          string
	HeartbeatInterval time.Duration
	RetryLimit        int
}

type Config struct {
	ServerPort        int    `toml:"serverPort"`
	ServerIp          string `toml:"serverIp"`
	HeartbeatInterval string `toml:"heartbeatInterval"`
	RetryLimit        int    `toml:"retryLimit"`
}

type WskConfigOptions func(*WskConfig)

func NewWskConfig(opts ...WskConfigOptions) *WskConfig {
	cfg := &WskConfig{
		ServerPort:        8080,
		ServerIp:          "127.0.0.1",
		HeartbeatInterval: 10 * time.Second,
		RetryLimit:        3,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(cfg)
		}
	}

	return cfg
}

func WithServerPort(port int) WskConfigOptions {
	return func(config *WskConfig) {
		if port > 0 {
			config.ServerPort = port
		}
	}
}

func WithServerIp(ip string) WskConfigOptions {
	return func(config *WskConfig) {
		if ip != "" {
			config.ServerIp = ip
		}
	}
}

func WithRetryLimit(limit int) WskConfigOptions {
	return func(config *WskConfig) {
		if limit > 0 {
			config.RetryLimit = limit
		}
	}
}

func WithHeartbeatInterval(interval string) WskConfigOptions {
	return func(config *WskConfig) {
		if d, err := time.ParseDuration(interval); err == nil {
			config.HeartbeatInterval = d
		}
	}
}

func InitConfig(cfgPath string) error {
	file, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	cfg := &Config{}
	if err = toml.NewDecoder(file).Decode(cfg); err != nil {
		return err
	}
	WskCfg = NewWskConfig(
		WithServerPort(cfg.ServerPort),
		WithServerIp(cfg.ServerIp),
		WithHeartbeatInterval(cfg.HeartbeatInterval),
		WithRetryLimit(cfg.RetryLimit))
	return nil
}
