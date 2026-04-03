package models

import (
	"context"
	"time"
)

// Message represents a chat message between users.
type Message struct {
	ID         int       `db:"id" json:"id"`
	SenderID   int       `db:"sender_id" json:"sender_id"`
	ReceiverID int       `db:"receiver_id" json:"receiver_id"`
	Username   string    `json:"username"` // Sender's username
	Body       string    `db:"body" json:"body"`
	ImageURL   *string   `db:"image_url" json:"image_url,omitempty"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

//--------------------------------------------------------------------------------------|

// ChatRepo defines the contract for persisting 1:1 chat messages.
type ChatRepo interface {
	WithTx(tx any) ChatRepo
	SaveMessage(ctx context.Context, senderID, receiverID int, body string, imageURL *string) (*Message, error)
	GetMessages(ctx context.Context, user1ID, user2ID, limit, offset int) ([]Message, error)
}

// ChatService defines the 1:1 chat business logic.
type ChatService interface {
	SendMessage(ctx context.Context, senderID, receiverID int, body string, imageURL *string) (*Message, error)
	GetChatHistory(ctx context.Context, user1ID, user2ID, limit, offset int) ([]Message, error)
}
