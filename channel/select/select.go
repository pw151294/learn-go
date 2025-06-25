package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator() chan int {
	out := make(chan int)
	go func() {
		i := 0
		for {
			time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
			out <- i
			i++
		}
	}()
	return out
}

func createWorker(id int) chan<- int {
	out := make(chan int)
	go func() {
		for val := range out {
			fmt.Printf("Worker %d received %d\n", id, val)
		}
	}()
	return out
}

func main() {
	var c1, c2 = generator(), generator()
	//select {
	//case val := <-c1:
	//	fmt.Printf("received from c1: %v\n", val)
	//case val := <-c2:
	//	fmt.Printf("received from c2: %v\n", val)
	//default:
	//	fmt.Println("No value received")
	//}

	var worker = createWorker(0)
	var values []int
	tm := time.After(10 * time.Second)
	tick := time.Tick(1 * time.Second)
	for {
		var activeWorker chan<- int
		var activeVal int
		if len(values) > 0 {
			activeWorker = worker
			activeVal = values[0]
		}

		select {
		case val := <-c1:
			values = append(values, val)
			fmt.Printf("received from c1: %v\n", val)
		case val := <-c2:
			values = append(values, val)
			fmt.Printf("received from c2: %v\n", val)
		case activeWorker <- activeVal:
			values = values[1:]
		case <-time.After(800 * time.Millisecond):
			fmt.Println("timed out")
		case <-tick:
			fmt.Printf("values len: %v\n", len(values))
		case <-tm:
			fmt.Println("bye bye ......")
			return
		}
	}
}
