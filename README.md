# noop
**No** **o**ption za**p** wrapper

## Example
```go
package main

import (
	"github.com/patsnapops/noop/log"
)

func main() {
	log.Default().WithFilename("app.log").Init()
	// or just
	// log.Default().Init()

	log.Debug("this is a simple debugging log")
	log.Infof("this is an info level log with string %s", "fmt")
}
```
