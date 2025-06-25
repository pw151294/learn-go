package main

import (
	"fmt"
	"reflect"
)

func MyAdd1(a int, b int) int {
	return a + b
}

func MyAdd2(a int, b int) int {
	return a - b
}

func callAdd(f func(a int, b int) int) {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		return
	}
	argv := make([]reflect.Value, 2)
	argv[0] = reflect.ValueOf(1)
	argv[1] = reflect.ValueOf(2)

	res := v.Call(argv)
	fmt.Println(res[0].Int())
}

func main() {
	callAdd(MyAdd1)
	callAdd(MyAdd2)
}
