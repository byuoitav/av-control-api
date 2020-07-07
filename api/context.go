package api

import (
	"context"
)

type contextKey int

const (
	_keyRequestID contextKey = iota
)

// CtxRequestID pulls a request ID from a context.Context.
func CtxRequestID(ctx context.Context) string {
	id, ok := ctx.Value(_keyRequestID).(string)
	if !ok {
		return ""
	}

	return id
}

// WithRequestID returns a new context.Context, based on ctx, with the request id set.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, _keyRequestID, id)
}
