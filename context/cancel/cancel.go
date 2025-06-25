package main

import (
	"context"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//创建协程
	f := func(ctx context.Context) chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case dst <- n:
					n++
				case <-ctx.Done():
					return
				}
			}
		}()

		return dst
	}

	for n := range f(ctx) {
		fmt.Println(n)
		if n == 5 {
			cancel()
		}
	}
}
