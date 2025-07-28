package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var uploadDir string

func main() {
	// Parse CLI flag
	flag.StringVar(&uploadDir, "upload-dir", "/data/uploads", "Directory to save uploaded files")
	flag.Parse()

	// Ensure the upload directory exists
	err := os.MkdirAll(uploadDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create upload dir: %v", err)
	}

	srv := &http.Server{
		Addr:    ":8888",
		Handler: router(),
	}

	log.Println("Starting ingestion service")
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Println(err)
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

	outPath := filepath.Join(uploadDir, safeName)
	log.Printf("Creating file %s", outPath)
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, r.Body)
	if err != nil {
		log.Printf("Failed to write file: %v", err)
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	go Pipeline{}.Process(outPath)

	log.Printf("Received file: %s -> %s", fileName, outPath)
	w.WriteHeader(http.StatusOK)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
