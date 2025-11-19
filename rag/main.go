package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"iflytek.com/weipan4/learn-go/rag/configs"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/embedding"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/mysql"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/weaviate"
	"iflytek.com/weipan4/learn-go/rag/internal/schedulers"
	"iflytek.com/weipan4/learn-go/rag/pkg/loggers"
	"iflytek.com/weipan4/learn-go/rag/web/route"
)

var configPath = flag.String("cnf", filepath.Join("./rag/configs", "config.toml"), "config file path")
var workerNum = flag.Int("workerNum", 4, "worker num")
var capacity = flag.Int("capacity", 10, "buffer for chan")
var period = flag.Int("period", 3000, "period of splitting text task in unit of million seconds")

func main() {
	flag.Parse()
	if err := configs.InitConfig(*configPath); err != nil {
		log.Fatalf("init global config failed: %v", err)
	}

	// 1. 初始化日志系统
	loggerConfig := configs.GetConfig().Log
	if err := loggers.InitGlobalLogger(&loggerConfig); err != nil {
		log.Fatalf("init logger system failed: %v", err)
	}

	// 2. 初始化infra
	if err := mysql.InitDB(); err != nil {
		log.Fatalf("create mysql client failed: %v", err)
	}
	if err := weaviate.InitWeaviateClient(); err != nil {
		log.Fatalf("create weaviate client failed: %v", err)
	}
	embeddingClient, err := embedding.NewOpenAIEmbeddingClient()
	if err != nil {
		log.Fatalf("create embedding client failed: %v", err)
	}

	// 3. 初始化调度器
	scheduler := schedulers.NewScheduler(mysql.DB, embeddingClient, *workerNum, *capacity, *period)
	go scheduler.Start()

	// 4. 初始化web层
	r := route.InitRouter()
	go func() {
		if err := r.Run("localhost:8080"); err != nil {
			log.Fatalf("start http server failed: %v", err)
		}
	}()

	// 程序优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	scheduler.Close()
	loggers.Info("shutting down gracefully ......")
}
