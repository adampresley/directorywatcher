Directory Watcher
=================

Go Directory Change Watcher. Use like so.

```go
package main

import (
   "os"
   "time"

   "github.com/adampresley/directorywatcher"
   "github.com/adampresley/logging"
)

func main() {
   log := logging.NewLogger("My Logger")
   watcher := directorywatcher.NewDirectoryWatcher("./www", log)

   watcher.Watch(func(path string, info os.FileInfo, startTime time.Time, modificationTime time.Time) error {
      log.Info("A file changed:", path)
   })
}
```