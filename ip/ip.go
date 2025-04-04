package main

import (
	"io"
	"net"
	"net/http"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/pkg/errors"
)

func ipGeoTool() *protocol.Tool {
	t := &protocol.Tool{
		Name:        "ip",
		Description: "Query the geo location of the ip address",
		InputSchema: protocol.InputSchema{
			Type: protocol.Object,
			Properties: map[string]any{
				"ip": map[string]string{
					"type":        "string",
					"description": "The ip address to query, default is the ip address of the current machine",
				},
			},
		},
	}
	return t
}

func ipGeoHandler(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	ip, ok := request.Arguments["ip"].(string)
	if !ok {
		return nil, errors.Errorf("ip is required")
	}

	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return nil, errors.Errorf("invalid ip address: %s", ip)
	}

	resp, err := http.Get("https://ip.rpcx.io/api/ip?ip=" + ip)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get ip geo location")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read ip geo location")
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			protocol.TextContent{
				Type: "text",
				Text: string(body),
			},
		},
	}, nil
}
