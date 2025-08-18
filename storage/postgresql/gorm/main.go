package main

import (
	"flag"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"iflytek.com/weipan4/learn-go/storage/postgresql/gorm/config"
	"iflytek.com/weipan4/learn-go/storage/postgresql/gorm/storage"
	"log"
)

var cfgPath = flag.String("cfgPath", "config.json", "path of postgresql config file")
var logPath = flag.String("logPath", "logs", "path of postgresql log file")

func main() {
	flag.Parse()
	zap.InitLogger(*logPath)
	if err := config.InitConfig(*cfgPath); err != nil {
		log.Fatalf("初始化postgresql配置失败, err: %v", err)
	}
	if err := storage.InitStorage(); err != nil {
		log.Fatalf("创建postgresql连接失败，err:%v", err)
	}
}
