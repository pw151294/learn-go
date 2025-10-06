package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
)

const (
	mcpServerName    = "Weather Demo"
	mcpServerVersion = "1.0.0"
	mcpServerAddr    = "localhost:8888"
)

func main() {
	s := server.NewMCPServer(
		mcpServerName,
		mcpServerVersion,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	s.AddTool(weatherForcastTool(), callForecast)
	httpServer := server.NewStreamableHTTPServer(s)
	if err := httpServer.Start(mcpServerAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
