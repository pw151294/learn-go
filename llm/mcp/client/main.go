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
	userMessage = "请帮我查询合肥的天气预报，合肥的经度是116.3974673500868，纬度是39.90873966065374"
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
	toolCalls, err := llmCli.listToolCalls()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(toolCalls)

	// 查询合肥的天气
	contents, err := llmCli.Chat(toolCalls, userMessage, sysMessage)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(contents)
}
