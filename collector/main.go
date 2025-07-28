package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var addr = flag.String("addr", "localhost:8888", "http service address")
var ingest = flag.String("ingest", "localhost:8889", "ingest service address")
var source = flag.String("source", "collector1", "name of collector")
var watchPath = flag.String("out", "/data/out", "core folder where collectors move files saved by plugin for sending")
var tmpPath = flag.String("tmp", "/data/tmp", "root folder where collector instructs plugin to store data")

var mu sync.Mutex

func main() {
	flag.Parse()

	if err := os.MkdirAll(*tmpPath, 0755); err != nil {
		log.Println(fmt.Errorf("failed to create directory: %w", err))
	}

	if err := os.MkdirAll(*watchPath, 0755); err != nil {
		log.Println(fmt.Errorf("failed to create directory: %w", err))
	}

	done := make(chan struct{})
	interrupt := make(chan os.Signal, 1)

	//sends notifications on interrupt signals
	// this allows the program to gracefully shut down when it receives an interrupt signal
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// kick off uploader/watcher
	uploadQueue := NewUploadQueue(10) // 10 workers
	watcher := NewWatcher(*watchPath, uploadQueue)

	messageHandler := NewMessageHandler(*addr, done, watcher)
	go messageHandler.Start()

	// run the initial refresh in nonblocking fashion
	go func() {
		err := refresh()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go watcher.Start()

	// close when done, sends HB message and handles interrups gracefully
	// eventLoop(c, ticker, done, interrupt, uploadQueue, watcher)
	eventLoop(done, interrupt, uploadQueue, watcher, messageHandler)
}

// func eventLoop(c *websocket.Conn, ticker *time.Ticker, done chan struct{}, interrupt chan os.Signal) {
func eventLoop(done chan struct{}, interrupt chan os.Signal, uploadQueue *UploadQueue, watcher *Watcher, mh *MessageHandler) {
	for {
		select {
		case <-done:
			// TODO: this part needs to change as not to close when the channel is closed as done is only closed when readhandler fails
			// attempt to reconnect should be made
			return
		case <-interrupt:
			mh.Stop()
			watcher.Stop()
			uploadQueue.Stop()
			return
		}
	}
}
