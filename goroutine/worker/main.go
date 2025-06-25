package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 示例：使用Channel同步多个Goroutine
func worker(id int, ready <-chan struct{}, done chan<- int) {
	<-ready
	fmt.Printf("worker %d starting\n", id)
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	fmt.Printf("worker %d done\n", id)
	done <- id
}

const workerCount = 5

func main() {
	ready := make(chan struct{})
	done := make(chan int, workerCount)

	for i := 0; i < workerCount; i++ {
		go worker(i, ready, done)
	}
	close(ready)

	for i := 0; i < workerCount; i++ {
		<-done
	}

	fmt.Println("all workers done")
}
