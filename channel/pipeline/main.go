package main

import "fmt"

// 示例：数据处理管道
func stage1(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- 2 * n
		}
		close(out)
	}()
	return out
}

func stage2(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n + 1
		}
		close(out)
	}()
	return out
}

func stage3(in <-chan int) <-chan string {
	out := make(chan string)
	go func() {
		for n := range in {
			out <- fmt.Sprintf("result: %d", n)
		}
		close(out)
	}()

	return out
}

func main() {
	in := make(chan int)
	ppl := stage3(stage2(stage1(in)))

	go func() {
		for i := 0; i < 10; i++ {
			in <- i
		}
		close(in)
	}()
	for l := range ppl {
		fmt.Println(l)
	}
}
