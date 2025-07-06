package redislock

import (
	"errors"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisLock struct {
	lock *redsync.Mutex
}

// Builder 使用build模式
type Builder struct {
	redisCli *redis.Client
	lockName string
	expire   time.Duration
}

func (builder *Builder) WithRedisClient(redisCli *redis.Client) *Builder {
	builder.redisCli = redisCli
	return builder
}

func (builder *Builder) WithLockName(lockName string) *Builder {
	builder.lockName = lockName
	return builder
}

func (builder *Builder) WithExpire(expire time.Duration) *Builder {
	builder.expire = expire
	return builder
}

func (builder *Builder) Build() (*RedisLock, error) {
	if builder.redisCli == nil {
		return nil, errors.New("redis client is nil")
	}
	if builder.lockName == "" {
		builder.lockName = "default"
	}
	if builder.expire == 0 {
		builder.expire = 3 * time.Second
	}
	pool := goredis.NewPool(builder.redisCli)
	rs := redsync.New(pool)
	return &RedisLock{
		lock: rs.NewMutex(builder.lockName, redsync.WithExpiry(builder.expire)),
	}, nil
}

func (l *RedisLock) TryLock() bool {
	err := l.lock.TryLock()
	if err != nil {
		return false
	}
	return true
}

func (l *RedisLock) Unlock() {
	ok, err := l.lock.Unlock()
	if err != nil || !ok {
		log.Printf("Redis lock unlock error: %v", err)
		//zap.GetLogger().Error("unlock failed", "message", err)
	}
	log.Printf("Redis lock unlock ok")
	//zap.GetLogger().Info("unlock success!")
}
