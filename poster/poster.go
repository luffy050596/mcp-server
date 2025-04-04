package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var defaultInputArgs = map[string]any{
	"title": map[string]string{
		"type":        "string",
		"description": "The title",
	},
	"sub_title": map[string]string{
		"type":        "string",
		"description": "The sub title",
	},
	"body_text": map[string]string{
		"type":        "string",
		"description": "The body text",
	},
	"prompt_text_zh": map[string]string{
		"type":        "string	",
		"description": "The prompt text in chinese",
	},
	"prompt_text_en": map[string]string{
		"type":        "string",
		"description": "The prompt text in english",
	},
	"wh_ratios": map[string]string{
		"type":        "string",
		"description": fmt.Sprintf("The wh ratios, optional, must be one of the following: %v", whRadios),
	},
	"lora_name": map[string]string{
		"type":        "string",
		"description": fmt.Sprintf("The lora name, optional, must be one of the following: %v", LoraNameMap),
	},
	"ctrl_step": map[string]string{
		"type":        "number",
		"description": "The control step",
	},
}

func posterHandler(request *protocol.CallToolRequest, generateMode string) (*protocol.CallToolResult, error) {
	var err error

	input := PosterInput{
		GenerateMode: generateMode,
	}

	titleArg, ok := request.Arguments["title"]
	if !ok {
		return nil, errors.Errorf("title is required")
	}
	input.Title = titleArg.(string)

	subTitleArg, ok := request.Arguments["sub_title"]
	if ok {
		input.SubTitle = subTitleArg.(string)
	}
	bodyTextArg, ok := request.Arguments["body_text"]
	if ok {
		input.BodyText = bodyTextArg.(string)
	}
	promptTextZhArg, ok := request.Arguments["prompt_text_zh"]
	if ok {
		input.PromptTextZh = promptTextZhArg.(string)
	}
	promptTextEnArg, ok := request.Arguments["prompt_text_en"]
	if ok {
		input.PromptTextEn = promptTextEnArg.(string)
	}
	whRatiosArg, ok := request.Arguments["wh_ratios"]
	if ok {
		input.WhRatios = whRatiosArg.(string)
	}
	loraNameArg, ok := request.Arguments["lora_name"]
	if ok {
		input.LoraName = loraNameArg.(string)
	}
	auxiliaryParamsArg, ok := request.Arguments["auxiliary_params"]
	if ok {
		input.AuxiliaryParams = auxiliaryParamsArg.(string)
	}

	if loraWeightArg, ok := request.Arguments["lora_weight"].(string); ok {
		input.LoraWeight, err = strconv.ParseFloat(loraWeightArg, 64)
		if err != nil {
			return nil, errors.Errorf("lora_weight is required")
		}
	}

	if ctrlRatioArg, ok := request.Arguments["ctrl_ratio"].(string); ok {
		input.CtrlRatio, err = strconv.ParseFloat(ctrlRatioArg, 64)
		if err != nil {
			return nil, errors.Errorf("ctrl_ratio is required")
		}
	}

	if ctrlStepArg, ok := request.Arguments["ctrl_step"].(string); ok {
		input.CtrlStep, err = strconv.ParseFloat(ctrlStepArg, 64)
		if err != nil {
			return nil, errors.Errorf("ctrl_step is required")
		}
	}

	if generateNumArg, ok := request.Arguments["generate_num"].(string); ok {
		input.GenerateNum, err = strconv.Atoi(generateNumArg)
		if err != nil {
			return nil, errors.Errorf("generate_num is required")
		}
	}

	if creativeTitleLayoutArg, ok := request.Arguments["creative_title_layout"].(string); ok {
		input.CreativeTitleLayout, err = strconv.ParseBool(creativeTitleLayoutArg)
		if err != nil {
			return nil, errors.Errorf("creative_title_layout is required")
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := NewPosterClient(key)

	createResp, err := client.CreatePosterTask(ctx, &CreateTaskRequest{
		Model:      posterModel,
		Input:      input,
		Parameters: map[string]any{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create poster task")
	}
	if createResp.Code != "" {
		return nil, errors.Errorf("code: %s, message: %s", createResp.Code, createResp.Message)
	}

	fmt.Println("createResp", createResp)

	taskID := createResp.Output.TaskID
	if taskID == "" {
		return nil, errors.Errorf("failed to create poster task")
	}

	reply := &posterReply{}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	eg := &errgroup.Group{}
	eg.Go(func() error {
		for {
			select {
			case <-ticker.C:
				ctx1, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				queryResp, err1 := client.QueryPosterTask(ctx1, taskID)
				cancel()

				if err1 != nil {
					return errors.Wrapf(err1, "failed to query poster task")
				}

				switch queryResp.Output.TaskStatus {
				case TaskStatusSucceeded:
					reply.RenderURLs = queryResp.Output.RenderURLs
					reply.AuxiliaryParams = queryResp.Output.AuxiliaryParams
					return nil
				case TaskStatusFailed:
					return errors.Errorf("code: %s, message: %s", queryResp.Output.Code, queryResp.Output.Message)
				case TaskStatusRunning, TaskStatusPending, TaskStatusSuspended:
					fmt.Printf("poster task is %s: taskID: %s\n", queryResp.Output.TaskStatus, taskID)
					continue
				default:
					return errors.Errorf("query task status failed: %s", queryResp.Output.TaskStatus)
				}

			case <-ctx.Done():
				return errors.Errorf("query task status timeout: taskID: %s", taskID)
			}
		}
	})

	if err2 := eg.Wait(); err2 != nil {
		return nil, err2
	}

	replyJSON, err := reply.Json()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal poster reply")
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: string(replyJSON),
			},
		},
	}, nil
}

type posterReply struct {
	RenderURLs      []string `json:"render_urls"`
	AuxiliaryParams []string `json:"auxiliary_parameters"`
	// BgURLs     []string `json:"bg_urls"`
	// ImageCount    int      `json:"image_count"`
	// SubmitTime    string   `json:"submit_time"`
	// ScheduledTime string   `json:"scheduled_time"`
	// EndTime       string   `json:"end_time"`
}

func (r *posterReply) Json() ([]byte, error) {
	return json.Marshal(r)
}
