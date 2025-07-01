package main

import (
	"flag"
	"iflytek.com/weipan4/learn-go/logger/zap"
	nacos1 "iflytek.com/weipan4/learn-go/middleware/nacos"
	cluster1 "iflytek.com/weipan4/learn-go/middleware/nacos/cluster"
	"iflytek.com/weipan4/learn-go/net/host"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfgPath = flag.String("cfgPath", "middleware/nacos/nacos.json", "nacos config file path")
	logPath = flag.String("logPath", "middleware/nacos/logs/subscriber.log", "nacos log path")
)

func main() {
	flag.Parse()

	zap.InitLogger(*logPath)
	host.InitHostInfo(false)
	nacos1.InitNacos(*cfgPath)

	err := nacos1.SubscribeService([]nacos1.SubScribeParamOptions{
		nacos1.WithSubscribeService(cluster1.Node),
		nacos1.WithSubscribeGroupName(cluster1.SvcGroup),
		nacos1.WithSubscribeCallback(nacos1.SubScribeCallback),
	})
	if err != nil {
		log.Fatalf("subscribe nacos service failed: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.GetLogger().Info("shutdown gracefully......")
}
