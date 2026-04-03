package cdn

import (
	"net/url"
	"strings"

	"iflytek.com/weipan4/learn-go/zoom/cloud_recording/configs"
)

type URLRewriter struct {
	cfg *configs.CDNConfig
}

func NewURLRewriter(cfg *configs.CDNConfig) *URLRewriter {
	return &URLRewriter{cfg: cfg}
}

// RewriteURL 将 MinIO URL 替换为 CDN URL
// 如果 CDN 未启用，直接返回原 URL
// 如果启用，提取路径部分，拼接 CDN 域名
func (r *URLRewriter) RewriteURL(minioURL string) string {
	if !r.cfg.Enabled {
		return minioURL
	}

	parsed, err := url.Parse(minioURL)
	if err != nil {
		return minioURL
	}

	cdnBase := strings.TrimRight(r.cfg.BaseURL, "/")
	return cdnBase + parsed.Path
}
