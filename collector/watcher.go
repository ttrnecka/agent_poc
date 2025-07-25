package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type RecursiveWatcher struct {
	watcher     *fsnotify.Watcher
	watchRoot   string
	queue       *UploadQueue
	debounceMap map[string]time.Time
	mutex       sync.Mutex
	done        chan struct{}
}

func NewRecursiveWatcher(path string, queue *UploadQueue) *RecursiveWatcher {
	w, _ := fsnotify.NewWatcher()
	return &RecursiveWatcher{
		watcher:     w,
		watchRoot:   path,
		queue:       queue,
		debounceMap: make(map[string]time.Time),
		done:        make(chan struct{}),
	}
}

func (rw *RecursiveWatcher) Start() {
	_ = filepath.WalkDir(rw.watchRoot, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			_ = rw.watcher.Add(path)
		}
		return nil
	})

	go rw.eventLoop()
}

func (rw *RecursiveWatcher) Stop() {
	close(rw.done)
	rw.watcher.Close()
}

func (rw *RecursiveWatcher) eventLoop() {
	for {
		select {
		case event := <-rw.watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create {
				rw.handleCreateEvent(event.Name)
			}
		case err := <-rw.watcher.Errors:
			fmt.Println("Watcher error:", err)
		case <-rw.done:
			return
		}
	}
}

func (rw *RecursiveWatcher) handleCreateEvent(path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if info.IsDir() {
		_ = rw.watcher.Add(path)
	} else {
		rw.mutex.Lock()
		lastSeen := rw.debounceMap[path]
		if time.Since(lastSeen) < 2*time.Second {
			rw.mutex.Unlock()
			return
		}
		rw.debounceMap[path] = time.Now()
		rw.mutex.Unlock()

		go func() {
			// Wait a bit to ensure file is fully written
			time.Sleep(2 * time.Second)
			rw.queue.Enqueue(path)
		}()
	}
}
