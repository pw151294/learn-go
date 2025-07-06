package main

import (
	"flag"
	"iflytek.com/weipan4/learn-go/logger/zap"
	nacos1 "iflytek.com/weipan4/learn-go/middleware/nacos"
	cluster1 "iflytek.com/weipan4/learn-go/middleware/nacos/cluster"
	"iflytek.com/weipan4/learn-go/net/host"
	go_redis "iflytek.com/weipan4/learn-go/storage/redis/go-redis"
	redis_config "iflytek.com/weipan4/learn-go/storage/redis/go-redis/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfgPath     = flag.String("cfgPath", "middleware/nacos/nacos.json", "nacos config file path")
	logPath     = flag.String("logPath", "middleware/nacos/logs/nacos.log", "nacos log path")
	serviceName = flag.String("serviceName", "", "service name")
	groupName   = flag.String("groupName", "", "service group name")
	port        = flag.Int("port", 0, "service port")
)

const redisConfigFile = "/Users/a123/Downloads/go/learn-go/learn-go/storage/redis/go-redis/config/redis.json"

func main() {
	flag.Parse()
	if *serviceName == "" || *groupName == "" || *port == 0 {
		log.Fatalf("invalid args of service")
	}

	// 完成一些初始化的配置
	zap.InitLogger(*logPath)
	host.InitHostInfo(false)
	if err := redis_config.InitConfig(redisConfigFile); err != nil {
		log.Fatalf("init redis config failed: %v", err)
	}
	if err := go_redis.InitRedis(); err != nil {
		log.Fatalf("init redis client failed: %v", err)
	}

	// 初始化nacos配置
	nacos1.InitNacos(*cfgPath)
	nacos1.Port = *port

	// 注册服务
	go cluster1.StartNode(*serviceName, *groupName, *port)

	// 实现优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.GetLogger().Info("shutting down gracefully")
}
