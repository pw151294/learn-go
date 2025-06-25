package main

import (
	"fmt"
	"time"
)

func RangeWorker(id int, c chan int) {
	for val := range c { //range来接收数据 直到channel被close掉之后就自动退出循环
		fmt.Printf("Worker %d received val:%d\n", id, val)
	}
}

func ifWorker(id int, c chan int) {
	for {
		if val, ok := <-c; !ok { //使用if来获取值 通过ok判断channel是否被close 发现Channel被关闭之后就自动退出循环
			fmt.Printf("Worker %d received val:%d\n", id, val)
		}
	}
}

func worker(id int, c chan int) {
	for {
		fmt.Printf("Worker %d received val:%d\n", id, <-c)
	}
}

func createWorker(id int) chan<- int {
	c := make(chan int)
	go func() {
		for {
			fmt.Printf("worker %d received %c\n", id, <-c)
		}
	}()
	return c
}

func chanDemo1() {
	var channels [10]chan int
	for i := 0; i < 10; i++ {
		channels[i] = make(chan int)
		go RangeWorker(i, channels[i])
	}
	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
	}
	for i := 0; i < 10; i++ {
		channels[i] <- 'A' + i
	}
	time.Sleep(time.Second)
}

func chanDemo2() {
	var channels [10]chan<- int
	for i := 0; i < 10; i++ {
		channels[i] = createWorker(i)
	}
	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
		channels[i] <- 'A' + i
	}
}

func buffChannelDemo() {
	c := make(chan int, 3)
	go worker(0, c)
	c <- 'a'
	c <- 'B'
	c <- 'c'
	c <- 'd'
	time.Sleep(time.Second)
}

func channelClose() {
	c := make(chan int)
	go worker(0, c)
	c <- 'a'
	c <- 'B'
	c <- 'c'
	c <- 'd'
	close(c)
	time.Sleep(time.Millisecond)
}

// 通过通信的方式来实现共享内存
func watch(c chan int) {
	if <-c == 1 {
		fmt.Println("channel closed")
	}
}

func main() {
	//fmt.Println("channel as first-class citizen")
	//chanDemo1()
	//chanDemo2()
	//fmt.Println("buffered channel ")
	//buffChannelDemo()
	//fmt.Println("channel close and range receive val")
	//channelClose()
	c := make(chan int)
	go watch(c)
	time.Sleep(time.Second)
	c <- 1
	time.Sleep(time.Second)

	// 通过select实现非阻塞的Channel
	c1 := make(chan int, 5)
	c2 := make(chan int, 5)
	select {
	case <-c1: //如果c1成功读取数据 则进行该case处理语句
		fmt.Println("c1")
	case c2 <- 1: //如果成功向c2写入数据 则进行该case的处理语句
		fmt.Println("c2")
	default: //如果上面都没有成功 则进入default处理流程
		fmt.Println("none")
	}
}
