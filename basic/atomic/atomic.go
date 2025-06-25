package main

import (
	"fmt"
	"sync"
	"time"
)

type atomicInt struct {
	val  int
	lock sync.Mutex
}

func (a *atomicInt) inc() {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.val++
}

func (a *atomicInt) get() int {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.val
}

func main() {
	var a atomicInt
	a.inc()
	go func() {
		a.inc()
	}()
	time.Sleep(time.Second)
	fmt.Println(a.get())
}
