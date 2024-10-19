package deferlog_test

import (
	"log/slog"

	"github.com/sower-proxy/deferlog/v2"
)

func ExampleInfo() {
	// deferlog.SetDefault(slog.Default())

	defer func() {
		deferlog.Info("hello world")
	}()

	slog.Info("hello world")

	return

	// Output: hello world
}
