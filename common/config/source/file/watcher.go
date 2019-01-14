package file

import (
	"teddy-backend/common/config"

	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	config.WatchResult
	fw *fsnotify.Watcher
}

func newWatcher(f *file, path []string) (*config.WatchResult, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw.Add(f.path)

	w := &watcher{
		fw: fw,
	}

	watchResult := &config.WatchResult{
		WatcherStopper: w,
		Results:        make(chan map[string]interface{}),
		Errors:         make(chan error),
	}

	go func() {
		defer func() {
			close(watchResult.Results)
			close(watchResult.Errors)
		}()
		for {
			select {
			case _, more := <-w.fw.Events:
				if !more {
					break
				}
				c, err := f.Read()
				if err != nil {
					watchResult.Errors <- err
					break
				}
				watchResult.Results <- c
			case err, more := <-w.fw.Errors:
				if !more {
					break
				}
				watchResult.Errors <- err
				break
			}
		}
		w.fw.Close()
	}()

	return watchResult, nil
}

func (w *watcher) Stop() error {
	return w.fw.Close()
}
