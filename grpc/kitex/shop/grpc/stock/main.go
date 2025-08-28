package main

import (
	"github.com/cloudwego/kitex/server"
	stock "iflytek.com/weipan4/learn-go/grpc/kitex/shop/kitex_gen/stock/stockservice"
	"log"
	"net"
)

const stockServiceAddr = "localhost:8890"

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", stockServiceAddr)
	svr := stock.NewServer(new(StockServiceImpl), server.WithServiceAddr(addr))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
