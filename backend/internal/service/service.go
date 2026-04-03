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
	User          models.UserService
	Post          *servicepost.PostService
	Group         models.GroupService
	Chat          models.ChatService
	Notifications models.NotificationService
}

func NewService(repo *repository.Repository, hub *chatsvc.Hub) *Service {
	notifSvc := notifications.NewService(repo.Notifications, hub)
	return &Service{
		User:          user.NewUserService(repo.User),
		Post:          servicepost.NewPostService(repo.Post),
		Group:         groupsvc.NewGroupService(repo.Group, notifSvc, repo.User, repo.DB),
		Chat:          chatsvc.NewChatService(repo.Chat),
		Notifications: notifSvc,
	}
}
