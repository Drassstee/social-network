package user

import (
	"social-network/internal/models"
)

type UserService struct {
	users    models.UserRepo
	sessions models.SessionRepo
	follows  models.FollowRepo
	posts    models.PostRepo
}

func NewUserService(ur models.UserRepo, sr models.SessionRepo, fr models.FollowRepo, pr models.PostRepo) *UserService {
	return &UserService{
		users:    ur,
		sessions: sr,
		follows:  fr,
		posts:    pr,
	}
}
