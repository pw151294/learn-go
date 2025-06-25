package main

import (
	"fmt"
	"sync"
	"time"
)

type Person struct {
	mu      sync.Mutex
	rwMutex sync.RWMutex
	salary  int
	level   int
}

func (p *Person) promote() {
	p.mu.Lock()         // 操作之前加锁
	defer p.mu.Unlock() // 操作之后或者离开函数之前解锁

	p.salary += 1000
	fmt.Println(p.salary)
	p.level++
	fmt.Println(p.level)
}

func (p *Person) printPerson() {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	fmt.Println(p.salary)
	fmt.Println(p.level)
}

func main() {
	p := Person{salary: 10000, level: 1}

	go p.promote()
	go p.promote()
	go p.promote()

	time.Sleep(1 * time.Second)
}
