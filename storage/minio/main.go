package main

import (
	"flag"
	"iflytek.com/weipan4/learn-go/storage"
	"iflytek.com/weipan4/learn-go/storage/minio/config"
	"log"
)

var configPath = flag.String("configFile", "storage/minio/config/config.toml", "config file of minio")

const (
	pkgOs        = "linux"
	pkgArch      = "amd64"
	pkgName      = "discover"
	filePath     = "/Users/a123/Downloads/working/2025.6/cmdb/discover"
	downloadPath = "/Users/a123/Downloads/working/2025.6/minio"
)

func main() {
	flag.Parse()
	err := config.InitConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	// 初始化storage
	initializer := storage.GetStorageInitFunc(StorageType)
	s := initializer()

	// 上传文件
	ulps := storage.NewPluginPkgUploadProps(
		storage.WithUploadOs(pkgOs),
		storage.WithUploadArch(pkgArch),
		storage.WithUploadFilepath(filePath))
	err = s.UploadFile(ulps)
	if err != nil {
		log.Fatalf("upload fail: %v", err)
	}

	// 下载文件
	dwlps := storage.NewPluginPkgDownloadProps(
		storage.WithDownloadOs(pkgOs),
		storage.WithDownloadArch(pkgArch),
		storage.WithDownloadPkgName(pkgName),
		storage.WithDestPath(downloadPath),
	)
	err = s.DownloadFile(dwlps)
	if err != nil {
		log.Fatalf("download fail: %v", err)
	}
}
