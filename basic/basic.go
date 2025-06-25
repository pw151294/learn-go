package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

var aa = 3
var ss = "kkk"
var bb = true

var (
	aaa = 3
	sss = "kkk"
	bbb = true
)

func variable() {
	var a int
	var s string
	fmt.Printf("%d %q\n", a, s) //%q进行格式化 可以将空字符串打印出来
}

func variableInitValue() {
	var a, b int = 3, 5
	var s string = "hello"
	fmt.Printf("%d %d %q\n", a, b, s)
}

func variableTypeDeduction() {
	var a, b, c, s = 3, 4, true, "hello"
	fmt.Printf("%d %d %t %q\n", a, b, c, s)
}

func variableShorter() {
	a, b, c, s := 3, 4, true, "hello"
	b = 5
	fmt.Println(a, b, c, s)
}

func euler() {
	println(cmplx.Pow(math.E, 1i*math.Pi) + 1)
	println(cmplx.Exp(1i*math.Pi) + 1)
}

func triangle() {
	var a, b int = 3, 4
	var c int
	c = int(math.Sqrt(float64(a*a + b*b)))
	fmt.Println(c)
}

const filename = "test.txt"

func enums() {
	const (
		cpp = iota
		_
		python
		golang
		javascript
	)
	fmt.Println(cpp, python, golang, javascript)

	const (
		b = 1 << (10 * iota)
		kb
		mb
		gb
		tb
		pb
	)
	fmt.Println(b, kb, mb, gb, tb, pb)
}

func constants() {
	const a, b = 3, 4
	const (
		s = "hello"
		d = true
	)
	var c int
	c = int(math.Sqrt(a*a + b*b))
	fmt.Println(filename, c)
}

func main() {
	fmt.Println("Hello World")
	variable()
	variableInitValue()
	variableTypeDeduction()
	variableShorter()
	fmt.Println(aa, ss, bb)
	fmt.Println(aaa, sss, bbb)
	euler()
	triangle()
	constants()
	enums()
}
