package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"iflytek.com/weipan4/learn-go/storage"
	"iflytek.com/weipan4/learn-go/storage/minio/config"
	"log"
	"os"
	"path/filepath"
)

const (
	StorageType  = "Minio"
	BucketPrefix = "nodeman:plugins"
)

type MinioStorage struct {
	Ctx    context.Context
	Client *minio.Client
}

func (ms *MinioStorage) UploadFile(props *storage.PluginPkgUploadProps) error {
	// 获取文件名
	stat, err := os.Stat(props.Filepath)
	if err != nil {
		return err
	}
	objName := stat.Name()

	// 创建bucket
	bktName := fmt.Sprintf("%s:%s:%s", BucketPrefix, props.Os, props.Arch)
	hash := md5.New()
	hash.Write([]byte(bktName))
	hashBytes := hash.Sum(nil)
	bktName = hex.EncodeToString(hashBytes[:])

	exists, err := ms.Client.BucketExists(ms.Ctx, bktName)
	if err != nil {
		return err
	}
	if !exists {
		if err = ms.Client.MakeBucket(ms.Ctx, bktName, minio.MakeBucketOptions{Region: "us-east-1"}); err != nil {
			return err
		}
	}
	// 上传文件
	_, err = ms.Client.FPutObject(ms.Ctx, bktName, objName, props.Filepath, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	return err
}

func (ms *MinioStorage) DownloadFile(props *storage.PluginPkgDownloadProps) error {
	// 创建目录
	if err := os.MkdirAll(props.DestPath, os.ModePerm); err != nil {
		return err
	}
	// 下载文件
	bktName := fmt.Sprintf("%s:%s:%s", BucketPrefix, props.Os, props.Arch)
	hash := md5.New()
	hash.Write([]byte(bktName))
	hashBytes := hash.Sum(nil)
	bktName = hex.EncodeToString(hashBytes[:])

	return ms.Client.FGetObject(ms.Ctx, bktName, props.PkgName,
		filepath.Join(props.DestPath, props.PkgName), minio.GetObjectOptions{})
}

func init() {
	storage.RegisterInitFunc(StorageType, func() storage.PluginPkgStorage {
		client, err := minio.New(config.MinioCfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(config.MinioCfg.AccessKeyID, config.MinioCfg.SecretAccessKey, ""),
			Secure: config.MinioCfg.UseSSL,
		})
		if err != nil {
			log.Fatal(err)
		}
		return &MinioStorage{
			Ctx:    context.Background(),
			Client: client,
		}
	})
}
