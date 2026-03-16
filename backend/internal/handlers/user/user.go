package user

import "social-network/internal/models"

type UserHandler struct {
	serv models.UserService
}

func NewUserHandler(serv models.UserHandler) *UserHandler {
	return &UserHandler{serv: serv}
}
