package user

import (
	"social-network/internal/models"
)

type UserService struct {
	repo models.UserRepo
}

func NewUserService(repo models.UserRepo) *UserService {
	return &UserService{repo: repo}
}
