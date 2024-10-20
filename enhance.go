package deferlog

import (
	"context"
	"os"
)

func InfoWarn(err error, msg string, args ...any) {
	if err != nil {
		Default().Warn(msg, append(args, "err", err)...)
	}
	Default().Info(msg, args...)
}

func InfoWarnContext(ctx context.Context, err error, msg string, args ...any) {
	if err != nil {
		Default().WarnContext(ctx, msg, append(args, "err", err)...)
	}
	Default().InfoContext(ctx, msg, args...)
}

func InfoError(err error, msg string, args ...any) {
	if err != nil {
		Default().Error(msg, append(args, "err", err)...)
	}
	Default().Info(msg, args...)
}

func InfoErrorContext(ctx context.Context, err error, msg string, args ...any) {
	if err != nil {
		Default().ErrorContext(ctx, msg, append(args, "err", err)...)
	}
	Default().InfoContext(ctx, msg, args...)
}

func InfoFatal(err error, msg string, args ...any) {
	if err != nil {
		Default().Error(msg, append(args, "err", err)...)
		os.Exit(1)
	}
	Default().Info(msg, args...)
}

func InfoFatalContext(ctx context.Context, err error, msg string, args ...any) {
	if err != nil {
		Default().ErrorContext(ctx, msg, append(args, "err", err)...)
		os.Exit(1)
	}
	Default().InfoContext(ctx, msg, args...)
}

func DebugWarn(err error, msg string, args ...any) {
	if err != nil {
		Default().Warn(msg, append(args, "err", err)...)
	}
	Default().Debug(msg, args...)
}

func DebugWarnContext(ctx context.Context, err error, msg string, args ...any) {
	if err != nil {
		Default().WarnContext(ctx, msg, append(args, "err", err)...)
	}
	Default().DebugContext(ctx, msg, args...)
}

func DebugError(err error, msg string, args ...any) {
	if err != nil {
		Default().Error(msg, append(args, "err", err)...)
	}
	Default().Debug(msg, args...)
}

func DebugErrorContext(ctx context.Context, err error, msg string, args ...any) {
	if err != nil {
		Default().ErrorContext(ctx, msg, append(args, "err", err)...)
	}
	Default().DebugContext(ctx, msg, args...)
}

func DebugFatal(err error, msg string, args ...any) {
	if err != nil {
		Default().Error(msg, append(args, "err", err)...)
		os.Exit(1)
	}
	Default().Debug(msg, args...)
}

func DebugFatalContext(ctx context.Context, err error, msg string, args ...any) {
	if err != nil {
		Default().ErrorContext(ctx, msg, append(args, "err", err)...)
		os.Exit(1)
	}
	Default().DebugContext(ctx, msg, args...)
}
