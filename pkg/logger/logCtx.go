package logger

import "context"

type logCtx struct {
	UserID string
}

type keyType int

const key = keyType(0)

func WithLogUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, key, logCtx{UserID: userID})
}
