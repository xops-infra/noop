# noop
**No** **o**ption za**p** wrapper

## Example
```go
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

```

Output:
```bash
[root@linux noop]# go run main.go 
2023-03-24T11:22:13.071+0800    INFO    noop/main.go:9  zap logger initialized
2023-03-24T11:22:13.071+0800    WARN    noop/main.go:16 this is a warning log with string fmt
2023-03-24T11:22:13.071+0800    ERROR   noop/main.go:17 this is an error level log with string fmt
main.main
        /data/github/noop/main.go:17
runtime.main
        /usr/local/go/src/runtime/proc.go:250
2023-03-24T11:22:13.071+0800    INFO    noop/main.go:18 this is an info level log with string fmt
# with debug level
[root@linux noop]# go run main.go 
2023-03-24T11:22:22.449+0800    INFO    noop/main.go:12 zap logger initialized
2023-03-24T11:22:22.449+0800    DEBUG   noop/main.go:12 zap logger debug level enabled
2023-03-24T11:22:22.449+0800    DEBUG   noop/main.go:16 this is a simple debugging log
2023-03-24T11:22:22.449+0800    WARN    noop/main.go:17 this is a warning log with string fmt
2023-03-24T11:22:22.449+0800    ERROR   noop/main.go:18 this is an error level log with string fmt
main.main
        /data/github/noop/main.go:18
runtime.main
        /usr/local/go/src/runtime/proc.go:250
2023-03-24T11:22:22.449+0800    INFO    noop/main.go:19 this is an info level log with string fmt

# with fields
[root@linux noop]# go run main.go 
2023-05-18T16:23:23.806+0800    DEBUG   noop/main.go:35 this is a simple debugging log  {"xx_key": "filed_value", "LocalTime": "2023-05-18 16:23:23.806"}
2023-05-18T16:23:23.807+0800    WARN    noop/main.go:36 this is a warning log with string fmt   {"xx_key": "filed_value", "LocalTime": "2023-05-18 16:23:23.806"}
2023-05-18T16:23:23.807+0800    ERROR   noop/main.go:37 this is an error level log with string fmt      {"xx_key": "filed_value", "LocalTime": "2023-05-18 16:23:23.806"}
main.main
        /Users/yao/GoLangWorkspace/MyselfProjects/github/noop/main.go:37
runtime.main
        /opt/homebrew/Cellar/go/1.20.3/libexec/src/runtime/proc.go:250
2023-05-18T16:23:23.807+0800    INFO    noop/main.go:38 this is an info level log with string fmt       {"xx_key": "filed_value", "LocalTime": "2023-05-18 16:23:23.806"}

```
