package controller

import (
	"github.com/kataras/iris/v12"
	"iflytek.com/weipan4/learn-go/net/iris/pkg/resp"
	"math/rand"
	"net/http"
	"time"
)

type GseController struct {
}

func (g *GseController) PostAgent(ctx iris.Context) {
	rand.NewSource(time.Now().UnixNano())
	if rand.Intn(2) == 1 {
		ctx.StatusCode(http.StatusOK)
		ctx.JSON(resp.S("agents operation success!", ""))
	} else {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(resp.E("agent oepration failed", ""))
	}
}
