package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"iflytek.com/weipan4/learn-go/config/viper/config"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	cfgPath = flag.String("config", "config/viper/config/st-agent.json", "config file path")
	logPath = flag.String("log", "config/viper/logs/viper.log", "log path")
)

func main() {
	flag.Parse()

	zap.InitLogger(*logPath)
	if err := config.InitConfig(*cfgPath); err != nil {
		log.Fatal(err)
	}

	// 定时任务打印配置信息
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				bytes, _ := json.Marshal(config.ConfigManager.GetConfig())
				zap.GetLogger().Info(fmt.Sprintf("agent config detail: %s", string(bytes)))
			}
		}
	}()

	// 程序优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.GetLogger().Info("shutting down gracefully......")
}
