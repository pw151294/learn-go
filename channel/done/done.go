package main

import (
	"fmt"
	"sync"
)

func doWork(id int, job <-chan int, res chan<- bool) {
	for val := range job {
		fmt.Printf("Worker %d received job %d\n", id, val)
		//res <- true
		go func() { res <- true }()
	}
}

func doWgWork(id int, job <-chan int, wg *sync.WaitGroup) {
	for val := range job {
		fmt.Printf("Worker %d received job %d\n", id, val)
		wg.Done()
	}
}

type Worker struct {
	job chan int
	res chan bool
}

type WgWorker struct {
	job chan int
	wg  *sync.WaitGroup
}

func createWorker(id int) Worker {
	worker := Worker{
		job: make(chan int),
		res: make(chan bool),
	}
	go doWork(id, worker.job, worker.res)
	return worker
}

func createWgWorker(id int, wg *sync.WaitGroup) WgWorker {
	worker := WgWorker{
		job: make(chan int),
		wg:  wg,
	}
	go doWgWork(id, worker.job, worker.wg)
	return worker
}

func channelDemo() {
	var workers [10]Worker
	for i := 0; i < 10; i++ {
		workers[i] = createWorker(i)
	}
	for i, worker := range workers {
		worker.job <- 'a' + i //channel的send操作是阻塞式的 发送的消息必须要有channel来接收
	}
	for i, worker := range workers {
		worker.job <- 'A' + i
	}

	for _, worker := range workers {
		<-worker.res
		<-worker.res
	}
}

func waitGroup() {
	var wg sync.WaitGroup
	var wgWorkers [10]WgWorker
	for i := 0; i < 10; i++ {
		wgWorkers[i] = createWgWorker(i, &wg)
	}

	wg.Add(20)
	for i, worker := range wgWorkers {
		worker.job <- 'a' + i
	}
	for i, wgWorker := range wgWorkers {
		wgWorker.job <- 'A' + i
	}
	wg.Wait()
}

func main() {
	//channelDemo()
	waitGroup()
}
