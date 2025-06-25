package main

import (
	"bufio"
	"fmt"
	"iflytek.com/weipan4/learn-go/functional/fib"
	"os"
)

func tryDefer() {
	for i := 0; i < 100; i++ {
		defer fmt.Println(i)
		if i == 30 {
			panic("print too many")
		}
	}
}

func writeFile(filename string) {
	//file, err := os.Create(filename)
	file, err := os.OpenFile(filename, os.O_EXCL|os.O_CREATE, 0666)

	// 自己实现一个新的error
	//err = errors.New("this is a new Error")
	if err != nil {
		if pathError, ok := err.(*os.PathError); !ok {
			panic(err)
		} else {
			fmt.Printf("%s %s,%s", pathError.Op, pathError.Path, pathError.Error())
		}
		//fmt.Println("file already exists")
		//return
		//panic(err)
	}
	defer file.Close()

	f := fib.Fibonacci()
	writer := bufio.NewWriter(file)
	defer writer.Flush() // 将写入内存中的内容导入文件中
	for i := 0; i < 20; i++ {
		fmt.Fprintln(writer, f())
	}
}

func main() {
	tryDefer()
	//writeFile("fib.txt")
}
