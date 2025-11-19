package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	mcpURL           = "http://localhost:8000/mcp"
	mcpClientName    = "go-agent"
	mcpClientVersion = "1.0"
)

const (
	userMessage = "请连接服务器172.30.34.73，查询/data/weipan4目录下的文件，获取结果并输出"
	sysMessage  = ""
)

func main() {
	// 初始化大模型工具类
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
	err := InitLLMClient()
	if err != nil {
		log.Fatal(err)
	}

	// 查询出所有的工具信息
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancelFunc()
	listReq := mcp.ListToolsRequest{}
	listRes, err := llmCli.mcpCli.ListTools(ctx, listReq)
	if err != nil {
		log.Fatal(err)
	}
	toolsBytes, err := json.Marshal(listRes)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("tools: %v", string(toolsBytes))

	// 发起对话
	contents, err := llmCli.Chat(userMessage, sysMessage)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(contents)
}
