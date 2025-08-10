package main

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Watcher struct {
	watchRoot string
	queue     *UploadQueue
	process   chan struct{}
	done      chan struct{}
	mu        sync.Mutex
}

func NewWatcher(path string) *Watcher {

	uploadQueue := NewUploadQueue(10) // 10 workers

	return &Watcher{
		watchRoot: path,
		queue:     uploadQueue,
		process:   make(chan struct{}),
		done:      make(chan struct{}),
	}
}

func (rw *Watcher) Start() {
	go rw.eventLoop()
}

func (rw *Watcher) Stop() {
	close(rw.done)
	rw.queue.Stop()
}

func (rw *Watcher) Process() {
	rw.process <- struct{}{}
}

func (rw *Watcher) eventLoop() {

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		walk := false
		select {
		case <-rw.process:
			logger.Info().Msg("Adhoc queue processing")
			walk = true
		case <-ticker.C:
			logger.Info().Msg("Scheduled queue processing")
			walk = true
		case <-rw.done:
			return
		}
		if walk {
			rw.walkAndEnqueue()
		}
	}
}

func (rw *Watcher) walkAndEnqueue() {
	// Read all entries in the source directory
	// only run it one at a time
	rw.mu.Lock()
	defer rw.mu.Unlock()

	entries, err := os.ReadDir(rw.watchRoot)
	if err != nil {
		logger.Error().Err(err).Str("path", rw.watchRoot).Msg("Error accessing path")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Skip directories
			continue
		}

		srcPath := filepath.Join(rw.watchRoot, entry.Name())
		logger.Debug().Str("file", srcPath).Msg("Queuing for upload")
		rw.queue.Enqueue((srcPath))

	}
}
