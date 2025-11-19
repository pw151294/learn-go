package test

import (
	"os"
	"path/filepath"
	"testing"

	"iflytek.com/weipan4/learn-go/rag/configs"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/weaviate"
	"iflytek.com/weipan4/learn-go/rag/pkg/loggers"
)

func TestCreateClient(t *testing.T) {
	// 假设配置文件在项目根目录下的 configs/config.toml
	configPath := filepath.Join("..", "configs", "config.toml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("配置文件不存在: %v", err)
	}

	if err := configs.InitConfig(configPath); err != nil {
		t.Fatalf("配置初始化失败: %v", err)
	}

	// 初始化日志记录器，这是 InitWeaviateClient 的一个依赖项
	loggerConfig := configs.GetConfig().Log
	if err := loggers.InitGlobalLogger(&loggerConfig); err != nil {
		t.Fatalf("日志系统初始化失败：%v", err)
	}

	err := weaviate.InitWeaviateClient()
	if err != nil {
		t.Fatalf("Weaviate 客户端初始化失败: %v", err)
	}

	client := weaviate.GetWeaviateClient()
	if client == nil {
		t.Fatal("Weaviate 客户端未正确创建")
	}
	t.Log("Weaviate 客户端创建成功")
}
