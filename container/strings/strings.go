package main

import (
	"fmt"
	"runtime"
	"unicode/utf8"
)

func main() {
	fmt.Println(runtime.GOARCH)

	s := "YES星际可观测"
	fmt.Println(len(s))

	fmt.Printf("%X\n", []byte(s))

	for _, b := range []byte(s) {
		fmt.Printf("%X ", b)
	}

	for i, ch := range s {
		fmt.Printf("(%d %X)\n", i, ch)
	}

	fmt.Printf("rune count: %d\n", utf8.RuneCount([]byte(s)))
	fmt.Printf("rune count in string: %d\n", utf8.RuneCountInString(s))

	bytes := []byte(s)
	for len(bytes) > 0 {
		ch, size := utf8.DecodeRune(bytes)
		bytes = bytes[size:]
		fmt.Printf("%c ", ch)
	}
	fmt.Println()
	for i, r := range []rune(s) {
		fmt.Printf("%d %c\n", i, r)
	}
}
