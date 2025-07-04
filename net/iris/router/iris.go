package router

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"iflytek.com/weipan4/learn-go/net/iris/controller"
	"iflytek.com/weipan4/learn-go/net/iris/middlewares"
)

func InitRouter() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("info")
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 注册gse控制器和重定向路由
	gseParty := app.Party("/gse", middlewares.FanOutFanIn)
	gse := mvc.New(gseParty)
	gse.Register(ctx)
	gse.Handle(new(controller.GseController))

	return app
}
