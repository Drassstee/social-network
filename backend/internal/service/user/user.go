package user

import (
	"social-network/internal/repository/user"
)

type UserService struct {
	repo *user.UserRepository
}

func NewUserService(repo *user.UserRepository) *UserService {
	return &UserService{repo: repo}
}
