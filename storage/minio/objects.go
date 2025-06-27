package main

type ObjectInfo struct {
	BucketName  string
	ObjectName  string
	FilePath    string
	ContentType string
}

type ObjectInfoOptions func(*ObjectInfo)

func NewObjectInfo(opts ...ObjectInfoOptions) *ObjectInfo {
	info := &ObjectInfo{
		ContentType: "application/octet-stream",
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(info)
		}
	}

	return info
}

func WithBucketName(bucketName string) ObjectInfoOptions {
	return func(info *ObjectInfo) {
		info.BucketName = bucketName
	}
}

func WithObjectName(objectName string) ObjectInfoOptions {
	return func(info *ObjectInfo) {
		info.ObjectName = objectName
	}
}

func WithFilePath(filePath string) ObjectInfoOptions {
	return func(info *ObjectInfo) {
		info.FilePath = filePath
	}
}

func WithContentType(contentType string) ObjectInfoOptions {
	return func(info *ObjectInfo) {
		info.ContentType = contentType
	}
}
