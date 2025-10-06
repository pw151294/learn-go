package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
)

var llmCli *LLMClient

// LLMClient 封装大模型工具类
type LLMClient struct {
	mcpCli    client.MCPClient
	openaiCli *openai.Client
}

func GetLLMClient() *LLMClient {
	return llmCli
}

func (l *LLMClient) callTool(name string, arguments map[string]interface{}) ([]mcp.Content, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancelFunc()

	request := mcp.CallToolRequest{}
	request.Params.Name = name
	request.Params.Arguments = arguments
	result, err := l.mcpCli.CallTool(ctx, request)
	if err != nil {
		return nil, err
	}

	contentBytes, err := json.Marshal(result.Content)
	if err != nil {
		return nil, err
	}
	log.Printf("content of call tool results: %v", string(contentBytes))
	return result.Content, nil
}

func (l *LLMClient) listToolCalls() ([]ToolCall, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancelFunc()

	listReq := mcp.ListToolsRequest{}
	listRes, err := l.mcpCli.ListTools(ctx, listReq)
	if err != nil {
		return nil, err
	}

	return convertToolCalls(listRes.Tools), nil
}

func (l *LLMClient) listToolNames() ([]string, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancelFunc()

	request := mcp.ListToolsRequest{}
	result, err := l.mcpCli.ListTools(ctx, request)
	if err != nil {
		return nil, err
	}
	if bytes, err := json.Marshal(result.Tools); err == nil {
		log.Println(string(bytes))
	}

	toolNames := make([]string, 0, len(result.Tools))
	for _, tool := range result.Tools {
		toolNames = append(toolNames, tool.Name)
	}
	return toolNames, nil
}

func (l *LLMClient) Chat(toolCalls []ToolCall, userMessage, systemMessage string) ([]string, error) {
	// 构建工具调用的参数
	tools := make([]openai.ChatCompletionToolParam, 0, len(toolCalls))
	for _, call := range toolCalls {
		properties := make(map[string]ToolVariable)
		requires := make([]string, 0, len(call.ToolVariables))
		for _, v := range call.ToolVariables {
			properties[v.Name] = v
			if v.Required {
				requires = append(requires, v.Name)
			}
		}
		schema := ToolSchema{
			Type:       "object",
			Properties: properties,
			Required:   requires,
		}
		params := make(map[string]interface{})
		paramsBytes, err := json.Marshal(schema)
		if err != nil {
			continue
		}
		if err = json.Unmarshal(paramsBytes, &params); err != nil {
			continue
		}
		tool := openai.ChatCompletionToolParam{
			Function: openai.FunctionDefinitionParam{
				Name:        call.ToolName,
				Strict:      param.NewOpt(true),
				Description: param.NewOpt(call.Description),
				Parameters:  shared.FunctionParameters(params),
			},
		}
		tools = append(tools, tool)
	}

	// 构建大模型对话的message
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(userMessage),
		openai.SystemMessage(systemMessage),
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancelFunc()

	// 对话
	res, err := l.openaiCli.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    getModelName(),
		Tools:    tools,
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}
	if len(res.Choices) == 0 {
		return nil, nil
	}
	msg := res.Choices[0].Message
	messages = append(messages, msg.ToParam())

	texts := make([]string, 0)
	if len(msg.ToolCalls) > 0 {
		for _, call := range msg.ToolCalls {
			toolName := call.Function.Name
			arguments := make(map[string]interface{})
			if err = json.Unmarshal([]byte(call.Function.Arguments), &arguments); err != nil {
				continue
			}
			contents, err := l.callTool(toolName, arguments)
			if err != nil {
				continue
			}
			for _, content := range contents {
				if text, ok := content.(mcp.TextContent); ok {
					texts = append(texts, text.Text)
				}
			}
		}
	}
	return texts, nil
}

func convertToolCalls(tools []mcp.Tool) []ToolCall {
	toolCalls := make([]ToolCall, 0, len(tools))
	for _, tool := range tools {
		toolCalls = append(toolCalls, convertToolCall(tool))
	}
	return toolCalls
}

func convertToolCall(tool mcp.Tool) ToolCall {
	toolCall := ToolCall{}
	toolCall.ToolName = tool.Name
	toolCall.Description = tool.Description
	schema := tool.InputSchema
	requiredVars := make(map[string]struct{})
	for _, vName := range schema.Required {
		requiredVars[vName] = struct{}{}
	}

	toolVars := make([]ToolVariable, 0, len(schema.Properties))
	for vName, vProp := range schema.Properties {
		varBytes, err := json.Marshal(vProp)
		if err != nil {
			continue
		}
		toolVar := ToolVariable{}
		if err := json.Unmarshal(varBytes, &toolVar); err != nil {
			continue
		}
		toolVar.Name = vName
		_, ok := requiredVars[toolVar.Name]
		toolVar.Required = ok
		toolVars = append(toolVars, toolVar)
	}

	toolCall.ToolVariables = toolVars
	return toolCall
}

func InitLLMClient() error {
	// 初始化mcp客户端
	mcpCli, err := client.NewStreamableHttpClient(mcpURL)
	if err != nil {
		return err
	}
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancelFunc()
	if err := mcpCli.Start(ctx); err != nil {
		return err
	}
	initReq := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    mcpClientName,
				Version: mcpClientVersion,
			},
		},
	}
	if _, err := mcpCli.Initialize(ctx, initReq); err != nil {
		return err
	}

	// 初始化openai客户端
	openaiCli := openai.NewClient(
		option.WithAPIKey(getAPIKey()),
		option.WithBaseURL(getBaseURL()))

	// 初始化大模型的工具类
	llmCli = &LLMClient{
		mcpCli:    mcpCli,
		openaiCli: &openaiCli,
	}
	return nil
}
