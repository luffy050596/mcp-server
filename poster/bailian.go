package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pkg/errors"
)

const (
	createURL   = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis"
	queryURL    = "https://dashscope.aliyuncs.com/api/v1/tasks/"
	asyncHeader = "X-DashScope-Async"
	authHeader  = "Authorization"
	contentType = "application/json"
	posterModel = "wanx-poster-generation-v1"
)

const (
	TaskStatusPending   = "PENDING"
	TaskStatusRunning   = "RUNNING"
	TaskStatusSuspended = "SUSPENDED"
	TaskStatusSucceeded = "SUCCEEDED"
	TaskStatusFailed    = "FAILED"
)

var LoraNameMap = map[string]struct{}{
	"2D 插画 1": {},
	"2D 插画 2": {},
	"浩瀚星云":    {},
	"浓郁色彩":    {},
	"光线粒子":    {},
	"透明玻璃":    {},
	"剪纸工艺":    {},
	"折纸工艺":    {},
	"中国水墨":    {},
	"中国刺绣":    {},
	"真实场景":    {},
	"2D 卡通":   {},
	"儿童水彩":    {},
	"赛博背景":    {},
	"浅蓝抽象":    {},
	"深蓝抽象":    {},
	"抽象点线":    {},
	"童话油画":    {},
}

var whRadios = map[string]struct{}{
	"横版": {},
	"竖版": {},
}

type PosterClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewPosterClient(key string) *PosterClient {
	return &PosterClient{
		apiKey:     key,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

type CreateTaskRequest struct {
	Model      string         `json:"model"`
	Input      PosterInput    `json:"input"`
	Parameters map[string]any `json:"parameters"`
}

type PosterInput struct {
	Title               string  `json:"title"`
	SubTitle            string  `json:"sub_title,omitempty"`
	BodyText            string  `json:"body_text,omitempty"`
	PromptTextZh        string  `json:"prompt_text_zh,omitempty"`
	PromptTextEn        string  `json:"prompt_text_en,omitempty"`
	WhRatios            string  `json:"wh_ratios,omitempty"`
	LoraName            string  `json:"lora_name,omitempty"`
	LoraWeight          float64 `json:"lora_weight,omitempty"`
	CtrlRatio           float64 `json:"ctrl_ratio,omitempty"`
	CtrlStep            float64 `json:"ctrl_step,omitempty"`
	GenerateMode        string  `json:"generate_mode,omitempty"`
	GenerateNum         int     `json:"generate_num,omitempty"`
	AuxiliaryParams     string  `json:"auxiliary_parameters,omitempty"`
	CreativeTitleLayout bool    `json:"creative_title_layout,omitempty"`
}

func (input *PosterInput) checkAndCompletion() error {
	if input.Title == "" {
		return errors.Errorf("title is required")
	}
	if input.GenerateMode != "" && input.GenerateMode != "generate" {
		if input.AuxiliaryParams == "" {
			return errors.Errorf("auxiliary_params is required if generate_mode is sr or hrf")
		}
	}

	promptTextEn := strings.Split(input.PromptTextEn, " ")
	zhCount := utf8.RuneCountInString(input.PromptTextZh)
	if len(promptTextEn)+zhCount > 50 {
		return errors.Errorf("the sum of the number of words in prompt_text_en and the number of characters in prompt_text_zh needs to be less than 50")
	}

	if _, ok := LoraNameMap[input.LoraName]; !ok {
		return errors.Errorf("lora_name must be one of the following: %v", LoraNameMap)
	}

	if _, ok := whRadios[input.WhRatios]; !ok {
		return errors.Errorf("wh_ratios must be one of the following: %v", whRadios)
	}

	return nil
}

type CreateTaskResponse struct {
	RequestID string `json:"request_id"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Output    struct {
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"`
	} `json:"output"`
}

func (c *PosterClient) CreatePosterTask(ctx context.Context, req CreateTaskRequest) (*CreateTaskResponse, error) {
	if err := req.Input.checkAndCompletion(); err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", createURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set(asyncHeader, "enable")
	httpReq.Header.Set(authHeader, "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", contentType)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, response: %s", resp.Status, body)
	}

	var result CreateTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &result, nil
}

type QueryTaskResponse struct {
	RequestID string `json:"request_id"`
	Output    struct {
		TaskID          string   `json:"task_id"`
		TaskStatus      string   `json:"task_status"`
		RenderURLs      []string `json:"render_urls"`
		BgURLs          []string `json:"bg_urls"`
		AuxiliaryParams []string `json:"auxiliary_parameters"`
		SubmitTime      string   `json:"submit_time"`
		ScheduledTime   string   `json:"scheduled_time"`
		EndTime         string   `json:"end_time"`
		Code            string   `json:"code"`
		Message         string   `json:"message"`
	} `json:"output"`
	Usage struct {
		ImageCount int `json:"image_count"`
	} `json:"usage"`
}

func (c *PosterClient) QueryPosterTask(ctx context.Context, taskID string) (*QueryTaskResponse, error) {
	url := fmt.Sprintf("%s/%s", queryURL, taskID)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set(authHeader, "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, response: %s", resp.Status, body)
	}

	var result QueryTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &result, nil
}
