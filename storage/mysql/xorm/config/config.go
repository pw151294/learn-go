package config

import (
	"encoding/json"
	"io"
	"os"
)

var dbCfg *DatabaseConfig

type DatabaseConfig struct {
	Drive       string `json:"drive"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
	Database    string `json:"database"`
	MaxOpenCons int    `json:"max_open_cons"`
	MaxIdleCons int    `json:"max_idle_cons"`
	ShowSql     bool   `json:"show_sql"`
	LogPath     string `json:"log_path"`
}

func Get() *DatabaseConfig {
	return dbCfg
}

func InitConfig(cfgPath string) error {
	file, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	cfg := &DatabaseConfig{}
	if err = json.NewDecoder(io.Reader(file)).Decode(cfg); err != nil {
		return err
	}
	dbCfg = cfg
	return nil
}
