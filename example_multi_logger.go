package main

import (
	"time"

	"github.com/xops-infra/noop/log"
)

func main() {
	// 方式1：使用默认实例（向后兼容）
	log.Default().WithFilename("app.log").Init()
	log.Info("using default logger")

	// 方式2：创建独立的日志实例
	// 服务1的日志器 - 记录调试信息
	service1Logger := log.New().
		WithFilename("service1.log").
		WithLevel(log.DebugLevel).
		WithFields(map[string]any{
			"service": "user-service",
			"version": "1.0.0",
		}).
		Init()

	// 服务2的日志器 - 只记录重要信息
	service2Logger := log.New().
		WithFilename("service2.log").
		WithLevel(log.InfoLevel).
		WithFields(map[string]any{
			"service": "order-service",
			"version": "2.1.0",
		}).
		Init()

	// 错误日志器 - 分离错误和警告日志
	errorLogger := log.New().
		WithFilename("main.log").
		WithLevel(log.DebugLevel).
		WithWarnLog("warn.log").
		WithErrorLog("error.log").
		WithHumanTime(time.Local).
		Init()

	// 使用不同的日志器
	service1Logger.Debug("user service debug message")
	service1Logger.Info("user service started")

	service2Logger.Debug("order service debug message") // 不会被记录，因为level是Info
	service2Logger.Info("order service started")
	service2Logger.Error("order service error")

	errorLogger.Debug("application debug")
	errorLogger.Warn("application warning")
	errorLogger.Error("application error")

	// 方式3：直接使用底层的 zap.Logger
	zapLogger := service1Logger.ZapLogger()
	zapLogger.Info("direct zap logger usage")

	// 方式4：在不同的goroutine中使用不同的日志器
	go func() {
		asyncLogger := log.New().
			WithFilename("async.log").
			WithFields(map[string]any{
				"goroutine": "worker-1",
			}).
			Init()
		
		for i := 0; i < 5; i++ {
			asyncLogger.Infof("async message %d", i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 等待异步日志完成
	time.Sleep(1 * time.Second)
	
	log.Info("all loggers demonstration completed")
}
