package main

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ttrnecka/agent_poc/webapi/ws"
)

type Message struct {
	ws.Message
	c *websocket.Conn
}

type MessageHandler struct {
	addr     string
	c        *websocket.Conn
	messages chan Message
	done     chan struct{}
	watcher  *Watcher
	ticker   *time.Ticker
}

func NewMessageHandler(addr string, done chan struct{}, watcher *Watcher) *MessageHandler {
	return &MessageHandler{
		addr:     addr,
		done:     done,
		messages: make(chan Message, 100),
		watcher:  watcher,
		ticker:   time.NewTicker(5 * time.Second),
	}
}

func (m *MessageHandler) Start() {
	//connect to the websocket server
	c, err := connectWS(m.addr)
	if err != nil {
		log.Fatal("dial:", err)
	}
	m.c = c
	defer c.Close()
	go m.readLoop()
	m.processLoop()
}

func (m *MessageHandler) Stop() {
	cleanShutdown(m.c, m.done)
}
func connectWS(host string) (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: host, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	return c, err
}

func (m *MessageHandler) readLoop() {
	defer close(m.done)

	for {
		mes := ws.NewMessage(0, "", "", "")

		err := m.c.ReadJSON(&mes)
		// _, mes, err := c.ReadMessage()
		if err != nil {
			log.Println("websocket read:", err)
			return
		}
		m.messages <- Message{Message: mes, c: m.c}
	}
}

func (m *MessageHandler) processLoop() {
	for {
		select {
		case <-m.ticker.C:
			err := sendHeartbeat(m.c)
			// TODO: this part needs to change
			// attempt attempt needs to be retried
			if err != nil {
				log.Printf("Failed to send heartbeat message: %v", err)
			}
		case msg, ok := <-m.messages:
			if !ok {
				return
			}
			// TODO needs to make this into pool of worker or as go routines
			log.Printf("websocket received: %v", msg)
			if msg.Destination == *source {
				if msg.Type == ws.MSG_REFRESH {
					refresh()
				}
				if msg.Type == ws.MSG_RUN {
					run(msg.Message, msg.c)
					log.Printf("Sending process message")
					m.watcher.Process()
					log.Printf("Sent process message")
				}
			}
		}
	}
}

func sendHeartbeat(c *websocket.Conn) error {
	mu.Lock()
	log.Println("Sending heartbeat message")
	err := c.WriteJSON(ws.NewMessage(ws.MSG_ONLINE, *source, "hub", "Collector is online"))
	mu.Unlock()
	if err != nil {
		log.Println("heartbeat error:", err)
		return err
	}
	log.Println("Heartbeat sent")
	return nil
}

func cleanShutdown(c *websocket.Conn, done chan struct{}) {
	log.Println("Interrupt received")

	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	// err := c.WriteMessage(websocket.TextMessage, []byte("OFFLINE"))
	mu.Lock()
	log.Println("Sending offline message")
	err := c.WriteJSON(ws.NewMessage(ws.MSG_OFFLINE, *source, "hub", "Collector is going offline"))
	if err != nil {
		log.Println("write:", err)
	}
	log.Println("Offline message sent. Sending WS close message")
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	mu.Unlock()
	if err != nil {
		log.Println("write close:", err)
	}
	log.Println("WS close message sent")
	select {
	case <-done:
	case <-time.After(time.Second):
	}
}
