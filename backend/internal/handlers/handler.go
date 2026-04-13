package handlers

import (
	"social-network/internal/models"

	"social-network/internal/handlers/chat"
	"social-network/internal/handlers/group"
	"social-network/internal/handlers/notifications"
	"social-network/internal/handlers/post"
	"social-network/internal/handlers/user"
	"social-network/internal/service"
	chatsvc "social-network/internal/service/chat"
)

type Handler struct {
	User          *user.UserHandler
	Post          *post.PostHandler
	Group         *group.GroupHandler
	Chat          *chat.ChatHandler
	Notifications *notifications.Handler
}

func NewHandler(service *service.Service, hub *chatsvc.Hub, userRepo models.UserRepo, chatRepo models.ChatRepo, uploader chat.ImageUploader) *Handler {
	return &Handler{
		User:          user.NewUserHandler(service.User),
		Post:          post.NewPostHandler(service.Post),
		Group:         group.NewGroupHandler(service.Group),
		Chat:          chat.NewChatHandler(chatRepo, hub, userRepo, uploader),
		Notifications: notifications.NewHandler(service.Notifications),
	}
}
