package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"path/filepath"
	"time"
)

type RedisConfig struct {
	Host        string `json:"host"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Port        int    `json:"port"`
	Database    int    `json:"database"`
	MaxIdle     int    `json:"maxIdle"`
	MaxActive   int    `json:"maxActive"`
	IdleTimeout int    `json:"idleTimeout"`
}

func readConfig() *RedisConfig {
	fn := filepath.Join("storage/redis/redigo", "redis.json")
	bytes, err := os.ReadFile(fn)
	if err != nil {
		log.Fatalf("read redis config file failed: %v", err)
	}
	cfg := &RedisConfig{}
	if err = json.Unmarshal(bytes, cfg); err != nil {
		log.Fatalf("load redis config failed: %v", err)
	}

	return cfg
}

func NewRedisConn() redis.Conn {
	cfg := readConfig()
	pool := &redis.Pool{
		IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, e := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
			if e != nil {
				return nil, e
			}
			return c, e
		},
	}

	return pool.Get()
}

// 执行CONFIG GET maxclients命令
func getMaxClients() (int, error) {
	c := NewRedisConn()
	defer c.Close()

	vals, err := redis.Values(c.Do("CONFIG", "GET", "maxclients"))
	if err != nil {
		return 0, fmt.Errorf("get max clients failed: %v", err)
	}

	var (
		cfgName    string
		maxClients int
	)
	if _, err := redis.Scan(vals, &cfgName, &maxClients); err != nil {
		return 0, fmt.Errorf("parse max clients result failed: %v", err)
	}

	return maxClients, nil
}

func getLogPath() (string, error) {
	c := NewRedisConn()
	defer c.Close()

	vals, err := redis.Values(c.Do("CONFIG", "GET", "logpath"))
	if err != nil {
		return "", fmt.Errorf("get logpath failed: %v", err)
	}

	var (
		cfgName string
		logPath string
	)
	if _, err := redis.Scan(vals, &cfgName, &logPath); err != nil {
		return "", fmt.Errorf("parse logpath result failed: %v", err)
	}

	return logPath, nil
}

func main() {
	maxClients, err := getMaxClients()
	if err != nil {
		log.Fatalf("get max clients failed: %v", err)
	}
	log.Printf("max clients: %d", maxClients)

	logPath, err := getLogPath()
	if err != nil {
		log.Fatalf("get logpath failed: %v", err)
	}
	log.Printf("logpath: %s", logPath)
}
