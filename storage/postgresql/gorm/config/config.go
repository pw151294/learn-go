package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DatabaseConfig struct {
	Drive    string `json:"drive"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Sslmode  string `json:"sslmode"`
}

func GetDataBaseConfig() *DatabaseConfig {
	return &DatabaseConfig{}
}

var dbConfig *DatabaseConfig

func InitConfig(cfgPath string) error {
	file, err := os.Open(filepath.Join(cfgPath, "config.json"))
	if err != nil {
		return err
	}
	defer file.Close()

	dbConfig = &DatabaseConfig{}
	if err = json.NewDecoder(file).Decode(dbConfig); err != nil {
		return err
	}
	return nil
}
