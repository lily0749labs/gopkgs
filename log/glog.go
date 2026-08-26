package log

// var (
// 	zapLogger *zap.Logger
// 	Debug     = zapLogger.Debug
// 	Info      = zapLogger.Info
// 	Warn      = zapLogger.Warn
// 	Error     = zapLogger.Error
// 	Fatal     = zapLogger.Fatal
// )

func InitLog(opt Option) bool {
	zapLogger = NewLog(opt)
	if zapLogger == nil {
		return false
	}
	Debug = zapLogger.Debug
	Info = zapLogger.Info
	Warn = zapLogger.Warn
	Error = zapLogger.Error
	Fatal = zapLogger.Fatal

	ZapLog = zapLogger
	Debug = zapLogger.Debug
	Info = zapLogger.Info
	Warn = zapLogger.Warn
	Error = zapLogger.Error
	Fatal = zapLogger.Fatal

	return true
}
