package user

import (
	"social-network/internal/models/follow"
	"social-network/internal/models/post"
	"social-network/internal/models/session"
	"social-network/internal/models/user"
)

type UserService struct {
	users    UserRepo
	sessions SessionRepo
	follows  FollowRepo
	posts    PostRepo
}

type UserRepo interface {
	CreateUser(*user.User) (int64, error)
	UpdateUser(*user.User) error
	GetByEmail(string) (*user.User, error)
	GetByID(id int64) (*user.User, error)
}

type SessionRepo interface {
	CreateSession(*session.Session) error
	Delete(int64) error
	GetUserID(string) (int64, error)
}

type FollowRepo interface {
	GetFollowers(int64) ([]follow.Follow, error)
	GetFollowing(int64) ([]follow.Follow, error)
}

type PostRepo interface {
	GetPosts(int64) ([]post.Post, error)
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{users: repo}
}
