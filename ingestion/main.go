package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	logging "github.com/ttrnecka/agent_poc/logger"
)

var logger zerolog.Logger

var config Config

type Config struct {
	uploadDir    string
	processedDir string
	failedDir    string
}

func init() {
	logger = logging.SetupLogger("ingesting_service")
	config = Config{
		uploadDir:    "/data/uploads",
		processedDir: "/data/processed",
		failedDir:    "/data/failed",
	}
}

func main() {

	// Ensure the directories exist
	for _, dir := range []string{config.uploadDir, config.processedDir, config.failedDir} {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Fatal().Err(err).Msgf("Failed to create directory %s", dir)
		}
	}

	srv := &http.Server{
		Addr:    ":8888",
		Handler: router(),
	}

	logger.Info().Msg("Starting ingestion service")
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start ingestion service")
	}
}

func router() http.Handler {
	mux := http.NewServeMux()

	// index page
	mux.HandleFunc("/", indexHandler)

	// upload page
	mux.HandleFunc("/upload", handleUpload)

	return mux
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.Header.Get("File-Name")
	if fileName == "" {
		http.Error(w, "Missing File-Name header", http.StatusBadRequest)
		return
	}

	// Sanitize filename: prevent path traversal
	safeName := filepath.Base(fileName)
	if strings.Contains(safeName, "..") || strings.HasPrefix(safeName, "/") {
		http.Error(w, "Invalid file name", http.StatusBadRequest)
		return
	}

	outPath := filepath.Join(config.uploadDir, safeName)
	logger.Info().Msgf("Creating file %s", outPath)
	outFile, err := os.Create(outPath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to save file")
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, r.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to write file")
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	go Pipeline{}.Process(outPath)

	logger.Info().Msgf("Received file: %s -> %s", fileName, outPath)
	w.WriteHeader(http.StatusOK)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
