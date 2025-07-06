package nacos

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/redis/go-redis/v9"
	"iflytek.com/weipan4/learn-go/encrypt"
	"iflytek.com/weipan4/learn-go/lock/redislock"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"iflytek.com/weipan4/learn-go/net/host"
	go_redis "iflytek.com/weipan4/learn-go/storage/redis/go-redis"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

var (
	NacosCfg *NacosConfig
	Port     int
)

type RegisterInstanceParamOptions func(param *vo.RegisterInstanceParam)

type NacosConfig struct {
	AppName     string `json:"app_name"`
	IpAddr      string `json:"ip_addr"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	ContextPath string `json:"context_path"`
	Port        uint64 `json:"port"`
	Scheme      string `json:"scheme"`
	NamespaceId string `json:"namespace_id"`
	LogPath     string `json:"log_path"`
	LogLevel    string `json:"log_level"`
}

func InitNacosConfig(cfgPath string) error {
	file, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	nacosCfg := &NacosConfig{}
	if err = json.NewDecoder(file).Decode(nacosCfg); err != nil {
		return err
	}
	NacosCfg = nacosCfg
	return nil
}

func newRegisterInstanceParam(opts ...RegisterInstanceParamOptions) vo.RegisterInstanceParam {
	param := vo.RegisterInstanceParam{
		Ip:          host.HostInfo.GetIP(),
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		ClusterName: "DEFAULT",
		GroupName:   svcGroup,
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(&param)
		}
	}

	return param
}

func WithInstanceIp(ip string) RegisterInstanceParamOptions {
	return func(param *vo.RegisterInstanceParam) {
		param.Ip = ip
	}
}

func WithInstancePort(port uint64) RegisterInstanceParamOptions {
	return func(param *vo.RegisterInstanceParam) {
		param.Port = port
	}
}

func WithServiceName(serviceName string) RegisterInstanceParamOptions {
	return func(param *vo.RegisterInstanceParam) {
		param.ServiceName = serviceName
	}
}

func WithClusterName(clusterName string) RegisterInstanceParamOptions {
	return func(param *vo.RegisterInstanceParam) {
		param.ClusterName = clusterName
	}
}

func WithGroupName(groupName string) RegisterInstanceParamOptions {
	return func(param *vo.RegisterInstanceParam) {
		param.GroupName = groupName
	}
}

func WithMetadata(metadata map[string]string) RegisterInstanceParamOptions {
	return func(param *vo.RegisterInstanceParam) {
		param.Metadata = metadata
	}
}

type SubScribeParamOptions func(param *vo.SubscribeParam)

func newSubscribeParam(opts ...SubScribeParamOptions) *vo.SubscribeParam {
	param := &vo.SubscribeParam{
		GroupName:         svcGroup,
		Clusters:          []string{"DEFAULT"},
		SubscribeCallback: SubScribeCallback,
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(param)
		}
	}

	return param
}

func SubScribeCallback(services []model.Instance, err error) {
	zap.GetLogger().Info("subscribed service list info changed", "port", Port)
	//获取instance实例id
	insIds := make([]string, 0)
	for _, svc := range services {
		if svc.Healthy {
			insIds = append(insIds, svc.InstanceId)
		}
	}
	sort.Strings(insIds) //排序 为了保证updateId的唯一性
	uid := encrypt.MD5Encode(strings.Join(insIds, ":"))
	zap.GetLogger().Info(fmt.Sprintf("update id: %s", uid))

	// 初始化分布式锁
	rand.NewSource(time.Now().UnixNano())
	redisCli := go_redis.GetClient().(*redis.Client)
	builder := &redislock.Builder{}
	rl, err := builder.WithRedisClient(redisCli).WithLockName(uid).WithExpire(3 * time.Second).Build()
	if err != nil {
		zap.GetLogger().Error("create redis lock failed", "message", err)
		return
	}
	if rl.TryLock() {
		defer rl.Unlock()
		zap.GetLogger().Info("lock success! begin handle subscribe callback")
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func WithSubscribeService(serviceName string) SubScribeParamOptions {
	return func(param *vo.SubscribeParam) {
		param.ServiceName = serviceName
	}
}

func WithSubscribeGroupName(groupName string) SubScribeParamOptions {
	return func(param *vo.SubscribeParam) {
		param.GroupName = groupName
	}
}

func WithSubscribeClusters(clusters []string) SubScribeParamOptions {
	return func(param *vo.SubscribeParam) {
		param.Clusters = clusters
	}
}

func WithSubscribeCallback(callback func(services []model.Instance, err error)) SubScribeParamOptions {
	return func(param *vo.SubscribeParam) {
		param.SubscribeCallback = callback
	}
}
