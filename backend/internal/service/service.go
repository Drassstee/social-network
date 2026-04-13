package service

import (
	"social-network/internal/models"
	"social-network/internal/repository"
	chatsvc "social-network/internal/service/chat"
	groupsvc "social-network/internal/service/group"
	"social-network/internal/service/notifications"
	servicepost "social-network/internal/service/post"
	"social-network/internal/service/user"
)

type Service struct {
	User          *user.UserService
	Post          *servicepost.PostService
	Group         models.GroupService
	Chat          models.ChatService
	Notifications models.NotificationService
}

func NewService(r *repository.Repository, hub *chatsvc.Hub) *Service {
	notifSvc := notifications.NewService(r.Notifications, hub)
	return &Service{
		User:          user.NewUserService(r.User, r.Session, r.Follow, r.Post),
		Post:          servicepost.NewPostService(r.Post),
		Group:         groupsvc.NewGroupService(r.Group, notifSvc, r.User, r.DB),
		Chat:          chatsvc.NewChatService(r.Chat),
		Notifications: notifSvc,
	}
}
