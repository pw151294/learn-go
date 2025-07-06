package go_redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"iflytek.com/weipan4/learn-go/lock/redislock"
	"iflytek.com/weipan4/learn-go/lock/syncmap"
	"iflytek.com/weipan4/learn-go/storage/redis/go-redis/config"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"
)

const (
	configFile        = "/Users/a123/Downloads/go/learn-go/learn-go/storage/redis/go-redis/config/redis.json"
	gseConnectionHash = "GSE:CONNECTIONS:%s:%d"
	instanceId        = "01fb611f-9a67-70eb-115b-48e0ab3eda68"
	ttl               = 15
	mutexName         = "test-lock"
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
	hash := fmt.Sprintf(gseConnectionHash, "172.30.34.73", 8092)
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

func TestGetKeys(t *testing.T) {
	config.InitConfig(configFile)
	InitRedis()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	keys, _, err := redisCli.ScanType(ctx, 0, "GSE:CONNECTIONS:*", 1000, "hash").Result()
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range keys {
		hKeys, _ := redisCli.HKeys(ctx, key).Result()
		t.Log("gse", key, "agent ids", hKeys)
	}
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

func TestDelKey(t *testing.T) {
	config.InitConfig(configFile)
	InitRedis()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	err := redisCli.Del(ctx, "GSE:CONNECTIONS:172.30.34.73:884").Err()
	if err != nil {
		t.Fatal(err)
	}
}

func TestHSetWithExpire(t *testing.T) {
	// 初始化配置
	config.InitConfig(configFile)
	InitRedis()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	insId := uuid.New().String()
	t.Log("instance id:", insId)

	hash := fmt.Sprintf(gseConnectionHash, "127.0.0.1", 8090)
	expire := time.Now().Add(ttl * time.Second).Format("2006-01-02 15:04:05")
	connMeta := &ConnectionMetadata{
		InstanceId: insId,
		LastPing:   time.Now().Format("2006-01-02 15:04:05"),
		ExpireTime: expire,
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	// 开启定时任务 更新数据
	go func() {
		for {
			select {
			case <-ticker.C:
				connMeta.LastPing = time.Now().Format("2006-01-02 15:04:05")
				if err := hashSet(hash, connMeta, ctx); err != nil {
					t.Errorf("set conection metadata failed: %v", err)
				} else {
					t.Log("refresh data success!")
				}
			case <-ctx.Done():
				t.Error("timeout!")
				return
			}
		}
	}()

	ts := syncmap.NewTimeStorage()
	timer := time.NewTimer(ttl * time.Second)
	ts.Set(insId, timer)
	go func() {
		for {
			select {
			case <-timer.C:
				err := redisCli.HDel(ctx, hash, insId).Err()
				if err != nil {
					return
				}
				ts.Del(insId)
				ticker.Stop()
				t.Log("生命周期已到，删除数据，停止计时器")
			case <-ctx.Done():
				t.Error("timeout!")
				return
			}
		}
	}()

	// 程序优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		return
	case <-ctx.Done():
		return
	}
}

func TestRedisLock(t *testing.T) {
	config.InitConfig(configFile)
	InitRedis()
	rand.NewSource(time.Now().UnixNano())

	rCli := redisCli.(*redis.Client)
	rlb := &redislock.Builder{}
	rl, err := rlb.WithRedisClient(rCli).WithLockName(mutexName).WithExpire(3 * time.Second).Build()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		t.Log("go routine1 try to acquire lock")
		if rl.TryLock() {
			t.Log("go routine1 get redis lock success!")
			time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
			//rl.Unlock()
			t.Log("go routine1 completed")
			return
		}
		t.Log("go routine1 get redis lock failed!")
	}()

	go func() {
		t.Log("go routine2 try to acquire lock")
		if rl.TryLock() {
			t.Log("go routine2 get redis lock success!")
			time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
			//rl.Unlock()
			t.Log("go routine2 completed")
			return
		}
		t.Log("go routine2 get redis lock failed!")
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFunc()
	LOOP:
		for {
			select {
			case <-ticker.C:
				if !rl.TryLock() {
					t.Log("go routine3 get redis lock failed!")
				} else {
					t.Log("go routine3 get redis lock success!")
					break LOOP
				}
			case <-ctx.Done():
				return
			}
		}
		time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
		t.Log("go routine3 completed")
	}()

	// 等待程序运行结束
	time.Sleep(10 * time.Second)
}

func hashSet(hash string, meta *ConnectionMetadata, ctx context.Context) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	bytes, err := json.Marshal(meta)
	if err != nil {
		return err
	}
	val := string(bytes)
	return redisCli.HSet(ctx, hash, meta.InstanceId, val).Err()
}

func TestBatchInsert(t *testing.T) {
	config.InitConfig(configFile)
	InitRedis()

	// 构造数据
	ports := []int{8881, 8882, 8883, 8884}
	connMeta := &ConnectionMetadata{
		LastPing:   time.Now().Format("2006-01-02 15:04:05"),
		ExpireTime: time.Now().Add(ttl * time.Second).Format("2006-01-02 15:04:05"),
	}

	// 插入数据
	for _, port := range ports {
		for i := 0; i < 10; i++ {
			connMeta.InstanceId = uuid.New().String()
			bytes, _ := json.Marshal(connMeta)
			val := string(bytes)
			hash := fmt.Sprintf(gseConnectionHash, "172.30.34.73", port)
			err := redisCli.HSet(context.Background(), hash, connMeta.InstanceId, val).Err()
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestHExpire(t *testing.T) {
	config.InitConfig(configFile)
	InitRedis()

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()

	insId := "2e87dace-54a7-4ce5-865e-5ec082e7d2c4"
	hash := fmt.Sprintf(gseConnectionHash, "172.30.34.73", 8881)
	err := redisCli.HExpire(ctx, hash, 5*time.Second, insId).Err()
	if err != nil {
		t.Error(err)
	}

	time.Sleep(5 * time.Second)

	exists, _ := redisCli.HExists(ctx, hash, insId).Result()
	t.Log(exists)
}
