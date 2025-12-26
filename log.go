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

// SetDefault makes l the default [Logger], which is used by
// the top-level functions [Info], [Debug] and so on.
// After this call, output from the log package's default Logger
// (as with [log.Print], etc.) will be logged using l's Handler,
// at a level controlled by [SetLogLoggerLevel].
func SetDefault(l *slog.Logger) {
	if l == nil {
		panic("deferlog: SetDefault(nil)")
	}

	h := l.Handler()
	if dw, ok := h.(*DeferWrap); ok {
		h = dw.Handler
	}

	slog.SetDefault(slog.New(NewDeferWrap(l.Handler(), 4)))
	defaultLogger.Store(slog.New(NewDeferWrap(l.Handler(), 6)))
}

// Debug calls [Logger.Debug] on the default logger.
func Debug(msg string, args ...any) {
	Default().Debug(msg, args...)
}

// DebugContext calls [Logger.DebugContext] on the default logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	Default().DebugContext(ctx, msg, args...)
}

// Info calls [Logger.Info] on the default logger.
func Info(msg string, args ...any) {
	Default().Info(msg, args...)
}

// InfoContext calls [Logger.InfoContext] on the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	Default().InfoContext(ctx, msg, args...)
}

// Warn calls [Logger.Warn] on the default logger.
func Warn(msg string, args ...any) {
	Default().Warn(msg, args...)
}

// WarnContext calls [Logger.WarnContext] on the default logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	Default().WarnContext(ctx, msg, args...)
}

// Error calls [Logger.Error] on the default logger.
func Error(msg string, args ...any) {
	Default().Error(msg, args...)
}

// ErrorContext calls [Logger.ErrorContext] on the default logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	Default().ErrorContext(ctx, msg, args...)
}

// Log calls [Logger.Log] on the default logger.
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	Default().Log(ctx, level, msg, args...)
}

// LogAttrs calls [Logger.LogAttrs] on the default logger.
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	Default().LogAttrs(ctx, level, msg, attrs...)
}
