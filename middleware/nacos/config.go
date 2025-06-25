package main

import (
	"encoding/json"
	"os"
)

var (
	cfg *NacosConfig
)

type NacosConfig struct {
	AppName     string `json:"app_name"`
	IpAddr      string `json:"ip_addr"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ContextPath string `json:"context_path"`
	Port        uint64 `json:"port"`
	Scheme      string `json:"scheme"`
	NamespaceId string `json:"namespace_id"`
	LogPath     string `json:"log_path"`
	LogLevel    string `json:"log_level"`
}

func InitNacosConfig(cfgPath string) error {
	file, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	nacosCfg := &NacosConfig{}
	if err = json.NewDecoder(file).Decode(nacosCfg); err != nil {
		return err
	}
	cfg = nacosCfg
	return nil
}
