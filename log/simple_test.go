package log

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewLoggerInstance(t *testing.T) {
	// 测试创建新的日志实例
	logger := New().WithLevel(InfoLevel).Init()
	
	if logger == nil {
		t.Error("New logger instance should not be nil")
	}
	
	// 测试日志方法
	logger.Info("test info message")
	logger.Error("test error message")
}

func TestLoggerMethods(t *testing.T) {
	logger := New().WithLevel(DebugLevel).Init()
	
	// 测试所有日志方法
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
	
	// 测试格式化方法
	logger.Debugf("debug %s", "formatted")
	logger.Infof("info %s", "formatted")
	logger.Warnf("warn %s", "formatted")
	logger.Errorf("error %s", "formatted")
}

func TestZapLoggerAccess(t *testing.T) {
	logger := New().Init()
	
	zapLogger := logger.ZapLogger()
	if zapLogger == nil {
		t.Error("ZapLogger() should return a valid zap.Logger")
	}
	
	// 测试直接使用 zap.Logger
	zapLogger.Info("direct zap message", zap.String("key", "value"))
}

func TestMultipleInstances(t *testing.T) {
	// 创建多个独立实例
	logger1 := New().WithLevel(DebugLevel).Init()
	logger2 := New().WithLevel(InfoLevel).Init()
	logger3 := New().WithLevel(ErrorLevel).Init()
	
	// 测试它们是独立的
	logger1.Debug("logger1 debug") // 应该显示
	logger2.Debug("logger2 debug") // 不应该显示，因为level是Info
	logger3.Debug("logger3 debug") // 不应该显示，因为level是Error
	
	logger1.Info("logger1 info")  // 应该显示
	logger2.Info("logger2 info")  // 应该显示
	logger3.Info("logger3 info")  // 不应该显示，因为level是Error
	
	logger1.Error("logger1 error") // 应该显示
	logger2.Error("logger2 error") // 应该显示
	logger3.Error("logger3 error") // 应该显示
}

func TestWithFields(t *testing.T) {
	fields := map[string]any{
		"service": "test-service",
		"version": "1.0.0",
	}
	
	logger := New().WithFields(fields).Init()
	logger.Info("message with fields")
}

func TestWithHumanTime(t *testing.T) {
	logger := New().WithHumanTime(time.UTC).Init()
	logger.Info("message with human time")
}

func TestBackwardCompatibilitySimple(t *testing.T) {
	// 测试向后兼容性
	defaultLogger := Default().Init()
	
	if defaultLogger == nil {
		t.Error("Default logger should not be nil")
	}
	
	// 使用全局函数
	Info("backward compatibility test")
	Debugf("formatted %s", "message")
}

func TestConcurrentAccess(t *testing.T) {
	logger1 := New().Init()
	logger2 := New().Init()
	
	done := make(chan bool, 2)
	
	// 并发写入测试
	go func() {
		for i := 0; i < 10; i++ {
			logger1.Infof("concurrent logger1 %d", i)
		}
		done <- true
	}()
	
	go func() {
		for i := 0; i < 10; i++ {
			logger2.Infof("concurrent logger2 %d", i)
		}
		done <- true
	}()
	
	// 等待完成
	<-done
	<-done
}

func BenchmarkNewLogger(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger := New().Init()
		logger.Info("benchmark message")
	}
}

func BenchmarkLoggerMethods(b *testing.B) {
	logger := New().Init()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("benchmark parallel message")
		}
	})
}
