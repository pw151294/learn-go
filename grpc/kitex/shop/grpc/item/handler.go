package main

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"iflytek.com/weipan4/learn-go/grpc/kitex/shop/kitex_gen/item"
	"iflytek.com/weipan4/learn-go/grpc/kitex/shop/kitex_gen/stock"
	"iflytek.com/weipan4/learn-go/grpc/kitex/shop/kitex_gen/stock/stockservice"
	"log"
)

// ItemServiceImpl implements the last service interface defined in the IDL.
type ItemServiceImpl struct {
	stockCli stockservice.Client
}

func NewStockClient(addr string) (stockservice.Client, error) {
	return stockservice.NewClient("stock", client.WithHostPorts(addr))
}

// GetItem implements the ItemServiceImpl interface.
func (s *ItemServiceImpl) GetItem(ctx context.Context, req *item.GetItemReq) (resp *item.GetItemResp, err error) {
	resp = item.NewGetItemResp()
	resp.Item = item.NewItem()
	resp.Item.Id = req.GetId()
	resp.Item.Title = "Kitex"
	resp.Item.Description = "Kitex is an excellent framwork!"

	stockReq := stock.NewGetItemStockReq()
	stockReq.ItemId = req.GetId()
	stockResp, err := s.stockCli.GetItemStock(ctx, stockReq)
	if err != nil {
		log.Println(err)
		stockResp.Stock = 0
	}
	resp.Item.Stock = stockResp.GetStock()

	return
}
