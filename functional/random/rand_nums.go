package main

import (
	"fmt"
	"math/rand"
	"time"
)

type RandomGenerator func() int

func squareRand(n int) RandomGenerator {
	rand.NewSource(time.Now().UnixNano())
	return func() int {
		return n*rand.Intn(n) + rand.Intn(n)
	}
}

func newRandomGenerator(m, n int) RandomGenerator {
	randGen := squareRand(m)
	for m < n {
		randGen = squareRand(m)
		m = m * m
	}

	return func() int {
		res := m
		for res >= n {
			res = randGen()
		}
		return res
	}
}

func main() {
	randGen := newRandomGenerator(3, 4)
	var c0, c1, c2, c3 int
	for i := 0; i < 1000000; i++ {
		switch randGen() {
		case 0:
			c0++
		case 1:
			c1++
		case 2:
			c2++
		case 3:
			c3++
		}
	}
	fmt.Println(c0, c1, c2, c3)
}
