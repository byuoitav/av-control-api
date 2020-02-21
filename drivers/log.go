package drivers

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// P is a plain zap logger
var Log *zap.SugaredLogger

// Config is the logger config used for P
var Config zap.Config

func init() {
	Config = zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "@",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, err := Config.Build()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}

	Log = l.Sugar()
	Log.Info("Zap logger started")
	_ = l.Sync()
}
