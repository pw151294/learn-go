package main

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"iflytek.com/weipan4/learn-go/logger/zap"
	"iflytek.com/weipan4/learn-go/net/iris/config"
	"iflytek.com/weipan4/learn-go/net/iris/controller"
	go_redis "iflytek.com/weipan4/learn-go/storage/redis/go-redis"
	"net/http"
	"strings"
	"time"
)

const (
	gseConnectionPrefix = "GSE:CONNECTIONS:"
)

func InitRouter() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("info")
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// 注册gse控制器和路由
	gseParty := app.Party("/gse", redirectMiddleware)
	gse := mvc.New(gseParty)
	gse.Register(ctx)
	gse.Handle(new(controller.GseController))

	return app
}

func redirectMiddleware(ctx iris.Context) {
	// 解析出实例id
	params := make(map[string]interface{})
	if err := ctx.ReadJSON(&params); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "invalid params"})
		return
	}
	insId := params["instanceId"].(string)
	if insId == "" {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "agent id can not be empty"})
		return
	}

	// 查询出该实例所对应的gse节点
	redisCli := go_redis.GetClient()
	timeOutCtx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	keys, _, err := redisCli.ScanType(timeOutCtx, 0, gseConnectionPrefix+"*", 10000, "hash").Result()
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": fmt.Errorf("get gse connection keys failed: %v", err)})
		return
	}
	var gseKey string
	for _, key := range keys {
		exists, _ := redisCli.HExists(timeOutCtx, key, insId).Result()
		if exists {
			gseKey = key
			break
		}
	}

	// 根据gseKey的结果判断是否重定向
	if gseKey == "" {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(iris.Map{"error": fmt.Sprintf("can not find gse connected with agent, instanceId: %s", insId)})
		return
	}
	gse := strings.TrimPrefix(gseKey, gseConnectionPrefix)
	if gse == fmt.Sprintf("%s:%d", config.Cfg.Server.Host, config.Cfg.Server.Port) { // agent连接的就是当前节点 直接执行后续逻辑
		ctx.Next()
		return
	}

	// 进行重定向
	redirectURL := fmt.Sprintf("http://%s", gse+ctx.Path())
	zap.GetLogger().Info(fmt.Sprintf("redirect to %s", redirectURL))
	if query := ctx.Request().URL.RawQuery; query != "" {
		redirectURL += "?" + query
	}
	ctx.Redirect(redirectURL, http.StatusTemporaryRedirect)
}
