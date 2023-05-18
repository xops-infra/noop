package log

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestInit(t *testing.T) {
	Default().WithFilename("app.log").WithLevel(DebugLevel).Init()
	Debug("this is a debug message")
	Info("this is a info message")
}

func TestConfig_WithFields(t *testing.T) {
	fields := []zap.Field{
		{
			Key:    "instanceID",
			Type:   zapcore.StringType,
			String: "1135286d-2fa7-4715-8b90-4937c0e49c2d",
		},
		{
			Key:    "LocalTime",
			Type:   zapcore.StringType,
			String: time.Now().In(time.Local).Format("2006-01-02 15:04:05.000"),
		},
	}
	Default().WithFilename("app.log").WithLevel(DebugLevel).WithFields(fields).Init()
	Debug("this is a debug message with fields")
	Info("this is a info message with fields")
}
