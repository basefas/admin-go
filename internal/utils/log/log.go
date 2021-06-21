package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logLevel = "debug"
var zapLog *zap.Logger

func Init() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.TimeKey = "time"
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	//encoder := zapcore.NewJSONEncoder(encoderConfig)
	writer := zapcore.AddSync(os.Stdout)
	level := getLoggerLevel(logLevel)

	core := zapcore.NewCore(encoder, writer, level)
	development := zap.Development()
	caller := zap.AddCaller()
	skip := zap.AddCallerSkip(1)

	zapLog = zap.New(core, development, caller, skip)
	defer zapLog.Sync()
}

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.DebugLevel
}

func Debug(msg string, fields ...zap.Field) {
	zapLog.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	zapLog.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	zapLog.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	zapLog.Error(msg, fields...)
}
