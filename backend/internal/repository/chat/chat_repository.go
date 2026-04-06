// Package chatrepo provides database persistence for chat messages.
package chatrepo

import (
	"context"
	"database/sql"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

type dbQuerier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

//--------------------------------------------------------------------------------------|

type sqlChatRepository struct {
	db dbQuerier
}

//--------------------------------------------------------------------------------------|

// NewChatRepository creates a new instance of the chat repository.
func NewChatRepository(db *sql.DB) models.ChatRepo {
	return &sqlChatRepository{db: db}
}

//--------------------------------------------------------------------------------------|

func (r *sqlChatRepository) WithTx(tx any) models.ChatRepo {
	if tx == nil {
		return r
	}
	if t, ok := tx.(*sql.Tx); ok {
		return &sqlChatRepository{db: t}
	}
	return r
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
