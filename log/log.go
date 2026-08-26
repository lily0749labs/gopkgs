package log

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dromara/carbon/v2"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLog(opt Option) *zap.Logger {

	allCore := []zapcore.Core{}
	if opt.IsFile {
		if _, err := os.Stat(opt.SavePath); os.IsNotExist(err) {
			os.Mkdir(opt.SavePath, 0777)
			os.Chmod(opt.SavePath, 0777)
		}

		filename := filepath.Join(opt.SavePath, fmt.Sprint("log_", carbon.Now().ToRfc3339String(), ".log"))
		// file, _ := os.Create(flienema)

		fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filename, //日志文件存放目录
			MaxSize:    100,      //文件大小限制,单位MB
			MaxBackups: 100,      //最大保留日志文件数量
			MaxAge:     30,       //日志文件保留天数
			Compress:   false,    //是否压缩处理
		})
		allCore = append(allCore, zapcore.NewCore(
			zapcore.NewConsoleEncoder(newEncoderConfig()),
			fileWriteSyncer,
			// zapcore.AddSync(file),
			// zapcore.NewMultiWriteSyncer(zapcore.AddSync(file)),
			zapcore.DebugLevel))
	}
	if opt.IsConsole {
		allCore = append(allCore, zapcore.NewCore(
			zapcore.NewConsoleEncoder(newEncoderConfig()),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel))

	}

	core := zapcore.NewTee(allCore...)

	// core := zapcore.NewCore(
	// 	zapcore.NewConsoleEncoder(newEncoderConfig()),
	// 	// zapcore.NewJSONEncoder(newEncoderConfig()),
	// 	// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(file)),
	// 	zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), fileWriteSyncer),
	// 	// zapcore.NewMultiWriteSyncer(zapcore.AddSync(file)),
	// 	zap.DebugLevel,
	// )

	return zap.New(core, zap.AddCaller())
}
