package go_redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"iflytek.com/weipan4/learn-go/storage/redis/go-redis/config"
	"strings"
	"time"
)

type RedisType int

const (
	standalone RedisType = iota
	cluster    RedisType = iota
	sentinel   RedisType = iota
)

var redisCli Redis

func GetClient() Redis {
	return redisCli
}

type Redis interface {
	Pipeline() redis.Pipeliner
	// string operator
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Incr(ctx context.Context, key string) *redis.IntCmd
	IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	ScanType(ctx context.Context, cursor uint64, match string, count int64, keyType string) *redis.ScanCmd

	// hset operator
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HSetEX(ctx context.Context, key string, fieldsAndValues ...string) *redis.IntCmd
	HExpire(ctx context.Context, key string, expiration time.Duration, fields ...string) *redis.IntSliceCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	HKeys(ctx context.Context, key string) *redis.StringSliceCmd
	HExists(ctx context.Context, key, field string) *redis.BoolCmd

	// list operator
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	LPop(ctx context.Context, key string) *redis.StringCmd
	RPop(ctx context.Context, key string) *redis.StringCmd
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd
	BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd
	RPopCount(ctx context.Context, key string, count int) *redis.StringSliceCmd
	LLen(ctx context.Context, key string) *redis.IntCmd

	Close() error
	Ping(ctx context.Context) *redis.StatusCmd
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	Watch(ctx context.Context, fn func(*redis.Tx) error, keys ...string) error
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd

	MGet(ctx context.Context, keys ...string) *redis.SliceCmd

	// set operator
	SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
}

func InitRedis() error {
	// 根据模式初始化不同的redis工具类
	switch RedisType(config.Cfg.RedisType) {
	case standalone:
		opts := &redis.Options{
			Addr:     config.Cfg.Addrs,
			Username: config.Cfg.Username,
			Password: config.Cfg.Password,
			DB:       config.Cfg.Db,
		}
		redisCli = redis.NewClient(opts)
	case cluster:
		opts := &redis.ClusterOptions{
			Addrs:    strings.Split(config.Cfg.Addrs, ","),
			Username: config.Cfg.Username,
			Password: config.Cfg.Password,
		}
		redisCli = redis.NewClusterClient(opts)
	case sentinel:
		opts := &redis.FailoverOptions{
			MasterName:       config.Cfg.MasterName,
			SentinelAddrs:    strings.Split(config.Cfg.Addrs, ","),
			Password:         config.Cfg.Password,
			Username:         config.Cfg.Username,
			DB:               config.Cfg.Db,
			SentinelUsername: config.Cfg.SentinelUsername,
			SentinelPassword: config.Cfg.SentinelPassword,
		}
		redisCli = redis.NewFailoverClient(opts)
	default:
		return fmt.Errorf("unsupported redis type: %d", config.Cfg.RedisType)
	}

	if err := redisCli.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("ping redis failed: %v", err)
	}
	return nil
}
