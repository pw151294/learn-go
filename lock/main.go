package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func add(p *int32) {
	//*p += 1
	atomic.AddInt32(p, 1)
}

func main() {
	c := int32(1)
	for i := 0; i < 1000; i++ {
		go add(&c)
	}
	time.Sleep(time.Second)
	fmt.Println(c)
}
