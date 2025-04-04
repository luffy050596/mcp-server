package main

import (
	"testing"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimestampHandler(t *testing.T) {
	reply, err := timestampHandler(&protocol.CallToolRequest{
		Arguments: map[string]any{
			"year":     "2025",
			"month":    "4",
			"day":      "4",
			"hour":     "12",
			"minute":   "30",
			"second":   "0",
			"timezone": "Asia/Tokyo",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, reply)
	assert.Equal(t, reply.Content[0].(protocol.TextContent).Text, "1743737400")
}
