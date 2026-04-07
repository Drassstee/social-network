package service

import "social-network/internal/service/user"

type Repository interface {
	user.SessionRepo
	user.UserRepo
	user.FollowRepo
	user.PostRepo
}

type Service struct {
	*user.UserService
}

func NewService(r Repository) *Service {
	return &Service{
		UserService: user.NewUserService(r, r, r, r),
	}
}
