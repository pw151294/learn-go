package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	go_redis "iflytek.com/weipan4/learn-go/storage/redis/go-redis"
	"net/http"
	"time"
)

const (
	gseConnectionPrefix = "GSE:CONNECTION:"
	gseConnectionKey    = "gseConn"
)

// 实现重定向
func redirect(c *gin.Context) {
	// 获取实例id
	params := make(map[string]interface{})
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid params"})
		return
	}
	insId := params["instanceId"].(string)
	if insId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "instance id can not be empty"})
		return
	}

	// 获取实例id所在的gse节点
	redisCli := go_redis.GetClient()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	keys, _, err := redisCli.ScanType(ctx, 0, gseConnectionPrefix+"*", 10000, "hash").Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": fmt.Errorf("get gse connection hash keys failed: %v", err)})
		return
	}
	for _, key := range keys {
		exists, _ := redisCli.HExists(ctx, key, insId).Result()
		if exists {
			c.Set(gseConnectionKey, key)
			c.Next()
		}
	}
	c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Errorf("no gse connection found for agent, instanceId: %s", insId)})
}
