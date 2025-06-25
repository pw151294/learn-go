package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pelletier/go-toml/v2"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"log"
)

const (
	dataId   = "config.toml"
	cfgGroup = "cmdb"
	svcGroup
)

var Cfg *Config

type Logging struct {
	Log          string `toml:"log"`
	Level        string `toml:"level"`
	DisableColor bool   `toml:"disable_color"`
	Rotate       int    `toml:"rotate"`
}

type Application struct {
	ServerAddr string `toml:"server_addr"`
	Port       int    `toml:"port"`
	AppName    string `toml:"app_name"`
}

type Config struct {
	Logging     Logging     `toml:"logging"`
	Application Application `toml:"application"`
}

var (
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
)

func InitNacos(cfgPath string) {
	if err := InitNacosConfig(cfgPath); err != nil {
		log.Fatalf("read nacos config file failed: %v", err)
	}
	CreateClient()
	InitConfig()
	RegisterService()
}

func CreateClient() {
	// 创建client config
	clientCfg := *constant.NewClientConfig(
		constant.WithNamespaceId(NacosCfg.NamespaceId),
		constant.WithUsername(NacosCfg.Username),
		constant.WithPassword(NacosCfg.Password),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(NacosCfg.LogPath))
	// 创建server configs
	serverCfgs := []constant.ServerConfig{
		*constant.NewServerConfig(
			NacosCfg.IpAddr,
			NacosCfg.Port,
			constant.WithScheme(NacosCfg.Scheme),
			constant.WithContextPath(NacosCfg.ContextPath)),
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
		if err = toml.Unmarshal([]byte(content), &Cfg); err != nil {
			log.Fatalf("load nacos config failed: %v", err)
		}
		zap.GetLogger().Info("load nacos config success!")
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
	instMeta := NewMetadata(
		WithInstanceId(fmt.Sprintf("%s-id", Cfg.Application.AppName)),
		WithPort(Cfg.Application.Port))

	regInstParam := newRegisterInstanceParam(
		WithInstanceIp(Cfg.Application.ServerAddr),
		WithInstancePort(uint64(Cfg.Application.Port)),
		WithServiceName(Cfg.Application.AppName),
		WithMetadata(instMeta),
	)

	_, err := namingClient.RegisterInstance(regInstParam)
	if err != nil {
		log.Fatalf("register service failed, service name:%s, message:%v", Cfg.Application.AppName, err)
	}

	zap.GetLogger().Info("register service success", "service name", Cfg.Application.AppName)
}
