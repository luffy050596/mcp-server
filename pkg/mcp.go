package pkg

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"github.com/pkg/errors"
)

const (
	ModeSSE   string = "sse"
	ModeStdio string = "stdio"
)

type Option func(*Options)

type Options struct {
	Addr string
}

func WithAddr(addr string) Option {
	return func(o *Options) {
		o.Addr = addr
	}
}

func Transport(mode string, opts ...Option) (transport.ServerTransport, error) {
	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	switch mode {
	case ModeSSE:
		return transport.NewSSEServerTransport(options.Addr)
	case ModeStdio:
		return transport.NewStdioServerTransport(), nil
	default:
		return nil, errors.Errorf("invalid mode: %s", mode)
	}
}

func Run(svr *server.Server) {
	errCh := make(chan error)
	go func() {
		slog.Info("MCP Server is running")
		errCh <- svr.Run()
	}()

	if err := signalWaiter(errCh); err != nil {
		slog.Error("server run failed", "error", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := svr.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}
}

func signalWaiter(errCh chan error) error {
	signalToNotify := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}
	if signal.Ignored(syscall.SIGHUP) {
		signalToNotify = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, signalToNotify...)
	<-signals

	select {
	case sig := <-signals:
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			slog.Info("Received signal", "signal", sig)
			// graceful shutdown
			return nil
		}
	case err := <-errCh:
		return err
	}

	return nil
}
