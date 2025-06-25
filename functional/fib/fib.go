package fib

import (
	"fmt"
	"iflytek.com/weipan4/learn-go/interface/readwrite"
	"io"
	"strings"
)

func Fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

type intGen func() int

func (g intGen) Read(p []byte) (n int, err error) {
	next := g()
	if next > 100 {
		return 0, io.EOF
	}
	s := fmt.Sprintf("%d\n", next)
	r := strings.NewReader(s)
	return r.Read(p)
}

func main() {
	f := Fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
	var fb intGen = Fibonacci()
	readwrite.PrintFileContents(fb)
}
