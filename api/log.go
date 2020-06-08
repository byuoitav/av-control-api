package api

import (
	"context"

	"go.uber.org/zap"
)

type Logger interface {
	Debug(string, ...zap.Field)
	Info(string, ...zap.Field)
	Warn(string, ...zap.Field)
	Error(string, ...zap.Field)

	With(...zap.Field) Logger
}

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

