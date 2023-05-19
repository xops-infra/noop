package main

import (
	"github.com/patsnapops/noop/log"
)

func main() {
	// set with filename
	// log.Default().WithFilename("app.log").Init()

	// set debug level
	// log.Default().WithFilename("app.log").WithLevel(log.DebugLevel).Init()

	// set fields, which will be printed in the log, the fields is a map,
	// key is the field name, value is the field value, the value can be any type, example:
	//fields := map[string]interface{}{
	//	"instance_id": "1135286d-2fa7-4715-8b90-4937c0e49c2d",
	//	"localtime":   time.Now().In(time.Local).Format("2006-01-02 15:04:05.000"),
	//}
	//log.Default().WithFilename("app.log").WithLevel(log.DebugLevel).WithFields(fields).Init()

	// set human time, which will be printed in the log, default is local time, example:
	log.Default().WithHumanTime(nil).Init()

	// or just
	// log.Default().Init()

	log.Debug("this is a simple debugging log")
	log.Warnf("this is a warning log with string %s", "fmt")
	log.Errorf("this is an error level log with string %s", "fmt")
	log.Infof("this is an info level log with string %s", "fmt")
}
