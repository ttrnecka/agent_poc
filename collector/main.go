package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	logging "github.com/ttrnecka/agent_poc/logger"
)

var addr = flag.String("addr", "localhost:8888", "http service address")
var ingest = flag.String("ingest", "localhost:8889", "ingest service address")
var source = flag.String("source", "collector1", "name of collector")
var watchPath = flag.String("out", "/data/out", "core folder where collectors move files saved by plugin for sending")
var tmpPath = flag.String("tmp", "/data/tmp", "root folder where collector instructs plugin to store data")

var logger zerolog.Logger

func init() {
	logger = logging.SetupLogger("collector")
}

func main() {
	flag.Parse()

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

	watcher := NewWatcher(*watchPath)

	// run the initial refresh in nonblocking fashion
	go func() {
		err := refresh()
		if err != nil {
			logger.Error().Err(err).Msg("")
		}
	}()

	go watcher.Start()

	messageHandler := NewMessageHandler(*addr, *source, *watchPath)

	go func() {
		reconnectDelay := 5 // seconds
		for {
			logger.Info().Msg("Starting MessageHandler and event loop...")
			messageHandler.Start()

			for range messageHandler.done {
				logger.Info().Msg("WebSocket connection closed or failed, will attempt to reconnect.")
				break
			}
			logger.Info().Msg(fmt.Sprintf("Retrying in %d seconds...", reconnectDelay))
			time.Sleep(time.Duration(reconnectDelay) * time.Second)
		}
	}()

	<-interrupt
	logger.Info().Msg("Received interrupt signal")
	messageHandler.Stop()
}
