package main

import (
	"flag"
	"iflytek.com/weipan4/learn-go/http/host"
	"iflytek.com/weipan4/learn-go/logger/zap"
	nacos1 "iflytek.com/weipan4/learn-go/middleware/nacos"
	cluster1 "iflytek.com/weipan4/learn-go/middleware/nacos/cluster"
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

func main() {
	flag.Parse()
	if *serviceName == "" || *groupName == "" || *port == 0 {
		log.Fatalf("invalid args of service")
	}

	// 完成一些初始化的配置
	zap.InitLogger(*logPath)
	host.InitHostInfo(false)

	// 初始化nacos配置
	nacos1.InitNacos(*cfgPath)

	// 注册服务
	go cluster1.StartNode(*serviceName, *groupName, *port)

	// 实现优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.GetLogger().Info("shutting down gracefully")
}
