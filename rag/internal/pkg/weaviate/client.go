package weaviate

import (
	"context"
	"fmt"
	"sync"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
	"iflytek.com/weipan4/learn-go/rag/configs"
	"iflytek.com/weipan4/learn-go/rag/pkg/loggers"
)

var (
	client     *weaviate.Client
	clientOnce sync.Once
)

// InitWeaviateClient 从配置初始化 Weaviate 客户端
func InitWeaviateClient() error {
	clientOnce.Do(func() {
		loggers.Info("开始初始化 Weaviate 客户端")
		conf := configs.GetConfig()
		if conf == nil {
			loggers.Error("Weaviate 初始化失败: 配置未初始化")
			return
		}

		weaviateConf := conf.Weaviate
		cfg := weaviate.Config{
			Host:       weaviateConf.Host,
			Scheme:     weaviateConf.Scheme,
			AuthConfig: auth.ApiKey{Value: weaviateConf.ApiKey},
		}

		c, e := weaviate.NewClient(cfg)
		if e != nil {
			loggers.Error("Weaviate 客户端创建失败", loggers.ErrorField(e))
			return
		}

		loggers.Info("Weaviate 客户端创建成功，开始检查服务状态")
		isReady, e := c.Misc().ReadyChecker().Do(context.Background())
		if e != nil {
			loggers.Error("Weaviate 状态检查失败", loggers.ErrorField(e))
			return
		}
		if !isReady {
			loggers.Error("Weaviate 服务未就绪")
			return
		}

		client = c
		loggers.Info("Weaviate 客户端初始化完成")
	})

	if client == nil {
		return fmt.Errorf("weaviate client initialization failed, see logs for details")
	}
	return nil
}

// GetWeaviateClient 返回已初始化的 Weaviate 客户端实例
func GetWeaviateClient() *weaviate.Client {
	if client == nil {
		return nil
	}
	return client
}
