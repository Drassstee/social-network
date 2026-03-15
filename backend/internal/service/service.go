package service

import (
	"social-network/internal/repository"
	"social-network/internal/service/user"
)

type Service struct {
	User *user.UserService
}

func NewRepo(repo *repository.Repository) *Service {
	return &Service{
		User: user.NewUserService(repo.User),
	}
}
