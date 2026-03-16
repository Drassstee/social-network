package service

import (
	"social-network/internal/models"
	"social-network/internal/repository"
	"social-network/internal/service/user"
)

type Service struct {
	User models.UserService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: user.NewUserService(repo.User),
	}
}
