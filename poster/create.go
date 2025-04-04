package main

import (
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func createTool() *protocol.Tool {
	t := &protocol.Tool{
		Name:        "create_poster",
		Description: "Create poster images",
		InputSchema: protocol.InputSchema{
			Type: protocol.Object,
			Properties: map[string]any{
				"generate_num": map[string]string{
					"type":        "number",
					"description": "The generate num",
				},
			},
		},
	}
	for k, v := range defaultInputArgs {
		t.InputSchema.Properties[k] = v
	}
	return t
}

func createHandler(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	return posterHandler(request, "generate")
}
