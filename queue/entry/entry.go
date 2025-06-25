package main

import (
	"fmt"
	"iflytek.com/weipan4/learn-go/queue"
)

func main() {
	q := queue.Queue{1}
	q.Push(1)
	q.Push(2)
	fmt.Println(q.Pop())
	fmt.Println(q.Pop())
	fmt.Println(q.IsEmpty())
	fmt.Println(q.Pop())
	fmt.Println(q.IsEmpty())
}
