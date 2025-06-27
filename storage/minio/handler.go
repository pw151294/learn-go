package main

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"iflytek.com/weipan4/learn-go/storage/minio/config"
)

var fHandler *fileHandler

type fileHandler struct {
	Ctx    context.Context
	Client *minio.Client
}

func InitFileHandler(ctx context.Context) error {
	client, err := minio.New(config.MinioCfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioCfg.AccessKeyID, config.MinioCfg.SecretAccessKey, ""),
		Secure: config.MinioCfg.UseSSL,
	})
	if err != nil {
		return err
	}

	fh := &fileHandler{
		Ctx:    ctx,
		Client: client,
	}
	fHandler = fh
	return nil
}

func (fh *fileHandler) doUpload(info ObjectInfo) (minio.UploadInfo, error) {
	ui, err := fh.Client.FPutObject(fh.Ctx, info.BucketName, info.ObjectName, info.FilePath,
		minio.PutObjectOptions{ContentType: info.ContentType})
	return ui, err
}

func (fh *fileHandler) doDownload(info ObjectInfo) error {
	return fh.Client.FGetObject(fh.Ctx, info.BucketName, info.ObjectName, info.FilePath, minio.GetObjectOptions{})
}

func (fh *fileHandler) doDelete(info ObjectInfo) {

}

func (fh *fileHandler) doCrete(bucketName string) error {
	if exists, err := fh.Client.BucketExists(fh.Ctx, bucketName); err == nil && exists {
		return nil
	}
	return fh.Client.MakeBucket(fh.Ctx, bucketName, minio.MakeBucketOptions{Region: location})
}
