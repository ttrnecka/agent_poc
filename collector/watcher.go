package main

import (
	"log"
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

func NewWatcher(path string, queue *UploadQueue) *Watcher {
	return &Watcher{
		watchRoot: path,
		queue:     queue,
		process:   make(chan struct{}),
		done:      make(chan struct{}),
	}
}

func (rw *Watcher) Start() {
	go rw.eventLoop()
}

func (rw *Watcher) Stop() {
	close(rw.done)
}

func (rw *Watcher) Process() {
	rw.process <- struct{}{}
}

func (rw *Watcher) eventLoop() {

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-rw.process:
			log.Printf("Adhoc queue processing")
			rw.walkAndEnqueue()
		case <-ticker.C:
			log.Printf("Scheduled queue processing")
			rw.walkAndEnqueue()
		case <-rw.done:
			return
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
		log.Printf("Error accessing path %q: %v\n", rw.watchRoot, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Skip directories
			continue
		}

		srcPath := filepath.Join(rw.watchRoot, entry.Name())
		rw.queue.Enqueue((srcPath))

	}
}
