package models

import (
	"context"
	"time"
)

//--------------------------------------------------------------------------------------|

// Notification represents a system alert for a user.
type Notification struct {
	ID               int       `json:"id"`
	UserID           int       `json:"user_id"`
	ActorID          int       `json:"actor_id"`
	ActorUsername    string    `json:"actor_username,omitempty"` // Joined field
	TargetType       string    `json:"target_type"`              // 'post', 'comment', 'group', etc.
	TargetID         int       `json:"target_id"`
	TargetTitle      string    `json:"target_title,omitempty"` // Joined field (e.g. post title)
	LinkID           int       `json:"link_id,omitempty"`      // ID to link to (e.g. post_id for comment)
	NotificationType string    `json:"notification_type"`      // 'like', 'comment', 'invite', etc.
	IsRead           bool      `json:"is_read"`
	CreatedAt        time.Time `json:"created_at"`
}

//--------------------------------------------------------------------------------------|

// NotificationRepo defines the interface for notification persistence.
type NotificationRepo interface {
	WithTx(tx any) NotificationRepo
	CreateNotification(ctx context.Context, userID, actorID int, targetType string, targetID int, notificationType string) error
	GetUserNotifications(ctx context.Context, userID, limit, offset int) ([]Notification, error)
	GetUnreadCount(ctx context.Context, userID int) (int, error)
	MarkAsRead(ctx context.Context, notificationID, userID int) error
	MarkAllAsRead(ctx context.Context, userID int) error
}

//--------------------------------------------------------------------------------------|

// NotificationService defines the interface for high-level notification management.
type NotificationService interface {
	Notify(ctx context.Context, userID, actorID int, actorUsername string, targetType string, targetID int, notifType string)
	GetNotifications(ctx context.Context, userID, limit, offset int) ([]Notification, error)
	GetUnreadCount(ctx context.Context, userID int) (int, error)
	MarkAsRead(ctx context.Context, notificationID, userID int) error
	MarkAllAsRead(ctx context.Context, userID int) error
}
