// Package model 定义视频上传系统的数据库实体。
package model

import "time"

// 上传任务状态
const (
	TaskStatusInit      = "init"      // 初始化
	TaskStatusUploading = "uploading" // 上传中
	TaskStatusCompleted = "completed" // 已完成
	TaskStatusFailed    = "failed"    // 失败
)

// 分片状态
const (
	ChunkStatusPending  = "pending"  // 待上传
	ChunkStatusUploaded = "uploaded" // 已上传
)

// UploadTask 上传任务表：记录每个文件的上传任务。
// file_md5 上建有唯一索引，用于秒传检测和防重复上传。
type UploadTask struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	TaskID      string `gorm:"uniqueIndex;size:36;not null"`
	FileName    string `gorm:"size:512;not null"`
	FileSize    int64  `gorm:"not null"`
	FileMD5     string `gorm:"uniqueIndex;size:32;not null"` // 防重 & 秒传
	UploadID    string `gorm:"size:512;not null;default:''"` // S3 Multipart UploadID
	Bucket      string `gorm:"size:256;not null;default:''"`
	ObjectKey   string `gorm:"size:512;not null;default:''"`
	FileURL     string `gorm:"size:1024;not null;default:''"` // 上传完成后的文件地址
	TotalChunks int    `gorm:"not null;default:0"`
	Status      string `gorm:"size:32;not null;default:'init'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (UploadTask) TableName() string { return "upload_tasks" }

// UploadChunk 上传分片表：记录每个分片的上传状态与完整性信息。
// (task_id, chunk_index) 联合唯一，保证幂等性。
type UploadChunk struct {
	ID         int64  `gorm:"primaryKey;autoIncrement"`
	TaskID     string `gorm:"index;size:36;not null"`
	ChunkIndex int    `gorm:"not null"` // 从 0 开始
	ChunkSize  int64  `gorm:"not null"`
	ChunkMD5   string `gorm:"size:32;not null"`             // 前端上报，用于校验
	ETag       string `gorm:"size:256;not null;default:''"` // S3 返回的 ETag
	Status     string `gorm:"size:32;not null;default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (UploadChunk) TableName() string { return "upload_chunks" }
