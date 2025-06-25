package main

import (
	"fmt"
	"time"
)

// 示例：生产者-消费者模型
func producer(items chan<- int) {
	for i := 0; i < 5; i++ {
		items <- i
		time.Sleep(1 * time.Second)
	}
	close(items)
}

func consumer(id int, items <-chan int) {
	for item := range items {
		fmt.Printf("consumer %d got %d\n", id, item)
	}
}

func main() {
	items := make(chan int)
	go producer(items)

	for i := 1; i <= 3; i++ {
		go consumer(i, items)
	}

	time.Sleep(10 * time.Second)
}
