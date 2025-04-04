package main

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/pkg/errors"
)

func timeTool() *protocol.Tool {
	t := &protocol.Tool{
		Name:        "time from timestamp",
		Description: "Get the time from a given timestamp",
		InputSchema: protocol.InputSchema{
			Type: protocol.Object,
			Properties: map[string]any{
				"timestamp": map[string]string{
					"type":        "string",
					"description": "The timestamp to get the time",
				},
				"timezone": map[string]string{
					"type":        "string",
					"description": "The timezone to get the time",
				},
			},
		},
	}
	return t
}

func timeHandler(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	timestamp, ok := request.Arguments["timestamp"].(string)
	if !ok {
		return nil, errors.Errorf("timestamp is required")
	}
	timezone, ok := request.Arguments["timezone"].(string)
	if !ok {
		return nil, errors.Errorf("timezone is required")
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load location: %s", timezone)
	}

	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse timestamp: %s", timestamp)
	}

	now := time.Unix(timestampInt, 0).In(loc)
	reply := newTimeReply(now)
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

type timeReply struct {
	Time    string `json:"time"`
	Weekday string `json:"weekday"`
}

func newTimeReply(now time.Time) *timeReply {
	return &timeReply{
		Time:    now.Format("2006-01-02 15:04:05"),
		Weekday: now.Weekday().String(),
	}
}

func (t *timeReply) Json() ([]byte, error) {
	return json.Marshal(t)
}
