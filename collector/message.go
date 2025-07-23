package main

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/ttrnecka/agent_poc/ws"
)

type Message struct {
	ws.Message
	c *websocket.Conn
}

var messages = make(chan Message, 100)

// reads messaes from websocket connections and sends them to the handler
func messageReader(c *websocket.Conn, done chan struct{}) {
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
}

// core message handler
func messageHandler() {
	for msg := range messages {
		log.Printf("recv: %v", msg)
		if msg.Destination == *source {
			if msg.Type == ws.MSG_REFRESH {
				refresh()
			}
			if msg.Type == ws.MSG_RUN {
				run(msg.Message, msg.c)
			}
		}
	}
}
