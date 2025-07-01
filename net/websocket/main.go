package main

import (
	"iflytek.com/weipan4/learn-go/logger/zap"
	"iflytek.com/weipan4/learn-go/net/host"
	"iflytek.com/weipan4/learn-go/net/websocket/ws_handler"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	wsUrlPrefix = "ws://"
	wsUrlHost   = "127.0.0.1:8081"
)

func main() {
	host.InitHostInfo(false)
	zap.InitLogger("net/websocket/logs/websocket.log")
	wsh := ws_handler.NewWebSocketHandler(
		ws_handler.WithWsUrl(wsUrlPrefix+wsUrlHost),
		ws_handler.WithOnMessage(func(bytes []byte) {
			zap.GetLogger().Info("receive message. begin handling messages......")
			time.Sleep(1 * time.Second)
			zap.GetLogger().Info("handling message completed")
		}),
		ws_handler.WithOnError(func(err error) {
			zap.GetLogger().Info("receive error. begin handling error......")
			time.Sleep(1 * time.Second)
			zap.GetLogger().Info("handling error completed")
		}))
	if err := wsh.Initialize(); err != nil {
		log.Fatal(err)
	}

	go func() {
		zap.GetLogger().Info("begin sending heartbeat......")
		wsh.StartPing()
	}()

	go func() {
		zap.GetLogger().Info("begin receiving message loop.....")
		wsh.ReceiveMessageLoop()
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.GetLogger().Info("shutting down gracefully......")
}
