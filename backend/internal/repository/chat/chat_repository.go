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
func (r *sqlChatRepository) SaveMessage(ctx context.Context, senderID, receiverID int64, body string, imageURL *string) (*models.Message, error) {
	msg := &models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Body:       body,
		ImageURL:   imageURL,
		IsRead:     false,
	}

	query := `
		INSERT INTO messages (sender_id, receiver_id, body, image_url) 
		VALUES (?, ?, ?, ?) 
		RETURNING id, created_at, (SELECT COALESCE(nickname, '') FROM users WHERE id = sender_id)`

	err := r.db.QueryRowContext(ctx, query, senderID, receiverID, body, imageURL).
		Scan(&msg.ID, &msg.CreatedAt, &msg.Username)

	if err != nil {
		return nil, err
	}

	return msg, nil
}

//--------------------------------------------------------------------------------------|

// GetMessages retrieves a paginated history of messages exchanged between two users.
func (r *sqlChatRepository) GetMessages(ctx context.Context, user1ID, user2ID int64, limit, offset int) ([]models.Message, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT m.id, m.sender_id, m.receiver_id, COALESCE(u.nickname, ''), m.body, m.image_url, m.is_read, m.created_at
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
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ReceiverID, &m.Username, &m.Body, &m.ImageURL, &m.IsRead, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}

//--------------------------------------------------------------------------------------|

// MarkAsRead marks all messages from sender to receiver as read.
func (r *sqlChatRepository) MarkAsRead(ctx context.Context, senderID, receiverID int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE messages SET is_read = 1 WHERE sender_id = ? AND receiver_id = ? AND is_read = 0`,
		senderID, receiverID)
	return err
}

//--------------------------------------------------------------------------------------|

// GetUnreadCounts returns a map of senderID -> unread count for the given user.
func (r *sqlChatRepository) GetUnreadCounts(ctx context.Context, userID int64) (map[int64]int, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT sender_id, COUNT(*) FROM messages WHERE receiver_id = ? AND is_read = 0 GROUP BY sender_id`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[int64]int)
	for rows.Next() {
		var senderID int64
		var count int
		if err := rows.Scan(&senderID, &count); err != nil {
			return nil, err
		}
		counts[senderID] = count
	}
	return counts, nil
}

