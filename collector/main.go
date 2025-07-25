package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ttrnecka/agent_poc/ws"
)

var addr = flag.String("addr", "localhost:8888", "http service address")
var ingest = flag.String("ingest", "localhost:8889", "ingest service address")
var source = flag.String("source", "collector1", "name of collector")
var watchPath = flag.String("out", "/data/out", "core folder where collectors move files saved by plugin for sending")
var tmpPath = flag.String("tmp", "/data/tmp", "root folder where collector instructs plugin to store data")

var mu sync.Mutex

func connectWS(host string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: host, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	return c, err
}

func main() {
	flag.Parse()

	if err := os.MkdirAll(*tmpPath, 0755); err != nil {
		log.Println(fmt.Errorf("failed to create directory: %w", err))
	}

	if err := os.MkdirAll(*watchPath, 0755); err != nil {
		log.Println(fmt.Errorf("failed to create directory: %w", err))
	}

	done := make(chan struct{})
	interrupt := make(chan os.Signal, 1)

	//sends notifications on interrupt signals
	// this allows the program to gracefully shut down when it receives an interrupt signal
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	//connect to the websocket server
	c, err := connectWS(*addr)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// set up go routines to read from the websocket connection and handle messages
	go func() {
		messageReader(c, done)
	}()

	go func() {
		messageHandler()
	}()

	// run the initial refresh in nonblocking fashion
	go func() {
		err = refresh()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// kick off uploader/watcher
	// uploadQueue := NewUploadQueue(10) // 10 workers
	// watcher := NewRecursiveWatcher(*watchPath, uploadQueue)

	// go watcher.Start()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// close when done, sends HB message and handles interrups gracefully
	// eventLoop(c, ticker, done, interrupt, uploadQueue, watcher)
	eventLoop(c, ticker, done, interrupt)
}

func sendHeartbeat(c *websocket.Conn) error {
	mu.Lock()
	err := c.WriteJSON(ws.NewMessage(ws.MSG_ONLINE, *source, "hub", "Collector is online"))
	mu.Unlock()
	if err != nil {
		log.Println("heartbeat error:", err)
		return err
	}
	return nil
}

func cleanShutdown(c *websocket.Conn, done chan struct{}) {
	log.Println("interrupt")

	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	// err := c.WriteMessage(websocket.TextMessage, []byte("OFFLINE"))
	mu.Lock()
	err := c.WriteJSON(ws.NewMessage(ws.MSG_OFFLINE, *source, "hub", "Collector is going offline"))
	if err != nil {
		log.Println("write:", err)
	}
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	mu.Unlock()
	if err != nil {
		log.Println("write close:", err)
	}
	select {
	case <-done:
	case <-time.After(time.Second):
	}
}

func eventLoop(c *websocket.Conn, ticker *time.Ticker, done chan struct{}, interrupt chan os.Signal) {
	// func eventLoop(c *websocket.Conn, ticker *time.Ticker, done chan struct{}, interrupt chan os.Signal, uploadQueue *UploadQueue, watcher *RecursiveWatcher) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err := sendHeartbeat(c)
			if err != nil {
				return
			}
		case <-interrupt:
			cleanShutdown(c, done)
			// watcher.Stop()
			// uploadQueue.Stop()
			return
		}
	}
}
