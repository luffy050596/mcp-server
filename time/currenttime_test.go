package main

import (
	"testing"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/stretchr/testify/require"
)

func TestCurrentTimeHandler(t *testing.T) {
	reply, err := currentTimeHandler(&protocol.CallToolRequest{
		Arguments: map[string]any{
			"timezone": "Asia/Tokyo",
		},
	})
	require.NoError(t, err)
	require.NotNil(t, reply)
}
