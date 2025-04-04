package main

import (
	"errors"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func refineTool() *protocol.Tool {
	t := &protocol.Tool{
		Name:        "refine_poster",
		Description: "Refine poster images",
		InputSchema: protocol.InputSchema{
			Type: protocol.Object,
			Properties: map[string]any{
				"refine_mode": map[string]string{
					"type":        "string",
					"description": "The refine mode, sr or hrf, default is sr",
					"enum":        "sr or hrf",
				},
				"auxiliary_params": map[string]string{
					"type":        "string",
					"description": "The auxiliary params, 海报生成的基础模式为generate，选择此模式会返回海报图片的url（render_urls）和与其一一对应的辅助参数（auxiliary_parameters）。用户可从返回的结果中，选择需要进行分辨率提升（或者高清修复）的海报，通过二次调用，输入选中的海报对应的辅助参数",
				},
			},
		},
	}
	for k, v := range defaultInputArgs {
		t.InputSchema.Properties[k] = v
	}
	return t
}

func refineHandler(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	refineMode := request.Arguments["refine_mode"].(string)
	if refineMode != "sr" && refineMode != "hrf" {
		return nil, errors.New("refine_mode must be sr or hrf")
	}
	return posterHandler(request, refineMode)
}
