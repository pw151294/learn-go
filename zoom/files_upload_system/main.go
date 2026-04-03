// 极简版视频上传系统入口。
// 启动顺序：加载配置 → 初始化各基础设施 → 注册路由 → 启动 HTTP 服务。
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/api"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/config"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/model"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/mq"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/repository"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/service"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/storage"
)

func main() {
	// ── 1. 加载配置（从环境变量） ──────────────────────────────
	config.Load()

	// ── 2. 初始化 PostgreSQL & 自动迁移表结构 ─────────────────
	db, err := gorm.Open(postgres.Open(config.Global.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("[main] connect postgres failed: %v", err)
	}
	if err = db.AutoMigrate(&model.UploadTask{}, &model.UploadChunk{}); err != nil {
		log.Fatalf("[main] auto migrate failed: %v", err)
	}
	log.Println("[main] postgres: connected and migrated")

	// ── 3. 初始化 Redis（分布式锁） ───────────────────────────
	if err = storage.InitRedis(); err != nil {
		log.Fatalf("[main] init redis failed: %v", err)
	}

	// ── 4. 初始化 S3/MinIO 客户端 ────────────────────────────
	if err = storage.InitS3Client(); err != nil {
		log.Fatalf("[main] init s3 failed: %v", err)
	}

	// ── 5. 初始化 Kafka 生产者 ────────────────────────────────
	if err = mq.InitKafkaProducer(); err != nil {
		log.Fatalf("[main] init kafka failed: %v", err)
	}
	defer mq.Close()

	// ── 6. 构建依赖注入链 ─────────────────────────────────────
	taskRepo := repository.NewUploadTaskRepo(db)
	chunkRepo := repository.NewChunkRepo(db)
	uploadSvc := service.NewUploadService(taskRepo, chunkRepo, storage.DefaultS3Client)
	uploadHandler := api.NewUploadHandler(uploadSvc)

	// ── 7. 启动 Gin HTTP 服务 ─────────────────────────────────
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	uploadHandler.RegisterRoutes(r)

	addr := ":" + config.Global.Port
	log.Printf("[main] server starting on %s", addr)
	if err = r.Run(addr); err != nil {
		log.Fatalf("[main] server error: %v", err)
	}
}
