package log

import (
	"os"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestMultipleLoggerInstances(t *testing.T) {
	// 获取当前日期格式
	expectedDate := time.Now().In(time.Local).Format("2006-01-02")
	
	// 清理测试文件
	defer func() {
		files := []string{
			"service1_" + expectedDate + ".log",
			"service2_" + expectedDate + ".log",
			"app_" + expectedDate + ".log",
			"service1.log", "service2.log", "app.log", // 备用清理
		}
		for _, file := range files {
			os.Remove(file)
		}
	}()

	// 创建第一个日志实例 - 用于服务1
	logger1 := New().WithFilename("service1.log").WithLevel(DebugLevel).Init()
	
	// 创建第二个日志实例 - 用于服务2
	logger2 := New().WithFilename("service2.log").WithLevel(InfoLevel).Init()

	// 测试两个实例独立工作
	logger1.Debug("service1 debug message")
	logger1.Info("service1 info message")
	
	logger2.Debug("service2 debug message") // 这个不应该被记录，因为level是Info
	logger2.Info("service2 info message")
	logger2.Error("service2 error message")

	// 验证文件是否创建（检查带日期的文件名）
	expectedFiles := []string{
		"service1_" + expectedDate + ".log",
		"service2_" + expectedDate + ".log",
	}
	
	for _, expectedFile := range expectedFiles {
		if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
			// 如果带日期的文件不存在，检查原始文件名
			originalFile := expectedFile[:len(expectedFile)-15] + ".log" // 移除日期部分
			if _, err := os.Stat(originalFile); os.IsNotExist(err) {
				t.Errorf("%s or %s should be created", expectedFile, originalFile)
			}
		}
	}
}

func TestLoggerWithFields(t *testing.T) {
	defer os.Remove("fields_test.log")

	fields := map[string]any{
		"service_id": "test-service-123",
		"version":    "1.0.0",
		"env":        "test",
	}

	logger := New().
		WithFilename("fields_test.log").
		WithLevel(DebugLevel).
		WithFields(fields).
		Init()

	logger.Info("test message with fields")
	logger.Error("test error with fields")

	// 验证文件创建
	if _, err := os.Stat("fields_test.log"); os.IsNotExist(err) {
		t.Error("fields_test.log should be created")
	}
}

func TestLoggerWithHumanTime(t *testing.T) {
	defer os.Remove("human_time_test.log")

	logger := New().
		WithFilename("human_time_test.log").
		WithLevel(DebugLevel).
		WithHumanTime(time.UTC).
		Init()

	logger.Info("test message with human time")

	// 验证文件创建
	if _, err := os.Stat("human_time_test.log"); os.IsNotExist(err) {
		t.Error("human_time_test.log should be created")
	}
}

func TestLoggerWithLevelFilter(t *testing.T) {
	defer func() {
		os.Remove("main_test.log")
		os.Remove("warn_test.log")
		os.Remove("error_test.log")
	}()

	logger := New().
		WithFilename("main_test.log").
		WithLevel(DebugLevel).
		WithWarnLog("warn_test.log").
		WithErrorLog("error_test.log").
		Init()

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// 验证所有文件都被创建
	files := []string{"main_test.log", "warn_test.log", "error_test.log"}
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("%s should be created", file)
		}
	}
}

func TestZapLoggerExposure(t *testing.T) {
	defer os.Remove("zap_test.log")

	logger := New().WithFilename("zap_test.log").Init()
	
	// 测试暴露的 zap.Logger
	zapLogger := logger.ZapLogger()
	if zapLogger == nil {
		t.Error("ZapLogger() should return a valid zap.Logger")
	}

	// 直接使用 zap.Logger
	zapLogger.Info("direct zap logger message", zap.String("key", "value"))
}

func TestBackwardCompatibility(t *testing.T) {
	// 测试向后兼容性 - 原有的使用方式应该仍然工作
	Default().WithFilename("backward_test.log").Init()

	// 使用全局函数
	Debug("backward compatibility debug")
	Info("backward compatibility info")
	Warnf("backward compatibility warn with %s", "format")
	Errorf("backward compatibility error with %s", "format")

	// 由于文件名会被自动添加日期，我们需要检查实际生成的文件名
	// 获取当前日期格式
	expectedDate := time.Now().In(time.Local).Format("2006-01-02")
	expectedFilename := "backward_test_" + expectedDate + ".log"
	
	// 验证文件创建
	if _, err := os.Stat(expectedFilename); os.IsNotExist(err) {
		// 如果带日期的文件不存在，检查原始文件名
		if _, err := os.Stat("backward_test.log"); os.IsNotExist(err) {
			t.Error("backward compatibility log file should be created")
		} else {
			// 清理原始文件
			defer os.Remove("backward_test.log")
		}
	} else {
		// 清理带日期的文件
		defer os.Remove(expectedFilename)
	}
}

func TestConcurrentLoggers(t *testing.T) {
	defer func() {
		os.Remove("concurrent1.log")
		os.Remove("concurrent2.log")
	}()

	logger1 := New().WithFilename("concurrent1.log").Init()
	logger2 := New().WithFilename("concurrent2.log").Init()

	// 并发测试
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			logger1.Infof("logger1 message %d", i)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			logger2.Infof("logger2 message %d", i)
		}
		done <- true
	}()

	// 等待两个goroutine完成
	<-done
	<-done

	// 验证文件创建
	files := []string{"concurrent1.log", "concurrent2.log"}
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("%s should be created", file)
		}
	}
}

func BenchmarkMultipleLoggers(b *testing.B) {
	logger1 := New().WithFilename("bench1.log").Init()
	logger2 := New().WithFilename("bench2.log").Init()
	
	defer func() {
		os.Remove("bench1.log")
		os.Remove("bench2.log")
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger1.Info("benchmark message 1")
			logger2.Info("benchmark message 2")
		}
	})
}
