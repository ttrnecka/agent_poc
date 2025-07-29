package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var addr = flag.String("addr", "localhost:8888", "http service address")
var ingest = flag.String("ingest", "localhost:8889", "ingest service address")
var source = flag.String("source", "collector1", "name of collector")
var watchPath = flag.String("out", "/data/out", "core folder where collectors move files saved by plugin for sending")
var tmpPath = flag.String("tmp", "/data/tmp", "root folder where collector instructs plugin to store data")

func main() {
	flag.Parse()

	if err := os.MkdirAll(*tmpPath, 0755); err != nil {
		log.Println(fmt.Errorf("failed to create directory: %w", err))
	}

	if err := os.MkdirAll(*watchPath, 0755); err != nil {
		log.Println(fmt.Errorf("failed to create directory: %w", err))
	}

	interrupt := make(chan os.Signal, 1)

	//sends notifications on interrupt signals
	// this allows the program to gracefully shut down when it receives an interrupt signal
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// kick off uploader/watcher
	uploadQueue := NewUploadQueue(10) // 10 workers
	watcher := NewWatcher(*watchPath, uploadQueue)

	// run the initial refresh in nonblocking fashion
	go func() {
		err := refresh()
		if err != nil {
			log.Println(err)
		}
	}()

	go watcher.Start()

	reconnectDelay := 5 // seconds
	for {
		log.Println("Starting MessageHandler and event loop...")
		messageHandler := NewMessageHandler(*addr, *source, watcher)
		messageHandler.Start()

		shouldExit := eventLoop(interrupt, uploadQueue, watcher, messageHandler)
		if shouldExit {
			log.Println("Shutting down main loop due to interrupt signal.")
			break
		}
		log.Printf("WebSocket connection lost, retrying in %d seconds...", reconnectDelay)
		time.Sleep(time.Duration(reconnectDelay) * time.Second)
	}
}

// eventLoop returns true if process should exit (interrupt), false if should reconnect
func eventLoop(interrupt chan os.Signal, uploadQueue *UploadQueue, watcher *Watcher, mh *MessageHandler) bool {
	for {
		select {
		case <-mh.done:
			log.Println("WebSocket connection closed or failed, will attempt to reconnect.")
			return false // signal to reconnect
		case <-interrupt:
			log.Println("Received interrupt signal")
			mh.Stop()
			watcher.Stop()
			uploadQueue.Stop()
			return true // signal to exit
		}
	}
}
