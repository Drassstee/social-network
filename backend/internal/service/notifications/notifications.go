package notifications

import (
	"context"
	"encoding/json"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

type Signaller interface {
	SendToUser(userID int64, data []byte)
}

//--------------------------------------------------------------------------------------|

type Service struct {
	repo models.NotificationRepo
	hub  Signaller
}

//--------------------------------------------------------------------------------------|

func NewService(repo models.NotificationRepo, hub Signaller) *Service {
	return &Service{
		repo: repo,
		hub:  hub,
	}
}

//--------------------------------------------------------------------------------------|

func (s *Service) Notify(ctx context.Context, userID, actorID int64, actorUsername string, targetType string, targetID int64, notifType string) {
	_ = s.repo.CreateNotification(ctx, userID, actorID, targetType, targetID, notifType)

	message := s.generateMessage(actorUsername, targetType, notifType)

	// 3. Send real-time signal via WebSocket

	if s.hub != nil {
		notif := map[string]interface{}{
			"type": "notification",
			"data": map[string]interface{}{
				"type":           notifType,
				"actor_username": actorUsername,
				"target_type":    targetType,
				"target_id":      targetID,
				"message":        message,
			},
		}
		b, _ := json.Marshal(notif)
		s.hub.SendToUser(userID, b)
	}
}


//--------------------------------------------------------------------------------------|

func (s *Service) GetNotifications(ctx context.Context, userID int64, limit, offset int) ([]models.Notification, error) {
	notifs, err := s.repo.GetUserNotifications(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range notifs {
		notifs[i].Message = s.generateMessage(notifs[i].ActorUsername, notifs[i].TargetType, notifs[i].NotificationType)
	}

	return notifs, nil
}

//--------------------------------------------------------------------------------------|

func (s *Service) generateMessage(actorUsername, targetType, notifType string) string {
	switch notifType {
	case "follow":
		return actorUsername + " started following you"
	case "request":
		if targetType == "group" {
			return actorUsername + " requested to join your group"
		}
		return actorUsername + " sent you a follow request"
	case "invite":
		return actorUsername + " invited you to a group"
	case "accept":
		if targetType == "group" {
			return actorUsername + " accepted your join request"
		}
		return actorUsername + " accepted your follow request"
	case "decline":
		return actorUsername + " declined your request"
	case "event":
		return actorUsername + " created a new event in your group"
	default:
		return "New notification from " + actorUsername
	}
}

//--------------------------------------------------------------------------------------|

func (s *Service) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}

//--------------------------------------------------------------------------------------|

func (s *Service) MarkAsRead(ctx context.Context, notificationID, userID int64) error {
	return s.repo.MarkAsRead(ctx, notificationID, userID)
}

//--------------------------------------------------------------------------------------|

func (s *Service) MarkAllAsRead(ctx context.Context, userID int64) error {
	return s.repo.MarkAllAsRead(ctx, userID)
}
