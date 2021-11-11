package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger

func init() {
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.CallerKey = "caller"
	encoderConfig.StacktraceKey = ""
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig = encoderConfig

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./ngonx-log/ngonx.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(config.EncoderConfig),
		w,
		zap.InfoLevel,
	)
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// LogInfo wrap for log.info
func LogInfo(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}

// LogDebug wrap for log.Debug
func LogDebug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}

// LogError wrap for log.Error
func LogError(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}

// LogWarn wrap for log.Warn
func LogWarn(message string, fields ...zap.Field) {
	log.Warn(message, fields...)
}
