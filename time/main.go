package main

import (
	"flag"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/luffy050596/mcp-server/pkg"
)

var (
	mode string
	addr string
)

func init() {
	flag.StringVar(&mode, "mode", pkg.ModeStdio, "mode")
	flag.StringVar(&addr, "addr", ":59001", "addr")
}

func main() {
	transport, err := pkg.Transport(mode, pkg.WithAddr(addr))
	if err != nil {
		panic(err)
	}

	svr, err := server.NewServer(transport, server.WithServerInfo(protocol.Implementation{
		Name:    "time",
		Version: "1.0.0",
	}))
	if err != nil {
		panic(err)
	}

	svr.RegisterTool(currentTimeTool(), currentTimeHandler)
	svr.RegisterTool(timestampTool(), timestampHandler)
	pkg.Run(svr)
}
