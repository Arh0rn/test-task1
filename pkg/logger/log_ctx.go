package logger

import "context"

type logCtx struct {
	UserID    string
	RequestID string
}

type keyType int

const key = keyType(0)

func WithLogUserID(ctx context.Context, userID string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.UserID = userID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{UserID: userID})
}

func WithLogRequestID(ctx context.Context, requestID string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.RequestID = requestID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{RequestID: requestID})
}
