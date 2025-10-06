package main

type ToolSchema struct {
	Type       string                  `json:"type"`
	Properties map[string]ToolVariable `json:"properties"`
	Required   []string                `json:"required"`
}

type ToolVariable struct {
	Name        string
	Type        string
	Description string
	Required    bool
}

type ToolCall struct {
	ToolName      string
	Description   string
	ToolVariables []ToolVariable
}
