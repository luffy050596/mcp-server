package main

import (
	"io"
	"net"
	"net/http"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/luffy050596/mcp-server/pkg"
	"github.com/pkg/errors"
)

func myIpTool() *protocol.Tool {
	t := &protocol.Tool{
		Name:        "my_ip",
		Description: "Get the ip address of the current machine",
		InputSchema: protocol.InputSchema{
			Type: protocol.Object,
		},
	}
	return t
}

func myIpHandler(request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	if mode != pkg.ModeStdio {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				protocol.TextContent{
					Type: "text",
					Text: "Only support stdio mode",
				},
			},
		}, nil
	}

	// get the ip address of the current machine
	resp, err := http.Get("https://ipinfo.io/ip")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "can not get the ip address of the current machine")
	}

	ipAddr := net.ParseIP(string(body))
	if ipAddr == nil {
		return nil, errors.Errorf("can not get the ip address of the current machine")
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
