package main

import (
	"fmt"
	"time"
)

func timeoutTask(result chan<- string) {
	time.Sleep(3 * time.Second)
	result <- "done"
}

func main() {
	result := make(chan string)
	go timeoutTask(result)
	t := time.NewTimer(5 * time.Second)
	select {
	case res := <-result:
		fmt.Printf("taks finished, result: %s\n", res)
	case <-t.C:
		fmt.Println("timeout")
	}
}
