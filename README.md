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

	// set fields, which will be printed in the log, the fields is a map,
	// key is the field name, value is the field value, the value can be any type, example:
	fields := map[string]interface{}{
		"instance_id": "1135286d-2fa7-4715-8b90-4937c0e49c2d",
		"localtime":   time.Now().In(time.Local).Format("2006-01-02 15:04:05.000"),
	}
	log.Default().WithFilename("app.log").WithLevel(log.DebugLevel).WithFields(fields).Init()

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
2023-05-19T00:36:25.652+0800    DEBUG   noop/main.go:27 this is a simple debugging log  {"instance_id": "1135286d-2fa7-4715-8b90-4937c0e49c2d", "localtime": "2023-05-19 00:36:25.652"}
2023-05-19T00:36:25.652+0800    WARN    noop/main.go:28 this is a warning log with string fmt   {"instance_id": "1135286d-2fa7-4715-8b90-4937c0e49c2d", "localtime": "2023-05-19 00:36:25.652"}
2023-05-19T00:36:25.652+0800    ERROR   noop/main.go:29 this is an error level log with string fmt      {"instance_id": "1135286d-2fa7-4715-8b90-4937c0e49c2d", "localtime": "2023-05-19 00:36:25.652"}
main.main
        /Users/longyao/GolangWorkspace/MyselfProjects/noop/main.go:29
runtime.main
        /opt/homebrew/opt/go/libexec/src/runtime/proc.go:250
2023-05-19T00:36:25.652+0800    INFO    noop/main.go:30 this is an info level log with string fmt       {"instance_id": "1135286d-2fa7-4715-8b90-4937c0e49c2d", "localtime": "2023-05-19 00:36:25.652"}
```
