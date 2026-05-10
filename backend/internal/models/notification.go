package models

import (
	"context"
	"time"
)

//--------------------------------------------------------------------------------------|

type Notification struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	ActorID          int64     `json:"actor_id"`
	ActorUsername    string    `json:"actor_username,omitempty"`
	TargetType       string    `json:"target_type"`
	TargetID         int64     `json:"target_id"`
	TargetTitle      string    `json:"target_title,omitempty"`
	LinkID           int64     `json:"link_id,omitempty"`
	NotificationType string    `json:"notification_type"`
	Message          string    `json:"message"`
	IsRead           bool      `json:"is_read"`

	CreatedAt        time.Time `json:"created_at"`
}

//--------------------------------------------------------------------------------------|

type NotificationRepo interface {
	WithTx(tx any) NotificationRepo
	CreateNotification(ctx context.Context, userID, actorID int64, targetType string, targetID int64, notificationType string) error
	GetUserNotifications(ctx context.Context, userID int64, limit, offset int) ([]Notification, error)
	GetUnreadCount(ctx context.Context, userID int64) (int, error)
	MarkAsRead(ctx context.Context, notificationID, userID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
}

//--------------------------------------------------------------------------------------|

type NotificationService interface {
	Notify(ctx context.Context, userID, actorID int64, actorUsername string, targetType string, targetID int64, notifType string)
	GetNotifications(ctx context.Context, userID int64, limit, offset int) ([]Notification, error)
	GetUnreadCount(ctx context.Context, userID int64) (int, error)
	MarkAsRead(ctx context.Context, notificationID, userID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
}
