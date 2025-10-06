package main

import "github.com/mark3labs/mcp-go/mcp"

const (
	JingDu = "longitude"
	WeiDu  = "latitude"
)

// 添加 forecast 工具
func weatherForcastTool() mcp.Tool {
	wfTool := mcp.NewTool("weather_forecast",
		mcp.WithDescription("获取某地的天气预报"),
		mcp.WithNumber(WeiDu, mcp.Required(), mcp.Description("地点的纬度")),
		mcp.WithNumber(JingDu, mcp.Required(), mcp.Description("地点的经度")))
	return wfTool
}
