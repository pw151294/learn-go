package main

import (
	"fmt"
	"sync"
)

type Person struct {
	rwMutex sync.RWMutex
	salary  int
	level   int
}

func (p *Person) printPerson(wg *sync.WaitGroup) {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	fmt.Println(p.salary)
	fmt.Println(p.level)
	wg.Done()
}

func main() {
	p := Person{salary: 10000, level: 1}

	wg := &sync.WaitGroup{}
	wg.Add(3) //WaitGroup增加3个等待任务

	go p.printPerson(wg)
	go p.printPerson(wg)
	go p.printPerson(wg)
	wg.Wait() // 相当于sleep 等待之前的所有协程都执行完了wg.Done()再结束等待

	m1 := sync.Mutex{}
	m1.Lock()
	m2 := m1
	m1.Unlock()
	m2.Lock()
}
