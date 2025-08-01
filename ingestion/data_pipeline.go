package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
)

type Pipeline struct {
	logger zerolog.Logger
}

func (p Pipeline) Process(file_path string) {

	p.logger = logger.With().
		Str("file_path", file_path).
		Logger()

	success := true
	go func() {
		p.PostProcess(file_path, success)
	}()

	p.logger.Info().Msgf("Data Pipeline process started for %s", file_path)
	headers, body, err := p.parseFile(file_path)
	if err != nil {
		success = false
		p.logger.Error().Err(err).Msg("Cannot parse file")
		return
	}

	err = p.saveToDb(headers, body)
	if err != nil {
		success = false
		p.logger.Error().Err(err).Msg("Cannot save file to DB")
		return
	}
}

func isHeaderLine(line string) bool {
	return strings.HasPrefix(line, "---") &&
		(len(line) == 3 || (len(line) > 3 && line[3] != '-'))
}

func (d Pipeline) parseFile(file_path string) (headers map[string]string, body string, err error) {
	file, err := os.Open(file_path)
	if err != nil {
		logger.Error().Err(err).Msgf("Cannot open file %s: %s", file_path, err)
		return
	}
	defer file.Close()

	var (
		bodyBuilder strings.Builder
		inBody      bool
	)

	headers = make(map[string]string)

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

	if err = scanner.Err(); err != nil {
		logger.Error().Err(err).Msg("Error reading file")
		return
	}
	body = bodyBuilder.String()
	return
}

func (p Pipeline) saveToDb(headers map[string]string, body string) error {
	db_path := "/data/db"

	// Required fields
	collector := headers["collector"]
	device := headers["device"]
	endpoint := headers["endpoint"]

	if collector == "" || device == "" || endpoint == "" {
		err := fmt.Errorf("missing required headers: collector, device, probe_id, or endpoint")
		return err
	}

	dirPath := filepath.Join(db_path, collector, device)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		err = fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	// Save body to file
	filePath := filepath.Join(dirPath, endpoint)
	if err := os.WriteFile(filePath, []byte(body), 0644); err != nil {
		err = fmt.Errorf("failed to write body to file: %w", err)
		return err
	}
	return nil
}

func (p Pipeline) PostProcess(srcPath string, success bool) {

	destDir := config.processedDir
	msg := "Data Pipeline process finished succesfully."
	if !success {
		destDir = config.failedDir
		msg = "Data Pipeline process failed."
	}
	msg = fmt.Sprintf("%s Moving file to %s", msg, destDir)
	p.logger.Info().Msg(msg)

	fileName := filepath.Base(srcPath)

	destPath := filepath.Join(destDir, fileName)

	err := os.Rename(srcPath, destPath)
	if err != nil {
		p.logger.Error().Err(err).Msg("Move failed")
		return
	}
	p.logger.Info().Msgf("Move succeeded")
}
