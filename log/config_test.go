package log

import (
	"testing"
)

func TestInit(t *testing.T) {
	Default().WithFilename("app.log").WithLevel(DebugLevel).Init()
	Debug("this is a debug message")
	Info("this is a info message")
}
