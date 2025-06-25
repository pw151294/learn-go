package main

import (
	"fmt"
	"time"
)

func tryRecover() {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			fmt.Println("Error occurred: ", err)
		} else {
			panic(fmt.Sprintf("I don't know what to do: %v", r))
		}
	}()
	//panic(errors.New("something bad happened"))
	b := 0
	a := 5 / b
	fmt.Println(a)
}

func main() {
	//tryRecover()
	defer fmt.Println("defer main g")

	go func() {
		//defer fmt.Println("defer g")
		defer func() {
			recover() //恢复程序的运行
		}()
		panic("panic")
		fmt.Println("end g")
	}()

	time.Sleep(time.Second)
	fmt.Println("end main g")
}
