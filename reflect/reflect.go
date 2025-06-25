package main

import (
	"fmt"
	"reflect"
)

func main() {
	s := "hello world"
	st := reflect.TypeOf(s)
	fmt.Println(st)
	sv := reflect.ValueOf(s)
	fmt.Println(sv)

	//根据类型和值还原数据
	newS := sv.Interface().(string)
	fmt.Println(newS)
}
