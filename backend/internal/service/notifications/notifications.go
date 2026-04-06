package notifications

import (
	"context"
	"encoding/json"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

// Signaller defines the interface for sending real-time signals to users (e.g., via WebSocket).
type Signaller interface {
	// SendToUser sends the provided data to a specific user.
	SendToUser(userID int, data []byte)
}

//--------------------------------------------------------------------------------------|

// Service implements models.NotificationService.
type Service struct {
	repo models.NotificationRepo
	hub  Signaller
}

//--------------------------------------------------------------------------------------|

// NewService creates a new notifications service.
func NewService(repo models.NotificationRepo, hub Signaller) *Service {
	return &Service{
		repo: repo,
		hub:  hub,
	}
}

//--------------------------------------------------------------------------------------|

// Notify creates a notification in the database and sends a real-time signal to the target user.
func (s *Service) Notify(ctx context.Context, userID, actorID int, actorUsername string, targetType string, targetID int, notifType string) {
	// 1. Create DB notification
	_ = s.repo.CreateNotification(ctx, userID, actorID, targetType, targetID, notifType)

	// 2. Send real-time signal via WebSocket
	if s.hub != nil {
		notif := map[string]interface{}{
			"type": "notification",
			"data": map[string]interface{}{
				"type":           notifType,
				"actor_username": actorUsername,
				"target_type":    targetType,
				"target_id":      targetID,
			},
		}
		b, _ := json.Marshal(notif)
		s.hub.SendToUser(userID, b)
	}
}

//--------------------------------------------------------------------------------------|

func (s *Service) GetNotifications(ctx context.Context, userID, limit, offset int) ([]models.Notification, error) {
	return s.repo.GetUserNotifications(ctx, userID, limit, offset)
}

//--------------------------------------------------------------------------------------|

func (s *Service) GetUnreadCount(ctx context.Context, userID int) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

//--------------------------------------------------------------------------------------|

func (s *Service) MarkAsRead(ctx context.Context, notificationID, userID int) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

//--------------------------------------------------------------------------------------|

func (s *Service) MarkAllAsRead(ctx context.Context, userID int) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}
