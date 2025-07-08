package main

import (
	"context"
	"encoding/json"
	"fmt"
	mysql "iflytek.com/weipan4/learn-go/storage/mysql/xorm/config"
	"iflytek.com/weipan4/learn-go/storage/mysql/xorm/datasource"
	"iflytek.com/weipan4/learn-go/storage/mysql/xorm/model"
	go_redis "iflytek.com/weipan4/learn-go/storage/redis/go-redis"
	redis "iflytek.com/weipan4/learn-go/storage/redis/go-redis/config"
	"testing"
	"time"
)

const (
	configFile = "/Users/a123/Downloads/go/learn-go/learn-go/storage/mysql/xorm/config/mysq.json"
	redisFile  = "/Users/a123/Downloads/go/learn-go/learn-go/storage/redis/go-redis/config/redis.json"
	port       = 8881
)

func TestBuildRequest(t *testing.T) {
	// 初始化MySQL和Redis数据源
	redis.InitConfig(redisFile)
	mysql.InitConfig(configFile)
	go_redis.InitRedis()
	datasource.InitEngine()
	datasource.GetEngine().Sync2(new(model.Node))

	// 查询出所有的节点实例
	nodeIds := []int64{99, 93, 5, 87, 53, 97, 7, 59, 55, 101, 69, 83}
	repository := model.NewNodeRepository()
	nodes, err := repository.SelectByIds(nodeIds)
	if err != nil {
		t.Fatal(err)
	}
	insIds := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if node.InstanceId.Valid {
			insIds = append(insIds, node.InstanceId.String)
		}
	}
	bytes, _ := json.Marshal(insIds)
	t.Logf("instance ids: %s", string(bytes))

	now := time.Now()
	redisCli := go_redis.GetClient()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	for i := range insIds {
		insId := insIds[i]
		hash := fmt.Sprintf("GSE:CONNECTIONS:172.30.34.73:%d", port+(i%4))
		connMeta := struct {
			InstanceId string `json:"instanceId"`
			LastPing   string `json:"lastPing"`
			ExpireTime string `json:"expireTime"`
		}{
			InstanceId: insId,
			LastPing:   now.Format("2006-01-02 15:04:05"),
			ExpireTime: now.Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
		}
		bytes, err := json.Marshal(connMeta)
		if err != nil {
			t.Errorf("marshal connMeta failed: %v", err)
			continue
		}
		val := string(bytes)
		err = redisCli.HSet(ctx, hash, insId, val).Err()
		if err != nil {
			t.Errorf("hset connMeta failed: %v", err)
		}
	}
}

func TestSelectNodes(t *testing.T) {
	mysql.InitConfig(configFile)
	datasource.InitEngine()
	datasource.GetEngine().Sync2(new(model.Node))

	nodeIds := []int64{99, 93, 5, 87, 53, 97, 7, 59, 55, 101, 69, 83}
	repository := model.NewNodeRepository()
	nodes, err := repository.SelectByIds(nodeIds)
	if err != nil {
		t.Fatal(err)
	}

	// 构建出映射关系
	type InstallHost struct {
		NodeId int64 `json:"node_id"`
	}
	hosts := make([]*InstallHost, 0)
	insIds := make([]string, 0)
	for _, node := range nodes {
		hosts = append(hosts, &InstallHost{NodeId: node.Id})
		if node.InstanceId.Valid {
			insIds = append(insIds, node.InstanceId.String)
		}
	}
	bytes, err := json.Marshal(hosts)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("request: %s", string(bytes))
}

func TestSelectByIds(t *testing.T) {
	mysql.InitConfig(configFile)
	datasource.InitEngine()
	datasource.GetEngine().Sync2(new(model.Node))

	nodeIds := []int64{99, 93, 5, 87, 53, 97, 7, 59, 55, 101, 69, 83}
	nodeList := make([]*model.Node, 0)
	err := datasource.GetEngine().Table(new(model.Node)).
		Select("id, instance_id").In("id", nodeIds).Find(&nodeList)
	if err != nil {
		t.Fatal(err)
	}

	bytes, _ := json.Marshal(nodeList)
	t.Logf("request: %s", string(bytes))
}
