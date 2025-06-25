package main

import (
	"flag"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"log"
)

var (
	cfgPath = flag.String("nacosCfg", "middleware/nacos/nacos.json", "nacos config file path")
)

func main() {
	flag.Parse()

	// 初始化nacos配置
	zap.InitLogger()
	if err := InitNacosConfig(*cfgPath); err != nil {
		log.Fatal(err)
	}

	// 完成服务注册和动态配置
	InitNacos()

	forever := make(chan struct{})
	<-forever
}
