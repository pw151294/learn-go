package main

import (
	"fmt"
	"github.com/mark3labs/mcp-go/server"
)

const (
	mcpServerName    = "Weather Demo"
	mcpServerVersion = "1.0.0"
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
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Error serving stdio: %s\n", err)
	}
}
