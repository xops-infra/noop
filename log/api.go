package log

import "go.uber.org/zap"

// 向后兼容的全局函数，使用默认实例
func Trace(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Debug(msg, fields...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Fatal(msg, fields...)
	}
}

func Panic(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Panic(msg, fields...)
	}
}

func Debugf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Sugar().Debugf(template, args...)
	}
}

func Infof(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Sugar().Infof(template, args...)
	}
}

func Warnf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Sugar().Warnf(template, args...)
	}
}

func Errorf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Sugar().Errorf(template, args...)
	}
}

func Fatalf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Sugar().Fatalf(template, args...)
	}
}

func Panicf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.zapLogger.Sugar().Panicf(template, args...)
	}
}

// not providing Xxxw such as Infow since structured logging should be typed, which Xxxw require reflect
