package main

import (
	"flag"
	"iflytek.com/weipan4/learn-go/config/cmd/config"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var configPath = flag.String("config", "config/cmd/config", "Path to the config file")

func main() {
	flag.Parse()

	zap.InitLogger()

	cfgPath := filepath.Join(*configPath, "config.toml")
	if err := config.InitConfig(cfgPath); err != nil {
		zap.GetLogger().Error("init config failed", "file path", cfgPath, "msg", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	zap.GetLogger().Info("shutdown gracefully......")
}
