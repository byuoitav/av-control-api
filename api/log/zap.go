package log

import (
	"github.com/byuoitav/av-control-api/api"
	"go.uber.org/zap"
)

type logger struct {
	*zap.Logger
}

func Wrap(l *zap.Logger) logger {
	return logger{
		Logger: l.WithOptions(zap.AddCallerSkip(1)),
	}
}

func (l logger) With(fields ...zap.Field) api.Logger {
	return logger{
		Logger: l.Logger.With(fields...),
	}
}

func (l logger) Debug(msg string, fields ...zap.Field) {
	if l.Logger != nil {
		l.Logger.Debug(msg, fields...)
	}
}

func (l logger) Info(msg string, fields ...zap.Field) {
	if l.Logger != nil {
		l.Logger.Info(msg, fields...)
	}
}

func (l logger) Warn(msg string, fields ...zap.Field) {
	if l.Logger != nil {
		l.Logger.Warn(msg, fields...)
	}
}

func (l logger) Error(msg string, fields ...zap.Field) {
	if l.Logger != nil {
		l.Logger.Error(msg, fields...)
	}
}
