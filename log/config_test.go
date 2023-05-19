package log

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	Default().WithFilename("app.log").WithLevel(DebugLevel).Init()
	Debug("this is a debug message")
	Info("this is a info message")
}

func TestConfig_WithFields(t *testing.T) {
	fields := map[string]any{
		"instance_id": "1135286d-2fa7-4715-8b90-4937c0e49c2d",
		"localtime":   time.Now().In(time.Local).Format("2006-01-02 15:04:05.000"),
	}
	Default().WithFilename("app.log").WithLevel(DebugLevel).WithFields(fields).Init()
	Debug("this is a debug message with fields")
	Info("this is a info message with fields")
}

func TestConfig_WithHumanTime(t *testing.T) {
	Default().WithFilename("app.log").WithLevel(DebugLevel).WithHumanTime(time.Local).Init()
	Debug("this is a debug message with human time")
	Info("this is a info message with human time")
}
