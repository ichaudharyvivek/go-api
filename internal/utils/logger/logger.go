package logger

import (
	"context"
	"os"

	"example.com/goapi/internal/config/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ctxKey struct{}

var globalLogger *zap.Logger

func init() {
	appEnv := env.GetString("APP_ENV", "dev")
	if appEnv == "prod" {
		globalLogger = newProductionLogger()
	} else {
		globalLogger = newDevelopmentLogger()
	}
}

// Production logger with file rotation
func newProductionLogger() *zap.Logger {
	// Configure file writer with rotation
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    100,  // MB
		MaxBackups: 7,    // Keep 7 days of logs
		MaxAge:     7,    // Days
		Compress:   true, // Compress rotated files
	})

	// Create core with file output
	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		fileWriter,
		zap.InfoLevel,
	)

	// Optional: Add console output in production
	var core zapcore.Core
	if env.GetBool("LOG_CONSOLE", false) {
		consoleCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zap.InfoLevel,
		)
		core = zapcore.NewTee(fileCore, consoleCore)
	} else {
		core = fileCore
	}

	return zap.New(core)
}

// Development logger with colored console
func newDevelopmentLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// Optional file logging in development
	if env.GetBool("LOG_FILE", false) {
		config.OutputPaths = append(config.OutputPaths, "logs/dev.log")
	}

	return zap.Must(config.Build())
}

// Existing context functions remain the same
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return l
	}
	return globalLogger
}

func NewContext(ctx context.Context, fields ...zap.Field) context.Context {
	return WithContext(ctx, FromContext(ctx).With(fields...))
}
