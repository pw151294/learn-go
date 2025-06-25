package main

import (
	"fmt"
	"math/rand"
	"time"
)

// <-chan 表示channel只能收数据
// chan<- 表示channel只能发数据
func msgGen1() chan string {
	c := make(chan string)
	go func() {
		i := 0
		for {
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			c <- fmt.Sprintf("message %d", i)
			i++
		}
	}()
	return c
}

func msgGen2(name string, done chan struct{}) chan string {
	c := make(chan string)
	go func() {
		i := 0
		for {
			select {
			case <-time.After(time.Duration(rand.Intn(5000)) * time.Millisecond):
				c <- fmt.Sprintf("message %s", name)
			case <-done:
				fmt.Println("cleaning up")
				time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
				fmt.Println("cleaning up")
				done <- struct{}{}
				return
			}
			i++
		}
	}()
	return c
}

func fanIn(c1 <-chan string, c2 <-chan string) chan string {
	c := make(chan string)
	go func() {
		for {
			c <- <-c1
		}
	}()
	go func() {
		for {
			c <- <-c2
		}
	}()
	return c
}

func fanInBySelect(c1, c2 chan string) chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case val := <-c1:
				c <- val
			case val := <-c2:
				c <- val

			}
		}
	}()
	return c
}

func fanInBatch(chs ...chan string) chan string {
	c := make(chan string)
	for _, ch := range chs {
		go func(in chan string) {
			for {
				c <- <-ch
			}
		}(ch)
	}
	return c
}

func nonBlockingWait(c chan string) (string, bool) {
	select {
	case val := <-c:
		return val, true
	default:
		return "", false
	}
}

func timeoutWait(c chan string, timeout time.Duration) (string, bool) {
	select {
	case val := <-c:
		return val, true
	case <-time.After(timeout):
		return "", false
	}
}

func main() {
	//m := msgGen1()
	//for {
	//	fmt.Println(<-m)
	//}

	//m1 := msgGen1()
	//m2 := msgGen1()
	//m := fanIn(m1, m2)
	//for {
	//	fmt.Println(<-m)
	//}

	//m1 := msgGen1()
	//m2 := msgGen1()
	//for {
	//	fmt.Println(<-m1)
	//	if m, ok := nonBlockingWait(m2); ok {
	//		fmt.Println(m)
	//	} else {
	//		fmt.Println("no message")
	//	}
	//}

	//m := msgGen1()
	//for {
	//	if m, ok := timeoutWait(m, time.Second*2); ok {
	//		fmt.Println(m)
	//	} else {
	//		fmt.Println("timeout")
	//	}
	//}

	done := make(chan struct{})
	m := msgGen2("service", done)
	for i := 0; i < 5; i++ {
		if m, ok := timeoutWait(m, time.Second); ok {
			fmt.Println(m)
		} else {
			fmt.Println("timeout")
		}
	}
	done <- struct{}{}
	<-done
	time.Sleep(1000 * time.Millisecond)
}
