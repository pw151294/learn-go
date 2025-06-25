package main

import (
	"fmt"
	"time"
)

type MyMutex chan struct{}

// NewMyMutex 使用一个缓冲大小为1的channel
func NewMyMutex() MyMutex {
	ch := make(chan struct{}, 1)
	return ch
}

// Lock 加锁时 向channel塞入一个数据
func (m *MyMutex) Lock() {
	*m <- struct{}{}
}

func (m *MyMutex) Unlock() {
	<-(*m)
}

func add(p *int32) {
	*p++
}

func addWithLock(p *int32, m MyMutex) {
	m.Lock()
	defer m.Unlock()
	*p++
}

func main() {
	c := int32(0)
	m := NewMyMutex()
	for i := 0; i < 1000; i++ {
		//go add(&c)
		go addWithLock(&c, m)
	}
	time.Sleep(time.Second)
	fmt.Println(c)
}
