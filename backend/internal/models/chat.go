package models

import "time"

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
