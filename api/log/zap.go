package log

import (
	"github.com/byuoitav/av-control-api/api"
	"go.uber.org/zap"
)

type Logger struct {
	log *zap.Logger
}

func Wrap(l *zap.Logger) Logger {
	var log Logger

	if l != nil {
		log.log = l.WithOptions(zap.AddCallerSkip(1))
	}

	return log
}

func (l Logger) With(fields ...zap.Field) api.Logger {
	var log Logger

	if l.log != nil {
		log.log = l.log.With(fields...)
	}

	return log
}

func (l Logger) Debug(msg string, fields ...zap.Field) {
	if l.log != nil {
		l.log.Debug(msg, fields...)
	}
}

func (l Logger) Info(msg string, fields ...zap.Field) {
	if l.log != nil {
		l.log.Info(msg, fields...)
	}
}

func (l Logger) Warn(msg string, fields ...zap.Field) {
	if l.log != nil {
		l.log.Warn(msg, fields...)
	}
}

func (l Logger) Error(msg string, fields ...zap.Field) {
	if l.log != nil {
		l.log.Error(msg, fields...)
	}
}
