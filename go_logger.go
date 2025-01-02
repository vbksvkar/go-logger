package gologger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerCtxKey struct{}

func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, logger)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(loggerCtxKey{}).(*zap.SugaredLogger)
	if !ok {
		return zap.NewNop().Sugar()
	}
	return logger
}

func New(appName string, appVersion string) (*zap.SugaredLogger, error) {
	logLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "@timestamp",
		NameKey:        "logger_name",
		CallerKey:      "caller",
		StacktraceKey:  "stack_trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	config := zap.Config{
		Level:            logLevel,
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	l, err := config.Build()
	if err != nil {
		return nil, err
	}
	return l.Sugar().With("app_name", appName, "app_version", appVersion), nil
}
