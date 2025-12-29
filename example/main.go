package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/sower-proxy/deferlog/v2"
)

func main() {
	deferlog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
	})))

	ctx := context.Background()
	defer func() {
		ctx = deferlog.CtxWithLogField(ctx, "trace_id", "I_am_Trace")
		deferlog.InfoContext(ctx, "hello world")
		slog.InfoContext(ctx, "hello world")
	}()

	slog.InfoContext(ctx, "hello world")
}
