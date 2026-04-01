// Package websocket provides real-time communication capabilities using WebSockets,
// including chat, notifications, and user status tracking.
package websocket

import (
	"context"
	"encoding/json"
	"log"
	"social-network/internal/notifications"
	"sync"
	"time"
)

//--------------------------------------------------------------------------------------|

const (
	dbTimeout = 5 * time.Second
)

//--------------------------------------------------------------------------------------|

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients by UserID.
	clients map[int]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	mu sync.RWMutex

	ChatRepo          ChatRepository
	NotificationsRepo notifications.Repository
	UserRepo          UserRepository
}

//--------------------------------------------------------------------------------------|

// NewHub creates a new instance of the Hub.
func NewHub(chatRepo ChatRepository, notificationsRepo notifications.Repository, userRepo UserRepository) *Hub {
	return &Hub{
		broadcast:         make(chan []byte),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		clients:           make(map[int]*Client),
		ChatRepo:          chatRepo,
		NotificationsRepo: notificationsRepo,
		UserRepo:          userRepo,
	}
}

//--------------------------------------------------------------------------------------|

type wsMessage struct {
	Type   string          `json:"type"`
	Sender int             `json:"sender,omitempty"`
	Data   json.RawMessage `json:"data"`
}

//--------------------------------------------------------------------------------------|

// Run starts the hub's main event loop. It handles client registration,
// unregistration, and message broadcasting. It runs until the context is canceled.
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Clean up all connections on context cancellation (e.g., server shutdown)
			h.mu.Lock()
			for _, client := range h.clients {
				close(client.send)
				delete(h.clients, client.UserID)
			}
			h.mu.Unlock()
			return
		case client := <-h.register:
			// Register a new client connection
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			// Notify others that this user is now online
			h.broadcastStatusUpdate(client.UserID, true)

		case client := <-h.unregister:
			// Remove a client connection and close its sending channel
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.send)
				h.mu.Unlock()
				// Notify others that this user is now offline
				h.broadcastStatusUpdate(client.UserID, false)
			} else {
				h.mu.Unlock()
			}

		case message := <-h.broadcast:
			// Process inbound messages from clients
			h.handleInbound(message)
		}
	}
}

//--------------------------------------------------------------------------------------|

// handleInbound unmarshals and routes incoming WebSocket messages based on their type.
func (h *Hub) handleInbound(message []byte) {
	var wsMsg wsMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		log.Printf("WebSocket unmarshal error: %v", err)
		return
	}

	switch wsMsg.Type {
	case "private_message":
		h.handlePrivateMessage(wsMsg)
	case "typing":
		h.handleTypingIndicator(wsMsg, true)
	case "stop_typing":
		h.handleTypingIndicator(wsMsg, false)
	default:
		// Broadcast generic messages to all connected clients
		h.doBroadcast(message)
	}
}

//--------------------------------------------------------------------------------------|

func (h *Hub) handlePrivateMessage(wsMsg wsMessage) {
	var data struct {
		ReceiverID int     `json:"receiver_id"`
		Body       string  `json:"body"`
		ImageURL   *string `json:"image_url"`
	}
	if err := json.Unmarshal(wsMsg.Data, &data); err != nil {
		log.Printf("HandlePrivateMessage unmarshal error: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	msg, err := h.ChatRepo.SaveMessage(ctx, wsMsg.Sender, data.ReceiverID, data.Body, data.ImageURL)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return
	}

	h.sendToUser(data.ReceiverID, "private_message", msg)
	h.sendToUser(wsMsg.Sender, "private_message", msg)
}

//--------------------------------------------------------------------------------------|

func (h *Hub) handleTypingIndicator(wsMsg wsMessage, isTyping bool) {
	var data struct {
		ReceiverID int    `json:"receiver_id"`
		SenderName string `json:"sender_name"`
	}
	if err := json.Unmarshal(wsMsg.Data, &data); err != nil {
		log.Printf("HandleTypingIndicator unmarshal error: %v", err)
		return
	}

	// Determine message type
	msgType := "stop_typing"
	if isTyping {
		msgType = "typing"
	}

	// Create typing indicator message
	typingMsg := map[string]interface{}{
		"sender_id":   wsMsg.Sender,
		"sender_name": data.SenderName,
	}

	// Send to receiver only
	h.sendToUser(data.ReceiverID, msgType, typingMsg)
}

//--------------------------------------------------------------------------------------|

func (h *Hub) broadcastStatusUpdate(userID int, online bool) {
	// Fetch username for the user
	username := ""
	if h.UserRepo != nil {
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()
		user, err := h.UserRepo.GetByID(ctx, userID)
		if err == nil && user != nil {
			username = user.Username
		}
	}

	h.broadcastJSON(map[string]interface{}{
		"type": "status_update",
		"data": map[string]interface{}{
			"user_id":  userID,
			"username": username,
			"online":   online,
		},
	})
}

//--------------------------------------------------------------------------------------|

func (h *Hub) sendToUser(userID int, msgType string, data any) {
	payload := map[string]any{
		"type": msgType,
		"data": data,
	}
	h.SendToUserJSON(userID, payload)
}

//--------------------------------------------------------------------------------------|

func (h *Hub) broadcastJSON(payload any) {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal broadcast payload: %v", err)
		return
	}
	h.doBroadcast(b)
}

//--------------------------------------------------------------------------------------|

func (h *Hub) doBroadcast(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, client := range h.clients {
		select {
		case client.send <- message:
		default:
		}
	}
}

//--------------------------------------------------------------------------------------|

// SendToUserJSON encodes the payload as JSON and delivers it to a specific user's WebSocket.
func (h *Hub) SendToUserJSON(userID int, payload interface{}) {
	b, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal payload for user %d: %v", userID, err)
		return
	}
	h.SendToUser(userID, b)
}

//--------------------------------------------------------------------------------------|

// SendToUser sends a raw byte message to a specific user if they are connected.
func (h *Hub) SendToUser(userID int, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[userID]; ok {
		select {
		case client.send <- message:
		default:
			// Non-blocking send
		}
	}
}

//--------------------------------------------------------------------------------------|

// GetOnlineUsers returns a slice of IDs for all currently connected users.
func (h *Hub) GetOnlineUsers() []int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	online := make([]int, 0, len(h.clients))
	for userID := range h.clients {
		online = append(online, userID)
	}
	return online
}
