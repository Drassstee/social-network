package notifications

import (
	"context"
	"database/sql"
	"fmt"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

// sqlRepository implements models.NotificationRepo.
type dbQuerier interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

//--------------------------------------------------------------------------------------|

type sqlRepository struct {
	db dbQuerier
}

//--------------------------------------------------------------------------------------|

// NewRepository creates a new instance of the notifications repository.
func NewRepository(db *sql.DB) models.NotificationRepo {
	return &sqlRepository{db: db}
}

//--------------------------------------------------------------------------------------|

func (r *sqlRepository) WithTx(tx any) models.NotificationRepo {
	if tx == nil {
		return r
	}
	if t, ok := tx.(*sql.Tx); ok {
		return &sqlRepository{db: t}
	}
	return r
}

//--------------------------------------------------------------------------------------|

func (r *sqlRepository) CreateNotification(ctx context.Context, userID, actorID int, targetType string, targetID int, notificationType string) error {
	if userID == actorID {
		return nil
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO notifications (user_id, actor_id, target_type, target_id, notification_type)
		 VALUES (?, ?, ?, ?, ?)`,
		userID, actorID, targetType, targetID, notificationType)
	return err
}

//--------------------------------------------------------------------------------------|

func (r *sqlRepository) GetUserNotifications(ctx context.Context, userID, limit, offset int) ([]models.Notification, error) {
	query := `
		SELECT 
			n.id, n.user_id, n.actor_id, u.username as actor_username,
			n.target_type, n.target_id,
			COALESCE(CASE 
				WHEN n.target_type = 'post' THEN p.title
				WHEN n.target_type = 'comment' THEN (SELECT p2.title FROM posts p2 JOIN comments c ON c.post_id = p2.id WHERE c.id = n.target_id)
				ELSE ''
			END, 'Deleted post/comment') as target_title,
			COALESCE(CASE 
				WHEN n.target_type = 'post' THEN n.target_id
				WHEN n.target_type = 'comment' THEN (SELECT post_id FROM comments WHERE id = n.target_id)
				ELSE 0
			END, 0) as link_id,
			n.notification_type, n.is_read, n.created_at
		FROM notifications n
		JOIN users u ON n.actor_id = u.id
		LEFT JOIN posts p ON n.target_type = 'post' AND n.target_id = p.id
		WHERE n.user_id = ?
		ORDER BY n.created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []models.Notification{}
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.ActorID, &n.ActorUsername,
			&n.TargetType, &n.TargetID, &n.TargetTitle, &n.LinkID,
			&n.NotificationType, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

//--------------------------------------------------------------------------------------|

func (r *sqlRepository) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0`,
		userID).Scan(&count)
	return count, err
}

//--------------------------------------------------------------------------------------|

func (r *sqlRepository) MarkAsRead(ctx context.Context, notificationID, userID int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?`,
		notificationID, userID)
	if err != nil {
		return fmt.Errorf("mark notification as read failed: %w", err)
	}
	return nil
}

//--------------------------------------------------------------------------------------|

func (r *sqlRepository) MarkAllAsRead(ctx context.Context, userID int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE notifications SET is_read = 1 WHERE user_id = ? AND is_read = 0`,
		userID)
	if err != nil {
		return fmt.Errorf("mark all as read failed: %w", err)
	}
	return nil
}
