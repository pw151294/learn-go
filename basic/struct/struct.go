package main

import (
	"fmt"
	"unsafe"
)

type People struct {
	Name string
}

func (p *People) walk() {
	fmt.Println("walking")
}

type Student struct {
	People
	Grade int
}

type K struct {
}

type F struct {
	k   K
	num int32
}

type S1 struct {
	num1 int32
	num2 int32
}

type S2 struct {
	num1 int16
	num2 int32
}

type User struct {
	A int32
	B []int32
	C string
	D bool
	E struct{}
}

func main() {
	student := new(Student)
	student.walk()
	fmt.Printf("%v\n", unsafe.Sizeof(int32(0)))
	i := 0
	p := &i
	fmt.Printf("size of pointer:%d\n", unsafe.Sizeof(p))

	// 空结构体 有指针但是没有长度
	a := K{}
	fmt.Printf("size of null struct:%d\n", unsafe.Sizeof(a))
	fmt.Printf("size of pointer of null struct:%d\n", unsafe.Sizeof(&a))
	fmt.Printf("address of null struct:%p\n", &a)

	// 所有独立的空结构体指针都指向同样的地址
	a1 := K{}
	a2 := 0
	a3 := K{}
	fmt.Printf("address of a1: %p\n", &a1)
	fmt.Printf("address of a2: %p\n", &a2)
	fmt.Printf("address of a3: %p\n", &a3)

	// 非独立的空结构体指针不会指向同样的地址
	f := F{}
	fmt.Printf("address of k: %p\n", &f.k)

	s1 := new(S1)
	s2 := new(S2)
	fmt.Printf("size of structs, s1:%d, s2:%d\n", unsafe.Sizeof(s1), unsafe.Sizeof(s2))

	user := new(User)
	fmt.Printf("size of user:%d, align of user:%d\n", unsafe.Sizeof(user), unsafe.Alignof(user))
}
