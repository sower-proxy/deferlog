package deferlog

import (
	"context"
	"log/slog"
	"runtime"
)

var _ slog.Handler = (*DeferWrap)(nil)

type DeferWrap struct {
	slog.Handler

	Skip int
}

func (d *DeferWrap) Handle(ctx context.Context, r slog.Record) error {
	pcs := [1]uintptr{r.PC}
	runtime.Callers(d.Skip, pcs[:])
	r.PC = pcs[0]
	return d.Handler.Handle(ctx, r)
}
