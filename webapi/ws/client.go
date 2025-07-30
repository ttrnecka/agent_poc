package ws

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 100 * 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub for one websocket connection
type Client struct {
	hub *Hub

	conn *websocket.Conn

	send chan []byte
}

// ServeWs handles WebSocket requests from clients. It upgrades the HTTP connection to a WebSocket,
// creates a new client instance with a buffered send channel, registers the client with the hub,
// and starts goroutines for reading from and writing to the WebSocket connection.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error().Err(err).Msg("")
		return
	}

	logger.Info().Str("client", conn.RemoteAddr().String()).Msg("Websocket client opened conenction")
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

// setupPongHandler sets a custom Pong handler for the given WebSocket connection.
// The handler logs the receipt of a Pong message from the client, including the client's remote address and any application data.
// It also updates the read deadline to ensure the connection remains alive.
// This function helps in detecting and handling client responsiveness in WebSocket communication.
func setupPongHandler(conn *websocket.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	conn.SetPongHandler(func(appData string) error {
		logger.Debug().Str("cliet", remoteAddr).Str("raw", appData).Msg("Received pong from client")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	setupPongHandler(c.conn)
	for {
		_, message, err := c.conn.ReadMessage()
		// if we cannot read message, we close the connection and drop the client list
		if err != nil {
			logger.Error().Err(err).Str("client", c.conn.RemoteAddr().String()).Msg("WS read error")
			break
		}
		// refresh deadline may be required if the client is not sending pongs to my pings
		// this may happend because the client is processing request in the same go routine that
		// reads the socket and the requets takes longer than deadline
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// right now every message is broadcasted to all clients (for POC purposes)
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		// Client received message to send
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// The hub closed the channel.
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Error().Err(err).Str("client", c.conn.RemoteAddr().String()).Msg("NextWriter error")
				return
			}
			// Write the message to the websocket connection.
			// do not return if there is failure as we want to run Close
			_, err = w.Write(message)
			if err != nil {
				logger.Error().Err(err).Str("client", c.conn.RemoteAddr().String()).Msg("NextWriter error")
				// return
			}
			// If you uncomment this all the json messages will be sent in one go and the client will have to handle it
			// in the future you may consider sending array of messages instead of one by one
			// frontend has already been updated to handle this by newline delimited messages
			// Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				logger.Error().Err(err).Str("client", c.conn.RemoteAddr().String()).Msg("WS close error")
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Error().Err(err).Str("client", c.conn.RemoteAddr().String()).Msg("WS write ping error")
				return
			}
		}
	}
}
