package common

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// zip all files in the directory without recursive, and deletes the original files
func ZipDirFlatAndDelete(srcFolder, zipFilePath string) error {
	absPath, err := filepath.Abs(srcFolder)
	if err != nil {
		return err
	}

	zipFileName := filepath.Base(zipFilePath)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Skip directories, no recursion here
			continue
		}

		if entry.Name() == zipFileName {
			// Skip the ZIP file itself to avoid recursion
			continue
		}

		filePath := filepath.Join(absPath, entry.Name())

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", filePath, err)
		}

		info, err := file.Stat()
		if err != nil {
			file.Close()
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			file.Close()
			return err
		}
		header.Name = entry.Name() // top-level file name only
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			file.Close()
			return err
		}

		if _, err := io.Copy(writer, file); err != nil {
			file.Close()
			return err
		}
		file.Close()

		// Delete file after adding it to ZIP
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to delete file %s: %w", filePath, err)
		}
	}
	return nil
}

// unzip a zip file to temp director and returns the directory
func UnzipToTemp(zipFilePath string) (string, error) {
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
