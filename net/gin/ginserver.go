package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	const requestId = "requestId"
	r.Use(func(c *gin.Context) {
		s := time.Now()

		c.Next()

		// log latency, response code
		logger.Info("incoming request",
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", time.Now().Sub(s)))
	}, func(c *gin.Context) {
		c.Set(requestId, rand.Int())

		c.Next()
	})

	r.GET("/test", func(c *gin.Context) {
		h := gin.H{
			"message": "hello world",
		}
		if rid, ok := c.Get(requestId); ok {
			h[requestId] = rid
		}
		c.JSON(http.StatusOK, h)
	})
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, "hello world")
	})

	r.Run(":8080")
}
