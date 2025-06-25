package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 示例：后台日志处理器
func logProcessor(logs <-chan string) {
	rand.NewSource(time.Now().UnixNano())
	for le := range logs {
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		fmt.Printf("processes log: %s\n", le)
	}
}

func main() {
	logCh := make(chan string, 100)
	go logProcessor(logCh)

	for i := 0; i < 10; i++ {
		logCh <- fmt.Sprintf("log entry %d", i)
	}
	close(logCh)

	time.Sleep(3 * time.Second)
}
