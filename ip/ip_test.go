package main

import (
	"testing"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIpGeoHandler(t *testing.T) {
	reply, err := ipGeoHandler(&protocol.CallToolRequest{
		Arguments: map[string]any{
			"ip": "8.8.8.8",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, reply)

	text := reply.Content[0].(protocol.TextContent).Text
	assert.Contains(t, text, "Google Public DNS")
}
