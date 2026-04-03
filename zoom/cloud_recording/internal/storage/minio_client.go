package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/configs"
)

type MinIOClient struct {
	client *minio.Client
	cfg    *configs.MinIOConfig
}

func NewMinIOClient(cfg *configs.MinIOConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("init minio client: %w", err)
	}
	return &MinIOClient{client: client, cfg: cfg}, nil
}

func (m *MinIOClient) EnsureBuckets(ctx context.Context) error {
	buckets := []string{m.cfg.HotBucket, m.cfg.WarmBucket, m.cfg.ColdBucket}
	for _, bucket := range buckets {
		exists, err := m.client.BucketExists(ctx, bucket)
		if err != nil {
			return fmt.Errorf("check bucket %s: %w", bucket, err)
		}
		if !exists {
			if err := m.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				return fmt.Errorf("create bucket %s: %w", bucket, err)
			}
		}
	}
	return nil
}

func (m *MinIOClient) UploadToHot(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(ctx, m.cfg.HotBucket, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("upload to hot bucket: %w", err)
	}
	return nil
}

func (m *MinIOClient) GetPresignedURL(ctx context.Context, bucket, objectKey string, expiry time.Duration) (string, error) {
	u, err := m.client.PresignedGetObject(ctx, bucket, objectKey, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("presign object %s/%s: %w", bucket, objectKey, err)
	}
	return u.String(), nil
}

func (m *MinIOClient) MigrateObject(ctx context.Context, srcBucket, dstBucket, objectKey string) error {
	dst := minio.CopyDestOptions{Bucket: dstBucket, Object: objectKey}
	src := minio.CopySrcOptions{Bucket: srcBucket, Object: objectKey}
	if _, err := m.client.CopyObject(ctx, dst, src); err != nil {
		return fmt.Errorf("copy object %s -> %s: %w", srcBucket, dstBucket, err)
	}
	if err := m.client.RemoveObject(ctx, srcBucket, objectKey, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("remove source object %s/%s: %w", srcBucket, objectKey, err)
	}
	return nil
}

func (m *MinIOClient) DeleteObject(ctx context.Context, bucket, objectKey string) error {
	if err := m.client.RemoveObject(ctx, bucket, objectKey, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("delete object %s/%s: %w", bucket, objectKey, err)
	}
	return nil
}
