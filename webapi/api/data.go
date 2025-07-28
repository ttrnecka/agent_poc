package api

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const baseDir = "/data/db"

// Utility: JSON response writer
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Handler entry point
func DataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	trimmedPath := strings.TrimPrefix(r.URL.Path, "/api/v1/data/collector")
	trimmedPath = strings.Trim(trimmedPath, "/")

	var parts []string
	if trimmedPath != "" {
		parts = strings.Split(trimmedPath, "/")
	}

	switch len(parts) {
	case 0:
		listCollectors(w, r)
	case 1:
		getCollector(w, r, parts[0])
	case 2:
		getDevice(w, r, parts[0], parts[1])
	case 3:
		getEndpoint(w, r, parts[0], parts[1], parts[2])
	default:
		http.NotFound(w, r)
	}
}

// /api/v1/data/collector
func listCollectors(w http.ResponseWriter, _ *http.Request) {
	dirs, err := os.ReadDir(baseDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var collectors []string
	for _, d := range dirs {
		if d.IsDir() {
			collectors = append(collectors, d.Name())
		}
	}
	writeJSON(w, collectors)
}

// /api/v1/data/collector/:collector
func getCollector(w http.ResponseWriter, _ *http.Request, collector string) {
	collectorPath := filepath.Join(baseDir, collector)
	dirs, err := os.ReadDir(collectorPath)
	if err != nil {
		http.Error(w, "Collector not found", http.StatusNotFound)
		return
	}
	var devices []string
	for _, d := range dirs {
		if d.IsDir() {
			devices = append(devices, d.Name())
		}
	}
	writeJSON(w, map[string]interface{}{"devices": devices})
}

// /api/v1/data/collector/:collector/:device
func getDevice(w http.ResponseWriter, _ *http.Request, collector, device string) {
	devicePath := filepath.Join(baseDir, collector, device)

	entries, err := os.ReadDir(devicePath)
	if err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	var endpoints []string
	for _, entry := range entries {
		if !entry.IsDir() {
			endpoints = append(endpoints, entry.Name())
		}
	}

	writeJSON(w, map[string][]string{"endpoints": endpoints})
}

// /api/v1/data/collector/:collector/:device/:endpoint
func getEndpoint(w http.ResponseWriter, _ *http.Request, collector, device, endpoint string) {
	filePath := filepath.Join(baseDir, collector, device, endpoint)

	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Endpoint not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
