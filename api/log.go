package api

import (
	"context"
)

type contextKey int

const (
	_keyRequestID contextKey = iota
)

func RequestID(ctx context.Context) string {
	id, ok := ctx.Value(_keyRequestID).(string)
	if !ok {
		return ""
	}

	return id
}

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, _keyRequestID, id)
}

