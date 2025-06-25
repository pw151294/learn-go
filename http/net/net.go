package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var buf [512]byte

	for {
		_, err := conn.Read(buf[0:])
		if err != nil {
			break
		}
		fmt.Printf("received message: %s\n", string(buf[0:]))

		_, err = conn.Write(buf[0:])
		if err != nil {
			break
		}
	}
}

func listenAndServe(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("listen failed, address:%s, err:%v\n", address, err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("accept failed, err:%v\n", err)
			return
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	//使用bufio标准库提供的缓冲区功能
	reader := bufio.NewReader(conn)
	for {
		//ReadString('\n')会一直阻塞直到遇见分隔符'\n'
		//遇到分隔符之后会返回上次遇到分隔符或连接建立后收到的所有数据，包括分隔符本身
		//如果遇见异常，ReadString('\n')会返回已经接收到的所有数据 和错误信息
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF { //通常遇到的错误是连接中断或者关闭 用io.EOF来表示
				log.Print("connection closed")
			} else {
				log.Printf("read message failed, message:%s, err:%v\n", msg, err)
			}
			break
		}
		fmt.Printf("received message: %s\n", msg)
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Printf("write message failed, err:%v\n", err)
			break
		}
	}
}

func main() {
	listenAndServe(":8888")
}
