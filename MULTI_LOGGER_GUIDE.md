# noop 多实例日志功能指南

## 概述

noop 日志库现在支持创建多个独立的日志实例，每个实例可以有不同的配置（文件名、日志级别、字段等）。这个功能通过暴露 `Logger` 结构体实现，同时保持向后兼容性。

## 主要特性

### 1. 多实例支持
- 可以创建多个独立的日志实例
- 每个实例有独立的配置和状态
- 支持并发使用

### 2. 向后兼容
- 原有的 `log.Default().Init()` 用法仍然有效
- 全局函数 `log.Info()`, `log.Debug()` 等仍然可用

### 3. 灵活配置
- 每个实例可以有不同的日志级别
- 支持不同的输出文件
- 可以设置不同的字段和时间格式

### 4. 底层访问
- 可以通过 `ZapLogger()` 方法访问底层的 zap.Logger
- 支持直接使用 zap 的高级功能

## 使用方法

### 基本用法

```go
package main

import (
    "github.com/xops-infra/noop/log"
)

func main() {
    // 方式1：使用默认实例（向后兼容）
    log.Default().WithFilename("app.log").Init()
    log.Info("using default logger")

    // 方式2：创建新的日志实例
    logger := log.New().
        WithFilename("service.log").
        WithLevel(log.DebugLevel).
        Init()
    
    logger.Info("using new logger instance")
}
```

### 多实例示例

```go
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

// 使用不同的日志器
service1Logger.Debug("user service debug message")
service1Logger.Info("user service started")

service2Logger.Debug("order service debug message") // 不会被记录
service2Logger.Info("order service started")
service2Logger.Error("order service error")
```

### 错误和警告分离

```go
errorLogger := log.New().
    WithFilename("main.log").
    WithLevel(log.DebugLevel).
    WithWarnLog("warn.log").
    WithErrorLog("error.log").
    WithHumanTime(time.Local).
    Init()

errorLogger.Debug("application debug")  // 写入 main.log
errorLogger.Warn("application warning") // 写入 warn.log
errorLogger.Error("application error")  // 写入 error.log
```

### 访问底层 zap.Logger

```go
logger := log.New().WithFilename("app.log").Init()

// 获取底层的 zap.Logger
zapLogger := logger.ZapLogger()
zapLogger.Info("direct zap logger usage", zap.String("key", "value"))
```

### 并发使用

```go
logger1 := log.New().WithFilename("concurrent1.log").Init()
logger2 := log.New().WithFilename("concurrent2.log").Init()

// 在不同的 goroutine 中使用
go func() {
    for i := 0; i < 100; i++ {
        logger1.Infof("logger1 message %d", i)
    }
}()

go func() {
    for i := 0; i < 100; i++ {
        logger2.Infof("logger2 message %d", i)
    }
}()
```

## API 参考

### 创建日志实例

```go
// 创建新的日志实例
func New() *Config

// 获取默认配置（向后兼容）
func Default() *Config
```

### 配置方法

```go
// 设置日志文件名
func (c *Config) WithFilename(filename string) *Config

// 设置日志级别
func (c *Config) WithLevel(level Level) *Config

// 设置字段
func (c *Config) WithFields(fields map[string]any) *Config

// 设置人类可读时间
func (c *Config) WithHumanTime(location *time.Location) *Config

// 设置警告日志文件
func (c *Config) WithWarnLog(filename string) *Config

// 设置错误日志文件
func (c *Config) WithErrorLog(filename string) *Config

// 初始化并返回 Logger 实例
func (c *Config) Init() *Logger
```

### Logger 方法

```go
// 结构化日志方法
func (l *Logger) Debug(msg string, fields ...zap.Field)
func (l *Logger) Info(msg string, fields ...zap.Field)
func (l *Logger) Warn(msg string, fields ...zap.Field)
func (l *Logger) Error(msg string, fields ...zap.Field)
func (l *Logger) Fatal(msg string, fields ...zap.Field)
func (l *Logger) Panic(msg string, fields ...zap.Field)

// 格式化日志方法
func (l *Logger) Debugf(template string, args ...any)
func (l *Logger) Infof(template string, args ...any)
func (l *Logger) Warnf(template string, args ...any)
func (l *Logger) Errorf(template string, args ...any)
func (l *Logger) Fatalf(template string, args ...any)
func (l *Logger) Panicf(template string, args ...any)

// 获取底层 zap.Logger
func (l *Logger) ZapLogger() *zap.Logger
```

## 使用场景

### 1. 微服务架构
每个服务组件使用独立的日志实例，便于日志分析和问题排查。

### 2. 不同日志级别
开发环境使用 Debug 级别，生产环境使用 Info 级别。

### 3. 日志分类
将不同类型的日志（访问日志、错误日志、业务日志）写入不同文件。

### 4. 多租户应用
为不同租户创建独立的日志实例。

## 性能考虑

- 每个 Logger 实例都是线程安全的
- 支持高并发写入
- 底层使用 zap，性能优异
- 建议复用 Logger 实例而不是频繁创建

## 迁移指南

### 从单实例迁移

原有代码：
```go
log.Default().WithFilename("app.log").Init()
log.Info("message")
```

新的多实例代码：
```go
logger := log.New().WithFilename("app.log").Init()
logger.Info("message")
```

### 保持向后兼容

如果不想修改现有代码，原有的全局函数仍然可用：
```go
log.Default().WithFilename("app.log").Init()
log.Info("message") // 仍然有效
```

## 测试

运行测试以验证功能：

```bash
# 运行所有测试
go test ./log -v

# 运行特定测试
go test ./log -run TestMultipleInstances -v

# 运行基准测试
go test ./log -bench=. -v
```

## 示例项目

查看 `example_multi_logger.go` 文件获取完整的使用示例。
