package main

import (
	"encoding/json"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/pkg/errors"
)

func currentTimeTool() *protocol.Tool {
	t := &protocol.Tool{
		Name:        "current time",
		Description: "Get the current time with timezone, Asia/Singapore is default",
		InputSchema: protocol.InputSchema{
			Type: protocol.Object,
			Properties: map[string]any{
				"timezone": map[string]string{
					"type":        "string",
					"description": "The timezone to get the current time",
				},
			},
		},
	}
	return t
}

func currentTimeHandler(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	timezone, ok := request.Arguments["timezone"].(string)
	if !ok {
		return nil, errors.Errorf("timezone is required")
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load location: %s", timezone)
	}

	now := time.Now().In(loc)
	reply := newCurrentTimeReply(now)
	json, err := reply.Json()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal current time reply")
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: string(json),
			},
		},
	}, nil
}

type currentTimeReply struct {
	Time    string `json:"time"`
	Weekday string `json:"weekday"`
}

func newCurrentTimeReply(now time.Time) *currentTimeReply {
	return &currentTimeReply{
		Time:    now.Format("2006-01-02 15:04:05"),
		Weekday: now.Weekday().String(),
	}
}

func (c *currentTimeReply) Json() ([]byte, error) {
	return json.Marshal(c)
}
