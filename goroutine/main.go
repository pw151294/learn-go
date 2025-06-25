package main

import (
	"fmt"
	"math"
	"time"
)

func do(i int, ch chan struct{}) {
	fmt.Println(i)
	time.Sleep(1 * time.Second)
	<-ch // 执行完业务方法之后将信号从channel里取出
}

func main() {
	//cnt := 0
	//for i := 0; i < 100000; i++ {
	//	go func(int) {
	//		fmt.Println(i)
	//		cnt++
	//	}(i)
	//}
	//time.Sleep(time.Second)
	//fmt.Println("cnt:", cnt)
	//var a [10]int
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		for {
	//			a[i]++
	//		}
	//	}()
	//}
	//time.Sleep(time.Millisecond)
	//fmt.Print(a)

	c := make(chan struct{}, 3000)
	for i := 0; i < math.MaxInt32; i++ { // 使用channel缓存区来解决协程数过多的问题 这样可以保证同时运行的协程数不超过3000
		c <- struct{}{}
		go do(i, c)
	}
}
