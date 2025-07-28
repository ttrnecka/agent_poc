package main

import (
	"log"
	"net/url"
	"sync"
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
	id       string
	c        *websocket.Conn
	messages chan Message
	done     chan struct{}
	watcher  *Watcher
	ticker   *time.Ticker
	mu       sync.Mutex
}

func NewMessageHandler(addr, id string, watcher *Watcher) *MessageHandler {
	return &MessageHandler{
		addr:     addr,
		id:       id,
		done:     make(chan struct{}),
		messages: make(chan Message, 100),
		watcher:  watcher,
		ticker:   time.NewTicker(5 * time.Second),
	}
}

// run it in goroutine else it will block
func (m *MessageHandler) Start() {
	err := m.connectWebSocket()
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer m.c.Close()
	go m.readLoop()
	m.processLoop()
}

func (m *MessageHandler) Stop() {
	log.Println("Stoping MessageHandler")
	m.closeWebSocket()
}

func (m *MessageHandler) readLoop() {
	defer close(m.done)

	for {
		mes := ws.Message{}

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
			err := m.sendHeartbeat()
			if err != nil {
				log.Printf("Failed to send heartbeat message: %v", err)
			}
		case msg, ok := <-m.messages:
			if !ok {
				return
			}
			// TODO needs to make this into pool of worker or as go routines
			log.Printf("websocket received: %v", msg)
			if msg.Destination == m.id {
				if msg.Type == ws.MSG_REFRESH {
					refresh()
				}
				if msg.Type == ws.MSG_RUN {
					run(msg.Message, m)
					log.Printf("Sending process message")
					m.watcher.Process()
					log.Printf("Sent process message")
				}
			}
		}
	}
}

func (m *MessageHandler) closeWebSocket() {
	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	// err := c.WriteMessage(websocket.TextMessage, []byte("OFFLINE"))
	log.Println("Sending offline message")

	err := m.SendMessage(ws.NewMessage(ws.MSG_OFFLINE, m.id, "hub", "Collector is going offline"))
	if err == nil {
		log.Println("Offline message sent. Sending WS close message")
	}

	//special message type
	m.mu.Lock()
	err = m.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	m.mu.Unlock()
	if err != nil {
		log.Println("write close:", err)
	} else {
		log.Println("WS close message sent")
	}
	select {
	case <-m.done:
	case <-time.After(time.Second):
	}
}

func (m *MessageHandler) sendHeartbeat() error {
	log.Println("Sending heartbeat message")

	err := m.SendMessage(ws.NewMessage(ws.MSG_ONLINE, m.id, "hub", "Collector is online"))
	if err != nil {
		return err
	}
	log.Println("Heartbeat sent")
	return nil
}

func (m *MessageHandler) connectWebSocket() error {
	u := url.URL{Scheme: "ws", Host: m.addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	m.c = c
	return nil
}

func (m *MessageHandler) SendMessage(message ws.Message) error {
	m.mu.Lock()
	log.Printf("Sending WS message")
	err := m.c.WriteJSON(message)
	m.mu.Unlock()
	if err != nil {
		log.Printf("Sending WS message failed: %v\n", err)
	} else {
		log.Printf("Sending WS message succeeded")
	}
	return err
}
