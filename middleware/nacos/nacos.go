package main

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"log"
)

const (
	SrvcGroup = "DEFAULT_GROUP"
	ClstrName = "DEFAULT"
	dataId    = "config.toml"
	cfgGroup  = "auto-discovery"
)

var (
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
)

func InitNacos() {
	CreateClient()
	InitConfig()
	RegisterService()
}

func CreateClient() {
	// 创建client config
	clientCfg := *constant.NewClientConfig(
		constant.WithNamespaceId(cfg.NamespaceId),
		constant.WithUsername(cfg.Username),
		constant.WithPassword(cfg.Password),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(cfg.LogPath))
	// 创建server configs
	serverCfgs := []constant.ServerConfig{
		*constant.NewServerConfig(
			cfg.IpAddr,
			cfg.Port,
			constant.WithScheme(cfg.Scheme),
			constant.WithContextPath(cfg.ContextPath)),
	}

	var err error
	// 创建服务发现和动态配置的客户端
	if namingClient, err = clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &clientCfg,
		ServerConfigs: serverCfgs,
	}); err != nil {
		log.Fatal(err)
	}
	if configClient, err = clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &clientCfg,
		ServerConfigs: serverCfgs,
	}); err != nil {
		log.Fatal(err)
	}
}

func InitConfig() {
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  cfgGroup,
	})
	if err != nil {
		log.Fatal(err)
	}
	if content == "" {
		zap.GetLogger().Warn("config content is empty", "dataId", dataId, "group", cfgGroup)
	} else {
		zap.GetLogger().Info("read nacos config success", "dataId", dataId, "group", cfgGroup, "content", content)
	}

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  cfgGroup,
		OnChange: func(namespace, group, dataId, data string) {
			zap.GetLogger().Info("nacos config changed",
				"namespace", namespace, "group", group, "dataId", dataId, "content", data)
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterService() {
	//读取配置
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "config.json",
		Group:  cfgGroup,
	})
	if err != nil {
		log.Fatal(err)
	}
	s := struct {
		AppName    string `json:"app_name"`
		ServerAddr string `json:"server_addr"`
		Port       int    `json:"port"`
	}{}
	if err = json.Unmarshal([]byte(content), &s); err != nil {
		log.Fatal(err)
	}

	_, err = namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          s.ServerAddr,
		Port:        uint64(s.Port),
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Metadata:    nil,
		ServiceName: s.AppName,
		GroupName:   cfgGroup,
	})
	if err != nil {
		log.Fatal(err)
	}
	zap.GetLogger().Info("register service success", "service name", s.AppName)
}
