package handlers

import (
	"social-network/internal/handlers/user"
	"social-network/internal/service"
)

type Handler struct {
	*user.UserHandler
}

func NewHandler(serv *service.Service) *Handler {
	return &Handler{
		UserHandler: user.NewUserHandler(serv),
	}
}
