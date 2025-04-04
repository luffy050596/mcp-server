package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/pkg/errors"
)

func timestampTool() *protocol.Tool {
	t := &protocol.Tool{
		Name:        "timestamp",
		Description: "Get the timestamp from a given time",
		InputSchema: protocol.InputSchema{
			Type: protocol.Object,
			Properties: map[string]any{
				"year": map[string]string{
					"type":        "string",
					"description": "The year of the time",
				},
				"month": map[string]string{
					"type":        "string",
					"description": "The month of the time",
				},
				"day": map[string]string{
					"type":        "string",
					"description": "The day of the time",
				},
				"hour": map[string]string{
					"type":        "string",
					"description": "The hour of the time",
				},
				"minute": map[string]string{
					"type":        "string",
					"description": "The minute of the time",
				},
				"second": map[string]string{
					"type":        "string",
					"description": "The second of the time",
				},
				"timezone": map[string]string{
					"type":        "string",
					"description": "The timezone of the time",
					"default":     "Asia/Singapore",
				},
			},
		},
	}
	return t
}

func timestampHandler(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	year, ok := request.Arguments["year"].(string)
	if !ok {
		return nil, errors.Errorf("year is required")
	}
	month, ok := request.Arguments["month"].(string)
	if !ok {
		return nil, errors.Errorf("month is required")
	}
	day, ok := request.Arguments["day"].(string)
	if !ok {
		return nil, errors.Errorf("day is required")
	}
	hour, ok := request.Arguments["hour"].(string)
	if !ok {
		return nil, errors.Errorf("hour is required")
	}
	minute, ok := request.Arguments["minute"].(string)
	if !ok {
		return nil, errors.Errorf("minute is required")
	}
	second, ok := request.Arguments["second"].(string)
	if !ok {
		return nil, errors.Errorf("second is required")
	}
	tz, ok := request.Arguments["timezone"].(string)
	if !ok {
		return nil, errors.Errorf("timezone is required")
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load location: %s", tz)
	}

	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert year to int: %s", year)
	}
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert month to int: %s", month)
	}
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert day to int: %s", day)
	}
	hourInt, err := strconv.Atoi(hour)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert hour to int: %s", hour)
	}
	minuteInt, err := strconv.Atoi(minute)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert minute to int: %s", minute)
	}
	secondInt, err := strconv.Atoi(second)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert second to int: %s", second)
	}
	tp := time.Date(yearInt, time.Month(monthInt), dayInt, hourInt, minuteInt, secondInt, 0, loc)

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: fmt.Sprintf("%d", tp.Unix()),
			},
		},
	}, nil
}
