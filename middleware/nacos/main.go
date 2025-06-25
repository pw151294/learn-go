package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"iflytek.com/weipan4/learn-go/http/host"
	"iflytek.com/weipan4/learn-go/logger/zap"
)

var (
	cfgPath = flag.String("nacosCfg", "middleware/nacos/nacos.json", "nacos config file path")
	logPath = flag.String("logPath", "middleware/nacos/logs/nacos.log", "application log file path")
)

func main() {
	flag.Parse()

	// 初始化基础配置
	zap.InitLogger(*logPath)
	host.InitHostInfo(false)

	// 完成服务注册和动态配置
	InitNacos(*cfgPath)

	r := gin.Default()
	r.Run(fmt.Sprintf(":%d", Cfg.Application.Port))
}
