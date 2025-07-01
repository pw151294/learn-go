package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"iflytek.com/weipan4/learn-go/net/host"
	"iflytek.com/weipan4/learn-go/net/websocket/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cfgPath = flag.String("cfgPath", "net/websocket/config/config.toml", "websocket config file path")

func main() {
	flag.Parse()

	// 初始化基本信息
	host.InitHostInfo(false)
	zap.InitLogger("net/websocket/logs/websocket.log")
	if err := config.InitConfig(*cfgPath); err != nil {
		log.Fatalf("read websocket config failed: %v", err)
	}

	// 初始化客户端
	wsc := NewWebSocketClient(
		WithMessageHandler(func(message []byte) error {
			zap.GetLogger().Info("begin handling messages......")
			time.Sleep(1 * time.Second)
			zap.GetLogger().Info("handling message completed")
			return nil
		}),
		WithDataHandler(func(data []byte) error {
			zap.GetLogger().Info("begin handling binary data......")
			time.Sleep(1 * time.Second)
			zap.GetLogger().Info("handling binary data completed")
			return nil
		}),
		WithErrorHandler(func(err error) {
			zap.GetLogger().Error("receive exception, begin handling......", "error", err)
			time.Sleep(1 * time.Second)
			zap.GetLogger().Info("handling error completed")
		}))

	// 初始化websocket监听端口
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "upgrade http request to websocket protocol failed", http.StatusBadRequest)
		}
		wsc.OpenWskConn(conn)
		go wsc.ReadLoop()
	})

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", config.WskCfg.ServerPort), nil); err != nil {
			log.Fatal("start http server failed", err)
		}
	}()
	zap.GetLogger().Info("start websocket server successfully!")

	// 程序优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.GetLogger().Info("shutting down gracefully......")
}
