package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v3"
)

var defaultLogLocation = "./logs/"
var defaultConfigLocation = "./config/loggers"

var loggers = make(map[string]zerolog.Logger)

type Config struct {
	Level      zerolog.Level `yaml:"level"`
	MaxSize    int           `yaml:"max_size"`
	MaxBackups int           `yaml:"max_backups"`
	MaxAge     int           `yaml:"max_age"`
	Compress   bool          `yaml:"compress"`
}

func GetLogLocation() string {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, defaultLogLocation)
}

func LogLocation(name string) {
	defaultLogLocation = name
}

func ConfigLocation(name string) {
	defaultConfigLocation = name
}

func DefaultConfig() Config {
	return Config{
		Level:      zerolog.DebugLevel,
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     28,   // days
		Compress:   true, // gzip
	}
}
func SetupLogger(name string) zerolog.Logger {
	os.Setenv("NO_COLOR", "true")
	val, ok := loggers[name]
	if ok {
		return val
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}
	log.Logger = log.Output(consoleWriter)

	var config Config

	configName := fmt.Sprintf("%s%s.yaml", defaultConfigLocation, name)
	logName := filepath.Join(GetLogLocation(), fmt.Sprintf("%s.log", name))
	if _, err := os.Stat(configName); err != nil {
		config = DefaultConfig()
		log.Info().Str("config", fmt.Sprintf("%+v", config)).Msg("Using default logging config")
	} else {
		f, err := os.ReadFile(configName)
		if err != nil {
			log.Fatal().Err(err).Str("file", configName).Msg("Cannot read log config file")
		}
		if err := yaml.Unmarshal(f, &config); err != nil {
			log.Fatal().Err(err).Str("file", configName).Msg("Cannot parse log config file")
		}
		log.Printf("Using custom logging config: %+v", config)
	}
	// Configure lumberjack for log rotation
	rotatingFile := &lumberjack.Logger{
		Filename:   logName,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,   // days
		Compress:   config.Compress, // gzip
	}

	// Combine stdout and file output

	multiWriter := zerolog.MultiLevelWriter(consoleWriter, rotatingFile)

	// Set global time format
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Initialize zerolog with desired log level
	logger := zerolog.New(multiWriter).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(config.Level)

	loggers[name] = logger

	return logger
}
