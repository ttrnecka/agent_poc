package main

import (
	"flag"
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
var source = flag.String("source", "collector1", "name of collector")

var mu sync.Mutex

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			mes := ws.NewMessage(0, "", "", "")

			err := c.ReadJSON(&mes)
			// _, mes, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			messages <- Message{Message: mes, c: c}
		}
	}()

	go func() {
		messageHandler()
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	err = refresh()
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// err := c.WriteMessage(websocket.TextMessage, []byte("ONLINE"))
			mu.Lock()
			err := c.WriteJSON(ws.NewMessage(ws.MSG_ONLINE, *source, "hub", "Collector is online"))
			mu.Unlock()
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			// err := c.WriteMessage(websocket.TextMessage, []byte("OFFLINE"))
			mu.Lock()
			err := c.WriteJSON(ws.NewMessage(ws.MSG_OFFLINE, *source, "hub", "Collector is going offline"))
			mu.Unlock()
			if err != nil {
				log.Println("write:", err)
				return
			}
			mu.Lock()
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			mu.Unlock()
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
