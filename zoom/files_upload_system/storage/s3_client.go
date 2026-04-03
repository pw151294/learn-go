// Package storage 封装 MinIO/S3 分片上传客户端。
// 使用 minio.Core（而非 minio.Client）以访问底层分片上传 API。
package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"iflytek.com/weipan4/learn-go/zoom/files_upload_system/config"
)

// S3Client 包装 MinIO Core 客户端，提供分片上传所需的四个原语：
// CreateMultipartUpload / UploadPart / CompleteMultipartUpload / AbortMultipartUpload。
// 使用 minio.Core 是因为 minio.Client 将分片上传方法设为私有；
// Core 是官方提供的底层访问入口，适合应用层自行管理分片。
type S3Client struct {
	core   *minio.Core
	client *minio.Client // 用于 BucketExists / MakeBucket 等管理操作
	bucket string
	ctx    context.Context
}

// DefaultS3Client 是全局 S3 客户端实例，由 InitS3Client() 初始化。
var DefaultS3Client *S3Client

// InitS3Client 初始化 MinIO/S3 客户端并确保目标 bucket 已存在。
func InitS3Client() error {
	cfg := config.Global
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.S3AccessKeyID, cfg.S3SecretAccessKey, ""),
		Secure: cfg.S3UseSSL,
	}

	// 普通客户端：用于 bucket 管理
	client, err := minio.New(cfg.S3Endpoint, opts)
	if err != nil {
		return fmt.Errorf("minio.New: %w", err)
	}

	// Core 客户端：暴露底层分片上传 API
	core, err := minio.NewCore(cfg.S3Endpoint, opts)
	if err != nil {
		return fmt.Errorf("minio.NewCore: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.S3Bucket)
	if err != nil {
		return fmt.Errorf("BucketExists: %w", err)
	}
	if !exists {
		if err = client.MakeBucket(ctx, cfg.S3Bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("MakeBucket: %w", err)
		}
		log.Printf("[s3] created bucket: %s", cfg.S3Bucket)
	}

	DefaultS3Client = &S3Client{core: core, client: client, bucket: cfg.S3Bucket, ctx: ctx}
	log.Printf("[s3] connected: endpoint=%s bucket=%s", cfg.S3Endpoint, cfg.S3Bucket)
	return nil
}

// Bucket 返回当前使用的 bucket 名称。
func (s *S3Client) Bucket() string { return s.bucket }

// CreateMultipartUpload 向 S3 发起分片上传请求，返回 UploadID。
// objectKey 格式建议：videos/{task_id}/{file_name}
func (s *S3Client) CreateMultipartUpload(objectKey string) (string, error) {
	uploadID, err := s.core.NewMultipartUpload(s.ctx, s.bucket, objectKey, minio.PutObjectOptions{
		ContentType: "video/mp4",
	})
	if err != nil {
		return "", fmt.Errorf("NewMultipartUpload: %w", err)
	}
	return uploadID, nil
}

// UploadPart 上传单个分片，返回 S3 响应的 ETag。
// chunkIndex 为前端上报的分片序号（从 0 开始），
// 内部转换为 partNumber（从 1 开始）以符合 S3 规范。
func (s *S3Client) UploadPart(objectKey, uploadID string, chunkIndex int, data io.Reader, size int64) (string, error) {
	part, err := s.core.PutObjectPart(
		s.ctx, s.bucket, objectKey, uploadID,
		chunkIndex+1, // S3 partNumber 从 1 开始
		data, size,
		minio.PutObjectPartOptions{},
	)
	if err != nil {
		return "", fmt.Errorf("PutObjectPart(index=%d): %w", chunkIndex, err)
	}
	return part.ETag, nil
}

// CompleteMultipartUpload 通知 S3 合并所有已上传分片，返回对象访问 URL。
func (s *S3Client) CompleteMultipartUpload(objectKey, uploadID string, parts []minio.CompletePart) (string, error) {
	_, err := s.core.CompleteMultipartUpload(s.ctx, s.bucket, objectKey, uploadID, parts, minio.PutObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("CompleteMultipartUpload: %w", err)
	}
	fileURL := fmt.Sprintf("s3://%s/%s", s.bucket, objectKey)
	return fileURL, nil
}

// AbortMultipartUpload 中止并清理未完成的分片上传（用于失败回滚）。
func (s *S3Client) AbortMultipartUpload(objectKey, uploadID string) error {
	return s.core.AbortMultipartUpload(s.ctx, s.bucket, objectKey, uploadID)
}
