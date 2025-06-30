package storage

type PluginPkgUploadProps struct {
	Os       string
	Arch     string
	Filepath string
}

type PluginPkgUploadPropsOptions func(props *PluginPkgUploadProps)

func NewPluginPkgUploadProps(opts ...PluginPkgUploadPropsOptions) *PluginPkgUploadProps {
	props := &PluginPkgUploadProps{
		Os:   "linux",
		Arch: "amd64",
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(props)
		}
	}

	return props
}

func WithUploadOs(os string) PluginPkgUploadPropsOptions {
	return func(props *PluginPkgUploadProps) {
		props.Os = os
	}
}

func WithUploadArch(arch string) PluginPkgUploadPropsOptions {
	return func(props *PluginPkgUploadProps) {
		props.Arch = arch
	}
}

func WithUploadFilepath(filepath string) PluginPkgUploadPropsOptions {
	return func(props *PluginPkgUploadProps) {
		props.Filepath = filepath
	}
}

type PluginPkgDownloadProps struct {
	Os       string
	Arch     string
	PkgName  string
	DestPath string
}

type PluginPkgDownloadPropsOptions func(props *PluginPkgDownloadProps)

func NewPluginPkgDownloadProps(opts ...PluginPkgDownloadPropsOptions) *PluginPkgDownloadProps {
	props := &PluginPkgDownloadProps{
		Os:   "linux",
		Arch: "amd64",
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(props)
		}
	}

	return props
}

func WithDownloadOs(os string) PluginPkgDownloadPropsOptions {
	return func(props *PluginPkgDownloadProps) {
		props.Os = os
	}
}

func WithDownloadArch(arch string) PluginPkgDownloadPropsOptions {
	return func(props *PluginPkgDownloadProps) {
		props.Arch = arch
	}
}

func WithDownloadPkgName(pkgName string) PluginPkgDownloadPropsOptions {
	return func(props *PluginPkgDownloadProps) {
		props.PkgName = pkgName
	}
}

func WithDestPath(destPath string) PluginPkgDownloadPropsOptions {
	return func(props *PluginPkgDownloadProps) {
		props.DestPath = destPath
	}
}
