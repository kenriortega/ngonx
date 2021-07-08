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

func LogInfo(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}
func LogDebug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}
func LogError(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}
func LogWarn(message string, fields ...zap.Field) {
	log.Warn(message, fields...)
}
