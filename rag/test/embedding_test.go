package test

import (
	"os"
	"path/filepath"
	"testing"

	"iflytek.com/weipan4/learn-go/rag/configs"
	"iflytek.com/weipan4/learn-go/rag/internal/pkg/embedding"
)

func TestCreateEmbeddingClient(t *testing.T) {
	// 假设配置文件在项目根目录下的 configs/config.toml
	configPath := filepath.Join("..", "configs", "config.toml")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("配置文件不存在: %v", err)
	}

	if err := configs.InitConfig(configPath); err != nil {
		t.Fatalf("配置初始化失败: %v", err)
	}

	client, err := embedding.NewOpenAIEmbeddingClient()
	if err != nil {
		t.Fatalf("Embedding 客户端初始化失败: %v", err)
	}

	if client == nil {
		t.Fatal("Embedding 客户端未正确创建")
	}
	t.Log("Embedding 客户端创建成功")
}
