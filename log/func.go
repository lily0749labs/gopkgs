package log

import (
	"fmt"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder, //timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

type Field = zap.Field

func Any(key string, value any) Field {
	return zap.Any(key, value)
}

func Args(args any) Field {
	return Any("Args", args)
}

func Cursor(cursor string) Field {
	return Any("Cursor", cursor)
}

func Data(data any) Field {
	return Any("DataName", data)
}

func Err(e string) Field {
	return Any("Error", e)
}

func E(err error) Field {
	return Any("Error", err)
}

func Line() Field {
	// 获取上层调用者PC，文件名，所在行
	funcName, file, line, ok := runtime.Caller(1)
	var info = ""
	if ok {
		info = fmt.Sprintf("file:%s func:%s line:%d ", file, runtime.FuncForPC(funcName).Name(), line)
	}
	return Any("Line", info)
}
