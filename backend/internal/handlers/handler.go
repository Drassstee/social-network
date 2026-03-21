package handlers

import (
	"social-network/internal/handlers/post"
	"social-network/internal/handlers/user"
	"social-network/internal/service"
)

type Handler struct {
	User *user.UserHandler
	Post *post.PostHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		User: user.NewUserHandler(service.User),
		Post: post.NewPostHandler(service.Post),
	}
}
