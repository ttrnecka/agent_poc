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
	h := &MessageHandler{
		addr:     addr,
		id:       id,
		done:     make(chan struct{}),
		messages: make(chan Message, 100),
		watcher:  watcher,
		ticker:   time.NewTicker(5 * time.Second),
	}

	// // Debug: watch for when h is GC'd
	// runtime.SetFinalizer(h, func(m *MessageHandler) {
	// 	log.Println("MessageHandler finalized (GC collected)")
	// })
	return h
}

// Start() opens a websocket connection and starts readLoop and processLoop in separate goroutines and returns immediately
// If there is issue opening websocket the loops will return and close the done channel
// Caller should wait for the done channel to be closed and then either close program or try to open new MessageHandler
// Caller should call Stop() if they want to close the handler. There is no need to call the Stop() if the done channel was closed

func (m *MessageHandler) Start() {
	err := m.connectWebSocket()
	if err != nil {
		log.Printf("error opening websocket: %s", err)
	}
	go m.readLoop()
	go m.processLoop()
}

// should only be called when wanting to intentially Stop the handler
// if you received closed handler.done then handler is cleaned up already
// Stop will notify about handler going offline
func (m *MessageHandler) Stop() {
	m.closeWebSocket()
	<-m.done
}

// closes all remaining channels and references
// does not closes the done as that is done only in readLoop
func (m *MessageHandler) cleanup() {
	close(m.messages)
	m.ticker.Stop()
	m.watcher = nil
	if m.c != nil {
		m.c.Close()
	}
}

// this loop closes when there is issue reading and closes done channel
func (m *MessageHandler) readLoop() {
	defer close(m.done)

	// if there is no websocket connection all loops get closed
	if m.c == nil {
		return
	}

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

// this loops closes when either done or messages are closed
func (m *MessageHandler) processLoop() {
	for {
		select {
		case <-m.ticker.C:
			err := m.sendHeartbeat()
			if err != nil {
				log.Printf("Failed to send heartbeat message: %v", err)
			}
		case <-m.done:
			// clean up rest of the resources
			m.cleanup()
			return
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
	// send offline message and WS close message
	// this result on server closing read pipe
	// which in turn closes done channel in readLoop

	if m.c == nil {
		return
	}

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
	// as long as this is called in the Stop only this does not make much sense
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
