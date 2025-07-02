package config

import (
	"encoding/json"
	"os"
)

var Cfg *RedisConfig

type RedisConfig struct {
	// 不是哨兵，多地址(逗号分隔)则为集群，否则为单节点
	Addrs    string `json:"addrs"`
	Username string `json:"username"`
	Password string `json:"password"`
	Db       int    `json:"db"`
	UseTLS   bool
	// 哨兵则该字段不能为空
	MasterName       string `json:"masterName"`
	RedisType        int    `json:"redisType"`
	SentinelUsername string `json:"sentinelUsername"`
	SentinelPassword string `json:"sentinelPassword"`
}

func InitConfig(cfgPath string) error {
	// 读取配置
	file, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	cfg := &RedisConfig{}
	if err = json.NewDecoder(file).Decode(cfg); err != nil {
		return err
	}
	Cfg = cfg
	return nil
}
