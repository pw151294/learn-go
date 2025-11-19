package utils

import (
	"sync"

	"github.com/yanyiwu/gojieba"
)

var (
	jb   *gojieba.Jieba
	once sync.Once
)

// InitSegmenter 使用 gojieba 内置词典初始化（无需提供自定义词典）
func InitSegmenter() {
	once.Do(func() {
		jb = gojieba.NewJieba()
	})
}

// SplitText 分词（精确模式 + HMM）
func SplitText(text string) []string {
	if jb == nil {
		InitSegmenter()
	}
	return jb.Cut(text, true)
}

// CloseSegmenter 释放资源（可在程序退出时调用）
func CloseSegmenter() {
	if jb != nil {
		jb.Free()
	}
}
