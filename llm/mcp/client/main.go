package main

import (
	"log"

	"github.com/joho/godotenv"
)

const (
	mcpURL           = "http://localhost:8888/mcp"
	mcpClientName    = "go-agent"
	mcpClientVersion = "1.0"
)

const (
	userMessage = "请查询/Users/panwei/Downloads/working/2025.10目录下的文件并展示"
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

	// 发起对话
	contents, err := llmCli.Chat(userMessage, sysMessage)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(contents)
}
