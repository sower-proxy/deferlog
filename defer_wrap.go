package deferlog

import (
	"context"
	"log/slog"
	"runtime"
)

type DeferWrap struct {
	slog.Handler

	Skip int
}

func NewDeferWrap(handler slog.Handler, skip int) *DeferWrap {
	return &DeferWrap{
		Handler: handler,
		Skip:    skip,
	}
}

var AutoLogCtxKeys []string

func (dw *DeferWrap) Handle(ctx context.Context, r slog.Record) error {
	attrs := make([]slog.Attr, 0, len(AutoLogCtxKeys))
	for _, key := range AutoLogCtxKeys {
		if val := ctx.Value(key); val != nil {
			attrs = append(attrs, slog.Any(key, val))
		}
	}
	r.AddAttrs(attrs...)

	pcs := [1]uintptr{r.PC}
	runtime.Callers(dw.Skip, pcs[:])
	r.PC = pcs[0]
	return dw.Handler.Handle(ctx, r)
}
