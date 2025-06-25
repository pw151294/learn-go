package config

import (
	"github.com/pelletier/go-toml/v2"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"os"
)

var Cfg *CmdbAgentConfig

func InitConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		zap.GetLogger().Error("read  config file failed", "file", configPath)
		return err
	}
	defer file.Close()

	cfg := &CmdbAgentConfig{}
	if err = toml.NewDecoder(file).Decode(cfg); err != nil {
		zap.GetLogger().Error("load cmdb agent config failed", "file", configPath)

	}
	Cfg = cfg

	zap.GetLogger().Info("init cmdb agent config success!")
	return nil
}
