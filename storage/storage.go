package storage

type PluginPkgStorage interface {
	UploadFile(props *PluginPkgUploadProps) error
	DownloadFile(props *PluginPkgDownloadProps) error
}

type PluginPkgStorageInitFunc func() PluginPkgStorage

var initFuncCache = make(map[string]PluginPkgStorageInitFunc)

func RegisterInitFunc(storageType string, initFunc PluginPkgStorageInitFunc) {
	initFuncCache[storageType] = initFunc
}

func GetStorageInitFunc(storageType string) PluginPkgStorageInitFunc {
	return initFuncCache[storageType]
}
