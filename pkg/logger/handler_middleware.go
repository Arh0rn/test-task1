package logger

import (
	"context"
	"log/slog"
)

// Don't confuse with router logging middleware.
// This is a middleware for slog standard handler.

type SlogHandlerMiddleware struct {
	next slog.Handler
}

func NewSlogHandlerMiddleware(next slog.Handler) *SlogHandlerMiddleware {
	return &SlogHandlerMiddleware{next: next}
}

func (m *SlogHandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return m.next.Enabled(ctx, rec)
}

func (m *SlogHandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		rec.Add("UserID", c.UserID)
	}
	return m.next.Handle(ctx, rec)
}

func (m *SlogHandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogHandlerMiddleware{next: m.next.WithAttrs(attrs)}
}

func (m *SlogHandlerMiddleware) WithGroup(name string) slog.Handler {
	return &SlogHandlerMiddleware{next: m.next.WithGroup(name)}
}
