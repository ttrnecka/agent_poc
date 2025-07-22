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

func messageHandler() {
	for {
		select {
		case msg := <-messages:
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
}
