package main

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
	"iflytek.com/weipan4/learn-go/grpc/kitex/shop/kitex_gen/item"
	"iflytek.com/weipan4/learn-go/grpc/kitex/shop/kitex_gen/item/itemservice"
	"log"
	"net/http"
	"time"
)

const (
	itemServiceName = "item"
	itemServiceAddr = "0.0.0.0:8888"
	serverAddr      = "localhost:8889"
)

var cli itemservice.Client

func main() {
	c, err := itemservice.NewClient(itemServiceName, client.WithHostPorts(itemServiceAddr))
	if err != nil {
		log.Fatalln(err)
	}
	cli = c

	hz := server.New(server.WithHostPorts(serverAddr))

	hz.GET("/api/item", Handler)

	if err := hz.Run(); err != nil {
		log.Fatalln(err)
	}
}

func Handler(ctx context.Context, c *app.RequestContext) {
	req := item.NewGetItemReq()
	req.Id = 1024
	resp, err := cli.GetItem(ctx, req, callopt.WithRPCTimeout(3*time.Second))
	if err != nil {
		log.Fatalln()
	}

	c.String(http.StatusOK, resp.String())
}
