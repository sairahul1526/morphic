package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var err error

func SetLevel(level string) {
	defaultConfig := zap.NewProductionConfig()
	defaultConfig.Level = zap.NewAtomicLevelAt(getZapLogLevelFromString(level))
	logger, err = zap.NewProductionConfig().Build()
	if err != nil {
		fmt.Println("Failed to create logger")
	}
}

func getZapLogLevelFromString(level string) zapcore.Level {
	l, err := zapcore.ParseLevel(level)
	if err != nil {
		l = zap.InfoLevel
	}
	return l
}

func GetLogger() *zap.Logger {
	return logger
}

func Info(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Info(msg, fields...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Debug(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Error(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Warn(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Fatal(msg, fields...)
	}
}

func Panic(msg string, fields ...zap.Field) {
	if logger != nil {
		logger.Panic(msg, fields...)
	}
}
