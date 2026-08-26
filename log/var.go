package log

import (
	"go.uber.org/zap"
)

var (
	zapLogger *zap.Logger
	ZapLog    *zap.Logger
	Debug     = zapLogger.Debug
	Info      = zapLogger.Info
	Warn      = zapLogger.Warn
	Error     = zapLogger.Error
	Fatal     = zapLogger.Fatal
)
