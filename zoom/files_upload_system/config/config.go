// Package config 负责从环境变量加载应用配置。
package config

import (
	"log"
	"os"
	"strconv"
)

// Config 汇聚所有基础设施的连接参数。
type Config struct {
	// HTTP 服务
	Port string

	// PostgreSQL DSN
	DBDSN string

	// Redis
	RedisAddr     string
	RedisPassword string

	// MinIO / S3 兼容对象存储
	S3Endpoint        string
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3Bucket          string
	S3UseSSL          bool

	// Kafka
	KafkaBrokers string // 逗号分隔，如 "localhost:9092,localhost:9093"
	KafkaTopic   string
}

// Global 是全局配置实例，由 Load() 初始化。
var Global Config

// Load 从环境变量读取配置，未设置时使用默认值（适合本地开发）。
func Load() {
	Global = Config{
		Port:  getEnv("PORT", "8080"),
		DBDSN: getEnv("DB_DSN", "host=localhost user=postgres password=postgres dbname=videoupload port=5432 sslmode=disable"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		S3Endpoint:        getEnv("S3_ENDPOINT", "localhost:9000"),
		S3AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", "minioadmin"),
		S3SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", "minioadmin"),
		S3Bucket:          getEnv("S3_BUCKET", "video-uploads"),
		S3UseSSL:          getEnvBool("S3_USE_SSL", false),

		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "video.transcode"),
	}
	log.Printf("[config] port=%s db=%s redis=%s s3=%s kafka=%s",
		Global.Port, Global.DBDSN, Global.RedisAddr, Global.S3Endpoint, Global.KafkaBrokers)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
