package models

import (
	"context"
	"time"
)

//--------------------------------------------------------------------------------------|

// Message represents a chat message between users.
type Message struct {
	ID         int64     `db:"id" json:"id"`
	SenderID   int64     `db:"sender_id" json:"sender_id"`
	ReceiverID int64     `db:"receiver_id" json:"receiver_id"`
	Username   string    `json:"username"` // Sender's username
	Body       string    `db:"body" json:"body"`
	ImageURL   *string   `db:"image_url" json:"image_url,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type OnlineUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

//--------------------------------------------------------------------------------------|

// ChatRepo defines the contract for persisting 1:1 chat messages.
type ChatRepo interface {
	WithTx(tx any) ChatRepo
	SaveMessage(ctx context.Context, senderID, receiverID int64, body string, imageURL *string) (*Message, error)
	GetMessages(ctx context.Context, user1ID, user2ID int64, limit, offset int) ([]Message, error)
}

//--------------------------------------------------------------------------------------|

// ChatService defines the 1:1 chat business logic.
type ChatService interface {
	SendMessage(ctx context.Context, senderID, receiverID int64, body string, imageURL *string) (*Message, error)
	GetChatHistory(ctx context.Context, user1ID, user2ID int64, limit, offset int) ([]Message, error)
	GetChatableOnlineUsers(ctx context.Context, userID int64, onlineIDs []int64) ([]OnlineUser, error)
}

//--------------------------------------------------------------------------------------|

// Hub defines the interface for real-time communication management.
type Hub interface {
	UpdateUserGroups(userID int64)
}
