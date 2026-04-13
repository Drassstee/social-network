package chatsvc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"social-network/internal/models"
	"social-network/internal/utils"
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
	clients    map[int]*Client // Registered clients by UserID.
	broadcast  chan []byte     // Inbound messages from the clients.
	register   chan *Client    // Register requests from the clients.
	unregister chan *Client    // Unregister requests from clients.

	mu sync.RWMutex

	ChatRepo  ChatRepository
	UserRepo  UserRepository
	GroupRepo models.GroupRepo

	// Optimization: Cache online group memberships and usernames
	userCache    *utils.Cache
	groupMembers map[int]map[int]bool // groupID -> set of userIDs
}

//--------------------------------------------------------------------------------------|

// NewHub creates a new instance of the Hub.
func NewHub(chatRepo ChatRepository, userRepo UserRepository, groupRepo models.GroupRepo) *Hub {
	return &Hub{
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		clients:      make(map[int]*Client),
		ChatRepo:     chatRepo,
		UserRepo:     userRepo,
		GroupRepo:    groupRepo,
		userCache:    utils.NewCache(),
		groupMembers: make(map[int]map[int]bool),
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

			// Optimization: Pre-load group memberships for faster broadcasting
			// We do this in a goroutine to avoid blocking the hub loop
			go h.updateClientGroupMemberships(client.UserID, true)

			h.mu.Unlock()
			// Notify others that this user is now online
			h.broadcastStatusUpdate(client.UserID, true)

		case client := <-h.unregister:
			// Remove a client connection and close its sending channel
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.send)

				// Optimization: Remove user from group tracking
				h.updateClientGroupMembershipsLocked(client.UserID, false)

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
	case "group_message":
		h.handleGroupMessage(wsMsg)
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

func (h *Hub) handleGroupMessage(wsMsg wsMessage) {
	var data struct {
		GroupID  int     `json:"group_id"`
		Body     string  `json:"body"`
		ImageURL *string `json:"image_url"`
	}
	if err := json.Unmarshal(wsMsg.Data, &data); err != nil {
		log.Printf("HandleGroupMessage unmarshal error: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// 1. Verify membership
	isMember, err := h.GroupRepo.IsMember(ctx, data.GroupID, wsMsg.Sender)
	if err != nil || !isMember {
		log.Printf("Unauthorized group message attempt from %d to group %d", wsMsg.Sender, data.GroupID)
		return
	}

	// 2. Save message
	msg := &models.GroupMessage{
		GroupID:  data.GroupID,
		SenderID: wsMsg.Sender,
		Body:     data.Body,
		ImageURL: data.ImageURL,
	}
	if err := h.GroupRepo.SaveGroupMessage(ctx, msg); err != nil {
		log.Printf("Failed to save group message: %v", err)
		return
	}

	// 3. Broadcast to online members using in-memory tracking
	h.mu.RLock()
	defer h.mu.RUnlock()

	if members, ok := h.groupMembers[data.GroupID]; ok {
		for userID := range members {
			h.sendToUser(userID, "group_message", msg)
		}
	}
}

//--------------------------------------------------------------------------------------|

func (h *Hub) updateClientGroupMemberships(userID int, join bool) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	groupIDs, err := h.GroupRepo.GetMemberGroupIDs(ctx, userID)
	if err != nil {
		log.Printf("Failed to fetch group IDs for user %d: %v", userID, err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.updateClientGroupMembershipsLocked(userID, join, groupIDs...)
}

//--------------------------------------------------------------------------------------|

func (h *Hub) updateClientGroupMembershipsLocked(userID int, join bool, groupIDs ...int) {
	if len(groupIDs) == 0 && !join {
		// If unregistering and we don't have groupIDs, we have to find which groups this user was in.
		// For simplicity, we can just iterate over all groups in the map.
		for gID, members := range h.groupMembers {
			delete(members, userID)
			if len(members) == 0 {
				delete(h.groupMembers, gID)
			}
		}
		return
	}

	for _, gID := range groupIDs {
		if join {
			if h.groupMembers[gID] == nil {
				h.groupMembers[gID] = make(map[int]bool)
			}
			h.groupMembers[gID][userID] = true
		} else {
			if members, ok := h.groupMembers[gID]; ok {
				delete(members, userID)
				if len(members) == 0 {
					delete(h.groupMembers, gID)
				}
			}
		}
	}
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
	// Fetch username for the user (using cache if available)
	username := ""
	cacheKey := fmt.Sprintf("u:%d", userID)
	if val, found := h.userCache.Get(cacheKey); found {
		username = val.(string)
	} else if h.UserRepo != nil {
		_, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()
		user, err := h.UserRepo.GetByID(int64(userID))
		if err == nil && user != nil {
			username = user.Nickname
			h.userCache.Set(cacheKey, username, 1*time.Hour)
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
