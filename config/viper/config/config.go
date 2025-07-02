package config

import (
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var ConfigManager *AgentConfigManager

type AgentConfigManager struct {
	config *AgentConfig
	sync.RWMutex
}

type Agent struct {
	AgentId  string `json:"agent_id"`
	Endpoint string `json:"endpoint"`
	Port     int    `json:"port"`
	IsHttps  bool   `json:"is_https"`
}

type Update struct {
	SelfUpdate            bool   `json:"self_update"`
	Endpoint              string `json:"endpoint"`
	Port                  int    `json:"port"`
	SelfUpdateScheduleDay int    `json:"self_update_schedule_day"`
	IsHttps               bool   `json:"is_https"`
}

type Path struct {
	DataPath   string `json:"data_path"`
	TempPath   string `json:"temp_path"`
	LogPath    string `json:"log_path"`
	PluginPath string `json:"plugin_path"`
}

type Plugins struct {
	FileTransferLimit int `json:"file_transfer_limit"`
}

type AgentConfig struct {
	Agent   Agent         `json:"agent"`
	Path    Path          `json:"path"`
	Update  Update        `json:"update"`
	Plugins Plugins       `json:"plugins"`
	Proxy   []interface{} `json:"proxy"`
}

func NewAgentConfigManager(config *AgentConfig) *AgentConfigManager {
	return &AgentConfigManager{
		config: config,
	}
}

func (cm *AgentConfigManager) SetConfig(config *AgentConfig) {
	cm.Lock()
	defer cm.Unlock()

	cm.config = config
}

func (cm *AgentConfigManager) GetConfig() *AgentConfig {
	cm.RLock()
	defer cm.RUnlock()

	return cm.config
}

func InitConfig(cfgPath string) error {
	// 读取配置
	file, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer file.Close()
	cfg := &AgentConfig{}
	if err = json.NewDecoder(file).Decode(cfg); err != nil {
		return err
	}
	ConfigManager = NewAgentConfigManager(cfg)

	// 解析出配置文件的信息
	fn := filepath.Base(cfgPath)
	ext := filepath.Ext(fn)
	cp := filepath.Dir(cfgPath)
	cn := strings.TrimSuffix(fn, ext)
	ct := ext[1:]

	// 使用viper实现配置监听
	viper.AddConfigPath(cp)
	viper.SetConfigName(cn)
	viper.SetConfigType(ct)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// 配置文件更新 重新读取配置
		zap.GetLogger().Info("config changed")
		file, err = os.Open(cfgPath)
		if err != nil {
			zap.GetLogger().Error("reload config failed", "message", err)
			return
		}
		defer file.Close()
		newCfg := &AgentConfig{}
		if err = json.NewDecoder(file).Decode(newCfg); err != nil {
			zap.GetLogger().Error("reload config failed", "message", err)
			return
		}
		ConfigManager.SetConfig(newCfg)
	})

	return nil
}
