// zoom/distributed_encoding_platform/main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/api"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/consumer"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/dlq"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/producer"
	"iflytek.com/weipan4/learn-go/zoom/distributed_encoding_platform/scheduler"
)

func main() {
	// 根 context，收到系统信号后取消
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 初始化各组件
	producer.Init()
	defer producer.Close()

	dlq.Init()
	defer dlq.Close()

	scheduler.Init()
	defer scheduler.Close()

	// 启动消费者组（后台 goroutine）
	go consumer.Start(ctx)

	// 注册 HTTP 路由
	mux := http.NewServeMux()
	api.Register(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 优雅关闭：等待 ctx 取消后关闭 HTTP server
	go func() {
		<-ctx.Done()
		log.Println("[Main] shutting down HTTP server...")
		server.Shutdown(context.Background())
	}()

	log.Println("[Main] transcoding platform started on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("[Main] server error: %v", err)
	}
	log.Println("[Main] shutdown complete")
}
