package reloader

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
)

func Watch(dir string, load func(), delay time.Duration) error {
	return watch(fsnotify.NewWatcher, dir, load, delay)
}

func watch(fn func() (*fsnotify.Watcher, error), dir string, load func(), delay time.Duration) error {
	watcher, err := fn()
	if err != nil {
		return fmt.Errorf("unable to initialize file system notifications: %v", err)
	}
	if err := watcher.Add(dir); err != nil {
		return fmt.Errorf("unable to watch directory: %v", err)
	}
	go func() {
		var cancel chan struct{}
		for {
			select {
			case e := <-watcher.Events:
				if e.Op != 0 {
					if cancel != nil {
						close(cancel)
					}
					cancel = make(chan struct{})
					go func(cancel chan struct{}) {
						select {
						case <-time.After(delay):
							load()
						case <-cancel:
						}
					}(cancel)
				}
			}
		}
	}()
	return nil
}
