package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Pipeline struct{}

func (p Pipeline) Process(file_path string) error {
	file, err := os.Open(file_path)
	if err != nil {
		log.Printf("Cannot open file %s: %s", file_path, err)
		return err
	}
	defer file.Close()

	headers := make(map[string]string)

	var (
		// headerBuilder strings.Builder
		bodyBuilder strings.Builder
		inBody      bool
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !inBody && isHeaderLine(line) {
			trimmed := strings.TrimSpace(line[3:])
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				headers[key] = value
			}
		} else {
			inBody = true
			bodyBuilder.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file: %s", err)
		return err
	}

	// Required fields
	collector := headers["collector"]
	device := headers["device"]
	// probeID := headers["probe_id"]
	endpoint := headers["endpoint"]

	if collector == "" || device == "" || endpoint == "" {
		err = fmt.Errorf("missing required headers: collector, device, probe_id, or endpoint")
		log.Printf("Parsing error: %s", err)
		return err
	}

	dirPath := filepath.Join("/data/db", collector, device)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		err = fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		log.Print(err)
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	// Save body to file
	filePath := filepath.Join(dirPath, endpoint)
	if err := os.WriteFile(filePath, []byte(bodyBuilder.String()), 0644); err != nil {
		err = fmt.Errorf("failed to write body to file: %w", err)
		log.Print(err)
		return err
	}

	fmt.Printf("Saved body to %s\n", filePath)
	return nil
}

func isHeaderLine(line string) bool {
	return strings.HasPrefix(line, "---") &&
		(len(line) == 3 || (len(line) > 3 && line[3] != '-'))
}
