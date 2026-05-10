package models

import (
	"context"
	"time"
)

//--------------------------------------------------------------------------------------|

type Message struct {
	ID         int64     `db:"id" json:"id"`
	SenderID   int64     `db:"sender_id" json:"sender_id"`
	ReceiverID int64     `db:"receiver_id" json:"receiver_id"`
	Username   string    `json:"username"`
	Body       string    `db:"body" json:"body"`
	ImageURL   *string   `db:"image_url" json:"image_url,omitempty"`
	IsRead     bool      `db:"is_read" json:"is_read"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}


type OnlineUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

//--------------------------------------------------------------------------------------|

type ChatRepo interface {
	WithTx(tx any) ChatRepo
	SaveMessage(ctx context.Context, senderID, receiverID int64, body string, imageURL *string) (*Message, error)
	GetMessages(ctx context.Context, user1ID, user2ID int64, limit, offset int) ([]Message, error)
	MarkAsRead(ctx context.Context, senderID, receiverID int64) error
	GetUnreadCounts(ctx context.Context, userID int64) (map[int64]int, error)
}


//--------------------------------------------------------------------------------------|

type ChatService interface {
	SendMessage(ctx context.Context, senderID, receiverID int64, body string, imageURL *string) (*Message, error)
	GetChatHistory(ctx context.Context, user1ID, user2ID int64, limit, offset int) ([]Message, error)
	GetChatableUsers(ctx context.Context, userID int64) ([]OnlineUser, error)
	MarkAsRead(ctx context.Context, senderID, receiverID int64) error
	GetUnreadCounts(ctx context.Context, userID int64) (map[int64]int, error)
}


//--------------------------------------------------------------------------------------|

type Hub interface {
	UpdateUserGroups(userID int64)
}
