package main

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/patsnapops/noop/log"
)

func main() {
	// set with filename
	// log.Default().WithFilename("app.log").Init()

	// set debug level
	// log.Default().WithFilename("app.log").WithLevel(log.DebugLevel).Init()

	// set fields
	log.Default().WithFilename("app.log").WithLevel(log.DebugLevel).WithFields([]zap.Field{
		{
			Key:    "xx_key",
			Type:   zapcore.StringType,
			String: "filed_value",
		},
		{
			Key:    "LocalTime",
			Type:   zapcore.StringType,
			String: time.Now().In(time.Local).Format("2006-01-02 15:04:05.000"),
		},
	}).Init()

	// or just
	// log.Default().Init()

	log.Debug("this is a simple debugging log")
	log.Warnf("this is a warning log with string %s", "fmt")
	log.Errorf("this is an error level log with string %s", "fmt")
	log.Infof("this is an info level log with string %s", "fmt")
}
