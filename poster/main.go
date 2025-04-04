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
	key  string
)

func init() {
	flag.StringVar(&mode, "mode", pkg.ModeStdio, "mode")
	flag.StringVar(&addr, "addr", ":59003", "addr")
	flag.StringVar(&key, "key", "", "bailian api key")
}

func main() {
	flag.Parse()

	transport, err := pkg.Transport(mode, pkg.WithAddr(addr))
	if err != nil {
		panic(err)
	}

	svr, err := server.NewServer(transport, server.WithServerInfo(protocol.Implementation{
		Name:    "poster",
		Version: "1.0.0",
	}))
	if err != nil {
		panic(err)
	}

	svr.RegisterTool(createTool(), createHandler)
	svr.RegisterTool(refineTool(), refineHandler)
	svr.RegisterTool(posterHelpTool(), posterHelpHandler)
	pkg.Run(svr)
}
