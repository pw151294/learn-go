package config

import (
	"encoding/json"
	"os"
)

var Cfg *AppConfig

type App struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Server struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Logging struct {
	Level string `json:"level"`
	Path  string `json:"path"`
}

type AppConfig struct {
	App     App     `json:"app"`
	Server  Server  `json:"server"`
	Logging Logging `json:"logging"`
}

func InitConfig(cfgPath string) error {
	file, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	cfg := &AppConfig{}
	if err = json.NewDecoder(file).Decode(cfg); err != nil {
		return err
	}
	Cfg = cfg
	return nil
}
