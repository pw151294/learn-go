package config

type CmdbAgentConfig struct {
	*OneAgent  `toml:"oneagent"`
	*OneMaster `toml:"onemaster"`
}

type OneAgent struct {
	ID           string `toml:"id"`
	Name         string `toml:"name"`
	LogPath      string `toml:"log_path"`
	MasterServer string `toml:"master_server"`
	UUID         string `toml:"uuid"`
}

type OneMaster struct {
	CmdbServer    string `toml:"cmdb_server"`
	Apikey        string `toml:"api_key"`
	ApiSecret     string `toml:"api_secret"`
	LicenseBindIp string `toml:"license_bind_ip"`
}
