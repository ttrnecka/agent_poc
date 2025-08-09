package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// moves file accross filesystems
func MoveFile(srcPath, dstDir string) error {
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
