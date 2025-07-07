package log

import "go.uber.org/zap"

// 向后兼容的全局函数，使用默认实例
func Trace(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Trace(msg, fields...)
	}
}

func Debug(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, fields...)
	}
}

func Info(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg, fields...)
	}
}

func Error(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, fields...)
	}
}

func Fatal(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Fatal(msg, fields...)
	}
}

func Panic(msg string, fields ...zap.Field) {
	if defaultLogger != nil {
		defaultLogger.Panic(msg, fields...)
	}
}

func Debugf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Debugf(template, args...)
	}
}

func Infof(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Infof(template, args...)
	}
}

func Warnf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Warnf(template, args...)
	}
}

func Errorf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Errorf(template, args...)
	}
}

func Fatalf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Fatalf(template, args...)
	}
}

func Panicf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Panicf(template, args...)
	}
}

// not providing Xxxw such as Infow since structured logging should be typed, which Xxxw require reflect
