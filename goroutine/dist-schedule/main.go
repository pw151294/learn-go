package main

import (
	"context"
	"fmt"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	zap.InitLogger(zap.LogPath)

	s := NewScheduler(3, 100)

	s.Start()

	go func() {
		for res := range s.resultCh {
			zap.GetLogger().Info("receive result", "val", res)
		}
	}()

	go func() {
		for e := range s.errCh {
			zap.GetLogger().Error("task failed", "err", e)
		}
	}()

	for i := 0; i < 100; i++ {
		s.Submit(&AdTask{
			id: strconv.Itoa(i),
			adFunc: func(ctx context.Context) (*AdResult, error) {
				time.Sleep(1 * time.Second)
				return &AdResult{
					Result: fmt.Sprintf("auto discovery task result: %d", i),
				}, nil
			},
		})
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.GetLogger().Info("shutting down gracefully")
	s.Stop()
}
