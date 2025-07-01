package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"
)

// EventStreamHandler 定义处理Event-Stream响应的函数类型
type EventStreamHandler func(http.ResponseWriter, *http.Request) <-chan string

// StartEventStream 启动处理Event-Stream响应的 goroutine
func StartEventStream(w http.ResponseWriter, r *http.Request, handler EventStreamHandler) {
	// 设置Event-Stream响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "Keep-Alive")

	// 防止浏览器缓存
	w.Header().Set("X-Accel-Buffering", "no")

	// 创建一个 goroutine 处理事件流
	go func() {
		// 设置状态码为200 OK
		w.WriteHeader(http.StatusOK)

		// 创建bufio.Writer来高效处理输出
		writer := bufio.NewWriter(w)
		defer writer.Flush()

		// 获取事件通道
		eventChan := handler(w, r)

		// 循环读取事件并写入响应
		for {
			select {
			case event, ok := <-eventChan:
				if !ok {
					// 事件通道关闭，退出循环
					return
				}
				// 写入事件数据
				if _, err := fmt.Fprintf(writer, "data: %s\n\n", event); err != nil {
					// 处理写入错误
					if err == io.EOF {
						// 客户端断开连接
						return
					}
					fmt.Printf("写入错误: %v\n", err)
					return
				}
				// 刷新输出以确保数据立即发送
				writer.Flush()

			case <-time.After(time.Second):
				// 每秒发送一个心跳以保持连接 alive
				if _, err := fmt.Fprintf(writer, "data: \n\n"); err != nil {
					if err == io.EOF {
						return
					}
					fmt.Printf("写入错误: %v\n", err)
					return
				}
				writer.Flush()
			}
		}
	}()
}

// ExampleHandler 是一个示例处理函数，生成事件数据
func ExampleHandler(w http.ResponseWriter, r *http.Request) <-chan string {
	// 创建一个 buffered 通道，用于发送事件
	eventChan := make(chan string, 5)
	go func() {
		// 发送一些示例事件
		eventChan <- "event1: data1"
		time.Sleep(2 * time.Second)
		eventChan <- "event2: data2"
		time.Sleep(2 * time.Second)
		eventChan <- "event3: data3"
		close(eventChan)
	}()
	return eventChan
}

func main() {
	// 定义一个简单的 HTTP 服务
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		// 启动处理Event-Stream的 goroutine
		StartEventStream(w, r, ExampleHandler)
	})

	// 启动 HTTP 服务器
	fmt.Println("Server 开始运行在 :8080")
	http.ListenAndServe(":8080", nil)
}
