// Package websocket provides real-time communication capabilities using WebSockets,
// including chat, notifications, and user status tracking.
package websocket

import (
	"context"
	"database/sql"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

// ChatRepository defines the interface for persisting and retrieving chat messages.
type ChatRepository interface {
	// SaveMessage stores a new chat message in the database.
	SaveMessage(ctx context.Context, senderID, receiverID int, body string, imageURL *string) (*models.Message, error)
	// GetMessages retrieves the chat history between two users, ordered by date.
	GetMessages(ctx context.Context, user1ID, user2ID, limit, offset int) ([]models.Message, error)
}

type sqlChatRepository struct {
	db *sql.DB
}

// NewChatRepository creates a new instance of the chat repository.
func NewChatRepository(db *sql.DB) ChatRepository {
	return &sqlChatRepository{db: db}
}

//--------------------------------------------------------------------------------------|

// SaveMessage stores a new chat message in the database and returns the message details.
func (r *sqlChatRepository) SaveMessage(ctx context.Context, senderID, receiverID int, body string, imageURL *string) (*models.Message, error) {
	msg := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Body:       body,
		ImageURL:   imageURL,
	}

	query := `
		INSERT INTO messages (sender_id, receiver_id, body, image_url) 
		VALUES (?, ?, ?, ?) 
		RETURNING id, created_at, (SELECT username FROM users WHERE id = sender_id)`

	err := r.db.QueryRowContext(ctx, query, senderID, receiverID, body, imageURL).
		Scan(&msg.ID, &msg.CreatedAt, &msg.Username)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

//--------------------------------------------------------------------------------------|

// GetMessages retrieves a paginated history of messages exchanged between two users,
// sorted by creation time (descending).
func (r *sqlChatRepository) GetMessages(ctx context.Context, user1ID, user2ID, limit, offset int) ([]models.Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT m.id, m.sender_id, m.receiver_id, u.username, m.body, m.image_url, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC, m.id DESC
		LIMIT ? OFFSET ?`,
		user1ID, user2ID, user2ID, user1ID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var m models.Message
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Username, &m.Body, &m.ImageURL, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
