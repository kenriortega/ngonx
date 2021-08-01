package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

func init() {
	var err error
	config := zap.NewProductionConfig()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.StacktraceKey = ""
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig = encoderConfig

	log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
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
