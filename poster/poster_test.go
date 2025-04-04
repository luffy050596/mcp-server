package main

import (
	"flag"
	"fmt"
	"testing"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func TestPoster(t *testing.T) {
	flag.Parse()

	ret, err := createHandler(&protocol.CallToolRequest{
		Name: "create_poster",
		Arguments: map[string]any{
			"title":          "春节快乐",
			"sub_title":      "家庭团聚，共享天伦之乐",
			"body_text":      "春节是中国最重要的传统节日之一，它象征着新的开始和希望",
			"prompt_text_zh": "灯笼，小猫，梅花",
			"wh_ratios":      "竖版",
			"lora_name":      "童话油画",
			"lora_weight":    0.8,
			"ctrl_ratio":     0.7,
			"ctrl_step":      0.7,
			"generate_num":   1,
			"generate_mode":  "generate",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(ret)
}
