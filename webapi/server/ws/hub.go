package ws

import (
	"encoding/json"
	"fmt"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

var hubInstance *Hub

func GetHub() *Hub {
	if hubInstance == nil {
		hubInstance = NewHub()
	}
	return hubInstance
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte, 100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

func (h *Hub) Run() {
	for {
		select {
		// register client
		case client := <-h.register:
			h.clients[client] = true
		//underegister client and close its send channel
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		// broadcast message to all clients
		case message := <-h.broadcast:
			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				logger.Error().Err(err).Str("raw", fmt.Sprintf("%+v", message)).Msg("Unmarshal error")
				continue
			}
			logger.Debug().Str("raw", fmt.Sprintf("%+v", msg)).Msg("Received message")

			// Send the message to all registered clients
			// if the client's send channel is full, close it and remove the client
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
