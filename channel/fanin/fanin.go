package main

import (
	"fmt"
	"math/rand"
	"time"
)

const producerNum = 3

// 示例：多个生产者，一个消费者
func produce(id int, nums chan<- int) {
	rand.NewSource(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		nums <- id*10 + i
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func consume(nums <-chan int, done chan<- struct{}) {
	for n := range nums {
		fmt.Printf("consumer received %d\n", n)
	}
	done <- struct{}{}
}

func main() {
	nums := make(chan int)
	done := make(chan struct{})

	for i := 0; i < producerNum; i++ {
		go produce(i, nums)
	}

	go func() {
		time.Sleep(2 * time.Second)
		close(nums)
	}()

	go consume(nums, done)

	<-done
}
