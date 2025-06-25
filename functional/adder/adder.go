package main

import "fmt"

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

type IAdder func(int) (int, IAdder)

func iAdder(base int) IAdder {
	return func(v int) (int, IAdder) {
		return base + v, iAdder(base + v)
	}
}

func main() {
	adder := adder()
	for i := 2; i < 100; i++ {
		fmt.Printf("0+1+...+%d=%d\n", i, adder(i))
	}

	iAdder := iAdder(0)
	for i := 0; i < 100; i++ {
		var sum int
		sum, iAdder = iAdder(i)
		fmt.Printf("0+1+...+%d=%d\n", i, sum)
	}
}
