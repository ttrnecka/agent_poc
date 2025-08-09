package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type UploadQueue struct {
	queue      chan string
	wg         sync.WaitGroup
	workers    int
	stopSignal chan struct{}
}

func NewUploadQueue(workers int) *UploadQueue {
	q := &UploadQueue{
		queue:      make(chan string, 1000),
		workers:    workers,
		stopSignal: make(chan struct{}),
	}
	q.startWorkers()
	return q
}

func (q *UploadQueue) startWorkers() {
	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker()
	}
}

func (q *UploadQueue) Stop() {
	close(q.stopSignal)
	close(q.queue)
	q.wg.Wait()
}

func (q *UploadQueue) Enqueue(filePath string) {
	select {
	case q.queue <- filePath:
	default:
		logger.Error().Str("file", filePath).Msg("Queue full. Dropping file")
	}
}

func (q *UploadQueue) worker() {
	defer q.wg.Done()
	for {
		select {
		case filePath, ok := <-q.queue:
			if !ok {
				return
			}
			q.uploadWithRetries(filePath, 3)
		case <-q.stopSignal:
			return
		}
	}
}

func (q *UploadQueue) uploadWithRetries(filePath string, retries int) {
	for attempt := 1; attempt <= retries; attempt++ {
		err := uploadFile(filePath)
		if err == nil {
			deleteFile(filePath)
			return
		}
		logger.Error().Err(err).Str("file", filePath).Int("retry", attempt).Msg("Failed to upload file")
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	logger.Error().Str("file", filePath).Int("retries", retries).Msg("Failed to upload file after several retries")
}

func deleteFile(path string) {
	logger.Info().Str("filepath", path).Msg("Deleting file")
	err := os.Remove(path)
	if err != nil {
		logger.Error().Err(err).Str("filepath", path).Msg("Error deleting file")
		return
	}
}

func uploadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ingestURL := fmt.Sprintf("http://%s/upload", *ingest)

	req, err := http.NewRequest("POST", ingestURL, file)
	if err != nil {
		return err
	}
	mimeType := mime.TypeByExtension(filepath.Ext(path))

	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("File-Name", filepath.Base(path))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("bad response: %d - %s", resp.StatusCode, string(body))
	}

	logger.Info().Str("file", path).Msg("Uploaded file")
	return nil
}
