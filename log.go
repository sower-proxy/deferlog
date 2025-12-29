package deferlog

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"
)

var defaultLogger atomic.Pointer[slog.Logger]

func Default() (l *slog.Logger) {
	if l = defaultLogger.Load(); l == nil {
		SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	}
	return defaultLogger.Load()
}

func SetDefault(l *slog.Logger) {
	if l == nil {
		panic("deferlog: SetDefault(nil)")
	}

	// Unwrap DeferWrap to get the underlying handler
	h := l.Handler()
	if dw, ok := h.(*DeferWrap); ok {
		h = dw.Handler
	}

	// Use the unwrapped handler to avoid double-wrapping
	slog.SetDefault(slog.New(NewDeferWrap(h, 4)))
	defaultLogger.Store(slog.New(NewDeferWrap(h, 6)))
}

func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	Default().DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	Default().InfoContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	Default().WarnContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	Default().ErrorContext(ctx, msg, args...)
}

func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	Default().Log(ctx, level, msg, args...)
}

func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	Default().LogAttrs(ctx, level, msg, attrs...)
}
