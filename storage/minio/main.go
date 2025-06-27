package main

import (
	"context"
	"flag"
	"github.com/minio/minio-go/v7"
	"iflytek.com/weipan4/learn-go/storage/minio/config"
	"log"
)

var configFile = flag.String("configFile", "storage/minio/config/config.toml", "path of minio config file")

const (
	bucketName   = "weipan4-bucket"
	location     = "us-east-1"
	objectName   = "install.pdf"
	filePath     = "/Users/a123/Downloads/working/2025.6/cmdb/discover"
	downloadPath = "/Users/a123/Downloads/working/2025.6/minio/auto_discovery"
)

func UploadFile(info ObjectInfo) (minio.UploadInfo, error) {
	return fHandler.doUpload(info)
}

func DownloadFile(info ObjectInfo) error {
	return fHandler.doDownload(info)
}

func main() {
	err := config.InitConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	if err = InitFileHandler(ctx); err != nil {
		log.Fatal(err)
	}

	objInfo := NewObjectInfo(
		WithBucketName(bucketName),
		WithObjectName(objectName),
		WithFilePath(filePath))
	fi, err := UploadFile(*objInfo)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("upload file %s successfully with %d bytes", objInfo.FilePath, fi.Size)

	objInfo = NewObjectInfo(
		WithBucketName(bucketName),
		WithObjectName(objectName),
		WithFilePath(downloadPath))
	err = DownloadFile(*objInfo)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("download file %s successfully with %d bytes", objInfo.FilePath, fi.Size)
}
