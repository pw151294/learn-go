package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/api"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/configs"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/auth"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/cdn"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/index"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/lifecycle"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/internal/storage"
)

func main() {
	configPath := flag.String("config", "configs/config.toml", "配置文件路径")
	flag.Parse()

	// 1. 加载配置
	cfg, err := configs.Load(*configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	ctx := context.Background()

	// 2. 初始化 MinIO 客户端
	minioClient, err := storage.NewMinIOClient(&cfg.MinIO)
	if err != nil {
		log.Fatalf("初始化 MinIO 失败: %v", err)
	}
	if err := minioClient.EnsureBuckets(ctx); err != nil {
		log.Fatalf("确保 bucket 存在失败: %v", err)
	}

	// 3. 初始化 ES 客户端
	esClient := index.NewESClient(&cfg.ES)
	if err := esClient.EnsureIndex(ctx); err != nil {
		log.Fatalf("确保 ES 索引存在失败: %v", err)
	}
	recIdx := index.NewRecordingIndex(esClient)

	// 4. 初始化冷热分层管理器
	tierMgr := storage.NewTierManager(&cfg.MinIO, &cfg.Lifecycle)

	// 5. 初始化生命周期管理器和调度器
	lcMgr := lifecycle.NewLifecycleManager(&cfg.Lifecycle, minioClient, recIdx, tierMgr)
	scheduler := lifecycle.NewScheduler(lcMgr, cfg.Lifecycle.ScanInterval)
	scheduler.Start()
	defer scheduler.Stop()

	// 6. 初始化鉴权和 CDN 模块
	signer := auth.NewSigner(&cfg.Auth)
	validator := auth.NewValidator(&cfg.Auth)
	rewriter := cdn.NewURLRewriter(&cfg.CDN)

	// 7. 初始化 HTTP 服务
	handler := api.NewHandler(cfg, minioClient, recIdx, signer, validator, rewriter)
	router := api.SetupRouter(handler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 8. 启动 HTTP 服务（后台）
	go func() {
		log.Printf("服务启动，监听 %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("启动 HTTP 服务失败: %v", err)
		}
	}()

	// 9. 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("关闭服务时出错: %v", err)
	}
	log.Println("服务已关闭")
}
