package service

import (
	"social-network/internal/models"
	"social-network/internal/repository"
	servicepost "social-network/internal/service/post"
	"social-network/internal/service/user"
)

type Service struct {
	User models.UserService
	Post *servicepost.PostService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: user.NewUserService(repo.User),
		Post: servicepost.NewPostService(repo.Post),
	}
}
