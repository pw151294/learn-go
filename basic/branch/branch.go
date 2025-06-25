package main

import (
	"fmt"
	"io/ioutil"
)

func readFile() {
	const filename = "/Users/a123/Downloads/ai/agent/kecheng/asset/test.txt"
	if contents, err := ioutil.ReadFile(filename); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", contents)
	}
}

func grade(score int) string {
	g := ""
	switch {
	case score < 0 || score > 100:
		panic(fmt.Sprint("bad score:", score))
	case score < 60:
		g = "F"
	case score < 80:
		g = "C"
	case score < 90:
		g = "B"
	default:
		g = "A"
	}

	return g
}

func eval(a, b int, op string) int {
	var result int
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		result = a / b
	default:
		panic("unsupported operation: " + op)
	}
	return result
}

func main() {
	readFile()
	const filename = "/Users/a123/Downloads/ai/agent/kecheng/asset/letter.txt"
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", contents)
	}
	fmt.Printf("result: %d", eval(3, 4, "+"))

	fmt.Println(grade(1000))
}
