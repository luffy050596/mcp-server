package main

import (
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

const text = `
{
  "headers": {
    "Authorization": {
      "type": "string",
      "required": true,
      "description": "推荐使用百炼 API-Key，也可填 DashScope API-Key。例如：Bearer d1xxx2a"
    },
    "X-DashScope-Async": {
      "type": "string", 
      "required": true,
      "description": "是否使用 DashScope 异步调用。HTTP 只支持异步调用",
      "value": "enable"
    },
    "Content-Type": {
      "type": "string",
      "required": true,
      "description": "请求内容类型",
      "value": "application/json"
    }
  },
  "body": {
    "model": {
      "type": "string",
      "required": true,
      "description": "调用模型",
      "value": "wanx-poster-generation-v1"
    },
    "parameters": {
      "type": "object",
      "required": true,
      "description": "其他模型调用参数，只需要输入一个空对象即可",
      "value": "{}"
    },
    "input": {
      "type": "object",
      "required": true,
      "description": "输入图像的基本信息",
      "properties": {
        "generate_mode": {
          "type": "string",
          "required": true,
          "description": "海报生成模式",
          "enum": ["generate", "sr", "hrf"],
          "default": "generate",
          "notes": [
            "generate：默认模式，生成新海报时使用",
            "sr：高分辨率模式，用于提升海报分辨率",
            "hrf：高清修复模式，用于修复海报模糊问题",
            "sr 和 hrf 模式下，需要输入辅助参数 auxiliary_parameters",
            "只能从["generate","sr","hrf"]中选择。海报生成的基础模式为"generate"，选择此模式会返回海报图片的url（render_urls）和与其一一对应的辅助参数（auxiliary_parameters）。用户可从返回的结果中，选择需要进行分辨率提升（或者高清修复）的海报，通过二次调用，输入选中的海报对应的辅助参数，将generate_mode设置为"sr"（或者"hrf"），得到对应的高分辨率（高清修复）结果。
          ]
        },
        "generate_num": {
          "type": "integer",
          "required": false,
          "description": "生成的海报数",
          "range": [1, 4],
          "default": 1,
          "notes": "仅在 generate_mode=generate 时有效"
        },
        "auxiliary_parameters": {
          "type": "string",
          "required": false,
          "description": "需要提升分辨率或者高清修复的海报图片对应的辅助参数",
          "notes": "当 generate_mode 为 sr 或 hrf 时为必选项"
        },
        "title": {
          "type": "string",
          "required": true,
          "description": "主标题",
          "maxLength": 30,
          "example": "春节快乐"
        },
        "sub_title": {
          "type": "string",
          "required": false,
          "description": "副标题",
          "maxLength": 30,
          "example": "家庭团聚，共享天伦之乐"
        },
        "body_text": {
          "type": "string",
          "required": false,
          "description": "正文",
          "maxLength": 50,
          "example": "春节是中国最重要的传统节日之一，它象征着新的开始和希望"
        },
        "prompt_text_zh": {
          "type": "string",
          "required": false,
          "description": "中文提示词",
          "example": "小朋友画的可爱的龙，白色背景",
          "notes": "中文和英文提示词至少二选一设置，两个字段字符数加起来最多 50 个字/单词"
        },
        "prompt_text_en": {
          "type": "string",
          "required": false,
          "description": "英文提示词",
          "example": "Children draw a lovely dragon, white background",
          "notes": "中文和英文提示词至少二选一设置，两个字段字符数加起来最多 50 个字/单词"
        },
        "wh_ratios": {
          "type": "string",
          "required": false,
          "description": "生成海报的版式",
          "enum": ["横版", "竖版"],
          "default": "横版"
        },
        "lora_name": {
          "type": "string",
          "required": false,
          "description": "海报风格名称",
          "default": "",
          "enum": [
            "2D 插画 1",
            "2D 插画 2",
            "浩瀚星云",
            "浓郁色彩",
            "光线粒子",
            "透明玻璃",
            "剪纸工艺",
            "折纸工艺",
            "中国水墨",
            "中国刺绣",
            "真实场景",
            "2D 卡通",
            "儿童水彩",
            "赛博背景",
            "浅蓝抽象",
            "深蓝抽象",
            "抽象点线",
            "童话油画"
          ]
        },
        "lora_weight": {
          "type": "float",
          "required": false,
          "description": "海报风格权重，需要与 lora_name 参数配合使用",
          "range": [0, 1],
          "default": 0.8,
          "notes": "取值越接近 1，海报风格越明显"
        },
        "ctrl_ratio": {
          "type": "float",
          "required": false,
          "description": "留白效果权重，用于控制海报留白效果",
          "range": [0, 1],
          "default": 0.7,
          "notes": "取值越接近 1，留白效果越好，但海报背景生成效果可能会受到负面影响"
        },
        "ctrl_step": {
          "type": "float",
          "required": false,
          "description": "留白步数比例，用于控制海报留白效果",
          "range": "(0, 1]",
          "default": 0.7,
          "notes": "取值越接近 1，留白效果越好，但是海报背景生成效果可能会受到负面影响"
        },
        "creative_title_layout": {
          "type": "boolean",
          "required": false,
          "description": "标题是否启用创意排版",
          "default": false
        }
      }
    }
  }
}
`

func posterHelpTool() *protocol.Tool {
	t := &protocol.Tool{
		Name:        "input_help",
		Description: "Query the input help json document of the poster",
		InputSchema: protocol.InputSchema{
			Type:       protocol.Object,
			Properties: map[string]any{},
		},
	}
	return t
}

func posterHelpHandler(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: text,
			},
		},
	}, nil
}
