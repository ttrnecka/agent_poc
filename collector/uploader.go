package main

import (
	"bytes"
	"fmt"
	"io"
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
		fmt.Println("Queue full. Dropping:", filePath)
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
			return
		}
		fmt.Printf("Retry %d for %s: %v\n", attempt, filePath, err)
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	fmt.Println("Failed to upload after retries:", filePath)
}

func uploadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	ingestURL := fmt.Sprintf("http://%s/upload", *ingest)

	req, err := http.NewRequest("POST", ingestURL, bytes.NewReader(content))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/octet-stream")
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

	fmt.Println("Uploaded:", path)
	return nil
}
