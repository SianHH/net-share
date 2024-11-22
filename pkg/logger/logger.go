package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type Logger struct {
	*zap.Logger
}

func NewLogger(logOutPath string, level int, console bool) *Logger {
	// 日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch level {
	case 0:
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case 1:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case 2:
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case 3:
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	case 4:
		atomicLevel.SetLevel(zapcore.DPanicLevel)
	case 5:
		atomicLevel.SetLevel(zapcore.PanicLevel)
	case 6:
		atomicLevel.SetLevel(zapcore.FatalLevel)
	}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "line",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 日志轮转
	writer := &lumberjack.Logger{
		// 日志名称
		Filename: logOutPath,
		// 日志大小限制，单位MB
		MaxSize: 5,
		// 历史日志文件保留天数
		MaxAge: 30,
		// 最大保留历史日志数量
		MaxBackups: 5,
		// 本地时区
		LocalTime: true,
		// 历史日志文件压缩标识
		Compress: false,
	}

	var core zapcore.Core
	if console {
		// 进入开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		// zapcore.NewTee可以实现日志多重输出，可以指定不同的日志级别输出到不同的位置
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), atomicLevel),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), atomicLevel)
	}

	return &Logger{zap.New(core, zap.AddCaller())}
}
