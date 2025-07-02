package go_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"iflytek.com/weipan4/learn-go/storage/redis/go-redis/config"
	"strings"
	"testing"
	"time"
)

const (
	configFile        = "/Users/a123/Downloads/go/learn-go/learn-go/storage/redis/go-redis/config/redis.json"
	gseConnectionHash = "GSE:CONNECTIONS:%s:%d"
	instanceId        = "01fb611f-9a67-70eb-115b-48e0ab3eda68"
)

type ConnectionMetadata struct {
	InstanceId string `json:"instanceId"`
	LastPing   string `json:"lastPing"`
	ExpireTime string `json:"expireTime"`
}

type ConnectionMetadataOptions func(*ConnectionMetadata)

func NewConnectionMetadata(opts ...ConnectionMetadataOptions) *ConnectionMetadata {
	connMeta := &ConnectionMetadata{
		LastPing:   time.Now().Format("2006-01-02 15:04:05"),
		ExpireTime: time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(connMeta)
		}
	}
	if connMeta.InstanceId == "" {
		connMeta.InstanceId = uuid.New().String()
	}

	return connMeta
}

func WithInstanceId(instanceId string) ConnectionMetadataOptions {
	return func(o *ConnectionMetadata) {
		o.InstanceId = instanceId
	}
}

func WithExpireTime(expireTime string) ConnectionMetadataOptions {
	return func(o *ConnectionMetadata) {
		if _, err := time.ParseDuration(expireTime); err == nil {
			o.ExpireTime = expireTime
		}
	}
}

func WithLastPing(lastPing string) ConnectionMetadataOptions {
	return func(o *ConnectionMetadata) {
		if _, err := time.ParseDuration(lastPing); err == nil {
			o.LastPing = lastPing
		}
	}
}

type InstanceMetadata struct {
	InstanceId        string `json:"instanceId"`
	Ip                string `json:"ip"`
	Port              int    `json:"port"`
	Status            string `json:"status"`
	ActiveConnections int    `json:"activeConnections"`
}

type InstanceStatus int

const (
	inactive InstanceStatus = iota
	active   InstanceStatus = iota
)

func NewInstanceMetadata(opts ...InstanceMetadataOptions) *InstanceMetadata {
	instMeta := &InstanceMetadata{
		Ip:                "127.0.0.1",
		Port:              8090,
		Status:            "active",
		ActiveConnections: 200,
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(instMeta)
		}
	}

	return instMeta
}

type InstanceMetadataOptions func(*InstanceMetadata)

func WithId(id string) InstanceMetadataOptions {
	return func(o *InstanceMetadata) {
		if id != "" {
			o.InstanceId = id
		}
	}
}

func WithIp(ip string) InstanceMetadataOptions {
	return func(o *InstanceMetadata) {
		if ip != "" {
			o.Ip = ip
		}
	}
}

func WithPort(port int) InstanceMetadataOptions {
	return func(o *InstanceMetadata) {
		if port > 0 {
			o.Port = port
		}
	}
}

func WithStatus(status InstanceStatus) InstanceMetadataOptions {
	return func(o *InstanceMetadata) {
		switch status {
		case inactive:
			o.Status = "inactive"
		case active:
			o.Status = "active"
		}
	}
}

func WithActiveConnections(activeConnections int) InstanceMetadataOptions {
	return func(o *InstanceMetadata) {
		if activeConnections > 0 {
			o.ActiveConnections = activeConnections
		}
	}
}

func TestHSet(t *testing.T) {
	// 初始化配置
	config.InitConfig(configFile)
	InitRedis()

	// 构建数据
	hash := fmt.Sprintf(gseConnectionHash, "127.0.0.1", 8092)
	expire := time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05")
	s := struct {
		InstanceId string `json:"instanceId"`
		LastPing   string `json:"lastPing"`
		ExpireTime string `json:"expireTime"`
	}{
		InstanceId: instanceId,
		LastPing:   time.Now().Format("2006-01-02 15:04:05"),
		ExpireTime: expire,
	}
	bytes, _ := json.Marshal(s)
	val := string(bytes)

	// 保存数据
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err := redisCli.HSet(ctx, hash, instanceId, val).Err()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("set data to redis success!", "hash", hash, "key", instanceId, "val", val)
}

func TestHGet(t *testing.T) {
	config.InitConfig(configFile)
	InitRedis()

	hash := fmt.Sprintf(gseConnectionHash, "127.0.0.1", 8090)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	ssMap, err := redisCli.HGetAll(ctx, hash).Result()
	if err != nil {
		t.Fatal(err)
	}
	s := ssMap[instanceId]
	connMeta := &ConnectionMetadata{}
	if err = json.NewDecoder(strings.NewReader(s)).Decode(connMeta); err != nil {
		t.Fatal(err)
	}

	t.Log("read connection metadata from redis success!")
}

func TestHDel(t *testing.T) {
	config.InitConfig(configFile)
	InitRedis()

	hash := fmt.Sprintf(gseConnectionHash, "127.0.0.1", 8090)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	err := redisCli.HDel(ctx, hash, instanceId).Err()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("delete connection metadata from redis success")
}

func TestScanType(t *testing.T) {
	config.InitConfig(configFile)
	InitRedis()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	keys, cursor, err := redisCli.ScanType(ctx, 0, "GSE:CONNECTIONS:*", 10000, "hash").Result()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("keys", keys, "cursor", cursor)
}

func TestFindNode(t *testing.T) {
	// 根据instanceId查询对应的gse节点
	config.InitConfig(configFile)
	InitRedis()
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	keys, _, err := redisCli.ScanType(ctx, 0, "GSE:CONNECTIONS:*", 10000, "hash").Result()
	if err != nil {
		t.Fatal(err)
	}

	var gseKey string
	for _, key := range keys {
		err := redisCli.HExists(ctx, key, instanceId).Err()
		exists, err := redisCli.HExists(ctx, key, instanceId).Result()
		if err != nil {
			t.Fatal(err)
		}
		if exists {
			gseKey = key
			break
		}
	}
	t.Log(gseKey)
}
