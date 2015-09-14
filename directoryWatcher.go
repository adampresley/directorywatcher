package directorywatcher

import (
	"os"
	"path/filepath"
	"time"

	"github.com/adampresley/logging"
)

/*
DirectoryWatcher provides a set of functions to watch for changes on a directory. It allows
you to specificy a base and a function to be called when a change on an file occurs.
*/
type DirectoryWatcher struct {
	BasePath  string
	PauseTime time.Duration
	Recurse   bool

	log       *logging.Logger
	startTime time.Time
}

/*
FileChangedEventFunc defines the function that is called when a change event occurs
*/
type FileChangedEventFunc func(path string, info os.FileInfo, startTime time.Time, modificationTime time.Time) error

/*
NewDirectoryWatcher creates a new instance of a directory watcher.
*/
func NewDirectoryWatcher(basePath string, log *logging.Logger) *DirectoryWatcher {
	return &DirectoryWatcher{
		BasePath:  basePath,
		PauseTime: 500,
		Recurse:   true,

		log:       log,
		startTime: time.Now(),
	}
}

/*
SetPauseTime determines how much time elapses (in milliseconds) between directory scans. The default is
500 milliseconds.
*/
func (watcher *DirectoryWatcher) SetPauseTime(pauseTime time.Duration) {
	watcher.PauseTime = pauseTime
}

/*
SetRecurse determines if the watcher will recurse down subdirectories. The default is true.
*/
func (watcher *DirectoryWatcher) SetRecurse(recurse bool) {
	watcher.Recurse = recurse
}

/*
Watch initiates a watch goroutine on the base directory. The function passed in will be called when
a file is changed.
*/
func (watcher *DirectoryWatcher) Watch(fileChangedEvent FileChangedEventFunc) {
	watcher.log.Info("Listening for directory changes on", watcher.BasePath)

	/*
	 * Walk the directory tree and scan for changes
	 */
	go func() {
		for {
			filepath.Walk(watcher.BasePath, func(path string, info os.FileInfo, err error) error {
				if !watcher.Recurse {
					if filepath.Dir(path) != "." {
						return filepath.SkipDir
					}
				}

				if info.ModTime().After(watcher.startTime) {
					result := fileChangedEvent(path, info, watcher.startTime, info.ModTime())
					watcher.startTime = time.Now()
					return result
				}

				return nil
			})

			time.Sleep(watcher.PauseTime * time.Millisecond)
		}
	}()
}
