package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
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
	createClient()
	//InitConfig(dataId, cfgGroup)
	//RegisterService()
}

func createClient() {
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

func InitConfig(dataId, cfgGroup string) (string, error) {
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  cfgGroup,
	})
	if err != nil {
		zap.GetLogger().Error("read nacos config failed", "dataId", dataId, "group", cfgGroup, "err", err)
		return "", err
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
		zap.GetLogger().Error("listen nacos config failed", "dataId", dataId, "group", cfgGroup, "err", err)
	}
	return content, nil
}

func RegisterService(opts1 []MetadataOptions, opts2 []RegisterInstanceParamOptions) error {
	instMeta := NewMetadata(opts1...)
	opts2 = append(opts2, WithMetadata(instMeta))
	regInstParam := newRegisterInstanceParam(opts2...)

	_, err := namingClient.RegisterInstance(regInstParam)
	if err != nil {
		zap.GetLogger().Error("register service failed",
			"service name", regInstParam.ServiceName, "group name", regInstParam.GroupName)
		return err
	}

	zap.GetLogger().Info("register service success!",
		"service name", regInstParam.ServiceName, "group name", regInstParam.GroupName)
	return nil
}

func SubscribeService(opts []SubScribeParamOptions) error {
	param := newSubscribeParam(opts...)
	if err := namingClient.Subscribe(param); err != nil {
		zap.GetLogger().Error("subscribe service failed",
			"service name", param.ServiceName, "group name", param.GroupName, "clusters", param.Clusters, "err", err)
		return err
	}
	zap.GetLogger().Info("subscribe service success",
		"service name", param.ServiceName, "group name", param.GroupName, "clusters", param.Clusters)
	return nil
}
