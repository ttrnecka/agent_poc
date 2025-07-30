package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger zerolog.Logger

func setupLogger() {
	// Configure lumberjack for log rotation
	rotatingFile := &lumberjack.Logger{
		Filename:   "./logs/collector.log",
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     28,   // days
		Compress:   true, // gzip
	}

	// Combine stdout and file output

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}

	multiWriter := zerolog.MultiLevelWriter(consoleWriter, rotatingFile)

	// Set global time format
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Initialize zerolog with desired log level
	logger = zerolog.New(multiWriter).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
