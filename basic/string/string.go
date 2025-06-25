package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	s := "imooc"
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	fmt.Printf("%q, len=%d\n", sh.Data, sh.Len)
}
