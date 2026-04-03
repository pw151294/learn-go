// Package storage 封装 Redis 客户端与基于 redsync 的分布式锁。
package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/config"
)

var (
	// RedisClient 是全局 Redis 客户端，可供其他包直接使用（如缓存）。
	RedisClient *goredislib.Client
	// rs 是 redsync 实例，提供分布式互斥锁能力。
	rs *redsync.Redsync
)

// InitRedis 初始化 Redis 连接并创建 redsync 实例。
func InitRedis() error {
	cfg := config.Global
	RedisClient = goredislib.NewClient(&goredislib.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	pool := goredis.NewPool(RedisClient)
	rs = redsync.New(pool)
	log.Printf("[redis] connected: addr=%s", cfg.RedisAddr)
	return nil
}

// AcquireLock 尝试获取以 fileMD5 为键的分布式锁（非阻塞）。
// 锁有效期为 ttl，调用方负责在操作完成后调用 mutex.Unlock() 释放锁。
// 若锁已被其他实例持有，立即返回错误（不轮询等待）。
func AcquireLock(fileMD5 string, ttl time.Duration) (*redsync.Mutex, error) {
	lockKey := fmt.Sprintf("lock:upload:%s", fileMD5)
	mutex := rs.NewMutex(lockKey,
		redsync.WithExpiry(ttl),
		redsync.WithTries(1), // 非阻塞：只尝试一次，失败立即返回
	)
	if err := mutex.Lock(); err != nil {
		return nil, fmt.Errorf("获取锁失败，可能有并发上传请求: %w", err)
	}
	return mutex, nil
}
