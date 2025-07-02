package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris/v12"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"iflytek.com/weipan4/learn-go/net/iris/config"
	go_redis "iflytek.com/weipan4/learn-go/storage/redis/go-redis"
	redis_config "iflytek.com/weipan4/learn-go/storage/redis/go-redis/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	cfgPath = flag.String("config", "net/iris/config/iris.json", "path to config file")
	logPath = flag.String("log", "net/iris/logs/iris.log", "path to log file")
)

func main() {
	flag.Parse()

	// 初始化项目配置和日志
	if err := config.InitConfig(*cfgPath); err != nil {
		log.Fatal(err)
	}
	zap.InitLogger(*logPath)

	// 初始化Redis模块配置
	if err := redis_config.InitConfig("storage/redis/go-redis/config/redis.json"); err != nil {
		log.Fatal(err)
	}
	if err := go_redis.InitRedis(); err != nil {
		log.Fatal(err)
	}

	// 启动服务器
	r1 := InitRouter()
	address := fmt.Sprintf("%s:%d", config.Cfg.Server.Host, 8090)
	go func() {
		if err := r1.Run(iris.Addr(address)); err != nil {
			log.Fatal(err)
		}
	}()

	r2 := InitRouter()
	address = fmt.Sprintf("%s:%d", config.Cfg.Server.Host, 8091)
	go func() {
		if err := r2.Run(iris.Addr(address)); err != nil {
			log.Fatal(err)
		}
	}()

	r3 := InitRouter()
	address = fmt.Sprintf("%s:%d", config.Cfg.Server.Host, 8092)
	go func() {
		if err := r3.Run(iris.Addr(address)); err != nil {
			log.Fatal(err)
		}
	}()

	// 程序优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.GetLogger().Info("shutting down gracefully")
}
