package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DisableTrace = true
	EnableTrace  = false
)

func InitLogger(disableCaller bool) *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	//config.Encoding = "json"
	config.DisableCaller = disableCaller
	config.DisableStacktrace = disableCaller
	config.EncoderConfig.CallerKey = "call"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.LevelKey = "level"
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
	logger.Sugar()
	return logger
}
