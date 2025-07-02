package controller

import (
	"github.com/kataras/iris/v12"
	"net/http"
)

type GseController struct {
}

func (g *GseController) PostAgent(ctx iris.Context) {
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(iris.Map{"message": "receive request!"})
}
