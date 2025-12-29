package deferlog

import (
	"context"
	"log/slog"
	"runtime"
	"slices"
)

type CtxKey string

var AutoLogCtxKeys []CtxKey

type DeferWrap struct {
	slog.Handler
	skip int
}

func NewDeferWrap(handler slog.Handler, skip int) *DeferWrap {
	return &DeferWrap{Handler: handler, skip: skip}
}

func (dw *DeferWrap) Handle(ctx context.Context, r slog.Record) error {
	if len(AutoLogCtxKeys) > 0 {
		attrs := make([]slog.Attr, 0, len(AutoLogCtxKeys))
		for _, key := range AutoLogCtxKeys {
			if val := ctx.Value(key); val != nil {
				attrs = append(attrs, slog.Any(string(key), val))
			}
		}
		if len(attrs) > 0 {
			r.AddAttrs(attrs...)
		}
	}

	var pcs [1]uintptr
	runtime.Callers(dw.skip, pcs[:])
	r.PC = pcs[0]

	return dw.Handler.Handle(ctx, r)
}

// CtxWithLogField wraps the context with a key-value pair and automatically
// adds the key to AutoLogCtxKeys if it's not already present.
func CtxWithLogField(ctx context.Context, key string, value any) context.Context {
	ctxKey := CtxKey(key)

	// Add key to AutoLogCtxKeys if not found
	if found := slices.Contains(AutoLogCtxKeys, ctxKey); !found {
		AutoLogCtxKeys = append(AutoLogCtxKeys, ctxKey)
	}

	return context.WithValue(ctx, ctxKey, value)
}
