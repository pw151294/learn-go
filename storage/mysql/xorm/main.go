package main

import (
	"flag"
	"iflytek.com/weipan4/learn-go/storage/mysql/xorm/config"
	"iflytek.com/weipan4/learn-go/storage/mysql/xorm/datasource"
	"iflytek.com/weipan4/learn-go/storage/mysql/xorm/model"
	"log"
)

var cfgPath = flag.String("config", "storage/mysql/xorm/config/mysq.json", "config file for mysql")

func main() {
	flag.Parse()

	// 初始化
	err := config.InitConfig(*cfgPath)
	if err != nil {
		log.Fatalf("读取配置文件失败：%v", err)
	}
	err = datasource.InitEngine()
	if err != nil {
		log.Fatalf("初始化MySQL数据源失败：%v", err)
	}
	// 初始化表
	err = datasource.GetEngine().Sync2(new(model.Node))
	if err != nil {
		log.Fatalf("初始化ORM映射失败：%v", err)
	}
}
