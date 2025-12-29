package deferlog

import (
	"context"
	"log/slog"
	"runtime"
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

	if dw.skip > 0 {
		var pcs [1]uintptr
		runtime.Callers(dw.skip, pcs[:])
		r.PC = pcs[0]
	}

	return dw.Handler.Handle(ctx, r)
}
