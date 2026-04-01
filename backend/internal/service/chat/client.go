// Package websocket provides real-time communication capabilities using WebSockets,
// including chat, notifications, and user status tracking.
package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

//--------------------------------------------------------------------------------------|

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
)

//--------------------------------------------------------------------------------------|

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		// For robustness, ensure origin matches host
		return strings.Contains(origin, r.Host)
	},
}

//--------------------------------------------------------------------------------------|

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	UserID int
}

//--------------------------------------------------------------------------------------|

// readPump reads messages from the WebSocket connection and pushes them to the Hub.
// It handles ping/pong heartbeat messages and processes incoming JSON commands.
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// In an SPA forum, we might receive messages like "typing" or "chat_message"
		var wsMsg struct {
			Type string          `json:"type"`
			Data json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			continue
		}

		// Attach sender ID for server-side routing
		outbound := struct {
			Type   string          `json:"type"`
			Sender int             `json:"sender"`
			Data   json.RawMessage `json:"data"`
		}{
			Type:   wsMsg.Type,
			Sender: c.UserID,
			Data:   wsMsg.Data,
		}

		b, _ := json.Marshal(outbound)
		c.Hub.broadcast <- b
	}
}

//--------------------------------------------------------------------------------------|

// writePump listens for messages from the Hub and writes them to the WebSocket connection.
// It also sends periodic ping messages to keep the connection alive.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

//--------------------------------------------------------------------------------------|

// ServeWs upgrades the HTTP connection to a WebSocket connection, creates a new Client,
// and starts the read and write pumps to handle real-time communication.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, userID int) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{Hub: hub, Conn: conn, send: make(chan []byte, 256), UserID: userID}
	client.Hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
