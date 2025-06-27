package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

var MinioCfg *MinioConfig

type MinioConfig struct {
	Endpoint        string `toml:"endpoint"`
	AccessKeyID     string `toml:"accessKeyID"`
	SecretAccessKey string `toml:"secretAccessKey"`
	UseSSL          bool   `toml:"useSSL"`
}

func InitConfig(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	MinioCfg = &MinioConfig{}
	err = toml.NewDecoder(file).Decode(MinioCfg)
	return err
}
