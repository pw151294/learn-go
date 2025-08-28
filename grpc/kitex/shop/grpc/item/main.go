package main

import (
	"iflytek.com/weipan4/learn-go/grpc/kitex/shop/kitex_gen/item/itemservice"
	"log"
)

const stockServiceAddr = "localhost:8890"

func main() {
	itemServiceImpl := new(ItemServiceImpl)
	stockCli, err := NewStockClient(stockServiceAddr)
	if err != nil {
		log.Fatalln(err)
	}
	itemServiceImpl.stockCli = stockCli

	svr := itemservice.NewServer(itemServiceImpl)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
