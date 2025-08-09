package main

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/ttrnecka/agent_poc/ingestion/parsers"
)

type Pipeline struct {
	logger zerolog.Logger
}

func (p Pipeline) Ingest(filePath string) {
	logger.Info().Msgf("Unzipping %s", filePath)
	tmpDir, err := unzipToTemp(filePath)
	if err != nil {
		logger.Error().Err(err).Msgf("Cannot unzip file %s", filePath)
	}
	logger.Info().Msgf("Going to read files in %s", tmpDir)
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		logger.Error().Err(err).Str("folder", tmpDir).Msg("Failed to read source folder")
	}
	for _, entry := range entries {
		if entry.IsDir() {
			// Skip subdirectories (you can recurse if needed)
			continue
		}
		srcPath := filepath.Join(tmpDir, entry.Name())
		p.Process(srcPath)
	}
}

func (p Pipeline) Process(file_path string) {

	p.logger = logger.With().
		Str("file_path", file_path).
		Logger()

	success := true
	defer func() {
		p.PostProcess(file_path, success)
	}()

	p.logger.Info().Msgf("Data Pipeline process started for %s", file_path)
	headers, body, err := p.parseFile(file_path)
	if err != nil {
		success = false
		p.logger.Error().Err(err).Msg("Cannot parse file")
		return
	}
	p.logger.Info().Msg("File headers and body read successfully")

	err = p.saveToDb(headers, body)
	if err != nil {
		success = false
		p.logger.Error().Err(err).Msg("Cannot save file to DB")
		return
	}
	p.logger.Info().Msg("File processed and saved to DB")
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
	policy := headers["policy"]

	if collector == "" || device == "" || endpoint == "" || policy == "" {
		err := fmt.Errorf("missing required headers: collector, device, policy, or endpoint")
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

	p.logger.Info().Msgf("Parsing body using policy %s", policy)
	parsed_json, err := p.parseBody(body, policy, endpoint)
	if err != nil {
		err = fmt.Errorf("failed to parse body: %w", err)
		return err
	}

	existingData := map[string]any{}
	filePath = filepath.Join(dirPath, "object")

	lockManager.Lock(filePath)
	p.logger.Info().Msgf("Lock acquired on %s", filePath)
	defer func() {
		lockManager.Unlock(filePath)
		p.logger.Info().Msgf("Lock released on %s", filePath)
	}()

	err = readExistingJson(filePath, existingData)
	if err != nil {
		err = fmt.Errorf("failed to read existing object: %w", err)
		return err
	}

	merged := mergeMaps(existingData, parsed_json)

	err = saveJson(filePath, merged)
	if err != nil {
		err = fmt.Errorf("failed to save json object: %w", err)
		return err
	}

	return nil
}

func (p Pipeline) parseBody(body, policy, endpoint string) (map[string]any, error) {
	result := make(map[string]any)

	if parsers.Parsers[policy] == nil {
		err := fmt.Errorf("no parser found for %s policy", policy)
		return nil, err
	}

	var extractors map[string]parsers.ExtractorConfig
	subkeys := make([]string, 0)
	// check endpoint parser matching name exactly
	if parsers.Parsers[policy].Extractors[endpoint] == nil {
		p.logger.Info().Msgf("No direct parser found for %s policy, endpoint %s, trying mappers", policy, endpoint)

		matched := false
		for key, mapper := range parsers.Parsers[policy].Mappers {
			p.logger.Debug().Msgf("Trying mapper %+v for %s", mapper, endpoint)
			if mapper.Pattern != "" {
				re := regexp.MustCompile(mapper.Pattern)
				names := re.SubexpNames()
				matches := re.FindStringSubmatch(endpoint)
				if len(matches) > 0 {
					if parsers.Parsers[policy].Extractors[key] != nil {
						p.logger.Info().Msgf("Mapper matched: using mapped endpoint %s", key)
						extractors = parsers.Parsers[policy].Extractors[key]
						matched = true

						// check if we capture groups or named capture groups
						if len(matches) > 1 {
							if len(names) > 1 {
								subkeys = append(subkeys, names[1])
							}
							subkeys = append(subkeys, matches[1])
						}
						break
					} else {
						p.logger.Error().Msgf("Mapper matched: but no endpoint %s exist in configuration", key)
					}
				}
			}
		}
		if !matched {
			extractors = make(map[string]parsers.ExtractorConfig)
		}
	} else {
		extractors = parsers.Parsers[policy].Extractors[endpoint]
	}

	for key, extractor := range extractors {
		fn, ok := parsers.Extractors[extractor.Method]
		if !ok {
			err := fmt.Errorf("no extractor for method %q", extractor.Method)
			return nil, err
		}
		val, err := fn(body, extractor)
		if err != nil {
			err := fmt.Errorf("error extracting %q: %v", key, err)
			return nil, err
		}
		if val != nil {
			p.logger.Debug().Msgf("Parsed value: %+v", val)
			tmpResult := result
			for _, subkey := range subkeys {
				if tmpResult[subkey] == nil {
					tmpResult[subkey] = make(map[string]any)
				}
				tmpResult = tmpResult[subkey].(map[string]any)
			}
			tmpResult[key] = val
		}
	}
	p.logger.Debug().Msgf("result: %+v", result)
	return result, nil
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

	// fileName := filepath.Base(srcPath)

	// destPath := filepath.Join(destDir, fileName)

	// err := os.Rename(srcPath, destPath)
	err := moveFile(srcPath, destDir)
	if err != nil {
		p.logger.Error().Err(err).Msg("Move failed")
		return
	}
	p.logger.Info().Msgf("Move succeeded")
}

// Merge two map[string]any, new overrides old
func mergeMaps(dst, src map[string]any) map[string]any {
	for k, v := range src {
		if vMap, ok := v.(map[string]any); ok {
			if dvMap, ok := dst[k].(map[string]any); ok {
				dst[k] = mergeMaps(dvMap, vMap)
			} else {
				dst[k] = vMap
			}
		} else {
			dst[k] = v
		}
	}
	return dst
}

func readExistingJson(filePath string, data map[string]any) error {
	if file, err := os.Open(filePath); err == nil {
		defer file.Close()
		byteValue, _ := io.ReadAll(file)
		json.Unmarshal(byteValue, &data)
	} else if !os.IsNotExist(err) {
		err = fmt.Errorf("error opening json file: %w", err)
		return err
	}
	return nil
}

func saveJson(filePath string, data map[string]any) error {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		err = fmt.Errorf("cannot marschal json: %w", err)
		return err
	}

	// Save body to file
	if err := os.WriteFile(filePath, out, 0644); err != nil {
		err = fmt.Errorf("failed to write parsed json to file: %w", err)
		return err
	}
	return nil
}

func unzipToTemp(zipFilePath string) (string, error) {
	// Open the zip archive for reading
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	// Create a temp directory to extract files into
	tempDir, err := os.MkdirTemp("", "unzipped-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Iterate through each file in the archive
	for _, f := range r.File {
		fPath := filepath.Join(tempDir, f.Name)

		// Get relative path from tempDir to fPath
		rel, err := filepath.Rel(tempDir, fPath)
		if err != nil {
			return "", fmt.Errorf("failed to get relative path: %w", err)
		}

		// If rel starts with ".." then fPath is outside tempDir => potential ZipSlip
		if strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
			return "", fmt.Errorf("illegal file path: %s", fPath)
		}

		if f.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(fPath, os.ModePerm); err != nil {
				return "", err
			}
			continue
		}

		// Open the file inside the zip
		srcFile, err := f.Open()
		if err != nil {
			return "", err
		}

		// Create destination file
		destFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			srcFile.Close()
			return "", err
		}

		// Copy contents
		_, err = io.Copy(destFile, srcFile)

		// Close files
		srcFile.Close()
		destFile.Close()

		if err != nil {
			return "", err
		}
	}

	return tempDir, nil
}

func moveFile(srcPath, dstDir string) error {
	// Make sure destination directory exists
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	dstPath := filepath.Join(dstDir, filepath.Base(srcPath))

	// Try to rename first (fast path)
	if err := os.Rename(srcPath, dstPath); err == nil {
		return nil
	}

	// If rename fails, fallback to copy + delete
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file data: %w", err)
	}

	// Flush to disk
	if err := dstFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	// Close files before deletion
	srcFile.Close()
	dstFile.Close()

	// Delete original file
	if err := os.Remove(srcPath); err != nil {
		return fmt.Errorf("failed to delete source file after copy: %w", err)
	}

	return nil
}
