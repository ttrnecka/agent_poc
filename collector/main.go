package main

import (
	"flag"
	"fmt"
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

	setupLogger()

	if err := os.MkdirAll(*tmpPath, 0755); err != nil {
		logger.Fatal().Err(fmt.Errorf("failed to create directory: %w", err)).Msg("")
	}

	if err := os.MkdirAll(*watchPath, 0755); err != nil {
		logger.Fatal().Err(fmt.Errorf("failed to create directory: %w", err)).Msg("")
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
			logger.Error().Err(err).Msg("")
		}
	}()

	go watcher.Start()

	reconnectDelay := 5 // seconds
	for {
		logger.Info().Msg("Starting MessageHandler and event loop...")
		messageHandler := NewMessageHandler(*addr, *source, watcher)
		messageHandler.Start()

		shouldExit := eventLoop(interrupt, uploadQueue, watcher, messageHandler)
		if shouldExit {
			logger.Info().Msg("Shutting down main loop due to interrupt signal.")
			break
		}
		logger.Info().Msg(fmt.Sprintf("WebSocket connection lost, retrying in %d seconds...", reconnectDelay))
		time.Sleep(time.Duration(reconnectDelay) * time.Second)
	}
}

// eventLoop returns true if process should exit (interrupt), false if should reconnect
func eventLoop(interrupt chan os.Signal, uploadQueue *UploadQueue, watcher *Watcher, mh *MessageHandler) bool {
	for {
		select {
		case <-mh.done:
			logger.Info().Msg("WebSocket connection closed or failed, will attempt to reconnect.")
			return false // signal to reconnect
		case <-interrupt:
			logger.Info().Msg("Received interrupt signal")
			mh.Stop()
			watcher.Stop()
			uploadQueue.Stop()
			return true // signal to exit
		}
	}
}
