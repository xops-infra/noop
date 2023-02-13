package log

import "go.uber.org/zap"

func Trace(msg string, fields ...zap.Field) {
	zapLogger.Debug(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	zapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	zapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	zapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	zapLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	zapLogger.Fatal(msg, fields...)
}

func Debugf(template string, args ...any) {
	zapLogger.Sugar().Debugf(template, args...)
}

func Infof(template string, args ...any) {
	zapLogger.Sugar().Infof(template, args...)
}

func Warnf(template string, args ...any) {
	zapLogger.Sugar().Warnf(template, args...)
}

func Errorf(template string, args ...any) {
	zapLogger.Sugar().Errorf(template, args...)
}

func Fatalf(template string, args ...any) {
	zapLogger.Sugar().Fatalf(template, args...)
}

// not providing Xxxw such as Infow since structured logging should be typed, which Xxxw require reflect
