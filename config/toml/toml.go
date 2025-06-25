package main

import (
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
	"path/filepath"
)

const configFile = "config.toml"

type OneAgent struct {
	ID           string `toml:"id"`
	Name         string `toml:"name"`
	LogPath      string `toml:"log_path"`
	MasterServer string `toml:"master_server"`
}

type OneMaster struct {
	CmdbServer    string `toml:"cmdb_server"`
	Apikey        string `toml:"apikey"`
	ApiSecret     string `toml:"api_secret"`
	LicenseBindIp string `toml:"license_bind_ip"`
}

type CmdbAgent struct {
	OneMaster OneMaster `toml:"OneMaster"`
	OneAgent  OneAgent  `toml:"OneAgent"`
}

func main() {
	fn := filepath.Join("config/toml", configFile)
	bytes, err := os.ReadFile(fn)
	if err != nil {
		log.Fatalf("read toml file failed: %v", err)
	}

	agent := CmdbAgent{}
	if err := toml.Unmarshal(bytes, &agent); err != nil {
		log.Fatalf("load toml file failed: %v", err)
	}

	bytes, _ = json.Marshal(agent)
	fmt.Println(string(bytes))
}
