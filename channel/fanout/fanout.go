package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 示例：一个生产者，多个消费者
func produce(nums chan<- int) {
	for i := 0; i < 10; i++ {
		nums <- i
	}
	close(nums)
}

func consumer(id int, nums <-chan int, done chan<- struct{}) {
	rand.NewSource(time.Now().UnixNano())
	for n := range nums {
		fmt.Printf("consumer %d received %d\n", id, n)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
	done <- struct{}{}
}

const numConsumers = 3

func main() {
	nums := make(chan int)
	done := make(chan struct{}, numConsumers)

	go produce(nums)

	for i := 0; i < numConsumers; i++ {
		go consumer(i, nums, done)
	}

	// 等待所有消费者完成
	for i := 0; i < numConsumers; i++ {
		<-done
	}
}
