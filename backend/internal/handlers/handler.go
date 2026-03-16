package handlers

import (
	"social-network/internal/handlers/user"
	"social-network/internal/models"
	"social-network/internal/service"
)

type Handler struct {
	User models.UserHandler
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{User: user.NewUserHandler(service.User)}
}
