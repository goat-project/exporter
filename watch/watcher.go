package watch

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
)

// Watcher file system notification using event channel.
type Watcher struct {
	Watcher   *fsnotify.Watcher
	EventChan chan fsnotify.Event
}

// Watch watches new files and directories. Files are added
// to event channel for processing and directories are
// added to Watcher for watching new files.
func (w Watcher) Watch() {
	// The watch does not support listing of added directories.
	// Check if the directory exists before the watching starts.
	// This issue should be fixed by developers of Watcher.
	// The watching does not detect error when no directory is added.

	for {
		select {
		case event, ok := <-w.Watcher.Events:
			if !ok {
				logrus.WithField("error", "not ok").Error("watcher is not set correctly")
				return
			}

			logrus.WithFields(logrus.Fields{"event": event.Name, "op": event.Op}).Debug("event")

			if event.Op&fsnotify.Create == fsnotify.Create {
				fi, err := os.Stat(event.Name)
				if err != nil {
					logrus.WithField("error", err).Error("error os state")
					continue
				}

				switch mode := fi.Mode(); {
				case mode.IsDir(): // add new directory to Watcher
					err = w.Watcher.Add(event.Name)
					if err != nil {
						logrus.WithField("error", err).Error("error add directory to watcher")
						continue
					}
					logrus.WithField("dir", event.Name).Debug("dir added to watcher")
				case mode.IsRegular():
					//w.EventChan <- event // add event to channel when file is created
				}
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				w.EventChan <- event // add event to channel when file is modified - writing was done
			}

		case err, ok := <-w.Watcher.Errors:
			if !ok {
				return
			}

			logrus.WithField("error", err).Error("watch error")
		}
	}
}

// AddRootWithSubDirs adds root and subdirectories recursively.
func (w Watcher) AddRootWithSubDirs(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			err = w.Watcher.Add(path)
			if err != nil {
				return err
			}

			logrus.WithField("dir", path).Debug("dir added to watcher")
		}

		return nil
	})
}
