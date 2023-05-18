package main

import (
	"github.com/patsnapops/noop/log"
)

func main() {
	// set with filename
	// log.Default().WithFilename("app.log").Init()
	// set debug level
	log.Default().WithFilename("app.log").WithLevel(log.DebugLevel).Init()
	// or just
	// log.Default().Init()

	log.Debug("this is a simple debugging log")
	log.Warnf("this is a warning log with string %s", "fmt")
	log.Errorf("this is an error level log with string %s", "fmt")
	log.Infof("this is an info level log with string %s", "fmt")
}
